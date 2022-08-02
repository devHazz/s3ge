package importer

import (
	"bytes"
	"converter"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	//"strconv"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	GraphicOne   = 0x3427B6EF1
	GraphicTwo   = 0x3427C2005
	GraphicThree = 0x3427CD119
	GraphicFour  = 0x3427D822D
	BodySize     = 0xAAC0
)

func getCompressionType(format byte) string {
	var strFormat string
	switch format {
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
	default:
		log.Panic("Compression type not valid")
	}
	return strFormat

}

func ModifyByPSG(g converter.Graphic, handle windows.Handle, oldData []byte, memAddress uintptr) error {
	err := os.WriteFile(fmt.Sprintf("./files/%s.dat", g.Name), oldData[0x248:], 0666)
	if err != nil {
		return err
	}
	var strFormat string
	r := bytes.NewReader(oldData)
	r.Seek(0x15C, io.SeekStart)
	imgFormat, _ := r.ReadByte()
	strFormat = getCompressionType(imgFormat)
	err = exec.Command(".\\bin\\RawtexCmd.exe", fmt.Sprintf("./files/%s.dat", g.Name), strFormat, "0", "256", "128").Run()
	if err != nil {
		return err
	}

	os.Remove("./files/" + g.Name + ".dat")

	ddsData, err := os.ReadFile("./files/" + g.Name + ".dds")
	if err != nil {
		return err
	}
	g.Buffer = ddsData
	if (len(g.Buffer) - 0x80) == 0x8000 {
		fmt.Println("[INFO] Unusual converted file size, continuing...")
		g.Buffer = append(g.Buffer, oldData[0x248 + 0x8000:]...)
	}
	//Data to make up for texture bleed
	err = ModifyByDDS(g.Buffer, handle, memAddress)
	if err != nil {
		return err
	}
	return nil
}

func ModifyByDDS(buffer []byte, handle windows.Handle, memAddress uintptr) error {
	file, err := DDSParse(buffer)
	if err != nil {
		return err
	}
	if file.Height == 128 && file.Width == 256 {
		if (file.Size - 0x80) == 0xAAC0 || (file.Size - 0x80) == 0x8000 { //Checking for files that have converted weirdly
			//Correct DDS File Length
			data := file.Buffer[0x80:]
			err = WriteProcMemory(handle, data, BodySize, memAddress+0x247)
			if err != nil {
				return err
			}
			return nil
		} else if (file.Size - 0x80) == 0xAAD0 {
			data := file.Buffer[0x90:]
			err = WriteProcMemory(handle, data, BodySize, memAddress+0x247)
			if err != nil {
				return err
			}
			return nil
		} else {
			return errors.New("graphic body size has incorrect read length")
		}
	} else {
		return errors.New("the current dimensions do not meet 256x128")
	}
}

func Attach(pid int) windows.Handle {
	handle, err := windows.OpenProcess(0x1F0FFF, true, uint32(pid))
	if err != nil {
		log.Panic(err)
	}
	return handle
}

func ReadProcMemory(handle windows.Handle, nSize int, address uintptr) ([]byte, error) {
	var data = make([]byte, nSize)
	var length uint32
	err := windows.ReadProcessMemory(handle, address, (*byte)(unsafe.Pointer(&data[0])), uintptr(nSize), (*uintptr)(unsafe.Pointer(&length)))
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("no buffer data found")
	}
	return data, nil
}

func WriteProcMemory(handle windows.Handle, data []byte, nSize int, address uintptr) error {
	var length uint32
	//var oldProtect uint32
	// err := windows.VirtualProtectEx(handle, address, uintptr(nSize), windows.PAGE_EXECUTE_READWRITE, &oldProtect)
	// if err != nil {
	// 	return err
	// }
	err := windows.WriteProcessMemory(handle, address, (*byte)(unsafe.Pointer(&data[0])), uintptr(nSize), (*uintptr)(unsafe.Pointer(&length)))
	if err != nil {
		return err
	}
	// err = windows.VirtualProtectEx(handle, address, uintptr(nSize), oldProtect, nil)
	// if err != nil {
	// 	return err
	// }
	return nil
}
