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
	s := ""

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

func TestQREncodeString(t *testing.T) {
	out, err := NewQRCodeBuilder().
		SetContent(StringContent(testString)).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := encode2QR(testString)

	if !reflect.DeepEqual(out.Bitmap(), expected.Bitmap()) {
		t.Fatalf("expected:\n	%#v\n	got:\n	%#v\n", expected, out)
	}
}

func TestQREncodeBytes(t *testing.T) {
	out, err := NewQRCodeBuilder().
		SetContent(BytesContent(testBytes)).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := encode2QR(testBytes)

	if !reflect.DeepEqual(out.Bitmap(), expected.Bitmap()) {
		t.Fatalf("expected:\n	%#v\n	got:\n	%#v\n", expected, out)
	}
}

func TestQREncodeWiFiNetwork(t *testing.T) {
	out, err := NewQRCodeBuilder().
		SetContent(
			WiFiNetworkContent(testWiFiSSID, testWiFiPassword, testWiFiSecurityType, testHiddenStatus),
		).Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := encode2QR(testWiFiNetwork())

	if !reflect.DeepEqual(out.Bitmap(), expected.Bitmap()) {
		t.Fatalf("expected:\n	%#v\n	got:\n	%#v\n", expected, out)
	}
}
