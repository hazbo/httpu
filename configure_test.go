package httpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigureFromFile(t *testing.T) {
	err := ConfigureFromFile("./projects/test_project")
	assert.Equal(t, nil, err, "it should be nil")
}
