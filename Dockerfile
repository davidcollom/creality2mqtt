FROM alpine:3.24
ARG TARGETPLATFORM
RUN apk add --no-cache ca-certificates
COPY ${TARGETPLATFORM}/creality2mqtt /creality2mqtt

ENTRYPOINT ["/creality2mqtt"]
CMD ["--help"]
