# scripts/delete-blog.sh 4

grpcurl -plaintext \
  -d '{"id": '"$1"'}' \
  localhost:8080 pb.Blogger/DeleteBlog
