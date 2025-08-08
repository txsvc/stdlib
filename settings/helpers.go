package settings

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	indentChar             = "  "
	filePerm   fs.FileMode = 0644
)

func ReadDialSettings(path string) (*DialSettings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ds := DialSettings{}
	if err := json.Unmarshal([]byte(data), &ds); err != nil {
		return nil, err
	}
	return &ds, nil
}

func WriteDialSettings(ds *DialSettings, path string) error {
	buf, err := json.MarshalIndent(ds, "", indentChar)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}
	}

	return os.WriteFile(path, buf, filePerm)
}

func ReadCredentials(path string) (*Credentials, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cred := Credentials{}
	if err := json.Unmarshal([]byte(data), &cred); err != nil {
		return nil, err
	}
	return &cred, nil
}

func WriteCredentials(cred *Credentials, path string) error {
	buf, err := json.MarshalIndent(cred, "", indentChar)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}
	}

	return os.WriteFile(path, buf, filePerm)
}
