# gogin-ip-filter

IP filter middleware for Gin Web Framework

### Usage
```go
import (
    "github.com/gin-gonic/gin"
    "github.com/allape/gogin-ip-filter"
)

func main() {
    prefixes, hosts, err := ip_filter.ReadFile("allowed-ips.txt")
	if err != nil {
        panic(err)
    }
	
    r := gin.Default()
	r.Use(ip_filter.New(prefixes, hosts))
	
	// ...
}
```
