FROM alpine:3.18
ARG TARGETPLATFORM
RUN apk add --no-cache ca-certificates
COPY ${TARGETPLATFORM}/creality2mqtt /creality2mqtt

ENTRYPOINT ["/creality2mqtt"]
CMD ["--help"]
