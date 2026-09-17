#!/usr/bin/env python3
"""Run a temporary HTTPS tunnel forwarding only the SePay bank webhook."""

import errno
import http.client
import re
import subprocess
import threading
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path

WEBHOOK_PATH = "/v2/payment/sepay/webhook"


class WebhookProxy(BaseHTTPRequestHandler):
    def do_POST(self):
        if self.path != WEBHOOK_PATH:
            self.send_error(404)
            return
        try:
            length = int(self.headers.get("Content-Length", "0"))
        except ValueError:
            self.send_error(400)
            return
        if not 0 < length <= 65536 or self.headers.get("Transfer-Encoding"):
            self.send_error(400)
            return
        self.connection.settimeout(15)
        connection = http.client.HTTPConnection("127.0.0.1", 8000, timeout=20)
        try:
            body = self.rfile.read(length)
            if len(body) != length:
                self.send_error(400)
                return
            connection.request("POST", WEBHOOK_PATH, body, {
                "Content-Type": self.headers.get("Content-Type", "application/json"),
                "Authorization": self.headers.get("Authorization", ""),
            })
            response = connection.getresponse()
            payload = response.read()
            self.send_response(response.status)
            self.send_header("Content-Type", response.getheader("Content-Type", "application/json"))
            self.send_header("Content-Length", str(len(payload)))
            self.end_headers()
            self.wfile.write(payload)
        except (OSError, http.client.HTTPException):
            self.send_error(502, "Local Gateway unavailable")
        finally:
            connection.close()

    def log_message(self, *_):
        pass


if __name__ == "__main__":
    root = Path(__file__).resolve().parents[2]
    cloudflared = root / ".tmp/development/tools/cloudflared"
    url_file = root / ".tmp/development/sepay-webhook-url.txt"
    if not cloudflared.is_file():
        raise SystemExit("Install cloudflared at .tmp/development/tools/cloudflared first.")
    url_file.unlink(missing_ok=True)
    try:
        server = ThreadingHTTPServer(("127.0.0.1", 18080), WebhookProxy)
    except OSError as error:
        if error.errno == errno.EADDRINUSE:
            raise SystemExit(
                "SePay webhook tunnel is already running on 127.0.0.1:18080. "
                "Stop the existing tunnel before starting another one."
            ) from error
        raise
    threading.Thread(target=server.serve_forever, daemon=True).start()
    print("Creating the SePay webhook tunnel. Waiting for the public HTTPS URL...", flush=True)
    process = None
    try:
        process = subprocess.Popen([
            str(cloudflared), "tunnel", "--no-autoupdate", "--protocol", "http2",
            "--url", "http://127.0.0.1:18080",
        ], stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True, bufsize=1)
        assert process.stdout is not None
        published = False
        for line in process.stdout:
            print(line, end="", flush=True)
            match = re.search(r"https://[a-z0-9-]+\.trycloudflare\.com", line)
            if match and not published:
                webhook_url = match.group(0) + WEBHOOK_PATH
                url_file.write_text(webhook_url + "\n", encoding="utf-8")
                print("\n============================================================", flush=True)
                print("COPY THIS URL INTO SEPAY:", flush=True)
                print(webhook_url, flush=True)
                print("Saved to: " + str(url_file), flush=True)
                print("Keep this terminal running while testing the webhook.", flush=True)
                print("============================================================\n", flush=True)
                published = True
        if process.wait() != 0:
            raise SystemExit("cloudflared stopped before the tunnel was ready.")
    finally:
        if process is not None and process.poll() is None:
            process.terminate()
        server.shutdown()
        server.server_close()
        # A quick-tunnel URL is valid only while this process is connected.
        # Removing it prevents operators from copying a stale endpoint into SePay.
        url_file.unlink(missing_ok=True)
