package marketplace

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"strings"
)

// maxIcon borne la taille d'une icône : une carte du catalogue n'a pas besoin
// de plus.
const maxIcon = 16 << 10

// CheckIcon vérifie l'icône d'un template : un SVG carré, sans script, sans
// gestionnaire d'événement ni lien hors du fichier. Le dashboard l'affiche par
// une balise <img>, qui n'exécute rien ; la règle garde le dépôt propre et
// l'icône sûre partout où elle sera reprise.
func CheckIcon(data []byte) []error {
	if len(data) == 0 {
		return []error{problem(IconFile, "fichier vide")}
	}
	if len(data) > maxIcon {
		return []error{problem(IconFile, "%d octets : 16 Kio au plus", len(data))}
	}

	var errs []error
	add := func(format string, args ...any) { errs = append(errs, problem(IconFile, format, args...)) }
	dec := xml.NewDecoder(bytes.NewReader(data))
	root := true
	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			add("SVG illisible : %v", err)
			break
		}
		switch t := tok.(type) {
		case xml.Directive:
			add("déclaration <!%s…> interdite", firstWord(string(t)))
		case xml.StartElement:
			if root {
				root = false
				if t.Name.Local != "svg" {
					add("l'élément racine doit être <svg>, pas <%s>", t.Name.Local)
				} else if !squareViewBox(t.Attr) {
					add("viewBox carré attendu (viewBox=\"0 0 64 64\")")
				}
			}
			switch strings.ToLower(t.Name.Local) {
			case "script", "foreignobject":
				add("élément <%s> interdit", t.Name.Local)
			}
			for _, a := range t.Attr {
				name := strings.ToLower(a.Name.Local)
				switch {
				case strings.HasPrefix(name, "on"):
					add("attribut %s interdit sur <%s>", a.Name.Local, t.Name.Local)
				case name == "href" && !strings.HasPrefix(a.Value, "#"):
					add("lien %q interdit : l'icône ne renvoie qu'à elle-même (href=\"#…\")", a.Value)
				}
			}
		}
	}
	if root && len(errs) == 0 {
		add("aucun élément <svg>")
	}
	return errs
}

func squareViewBox(attrs []xml.Attr) bool {
	for _, a := range attrs {
		if a.Name.Local != "viewBox" {
			continue
		}
		f := strings.Fields(strings.ReplaceAll(a.Value, ",", " "))
		if len(f) != 4 {
			return false
		}
		w, errW := strconv.ParseFloat(f[2], 64)
		h, errH := strconv.ParseFloat(f[3], 64)
		return errW == nil && errH == nil && w > 0 && w == h
	}
	return false
}

func firstWord(s string) string {
	if f := strings.Fields(s); len(f) > 0 {
		return f[0]
	}
	return ""
}
