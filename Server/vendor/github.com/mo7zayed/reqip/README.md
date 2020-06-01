# reqip
A simple tool for retrieving a request's IP address on the server. Inspired from [request-ip](https://github.com/pbojinov/request-ip)

# Installation
Via `go get`

```bash
go get github.com/mo7zayed/reqip
```

# How to use
```golang
package main

import (
	"fmt"
	"net/http"
	"github.com/mo7zayed/reqip"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(
			w,
			fmt.SprintF("your ip is %s", reqip.GetClientIP(r)) // reqip.GetClientIP receives a *http.Request var type
		)
	})

	fmt.Println("server started on: http://127.0.1.1:8000")
	http.ListenAndServe(":8000", nil)
}
```

## How It Works

It looks for specific headers in the request and falls back to some defaults if they do not exist.

The user ip is determined by the following order:

1. `X-Client-IP`
2. `X-Forwarded-For` (Header may return multiple IP addresses in the format: "client IP, proxy 1 IP, proxy 2 IP", so we take the the first one.)
3. `CF-Connecting-IP` (Cloudflare)
4. `Fastly-Client-Ip` (Fastly CDN and Firebase hosting header when forwared to a cloud function)
5. `True-Client-Ip` (Akamai and Cloudflare)
6. `X-Real-IP` (Nginx proxy/FastCGI)
7. `X-Cluster-Client-IP` (Rackspace LB, Riverbed Stingray)
8. `X-Forwarded`, `Forwarded-For` and `Forwarded` (Variations of #2)
9. `http.Request.RemoteAddr`

## License
The MIT License (MIT) - 2020