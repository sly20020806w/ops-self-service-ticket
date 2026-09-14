FROM golang:1.22-alpine AS build
WORKDIR /src
RUN apk add --no-cache gcc musl-dev
COPY go.mod ./
COPY . .
RUN go mod tidy && CGO_ENABLED=1 go build -o /out/server ./cmd/server

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=build /out/server /app/server
COPY configs /app/configs
RUN mkdir -p /app/data
EXPOSE 8080
CMD ["/app/server"]
