package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"easycoding/internal/app"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func initServe() error {
	rootCmd.AddCommand(serveCmd)
	return nil
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Server",
	Long:  `This command boots the web server and serves the application to the local network.`,
	Run:   runServer,
}

func runServer(_ *cobra.Command, _ []string) {
	instance := boot()

	go instance.ListenGrpc()
	go instance.ListenGrpcGateway()
	go instance.ListenSwagger()
	go instance.StartDaemons()

	graceful(instance, 30*time.Second, instance.Log)
}

func graceful(instance *app.Kernel, timeout time.Duration, logger *logrus.Logger) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := instance.Shutdown(ctx); err != nil {
		logger.Errorf("application shutdown error: %v\n", err)
	} else {
		logger.Infoln("application stopped")
	}
}
