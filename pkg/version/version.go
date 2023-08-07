/*
 * Copyright (c) 2023 The nebula-contrib Authors.
 * Licensed under the Apache License, GitVersion 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package version

import "fmt"

var (
	VerMajor = 0
	VerMinor = 0
	VerPatch = 1
	VerName  = "Nebula Operator Command Line Tool"
	GitSha   = "UNKNOWN"
	GitRef   = "UNKNOWN"
)

func GetVersion() string {
	return fmt.Sprintf(`%s,V-%d.%d.%d [GitSha: %s GitRef: %s]`,
		VerName, VerMajor, VerMinor, VerPatch, GitSha, GitRef)
}
