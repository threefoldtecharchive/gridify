// Package cmd for handling command line arguments
package cmd

import (
	"os"
	"os/exec"

	"github.com/rawdaGastan/gridify/internal/config"
	"github.com/rawdaGastan/gridify/internal/deployer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the project in the current directory from threefold grid",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.LoadConfigData()
		if err != nil {
			log.Error().Err(err).Msg("failed to load configuration try logging again")
			return err
		}

		repoURL, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
		if err != nil {
			log.Error().Err(err).Msg("failed to get remote repository url")
			return err
		}

		logLevel := zerolog.InfoLevel
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			log.Error().Err(err).Send()
			return err
		}
		if debug {
			logLevel = zerolog.DebugLevel
		}

		logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).Level(logLevel).
			With().
			Timestamp().
			Logger()

		deployer, err := deployer.NewDeployer(config.Mnemonics, config.Network, string(repoURL), logger)
		if err != nil {
			log.Error().Err(err).Send()
			return err
		}

		err = deployer.Destroy()
		if err != nil {
			log.Error().Err(err).Send()
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}