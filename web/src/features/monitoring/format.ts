/**
 * Mise en forme des mesures, en français : virgule décimale, unités binaires
 * (Ko, Mo, Go) comme `docker stats` et `free`.
 */

const units = ['o', 'Ko', 'Mo', 'Go', 'To'];

function decimal(value: number, digits: number): string {
  return value.toLocaleString('fr-FR', {
    minimumFractionDigits: digits,
    maximumFractionDigits: digits,
  });
}

/** 5 368 709 120 → « 5,0 Go » ; une décimale sous 10, aucune au-delà. */
export function formatBytes(bytes: number): string {
  if (!Number.isFinite(bytes) || bytes < 0) return '—';
  let value = bytes;
  let unit = 0;
  while (value >= 1024 && unit < units.length - 1) {
    value /= 1024;
    unit++;
  }
  if (unit === 0) return `${Math.round(value)} o`;
  return `${decimal(value, value < 10 ? 1 : 0)} ${units[unit]}`;
}

/** Débit en octets par seconde : « 1,2 Mo/s ». */
export function formatRate(bytesPerSecond: number | undefined): string {
  if (bytesPerSecond === undefined) return '—';
  return `${formatBytes(bytesPerSecond)}/s`;
}

/** 37,42 → « 37 % » ; sous 10, une décimale (« 6,2 % »). */
export function formatPercent(value: number | undefined): string {
  if (value === undefined || !Number.isFinite(value)) return '—';
  return `${decimal(value, value < 10 ? 1 : 0)} %`;
}

/** 1 054 800 s → « 12 j 4 h » ; sous un jour « 3 h 12 min ». */
export function formatUptime(seconds: number): string {
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);
  if (days > 0) return `${days} j ${hours % 24} h`;
  if (hours > 0) return `${hours} h ${minutes % 60} min`;
  return `${minutes} min`;
}

/** Ancienneté d'une mesure : « il y a 6 min », « il y a moins d’une minute ». */
export function formatAgo(ms: number): string {
  const minutes = Math.round(ms / 60_000);
  if (minutes < 1) return 'il y a moins d’une minute';
  if (minutes < 60) return `il y a ${minutes} min`;
  return `il y a ${Math.round(minutes / 60)} h`;
}

/** Heure de la dernière mesure : « 10 h 42 ». */
export function formatClock(ms: number): string {
  const d = new Date(ms);
  return `${d.getHours()} h ${String(d.getMinutes()).padStart(2, '0')}`;
}
