FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/agrismart-api ./cmd/api

FROM alpine:3.20
RUN adduser -D -H -u 10001 agrismart
USER agrismart
COPY --from=build /out/agrismart-api /agrismart-api
EXPOSE 8080
ENTRYPOINT ["/agrismart-api"]
