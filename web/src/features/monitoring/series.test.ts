import { describe, expect, it } from 'vitest';

import { formatAgo, formatBytes, formatPercent, formatRate, formatUptime } from './format';
import { HISTORY_MS, merge, peak, rate, rates, recent, segments, toneOf } from './series';

describe('mise en forme', () => {
  it('écrit les tailles en unités binaires, à la française', () => {
    expect(formatBytes(512)).toBe('512 o');
    expect(formatBytes(1536)).toBe('1,5 Ko');
    expect(formatBytes(5 * 1024 ** 3)).toBe('5,0 Go');
    expect(formatBytes(128 * 1024 ** 3)).toBe('128 Go');
    expect(formatBytes(-1)).toBe('—');
    expect(formatRate(1.2 * 1024 ** 2)).toBe('1,2 Mo/s');
    expect(formatRate(undefined)).toBe('—');
  });

  it('écrit les pourcentages, les durées et les anciennetés', () => {
    expect(formatPercent(37.4)).toBe('37 %');
    expect(formatPercent(6.24)).toBe('6,2 %');
    expect(formatUptime(12 * 86400 + 4 * 3600 + 59)).toBe('12 j 4 h');
    expect(formatUptime(3 * 3600 + 12 * 60)).toBe('3 h 12 min');
    expect(formatAgo(20_000)).toBe('il y a moins d’une minute');
    expect(formatAgo(6 * 60_000)).toBe('il y a 6 min');
  });
});

describe('courbes', () => {
  const now = 10 * HISTORY_MS;

  it('raccorde des points sans doublon et oublie ce qui dépasse l’heure', () => {
    const old = { t: now - HISTORY_MS - 1, v: 1 };
    const merged = merge(
      [old, { t: now - 10_000, v: 2 }],
      [
        { t: now - 10_000, v: 3 },
        { t: now - 20_000, v: 4 },
      ],
      now
    );
    expect(merged).toEqual([
      { t: now - 20_000, v: 4 },
      { t: now - 10_000, v: 3 },
    ]);
    expect(recent(merged, 15_000, now)).toEqual([{ t: now - 10_000, v: 3 }]);
  });

  it('tire un débit d’un compteur cumulé, sauf après une remise à zéro', () => {
    const counter = [
      { t: 0, v: 1000 },
      { t: 5000, v: 6000 },
      { t: 10_000, v: 100 }, // redémarrage
      { t: 15_000, v: 600 },
    ];
    expect(rates(counter)).toEqual([
      { t: 5000, v: 1000 },
      { t: 15_000, v: 100 },
    ]);
    expect(rate(undefined, { t: 1, v: 1 })).toBeUndefined();
    expect(rate({ t: 0, v: 0 }, { t: 2000, v: 4096 })).toBe(2048);
  });

  it('coupe une courbe là où les mesures manquent', () => {
    const points = [0, 5, 10, 15, 60, 65].map((s) => ({ t: s * 1000, v: 1 }));
    expect(segments(points).map((s) => s.length)).toEqual([4, 2]);
    expect(segments([])).toEqual([]);
    const mixed = [0, 10, 20, 30, 32, 34, 36].map((s) => ({ t: s * 1000, v: 1 }));
    expect(segments(mixed)).toHaveLength(1); // historique à 10 s puis direct à 2 s
  });

  it('trouve la pointe et applique les seuils de la mémoire et du disque', () => {
    expect(
      peak([
        { t: 1, v: 3 },
        { t: 2, v: 81 },
        { t: 3, v: 12 },
      ])
    ).toEqual({ t: 2, v: 81 });
    expect(toneOf(84.9)).toBe('ok');
    expect(toneOf(85)).toBe('warn');
    expect(toneOf(95)).toBe('danger');
    expect(toneOf(undefined)).toBe('ok');
  });
});
