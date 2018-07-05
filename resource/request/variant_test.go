package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNames(t *testing.T) {
	vs := Variants{
		Variant{
			Name: "test-variant-1",
		},
		Variant{
			Name: "test-variant-2",
		},
	}

	names := vs.Names()

	assert.Equal(t, 2, len(names), "Expecting 2 variant names.")
	assert.Equal(t, "test-variant-1", names[0], "Expecting test-variant-1.")
	assert.Equal(t, "test-variant-2", names[1], "Expecting test-variant-1.")
}
