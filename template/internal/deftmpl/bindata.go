// Code generated by go-bindata.
// sources:
// template/default.tmpl
// DO NOT EDIT!

package deftmpl

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templateDefaultTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x3b\x6b\x6f\xdb\x36\xbb\xdf\xf5\x2b\x9e\x69\x38\x58\x03\x58\x96\x93\x6e\xc5\xe2\xd8\x39\x70\x1d\xa5\x11\x8e\x23\x07\xb2\xd2\xae\x18\x86\x82\x96\x68\x9b\xad\x44\x6a\x24\x95\xc4\xcb\xfc\xdf\x0f\x48\xc9\x17\xd9\x72\xea\x14\x5d\xe2\xf7\x5d\x12\xb4\x91\x28\x3e\xf7\x2b\x45\xea\xfe\x1e\x22\x3c\x22\x14\x83\xf9\xe9\x13\x8a\x31\x97\x09\xa2\x68\x8c\xb9\x09\xb3\x59\x47\xdd\x5f\xe6\xf7\xf7\xf7\x80\x69\x04\xb3\x99\xb1\x15\xe4\xda\xef\x29\xa8\xfb\x7b\xa8\x3b\x77\x12\x73\x8a\xe2\x6b\xbf\x07\xb3\x99\xfd\xa3\xad\xe7\x89\xff\xe5\x38\xc4\xe4\x06\xf3\xb6\x9a\xe4\x17\x37\x39\x4c\x81\xbd\x8c\x5e\x64\xc3\xcf\x38\x94\x0a\xed\xef\x0a\x64\x20\x91\xcc\x04\xfc\x0d\x92\x5d\xa7\xe9\x1c\x94\x8c\x00\xff\xb9\x78\x68\x8e\x08\x27\x74\xac\x60\x9a\x0a\x46\x4b\x21\xea\xe7\x7a\x14\xfe\x86\x18\xd3\x55\x8a\x7f\x80\x9a\xf4\x8e\xb3\x2c\xed\xa1\x21\x8e\x45\x7d\xc0\xb8\xc4\xd1\x15\x22\x5c\xd4\xdf\xa3\x38\xc3\x8a\xe0\x67\x46\x28\x98\xa0\xb0\x42\x4e\x72\x2c\xe1\x95\xc2\x55\xef\xb2\x24\x61\x34\x07\x3e\x28\xc6\x56\xf0\x1d\xc0\x6c\xf6\xea\xfe\x1e\x6e\x89\x9c\x94\x27\xd7\x7d\x9c\xb0\x1b\x5c\xa6\xee\xa1\x04\x8b\x42\x8d\x55\xd4\x17\x8c\x1f\x2c\xae\xb6\xd8\x26\xc2\x22\xe4\x24\x95\x84\x51\xf3\x01\x1d\x4b\x7c\x27\x73\x3b\x7e\x8a\x89\x90\xc5\x54\x8e\xe8\x18\x43\x1d\x66\xb3\x9c\xaf\xa6\xb1\x1c\xdc\xd4\x93\xd2\x8a\xa5\x15\xa9\xd8\x57\x77\x6d\x58\x08\x50\x30\x96\x13\xef\x50\xca\x24\x52\x3c\x95\x50\xae\x0c\x7f\x1b\xde\x01\xcb\x78\x88\x9b\xb9\x31\x31\xc5\x1c\x49\xc6\x73\xf7\x33\x2a\x14\x55\xd2\x81\x88\x51\xf8\xa5\x1e\xe1\x11\xca\x62\x59\x97\x44\xc6\xb8\xd0\x82\xc4\x49\x1a\x23\x59\xf6\xc5\xfa\x36\x95\x97\xf1\x64\x42\x85\x40\x52\x85\xaa\x1c\x68\x3b\xe2\x1b\xa1\x38\x1e\xa2\xf0\xcb\x06\xbe\x4a\xf6\x15\x52\xf8\x1b\xbe\x36\x31\x26\xf4\xcb\xce\x1c\xa4\x1c\x2b\x67\x31\x77\x9b\xbd\x82\xff\x41\x05\xe8\xb4\xb1\x23\x07\x24\x64\x14\x27\xec\x33\xd9\x91\x07\x35\x3f\xe3\xf1\xae\x1c\x6f\x08\x57\x72\x93\x09\x49\xc3\x09\x92\x4b\x83\x70\x96\x7c\xbb\x71\xd7\xb1\x25\x58\x08\x34\x7e\x84\xe3\x95\x78\x4b\x15\xb5\x28\x93\xd3\x05\xbe\xcd\xe8\x7f\x9c\x33\x6f\x62\x0c\x63\x82\xa9\xfc\x76\x89\xb7\x61\x5c\xd6\x8d\x6f\x73\x91\x4d\xbc\x84\x0a\x89\x68\x88\x45\x05\xde\x8d\x74\xf7\x80\x56\x59\x2a\xc6\x98\x12\xfc\xed\x46\x7a\x08\xd9\xa6\x85\x8a\xea\xb0\x25\x19\x56\x96\x03\x63\xad\x18\x95\xaa\xdd\x01\x34\xc0\x9a\xcd\x8c\x7c\x10\xf2\x41\x9d\x76\x1f\xd6\x48\xb9\x64\x6a\x22\xd6\x8a\x44\x15\xf4\x7c\x2c\x58\x7c\x83\xa3\x35\x8a\xf3\xe1\xdd\x69\xce\x21\x36\xa8\x5a\xbb\xa8\x54\xe8\x2a\xf0\x78\x6f\x2a\x59\xfd\x86\x84\x92\x71\x96\x8a\x25\x5a\x89\x24\xfe\xb4\xa3\xf1\xd7\xb3\xee\x63\x5c\x79\x93\x74\xc2\x28\x91\x4c\xd9\xe1\x93\x64\x2c\x7e\x64\xf4\x95\xe4\xc2\x09\x22\xf1\x52\xa6\x65\x6b\xf5\x68\x57\x2e\x63\x9a\xc8\x44\xf3\x65\xb4\x7e\x38\xeb\x77\x83\x8f\x57\x0e\xa8\x21\xb8\xba\x7e\xdb\x73\xbb\x60\x5a\xb6\xfd\xe1\x75\xd7\xb6\xcf\x82\x33\xf8\xed\x22\xb8\xec\xc1\x61\xbd\x01\x01\x47\x54\x10\xe5\xe4\x28\xb6\x6d\xc7\x33\xc1\x9c\x48\x99\x36\x6d\xfb\xf6\xf6\xb6\x7e\xfb\xba\xce\xf8\xd8\x0e\x7c\xfb\x4e\xe1\x3a\x54\xc0\xc5\xa5\x25\x57\x20\xeb\x91\x8c\xcc\x53\xa3\xf5\x83\x65\x19\x03\x39\x8d\x31\x20\x1a\x81\x26\x12\x61\x4e\x94\x23\xa9\x34\x0d\x0a\xb5\x68\xda\xf6\x98\xc8\x49\x36\xac\x87\x2c\xb1\x95\x0c\xe3\x8c\xda\x1a\x1d\x0a\x73\x7c\x96\x16\xcd\x9a\xab\x43\x18\x86\x11\x4c\x30\x5c\xba\x01\xf4\x48\x88\xa9\xc0\xf0\xea\xd2\x0d\x0e\x0c\xa3\xcb\xd2\x29\x27\xe3\x89\x84\x57\xe1\x01\x1c\x35\x0e\x7f\x86\xcb\x1c\xa3\x61\x5c\x61\x9e\x10\x21\x08\xa3\x40\x04\x4c\x30\xc7\xc3\x29\x8c\x39\xa2\x12\x47\x35\x18\x71\x8c\x81\x8d\x20\x9c\x20\x3e\xc6\x35\x90\x0c\x10\x9d\x42\x8a\xb9\x60\x14\xd8\x50\x22\x42\x55\xdc\x21\x08\x59\x3a\x35\xd8\x08\xe4\x84\x08\x10\x6c\x24\x6f\x11\xcf\x25\x44\x42\xb0\x90\x20\x89\x23\x88\x58\x98\x25\x98\xe6\x09\x03\x46\x24\xc6\x02\x5e\xc9\x09\x06\x73\x50\x40\x98\x07\x9a\x48\x84\x51\x6c\x10\x0a\xea\xd9\xfc\x91\xee\x4a\x59\x26\x81\x63\x21\x39\xd1\x5a\xa8\x01\xa1\x61\x9c\x45\x8a\x87\xf9\xe3\x98\x24\xa4\xa0\xa0\xc0\xb5\xe0\xc2\x90\x0c\x32\x81\x6b\x9a\xcf\x1a\x24\x2c\x22\x23\xf5\x17\x6b\xb1\xd2\x6c\x18\x13\x31\xa9\x41\x44\x14\xea\x61\x26\x71\x0d\x84\x1a\xd4\x7a\xac\x29\x39\x6c\xc6\x41\xe0\x38\x36\x42\x96\x12\x2c\x40\xcb\xba\xe4\x4e\xcf\x51\xac\xa7\x4a\xa1\xb2\x50\x91\x50\x23\xb7\x13\x96\x94\x25\x21\xc2\x18\x65\x9c\x12\x31\xc1\x1a\x26\x62\x20\x98\xa6\xa8\xbc\x59\x8d\xa8\xe9\x23\x16\xc7\xec\x56\x89\x16\x32\x1a\x91\xa2\x11\xd5\x46\x46\x43\xd5\x8c\x87\x0b\xbb\x52\x26\x49\x98\xab\x5b\x1b\x20\x5d\x5a\xb5\x78\x24\x26\x28\x8e\x61\x88\x0b\x85\xe1\x08\x08\x05\xb4\x22\x0e\x57\xe4\x55\x2d\x92\x04\xc5\x90\x32\xae\xe9\xad\x8b\x59\x37\x8c\xe0\xc2\x81\x41\xff\x3c\xf8\xd0\xf1\x1d\x70\x07\x70\xe5\xf7\xdf\xbb\x67\xce\x19\x98\x9d\x01\xb8\x03\xb3\x06\x1f\xdc\xe0\xa2\x7f\x1d\xc0\x87\x8e\xef\x77\xbc\xe0\x23\xf4\xcf\xa1\xe3\x7d\x84\xff\x73\xbd\xb3\x1a\x38\xbf\x5d\xf9\xce\x60\x00\x7d\xdf\x70\x2f\xaf\x7a\xae\x73\x56\x03\xd7\xeb\xf6\xae\xcf\x5c\xef\x1d\xbc\xbd\x0e\xc0\xeb\x07\xd0\x73\x2f\xdd\xc0\x39\x83\xa0\x0f\x8a\x60\x81\xca\x75\x06\x0a\xd9\xa5\xe3\x77\x2f\x3a\x5e\xd0\x79\xeb\xf6\xdc\xe0\x63\xcd\x38\x77\x03\x4f\xe1\x3c\xef\xfb\xd0\x81\xab\x8e\x1f\xb8\xdd\xeb\x5e\xc7\x87\xab\x6b\xff\xaa\x3f\x70\xa0\xe3\x9d\x81\xd7\xf7\x5c\xef\xdc\x77\xbd\x77\xce\xa5\xe3\x05\x75\x70\x3d\xf0\xfa\xe0\xbc\x77\xbc\x00\x06\x17\x9d\x5e\x4f\x91\x32\x3a\xd7\xc1\x45\xdf\x57\xfc\x41\xb7\x7f\xf5\xd1\x77\xdf\x5d\x04\x70\xd1\xef\x9d\x39\xfe\x00\xde\x3a\xd0\x73\x3b\x6f\x7b\x4e\x4e\xca\xfb\x08\xdd\x5e\xc7\xbd\xac\xc1\x59\xe7\xb2\xf3\xce\xd1\x50\xfd\xe0\xc2\xf1\x0d\x35\x2d\xe7\x0e\x3e\x5c\x38\x6a\x48\xd1\xeb\x78\xd0\xe9\x06\x6e\xdf\x53\x62\x74\xfb\x5e\xe0\x77\xba\x41\x0d\x82\xbe\x1f\x2c\x40\x3f\xb8\x03\xa7\x06\x1d\xdf\x1d\x28\x85\x9c\xfb\xfd\xcb\x9a\xa1\xd4\xd9\x3f\x57\x53\x5c\x4f\xc1\x79\x4e\x8e\x45\xa9\x1a\x4a\x16\xe9\xfb\xfa\xfe\x7a\xe0\x2c\x10\xc2\x99\xd3\xe9\xb9\xde\xbb\x81\x02\x56\x22\xce\x27\xd7\x0d\xcb\x3a\x35\x5a\x3a\x05\xde\x25\x31\x15\xed\x8a\xc4\x76\x78\x7c\x7c\x9c\xe7\x33\x73\xb7\x49\x42\x25\xb7\xb6\x39\x62\x54\x5a\x23\x94\x90\x78\xda\x84\x9f\x2e\x70\x7c\x83\x25\x09\x11\x78\x38\xc3\x3f\xd5\x60\x31\x50\x83\x0e\x27\x28\xae\x81\x40\x54\x58\x02\x73\x32\x3a\x81\x21\xbb\xb3\x04\xf9\x4b\xf5\x00\x30\x64\x3c\xc2\xdc\x1a\xb2\xbb\x13\xd0\x48\x05\xf9\x0b\x37\xe1\xf0\xe7\xf4\xee\x04\x12\xc4\xc7\x84\x36\xa1\x71\xa2\x72\xeb\x04\xa3\xe8\x39\xe9\x27\x58\x22\x50\x2b\xa9\xb6\x79\x43\xf0\xad\x8a\x22\x53\x45\xaf\xc4\x54\xb6\xcd\x5b\x12\xc9\x49\x3b\xc2\x37\x24\xc4\x96\xbe\x79\x3e\x65\x81\x3d\x67\x57\x19\xd3\xc2\x7f\x66\xe4\xa6\x6d\x76\x73\x56\xad\x60\x9a\xe2\x15\xc6\x55\x0b\x64\x2b\xe3\x9e\xe8\x4a\x20\xb0\x6c\x5f\x07\xe7\xd6\xaf\xcf\xcc\xbe\x5e\xb6\x3d\x9f\xb9\x1f\xea\x45\x5a\xb6\x66\xee\xd4\x30\x5a\xb6\x72\x4a\x75\x31\x64\xd1\x14\x88\xc4\x89\x08\x59\x8a\xdb\xa6\xa9\x6f\xe4\x54\x5d\x17\x11\x25\xc2\x09\x4e\x90\x8e\x28\x47\x55\xf7\xcb\x79\x1f\xf7\xa4\x42\x5a\xb7\x78\xf8\x85\x48\x2b\x7f\x90\x30\x26\x27\x1a\x28\xaf\x0d\x04\x09\x1c\x2d\x27\x29\xdf\xd0\xd0\x16\x8a\x3e\x67\x42\x36\x81\x32\x8a\x4f\x60\x82\x55\x65\x6a\xc2\x61\xa3\xf1\x3f\x27\x10\x13\x8a\xad\xc5\x50\xfd\x0d\x4e\x4e\x40\x47\x40\x3e\x01\x7e\x20\x89\x0a\x16\x44\xe5\x09\x0c\x51\xf8\x65\xcc\x59\x46\x23\x2b\x64\x31\xe3\x4d\xf8\x71\xf4\x46\xfd\xae\xaa\x1f\x52\x14\x45\x9a\x2b\xe5\x0d\xc3\xb1\x9e\xd9\x36\x8b\x99\xa6\xd2\xb7\x44\xc3\xa7\x76\x8f\x15\x91\x76\x94\xa3\x92\x77\x80\x96\xe4\xcf\x98\xc7\x00\x14\x07\x4f\x9c\x49\x6f\x30\x57\x48\x62\x0b\xc5\x64\x4c\x9b\x20\x59\x5a\x56\xd4\x8d\x7e\xd0\x36\x25\x4b\xcd\xd3\x96\x2d\xa3\x25\xa3\x79\x66\x35\xdf\x34\x1a\x4f\x1c\x2a\x95\x4c\x47\x44\xa4\x31\x9a\x36\x61\x18\xb3\xf0\x4b\xc9\xb7\x13\x74\x67\x15\x4e\xf2\xa6\xd1\x48\xef\x4a\x0f\xc3\x18\x23\xae\x08\xca\x49\x69\x7c\x5b\xa0\x2c\x94\x03\x28\x93\x6c\x2d\x24\x4a\xda\xd2\x8a\x02\x68\x45\xe4\xe6\xa9\xdd\xaa\x2c\xef\xba\x72\x1e\x16\x62\xce\xb7\x32\xb2\x0e\xe6\xc2\xce\x4a\x13\x26\x84\x38\x8e\x8b\xd9\x6d\xb3\x91\xdf\x8b\x14\x85\xf3\xfb\x27\x15\xb4\x78\xc8\x51\x44\x32\xd1\x84\xd7\x7a\xac\x22\x01\x8c\x46\xa5\x2c\x96\x83\x35\xe1\x30\xbd\x03\xc1\x62\x12\xc1\x8f\xf8\x58\xfd\x96\x13\xc3\x68\xb4\xa2\x8b\x7d\xc8\x0e\x4b\x4e\x9e\x2e\x4b\xbc\xd9\x1a\x70\x25\xed\x6a\x90\xdb\xa2\xd4\xfc\xd2\x68\x9c\x80\x2e\x51\xc5\xfc\x10\x53\x89\x79\x95\xbd\xf4\xbf\x86\x36\xca\xa6\xdd\x9c\x37\xbf\x1c\x1d\x75\xab\x0b\xd0\x91\xf2\x6b\x13\x8a\x78\xcb\x09\xac\x5a\x2f\x87\xad\x8e\xc8\xf9\xcf\x72\xf7\x67\xb1\xed\x03\xfa\x6d\x49\xe5\x3b\xac\x03\x38\x84\xd9\x4c\x2c\x5e\x78\xc0\x88\x71\x58\xee\x50\x6c\xd9\x21\x82\xd9\x6c\x8d\x2a\xac\xee\x57\xb4\x4b\xbb\x15\x1b\xd3\x8a\x57\x2b\x25\xe3\x2f\x72\xf0\xe2\x9e\xbf\xb8\xe9\x2e\xc5\x6c\xe9\x3c\x87\xb9\xf3\x3c\xe4\x1b\x7b\x9f\xfb\xb6\xaa\x7d\xbf\x9c\x60\xdf\x5d\xa1\x01\x8d\x79\x2e\x79\xc8\x1d\x0a\x31\x10\x4c\x38\x1e\xb5\xcd\x5d\x5e\xe0\x3e\xb1\x3f\xcc\x93\xe6\xf9\xf9\x79\x91\x7c\x23\x1c\x32\xae\xdf\xc9\xcd\x97\x07\xa5\x05\xc1\x91\x5a\x0e\x94\xf2\xf6\x90\xc5\x51\x75\xe2\x0e\x33\x2e\x14\xf6\x94\x91\x7c\x60\xd1\x50\x10\xaa\x91\x16\x7d\xc5\x5a\x82\xff\x45\x31\xa6\xf1\xe9\x97\xa8\x23\xc6\x93\x26\x84\x28\x25\x12\xc5\xe4\x2f\x5c\x99\xf4\x5f\xff\xfc\x2b\x8e\x50\x45\xbd\xde\x98\x51\x0c\x6b\x2d\x37\xf3\x42\xbe\x18\x5c\x74\x6f\xe9\x5d\x61\xde\xd3\xf7\x04\xdf\x02\xa1\x0f\xbd\x7c\x9f\x2f\x23\x51\xa5\x0f\xaf\x25\xde\xea\xf4\x9b\xff\x7c\x6d\xd3\xa5\xa2\x28\xbc\x84\xec\x3f\x13\xb2\x42\x72\x46\xc7\xcf\xa7\xda\xdf\xb7\x9f\x31\xf9\xa3\xd8\x71\x6b\xd9\x39\x93\xdf\xc1\xeb\x2a\x1a\x86\xe2\xc9\xfc\x20\xc5\xfa\xd6\xdd\x8b\x1f\xfe\x3b\xfc\x30\x6f\x4d\x17\xae\xd6\x1a\x3e\x9f\x99\xc1\xae\xd6\xd1\x57\x4e\x10\x6d\x3f\xe6\xf3\xcc\xc2\x6c\x8f\x3b\xa8\xa8\x05\xcb\xcd\xfb\xbc\x12\x3c\xbb\x67\xac\x70\xb4\x2f\xee\xf1\x55\x8d\x7e\xf5\x58\xd8\x7f\xa8\xb3\xac\x76\x98\xeb\xe7\xd4\x9e\xa9\xa1\x9c\xb7\x5b\x1b\x3d\x65\x46\x23\xcc\x55\xf7\x57\x76\xa7\xfc\xa4\x9d\x6a\xa2\xf6\x2f\xc7\x7c\x5b\x35\xdd\xb1\xbd\x5b\x3d\xe3\x52\x69\xde\x97\xae\x70\x6f\xaa\xf1\xde\x79\x26\x40\x6b\xb2\x87\x3c\xed\x9d\x9e\x1e\x13\xc1\x0f\x75\xc4\x2f\x81\xf5\xdf\xd9\xe6\xae\x2e\xb7\x16\x67\x05\x97\x0b\xae\xf9\xd0\x33\x2c\xb9\x56\x4f\x2e\xbe\x78\xe3\xbf\xc3\x1b\x5f\x16\x5d\x2f\x8b\xae\x97\x45\xd7\xbe\x3b\xcb\xcb\xa2\x6b\x6f\x5a\xb6\x6d\x86\x6a\xd9\x7a\x3f\xee\xf4\x11\x5b\xa1\x0b\x90\xe5\xc8\x93\x9f\xc4\x28\x1d\x4d\x5a\x39\x69\xb2\x34\xf4\xf1\xf1\xf1\x43\x1b\xdc\xe5\x9d\xdd\xcd\x2d\xc9\xfd\x68\x1a\xf6\xa9\x7d\x79\xca\xd6\xe5\x68\x6b\xeb\x52\xb9\x89\xf6\x35\x93\xaf\xf4\x36\x6b\xe7\x1a\xca\xa7\xb0\x56\xd3\x55\xf9\x4b\xda\xa7\x73\x88\xa3\xd5\x6c\xa5\x25\xda\x39\x55\x61\x2a\x61\x38\xdd\x6d\x1f\x6e\x33\x77\x6c\x9c\x77\x58\xcf\x0c\x2d\x3b\x22\x37\xa7\xf9\xff\x46\x39\x4d\xec\x5b\x5b\xbb\xe5\x78\x5d\x2e\xe2\x32\x7f\xb5\xec\x21\x8b\xa6\x6a\x64\x22\x93\xf8\xd4\x30\xaa\x3f\xd5\x4d\x33\x31\x61\x37\x98\x7f\x87\x2f\x55\x37\x50\x95\xbf\x6d\xfa\x27\xbe\x43\xfb\x3e\x9f\xa1\xed\xfe\x15\xda\xf7\xfb\x08\x6d\x85\xe6\x0e\x9a\x5c\x7e\x6e\xfa\x88\x4f\xc0\xfe\x3f\x00\x00\xff\xff\x3c\xf7\xad\xcb\x87\x3f\x00\x00")

func templateDefaultTmplBytes() ([]byte, error) {
	return bindataRead(
		_templateDefaultTmpl,
		"template/default.tmpl",
	)
}

func templateDefaultTmpl() (*asset, error) {
	bytes, err := templateDefaultTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "template/default.tmpl", size: 16263, mode: os.FileMode(420), modTime: time.Unix(1, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"template/default.tmpl": templateDefaultTmpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"template": &bintree{nil, map[string]*bintree{
		"default.tmpl": &bintree{templateDefaultTmpl, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
