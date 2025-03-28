// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package test provides e2e tests for Zarf.
package test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"testing"

	"github.com/defenseunicorns/zarf/src/pkg/utils"
	"github.com/defenseunicorns/zarf/src/pkg/utils/exec"
	dconfig "github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/stretchr/testify/require"
)

// ZarfE2ETest Struct holding common fields most of the tests will utilize.
type ZarfE2ETest struct {
	ZarfBinPath     string
	Arch            string
	ApplianceMode   bool
	RunClusterTests bool
}

var logRegex = regexp.MustCompile(`Saving log file to (?P<logFile>.*?\.log)`)

// GetCLIName looks at the OS and CPU architecture to determine which Zarf binary needs to be run.
func GetCLIName() string {
	var binaryName string
	if runtime.GOOS == "linux" {
		binaryName = "zarf"
	} else if runtime.GOOS == "darwin" {
		if runtime.GOARCH == "arm64" {
			binaryName = "zarf-mac-apple"
		} else {
			binaryName = "zarf-mac-intel"
		}
	} else if runtime.GOOS == "windows" {
		if runtime.GOARCH == "amd64" {
			binaryName = "zarf.exe"
		}
	}
	return binaryName
}

// SetupWithCluster performs actions for each test that requires a K8s cluster.
func (e2e *ZarfE2ETest) SetupWithCluster(t *testing.T) {
	if !e2e.RunClusterTests {
		t.Skip("")
	}
	_ = exec.CmdWithPrint("sh", "-c", fmt.Sprintf("%s tools kubectl describe nodes | grep -A 99 Non-terminated", e2e.ZarfBinPath))
}

// Zarf executes a Zarf command.
func (e2e *ZarfE2ETest) Zarf(args ...string) (string, string, error) {
	return exec.CmdWithContext(context.TODO(), exec.PrintCfg(), e2e.ZarfBinPath, args...)
}

// Kubectl executes `zarf tools kubectl ...`
func (e2e *ZarfE2ETest) Kubectl(args ...string) (string, string, error) {
	tk := []string{"tools", "kubectl"}
	args = append(tk, args...)
	return e2e.Zarf(args...)
}

// CleanFiles removes files and directories that have been created during the test.
func (e2e *ZarfE2ETest) CleanFiles(files ...string) {
	for _, file := range files {
		_ = os.RemoveAll(file)
	}
}

// GetMismatchedArch determines what architecture our tests are running on,
// and returns the opposite architecture.
func (e2e *ZarfE2ETest) GetMismatchedArch() string {
	switch e2e.Arch {
	case "arm64":
		return "amd64"
	default:
		return "arm64"
	}
}

// GetLogFileContents gets the log file contents from a given run's std error.
func (e2e *ZarfE2ETest) GetLogFileContents(t *testing.T, stdErr string) string {
	get, err := utils.MatchRegex(logRegex, stdErr)
	require.NoError(t, err)
	logFile := get("logFile")
	logContents, err := os.ReadFile(logFile)
	require.NoError(t, err)
	return string(logContents)
}

// SetupDockerRegistry uses the host machine's docker daemon to spin up a local registry for testing purposes.
func (e2e *ZarfE2ETest) SetupDockerRegistry(t *testing.T, port int) *configfile.ConfigFile {
	// spin up a local registry
	registryImage := "registry:2.8.2"
	err := exec.CmdWithPrint("docker", "run", "-d", "--restart=always", "-p", fmt.Sprintf("%d:5000", port), "--name", "registry", registryImage)
	require.NoError(t, err)

	// docker config folder
	cfg, err := dconfig.Load(dconfig.Dir())
	require.NoError(t, err)
	if !cfg.ContainsAuth() {
		// make a docker config file w/ some blank creds
		_, _, err := e2e.Zarf("tools", "registry", "login", "--username", "zarf", "-p", "zarf", "localhost:6000")
		require.NoError(t, err)
	}

	return cfg
}
