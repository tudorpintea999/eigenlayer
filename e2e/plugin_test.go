package e2e

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/NethermindEth/eigenlayer/e2e/docker"
	"github.com/NethermindEth/eigenlayer/internal/data"
	"github.com/NethermindEth/eigenlayer/internal/package_handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestPlugin_Install_RemoteContext(t *testing.T) {
	// Test context
	var (
		testDir    = t.TempDir()
		pkgDir     = filepath.Join(testDir, "mock-avs")
		installErr error
	)
	e2eTest := newE2ETestCase(
		t,
		// Arrange
		func(t *testing.T, egnPath string) error {
			err := os.MkdirAll(pkgDir, 0o755)
			if err != nil {
				return err
			}
			err = runCommand(t, "git", "clone", "--single-branch", "-b", latestMockAVSVersion, mockAVSRepo, pkgDir)
			if err != nil {
				return err
			}
			// remove .git folder
			return os.RemoveAll(filepath.Join(pkgDir, ".git"))
		},
		// Act
		func(t *testing.T, egnPath string) {
			installErr = runCommand(t, egnPath, "local-install", pkgDir, "--profile", "option-returner", "--run")
		},
		// Assert
		func(t *testing.T) {
			assert.NoError(t, installErr, "local-install command should succeed")
			instance := getInstance(t, "mock-avs-default")
			assert.NotNil(t, instance.Plugin, "plugin should be installed")
			assert.Equal(t, instance.Plugin.Type, data.PluginTypeRemoteContext, "plugin type should be ", data.PluginTypeRemoteContext)
			assert.Equal(t, instance.Plugin.Src, "https://github.com/NethermindEth/mock-avs.git#main:plugin")
		},
	)
	e2eTest.run()
}

func TestPlugin_Install_LocalContext(t *testing.T) {
	// Test context
	var (
		testDir    = t.TempDir()
		pkgDir     = filepath.Join(testDir, "mock-avs")
		installErr error
		dataDir    string
	)
	e2eTest := newE2ETestCase(
		t,
		// Arrange
		func(t *testing.T, egnPath string) error {
			var err error
			dataDir, err = dataDirPath()
			if err != nil {
				return err
			}
			if err := os.MkdirAll(pkgDir, 0o755); err != nil {
				return err
			}
			err = runCommand(t, "git", "clone", "--single-branch", "-b", latestMockAVSVersion, mockAVSRepo, pkgDir)
			if err != nil {
				return err
			}
			// remove .git folder
			if err := os.RemoveAll(filepath.Join(pkgDir, ".git")); err != nil {
				return err
			}
			// modify manifest to use local context
			manifestFile, err := os.OpenFile(filepath.Join(pkgDir, "pkg", "manifest.yml"), os.O_RDWR, 0o755)
			if err != nil {
				return err
			}
			manifestData, err := io.ReadAll(manifestFile)
			if err != nil {
				return err
			}
			if err = manifestFile.Close(); err != nil {
				return err
			}
			var manifest package_handler.Manifest
			if err = yaml.Unmarshal(manifestData, &manifest); err != nil {
				return err
			}
			manifest.Plugin.BuildFrom = "../plugin"
			manifestData, err = yaml.Marshal(manifest)
			if err != nil {
				return err
			}
			manifestFile, err = os.OpenFile(filepath.Join(pkgDir, "pkg", "manifest.yml"), os.O_RDWR|os.O_TRUNC, 0o755)
			if err != nil {
				return err
			}
			if _, err = manifestFile.Write(manifestData); err != nil {
				return err
			}
			return manifestFile.Close()
		},
		// Act
		func(t *testing.T, egnPath string) {
			installErr = runCommand(t, egnPath, "local-install", pkgDir, "--profile", "option-returner", "--run")
		},
		// Assert
		func(t *testing.T) {
			assert.NoError(t, installErr, "local-install command should succeed")
			instance := getInstance(t, "mock-avs-default")
			assert.NotNil(t, instance.Plugin, "plugin should be installed")
			assert.Equal(t, instance.Plugin.Type, data.PluginTypeLocalContext, "plugin type should be ", data.PluginTypeLocalContext)
			assert.Equal(t, instance.Plugin.Src, "mock-avs-default")
			assert.FileExists(t, filepath.Join(dataDir, "plugin", "mock-avs-default.tar"))
		},
	)
	e2eTest.run()
}

