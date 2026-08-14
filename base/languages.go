package base

import (
	"errors"
	"strings"
)

type Language struct {
	Portuguese string
	Italian    string
	English    string
	French     string
	Spanish    string
	German     string
}

var Languages = &Language{
	Portuguese: "Portuguese",
	Italian:    "Italian",
	English:    "English",
	French:     "French",
	Spanish:    "Spanish",
	German:     "German",
}

func LanguageFromCountryCode(code string) (string, error) {
	lowerCode := strings.ToLower(code)
	languageCode := strings.Split(lowerCode, "_")[0]

	switch languageCode {
	case "pt":
		return Languages.Portuguese, nil
	case "it":
		return Languages.Italian, nil
	case "en":
		return Languages.English, nil
	case "fr":
		return Languages.French, nil
	case "es":
		return Languages.Spanish, nil
	case "de":
		return Languages.German, nil
	default:
		return "", errors.New("invalid language code: " + code)
	}
}
