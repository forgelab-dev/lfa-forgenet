package mail

import (
	"context"
	"encoding/base64"
	"html"
	"io"
	"mime"
	"mime/multipart"
	netmail "net/mail"
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
)

const jeton = "https://forgenet.exemple.org/acces#Zm9yZ2VuZXQtamV0b24tZGUtdGVzdC1xdWktZmFpdC00My1j"

// allMessages couvre chaque e-mail de ForgeNet, avec et sans tunnel.
func allMessages(url string) map[string]Message {
	return map[string]Message{
		"premier accès":            FirstLoginMessage("ana@exemple.org", "ana", url, false, false),
		"premier accès, tunnel":    FirstLoginMessage("ana@exemple.org", "ana", url, true, false),
		"lien refait":              FirstLoginMessage("ana@exemple.org", "ana", url, false, true),
		"mot de passe oublié":      PasswordResetMessage("ana@exemple.org", "ana", url, false),
		"mot de passe oublié, ssh": PasswordResetMessage("ana@exemple.org", "ana", url, true),
		"invitation":               InvitationMessage("bio@exemple.org", "bio", "ana", url, false),
		"invitation, tunnel":       InvitationMessage("bio@exemple.org", "bio", "ana", url, true),
		"adresse changée":          EmailChangedMessage("a@exemple.org", "bio", "b@autre.org", "ana"),
		"essai":                    TestMessage("ana@exemple.org", "ana", "smtp.exemple.org"),
	}
}

var (
	hrefRe = regexp.MustCompile(`href="([^"]*)"`)
	srcRe  = regexp.MustCompile(`src="([^"]*)"`)
	urlRe  = regexp.MustCompile(`https?://[^\s"<>]+`)
	tagRe  = regexp.MustCompile(`(?s)<style.*?</style>|<[^>]+>`)
)

func flat(s string) string { return strings.Join(strings.Fields(s), " ") }

// visibleText est ce qu'affiche le HTML : sans balises, entités décodées.
func visibleText(h string) string { return flat(html.UnescapeString(tagRe.ReplaceAllString(h, " "))) }

func TestLetters_SameWordsInBothVersions(t *testing.T) {
	for name, msg := range allMessages(jeton) {
		t.Run(name, func(t *testing.T) {
			require.NotEmpty(t, msg.HTML)
			shown := visibleText(msg.HTML)
			// Chaque paragraphe du texte brut se retrouve, mot pour mot, dans le HTML.
			for _, p := range strings.Split(msg.Body, "\n\n") {
				if p = flat(p); p != "" {
					assert.Contains(t, shown, p)
				}
			}
		})
	}
}

func TestLetters_TheLinkIsTheOnlyAddress(t *testing.T) {
	// ADR-0014 et ADR-0016 : le jeton reste dans le fragment ; le HTML ne
	// porte aucune autre adresse (pas de pixel, pas de lien réécrit).
	for name, msg := range allMessages(jeton) {
		t.Run(name, func(t *testing.T) {
			hasLink := strings.Contains(msg.Body, jeton)
			for _, m := range hrefRe.FindAllStringSubmatch(msg.HTML, -1) {
				assert.Equal(t, jeton, m[1], "un lien du HTML mène ailleurs")
			}
			for _, u := range urlRe.FindAllString(msg.HTML, -1) {
				assert.Equal(t, jeton, u, "adresse inattendue dans le HTML")
			}
			for _, m := range srcRe.FindAllStringSubmatch(msg.HTML, -1) {
				assert.Equal(t, "cid:"+logoCID, m[1], "seul le logo joint est appelé")
			}
			if hasLink {
				// L'intitulé du bouton porte sa couleur lui-même (Thunderbird
				// impose sa couleur de lien à <a>).
				assert.Regexp(t, `<a href="[^"]+" style="[^"]*"><span style="color:#090c11;">Choisir [^<]+</span></a>`, msg.HTML)
				assert.Contains(t, msg.Body, "  "+jeton+"\n", "le texte garde le lien entier sur sa ligne")
				assert.Len(t, hrefRe.FindAllString(msg.HTML, -1), 2, "le bouton et le lien en clair")
			} else {
				assert.Empty(t, hrefRe.FindAllString(msg.HTML, -1))
			}
		})
	}
}

func TestLetters_AccountDataIsEscaped(t *testing.T) {
	msg := InvitationMessage("bio@exemple.org", "bio", `<img src=x onerror="alert(1)">`, jeton, false)
	assert.NotContains(t, msg.HTML, "<img src=x")
	assert.Contains(t, msg.HTML, "&lt;img src=x")

	// Un lien qui ne serait pas http(s) ne devient jamais un lien actif.
	bad := PasswordResetMessage("ana@exemple.org", "ana", "javascript:alert(1)", false)
	assert.NotContains(t, bad.HTML, `href="javascript:`)
}

func TestLetters_TextLinesFit(t *testing.T) {
	for name, msg := range allMessages(jeton) {
		for _, line := range strings.Split(msg.Body, "\n") {
			if strings.HasPrefix(line, "  ") {
				continue // lien ou commande, jamais coupés
			}
			assert.LessOrEqual(t, utf8.RuneCountInString(line), textWidth, "%s : %q", name, line)
		}
	}
}

