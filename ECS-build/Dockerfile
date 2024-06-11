# Use Ubuntu 20.04 (Focal Fossa) as the base image
FROM ubuntu:focal

# Update and install necessary packages
RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y wget curl apt-transport-https ca-certificates software-properties-common

# Add Docker's official GPG key and repository, then install Docker
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - \
    && add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" \
    && apt-get update \
    && apt-get install -y docker-ce

# Install Git
RUN apt-get install -y git

# Set the working directory
WORKDIR /home/app

# Copy all files from the host to the container
COPY . .

# Make the script executable
RUN chmod +x script.sh

# Check Docker version
RUN docker --version

# Export environment variables
ENV DOCKER_IMAGE_NAME=g4ze/cattodb
ENV PORT=8000

# Set the entry point to the script
ENTRYPOINT [ "/home/app/script.sh" ]
