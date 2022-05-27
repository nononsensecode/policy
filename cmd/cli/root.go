package cli

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.zhservices.org/DigitalHealthPlatform/services/blue-go-base.git/configs"
	"gitlab.zhservices.org/DigitalHealthPlatform/services/blue-go-base.git/infrastructure/mongodb"
	"gitlab.zhservices.org/DigitalHealthPlatform/services/blue-go-base.git/interfaces/httpsrvr"
	"nononsensecode.com/policy/pkg/application"
	"nononsensecode.com/policy/pkg/domain/model/policy"
	"nononsensecode.com/policy/pkg/infrastructure/nosql/mongo"
	"nononsensecode.com/policy/pkg/interfaces/openapi"
)

var (
	rootCmd = cobra.Command{
		Use: "policy",
		Run: func(cmd *cobra.Command, args []string) {
			initServer()
		},
	}
	cfg     *configs.Config
	cfgFile string
)

func Execute() error {
	return rootCmd.Execute()
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	fmt.Printf("Reading configuration file %s...\n", cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config file \"%s\" cannot be read: %w", cfgFile, err))
	}

	cfg = new(configs.Config)
	fmt.Println("Loading configuration to memory....")
	if err := viper.Unmarshal(cfg); err != nil {
		panic(fmt.Errorf("configuration is invalid: %w", err))
	}

	fmt.Println("configuration loaded successfully to memory")
}

func initServer() {
	fmt.Println("Initializing configuration...")
	cfg.Init()

	fmt.Println("initialzing the http server...")
	var policyRepo policy.Repository

	if cfg.IsMongoDbEnabled() {
		provider := mongodb.DefaultClientProvider{}
		policyRepo = mongo.NewPolicyRepository(provider)
	}

	server := openapi.NewServer(application.NewPolicyService(policyRepo))
	fmt.Println("starting server...")
	httpsrvr.RunHTTPServer(
		cfg.HttpAddress(),
		func(router chi.Router) http.Handler {
			return openapi.HandlerFromMux(server, router)
		},
		[]func(http.Handler) http.Handler{},
		cfg.HttpCorsOrigins(),
		cfg.HttpApiPrefix(),
	)
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().
		StringVar(&cfgFile, "config", "",
			"configuration file in yaml format (Default is ./config.yaml)")
}
