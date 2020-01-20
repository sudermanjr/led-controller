package main

import (
	"testing"
)

// workaround for some weird go 1.13 testing thing with flags
// https://stackoverflow.com/questions/29699982/go-test-flag-flag-provided-but-not-defined
// https://github.com/golang/go/issues/31859
var _ = func() bool {
	testing.Init()
	return true
}()

func init() {
	minBrightness = 30
	maxBrightness = 200
}
