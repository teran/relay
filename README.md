[![Go Report](https://goreportcard.com/badge/github.com/teran/relay)](https://goreportcard.com/report/github.com/teran/relay)
[![License](https://img.shields.io/github/license/teran/relay.svg)](https://github.com/teran/relay/blob/master/LICENSE)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fteran%2Frelay.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fteran%2Frelay?ref=badge_shield)

[![Layers size](https://images.microbadger.com/badges/image/teran/relay.svg)](https://hub.docker.com/r/teran/relay/)
[![Recent build commit](https://images.microbadger.com/badges/commit/teran/relay.svg)](https://hub.docker.com/r/teran/relay/)
[![Docker Automated build](https://img.shields.io/docker/automated/teran/relay.svg)](https://hub.docker.com/r/teran/relay/)
[![GoDoc](https://godoc.org/github.com/teran/relay?status.svg)](https://godoc.org/github.com/teran/relay)

# mail-relay

SMTP server to use for forwarding messages via HTTP-based services in environments not supporting SMTP outbound directly

# How this works
Relay accepts SMTP connection to handle message with invoking Mailgun API with the message it got. Simply :)
In depth in current implementation relay sends message synchronously via Mailgun Go's client.

# How to use
```
docker run -it \
  -e RELAY_ADDR=:25 \
  -e RELAY_DOMAIN="<domain>" \
  -e RELAY_MAILGUN_PRIVATE_KEY="<mailgun private key>" \
  -e RELAY_MAILGUN_PUBLIC_KEY="<mailgun public key>" \
  -e RELAY_MAX_IDLE_SECONDS=300 \
  -e RELAY_MAX_MESSAGE_BYTES=1048576 \
  -e RELAY_MAX_RECIPIENTS=50 \
  teran/relay
```


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fteran%2Frelay.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fteran%2Frelay?ref=badge_large)
