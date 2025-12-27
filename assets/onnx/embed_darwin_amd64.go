//go:build darwin && amd64

package onnx

import _ "embed"

var (
	//go:embed onnxruntime-osx-x86_64-1.23.2.tgz
	CompressedLib []byte

	IsZip = false
	IsTgz = true
)
