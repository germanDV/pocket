package pocket

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestHomeDir(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "HOME env var set",
			envVars: map[string]string{
				"HOME": "/home/testuser",
			},
			wantErr: false,
		},
		{
			name:    "HOME not set (has fallback)",
			envVars: map[string]string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear relevant env vars
			os.Unsetenv("HOME")
			os.Unsetenv("USERPROFILE")
			os.Unsetenv("HOMEDRIVE")
			os.Unsetenv("HOMEPATH")

			// Set test env vars
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			got, err := HomeDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("HomeDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Error("HomeDir() returned empty string without error")
			}
		})
	}
}

func TestHomeDirWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test")
	}

	tests := []struct {
		name     string
		envVars  map[string]string
		expected string
		wantErr  bool
	}{
		{
			name: "HOME takes priority",
			envVars: map[string]string{
				"HOME":        "/home/test",
				"USERPROFILE": "C:\\Users\\test",
			},
			expected: "/home/test",
			wantErr:  false,
		},
		{
			name: "USERPROFILE fallback",
			envVars: map[string]string{
				"USERPROFILE": "C:\\Users\\test",
			},
			expected: "C:\\Users\\test",
			wantErr:  false,
		},
		{
			name: "HOMEDRIVE and HOMEPATH fallback",
			envVars: map[string]string{
				"HOMEDRIVE": "C:",
				"HOMEPATH":  "\\Users\\test",
			},
			expected: "C:\\Users\\test",
			wantErr:  false,
		},
		{
			name:     "No env vars set",
			envVars:  map[string]string{},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env vars
			os.Unsetenv("HOME")
			os.Unsetenv("USERPROFILE")
			os.Unsetenv("HOMEDRIVE")
			os.Unsetenv("HOMEPATH")

			// Set test env vars
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			got, err := homeDirWindows()
			if (err != nil) != tt.wantErr {
				t.Errorf("homeDirWindows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("homeDirWindows() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHomeDirUnix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping Unix-specific test on Windows")
	}

	t.Run("HOME env var set", func(t *testing.T) {
		t.Setenv("HOME", "/home/testuser")

		got, err := homeDirUnix()
		if err != nil {
			t.Errorf("homeDirUnix() error = %v", err)
			return
		}
		if got != "/home/testuser" {
			t.Errorf("homeDirUnix() = %v, want /home/testuser", got)
		}
	})
}

func TestConfigDir(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "XDG_CONFIG_HOME set",
			envVars: map[string]string{
				"XDG_CONFIG_HOME": "/custom/config",
			},
			wantErr: false,
		},
		{
			name:    "XDG_CONFIG_HOME not set",
			envVars: map[string]string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env vars
			os.Unsetenv("XDG_CONFIG_HOME")
			os.Unsetenv("APPDATA")

			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			got, err := ConfigDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Error("ConfigDir() returned empty string without error")
			}
		})
	}
}

func TestConfigDirUnix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping Unix-specific test on Windows")
	}

	t.Run("with existing .config directory", func(t *testing.T) {
		// Get actual home directory
		home, err := os.UserHomeDir()
		if err != nil {
			t.Skip("Cannot get user home directory:", err)
		}

		configDir := filepath.Join(home, ".config")

		// Check if .config exists, skip if it doesn't
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			t.Skip("$HOME/.config does not exist")
		}

		got, err := configDirUnix()
		if err != nil {
			t.Errorf("configDirUnix() error = %v", err)
			return
		}
		if got != configDir {
			t.Errorf("configDirUnix() = %v, want %v", got, configDir)
		}
	})
}

func TestConfigDirWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test")
	}

	tests := []struct {
		name     string
		envVars  map[string]string
		expected string
		wantErr  bool
	}{
		{
			name: "APPDATA set",
			envVars: map[string]string{
				"APPDATA": "C:\\Users\\test\\AppData\\Roaming",
			},
			expected: "C:\\Users\\test\\AppData\\Roaming",
			wantErr:  false,
		},
		{
			name:     "APPDATA not set",
			envVars:  map[string]string{},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("APPDATA")

			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			got, err := configDirWindows()
			if (err != nil) != tt.wantErr {
				t.Errorf("configDirWindows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("configDirWindows() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDataDir(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "XDG_DATA_HOME set",
			envVars: map[string]string{
				"XDG_DATA_HOME": "/custom/data",
			},
			wantErr: false,
		},
		{
			name:    "XDG_DATA_HOME not set",
			envVars: map[string]string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env vars
			os.Unsetenv("XDG_DATA_HOME")
			os.Unsetenv("LOCALAPPDATA")
			os.Unsetenv("APPDATA")

			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			got, err := DataDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("DataDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Error("DataDir() returned empty string without error")
			}
		})
	}
}

func TestDataDirUnix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping Unix-specific test on Windows")
	}

	t.Run("with existing .local/share directory", func(t *testing.T) {
		// Get actual home directory
		home, err := os.UserHomeDir()
		if err != nil {
			t.Skip("Cannot get user home directory:", err)
		}

		dataDir := filepath.Join(home, ".local", "share")

		// Check if .local/share exists, skip if it doesn't
		if _, err := os.Stat(dataDir); os.IsNotExist(err) {
			t.Skip("$HOME/.local/share does not exist")
		}

		got, err := dataDirUnix()
		if err != nil {
			t.Errorf("dataDirUnix() error = %v", err)
			return
		}
		if got != dataDir {
			t.Errorf("dataDirUnix() = %v, want %v", got, dataDir)
		}
	})
}

func TestDataDirWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test")
	}

	tests := []struct {
		name     string
		envVars  map[string]string
		expected string
		wantErr  bool
	}{
		{
			name: "LOCALAPPDATA set",
			envVars: map[string]string{
				"LOCALAPPDATA": "C:\\Users\\test\\AppData\\Local",
			},
			expected: "C:\\Users\\test\\AppData\\Local",
			wantErr:  false,
		},
		{
			name: "LOCALAPPDATA not set, APPDATA fallback",
			envVars: map[string]string{
				"APPDATA": "C:\\Users\\test\\AppData\\Roaming",
			},
			expected: "C:\\Users\\test\\AppData\\Roaming",
			wantErr:  false,
		},
		{
			name:     "Neither set",
			envVars:  map[string]string{},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("LOCALAPPDATA")
			os.Unsetenv("APPDATA")

			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			got, err := dataDirWindows()
			if (err != nil) != tt.wantErr {
				t.Errorf("dataDirWindows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("dataDirWindows() = %v, want %v", got, tt.expected)
			}
		})
	}
}
