FROM golang:1.21.0-alpine AS builder
RUN apk update && apk upgrade \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates 2>/dev/null || true
RUN apk add --no-cache make git gcc g++ libc-dev
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
WORKDIR /build
ADD go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags "-s -w" -o aoa .

FROM alpine:3.18.0
RUN apk update && apk upgrade \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates 2>/dev/null || true
WORKDIR /app
EXPOSE 8080
COPY --from=builder /build/aoa .
ENTRYPOINT ["/app/aoa"]