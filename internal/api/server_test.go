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

package api_test

import (
	"context"
	"github.com/m4dcoder/nihao/internal/api"
	"github.com/m4dcoder/nihao/internal/testing/fixtures"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

func mockHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP request is handled.")
}

var mockRoutes = api.Routes{
	api.Route{
		Name:        "HandleMock",
		Method:      "POST",
		Pattern:     "/mock",
		HandlerFunc: mockHandler,
	},
}

func TestRunServer(t *testing.T) {
	server := api.NewServer("8078", api.NewRouter(mockRoutes), false, "", "")

	err := server.Run()
	assert.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	assert.Nil(t, err)
}

func TestRunServerTLS(t *testing.T) {
	tlsCrtPath := fixtures.GetPath("certs/ssl/nihao.crt")
	tlsKeyPath := fixtures.GetPath("certs/ssl/nihao.key")

	server := api.NewServer("8078", api.NewRouter(mockRoutes), true, tlsCrtPath, tlsKeyPath)

	err := server.Run()
	assert.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	assert.Nil(t, err)
}

func TestRunServerTLSBadCrtPath(t *testing.T) {
	tlsCrtPath := fixtures.GetPath("certs/ssl/nihao.foo")
	tlsKeyPath := fixtures.GetPath("certs/ssl/nihao.key")

	server := api.NewServer("8078", api.NewRouter(mockRoutes), true, tlsCrtPath, tlsKeyPath)

	err := server.Run()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "does not exist or is not a file")
}

func TestRunServerTLSBadKeyPath(t *testing.T) {
	tlsCrtPath := fixtures.GetPath("certs/ssl/nihao.crt")
	tlsKeyPath := fixtures.GetPath("certs/ssl/nihao.foo")

	server := api.NewServer("8078", api.NewRouter(mockRoutes), true, tlsCrtPath, tlsKeyPath)

	err := server.Run()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "does not exist or is not a file")
}

func TestRunServerTLSCrtPathIsDir(t *testing.T) {
	tlsCrtPath := fixtures.GetPath("certs/ssl")
	tlsKeyPath := fixtures.GetPath("certs/ssl/nihao.key")

	server := api.NewServer("8078", api.NewRouter(mockRoutes), true, tlsCrtPath, tlsKeyPath)

	err := server.Run()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "does not exist or is not a file")
}

func TestRunServerTLSKeyPathIsDir(t *testing.T) {
	tlsCrtPath := fixtures.GetPath("certs/ssl/nihao.crt")
	tlsKeyPath := fixtures.GetPath("certs/ssl")

	server := api.NewServer("8078", api.NewRouter(mockRoutes), true, tlsCrtPath, tlsKeyPath)

	err := server.Run()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "does not exist or is not a file")
}

func TestRunServerTLSBadCertContent(t *testing.T) {
	file, _ := ioutil.TempFile("/tmp", "cacert")
	defer os.Remove(file.Name())

	// Write some bad content into the cacert file
	f, _ := os.OpenFile(file.Name(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString("你好, Metaverse!\n")

	tlsCrtPath := file.Name()
	tlsKeyPath := fixtures.GetPath("certs/ssl/nihao.key")

	server := api.NewServer("8078", api.NewRouter(mockRoutes), true, tlsCrtPath, tlsKeyPath)

	err := server.Run()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to find any PEM data in certificate input")
}

func TestRunServerLaunchFailed(t *testing.T) {
	server := api.NewServer("1000000", api.NewRouter(mockRoutes), false, "", "")

	err := server.Run()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "listen tcp: address 1000000: invalid port")
}
