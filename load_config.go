package pocket

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

// LoadConfigFromEnv returns a config struct populated with environment variables.
//
// It uses the `env` struct tag to determine the environment variable name
// and the `default` tag to determine the default value if the environment variable is not set.
// It casts the value to the type specified in the struct field.
//
// Example:
//
//		  type AppConfig struct {
//			Port     int           `env:"PORT" default:"8080"`
//		    LogLevel string        `env:"LOG_LEVEL"`
//	       Timeout  time.Duration `env:"TIMEOUT" default:"10s"`
//		  }
//
//		  config, err := pocket.LoadConfigFromEnv[AppConfig]()
func LoadConfigFromEnv[T any]() (*T, error) {
	config := new(T)

	v := reflect.TypeOf(*config)

	for i := 0; i < v.NumField(); i++ {
		structField := v.Field(i).Name
		structFieldType := v.Field(i).Type
		envVarName := v.Field(i).Tag.Get("env")
		defaultValue := v.Field(i).Tag.Get("default")

		envVarValue, ok := os.LookupEnv(envVarName)
		if !ok {
			if defaultValue == "" {
				return nil, fmt.Errorf("missing env var %v (no default provided)", envVarName)
			}
			envVarValue = defaultValue
		}

		value, err := cast(structFieldType.Name(), envVarValue)
		if err != nil {
			return nil, err
		}

		reflect.ValueOf(config).Elem().FieldByName(structField).Set(value)
	}

	return config, nil
}

func cast(fieldType string, fieldValue string) (reflect.Value, error) {
	switch fieldType {
	case "string":
		return reflect.ValueOf(fieldValue), nil
	case "int":
		v, err := strconv.Atoi(fieldValue)
		if err != nil {
			e := fmt.Errorf("cannot parse %s as int: %w", fieldValue, err)
			return reflect.ValueOf(nil), e
		}
		return reflect.ValueOf(v), nil
	case "bool":
		v, err := strconv.ParseBool(fieldValue)
		if err != nil {
			e := fmt.Errorf("cannot parse %s as bool: %w", fieldValue, err)
			return reflect.ValueOf(nil), e
		}
		return reflect.ValueOf(v), nil
	case "Duration":
		v, err := time.ParseDuration(fieldValue)
		if err != nil {
			e := fmt.Errorf("cannot parse %s as time.Duration: %w", fieldValue, err)
			return reflect.ValueOf(nil), e
		}
		return reflect.ValueOf(v), nil
	default:
		return reflect.ValueOf(nil), fmt.Errorf("unsupported type %s", fieldType)
	}
}
