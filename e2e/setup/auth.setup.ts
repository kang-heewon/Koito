import { test as setup, expect } from '@playwright/test';
import * as fs from 'fs';
import * as path from 'path';

const authFile = path.resolve(__dirname, '..', '.auth', 'user.json');

setup('authenticate', async ({ request }) => {
  const loginResponse = await request.post('/apis/web/v1/login', {
    form: {
      username: 'testuser',
      password: 'testpass123',
    },
  });

  expect(loginResponse.status()).toBe(204);

  const headers = await loginResponse.headers();
  const setCookie = headers['set-cookie'];
  expect(setCookie).toBeDefined();

  const match = setCookie?.match(/koito_session=([^;]+)/);
  expect(match).toBeTruthy();
  const sessionValue = match?.[1];
  expect(sessionValue).toBeTruthy();

  const authData = {
    cookies: [
      {
        name: 'koito_session',
        value: sessionValue,
        domain: 'localhost',
        path: '/',
        httpOnly: true,
        secure: false,
        sameSite: 'Lax',
      },
    ],
    origins: [],
  };

  fs.mkdirSync(path.dirname(authFile), { recursive: true });
  fs.writeFileSync(authFile, JSON.stringify(authData, null, 2));
});
