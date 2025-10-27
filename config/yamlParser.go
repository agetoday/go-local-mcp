package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type DNSResolverConfig struct {
	InputFile  string `yaml:"input"`
	OutputFile string `yaml:"output"`
	Workers    int    `yaml:"workers"`
}

func ParseYAMLConfig(filePath string) (*DNSResolverConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	config := &DNSResolverConfig{Workers: 10} // 默认10个worker
	if err := yaml.Unmarshal(data, config); err != nil {
		fmt.Printf("Error parsing YAML: %v\n", err)
		return nil, err
	}
	fmt.Printf("Parsed config: %+v\n", config)
	return config, nil
}
