package cli

import (
	"encoding/json"
	"fmt"
	"github.com/goccy/go-yaml"
)

type OutputMode string

const (
	OutputModeText OutputMode = "text"
	OutputModeJson OutputMode = "json"
	OutputModeYaml OutputMode = "yaml"
)

func (o *OutputMode) Set(value string) error {
	switch value {
	case string(OutputModeText), string(OutputModeJson), string(OutputModeYaml):
		*o = OutputMode(value)
		return nil
	default:
		return fmt.Errorf("invalid output mode: %s", value)
	}
}

func (o *OutputMode) String() string {
	return string(*o)
}

func (o *OutputMode) Type() string {
	return "outputMode"
}

func WriteOutput(v interface{}) error {
	switch FlagOutputMode {
	case OutputModeJson:
		return writeJsonOutput(v)
	case OutputModeYaml:
		return writeYamlOutput(v)
	case OutputModeText:
		fallthrough
	default:
		return writeTextOutput(v)
	}
}

func writeTextOutput(v interface{}) error {
	PrettyPrint(v)
	return nil
}

func writeJsonOutput(v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))
	return nil
}

func writeYamlOutput(v interface{}) error {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))
	return nil
}
