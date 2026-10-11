FROM --platform=$BUILDPLATFORM alpine:latest AS timezone_build
RUN apk --no-cache add tzdata ca-certificates

FROM --platform=$BUILDPLATFORM golang:1.27.1-alpine3.24 AS builder

ARG TARGETOS
ARG TARGETARCH

RUN apk --no-cache add tzdata ca-certificates

WORKDIR /go/api

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

ADD . /go/api

RUN rm -rf deploy
RUN mkdir deploy

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -o go_app cmd/main.go
RUN mv go_app ./deploy/go_app
RUN sed 's/\$SERVER_PORT/8080/g' config.example.yaml > ./deploy/config.yaml

FROM scratch AS production

COPY --from=timezone_build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=timezone_build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go/api/deploy /api/

WORKDIR /api

ENTRYPOINT ["/api/go_app"]
CMD ["serve"]
