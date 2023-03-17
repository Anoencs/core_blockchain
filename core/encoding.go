package core

import (
	"crypto/elliptic"
	"encoding/gob"
	"io"
)

func init() {
	gob.Register(elliptic.P256())
}

type Encoder[T any] interface {
	Encode(T) error
}

type Decoder[T any] interface {
	Decode(T) error
}

type GobTxEncoder struct {
	w io.Writer
}

func NewGobTxEncoder(w io.Writer) *GobTxEncoder {
	return &GobTxEncoder{
		w: w,
	}
}

func (e *GobTxEncoder) Encode(tx *Transaction) error {
	return gob.NewEncoder(e.w).Encode(tx)

}

type GobTxDecoder struct {
	r io.Reader
}

func NewGobTxDecoder(r io.Reader) *GobTxDecoder {
	return &GobTxDecoder{
		r: r,
	}
}

func (d *GobTxDecoder) Decode(tx *Transaction) error {
	return gob.NewDecoder(d.r).Decode(tx)
}

type GobBlockEncoder struct {
	w io.Writer
}

func NewBlockEncoder(w io.Writer) *GobBlockEncoder {
	return &GobBlockEncoder{
		w: w,
	}
}

func (e *GobBlockEncoder) Encode(b *Block) error {
	return gob.NewEncoder(e.w).Encode(b)
}

type GobBlockDecoder struct {
	r io.Reader
}

func NewBlockDecoder(r io.Reader) *GobBlockDecoder {
	return &GobBlockDecoder{
		r: r,
	}
}

func (d *GobBlockDecoder) Decode(b *Block) error {
	return gob.NewDecoder(d.r).Decode(b)
}
