ARG BASE_IMAGE=gcr.io/distroless/static-debian12:nonroot

FROM --platform=$BUILDPLATFORM golang:1.24.4-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a \
    -o random-quote \
    .

FROM ${BASE_IMAGE}

COPY --from=builder /build/random-quote /app/random-quote

WORKDIR /app

USER nonroot:nonroot

ENTRYPOINT ["/app/random-quote"]

