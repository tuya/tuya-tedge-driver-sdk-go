package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/common"
)

type (
	// LogConfig logger common
	LogConfig struct {
		FileName string
		LogLevel string
	}

	// RPCConfig internal grpc server common
	RPCConfig struct {
		Address  string
		UseTLS   bool
		CertFile string
		KeyFile  string
	}

	// ClientInfo provides the host and port of another service in tedge.
	ClientInfo struct {
		Address      string
		UseTLS       bool
		CertFilePath string
	}

	ServiceInfo struct {
		ID     string
		Name   string
		Server RPCConfig
	}

	DriverConfig struct {
		Logger       LogConfig
		Clients      map[string]ClientInfo
		Service      ServiceInfo
		CustomConfig map[string]interface{}
	}
)

func ParseConfig(configPath string, dc *DriverConfig) error {
	var (
		err      error
		contents []byte
	)
	absPath, _ := filepath.Abs(configPath)
	fmt.Printf("read config file: %s\n", absPath)

	if contents, err = ioutil.ReadFile(configPath); err != nil {
		return fmt.Errorf("could not load configuration file: %s", err.Error())
	}

	if err = toml.Unmarshal(contents, &dc); err != nil {
		return fmt.Errorf("unmarshal configuration file contents error: %s", err.Error())
	}
	return nil
}

func (dc *DriverConfig) GetCustomConfig() map[string]interface{} {
	return dc.CustomConfig
}

func (dc *DriverConfig) ValidateClientConfig() error {
	resClient := dc.Clients[common.Resource]
	if len(resClient.Address) == 0 {
		return fmt.Errorf("resource address client not configured")
	}

	if resClient.UseTLS {
		if len(resClient.CertFilePath) == 0 {
			return fmt.Errorf("resource cert file not configured")
		}
	}
	return nil
}

func (dc *DriverConfig) WriteToFile(configPath string) error {
	var (
		err  error
		buff bytes.Buffer
	)
	e := toml.NewEncoder(&buff)
	if err = e.Encode(dc); err != nil {
		return err
	}

	absPath, _ := filepath.Abs(configPath)
	f, err := os.OpenFile(absPath, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(buff.Bytes())
	return err
}
