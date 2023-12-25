package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RSO-project-Prepih/get-photo-info/health"
	"github.com/RSO-project-Prepih/get-photo-info/prometheus"
	"github.com/RSO-project-Prepih/get-photo-info/server"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	pb "github.com/RSO-project-Prepih/get-photo-info/github.com/RSO-project-Prepih/get-photo-info"
)

// Load .env file for environment variables
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	log.Println("Starting the application...")

	// Create a listener on TCP port
	log.Println("gRPC server listening on port 50051...")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Start the gRPC server
	log.Println("Starting gRPC server...")
	go func() {
		log.Println("Starting gRPC server...")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// Register the server
	photoServer := server.NewServer()
	pb.RegisterPhotoServiceServer(grpcServer, photoServer)

	log.Println("Starting the HTTP server...")
	r := gin.Default()
	// get health
	liveHandler, readyHandler := health.HealthCheckHandler()
	r.GET("/live", gin.WrapH(liveHandler))
	r.GET("/ready", gin.WrapH(readyHandler))

	// get metrics
	r.GET("/metrics", gin.WrapH(prometheus.GetMetrics()))

	srver := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srver.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	qoit := make(chan os.Signal, 1)
	signal.Notify(qoit, syscall.SIGINT, syscall.SIGTERM)
	<-qoit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srver.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")

}
