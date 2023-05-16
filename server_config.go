package main

import (
	"flag"
	"os/user"
	"path/filepath"
	//"os/signal"
)

type ServerConfig struct {
	port           string
	configFullPath string
	uiDir          string
}

func defaultFilepath() (string, string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", "", err
	}
	fp := filepath.Join(currentUser.HomeDir, "lattice_cfg.json")
	help := `path_to_conf:\n\tfull path, including filename, 
	of the desired config. default: ` + fp
	return fp, help, nil
}

func NewServerConfig() (*ServerConfig, error) {
	defaultConfigFilepath, help, err := defaultFilepath()
	if err != nil {
		return &ServerConfig{}, err
	}
	configFilepath := flag.String("path_to_conf", defaultConfigFilepath, help)
	flag.Parse()

	return &ServerConfig{
		port:           ":8888",
		configFullPath: *configFilepath,
		uiDir:          "static",
	}, nil
}

func (sc *ServerConfig) Port() string { return sc.port }

func (sc *ServerConfig) PathToConfig() string { return sc.configFullPath }

func (sc *ServerConfig) PathToWebUI() string { return sc.uiDir }
