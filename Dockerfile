# Use a official Go image as a base
FROM golang:1.20

# Set the working directory inside the Docker image
WORKDIR /app

# Copy Go code to the container
COPY . .

# Download Go dependencies
RUN go mod download

# Compile the application
RUN go build -o main .

# Expose the port the app listens on
EXPOSE 8080

# Command to run when starting the Docker container
ENTRYPOINT ["./main"]