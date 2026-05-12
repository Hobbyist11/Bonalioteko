package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	homedir, _ = os.UserHomeDir()
	ebookdir   = filepath.Join(homedir, "Downloads/Ebooks/")
)

// AppDir is the name of the directory where the config file is stored.
const AppDir = "Bonalioteko"

// ConfigFileName is the name of the config file that gets created.
const (
	ConfigFileName  = "config.yml"
	RecentFilesName = "recents.json"
	maxRecent       = 5
)

// SettingsConfig struct represents the config for the settings.
type SettingsConfig struct {
	EbookDir string `yaml:"start_dir"`
}

type Config struct {
	Settings    SettingsConfig `yaml:"settings"`
	RecentFiles []string       `json:"recent_files"`
}

// configError represents an error that occurred while parsing the config file.
type configError struct {
	configDir string
	parser    ConfigParser
	err       error
}

// ConfigParser is the parser for the config file.
type ConfigParser struct{}

// getDefaultConfig returns the default config for the application.
func (parser ConfigParser) getDefaultConfig() Config {
	return Config{
		Settings: SettingsConfig{
			EbookDir: ebookdir,
		},
	}
}

// getDefaultConfigYamlContents returns the default config file contents.
func (parser ConfigParser) getDefaultConfigYamlContents() string {
	defaultConfig := parser.getDefaultConfig()
	yaml, _ := yaml.Marshal(defaultConfig)

	return string(yaml)
}

// Error returns the error message for when a config file is not found.
func (e configError) Error() string {
	return fmt.Sprintf(
		`Couldn't find a config.yml configuration file.
Create one under: %s
Example of a config.yml file:
%s
For more info, go to https://github.com/mistakenelf/fm
press q to exit.
Original error: %v`,
		path.Join(e.configDir, AppDir, ConfigFileName),
		e.parser.getDefaultConfigYamlContents(),
		e.err,
	)
}

// writeDefaultConfigContents writes the default config file contents to the given file.
func (parser ConfigParser) writeDefaultConfigContents(newConfigFile *os.File) error {
	_, err := newConfigFile.WriteString(parser.getDefaultConfigYamlContents())
	if err != nil {
		return err
	}

	return nil
}

// createConfigFileIfMissing creates the config file if it doesn't exist.
func (parser ConfigParser) createConfigFileIfMissing(configFilePath string) error {
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		newConfigFile, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o666)
		if err != nil {
			return err
		}

		defer newConfigFile.Close()
		return parser.writeDefaultConfigContents(newConfigFile)
	}

	return nil
}

// createRecentFileIfMissing creates the recents file if it doesn't exist.
func (parser ConfigParser) createRecentFileIfMissing(configFilePath string) error {
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		newRecentFile, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o666)
		if err != nil {
			if os.IsExist(err) {
				return nil // File already exists, no need to initialize
			}
			return err
		}
		defer newRecentFile.Close()

		if _, err := newRecentFile.WriteString("[]"); err != nil {
			return err
		}
	}

	return nil
}

func AddToRecentFileList(newFile string) error {
	var recent []string
	configFilePath, err := os.UserConfigDir()
	fullPath := filepath.Join(configFilePath, AppDir, RecentFilesName)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(fullPath)
	if err == nil {
		_ = json.Unmarshal(data, &recent)
	}

	updated := []string{newFile}
	for _, f := range recent {
		if f != newFile {
			updated = append(updated, f)
		}
	}

	if len(updated) > maxRecent {
		updated = updated[:maxRecent]
	}

	newData, err := json.MarshalIndent(updated, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(fullPath, newData, 0o644)
}

// getRecentsFileOrCreate is missing returns the config file path or creates the config file if it doesn't exist.
func (parser ConfigParser) getRecentsFileOrCreate() (*string, error) {
	var err error
	configDir := os.Getenv("XDG_CONFIG_HOME")

	if configDir == "" {
		configDir, err = os.UserConfigDir()
		if err != nil {
			return nil, configError{parser: parser, configDir: RecentFilesName, err: err}
		}
	}

	prsConfigDir := filepath.Join(configDir, AppDir)
	err = os.MkdirAll(prsConfigDir, os.ModePerm)
	if err != nil {
		return nil, configError{parser: parser, configDir: RecentFilesName, err: err}
	}

	configFilePath := filepath.Join(prsConfigDir, RecentFilesName)
	err = parser.createRecentFileIfMissing(configFilePath)
	if err != nil {
		return nil, configError{parser: parser, configDir: RecentFilesName, err: err}
	}

	return &configFilePath, nil
}

// getConfigFileOrCreateIfMissing returns the config file path or creates the config file if it doesn't exist.
func (parser ConfigParser) getConfigFileOrCreateIfMissing() (*string, error) {
	var err error
	configDir := os.Getenv("XDG_CONFIG_HOME")

	if configDir == "" {
		configDir, err = os.UserConfigDir()
		if err != nil {
			return nil, configError{parser: parser, configDir: configDir, err: err}
		}
	}

	prsConfigDir := filepath.Join(configDir, AppDir)
	err = os.MkdirAll(prsConfigDir, os.ModePerm)
	if err != nil {
		return nil, configError{parser: parser, configDir: configDir, err: err}
	}

	configFilePath := filepath.Join(prsConfigDir, ConfigFileName)
	err = parser.createConfigFileIfMissing(configFilePath)
	if err != nil {
		return nil, configError{parser: parser, configDir: configDir, err: err}
	}

	return &configFilePath, nil
}

// parsingError represents an error that occurred while parsing the config file.
type parsingError struct {
	err error
}

// Error represents an error that occurred while parsing the config file.
func (e parsingError) Error() string {
	return fmt.Sprintf("failed parsing config.yml: %v", e.err)
}

// readConfigFile reads the config file and returns the config.
func (parser ConfigParser) readConfigFile(path string) (Config, error) {
	config := parser.getDefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		return config, configError{parser: parser, configDir: path, err: err}
	}

	err = yaml.Unmarshal((data), &config)
	return config, err
}

// readRecentsFile reads the recents file and returns the recentFiles.
func (parser ConfigParser) readRecentsFile(path string) ([]string, error) {
	var recentFiles []string
	data, err := os.ReadFile(path)
	if err != nil {
		return recentFiles, configError{parser: parser, configDir: path, err: err}
	}

	err = json.Unmarshal((data), &recentFiles)
	return recentFiles, err
}

// initParser initializes the parser.
func initParser() ConfigParser {
	return ConfigParser{}
}

// ParseRecent parses the config file and returns the config.
func ParseRecent() ([]string, error) {
	var recentFiles []string
	var err error

	parser := initParser()

	configFilePath, err := parser.getRecentsFileOrCreate()
	if err != nil {
		return recentFiles, parsingError{err: err}
	}

	recentFiles, err = parser.readRecentsFile(*configFilePath)
	if err != nil {
		return recentFiles, parsingError{err: err}
	}

	return recentFiles, nil
}

func GetRecentsSlice() ([]string, error) {
	cfg, err := ParseRecent()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// ParseConfig parses the config file and returns the config.
func ParseConfig() (Config, error) {
	var config Config
	var err error

	parser := initParser()

	configFilePath, err := parser.getConfigFileOrCreateIfMissing()
	if err != nil {
		return config, parsingError{err: err}
	}

	config, err = parser.readConfigFile(*configFilePath)
	if err != nil {
		return config, parsingError{err: err}
	}

	return config, nil
}
