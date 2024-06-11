# mail-relay

[![Verify](https://github.com/teran/relay/actions/workflows/verify.yml/badge.svg?branch=master)](https://github.com/teran/relay/actions/workflows/verify.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/teran/relay)](https://goreportcard.com/report/github.com/teran/relay)
[![Go Reference](https://pkg.go.dev/badge/github.com/teran/relay.svg)](https://pkg.go.dev/github.com/teran/relay)

SMTP server to use for forwarding messages via HTTP-based services in environments
not supporting SMTP outbound directly

## How this works

Relay accepts SMTP connection to handle message with invoking Mailgun API with
the message it got. Simply :)
In depth in current implementation relay sends message synchronously via Mailgun
Go's client.

## How to use

```shell
docker run -it \
  -e RELAY_MAILGUN_API_KEY="<MAILGUN_API_KEY> \
  -e RELAY_DOMAIN="<domain>" \
  -e RELAY_MAILGUN_URL="<url>" \
  -e RELAY_ADDR=:25 \
  -e RELAY_MAX_IDLE_SECONDS=300 \
  -e RELAY_MAX_MESSAGE_BYTES=1048576 \
  -e RELAY_MAX_RECIPIENTS=50 \
  teran/relay
```
