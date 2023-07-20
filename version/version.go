package version

type Version struct {
	Value string
	major int
	minor int
	patch int
}

func NewVersion(version string) (Version, error) {
	v, err := parseVersion(version)
	if err != nil {
		return Version{}, err
	}
	return v, nil
}

func (receiver Version) Major() int {
	return receiver.major
}

func CompareDesc(v1, v2 Version) bool {
	return v1.major > v2.major ||
		(v1.major == v2.major && v1.minor > v2.minor) ||
		(v1.major == v2.major && v1.minor == v2.minor && v1.patch > v2.patch)
}

func CompareAsc(v1, v2 Version) bool {
	return CompareDesc(v2, v1)
}
