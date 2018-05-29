package ipfs

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"testing"

	"github.com/cheekybits/is"
)

const (
	shellUrl    = "localhost:5001"
	exampleHash = "Qmd286K6pohQcTKYqnS1YhWrCiS4gz7Xi34sdwMe9USZ7u"
)

func TestCat(t *testing.T) {
	is := is.New(t)
	s := New(shellUrl)

	rc, err := s.Cat(exampleHash)
	is.Nil(err)

	md5 := md5.New()
	_, err = io.Copy(md5, rc)
	is.Nil(err)
	is.Equal(fmt.Sprintf("%x", md5.Sum(nil)), "b84d6366deec053ff3fa77df01a54464")
}

func TestAdd(t *testing.T) {
	is := is.New(t)
	s := New(shellUrl)

	mhash, err := s.Add(bytes.NewBufferString("Hello IPFS Shell tests"))
	is.Nil(err)
	is.Equal(mhash, "QmUfZ9rAdhV5ioBzXKdUTh2ZNsz9bzbkaLVyQ8uc8pj21F")
}
