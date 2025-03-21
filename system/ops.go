package system

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	defaultReadHeaderTimeout = 10 * time.Second
	defaultShutdownTimeout   = 30 * time.Second
)

func ServeMetricsAndHealthChecks(ctx context.Context, appName, address string) error {
	opServer := &http.Server{
		Addr:              address,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		slog.InfoContext(ctx, "ops server started", slog.String("address", address))
		if err := opServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("ops ListenAndServer: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		<-ctx.Done()
		slog.InfoContext(ctx, "stopping operational server")
		sCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
		defer cancel()
		return opServer.Shutdown(sCtx)
	})
	return eg.Wait()
}
