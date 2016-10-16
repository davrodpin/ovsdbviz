package main

import (
	"errors"
	"fmt"
	"os"
)

type schemaValue string

func (schema *schemaValue) Set(value string) error {
	if *schema != "" {
		return errors.New("OVSDB Schema File Path already set")
	}

	if _, err := os.Stat(value); os.IsNotExist(err) {
		return fmt.Errorf("OVSDB Schema File %s does not exists", value)
	}

	*schema = schemaValue(value)

	return nil
}

func (schema *schemaValue) String() string {
	return string(*schema)
}

type outputValue string

func (output *outputValue) Set(value string) error {
	if *output != "" {
		return errors.New("Output File Path Variable already set")
	}

	if _, err := os.Stat(value); err == nil {
		return fmt.Errorf("Output File %s already exists", value)
	}

	*output = outputValue(value)

	return nil
}

func (output *outputValue) String() string {
	return string(*output)
}
