package suite

import (
	"context"
	"gRPC/internal/config"
	"net"
	"strconv"
	"testing"
	"time"

	ssov1 "github.com/PoliakovIvan2606/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcHost = "localhost"
)

type Suite struct {
	*testing.T
	Cfg *config.Config
	AuthClient ssov1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MastLoadByPath("../config/local.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(ctx,
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()), // Используем insecure-коннект для тестов
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T: t,
		Cfg: cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}