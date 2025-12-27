//go:build darwin && arm64

package onnx

import _ "embed"

var (
	//go:embed onnxruntime-osx-arm64-1.23.2.tgz
	CompressedLib []byte

	isZip           = false
	isTgz           = true
	runtimeVersion  = "1.23.2"
	runtimePlatform = "osx-arm64"
	sharedLibName   = "libonnxruntime.1.23.2.dylib"
)
