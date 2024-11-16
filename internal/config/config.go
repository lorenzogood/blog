package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/ardanlabs/conf/v3"
	"go.uber.org/zap"
)

func Parse(prefix string, c any) {
	help, err := conf.Parse(prefix, c)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			os.Exit(0)
		}

		zap.L().Error("failed to parse configuration", zap.Error(err))
		os.Exit(1)
	}
}
