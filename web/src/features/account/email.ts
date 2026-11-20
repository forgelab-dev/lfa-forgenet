/** Saisies du formulaire de changement d'adresse (#100). */
export interface EmailForm {
  email: string;
  password: string;
}

export type EmailProblems = Partial<Record<keyof EmailForm, string>>;

/** Même borne que le serveur (internal/auth, NormalizeEmail). */
export const MAX_EMAIL_CHARS = 254;

/**
 * Contrôle le formulaire avant l'envoi. Le serveur refait chaque contrôle,
 * plus strict (adresse nue, sans nom affiché) : celui-ci évite seulement un
 * aller-retour pour une erreur visible d'ici.
 */
export function checkEmailForm(form: EmailForm, currentEmail: string): EmailProblems {
  const problems: EmailProblems = {};
  const email = checkNewEmail(form.email, currentEmail);
  if (email) {
    problems.email = email;
  }
  if (!form.password) {
    problems.password = 'Saisissez votre mot de passe pour confirmer.';
  }
  return problems;
}

/**
 * Contrôle une nouvelle adresse, pour la page Compte comme pour la page
 * Équipe ; renvoie le problème trouvé, ou undefined.
 */
export function checkNewEmail(value: string, currentEmail: string): string | undefined {
  const email = value.trim().toLowerCase();
  if (!email) {
    return 'Saisissez la nouvelle adresse.';
  }
  if (email.length > MAX_EMAIL_CHARS || !/^[^\s@]+@[^\s@]+$/.test(email)) {
    return 'Cette adresse n’est pas valide : elle s’écrit nom@domaine, sans espace.';
  }
  return email === currentEmail.toLowerCase() ? 'C’est déjà l’adresse de ce compte.' : undefined;
}
