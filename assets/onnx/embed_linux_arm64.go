//go:build linux && arm64

package onnx

import _ "embed"

var (
	//go:embed onnxruntime-linux-aarch64-1.23.2.tgz
	CompressedLib []byte

	IsZip = false
	IsTgz = true
)
