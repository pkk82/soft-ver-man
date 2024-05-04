package software

import (
	"errors"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/test"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func FileHandler(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")
	fileName := pathSegments[len(pathSegments)-1]
	file, err := os.Open(filepath.Join("testdata", fileName))
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)
	w.Header().Set("Content-Type", "application/x-gzip")
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func Test_fetch(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(FileHandler))
	defer svr.Close()
	type args struct {
		plugin         domain.Plugin
		inputVersion   string
		verifyChecksum bool
	}
	tests := []struct {
		name    string
		args    args
		want    domain.FetchedPackage
		wantErr bool
	}{
		{
			name: "direct fetch",
			args: args{
				plugin: domain.Plugin{
					Name: "direct",
					GetAvailableAssets: func() ([]domain.Asset, error) {
						return []domain.Asset{}, errors.New("unsupported")
					},
					CalculateDownloadUrl: func(version domain.Version, os, arch string) (string, domain.Type) {
						return svr.URL + "/artifacts/artifact.tar.gz", domain.TAR_GZ
					},
					CalculateDownloadedFileName: func(asset domain.Asset) string {
						return "artifact.tar.gz"
					},
				},
				inputVersion: "1.0.0",
			},
			wantErr: false,
			want: domain.FetchedPackage{
				Version:  domain.Ver("1.0.0", t),
				FilePath: "direct/artifact.tar.gz",
				Type:     domain.TAR_GZ,
			},
		},
		{
			name: "assets fetch first",
			args: args{
				plugin: domain.Plugin{
					Name: "asset",
					GetAvailableAssets: func() ([]domain.Asset, error) {
						return []domain.Asset{
							{Version: "1.0.0", Type: domain.TAR_GZ, Url: svr.URL + "/artifacts/artifact.tar.gz"},
						}, nil
					},
					CalculateDownloadUrl: func(version domain.Version, os, arch string) (string, domain.Type) {
						return svr.URL + "/artifacts/artifact.tar.gz", domain.TAR_GZ
					},
					CalculateDownloadedFileName: func(asset domain.Asset) string {
						return "artifact.tar.gz"
					},
				},
				inputVersion: "1.0.0",
			},
			wantErr: false,
			want: domain.FetchedPackage{
				Version:  domain.Ver("1.0.0", t),
				FilePath: "asset/artifact.tar.gz",
				Type:     domain.TAR_GZ,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := test.CreateTestDir(t)
			got, err := Fetch(tt.args.plugin, tt.args.inputVersion, testDir, tt.args.verifyChecksum)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			want := domain.FetchedPackage{Version: tt.want.Version, FilePath: path.Join(testDir, tt.want.FilePath), Type: tt.want.Type}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("fetch() got = %v, want %v", got, want)
			}
		})
	}
}
