
# scripts/get-blogs.sh 10 2 "Id asc"
# scripts/get-blogs.sh 10 2
# scripts/get-blogs.sh 10

if [ "$#" -eq 3 ]; then
  grpcurl -plaintext \
    -d '{"limit": '"$1"', "page": '"$2"', "sort": "'"$3"'"}' \
    localhost:8080 pb.Blogger/GetBlogs
elif [ "$#" -eq 2 ]; then
  grpcurl -plaintext \
    -d '{"limit": '"$1"', "page": '"$2"'}' \
    localhost:8080 pb.Blogger/GetBlogs
else
  grpcurl -plaintext \
    -d '{"limit": '"$1"'}' \
    localhost:8080 pb.Blogger/GetBlogs
fi
