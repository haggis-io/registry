FROM scratch
EXPOSE 8080
ADD bin/linux/amd64/registry /registry
ENTRYPOINT ["/registry"]
VOLUME ["/tmp"]