package base

import "strings"

type Language struct {
	Portuguese string
	Italian    string
	English    string
	French     string
}

func LanguageFromCountryCode(code string) string {
	switch strings.ToLower(code) {
	case "pt", "pt-pt", "pt-br", "br":
		return Languages.Portuguese
	case "it", "it-it":
		return Languages.Italian
	case "en", "en-us", "en-gb", "us":
		return Languages.English
	case "fr", "fr-fr":
		return Languages.French
	default:
		return ""
	}
}

var Languages = &Language{
	Portuguese: "Portuguese",
	Italian:    "Italian",
	English:    "English",
	French:     "French",
}
