/*
Copyright © 2025 Ulrich Wisser

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/apex/log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rpkistats",
	Short: "RPKI statistics for domain names",
	Long: `RPKI statistics for domain names`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().CountP(VERBOSE, VERBOSE_SHORT, "repeat for more verbose printouts")
	rootCmd.PersistentFlags().StringP(CONFIG_FILE, CONFIG_FILE_SHORT, "", "config file (default is $HOME/.tldrpki)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Use flags for viper values
	viper.BindPFlags(rootCmd.Flags())
	viper.BindPFlags(rootCmd.PersistentFlags())	
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if viper.GetString(CONFIG_FILE) != "" {
		// Use config file from the flag.
		viper.SetConfigFile(viper.GetString(CONFIG_FILE))
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".tldrpki" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".tldrpki")
	}

	// read in environment variables that match
	viper.SetEnvPrefix("RPKISTATS")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// init log level
	switch viper.GetInt(VERBOSE) {
	case VERBOSE_QUIET:
		log.SetLevel(log.FatalLevel)
	case VERBOSE_ERROR:
		log.SetLevel(log.ErrorLevel)
	case VERBOSE_WARNING:
		log.SetLevel(log.WarnLevel)
	case VERBOSE_INFO:
		log.SetLevel(log.InfoLevel)
	case VERBOSE_DEBUG:
		log.SetLevel(log.DebugLevel)
	default:
		if viper.GetInt(VERBOSE) > VERBOSE_DEBUG {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.ErrorLevel)
		}
	}

}




