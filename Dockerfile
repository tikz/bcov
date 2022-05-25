FROM golang:1.18-alpine as builder
RUN apk update && apk add --no-cache git gcc musl-dev nodejs npm make

WORKDIR /bcov
COPY . .
RUN make


FROM scratch
COPY --from=builder /bcov/bcov /bcov
COPY --from=builder /bcov/web/build /web/build
COPY --from=builder /bcov/test.db /test.db

ENTRYPOINT ["/bcov", "-web"]