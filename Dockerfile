# build final image
FROM alpine:latest
WORKDIR /
COPY cld /cld
USER 65532:65532
ENTRYPOINT ["/cld"]
