package utils

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

type Theme struct {
	Background  string `yaml:"Background"`
	Text        string `yaml:"Text"`
	Button      string `yaml:"Button"`
	ButtonText  string `yaml:"ButtonText"`
	ButtonHover string `yaml:"ButtonHover"`
}

type Social struct {
	Icon string `yaml:"Icon"`
	URL  string `yaml:"URL"`
}

type Link struct {
	Name string `yaml:"Name"`
	URL  string `yaml:"URL"`
}

type Params struct {
	Avatar   string   `yaml:"Avatar"`
	Name     string   `yaml:"Name"`
	Headline string   `yaml:"Headline"`
	Theme    Theme    `yaml:"Theme"`
	Socials  []Social `yaml:"Socials"`
	Links    []Link   `yaml:"Links"`
}

type Config struct {
	Params Params `yaml:"Params"`
}

func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing config file: %w", err)
	}

	return config, nil
}

// CopyDir recursively copies a directory from src to dst, skipping .gitkeep files
func CopyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Skip .gitkeep files
		if info.Name() == ".gitkeep" {
			return nil
		}

		destPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		return err
	})
}
