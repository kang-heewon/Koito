import { test, expect } from '../fixtures/auth.fixture';
import { seedDefaultData, generateRandomListens } from '../fixtures/test-data';
import { TEST_TRACKS } from '../fixtures/test-constants';

test.describe('Browse and Navigation', () => {
  test.beforeEach(async ({ authPage }) => {
    const cookies = await authPage.context().cookies();
    const sessionCookie = cookies.find(c => c.name === 'session')?.value;
    if (sessionCookie) {
      await seedDefaultData(`session=${sessionCookie}`, 15);
    }
  });

  test.describe('Home Dashboard', () => {
    test('should display AllTimeStats section', async ({ authPage }) => {
      await authPage.goto('/');

      await expect(authPage.getByText(/all time stats/i)).toBeVisible();
      await expect(authPage.getByText(/minutes listened/i)).toBeVisible();
      await expect(authPage.getByText(/plays/i)).toBeVisible();
      await expect(authPage.getByText(/tracks/i)).toBeVisible();
      await expect(authPage.getByText(/albums/i)).toBeVisible();
      await expect(authPage.getByText(/artists/i)).toBeVisible();
    });

    test('should display ActivityGrid section', async ({ authPage }) => {
      await authPage.goto('/');

      await expect(authPage.getByText(/activity/i)).toBeVisible();
    });

    test('should display PeriodSelector with all options', async ({ authPage }) => {
      await authPage.goto('/');

      await expect(authPage.getByText(/showing stats for/i)).toBeVisible();
      await expect(authPage.getByRole('button', { name: /day/i })).toBeVisible();
      await expect(authPage.getByRole('button', { name: /week/i })).toBeVisible();
      await expect(authPage.getByRole('button', { name: /month/i })).toBeVisible();
      await expect(authPage.getByRole('button', { name: /year/i })).toBeVisible();
      await expect(authPage.getByRole('button', { name: /all time/i })).toBeVisible();
    });

    test('should switch period and update content', async ({ authPage }) => {
      await authPage.goto('/');

      const monthButton = authPage.getByRole('button', { name: /month/i });
      await monthButton.click();

      await expect(authPage.getByRole('button', { name: /month/i })).toBeDisabled();
    });

    test('should display TopArtists section', async ({ authPage }) => {
      await authPage.goto('/');

      await expect(authPage.getByRole('heading', { name: /top artists/i })).toBeVisible();
    });

    test('should display TopAlbums section', async ({ authPage }) => {
      await authPage.goto('/');

      await expect(authPage.getByRole('heading', { name: /top albums/i })).toBeVisible();
    });

    test('should display TopTracks section', async ({ authPage }) => {
      await authPage.goto('/');

      await expect(authPage.getByRole('heading', { name: /top tracks/i })).toBeVisible();
    });

    test('should display LastPlays section', async ({ authPage }) => {
      await authPage.goto('/');

      await expect(authPage.getByRole('heading', { name: /last played/i })).toBeVisible();
    });
  });

  test.describe('Artist Detail Page', () => {
    test('should navigate to artist page from home', async ({ authPage }) => {
      await authPage.goto('/');

      const artistLink = authPage.getByRole('link', { name: /test artist/i }).first();
      await expect(artistLink).toBeVisible();
      await artistLink.click();

      await expect(authPage).toHaveURL(/\/artist\/\d+/);
    });

    test('should display artist information', async ({ authPage }) => {
      await authPage.goto('/');

      const artistLink = authPage.getByRole('link', { name: /test artist/i }).first();
      await artistLink.click();

      await expect(authPage.getByRole('heading', { name: /artist/i, level: 3 })).toBeVisible();
      await expect(authPage.getByText(/plays/i)).toBeVisible();
    });

    test('should display artist albums section', async ({ authPage }) => {
      await authPage.goto('/');

      const artistLink = authPage.getByRole('link', { name: /test artist/i }).first();
      await artistLink.click();

      await expect(authPage.getByText(/albums/i)).toBeVisible();
    });

    test('should display artist tracks section', async ({ authPage }) => {
      await authPage.goto('/');

      const artistLink = authPage.getByRole('link', { name: /test artist/i }).first();
      await artistLink.click();

      await expect(authPage.getByRole('heading', { name: /top tracks/i })).toBeVisible();
    });
  });

  test.describe('Album Detail Page', () => {
    test('should navigate to album page from artist page', async ({ authPage }) => {
      await authPage.goto('/');

      const artistLink = authPage.getByRole('link', { name: /test artist/i }).first();
      await artistLink.click();

      await expect(authPage).toHaveURL(/\/artist\/\d+/);

      const albumLink = authPage.getByRole('link', { name: /test album/i }).first();
      await expect(albumLink).toBeVisible();
      await albumLink.click();

      await expect(authPage).toHaveURL(/\/album\/\d+/);
    });

    test('should display album information', async ({ authPage }) => {
      await authPage.goto('/');

      const artistLink = authPage.getByRole('link', { name: /test artist/i }).first();
      await artistLink.click();

      const albumLink = authPage.getByRole('link', { name: /test album/i }).first();
      await albumLink.click();

      await expect(authPage.getByRole('heading', { name: /album/i, level: 3 })).toBeVisible();
      await expect(authPage.getByText(/plays/i)).toBeVisible();
    });

    test('should display album track list', async ({ authPage }) => {
      await authPage.goto('/');

      const artistLink = authPage.getByRole('link', { name: /test artist/i }).first();
      await artistLink.click();

      const albumLink = authPage.getByRole('link', { name: /test album/i }).first();
      await albumLink.click();

      await expect(authPage.getByRole('heading', { name: /top tracks/i })).toBeVisible();
    });
  });

  test.describe('Track Detail Page', () => {
    test('should navigate to track page from album page', async ({ authPage }) => {
      await authPage.goto('/');

      const artistLink = authPage.getByRole('link', { name: /test artist/i }).first();
      await artistLink.click();

      const albumLink = authPage.getByRole('link', { name: /test album/i }).first();
      await albumLink.click();

      const trackLink = authPage.getByRole('link', { name: /test track/i }).first();
      await expect(trackLink).toBeVisible();
      await trackLink.click();

      await expect(authPage).toHaveURL(/\/track\/\d+/);
    });

    test('should display track information', async ({ authPage }) => {
      await authPage.goto('/');

      const artistLink = authPage.getByRole('link', { name: /test artist/i }).first();
      await artistLink.click();

      const albumLink = authPage.getByRole('link', { name: /test album/i }).first();
      await albumLink.click();

      const trackLink = authPage.getByRole('link', { name: /test track/i }).first();
      await trackLink.click();

      await expect(authPage.getByRole('heading', { name: /track/i, level: 3 })).toBeVisible();
      await expect(authPage.getByText(/appears on/i)).toBeVisible();
      await expect(authPage.getByText(/plays/i)).toBeVisible();
    });
  });

  test.describe('Listens Page', () => {
    test('should navigate to listens page from navigation', async ({ authPage }) => {
      await authPage.goto('/');

      await authPage.goto('/listens');

      await expect(authPage).toHaveURL(/\/listens/);
    });

    test('should display listens list', async ({ authPage }) => {
      await authPage.goto('/listens');

      await expect(authPage.getByRole('heading', { name: /last played/i })).toBeVisible();

      const table = authPage.locator('table');
      await expect(table).toBeVisible();
    });

    test('should display listen entries with track and artist info', async ({ authPage }) => {
      await authPage.goto('/listens');

      const table = authPage.locator('table');
      await expect(table).toBeVisible();

      const rows = table.locator('tbody tr');
      await expect(rows.first()).toBeVisible();

      const firstRow = rows.first();
      await expect(firstRow.getByRole('link')).toBeVisible();
    });

    test('should have pagination controls', async ({ authPage }) => {
      await authPage.goto('/listens');

      const prevButton = authPage.getByRole('button', { name: /prev/i });
      const nextButton = authPage.getByRole('button', { name: /next/i });

      await expect(prevButton).toBeVisible();
      await expect(nextButton).toBeVisible();

      await expect(prevButton).toBeDisabled();
    });
  });

  test.describe('Navigation Flow', () => {
    test('should complete full navigation flow: Home → Artist → Album → Track', async ({ authPage }) => {
      await authPage.goto('/');
      await expect(authPage).toHaveURL('/');

      const artistLink = authPage.getByRole('link', { name: /test artist/i }).first();
      await artistLink.click();
      await expect(authPage).toHaveURL(/\/artist\/\d+/);

      const albumLink = authPage.getByRole('link', { name: /test album/i }).first();
      await albumLink.click();
      await expect(authPage).toHaveURL(/\/album\/\d+/);

      const trackLink = authPage.getByRole('link', { name: /test track/i }).first();
      await trackLink.click();
      await expect(authPage).toHaveURL(/\/track\/\d+/);
    });

    test('should navigate back using browser back button', async ({ authPage }) => {
      await authPage.goto('/');

      const artistLink = authPage.getByRole('link', { name: /test artist/i }).first();
      await artistLink.click();
      await expect(authPage).toHaveURL(/\/artist\/\d+/);

      await authPage.goBack();
      await expect(authPage).toHaveURL('/');
    });
  });
});
