FROM golang:1.11 as builder
WORKDIR /go/src/go-slack-bot
COPY ./ .
#RUN go build main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/go-slack-bot/app .
COPY --from=builder /go/src/go-slack-bot/.env.production.yaml .
CMD ["./app","-e production"]  

