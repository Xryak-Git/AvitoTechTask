# Step 1: Modules caching
FROM golang:alpine as builder
WORKDIR /usr/local/src
RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY  ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . ./
RUN go build -o ./bin/app cmd/avito-tech/main.go

FROM alpine

COPY --from=builder  /usr/local/src/bin/app /
COPY --from=builder /usr/local/src/config/ /config
ENV CONFIG_PATH="/config/local.yaml"
CMD ["/app"]