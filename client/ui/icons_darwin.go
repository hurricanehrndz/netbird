package main

import (
	_ "embed"
)

//go:embed assets/mac/netbird.png
var iconAbout []byte

//go:embed assets/mac/netbird-disconnected.png
var iconAboutDisconnected []byte

//go:embed assets/mac/netbird-systemtray-connected-macos.svg
var iconConnected []byte
var iconConnectedDark  = iconConnected

//go:embed assets/mac/netbird-systemtray-disconnected-macos.svg
var iconDisconnected []byte

//go:embed assets/mac/netbird-systemtray-update-disconnected-macos.svg
var iconUpdateDisconnected []byte
var iconUpdateDisconnectedDark = iconUpdateDisconnected

//go:embed assets/mac/netbird-systemtray-update-connected-macos.svg
var iconUpdateConnected []byte
var iconUpdateConnectedDark = iconUpdateConnected

//go:embed assets/mac/netbird-systemtray-connecting-macos.svg
var iconConnecting []byte
var iconConnectingDark  = iconConnecting

//go:embed assets/mac/netbird-systemtray-error-macos.svg
var iconError []byte
var iconErrorDark  = iconError

//go:embed assets/mac/connected.svg
var iconConnectedDot []byte

//go:embed assets/mac/disconnected.svg
var iconDisconnectedDot []byte

//go:embed assets/mac/netbird-systemtray-connecting-frame1-macos.svg
var iconConnectingFrame1 []byte
var iconConnectingDarkFrame1  = iconConnectingFrame1

//go:embed assets/mac/netbird-systemtray-connecting-frame2-macos.svg
var iconConnectingFrame2 []byte
var iconConnectingDarkFrame2  = iconConnectingFrame2

//go:embed assets/mac/netbird-systemtray-connecting-frame3-macos.svg
var iconConnectingFrame3 []byte
var iconConnectingDarkFrame3  = iconConnectingFrame3

//go:embed assets/mac/netbird-systemtray-connecting-frame4-macos.svg
var iconConnectingFrame4 []byte
var iconConnectingDarkFrame4  = iconConnectingFrame4
