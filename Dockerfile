FROM golang:1.6-alpine
COPY bin/rr /
ENTRYPOINT ["/rr"]
