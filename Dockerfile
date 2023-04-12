FROM golang:1.19-alpine3.17 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /app/s3d .

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/s3d .
CMD ["/app/s3d"]