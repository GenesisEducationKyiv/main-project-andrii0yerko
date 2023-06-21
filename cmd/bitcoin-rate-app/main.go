package main

import (
	"bitcoinrateapp/pkg/core"
	"bitcoinrateapp/pkg/http"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Load app configuration from both config file and command line args
// Custom config path can be passed with `--config path/to/my/config.yaml`
// By default - searching for file `./config.yaml`
func parseConfiguration() {
	// Set up the command line flags
	pflag.String("config", "config.yaml", "Config file name. Supported types are yaml, json, toml, ini, env")

	pflag.String("sender.smtpPort", "", "SMTP port")
	pflag.String("sender.smtpHost", "", "SMTP host")
	pflag.String("sender.from", "", "From email address")
	pflag.String("sender.password", "", "SMTP password (optional)")
	pflag.String("storage.filename", "emails.dat", "Filename for emails storage. Default is emails.dat")
	pflag.String("server.host", "0.0.0.0", "Host to serve HTTP api. Default is 0.0.0.0")
	pflag.String("server.port", "3333", "Post to serve HTTP api. Default is 3333")

	pflag.Parse()

	// Bind command line flags to Viper
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatalf("Error binding flags: %s", err)
	}

	viper.SetEnvPrefix("BTCAPP")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set up the config file name and path
	viper.SetConfigFile(viper.GetString("config"))
	viper.AddConfigPath(".")
	if err = viper.ReadInConfig(); err != nil {
		var notExists *fs.PathError
		if errors.As(err, &notExists) {
			log.Printf("Warning: %s\n", err)
		} else {
			log.Fatalf("Error reading config file: %s\n", err)
		}
	}

	// Ensure that all required values are given
	for _, field := range []string{
		"sender.smtpHost", "sender.smtpPort", "sender.from", "storage.filename", "server.host", "server.port",
	} {
		if viper.GetString(field) == "" {
			log.Fatalf(
				"\"%s\" value is missing! Please pass it as CLI arg with \"--%s value\","+
					" or add it to the config file with the same key name!",
				field,
				field,
			)
		}
	}
}

func main() {
	parseConfiguration()
	smtpPort := viper.GetString("sender.smtpPort")
	smtpHost := viper.GetString("sender.smtpHost")
	from := viper.GetString("sender.from")
	password := viper.GetString("sender.password")
	filename := viper.GetString("storage.filename")

	addr := fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port"))
	controller := core.NewController(smtpPort, smtpHost, from, password, filename)
	handler := http.NewExchangeRateHandler(controller)
	server := http.NewServer(handler, addr)

	err := server.Start()
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("server closed")
	} else if err != nil {
		log.Fatalf("error starting server: %s", err)
	}
}
