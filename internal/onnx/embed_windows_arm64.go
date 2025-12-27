//go:build windows && arm64

package onnx

import _ "embed"

var (
	//go:embed onnxruntime-win-arm64-1.23.2.zip
	CompressedLib []byte

	isZip           = true
	isTgz           = false
	runtimeVersion  = "1.23.2"
	runtimePlatform = "windows-arm64"
	sharedLibName   = "onnxruntime.dll"
)
