/**
 * Le dashboard est servi par le serveur qu'il administre (ADR-0006) : son
 * origine dit donc à quel serveur on parle, et si le mot de passe y circule
 * chiffré.
 */
export type ChannelSecurity = 'encrypted' | 'local' | 'clear';

export interface ServerOrigin {
  host: string;
  security: ChannelSecurity;
}

export function describeOrigin(
  location: Pick<Location, 'protocol' | 'host' | 'hostname'>
): ServerOrigin {
  let security: ChannelSecurity = 'clear';
  if (location.protocol === 'https:') {
    security = 'encrypted';
  } else if (isLoopback(location.hostname)) {
    security = 'local';
  }
  return { host: location.host, security };
}

function isLoopback(hostname: string): boolean {
  return (
    hostname === 'localhost' ||
    hostname.endsWith('.localhost') ||
    hostname === '[::1]' ||
    hostname === '::1' ||
    /^127\.\d{1,3}\.\d{1,3}\.\d{1,3}$/.test(hostname)
  );
}
