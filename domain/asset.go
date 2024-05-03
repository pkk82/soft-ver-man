package domain

type Asset struct {
	Version           string
	Name              string
	Url               string
	ExternalReference string
	Type              Type
}
