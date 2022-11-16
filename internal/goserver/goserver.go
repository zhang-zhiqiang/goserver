package goserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/marmotedu/goserver/internal/goserver/store"
	"github.com/marmotedu/goserver/internal/pkg/known"
	"github.com/marmotedu/goserver/internal/pkg/log"
	"github.com/marmotedu/goserver/internal/pkg/middleware"
	"github.com/marmotedu/goserver/pkg/db"
	"github.com/marmotedu/goserver/pkg/token"
	"github.com/marmotedu/goserver/pkg/util/homedir"
	"github.com/marmotedu/goserver/pkg/version/verflag"
)

const (
	// recommendedHomeDir defines the default directory used to place all goserver service configurations.
	recommendedHomeDir = ".goserver"

	// defaultConfigName specify the default configuration file for goserver.
	defaultConfigName = "goserver.yaml"
)

var cfgFile string

// NewGoServerCommand creates a *cobra.Command object with default parameters.
func NewGoServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goserver",
		Short: "A good Go practical project",
		Long: `A good Go practical project, used to create user with basic information.

Find more goserver information at:
    https://github.com/marmotedu/goserver/blob/master/docs/master/goserver.md`,

		// stop printing usage when the command errors
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			verflag.PrintAndExitIfRequested()

			log.Init(logOptions())
			defer log.Flush()

			return run()
		},
		PostRun: func(cmd *cobra.Command, args []string) {},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		},
	}

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.test.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Add --version flag to goserver command line
	verflag.AddFlags(cmd.PersistentFlags())

	return cmd
}

// Run runs the specified APIServer. This should never exit.
func run() error {
	// Do some initialization work
	// init store layer
	if err := initStore(); err != nil {
		return err
	}

	// init secret key in token package
	token.Init(viper.GetString("jwt_secret"), known.XUsernameKey)

	// set gin mode
	gin.SetMode(viper.GetString("runmode"))

	// create the gin engine
	g := gin.New()

	// load routers
	loadRouter(g, middleware.RequestID())

	// create http server instance
	insecureServer := &http.Server{
		Addr:    viper.GetString("addr"),
		Handler: g,
	}

	// create https server instance
	secureServer := &http.Server{
		Addr:    viper.GetString("tls.addr"),
		Handler: g,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
		if err := insecureServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Errorf("listen: %s\n", err)
		}
	}()

	// Start to listening the incoming requests.
	cert, key := viper.GetString("tls.cert"), viper.GetString("tls.key")
	if cert != "" && key != "" {
		log.Infof("Start to listening the incoming requests on https address: %s", viper.GetString("tls.addr"))
		go func() {
			if err := secureServer.ListenAndServeTLS(cert, key); err != nil && errors.Is(err, http.ErrServerClosed) {
				log.Errorf("listen: %s\n", err)
			}
		}()
	}

	// The context is used to inform the server goroutine it has 10 seconds to finish
	// the request it is currently handling and the ping process
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infof("Shutting down server...")

	if err := insecureServer.Shutdown(ctx); err != nil {
		log.Error("Insecure Server forced to shutdown:", log.Err(err))

		return err
	}

	if err := secureServer.Shutdown(ctx); err != nil {
		log.Error("Secure Server forced to shutdown:", log.Err(err))

		return err
	}

	log.Infof("Server exiting")

	return nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(homedir.HomeDir(), recommendedHomeDir))
		viper.SetConfigName(defaultConfigName)
	}

	viper.SetConfigType("yaml")    // 设置配置文件格式为YAML
	viper.AutomaticEnv()           // 读取匹配的环境变量
	viper.SetEnvPrefix("GOSERVER") // 读取环境变量的前缀为APISERVER
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Fprintln(os.Stdout, "Using config file:", viper.ConfigFileUsed())
}

func logOptions() *log.Options {
	return &log.Options{
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
		ErrorOutputPaths:  viper.GetStringSlice("log.error-output-paths"),
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		EnableColor:       viper.GetBool("log.enable-color"),
		Development:       viper.GetBool("log.development"),
		Name:              viper.GetString("log.name"),
	}
}

func initStore() error {
	dbOptions := &db.MySQLOptions{
		Host:                  viper.GetString("db.host"),
		Username:              viper.GetString("db.username"),
		Password:              viper.GetString("db.password"),
		Database:              viper.GetString("db.database"),
		MaxIdleConnections:    viper.GetInt("db.max-idle-connections"),
		MaxOpenConnections:    viper.GetInt("db.max-open-connections"),
		MaxConnectionLifeTime: viper.GetDuration("db.max-connection-life-time"),
		LogLevel:              viper.GetInt("db.log-level"),
	}

	ins, err := db.NewMySQL(dbOptions)
	if err != nil {
		return err
	}

	_ = store.NewStore(ins)

	return nil
}
