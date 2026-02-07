package pocket

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// HomeDir returns the home directory of the current user.
// On Unix, it looks for HOME, defaults to the current user's home directory.
// On Windows, it checks USERPROFILE first, then HOME, then HOMEDRIVE+HOMEPATH.
func HomeDir() (string, error) {
	if runtime.GOOS == "windows" {
		return homeDirWindows()
	}
	return homeDirUnix()
}

// ConfigDir returns the configuration directory of the current user.
// On Unix, it looks for XDG_CONFIG_HOME, defaults to $HOME/.config.
// On Windows, it checks APPDATA.
func ConfigDir() (string, error) {
	if runtime.GOOS == "windows" {
		return configDirWindows()
	}
	return configDirUnix()
}

// DataDir returns the data directory of the current user.
// On Unix, it looks for XDG_DATA_HOME, defaults to $HOME/.local/share.
// On Windows, it checks LOCALAPPDATA.
func DataDir() (string, error) {
	if runtime.GOOS == "windows" {
		return dataDirWindows()
	}
	return dataDirUnix()
}

func homeDirWindows() (string, error) {
	if home := os.Getenv("USERPROFILE"); home != "" {
		return home, nil
	}

	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, or USERPROFILE are blank")
	}

	return home, nil
}

func homeDirUnix() (string, error) {
	homeEnv := "HOME"
	if runtime.GOOS == "plan9" {
		// On plan9, env vars are lowercase.
		homeEnv = "home"
	}

	// First prefer the HOME environmental variable
	if home := os.Getenv(homeEnv); home != "" {
		return home, nil
	}

	var stdout bytes.Buffer

	// If that fails, try OS specific commands
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("sh", "-c", `dscl -q . -read /Users/"$(whoami)" NFSHomeDirectory | sed 's/^[^ ]*: //'`)
		cmd.Stdout = &stdout
		if err := cmd.Run(); err == nil {
			result := strings.TrimSpace(stdout.String())
			if result != "" {
				return result, nil
			}
		}
	} else {
		cmd := exec.Command("getent", "passwd", strconv.Itoa(os.Getuid()))
		cmd.Stdout = &stdout
		if err := cmd.Run(); err != nil {
			// If the error is ErrNotFound, we ignore it. Otherwise, return it.
			if err != exec.ErrNotFound {
				return "", err
			}
		} else {
			if passwd := strings.TrimSpace(stdout.String()); passwd != "" {
				// username:password:uid:gid:gecos:home:shell
				passwdParts := strings.SplitN(passwd, ":", 7)
				if len(passwdParts) > 5 {
					return passwdParts[5], nil
				}
			}
		}
	}

	// If all else fails, try the shell
	stdout.Reset()
	cmd := exec.Command("sh", "-c", "cd && pwd")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func configDirUnix() (string, error) {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		return xdgConfigHome, nil
	}

	home, err := homeDirUnix()
	if err != nil {
		return "", err
	}

	configDir := home + "/.config"
	if _, err := os.Stat(configDir); err != nil {
		return "", errors.New("$HOME/.config directory does not exist")
	}

	return configDir, nil
}

func configDirWindows() (string, error) {
	// First prefer the APPDATA environmental variable
	appData := os.Getenv("APPDATA")
	if appData != "" {
		return appData, nil
	}

	return "", errors.New("APPDATA is blank")
}

func dataDirUnix() (string, error) {
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome != "" {
		return xdgDataHome, nil
	}

	home, err := homeDirUnix()
	if err != nil {
		return "", err
	}

	dataDir := home + "/.local/share"
	if _, err := os.Stat(dataDir); err != nil {
		return "", errors.New("$HOME/.local/share directory does not exist")
	}

	return dataDir, nil
}

func dataDirWindows() (string, error) {
	// First prefer the LOCALAPPDATA environment variable
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData != "" {
		return localAppData, nil
	}

	// Fall back to APPDATA if LOCALAPPDATA is not set
	appData := os.Getenv("APPDATA")
	if appData != "" {
		return appData, nil
	}

	return "", errors.New("LOCALAPPDATA and APPDATA are blank")
}
