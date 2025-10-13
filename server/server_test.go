package server

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/susana-garcia/go-crud/pb"
	"github.com/susana-garcia/go-crud/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	crudtesting "github.com/susana-garcia/go-crud/internal/testing"
)

func TestCreateBlog(t *testing.T) {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	db, dbErr := crudtesting.SetupDatabase(logger)
	assert.NoError(t, dbErr)
	dbErr = crudtesting.CleanUpDatabaseEntries(db, logger)
	assert.NoError(t, dbErr)

	tEnv := crudtesting.NewTestEnvWithRegistration(ctx, t, func(reg grpc.ServiceRegistrar) {
		bService := service.New(db, logger)
		srv := New(bService, logger)
		srv.Register(reg)
	})
	defer tEnv.Cancel()

	tests := []struct {
		name      string
		request   service.Blog
		wantError *status.Status
	}{
		{
			name: "should create blog successfully",
			request: service.Blog{
				Title: "new title",
				Body:  "new body",
			},
			wantError: nil,
		},
		{
			name: "should fail when title is empty",
			request: service.Blog{
				Title: "",
				Body:  "new body",
			},
			wantError: status.New(codes.InvalidArgument, "title: value length must be at least 3 characters [string.min_len]"),
		},
		{
			name: "should fail when body is empty",
			request: service.Blog{
				Title: "new title",
				Body:  "",
			},
			wantError: status.New(codes.InvalidArgument, "body: value length must be at least 3 characters [string.min_len]"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tEnv.Client.CreateBlog(ctx, &pb.CreateBlogRequest{
				Title: tt.request.Title,
				Body:  tt.request.Body,
			})
			if tt.wantError != nil {
				s, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantError.Code(), s.Code())
				assert.Contains(t, s.Message(), tt.wantError.Message())
			} else {
				assert.Nil(t, err)
				assert.True(t, resp.Id > 0)
			}
		})
	}
}

