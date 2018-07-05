package ui

import (
	"testing"
)

func TestDimensions(t *testing.T) {
	uis := NewUiSpec(100, 100)

	rvsd := uis.RequestViewSpec.Dimensions()
	if rvsd[0] != 0 || rvsd[1] != 0 || rvsd[2] != 48 || rvsd[3] != 96 {
		t.Error("Invalid dimensions for RequestViewSpec")
	}

	rsvsd := uis.ResponseViewSpec.Dimensions()
	if rsvsd[0] != 52 || rsvsd[1] != 0 || rsvsd[2] != 99 || rsvsd[3] != 96 {
		t.Error("Invalid dimensions for ResponseViewSpec")
	}

	cvsd := uis.CmdBarViewSpec.Dimensions()
	if cvsd[0] != 0 || cvsd[1] != 97 || cvsd[2] != 48 || cvsd[3] != 99 {
		t.Error("Invalid dimensions for CmdBarViewSpec")
	}

	svsd := uis.StatusCodeViewSpec.Dimensions()
	if svsd[0] != 52 || svsd[1] != 97 || svsd[2] != 61 || svsd[3] != 99 {
		t.Error("Invalid dimensions for StatusCodeViewSpec")
	}
}
