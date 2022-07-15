# Server build
FROM golang:1.17-alpine as server-builder

RUN apk add --no-cache \
    alpine-sdk

# Force the go compiler to use modules
ENV GO111MODULE=on

ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -o pilarm .

# Final image
FROM alpine:3.16

RUN apk add --no-cache \
    ca-certificates \
    alsa-utils \
    tzdata

# copy server files
COPY --from=server-builder /app/pilarm /app/ressources /app/
ADD ./ressources /app/ressources

ENTRYPOINT ["/app/pilarm"]