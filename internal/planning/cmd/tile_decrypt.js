/**
 * tile_decrypt.js
 *
 * Client-side AES-256-CTR tile decryption for ci_* tilesets.
 *
 * Server format (EncryptBytesCTR):
 *   [nonce: 16 bytes] + [ciphertext]   (same length as plaintext)
 *
 * Flow:
 *   1. POST /v2/user/profile/session  → { sessionK, sessionEncryptKey }
 *   2. Fetch tile with header X-Session-Id: {sessionK}
 *   3. If X-Content-Encrypted: aes-256-ctr → decrypt body
 */

const TILE_ENCRYPTED_HEADER = "X-Content-Encrypted";
const TILE_ENCRYPTED_HEADER_VAL = "aes-256-ctr";
const TILE_SESSION_ID_HEADER = "X-Session-Id";
const CTR_NONCE_SIZE = 16;

/**
 * @param {string} base64Key - base64-encoded 32-byte AES key
 * @returns {Promise<CryptoKey>}
 */
async function importAESKey(base64Key) {
  const rawKey = Uint8Array.from(atob(base64Key), (c) => c.charCodeAt(0));
  if (rawKey.length !== 32) {
    throw new Error(`AES key must be 32 bytes, got ${rawKey.length}`);
  }
  return crypto.subtle.importKey("raw", rawKey, { name: "AES-CTR" }, false, [
    "decrypt",
  ]);
}

/**
 * @param {ArrayBuffer} encryptedBuffer
 * @param {CryptoKey|string} key
 * @returns {Promise<ArrayBuffer>}
 */
async function decryptTileBody(encryptedBuffer, key) {
  const cryptoKey =
    typeof key === "string" ? await importAESKey(key) : key;

  const data = new Uint8Array(encryptedBuffer);
  if (data.length < CTR_NONCE_SIZE) {
    throw new Error(
      `Encrypted tile too short: ${data.length} bytes (min ${CTR_NONCE_SIZE})`
    );
  }

  const counter = data.slice(0, CTR_NONCE_SIZE);
  const ciphertext = data.slice(CTR_NONCE_SIZE);

  return crypto.subtle.decrypt(
    { name: "AES-CTR", counter, length: 128 },
    cryptoKey,
    ciphertext
  );
}

/**
 * @param {string} tileUrl
 * @param {string} sessionId - sessionK string
 * @param {CryptoKey|string} key
 * @param {AbortSignal} [signal]
 * @returns {Promise<ArrayBuffer>}
 */
async function fetchTile(tileUrl, sessionId, key, signal) {
  const res = await fetch(tileUrl, {
    headers: { [TILE_SESSION_ID_HEADER]: sessionId },
    signal,
  });

  if (res.status === 401) {
    throw new Error("Tile session invalid or expired (401). Re-create session.");
  }
  if (!res.ok) {
    throw new Error(`Tile fetch failed: ${res.status} ${res.statusText}`);
  }

  const buffer = await res.arrayBuffer();
  const isEncrypted =
    res.headers.get(TILE_ENCRYPTED_HEADER) === TILE_ENCRYPTED_HEADER_VAL;

  return isEncrypted ? decryptTileBody(buffer, key) : buffer;
}

class TileSession {
  constructor(sessionApiUrl, fetchOptions = {}) {
    this._sessionApiUrl = sessionApiUrl;
    this._fetchOptions = fetchOptions;
    this._sessionId = null;
    this._cryptoKey = null;
    this._expiresAt = null;
    this._refreshPromise = null;
  }

  async init() {
    if (this._refreshPromise) return this._refreshPromise;

    this._refreshPromise = (async () => {
      const res = await fetch(this._sessionApiUrl, {
        method: "POST",
        ...this._fetchOptions,
      });
      if (!res.ok) {
        throw new Error(`Session API error: ${res.status} ${res.statusText}`);
      }

      const body = await res.json();
      const sessionId = body.sessionK ?? body.sessionId;
      const aesKey = body.sessionEncryptKey ?? body.aesKey;
      if (!sessionId || !aesKey) {
        throw new Error("Session API missing sessionK or sessionEncryptKey");
      }

      this._sessionId = String(sessionId);
      this._cryptoKey = await importAESKey(aesKey);
      const expiresIn = body.expiresIn;
      this._expiresAt =
        typeof expiresIn === "number"
          ? new Date(Date.now() + expiresIn * 1000)
          : body.expiresAt
            ? new Date(body.expiresAt)
            : null;
    })();

    try {
      await this._refreshPromise;
    } finally {
      this._refreshPromise = null;
    }
  }

  get isValid() {
    if (!this._sessionId || !this._cryptoKey) return false;
    if (!this._expiresAt) return true;
    return Date.now() < this._expiresAt.getTime() - 30_000;
  }

  async fetchTile(tileUrl, signal) {
    if (!this.isValid) await this.init();

    try {
      return await fetchTile(tileUrl, this._sessionId, this._cryptoKey, signal);
    } catch (err) {
      if (err.message.includes("401")) {
        await this.init();
        return fetchTile(tileUrl, this._sessionId, this._cryptoKey, signal);
      }
      throw err;
    }
  }

  get sessionId() {
    return this._sessionId;
  }
}

function registerCIProtocol(maplibregl, session) {
  maplibregl.addProtocol("ci", async (params, abortController) => {
    const url = params.url.replace(/^ci:\/\//, "http://");
    const data = await session.fetchTile(url, abortController.signal);
    return { data };
  });
}

if (typeof module !== "undefined" && module.exports) {
  module.exports = {
    importAESKey,
    decryptTileBody,
    fetchTile,
    TileSession,
    registerCIProtocol,
  };
} else if (typeof window !== "undefined") {
  window.TileDecrypt = {
    importAESKey,
    decryptTileBody,
    fetchTile,
    TileSession,
    registerCIProtocol,
  };
}
