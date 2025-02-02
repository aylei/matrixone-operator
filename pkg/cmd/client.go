// Copyright 2023 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func GetCmdStatus(host string, port int) (*Status, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/status", host, port))
	if err != nil {
		return nil, errors.Wrap(err, "error polling backup status")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("error polling backup status, status code %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error polling backup status")
	}
	status := &Status{}
	err = json.Unmarshal(body, status)
	if err != nil {
		return nil, errors.Wrap(err, "error polling backup status")
	}
	return status, nil
}

func Stop(host string, port int) error {
	_, err := http.Post(fmt.Sprintf("http://%s:%d/shutdown", host, port), "", nil)
	if err != nil {
		return errors.Wrap(err, "error stopping command")
	}
	return nil
}
