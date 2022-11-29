package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"easycoding/internal/app"

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
	kernel := boot()

	go kernel.ListenGrpc()
	go kernel.ListenGrpcGateway()
	go kernel.ListenSwagger()
	go kernel.StartDaemons()

	graceful(kernel, 30*time.Second)
}

func graceful(instance *app.Kernel, timeout time.Duration) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := instance.Shutdown(ctx); err != nil {
		log.Fatalf("application shutdown error: %v\n", err)
	} else {
		log.Println("application stopped")
	}
}
