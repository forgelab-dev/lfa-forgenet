package update

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Source lit les releases de ForgeNet sur l'API GitHub. C'est le seul appel
// sortant de la mise à jour : il ne transmet rien sur l'installation, à part un
// User-Agent générique (ADR-0011).
type Source struct {
	HTTP      *http.Client
	API       string // https://api.github.com par défaut
	Repo      string // forgelab-dev/lfa-forgenet par défaut
	Raw       string // https://raw.githubusercontent.com par défaut
	UserAgent string
}

// FetchFile télécharge un fichier d'installation tel qu'il est dans le tag de
// la release. Son empreinte est vérifiée par le moteur, contre la fiche.
func (s Source) FetchFile(ctx context.Context, version, name string) ([]byte, error) {
	raw := s.Raw
	if raw == "" {
		raw = "https://raw.githubusercontent.com"
	}
	return s.get(ctx, fmt.Sprintf("%s/%s/v%s/deploy/%s", strings.TrimRight(raw, "/"), s.repo(), version, name), maxFileBytes)
}

const (
	defaultAPI  = "https://api.github.com"
	defaultRepo = "forgelab-dev/lfa-forgenet"
	// Bornes de lecture : la liste des releases et une fiche restent petites.
	maxListBytes = 2 << 20
	maxFileBytes = 64 << 10
)

// ErrNoRelease : aucune release publiée n'a de fiche (ou aucune ne convient).
var ErrNoRelease = errors.New("aucune release avec sa fiche")

type ghRelease struct {
	TagName string `json:"tag_name"`
	Draft   bool   `json:"draft"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

func (s Source) client() *http.Client {
	if s.HTTP != nil {
		return s.HTTP
	}
	return &http.Client{Timeout: 30 * time.Second}
}

func (s Source) api() string {
	if s.API != "" {
		return strings.TrimRight(s.API, "/")
	}
	return defaultAPI
}

func (s Source) repo() string {
	if s.Repo != "" {
		return s.Repo
	}
	return defaultRepo
}

// Latest renvoie la fiche de la plus haute version publiée. Les pré-versions
// (1.2.0-rc.1) ne sont retenues qu'avec withPre.
func (s Source) Latest(ctx context.Context, withPre bool) (Release, error) {
	return s.pick(ctx, func(v Version) bool { return withPre || v.Pre == "" })
}

// ByVersion renvoie la fiche d'une version précise.
func (s Source) ByVersion(ctx context.Context, want Version) (Release, error) {
	return s.pick(ctx, func(v Version) bool { return v.Compare(want) == 0 })
}

func (s Source) pick(ctx context.Context, keep func(Version) bool) (Release, error) {
	var list []ghRelease
	if err := s.getJSON(ctx, fmt.Sprintf("%s/repos/%s/releases?per_page=30", s.api(), s.repo()), &list); err != nil {
		return Release{}, err
	}
	var (
		best    *Version
		fileURL string
	)
	for _, r := range list {
		v, err := ParseVersion(r.TagName)
		if r.Draft || err != nil || !keep(v) {
			continue
		}
		for _, a := range r.Assets {
			if a.Name == ReleaseFile && (best == nil || v.Compare(*best) > 0) {
				vv := v
				best, fileURL = &vv, a.URL
			}
		}
	}
	if best == nil {
		return Release{}, ErrNoRelease
	}
	data, err := s.get(ctx, fileURL, maxFileBytes)
	if err != nil {
		return Release{}, err
	}
	rel, err := ParseRelease(data)
	if err != nil {
		return Release{}, err
	}
	// La fiche doit parler de la release qui la porte.
	if v, _ := ParseVersion(rel.Version); v.Compare(*best) != 0 {
		return Release{}, fmt.Errorf("%w : la fiche de v%s annonce %s", ErrInvalidRelease, best, rel.Version)
	}
	return rel, nil
}

func (s Source) getJSON(ctx context.Context, u string, into any) error {
	data, err := s.get(ctx, u, maxListBytes)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, into); err != nil {
		return fmt.Errorf("réponse de GitHub illisible : %w", err)
	}
	return nil
}

// get refuse toute adresse qui n'est pas en HTTPS (sauf pour les tests, qui
// fournissent leur propre API locale) et borne la taille lue.
func (s Source) get(ctx context.Context, raw string, limit int64) ([]byte, error) {
	u, err := url.Parse(raw)
	local := strings.HasPrefix(raw, s.api()) || (s.Raw != "" && strings.HasPrefix(raw, s.Raw))
	if err != nil || (u.Scheme != "https" && !local) {
		return nil, fmt.Errorf("adresse refusée : %q", raw)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		return nil, err
	}
	ua := s.UserAgent
	if ua == "" {
		ua = "forgenet"
	}
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "application/vnd.github+json, application/json")
	resp, err := s.client().Do(req)
	if err != nil {
		return nil, fmt.Errorf("GitHub injoignable : %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub répond %d pour %s", resp.StatusCode, u.Path)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > limit {
		return nil, fmt.Errorf("réponse trop grande pour %s", u.Path)
	}
	return data, nil
}
