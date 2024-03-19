package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/jimeh/undent"
)

var logger = hclog.New(&hclog.LoggerOptions{
	Name:       "schemafetcher",
	Level:      hclog.Info,
	Output:     os.Stderr,
	Mutex:      &sync.Mutex{},
	TimeFormat: time.RFC3339,
	Color:      hclog.AutoColor,
})

const (
	defaultBaseURL          = "https://api.katapult.io"
	defaultSchemaTemplate   = "{{ .BaseURL }}/{{ .Name }}/{{ .Version }}/schema"
	defaultFileNameTemplate = "{{ .Name }}/{{ .Version }}.json"
)

type configuration struct {
	Name             string
	Version          string
	BaseURL          string
	outputDir        string
	forceUpdate      bool
	urlTemplate      string
	fileNameTemplate string
	logLevel         string
}

func configure() (*configuration, *flag.FlagSet, error) {
	config := &configuration{}

	wd, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	fs := flag.NewFlagSet("codegen", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(),
			undent.String(`
				usage: %s [<options>]
				Generate resources based on Katapult API schemas.

				Options:
				`,
			),
			fs.Name(),
		)
		fs.PrintDefaults()
	}

	fs.StringVar(&config.Name, "n", "",
		"name of schema to fetch (required)")
	fs.StringVar(&config.Version, "v", "",
		"version of schema to fetch (required)")
	fs.StringVar(&config.BaseURL, "u", defaultBaseURL,
		"base URL of Katapult API to fetch schema from")
	fs.StringVar(&config.outputDir, "o", wd,
		"output directory to write schema files to")
	fs.BoolVar(&config.forceUpdate, "F", false, "force update existing file")
	fs.StringVar(&config.urlTemplate, "s", defaultSchemaTemplate,
		"schema URL template")
	fs.StringVar(&config.fileNameTemplate, "f", defaultFileNameTemplate,
		"schema filename template")
	fs.StringVar(&config.logLevel, "l", "info", "log level")

	err = fs.Parse(os.Args[1:])
	if err != nil {
		return nil, nil, err
	}

	switch strings.ToLower(os.Getenv("SCHEMA_FORCE_UPDATE")) {
	case "yes", "true", "1":
		config.forceUpdate = true
	}

	if v := os.Getenv("SCHEMA_BASE_URL"); v != "" {
		config.BaseURL = v
	}

	return config, fs, nil
}

func main() {
	err := mainE()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(127)
	}
}

func mainE() error {
	config, fs, err := configure()
	if err != nil {
		return err
	}

	logLevel := hclog.LevelFromString(config.logLevel)
	if logLevel == hclog.NoLevel {
		return fmt.Errorf("invalid log level \"%s\"", config.logLevel)
	}
	logger.SetLevel(logLevel)

	if config.Name == "" || config.Version == "" {
		fs.Usage()
		os.Exit(1)
	}

	fileTpl, err := template.New("filename").Parse(config.fileNameTemplate)
	if err != nil {
		return err
	}
	urlTpl, err := template.New("filename").Parse(config.urlTemplate)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	err = fileTpl.Execute(buf, config)
	if err != nil {
		return err
	}
	filename := buf.String()

	targetFile := filepath.Join(config.outputDir, filename)

	if !config.forceUpdate && fileExists(targetFile) {
		logger.Info("no action: schema file already exists, "+
			"set -u or SCHEMA_FORCE_UPDATE env var to update",
			"file", targetFile,
		)

		return nil
	}

	logger.Info("schema file does not exist",
		"file", targetFile,
	)

	buf = &bytes.Buffer{}
	err = urlTpl.Execute(buf, config)
	if err != nil {
		return err
	}
	schemaURL := buf.String()

	logger.Info("fetching schema", "url", schemaURL)
	b, err := getSchema(schemaURL)
	if err != nil {
		return err
	}

	buf = &bytes.Buffer{}
	err = json.Indent(buf, b, "", "  ")
	if err != nil {
		return err
	}

	logger.Info("writing schema file", "file", targetFile, "size", buf.Len())
	err = os.MkdirAll(filepath.Dir(targetFile), 0o755)
	if err != nil {
		return err
	}

	//nolint:gosec
	err = os.WriteFile(targetFile, buf.Bytes(), 0o644)
	if err != nil {
		return err
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

func getSchema(u string) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	if t := os.Getenv("KATAPULT_API_KEY"); t != "" {
		req.Header.Set("Authorization", "Bearer "+t)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received HTTP %d status", resp.StatusCode)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}
