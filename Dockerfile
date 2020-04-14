FROM golang:alpine as builder
RUN apk add ca-certificates git

RUN mkdir -p /root
WORKDIR /root

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine

#RUN apk add ca-certificates rsync openssh

WORKDIR /root

COPY --from=builder /root/promotion-management-api /root/promotion-management-api
COPY .env .env
COPY firebase.json firebase.json

EXPOSE 80

ENTRYPOINT ["./promotion-management-api"]

# sudo docker build -t swd391 .
# sudo docker run -dit -p 8081:80 swd391:latest