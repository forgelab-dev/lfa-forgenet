import { MailTLS, type GetMailStatusResponse } from '../../gen/forgenet/v1/mail_pb';

/** Saisies du formulaire du relais SMTP (#102). Le port reste du texte tant qu'il est saisi. */
export interface RelayForm {
  host: string;
  port: string;
  tls: MailTLS;
  username: string;
  password: string;
  from: string;
}

export type RelayProblems = Partial<Record<keyof RelayForm, string>>;

/** Port habituel de chaque protection : changer de protection le propose. */
export const defaultPort: Record<MailTLS, number> = {
  [MailTLS.MAIL_TLS_STARTTLS]: 587,
  [MailTLS.MAIL_TLS_IMPLICIT]: 465,
  [MailTLS.MAIL_TLS_NONE]: 25,
  [MailTLS.MAIL_TLS_UNSPECIFIED]: 587,
};

/** Relais acceptés sans protection, comme le serveur (config.SMTPConfig). */
const localHosts = new Set(['localhost', '127.0.0.1', '::1']);

/** Le formulaire part des réglages en vigueur, sans le mot de passe. */
export function formFromStatus(status: GetMailStatusResponse): RelayForm {
  if (!status.host) {
    return {
      host: '',
      port: '587',
      tls: MailTLS.MAIL_TLS_STARTTLS,
      username: '',
      password: '',
      from: '',
    };
  }
  return {
    host: status.host,
    port: String(status.port),
    tls: status.tls === MailTLS.MAIL_TLS_UNSPECIFIED ? MailTLS.MAIL_TLS_STARTTLS : status.tls,
    username: status.username,
    password: '',
    from: status.from,
  };
}

/** Un mot de passe enregistré peut être gardé : le champ vide le garde. */
export function canKeepPassword(status: GetMailStatusResponse): boolean {
  return status.authenticated && !status.passwordUnreadable;
}

/**
 * Contrôle le formulaire avant l'envoi. Le serveur refait chaque contrôle :
 * celui-ci évite un aller-retour, et un e-mail d'essai, pour une erreur
 * visible d'ici.
 */
export function checkRelayForm(form: RelayForm, keepable: boolean): RelayProblems {
  const problems: RelayProblems = {};
  const host = form.host.trim().toLowerCase();
  if (!host) {
    problems.host = 'Indiquez le relais de votre fournisseur, par exemple smtp.exemple.org.';
  } else if (/[\s/@]/.test(host)) {
    problems.host = 'Le relais s’écrit comme un nom d’hôte (smtp.exemple.org), sans espace ni @.';
  }
  const port = Number(form.port.trim());
  if (!Number.isInteger(port) || port < 1 || port > 65535) {
    problems.port = 'Le port est un nombre entre 1 et 65535 (587 ou 465 le plus souvent).';
  }
  if (form.tls === MailTLS.MAIL_TLS_NONE && host && !localHosts.has(host)) {
    problems.tls =
      'Sans protection, le mot de passe et les e-mails passeraient en clair : réservé à un relais sur ce serveur (localhost).';
  }
  if (form.username.trim() && !form.password && !keepable) {
    problems.password = 'Saisissez le mot de passe du relais.';
  }
  if (/[\r\n]/.test(form.password)) {
    problems.password = 'Le mot de passe ne contient pas de saut de ligne.';
  }
  const from = form.from.trim();
  if (!from) {
    problems.from = 'Indiquez l’adresse d’expéditeur des e-mails.';
  } else if (!/^[^\s@]+@[^\s@]+$/.test(from) && !/<[^\s@]+@[^\s@]+>$/.test(from)) {
    problems.from = 'L’expéditeur est une adresse : forgenet@exemple.org.';
  }
  return problems;
}

/** Demande envoyée au serveur ; le mot de passe vide est gardé s'il existe. */
export function requestFromForm(form: RelayForm, keepable: boolean) {
  const username = form.username.trim();
  const keepPassword = username !== '' && form.password === '' && keepable;
  return {
    host: form.host.trim().toLowerCase(),
    port: Number(form.port.trim()),
    tls: form.tls,
    username,
    from: form.from.trim(),
    keepPassword,
    password: keepPassword || username === '' ? '' : form.password,
  };
}
