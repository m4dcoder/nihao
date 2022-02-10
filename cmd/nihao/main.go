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

package main

import (
	"context"
	"github.com/m4dcoder/nihao/internal/api"
	"github.com/m4dcoder/nihao/internal/api/handlers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	port string = "6688"
	tls  bool   = false
	cert string = ""
	key  string = ""
)

var server *api.Server

var routes = api.Routes{
	api.Route{
		Name:        "ReplyHello",
		Method:      "GET",
		Pattern:     "/hello",
		HandlerFunc: handlers.ReplyHello,
	},
}

func setup(flags *pflag.FlagSet) {
	log.Infof("Starting the API server on port %s.", port)

	server = api.NewServer(port, api.NewRouter(routes), tls, cert, key)

	if err := server.Run(); err != nil {
		log.Errorf("Unable to start the API server on port %s because %s.", port, err.Error())
		os.Exit(1)
	}

	log.Infof("Successfully started the API server on port %s.", port)
}

func waitShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signals
	sig := <-sigs
	log.Infof("Shutdown request received: %s", sig)
}

func cleanup() {
	log.Info("Shutting down the API server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the API server.
	if server != nil {
		if err := server.Shutdown(ctx); err != nil {
			log.Errorf("Failed to shutdown the API server because %s.", err.Error())
			return
		}
	}

	log.Info("Successfully shutdown the API server.")
}

func main() {
	// Run setup.
	setup(pflag.CommandLine)

	// Wait for shutdown.
	waitShutdown()

	// Run cleanup.
	cleanup()
}
