package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/jimeh/undent"
	"github.com/krystal/go-katapult/tools/codegen/gen"
)

const schemaURL = "https://api.katapult.io/%s/schema"

var logger = hclog.New(&hclog.LoggerOptions{
	Name:       "codegen",
	Level:      hclog.Info,
	Output:     os.Stderr,
	Mutex:      &sync.Mutex{},
	TimeFormat: time.RFC3339,
	Color:      hclog.AutoColor,
})

type indexMap map[string]bool

func (m indexMap) String() string {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}

	return strings.Join(keys, ",")
}

func (m indexMap) Set(value string) error {
	m[value] = true

	return nil
}

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)

	return nil
}

type configuration struct {
	GenTypes          indexMap
	PkgName           string
	SchemaFiles       stringSlice
	SchemaNames       stringSlice
	SchemaIncludePath string
	SchemaExcludePath string
	OutputDir         string
	LogLevel          string
	SkipCache         bool
}

func configure() (*configuration, *flag.FlagSet, error) {
	config := &configuration{
		GenTypes: indexMap{},
	}

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

	fs.Var(&config.GenTypes, "t", "type of files to generate (repeatable)")
	fs.Var(&config.SchemaFiles, "f", "path to schema files (repeatable)")
	fs.Var(&config.SchemaNames, "n", "APIs schema name to fetch (repeatable)")
	fs.StringVar(&config.SchemaIncludePath,
		"i", ".*", "regexp matching object IDs to include",
	)
	fs.StringVar(&config.SchemaExcludePath,
		"e", "", "regexp matching object IDs to exclude",
	)
	fs.StringVar(&config.OutputDir, "o", wd, "")
	fs.StringVar(&config.PkgName, "p", "", "output package name")
	fs.StringVar(&config.LogLevel, "l", "info", "log level")
	fs.BoolVar(&config.SkipCache, "s", false, "skip schema cache")

	err = fs.Parse(os.Args[1:])
	if err != nil {
		return nil, nil, err
	}

	return config, fs, nil
}

func main() {
	config, fs, err := configure()
	if err != nil {
		fatal(err)
	}

	logLevel := hclog.LevelFromString(config.LogLevel)
	if logLevel == hclog.NoLevel {
		log.Fatal(fmt.Errorf("invalid log level \"%s\"", config.LogLevel))
	}
	logger.SetLevel(logLevel)

	if len(config.GenTypes) == 0 || config.PkgName == "" ||
		(len(config.SchemaFiles) == 0 && len(config.SchemaNames) == 0) {
		fs.Usage()
		os.Exit(1)
	}

	if len(config.SchemaNames) > 0 {
		for _, n := range config.SchemaNames {
			filename, err2 := getSchema(n, config.SkipCache)
			if err2 != nil {
				fatal(err2)
			}

			config.SchemaFiles = append(config.SchemaFiles, filename)
		}
	}

	generator := &gen.Generator{
		PkgName:           config.PkgName,
		OutputDir:         config.OutputDir,
		SchemaIncludePath: config.SchemaIncludePath,
		SchemaExcludePath: config.SchemaExcludePath,
		SchemaFiles:       config.SchemaFiles,
		Logger:            logger,
	}

	for t := range config.GenTypes {
		switch t {
		case "errors":
			err = generator.Errors()
			if err != nil {
				fatal(err)
			}
		default:
			logger.Error("invalid option", "-t", t)
		}
	}
}

func getSchema(name string, skipCache bool) (string, error) {
	u := fmt.Sprintf(schemaURL, name)

	// Check for local cache file to avoid repeatedly downloading schema.
	h := fmt.Sprintf("%x", sha256.Sum256([]byte(u)))
	tmpFile := filepath.Join(os.TempDir(), "katapult-api-schema-"+h+".json")
	stat, err := os.Stat(tmpFile)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	// Return cache file if it exists and is less than 10 minutes old.
	if !skipCache && stat != nil && stat.Mode().IsRegular() {
		exp := stat.ModTime().Add(10 * time.Minute)

		if time.Now().Before(exp) {
			logger.Info(
				"schema cache found",
				"url", u, "cache_file", tmpFile,
			)

			return tmpFile, nil
		}

		logger.Info(
			"schema cache expired",
			"url", u, "cache_file", tmpFile,
		)
	}

	req, err := http.NewRequestWithContext(context.Background(), "GET", u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")

	if t := os.Getenv("KATAPULT_API_KEY"); t != "" {
		req.Header.Set("Authorization", "Bearer "+t)
	}

	logger.Info("downloading schema", "url", u)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("received HTTP %d status", resp.StatusCode)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	logger.Info(
		"writing schema cache",
		"cache_file", tmpFile, "size", hclog.Fmt("%d bytes", len(b)),
	)
	//nolint:gosec
	err = os.WriteFile(tmpFile, b, 0o644)
	if err != nil {
		return "", err
	}

	return tmpFile, nil
}

func fatal(err error) {
	logger.Error(err.Error())
	os.Exit(1)
}
