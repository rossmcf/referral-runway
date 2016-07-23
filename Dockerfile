FROM golang:1.6-alpine
EXPOSE 8080
COPY bin/rr /
ADD index.html /
ENTRYPOINT ["/rr"]
