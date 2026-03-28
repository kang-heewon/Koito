import { test as base, Page } from '@playwright/test';
import * as path from 'path';

type TestFixtures = {
  authPage: Page;
};

export const test = base.extend<TestFixtures>({
  authPage: async ({ browser }, use) => {
    const context = await browser.newContext({
      storageState: path.resolve(__dirname, '..', '.auth', 'user.json'),
    });
    const page = await context.newPage();
    await use(page);
    await context.close();
  },
});

export { expect } from '@playwright/test';
