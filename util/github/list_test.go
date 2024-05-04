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
