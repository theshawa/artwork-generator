package main

import (
	"errors"
	"image/color"
	"io/ioutil"
	"os"
)

func FormatPath(path string) string {
	if len(path) > 0 && string(path[len(path)-1]) != "/" {
		return path + "/"
	}
	return path
}
func ValidathPath(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	// attempting to create it
	var d []byte
	if err := ioutil.WriteFile(path, d, 0644); err == nil {
		// deleting created dire
		os.Remove(path)
		return true
	}

	return false
}

func ParseHexColor(s string) (c color.RGBA, err error) {

	errInvalidFormat := errors.New("invalid background color provided!!! please add like `#f2d3a2`")
	c.A = 0xff

	if s[0] != '#' {
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errInvalidFormat
	}
	return
}
