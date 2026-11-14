package metrics

import "time"

// Kind dit comment résumer plusieurs mesures en un point.
type Kind int

const (
	// Gauge : une valeur instantanée (pourcentage, octets occupés), résumée
	// par sa moyenne.
	Gauge Kind = iota
	// Counter : un compteur cumulé (octets reçus), résumé par sa dernière
	// valeur ; le client en tire un débit par différence.
	Counter
)

// Point est une valeur datée.
type Point struct {
	At    time.Time
	Value float64
}

// MaxPoints est le plafond de points par série d'une réponse d'historique.
const MaxPoints = 1000

// Step choisit le pas de regroupement d'une requête d'historique : le pas
// demandé, sans descendre sous le rythme du relevé (requested nul) ni produire
// plus de maxPoints points sur [from, to]. Le pas est arrondi à la seconde.
func Step(from, to time.Time, requested, sampleEvery time.Duration, maxPoints int) time.Duration {
	if maxPoints <= 0 || maxPoints > MaxPoints {
		maxPoints = MaxPoints
	}
	step := requested
	if step <= 0 {
		step = sampleEvery
	}
	if span := to.Sub(from); span > 0 {
		least := (span + time.Duration(maxPoints) - 1) / time.Duration(maxPoints)
		step = max(step, least)
	}
	step = (step + time.Second - 1).Truncate(time.Second)
	return max(step, time.Second)
}

// Downsample regroupe des points triés par date en tranches de step, alignées
// sur l'époque Unix : deux requêtes successives donnent les mêmes tranches, et
// le client peut les raccorder. Chaque tranche donne un point daté de sa
// dernière mesure.
func Downsample(points []Point, step time.Duration, kind Kind) []Point {
	if len(points) == 0 || step <= 0 {
		return nil
	}
	out := make([]Point, 0, len(points))
	var bucket int64
	var sum float64
	var n int
	flush := func(last Point) {
		v := last.Value
		if kind == Gauge {
			v = sum / float64(n)
		}
		out = append(out, Point{At: last.At, Value: v})
	}
	var prev Point
	for i, p := range points {
		b := p.At.UnixNano() / int64(step)
		if i > 0 && b != bucket {
			flush(prev)
			sum, n = 0, 0
		}
		bucket, prev = b, p
		sum += p.Value
		n++
	}
	flush(prev)
	return out
}
