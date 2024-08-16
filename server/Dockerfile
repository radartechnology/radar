FROM sqlc/sqlc:1.27.0 AS sqlc

WORKDIR /app

COPY migrations ./migrations
COPY query.sql sqlc.yml ./

RUN ["/workspace/sqlc", "generate"]

FROM golang:1.22.6-alpine3.20 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY --from=sqlc /app/db /app/db

RUN go build -o /server

FROM scratch

COPY migrations /migrations
COPY --from=build /server /server

ENTRYPOINT ["/server"]

ENV WEBSOCKET_PORT=1887
EXPOSE $WEBSOCKET_PORT
