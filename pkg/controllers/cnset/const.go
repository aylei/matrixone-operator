// Copyright 2022 Matrix Origin
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

package cnset

const (
	servicePort   = 1001
	listenAddress = ""
	logLevel      = "debug"
	configFile    = "cn-config.toml"
	logFormatType = "json"
	logMaxSize    = "512"
	backendType   = "s3"
	hostSize      = 1000
	guestSize     = 2000
	operatorSize  = 3000
	configVolume  = "config"
	dataVolume    = "cndata"
	configPath    = ""
	Entrypoint    = ""
	dataPath      = '"'
	batchRow      = 300
	batchSize     = 400
)
