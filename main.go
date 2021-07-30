package main

import (
	"archive/tar"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

var (
	r    = flag.String("ref", "gcr.io/imjasonh/data", "ref to push to")
	oci  = flag.Bool("oci", false, "if true, push as OCI image")
	size = flag.Int("size", 10, "size of bytes to send")
)

func main() {
	flag.Parse()

	ref, err := name.ParseReference(*r)
	if err != nil {
		log.Fatal(err)
	}

	var img v1.Image = empty.Image
	if *oci {
		img = mutate.MediaType(img, types.OCIManifestSchema1)
	}
	mt := types.DockerLayer
	if *oci {
		mt = types.OCILayer
	}

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	if err := tw.WriteHeader(&tar.Header{
		Name: "hello.txt",
		Mode: 0600,
		Size: int64(*size),
	}); err != nil {
		log.Fatal(err)
	}
	if _, err := tw.Write([]byte(strings.Repeat(".", *size))); err != nil {
		log.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		log.Fatal(err)
	}

	img, err = mutate.AppendLayers(img, &staticLayer{
		b:  buf.Bytes(),
		mt: mt,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := remote.Write(ref, img, remote.WithAuthFromKeychain(authn.DefaultKeychain)); err != nil {
		log.Fatal(err)
	}

	d, err := img.Digest()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ref.Context().Digest(d.String()))
}

// copied from cosign:

type staticLayer struct {
	b  []byte
	mt types.MediaType
}

func (l *staticLayer) Digest() (v1.Hash, error) {
	h, _, err := v1.SHA256(bytes.NewReader(l.b))
	return h, err
}

// DiffID returns the Hash of the uncompressed layer.
func (l *staticLayer) DiffID() (v1.Hash, error) {
	h, _, err := v1.SHA256(bytes.NewReader(l.b))
	return h, err
}

// Compressed returns an io.ReadCloser for the compressed layer contents.
func (l *staticLayer) Compressed() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewReader(l.b)), nil
}

// Uncompressed returns an io.ReadCloser for the uncompressed layer contents.
func (l *staticLayer) Uncompressed() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewReader(l.b)), nil
}

// Size returns the compressed size of the Layer.
func (l *staticLayer) Size() (int64, error) {
	return int64(len(l.b)), nil
}

// MediaType returns the media type of the Layer.
func (l *staticLayer) MediaType() (types.MediaType, error) {
	return l.mt, nil
}

// Descriptor returns a Descriptor that includes the embedded Data.
func (l *staticLayer) Descriptor() (*v1.Descriptor, error) {
	digest, err := l.Digest()
	if err != nil {
		return nil, err
	}
	desc := &v1.Descriptor{
		MediaType: l.mt,
		Data:      l.b,
		Size:      int64(len(l.b)),
		Digest:    digest,
	}

	return desc, nil
}
