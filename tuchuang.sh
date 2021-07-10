export GOROOT=/usr/local/go
export GO111MODULE=on

cd /go-tuchuang

git pull

kill $(ps axu | grep "/tmp/go.*/exe/tuchuang.*" | grep -v grep | awk '{print $2}')
nohup go run src/tuchuang.go -port=8089 1>>/var/log/go_tuchuang.log 2>&1 &
echo "Inspur Check Server run on Port \
$(ps aux | grep "go run.*tuchuang.go" | grep -v grep | awk '{print $2}')"

echo "========================================================================="
echo ""
echo "server log here, press ctrl + c to exit (the server running normally) "
echo ""
echo "========================================================================="