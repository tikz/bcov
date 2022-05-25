FROM golang:alpine
RUN apk update && apk add --no-cache git gcc musl-dev nodejs npm
WORKDIR /bcov
COPY . .

ENV GO111MODULE=on
RUN go mod tidy
RUN go build

WORKDIR /bcov/web
RUN npm install --legacy-peer-deps
RUN npm run build

WORKDIR /bcov
ENTRYPOINT ["/bcov/bcov"]