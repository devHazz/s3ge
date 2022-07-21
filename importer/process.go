package importer

import (
	"bytes"
	"converter"
	"fmt"
	"io"
	"log"
	"os"
	"unsafe"
	"golang.org/x/sys/windows"
	"os/exec"
	"strconv"
)
const (
	GraphicOne = 0x3427B6EF1
	GraphicTwo = 0x3427CD119
	GraphicThree = 0x3427D822D
	GraphicFour = 0x4427C21C0
)

func ModifyGraphic(handle windows.Handle, file string, address uintptr) {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Panic(err)
	}
	
	r := bytes.NewReader(data)
	var g converter.Graphic
	g.Offset = int64(address)

	r.Seek(0x44, io.SeekCurrent)
	g.HeadSize = converter.ReadBEUint32(r)

	r.Seek(0x24, io.SeekCurrent)
	g.BodySize = converter.ReadBEUint32(r)
	g.Size = g.HeadSize + g.BodySize

	// buffer := ReadProcMemory(handle, int(g.Size), address)
	// fmt.Println(string(buffer))
	r.Seek(0x15C, io.SeekStart)
	//Image Format Byte
	imgFormat, _ := r.ReadByte()
	r.Seek(0x164, io.SeekStart)
	width := converter.ReadBEUint16(r)
	height := converter.ReadBEUint16(r)
	text := data[0x238:]

	r.Seek(0x1BC, io.SeekStart)
	strBuffer := make([]byte, 8)
	br, _ := r.Read(strBuffer)
	g.Name = string(strBuffer[:br])

	fmt.Printf("Width: %d | Height: %d\n", width, height)
	fmt.Printf("Head Size: %d | Body Size: %d\n", g.HeadSize, g.BodySize)

	err = os.WriteFile(fmt.Sprintf("./files/%s.dat", g.Name), text, 0666)
	if err != nil {
		log.Panic(err)
	}
	var strFormat string
	switch imgFormat {
	case 0xA6:
		strFormat = "DXT1"
	case 0x86:
		strFormat = "DXT1"
	case 0x87:
		strFormat = "DXT3"
	case 0x88:
		strFormat = "DXT5"
	case 0xA5:
		strFormat = "87"
	}


	err = exec.Command(".\\bin\\RawtexCmd.exe", fmt.Sprintf("./files/%s.dat", g.Name), strFormat, "0", strconv.Itoa(int(width)), strconv.Itoa(int(height))).Run()
		if err != nil {
			log.Panic(err)
		}
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				log.Fatal()
			}
		}("./files/" + g.Name + ".dat")

		dds, _ := os.ReadFile("./files/"+g.Name+".dds")
	WriteProcMemory(handle, dds[0x80:], 0xAAD0, address + 0x237)
	//windows.CloseHandle(handle)
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
	var length uint32

	err := windows.WriteProcessMemory(handle, address, (*byte)(unsafe.Pointer(&data[0])), uintptr(nSize), (*uintptr)(unsafe.Pointer(&length)))
	if err != nil {
		log.Panic(err)
	}
	//   ret, _, e := write.Call(uintptr(handle), address,
    // 	uintptr(unsafe.Pointer(&data[0])),
    // 	uintptr(nSize), uintptr(unsafe.Pointer(&length)))
	// 	if ret == 0 {
	// 		log.Panic(e)
	// 	}
	windows.CloseHandle(handle)
}