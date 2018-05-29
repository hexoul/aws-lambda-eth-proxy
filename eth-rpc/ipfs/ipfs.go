package ipfs

import (
	"io/ioutil"

	"github.com/ipfs/go-ipfs-api"
)

func New(url string) *shell.Shell {
	s := shell.NewShell(url)
	return s
}

func Cat(s *shell.Shell, path string) (ret string) {
	rc, err := s.Cat(path)
	if err != nil {
		return
	}
	if b, err := ioutil.ReadAll(rc); err != nil {
		ret = string(b)
		rc.Close()
	}
	return
}
