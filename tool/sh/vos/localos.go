package vos

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

type LocalSystem struct {
	root string
	dir  string
	env  map[string]any
	mu   sync.RWMutex

	// exit call back
	Exitf func(int)
}

func NewLocalSystem(root string) *LocalSystem {
	return &LocalSystem{
		root: root,
		dir:  root,
		env:  make(map[string]any),
	}
}

func (s *LocalSystem) Command(name string, arg ...string) *exec.Cmd {
	e := exec.Command(name, arg...)
	e.Env = s.Env()
	e.Dir = s.dir
	return e
}

func (s *LocalSystem) validatePath(path string) (string, error) {
	path = filepath.Clean(path)
	if !strings.HasPrefix(path, s.root) {
		return "", fmt.Errorf("invalid local fs path: %s", path)
	}
	rel := strings.TrimPrefix(path, s.root)
	abs, err := filepath.Abs(filepath.Join(s.root, rel))
	if err != nil {
		return "", fmt.Errorf("invalid path %q: %w", path, err)
	}
	return abs, nil
}

func (s *LocalSystem) Chdir(path string) error {
	abs, err := s.validatePath(path)
	if err != nil {
		return err
	}
	if info, err := os.Stat(abs); err != nil || !info.IsDir() {
		return fmt.Errorf("invalid directory: %v", err)
	}
	s.dir = abs
	return os.Chdir(abs)
}

func (s *LocalSystem) Getwd() (string, error) {
	// return os.Getwd()
	return s.dir, nil
}

// Env returns all environment variables as a name=value list.
// It converts complex and nested data structures into a JSON string
// representation when necessary.
// For basic types such as strings, integers, floats, and booleans,
// a direct conversion to a string format is applied.
func (s *LocalSystem) Env() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var env []string
	for k, v := range s.env {
		kv := fmt.Sprintf("%s=%v", k, stringify(v))
		env = append(env, kv)
	}
	return env
}

func stringify(v interface{}) string {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.String, reflect.Int, reflect.Float64, reflect.Bool:
		return fmt.Sprintf("%v", v)
	default:
		jsonStr, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("error: %+v", err)
		}
		return string(jsonStr)
	}
}

// Return all environment variables as a map.
func (s *LocalSystem) Environ() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	env := make(map[string]any)
	maps.Copy(env, s.env)
	return env
}

func (s *LocalSystem) Getenv(key string) any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.env[key]
}

func (s *LocalSystem) Setenv(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.env[key] = value
}

func (s *LocalSystem) Exit(code int) {
	if s.Exitf != nil {
		s.Exitf(code)
	}
}
