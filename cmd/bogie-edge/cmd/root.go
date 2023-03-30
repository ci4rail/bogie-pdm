/*
Copyright © 2022 Ci4Rail GmbH <engineering@ci4rail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/ci4rail/bogie-pdm/cmd/bogie-edge/internal/daprpubsub"
	"github.com/ci4rail/bogie-pdm/cmd/bogie-edge/internal/export"
	"github.com/ci4rail/bogie-pdm/cmd/bogie-edge/internal/gnss"
	"github.com/ci4rail/bogie-pdm/cmd/bogie-edge/internal/metrics"
	"github.com/ci4rail/bogie-pdm/cmd/bogie-edge/internal/nats"
	"github.com/ci4rail/bogie-pdm/cmd/bogie-edge/internal/position"
	"github.com/ci4rail/bogie-pdm/cmd/bogie-edge/internal/sensor"
	"github.com/ci4rail/bogie-pdm/cmd/bogie-edge/internal/steadydrive"
	"github.com/ci4rail/bogie-pdm/cmd/bogie-edge/internal/triggerunit"
	"github.com/cskr/pubsub"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type gloabalConfiguration struct {
	NatsAddress   string // nats address. If set, publish directly to nats. Otherwise, publish to dapr pubsub
	NatsCredsPath string // nats credentials file path
	NetworkName   string // edgefarm network name, required for dapr pubsub
}

var (
	cfgFile   string
	globalCfg gloabalConfiguration
)

var rootCmd = &cobra.Command{
	Use:   "bogie-edge",
	Short: "Edge Component of the Bogie Project",
	Long:  `Edge Component of the Bogie Project`,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05.999Z07:00"})
	ps := pubsub.New(100)

	err := viper.Unmarshal(&globalCfg)
	if err != nil {
		log.Fatal().Msgf("unmarshal global config %s", err)
	}
	var exporter export.Exporter
	nodeID := os.Getenv("NODE_NAME")
	if globalCfg.NatsAddress != "" {
		credsPath := ""
		if globalCfg.NatsCredsPath != "" {
			credsPath = fmt.Sprintf("%s/%s.creds", globalCfg.NatsCredsPath, globalCfg.NetworkName)
			log.Info().Msgf("using nats creds file %s", credsPath)
		}
		exporter, err = nats.Connect(globalCfg.NatsAddress, nodeID, credsPath)
		if err != nil {
			log.Fatal().Msgf("nats: %s", err)
		}
	} else {
		address := os.Getenv("DAPR_GRPC_ADDRESS")
		if address == "" {
			log.Fatal().Msg("DAPR_GRPC_ADDRESS not set")
		}

		if nodeID == "" {
			log.Fatal().Msg("NODE_NAME not set")
		}
		exporter, err = daprpubsub.New(address, nodeID, globalCfg.NetworkName)
		if err != nil {
			log.Fatal().Msgf("dapr pubsub: %s", err)
		}
	}
	var steadyDrive *steadydrive.SteadyDrive

	sdConfig := viper.Sub("steadydrive")
	if sdConfig != nil {
		steadyDrive, err = steadydrive.New(sdConfig, ps)
		if err != nil {
			log.Fatal().Msgf("steadydrive: %s", err)
		}
	}
	positionunit, err := position.NewFromViper(viper.Sub("position"), ps)
	if err != nil {
		log.Fatal().Msgf("positionunit: %s", err)
	}
	triggerunit, err := triggerunit.NewFromViper(viper.Sub("triggerunit"), ps)
	if err != nil {
		log.Fatal().Msgf("triggerunit: %s", err)
	}
	var sensorUnit *sensor.Unit
	sensorcfg := viper.Sub("sensor")
	if sensorcfg != nil {
		sensorUnit, err = sensor.NewFromViper(sensorcfg, ps, exporter)
		if err != nil {
			log.Fatal().Msgf("sensorunit: %s", err)
		}
	}
	gnssunit, err := gnss.NewFromViper(viper.Sub("gnss"), ps)
	if err != nil {
		log.Fatal().Msgf("gnss: %s", err)
	}
	metricsunit, err := metrics.NewFromViper(viper.Sub("metrics"), ps, exporter)
	if err != nil {
		log.Fatal().Msgf("metricsunit: %s", err)
	}

	if steadyDrive != nil {
		go steadyDrive.Run()
	}
	go positionunit.Run()
	go triggerunit.Run()
	gnssunit.Run()
	if sensorUnit != nil {
		sensorUnit.Run()
	}
	go metricsunit.Run()
	select {}
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Msgf("Execute Root cmd: %s", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", ".bogie-edge-config.yaml", "config file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName(cfgFile) // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/") // call multiple times to add many search paths
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		log.Fatal().Msgf("fatal error config file: %s", err)
	}
}
