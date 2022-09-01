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
	"os"

	"github.com/cskr/pubsub"
	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/gnss"
	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/metrics"
	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/nats"
	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/position"
	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/sensor"
	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/steadydrive"
	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/triggerunit"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type gloabalConfiguration struct {
	NatsAddress string
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
	ps := pubsub.New(10)

	err := viper.Unmarshal(&globalCfg)
	if err != nil {
		log.Fatal().Msgf("unmarshal global config %s", err)
	}

	natsConn, err := nats.Connect(globalCfg.NatsAddress)
	if err != nil {
		log.Fatal().Msgf("nats: %s", err)
	}
	steadydrive, err := steadydrive.New(viper.Sub("steadydrive"), ps)
	if err != nil {
		log.Fatal().Msgf("steadydrive: %s", err)
	}
	positionunit, err := position.NewFromViper(viper.Sub("position"), ps)
	if err != nil {
		log.Fatal().Msgf("positionunit: %s", err)
	}
	triggerunit, err := triggerunit.NewFromViper(viper.Sub("triggerunit"), ps)
	if err != nil {
		log.Fatal().Msgf("triggerunit: %s", err)
	}
	sensorunit, err := sensor.NewFromViper(viper.Sub("sensor"), ps, natsConn)
	if err != nil {
		log.Fatal().Msgf("sensorunit: %s", err)
	}
	gnssunit, err := gnss.NewFromViper(viper.Sub("gnss"), ps)
	if err != nil {
		log.Fatal().Msgf("gnss: %s", err)
	}
	metricsunit, err := metrics.NewFromViper(viper.Sub("metrics"), ps, natsConn)
	if err != nil {
		log.Fatal().Msgf("metricsunit: %s", err)
	}

	go steadydrive.Run()
	go positionunit.Run()
	go triggerunit.Run()
	gnssunit.Run()
	sensorunit.Run()
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
