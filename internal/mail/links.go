package mail

import (
	"fmt"
	"strings"
)

// LocalDashboardURL est l'adresse du dashboard sans domaine : il n'écoute que
// sur la machine, et s'ouvre par un tunnel SSH (ADR-0009).
const LocalDashboardURL = "http://127.0.0.1:8080"

// tunnelHint dit comment ouvrir le dashboard d'un serveur sans domaine.
func tunnelHint() []block {
	return []block{
		note("Ce serveur n'a pas de domaine : ouvrez d'abord un tunnel SSH depuis votre ordinateur,"),
		cmd("ssh -N -L 8080:127.0.0.1:8080 utilisateur@adresse-du-serveur"),
		note("puis le lien ci-dessus dans votre navigateur."),
	}
}

// Intitulés des boutons de la version HTML (ADR-0016).
const (
	labelChoose = "Choisir mon mot de passe"
	labelRenew  = "Choisir un nouveau mot de passe"
)

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

// firstLoginBlocks accompagne un lien de première connexion. hasPassword
// dit si le compte a déjà un mot de passe : le lien en fait alors choisir un
// nouveau, qui remplace l'actuel.
func firstLoginBlocks(username, url string, local, hasPassword bool) []block {
	var bl []block
	if hasPassword {
		bl = append(bl,
			para(fmt.Sprintf("Le compte %s a déjà un mot de passe : ce lien en fait choisir un nouveau, qui remplacera "+
				"l'actuel et fermera les sessions du compte. Il est valable 24 heures et une seule fois :", username)),
			link(labelRenew, url))
	} else {
		bl = append(bl,
			para(fmt.Sprintf("Le compte %s attend son mot de passe. Choisissez-le en ouvrant ce lien, "+
				"valable 24 heures et une seule fois :", username)),
			link(labelChoose, url))
	}
	if local {
		bl = append(bl, tunnelHint()...)
	}
	return append(bl,
		note("Un nouveau lien annule celui-ci. Pour en refaire un, lancez sur le serveur :"),
		cmd("sudo forgenet user login-link --username "+username))
}

// FirstLoginText est le texte qui accompagne un lien de première connexion
// dans le terminal ou le fichier premier-acces ; l'e-mail dit la même chose.
func FirstLoginText(username, url string, local, hasPassword bool) string {
	return strings.TrimSuffix(blocksText(firstLoginBlocks(username, url, local, hasPassword)), "\n")
}

// FirstLoginMessage est l'e-mail de première connexion de l'administrateur,
// ou d'un lien refait pour un compte qui a déjà son mot de passe.
func FirstLoginMessage(to, username, url string, local, hasPassword bool) Message {
	intro := "ForgeNet est installé sur votre serveur."
	if hasPassword {
		intro = "Voici un lien d'accès à votre compte ForgeNet, demandé sur le serveur."
	}
	return letter{
		subject:  "Votre accès à ForgeNet",
		greeting: fmt.Sprintf("Bonjour %s,", username),
		blocks:   append([]block{para(intro)}, firstLoginBlocks(username, url, local, hasPassword)...),
	}.message(to)
}

// PasswordResetMessage est l'e-mail d'un mot de passe oublié. Il ne dit
// jamais l'ancien mot de passe, que ForgeNet ne connaît d'ailleurs pas.
func PasswordResetMessage(to, username, url string, local bool) Message {
	bl := []block{
		para(fmt.Sprintf("Quelqu'un, sans doute vous, a demandé un nouveau mot de passe pour le compte ForgeNet %s. "+
			"Choisissez-le en ouvrant ce lien, valable une heure et une seule fois :", username)),
		link(labelRenew, url),
	}
	if local {
		bl = append(bl, tunnelHint()...)
	}
	bl = append(bl, note("Si vous n'êtes pas à l'origine de cette demande, ignorez cet e-mail : votre mot de passe "+
		"ne change pas, et le lien expirera de lui-même."))
	return letter{
		subject:  "Nouveau mot de passe ForgeNet",
		greeting: fmt.Sprintf("Bonjour %s,", username),
		blocks:   bl,
	}.message(to)
}

// InvitationMessage est l'e-mail d'invitation d'un membre, envoyé au nom de
// l'administrateur qui invite.
func InvitationMessage(to, username, inviter, url string, local bool) Message {
	bl := []block{
		para(fmt.Sprintf("%s vous invite sur son serveur ForgeNet, avec le compte %s. "+
			"Choisissez votre mot de passe en ouvrant ce lien, valable 7 jours et une seule fois :", inviter, username)),
		link(labelChoose, url),
	}
	if local {
		bl = append(bl, tunnelHint()...)
	}
	bl = append(bl, note(fmt.Sprintf("Lien expiré ? Demandez-en un nouveau à %s. "+
		"Si vous n'attendiez pas cette invitation, ignorez cet e-mail.", inviter)))
	return letter{
		subject:  inviter + " vous invite sur ForgeNet",
		greeting: fmt.Sprintf("Bonjour %s,", username),
		blocks:   bl,
	}.message(to)
}

// EmailChangedMessage prévient l'ancienne adresse d'un compte que la sienne a
// changé (#100). actor est l'administrateur qui l'a changée ; le nom du
// compte s'il l'a changée lui-même ; vide sur la machine du serveur. La
// nouvelle adresse est masquée : l'ancienne était peut-être celle d'un
// inconnu.
func EmailChangedMessage(to, username, newEmail, actor string) Message {
	var what string
	switch actor {
	case username:
		what = fmt.Sprintf("Vous venez de changer l'adresse e-mail de votre compte ForgeNet %s.", username)
	case "":
		what = fmt.Sprintf("L'adresse e-mail de votre compte ForgeNet %s a été changée sur le serveur.", username)
	default:
		what = fmt.Sprintf("%s vient de changer l'adresse e-mail de votre compte ForgeNet %s.", actor, username)
	}
	warning := "Si vous n'attendiez pas ce changement, prévenez au plus vite un administrateur de " +
		"votre serveur : il peut rétablir l'adresse et fermer les sessions du compte."
	if actor == username {
		warning = "Si ce n'est pas vous, quelqu'un connaît votre mot de passe : prévenez au plus vite " +
			"un administrateur de votre serveur, qui peut rétablir l'adresse et fermer les sessions du compte."
	}
	return letter{
		subject:  "L'adresse de votre compte ForgeNet a changé",
		greeting: fmt.Sprintf("Bonjour %s,", username),
		blocks: []block{
			para(what + " " + fmt.Sprintf("Les prochains e-mails du compte partiront à %s ; celle-ci n'en recevra plus.",
				MaskAddress(newEmail))),
			note(warning),
		},
	}.message(to)
}

// TestMessage est l'e-mail d'essai envoyé à l'administrateur connecté, depuis
// la page Réglages, avec le relais host.
func TestMessage(to, username, host string) Message {
	return letter{
		subject:  "Essai d'envoi de ForgeNet",
		greeting: fmt.Sprintf("Bonjour %s,", username),
		blocks: []block{
			para(fmt.Sprintf("Cet e-mail confirme que votre serveur ForgeNet envoie bien ses messages par le relais %s. "+
				"Vous n'avez rien à faire.", host)),
		},
	}.message(to)
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
