package qrcode

import (
	"image/color"
)

type ErrorCorrectionLevel uint8

const (
	// Level L: 7% error recovery.
	ErrorCorrectionLevelLow ErrorCorrectionLevel = iota
	// Level M: 15% error recovery. Good default choice.
	ErrorCorrectionLevelMedium
	// Level Q: 25% error recovery.
	ErrorCorrectionLevelQuartile
	// Level H: 30% error recovery.
	ErrorCorrectionLevelHigh

	DefaultRecoveryLevel = ErrorCorrectionLevelHigh
)

type Version uint8

const (
	DefaultVersion Version = iota
	V1
	V2
	V3
	V4
	V5
	V6
	V7
	V8
	V9
	V10
	V11
	V12
	V13
	V14
	V15
	V16
	V17
	V18
	V19
	V20
	V21
	V22
	V23
	V24
	V25
	V26
	V27
	V28
	V29
	V30
	V31
	V32
	V33
	V34
	V35
	V36
	V37
	V38
	V39
	V40
)

type WiFiSecurityType string

const (
	WiFiSecurityTypeWPA  WiFiSecurityType = "WPA"
	WiFiSecurityTypeWPA2 WiFiSecurityType = "WPA2"

	WiFiSecurityTypeWEP        WiFiSecurityType = "WEP"
	WiFiSecurityTypeNoPassword WiFiSecurityType = "nopass"
)

const (
	DefaultWidth  = 14
	DefaultHeight = 14
)

var (
	DefaultForeground color.Color = color.Black
	DefaultBackground color.Color = color.White
)
