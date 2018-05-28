package ipfs

import (
	"fmt"

	"github.com/ipfs/go-ipfs-api"
)

func New() *shell.Shell {
	s := shell.NewLocalShell()
	return s
}