func TestPlugin_Install_RemoteImage(t *testing.T) {
	// Test context
	var (
		testDir    = t.TempDir()
		pkgDir     = filepath.Join(testDir, "mock-avs")
		installErr error
	)
	e2eTest := newE2ETestCase(
		t,
		// Arrange
		func(t *testing.T, egnPath string) error {
			if err := os.MkdirAll(pkgDir, 0o755); err != nil {
				return err
			}
			err := runCommand(t, "git", "clone", "--single-branch", "-b", latestMockAVSVersion, mockAVSRepo, pkgDir)
			if err != nil {
				return err
			}
			// remove .git folder
			if err := os.RemoveAll(filepath.Join(pkgDir, ".git")); err != nil {
				return err
			}
			// modify manifest to use local context
			manifestFile, err := os.OpenFile(filepath.Join(pkgDir, "pkg", "manifest.yml"), os.O_RDWR, 0o755)
			if err != nil {
				return err
			}
			manifestData, err := io.ReadAll(manifestFile)
			if err != nil {
				return err
			}
			if err = manifestFile.Close(); err != nil {
				return err
			}
			var manifest package_handler.Manifest
			if err = yaml.Unmarshal(manifestData, &manifest); err != nil {
				return err
			}
			manifest.Plugin.BuildFrom = ""
			manifest.Plugin.Image = "busybox:1.36"
			manifestData, err = yaml.Marshal(manifest)
			if err != nil {
				return err
			}
			manifestFile, err = os.OpenFile(filepath.Join(pkgDir, "pkg", "manifest.yml"), os.O_RDWR|os.O_TRUNC, 0o755)
			if err != nil {
				return err
			}
			if _, err = manifestFile.Write(manifestData); err != nil {
				return err
			}
			return manifestFile.Close()
		},
		// Act
		func(t *testing.T, egnPath string) {
			installErr = runCommand(t, egnPath, "local-install", pkgDir, "--profile", "option-returner", "--run")
		},
		// Assert
		func(t *testing.T) {
			assert.NoError(t, installErr, "local-install command should succeed")
			instance := getInstance(t, "mock-avs-default")
			assert.NotNil(t, instance.Plugin, "plugin should be installed")
			assert.Equal(t, instance.Plugin.Type, data.PluginTypeRemoteImage, "plugin type should be ", data.PluginTypeRemoteImage)
			assert.Equal(t, instance.Plugin.Src, "busybox:1.36")
		},
	)
	e2eTest.run()
}

func TestPlugin_Run_RemoteContext(t *testing.T) {
	// Test context
	var (
		testDir      = t.TempDir()
		pkgDir       = filepath.Join(testDir, "mock-avs")
		installErr   error
		runPluginErr error
		eventsSince  time.Time
		eventsUntil  time.Time
	)
	e2eTest := newE2ETestCase(
		t,
		// Arrange
		func(t *testing.T, egnPath string) error {
			err := os.MkdirAll(pkgDir, 0o755)
			if err != nil {
				return err
			}
			err = runCommand(t, "git", "clone", "--single-branch", "-b", latestMockAVSVersion, mockAVSRepo, pkgDir)
			if err != nil {
				return err
			}
			// remove .git folder
			if err := os.RemoveAll(filepath.Join(pkgDir, ".git")); err != nil {
				return err
			}
			return runCommand(t, egnPath, "local-install", pkgDir, "--profile", "option-returner", "--run")
		},
		// Act
		func(t *testing.T, egnPath string) {
			eventsSince = time.Now()
			runPluginErr = runCommand(t, egnPath, "plugin", "mock-avs-default")
			eventsUntil = time.Now()
		},
		// Assert
		func(t *testing.T) {
			assert.NoError(t, installErr, "local-install command should succeed")
			assert.NoError(t, runPluginErr, "plugin command should succeed")

			// Check docker events
			pluginContainerID := ""
			networkID, err := getNetworkIDByName("eigenlayer")
			require.NoError(t, err, "getNetworkIDByName should succeed")
			events, err := docker.EventsRange(context.Background(), eventsSince, eventsUntil)
			require.NoError(t, err, "docker events should succeed")

			events.CheckInOrder(t,
				docker.NewImageTagged("mock-avs-default", "latest"),
				docker.NewContainerCreated("eigen-plugin-mock-avs-default", &pluginContainerID),
				docker.NewNetworkConnect(&pluginContainerID, &networkID),
				docker.NewNetworkDisconnect(&pluginContainerID, &networkID),
				docker.NewContainerDies(&pluginContainerID),
				docker.NewContainerDestroy(&pluginContainerID),
			)
		},
	)
	e2eTest.run()
}
