FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct \
 && go mod tidy \
 && CGO_ENABLED=0 go build -o /out/server ./cmd/server

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates \
 && mkdir -p /app/data
COPY --from=build /out/server /app/server
ENV ADDR=:8080
ENV DB_PATH=/app/data/ticket.db
ENV SEED_DEMO=true
EXPOSE 8080
CMD ["/app/server"]
