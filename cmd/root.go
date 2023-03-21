package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// const slack_config string = "hydra_slackbot.yml"

// var slackcfg configuration.BottyConfig

var rootCmd = &cobra.Command{
	Use:   "discordbot",
	Short: "Discord bot",
	Long:  `Discord bot`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
