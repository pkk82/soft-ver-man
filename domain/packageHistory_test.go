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
	"github.com/yudai/gojsondiff"
	"reflect"
	"testing"
)

func TestParseHistory(t *testing.T) {
	type args struct {
		repr string
	}
	tests := []struct {
		name    string
		args    args
		want    PackageHistory
		wantErr bool
	}{
		{
			name:    "empty string",
			args:    args{repr: ""},
			want:    PackageHistory{},
			wantErr: false,
		},
		{
			name:    "empty json",
			args:    args{repr: "{}"},
			want:    PackageHistory{},
			wantErr: false,
		},
		{
			name: "domain",
			args: args{repr: `{"name": "node", "items": [{"version": "v20.1.3", "path": "/home/user/pf/node/node-v12.22.12-linux-x64", "main": true, "installedOn": 1689017267000}]}`},
			want: PackageHistory{
				Name: "node",
				Items: []PackageHistoryItem{
					{
						Version:     "v20.1.3",
						Path:        "/home/user/pf/node/node-v12.22.12-linux-x64",
						Main:        true,
						InstalledOn: 1689017267000,
					},
				}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHistory(tt.args.repr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseHistory() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackageHistory_Serialize(t *testing.T) {

	tests := []struct {
		name string
		ph   PackageHistory
		want string
	}{
		{
			name: "empty",
			ph:   PackageHistory{},
			want: `{"name":"","items":null}`,
		},
		{
			name: "name only",
			ph:   PackageHistory{Name: "node"},
			want: `{"name":"node","items":null}`,
		},
		{
			name: "domain",
			ph: PackageHistory{Name: "node", Items: []PackageHistoryItem{
				{
					Version:     "v20.1.3",
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Main:        true,
					InstalledOn: 1689017267000,
				}, {
					Version:     "v20.1.4",
					Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
					Main:        false,
					InstalledOn: 1689017268000,
				},
			}},
			want: `{"name":"node","items":[
{"version":"v20.1.3","path":"/home/user/pf/node/node-v20.1.3-linux-x64","installedOn":1689017267000,"main":true},
{"version":"v20.1.4","path":"/home/user/pf/node/node-v20.1.4-linux-x64","installedOn":1689017268000,"main":false}]}`,
		},
	}

	differ := gojsondiff.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := tt.ph
			got := ph.Serialize()
			compare, err := differ.Compare([]byte(got), []byte(tt.want))
			if err != nil {
				t.Errorf("differ.Compare() error = %v", err)
			}
			if len(compare.Deltas()) > 0 {
				t.Errorf("Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackageHistory_Add(t *testing.T) {

	tests := []struct {
		name string
		ph   PackageHistory
		want []PackageHistoryItem
	}{
		{
			name: "empty",
			ph:   PackageHistory{Name: "node", Items: []PackageHistoryItem{}},
			want: []PackageHistoryItem{
				{
					Version:     "v20.1.4",
					Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
					Main:        false,
					InstalledOn: 1689017268000,
				},
			},
		},
		{
			name: "with-one-item",
			ph: PackageHistory{Name: "node", Items: []PackageHistoryItem{
				{
					Version:     "v20.1.3",
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Main:        true,
					InstalledOn: 1689017267000,
				},
			}},
			want: []PackageHistoryItem{{
				Version:     "v20.1.3",
				Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
				Main:        true,
				InstalledOn: 1689017267000,
			}, {
				Version:     "v20.1.4",
				Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
				Main:        false,
				InstalledOn: 1689017268000,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := tt.ph
			got := ph.Add(InstalledPackage{
				Version: Version{
					Value: "v20.1.4",
				},
				Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
				InstalledOn: 1689017268000,
			})
			if !reflect.DeepEqual(got.Items, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestPackageHistory_Exists(t *testing.T) {

	tests := []struct {
		name       string
		ph         PackageHistory
		newVersion Version
		want       bool
	}{
		{
			name:       "empty - false",
			ph:         PackageHistory{Name: "node", Items: []PackageHistoryItem{}},
			newVersion: Version{Value: "v20.1.4"},
			want:       false,
		},
		{
			name: "one element - false",
			ph: PackageHistory{Name: "node", Items: []PackageHistoryItem{
				{
					Version:     "v20.1.3",
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Main:        true,
					InstalledOn: 1689017267000,
				},
			}},
			newVersion: Version{Value: "v20.1.3"},
			want:       true,
		}, {
			name: "one element -  true",
			ph: PackageHistory{Name: "node", Items: []PackageHistoryItem{
				{
					Version:     "v20.1.3",
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Main:        true,
					InstalledOn: 1689017267000,
				},
			}},
			newVersion: Version{Value: "v20.1.4"},
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := tt.ph
			got := ph.IsInstalled(tt.newVersion)
			if got != tt.want {
				t.Errorf("IsInstalled() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestPackageHistory_Remove(t *testing.T) {

	tests := []struct {
		name        string
		ph          PackageHistory
		want        PackageHistory
		wantRemoved *PackageHistoryItem
	}{
		{
			name:        "empty",
			ph:          PackageHistory{Name: "node", Items: []PackageHistoryItem{}},
			want:        PackageHistory{Name: "node", Items: []PackageHistoryItem{}},
			wantRemoved: nil,
		},
		{
			name: "with-one-item",
			ph: PackageHistory{Name: "node", Items: []PackageHistoryItem{
				{
					Version:     "v20.1.3",
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Main:        true,
					InstalledOn: 1689017267000,
				},
			}},
			want: PackageHistory{Name: "node", Items: nil},
			wantRemoved: &PackageHistoryItem{
				Version:     "v20.1.3",
				Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
				Main:        true,
				InstalledOn: 1689017267000,
			},
		},
		{
			name: "with-two-items",
			ph: PackageHistory{Name: "node", Items: []PackageHistoryItem{
				{
					Version:     "v20.1.3",
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Main:        true,
					InstalledOn: 1689017267000,
				},
				{
					Version:     "v20.1.4",
					Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
					Main:        false,
					InstalledOn: 1689017267000,
				},
			}},
			want: PackageHistory{Name: "node", Items: []PackageHistoryItem{
				{
					Version:     "v20.1.4",
					Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
					Main:        false,
					InstalledOn: 1689017267000,
				},
			}},
			wantRemoved: &PackageHistoryItem{
				Version:     "v20.1.3",
				Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
				Main:        true,
				InstalledOn: 1689017267000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := tt.ph
			got, removed := ph.RemoveByVersion(
				Version{
					Value: "v20.1.3",
				})

			if len(got.Items) != 0 && !reflect.DeepEqual(*got, tt.want) {
				//t.Errorf("RemoveByVersion().history = %v, want %v", *got, tt.want)
				t.Errorf("RemoveByVersion().history = %v, want %v", *got, tt.want)
			}
			if got.Name != tt.want.Name || (len(got.Items) != len(tt.want.Items)) {
				t.Errorf("RemoveByVersion().history = %v, want %v", *got, tt.want)
			}

			if !reflect.DeepEqual(removed, tt.wantRemoved) {
				t.Errorf("RemoveByVersion().removed = %v, want %v", removed, tt.wantRemoved)
			}
		})
	}

}

func TestToInstalledPackages(t *testing.T) {
	tests := []struct {
		name    string
		arg     PackageHistory
		want    []InstalledPackage
		wantErr bool
	}{
		{
			name: "empty string",
			arg: PackageHistory{Name: "node", Items: []PackageHistoryItem{
				{
					Version:     "v20.1.3",
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Main:        true,
					InstalledOn: 1689017267000,
				},
				{
					Version:     "v20.1.4",
					Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
					Main:        false,
					InstalledOn: 1689017268000,
				},
			}},
			want: []InstalledPackage{
				{
					InstalledOn: 1689017267000,
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Version:     toVersion("v20.1.3", t),
					Main:        true,
				},
				{
					InstalledOn: 1689017268000,
					Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
					Version:     toVersion("v20.1.4", t),
					Main:        false,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arg.ToInstalledPackages()
			if (err != nil) != tt.wantErr {
				t.Errorf("PackageHistory.ToInstalledPackages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ackageHistory.ToInstalledPackages() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func toVersion(version string, t *testing.T) Version {
	v, err := NewVersion(version)
	if err != nil {
		t.Fatal(err)
	}
	return v

}
