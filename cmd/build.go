/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"
)

func copyDir(src string, dst string) error {
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

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the ProfileHub static files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Building dist folder...")
		os.MkdirAll("dist/static", 0755)
		os.MkdirAll("dist/webfonts", 0755)

		// Render index.html as Go template
		// Load config.yml
		configPath := "config/config.yml"
		configFile, err := os.Open(configPath)
		if err != nil {
			fmt.Println("Error opening config.yml:", err)
			return
		}
		defer configFile.Close()
		configData := make(map[string]interface{})
		{
			// Parse YAML
			importYaml := func() error {
				decoder := yaml.NewDecoder(configFile)
				return decoder.Decode(&configData)
			}
			if err := importYaml(); err != nil {
				fmt.Println("Error parsing config.yml:", err)
				return
			}
		}

		// Parse templates (index.html + partials)
		layoutGlob := "src/layout/*.html"
		tmpl, err := template.ParseFiles("src/index.html")
		if err == nil {
			// Parse partials
			if layoutTmpl, err := tmpl.ParseGlob(layoutGlob); err == nil {
				tmpl = layoutTmpl
			}
		}
		if err != nil {
			fmt.Println("Error parsing templates:", err)
			return
		}

		out, err := os.Create("dist/index.html")
		if err != nil {
			fmt.Println("Error creating dist/index.html:", err)
			return
		}
		defer out.Close()
		if err := tmpl.Execute(out, configData); err != nil {
			fmt.Println("Error rendering template:", err)
			return
		}

		// Copy layout (templates)
		err = copyDir("src/layout", "dist/layout")
		if err != nil {
			fmt.Println("Error copying layout files:", err)
		}
		// Copy static assets
		err = copyDir("src/static", "dist/static")
		if err != nil {
			fmt.Println("Error copying static files:", err)
		}
		// Copy webfonts
		err = copyDir("src/static/webfonts", "dist/webfonts")
		if err != nil {
			fmt.Println("Error copying webfonts:", err)
		}
		fmt.Println("Build complete. Files are in dist/")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
