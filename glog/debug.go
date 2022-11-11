//go:build debug
// +build debug

package glog

func needCheckNil() bool         { return true }
func glogInternalError(s string) { panic(s) }
