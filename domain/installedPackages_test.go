/*
 * Copyright © 2024 Piotr Kozak <piotrkrzysztofkozak@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */
package domain

import (
	"reflect"
	"testing"
)

func Test_InstalledPackages_IsInstalled(t *testing.T) {

	tests := []struct {
		name              string
		installedPackages InstalledPackages
		version           Version
		want              bool
	}{
		{
			name:              "empty - false",
			installedPackages: InstalledPackages{Plugin: Plugin{Name: "node"}, Items: []InstalledPackage{}},
			version:           Ver("v20.1.4", t),
			want:              false,
		},
		{
			name: "one element - true",
			installedPackages: InstalledPackages{Plugin: Plugin{Name: "node"}, Items: []InstalledPackage{
				{
					Version:     Ver("v20.1.3", t),
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Main:        true,
					InstalledOn: 1689017267000,
				},
			}},
			version: Ver("v20.1.3", t),
			want:    true,
		}, {
			name: "one element -  false",
			installedPackages: InstalledPackages{Plugin: Plugin{Name: "node"}, Items: []InstalledPackage{
				{
					Version:     Ver("v20.1.3", t),
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Main:        true,
					InstalledOn: 1689017267000,
				},
			}},
			version: Version{Value: "v20.1.4"},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installedPackages := tt.installedPackages
			got := installedPackages.IsInstalled(tt.version)
			if got != tt.want {
				t.Errorf("InstalledPackages.IsInstalled() = %v, want %v", got, tt.want)
			}
		})
	}

}

