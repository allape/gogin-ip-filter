package ipfilter

import (
	"net/netip"
	"testing"
)

var TestSet = map[string]bool{
	"192.168.1.1":     true,
	"::1":             true,
	"localhost":       true,
	"172.123.122.1":   false,
	"10.1.10.123":     false,
	"0000::1":         true,
	"1234:1234::1111": false,
	"abcd::ef12":      true,
	"abcd:1234::ef12": false,
	"7890:1234::ef12": true,
}

func TestCanPass(t *testing.T) {
	var okIPSet = []string{
		"192.168.0.0/16",
		"::1",
		"::1/128",
		"localhost",
		"127.0.0.1",
		"abcd::ef12/128",
		"7890::0/16",
	}
	var hosts []string
	var prefixes []netip.Prefix
	prefixes, hosts = SetupIPFilter(okIPSet)

	for addr, expected := range TestSet {
		passed, _, _, err := CanPass(addr, prefixes, hosts)
		if err != nil {
			t.Error(err)
		}
		if passed != expected {
			t.Errorf("expected %v, got %v for %s", expected, passed, addr)
		}

	}
}

//func main() {
//	reader := bufio.NewReader(os.Stdin)
//	for {
//		fmt.Printf("Please input an IP: ")
//		text, err := reader.ReadString('\n')
//		if err != nil {
//			panic(err)
//		}
//		passed, prefix, host, err := CanPass(text, Prefixes, Hosts)
//		if err != nil {
//			fmt.Println(text, "is not a valid IP")
//			continue
//		}
//		text = strings.TrimSpace(text)
//		if passed {
//			if host != "" {
//				fmt.Println("found exact match for", text)
//				continue
//			} else if prefix != nil {
//				fmt.Println("found prefix match for", text, "with prefix", prefix.String())
//				continue
//			} else {
//				fmt.Println("impossible that you can see this message")
//				continue
//			}
//		}
//		fmt.Println("no match found")
//	}
//}
