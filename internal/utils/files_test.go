//go:build unit

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

package utils_test

import (
	"github.com/m4dcoder/nihao/internal/utils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestFilePathExists(t *testing.T) {
	file, _ := ioutil.TempFile("/tmp", "foobar")
	defer os.Remove(file.Name())
	assert.True(t, utils.FilePathExists(file.Name()))
}

func TestNotFilePathExists(t *testing.T) {
	assert.False(t, utils.FilePathExists("/this/path/does/not/exist"))
}

func TestIsDir(t *testing.T) {
	assert.True(t, utils.IsDir("/tmp"))
}

func TestNotIsDir(t *testing.T) {
	file, _ := ioutil.TempFile("/tmp", "foobar")
	defer os.Remove(file.Name())
	assert.False(t, utils.IsDir(file.Name()))
}

func TestNotExistIsDir(t *testing.T) {
	assert.False(t, utils.IsDir("/this/path/does/not/exist"))
}

func TestIsFile(t *testing.T) {
	file, _ := ioutil.TempFile("/tmp", "foobar")
	defer os.Remove(file.Name())
	assert.True(t, utils.IsFile(file.Name()))
}

func TestNotIsFile(t *testing.T) {
	assert.False(t, utils.IsFile("/tmp"))
}

func TestNotExistIsFile(t *testing.T) {
	assert.False(t, utils.IsFile("/this/path/does/not/exist"))
}
