FROM golang:1.25-alpine AS build
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN apk add --no-cache ca-certificates git tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o payment-gateway .

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /app/payment-gateway /app/payment-gateway
COPY --from=build /app/config /app/config
EXPOSE 5003 6003 5103
USER nonroot:nonroot
ENTRYPOINT ["/app/payment-gateway"]
