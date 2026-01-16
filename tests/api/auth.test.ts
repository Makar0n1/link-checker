import { describe, it, expect, beforeAll } from 'vitest';
import { authApi, generateTestEmail, createTestUser } from './helpers/api-client';

describe('Auth Service API', () => {
  describe('POST /api/v1/auth/register', () => {
    it('should register a new user and return 201', async () => {
      const email = generateTestEmail();
      const result = await authApi.register({
        email,
        password: 'ValidPassword123!',
        name: 'Test User',
      });

      expect(result.status).toBe(201);
      expect(result.data.email).toBe(email);
      expect(result.data.name).toBe('Test User');
      expect(result.data.id).toBeDefined();
      expect(result.data.role).toBe('user');
      expect(result.data.created_at).toBeDefined();
    });

    it('should return 409 when registering duplicate email', async () => {
      const email = generateTestEmail();

      // First registration should succeed
      const firstResult = await authApi.register({
        email,
        password: 'ValidPassword123!',
        name: 'Test User',
      });
      expect(firstResult.status).toBe(201);

      // Second registration with same email should fail
      const secondResult = await authApi.register({
        email,
        password: 'AnotherPassword123!',
        name: 'Another User',
      });
      expect(secondResult.status).toBe(409);
    });

    it('should return 400 for invalid request body', async () => {
      const result = await authApi.register({
        email: '',
        password: 'ValidPassword123!',
        name: 'Test User',
      });
      expect(result.status).toBe(400);
    });

    it('should return 400 for short password', async () => {
      const result = await authApi.register({
        email: generateTestEmail(),
        password: 'short',
        name: 'Test User',
      });
      expect(result.status).toBe(400);
    });

    it('should return 400 for missing name', async () => {
      const result = await authApi.register({
        email: generateTestEmail(),
        password: 'ValidPassword123!',
        name: '',
      });
      expect(result.status).toBe(400);
    });
  });

  describe('POST /api/v1/auth/login', () => {
    let testEmail: string;
    const testPassword = 'ValidPassword123!';

    beforeAll(async () => {
      testEmail = generateTestEmail();
      await authApi.register({
        email: testEmail,
        password: testPassword,
        name: 'Login Test User',
      });
    });

    it('should login successfully and return 200 with tokens', async () => {
      const result = await authApi.login({
        email: testEmail,
        password: testPassword,
      });

      expect(result.status).toBe(200);
      expect(result.data.access_token).toBeDefined();
      expect(result.data.refresh_token).toBeDefined();
      expect(result.data.expires_in).toBeDefined();
      expect(result.data.token_type).toBe('Bearer');
    });

    it('should return 401 for wrong password', async () => {
      const result = await authApi.login({
        email: testEmail,
        password: 'WrongPassword123!',
      });

      expect(result.status).toBe(401);
    });

    it('should return 401 for non-existent user', async () => {
      const result = await authApi.login({
        email: 'nonexistent@example.com',
        password: 'AnyPassword123!',
      });

      expect(result.status).toBe(401);
    });

    it('should return 400 for missing email', async () => {
      const result = await authApi.login({
        email: '',
        password: testPassword,
      });

      expect(result.status).toBe(400);
    });

    it('should return 400 for missing password', async () => {
      const result = await authApi.login({
        email: testEmail,
        password: '',
      });

      expect(result.status).toBe(400);
    });
  });

  describe('GET /api/v1/auth/me', () => {
    let accessToken: string;
    let testEmail: string;

    beforeAll(async () => {
      const user = await createTestUser();
      accessToken = user.token;
      testEmail = user.email;
    });

    it('should return current user with valid token', async () => {
      const result = await authApi.me(accessToken);

      expect(result.status).toBe(200);
      expect(result.data.email).toBe(testEmail);
      expect(result.data.id).toBeDefined();
      expect(result.data.name).toBeDefined();
      expect(result.data.role).toBeDefined();
    });

    it('should return 401 without token', async () => {
      const result = await authApi.me('');

      expect(result.status).toBe(401);
    });

    it('should return 401 with invalid token', async () => {
      const result = await authApi.me('invalid-token');

      expect(result.status).toBe(401);
    });

    it('should return 401 with malformed token', async () => {
      const result = await authApi.me('not.a.valid.jwt.token');

      expect(result.status).toBe(401);
    });
  });

  describe('POST /api/v1/auth/refresh', () => {
    let refreshToken: string;

    beforeAll(async () => {
      const user = await createTestUser();
      refreshToken = user.refreshToken;
    });

    it('should refresh tokens successfully and return 200', async () => {
      const result = await authApi.refresh(refreshToken);

      expect(result.status).toBe(200);
      expect(result.data.access_token).toBeDefined();
      expect(result.data.refresh_token).toBeDefined();
      expect(result.data.expires_in).toBeDefined();
      expect(result.data.token_type).toBe('Bearer');
    });

    it('should return 401 for invalid refresh token', async () => {
      const result = await authApi.refresh('invalid-refresh-token');

      expect(result.status).toBe(401);
    });

    it('should return 400 for empty refresh token', async () => {
      const result = await authApi.refresh('');

      expect(result.status).toBe(400);
    });
  });

  describe('POST /api/v1/auth/logout', () => {
    it('should logout successfully and return 200', async () => {
      const user = await createTestUser();
      const result = await authApi.logout(user.refreshToken);

      expect(result.status).toBe(200);
      expect(result.data.message).toContain('logout');
    });

    it('should return 400 for empty refresh token', async () => {
      const result = await authApi.logout('');

      expect(result.status).toBe(400);
    });
  });

  describe('GET /health', () => {
    it('should return health status', async () => {
      const result = await authApi.health();

      expect(result.status).toBe(200);
      expect(result.data.status).toBe('ok');
    });
  });
});
