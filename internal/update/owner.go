package update

import (
	"os"
	"syscall"
)

// sameOwner donne à path le propriétaire de ref. Le conteneur de mise à jour
// tourne en root, le serveur en utilisateur 10001 : une base restaurée ou un
// état écrit par root seraient illisibles pour le serveur sans cela.
func sameOwner(path, ref string) error {
	info, err := os.Stat(ref)
	if err != nil {
		return err
	}
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return nil
	}
	if err := os.Lchown(path, int(st.Uid), int(st.Gid)); err != nil && !os.IsPermission(err) {
		return err
	}
	return nil
}
