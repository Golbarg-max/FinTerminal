# syntax=docker/dockerfile:1
FROM golang:1.26.3 AS build

WORKDIR /app 
# Copy dependency manifests first to leverage Docker caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# copy all go files
COPY . .

# Build the Go app as a static binary
# CGO_ENABLED=0 prevents dynamic linking to C libraries
# GOOS=linux ensures the binary runs natively on Linux containers
RUN CGO_ENABLED=0 GOOS=linux go build -o /FINTERMINAL .

FROM alpine

# Copy the compiled binary from the build stage
COPY --from=build /FINTERMINAL /FINTERMINAL

# Copy the static frontend (it lived at /app/web in the build stage)
COPY --from=build /app/web /app/web

# Tell the server where the frontend is (absolute path, so the working dir doesn't matter)
ENV FINTERMINAL_WEB_DIR=/app/web

# Expose the application port (change to match your app's port)
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/FINTERMINAL"] 
CMD ["server"]
