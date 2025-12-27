//go:build linux && arm64

package onnx

import _ "embed"

var (
	//go:embed onnxruntime-linux-aarch64-1.23.2.tgz
	CompressedLib []byte

	isZip           = false
	isTgz           = true
	runtimeVersion  = "1.23.2"
	runtimePlatform = "linux-arm64"
	sharedLibName   = "libonnxruntime.so.1.23.2"
)
