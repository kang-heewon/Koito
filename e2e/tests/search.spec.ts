import { test, expect } from '@playwright/test';
import { TEST_ARTISTS, API_BASE } from '../fixtures/test-constants';

test.describe('SearchModal', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.waitForLoadState('networkidle');
  });

  test('검색 모달 열기/닫기', async ({ page }) => {
    const modal = page.getByRole('dialog');
    const searchInput = page.getByRole('textbox');

    await expect(modal).not.toBeVisible();

    await page.keyboard.press('/');
    await expect(modal).toBeVisible();
    await expect(searchInput).toBeFocused();

    await page.keyboard.press('Escape');
    await expect(modal).not.toBeVisible();
  });

  test('검색 결과 표시', async ({ page }) => {
    const testArtist = TEST_ARTISTS[0].name;

    await page.keyboard.press('/');
    await page.waitForSelector('[role="dialog"]', { state: 'visible' });

    const searchInput = page.getByRole('textbox');
    await searchInput.fill(testArtist);

    const responsePromise = page.waitForResponse(
      (resp) => resp.url().includes('/apis/web/v1/search') && resp.status() === 200
    );
    await responsePromise;

    const results = page.locator('[role="dialog"] a, [role="dialog"] [role="link"], [role="dialog"] .result, [role="dialog"] [class*="result"]').first();
    await expect(results).toBeVisible({ timeout: 5000 });
  });

  test('빈 검색 결과', async ({ page }) => {
    const nonExistentName = 'xyznonexistent12345abc';

    await page.keyboard.press('/');
    await page.waitForSelector('[role="dialog"]', { state: 'visible' });

    const searchInput = page.getByRole('textbox');
    await searchInput.fill(nonExistentName);

    const responsePromise = page.waitForResponse(
      (resp) => resp.url().includes('/apis/web/v1/search') && resp.status() === 200
    );
    await responsePromise;

    await page.waitForFunction(() => {
      const dialog = document.querySelector('[role="dialog"]');
      return dialog && dialog.querySelectorAll('a').length === 0;
    });

    const results = page.locator('[role="dialog"] a, [role="dialog"] [role="link"], [role="dialog"] .result, [role="dialog"] [class*="result"]');
    const count = await results.count();
    expect(count).toBe(0);

    const emptyState = page
      .locator('[role="dialog"]')
      .locator('text=/no results|empty|not found|결과 없음/i')
      .or(page.locator('[role="dialog"] [class*="empty"]'))
      .or(page.locator('[role="dialog"] [class*="no-result"]'));

    const hasEmptyState = await emptyState.isVisible().catch(() => false);
    const hasNoResults = count === 0;
    expect(hasEmptyState || hasNoResults).toBeTruthy();
  });

  test('검색 후 네비게이션', async ({ page }) => {
    const testArtist = TEST_ARTISTS[0].name;

    await page.keyboard.press('/');
    await page.waitForSelector('[role="dialog"]', { state: 'visible' });

    const searchInput = page.getByRole('textbox');
    await searchInput.fill(testArtist);

    const responsePromise = page.waitForResponse(
      (resp) => resp.url().includes('/apis/web/v1/search') && resp.status() === 200
    );
    const searchResponse = await responsePromise;
    const searchData = await searchResponse.json();

    const resultLink = page.locator('[role="dialog"] a').first();
    await expect(resultLink).toBeVisible({ timeout: 5000 });

    const href = await resultLink.getAttribute('href');

    await resultLink.click();

    const modal = page.getByRole('dialog');
    await expect(modal).not.toBeVisible();

    await expect(page).toHaveURL(/\/(artist|album|track)\/.+/);

    if (searchData.artists && searchData.artists.length > 0) {
      const artistId = searchData.artists[0].id;
      await expect(page).toHaveURL(new RegExp(`/artist/${artistId}`));
    }
  });
});
