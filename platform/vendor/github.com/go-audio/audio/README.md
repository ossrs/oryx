# audio

[![GoDoc](http://godoc.org/github.com/go-audio/audio?status.svg)](http://godoc.org/github.com/go-audio/audio)

`audio` is a generic Go package designed to define a common interface to analyze
and/or process audio data.

At the heart of the package is the `Buffer` interface and its implementations:

* `FloatBuffer`
* `Float32Buffer`
* `IntBuffer`

Decoders, encoders, processors, analyzers and transformers can be written to
accept or return these types and share a common interface.

The idea is that audio libraries can define this interface or its
implementations as input and return an `audio.Buffer` interface allowing all
audio libraries to be chainable.

## Performance

The buffer implementations are designed so a buffer can be reused and mutated
avoiding allocation penalties.

It is recommended to avoid using `Float32Buffer` unless performance is critical.
The major drawback of using float32s is that the Go stdlib was designed to work
with float64 and therefore the access to standard packages is limited.

## Usage

Examples of how to use this interface is available under the
[go-audio](https://github.com/go-audio) organization.
