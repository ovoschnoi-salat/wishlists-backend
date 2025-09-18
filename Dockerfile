FROM golang:1.24 AS build-stage

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

RUN go install tool

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app cmd/main.go


# Deploy the application binary into a lean image
#FROM gcr.io/distroless/base-debian11 AS build-release-stage
#FROM gcr.io/distroless/base-debian11 AS build-release-stage
FROM debian:11-slim as build-release-stage

WORKDIR /

COPY entrypoint.sh /
COPY migrations /migrations
COPY docs /docs


COPY --from=build-stage /app /app

RUN ls -la
#USER nonroot:nonroot

ENTRYPOINT ["/entrypoint.sh"]

CMD ["/app"]