func Test_InstalledPackages_Add(t *testing.T) {

	tests := []struct {
		name              string
		installedPackages InstalledPackages
		newPackage        InstalledPackage
		want              []InstalledPackage
	}{
		{
			name:              "empty",
			installedPackages: InstalledPackages{Plugin: Plugin{Name: "node"}, Items: []InstalledPackage{}},
			newPackage: InstalledPackage{
				Version:     Ver("v20.1.4", t),
				Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
				InstalledOn: 1689017268000,
				Main:        false,
			},
			want: []InstalledPackage{
				{
					Version:     Ver("v20.1.4", t),
					Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
					Main:        false,
					InstalledOn: 1689017268000,
				},
			},
		},
		{
			name:              "empty - add main",
			installedPackages: InstalledPackages{Plugin: Plugin{Name: "node"}, Items: []InstalledPackage{}},
			newPackage: InstalledPackage{
				Version:     Ver("v20.1.4", t),
				Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
				InstalledOn: 1689017268000,
				Main:        true,
			},
			want: []InstalledPackage{
				{
					Version:     Ver("v20.1.4", t),
					Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
					Main:        true,
					InstalledOn: 1689017268000,
				},
			},
		},
		{
			name: "with-one-item ",
			installedPackages: InstalledPackages{Plugin: Plugin{Name: "node"}, Items: []InstalledPackage{
				{
					Version:     Ver("v20.1.3", t),
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Main:        true,
					InstalledOn: 1689017267000,
				},
			}},
			newPackage: InstalledPackage{
				Version:     Ver("v20.1.4", t),
				Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
				InstalledOn: 1689017268000,
				Main:        false,
			},
			want: []InstalledPackage{{
				Version:     Ver("v20.1.3", t),
				Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
				Main:        true,
				InstalledOn: 1689017267000,
			}, {
				Version:     Ver("v20.1.4", t),
				Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
				Main:        false,
				InstalledOn: 1689017268000,
			}},
		}, {
			name: "with-one-item - add main",
			installedPackages: InstalledPackages{Plugin: Plugin{Name: "node"}, Items: []InstalledPackage{
				{
					Version:     Ver("v20.1.3", t),
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Main:        true,
					InstalledOn: 1689017267000,
				},
			}},
			newPackage: InstalledPackage{
				Version:     Ver("v20.1.4", t),
				Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
				InstalledOn: 1689017268000,
				Main:        true,
			},
			want: []InstalledPackage{{
				Version:     Ver("v20.1.3", t),
				Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
				Main:        false,
				InstalledOn: 1689017267000,
			}, {
				Version:     Ver("v20.1.4", t),
				Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
				Main:        true,
				InstalledOn: 1689017268000,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installedPackages := tt.installedPackages
			installedPackages.Add(tt.newPackage)
			if !reflect.DeepEqual(installedPackages.Items, tt.want) {
				t.Errorf("InstalledPackages.Add() = %v, want %v", installedPackages.Items, tt.want)
			}
		})
	}

}

func Test_InstalledPackages_Remove(t *testing.T) {

	tests := []struct {
		name              string
		installedPackages InstalledPackages
		want              []InstalledPackage
		wantRemoved       *InstalledPackage
	}{
		{
			name:              "empty",
			installedPackages: InstalledPackages{Plugin: Plugin{Name: "node"}, Items: []InstalledPackage{}},
			want:              []InstalledPackage{},
			wantRemoved:       nil,
		},
		{
			name: "with-one-item",
			installedPackages: InstalledPackages{Plugin: Plugin{Name: "node"}, Items: []InstalledPackage{
				{
					Version:     Ver("v20.1.3", t),
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Main:        true,
					InstalledOn: 1689017267000,
				},
			}},
			want: []InstalledPackage{},
			wantRemoved: &InstalledPackage{
				Version:     Ver("v20.1.3", t),
				Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
				Main:        true,
				InstalledOn: 1689017267000,
			},
		},
		{
			name: "with-two-items",
			installedPackages: InstalledPackages{Plugin: Plugin{Name: "node"}, Items: []InstalledPackage{
				{
					Version:     Ver("v20.1.3", t),
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Main:        true,
					InstalledOn: 1689017267000,
				},
				{
					Version:     Ver("v20.1.4", t),
					Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
					Main:        false,
					InstalledOn: 1689017267000,
				},
			}},
			want: []InstalledPackage{
				{
					Version:     Ver("v20.1.4", t),
					Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
					Main:        false,
					InstalledOn: 1689017267000,
				},
			},
			wantRemoved: &InstalledPackage{
				Version:     Ver("v20.1.3", t),
				Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
				Main:        true,
				InstalledOn: 1689017267000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			removed := tt.installedPackages.RemoveByVersion(Ver("v20.1.3", t))
			if !reflect.DeepEqual(tt.installedPackages.Items, tt.want) {
				t.Errorf("RemoveByVersion().Items = %v, want %v", tt.installedPackages.Items, tt.want)
			}

			if !reflect.DeepEqual(removed, tt.wantRemoved) {
				t.Errorf("RemoveByVersion().Items = %v, want %v", removed, tt.wantRemoved)
			}
		})
	}
}

func Test_InstalledPackages_FoundMain(t *testing.T) {
	tests := []struct {
		name    string
		arg     InstalledPackages
		want    *InstalledPackage
		wantErr bool
	}{
		{
			name:    "empty",
			arg:     InstalledPackages{Plugin: Plugin{Name: "node"}, Items: []InstalledPackage{}},
			wantErr: true,
			want:    nil,
		},
		{
			name: "one main",
			arg: InstalledPackages{Plugin: Plugin{Name: "node"}, Items: []InstalledPackage{
				{
					Version:     Ver("v20.1.3", t),
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Main:        true,
					InstalledOn: 1689017267000,
				},
			}},
			wantErr: false,
			want: &InstalledPackage{
				Version:     Ver("v20.1.3", t),
				Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
				Main:        true,
				InstalledOn: 1689017267000,
			},
		},
		{
			name: "many mains",
			arg: InstalledPackages{Plugin: Plugin{Name: "node"}, Items: []InstalledPackage{
				{
					Version:     Ver("v20.1.1", t),
					Path:        "/home/user/pf/node/node-v20.1.1-linux-x64",
					Main:        true,
					InstalledOn: 1689017266000,
				}, {
					Version:     Ver("v20.1.2", t),
					Path:        "/home/user/pf/node/node-v20.1.2-linux-x64",
					Main:        true,
					InstalledOn: 1689017269000,
				},
				{
					Version:     Ver("v20.1.3", t),
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Main:        false,
					InstalledOn: 1689017268000,
				},
				{
					Version:     Ver("v20.1.4", t),
					Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
					Main:        true,
					InstalledOn: 1689017267000,
				},
			}},
			wantErr: false,
			want: &InstalledPackage{
				Version:     Ver("v20.1.2", t),
				Path:        "/home/user/pf/node/node-v20.1.2-linux-x64",
				Main:        true,
				InstalledOn: 1689017269000,
			},
		}, {
			name: "no mains",
			arg: InstalledPackages{Plugin: Plugin{Name: "node"}, Items: []InstalledPackage{
				{
					Version:     Ver("v20.1.1", t),
					Path:        "/home/user/pf/node/node-v20.1.1-linux-x64",
					Main:        false,
					InstalledOn: 1689017266000,
				}, {
					Version:     Ver("v20.1.2", t),
					Path:        "/home/user/pf/node/node-v20.1.2-linux-x64",
					Main:        false,
					InstalledOn: 1689017269000,
				},
				{
					Version:     Ver("v20.1.3", t),
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Main:        false,
					InstalledOn: 1689017268000,
				},
				{
					Version:     Ver("v20.1.4", t),
					Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
					Main:        true,
					InstalledOn: 1689017267000,
				},
			}},
			wantErr: false,
			want: &InstalledPackage{
				Version:     Ver("v20.1.4", t),
				Path:        "/home/user/pf/node/node-v20.1.4-linux-x64",
				Main:        true,
				InstalledOn: 1689017267000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arg.FoundMain()
			if (err != nil) != tt.wantErr {
				t.Errorf("FoundMain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == nil && got != nil {
				t.Errorf("FoundMain() got = %v, want %v", got, tt.want)
			}
			if tt.want != nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FoundMain() got = %v, want %v", got, tt.want)
			}
		})
	}
}
