package mail

import (
	"bytes"
	_ "embed"
	"html/template"
	"strings"
	"unicode/utf8"
)

// Un e-mail de ForgeNet se décrit une fois, en blocs, et se rend en deux
// versions : le texte brut et le HTML aux couleurs du dashboard (ADR-0016).
// Les deux disent donc forcément la même chose.

// logoPNG est le logo de 96 pixels joint à chaque e-mail HTML et appelé par
// cid: : ni SVG en ligne (ignoré par Gmail et Outlook) ni image distante (elle
// ferait savoir quand et d'où le message est lu).
//
//go:embed assets/forgenet-logo.png
var logoPNG []byte

// logoCID relie l'image jointe à la balise <img> du modèle.
const logoCID = "forgenet-logo@forgenet"

// textWidth est la largeur des lignes de la version texte : 76 caractères,
// en deçà des 78 que recommande la RFC 5322.
const textWidth = 76

type blockKind string

const (
	paragraph blockKind = "paragraph"
	// aside est une remarque secondaire, plus discrète dans le HTML.
	aside blockKind = "aside"
	// action est le lien à ouvrir : un bouton dans le HTML, une ligne à part
	// dans le texte.
	action blockKind = "action"
	// command est une commande à lancer, sur une seule ligne.
	command blockKind = "command"
)

type block struct {
	Kind  blockKind
	Text  string // paragraphe, remarque ou commande
	URL   string // action
	Label string // action : intitulé du bouton
}

func para(text string) block { return block{Kind: paragraph, Text: text} }
func note(text string) block { return block{Kind: aside, Text: text} }
func cmd(text string) block  { return block{Kind: command, Text: text} }
func link(label, url string) block {
	return block{Kind: action, URL: url, Label: label}
}

// letter est un e-mail avant rendu.
type letter struct {
	subject  string
	greeting string
	blocks   []block
}

// message rend les deux versions. Si le HTML ne se rend pas, ce qui
// n'arrive qu'avec un modèle cassé, le texte part seul.
func (l letter) message(to string) Message {
	html, err := l.html()
	if err != nil {
		html = ""
	}
	return Message{To: to, Subject: l.subject, Body: l.text(), HTML: html}
}

const signature = "— ForgeNet"

func (l letter) text() string {
	var b strings.Builder
	if l.greeting != "" {
		b.WriteString(l.greeting + "\n\n")
	}
	b.WriteString(blocksText(l.blocks))
	b.WriteString(signature + "\n")
	return b.String()
}

// blocksText rend les blocs en texte, séparés par une ligne vide. Les liens
// et les commandes ne sont jamais coupés : ils doivent se copier entiers.
func blocksText(blocks []block) string {
	var b strings.Builder
	for _, bl := range blocks {
		switch bl.Kind {
		case action:
			b.WriteString("  " + bl.URL + "\n\n")
		case command:
			b.WriteString("  " + bl.Text + "\n\n")
		default:
			b.WriteString(wrap(bl.Text, textWidth) + "\n\n")
		}
	}
	return b.String()
}

// wrap coupe text en lignes d'au plus width caractères, entre les mots. Une
// ponctuation haute précédée d'une espace (« : », « ? ») reste avec le mot
// qui la précède, comme le veut la typographie française.
func wrap(text string, width int) string {
	var words []string
	for _, w := range strings.Fields(text) {
		if len(words) > 0 && strings.Trim(w, ":;?!»") == "" {
			words[len(words)-1] += " " + w
			continue
		}
		words = append(words, w)
	}
	var b strings.Builder
	lineLen := 0
	for _, w := range words {
		n := utf8.RuneCountInString(w)
		switch {
		case lineLen == 0:
		case lineLen+1+n > width:
			b.WriteByte('\n')
			lineLen = 0
		default:
			b.WriteByte(' ')
			lineLen++
		}
		b.WriteString(w)
		lineLen += n
	}
	return b.String()
}

// Couleurs du dashboard (web/src/index.css), jetons OKLCH convertis en
// hexadécimal : seule écriture que toutes les messageries comprennent.
// Contrastes sur la carte (iron-900) : texte 16:1, texte secondaire 8,6:1,
// lien 11:1 ; bouton or sur texte iron-950 : 9,7:1.
// Les valeurs sont des template.CSS : constantes du paquet, jamais une
// donnée de compte.
type palette struct {
	Page, Card, Border, Ink, Muted, Gold, GoldLight template.CSS
}

var colors = palette{
	Page:      "#090c11", // iron-950
	Card:      "#10141a", // iron-900
	Border:    "#21272e", // iron-800
	Ink:       "#eceff2", // iron-100
	Muted:     "#abb2ba", // iron-300
	Gold:      "#e3ae28", // gold-400
	GoldLight: "#f2cd6f", // gold-300
}

// Les polices du dashboard ne sont pas chargées (aucune ressource distante) :
// elles servent si elles sont installées, sinon la police système prend le
// relais.
const (
	sansStack template.CSS = "'Bricolage Grotesque','Segoe UI',Helvetica,Arial,sans-serif"
	monoStack template.CSS = "'JetBrains Mono',Consolas,Menlo,monospace"
)

// html/template retire les commentaires : les commentaires conditionnels
// d'Outlook pour Windows, qui ignore max-width, passent donc en valeurs sûres.
const (
	msoOpen  = `<!--[if mso]><table role="presentation" width="560" cellpadding="0" cellspacing="0" border="0"><tr><td><![endif]-->`
	msoClose = `<!--[if mso]></td></tr></table><![endif]-->`
)

