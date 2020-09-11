package utils

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Params struct {
	src string
	dst string
}

type FileCopy struct {
	Dir  chan *Params
	File chan *Params
	Wg   sync.WaitGroup
}

func (f *FileCopy) Run(src string, dst string) (err error) {
	f.Dir <- &Params{
		src: src,
		dst: dst,
	}

	for {
		select {
		case dir := <-f.Dir:
			go func() {
				err := f.copyDir(dir.src, dir.dst)
				if err != nil {
					log.Println(err)
				}
			}()

		case src := <-f.File:
			go func(src *Params) {
				err := f.copyFile(src.src, src.dst)
				if err != nil {
					log.Println(err)
				}
			}(src)
		}
	}
}

// Deep copy
func (f *FileCopy) copyDir(src string, dst string) error {
	// Use trimleft ,because src may be root path `/xx`
	src = strings.TrimRight(src, string(os.PathSeparator))
	dst = strings.TrimRight(dst, string(os.PathSeparator))

	fileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		return errors.New("The src is not dir")
	}

	// Foreach all element
	srcInfo, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range srcInfo {

		newSrc := strings.Join([]string{src, file.Name()}, string(os.PathSeparator))
		newDst := strings.Join([]string{dst, file.Name()}, string(os.PathSeparator))

		if file.IsDir() {
			f.Dir <- &Params{
				src: newSrc,
				dst: newDst,
			}
		} else {
			f.File <- &Params{
				src: newSrc,
				dst: newDst,
			}
		}
	}

	return nil
}

func (f *FileCopy) copyFile(src string, dst string) error {
	_, err := os.Stat(dst)

	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(dst), 0755)
			if err != nil {
				return err
			}
		}
	}

	srcFile, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = dstFile.Write(srcFile)
	if err != nil {
		return err
	}

	return nil
}
