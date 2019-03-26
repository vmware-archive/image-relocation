package image

import "github.com/opencontainers/go-digest"

// Digest provides a CAS address of an image.
type Digest struct {
	dig digest.Digest
}

func NewDigest(dig string) Digest {
	return Digest{digest.Digest(dig)}
}

var EmptyDigest Digest

func init() {
	EmptyDigest = Digest{""}
}

func (d Digest) String() string {
	return string(d.dig)
}
