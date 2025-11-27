package main

import (
	_ "embed"
)

//go:embed assets/win/netbird.ico
var iconAbout []byte

//go:embed assets/win/netbird-disconnected.ico
var iconAboutDisconnected []byte

//go:embed assets/win/netbird-systemtray-connected.ico
var iconConnected []byte

//go:embed assets/win/netbird-systemtray-connected-dark.ico
var iconConnectedDark []byte

//go:embed assets/win/netbird-systemtray-disconnected.ico
var iconDisconnected []byte

//go:embed assets/win/netbird-systemtray-update-disconnected.ico
var iconUpdateDisconnected []byte

//go:embed assets/win/netbird-systemtray-update-disconnected-dark.ico
var iconUpdateDisconnectedDark []byte

//go:embed assets/win/netbird-systemtray-update-connected.ico
var iconUpdateConnected []byte

//go:embed assets/win/netbird-systemtray-update-connected-dark.ico
var iconUpdateConnectedDark []byte

//go:embed assets/win/netbird-systemtray-connecting.ico
var iconConnecting []byte

//go:embed assets/win/netbird-systemtray-connecting-dark.ico
var iconConnectingDark []byte

//go:embed assets/win/netbird-systemtray-error.ico
var iconError []byte

//go:embed assets/win/netbird-systemtray-error-dark.ico
var iconErrorDark []byte

//go:embed assets/win/netbird-systemtray-connecting-frame1.ico
var iconConnectingFrame1 []byte

//go:embed assets/win/netbird-systemtray-connecting-dark-frame1.ico
var iconConnectingDarkFrame1 []byte

//go:embed assets/win/netbird-systemtray-connecting-frame2.ico
var iconConnectingFrame2 []byte

//go:embed assets/win/netbird-systemtray-connecting-dark-frame2.ico
var iconConnectingDarkFrame2 []byte

//go:embed assets/win/netbird-systemtray-connecting-frame3.ico
var iconConnectingFrame3 []byte

//go:embed assets/win/netbird-systemtray-connecting-dark-frame3.ico
var iconConnectingDarkFrame3 []byte

//go:embed assets/win/netbird-systemtray-connecting-frame4.ico
var iconConnectingFrame4 []byte

//go:embed assets/win/netbird-systemtray-connecting-dark-frame4.ico
var iconConnectingDarkFrame4 []byte
