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

package node

const Name = "node"
const LongName = "Node.js"

var Aliases = []string{"node", "nodejs", "node.js"}

const DistURL = "https://nodejs.org/dist"
const JsonFileURL = "https://nodejs.org/dist/index.json"
const ShaSumFileName = "SHASUMS256.txt"
const ShaSumSigFileName = "SHASUMS256.txt.sig"

const gpgKeys = `
Beth Griggs <bethanyngriggs@gmail.com> 4ED778F539E3634C779C87C6D7062848A1AB005C
Bryan English <bryan@bryanenglish.com> 141F07595B7B3FFE74309A937405533BE57C7D57
Danielle Adams <adamzdanielle@gmail.com> 74F12602B6F1C4E913FAA37AD3A89613643B6201
Juan José Arboleda <soyjuanarbol@gmail.com> DD792F5973C6DE52C432CBDAC77ABFA00DDBF2B7
Michaël Zasso <targos@protonmail.com> 8FCCA13FEF1D0C2E91008E09770F7A9A5AE15600
Myles Borins <myles.borins@gmail.com> C4F0DFFF4E8C1A8236409D08E73BC641CC11F4C8
RafaelGSS <rafael.nunu@hotmail.com> 890C08DB8579162FEE0DF9DB8BEAB4DFCF555EF4
Richard Lau <rlau@redhat.com> C82FA3AE1CBEDC6BE46B9360C43CEC45C17AB93C
Ruy Adorno <ruyadorno@hotmail.com> 108F52B48DB57BB0CC439B2997B01419BD92F80A
`
