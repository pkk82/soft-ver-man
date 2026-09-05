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

const Name = "gradle"
const EnvNamePrefix = "GRADLE"
const EnvNameSuffix = "_HOME"
const LongName = "Gradle"

var Aliases = []string{}

// VersionsURL lists every published Gradle distribution together with its download and checksum
const VersionsURL = "https://services.gradle.org/versions/all"

// DownloadURL is the fallback location of the binary-only distribution when the versions API is unreachable
const DownloadURL = "https://services.gradle.org/distributions/gradle-%s-bin.zip"

// ChecksumURLSuffix appended to a distribution URL gives the location of its SHA-256 file
const ChecksumURLSuffix = ".sha256"

const Sha256Property = "sha256"
