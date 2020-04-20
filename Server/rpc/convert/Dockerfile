FROM golang as builder

WORKDIR /home
ADD . /home
RUN go build -o app

FROM zhenshaw/libreoffice
COPY --from=builder /home/app /home/
CMD ["./home/app"]