package mail

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"io"
	"math/big"
	"mime/quotedprintable"
	"net"
	"net/netip"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
)

// fakeRelay est un relais SMTP minimal : EHLO, STARTTLS, AUTH PLAIN, MAIL,
// RCPT, DATA, QUIT. Il garde ce qu'il a reçu pour les assertions.
type fakeRelay struct {
	t        *testing.T
	ln       net.Listener
	cert     tls.Certificate
	implicit bool // TLS dès la connexion (port 465)
	starttls bool // propose STARTTLS
	user     string
	pass     string

	mu       sync.Mutex
	data     string
	rcpt     string
	sawTLS   bool
	authUser string
}

func newRelay(t *testing.T, implicit, starttls bool) (*fakeRelay, *tls.Config) {
	t.Helper()
	cert, pool := selfSigned(t)
	r := &fakeRelay{t: t, cert: cert, implicit: implicit, starttls: starttls, user: "forgenet", pass: "secret-du-relais"}
	var err error
	if implicit {
		r.ln, err = tls.Listen("tcp", "127.0.0.1:0", &tls.Config{Certificates: []tls.Certificate{cert}})
	} else {
		r.ln, err = net.Listen("tcp", "127.0.0.1:0")
	}
	require.NoError(t, err)
	t.Cleanup(func() { _ = r.ln.Close() })
	go r.serve()
	return r, &tls.Config{RootCAs: pool, ServerName: "127.0.0.1", MinVersion: tls.VersionTLS12}
}

func (r *fakeRelay) port() int { return r.ln.Addr().(*net.TCPAddr).Port }

func (r *fakeRelay) serve() {
	for {
		conn, err := r.ln.Accept()
		if err != nil {
			return
		}
		go r.handle(conn)
	}
}

func (r *fakeRelay) handle(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	secure := r.implicit
	rd, wr := bufio.NewReader(conn), conn
	say := func(s string) { _, _ = wr.Write([]byte(s + "\r\n")) }
	say("220 relais d'essai")
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			return
		}
		cmd := strings.TrimRight(line, "\r\n")
		verb := strings.ToUpper(strings.SplitN(cmd, " ", 2)[0])
		switch verb {
		case "EHLO", "HELO":
			lines := []string{"250-relais"}
			if r.starttls && !secure {
				lines = append(lines, "250-STARTTLS")
			}
			lines = append(lines, "250-AUTH PLAIN", "250 OK")
			for _, l := range lines {
				say(l)
			}
		case "STARTTLS":
			say("220 allez-y")
			tc := tls.Server(conn, &tls.Config{Certificates: []tls.Certificate{r.cert}})
			if tc.Handshake() != nil {
				return
			}
			conn, secure = tc, true
			rd, wr = bufio.NewReader(tc), tc
		case "AUTH":
			raw, _ := base64.StdEncoding.DecodeString(strings.TrimSpace(strings.TrimPrefix(cmd, "AUTH PLAIN")))
			parts := strings.Split(string(raw), "\x00")
			if len(parts) == 3 && parts[1] == r.user && parts[2] == r.pass {
				r.mu.Lock()
				r.authUser = parts[1]
				r.mu.Unlock()
				say("235 bienvenue")
			} else {
				say("535 identifiants refusés")
			}
		case "MAIL":
			say("250 OK")
		case "RCPT":
			r.mu.Lock()
			r.rcpt = cmd
			r.mu.Unlock()
			say("250 OK")
		case "DATA":
			say("354 envoyez")
			var b strings.Builder
			for {
				l, err := rd.ReadString('\n')
				if err != nil {
					return
				}
				if l == ".\r\n" {
					break
				}
				b.WriteString(l)
			}
			r.mu.Lock()
			r.data, r.sawTLS = b.String(), secure
			r.mu.Unlock()
			say("250 reçu")
		case "QUIT":
			say("221 au revoir")
			return
		default:
			say("502 inconnu")
		}
	}
}

func selfSigned(t *testing.T) (tls.Certificate, *x509.CertPool) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "relais d'essai"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		IPAddresses:  []net.IP{net.IP(netip.MustParseAddr("127.0.0.1").AsSlice())},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	require.NoError(t, err)
	leaf, err := x509.ParseCertificate(der)
	require.NoError(t, err)
	pool := x509.NewCertPool()
	pool.AddCert(leaf)
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key, Leaf: leaf}, pool
}

func sender(r *fakeRelay, tlsCfg *tls.Config, mode string) *SMTP {
	s := NewSMTP(config.SMTPConfig{
		Host: "127.0.0.1", Port: r.port(), Username: "forgenet", Password: "secret-du-relais",
		From: "ForgeNet <forgenet@exemple.org>", TLS: mode,
	})
	s.tlsConfig = tlsCfg
	s.now = func() time.Time { return time.Date(2026, 11, 18, 12, 0, 0, 0, time.UTC) }
	return s
}

