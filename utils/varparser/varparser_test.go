package varparser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	vp := New("stash")
	assert.Equal(t, "stash", vp.kind)
}

type rep struct{}

func (r rep) Replace(s string) string {
	return fmt.Sprintf("test-variable-%s", s)
}

func TestParse(t *testing.T) {
	r := rep{}
	vp1 := VarParser{kind: "stash"}
	res1 := vp1.Parse("${stash[test]}", r)
	assert.Equal(t, "test-variable-test", res1)

	vp2 := VarParser{kind: "env"}
	res2 := vp2.Parse("${env[hello]}", r)
	assert.Equal(t, "test-variable-hello", res2)

	vp3 := VarParser{kind: "multi"}
	res3 := vp3.Parse("${multi[hello]} string ${multi[world]}", r)
	assert.Equal(t, "test-variable-hello string test-variable-world", res3)

	vp4 := VarParser{kind: "ignore"}
	res4 := vp4.Parse("${multi[hello]} string ${multi[world]}", r)
	assert.Equal(t, "${multi[hello]} string ${multi[world]}", res4)
}
