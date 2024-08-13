FROM golang:1.22.6-alpine3.20 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0
RUN sqlc generate

RUN go build -o /server

FROM scratch

COPY migrations /migrations
COPY --from=build /server /server

ENTRYPOINT ["/server"]

ENV WEBSOCKET_PORT=1887
EXPOSE $WEBSOCKET_PORT
