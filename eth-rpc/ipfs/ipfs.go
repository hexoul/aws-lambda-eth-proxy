package ipfs

import (
	"bytes"
	_ "crypto/md5"
	"io/ioutil"

	"github.com/ipfs/go-ipfs-api"
)

type Ipfs struct {
	s *shell.Shell
}

func New(url string) *Ipfs {
	ns := shell.NewShell(url)
	return &Ipfs{
		s: ns,
	}
}

func (ipfs *Ipfs) Cat(path string) (ret string) {
	rc, err := ipfs.s.Cat(path)
	if err != nil {
		return
	}
	if b, err := ioutil.ReadAll(rc); err == nil {
		ret = string(b)
		rc.Close()
	}
	return
}

func (ipfs *Ipfs) Add(data string) (string, error) {
	ndata := bytes.NewBufferString(data)
	return ipfs.s.Add(ndata)
}
