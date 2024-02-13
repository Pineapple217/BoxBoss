ARG GO_VERSION=1.21.6
FROM golang:${GO_VERSION} AS build
WORKDIR /src

ENV CGO_ENABLED=1

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    go build -ldflags='-s -w -extldflags "-static"' -o /bin/server .
    # static linking is necessary because of CGO dependency
    # -s -w removes debug info for smaller bin

FROM alpine:latest AS final

RUN echo "http://ftp.halifax.rwth-aachen.de/alpine/v3.19/main" >> /etc/apk/repositories \
    && echo "http://ftp.halifax.rwth-aachen.de/alpine/v3.19/community" >> /etc/apk/repositories \
    && apk update \
    && apk add docker-compose

WORKDIR /app

RUN mkdir ./static
COPY ./static ./static

COPY --from=build /bin/server /app/server

EXPOSE 3000

CMD [ "/app/server" ]
