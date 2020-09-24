package main

import (
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/life4/flakehell-online/wasm/statik"

	"github.com/rakyll/statik/fs"
)

type Scripts struct {
	sfs http.FileSystem
}

func (sc *Scripts) read(fname string) string {
	file, err := sc.sfs.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func (sc *Scripts) ReadFlakeHell() string {
	return sc.read("/flakehell.py")
}

func (sc *Scripts) ReadExtract() string {
	return sc.read("/extract.py")
}

func NewScripts() Scripts {
	sfs, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	return Scripts{sfs: sfs}
}
