# netw
Small command-line utility for getting data from [bgp.tools](https://bgp.tools) and other sources

![Demo](/demo.gif)

# Features
## Implemented
- Getting your IP (also the prefix and autonomous system)
## Planned
- Getting information about an autonomous system (class, tags, prefixes list, whois data)
- Getting information about an IP range

# Install
Currently, only installing with Go is supported:
```sh
go install github.com/laptopcat/netw/cmd/netw@latest
```
