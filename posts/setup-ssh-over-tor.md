# Setup SSH over TOR
---
title: "Setup SSH Over TOR"
date: 2023-07-02T08:54:00+02:00
toc: true
tags: ["tor", "tech"]
---

This is tutorial for Debian systems.

## Setting up Tor - Server

First install tor.
```bash
sudo apt install tor
```

Then edit torrc file.
```bash
sudo vim /etc/tor/torrc
```
Add this to it
```
HiddenServiceDir /var/lib/tor/hidden-service-example/
HiddenServicePort 22 127.0.0.1:22
HiddenServiceAuthorizeClient stealth hidden-service-example
```
Restart tor
```bash
sudo systemctl restart tor
```

That's it for the server


## Setting up Tor - Client

First install tor.
```bash
sudo apt install tor
sudo systemctl start tor
```

Then ssh into your machine
```bash
torify user@hostname_of_your_service
```

To get hostname run this.
```bash
sudo cat /var/lib/tor/hidden-service-example/hostname
```

