# infinitude-prometheus

Minimal web service to expose Prometheus-formatted metrics from an [Infinitude](https://github.com/nebulous/infinitude) server.

## Usage

### Docker
```sh
$ docker run \
   --name infinitude-prometheus \
   -p 8080:8080 \
   -e INFINITUDE_BASE_URL=http://infinitude:3000 \
   ghcr.io/jeremyhayes/infinitude-prometheus:v0.0.1
```

### Docker Compose
```yml
version: '3.8'

services:

  infinitude:
    image: nebulous/infinitude:latest
    ports:
      - 3000:3000

  infinitude-prometheus:
    image: ghcr.io/jeremyhayes/infinitude-prometheus:v0.0.1
    ports:
      - 8080:8080
    environment:
      - INFINITUDE_BASE_URL=http://infinitude:3000
```
