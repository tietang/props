package helper

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// DisableCache will disable caching of the home directory. Caching is enabled
// by default.
var DisableCache bool

var homedirCache string
var cacheLock sync.RWMutex

// Dir returns the home directory for the executing user.
//
// This uses an OS-specific method for discovering the home directory.
// An error is returned if a home directory cannot be detected.
func dir() (string, error) {
	if !DisableCache {
		cacheLock.RLock()
		cached := homedirCache
		cacheLock.RUnlock()
		if cached != "" {
			return cached, nil
		}
	}

	cacheLock.Lock()
	defer cacheLock.Unlock()

	var result string
	var err error
	if runtime.GOOS == "windows" {
		result, err = homeWindows()
	} else {
		// Unix-like system, so just assume Unix
		result, err = homeUnix()
	}

	if err != nil {
		return "", err
	}
	homedirCache = result
	return result, nil
}

// Expand expands the path to include the home directory if the path
// is prefixed with `~`. If it isn't prefixed with `~`, the path is
// returned as-is.
func Expand(path string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}

	if path[0] != '~' {
		return path, nil
	}

	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return "", errors.New("cannot expand user-specific home dir")
	}

	dir, err := dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, path[1:]), nil
}

// Reset clears the cache, forcing the next call to Dir to re-detect
// the home directory. This generally never has to be called, but can be
// useful in tests if you're modifying the home directory via the HOME
// env var or something.
func Reset() {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	homedirCache = ""
}

func homeUnix() (string, error) {
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

func homeWindows() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

func UserHomeDir() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func UserConfigDir() (string, error) {
	var dir string
	homeDir, _ := UserHomeDir()
	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("AppData")
		if dir == "" {
			dir = filepath.Join(homeDir, "AppData")
		}
		if dir == "" {
			return "", errors.New("%AppData% is not defined")
		}
	case "darwin", "ios":
		home := os.Getenv("HOME")
		if home == "" {
			home = homeDir
		}
		if home == "" {
			return "", errors.New("$HOME is not defined")
		}
		dir = filepath.Join(home, "/Library/Application Support")

	case "plan9":
		home := os.Getenv("home")
		if home == "" {
			home = homeDir
		}
		if home == "" {
			return "", errors.New("$home is not defined")
		}
		dir = filepath.Join(home, "/lib")
	default: // Unix
		home := os.Getenv("XDG_CONFIG_HOME")
		if home == "" {
			home = os.Getenv("HOME")
			if home == "" {
				home = homeDir
				return "", errors.New("neither $XDG_CONFIG_HOME nor $HOME are defined")
			}
			if home == "" {
				return "", errors.New("neither $XDG_CONFIG_HOME nor $HOME are defined")
			}
			dir = filepath.Join(home, "/.config")
		}
	}

	return dir, nil
}
