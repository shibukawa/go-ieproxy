package ieproxy

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
	"sync"
)

var once sync.Once
var darwinProxyConf ProxyConf

func getConf() ProxyConf {
	once.Do(writeConf)
	return darwinProxyConf
}

func writeConf() {
	cmd := exec.Command("scutil", "--proxy")
	result, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	darwinProxyConf = parseScutil(result)
}

func parseScutil(result []byte) ProxyConf {
	osPref := make(map[string]string)
	var noProxy []string

	scanner := bufio.NewScanner(bytes.NewReader(result))
	inExceptionList := false
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 1 {
			continue
		}
		if inExceptionList {
			if fields[0] == "}" {
				inExceptionList = false
			} else if len(fields) == 3 && fields[1] == ":" {
				noProxy = append(noProxy, fields[2])
			}
		} else {
			if fields[0] == "ExceptionsList" {
				inExceptionList = true
			} else if len(fields) == 3 && fields[1] == ":" {
				if strings.HasPrefix(fields[0], "HTTP") {
					osPref[fields[0]] = fields[2]
				}
			}
		}
	}
	protocol := make(map[string]string)
	for _, scheme := range []string{"HTTP", "HTTPS"} {
		if v, ok := osPref[scheme+"Enable"]; ok && v == "1" {
			protocol[strings.ToLower(scheme)] = osPref[scheme+"Proxy"] + ":" + osPref[scheme+"Port"]
		}
	}
	return ProxyConf{
		Static: StaticProxyConf{
			Active:    len(protocol) > 0,
			Protocols: protocol,
			NoProxy:   strings.Join(noProxy, ","),
		},
	}
}

func overrideEnvWithStaticProxy(conf ProxyConf, setenv envSetter) {
	// todo
}
