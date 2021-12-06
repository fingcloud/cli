package session

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/thoas/go-funk"
)

// Config stores all logged in sessions
type Config []Session

// Session stores user credentials
type Session struct {
	Token      string     `json:"token"`
	Email      string     `json:"email"`
	Default    bool       `json:"default"`
	Alias      string     `json:"alias"`
	LoginAt    *time.Time `json:"login_at,omitempty"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
}

// ReadConfig loads sessions from auth.json
func ReadConfig() (Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	cfg := Config{}

	configPath := path.Join(home, ".fing")
	sessPath := path.Join(configPath, "auth.json")

	bs, err := os.ReadFile(sessPath)
	if err != nil {
		WriteConfig(cfg)
		return cfg, nil
	}

	err = json.Unmarshal(bs, &cfg)
	if err != nil {
		WriteConfig(cfg)
		return nil, err
	}

	return cfg, nil
}

// WriteConfig writes session data to auth.json
func WriteConfig(cfg Config) error {
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
	cfg, err := ReadConfig()
	if err != nil {
		return err
	}

	for _, s := range cfg {
		s.Default = false
	}
	sess.Default = true
	now := time.Now()
	sess.LoginAt = &now
	sess.LastUsedAt = &now

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

	return WriteConfig(cfg)
}

func CurrentSession() (Session, error) {
	cfg, err := ReadConfig()
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
	cfg, err := ReadConfig()
	if err != nil {
		return err
	}

	for _, sess := range cfg {
		if sess.Email == email || sess.Alias == email {
			sess.Default = true
			now := time.Now()
			sess.LastUsedAt = &now
			break
		}
	}

	return WriteConfig(cfg)
}

func RemoveSession() (Session, error) {
	cfg, err := ReadConfig()
	if err != nil {
		return Session{}, err
	}

	for i, sess := range cfg {
		if sess.Default {
			cfg = append(cfg[:i], cfg[i+1:]...)
			return sess, WriteConfig(cfg)
		}
	}

	return Session{}, errors.New("No default session found")
}

func AllSessions() ([]Session, error) {
	return ReadConfig()
}
