FROM golang:alpine
RUN apk add ca-certificates git

RUN mkdir -p /root/src/go
WORKDIR /root/src/go
COPY . .
RUN go mod download

EXPOSE 80

ENTRYPOINT ["go","run","main.go"]

# sudo docker build -t swd391/dev -f dev.Dockerfile .
# sudo docker run -dit -p 80:80 swd391/dev:latest