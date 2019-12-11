# 构建适合部署于阿里云函数计算的二进制文件
# https://help.aliyun.com/document_detail/132053.html

docker run -it -v "$PWD":/tmp/code -w /tmp/code golang:1.12.9-stretch \
        bash -c "go build -mod=vendor -o /tmp/code/build/bootstrap /tmp/code/main.go"
