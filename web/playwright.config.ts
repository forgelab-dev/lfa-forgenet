import { defineConfig, devices } from '@playwright/test';

/**
 * Tests de bout en bout du dashboard sur le vrai serveur (#116, lot 4).
 *
 * `e2e/server.ts` compile forgenet-server, crée un administrateur d'essai et
 * lance le serveur, qui sert `web/dist` : lancer `npm run build` avant
 * `npm run e2e`. Le test installe une vraie stack Docker, d'où un seul worker
 * et des délais larges.
 *
 * En local, FORGENET_E2E_CHROME=1 utilise le Chrome déjà installé au lieu du
 * Chromium de Playwright (pas de téléchargement).
 */
export default defineConfig({
  testDir: './e2e',
  globalSetup: './e2e/server.ts',
  timeout: 180_000,
  expect: { timeout: 15_000 },
  fullyParallel: false,
  workers: 1,
  forbidOnly: !!process.env.CI,
  retries: 0,
  reporter: process.env.CI ? [['list'], ['html', { open: 'never' }]] : 'list',
  // L'adresse du serveur n'est connue qu'après le globalSetup : chaque test
  // la reprend par test.use({ baseURL }).
  use: {
    locale: 'fr-FR',
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
  },
  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
        ...(process.env.FORGENET_E2E_CHROME ? { channel: 'chrome' } : {}),
      },
    },
  ],
});
