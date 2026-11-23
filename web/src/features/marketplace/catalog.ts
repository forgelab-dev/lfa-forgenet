import type { Template } from '../../gen/forgenet/v1/marketplace_pb';

/**
 * Libellés de la liste fermée des catégories (internal/marketplace), dans
 * l'ordre où le filtre les propose.
 */
const categories: [string, string][] = [
  ['database', 'Bases de données'],
  ['cache', 'Cache'],
  ['storage', 'Stockage'],
  ['cms', 'CMS'],
  ['devtools', 'Outils de développement'],
  ['monitoring', 'Surveillance'],
  ['communication', 'Communication'],
  ['other', 'Autres'],
];

const labels = new Map(categories);
const order = new Map(categories.map(([id], i) => [id, i]));

/** Nom français d'une catégorie ; une catégorie ajoutée côté serveur garde son identifiant. */
export function categoryLabel(category: string): string {
  return labels.get(category) ?? category;
}

/** Les catégories présentes dans le catalogue : le filtre ne propose jamais un choix vide. */
export function categoriesOf(templates: Template[]): string[] {
  const present = [...new Set(templates.map((t) => t.category))];
  return present.sort((a, b) => (order.get(a) ?? order.size) - (order.get(b) ?? order.size));
}

/** Les templates d'une catégorie, ou tous sans filtre. */
export function inCategory(templates: Template[], category: string | undefined): Template[] {
  return category ? templates.filter((t) => t.category === category) : templates;
}

/**
 * Source d'une balise <img> pour l'icône du template. Une image SVG n'exécute
 * ni script ni lien, même si le fichier en contenait : c'est plus sûr que de
 * l'insérer dans la page, que le serveur l'ait vérifié ou non.
 */
export function iconSource(svg: string): string {
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`;
}
