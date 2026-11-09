import { Laptop, LockKeyhole, TriangleAlert } from 'lucide-react';

import { describeOrigin, type ChannelSecurity } from '../../lib/origin';

const channel: Record<ChannelSecurity, { label: string; tone: string; Icon: typeof LockKeyhole }> =
  {
    encrypted: { label: 'Chiffrée (HTTPS)', tone: 'text-ok-400', Icon: LockKeyhole },
    local: { label: 'Locale, sur cette machine', tone: 'text-iron-300', Icon: Laptop },
    clear: {
      label: 'Non chiffrée : le mot de passe circule en clair sur le réseau',
      tone: 'text-warn-400',
      Icon: TriangleAlert,
    },
  };

interface Props {
  location?: Pick<Location, 'protocol' | 'host' | 'hostname'>;
  className?: string;
}

/**
 * Plaque du serveur : le dashboard est servi par le serveur qu'il administre,
 * la plaque dit donc lequel, et si la liaison est chiffrée.
 */
export function ServerPlate({ location = window.location, className = '' }: Props) {
  const origin = describeOrigin(location);
  const { label, tone, Icon } = channel[origin.security];
  return (
    <dl
      className={`grid grid-cols-[auto_1fr] gap-x-5 gap-y-2 rounded-lg border border-iron-700 bg-iron-900/80 px-4 py-3 text-sm ${className}`}
    >
      <dt className="text-iron-300">Serveur</dt>
      <dd className="min-w-0 truncate font-mono text-iron-100">{origin.host}</dd>
      <dt className="text-iron-300">Liaison</dt>
      <dd className={`flex items-start gap-1.5 ${tone}`}>
        <Icon aria-hidden="true" className="mt-0.5 size-4 shrink-0" />
        <span>{label}</span>
      </dd>
    </dl>
  );
}
