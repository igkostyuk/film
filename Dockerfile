# build stage
ARG ALPINE_VERSION=3.16

FROM golang:1.18-alpine${ALPINE_VERSION} AS builder

ARG SRC=/build/
COPY . ${SRC}
WORKDIR ${SRC}

RUN go mod download && go mod verify
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/film-service cmd/main.go

# final stage
FROM alpine:${ALPINE_VERSION}
COPY --from=builder /bin/ /usr/local/bin/

EXPOSE 9090
CMD [ "film-service" ]