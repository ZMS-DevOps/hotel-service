FROM golang:1.22.1-alpine AS builder
RUN apk add --no-progress --no-cache gcc musl-dev
WORKDIR /build
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .

RUN go build -tags musl -ldflags '-extldflags "-static"' -o /build/main

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /build/main .
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown appuser:appgroup /app/main && chmod 555 /app/main

RUN mkdir -p /app/uploads
RUN chown -R appuser:appgroup /app/uploads && chmod -R 0777 /app/uploads

USER appuser
ENTRYPOINT ["/app/main"]