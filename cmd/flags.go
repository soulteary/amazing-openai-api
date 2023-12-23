package cmd

import (
	"github.com/soulteary/amazing-openai-api/internal/fn"
	AoaModel "github.com/soulteary/amazing-openai-api/internal/model"
)

// refs: https://github.com/soulteary/flare/blob/main/cmd/flags.go
func parseEnvVars() AoaModel.Flags {
	// use default values
	flags := AoaModel.Flags{
		DebugMode:   false,
		ShowVersion: false,
		ShowHelp:    false,

		Type: _DEFAULT_TYPE,
		Port: _DEFAULT_PORT,
		Host: _DEFAULT_HOST,
	}

	// check and set port
	flags.Port = fn.GetIntOrDefaultFromEnv(_ENV_KEY_NAME_PORT, _DEFAULT_PORT)
	if flags.Port <= 0 || flags.Port > 65535 {
		flags.Port = _DEFAULT_PORT
	}

	// check and set host
	flags.Host = fn.GetStringOrDefaultFromEnv(_ENV_KEY_NAME_HOST, _DEFAULT_HOST)
	if !fn.IsValidIPAddress(flags.Host) {
		flags.Host = _DEFAULT_HOST
	}

	// check and set type
	flags.Type = fn.GetStringOrDefaultFromEnv(_ENV_KEY_SERVICE_TYPE, _DEFAULT_TYPE)
	// TODO support all types
	if flags.Type != "azure" &&
		flags.Type != "yi" {
		flags.Type = _DEFAULT_TYPE
	}
	return flags
}

// func parseCLI() {
// TODO: parse command line flags
// }
