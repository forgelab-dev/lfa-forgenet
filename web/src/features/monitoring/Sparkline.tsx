import { segments, type Point } from './series';

export interface Line {
  points: Point[];
  /** Couleur CSS du tracé, de préférence un token (`var(--color-ok-400)`). */
  color: string;
  /** Remplissage léger sous la courbe. */
  area?: boolean;
}

interface Props {
  lines: Line[];
  /** Début et fin de la fenêtre affichée, en millisecondes. */
  from: number;
  to: number;
  /** Haut de l'échelle ; par défaut la plus haute valeur, avec une marge. */
  max?: number;
  /**
   * Haut d'échelle minimal quand `max` est automatique : une courbe presque
   * plate (0,3 % de processeur) ne doit pas remplir toute la hauteur.
   */
  floor?: number;
  height?: number;
  /** Ce que montre la courbe, pour les lecteurs d'écran. */
  label: string;
  className?: string;
}

const WIDTH = 300;

/**
 * Courbe en SVG, sans bibliothèque : l'axe des temps est réel, si bien qu'une
 * coupure de mesures se voit comme un trou plutôt que d'être lissée.
 */
export function Sparkline({
  lines,
  from,
  to,
  max,
  floor = 1,
  height = 44,
  label,
  className = '',
}: Props) {
  const top = max ?? Math.max(floor, ...lines.flatMap((l) => l.points.map((p) => p.v))) * 1.15;
  const x = (t: number) => ((t - from) / Math.max(1, to - from)) * WIDTH;
  const y = (v: number) => height - 1 - (Math.min(v, top) / top) * (height - 2);

  return (
    <svg
      viewBox={`0 0 ${WIDTH} ${height}`}
      preserveAspectRatio="none"
      role="img"
      aria-label={label}
      className={`block w-full ${className}`}
      style={{ height }}
    >
      {lines.map((line, i) =>
        segments(line.points.filter((p) => p.t >= from)).map((seg, j) => {
          const d = seg
            .map((p, k) => `${k ? 'L' : 'M'}${x(p.t).toFixed(1)},${y(p.v).toFixed(1)}`)
            .join('');
          const first = seg[0]!;
          const last = seg[seg.length - 1]!;
          return (
            <g key={`${i}-${j}`}>
              {line.area && seg.length > 1 && (
                <path
                  d={`${d}L${x(last.t).toFixed(1)},${height}L${x(first.t).toFixed(1)},${height}Z`}
                  fill={line.color}
                  opacity={0.12}
                />
              )}
              <path
                d={seg.length > 1 ? d : `${d}l0.1,0`}
                fill="none"
                stroke={line.color}
                strokeWidth={1.5}
                strokeLinejoin="round"
                strokeLinecap="round"
                vectorEffect="non-scaling-stroke"
              />
            </g>
          );
        })
      )}
    </svg>
  );
}
