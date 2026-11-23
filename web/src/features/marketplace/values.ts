import type { Template, TemplateVariable } from '../../gen/forgenet/v1/marketplace_pb';

/** Taille maximale d'une valeur, comme internal/marketplace (maxValueSize). */
export const maxValueBytes = 4096;

/** Valeurs de départ du formulaire : celles par défaut du manifeste. */
export function initialValues(variables: TemplateVariable[]): Record<string, string> {
  return Object.fromEntries(variables.map((v) => [v.name, v.defaultValue]));
}

/** Vrai quand le serveur sait remplir un champ laissé vide. */
function mayBeEmpty(v: TemplateVariable): boolean {
  return v.defaultValue !== '' || v.generate > 0 || v.optional;
}

/**
 * Vérifie une valeur comme le fera le serveur (Manifest.Resolve), pour le
 * dire avant l'envoi. Le serveur reste juge : un motif écrit pour Go que le
 * navigateur ne sait pas lire est laissé à sa vérification.
 */
export function checkValue(v: TemplateVariable, value: string): string | undefined {
  if (value === '') {
    return mayBeEmpty(v) ? undefined : 'Saisissez une valeur.';
  }
  if (new TextEncoder().encode(value).length > maxValueBytes) {
    return '4 096 octets au plus.';
  }
  if (value.includes('\u0000')) {
    return 'Le caractère nul est interdit.';
  }
  if (v.pattern && compile(v.pattern)?.test(value) === false) {
    return `La valeur ne respecte pas le format attendu : ${v.pattern}`;
  }
  return undefined;
}

function compile(pattern: string): RegExp | undefined {
  try {
    return new RegExp(pattern, 'u');
  } catch {
    return undefined;
  }
}

/** Ce que devient le champ s'il reste vide ; rien quand il doit être rempli. */
export function emptyHint(v: TemplateVariable): string | undefined {
  if (v.generate > 0) {
    return `Laissez vide : un mot de passe de ${v.generate} caractères sera généré.`;
  }
  if (v.defaultValue !== '') {
    return `Laissé vide, il vaut « ${v.defaultValue} ».`;
  }
  return v.optional ? 'Facultatif.' : undefined;
}

/**
 * Le serveur nomme la variable qu'il refuse, jamais sa valeur :
 * « valeur d'installation refusée : NOM : raison ». Le refus s'affiche alors
 * sous ce champ, sans le préfixe ; les autres restent un message d'ensemble.
 */
export function refusal(
  message: string,
  variables: TemplateVariable[]
): { field?: string; message: string } {
  const match = / : ([A-Z][A-Z0-9_]*) : (.+)$/s.exec(message);
  if (match && variables.some((v) => v.name === match[1])) {
    const reason = match[2].trim();
    return { field: match[1], message: reason.charAt(0).toUpperCase() + reason.slice(1) };
  }
  return { message };
}

/** Ce que l'administrateur doit savoir des valeurs de ce template-ci. */
export function introduction(template: Template): string {
  if (template.variables.length === 0) {
    return 'Ce template ne demande aucune valeur : choisissez le nom de l’application et installez.';
  }
  const stored = 'Les valeurs sont chiffrées sur le serveur et ne se réaffichent pas.';
  return template.variables.some((v) => v.generate > 0)
    ? `${stored} Un mot de passe laissé vide est tiré au hasard et vous est montré une seule fois, à la fin de l’installation.`
    : stored;
}
