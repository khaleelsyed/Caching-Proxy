# Caching-Proxy

The following caching proxy will forward all GET requests to the defined destination host (defined with the command line argument `destination`). Requests using any other method (`POST`, `PUT`, ...) will be returned with a `405` status.

For help on running the program (Using TaskFile)

```bash
task run -- -h
```

Or run the following commands:

```bash
go build -o .bin/caching-proxy
./caching-proxy -h
```
