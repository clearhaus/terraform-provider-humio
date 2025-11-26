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
	"os"
	"strconv"
	"testing"

	"github.com/clearhaus/terraform-provider-humio/humio/acceptance"
)

func TestProviderInternalValidation(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestMain(m *testing.M) {
	if tfAccVal, ok := os.LookupEnv("TF_ACC"); ok {
		// Check for presence in the environment
		_, addrSet := os.LookupEnv("HUMIO_ADDR")
		_, tokenSet := os.LookupEnv("HUMIO_API_TOKEN")
		manuallySet := addrSet || tokenSet

		if shouldRun, _ := strconv.ParseBool(tfAccVal); shouldRun {
			// If externally configured, assume that spinning up a Docker
			// instance of Humio is wasteful
			if manuallySet {
				m.Run()
			} else {
				acceptance.RunWithInstance(func(addr string, token string) int {
					_ = os.Setenv("HUMIO_ADDR", addr)
					_ = os.Setenv("HUMIO_API_TOKEN", token)

					return m.Run()
				})
			}
		}
	}
}
