package mail

import (
	"fmt"
	"strings"
)

// LocalDashboardURL est l'adresse du dashboard sans domaine : il n'écoute que
// sur la machine, et s'ouvre par un tunnel SSH (ADR-0009).
const LocalDashboardURL = "http://127.0.0.1:8080"

// tunnelHint dit comment ouvrir le dashboard d'un serveur sans domaine.
const tunnelHint = "Ce serveur n'a pas de domaine : ouvrez d'abord un tunnel SSH depuis votre ordinateur,\n" +
	"  ssh -N -L 8080:127.0.0.1:8080 utilisateur@adresse-du-serveur\n" +
	"puis le lien ci-dessus dans votre navigateur.\n"

// DashboardURL est l'adresse publique du dashboard : son domaine en HTTPS,
// ou l'adresse locale à ouvrir par un tunnel.
func DashboardURL(dashboardHost string) string {
	if h := strings.TrimSpace(dashboardHost); h != "" {
		return "https://" + h
	}
	return LocalDashboardURL
}

// AccessURL est l'adresse d'un lien à usage unique. Le jeton voyage dans le
// fragment, que le navigateur n'envoie jamais au serveur : il n'entre dans
// aucun journal HTTP, ni celui de Caddy ni celui de ForgeNet (ADR-0014).
func AccessURL(dashboardURL, token string) string {
	return strings.TrimRight(dashboardURL, "/") + "/acces#" + token
}

// FirstLoginText est le texte qui accompagne un lien de première connexion,
// dans un e-mail comme dans le terminal ou le fichier premier-acces.
// hasPassword dit si le compte a déjà un mot de passe : le lien en fait alors
// choisir un nouveau, qui remplace l'actuel.
func FirstLoginText(username, url string, local, hasPassword bool) string {
	var b strings.Builder
	if hasPassword {
		fmt.Fprintf(&b, "Le compte %s a déjà un mot de passe : ce lien en fait choisir un nouveau, qui remplacera\n", username)
		b.WriteString("l'actuel et fermera les sessions du compte. Il est valable 24 heures et une seule fois :\n\n")
	} else {
		fmt.Fprintf(&b, "Le compte %s attend son mot de passe. Choisissez-le en ouvrant ce lien,\n", username)
		b.WriteString("valable 24 heures et une seule fois :\n\n")
	}
	fmt.Fprintf(&b, "  %s\n\n", url)
	if local {
		b.WriteString(tunnelHint + "\n")
	}
	b.WriteString("Un nouveau lien annule celui-ci. Pour en refaire un, lancez sur le serveur :\n")
	fmt.Fprintf(&b, "  sudo forgenet user login-link --username %s\n", username)
	return b.String()
}

// FirstLoginMessage est l'e-mail de première connexion de l'administrateur,
// ou d'un lien refait pour un compte qui a déjà son mot de passe.
func FirstLoginMessage(to, username, url string, local, hasPassword bool) Message {
	intro := "ForgeNet est installé sur votre serveur.\n"
	if hasPassword {
		intro = "Voici un lien d'accès à votre compte ForgeNet, demandé sur le serveur.\n"
	}
	return Message{
		To:      to,
		Subject: "Votre accès à ForgeNet",
		Body: fmt.Sprintf("Bonjour %s,\n\n", username) + intro +
			FirstLoginText(username, url, local, hasPassword) +
			"\n— ForgeNet\n",
	}
}

// PasswordResetMessage est l'e-mail d'un mot de passe oublié. Il ne dit
// jamais l'ancien mot de passe, que ForgeNet ne connaît d'ailleurs pas.
func PasswordResetMessage(to, username, url string, local bool) Message {
	var b strings.Builder
	fmt.Fprintf(&b, "Bonjour %s,\n\n", username)
	fmt.Fprintf(&b, "Quelqu'un, sans doute vous, a demandé un nouveau mot de passe pour le compte ForgeNet %s.\n", username)
	b.WriteString("Choisissez-le en ouvrant ce lien, valable une heure et une seule fois :\n\n")
	fmt.Fprintf(&b, "  %s\n\n", url)
	if local {
		b.WriteString(tunnelHint + "\n")
	}
	b.WriteString("Si vous n'êtes pas à l'origine de cette demande, ignorez cet e-mail : votre mot de passe\n")
	b.WriteString("ne change pas, et le lien expirera de lui-même.\n\n— ForgeNet\n")
	return Message{To: to, Subject: "Nouveau mot de passe ForgeNet", Body: b.String()}
}

// InvitationMessage est l'e-mail d'invitation d'un membre, envoyé au nom de
// l'administrateur qui invite.
func InvitationMessage(to, username, inviter, url string, local bool) Message {
	var b strings.Builder
	fmt.Fprintf(&b, "Bonjour %s,\n\n", username)
	fmt.Fprintf(&b, "%s vous invite sur son serveur ForgeNet, avec le compte %s.\n", inviter, username)
	b.WriteString("Choisissez votre mot de passe en ouvrant ce lien, valable 7 jours et une seule fois :\n\n")
	fmt.Fprintf(&b, "  %s\n\n", url)
	if local {
		b.WriteString(tunnelHint + "\n")
	}
	fmt.Fprintf(&b, "Lien expiré ? Demandez-en un nouveau à %s.\n", inviter)
	b.WriteString("Si vous n'attendiez pas cette invitation, ignorez cet e-mail.\n\n— ForgeNet\n")
	return Message{To: to, Subject: inviter + " vous invite sur ForgeNet", Body: b.String()}
}

// EmailChangedMessage prévient l'ancienne adresse d'un compte que la sienne a
// changé (#100). actor est l'administrateur qui l'a changée ; le nom du
// compte s'il l'a changée lui-même ; vide sur la machine du serveur. La
// nouvelle adresse est masquée : l'ancienne était peut-être celle d'un
// inconnu.
func EmailChangedMessage(to, username, newEmail, actor string) Message {
	var b strings.Builder
	fmt.Fprintf(&b, "Bonjour %s,\n\n", username)
	switch actor {
	case username:
		fmt.Fprintf(&b, "Vous venez de changer l'adresse e-mail de votre compte ForgeNet %s.\n", username)
	case "":
		fmt.Fprintf(&b, "L'adresse e-mail de votre compte ForgeNet %s a été changée sur le serveur.\n", username)
	default:
		fmt.Fprintf(&b, "%s vient de changer l'adresse e-mail de votre compte ForgeNet %s.\n", actor, username)
	}
	fmt.Fprintf(&b, "Les prochains e-mails du compte partiront à %s ; celle-ci n'en recevra plus.\n\n", MaskAddress(newEmail))
	if actor == username {
		b.WriteString("Si ce n'est pas vous, quelqu'un connaît votre mot de passe : prévenez au plus vite\n")
		b.WriteString("un administrateur de votre serveur, qui peut rétablir l'adresse et fermer les sessions du compte.\n")
	} else {
		b.WriteString("Si vous n'attendiez pas ce changement, prévenez au plus vite un administrateur de\n")
		b.WriteString("votre serveur : il peut rétablir l'adresse et fermer les sessions du compte.\n")
	}
	b.WriteString("\n— ForgeNet\n")
	return Message{To: to, Subject: "L'adresse de votre compte ForgeNet a changé", Body: b.String()}
}

// MaskAddress ne garde d'une adresse que sa première lettre et son domaine :
// assez pour la reconnaître, pas pour la lire en entier.
func MaskAddress(email string) string {
	at := strings.LastIndex(email, "@")
	if at < 1 {
		return "***"
	}
	r := []rune(email[:at])
	return string(r[0]) + "***" + email[at:]
}
