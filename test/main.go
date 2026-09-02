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

// Command test runs the provider acceptance tests against a throwaway Humio
// instance started in Docker.
//
// It lives in its own module so that testcontainers, and the container runtime
// client it drags in, stay out of the provider's dependency graph. Any
// arguments are passed on to the underlying "go test" invocation:
//
//	cd test && go run . -run TestAccRepository
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	humioJvmArgs = "-Xss2M"
	humioPort    = "8080/tcp"

	// providerModule is the module the acceptance tests belong to. The harness
	// looks upwards from the working directory for it, so it can be started
	// from the repository root or from this directory.
	providerModule = "github.com/clearhaus/terraform-provider-humio"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(testArgs []string) int {
	root, err := providerModuleRoot()
	if err != nil {
		log.Print(err)
		return 1
	}

	commonIdentifier := randomIdentifier()

	// Containers definitely run in the background
	ctx := context.Background()

	// Define network
	nw, err := network.New(ctx, network.WithDriver("bridge"))
	if err != nil {
		log.Print("Could not create Docker network: ", err)
		return 1
	}
	defer func() { _ = nw.Remove(ctx) }()

	log.Println("Container network created: " + nw.Name)

	// Start container(s) in order
	startedContainers, err := startContainers(ctx, humioRequest(commonIdentifier, nw.Name))

	// Stop containers after the test run, including a partially started set
	defer func() {
		log.Println("Tearing down containers")
		for _, c := range startedContainers {
			_ = c.Terminate(ctx)
		}
	}()

	if err != nil {
		log.Print(err)
		return 1
	}

	// Expose mapped Humio port
	mapped, err := startedContainers[0].MappedPort(ctx, humioPort)
	if err != nil {
		log.Print("Could not get mapped port of Humio container: ", err)
		return 1
	}

	addr := fmt.Sprintf("http://localhost:%s", mapped.Port())

	// Fixed auth credentials
	token, err := fetchDeveloperToken(addr, commonIdentifier)
	if err != nil {
		log.Printf("Could not get token for user 'developer': %v", err)
		return 1
	}

	// Run the actual tests
	log.Printf("Humio container running at %s", addr)
	return goTest(root, addr, *token, testArgs)
}

// goTest runs the provider's acceptance tests against the given Humio instance
// and reports the exit code of "go test".
func goTest(root, addr, token string, testArgs []string) int {
	args := append([]string{"test", "./humio/..."}, testArgs...)

	cmd := exec.Command("go", args...)
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"TF_ACC=1",
		"HUMIO_ADDR="+addr,
		"HUMIO_API_TOKEN="+token,
	)

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		log.Printf("Could not run acceptance tests: %v", err)
		return 1
	}

	return 0
}

// providerModuleRoot walks upwards from the working directory looking for the
// go.mod of the provider module.
func providerModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not determine working directory: %w", err)
	}

	for {
		content, err := os.ReadFile(filepath.Join(dir, "go.mod"))
		if err == nil && strings.HasPrefix(string(content), "module "+providerModule+"\n") {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find the %s module above the working directory", providerModule)
		}
		dir = parent
	}
}

func humioRequest(identifier, networkName string) testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Name:         identifier,
		Image:        "humio/humio:stable", // TODO: Support multiple versions?
		ExposedPorts: []string{humioPort},
		Env: map[string]string{
			"HUMIO_JVM_ARGS":        humioJvmArgs,
			"AUTHENTICATION_METHOD": "single-user",
			"SINGLE_USER_PASSWORD":  identifier,
		},
		Networks:   []string{networkName},
		WaitingFor: wait.ForListeningPort(humioPort),
	}
}

func startContainers(ctx context.Context, reqs ...testcontainers.ContainerRequest) ([]testcontainers.Container, error) {
	var containers []testcontainers.Container
	for _, req := range reqs {
		c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})

		if err != nil {
			return containers, fmt.Errorf("could not start container %s: %w", req.Name, err)
		}

		containers = append(containers, c)
	}

	return containers, nil
}

func fetchDeveloperToken(addr string, password string) (*string, error) {
	var token = ""

	payload := fmt.Sprintf(`{"login": "developer", "password": "%s"}`, password)
	path := fmt.Sprintf("%s%s", addr, "/api/v1/login")

	res, err := http.Post(path, "application/json", strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("could not perform developer login: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("got a non-OK status code while logging in: %d", res.StatusCode)
	}

	m := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("could not decode Humio login response: %w", err)
	}

	token = m["token"].(string)

	return &token, nil
}

func randomIdentifier() string {
	return fmt.Sprintf("tf-acc-humio-%d", rand.IntN(9999)+1)
}
