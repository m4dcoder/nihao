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

package api

import (
	"fmt"
	"github.com/coreos/go-systemd/activation"
	"github.com/gorilla/mux"
	"github.com/m4dcoder/nihao/internal/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Server defines the API server.
type Server struct {
	http.Server
	tls           bool
	tlsCrtPath    string
	tlsKeyPath    string
	lastError     error
	lastErrorTime time.Time
}

// Run starts the API server in the background.
func (s *Server) Run() error {
	if s.tls {
		if !(utils.FilePathExists(s.tlsCrtPath) && utils.IsFile(s.tlsCrtPath)) {
			return fmt.Errorf("the tls cert \"%s\" does not exist or is not a file", s.tlsCrtPath)
		}

		if !(utils.FilePathExists(s.tlsKeyPath) && utils.IsFile(s.tlsKeyPath)) {
			return fmt.Errorf("the tls key \"%s\" does not exist or is not a file", s.tlsKeyPath)
		}
	}

	go func() {
		serverType := "HTTP"

		if s.tls {
			serverType = "HTTPS"
		}

		log.Infof("%s server will listen at %s.", serverType, s.Addr)

		for {
			// Use socket activation for the http listener
			listeners, sockErr := activation.Listeners()

			if sockErr != nil {
				log.Errorf("Unable to activate listeners: %s", sockErr.Error())
				time.Sleep(3 * time.Second)
				continue
			}

			socketActivated := (len(listeners) > 0)

			if !socketActivated {
				message := "There are no socket activated. Please make sure that the systemd " +
					"socket is activated, either via the service.socket deployed to systemd or " +
					"manually using systemd-socket-activate. The server will continue to run " +
					"without socket activation."
				log.Warn(message)
			}

			// Run the http(s) server
			var err error

			if s.tls && socketActivated {
				log.Info("Launch the https server with socket activation.")
				err = s.ServeTLS(listeners[0], s.tlsCrtPath, s.tlsKeyPath)
			} else if s.tls && !socketActivated {
				log.Info("Launch the https server without socket activation.")
				err = s.ListenAndServeTLS(s.tlsCrtPath, s.tlsKeyPath)
			} else if socketActivated {
				log.Info("Launch the http server with socket activation.")
				err = s.Serve(listeners[0])
			} else {
				log.Info("Launch the http server without socket activation.")
				err = s.ListenAndServe()
			}

			if err != nil {
				s.lastErrorTime = time.Now().UTC()
				s.lastError = err
				log.Warnf("%s server returns: %s", serverType, err.Error())
			}

			log.Warnf("%s server stopped running.", serverType)
			time.Sleep(3 * time.Second)
		}
	}()

	time.Sleep(250 * time.Millisecond)

	if s.lastError != nil {
		return s.lastError
	}

	return nil
}

//NewServer returns a new instance of Server.
func NewServer(
	port string, router *mux.Router, tls bool,
	tlsCrtPath string, tlsKeyPath string) *Server {

	// Create a new API server instance.
	s := &Server{
		Server: http.Server{
			Addr:         ":" + port,
			ReadTimeout:  1000 * time.Second,
			WriteTimeout: 1000 * time.Second,
			Handler:      router,
		},
		tls:        tls,
		tlsCrtPath: tlsCrtPath,
		tlsKeyPath: tlsKeyPath,
	}

	return s
}
