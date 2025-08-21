/*
Copyright © 2025 Victoria Cheng
*/
package cmd

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"profilehub/utils"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Run the development server",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := utils.LoadConfig("config/config.yml")
		if err != nil {
			log.Fatalf("Error reading config.yml: %v", err)
		}

		tmpl := template.Must(template.ParseFiles("src/index.html"))
		tmpl = template.Must(tmpl.ParseGlob("src/layout/*.html"))

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