func TestUpdateBlog(t *testing.T) {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	db, dbErr := crudtesting.SetupDatabase(logger)
	assert.NoError(t, dbErr)
	dbErr = crudtesting.CleanUpDatabaseEntries(db, logger)
	assert.NoError(t, dbErr)

	tEnv := crudtesting.NewTestEnvWithRegistration(ctx, t, func(reg grpc.ServiceRegistrar) {
		bService := service.New(db, logger)
		srv := New(bService, logger)
		srv.Register(reg)
	})
	defer tEnv.Cancel()

	// prepare test by creating a new entry
	resp, err := tEnv.Client.CreateBlog(ctx, &pb.CreateBlogRequest{
		Title: "title",
		Body:  "body",
	})
	assert.NoError(t, err)

	tests := []struct {
		name      string
		request   service.Blog
		wantError *status.Status
	}{
		{
			name: "should update blog successfully",
			request: service.Blog{
				ID:    uint(resp.Id),
				Title: "new title",
				Body:  "new body",
			},
			wantError: nil,
		},
		{
			name: "should update blog title successfully",
			request: service.Blog{
				ID:    uint(resp.Id),
				Title: "new title",
				Body:  "",
			},
			wantError: nil,
		},
		{
			name: "should update blog body successfully",
			request: service.Blog{
				ID:    uint(resp.Id),
				Title: "",
				Body:  "new body",
			},
			wantError: nil,
		},
		{
			name: "should fail when title and body are empty",
			request: service.Blog{
				ID:    uint(resp.Id),
				Title: "",
				Body:  "",
			},
			wantError: status.New(codes.InvalidArgument, "at_least_one_param"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tEnv.Client.UpdateBlog(ctx, &pb.UpdateBlogRequest{
				Id:    uint32(resp.Id),
				Title: tt.request.Title,
				Body:  tt.request.Body,
			})
			if tt.wantError != nil {
				s, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantError.Code(), s.Code())
				assert.Contains(t, s.Message(), tt.wantError.Message())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetBlog(t *testing.T) {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	db, dbErr := crudtesting.SetupDatabase(logger)
	assert.NoError(t, dbErr)
	dbErr = crudtesting.CleanUpDatabaseEntries(db, logger)
	assert.NoError(t, dbErr)

	tEnv := crudtesting.NewTestEnvWithRegistration(ctx, t, func(reg grpc.ServiceRegistrar) {
		bService := service.New(db, logger)
		srv := New(bService, logger)
		srv.Register(reg)
	})
	defer tEnv.Cancel()

	title := "title"
	// prepare test by creating a new entry
	resp, err := tEnv.Client.CreateBlog(ctx, &pb.CreateBlogRequest{
		Title: title,
		Body:  "body",
	})
	assert.NoError(t, err)

	tests := []struct {
		name      string
		request   *pb.GetBlogRequest
		wantError *status.Status
	}{
		{
			name: "should get blog by id successfully",
			request: &pb.GetBlogRequest{
				Value: &pb.GetBlogRequest_Id{
					Id: resp.Id,
				},
			},
			wantError: nil,
		},
		{
			name: "should get blog by title successfully",
			request: &pb.GetBlogRequest{
				Value: &pb.GetBlogRequest_Title{
					Title: title,
				},
			},
			wantError: nil,
		},
		{
			name: "should not get blog that does not exist",
			request: &pb.GetBlogRequest{
				Value: &pb.GetBlogRequest_Title{
					Title: "title that does not exist",
				},
			},
			wantError: status.New(codes.NotFound, ""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blog, err := tEnv.Client.GetBlog(ctx, tt.request)
			if tt.wantError != nil {
				s, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantError.Code(), s.Code())
				assert.Contains(t, s.Message(), tt.wantError.Message())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, resp.Id, blog.Item.Id)
			}
		})
	}
}

func TestGetBlogs(t *testing.T) {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	db, dbErr := crudtesting.SetupDatabase(logger)
	assert.NoError(t, dbErr)
	dbErr = crudtesting.CleanUpDatabaseEntries(db, logger)
	assert.NoError(t, dbErr)

	tEnv := crudtesting.NewTestEnvWithRegistration(ctx, t, func(reg grpc.ServiceRegistrar) {
		bService := service.New(db, logger)
		srv := New(bService, logger)
		srv.Register(reg)
	})
	defer tEnv.Cancel()

	title := "title"
	// prepare test by creating a new entries
	for i := 1; i < 11; i++ {
		_, err := tEnv.Client.CreateBlog(ctx, &pb.CreateBlogRequest{
			Title: title,
			Body:  "body",
		})
		assert.NoError(t, err)
	}

	tests := []struct {
		name           string
		request        *pb.GetBlogsRequest
		expectedAmount int
	}{
		{
			name: "should get blogs with pagination successfully",
			request: &pb.GetBlogsRequest{
				Limit: 3,
				Page:  1,
				Sort:  "Id asc",
			},
			expectedAmount: 3,
		},
		{
			name: "should get all blogs successfully",
			request: &pb.GetBlogsRequest{
				Limit: 12,
				Page:  1,
				Sort:  "Id asc",
			},
			expectedAmount: 10,
		},
		{
			name: "should not get blogs when no more records",
			request: &pb.GetBlogsRequest{
				Limit: 10,
				Page:  100,
				Sort:  "Id asc",
			},
			expectedAmount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tEnv.Client.GetBlogs(ctx, tt.request)
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedAmount, len(resp.Items))
		})
	}
}

func TestDeleteBlog(t *testing.T) {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	db, dbErr := crudtesting.SetupDatabase(logger)
	assert.NoError(t, dbErr)
	dbErr = crudtesting.CleanUpDatabaseEntries(db, logger)
	assert.NoError(t, dbErr)

	tEnv := crudtesting.NewTestEnvWithRegistration(ctx, t, func(reg grpc.ServiceRegistrar) {
		bService := service.New(db, logger)
		srv := New(bService, logger)
		srv.Register(reg)
	})
	defer tEnv.Cancel()

	// prepare test by creating a new entry
	resp, err := tEnv.Client.CreateBlog(ctx, &pb.CreateBlogRequest{
		Title: "title",
		Body:  "body",
	})
	assert.NoError(t, err)

	tests := []struct {
		name      string
		request   *pb.DeleteBlogRequest
		wantError *status.Status
	}{
		{
			name: "should delete blog successfully",
			request: &pb.DeleteBlogRequest{
				Id: resp.Id,
			},
			wantError: nil,
		},
		{
			name: "should not get error when deleting blog that does not exist",
			request: &pb.DeleteBlogRequest{
				Id: resp.Id,
			},
			wantError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tEnv.Client.DeleteBlog(ctx, tt.request)
			assert.Nil(t, err)
		})
	}
}
