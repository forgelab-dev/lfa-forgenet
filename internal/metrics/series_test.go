package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStep(t *testing.T) {
	from := time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name      string
		span      time.Duration
		requested time.Duration
		maxPoints int
		want      time.Duration
	}{
		{"rythme du relevé par défaut", 15 * time.Minute, 0, 0, 5 * time.Second},
		{"pas demandé respecté", time.Hour, time.Minute, 0, time.Minute},
		{"plafond de points", time.Hour, 0, 60, time.Minute},
		{"plafond du serveur au-delà de 1 000", 10 * time.Hour, 0, 5000, 36 * time.Second},
		{"arrondi à la seconde supérieure", time.Hour, 0, 7, 515 * time.Second},
		{"jamais sous la seconde", time.Second, 100 * time.Millisecond, 0, time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Step(from, from.Add(tt.span), tt.requested, 5*time.Second, tt.maxPoints)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDownsample(t *testing.T) {
	base := time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC)
	points := []Point{
		{base.Add(5 * time.Second), 10},
		{base.Add(10 * time.Second), 20},
		{base.Add(15 * time.Second), 30},
		{base.Add(65 * time.Second), 50},
	}

	gauge := Downsample(points, time.Minute, Gauge)
	require.Len(t, gauge, 2, "deux tranches d'une minute alignées sur l'époque")
	assert.Equal(t, Point{base.Add(15 * time.Second), 20}, gauge[0], "moyenne, datée de la dernière mesure")
	assert.Equal(t, Point{base.Add(65 * time.Second), 50}, gauge[1])

	counter := Downsample(points, time.Minute, Counter)
	assert.Equal(t, 30.0, counter[0].Value, "un compteur garde sa dernière valeur")

	same := Downsample(points, 5*time.Second, Gauge)
	assert.Equal(t, points, same, "un pas égal au relevé ne change rien")

	assert.Nil(t, Downsample(nil, time.Minute, Gauge))
}
