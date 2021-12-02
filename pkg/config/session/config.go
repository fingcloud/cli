package session

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/thoas/go-funk"
)

// Config stores all logged in sessions
type Config []Session

// Session stores user credentials
type Session struct {
	Token   string `json:"token"`
	Email   string `json:"email"`
	Default bool   `json:"default"`
	Alias   string `json:"alias"`
}

// readConfig loads sessions from auth.json
func readConfig() (Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	cfg := Config{}

	configPath := path.Join(home, ".fing")
	sessPath := path.Join(configPath, "auth.json")

	bs, err := os.ReadFile(sessPath)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(bs, &cfg)
	if err != nil {
		writeConfig(cfg)
		return nil, err
	}

	return cfg, nil
}

// writeConfig writes session data to auth.json
func writeConfig(cfg Config) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	bs, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	configPath := path.Join(home, ".fing")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.MkdirAll(configPath, os.ModePerm)
	}
	sessPath := path.Join(configPath, "auth.json")

	return os.WriteFile(sessPath, bs, 0644)
}

// AddSession appends an additional session to auth.json
func AddSession(sess Session) error {
	cfg, err := readConfig()
	if err != nil {
		return err
	}

	for _, s := range cfg {
		s.Default = false
	}
	sess.Default = true

	var exists bool
	for i, s := range cfg {
		if s.Email == sess.Email {
			cfg[i] = sess
			exists = true
			break
		}
	}

	if !exists {
		cfg = append(cfg, sess)
	}

	return writeConfig(cfg)
}

func CurrentSession() (Session, error) {
	cfg, err := readConfig()
	if err != nil {
		return Session{}, err
	}

	found := funk.Find(cfg, func(sess Session) bool { return sess.Default })
	if sess, ok := found.(Session); ok {
		return sess, nil
	}

	return Session{}, errors.New("no default session")
}

func UseSession(email string) error {
	cfg, err := readConfig()
	if err != nil {
		return err
	}

	for _, sess := range cfg {
		if sess.Email == email || sess.Alias == email {
			sess.Default = true
			break
		}
	}

	return writeConfig(cfg)
}

func RemoveSession() (Session, error) {
	cfg, err := readConfig()
	if err != nil {
		return Session{}, err
	}

	for i, sess := range cfg {
		if sess.Default {
			cfg = append(cfg[:i], cfg[i+1:]...)
			return sess, writeConfig(cfg)
		}
	}

	return Session{}, errors.New("No default session found")
}

func Sessions() ([]Session, error) {
	return readConfig()
}
