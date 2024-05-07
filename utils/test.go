package utils

import (
	"fmt"
	"strings"
)

func ShouldContain(s string) func(string) error {
	return func(t string) error {
		if !strings.Contains(t, s) {
			return fmt.Errorf("expected \n\"\"\"\n%s\"\"\" to contain \"%s\"", t, s)
		}
		return nil
	}
}

func ShouldntContain(s string) func(string) error {
	return func(t string) error {
		if strings.Contains(t, s) {
			return fmt.Errorf("expected \n\"\"\"\n%s\"\"\" not to contain \"%s\"", t, s)
		}
		return nil
	}
}
