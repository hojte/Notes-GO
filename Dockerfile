# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Mathias <mhoo@itu.dk>"

RUN mkdir /build

ADD . /build

# Set the Current Working Directory inside the container
WORKDIR /build

RUN go build -o main .

# Run the app
ENTRYPOINT [ "/build/main" ]
