FROM golang:1.16.10 AS build-env

WORKDIR /project

ENV GO111MODULE="on" \
    GOSUMDB=off \
    GOPROXY="https://goproxy.cn,direct"

COPY go.mod go.sum ./
RUN time go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    && go build -o main

FROM golang:1.16.10

WORKDIR /app
RUN apk add curl
COPY --from=build-env /project/main /app/

CMD ["/app/main"]