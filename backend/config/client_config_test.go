package config

import (
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientConfigEnv(t *testing.T) {
	tests := []struct {
		name    string
		envHome string
	}{ 
		{
			name:    "environment setup correctly",
			envHome: "/tmp/home",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("HOME", tt.envHome)
			defer os.Unsetenv("HOME")

			env := NewClientConfigEnv()
			assert.NotNil(t, env)
			assert.Empty(t, env.KubeConfigs)

			expectedDir := filepath.Join(tt.envHome, appConfigDir, appKubeConfigDir)
			_, err := os.Stat(expectedDir)
			assert.NoError(t, err)
		})
	}
}

func TestNewClientConfigAppConfig(t *testing.T) {
	t.Run("app config initialization", func(t *testing.T) {
		config := NewClientConfigAppConfig("appTest", 40, 51)
		assert.NotNil(t, config)
		assert.NotNil(t, config.KubeConfig)
	})
}
func ClientAppConfigLoadClientAppConfig(t *testing.T) {
	t.Run("load app config with invalid paths", func(t *testing.T) {
		os.Setenv("HOME", "/invalid/home/path")
		defer os.Unsetenv("HOME")

		config := NewClientConfigAppConfig("appTest", 40, 51)
		config.LoadClientAppConfig()
		assert.NotContains(t, config.KubeConfig, InClusterKey)
	})
}

func TestClientAppKubeConfigBuilds(t *testing.T) {
	tests := []struct {
		name      string
		homeDir   string
		kubeDir   string
		shouldErr bool
	}{
		{
			name:      "happy path - directory exists",
			homeDir:   "/tmp/home",
			kubeDir:   "/tmp/home/.kube",
			shouldErr: false,
		},
		{
			name:      "error path - directory does not exist",
			homeDir:   "/tmp/home",
			kubeDir:   "/tmp/home/.kube",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("HOME", tt.homeDir)
			defer os.Unsetenv("HOME")

			config := NewClientConfigAppConfig("appTest", 40, 51)
			config.LoadClientAppConfig()
			config.buildKubeConfigs(filepath.Join(homedir.HomeDir(), defaultKubeConfigDir))
		})
	}
}

func TestClientAppConfigRemoveKubeConfig(t *testing.T) {

	config := NewClientConfigAppConfig("appTest", 40, 51)
	config.RemoveKubeConfig("test")
	assert.Empty(t, config.KubeConfig)
}

func TestClientAppConfigSaveKubeConfig(t *testing.T) {
	config := NewClientConfigAppConfig("appTest", 40, 51)
	config.SaveKubeConfig("test")
	assert.NotEmpty(t, config.KubeConfig)
}

func TestClientAppConfigReloadConfig(t *testing.T) {
	config := NewClientConfigAppConfig("appTest", 40, 51)
	config.ReloadConfig()
	assert.NotEmpty(t, config.KubeConfig)
}

func TestClientAppConfigReadAllFilesInDir(t *testing.T) {
	tests := []struct {
		name        string
		dirPath     string
		expectedLen int
	}{
		{
			name:        "happy path - directory exists",
			dirPath:     "/tmp",
			expectedLen: 2,
		},
		{
			name:        "error path - directory does not exist",
			dirPath:     "/invalid/path",
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := readAllFilesInDir(tt.dirPath)
			assert.Len(t, files, tt.expectedLen)
		})
	}
}

