package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"k8s.io/klog"
)

const chunkSize = 64000

//DeepCompareFiles compares two files
func DeepCompareFiles(file1, file2 string) bool {
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

//IPAddress returns the current IP Address
func IPAddress() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			klog.V(8).Infof("IP Found: %s", ip.String())

			switch ip.String() {
			case "127.0.0.1":
				continue
			case "::1":
				continue
			}
			return ip.String(), nil
		}
	}
	return "", fmt.Errorf("Blank IP found")
}

// ScaleBrightness converts a 0-100 homekit brightness
// to the scale of the controller (min - max)
// math isn't as easy as it used to be for me:
// https://stackoverflow.com/questions/5294955/how-to-scale-down-a-range-of-numbers-with-a-known-min-and-max-value
func ScaleBrightness(value int, minArray int, maxArray int) int {
	min := 0
	max := 100
	a := minArray
	b := maxArray

	new := ((b-a)*(value-min))/(max-min) + a

	return new
}
