FROM golang:1.26-alpine AS build
WORKDIR /src/api

COPY api/go.mod api/go.sum ./
RUN go mod download

COPY api/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server .

FROM alpine:3.21
RUN adduser -D -H appuser
COPY --from=build --chown=appuser:appuser /out/server /server

USER appuser
EXPOSE 8081
ENTRYPOINT ["/server"]