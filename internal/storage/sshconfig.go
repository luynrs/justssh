package storage

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/luynrs/justssh/internal/models"
)

// only handles a flat ~/.ssh/config, skips wildcard match
func ImportSSHConfig() ([]models.Server, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filepath.Join(home, ".ssh", "config"))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var servers []models.Server
	var cur *models.Server

	flush := func() {
		if cur != nil && cur.Host != "" {
			servers = append(servers, *cur)
		}
		cur = nil
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		key, val, ok := splitDirective(scanner.Text())
		if !ok {
			continue
		}

		switch strings.ToLower(key) {
		case "host":
			flush()
			if isPattern(val) {
				continue
			}
			cur = &models.Server{Name: val, Host: val, Port: 22}
		case "match":
			flush()
		case "hostname":
			if cur != nil {
				cur.Host = val
			}
		case "user":
			if cur != nil {
				cur.User = val
			}
		case "port":
			if cur != nil {
				if p, err := strconv.Atoi(val); err == nil {
					cur.Port = p
				}
			}
		case "identityfile":
			if cur != nil {
				cur.Key = val
			}
		}
	}
	flush()

	return servers, scanner.Err()
}

func splitDirective(line string) (key, val string, ok bool) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false
	}

	fields := strings.Fields(line)
	if len(fields) < 2 {
		return "", "", false
	}
	return fields[0], strings.Join(fields[1:], " "), true
}

func isPattern(host string) bool {
	return strings.ContainsAny(host, "*?!") || strings.Contains(host, " ")
}
