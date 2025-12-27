//go:build windows && amd64

package onnx

import _ "embed"

var (
	//go:embed onnxruntime-win-x64-1.23.2.zip
	CompressedLib []byte

	isZip           = true
	isTgz           = false
	runtimeVersion  = "1.23.2"
	runtimePlatform = "windows-amd64"
	sharedLibName   = "onnxruntime.dll"
)
