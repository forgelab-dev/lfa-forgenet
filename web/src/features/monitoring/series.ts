import { create } from '@bufbuild/protobuf';
import { timestampMs } from '@bufbuild/protobuf/wkt';

import {
  MetricKeySchema,
  type MetricKey,
  type MetricName,
  type MetricSeries,
} from '../../gen/forgenet/v1/metrics_pb';

/** Un point d'une courbe : date en millisecondes, valeur. */
export interface Point {
  t: number;
  v: number;
}

/** Le serveur garde une heure de mesures (ADR-0010) ; l'écran aussi. */
export const HISTORY_MS = 60 * 60_000;

/** Fenêtres proposées à l'écran. */
export const RANGES = [
  { label: '15 min', ms: 15 * 60_000 },
  { label: '1 h', ms: HISTORY_MS },
] as const;

export function key(name: MetricName): MetricKey {
  return create(MetricKeySchema, { key: { case: 'known', value: name } });
}

/** Les séries d'une réponse d'historique, rangées par métrique. */
export function byName(series: MetricSeries[]): Map<MetricName, Point[]> {
  const out = new Map<MetricName, Point[]>();
  for (const s of series) {
    if (s.key?.key.case !== 'known') continue;
    out.set(
      s.key.key.value,
      s.samples.flatMap((p) => (p.timestamp ? [{ t: timestampMs(p.timestamp), v: p.value }] : []))
    );
  }
  return out;
}

/**
 * Raccorde des points à une courbe : tri par date, un seul point par date,
 * et rien de plus ancien que l'heure gardée.
 */
export function merge(points: Point[], added: Point[], now = Date.now()): Point[] {
  const floor = now - HISTORY_MS;
  const byTime = new Map<number, number>();
  for (const p of points) if (p.t >= floor) byTime.set(p.t, p.v);
  for (const p of added) if (p.t >= floor) byTime.set(p.t, p.v);
  return [...byTime.entries()].sort((a, b) => a[0] - b[0]).map(([t, v]) => ({ t, v }));
}

/** Les points des `ms` dernières millisecondes. */
export function recent(points: Point[], ms: number, now = Date.now()): Point[] {
  const floor = now - ms;
  const first = points.findIndex((p) => p.t >= floor);
  return first < 0 ? [] : points.slice(first);
}

/**
 * Débit d'un compteur cumulé, en unités par seconde, entre deux points
 * successifs. Un compteur remis à zéro (redémarrage) ne donne pas de point.
 */
export function rates(counter: Point[]): Point[] {
  const out: Point[] = [];
  for (let i = 1; i < counter.length; i++) {
    const a = counter[i - 1]!;
    const b = counter[i]!;
    const seconds = (b.t - a.t) / 1000;
    if (seconds <= 0 || b.v < a.v) continue;
    out.push({ t: b.t, v: (b.v - a.v) / seconds });
  }
  return out;
}

/** Débit entre deux relevés d'un compteur, ou rien si le compteur est reparti de zéro. */
export function rate(prev: Point | undefined, cur: Point): number | undefined {
  if (!prev) return undefined;
  const seconds = (cur.t - prev.t) / 1000;
  if (seconds <= 0 || cur.v < prev.v) return undefined;
  return (cur.v - prev.v) / seconds;
}

/** Valeur la plus haute d'une courbe, et sa date. */
export function peak(points: Point[]): Point | undefined {
  let best: Point | undefined;
  for (const p of points) if (!best || p.v > best.v) best = p;
  return best;
}

/**
 * Seuils d'alerte de la mémoire et du disque : orange à 85 %, rouge à 95 %.
 * Le processeur n'en a pas : ses pointes sont normales.
 */
export type Tone = 'ok' | 'warn' | 'danger';

export function toneOf(percent: number | undefined): Tone {
  if (percent === undefined) return 'ok';
  if (percent >= 95) return 'danger';
  if (percent >= 85) return 'warn';
  return 'ok';
}

/**
 * Écart en dessous duquel deux mesures se suivent toujours : l'historique
 * arrive toutes les 10 s, le direct d'un conteneur toutes les 2 s, et leur
 * raccord n'est pas un trou.
 */
const MIN_GAP_MS = 30_000;

/**
 * Coupe une courbe là où les mesures manquent : un écart de plus de 30 s, et
 * plus de trois fois plus long que l'écart habituel, est un trou (serveur
 * arrêté, liaison coupée).
 */
export function segments(points: Point[]): Point[][] {
  if (points.length < 2) return points.length ? [points] : [];
  const gaps = points
    .slice(1)
    .map((p, i) => p.t - points[i]!.t)
    .sort((a, b) => a - b);
  const usual = gaps[Math.floor(gaps.length / 2)]!;
  const out: Point[][] = [[points[0]!]];
  for (let i = 1; i < points.length; i++) {
    if (points[i]!.t - points[i - 1]!.t > Math.max(3 * usual, MIN_GAP_MS)) out.push([]);
    out[out.length - 1]!.push(points[i]!);
  }
  return out;
}
