package sign

import (
	"io/ioutil"
)

const outputDirectory = "/var/lib/coredns"

// write writes out the zone file to a temporary file which can then be moved into place.
func write() error {
	tmpFile, err := ioutil.TempFile(outputDirectory, "signed-")
	if err != nil {
		return err
	}
	defer tmpFile.Close()
	// write it out

	return nil
}
