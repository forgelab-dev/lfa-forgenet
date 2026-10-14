package ids_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/forgelab-dev/lfa-forgenet/internal/ids"
)

func TestNew_PrefixAndUniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for range 1000 {
		id := ids.New("usr")
		assert.Regexp(t, `^usr_[0-9a-f]{16}$`, id)
		assert.False(t, seen[id], "identifiant en double : %s", id)
		seen[id] = true
	}
}
