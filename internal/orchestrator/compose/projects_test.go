package compose

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseProjects(t *testing.T) {
	projects, err := parseProjects([]byte(`[{"Name":"forgenet-a","Status":"exited(1), running(2)","ConfigFiles":"-"},{"Name":"autre","Status":"created(1)","ConfigFiles":"/srv/x.yml"}]`))
	assert.NoError(t, err)
	assert.Equal(t, []Project{{Name: "forgenet-a", Containers: 3, Running: 2}, {Name: "autre", Containers: 1, Running: 0}}, projects)

	empty, err := parseProjects([]byte("[]\n"))
	assert.NoError(t, err)
	assert.Empty(t, empty)

	_, err = parseProjects([]byte("pas du json"))
	assert.Error(t, err)
}
