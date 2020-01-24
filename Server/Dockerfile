FROM golang:alpine

WORKDIR /home
ADD . /home
RUN go build --mod=vendor main.go

FROM alpine
COPY --from=0 /home/main /home
COPY --from=0 /home/config.yaml /home
WORKDIR /home

CMD ["/home/main"]