package main

type DomainConfig struct {
	Zone    string   `yaml:"zone"`
	Records []string `yaml:"records"`
}

type Config struct {
	Domains []DomainConfig `yaml:"domains"`
}