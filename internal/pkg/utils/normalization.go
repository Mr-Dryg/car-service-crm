package utils

import (
	"strings"
)

func NormalizePhone(phone string) string {
	cleaned := strings.NewReplacer(" ", "", "-", "", "(", "", ")", "").Replace(phone)

	if len(cleaned) == 0 {
		return ""
	}

	if len(cleaned) == 11 && cleaned[0] == '8' {
		return "+7" + cleaned[1:]
	}

	if len(cleaned) == 11 && cleaned[0] == '7' {
		return "+" + cleaned
	}

	if len(cleaned) == 10 {
		return "+7" + cleaned
	}

	if cleaned[0] == '+' {
		return cleaned
	}

	return "+" + cleaned
}

func NormalizeCarPlate(plate string) string {
	upper := strings.ToUpper(plate)

	replacer := strings.NewReplacer(
		"A", "А",
		"B", "В",
		"E", "Е",
		"K", "К",
		"M", "М",
		"H", "Н",
		"O", "О",
		"P", "Р",
		"C", "С",
		"T", "Т",
		"Y", "У",
		"X", "Х",
	)

	return replacer.Replace(upper)
}
