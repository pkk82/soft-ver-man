package domain

import (
	"github.com/yudai/gojsondiff"
	"reflect"
	"testing"
)

func Test_DeserializeInstalledPackages(t *testing.T) {
	type args struct {
		name string
		json string
	}
	tests := []struct {
		name    string
		args    args
		want    InstalledPackages
		wantErr bool
	}{
		{
			name:    "empty string",
			args:    args{name: "node", json: ""},
			want:    InstalledPackages{Items: []InstalledPackage{}},
			wantErr: false,
		},
		{
			name:    "empty json",
			args:    args{name: "node", json: "{}"},
			want:    InstalledPackages{Items: []InstalledPackage{}},
			wantErr: false,
		},
		{
			name: "domain",
			args: args{name: "node", json: `{"name": "node", "items": [{"version": "v20.1.3", "path": "/home/user/pf/node/node-v12.22.12-linux-x64", "main": true, "installedOn": 1689017267000}]}`},
			want: InstalledPackages{
				Items: []InstalledPackage{
					{
						Version:     Ver("v20.1.3", t),
						Path:        "/home/user/pf/node/node-v12.22.12-linux-x64",
						Main:        true,
						InstalledOn: 1689017267000,
					},
				}},
			wantErr: false,
		}, {
			name:    "wrong name",
			args:    args{name: "go", json: `{"name": "node", "items": [{"version": "v20.1.3", "path": "/home/user/pf/node/node-v12.22.12-linux-x64", "main": true, "installedOn": 1689017267000}]}`},
			want:    InstalledPackages{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeserializeInstalledPackages(tt.args.name, tt.args.json)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializeInstalledPackages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeserializeInstalledPackages() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_InstalledPackages_Serialize(t *testing.T) {

	tests := []struct {
		name              string
		installedPackages InstalledPackages
		want              string
		wantErr           bool
	}{
		{
			name:              "empty",
			installedPackages: InstalledPackages{},
			want:              `{"name":"","items":[]}`,
			wantErr:           false,
		},
		{
			name:              "name only",
			installedPackages: InstalledPackages{Plugin: Plugin{Name: "node"}},
			want:              `{"name":"node","items":[]}`,
			wantErr:           false,
		},
		{
			name: "domain",
			installedPackages: InstalledPackages{Plugin: Plugin{Name: "node"}, Items: []InstalledPackage{
				{
					Version:     Ver("v20.1.3", t),
					Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
					Main:        true,
					InstalledOn: 1689017267000,
				}, {
					Version:     Ver("v20.1.4", t),
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
			got, err := tt.installedPackages.SerializeInstalledPackages()
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeInstalledPackages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			compare, err := differ.Compare([]byte(got), []byte(tt.want))
			if err != nil {
				t.Errorf("differ.Compare() error = %v", err)
			}
			if len(compare.Deltas()) > 0 {
				t.Errorf("SerializeInstalledPackages() = %v, want %v", got, tt.want)
			}
		})
	}
}
