FROM golang:1.17.7 as builder

WORKDIR /workspace
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -mod=vendor -o /tmp/heartbeat cmd/heartbeat.go

FROM gcr.io/distroless/base
WORKDIR /
COPY --from=builder /tmp/heartbeat /heartbeat
ENTRYPOINT ["/heartbeat"]