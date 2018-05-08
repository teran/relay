# mail-relay
SMTP server to use for forwarding messages via HTTP-based services in environments not supporting SMTP outbound directly

# How this works
relay accepts SMTP connection to handle message with invoking Mailgun API with the message it got. Simply :)

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
