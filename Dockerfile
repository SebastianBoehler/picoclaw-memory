FROM golang:1.21-alpine AS build

WORKDIR /src

ARG TARGETARCH

COPY go.mod ./
RUN go mod download

COPY . .

RUN target_arch="${TARGETARCH:-amd64}" && \
    CGO_ENABLED=0 GOOS=linux GOARCH="$target_arch" go build -o /out/picoclaw-memory ./cmd/server

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=build /out/picoclaw-memory /app/picoclaw-memory

ENV LISTEN_ADDR=:8080
ENV DATA_DIR=/data
ENV SQLITE_PATH=/data/memory.db

VOLUME ["/data"]

EXPOSE 8080

ENTRYPOINT ["/app/picoclaw-memory"]
