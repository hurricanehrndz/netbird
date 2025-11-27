//go:build !(linux && 386) && !windows && !darwin

package main

import (
	_ "embed"
)

//go:embed assets/other/netbird.png
var iconAbout []byte

//go:embed assets/other/netbird-disconnected.png
var iconAboutDisconnected []byte

//go:embed assets/other/netbird-systemtray-connected.png
var iconConnected []byte

//go:embed assets/other/netbird-systemtray-connected-dark.png
var iconConnectedDark []byte

//go:embed assets/other/netbird-systemtray-disconnected.png
var iconDisconnected []byte

//go:embed assets/other/netbird-systemtray-update-disconnected.png
var iconUpdateDisconnected []byte

//go:embed assets/other/netbird-systemtray-update-disconnected-dark.png
var iconUpdateDisconnectedDark []byte

//go:embed assets/other/netbird-systemtray-update-connected.png
var iconUpdateConnected []byte

//go:embed assets/other/netbird-systemtray-update-connected-dark.png
var iconUpdateConnectedDark []byte

//go:embed assets/other/netbird-systemtray-connecting.png
var iconConnecting []byte

//go:embed assets/other/netbird-systemtray-connecting-dark.png
var iconConnectingDark []byte

//go:embed assets/other/netbird-systemtray-error.png
var iconError []byte

//go:embed assets/other/netbird-systemtray-error-dark.png
var iconErrorDark []byte

//go:embed assets/other/connected.png
var iconConnectedDot []byte

//go:embed assets/other/disconnected.png
var iconDisconnectedDot []byte

//go:embed assets/other/netbird-systemtray-connecting-frame1.png
var iconConnectingFrame1 []byte

//go:embed assets/other/netbird-systemtray-connecting-dark-frame2.png
var iconConnectingDarkFrame1 []byte

//go:embed assets/other/netbird-systemtray-connecting-frame2.png
var iconConnectingFrame2 []byte

//go:embed assets/other/netbird-systemtray-connecting-dark-frame2.png
var iconConnectingDarkFrame2 []byte

//go:embed assets/other/netbird-systemtray-connecting-frame3.png
var iconConnectingFrame3 []byte

//go:embed assets/other/netbird-systemtray-connecting-dark-frame3.png
var iconConnectingDarkFrame3 []byte

//go:embed assets/other/netbird-systemtray-connecting-frame4.png
var iconConnectingFrame4 []byte

//go:embed assets/other/netbird-systemtray-connecting-dark-frame4.png
var iconConnectingDarkFrame4 []byte
