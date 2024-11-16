package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate [flags] folder",
	Short: "Validate the items in a given folder",
	Args:  cobra.ExactArgs(1),
	Run:   validate,
}

func validate(cmd *cobra.Command, args []string) {
	// Load schemas
	err := loadSchemas()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	var errors int
	// Validate files
	instanceFile := args[0]
	if err = filepath.Walk(instanceFile, validateFile); err != nil {
		errors++
	}
	slog.Info("errors count", slog.Int("errors", errors))
}

func validateFile(path string, info os.FileInfo, err error) error {
	// Skip folders
	if info.IsDir() {
		return nil
	}

	// Read and unmarshal the file
	fc, err := os.ReadFile(path)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	item := make(map[string]interface{})
	if err = yaml.Unmarshal([]byte(fc), &item); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// Validate the file with the correct validator
	switch item["kind"] {
	case "event.openevents.tech/v1alpha1":
		err = schEventV1A1.Validate(item)
	case "organizer.openevents.tech/v1alpha1":
		err = schOrganizerV1A1.Validate(item)
	default:
		err = fmt.Errorf("invalid kind")
	}

	slog.Debug("file validity",
		slog.String("path", path),
		slog.Bool("valid", err == nil),
	)

	if err != nil {
		slog.Info("invalid file", slog.Any("error", err))
	}
	return err
}
