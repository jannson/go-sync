package indexbuilder

import (
	"bytes"
	"github.com/Redundancy/go-sync/chunks"
	"github.com/Redundancy/go-sync/filechecksum"
	"github.com/Redundancy/go-sync/index"
	"io"
)

// Generates an index from a reader
// This is mostly a utility function to avoid being overly verbose in tests that need
// an index to work, but don't want to construct one by hand in order to avoid the dependencies
// obviously this means that those tests are likely to fail if there are issues with any of the other
// modules, which is not ideal.
// TODO: move to util?
func BuildChecksumIndex(check *filechecksum.FileChecksumGenerator, r io.Reader) (
	fcheck []byte,
	i *index.ChecksumIndex,
	err error,
) {
	b := bytes.NewBuffer(nil)
	fcheck, err = check.GenerateChecksums(r, b)

	if err != nil {
		return
	}

	weakSize := check.WeakRollingHash.Size()
	strongSize := check.GetStrongHash().Size()
	readChunks, err := chunks.LoadChecksumsFromReader(b, weakSize, strongSize)

	if err != nil {
		return
	}

	i = index.MakeChecksumIndex(readChunks)

	return
}

func BuildIndexFromString(generator *filechecksum.FileChecksumGenerator, reference string) (
	fileCheckSum []byte,
	referenceIndex *index.ChecksumIndex,
	err error,
) {
	return BuildChecksumIndex(generator, bytes.NewBufferString(reference))
}
