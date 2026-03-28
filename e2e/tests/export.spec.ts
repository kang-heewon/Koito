import { test, expect } from '@playwright/test';
import { TEST_USER } from '../fixtures/test-constants';

test.describe('Export', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.waitForLoadState('networkidle');

    await page.getByPlaceholder(/username/i).fill(TEST_USER.username);
    await page.getByPlaceholder(/password/i).fill(TEST_USER.password);
    await page.getByRole('button', { name: /login/i }).click();

    await expect(page.getByText(/all.*stats/i)).toBeVisible();
  });

  test('Export 모달 열기', async ({ page }) => {
    await page.keyboard.press('Backslash');

    const settingsModal = page.locator('[role="dialog"]').filter({ hasText: /settings|설정/i });
    await expect(settingsModal).toBeVisible();

    const exportTab = page.getByRole('tab', { name: /export/i });
    await exportTab.click();

    await expect(page.getByRole('heading', { name: /export/i })).toBeVisible();
    await expect(page.getByRole('button', { name: /export data/i })).toBeVisible();
  });

  test('Export 다운로드', async ({ page }) => {
    await page.keyboard.press('Backslash');

    const exportTab = page.getByRole('tab', { name: /export/i });
    await exportTab.click();

    await expect(page.getByRole('button', { name: /export data/i })).toBeVisible();

    const [download] = await Promise.all([
      page.waitForEvent('download'),
      page.getByRole('button', { name: /export data/i }).click()
    ]);

    const path = await download.path();
    expect(path).toBeTruthy();

    const fileName = download.suggestedFilename();
    expect(fileName).toMatch(/\.json$/);
  });
});
