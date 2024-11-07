# Use a valid and supported Golang version
FROM golang:1.20


COPY ./netrc /root/.netrc
RUN chmod 600 /root/.netrc

WORKDIR /app

RUN git config --global user.email "sim.lesueur81@gmail.com" && git config --global user.name "Camanar"
# Create a directory to store the SSH key

COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the rest of the application to the container
COPY . .

# Build the Go application
RUN go build -o app .

# Expose the port that your micro-service listens on
EXPOSE 1170

# Define the command to run your micro-service
CMD ["./app"]