var letterTemplate = template.Must(template.New("letter").Parse(letterHTML))

type letterView struct {
	Subject, Greeting, Preheader string
	Blocks                       []block
	C                            palette
	Sans, Mono                   template.CSS
	MSOOpen, MSOClose            template.HTML
}

func (l letter) html() (string, error) {
	view := letterView{
		Subject:  l.subject,
		Greeting: l.greeting,
		Blocks:   l.blocks,
		C:        colors,
		Sans:     sansStack,
		Mono:     monoStack,
		MSOOpen:  template.HTML(msoOpen),  //nolint:gosec // constante du paquet, sans donnée de compte
		MSOClose: template.HTML(msoClose), //nolint:gosec // idem
	}
	for _, bl := range l.blocks {
		if bl.Kind == paragraph {
			view.Preheader = bl.Text
			break
		}
	}
	var b bytes.Buffer
	if err := letterTemplate.Execute(&b, view); err != nil {
		return "", err
	}
	return b.String(), nil
}

// letterHTML est le seul modèle HTML : tableaux et styles en ligne, les seuls
// que lisent toutes les messageries. Le bouton est une cellule colorée, que
// même Outlook pour Windows affiche. Son intitulé porte sa couleur sur un
// <span> : Thunderbird impose sa couleur de lien à <a>, ce qui rendait l'or
// clair illisible sur le bouton or.
const letterHTML = `<!DOCTYPE html>
<html lang="fr">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="color-scheme" content="light dark">
<meta name="supported-color-schemes" content="light dark">
<title>{{.Subject}}</title>
<style>
:root { color-scheme: light dark; supported-color-schemes: light dark; }
@media (max-width: 480px) { .fn-card { padding: 24px 18px !important; } }
</style>
</head>
<body style="margin:0;padding:0;background-color:{{.C.Page}};">
<div style="display:none;max-height:0;overflow:hidden;mso-hide:all;">{{.Preheader}}</div>
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" bgcolor="{{.C.Page}}" style="background-color:{{.C.Page}};">
<tr><td align="center" style="padding:32px 12px;">
{{.MSOOpen}}
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="max-width:560px;">
<tr><td style="padding:0 4px 20px;">
<table role="presentation" cellpadding="0" cellspacing="0" border="0"><tr>
<td style="vertical-align:middle;"><img src="cid:` + logoCID + `" width="32" height="32" alt="" style="display:block;border:0;width:32px;height:32px;"></td>
<td style="vertical-align:middle;padding-left:10px;font-family:{{.Sans}};font-size:19px;line-height:24px;font-weight:700;color:{{.C.Ink}};">Forge<span style="font-weight:400;color:{{.C.Muted}};">Net</span></td>
</tr></table>
</td></tr>
<tr><td class="fn-card" bgcolor="{{.C.Card}}" style="background-color:{{.C.Card}};border:1px solid {{.C.Border}};border-radius:12px;padding:32px 28px;font-family:{{.Sans}};color:{{.C.Ink}};">
{{- if .Greeting}}
<p style="margin:0 0 16px;font-size:20px;line-height:28px;font-weight:700;color:{{.C.Ink}};">{{.Greeting}}</p>
{{- end}}
{{- range .Blocks}}
{{- if eq .Kind "action"}}
<table role="presentation" cellpadding="0" cellspacing="0" border="0" style="margin:8px 0 14px;"><tr>
<td bgcolor="{{$.C.Gold}}" style="background-color:{{$.C.Gold}};border-radius:8px;"><a href="{{.URL}}" style="display:inline-block;padding:13px 24px;font-family:{{$.Sans}};font-size:15px;line-height:20px;font-weight:700;color:{{$.C.Page}};text-decoration:none;border-radius:8px;"><span style="color:{{$.C.Page}};">{{.Label}}</span></a></td>
</tr></table>
<p style="margin:0 0 4px;font-size:13px;line-height:20px;color:{{$.C.Muted}};">ou copiez ce lien :</p>
<p style="margin:0 0 20px;font-family:{{$.Mono}};font-size:13px;line-height:20px;word-break:break-all;"><a href="{{.URL}}" style="color:{{$.C.GoldLight}};text-decoration:underline;">{{.URL}}</a></p>
{{- else if eq .Kind "command"}}
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="margin:0 0 16px;"><tr>
<td bgcolor="{{$.C.Page}}" style="background-color:{{$.C.Page}};border:1px solid {{$.C.Border}};border-radius:8px;padding:10px 14px;font-family:{{$.Mono}};font-size:13px;line-height:20px;color:{{$.C.Ink}};word-break:break-all;">{{.Text}}</td>
</tr></table>
{{- else if eq .Kind "aside"}}
<p style="margin:0 0 16px;font-size:14px;line-height:22px;color:{{$.C.Muted}};">{{.Text}}</p>
{{- else}}
<p style="margin:0 0 16px;font-size:15px;line-height:24px;color:{{$.C.Ink}};">{{.Text}}</p>
{{- end}}
{{- end}}
<p style="margin:24px 0 0;font-size:15px;line-height:24px;color:{{.C.Muted}};">` + signature + `</p>
</td></tr>
</table>
{{.MSOClose}}
</td></tr></table>
</body>
</html>
`
