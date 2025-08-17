package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

func printYAML(key string, value interface{}, indent int) {
	pad := ""
	for i := 0; i < indent; i++ {
		pad += "  "
	}
	switch v := value.(type) {
	case map[string]interface{}:
		fmt.Printf("%s%v:\n", pad, key)
		for k, val := range v {
			printYAML(k, val, indent+1)
		}
	case map[string]string:
		fmt.Printf("%s%v:\n", pad, key)
		for k, val := range v {
			fmt.Printf("%s  %s: %s\n", pad, k, val)
		}
	case []interface{}:
		fmt.Printf("%s%v:\n", pad, key)
		for i, item := range v {
			printYAML(fmt.Sprintf("[%d]", i), item, indent+1)
		}
	default:
		fmt.Printf("%s%v: %v\n", pad, key, v)
	}
}

func loadConfig(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config map[string]interface{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return config, nil
}

func loadThemes(path string) (map[string]map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading themes file: %w", err)
	}

	var themes map[string]map[string]string
	err = yaml.Unmarshal(data, &themes)
	if err != nil {
		return nil, fmt.Errorf("error parsing themes file: %w", err)
	}

	return themes, nil
}

func main() {
	config, err := loadConfig("src/config/config.yml")
	if err != nil {
		log.Fatalf("Error reading config.yml: %v", err)
	}

	themeName, ok := config["theme"].(string)
	if !ok {
		log.Fatalf("Theme not found or not a string in config.yml")
	}

	themes, err := loadThemes("src/themes/themes.yml")
	if err != nil {
		log.Fatalf("Error reading themes.yml: %v", err)
	}

	themeColors, ok := themes[themeName]
	if !ok {
		log.Fatalf("Theme '%s' not found in themes.yml", themeName)
	}

	config["theme"] = themeColors

	tmpl := template.Must(template.ParseFiles("src/templates/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
