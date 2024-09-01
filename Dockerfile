# Start from the latest Go 1.22 image
FROM golang:1.22-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project
COPY . .

# Build the application
RUN go build -o bin/big-john ./cmd

# Start a new stage from scratch
FROM alpine:latest

# Copy the binary from the build stage
COPY --from=build /app/bin/big-john /big-john

# Copy static files
COPY --from=build /app/home.html /home.html
COPY app.env .

# Set environment variables with default values
ENV LOG_LEVEL=0
ENV APP_ENV=production
ENV PORT=5001

EXPOSE 5001

# Run the binary
CMD ["/big-john"]
