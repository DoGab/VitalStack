package main

import (
	"log/slog"
	"os"

	"github.com/dogab/macroguard/api/cmd"
)

func main() {
	logger := slog.Default()
	err := cmd.Execute()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
