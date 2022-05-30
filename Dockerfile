# Base image 
FROM golang:1.17.6-buster

# Selects everything in current directory and copies to /apicontainer
ADD ./ /apicontainer

# Move to working directory /apicontainer/server
WORKDIR /apicontainer/server

# Build the application
RUN go build -o apiapp

# Expose necessary port
EXPOSE 8080

# Command to run when starting the container
CMD ["./apiapp"]
