package configloader

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

func Load(config any, envPath string) error {
	if envPath != "" {
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(envPath)

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				fmt.Println("No .env file found. Using environment variables.")
			} else {
				return fmt.Errorf("error reading config file: %w", err)
			}
		}
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	bindEnvVarsFromStruct(config)

	if err := viper.Unmarshal(config); err != nil {
		return fmt.Errorf("unable to decode config into struct: %w", err)
	}

	if err := validateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	return nil
}

func bindEnvVarsFromStruct(cfg any) {
	t := reflect.TypeOf(cfg)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag != "" {
			if err := viper.BindEnv(tag); err != nil {
				fmt.Printf("Warning: Failed to bind env var %s: %v\n", tag, err)
			}
		}
	}
}

func validateConfig(cfg any) error {
	var errors []string

	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		value := v.Field(i)
		tag := field.Tag.Get("mapstructure")

		if value.Kind() != reflect.Bool && value.IsZero() {
			errors = append(errors, fmt.Sprintf("%s is required but was empty or zero", tag))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("config validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}
