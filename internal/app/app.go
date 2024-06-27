package app

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/pkg/errors"
	slogfiber "github.com/samber/slog-fiber"

	"github.com/AlexandrKobalt/trip-track_file-server/config"
	"github.com/AlexandrKobalt/trip-track_file-server/pkg/fiberapp"
	grpcserver "github.com/AlexandrKobalt/trip-track_file-server/pkg/grpc/server"
	"github.com/AlexandrKobalt/trip-track_file-server/pkg/lifecycle"
	fileserverproto "github.com/AlexandrKobalt/trip-track_proto/fileserver"

	filedeliverygrpc "github.com/AlexandrKobalt/trip-track_file-server/internal/file/delivery/grpc"
	filedeliveryhttp "github.com/AlexandrKobalt/trip-track_file-server/internal/file/delivery/http"
	filehandler "github.com/AlexandrKobalt/trip-track_file-server/internal/file/handler"
	fileservice "github.com/AlexandrKobalt/trip-track_file-server/internal/file/service"
)

var (
	ErrStartTimeout    = errors.New("start timeout")
	ErrShutdownTimeout = errors.New("shutdown timeout")
)

type (
	App struct {
		cfg    *config.Config
		cmps   []cmp
		logger *slog.Logger
	}
	cmp struct {
		Service lifecycle.Lifecycle
		Name    string
	}
)

func New(cfg *config.Config, logger *slog.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Start(ctx context.Context) error {
	fiberApp := fiberapp.New(a.cfg.FiberApp, a.logger)
	fiberApp.App.Use(slogfiber.New(a.logger))

	grpcServer, err := grpcserver.New(a.cfg.GRPC)
	if err != nil {
		log.Fatalf("failed to initialize gRPC server: %s", err.Error())
	}

	fileService := fileservice.New(a.cfg.Service.File)
	fileHandler := filehandler.New(a.cfg.Service.File.UploadDirectory)
	fileDeliveryGRPC := filedeliverygrpc.New(fileService)
	filedeliveryhttp.Map(fiberApp.App.Group("/files"), fileHandler)

	fileserverproto.RegisterFileServer(grpcServer.App, fileDeliveryGRPC)

	a.cmps = append(
		a.cmps,
		cmp{grpcServer, "gRPC Server"},
		cmp{fiberApp, "Fiber v2"},
	)

	okCh, errCh := make(chan any), make(chan error)

	go func() {
		err := error(nil)
		for _, c := range a.cmps {
			a.logger.Info("starting", "service", c.Name)

			err = c.Service.Start(context.Background())
			if err != nil {
				a.logger.Error("error on starting %s:\n%s", c.Name, err.Error())
				errCh <- errors.Wrapf(err, "cannot start %s", c.Name)

				return
			}
		}
		okCh <- nil
	}()

	select {
	case <-ctx.Done():
		return ErrStartTimeout
	case err := <-errCh:
		return err
	case <-okCh:
		a.logger.Info("application started!")
		return nil
	}
}

func (a *App) Stop(ctx context.Context) error {
	a.logger.Info("shutting down service...")
	okCh, errCh := make(chan struct{}), make(chan error)

	go func() {
		for i := len(a.cmps) - 1; i > 0; i-- {
			c := a.cmps[i]
			a.logger.Info("stopping", "service", c.Name)

			if err := c.Service.Stop(ctx); err != nil {
				a.logger.Error(err.Error())
				errCh <- err

				return
			}
		}
		okCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ErrShutdownTimeout
	case err := <-errCh:
		return err
	case <-okCh:
		a.logger.Info("Application stopped!")
		return nil
	}
}
func (a *App) GetStartTimeout() time.Duration { return a.cfg.StartTimeout.Duration }
func (a *App) GetStopTimeout() time.Duration  { return a.cfg.StopTimeout.Duration }