func TestWrap(t *testing.T) {
	assert.Equal(t, "un deux\ntrois", wrap("un deux trois", 8))
	// La ponctuation haute reste avec le mot qui la précède.
	assert.Equal(t, "valable une\nfois :", wrap("valable une fois :", 13))
	assert.Equal(t, "mot-très-long-sans-espace", wrap("mot-très-long-sans-espace", 5))
	assert.Equal(t, "", wrap("", 10))
}

func TestFirstLoginText_EndsWithTheCommand(t *testing.T) {
	text := FirstLoginText("ana", jeton, false, false)
	assert.True(t, strings.HasSuffix(text, "  sudo forgenet user login-link --username ana\n"), text)
	assert.Contains(t, text, "  "+jeton+"\n")
}

func TestTestMessage(t *testing.T) {
	msg := TestMessage("ana@exemple.org", "ana", "smtp.exemple.org")
	assert.Equal(t, "ana@exemple.org", msg.To)
	assert.Equal(t, "Essai d'envoi de ForgeNet", msg.Subject)
	assert.Contains(t, flat(msg.Body), "par le relais smtp.exemple.org. Vous n'avez rien à faire.")
	assert.NotContains(t, msg.HTML, "href=")
}

// TestSend_TwoVersionsWithTheLogo vérifie la structure MIME remise au relais :
// alternative(texte, related(html, logo)).
func TestSend_TwoVersionsWithTheLogo(t *testing.T) {
	r, tlsCfg := newRelay(t, false, true)
	msg := PasswordResetMessage("ana@exemple.org", "ana", jeton, false)
	require.NoError(t, sender(r, tlsCfg, config.SMTPStartTLS).Send(context.Background(), msg))
	r.mu.Lock()
	data := r.data
	r.mu.Unlock()

	for _, line := range strings.Split(data, "\r\n") {
		assert.LessOrEqual(t, len(line), 998, "ligne trop longue pour SMTP")
	}

	parsed, err := netmail.ReadMessage(strings.NewReader(data))
	require.NoError(t, err)
	assert.Equal(t, "1.0", parsed.Header.Get("MIME-Version"))
	mediaType, params, err := mime.ParseMediaType(parsed.Header.Get("Content-Type"))
	require.NoError(t, err)
	require.Equal(t, "multipart/alternative", mediaType)

	alt := multipart.NewReader(parsed.Body, params["boundary"])
	text, err := alt.NextPart()
	require.NoError(t, err)
	assert.Equal(t, "text/plain; charset=UTF-8", text.Header.Get("Content-Type"))
	plain, err := io.ReadAll(text)
	require.NoError(t, err)
	assert.Equal(t, strings.ReplaceAll(msg.Body, "\n", "\r\n"), string(plain))

	related, err := alt.NextPart()
	require.NoError(t, err)
	relType, relParams, err := mime.ParseMediaType(related.Header.Get("Content-Type"))
	require.NoError(t, err)
	require.Equal(t, "multipart/related", relType)
	assert.Equal(t, "text/html", relParams["type"])

	rel := multipart.NewReader(related, relParams["boundary"])
	htmlPart, err := rel.NextPart()
	require.NoError(t, err)
	assert.Equal(t, "text/html; charset=UTF-8", htmlPart.Header.Get("Content-Type"))
	body, err := io.ReadAll(htmlPart)
	require.NoError(t, err)
	assert.Equal(t, strings.ReplaceAll(msg.HTML, "\n", "\r\n"), string(body))
	assert.Contains(t, string(body), `href="`+jeton+`"`)
	assert.Contains(t, string(body), `src="cid:`+logoCID+`"`)

	logo, err := rel.NextPart()
	require.NoError(t, err)
	assert.Equal(t, "image/png", logo.Header.Get("Content-Type"))
	assert.Equal(t, "<"+logoCID+">", logo.Header.Get("Content-Id"))
	assert.Contains(t, logo.Header.Get("Content-Disposition"), "inline")
	raw, err := io.ReadAll(logo)
	require.NoError(t, err)
	png, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(string(raw), "\r\n", ""))
	require.NoError(t, err)
	assert.Equal(t, logoPNG, png)

	_, err = rel.NextPart()
	assert.ErrorIs(t, err, io.EOF, "rien d'autre dans la partie related")
	_, err = alt.NextPart()
	assert.ErrorIs(t, err, io.EOF, "deux versions, pas une de plus")
}

func TestSend_TextOnlyWithoutHTML(t *testing.T) {
	// Un Message sans HTML reste en texte brut seul, sans logo.
	r, tlsCfg := newRelay(t, false, true)
	require.NoError(t, sender(r, tlsCfg, config.SMTPStartTLS).Send(context.Background(), essai))
	r.mu.Lock()
	defer r.mu.Unlock()
	assert.NotContains(t, r.data, "multipart")
	assert.NotContains(t, r.data, "image/png")
}
