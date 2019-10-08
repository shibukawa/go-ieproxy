// +build darwin

package ieproxy

import (
	"reflect"
	"testing"
)

var proxyDisabled = `
<dictionary> {
  HTTPEnable : 0
  HTTPSEnable : 0
  HTTPSUser : macuser
  HTTPUser : macuser
}
`

var proxyEnabled = `
<dictionary> {
  HTTPEnable : 1
  HTTPPort : 8080
  HTTPProxy : proxy.example.com
  HTTPSEnable : 1
  HTTPSPort : 8080
  HTTPSProxy : proxy.example.com
  HTTPSUser : macuser
  HTTPUser : macuser
}
`

var proxyWithNoProxy = `
<dictionary> {
  ExceptionsList : <array> {
    0 : test1.example.com
    1 : test2.example.com
  }
  HTTPEnable : 1
  HTTPPort : 8080
  HTTPProxy : proxy.example.com
  HTTPSEnable : 1
  HTTPSPort : 8080
  HTTPSProxy : proxy.example.com
  HTTPSUser : macuser
  HTTPUser : macuser
}
`

func TestParseScutil(t *testing.T) {
	testcases := []struct {
		name string
		in   string
		out  ProxyConf
	}{
		{
			name: "proxy disabled",
			in:   proxyDisabled,
			out: ProxyConf{
				Static: StaticProxyConf{
					Protocols: map[string]string{},
				},
			},
		},
		{
			name: "proxy enabled",
			in:   proxyEnabled,
			out: ProxyConf{
				Static: StaticProxyConf{
					Active: true,
					Protocols: map[string]string{
						"http":  "proxy.example.com:8080",
						"https": "proxy.example.com:8080",
					},
					NoProxy: "",
				},
			},
		},
		{
			name: "proxy with noProxy",
			in:   proxyWithNoProxy,
			out: ProxyConf{
				Static: StaticProxyConf{
					Active: true,
					Protocols: map[string]string{
						"http":  "proxy.example.com:8080",
						"https": "proxy.example.com:8080",
					},
					NoProxy: "test1.example.com,test2.example.com",
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			parsed := parseScutil([]byte(tc.in))
			if !reflect.DeepEqual(tc.out, parsed) {
				t.Errorf("invalid: %v != %v", tc.out, parsed)
			}
		})
	}
}
