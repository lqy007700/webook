FROM ubuntu:latest


RUN mkdir -p /app
COPY webook /app
WORKDIR /app
ENTRYPOINT ["/app/webook"]