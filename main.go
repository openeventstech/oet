package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use: "oetdctl",
	}
)

var logger *slog.Logger

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "explain what is being done")
	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()

	programLevel := new(slog.LevelVar)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel}))
	slog.SetDefault(logger)
	if viper.GetBool("verbose") {
		programLevel.Set(slog.LevelDebug)
	}
	slog.Debug("logging", slog.Any("level", programLevel))
}
