/*
Package audio defines a common
interface to analyze and/or process audio data.

At the heart of the package is the Buffer interface and its implementations:
FloatBuffer and IntBuffer.
Decoders, encoders, processors, analyzers and transformers can be written to
accept or return these types and share a common interface.
*/
package audio
