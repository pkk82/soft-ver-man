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

package golang

import (
	"github.com/pkk82/soft-ver-man/domain"
	"reflect"
	"testing"
)

func TestSupportedFiles(t *testing.T) {
	filesPerVersions := []ApiPackagesPerVersion{
		{Version: "go1.25rc2", Files: []ApiPackage{
			{Version: "go1.25rc2", Filename: "go1.25rc2.src.tar.gz", Kind: "source", Os: "", Arch: "", Sha256: "e6314a3234c4c03b8d006bcd1d14594d928170d1844b6dffdd5fa5a2333449e1"},
			{Version: "go1.25rc2", Filename: "go1.25rc2.aix-ppc64.tar.gz", Kind: "archive", Os: "aix", Arch: "ppc64", Sha256: "8ab16ba3cf322d9fa2dabb249a32dedb2c5c34f2ce9a38d01c9349ed9a1deec9"},
			{Version: "go1.25rc2", Filename: "go1.25rc2.darwin-amd64.tar.gz", Kind: "archive", Os: "darwin", Arch: "amd64", Sha256: "a09e19acff22863dfad464aee1b8b83689b75b233d65406b1d9e63d1dea21296"},
			{Version: "go1.25rc2", Filename: "go1.25rc2.darwin-amd64.pkg", Kind: "installer", Os: "darwin", Arch: "amd64", Sha256: "534dc80395ae11fc9a0b2a136fec48d70b2cf042afddbd16beb191980cccf636"},
			{Version: "go1.25rc2", Filename: "go1.25rc2.darwin-arm64.tar.gz", Kind: "archive", Os: "darwin", Arch: "arm64", Sha256: "995979864ed7a80a81fc5c20381da9973f8528f3776ddcebf542d888e72991d2"},
			{Version: "go1.25rc2", Filename: "go1.25rc2.darwin-arm64.pkg", Kind: "installer", Os: "darwin", Arch: "arm64", Sha256: "b8f1a1bc317d86dbf319fa382eb2b5552a493143d94f0aa14ea0ddb017f47314"},
			{Version: "go1.25rc2", Filename: "go1.25rc2.linux-386.tar.gz", Kind: "archive", Os: "linux", Arch: "386", Sha256: "918ed022a37e1f8e6cf11d2285084b45c6f83818592b2c7cd925449d756e873a"},
			{Version: "go1.25rc2", Filename: "go1.25rc2.linux-amd64.tar.gz", Kind: "archive", Os: "linux", Arch: "amd64", Sha256: "efcd3a151b174ffebde86b9d391ad59084300a4c5e9ea8c1d5dff90bbac38820"},
			{Version: "go1.25rc2", Filename: "go1.25rc2.linux-arm64.tar.gz", Kind: "archive", Os: "linux", Arch: "arm64", Sha256: "e9a077cef12d5c4a82df6d85a76f5bb7a4abd69c7d0fbab89d072591ef219ed3"},
			{Version: "go1.25rc2", Filename: "go1.25rc2.windows-386.zip", Kind: "archive", Os: "windows", Arch: "386", Sha256: "f1eba3975419b0b693bb24064f8bd775ec0b7f3755413fb64d133b48f517279e"},
			{Version: "go1.25rc2", Filename: "go1.25rc2.windows-386.msi", Kind: "installer", Os: "windows", Arch: "386", Sha256: "f0096f64f17e8f5c96ba868b89b32d85780141819820ed04b97d8efea7ec0c1f"},
			{Version: "go1.25rc2", Filename: "go1.25rc2.windows-amd64.zip", Kind: "archive", Os: "windows", Arch: "amd64", Sha256: "f0dadd0cdebbf56ad4ae7a010a5e980b64d6e4ddd13d6d294c0dad9e42285d09"},
			{Version: "go1.25rc2", Filename: "go1.25rc2.windows-amd64.msi", Kind: "installer", Os: "windows", Arch: "amd64", Sha256: "1e531e7318fce43040bdac041009b6ab2d1b8a8581bc94939aa101354d33ed61"},
			{Version: "go1.25rc2", Filename: "go1.25rc2.windows-arm64.zip", Kind: "archive", Os: "windows", Arch: "arm64", Sha256: "1d1603d5c3250da2c10b3decf0da279e324e9effe45fd4e775005c4a8022bf6a"},
			{Version: "go1.25rc2", Filename: "go1.25rc2.windows-arm64.msi", Kind: "installer", Os: "windows", Arch: "arm64", Sha256: "70480f81c941a8f51e2a796db086c21f4e7a157f076fe375f8e063850530e36f"},
		}},

		{Version: "go1.24.5", Files: []ApiPackage{
			{Version: "go1.24.5", Filename: "go1.24.5.src.tar.gz", Kind: "source", Os: "", Arch: "", Sha256: "74fdb09f2352e2b25b7943e56836c9b47363d28dec1c8b56c4a9570f30b8f59f"},
			{Version: "go1.24.5", Filename: "go1.24.5.aix-ppc64.tar.gz", Kind: "archive", Os: "aix", Arch: "ppc64", Sha256: "94c73c05660a7f8fa769ee89ee80a6792442872ce5e9d5b9999f4baae9aaeba2"},
			{Version: "go1.24.5", Filename: "go1.24.5.darwin-amd64.tar.gz", Kind: "archive", Os: "darwin", Arch: "amd64", Sha256: "2fe5f3866b8fbcd20625d531f81019e574376b8a840b0a096d8a2180308b1672"},
			{Version: "go1.24.5", Filename: "go1.24.5.darwin-amd64.pkg", Kind: "installer", Os: "darwin", Arch: "amd64", Sha256: "8dd21674c7845306b955d82d7c67696448e04f8ba31c0f701d623d71cbfa3d55"},
			{Version: "go1.24.5", Filename: "go1.24.5.darwin-arm64.tar.gz", Kind: "archive", Os: "darwin", Arch: "arm64", Sha256: "92d30a678f306c327c544758f2d2fa5515aa60abe9dba4ca35fbf9b8bfc53212"},
			{Version: "go1.24.5", Filename: "go1.24.5.darwin-arm64.pkg", Kind: "installer", Os: "darwin", Arch: "arm64", Sha256: "7fbdc0775c1cee6b1a1625b9babeea10c0ee03c520a1c61b354f4afee35b7977"},
			{Version: "go1.24.5", Filename: "go1.24.5.linux-386.tar.gz", Kind: "archive", Os: "linux", Arch: "386", Sha256: "ddcd926755a9e1aa66baf16c42cf705fc00defa4bdf3225f1676b7672c9a46fa"},
			{Version: "go1.24.5", Filename: "go1.24.5.linux-amd64.tar.gz", Kind: "archive", Os: "linux", Arch: "amd64", Sha256: "10ad9e86233e74c0f6590fe5426895de6bf388964210eac34a6d83f38918ecdc"},
			{Version: "go1.24.5", Filename: "go1.24.5.linux-arm64.tar.gz", Kind: "archive", Os: "linux", Arch: "arm64", Sha256: "0df02e6aeb3d3c06c95ff201d575907c736d6c62cfa4b6934c11203f1d600ffa"},
			{Version: "go1.24.5", Filename: "go1.24.5.windows-386.zip", Kind: "archive", Os: "windows", Arch: "386", Sha256: "f1eba3975419b0b693bb24064f8bd775ec0b7f3755413fb64d133b48f517279e"},
			{Version: "go1.24.5", Filename: "go1.24.5.windows-386.msi", Kind: "installer", Os: "windows", Arch: "386", Sha256: "eef0329bb25770ce2392c56a70e6074baff3c2ec38bb2619a6a199048617563a"},
			{Version: "go1.24.5", Filename: "go1.24.5.windows-amd64.zip", Kind: "archive", Os: "windows", Arch: "amd64", Sha256: "658f432689106d4e0a401a2ebb522b1213f497bc8357142fe8def18d79f02957"},
			{Version: "go1.24.5", Filename: "go1.24.5.windows-amd64.msi", Kind: "installer", Os: "windows", Arch: "amd64", Sha256: "3394ab0a830727764b6a576cd84e1fe7640919730b706d42e5948af35b46cfc8"},
			{Version: "go1.24.5", Filename: "go1.24.5.windows-arm64.zip", Kind: "archive", Os: "windows", Arch: "arm64", Sha256: "cd2955c4e3166a0cef4b76830025e4cc6e9ecccff32c02979a63f534d83c2e66"},
			{Version: "go1.24.5", Filename: "go1.24.5.windows-arm64.msi", Kind: "installer", Os: "windows", Arch: "arm64", Sha256: "f236821b94b4db3febc7ab36ac9fc477c14d449fa53ec6debb7dbd3c10f9128a"},
		}},
		{Version: "go1.23.11", Files: []ApiPackage{
			{Version: "go1.23.11", Filename: "go1.23.11.src.tar.gz", Kind: "source", Os: "", Arch: "", Sha256: "296381607a483a8a8667d7695331752f94a1f231c204e2527d2f22e1e3d1247d"},
			{Version: "go1.23.11", Filename: "go1.23.11.aix-ppc64.tar.gz", Kind: "archive", Os: "aix", Arch: "ppc64", Sha256: "cabf933f78628c020d8719286dd75fd85a9cc8a9b1ec2f91aed2429bccd11e4c"},
			{Version: "go1.23.11", Filename: "go1.23.11.darwin-amd64.tar.gz", Kind: "archive", Os: "darwin", Arch: "amd64", Sha256: "804538b068ebf449789e060d221c7be94d92d5f3e86842071cc70148d677f84d"},
			{Version: "go1.23.11", Filename: "go1.23.11.darwin-amd64.pkg", Kind: "installer", Os: "darwin", Arch: "amd64", Sha256: "0d15b2ceebb25ce3ee8f40c165c15a76b9fd0f39711e58722f487734676d1931"},
			{Version: "go1.23.11", Filename: "go1.23.11.darwin-arm64.tar.gz", Kind: "archive", Os: "darwin", Arch: "arm64", Sha256: "d3c2c69a79eb3e2a06e5d8bbca692c9166b27421f7251ccbafcada0ba35a05ee"},
			{Version: "go1.23.11", Filename: "go1.23.11.darwin-arm64.pkg", Kind: "installer", Os: "darwin", Arch: "arm64", Sha256: "963172d6feea43221df25e4449197e65f7e3f5c779da6d7c5aa8795246b9ebdc"},
			{Version: "go1.23.11", Filename: "go1.23.11.linux-386.tar.gz", Kind: "archive", Os: "linux", Arch: "386", Sha256: "3d3ed184a639e71b91450d24d1a2dab050e7cf0c8ed4c83f6a2daba8cd5fe1be"},
			{Version: "go1.23.11", Filename: "go1.23.11.linux-amd64.tar.gz", Kind: "archive", Os: "linux", Arch: "amd64", Sha256: "80899df77459e0b551d2eb8800ad6eb47023b99cccbf8129e7b5786770b948c5"},
			{Version: "go1.23.11", Filename: "go1.23.11.linux-arm64.tar.gz", Kind: "archive", Os: "linux", Arch: "arm64", Sha256: "1085c6ff805ec1f4893fa92013d16e58f74aeac830b1b9919b6908f3ed1a85c5"},
			{Version: "go1.23.11", Filename: "go1.23.11.windows-386.zip", Kind: "archive", Os: "windows", Arch: "386", Sha256: "64abdd62ad1c9cf49b41fe457be3481537366a81295ed35ca4d1d08784b35188"},
			{Version: "go1.23.11", Filename: "go1.23.11.windows-386.msi", Kind: "installer", Os: "windows", Arch: "386", Sha256: "a06ed55d6f0e9c71fa2c99fccd31f29bbd9f421931e6641d9983a16d276c4b48"},
			{Version: "go1.23.11", Filename: "go1.23.11.windows-amd64.zip", Kind: "archive", Os: "windows", Arch: "amd64", Sha256: "1dbcf0b4183066550964b22890fe119b0b867b51f12c1eea4445c71494d98cbb"},
			{Version: "go1.23.11", Filename: "go1.23.11.windows-amd64.msi", Kind: "installer", Os: "windows", Arch: "amd64", Sha256: "69917156ff828e7c8c0392afc862edf19c5def6f622a06a47cc54a79be8f0d9c"},
			{Version: "go1.23.11", Filename: "go1.23.11.windows-arm64.zip", Kind: "archive", Os: "windows", Arch: "arm64", Sha256: "83b81a03b51b92a26e711cf28d0633266c82d503d1c0d4f8fcaecf096b9df42d"},
			{Version: "go1.23.11", Filename: "go1.23.11.windows-arm64.msi", Kind: "installer", Os: "windows", Arch: "arm64", Sha256: "608583ca95134f3afb20283ee4a8ac639e503f65660c045735f4d9ebdc660461"},
		}}}

	actual := supportedPackages(&filesPerVersions, "linux", "amd64")
	expected := []Package{
		{Version: "1.25.rc2", Filename: "go1.25rc2.linux-amd64.tar.gz", Type: domain.TAR_GZ, DownloadLink: "https://go.dev/dl/go1.25rc2.linux-amd64.tar.gz", Sha256: "efcd3a151b174ffebde86b9d391ad59084300a4c5e9ea8c1d5dff90bbac38820"},
		{Version: "1.24.5", Filename: "go1.24.5.linux-amd64.tar.gz", Type: domain.TAR_GZ, DownloadLink: "https://go.dev/dl/go1.24.5.linux-amd64.tar.gz", Sha256: "10ad9e86233e74c0f6590fe5426895de6bf388964210eac34a6d83f38918ecdc"},
		{Version: "1.23.11", Filename: "go1.23.11.linux-amd64.tar.gz", Type: domain.TAR_GZ, DownloadLink: "https://go.dev/dl/go1.23.11.linux-amd64.tar.gz", Sha256: "80899df77459e0b551d2eb8800ad6eb47023b99cccbf8129e7b5786770b948c5"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v", expected, actual)
	}

	actual = supportedPackages(&filesPerVersions, "windows", "amd64")
	expected = []Package{
		{Version: "1.25.rc2", Filename: "go1.25rc2.windows-amd64.zip", Type: domain.ZIP, DownloadLink: "https://go.dev/dl/go1.25rc2.windows-amd64.zip", Sha256: "f0dadd0cdebbf56ad4ae7a010a5e980b64d6e4ddd13d6d294c0dad9e42285d09"},
		{Version: "1.24.5", Filename: "go1.24.5.windows-amd64.zip", Type: domain.ZIP, DownloadLink: "https://go.dev/dl/go1.24.5.windows-amd64.zip", Sha256: "658f432689106d4e0a401a2ebb522b1213f497bc8357142fe8def18d79f02957"},
		{Version: "1.23.11", Filename: "go1.23.11.windows-amd64.zip", Type: domain.ZIP, DownloadLink: "https://go.dev/dl/go1.23.11.windows-amd64.zip", Sha256: "1dbcf0b4183066550964b22890fe119b0b867b51f12c1eea4445c71494d98cbb"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v", expected, actual)
	}

	actual = supportedPackages(&filesPerVersions, "darwin", "amd64")
	expected = []Package{
		{Version: "1.25.rc2", Filename: "go1.25rc2.darwin-amd64.tar.gz", Type: domain.TAR_GZ, DownloadLink: "https://go.dev/dl/go1.25rc2.darwin-amd64.tar.gz", Sha256: "a09e19acff22863dfad464aee1b8b83689b75b233d65406b1d9e63d1dea21296"},
		{Version: "1.24.5", Filename: "go1.24.5.darwin-amd64.tar.gz", Type: domain.TAR_GZ, DownloadLink: "https://go.dev/dl/go1.24.5.darwin-amd64.tar.gz", Sha256: "2fe5f3866b8fbcd20625d531f81019e574376b8a840b0a096d8a2180308b1672"},
		{Version: "1.23.11", Filename: "go1.23.11.darwin-amd64.tar.gz", Type: domain.TAR_GZ, DownloadLink: "https://go.dev/dl/go1.23.11.darwin-amd64.tar.gz", Sha256: "804538b068ebf449789e060d221c7be94d92d5f3e86842071cc70148d677f84d"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v", expected, actual)
	}

	actual = supportedPackages(&filesPerVersions, "darwin", "arm64")
	expected = []Package{
		{Version: "1.25.rc2", Filename: "go1.25rc2.darwin-arm64.tar.gz", Type: domain.TAR_GZ, DownloadLink: "https://go.dev/dl/go1.25rc2.darwin-arm64.tar.gz", Sha256: "995979864ed7a80a81fc5c20381da9973f8528f3776ddcebf542d888e72991d2"},
		{Version: "1.24.5", Filename: "go1.24.5.darwin-arm64.tar.gz", Type: domain.TAR_GZ, DownloadLink: "https://go.dev/dl/go1.24.5.darwin-arm64.tar.gz", Sha256: "92d30a678f306c327c544758f2d2fa5515aa60abe9dba4ca35fbf9b8bfc53212"},
		{Version: "1.23.11", Filename: "go1.23.11.darwin-arm64.tar.gz", Type: domain.TAR_GZ, DownloadLink: "https://go.dev/dl/go1.23.11.darwin-arm64.tar.gz", Sha256: "d3c2c69a79eb3e2a06e5d8bbca692c9166b27421f7251ccbafcada0ba35a05ee"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v", expected, actual)
	}
}
