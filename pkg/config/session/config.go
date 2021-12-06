package session

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"time"

	"github.com/mitchellh/go-homedir"
)

// Session stores user credentials
type Session struct {
	Token      string     `json:"token"`
	Email      string     `json:"email"`
	Default    bool       `json:"default"`
	Alias      string     `json:"alias"`
	LoginAt    *time.Time `json:"login_at,omitempty"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
}

// Read loads sessions from auth.json
func Read() ([]*Session, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	sessions := make([]*Session, 0)

	configPath := path.Join(home, ".fing")
	sessPath := path.Join(configPath, "auth.json")

	bs, err := os.ReadFile(sessPath)
	if err != nil {
		Write(sessions)
		return sessions, nil
	}

	err = json.Unmarshal(bs, &sessions)
	if err != nil {
		Write(sessions)
		return nil, err
	}

	return sessions, nil
}

// Write writes session data to auth.json
func Write(sessions []*Session) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	bs, err := json.Marshal(sessions)
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
func AddSession(sess *Session) error {
	sessions, err := Read()
	if err != nil {
		return err
	}

	for _, s := range sessions {
		s.Default = false
	}
	sess.Default = true
	now := time.Now()
	sess.LoginAt = &now
	sess.LastUsedAt = &now

	var exists bool
	for i, s := range sessions {
		if s.Email == sess.Email {
			sessions[i] = sess
			exists = true
			break
		}
	}

	if !exists {
		sessions = append(sessions, sess)
	}

	return Write(sessions)
}

func CurrentSession() (*Session, error) {
	sessions, err := Read()
	if err != nil {
		return nil, err
	}

	for _, sess := range sessions {
		if sess.Default {
			now := time.Now()
			sess.LastUsedAt = &now
			return sess, Write(sessions)
		}
	}

	return nil, errors.New("you need to login")
}

func UseSession(email string) error {
	sessions, err := Read()
	if err != nil {
		return err
	}

	for _, sess := range sessions {
		sess.Default = false
	}

	for _, sess := range sessions {
		if sess.Email == email || sess.Alias == email {
			sess.Default = true
			now := time.Now()
			sess.LastUsedAt = &now
			break
		}
	}

	return Write(sessions)
}

func RemoveSession() (*Session, error) {
	sessions, err := Read()
	if err != nil {
		return nil, err
	}

	for i, sess := range sessions {
		if sess.Default {
			sessions = append(sessions[:i], sessions[i+1:]...)
			if i < len(sessions) {
				sessions[i].Default = true
			} else {
				sessions[0].Default = true
			}
			return sess, Write(sessions)
		}
	}

	return nil, nil
}
