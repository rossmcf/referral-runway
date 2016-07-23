#!/bin/bash

function success() {
    prompt="$1"
    echo -e -n "\033[1;32m$prompt"
    echo -e -n '\033[0m'
    echo -e -n "\n"
}
function error() {
    prompt="$1"
    echo -e -n "\033[1;31m$prompt"
    echo -e -n '\033[0m'
    echo -e -n "\n"
}
function info() {
    prompt="$1"
    echo -e -n "\033[1;36m$prompt"
    echo -e -n '\033[0m'
    echo -e -n "\n"
}

godep get
godep save
go test

#Run Docker ps to make sure that docker is installed
#As well as that the Daemon is connected.
docker ps &>/dev/null
if [ $? -gt 0 ]; then
    error "Docker is either not installed, or the Docker Daemon is not currently connected."
    exit 1
fi
if [ "$DOCKER_HOST" ]; then
    DOCKERHOST=$(docker-machine ip $(docker-machine ls --filter state=Running -q)) >/dev/null
else
    DOCKERHOST=localhost
fi

APPLICATION=rr

# make binary directory.
mkdir bin/ >/dev/null 2>/dev/null
# delete previously generated binary.
rm bin/$APPLICATION 2>/dev/null

info "Compiling application.."
$(env GOOS=linux go build -o bin/$APPLICATION)
if [ $? -gt 0 ]; then
  error "Error building application"
  exit 42
fi
success "OK"

#info "Starting environment.."
#docker-compose -f compose-local.yml up -d --build
#if [ $? -gt 0 ]; then
#  error "Error composing"
#  exit 42
#fi
#success "OK"
exit 0
