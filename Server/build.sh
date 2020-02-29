# 构建适合部署于阿里云函数计算的二进制文件
# https://help.aliyun.com/document_detail/132053.html

docker run -it --rm -v "$PWD":/tmp/code -w /tmp/code golang \
        bash -c "go build -mod=vendor -o /tmp/code/build/bootstrap /tmp/code/main.go"

# docker run -it --rm -v "$PWD":/tmp/code -w /tmp/code zhenshaw/gzhupi:builder \
#         bash -c CGO_ENABLED=1 go build  -a -ldflags '-extldflags "-static"' -mod=vendor -o /tmp/code/build/bootstrap /tmp/code/main.go

cp config.toml ./build
