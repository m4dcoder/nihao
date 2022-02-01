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

package handlers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// ReplyHello replies to the hello endpoint
func ReplyHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	resp := make(map[string]string)
	resp["message"] = "你好, Metaverse!"
	body, err := json.Marshal(resp)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err.Error())
		return
	}

	if _, err = w.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err.Error())
	}
}
