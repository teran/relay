# mail-relay
SMTP server to use for forwarding messages via HTTP-based services in environments not supporting SMTP outbound directly

# How this works
relay accepts SMTP connection to handle message with invoking Mailgun API with the message it got. Simply :)

# How to use
```
docker run -it \
  -e RELAY_DOMAIN="<domain>" \
  -e MAILGUN_PRIVATE_KEY="<mailgun private key>" \
  -e MAILGUN_PUBLIC_KEY="<mailgun public key>" \
  teran/mail-relay
```
