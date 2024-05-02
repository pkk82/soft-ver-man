package domain

type EnvVariable struct {
	Name           string
	PrefixVariable *EnvVariable
	SuffixValue    string
}

type EnvVariables struct {
	Variables              []EnvVariable
	MainVariable           *EnvVariable
	ExecutableRelativePath string
}
