/*
Copyright © 2025 Victoria Cheng
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"html/template"
	"os"
	"profilehub/utils"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the ProfileHub static files",
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.MkdirAll("dist/static", 0755); err != nil {
			fmt.Println("Error creating dist/static:", err)
			return
		}

		config, err := utils.LoadConfig("config/config.yml")
		if err != nil {
			fmt.Println("Error loading config.yml:", err)
			return
		}

		tmpl, err := template.ParseFiles("src/index.html")
		if err == nil {
			if layoutTmpl, err := tmpl.ParseGlob("src/layout/*.html"); err == nil {
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

		if err := tmpl.Execute(out, config); err != nil {
			fmt.Println("Error rendering template:", err)
			return
		}

		err = utils.CopyDir("src/static", "dist/static")
		if err != nil {
			fmt.Println("Error copying static files:", err)
		}

		fmt.Println("Build complete. Files are in dist/")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
