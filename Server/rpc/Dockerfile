#FROM golang:latest as builder
## 安装依赖
#RUN apt-get update -qq && apt-get install -y -qq libtesseract-dev libleptonica-dev

FROM zhenshaw/golang:tesseract as builder
WORKDIR /home
ADD . /home
RUN go build -mod=vendor rpc


# 运行容器
#FROM ubuntu:latest
#RUN apt-get update -qq && apt-get install -y -qq libtesseract-dev libleptonica-dev

# 下载数据文件，并设置数据文件目录
#RUN apt-get install -y -qq  tesseract-ocr-eng
#ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/4.00/tessdata

FROM zhenshaw/ubuntu:tesseract
COPY --from=builder /home/rpc /home/
CMD ["./home/rpc"]


#docker run -it -v "$PWD":/tmp/code -w /tmp/code zhenshaw/golang:tesseract \
#        bash -c "go build -mod=vendor -o /tmp/code/build/bootstrap /tmp/code/rpc.go"