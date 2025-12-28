package sounds

import _ "embed"

// InOutSound represents a pair of input and output sounds. Great for start/stop actions.
type InOutSound struct {
	ID     string
	Input  []byte
	Output []byte
}

// InOutSounds is a collection of input/output sound pairs. Great for start/stop actions.
var InOutSounds = []InOutSound{
	{ID: "1", Input: in1, Output: out1},
	{ID: "2", Input: in2, Output: out2},
	{ID: "3", Input: in3, Output: out3},
	{ID: "4", Input: in4, Output: out4},
	{ID: "5", Input: in5, Output: out5},
	{ID: "6", Input: in6, Output: out6},
	{ID: "7", Input: in7, Output: out7},
	{ID: "8", Input: in8, Output: out8},
	{ID: "9", Input: in9, Output: out9},
}

// SuccessSound represents a sound played on successful operations.
type SuccessSound struct {
	ID    string
	Sound []byte
}

// SuccessSounds is a collection of sounds played on successful operations.
var SuccessSounds = []SuccessSound{
	{ID: "1", Sound: success1},
	{ID: "2", Sound: success2},
	{ID: "3", Sound: success3},
	{ID: "4", Sound: success4},
}

// ErrorSound represents a sound played on error events.
type ErrorSound struct {
	ID    string
	Sound []byte
}

// ErrorSounds is a collection of sounds played on error events.
var ErrorSounds = []ErrorSound{
	{ID: "1", Sound: error1},
	{ID: "2", Sound: error2},
	{ID: "3", Sound: error3},
	{ID: "4", Sound: error4},
	{ID: "5", Sound: error5},
	{ID: "6", Sound: error6},
	{ID: "7", Sound: error7},
	{ID: "8", Sound: error8},
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

	//go:embed success_1.wav
	success1 []byte
	//go:embed success_2.wav
	success2 []byte
	//go:embed success_3.wav
	success3 []byte
	//go:embed success_4.wav
	success4 []byte

	//go:embed error_1.wav
	error1 []byte
	//go:embed error_2.wav
	error2 []byte
	//go:embed error_3.wav
	error3 []byte
	//go:embed error_4.wav
	error4 []byte
	//go:embed error_5.wav
	error5 []byte
	//go:embed error_6.wav
	error6 []byte
	//go:embed error_7.wav
	error7 []byte
	//go:embed error_8.wav
	error8 []byte
)
