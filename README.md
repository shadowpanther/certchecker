# certchecker

A simple tool to check whether the DNS challenge for the LetsEncrypt certificate is set up correctly.

To build:
```
go build -o ~/your/bin/directory/
```

```
Usage of certchecker:
certchecker [-h] [-v] [-dns 1.2.3.4] host.name...
  -dns string
        set the DNS server to query from (default "1.1.1.1")
  -h    show this help
  -v    output raw request results
```
