package mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDashboardURL(t *testing.T) {
	assert.Equal(t, "https://forgenet.exemple.org", DashboardURL(" forgenet.exemple.org "))
	assert.Equal(t, "http://127.0.0.1:8080", DashboardURL(""))
}

func TestAccessURL_KeepsTheTokenInTheFragment(t *testing.T) {
	assert.Equal(t, "https://forgenet.exemple.org/acces#jeton", AccessURL("https://forgenet.exemple.org/", "jeton"))
}

func TestFirstLoginMessage(t *testing.T) {
	msg := FirstLoginMessage("ana@exemple.org", "ana", "https://f.exemple.org/acces#jeton", false, false)
	assert.Equal(t, "ana@exemple.org", msg.To)
	assert.Contains(t, msg.Body, "Bonjour ana,")
	assert.Contains(t, msg.Body, "  https://f.exemple.org/acces#jeton\n")
	assert.Contains(t, msg.Body, "sudo forgenet user login-link --username ana")
	assert.NotContains(t, msg.Body, "tunnel")

	local := FirstLoginText("ana", "http://127.0.0.1:8080/acces#jeton", true, false)
	assert.Contains(t, local, "ssh -N -L 8080:127.0.0.1:8080")
}

func TestPasswordResetMessage(t *testing.T) {
	msg := PasswordResetMessage("ana@exemple.org", "ana", "https://f.exemple.org/acces#jeton", false)
	assert.Equal(t, "ana@exemple.org", msg.To)
	assert.Contains(t, msg.Body, "  https://f.exemple.org/acces#jeton\n")
	assert.Contains(t, msg.Body, "valable une heure")
	assert.Contains(t, msg.Body, "ignorez cet e-mail")
	assert.NotContains(t, msg.Body, "ssh -N")
	assert.Contains(t, PasswordResetMessage("a@b.org", "a", "u", true).Body, "ssh -N")
}

func TestInvitationMessage(t *testing.T) {
	msg := InvitationMessage("bio@exemple.org", "bio", "ana", "https://f.exemple.org/acces#jeton", false)
	assert.Equal(t, "bio@exemple.org", msg.To)
	assert.Equal(t, "ana vous invite sur ForgeNet", msg.Subject)
	assert.Contains(t, msg.Body, "ana vous invite sur son serveur ForgeNet, avec le compte bio.")
	assert.Contains(t, msg.Body, "  https://f.exemple.org/acces#jeton\n")
	assert.Contains(t, msg.Body, "valable 7 jours")
	assert.NotContains(t, msg.Body, "ssh -N")
	assert.Contains(t, InvitationMessage("a@b.org", "a", "c", "u", true).Body, "ssh -N")
}

func TestFirstLoginText_AccountWithAPassword(t *testing.T) {
	// #100 : le lien remplace un mot de passe existant, il ne faut pas dire
	// que le compte en attend un.
	text := FirstLoginText("ana", "https://f.exemple.org/acces#jeton", false, true)
	assert.NotContains(t, text, "attend son mot de passe")
	assert.Contains(t, text, "a déjà un mot de passe")
	assert.Contains(t, text, "  https://f.exemple.org/acces#jeton\n")

	msg := FirstLoginMessage("ana@exemple.org", "ana", "u", false, true)
	assert.NotContains(t, msg.Body, "est installé")
	assert.Contains(t, msg.Body, "a déjà un mot de passe")
}

func TestEmailChangedMessage(t *testing.T) {
	msg := EmailChangedMessage("ancienne@exemple.org", "bio", "nouvelle@autre.org", "ana")
	assert.Equal(t, "ancienne@exemple.org", msg.To, "l'avis part à l'ancienne adresse")
	assert.Contains(t, msg.Body, "Bonjour bio,")
	assert.Contains(t, msg.Body, "ana vient de changer l'adresse e-mail de votre compte ForgeNet bio.")
	assert.Contains(t, msg.Body, "n***@autre.org", "la nouvelle adresse est masquée")
	assert.NotContains(t, msg.Body, "nouvelle@autre.org")
	assert.Contains(t, msg.Body, "prévenez")

	self := EmailChangedMessage("a@exemple.org", "bio", "b@autre.org", "bio")
	assert.Contains(t, self.Body, "Vous venez de changer l'adresse e-mail de votre compte ForgeNet bio.")
	assert.Contains(t, self.Body, "votre mot de passe")

	server := EmailChangedMessage("a@exemple.org", "bio", "b@autre.org", "")
	assert.Contains(t, server.Body, "a été changée sur le serveur")
}

func TestMaskAddress(t *testing.T) {
	assert.Equal(t, "n***@autre.org", MaskAddress("nouvelle@autre.org"))
	assert.Equal(t, "b***@autre.org", MaskAddress("b@autre.org"))
	assert.Equal(t, "***", MaskAddress("sans-arobase"))
}
