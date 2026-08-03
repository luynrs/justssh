package storage

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/luynrs/justssh/internal/models"
)

type config struct {
	Servers []models.Server `yaml:"servers"`
}

type Store struct {
	path string
}

func NewStore() (*Store, error) {
	path := os.Getenv("JUSTSSH_CONFIG")
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(home, ".config", "justssh", "config.yaml")
	}
	return &Store{path: path}, nil
}

func (s *Store) Load() ([]models.Server, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		// first run, nothing saved, trying to pull hosts from ~/.ssh/config
		servers, importErr := ImportSSHConfig()
		if importErr != nil {
			servers = nil
		}
		return servers, s.Save(servers)
	}
	if err != nil {
		return nil, err
	}

	var cfg config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return cfg.Servers, nil
}

func (s *Store) Save(servers []models.Server) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	data, err := yaml.Marshal(config{Servers: servers})
	if err != nil {
		return err
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
