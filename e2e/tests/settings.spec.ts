import { test, expect } from '@playwright/test';
import { TEST_USER } from '../fixtures/test-constants';

test.describe('Settings', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.getByPlaceholder(/username/i).fill(TEST_USER.username);
    await page.getByPlaceholder(/password/i).fill(TEST_USER.password);
    await page.getByRole('button', { name: /login/i }).click();
    await expect(page.getByText(/all.*stats/i)).toBeVisible();
  });

  test('open and close settings modal with keyboard shortcuts', async ({ page }) => {
    await page.keyboard.press('Backslash');
    await expect(page.getByRole('tab', { name: /appearance/i })).toBeVisible();
    await expect(page.getByRole('tab', { name: /account/i })).toBeVisible();
    await expect(page.getByRole('tab', { name: /api keys/i })).toBeVisible();
    await expect(page.getByRole('tab', { name: /export/i })).toBeVisible();

    await page.keyboard.press('Escape');
    await expect(page.getByRole('tab', { name: /appearance/i })).not.toBeVisible();
  });

  test('account settings - display current username', async ({ page }) => {
    await page.keyboard.press('Backslash');
    await page.getByRole('tab', { name: /account/i }).click();

    await expect(page.getByText(new RegExp(`you.*logged in as.*${TEST_USER.username}`, 'i'))).toBeVisible();
    await expect(page.getByRole('button', { name: /logout/i })).toBeVisible();
  });

  test('account settings - update username', async ({ page }) => {
    const newUsername = 'updateduser';

    await page.keyboard.press('Backslash');
    await page.getByRole('tab', { name: /account/i }).click();

    await page.getByPlaceholder(/update username/i).fill(newUsername);
    await page.getByRole('button', { name: /^submit$/i }).first().click();

    await expect(page.getByText(/successfully updated user/i)).toBeVisible();
    await expect(page.getByText(new RegExp(`you.*logged in as.*${newUsername}`, 'i'))).toBeVisible();

    await page.getByPlaceholder(/update username/i).fill(TEST_USER.username);
    await page.getByRole('button', { name: /^submit$/i }).first().click();
    await expect(page.getByText(/successfully updated user/i)).toBeVisible();
  });

  test('api keys - create, display and delete', async ({ page }) => {
    const keyLabel = 'test-api-key';

    await page.keyboard.press('Backslash');
    await page.getByRole('tab', { name: /api keys/i }).click();

    await expect(page.getByRole('heading', { name: /api keys/i })).toBeVisible();

    await page.getByPlaceholder(/add a label for a new api key/i).fill(keyLabel);
    await page.getByRole('button', { name: /create/i }).click();

    const keyElement = page.locator('div').filter({ hasText: new RegExp(`... ${keyLabel}$`) });
    await expect(keyElement).toBeVisible();

    const copyButton = page.getByRole('button').filter({ has: page.locator('svg') }).first();
    await expect(copyButton).toBeVisible();

    const deleteButton = page.getByRole('button', { name: '' }).nth(1);
    await deleteButton.click();

    await expect(keyElement).not.toBeVisible();
  });

  test('appearance tab - theme selection and reset', async ({ page }) => {
    await page.keyboard.press('Backslash');
    await page.getByRole('tab', { name: /appearance/i }).click();

    await expect(page.getByRole('heading', { name: /select theme/i })).toBeVisible();
    await expect(page.getByRole('button', { name: /reset/i })).toBeVisible();

    const themeButtons = page.locator('[role="button"]').filter({ has: page.locator('div[class*="theme"]') });
    await expect(themeButtons.first()).toBeVisible();

    await page.getByRole('heading', { name: /use custom theme/i }).click();
    await expect(page.locator('textarea#custom-theme-input')).toBeVisible();

    const bgColor = await page.evaluate(() => {
      return getComputedStyle(document.documentElement).getPropertyValue('--color-bg');
    });
    expect(bgColor).toBeTruthy();
  });
});
