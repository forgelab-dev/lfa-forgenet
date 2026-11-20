// Package mail envoie les e-mails de ForgeNet par un relais SMTP (ADR-0014).
//
// Bibliothèque standard seule (net/smtp, crypto/tls, mime/multipart) : un
// destinataire, le texte brut UTF-8 et sa version HTML avec le logo joint
// (ADR-0016). Le paquet ne journalise rien : les appelants notent le type de
// message et le compte, jamais le contenu.
package mail

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strconv"
	"strings"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
)

// ErrDisabled est retourné quand aucun relais n'est configuré.
var ErrDisabled = errors.New("l'envoi d'e-mails n'est pas configuré sur ce serveur (section smtp)")

// DefaultTimeout borne une remise complète, connexion comprise.
const DefaultTimeout = 30 * time.Second

// Message est un e-mail pour un seul destinataire. Body est le texte brut,
// toujours présent ; HTML, s'il n'est pas vide, en est la version mise en
// forme, qui appelle le logo par cid:.
type Message struct {
	To      string
	Subject string
	Body    string
	HTML    string
}

// Sender remet un message. Les tests en fournissent une doublure.
type Sender interface {
	Send(ctx context.Context, msg Message) error
}

// SMTP remet les messages au relais configuré.
type SMTP struct {
	cfg config.SMTPConfig
	// tlsConfig remplace la configuration TLS, pour les tests seulement.
	tlsConfig *tls.Config
	now       func() time.Time
}

// NewSMTP prépare l'envoi ; la configuration a déjà été validée au chargement.
func NewSMTP(cfg config.SMTPConfig) *SMTP {
	return &SMTP{cfg: cfg, now: time.Now}
}

// Enabled dit si un relais est configuré.
func (s *SMTP) Enabled() bool { return s.cfg.Enabled() }

// Send remet msg au relais. L'erreur reprend la réponse du relais, utile à
// l'administrateur, sans jamais ses identifiants.
func (s *SMTP) Send(ctx context.Context, msg Message) error {
	if !s.Enabled() {
		return ErrDisabled
	}
	to, err := mail.ParseAddress(msg.To)
	if err != nil || strings.ContainsAny(msg.To, "\r\n") {
		return fmt.Errorf("adresse du destinataire invalide : %q", msg.To)
	}
	if strings.ContainsAny(msg.Subject, "\r\n") {
		return errors.New("le sujet ne doit pas contenir de saut de ligne")
	}
	from, err := mail.ParseAddress(s.cfg.From)
	if err != nil {
		return fmt.Errorf("expéditeur invalide : %w", err)
	}
	data, err := s.build(from, to, msg)
	if err != nil {
		return err
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, DefaultTimeout)
		defer cancel()
	}
	c, err := s.dial(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = c.Close() }()
	return s.deliver(c, from.Address, to.Address, data)
}

func (s *SMTP) tlsFor() *tls.Config {
	if s.tlsConfig != nil {
		return s.tlsConfig.Clone()
	}
	return &tls.Config{ServerName: s.cfg.Host, MinVersion: tls.VersionTLS12}
}

// dial ouvre la connexion, en TLS dès le départ pour le mode implicit. Le
// délai du contexte s'applique à toute la conversation.
func (s *SMTP) dial(ctx context.Context) (*smtp.Client, error) {
	addr := net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port))
	var (
		conn net.Conn
		err  error
	)
	if s.cfg.TLS == config.SMTPImplicit {
		d := &tls.Dialer{Config: s.tlsFor()}
		conn, err = d.DialContext(ctx, "tcp", addr)
	} else {
		var d net.Dialer
		conn, err = d.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return nil, fmt.Errorf("connexion au relais %s : %w", addr, err)
	}
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	}
	c, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("accueil du relais %s : %w", addr, err)
	}
	return c, nil
}

func (s *SMTP) deliver(c *smtp.Client, from, to string, data []byte) error {
	if err := c.Hello(helloName(from)); err != nil {
		return fmt.Errorf("présentation au relais : %w", err)
	}
	if s.cfg.TLS == config.SMTPStartTLS {
		// Exigé, jamais facultatif : sans lui, le mot de passe du relais et le
		// lien du message passeraient en clair.
		if ok, _ := c.Extension("STARTTLS"); !ok {
			return errors.New("le relais ne propose pas STARTTLS : choisissez le port 465 (smtp.tls=implicit) ou un autre relais")
		}
		if err := c.StartTLS(s.tlsFor()); err != nil {
			return fmt.Errorf("STARTTLS : %w", err)
		}
	}
	if s.cfg.Username != "" {
		// PlainAuth refuse lui-même d'envoyer le mot de passe hors TLS, sauf
		// vers localhost.
		if err := c.Auth(smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)); err != nil {
			return fmt.Errorf("authentification refusée par le relais : %w", err)
		}
	}
	if err := c.Mail(from); err != nil {
		return fmt.Errorf("expéditeur refusé par le relais : %w", err)
	}
	if err := c.Rcpt(to); err != nil {
		return fmt.Errorf("destinataire refusé par le relais : %w", err)
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("envoi du message : %w", err)
	}
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("envoi du message : %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("message refusé par le relais : %w", err)
	}
	return c.Quit()
}

