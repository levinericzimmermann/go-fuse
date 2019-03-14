// Copyright 2019 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nodefs

import (
	"io"
	"os"

	"github.com/hanwen/go-fuse/fuse"
)

type dirArray struct {
	Entries []fuse.DirEntry
}

func (a *dirArray) HasNext() bool {
	return len(a.Entries) > 0
}

func (a *dirArray) Next() (fuse.DirEntry, fuse.Status) {
	e := a.Entries[0]
	a.Entries = a.Entries[1:]
	return e, fuse.OK
}

func (a *dirArray) Close() {

}

func NewLoopbackDirStream(nm string) (DirStream, fuse.Status) {
	f, err := os.Open(nm)
	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	defer f.Close()

	var entries []fuse.DirEntry
	for {
		want := 100
		infos, err := f.Readdir(want)
		for _, info := range infos {
			s := fuse.ToStatT(info)
			if s == nil {
				continue
			}

			entries = append(entries, fuse.DirEntry{
				Name: info.Name(),
				Mode: uint32(s.Mode),
				Ino:  s.Ino,
			})
		}
		if len(infos) < want || err == io.EOF {
			break
		}

		if err != nil {
			return nil, fuse.ToStatus(err)
		}
	}

	return &dirArray{entries}, fuse.OK
}