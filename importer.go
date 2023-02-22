package main

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"io"
	"os"
	"os/exec"
	"syscall"

	//"strconv"
	"unsafe"
)

var graphicOffsets = [4]int{
	0x3427B6EF1,
	0x3427C2005,
	0x3427CD119,
	0x3427D822D,
}

var compressionType = map[byte]string{
	0xA6: "DXT1",
	0x86: "DXT1",
	0x87: "DXT3",
	0x88: "DXT5",
	0xA5: "87",
}

// TODO: Clean up DDS and PSG funcs
func ModifyByPSG(g Graphic, handle windows.Handle, oldData []byte, memAddress uintptr) error {
	name := g.Id + ".Texture"
	err := os.WriteFile(fmt.Sprintf("./files/%s.dat", name), oldData[0x248:], 0666)
	if err != nil {
		return err
	}
	r := bytes.NewReader(oldData)
	r.Seek(0x15C, io.SeekStart)

	//Get formats to pass along to converter
	imgFormat, _ := r.ReadByte()
	format := compressionType[imgFormat]

	//TODO: Write custom util package to convert textures
	err = exec.Command(".\\bin\\RawtexCmd.exe", fmt.Sprintf("./files/%s.dat", name), format, "0", "256", "128").Run()
	if err != nil {
		return err
	}

	//Remove data which has been passed to DDS File
	os.Remove("./files/" + name + ".dat")

	//Read the data from new DDS File
	ddsData, err := os.ReadFile("./files/" + name + ".dds")
	if err != nil {
		return err
	}

	data := ddsData
	//Check for the correct image buffer size
	if (len(data) - 0x80) == 0x8000 {
		fmt.Println("[INFO] Unusual converted file size, continuing...")
		data = append(data, oldData[0x248+0x8000:]...)
	}
	err = ModifyByDDS(data, handle, memAddress)
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
		if (file.Size-0x80) == 0xAAC0 || (file.Size-0x80) == 0x8000 { //Checking for files that have converted weirdly
			data := file.Buffer[0x80:]
			err = WriteProcMemory(handle, data, 0xAAC0, memAddress+0x247)
			if err != nil {
				return err
			}
			return nil
		} else if (file.Size - 0x80) == 0xAAD0 {
			data := file.Buffer[0x90:]
			err = WriteProcMemory(handle, data, 0xAAC0, memAddress+0x247)
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

func Attach() (windows.Handle, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, err
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	err = windows.Process32First(snapshot, &entry)
	if err != nil {
		return 0, err
	}

	var pid uint32
	for {
		err = windows.Process32Next(snapshot, &entry)
		end := 0
		for {
			if entry.ExeFile[end] == 0 {
				break
			}
			end++
		}
		processName := syscall.UTF16ToString(entry.ExeFile[:end])
		if processName == "rpcs3.exe" {
			fmt.Printf("Found %s process with pid:%d", processName, entry.ProcessID)
			pid = entry.ProcessID
			break
		}
	}

	handle, err := windows.OpenProcess(0x1F0FFF, true, pid)
	if err != nil {
		return windows.Handle(0), err
	}
	return handle, nil
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
	err := windows.WriteProcessMemory(handle, address, (*byte)(unsafe.Pointer(&data[0])), uintptr(nSize), (*uintptr)(unsafe.Pointer(&length)))
	if err != nil {
		return err
	}
	return nil
}
