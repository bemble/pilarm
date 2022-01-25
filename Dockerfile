# Server build
FROM golang:1.17-alpine as server-builder

RUN apk add --no-cache \
    alpine-sdk \
    ca-certificates \
    tzdata

# Force the go compiler to use modules
ENV GO111MODULE=on

ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -o pilarm .

# Final image
FROM scratch

# copy server files
COPY --from=server-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=server-builder /usr/share/zoneinfo /usr/share/zoneinfov
COPY --from=server-builder /app/pilarm /app/ressources /app/
ADD ./ressources /app/ressources

ENTRYPOINT ["/app/pilarm"]