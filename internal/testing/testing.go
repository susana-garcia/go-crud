package testing

import (
	"context"
	"log/slog"
	"net"
	"os"
	"testing"

	"buf.build/go/protovalidate"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/susana-garcia/go-crud/config"
	"github.com/susana-garcia/go-crud/pb"
	"github.com/susana-garcia/go-crud/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

const bufSize = 1024 * 1024

var database = config.Database{
	Port:     config.GetEnv("DB_PORT", "5435"),
	Host:     config.GetEnv("DB_HOST", "localhost"),
	Name:     config.GetEnv("DB_NAME", "gocrudtest"),
	User:     config.GetEnv("DB_USER", "gocrudtest"),
	Password: config.GetEnv("DB_PASSWORD", "gocrudtest"),
}

type TestEnv struct {
	Client      pb.BloggerClient
	CancelFuncs []func()
}

func (te *TestEnv) Cancel() {
	for _, f := range te.CancelFuncs {
		f()
	}
}

func NewTestEnvWithRegistration(ctx context.Context, t *testing.T, register func(grpc.ServiceRegistrar)) *TestEnv {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	// create protovalidate validator
	validator, err := protovalidate.New()
	if err != nil {
		logger.Error("failed to create validator", "error", err)
		os.Exit(1)
	}

	// create gRPC server with validation interceptor
	s := grpc.NewServer(
		grpc.UnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
			// validate request using protovalidate if it's a protobuf message
			if msg, ok := req.(proto.Message); ok {
				if err := validator.Validate(msg); err != nil {
					return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
				}
			}
			return handler(ctx, req)
		}),
	)
	register(s)
	lis := bufconn.Listen(bufSize)
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Logf("Server exited with error: %v\n", err)
			os.Exit(1)
		}
	}()

	conn, err := grpcConn(ctx, lis)
	assert.NoError(t, err)

	client := pb.NewBloggerClient(conn)

	return &TestEnv{
		CancelFuncs: []func(){
			func() { _ = conn.Close() },
		},
		Client: client,
	}
}

func grpcConn(ctx context.Context, lis *bufconn.Listener) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer(lis)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func bufDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, s string) (net.Conn, error) {
		return lis.DialContext(ctx)
	}
}

// SetupDatabase opens an SQL connection and runs automigrate
func SetupDatabase(logger *slog.Logger) (*gorm.DB, error) {
	db := config.OpenConnection(database)

	// run DB migration
	logger.Info("running database migration for blogs table")
	err := db.AutoMigrate(&service.Blog{})
	if err != nil {
		logger.Error("error running auto migrate", "err", err)
		os.Exit(1)
	}

	logger.Info("database migration completed successfully")
	return db, nil
}

// CleanUpDatabaseEntries deletes previous entries
func CleanUpDatabaseEntries(db *gorm.DB, logger *slog.Logger) error {
	tx := db.Exec("DELETE FROM blogs")
	logger.Info("deleted", "rows", tx.RowsAffected)
	return tx.Error
}
