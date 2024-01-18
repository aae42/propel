/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"log/slog"
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
		Run: func(cmd *cobra.Command, args []string) {
			if err := runServer(); err != nil {
				panic(err)
			}
			os.Exit(0)
		},
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

func runServer() error {
	port := viper.GetInt("port")
	slog.Info("running on port " + strconv.Itoa(port))

	endpointsParsed := make(map[string]*EndpointConfig)
	err := viper.UnmarshalKey("endpoints", &endpointsParsed)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	// logFile, err := os.OpenFile("propel_log.ndjson", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	// if err != nil {
	// 	panic(err)
	// }
	// mw := io.MultiWriter(os.Stderr, logFile)
	// slog.SetOutput(mw)

	for endpoint := range endpointsParsed {
		startIn := endpointsParsed[endpoint].StartIn
		command := endpointsParsed[endpoint].Command
		endpointToPass := endpoint
		slog.Info("adding '" + endpoint + "' endpoint...")
		mux.HandleFunc(
			("/" + endpoint),
			func(w http.ResponseWriter, r *http.Request) {
				err = executeCommand(endpointToPass, command, startIn)
				if err != nil {
					slog.Error(err.Error())
				}
			},
		)
	}

	err = http.ListenAndServe(":"+strconv.Itoa(port), mux)
	if err != nil {
		return err
	}
	return nil
}

func executeCommand(endpoint string, command string, startIn string) error {
	slog.Debug("POST to " + endpoint + ", running command '" + command + "'...")
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = startIn

	pipe, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	done := make(chan struct{})

	scanner := bufio.NewScanner(pipe)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			slog.Info("output",
				"endpoint", endpoint,
				"command", command,
				"output", line,
			)
		}

		done <- struct{}{}
	}()
	err := cmd.Start()
	if err != nil {
		return err
	}

	<-done

	err = cmd.Wait()
	if err != nil {
		return err
	}
	slog.Info("run complete",
		"endpoint", endpoint,
		"command", command,
	)
	return nil
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
		slog.Info("Using config file: " + viper.ConfigFileUsed())
	}
}
