# scripts/create-blog.sh "new blog"

grpcurl -plaintext \
  -d '{"title": "'"$1"'", "body": "some body"}' \
  localhost:8080 pb.Blogger/CreateBlog
