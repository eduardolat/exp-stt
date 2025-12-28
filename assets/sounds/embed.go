package sounds

import _ "embed"

type Sound struct {
	Input  []byte
	Output []byte
}

var Sounds = []Sound{
	{Input: in1, Output: out1},
	{Input: in2, Output: out2},
	{Input: in3, Output: out3},
	{Input: in4, Output: out4},
	{Input: in5, Output: out5},
	{Input: in6, Output: out6},
	{Input: in7, Output: out7},
	{Input: in8, Output: out8},
	{Input: in9, Output: out9},
}

var (
	//go:embed in_1.wav
	in1 []byte
	//go:embed in_2.wav
	in2 []byte
	//go:embed in_3.wav
	in3 []byte
	//go:embed in_4.wav
	in4 []byte
	//go:embed in_5.wav
	in5 []byte
	//go:embed in_6.wav
	in6 []byte
	//go:embed in_7.wav
	in7 []byte
	//go:embed in_8.wav
	in8 []byte
	//go:embed in_9.wav
	in9 []byte

	//go:embed out_1.wav
	out1 []byte
	//go:embed out_2.wav
	out2 []byte
	//go:embed out_3.wav
	out3 []byte
	//go:embed out_4.wav
	out4 []byte
	//go:embed out_5.wav
	out5 []byte
	//go:embed out_6.wav
	out6 []byte
	//go:embed out_7.wav
	out7 []byte
	//go:embed out_8.wav
	out8 []byte
	//go:embed out_9.wav
	out9 []byte
)
