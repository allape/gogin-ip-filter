package ip_filter

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/netip"
	"os"
	"slices"
	"strings"
)

func SetupIPFilter(okIPSet []string) (prefixes []netip.Prefix, hosts []string) {
	for _, ip := range okIPSet {
		prefix, err := netip.ParsePrefix(ip)
		if err != nil {
			hosts = append(hosts, ip)
			continue
		}
		prefixes = append(prefixes, prefix)
	}
	return prefixes, hosts
}

func CanPass(addr string, prefixes []netip.Prefix, hosts []string) (passed bool, prefix *netip.Prefix, host string, err error) {
	addr = strings.TrimSpace(addr)
	if slices.Contains(hosts, addr) {
		return true, nil, addr, nil
	}
	for _, prefix := range prefixes {
		ip, err := netip.ParseAddr(addr)
		if err != nil {
			return false, nil, "", err
		} else if prefix.Contains(ip) {
			return true, &prefix, "", nil
		}
	}
	return false, nil, "", nil
}

type OnHandle func(ctx *gin.Context, passed bool, prefix *netip.Prefix, host string, err error)

// New middleware for Gin Web Framework
func New(prefixes []netip.Prefix, hosts []string, onHandle OnHandle) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		passed, prefix, host, err := CanPass(ctx.ClientIP(), prefixes, hosts)
		if onHandle == nil {
			ctx.AbortWithStatus(http.StatusForbidden)
		} else {
			onHandle(ctx, passed, prefix, host, err)
		}
	}
}

// ReadFile reads the config file and sets up the IP filter.
// Line starts with # will be ignored
// It returns the prefixes and hosts that are allowed.
func ReadFile(filename string) (prefixes []netip.Prefix, hosts []string, err error) {
	ips, err := os.ReadFile(filename)

	if err != nil {
		return nil, nil, err
	}

	var okIpSet []string

	for _, ip := range strings.Split(string(ips), "\n") {
		ip = strings.TrimSpace(ip)
		if ip != "" && !strings.HasPrefix(ip, "#") {
			okIpSet = append(okIpSet, ip)
		}
	}

	prefixes, hosts = SetupIPFilter(okIpSet)

	return prefixes, hosts, nil
}
