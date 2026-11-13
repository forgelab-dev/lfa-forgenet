import { describe, expect, it } from 'vitest';

import { byteLength, projectIdPattern, suggestProjectId, toProjectId } from './projectId';

describe('toProjectId', () => {
  it('rend un nom libre conforme à la règle du serveur', () => {
    expect(toProjectId('Ma Boutique (prod)')).toBe('ma-boutique-prod');
    expect(toProjectId('Café Médina')).toBe('cafe-medina');
    expect(toProjectId('--api_v2--')).toBe('api_v2');
    expect(toProjectId('x'.repeat(80))).toHaveLength(63);
    for (const name of ['Ma Boutique (prod)', 'Café Médina', '--api_v2--', 'x'.repeat(80)]) {
      expect(toProjectId(name)).toMatch(projectIdPattern);
    }
    expect(toProjectId('!!!')).toBe('');
  });
});

describe('suggestProjectId', () => {
  it('préfère le champ name: du fichier', () => {
    expect(suggestProjectId('name: "Boutique Dakar"\nservices:\n  web: {}\n', 'stack.yml')).toBe(
      'boutique-dakar'
    );
    expect(suggestProjectId('name: vitrine # le site\nservices: {}\n')).toBe('vitrine');
  });

  it('ignore un name: imbriqué, propre à un service', () => {
    expect(suggestProjectId('services:\n  web:\n    name: pas-moi\n', 'Ma Boutique.yml')).toBe(
      'ma-boutique'
    );
  });

  it('se rabat sur le nom du fichier, sauf s’il est générique', () => {
    expect(suggestProjectId('services: {}', 'blog.prod.yaml')).toBe('blog');
    expect(suggestProjectId('services: {}', 'docker-compose.yml')).toBe('');
    expect(suggestProjectId('services: {}')).toBe('');
  });
});

describe('byteLength', () => {
  it('compte les octets, pas les caractères', () => {
    expect(byteLength('é')).toBe(2);
  });
});
