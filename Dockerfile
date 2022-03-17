FROM golang:1.17.8 as builder
WORKDIR /workspace
COPY . /workspace
RUN CGO_ENABLED=0 GOOS=linux go build -o main && chmod +x ./main

FROM alpine:3.15
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /workspace/secrets/ ./secrets/
COPY --from=builder /workspace/main ./
CMD ["./main"]
