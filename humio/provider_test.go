// Copyright © 2020 Humio Ltd.
//
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

package humio

import (
	"log"
	"os"
	"strconv"
	"testing"
)

func TestProviderInternalValidation(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMain(m *testing.M) {
	// The acceptance tests skip themselves unless TF_ACC is set, but when it is
	// set they need a Humio instance to talk to. One can be started for the
	// duration of the run with "cd test && go run ."; see the README.
	if tfAcc, _ := strconv.ParseBool(os.Getenv("TF_ACC")); tfAcc {
		if os.Getenv("HUMIO_ADDR") == "" || os.Getenv("HUMIO_API_TOKEN") == "" {
			log.Fatal("HUMIO_ADDR and HUMIO_API_TOKEN must be set to run the acceptance tests; " +
				"run \"cd test && go run .\" to start a throwaway Humio instance in Docker")
		}
	}

	os.Exit(m.Run())
}
