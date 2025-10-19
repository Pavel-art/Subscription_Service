FROM ubuntu:latest
LABEL authors="pole"

ENTRYPOINT ["top", "-b"]