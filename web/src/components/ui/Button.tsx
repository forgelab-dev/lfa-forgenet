import type { ButtonHTMLAttributes } from 'react';

type Variant = 'primary' | 'quiet' | 'danger';

interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
}

const base =
  'inline-flex min-h-11 items-center justify-center gap-2 rounded-md px-4 text-sm font-semibold ' +
  'transition-[background-color,border-color,color,transform] duration-150 ease-out-quart ' +
  'active:translate-y-px disabled:cursor-not-allowed disabled:opacity-55 disabled:active:translate-y-0';

const variants: Record<Variant, string> = {
  primary: 'bg-gold-400 text-iron-950 hover:bg-gold-300 disabled:hover:bg-gold-400',
  quiet:
    'border border-iron-700 text-iron-100 hover:border-iron-500 hover:bg-iron-850 ' +
    'disabled:hover:border-iron-700 disabled:hover:bg-transparent',
  danger: 'bg-danger-400 text-iron-950 hover:bg-danger-400/85 disabled:hover:bg-danger-400',
};

/**
 * Bouton du dashboard : `primary` pour l'action principale d'un écran, une
 * seule à la fois ; `danger` pour confirmer une action qui retire quelque chose.
 */
export function Button({ variant = 'primary', className = '', type = 'button', ...props }: Props) {
  return <button type={type} className={`${base} ${variants[variant]} ${className}`} {...props} />;
}
