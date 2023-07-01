package node

import "testing"

func TestGetFingerprints(t *testing.T) {
	actual := getFingerprints()
	expected := []string{
		"4ED778F539E3634C779C87C6D7062848A1AB005C",
		"141F07595B7B3FFE74309A937405533BE57C7D57",
		"74F12602B6F1C4E913FAA37AD3A89613643B6201",
		"DD792F5973C6DE52C432CBDAC77ABFA00DDBF2B7",
		"8FCCA13FEF1D0C2E91008E09770F7A9A5AE15600",
		"C4F0DFFF4E8C1A8236409D08E73BC641CC11F4C8",
		"890C08DB8579162FEE0DF9DB8BEAB4DFCF555EF4",
		"C82FA3AE1CBEDC6BE46B9360C43CEC45C17AB93C",
		"108F52B48DB57BB0CC439B2997B01419BD92F80A",
	}
	if !areSlicesEqual(actual, expected) {
		t.Errorf("Expected: %v, but got: %v", expected, actual)
	}

}

func areSlicesEqual(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := 0; i < len(slice1); i++ {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}
