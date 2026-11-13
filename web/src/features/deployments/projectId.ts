/** Règle du serveur : c'est le nom de la stack Docker Compose. */
export const projectIdPattern = /^[a-z0-9][a-z0-9_-]{0,62}$/;

/** Taille maximale d'un fichier Compose acceptée par le serveur (1 Mio). */
export const maxComposeBytes = 1 << 20;

/** Noms de fichier génériques, qui ne disent rien du projet. */
const genericNames = new Set(['docker-compose', 'compose', 'stack', 'docker-stack']);

/**
 * Transforme un nom libre en identifiant de projet valide :
 * « Ma Boutique (prod) » devient « ma-boutique-prod ».
 */
export function toProjectId(name: string): string {
  return name
    .normalize('NFD')
    .replace(/[̀-ͯ]/g, '')
    .toLowerCase()
    .replace(/[^a-z0-9_-]+/g, '-')
    .replace(/-{2,}/g, '-')
    .replace(/^[-_]+/, '')
    .slice(0, 63)
    .replace(/[-_]+$/, '');
}

/**
 * Propose un identifiant : le champ `name:` du fichier Compose s'il existe,
 * sinon le nom du fichier importé, sauf s'il est générique.
 */
export function suggestProjectId(compose: string, fileName?: string): string {
  const declared = /^name:[ \t]*["']?([^"'#\r\n]+?)["']?[ \t]*(?:#.*)?$/m.exec(compose);
  if (declared) {
    const id = toProjectId(declared[1]);
    if (id) {
      return id;
    }
  }
  if (fileName) {
    const base = fileName.replace(/\.(ya?ml)$/i, '').replace(/\.(prod|production|dev)$/i, '');
    if (!genericNames.has(base.toLowerCase())) {
      return toProjectId(base);
    }
  }
  return '';
}

/** Taille en octets du texte une fois encodé, comme le compte le serveur. */
export function byteLength(text: string): number {
  return new TextEncoder().encode(text).length;
}
