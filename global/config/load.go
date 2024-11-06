package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	errs "github.com/pkg/errors"
	"github.com/spf13/viper"
)

func Load(filePath string) error {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found. Proceeding without it.")
	}

	viper.SetConfigFile(filePath)

	if err := viper.ReadInConfig(); err != nil {
		return errs.Wrap(err, "failed to load config")
	}

	if err := applyEnv(); err != nil {
		return errs.Wrap(err, "apply env failed")
	}

	var config RuntimeConfig
	if err := viper.UnmarshalExact(&config); err != nil {
		return errs.Wrap(err, "config unmarshaling failed")
	}

	conf = config

	return nil
}

func applyEnv() (err error) {
	for _, key := range viper.AllKeys() {
		val := viper.GetString(key)

		if strings.HasPrefix(val, "${") && strings.HasSuffix(val, "}") {
			k := val

			val = os.ExpandEnv(val)

			if val == "" || val == k {
				err = errors.Join(err, fmt.Errorf("%s: env var does not exist", k))
			}

			viper.Set(key, val)
		}
	}
	return
}
