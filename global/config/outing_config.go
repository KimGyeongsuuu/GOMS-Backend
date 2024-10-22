package config

import (
	"os"
	"strconv"
)

type OutingProperties struct {
	OutingExp int64
}

func LoadOutingProperties() (*OutingProperties, error) {
	outingExp, err := strconv.ParseInt(os.Getenv("OUTING_EXP"), 10, 64)
	if err != nil {
		return nil, err
	}

	return &OutingProperties{
		OutingExp: outingExp,
	}, nil
}
