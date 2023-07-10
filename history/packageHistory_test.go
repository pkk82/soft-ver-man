package history

import (
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
