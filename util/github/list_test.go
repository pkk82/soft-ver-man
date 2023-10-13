package github

import "testing"

func Test_nextPageLink(t *testing.T) {
	type args struct {
		linkHeader string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "next page link exists",
			args: args{
				linkHeader: "<https://api.github.com/repositories/24292831/releases?page=2>; rel=\"next\", <https://api.github.com/repositories/24292831/releases?page=2>; rel=\"last\"",
			},
			want: "https://api.github.com/repositories/24292831/releases?page=2",
		}, {
			name: "next page link does not exist",
			args: args{
				linkHeader: "<https://api.github.com/repositories/3432266/releases?per_page=100&page=2>; rel=\"prev\", <https://api.github.com/repositories/3432266/releases?per_page=100&page=1>; rel=\"first\"",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextPageLink(tt.args.linkHeader); got != tt.want {
				t.Errorf("nextPageLink() = %v, want %v", got, tt.want)
			}
		})
	}
}
