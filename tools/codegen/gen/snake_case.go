package gen

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var camelAcronyms = map[string]bool{
	"AAAA":  true,
	"ALIAS": true,
	"API":   true,
	"CAA":   true,
	"CNAME": true,
	"DNS":   true,
	"IP":    true,
	"IPS":   true,
	"ISO":   true,
	"JSON":  true,
	"MX":    true,
	"NS":    true,
	"SRV":   true,
	"SSH":   true,
	"SSHFP": true,
	"TXT":   true,
	"XML":   true,
	"YAML":  true,
}

var titler = cases.Title(language.English)

func snakeToPascal(snake string) string {
	parts := strings.Split(snake, "_")
	r := []string{}

	for _, part := range parts {
		if part == "" {
			continue
		}

		upper := strings.ToUpper(part)
		if camelAcronyms[upper] {
			r = append(r, upper)
		} else {
			r = append(r, titler.String(part))
		}
	}

	return strings.Join(r, "")
}
