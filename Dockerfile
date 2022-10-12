# Build from Go and Webpack sources
FROM golang:1.19-alpine as builder
RUN apk update && apk add --no-cache git gcc musl-dev nodejs npm make binutils-gold

WORKDIR /bcov
COPY . .
RUN make


# Copy binary, transpiled TS and assets into a scratch container containing just the essentials, production ready
FROM scratch
COPY --from=builder /bcov/bcov /bcov
COPY --from=builder /bcov/web/build /web/build

ENTRYPOINT ["/bcov", "-web"]