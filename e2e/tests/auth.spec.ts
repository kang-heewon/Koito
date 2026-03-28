import { test, expect } from '@playwright/test';
import { TEST_USER } from '../fixtures/test-constants';

test.describe('Authentication', () => {
  test('successful login redirects to dashboard', async ({ page }) => {
    await page.goto('/');

    await expect(page.getByRole('heading', { name: /log in/i })).toBeVisible();

    await page.getByPlaceholder(/username/i).fill(TEST_USER.username);
    await page.getByPlaceholder(/password/i).fill(TEST_USER.password);
    await page.getByRole('button', { name: /login/i }).click();

    await expect(page.getByText(/all.*stats/i)).toBeVisible();
    await expect(page.getByText(/top artists/i)).toBeVisible();
    await expect(page.getByText(/activity/i)).toBeVisible();
  });

  test('invalid credentials show error message', async ({ page }) => {
    await page.goto('/');

    await expect(page.getByRole('heading', { name: /log in/i })).toBeVisible();

    await page.getByPlaceholder(/username/i).fill('wronguser');
    await page.getByPlaceholder(/password/i).fill('wrongpassword');
    await page.getByRole('button', { name: /login/i }).click();

    await expect(page.getByText(/invalid/i)).toBeVisible();
    await expect(page.getByRole('heading', { name: /log in/i })).toBeVisible();
  });

  test('logout returns to login form', async ({ page }) => {
    await page.goto('/');

    await page.getByPlaceholder(/username/i).fill(TEST_USER.username);
    await page.getByPlaceholder(/password/i).fill(TEST_USER.password);
    await page.getByRole('button', { name: /login/i }).click();

    await expect(page.getByText(/all.*stats/i)).toBeVisible();

    await page.keyboard.press('Backslash');
    await expect(page.getByRole('tab', { name: /account/i })).toBeVisible();

    await page.getByRole('button', { name: /logout/i }).click();

    await expect(page.getByRole('heading', { name: /log in/i })).toBeVisible();
  });

  test('session persists after page reload', async ({ page, context }) => {
    await page.goto('/');

    await page.getByPlaceholder(/username/i).fill(TEST_USER.username);
    await page.getByPlaceholder(/password/i).fill(TEST_USER.password);
    await page.getByRole('button', { name: /login/i }).click();

    await expect(page.getByText(/all.*stats/i)).toBeVisible();

    await page.reload();

    await expect(page.getByText(/all.*stats/i)).toBeVisible();
    await expect(page.getByText(/top artists/i)).toBeVisible();
    await expect(page.getByRole('heading', { name: /log in/i })).not.toBeVisible();
  });
});
