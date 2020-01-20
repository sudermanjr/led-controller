package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"k8s.io/klog"
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

const chunkSize = 64000

func deepCompare(file1, file2 string) bool {
	// Check file size ...

	f1, err := os.Open(file1)
	if err != nil {
		klog.Fatal(err)
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		klog.Fatal(err)
	}
	defer f2.Close()

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}
