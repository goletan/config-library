# Goletan Configuration Library

The Goletan Configuration Library simplifies and standardizes configuration management across Nemetons. It provides an easy and consistent way to load, validate, and access configuration settings in Go services, while being adaptable to different domains through its flexibility.

## Installation

To use this library, add it to your Go module by running: `go install github.com/goletan/config`

## Features

Environment Overrides: Easily override settings with environment variables.
Domain-Specific Configuration: Load and manage configuration specific to a particular domain or Nemeton.
Flexible Sources: Supports loading from configuration files (YAML/JSON) and environment variables.

## Getting Started

To start using the Goletan Configuration Library, import the package and load the configuration with `LoadConfig()`, providing the configuration file name and the list of paths where it might be located.

## How to Customize

Modify the config.yaml file to set the appropriate values for your environment. Set environment variables to override the file settings if needed. Define your own configuration struct to map the loaded settings to domain-specific values.

## License

This library is released under the Apache 2.0 License.
