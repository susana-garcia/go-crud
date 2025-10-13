# scripts/update-blog.sh "new blog title"

grpcurl -plaintext \
  -d '{"id": 1, "title": "'"$1"'"}' \
  localhost:8080 pb.Blogger/UpdateBlog
