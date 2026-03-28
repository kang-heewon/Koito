import { FullConfig } from '@playwright/test';

async function globalSetup(config: FullConfig) {
  const maxAttempts = 30;
  const intervalMs = 2000;
  const timeoutMs = 60000;

  const startTime = Date.now();

  for (let attempt = 1; attempt <= maxAttempts; attempt++) {
    try {
      const response = await fetch('http://localhost:4110/apis/web/v1/health', {
        method: 'GET',
      });

      if (response.status === 200) {
        console.log('API is ready');
        return;
      }
    } catch {
    }

    const elapsed = Date.now() - startTime;
    if (elapsed >= timeoutMs) {
      throw new Error(`API health check timed out after ${timeoutMs}ms`);
    }

    console.log(`Waiting for API... (attempt ${attempt}/${maxAttempts})`);
    await new Promise(resolve => setTimeout(resolve, intervalMs));
  }

  throw new Error(`API health check failed after ${maxAttempts} attempts`);
}

export default globalSetup;
