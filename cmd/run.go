/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Get RPKI statistics",
	Long: `Get RPKI statistics for a single domain name or a list of domain names and optionally save to MariaDB`,
	Run: execRun,
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")
	runCmd.Flags().StringP(DOMAIN_FILE, DOMAIN_FILE_SHORT, "", "name of file with the TLD list")
	runCmd.Flags().StringP(COUNTRIES_FILE, COUNTRIES_FILE_SHORT, "", "name of file with list of domain by country")

	// Use flags for viper values
	viper.BindPFlags(submitCmd.Flags())
}

func execRun(cmd *cobra.Command, args []string) {
	fmt.Println("run called")
}