package config

import (
	_ "embed"
	"os"

	"github.com/sergei-dyshel/claug/internal/utils"

	"sigs.k8s.io/yaml"
)

var (
	//go:embed default.yml
	Default string

	current *Config
)

func Read(fname string) error {
	default_ := &Config{}
	err := yaml.UnmarshalStrict([]byte(Default), default_)
	utils.AssertErr(err)
	if fname == "" {
		current = default_
		return nil
	}
	config := &Config{}
	bytes, err := os.ReadFile(fname)
	if err != nil {
		return utils.Wrapf(err, "Failed to read %s", fname)
	}
	if err := yaml.UnmarshalStrict(bytes, config); err != nil {
		return utils.Wrapf(err, "Failed to parse YAML from %s", fname)
	}
	current = config
	return nil
}

func merge[T any](dst **T, src *T) {
	if *dst == nil {
		*dst = src
	}
}

func mergeOptions(dst *Options, src *Options) {
	merge(&dst.BracketedPaste, src.BracketedPaste)
	merge(&dst.History.Disabled, src.History.Disabled)
	merge(&dst.Selector, src.Selector)
}

func Get() *Config {
	if current == nil {
		utils.Panicf("Must read config before use")
	}
	return current
}
