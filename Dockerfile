FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/neighbor-relay ./cmd/server
FROM debian:bookworm-slim
RUN useradd --system --uid 10001 app && mkdir -p /data && chown app:app /data
WORKDIR /app
COPY --from=build /out/neighbor-relay /app/neighbor-relay
COPY --from=build /src/migrations /app/migrations
USER app
ENV ADDR=:8080 DB_PATH=/data/neighbor-relay.db
EXPOSE 8080
HEALTHCHECK --interval=5s --timeout=2s --retries=10 CMD /app/neighbor-relay healthcheck || exit 1
ENTRYPOINT ["/app/neighbor-relay"]
