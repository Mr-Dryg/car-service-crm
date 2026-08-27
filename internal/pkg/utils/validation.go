package utils

import "regexp"

var (
	carPlateRegex = regexp.MustCompile(
		`^[АВЕКМНОРСТУХавекмнорстухABEKMHOPCTYXabekmhopctyx]\d{3}[АВЕКМНОРСТУХавекмнорстухABEKMHOPCTYXabekmhopctyx]{2}\d{2,3}$`,
	)
	phoneRegex = regexp.MustCompile(`^\+?\d{10,15}$`)
)

func IsValidCarPlate(plate string) bool {
	return carPlateRegex.MatchString(plate)
}

func IsValidPhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}
