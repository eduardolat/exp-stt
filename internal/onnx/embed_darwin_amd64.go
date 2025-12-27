//go:build darwin && amd64

package onnx

import _ "embed"

var (
	//go:embed onnxruntime-osx-x86_64-1.23.2.tgz
	CompressedLib []byte

	isZip           = false
	isTgz           = true
	runtimeVersion  = "1.23.2"
	runtimePlatform = "osx-amd64"
	sharedLibName   = "libonnxruntime.1.23.2.dylib"
)
