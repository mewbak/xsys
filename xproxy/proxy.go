package xproxy

import (
	"github.com/getlantern/sysproxy"
)

func SystemProxyOn(proxyHostAddr string) error {
	return sysproxy.On(proxyHostAddr)
}

func SystemProxyOff(proxyHostAddr string) error {
	return sysproxy.Off(proxyHostAddr)
}
