package cmd

import (
	"os"
	"os/exec"

	"github.com/rawdaGastan/gridify/internal/config"
	"github.com/rawdaGastan/gridify/internal/deployer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Destroy(debug bool) error {
	config, err := config.LoadConfigData()
	if err != nil {
		log.Error().Err(err).Msg("failed to load configuration try logging again using gridify login")
		return err
	}

	repoURL, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
	if err != nil {
		log.Error().Err(err).Msg("failed to get remote repository url")
		return err
	}

	logLevel := zerolog.InfoLevel
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
}
