FROM golang:1.25 AS build-stage

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app cmd/main.go


# Deploy the application binary into a lean image
#FROM gcr.io/distroless/base-debian11 AS build-release-stage
FROM debian:11-slim AS build-release-stage

WORKDIR /

COPY migrations /migrations

COPY --from=build-stage /app /app

#USER nonroot:nonroot

ENTRYPOINT ["/app"]