/*
Copyright © 2023 Piotr Kozak <piotrkrzysztofkozak@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package domain

import (
	"reflect"
	"testing"
)

func TestFetchedPackage(t *testing.T) {

	fetchedPackages := []FetchedPackage{
		{
			Version: Version{
				Value: "v0.8.6",
			},
			Type:     TAR_GZ,
			FilePath: "/tmp/node/node-v0.8.6-darwin-x64.tar.gz",
		},
		{
			Version: Version{
				Value: "v0.8.6",
			},
			Type:     TAR_GZ,
			FilePath: "/tmp/node/node-v0.8.6-darwin-x64.tgz",
		},
		{
			Version: Version{
				Value: "v20.4.0",
			},
			Type:     ZIP,
			FilePath: "/tmp/node/node-v20.4.0-win-x64.zip",
		},
	}

	expected := []string{"node-v0.8.6-darwin-x64", "node-v0.8.6-darwin-x64", "node-v20.4.0-win-x64"}

	for i, fetchPackage := range fetchedPackages {
		actual := fetchPackage.getDirName()
		if actual != expected[i] {
			t.Errorf("Expected: %v, got: %v", expected[i], actual)
		}
	}
}

func TestRoundVersion(t *testing.T) {
	installedPackage := InstalledPackage{
		InstalledOn: 1689017267000,
		Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
		Version:     toVersion("v20.1.3", t),
		Main:        true,
	}
	type args struct {
		installedPackage InstalledPackage
		granularity      VersionGranularity
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    Version
	}{
		{
			name: "major",
			args: args{
				installedPackage: installedPackage,
				granularity:      VersionGranularityMajor,
			},
			wantErr: false,
			want:    toVersion("20", t),
		},
		{
			name: "minor",
			args: args{
				installedPackage: installedPackage,
				granularity:      VersionGranularityMinor,
			},
			wantErr: false,
			want:    toVersion("20.1", t),
		},
		{
			name: "minor",
			args: args{
				installedPackage: installedPackage,
				granularity:      "unknown",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.installedPackage.RoundVersion(tt.args.granularity)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstalledPackage.RoundVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstalledPackage.RoundVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}
