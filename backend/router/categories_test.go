package router_test

import (
	"testing"

	"github.com/0xstxrless/punkt/backend/router"
)

func TestIsValidHex(t *testing.T) {
	tests := []struct {
		Name   string
		Expect bool
		Got    bool
	}{
		{
			Name:   "Valid Hex",
			Expect: true,
			Got:    router.IsValidHex("#1A2B3C"),
		},
		{
			Name:   "Invalid Hex Values",
			Expect: false,
			Got:    router.IsValidHex("#INVALID"),
		},
		{
			Name:   "Missing # symbol",
			Expect: false,
			Got:    router.IsValidHex("A1A2B3C"),
		},
		{
			Name:   "Too short hex number",
			Expect: false,
			Got:    router.IsValidHex("#1B3"),
		},
		{
			Name:   "Hex only digits",
			Expect: true,
			Got:    router.IsValidHex("#112233"),
		},
		{
			Name:   "Hex camelCase string",
			Expect: true,
			Got:    router.IsValidHex("#FfAaFf"),
		},
		{
			Name:   "Hex only lowercase chars",
			Expect: true,
			Got:    router.IsValidHex("#aabbaa"),
		},
		{
			Name:   "Hex only uppercase chars",
			Expect: true,
			Got:    router.IsValidHex("#CCAACC"),
		},
		{
			Name:   "Too long hex",
			Expect: false,
			Got:    router.IsValidHex("#AABBCCDD"),
		},
		{
			Name:   "Empty string",
			Expect: false,
			Got:    router.IsValidHex(""),
		},
	}

	for _, test := range tests {
		if test.Expect != test.Got {
			t.Errorf("Failed test isValidHex: %v\nGot: %v\tExpected: %v\n", test.Name, test.Got, test.Expect)
		}
	}
}
