FROM golang:latest as builder

WORKDIR /home
ADD . /home

# 安装依赖
RUN apt-get update -qq && apt-get install -y -qq libtesseract-dev libleptonica-dev
RUN go build -mod=vendor rpc


# 运行容器
FROM ubuntu
COPY --from=builder /home/rpc /home/
RUN apt-get update -qq && apt-get install -y -qq libtesseract-dev libleptonica-dev

# 下载数据文件，并设置数据文件目录
RUN apt-get install -y -qq  tesseract-ocr-eng
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/4.00/tessdata

CMD ["./home/rpc"]