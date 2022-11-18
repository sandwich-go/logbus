package logbus

import (
	"errors"
	"testing"
)

func TestSetGlobalGLogger(t *testing.T) {
	SetGlobalGLogger(nil, "siid", false)
	Debug("debug array", Strings("strings", []string{"a", "b", "c"}), Uint64s("uints64", []uint64{1, 2, 3}), ErrorField(errors.New("an error")))
}

func TestPrintAsError(t *testing.T) {
	SetGlobalGLogger(nil, "printAsErr", true)
	Debug("debug array", Strings("strings", []string{"a", "b", "c"}), Uint64s("uints64", []uint64{1, 2, 3}), ErrorField(errors.New("an error")))
}
