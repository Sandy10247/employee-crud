
# Build the application from source
FROM golang:1.25.5-bookworm AS build-stage

WORKDIR /app

ENV HOSTNAME=0.0.0.0

COPY . .

WORKDIR /app/server

RUN go mod download

RUN ls 

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /server /app/cmd/main.go

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

# Copy Built Server Binary
COPY --from=build-stage ./server /server

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/server"]