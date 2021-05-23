# Simple Mail API server

This is a simple SMTP relay server written in Go which can be deployed in a public setting.

## Features

- Sends e-mails via a SMTP server
- Allows rate limiting on a IP basis
- Has a built-in allow list which contact e-mails may be used
- Offers a way to use a secret based authentication
- Set optional 'Reply-To' header of contact forms

## Spam protection

Deploying a public mail server opens up a lot of possible ways to spam. This server combats this in a few different ways.

### Rate limiting

The server rate limits IPs by only allowing a request every `RATE_LIMIT_SECONDS` seconds. You can disable this by setting `RATE_LIMIT_SECONDS` to zero.

Only the `/send` api is rate limited.

### Allow list

The `to` field of the send mail request may only contain e-mails on the `ALLOW_LIST`, all others will be denied.

That way only you might be spammed.

This is always enabled.

### Nonce and secret

The last one tries to protect the api by utilizing a secret - on a public website this of course is not bullet proof.

If enabled (via `REQUIRE_NONCE`) the client has to retrieve a "nonce" via `/nonce` before sending an e-mail.
The nonce has to be send with the e-mail request together with a sha256 hash of `{NONCE}{SECRET}`.

A pure secret based authentication does not work since api requests are public when done via your web browser.

The secret needs to be difficult to access, you might want to use environment variables when building and minify your code.

## Deployment

The server is available as a [docker image](https://hub.docker.com/r/inheadendev/simple-mail-relay).

Otherwise simply run

```
cp .env.example .env
make
./service
```

## Configuration

```
# Port the server will listen on
PORT=8000

# A list of emails (comma separated) - required
ALLOW_LIST=test@example.com,test2@example.com

# URL of your SMTP server - required
SMTP_URL=smtp.example.com
# Port your SMTP server listens to
SMTP_PORT=465
# The email address to use in the "From" field
SMTP_FROM=your@email.com
# Authentication for SMTP server - required
SMTP_USERNAME=your@email.com
SMTP_PASSWORD=XXX

# Log level (debug,info,warn,error)
LOG_LEVEL=debug

# Header where the ip can be found (e.g. x-real-ip)
IP_HEADER=

# Wheather to require the nonce hash
REQUIRE_NONCE=true

# How often one ip may send mail requests
RATE_LIMIT_SECONDS=60

# How long we keep ips and nonces around
MAX_AGE_SECONDS=60

# The secret to use for authentication - required
SECRET=secret
```

## API

### `POST /send`

Main method, will send an email to the specified address.

#### Request

```json
// headers
{
  "Content-Type": "application/json"
}

// body
{
  "to": "email-to-send-mail-to@example.com",
  "subject": "Your subject",
  "body": "Your message",
  "nonce": "Nonce (optional)",
  "hash": "Hash (optional)",
  "from": "email-to-reply-to@example.com (optional)"
}
```

#### Response

- `204`: Success, message was sent
- `400`: Bad Request, body is not valid
- `403`: Forbidden, address is not allowed or hash is not correct
- `429`: Too many requests (IP was rate limited)
- `500`: Internal error

### `GET /nonce`

Returns a new random nonce.

#### Response

- `200`: Succes

```json
{
  "nonce": "The nonce"
}
```

### `GET /health`

Health check.

#### Response

- `204`: Success, service is running

## Development

The server is written in Go 1.16.

```sh
cp .env.example .env
go run main.go
```

## License

This project is licensed under [MIT](/LICENSE.md).
