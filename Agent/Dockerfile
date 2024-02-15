# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang:1.21-alpine base image
FROM golang:1.21-alpine

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY . .

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./main"]