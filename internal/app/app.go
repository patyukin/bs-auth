package app

import (
	"context"
	"fmt"
	"github.com/patyukin/bs-auth/internal/config"
	"github.com/patyukin/bs-auth/internal/ratelimiter"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/patyukin/bs-auth/internal/closer"
	"github.com/patyukin/bs-auth/internal/interceptor"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	descAuth "github.com/patyukin/bs-auth/pkg/auth_v1"
	descUser "github.com/patyukin/bs-auth/pkg/user_v1"
	_ "github.com/patyukin/bs-auth/statik"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
}

func NewApp(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx, cfg, logger)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := a.RunGRPCServer()
		if err != nil {
			log.Fatalf("failed to run GRPC server: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			log.Fatalf("failed to run Swagger server: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context, cfg *config.Config, logger *slog.Logger) error {
	inits := []func(context.Context, *config.Config, *slog.Logger) error{
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
	}

	for _, f := range inits {
		err := f(ctx, cfg, logger)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context, cfg *config.Config, logger *slog.Logger) error {
	a.serviceProvider = newServiceProvider(cfg, logger, ratelimiter.NewRequestCounter())

	return nil
}

func (a *App) initGRPCServer(ctx context.Context, _ *config.Config, logger *slog.Logger) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)

	reflection.Register(a.grpcServer)

	descUser.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImpl(ctx))
	descAuth.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImpl(ctx))
	logger.Info("gRPC server initialized")

	return nil
}

func (a *App) initHTTPServer(ctx context.Context, cfg *config.Config, logger *slog.Logger) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	grpcAddress := cfg.Server.GRPC.Host + ":" + cfg.Server.GRPC.Port
	err := descUser.RegisterUserV1HandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		return fmt.Errorf("failed to descUser.RegisterUserV1HandlerFromEndpoint, error: %v", err)
	}

	err = descAuth.RegisterAuthV1HandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		return fmt.Errorf("failed to descAuth.RegisterAuthV1HandlerFromEndpoint, error: %v", err)
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	httpAddress := cfg.Server.HTTP.Host + ":" + cfg.Server.HTTP.Port
	a.httpServer = &http.Server{
		Addr:    httpAddress,
		Handler: corsMiddleware.Handler(mux),
	}
	logger.Info("HTTP server initialized")

	return nil
}

func (a *App) initSwaggerServer(_ context.Context, cfg *config.Config, logger *slog.Logger) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json", logger))

	swaggerAddress := cfg.Server.Swagger.Host + ":" + cfg.Server.Swagger.Port
	a.swaggerServer = &http.Server{
		Addr:    swaggerAddress,
		Handler: mux,
	}
	logger.Info("Swagger server initialized")

	return nil
}

func (a *App) RunGRPCServer() error {
	grpcAddress := a.serviceProvider.config.Server.GRPC.Host + ":" + a.serviceProvider.config.Server.GRPC.Port
	log.Printf("GRPC server is running on %s", grpcAddress)

	list, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runHTTPServer() error {
	httpAddress := a.serviceProvider.config.Server.HTTP.Host + ":" + a.serviceProvider.config.Server.HTTP.Port
	a.serviceProvider.logger.Info("HTTP server is running on %s", httpAddress)

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	swaggerAddress := a.serviceProvider.config.Server.Swagger.Host + ":" + a.serviceProvider.config.Server.Swagger.Port
	a.serviceProvider.logger.Info("Swagger server is running on %s", swaggerAddress)

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer func(file http.File) {
			err = file.Close()
			if err != nil {
				logger.Info("Failed to close file: %s", path)
			}
		}(file)

		logger.Info("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Served swagger file: %s", path)
	}
}
