package importer

import (
	"bytes"
	"converter"
	"io"
	"log"
	"os"
	"unsafe"
	"golang.org/x/sys/windows"
)
const (
	GraphicOne = 0x3427B6EF1
	GraphicTwo = 0x3427CD119
	GraphicThree = 0x3427D822D
	GraphicFour = 0x4427C21C0
)

func ModifyGraphic(handle windows.Handle, file string, nsize int, address uintptr) {
	//We make two memory reads, one to get the header and body size of the graphic, two to read the whole graphic to a buffer
	data, err := os.ReadFile(file)
	if err != nil {
		log.Panic(err)
	}
	
	buffer := ReadProcMemory(handle, nsize, address)
	r := bytes.NewReader(buffer)
	var g converter.Graphic
	g.Offset = int64(address)

	r.Seek(0x43, io.SeekStart)
	g.HeadSize = converter.ReadBEUint32(r)

	r.Seek(0x6B, io.SeekStart)
	g.BodySize = converter.ReadBEUint32(r)
	g.Size = g.HeadSize + g.BodySize

	WriteProcMemory(handle, data[g.HeadSize:g.Size], int(g.BodySize), address + uintptr(g.HeadSize))
	windows.CloseHandle(handle)
}

func Attach(pid int) (windows.Handle) {
	handle, err := windows.OpenProcess(0x1F0FFF, false, uint32(pid))
	if err != nil {
		log.Panic(err)
	}
	return handle
}

func ReadProcMemory(handle windows.Handle, nSize int, address uintptr) ([]byte) {
	read := windows.MustLoadDLL("kernel32.dll").MustFindProc("ReadProcessMemory")
	var data = make([]byte, nSize)
	var length uint32

	  ret, _, e := read.Call(uintptr(handle), address,
    	uintptr(unsafe.Pointer(&data[0])),
    	uintptr(nSize), uintptr(unsafe.Pointer(&length)))
		if ret == 0 {
			log.Panic(e)
		}
	return data
}


func WriteProcMemory(handle windows.Handle, data []byte,  nSize int, address uintptr) {
	write := windows.MustLoadDLL("kernel32.dll").MustFindProc("WriteProcessMemory")
	var length uint32

	  ret, _, e := write.Call(uintptr(handle), address,
    	uintptr(unsafe.Pointer(&data[0])),
    	uintptr(nSize), uintptr(unsafe.Pointer(&length)))
		if ret == 0 {
			log.Panic(e)
		}
	windows.CloseHandle(handle)
}