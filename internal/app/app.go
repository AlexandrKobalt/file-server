package app

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/pkg/errors"

	"github.com/AlexandrKobalt/trip-track/backend/file-server/config"
	grpcserver "github.com/AlexandrKobalt/trip-track/backend/file-server/pkg/grpc/server"
	"github.com/AlexandrKobalt/trip-track/backend/file-server/pkg/lifecycle"
	fileserverproto "github.com/AlexandrKobalt/trip-track/backend/proto/fileserver"

	filedelivery "github.com/AlexandrKobalt/trip-track/backend/file-server/internal/file/delivery/grpc"
	fileservice "github.com/AlexandrKobalt/trip-track/backend/file-server/internal/file/service"
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
	grpcServer, err := grpcserver.New(a.cfg.GRPC)
	if err != nil {
		log.Fatalf("failed to initialize gRPC server: %s", err.Error())
	}

	fileService := fileservice.New(a.cfg.Service.File)
	fileDelivery := filedelivery.New(fileService)

	fileserverproto.RegisterFileServer(grpcServer.App, fileDelivery)

	a.cmps = append(
		a.cmps,
		cmp{grpcServer, "gRPC Server"},
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
