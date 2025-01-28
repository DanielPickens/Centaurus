FROM golang:1.23

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/danielpickens/centarus

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .



# Run the executable
CMD ["Centarus"]