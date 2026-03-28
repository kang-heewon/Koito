import { test, expect } from '../fixtures/auth.fixture';
import { seedDefaultData, generateRandomListens, Listen } from '../fixtures/test-data';
import { API_BASE, TEST_ARTISTS, TEST_TRACKS } from '../fixtures/test-constants';

test.describe('CRUD Operations', () => {
  test.beforeEach(async ({ authPage }) => {
    const cookies = await authPage.context().cookies();
    const sessionCookie = cookies.find(c => c.name === 'session')?.value;
    if (sessionCookie) {
      await seedDefaultData(`session=${sessionCookie}`, 15);
    }
  });

  test.describe('앨범 Various Artists 토글', () => {
    test('앨범의 Various Artists 상태 토글', async ({ authPage }) => {
      await authPage.goto('/');

      const artistLink = authPage.getByRole('link', { name: /test artist/i }).first();
      await artistLink.click();
      await expect(authPage).toHaveURL(/\/artist\/\d+/);

      const albumLink = authPage.getByRole('link', { name: /test album/i }).first();
      await albumLink.click();
      await expect(authPage).toHaveURL(/\/album\/\d+/);

      const editButton = authPage.getByRole('button', { name: /edit/i });
      await editButton.click();

      const dialog = authPage.getByRole('dialog');
      await expect(dialog).toBeVisible();

      const variousArtistsToggle = dialog.locator('input[type="checkbox"]').filter({ hasText: /various artists|various/i });
      const checkbox = dialog.locator('input[name="is_various_artists"], [data-testid="various-artists-toggle"]').or(dialog.locator('input[type="checkbox"]').first());

      const initialChecked = await checkbox.isChecked().catch(() => false);

      const patchResponse = authPage.waitForResponse(
        (resp) => resp.url().includes('/apis/web/v1/album') && resp.request().method() === 'PATCH' && resp.status() === 200
      );

      await checkbox.click();
      await patchResponse;

      await authPage.getByRole('button', { name: /save|확인|save changes/i }).click();

      await expect(dialog).not.toBeVisible();
    });
  });

  test.describe('트랙 병합', () => {
    test('트랙 상세 페이지에서 병합', async ({ authPage }) => {
      await authPage.goto('/');

      const artistLink = authPage.getByRole('link', { name: /test artist/i }).first();
      await artistLink.click();

      const albumLink = authPage.getByRole('link', { name: /test album/i }).first();
      await albumLink.click();

      const trackLink = authPage.getByRole('link', { name: /test track/i }).first();
      await trackLink.click();
      await expect(authPage).toHaveURL(/\/track\/\d+/);

      const mergeButton = authPage.getByRole('button', { name: /merge/i });
      await mergeButton.click();

      const dialog = authPage.getByRole('dialog');
      await expect(dialog).toBeVisible();

      const searchInput = dialog.getByRole('textbox');
      await searchInput.fill('Test Track Two');

      const searchResponse = authPage.waitForResponse(
        (resp) => resp.url().includes('/apis/web/v1/search') && resp.status() === 200
      );
      await searchResponse;

      const resultItem = dialog.getByRole('listitem').first().or(dialog.locator('a, [role="button"]').filter({ hasText: /Test Track/ })).first();
      await resultItem.click();

      const mergeResponse = authPage.waitForResponse(
        (resp) => resp.url().includes('/apis/web/v1/merge/tracks') && resp.status() === 200
      );

      const confirmButton = dialog.getByRole('button', { name: /merge|confirm/i });
      await confirmButton.click();
      await mergeResponse;

      await expect(dialog).not.toBeVisible();
    });
  });

  test.describe('아티스트 삭제', () => {
    test('아티스트 상세 페이지에서 삭제', async ({ authPage }) => {
      await authPage.goto('/');

      const artistLink = authPage.getByRole('link', { name: /test artist/i }).first();
      await artistLink.click();
      await expect(authPage).toHaveURL(/\/artist\/\d+/);

      const deleteButton = authPage.getByRole('button', { name: /delete/i });
      await deleteButton.click();

      const dialog = authPage.getByRole('dialog');
      await expect(dialog).toBeVisible();

      const deleteResponse = authPage.waitForResponse(
        (resp) => resp.url().includes('/apis/web/v1/artist') && resp.request().method() === 'DELETE' && resp.status() === 200
      );

      const confirmButton = dialog.getByRole('button', { name: /delete|confirm|yes/i });
      await confirmButton.click();
      await deleteResponse;

      await expect(dialog).not.toBeVisible();
      await expect(authPage).toHaveURL('/');
    });
  });

  test.describe('리슨 삭제', () => {
    test('Listens 페이지에서 리슨 삭제', async ({ authPage }) => {
      await authPage.goto('/listens');

      const table = authPage.locator('table');
      await expect(table).toBeVisible();

      const rows = table.locator('tbody tr');
      await expect(rows.first()).toBeVisible();

      const firstRow = rows.first();
      const deleteButton = firstRow.getByRole('button', { name: /delete|remove/i }).or(firstRow.locator('button').filter({ has: authPage.locator('svg') })).first();

      await deleteButton.click();

      const dialog = authPage.getByRole('dialog');
      await expect(dialog).toBeVisible();

      const deleteResponse = authPage.waitForResponse(
        (resp) => resp.url().includes('/apis/web/v1/listen') && resp.request().method() === 'DELETE' && resp.status() === 200
      );

      const confirmButton = dialog.getByRole('button', { name: /delete|confirm|yes/i });
      await confirmButton.click();
      await deleteResponse;

      await expect(dialog).not.toBeVisible();
    });
  });

  test.describe('수동 리슨 추가', () => {
    test('AddListen 모달에서 트랙 정보 입력', async ({ authPage }) => {
      await authPage.goto('/listens');

      const addButton = authPage.getByRole('button', { name: /add listen|manual add|add track/i });
      await addButton.click();

      const dialog = authPage.getByRole('dialog');
      await expect(dialog).toBeVisible();

      const inputs = dialog.getByRole('textbox');
      const artistInput = inputs.filter({ has: authPage.locator('[placeholder*="artist"], [name*="artist"]').first() }).first();
      const trackInput = inputs.filter({ has: authPage.locator('[placeholder*="track"], [name*="track"]').first() }).first();
      const albumInput = inputs.filter({ has: authPage.locator('[placeholder*="album"], [name*="album"]').first() }).first();

      await artistInput.fill('Manual Test Artist');
      await trackInput.fill('Manual Test Track');
      await albumInput.fill('Manual Test Album');

      const submitResponse = authPage.waitForResponse(
        (resp) => resp.url().includes('/apis/web/v1/listen') && resp.request().method() === 'POST' && resp.status() === 200
      );

      const submitButton = dialog.getByRole('button', { name: /add|submit|save/i });
      await submitButton.click();
      await submitResponse;

      await expect(dialog).not.toBeVisible();

      await authPage.reload();
      await expect(authPage.getByText('Manual Test Track')).toBeVisible();
    });
  });
});
