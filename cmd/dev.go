/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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

func loadConfig(path string) (Config, error) {
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

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Run the development server",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig("config/config.yml")
		if err != nil {
			log.Fatalf("Error reading config.yml: %v", err)
		}

		tmpl := template.Must(template.ParseGlob("src/layout/*.html"))

		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("src/static"))))

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			err := tmpl.Execute(w, config)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		fmt.Println("Server running at http://localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
