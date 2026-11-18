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
	msg := FirstLoginMessage("ana@exemple.org", "ana", "https://f.exemple.org/acces#jeton", false)
	assert.Equal(t, "ana@exemple.org", msg.To)
	assert.Contains(t, msg.Body, "Bonjour ana,")
	assert.Contains(t, msg.Body, "  https://f.exemple.org/acces#jeton\n")
	assert.Contains(t, msg.Body, "sudo forgenet user login-link --username ana")
	assert.NotContains(t, msg.Body, "tunnel")

	local := FirstLoginText("ana", "http://127.0.0.1:8080/acces#jeton", true)
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
