/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type EndpointConfig struct {
	StartIn string `mapstructure:"start_in"`
	Command string `mapstructure:"command"`
}

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "propel",
		Short: "tool for triggering pre-defined actions on a remote system",
		Long: `Use propel as part of your remote management system.

It can be particularly useful for ad hoc management tasks and deployment processes.
It exposes webhooks that can be used to run a predetermined command on the local
system.`,
		Run: func(cmd *cobra.Command, args []string) {},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().Int("port", 42424, "port number to listen on")
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
}

func executeCommand(command string, startIn string) {
	fmt.Printf("running command %s...\n", command)
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = startIn

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("propel_config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	port := viper.GetInt("port")
	fmt.Println("port number is: ", port)

	endpointsParsed := make(map[string]*EndpointConfig)
	err := viper.UnmarshalKey("endpoints", &endpointsParsed)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	for endpoint := range endpointsParsed {
		startIn := endpointsParsed[endpoint].StartIn
		command := endpointsParsed[endpoint].Command
		fmt.Printf(
			"endpoint: '%s' start_in: '%s' command: '%s'...\n",
			endpoint,
			startIn,
			command,
		)
		mux.HandleFunc(
			("/" + endpoint),
			func(w http.ResponseWriter, r *http.Request) {
				executeCommand(command, startIn)
			},
		)
	}

	err = http.ListenAndServe(":"+strconv.Itoa(port), mux)
	if err != nil {
		panic(err)
	}
}
