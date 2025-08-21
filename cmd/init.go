/*
Copyright © 2025 Victoria Cheng
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"profilehub/templates"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new ProfileHub project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		template, _ := cmd.Flags().GetString("template")
		
		// Default to "default" template if not specified
		if template == "" {
			template = "default"
		}
		
		// Validate template dynamically against embedded templates
		if _, err := templates.FS.ReadDir(template); err != nil {
			fmt.Printf("❌ Unknown template: %s\n", template)
			// List available templates from embedded FS root
			if entries, err2 := templates.FS.ReadDir("."); err2 == nil {
				fmt.Println("Available templates:")
				for _, e := range entries {
					if e.IsDir() {
						fmt.Printf("  - %s\n", e.Name())
					}
				}
			}
			return
		}
		
		fmt.Printf("Creating new ProfileHub project: %s\n", projectName)
		fmt.Printf("Using template: %s\n", template)
		
		// Create project directory
		if err := os.MkdirAll(projectName, 0755); err != nil {
			fmt.Printf("Error creating project directory: %v\n", err)
			return
		}
		
		// Extract embedded files
		if err := extractTemplate(projectName, template); err != nil {
			fmt.Printf("Error extracting template files: %v\n", err)
			return
		}
		
		fmt.Printf("✅ Project '%s' created successfully!\n", projectName)
		fmt.Printf("\nNext steps:\n")
		fmt.Printf("  cd %s\n", projectName)
		fmt.Printf("  # Edit config/config.yml with your information\n")
		fmt.Printf("  profilehub dev    # Start development server\n")
		fmt.Printf("  profilehub build  # Build static site\n")
	},
}

func extractTemplate(projectDir, templateName string) error {
	// In embedded FS, the root is the templates directory; template folders are under "default".
	templatePath := templateName
	// Ensure the template directory exists in the embedded FS
	if _, err := templates.FS.ReadDir(templatePath); err != nil {
		return fmt.Errorf("template '%s' not found in embedded assets: %w", templateName, err)
	}
	// Extract the entire template directory recursively
	return extractDir(templatePath, projectDir)
}

func extractDir(embedPath, destPath string) error {
	entries, err := templates.FS.ReadDir(embedPath)
	if err != nil {
		return err
	}
	
	for _, entry := range entries {
		srcPath := filepath.Join(embedPath, entry.Name())
		dstPath := filepath.Join(destPath, entry.Name())
		
		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := extractDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			content, err := templates.FS.ReadFile(srcPath)
			if err != nil {
				return err
			}
			
			if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
				return err
			}
			
			if err := os.WriteFile(dstPath, content, 0644); err != nil {
				return err
			}
		}
	}
	
	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