// build écrit le message : en-têtes encodés, fins de ligne CRLF comme le
// veut SMTP. Sans HTML, le texte seul ; avec, la structure que lisent toutes
// les messageries :
//
//	multipart/alternative
//	├── text/plain
//	└── multipart/related
//	    ├── text/html
//	    └── image/png (logo, Content-ID)
func (s *SMTP) build(from, to *mail.Address, msg Message) ([]byte, error) {
	var b bytes.Buffer
	header := func(k, v string) { b.WriteString(k + ": " + v + "\r\n") }
	header("From", from.String())
	header("To", to.String())
	header("Subject", mime.QEncoding.Encode("utf-8", msg.Subject))
	header("Date", s.now().Format(time.RFC1123Z))
	header("Message-ID", messageID(from.Address))
	header("MIME-Version", "1.0")

	if msg.HTML == "" {
		header("Content-Type", "text/plain; charset=UTF-8")
		header("Content-Transfer-Encoding", "quoted-printable")
		b.WriteString("\r\n")
		if err := writeQP(&b, msg.Body); err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	}

	alt := multipart.NewWriter(&b)
	header("Content-Type", "multipart/alternative; boundary="+alt.Boundary())
	b.WriteString("\r\n")
	if err := writeAlternative(alt, msg); err != nil {
		return nil, fmt.Errorf("écriture du message : %w", err)
	}
	return b.Bytes(), nil
}

func writeAlternative(alt *multipart.Writer, msg Message) error {
	text, err := alt.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"text/plain; charset=UTF-8"},
		"Content-Transfer-Encoding": {"quoted-printable"},
	})
	if err != nil {
		return err
	}
	if err := writeQP(text, msg.Body); err != nil {
		return err
	}

	// La partie related s'écrit d'abord à part : son en-tête doit porter la
	// frontière que choisit son propre writer.
	var inner bytes.Buffer
	rel := multipart.NewWriter(&inner)
	html, err := rel.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"text/html; charset=UTF-8"},
		"Content-Transfer-Encoding": {"quoted-printable"},
	})
	if err != nil {
		return err
	}
	if err := writeQP(html, msg.HTML); err != nil {
		return err
	}
	logo, err := rel.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"image/png"},
		"Content-Transfer-Encoding": {"base64"},
		"Content-Id":                {"<" + logoCID + ">"},
		"Content-Disposition":       {`inline; filename="forgenet.png"`},
	})
	if err != nil {
		return err
	}
	if err := writeBase64(logo, logoPNG); err != nil {
		return err
	}
	if err := rel.Close(); err != nil {
		return err
	}

	related, err := alt.CreatePart(textproto.MIMEHeader{
		"Content-Type": {`multipart/related; type="text/html"; boundary=` + rel.Boundary()},
	})
	if err != nil {
		return err
	}
	if _, err := related.Write(inner.Bytes()); err != nil {
		return err
	}
	return alt.Close()
}

// writeQP écrit text en quoted-printable, fins de ligne en CRLF.
func writeQP(w io.Writer, text string) error {
	qp := quotedprintable.NewWriter(w)
	body := strings.ReplaceAll(strings.ReplaceAll(text, "\r\n", "\n"), "\n", "\r\n")
	if _, err := qp.Write([]byte(body)); err != nil {
		return fmt.Errorf("écriture du message : %w", err)
	}
	if err := qp.Close(); err != nil {
		return fmt.Errorf("écriture du message : %w", err)
	}
	return nil
}

// writeBase64 écrit data en base64, en lignes de 76 caractères (RFC 2045).
func writeBase64(w io.Writer, data []byte) error {
	enc := base64.StdEncoding.EncodeToString(data)
	for len(enc) > 76 {
		if _, err := fmt.Fprint(w, enc[:76]+"\r\n"); err != nil {
			return err
		}
		enc = enc[76:]
	}
	_, err := fmt.Fprint(w, enc+"\r\n")
	return err
}

// helloName se présente sous le domaine de l'expéditeur : « localhost »,
// la valeur par défaut, fait mauvaise impression auprès de certains relais.
func helloName(from string) string {
	if i := strings.LastIndexByte(from, '@'); i >= 0 && i < len(from)-1 {
		return from[i+1:]
	}
	return "localhost"
}

func messageID(from string) string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return "<" + hex.EncodeToString(b) + "@" + helloName(from) + ">"
}
