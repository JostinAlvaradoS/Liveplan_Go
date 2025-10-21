# Dockerfile para backend Go
FROM golang:1.21-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o liveplan_backend_go main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/liveplan_backend_go ./
COPY --from=build /app/docker-compose.yml ./docker-compose.yml
EXPOSE 8080
CMD ["./liveplan_backend_go"]
