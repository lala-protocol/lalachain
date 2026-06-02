# Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR /app
COPY chain/go.mod chain/go.sum ./
RUN go mod download

COPY chain/ ./
RUN CGO_ENABLED=1 go build -o /lalachaind ./cmd/lalachaind

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates jq curl bash

# Run as non-root user
RUN addgroup -g 1000 lala && adduser -D -u 1000 -G lala lala

COPY --from=builder /lalachaind /usr/local/bin/lalachaind

USER lala

EXPOSE 26656 26657 1317 9090 9091

ENTRYPOINT ["lalachaind"]
CMD ["start"]
