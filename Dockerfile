FROM golang:alpine
RUN apk add ca-certificates git

RUN mkdir -p /root/src/go
WORKDIR /root/src/go
COPY . .
RUN go mod download

EXPOSE 80

ENTRYPOINT ["go","run","main.go"]