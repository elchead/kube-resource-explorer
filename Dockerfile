FROM golang:1.17-alpine AS builder
RUN apk add --no-cache git make
RUN apk add build-base

WORKDIR /code

COPY . .
RUN go mod download
RUN go build ./cmd/memreq

FROM golang:1.17-alpine
COPY --from=builder /code/memreq /

CMD ["/memreq", "--isLocal=false"]
