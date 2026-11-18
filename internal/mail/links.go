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
func FirstLoginText(username, url string, local bool) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Le compte %s attend son mot de passe. Choisissez-le en ouvrant ce lien,\n", username)
	b.WriteString("valable 24 heures et une seule fois :\n\n")
	fmt.Fprintf(&b, "  %s\n\n", url)
	if local {
		b.WriteString(tunnelHint + "\n")
	}
	b.WriteString("Un nouveau lien annule celui-ci. Pour en refaire un, lancez sur le serveur :\n")
	fmt.Fprintf(&b, "  sudo forgenet user login-link --username %s\n", username)
	return b.String()
}

// FirstLoginMessage est l'e-mail de première connexion de l'administrateur.
func FirstLoginMessage(to, username, url string, local bool) Message {
	return Message{
		To:      to,
		Subject: "Votre accès à ForgeNet",
		Body: fmt.Sprintf("Bonjour %s,\n\nForgeNet est installé sur votre serveur.\n", username) +
			FirstLoginText(username, url, local) +
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
