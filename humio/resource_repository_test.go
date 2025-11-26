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
	"fmt"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	humio "github.com/humio/terraform-provider-humio/internal/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccRepositoryRequiredFields(t *testing.T) {
	config := repositoryEmpty
	accTestCase(t, []resource.TestStep{
		{
			Config:      config,
			ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found.`),
		},
		/* TODO: Figure out why it doesn't complain about retention being required as it is marked as required.
		         Perhaps we should not have it required in the first place?
		{
			Config:      config,
			ExpectError: regexp.MustCompile(`The argument "retention" is required, but no definition was found.`),
		},
		*/
	}, nil)
}

func TestAccRepositoryInvalidInputs(t *testing.T) {
	config := repositoryInvalidInputs
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "name"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "description"`)},
		{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "allow_data_deletion"`)},
		//{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "retention"`)},
	}, nil)
}

func TestAccRepositoryInvalidRetentionSettings(t *testing.T) {
	config := repositoryInvalidRetentionSettings
	accTestCase(t, []resource.TestStep{
		{Config: config, ExpectError: regexp.MustCompile(`Name:time_in_days}] to be at least \(0\.000000\), got -30\.000000`)},
		{Config: config, ExpectError: regexp.MustCompile(`Name:ingest_size_in_gb}] to be at least \(0\.000000\), got -10\.000000`)},
		{Config: config, ExpectError: regexp.MustCompile(`Name:storage_size_in_gb}] to be at least \(0\.000000\), got -5\.000000`)},
		//{Config: config, ExpectError: regexp.MustCompile(`Inappropriate value for attribute "retention"`)},
	}, nil)
}

func TestAccRepositoryBasic(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: repositoryBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_repository.test", "name", "repository-test"),
				resource.TestCheckResourceAttr("humio_repository.test", "description", ""),
				resource.TestCheckResourceAttr("humio_repository.test", "allow_data_deletion", "false"),
				resource.TestCheckResourceAttr("humio_repository.test", "retention.#", "1"), // TODO: Figure out if we want to require this set by the user. If not, how can we ensure this is not put in state?
				resource.TestCheckNoResourceAttr("humio_repository.test", "retention.0.time_in_days"),
				resource.TestCheckNoResourceAttr("humio_repository.test", "retention.0.ingest_size_in_gb"),
				resource.TestCheckNoResourceAttr("humio_repository.test", "retention.0.storage_size_in_gb"),
				resource.TestCheckNoResourceAttr("humio_repository.test", "time_in_days"),
				resource.TestCheckNoResourceAttr("humio_repository.test", "ingest_size_in_gb"),
				resource.TestCheckNoResourceAttr("humio_repository.test", "storage_size_in_gb"),
			),
		},
	}, testAccCheckRepositoryDestroy)
}

func TestAccRepositoryBasicToFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: repositoryBasic,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_repository.test", "name", "repository-test"),
				resource.TestCheckResourceAttr("humio_repository.test", "description", ""),
				resource.TestCheckResourceAttr("humio_repository.test", "allow_data_deletion", "false"),
				resource.TestCheckResourceAttr("humio_repository.test", "retention.#", "1"), // TODO: Figure out if we want to require this set by the user. If not, how can we ensure this is not put in state?
				resource.TestCheckNoResourceAttr("humio_repository.test", "retention.0.time_in_days"),
				resource.TestCheckNoResourceAttr("humio_repository.test", "retention.0.ingest_size_in_gb"),
				resource.TestCheckNoResourceAttr("humio_repository.test", "retention.0.storage_size_in_gb"),
				resource.TestCheckNoResourceAttr("humio_repository.test", "time_in_days"),
				resource.TestCheckNoResourceAttr("humio_repository.test", "ingest_size_in_gb"),
				resource.TestCheckNoResourceAttr("humio_repository.test", "storage_size_in_gb"),
			),
		},
		{
			Config: repositoryFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_repository.test", "name", "repository-test"),
				resource.TestCheckResourceAttr("humio_repository.test", "description", "some description"),
				resource.TestCheckResourceAttr("humio_repository.test", "allow_data_deletion", "true"),
				resource.TestCheckResourceAttr("humio_repository.test", "retention.#", "1"),
				resource.TestCheckResourceAttr("humio_repository.test", "retention.0.time_in_days", "30"),
				resource.TestCheckResourceAttr("humio_repository.test", "retention.0.ingest_size_in_gb", "10"),
				resource.TestCheckResourceAttr("humio_repository.test", "retention.0.storage_size_in_gb", "5"),
				resource.TestCheckNoResourceAttr("humio_repository.test", "time_in_days"),
				resource.TestCheckNoResourceAttr("humio_repository.test", "ingest_size_in_gb"),
				resource.TestCheckNoResourceAttr("humio_repository.test", "storage_size_in_gb"),
			),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
		},
		{
			Config: repositoryFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_repository.test", "name", "repository-test"),
				resource.TestCheckResourceAttr("humio_repository.test", "description", "some description"),
				resource.TestCheckResourceAttr("humio_repository.test", "allow_data_deletion", "true"),
				resource.TestCheckResourceAttr("humio_repository.test", "retention.#", "1"),
				resource.TestCheckResourceAttr("humio_repository.test", "retention.0.time_in_days", "30"),
				resource.TestCheckResourceAttr("humio_repository.test", "retention.0.ingest_size_in_gb", "10"),
				resource.TestCheckResourceAttr("humio_repository.test", "retention.0.storage_size_in_gb", "5"),
				resource.TestCheckNoResourceAttr("humio_repository.test", "time_in_days"),
				resource.TestCheckNoResourceAttr("humio_repository.test", "ingest_size_in_gb"),
				resource.TestCheckNoResourceAttr("humio_repository.test", "storage_size_in_gb"),
			),
		},
	}, testAccCheckRepositoryDestroy)
}

func TestAccRepositoryFull(t *testing.T) {
	accTestCase(t, []resource.TestStep{
		{
			Config: repositoryFull,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("humio_repository.test", "name", "repository-test"),
				resource.TestCheckResourceAttr("humio_repository.test", "description", "some description"),
				resource.TestCheckResourceAttr("humio_repository.test", "allow_data_deletion", "true"),
				resource.TestCheckResourceAttr("humio_repository.test", "retention.#", "1"),
				resource.TestCheckResourceAttr("humio_repository.test", "retention.0.time_in_days", "30"),
				resource.TestCheckResourceAttr("humio_repository.test", "retention.0.ingest_size_in_gb", "10"),
				resource.TestCheckResourceAttr("humio_repository.test", "retention.0.storage_size_in_gb", "5"),
			),
		},
	}, testAccCheckRepositoryDestroy)
}

func testAccCheckRepositoryDestroy(s *terraform.State) error {
	conn := testAccProviders["humio"].Meta().(*humio.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "humio_repository" {
			continue
		}
		// TODO: Use rs.Primary.ID to figure out if repository exists, and not just list all repositories.
		resp, err := conn.Repositories().List()
		if err == nil {
			if len(resp) > 4 { // only consider repositories not built in by default
				return fmt.Errorf("repositories still exist: %#+v", resp)
			}
		}
	}
	return nil
}

const repositoryEmpty = `
resource "humio_repository" "test" {}
`

const repositoryInvalidInputs = `
resource "humio_repository" "test" {
    name                = ["invalid"]
    description         = ["invalid"]
    allow_data_deletion = ["invalid"]
    retention           = "invalid"
}
`

const repositoryInvalidRetentionSettings = `
resource "humio_repository" "test" {
    name = "repository-invalid-retention"
    retention {
        storage_size_in_gb = -5
        ingest_size_in_gb  = -10
        time_in_days       = -30
    }
}
`

const repositoryBasic = `
resource "humio_repository" "test" {
    name = "repository-test"
    retention {}
}
`

const repositoryFull = `
resource "humio_repository" "test" {
    name                = "repository-test"
    description         = "some description"
    allow_data_deletion = true
    retention {
        storage_size_in_gb = 5
        ingest_size_in_gb  = 10
        time_in_days       = 30
    }
}
`

var wantRepository = humio.Repository{
	Name:          "test-repository",
	Description:   "important",
	RetentionDays: 30,
}

func TestEncodeDecodeRepositoryResource(t *testing.T) {
	res := resourceRepository()
	data := res.TestResourceData()
	resourceDataFromRepository(&wantRepository, data)
	got, err := repositoryFromResourceData(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(wantRepository, got) {
		t.Error(cmp.Diff(wantRepository, got))
	}
}
