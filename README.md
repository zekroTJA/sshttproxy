# sshttproxy

sshttproxy is an SSH subsystem that allows proxying HTTP requests through SSH to another HTTP server.

## Why?

Using SSH subsystems instead of hosting a separate SSH server has two main advantages:

- You can simply use the default SSH configuration and credentials for user authentication.
- Authentication is done though the actual SSH server instead of through a custom implementation.

But it also has some disadvantages:

- For each SSH connection, a new instance of the service is spawned instead of using a single service instance for all incomming connections.
- For implementation convenience, all requests are first parsed by the Go HTTP server and then relayed to the upstream using the Go's HTTP client instead of passing the requests directly via a TCP tunnel. This introduces some overhead on each request proxied.

If you are looking for a solution that hosts its own SSH server, feel free to take a look at the [sshified project](https://github.com/hoffie/sshified).

## Setup

Simply downlaod [the latest release](https://github.com/zekroTJA/sshttproxy/releases) to your server.

After that, add the following to your sshd config (commonly found at `/etc/ssh/sshd_config`).
```
Subsystem http /usr/sbin/http-subsystem-proxy /etc/sshttproxy.env
```

After that, create a config file at `/etc/sshttproxy.env`. There you can define the log leve, a log file directory as well as the target upstream HTTP server to proxy the requests to.
```bash
SSHTTPROXY_TARGET="http://localhost:8080"
SSHTTPROXY_LOGFILE="/var/log/proxy.log"
SSHTTPROXY_LOGLEVEL="info"
```

After that, restart the sshd service.
```bash
sudo systemctl restart sshd
```

## Client

To connect to the proxy, you can use the client implementation in [pkg/client](pkg/client). In [examples/client](examples/client), you can find a simple example on how to use the client implementation.
