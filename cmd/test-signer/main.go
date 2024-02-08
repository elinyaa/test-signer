package main

import (
	"context"
	"flag"
	"log"
	gohttp "net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/elinyaa/test-signer/cmd/test-signer/conf"
	"github.com/elinyaa/test-signer/internal"
	redisInfra "github.com/elinyaa/test-signer/internal/infrastructure/redis"
	"github.com/elinyaa/test-signer/internal/interface/http"
	"github.com/elinyaa/test-signer/internal/usecase"
	"github.com/elinyaa/test-signer/pkg/lib"
)

func main() {
	logTarget := flag.String("log-target", "stdout",
		"Log target: stdout, stderr, or a file path prefixed with 'file:'")
	logPrefix := flag.String("log-prefix", "", "Log prefix")
	logFlags := flag.Int("log-flags", log.LstdFlags, "Log flags")

	httpAddr := flag.String("http-addr", ":8080", "HTTP server address")
	flag.Parse()

	secret, exists := os.LookupEnv("JWT_SECRET")
	if !exists {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	logger := conf.BuildLogger(*logTarget, *logPrefix, *logFlags)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer stop()

	redisClient, err := lib.InitializeAndTestRedisClient()
	if err != nil {
		logger.Fatalf("Error initializing Redis client: %v", err)
	}

	// consider commenting in the following line and commenting out the next line to use the in-memory repository
	// userRepository := inmem.NewUserRepository()
	userRepository := redisInfra.NewUserRepository(redisClient)

	signAnswerUsecase := usecase.NewSignAnswer(userRepository)
	verifySignatureUsecase := usecase.NewVerifySignature(userRepository)

	app := internal.NewApp(logger, signAnswerUsecase, verifySignatureUsecase)

	server := http.NewServer(app, logger, *httpAddr, secret)
	if err := server.Start(ctx); err != gohttp.ErrServerClosed {
		logger.Fatalf("HttpServer ended unexpectedly with error: %v", err)
	}
}
