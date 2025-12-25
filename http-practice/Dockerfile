# STAGE 1: The Builder
# Use the official Go image to compile your code
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy your Go files from your Mac into the container
COPY . .

# Build the binary inside the container (Notice the variables are here now!)
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# STAGE 2: The Final Image
# Start fresh with a tiny 'scratch' image
FROM scratch

# Copy the binary from the 'builder' stage to this final image
COPY --from=builder /main /main

# Run the app
EXPOSE 8080
ENTRYPOINT ["/main"]