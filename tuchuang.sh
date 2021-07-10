export GOROOT=/usr/local/go
export GO111MODULE=on

cd /go-tuchuang/src

git pull

kill $(ps axu | grep "/tmp/go.*/exe/tuchuang.*" | grep -v grep | awk '{print $2}')
nohup go run tuchuang.go 1>>/var/log/go_tuchuang.log 2>&1 &
echo "goTuchuang run on Port \
$(ps aux | grep "go run.*tuchuang.go" | grep -v grep | awk '{print $2}')"

echo "========================================================================="
echo ""
echo "server log here, press ctrl + c to exit (the server running normally) "
echo ""
echo "========================================================================="