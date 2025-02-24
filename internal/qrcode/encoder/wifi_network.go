package encoder

import (
	"fmt"
	"strings"

	"github.com/sudosz/qmars/internal/qrcode"
)

// WIFI:T:<TYPE[WPA/WEP/nopass]>;S:<SSID>;P:<PASSWORD>;<HIDDEN_STATUS>;
// e.g. WIFI:T:WPA;S:TP-Link;P:12345;H:true;
const wifiFormat = "WIFI:T:%s;S:%s;P:%s;%s;"

type wifiNetworkContent struct {
	ssid         string
	password     string
	securityType qrcode.WiFiSecurityType
	hiddenStatus bool
}

func WiFiNetworkContent(ssid, password string, securityType qrcode.WiFiSecurityType, hiddenStatus ...bool) Content {
	return newWiFiNetworkContent(ssid, password, securityType, hiddenStatus...)
}

func WiFiNetworkNoPasswordContent(ssid string, hiddenStatus ...bool) Content {
	return newWiFiNetworkContent(ssid, "", qrcode.NoPassword, hiddenStatus...)
}

func newWiFiNetworkContent(ssid, password string, securityType qrcode.WiFiSecurityType, hiddenStatus ...bool) Content {
	h := false
	if len(hiddenStatus) > 0 {
		h = hiddenStatus[0]
	}
	return &wifiNetworkContent{
		ssid:         ssid,
		password:     password,
		securityType: securityType,
		hiddenStatus: h,
	}
}

func (c *wifiNetworkContent) SetSSID(ssid string) Content {
	c.ssid = ssid
	return c
}

func (c *wifiNetworkContent) SetPassword(password string) Content {
	c.password = password
	return c
}

func (c *wifiNetworkContent) SetNoPassword() Content {
	c.password = ""
	c.securityType = qrcode.NoPassword
	return c
}

func (c *wifiNetworkContent) SetSecurityType(securityType qrcode.WiFiSecurityType) Content {
	c.securityType = securityType
	return c
}

func (c *wifiNetworkContent) SetHidden(hiddenStatus bool) Content {
	c.hiddenStatus = hiddenStatus
	return c
}

func escapePassword(p string) string {
	return strings.ReplaceAll(p, ";", "\\;")
}

func (c *wifiNetworkContent) Get() string {
	hidden := ""
	if c.hiddenStatus {
		hidden = "H:true"
	}

	return fmt.Sprintf(wifiFormat, string(c.securityType), c.ssid, escapePassword(c.password), hidden)
}
