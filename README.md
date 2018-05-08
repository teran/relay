# mail-relay
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fteran%2Frelay.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fteran%2Frelay?ref=badge_shield)

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


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fteran%2Frelay.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fteran%2Frelay?ref=badge_large)