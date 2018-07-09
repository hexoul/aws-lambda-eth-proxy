package ipfs

import (
	"crypto/md5"
	"fmt"
	"testing"

	"github.com/cheekybits/is"
)

const (
	exampleHash = "Qmd286K6pohQcTKYqnS1YhWrCiS4gz7Xi34sdwMe9USZ7u"
)

func TestPinByCluster(t *testing.T) {
	s := GetInstance()
	err := s.PinByCluster(exampleHash)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCat(t *testing.T) {
	is := is.New(t)
	s := GetInstance()

	ret := s.Cat(exampleHash)
	is.NotNil(ret)

	md5 := md5.New()
	md5.Write([]byte(ret))
	is.Equal(fmt.Sprintf("%x", md5.Sum(nil)), "b84d6366deec053ff3fa77df01a54464")
}

func TestAdd(t *testing.T) {
	is := is.New(t)
	s := GetInstance()

	mhash, err := s.Add("Hello IPFS Shell tests")
	is.Nil(err)
	is.Equal(mhash, "QmUfZ9rAdhV5ioBzXKdUTh2ZNsz9bzbkaLVyQ8uc8pj21F")
}

func TestAddnCat(t *testing.T) {
	is := is.New(t)
	s := GetInstance()

	testMsg := "TestTestTest"
	mhash, err := s.Add(testMsg)
	is.Nil(err)
	val := s.Cat(mhash)
	is.NotNil(val)
	t.Logf("Hash: %s, val: %s\n", mhash, val)
	is.Equal(testMsg, val)
}
