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
	"github.com/cskr/pubsub"
	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/steadydrive"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "bogie-edge",
	Short: "Edge Component of the Bogie Project",
	Long:  `Edge Component of the Bogie Project`,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	ps := pubsub.New(10)

	steadydrive, err := steadydrive.New(viper.Sub("steadydrive"), ps)
	if err != nil {
		log.Fatal().Msgf("steadydrive: %s", err)
	}
	steadydrive.Run()
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
