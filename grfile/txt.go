package grfile

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ReadFile(filename string) (contents []byte, err error) {

	fd, err := os.Open(filename)

	if err != nil {

		err = fmt.Errorf("LoadConfig: Error: Counld not open %q for reading: %s\n ", filename, err)

		return
	}
	defer fd.Close()

	contents, err = ioutil.ReadAll(fd)
	if err != nil {
		err = fmt.Errorf("LoadConfig: Error: Could not open %q: %s \n", filename, err)

		return
	}

	return

}
