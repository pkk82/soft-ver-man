package domain

type Asset struct {
	Version         string
	Name            string
	Url             string
	Type            Type
	ExtraProperties map[string]string
}
