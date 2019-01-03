# Use golang base image
FROM golang:latest

# Make the working directory for the executable
RUN mkdir /app

# Add app as the working directory
ADD . /app/

# Declare the working directory
WORKDIR /app

# Expose the port
EXPOSE 8000

# Import MS Sql Driver
RUN go get github.com/denisenkom/go-mssqldb

# Import Gorilla mux
RUN go get github.com/gorilla/mux

# Import respond package
RUN go get gopkg.in/matryer/respond.v1

# Build the go executable from the main file
RUN go build main.go

# Give the executable run permission
RUN ["chmod", "+x", "main"]

# Run the executable
CMD ["./main"]
