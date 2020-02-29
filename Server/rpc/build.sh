# 构建适合部署于阿里云函数计算的二进制文件
# https://help.aliyun.com/document_detail/132053.html

docker run -it --rm -v "$PWD":/tmp/code -w /tmp/code zhenshaw/golang:tesseract \
      bash -c "go build -mod=vendor -o /tmp/code/build/bootstrap /tmp/code/rpc.go"


cp config.toml ./build
