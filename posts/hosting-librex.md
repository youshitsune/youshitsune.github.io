# Hosting LibreX
## What is LibreX?
LibreX is meta-search engine which gives you results from Google, Qwant, Ahmia and popular torrent sites without spying on you.
LibreX doesn't save any type of data about the user, there are no logs, no caches.

## How to host LibreX in Docker container?
To host librex in a Docker container, just use this command:
```bash
docker run -d --name librex \
    -e TZ="America/New_York" \
    -e CONFIG_GOOGLE_DOMAIN="com" \
    -e CONFIG_GOOGLE_LANGUAGE="en" \
    -p 8080:8080 \
    librex/librex:latest
```

## How to host LibreX with Docker compose?
This is your **docker-compose.yml**
```yaml
version: "2.1"
services:
  librex:
    image: librex/librex:latest
    container_name: librex
    network_mode: bridge
    ports:
      - 8080:8080
    environment:
      - PUID=1000
      - PGID=1000
      - VERSION=docker
      - TZ="America/New_York"
      - CONFIG_GOOGLE_DOMAIN="com"
      - CONFIG_GOOGLE_LANGUAGE="en"
    volumes:
      - ./nginx_logs:/var/log/nginx
      - ./php_logs:/var/log/php7
    restart: unless-stopped
```
