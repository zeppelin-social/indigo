# Run this dockerfile from the top level of the indigo git repository like:
#
#   podman build -f ./cmd/relay/Dockerfile -t relay .

### Compile stage
FROM golang:1.24-alpine3.22 AS build-env
RUN apk add --no-cache build-base make git

ADD . /dockerbuild
WORKDIR /dockerbuild

# timezone data for alpine builds
ENV GOEXPERIMENT=loopvar
RUN GIT_VERSION=$(git describe --tags --long --always) && \
    go build -tags timetzdata -o /relay ./cmd/relay

### Build Frontend stage
FROM node:18-alpine as web-builder

WORKDIR /app

COPY cmd/relay/relay-admin-ui /app/

RUN yarn install --frozen-lockfile

RUN yarn build

### Run stage
FROM alpine:3.22

RUN apk add --no-cache --update dumb-init ca-certificates runit
ENTRYPOINT ["dumb-init", "--"]

WORKDIR /
RUN mkdir -p data/relay
COPY --from=build-env /relay /
COPY --from=web-builder /app/dist/ public/

# small things to make golang binaries work well under alpine
ENV GODEBUG=netdns=go
ENV TZ=Etc/UTC

EXPOSE 2470

CMD ["/relay", "serve"]

LABEL org.opencontainers.image.source=https://github.com/bluesky-social/indigo
LABEL org.opencontainers.image.description="atproto Relay"
LABEL org.opencontainers.image.licenses=MIT
