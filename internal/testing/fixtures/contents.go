// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fixtures

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetDir returns the directory where test data for this project resides.
func GetDir() string {
	cwd, _ := os.Getwd()
	parts := strings.Split(cwd, "/")
	backtrackCount := 0

	for i, v := range parts {
		if v == "api" || v == "cmd" || v == "pkg" || v == "internal" || v == "test" {
			backtrackCount = len(parts) - i
			break
		}
	}

	backtrackedPath := cwd
	for i := 0; i < backtrackCount; i++ {
		backtrackedPath = fmt.Sprintf("%s/..", backtrackedPath)
	}

	path, _ := filepath.Abs(fmt.Sprintf("%s/test/data", backtrackedPath))

	return path
}

// GetPath returns the absolute path where the specific test data file resides.
func GetPath(relativePath string) string {
	return fmt.Sprintf("%s/%s", GetDir(), relativePath)
}
