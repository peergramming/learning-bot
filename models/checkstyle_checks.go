package models

var Checks = map[string]CheckDesc{
	"ArrayTrailingComma": CheckDesc{Category: "coding", Description: "Checks that array initialisation contains a trailing comma."},
}

type CheckDesc struct {
	Category    string
	Description string
}
