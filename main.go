package main

import (
	"flag"
	"fmt"
	human "github.com/dustin/go-humanize"
	"syscall"
	"unsafe"
)

var MAX = flag.Uint64("max", 2048, "-max=2048 value is in MB")
var MIN = flag.Uint64("min", 0, "-min=2048 value is in MB")

const (
	FILE_CACHE_MAX_HARD_ENABLE = 0x1
	FILE_CACHE_MIN_HARD_ENABLE = 0x4
)

func abort(funcname string, err error) {
	panic(fmt.Sprintf("%s failed: %v", funcname, err))
}

func GetSizes() string {
	kernel32, err := syscall.LoadDLL("kernel32.dll")

	if nil != err {
		abort("loadLibrary", err)
	}
	defer kernel32.Release()
	get, err := kernel32.FindProc("GetSystemFileCacheSize")
	if nil != err {
		abort("GetProcAddress", err)
	}
	var minFileCache uint64
	var maxFileCache uint64
	var lpFlags uint32
	res, _, err := get.Call(uintptr(unsafe.Pointer(&minFileCache)), uintptr(unsafe.Pointer(&maxFileCache)), uintptr(unsafe.Pointer(&lpFlags)))
	if res == 0 {
		abort("getSystemFileCacheSize", err)
	}
	return fmt.Sprintf("Min: %v Max: %v Flags: %v", human.Bytes(minFileCache), human.Bytes(maxFileCache), lpFlags)

}

func setSizes() {
	kernel32, err := syscall.LoadDLL("kernel32.dll")
	if nil != err {
		abort("loadLibrary", err)
	}
	defer kernel32.Release()
	set, err := kernel32.FindProc("SetSystemFileCacheSize")
	if nil != err {
		abort("GetProcAddress", err)
	}
	var lpFlags uint32
	lpFlags = FILE_CACHE_MAX_HARD_ENABLE
	max := *MAX * uint64(1000) * uint64(1000)
	min := *MIN * uint64(1000) * uint64(1000)
	res, _, err := set.Call(uintptr(min), uintptr(max), uintptr(lpFlags))
	if res == 0 {
		abort("SetSystemFileCacheSize", err)
	}
}

func main() {
	flag.Parse()

	fmt.Println("BEFORE: ", GetSizes())
	setSizes()
	fmt.Println("AFTER: ", GetSizes())

}
