package vos

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sync"
)

type System interface {
	Command(name string, arg ...string) *exec.Cmd
	Chdir(dir string) error
	Getwd() (string, error)
	Environ() []string
	Getenv(string) string
	Setenv(string, string)
	Exit(int)
}

type LocalSystem struct {
	root string
	dir  string
	env  map[string]any
	mu   sync.RWMutex
}

func NewLocalSystem(root string) *LocalSystem {
	dir, _ := os.Getwd()
	if v, err := filepath.Rel(root, dir); err != nil {
		dir = ""
	} else {
		dir = v
	}
	return &LocalSystem{
		root: root,
		dir:  dir,
		env:  make(map[string]any),
	}
}

func (s *LocalSystem) Command(name string, arg ...string) *exec.Cmd {
	e := exec.Command(name, arg...)
	e.Env = s.Env()
	e.Dir = s.dir
	return e
}

func (s *LocalSystem) Chdir(dir string) error {
	abs := filepath.Join(s.root, dir)
	if info, err := os.Stat(abs); err != nil || !info.IsDir() {
		return fmt.Errorf("invalid directory: %v", err)
	}
	s.dir = dir
	return nil
}

func (s *LocalSystem) Getwd() (string, error) {
	return s.dir, nil
}

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
	// os.Exit(code)
}
