import { test, expect } from '@playwright/test';
import { TEST_USER } from '../fixtures/test-constants';

test.describe('Charts', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login');
    await page.getByLabel('Username').fill(TEST_USER.username);
    await page.getByLabel('Password').fill(TEST_USER.password);
    await page.getByRole('button', { name: /log in/i }).click();
    await page.waitForURL('/');
  });

  test('Top Artists 차트 표시 및 기간 필터 변경', async ({ page }) => {
    await page.goto('/chart/top-artists');

    await expect(page.getByRole('heading', { name: 'Top Artists' })).toBeVisible();

    const listItems = page.getByRole('listitem');
    await expect(listItems.count()).resolves.toBeGreaterThan(0);

    const allTimeButton = page.getByRole('button', { name: 'All Time' });
    const weekButton = page.getByRole('button', { name: 'Week' });
    const monthButton = page.getByRole('button', { name: 'Month' });
    const yearButton = page.getByRole('button', { name: 'Year' });

    await expect(allTimeButton).toBeVisible();
    await expect(weekButton).toBeVisible();
    await expect(monthButton).toBeVisible();
    await expect(yearButton).toBeVisible();

    await monthButton.click();
    await page.waitForLoadState('networkidle');
    await expect(listItems.count()).resolves.toBeGreaterThan(0);

    await yearButton.click();
    await page.waitForLoadState('networkidle');
    await expect(listItems.count()).resolves.toBeGreaterThan(0);
  });

  test('Top Albums 차트 표시 및 기간 필터 변경', async ({ page }) => {
    await page.goto('/chart/top-albums');

    await expect(page.getByRole('heading', { name: 'Top Albums' })).toBeVisible();

    const listItems = page.getByRole('listitem');
    await expect(listItems.count()).resolves.toBeGreaterThan(0);

    const weekButton = page.getByRole('button', { name: 'Week' });
    const monthButton = page.getByRole('button', { name: 'Month' });

    await expect(weekButton).toBeVisible();
    await expect(monthButton).toBeVisible();

    await weekButton.click();
    await page.waitForLoadState('networkidle');
    await expect(listItems.count()).resolves.toBeGreaterThan(0);

    await monthButton.click();
    await page.waitForLoadState('networkidle');
    await expect(listItems.count()).resolves.toBeGreaterThan(0);
  });

  test('Top Tracks 차트 표시 및 기간 필터 변경', async ({ page }) => {
    await page.goto('/chart/top-tracks');

    await expect(page.getByRole('heading', { name: 'Top Tracks' })).toBeVisible();

    const listItems = page.getByRole('listitem');
    await expect(listItems.count()).resolves.toBeGreaterThan(0);

    const allTimeButton = page.getByRole('button', { name: 'All Time' });
    const dayButton = page.getByRole('button', { name: 'Day' });

    await expect(allTimeButton).toBeVisible();
    await expect(dayButton).toBeVisible();

    await dayButton.click();
    await page.waitForLoadState('networkidle');
    await expect(listItems.count()).resolves.toBeGreaterThan(0);

    await allTimeButton.click();
    await page.waitForLoadState('networkidle');
    await expect(listItems.count()).resolves.toBeGreaterThan(0);
  });

  test('차트 항목 클릭 시 상세 페이지로 이동', async ({ page }) => {
    await page.goto('/chart/top-artists');

    const listItems = page.getByRole('listitem');
    await expect(listItems.count()).resolves.toBeGreaterThan(0);

    const firstItem = listItems.first();
    const link = firstItem.locator('a').first();
    await expect(link).toBeVisible();

    const href = await link.getAttribute('href');
    expect(href).toMatch(/^\/artist\/\d+$/);

    await link.click();
    await page.waitForURL(/\/artist\/\d+$/);
    await expect(page.url()).toMatch(/\/artist\/\d+$/);
  });

  test('Top Albums 항목 클릭 시 앨범 상세 페이지로 이동', async ({ page }) => {
    await page.goto('/chart/top-albums');

    const listItems = page.getByRole('listitem');
    await expect(listItems.count()).resolves.toBeGreaterThan(0);

    const firstItem = listItems.first();
    const link = firstItem.locator('a').first();
    await expect(link).toBeVisible();

    const href = await link.getAttribute('href');
    expect(href).toMatch(/^\/album\/\d+$/);

    await link.click();
    await page.waitForURL(/\/album\/\d+$/);
    await expect(page.url()).toMatch(/\/album\/\d+$/);
  });

  test('Top Tracks 항목 클릭 시 트랙 상세 페이지로 이동', async ({ page }) => {
    await page.goto('/chart/top-tracks');

    const listItems = page.getByRole('listitem');
    await expect(listItems.count()).resolves.toBeGreaterThan(0);

    const firstItem = listItems.first();
    const link = firstItem.locator('a').first();
    await expect(link).toBeVisible();

    const href = await link.getAttribute('href');
    expect(href).toMatch(/^\/track\/\d+$/);

    await link.click();
    await page.waitForURL(/\/track\/\d+$/);
    await expect(page.url()).toMatch(/\/track\/\d+$/);
  });

  test('네비게이션에서 Top Artists 차트로 이동', async ({ page }) => {
    await page.goto('/');

    const topArtistsLink = page.getByRole('link', { name: /top artists/i });
    await expect(topArtistsLink).toBeVisible();

    await topArtistsLink.click();
    await page.waitForURL('/chart/top-artists');
    await expect(page.getByRole('heading', { name: 'Top Artists' })).toBeVisible();
  });
});