var essai = Message{To: "ana@exemple.org", Subject: "Essai d'envoi — ForgeNet", Body: "Bonjour,\nl'envoi marche.\n"}

func TestSend_StartTLS(t *testing.T) {
	r, tlsCfg := newRelay(t, false, true)
	require.NoError(t, sender(r, tlsCfg, config.SMTPStartTLS).Send(context.Background(), essai))

	r.mu.Lock()
	defer r.mu.Unlock()
	assert.True(t, r.sawTLS, "le message part sous TLS")
	assert.Equal(t, "forgenet", r.authUser)
	assert.Equal(t, "RCPT TO:<ana@exemple.org>", r.rcpt)
	assert.Contains(t, r.data, "From: \"ForgeNet\" <forgenet@exemple.org>\r\n")
	assert.Contains(t, r.data, "Subject: =?utf-8?q?Essai_d'envoi_=E2=80=94_ForgeNet?=\r\n")
	assert.Contains(t, r.data, "Date: Wed, 18 Nov 2026 12:00:00 +0000\r\n")
	assert.Contains(t, r.data, "Content-Type: text/plain; charset=UTF-8\r\n")
	assert.Regexp(t, `Message-ID: <[0-9a-f]{32}@exemple\.org>`, r.data)

	_, body, ok := strings.Cut(r.data, "\r\n\r\n")
	require.True(t, ok)
	decoded, err := io.ReadAll(quotedprintable.NewReader(strings.NewReader(body)))
	require.NoError(t, err)
	assert.Equal(t, "Bonjour,\r\nl'envoi marche.\r\n", string(decoded))
}

func TestSend_ImplicitTLS(t *testing.T) {
	r, tlsCfg := newRelay(t, true, false)
	require.NoError(t, sender(r, tlsCfg, config.SMTPImplicit).Send(context.Background(), essai))
	r.mu.Lock()
	defer r.mu.Unlock()
	assert.True(t, r.sawTLS)
}

func TestSend_RefusesARelayWithoutSTARTTLS(t *testing.T) {
	r, tlsCfg := newRelay(t, false, false)
	err := sender(r, tlsCfg, config.SMTPStartTLS).Send(context.Background(), essai)
	require.ErrorContains(t, err, "ne propose pas STARTTLS")
	r.mu.Lock()
	defer r.mu.Unlock()
	assert.Empty(t, r.authUser, "aucun identifiant envoyé en clair")
	assert.Empty(t, r.data)
}

func TestSend_RefusesAnUntrustedCertificate(t *testing.T) {
	r, _ := newRelay(t, false, true)
	err := sender(r, nil, config.SMTPStartTLS).Send(context.Background(), essai)
	require.ErrorContains(t, err, "STARTTLS")
}

func TestSend_ReportsTheRelayRefusalWithoutCredentials(t *testing.T) {
	r, tlsCfg := newRelay(t, false, true)
	r.pass = "autre"
	err := sender(r, tlsCfg, config.SMTPStartTLS).Send(context.Background(), essai)
	require.ErrorContains(t, err, "535")
	assert.Contains(t, err.Error(), "identifiants refusés")
	assert.NotContains(t, err.Error(), "secret-du-relais")
}

func TestSend_RejectsHeaderInjection(t *testing.T) {
	s := NewSMTP(config.SMTPConfig{Host: "127.0.0.1", Port: 1, From: "forgenet@exemple.org", TLS: config.SMTPStartTLS})
	for _, msg := range []Message{
		{To: "ana@exemple.org\r\nBcc: tous@exemple.org", Subject: "x"},
		{To: "ana@exemple.org", Subject: "x\r\nBcc: tous@exemple.org"},
		{To: "pas une adresse", Subject: "x"},
	} {
		assert.Error(t, s.Send(context.Background(), msg))
	}
}

func TestSend_Disabled(t *testing.T) {
	assert.ErrorIs(t, NewSMTP(config.SMTPConfig{}).Send(context.Background(), essai), ErrDisabled)
}

func TestSend_UnreachableRelay(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := ln.Addr().(*net.TCPAddr).Port
	require.NoError(t, ln.Close())
	s := NewSMTP(config.SMTPConfig{Host: "127.0.0.1", Port: port, From: "forgenet@exemple.org", TLS: config.SMTPStartTLS})
	err = s.Send(context.Background(), essai)
	require.ErrorContains(t, err, "connexion au relais 127.0.0.1:"+strconv.Itoa(port))
}
