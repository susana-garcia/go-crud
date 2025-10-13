package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	logger *slog.Logger
}

func New(db *gorm.DB, logger *slog.Logger) *Service {
	return &Service{
		db:     db,
		logger: logger,
	}
}

type Blogger interface {
	GetAllBlogs(ctx context.Context, pagination *Pagination) (*Pagination, error)
	GetBlogByIDOrTitle(ctx context.Context, id uint, title string) (*Blog, error)
	CreateBlog(ctx context.Context, blog Blog) (uint, error)
	UpdateBlog(ctx context.Context, blog Blog) error
	DeleteBlog(ctx context.Context, id uint) error
}

type Blog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Body      string    `gorm:"type:text" json:"body"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for the Blog model
func (Blog) TableName() string {
	return "blogs"
}

func (s *Service) GetAllBlogs(ctx context.Context, pagination *Pagination) (*Pagination, error) {
	var blogs []Blog
	result := s.db.Scopes(paginate(blogs, pagination, s.db)).Find(&blogs)
	s.logger.Info(fmt.Sprintf("found %d blogs", result.RowsAffected))
	if result.Error != nil {
		s.logger.Error("unable to get all blogs", "error", result.Error)
		return nil, result.Error
	}
	pagination.Items = blogs
	return pagination, result.Error
}

func (s *Service) GetBlogByIDOrTitle(ctx context.Context, id uint, title string) (*Blog, error) {
	query := gorm.G[Blog](s.db)
	var blog Blog
	var err error
	if id > 0 {
		blog, err = query.Where("id = ?", id).First(ctx)
	} else {
		blog, err = query.Where("title LIKE ?", fmt.Sprintf("%%%s%%", title)).First(ctx)
	}
	if err != nil {
		s.logger.Error("unable to get blog", "id", id, "error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "blog not found")
		}
	}
	return &blog, err
}

func (s *Service) CreateBlog(ctx context.Context, blog Blog) (uint, error) {
	err := gorm.G[Blog](s.db).Create(ctx, &blog)
	if err != nil {
		s.logger.Error("unable to create blog", "id", blog.ID, "error", err)
	}
	return blog.ID, err
}

func (s *Service) UpdateBlog(ctx context.Context, blog Blog) error {
	rows, err := gorm.G[Blog](s.db).Updates(ctx, blog)
	if err != nil {
		s.logger.Error("unable to update blog", "id", blog.ID, "error", err)
	} else {
		s.logger.Info("updated", "rows", rows)
	}
	return err
}

func (s *Service) DeleteBlog(ctx context.Context, id uint) error {
	rows, err := gorm.G[Blog](s.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		s.logger.Error("unable to delete blog", "id", id, "error", err)
	} else {
		s.logger.Info("deleted", "rows", rows)
	}
	return err
}
