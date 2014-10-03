package go_common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var BASE_PATH = ""

type Configuration interface{}
type BaseConfig struct {
	Description    string
	Version        string
	IntraserverKey string
	LogPath        string
}
type StatsDConfiguration struct {
	Host       string
	Port       string
	ClientName string
	Interval   string
}
type PostgreSQLConfiguration struct {
	Name           string
	User           string
	Password       string
	Host           string
	Port           string
	SSLMode        string
	ConnectTimeout string
}

/* Attempts to retreive the configuration file with the given name and load
it. On failure, prints a warning and loads the base_config.json file. */
func GetConfig(filename string, cfg Configuration) (string, error) {
	base_config_err := loadBaseConfigFile(cfg)
	if base_config_err != nil {
		return "", base_config_err
	}

	config_file, err := os.OpenFile(BASE_PATH+"etc/"+filename+".json",
		os.O_RDONLY, 0655)
	if err != nil {
		// Use the base config
		fmt.Printf("WARNING: configuration file '%s' was not found in maestro/etc. "+
			"Using default instead (base_config.json).\n", filename)
		base_config_err := loadBaseConfigFile(cfg)
		if base_config_err != nil {
			return "", base_config_err
		}
		return "base_config.json", nil
	}

	// Have a good config file, override base config loaded earlier
	parse_err := parseConfigFile(config_file, cfg)
	if parse_err != nil {
		return "", parse_err
	}

	return filename + ".json", nil
}

func parseConfigFile(config_file *os.File, cfg Configuration) error {
	bytes, io_err := ioutil.ReadAll(config_file)
	if io_err != nil {
		return io_err
	}
	unmarshal_err := json.Unmarshal(bytes, &cfg)
	if unmarshal_err != nil {
		return unmarshal_err
	}

	return nil
}

func loadBaseConfigFile(cfg Configuration) error {
	base_config_file, err := os.OpenFile(BASE_PATH+"etc/base_config.json",
		os.O_RDONLY, 0655)
	if err != nil {
		return err
	}
	parse_err := parseConfigFile(base_config_file, cfg)
	if parse_err != nil {
		return parse_err
	}
	return nil
}

func SetConfig(config_filename *string, executable_name string,
	cfg Configuration) (string, error) {

	BASE_PATH = GetBasePathForExecutable(executable_name)

	used_config_filename, config_err := GetConfig(*config_filename, cfg)
	if config_err != nil {
		return "", config_err
	}

	return used_config_filename, nil
}
