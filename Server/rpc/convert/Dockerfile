FROM golang as builder

WORKDIR /home
ADD . /home
RUN go build -o app


# build zhenshaw/libreoffice
#FROM ubuntu:18.04
#
## 拷贝windows字体
#ADD path-to-win-fonts/* /usr/share/fonts/
#
#RUN apt update && apt install -y language-pack-zh* chinese* libreoffice

FROM zhenshaw/libreoffice
COPY --from=builder /home/app /home/
CMD ["./home/app"]