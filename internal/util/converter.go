package util

import (
	"strconv"
	"strings"
)

// PercentToDMX converts 0–100% → 0–255
func PercentToDMX(percent int) byte {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	return byte(percent * 255 / 100)
}

// HexToRGB converts "#RRGGBB" → RGB bytes
func HexToRGB(hex string) (byte, byte, byte, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 0, 0, 0, strconv.ErrSyntax
	}

	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return 0, 0, 0, err
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return 0, 0, 0, err
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return 0, 0, 0, err
	}

	return byte(r), byte(g), byte(b), nil
}
