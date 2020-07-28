FROM golang:1.14-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
RUN go build github.com/agilestacks/tls-host-controller/cmd/tls-host-controller


FROM alpine:3.11

RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/tls-host-controller ./
EXPOSE 4443
ENTRYPOINT [ "./tls-host-controller" ]
