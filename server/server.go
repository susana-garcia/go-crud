package server

import (
	"context"
	"log/slog"

	"github.com/susana-garcia/go-crud/pb"
	"github.com/susana-garcia/go-crud/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	service service.Blogger
	logger  *slog.Logger
	pb.UnimplementedBloggerServer
}

func New(service service.Blogger, logger *slog.Logger) *Server {
	return &Server{
		service: service,
		logger:  logger,
	}
}

func (s *Server) Register(server grpc.ServiceRegistrar) {
	pb.RegisterBloggerServer(server, s)
}

func (s *Server) GetBlogs(ctx context.Context, req *pb.GetBlogsRequest) (*pb.GetBlogsResponse, error) {
	pagination := service.Pagination{
		Limit: int(req.GetLimit()),
		Page:  int(req.GetPage()),
		Sort:  req.GetSort(),
	}
	sRes, err := s.service.GetAllBlogs(ctx, &pagination)
	if err != nil {
		s.logger.Error("got service error ", "error", err)
		return nil, err
	}
	var res pb.GetBlogsResponse
	for _, blog := range sRes.Items.([]service.Blog) {
		res.Items = append(res.Items, &pb.Blog{
			Id:        uint32(blog.ID),
			Title:     blog.Title,
			Body:      blog.Body,
			CreatedAt: timestamppb.New(blog.CreatedAt),
			UpdatedAt: timestamppb.New(blog.UpdatedAt),
		})
	}
	res.Limit = int32(sRes.Limit)
	res.Page = int32(sRes.Page)
	res.Sort = sRes.Sort
	res.TotalItems = int32(sRes.TotalItems)
	res.TotalPages = int32(sRes.TotalPages)
	return &res, nil
}

func (s *Server) GetBlog(ctx context.Context, req *pb.GetBlogRequest) (*pb.GetBlogResponse, error) {
	sRes, err := s.service.GetBlogByIDOrTitle(ctx, uint(req.GetId()), req.GetTitle())
	if err != nil {
		s.logger.Error("got service error ", "error", err)
		return nil, err
	}
	return &pb.GetBlogResponse{
		Item: &pb.Blog{
			Id:        uint32(sRes.ID),
			Title:     sRes.Title,
			Body:      sRes.Body,
			CreatedAt: timestamppb.New(sRes.CreatedAt),
			UpdatedAt: timestamppb.New(sRes.UpdatedAt),
		},
	}, err
}

func (s *Server) CreateBlog(ctx context.Context, req *pb.CreateBlogRequest) (*pb.CreateBlogResponse, error) {
	id, err := s.service.CreateBlog(ctx, service.Blog{
		Title: req.GetTitle(),
		Body:  req.GetBody(),
	})
	if err != nil {
		s.logger.Error("got service error ", "error", err)
		return nil, err
	}
	return &pb.CreateBlogResponse{
		Id: uint32(id),
	}, nil
}

func (s *Server) UpdateBlog(ctx context.Context, req *pb.UpdateBlogRequest) (*emptypb.Empty, error) {
	err := s.service.UpdateBlog(ctx, service.Blog{
		ID:    uint(req.GetId()),
		Title: req.GetTitle(),
		Body:  req.GetBody(),
	})
	if err != nil {
		s.logger.Error("got service error ", "error", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteBlog(ctx context.Context, req *pb.DeleteBlogRequest) (*emptypb.Empty, error) {
	err := s.service.DeleteBlog(ctx, uint(req.GetId()))
	if err != nil {
		s.logger.Error("got service error ", "error", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
