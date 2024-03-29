go install github.com/cloudwego/hertz/cmd/hz@latest
go install github.com/cloudwego/thriftgo@latest

hz new -module abcp -idl idl/main.thrift
# hz update -module abcp -idl idl/main.thrift

go mod tidy

bash build.sh

./output/bin/hertz_service
