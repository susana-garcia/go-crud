
# scripts/get-blog.sh "new blog"
# scripts/get-blog.sh 1

re='^[0-9]+$' # number -> for id
if ! [[ $1 =~ $re ]] ; then
  echo " search by title"
  grpcurl -plaintext \
    -d '{"title": "'"$1"'"}' \
    localhost:8080 pb.Blogger/GetBlog
else
  echo "search by id"
  grpcurl -plaintext \
    -d '{"id": '"$1"'}' \
    localhost:8080 pb.Blogger/GetBlog
fi
