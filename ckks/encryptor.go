package ckks

import (
	"github.com/tuneinsight/lattigo/v3/rlwe"
)

// Encryptor an encryption interface for the CKKS scheme.
type Encryptor interface {
	Encrypt(plaintext *Plaintext, ciphertext *Ciphertext)
	EncryptNew(plaintext *Plaintext) *Ciphertext
	EncryptZero(ciphertext *Ciphertext)
	EncryptZeroNew(level int, scale float64) *Ciphertext
	ShallowCopy() Encryptor
	WithKey(key interface{}) Encryptor
}

type encryptor struct {
	rlwe.Encryptor
	params Parameters
}

// NewEncryptor instantiates a new Encryptor for the CKKS scheme. The key argument can
// be *rlwe.PublicKey, *rlwe.SecretKey or nil.
func NewEncryptor(params Parameters, key interface{}) Encryptor {
	return &encryptor{rlwe.NewEncryptor(params.Parameters, key), params}
}

// Encrypt encrypts the input plaintext and write the result on ciphertext.
// The level of the output ciphertext is min(plaintext.Level(), ciphertext.Level()).
func (enc *encryptor) Encrypt(plaintext *Plaintext, ciphertext *Ciphertext) {
	enc.Encryptor.Encrypt(plaintext.Plaintext, ciphertext.Ciphertext)
	ciphertext.Scale = plaintext.Scale
}

// EncryptNew encrypts the input plaintext returns the result as a newly allocated ciphertext.
// The level of the output ciphertext is min(plaintext.Level(), ciphertext.Level()).
func (enc *encryptor) EncryptNew(plaintext *Plaintext) (ciphertext *Ciphertext) {
	ciphertext = NewCiphertext(enc.params, 1, plaintext.Level(), plaintext.Scale)
	enc.Encryptor.Encrypt(plaintext.Plaintext, ciphertext.Ciphertext)
	return
}

// EncryptZero generates an encryption of zero at the level and scale of ct, and writes the result on ctOut.
// Note that the Scale field of an encryption of zero can be changed arbitrarily, without requiring a Rescale.
func (enc *encryptor) EncryptZero(ciphertext *Ciphertext) {
	enc.Encryptor.EncryptZero(ciphertext.Ciphertext)
}

// EncryptZero generates an encryption of zero at the given level and scale and returns the
// result as a newly allocated ciphertext.
// Note that the Scale field of an encryption of zero can be changed arbitrarily, without requiring a Rescale.
func (enc *encryptor) EncryptZeroNew(level int, scale float64) *Ciphertext {
	ct := NewCiphertext(enc.params, 1, level, scale)
	enc.Encryptor.EncryptZero(ct.Ciphertext)
	return ct
}

// ShallowCopy creates a shallow copy of this encryptor in which all the read-only data-structures are
// shared with the receiver and the temporary buffers are reallocated. The receiver and the returned
// Encryptors can be used concurrently.
func (enc *encryptor) ShallowCopy() Encryptor {
	return &encryptor{enc.Encryptor.ShallowCopy(), enc.params}
}

// WithKey creates a shallow copy of this encryptor with a new key in which all the read-only data-structures are
// shared with the receiver and the temporary buffers are reallocated. The receiver and the returned
// Encryptors can be used concurrently.
// Key can be *rlwe.PublicKey or *rlwe.SecretKey.
func (enc *encryptor) WithKey(key interface{}) Encryptor {
	return &encryptor{enc.Encryptor.WithKey(key), enc.params}
}
