2023/07/22
# Configuring weechat proxy to use tor
If you are lover of privacy like I am, you'll love this configuration.

### TOR

Firstly, you need to start tor.
```bash
sudo systemctl start tor
```

Then add proxy with this command
```bash
/proxy add tor socks5 127.0.0.1 9050
```

If you want to setup when you connect to any irc server to use proxy, execute this command
```bash
/set irc.server_default.proxy "tor"
```
But, if you want for specific server run this
```bash
/set irc.server.server_name.proxy "tor"
```
Of course you should change server_name with server name.


Tor and ssl_verify don't like each other so you must turn ssl_verify off.
```bash
/set irc.server.server_name.ssl_verify off
```
You can also set this for every irc server, but I don't recommend it.
```bash
/set irc.server_default.ssl_verify off
```

### Enhance your privacy
If you want to enhance your privacy run this commands.
```bash
/set irc.server_default.msg_part ""
/set irc.server_default.msg_quit ""
/set irc.ctcp.clientinfo ""
/set irc.ctcp.finger ""
/set irc.ctcp.source ""
/set irc.ctcp.time ""
/set irc.ctcp.userinfo ""
/set irc.ctcp.version ""
/set irc.ctcp.ping ""
/plugin unload xfer
/set weechat.plugin.autoload "*,!xfer"
```

After all this you just need to run one more command
```bash
/save
```

And that's it, you have irc tor proxy.
