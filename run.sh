if [ ! -d "conf/secret/rsa" ]; then
  mkdir conf/secret/rsa
  openssl genrsa -out conf/secret/rsa/rsa_private_key.pem 1024
  openssl rsa -in conf/secret/rsa/rsa_private_key.pem -pubout -out conf/secret/rsa/rsa_public_key.pem
fi

if [ "$1" == '--s' ]; then
    # 后台启动
    go run cmd/main.go -log_dir=../log -alsologtostderr > out.out 2>&1 &
else 
    # 终端启动
    go run cmd/main.go
fi