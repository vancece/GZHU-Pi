package reqip

import (
	"net"
	"net/http"
	"strings"
)

// IsIP : Check if the given ip address is valid
func isIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// GetClientIPFromXForwardedFor : Parse x-forwarded-for headers.
func getClientIPFromXForwardedFor(header string) string {
	if header == "" {
		return ""
	}

	// x-forwarded-for may return multiple IP addresses in the format
	// @see https://en.wikipedia.org/wiki/X-Forwarded-For#Format
	proxies := strings.Split(header, ", ")

	var ips []string

	if len(proxies) > 0 {
		for _, proxy := range proxies {
			ip := proxy
			// make sure we only use this if it's ipv4 (ip:port)
			if strings.Contains(ip, ":") {
				splitted := strings.Split(ip, ":")
				ips = append(ips, splitted[0])
				continue
			}
			ips = append(ips, ip)
		}
	}

	// Sometimes IP addresses in this header can be 'unknown' (http://stackoverflow.com/a/11285650).
	// Therefore taking the left-most IP address that is not unknown
	// A Squid configuration directive can also set the value to "unknown" (http://www.squid-cache.org/Doc/config/forwarded_for/)
	for _, ip := range ips {
		if isIP(ip) {
			return ip
		}
	}

	return ""
}

// GetClientIP : Parse all headers.
func GetClientIP(r *http.Request) string {
	headers := r.Header

	if len(headers) > 0 {
		checklist := []string{
			"x-client-ip", // Standard headers used by Amazon EC2, Heroku, and others.
			"x-forwarded-for", // Load-balancers (AWS ELB) or proxies.
			"cf-connecting-ip", // @see https://support.cloudflare.com/hc/en-us/articles/200170986-How-does-Cloudflare-handle-HTTP-Request-headers-
			"fastly-client-ip", // Fastly and Firebase hosting header (When forwared to cloud function)
			"true-client-ip", // Akamai and Cloudflare: True-Client-IP.
			"x-real-ip", // Default nginx proxy/fcgi; alternative to x-forwarded-for, used by some proxies.
			"x-cluster-client-ip", // (Rackspace LB and Riverbed's Stingray) http://www.rackspace.com/knowledge_center/article/controlling-access-to-linux-cloud-sites-based-on-the-client-ip-address
			"x-forwarded",
			"forwarded-for",
			"forwarded",
		}

		for _, h := range checklist {
			if h == "x-forwarded-for" {
				if ip := getClientIPFromXForwardedFor(r.Header.Get(h)); isIP(ip) {
					return ip
				}
				continue
			}

			if ip := r.Header.Get(h); isIP(ip) {
				return ip
			}
		}
	}

	if ip := r.RemoteAddr; isIP(ip) {
		return ip
	}

	return ""
}
