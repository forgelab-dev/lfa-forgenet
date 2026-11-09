import tailwindcss from '@tailwindcss/vite';
import react from '@vitejs/plugin-react';
import { defineConfig } from 'vitest/config';

// Serveur ForgeNet local (port HTTP par défaut de forgenet-server).
const api = 'http://127.0.0.1:8080';

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    // Le dashboard et l'API partagent une origine (ADR-0006) : en développement,
    // Vite relaie les appels Connect et les flux WebSocket vers le serveur.
    proxy: {
      '^/forgenet\\.v1\\.': { target: api, changeOrigin: true },
      '/ws/': {
        target: api,
        ws: true,
        changeOrigin: true,
        // Le serveur refuse une origine différente de son hôte (ADR-0005) :
        // la poignée de main relayée porte celle du serveur, comme en production.
        configure: (proxy) => {
          proxy.on('proxyReqWs', (proxyReq) => proxyReq.setHeader('Origin', api));
        },
      },
    },
  },
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    include: ['src/**/*.test.{ts,tsx}'],
  },
});
