Simple HTTP proxy for domain filter
===

It's a very simple http proxy server with domain filtering.

## block list
A block list is a text file that includes a list of domains which you want to block (discard) to access.

By default, the server reads a list in conf/ directory. (see Makefile in go/ directory)


## Usage (docker)

0. install docker and docker-compose

1. build docker image
    ```bash
    docker build -t proxy_server .
    ```

2. start docker using docker-compose
    ```bash
    docker-compose up -d
    ```

3. set HTTP_PROXY and HTTPS_PROXY on your deveice to the docker container of port 8080


## Usage (command line)

0. setup golang development environment

    You also need "make" command.

1. run make command
```bash
make build
```