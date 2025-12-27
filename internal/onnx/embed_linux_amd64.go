//go:build linux && amd64

package onnx

import _ "embed"

var (
	//go:embed onnxruntime-linux-x64-1.23.2.tgz
	CompressedLib []byte

	IsZip = false
	IsTgz = true
)
