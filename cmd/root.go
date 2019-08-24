package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"github.com/logrusorgru/aurora"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	bolt "go.etcd.io/bbolt"
)

var db bolt.DB
var cfgFile string
var Verbose bool
var C aurora.Aurora
var color bool
var DataDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "archivist",
	Short: "The simple tool to back up your online presence.",
	Long:  `The Archivist is an all in one tool intended to allow the user to backup their online presence in a single place.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	C = aurora.NewAurora(color)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.the-archivist.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Speak loudly.")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	rootCmd.PersistentFlags().BoolVar(&color, "colors", true, "Don't be flashy")
	viper.BindPFlag("colors", rootCmd.PersistentFlags().Lookup("colors"))
	rootCmd.PersistentFlags().StringVarP(&DataDir, "data", "d", "./data", "Location to store data.")
	viper.BindPFlag("data", rootCmd.PersistentFlags().Lookup("data"))

	db, err := bolt.Open(DataDir+"/archivist.db", 0600, nil)
	if err != nil {
		fmt.Println("Cannot open directory")
	}
	defer db.Close()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".the-archivist" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".the-archivist")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if Verbose {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}
