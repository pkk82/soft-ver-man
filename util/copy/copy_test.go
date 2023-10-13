package copy

import (
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/test"
	"path"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCopy(t *testing.T) {
	type args struct {
		fetchedPackage domain.FetchedPackage
	}
	tests := []struct {
		name         string
		args         args
		wantPath     string
		wantFileName string
		wantContent  string
	}{
		{
			name: "copy file",
			args: args{
				fetchedPackage: domain.FetchedPackage{
					Version:  domain.Version{Value: "v0.5.0"},
					FilePath: filepath.Join("testdata", "soft-ver-man-v0.5.0"),
					Type:     domain.RAW,
				},
			},
			wantPath:     "soft-ver-man-v0.5.0-dir",
			wantFileName: "svm",
			wantContent:  "soft-ver-man-v0.5.0-content",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := test.CreateTestDir(t)
			got, err := Copy(tt.args.fetchedPackage, path.Join(testDir, "soft-ver-man-v0.5.0-dir"), "svm")
			if err != nil {
				t.Errorf("Extract() error = %v", err)
				return
			}
			expected := CopiedPackage{
				Version:    tt.args.fetchedPackage.Version,
				PathToFile: path.Join(testDir, tt.wantPath),
				FileName:   tt.wantFileName,
			}
			if !reflect.DeepEqual(got, expected) {
				t.Errorf("Copy() got = %v, want %v", got, expected)
			}
			test.AssertFileContent(expected.PathToFile, expected.FileName, []string{tt.wantContent}, t)
			test.AssertFileMode(expected.PathToFile, expected.FileName, 0755, t)
		})
	}
}
