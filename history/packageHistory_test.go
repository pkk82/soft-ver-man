package history

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
			name: "history",
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
			name: "history",
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
