package encoder

import (
	"encoding/base64"
	"fmt"
	"image/color"
	"testing"

	"github.com/makiuchi-d/gozxing"
	gqrcode "github.com/makiuchi-d/gozxing/qrcode"
	"github.com/stretchr/testify/assert"
	"github.com/sudosz/qmars/internal/qrcode"
)

var (
	testString           = "Hello, world!"
	testBytes            = []byte(testString)
	testWiFiSecurityType = qrcode.WiFiSecurityTypeWPA
	testWiFiSSID         = "TP-Link"
	testWiFiPassword     = "12345"
	testHiddenStatus     = true

	globalLevel      = qrcode.ErrorCorrectionLevelLow
	globalForeground = color.RGBA{
		R: 51,
		G: 51,
		B: 51,
	}
	globalBackground = color.RGBA{
		R: 254,
		G: 254,
		B: 254,
	}
	globalMarginSize = qrcode.DefaultMarginSize
)

func encode2QR(data any) *qrcode.QRCode {
	var s string

	switch v := data.(type) {
	case string:
		s = v
	case []byte:
		s = base64.StdEncoding.EncodeToString(v)
	}

	enc := gqrcode.NewQRCodeWriter()
	hints := map[gozxing.EncodeHintType]interface{}{
		gozxing.EncodeHintType_ERROR_CORRECTION: ecLevelToGozxingECLevel(globalLevel),
	}

	b, _ := enc.Encode(s, gozxing.BarcodeFormat_QR_CODE, 0, 0, hints)
	return qrcode.NewQRCode(b, false, globalMarginSize, globalForeground, globalBackground)
}

func formatWiFiNetwork(securityType qrcode.WiFiSecurityType, ssid, password string, hidden bool) string {
	hiddenStatus := ""
	if hidden {
		hiddenStatus = "H:true"
	}
	return fmt.Sprintf(wifiFormat, string(securityType), ssid, password, hiddenStatus)
}

func TestQREncode(t *testing.T) {
	tests := []struct {
		name     string
		content  Content
		expected *qrcode.QRCode
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
			expected: encode2QR(formatWiFiNetwork(testWiFiSecurityType, testWiFiSSID, testWiFiPassword, testHiddenStatus)),
		},
		{
			name: "WiFiNetworkNoPasswordContent",
			content: WiFiNetworkNoPasswordContent(
				testWiFiSSID, testHiddenStatus,
			),
			expected: encode2QR(formatWiFiNetwork(qrcode.WiFiSecurityTypeNoPassword, testWiFiSSID, "", testHiddenStatus)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := NewQRCodeBuilder().
				SetErrorCorrectionLevel(globalLevel).
				SetForeground(globalForeground).
				SetBackground(globalBackground).
				SetMarginSize(globalMarginSize).
				SetContent(tt.content).
				Build()
			if assert.NoError(t, err) {
				assert.Equal(t, tt.expected, out)
			}
		})
	}
}
