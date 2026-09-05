/*
 * Copyright © 2026 Piotr Kozak <piotrkrzysztofkozak@gmail.com>
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

package gradle

import (
	"github.com/pkk82/soft-ver-man/domain"
	"reflect"
	"testing"
)

func TestSupportedPackages(t *testing.T) {
	apiVersions := []ApiVersion{
		{Version: "9.8.0-20260905023534+0000", Snapshot: true, ReleaseNightly: true, DownloadUrl: "https://services.gradle.org/distributions-snapshots/gradle-9.8.0-20260905023534+0000-bin.zip", Checksum: "98a74620703a4ff671a893b460b5b54c6e1586df04bf53eed25a2c4f1f6de21a"},
		{Version: "9.9.0-20260905003642+0000", Snapshot: true, Nightly: true, DownloadUrl: "https://services.gradle.org/distributions-snapshots/gradle-9.9.0-20260905003642+0000-bin.zip", Checksum: "09fb0c4b6aa054ecea85b75df91cd186d38079ce626a82f913642eed0066e1d5"},
		{Version: "9.8.0-milestone-2", MilestoneFor: "9.8.0", DownloadUrl: "https://services.gradle.org/distributions/gradle-9.8.0-milestone-2-bin.zip", Checksum: "1111111111111111111111111111111111111111111111111111111111111111"},
		{Version: "8.14.5", DownloadUrl: "https://services.gradle.org/distributions/gradle-8.14.5-bin.zip", Checksum: "2222222222222222222222222222222222222222222222222222222222222222"},
		{Version: "9.7.1", DownloadUrl: "https://services.gradle.org/distributions/gradle-9.7.1-bin.zip", Checksum: "acd53f1edaf02f1a8ff99879f8a34b302661a057d9b063ae9e35b552f804d20a"},
		{Version: "9.7.0-rc-3", RcFor: "9.7.0", DownloadUrl: "https://services.gradle.org/distributions/gradle-9.7.0-rc-3-bin.zip", Checksum: "3333333333333333333333333333333333333333333333333333333333333333"},
		{Version: "9.7.0", DownloadUrl: "https://services.gradle.org/distributions/gradle-9.7.0-bin.zip", Checksum: "4444444444444444444444444444444444444444444444444444444444444444"},
		{Version: "8.14", DownloadUrl: "https://services.gradle.org/distributions/gradle-8.14-bin.zip", Checksum: "5555555555555555555555555555555555555555555555555555555555555555"},
		{Version: "1.0-milestone-9", DownloadUrl: "https://services.gradle.org/distributions/gradle-1.0-milestone-9-bin.zip", Checksum: "6666666666666666666666666666666666666666666666666666666666666666"},
		{Version: "1.0-milestone-4", Broken: true, DownloadUrl: "https://services.gradle.org/distributions/gradle-1.0-milestone-4-bin.zip", Checksum: "7777777777777777777777777777777777777777777777777777777777777777"},
		{Version: "0.9", DownloadUrl: "https://services.gradle.org/distributions/gradle-0.9-bin.zip", Checksum: "8888888888888888888888888888888888888888888888888888888888888888"},
	}

	actual := supportedPackages(apiVersions)
	expected := []Package{
		{Version: domain.Ver("9.7.1", t), DownloadLink: "https://services.gradle.org/distributions/gradle-9.7.1-bin.zip", Sha256: "acd53f1edaf02f1a8ff99879f8a34b302661a057d9b063ae9e35b552f804d20a"},
		{Version: domain.Ver("9.7.0", t), DownloadLink: "https://services.gradle.org/distributions/gradle-9.7.0-bin.zip", Sha256: "4444444444444444444444444444444444444444444444444444444444444444"},
		{Version: domain.Ver("8.14.5", t), DownloadLink: "https://services.gradle.org/distributions/gradle-8.14.5-bin.zip", Sha256: "2222222222222222222222222222222222222222222222222222222222222222"},
		{Version: domain.Ver("8.14", t), DownloadLink: "https://services.gradle.org/distributions/gradle-8.14-bin.zip", Sha256: "5555555555555555555555555555555555555555555555555555555555555555"},
		{Version: domain.Ver("0.9", t), DownloadLink: "https://services.gradle.org/distributions/gradle-0.9-bin.zip", Sha256: "8888888888888888888888888888888888888888888888888888888888888888"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v", expected, actual)
	}
}

func TestSupportedPackages_empty(t *testing.T) {
	actual := supportedPackages([]ApiVersion{})
	if len(actual) != 0 {
		t.Errorf("Expected no packages, got: %v", actual)
	}
}
