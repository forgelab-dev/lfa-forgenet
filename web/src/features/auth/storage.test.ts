import { beforeEach, describe, expect, it } from 'vitest';

import { clearSession, createTokenStore, readSession, writeSession } from './storage';

describe('conservation de la session', () => {
  beforeEach(() => sessionStorage.clear());

  it('relit une session encore valable', () => {
    writeSession({ token: 'jeton', expiresAt: 2_000 });

    expect(readSession(1_000)).toEqual({ token: 'jeton', expiresAt: 2_000 });
  });

  it('oublie une session échue', () => {
    writeSession({ token: 'jeton', expiresAt: 2_000 });

    expect(readSession(2_000)).toBeUndefined();
    expect(sessionStorage.getItem('forgenet.session')).toBeNull();
  });

  it('ignore une valeur illisible ou incomplète', () => {
    sessionStorage.setItem('forgenet.session', '{pas du json');
    expect(readSession()).toBeUndefined();

    sessionStorage.setItem('forgenet.session', JSON.stringify({ token: '', expiresAt: 9e15 }));
    expect(readSession()).toBeUndefined();
  });

  it('tient le jeton en mémoire et dans le navigateur', () => {
    writeSession({ token: 'ancien', expiresAt: Date.now() + 60_000 });
    const store = createTokenStore();
    expect(store.get()).toBe('ancien');

    store.set({ token: 'nouveau', expiresAt: Date.now() + 60_000 });
    expect(store.get()).toBe('nouveau');
    expect(readSession()?.token).toBe('nouveau');

    store.clear();
    expect(store.get()).toBeUndefined();
    clearSession();
    expect(readSession()).toBeUndefined();
  });
});
