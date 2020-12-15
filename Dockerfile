# Get latest golang docker image.
FROM golang:latest
 
# Create a directory for the app
RUN mkdir /app
 
# Copy all files from the current directory to the app directory
COPY . /app
 
# Set working directory
WORKDIR /app
 
# Run command as described:
# go build will build an executable file named server in the current directory
RUN go build -o server . 

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the server executable
CMD [ "/app/server" ]