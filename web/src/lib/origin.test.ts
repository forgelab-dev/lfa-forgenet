import { describe, expect, it } from 'vitest';

import { describeOrigin } from './origin';

const at = (protocol: string, host: string) => ({
  protocol,
  host,
  hostname: host.replace(/:\d+$/, ''),
});

describe('describeOrigin', () => {
  it('reconnaît une liaison chiffrée', () => {
    expect(describeOrigin(at('https:', 'forgenet.example.com'))).toEqual({
      host: 'forgenet.example.com',
      security: 'encrypted',
    });
  });

  it.each(['localhost:5173', '127.0.0.1:8080', 'app.localhost'])(
    'tient %s pour une liaison locale',
    (host) => {
      expect(describeOrigin(at('http:', host)).security).toBe('local');
    }
  );

  it('signale un mot de passe envoyé en clair sur le réseau', () => {
    expect(describeOrigin(at('http:', '192.168.1.20:8080')).security).toBe('clear');
  });
});
