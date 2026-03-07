FROM alpine:3.23.3 AS certificates

RUN apk add --update --no-cache \
  ca-certificates=20251003-r0

FROM scratch

LABEL org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/teran/relay"

COPY --from=certificates /etc/ssl/cert.pem /etc/ssl/cert.pem
COPY --chmod=0755 --chown=root:root dist/relay_linux_amd64_v3/relay /relay

ENTRYPOINT [ "/relay" ]
