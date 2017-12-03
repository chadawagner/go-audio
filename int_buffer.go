package audio

import "math"

// IntBuffer is an audio buffer with its PCM data formatted as int.
type IntBuffer struct {
	// Format is the representation of the underlying data format
	Format *Format
	// Data is the buffer PCM data as ints
	Data []int
	// SourceBitDepth helps us know if the source was encoded on
	// 1 (int8), 2 (int16), 3(int24), 4(int32), 8(int64) bytes.
	SourceBitDepth int
}

// PCMFormat returns the buffer format information.
func (buf *IntBuffer) PCMFormat() *Format { return buf.Format }

// AsFloatBuffer returns a copy of this buffer but with data converted to floats.
func (buf *IntBuffer) AsFloatBuffer() *FloatBuffer {
	newB := &FloatBuffer{}
	newB.Data = make([]float64, len(buf.Data))
	for i := 0; i < len(buf.Data); i++ {
		newB.Data[i] = float64(buf.Data[i])
	}
	newB.Format = &Format{
		NumChannels: buf.Format.NumChannels,
		SampleRate:  buf.Format.SampleRate,
	}
	return newB
}

// AsFloat32Buffer returns a copy of this buffer but with data converted to float 32.
func (buf *IntBuffer) AsFloat32Buffer() *Float32Buffer {
	newB := &Float32Buffer{}
	newB.Data = make([]float32, len(buf.Data))
	max := 0
	bitDepth := buf.SourceBitDepth
	// try to guess the bit depths without knowing the source
	if bitDepth == 0 {
		for _, s := range buf.Data {
			if s > max {
				max = s
			}
		}
		bitDepth = 8
		if max > 127 {
			bitDepth = 16
		}
		// greater than int16, expecting int24
		if max > 32767 {
			bitDepth = 24
		}
		// int 32
		if max > 8388607 {
			bitDepth = 32
		}
		// int 64
		if max > 4294967295 {
			bitDepth = 64
		}
	}
	factor := math.Pow(2, 8*float64(bitDepth/8)-1)
	for i := 0; i < len(buf.Data); i++ {
		newB.Data[i] = float32(float64(int64(buf.Data[i])) / factor)
	}
	newB.Format = &Format{
		NumChannels: buf.Format.NumChannels,
		SampleRate:  buf.Format.SampleRate,
	}
	return newB
}

// AsIntBuffer implements the Buffer interface and returns itself.
func (buf *IntBuffer) AsIntBuffer() *IntBuffer { return buf }

// NumFrames returns the number of frames contained in the buffer.
func (buf *IntBuffer) NumFrames() int {
	if buf == nil || buf.Format == nil {
		return 0
	}
	numChannels := buf.Format.NumChannels
	if numChannels == 0 {
		numChannels = 1
	}

	return len(buf.Data) / numChannels
}

// Clone creates a clean clone that can be modified without
// changing the source buffer.
func (buf *IntBuffer) Clone() Buffer {
	if buf == nil {
		return nil
	}
	newB := &IntBuffer{}
	newB.Data = make([]int, len(buf.Data))
	copy(newB.Data, buf.Data)
	newB.Format = &Format{
		NumChannels: buf.Format.NumChannels,
		SampleRate:  buf.Format.SampleRate,
	}
	return newB
}
