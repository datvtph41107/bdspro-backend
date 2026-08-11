import crypto from 'node:crypto';

const secret = process.env.QHPRO_LAB_JWT_SECRET;
const profileId = Number(process.env.QHPRO_LAB_PROFILE_ID || '42');
if (!secret || !Number.isSafeInteger(profileId) || profileId <= 0) {
  throw new Error('QHPRO_LAB_JWT_SECRET and a positive QHPRO_LAB_PROFILE_ID are required');
}

const now = Math.floor(Date.now() / 1000);
const encode = value => Buffer.from(JSON.stringify(value)).toString('base64url');
const header = encode({alg: 'HS256', typ: 'JWT'});
const payload = encode({
  sub: profileId,
  profile: profileId,
  origin: profileId,
  type: 'ACCESS',
  role: 'user',
  session: 1,
  iat: now,
  exp: now + 3600,
});
const unsigned = `${header}.${payload}`;
const signature = crypto.createHmac('sha256', secret).update(unsigned).digest('base64url');
process.stdout.write(`${unsigned}.${signature}`);
