FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache gcc musl-dev sqlite-dev
RUN CGO_ENABLED=1 GOOS=linux go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o forum .

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite-libs
WORKDIR /root/
COPY --from=builder /app/forum .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/database ./database
COPY --from=builder /app/data ./data
RUN mkdir -p /root/data
EXPOSE 8080
CMD ["./forum"]