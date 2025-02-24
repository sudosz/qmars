package qrcode

import "image"

type QRCode interface {
	Bitmap() [][]bool
	Image(int) image.Image
}

type RecoveryLevel uint8

const (
	// Level L: 7% error recovery.
	Low RecoveryLevel = iota
	// Level M: 15% error recovery. Good default choice.
	Medium
	// Level Q: 25% error recovery.
	Quartile
	// Level H: 30% error recovery.
	High

	DefaultRecoveryLevel = High
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

type WIFISecurityType string

const (
	WPA  WIFISecurityType = "WPA"
	WPA2 WIFISecurityType = "WPA2"

	WEP        WIFISecurityType = "WEP"
	NoPassword WIFISecurityType = "nopass"
)
