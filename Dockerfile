# Build image

FROM golang:1.14.3-stretch AS builder
ENV GO111MODULE=on CGO_ENABLED=1
WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build

WORKDIR /dist
RUN mv /build/mem-limit ./mem-limit

# Runtime image

FROM scratch
COPY --chown=0:0 --from=builder /dist /

USER 65534
WORKDIR /
ENTRYPOINT ["/mem-limit"]
