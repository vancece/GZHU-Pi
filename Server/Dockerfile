FROM golang:alpine

WORKDIR /home
ADD . /home
RUN go build --mod=vendor main.go

#FROM zhenshaw/frp
FROM alpine
COPY --from=0 /home/main /home
COPY --from=0 /home/config.toml /home
WORKDIR /home

CMD ["/home/main"]