import { test, expect } from '../fixtures/auth.fixture';
import { seedDefaultData } from '../fixtures/test-data';

test.describe('Wrapped and Rewind', () => {
  test.beforeEach(async ({ authPage }) => {
    const cookies = await authPage.context().cookies();
    const sessionCookie = cookies.find(c => c.name === 'session')?.value;
    if (sessionCookie) {
      await seedDefaultData(`session=${sessionCookie}`, 15);
    }
  });

  test.describe('Wrapped Page', () => {
    test('should display Wrapped page with year selector', async ({ authPage }) => {
      await authPage.goto('/wrapped');

      await expect(authPage).toHaveURL(/\/wrapped/);

      await expect(authPage.getByText(/Koito Wrapped/i)).toBeVisible();
      await expect(authPage.getByText(/Your Year in Music/i)).toBeVisible();

      await expect(authPage.getByText(/Choose year/i)).toBeVisible();
    });

    test('should display Wrapped sections', async ({ authPage }) => {
      await authPage.goto('/wrapped');

      await expect(authPage.getByText(/Totals/i)).toBeVisible();
      await expect(authPage.getByText(/Everything added up/i)).toBeVisible();
    });

    test('should display Top Tracks section', async ({ authPage }) => {
      await authPage.goto('/wrapped');

      await expect(authPage.getByText(/Top Tracks/i)).toBeVisible();
      await expect(authPage.getByText(/These songs owned your year/i)).toBeVisible();
    });

    test('should display Top Artists section', async ({ authPage }) => {
      await authPage.goto('/wrapped');

      await expect(authPage.getByText(/Top Artists/i)).toBeVisible();
      await expect(authPage.getByText(/The voices that stayed closest/i)).toBeVisible();
    });

    test('should display Top Albums section', async ({ authPage }) => {
      await authPage.goto('/wrapped');

      await expect(authPage.getByText(/Top Albums/i)).toBeVisible();
      await expect(authPage.getByText(/Front-to-back favorites/i)).toBeVisible();
    });

    test('should display Listening Hours section', async ({ authPage }) => {
      await authPage.goto('/wrapped');

      await expect(authPage.getByText(/Listening Hours/i)).toBeVisible();
      await expect(authPage.getByText(/Your day had a soundtrack/i)).toBeVisible();
    });

    test('should display Discovery section', async ({ authPage }) => {
      await authPage.goto('/wrapped');

      await expect(authPage.getByText(/Discovery/i)).toBeVisible();
      await expect(authPage.getByText(/You still made room for surprises/i)).toBeVisible();
    });

    test('should display Busiest Week section', async ({ authPage }) => {
      await authPage.goto('/wrapped');

      await expect(authPage.getByText(/Busiest Week/i)).toBeVisible();
      await expect(authPage.getByText(/One week went louder than the rest/i)).toBeVisible();
    });

    test('should display Most Replayed section', async ({ authPage }) => {
      await authPage.goto('/wrapped');

      await expect(authPage.getByText(/Most Replayed/i)).toBeVisible();
      await expect(authPage.getByText(/This one kept coming back/i)).toBeVisible();
    });

    test('should display Concentration section', async ({ authPage }) => {
      await authPage.goto('/wrapped');

      await expect(authPage.getByText(/Concentration/i)).toBeVisible();
      await expect(authPage.getByText(/You knew exactly what you loved/i)).toBeVisible();
    });
  });

  test.describe('Rewind Page', () => {
    test('should display Rewind page with navigation controls', async ({ authPage }) => {
      await authPage.goto('/rewind');

      await expect(authPage).toHaveURL(/\/rewind/);

      await expect(authPage.getByText(/Month/i)).toBeVisible();
      await expect(authPage.getByText(/Year/i)).toBeVisible();

      await expect(authPage.getByRole('button', { name: /Month 이전으로 이동/i })).toBeVisible();
      await expect(authPage.getByRole('button', { name: /Month 다음으로 이동/i })).toBeVisible();
      await expect(authPage.getByRole('button', { name: /Year 이전으로 이동/i })).toBeVisible();
      await expect(authPage.getByRole('button', { name: /Year 다음으로 이동/i })).toBeVisible();
    });

    test('should display Full Year as default month', async ({ authPage }) => {
      await authPage.goto('/rewind');

      await expect(authPage.getByText(/Full Year/i)).toBeVisible();
    });

    test('should navigate to specific year and month', async ({ authPage }) => {
      await authPage.goto('/rewind/2024/1');

      await expect(authPage).toHaveURL(/\/rewind\/2024\/1/);

      await expect(authPage.getByText(/January/i)).toBeVisible();
      await expect(authPage.getByText(/2024/i)).toBeVisible();
    });

    test('should navigate to year only', async ({ authPage }) => {
      await authPage.goto('/rewind/2024');

      await expect(authPage).toHaveURL(/\/rewind\/2024/);

      await expect(authPage.getByText(/Full Year/i)).toBeVisible();
    });

    test('should switch month using navigation', async ({ authPage }) => {
      await authPage.goto('/rewind');

      const nextMonthButton = authPage.getByRole('button', { name: /Month 다음으로 이동/i });
      await nextMonthButton.click();

      await expect(authPage.getByText(/January/i)).toBeVisible();
    });

    test('should switch year using navigation', async ({ authPage }) => {
      await authPage.goto('/rewind');

      const prevYearButton = authPage.getByRole('button', { name: /Year 이전으로 이동/i });
      await prevYearButton.click();

      await expect(authPage.getByText(/2024/i)).toBeVisible();
    });

    test('should display Rewind stats sections', async ({ authPage }) => {
      await authPage.goto('/rewind');

      await expect(authPage.getByText(/Total Listens/i)).toBeVisible();
      await expect(authPage.getByText(/Unique Artists/i)).toBeVisible();
    });

    test('should display Top Artists in Rewind', async ({ authPage }) => {
      await authPage.goto('/rewind');

      await expect(authPage.getByRole('heading', { name: /Top Artists/i })).toBeVisible();
    });

    test('should display Top Tracks in Rewind', async ({ authPage }) => {
      await authPage.goto('/rewind');

      await expect(authPage.getByRole('heading', { name: /Top Tracks/i })).toBeVisible();
    });
  });

  test.describe('Empty Data Periods', () => {
    test('should handle data from past year gracefully', async ({ authPage }) => {
      await authPage.goto('/rewind/2020');

      await expect(authPage).toHaveURL(/\/rewind\/2020/);

      await expect(authPage.getByText(/Total Listens/i)).toBeVisible();
      await expect(authPage.getByText(/Unique Artists/i)).toBeVisible();
    });

    test('should handle month with no data', async ({ authPage }) => {
      await authPage.goto('/rewind/2020/6');

      await expect(authPage).toHaveURL(/\/rewind\/2020\/6/);

      await expect(authPage.getByText(/Total Listens/i)).toBeVisible();
    });

    test('should display empty state for future periods', async ({ authPage }) => {
      await authPage.goto('/rewind/2030');

      await expect(authPage).toHaveURL(/\/rewind\/2030/);

      await expect(authPage.getByText(/Total Listens/i)).toBeVisible();
    });
  });
});
