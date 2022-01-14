FROM golang:1.17-alpine AS builder
RUN apk add --no-cache git make

WORKDIR /code

COPY . .
RUN go mod download
RUN make build

FROM scratch
COPY --from=builder /code/out/kube-resource-explorer /

ENTRYPOINT ["/kube-resource-explorer"]
