package encoder

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"

	gqrcode "github.com/skip2/go-qrcode"
	"github.com/sudosz/qmars/internal/qrcode"
)

var (
	testString           = "Hello, world!"
	testBytes            = []byte(testString)
	testWiFiSecurityType = qrcode.WPA
	testWiFiSSID         = "TP-Link"
	testWiFiPassword     = "12345"
	testHiddenStatus     = true

	globalLevel = gqrcode.Highest
)

func encode2QR(data any) qrcode.QRCode {
	var s string

	switch v := data.(type) {
	case string:
		s = v
	case []byte:
		s = base64.StdEncoding.EncodeToString(v)
	}

	enc, _ := gqrcode.New(s, globalLevel)
	return enc
}

func testWiFiNetwork() string {
	hidden := ""
	if testHiddenStatus {
		hidden = "H:true"
	}
	return fmt.Sprintf(wifiFormat, string(testWiFiSecurityType), testWiFiSSID, testWiFiPassword, hidden)
}

func testWiFiNetworkNoPassword() string {
	hidden := ""
	if testHiddenStatus {
		hidden = "H:true"
	}
	return fmt.Sprintf(wifiFormat, string(qrcode.NoPassword), testWiFiSSID, "", hidden)
}

func TestQREncode(t *testing.T) {
	tests := []struct {
		name     string
		content  Content
		expected qrcode.QRCode
	}{
		{
			name:     "StringContent",
			content:  StringContent(testString),
			expected: encode2QR(testString),
		},
		{
			name:     "BytesContent",
			content:  BytesContent(testBytes),
			expected: encode2QR(testBytes),
		},
		{
			name: "WiFiNetworkContent",
			content: WiFiNetworkContent(
				testWiFiSSID, testWiFiPassword, testWiFiSecurityType, testHiddenStatus,
			),
			expected: encode2QR(testWiFiNetwork()),
		},
		{
			name: "WiFiNetworkNoPasswordContent",
			content: WiFiNetworkNoPasswordContent(
				testWiFiSSID, testHiddenStatus,
			),
			expected: encode2QR(testWiFiNetworkNoPassword()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := NewQRCodeBuilder().SetContent(tt.content).Build()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(out.Bitmap(), tt.expected.Bitmap()) {
				t.Fatalf("expected:\n	%#v\n	got:\n	%#v\n", tt.expected, out)
			}
		})
	}
}
