/** Saisies du formulaire de changement de mot de passe. */
export interface PasswordForm {
  current: string;
  next: string;
  confirm: string;
}

export type PasswordProblems = Partial<Record<keyof PasswordForm, string>>;

/** Mêmes bornes que le serveur (internal/auth/password.go). */
export const MIN_PASSWORD_CHARS = 8;
export const MAX_PASSWORD_BYTES = 72;

/**
 * Contrôle le formulaire avant l'envoi. Le serveur refait chaque contrôle :
 * celui-ci évite seulement un aller-retour pour une erreur visible d'ici.
 * Le minimum se compte en caractères (« é » en vaut un), le maximum en octets,
 * qui est la limite réelle de bcrypt.
 */
export function checkPasswordForm(form: PasswordForm): PasswordProblems {
  const problems: PasswordProblems = { ...checkChosenPassword(form.next, form.confirm) };
  if (!form.current) {
    problems.current = 'Saisissez votre mot de passe actuel.';
  }
  if (!problems.next && form.current && form.next === form.current) {
    problems.next = 'Le nouveau mot de passe doit être différent de l’actuel.';
    delete problems.confirm;
  }
  return problems;
}

/**
 * Contrôle un mot de passe choisi et sa confirmation, sans mot de passe
 * actuel : c'est le cas d'un lien à usage unique (page /acces).
 */
export function checkChosenPassword(
  next: string,
  confirm: string
): Partial<Record<'next' | 'confirm', string>> {
  if ([...next].length < MIN_PASSWORD_CHARS) {
    return {
      next: `Le nouveau mot de passe doit comporter au moins ${MIN_PASSWORD_CHARS} caractères.`,
    };
  }
  if (new TextEncoder().encode(next).length > MAX_PASSWORD_BYTES) {
    return {
      next: `Le nouveau mot de passe doit comporter au plus ${MAX_PASSWORD_BYTES} octets (environ ${MAX_PASSWORD_BYTES} caractères sans accent).`,
    };
  }
  return confirm === next ? {} : { confirm: 'Les deux saisies diffèrent.' };
}

/** Phrase de réussite, avec le nombre de sessions fermées ailleurs. */
export function successMessage(revoked: number): string {
  if (revoked <= 0) {
    return 'Mot de passe changé. Aucune autre session n’était ouverte.';
  }
  const sessions =
    revoked === 1 ? '1 autre session a été fermée' : `${revoked} autres sessions ont été fermées`;
  return `Mot de passe changé. ${sessions} : il faudra s’y reconnecter avec le nouveau.`;
}
