//go:build !debug
// +build !debug

package glog

func needCheckNil() bool         { return false }
func glogInternalError(s string) {}
