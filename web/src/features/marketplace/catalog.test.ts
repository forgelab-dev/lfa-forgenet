import { create } from '@bufbuild/protobuf';
import { describe, expect, it } from 'vitest';

import { TemplateSchema } from '../../gen/forgenet/v1/marketplace_pb';
import { categoriesOf, categoryLabel, iconSource, inCategory } from './catalog';

const templates = [
  create(TemplateSchema, { id: 'redis', name: 'Redis', category: 'cache' }),
  create(TemplateSchema, { id: 'postgres', name: 'PostgreSQL', category: 'database' }),
  create(TemplateSchema, { id: 'mysql', name: 'MySQL', category: 'database' }),
];

describe('catalog', () => {
  it('nomme chaque catégorie en français, et garde une inconnue telle quelle', () => {
    expect(categoryLabel('database')).toBe('Bases de données');
    expect(categoryLabel('devtools')).toBe('Outils de développement');
    expect(categoryLabel('nouvelle')).toBe('nouvelle');
  });

  it('ne propose que les catégories présentes, dans l’ordre de la liste fermée', () => {
    expect(categoriesOf(templates)).toEqual(['database', 'cache']);
  });

  it('filtre par catégorie, et garde tout sans filtre', () => {
    expect(inCategory(templates, 'database').map((t) => t.id)).toEqual(['postgres', 'mysql']);
    expect(inCategory(templates, undefined)).toHaveLength(3);
  });

  it('lit l’icône comme une image, jamais comme du code de la page', () => {
    const src = iconSource('<svg viewBox="0 0 1 1"><path d="M0 0h1"/></svg>');
    expect(src.startsWith('data:image/svg+xml;charset=utf-8,')).toBe(true);
    expect(decodeURIComponent(src.slice(src.indexOf(',') + 1))).toBe(
      '<svg viewBox="0 0 1 1"><path d="M0 0h1"/></svg>'
    );
    expect(src).not.toContain('<');
  });
});
