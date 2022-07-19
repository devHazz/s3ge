package main

import (
	"bufio"
	"fmt"
	"importer"
	"log"
	"os"
	"strconv"
	"strings"
)

func CleanLineString(text string) (string) {
	clean := text
	clean = strings.TrimSuffix(clean, "\n")
	clean = strings.TrimSuffix(clean, "\r")
	return clean
}

func Select(reader *bufio.Reader) (int, int) {
	fmt.Println("Which graphic would you like to assign this image to? (1-4)")
	text, _ := reader.ReadString('\n')
	text = CleanLineString(text)
	if text == "1" || text == "2" || text == "3" || text == "4" {
		fmt.Printf("You have selected graphic: %s\n", text)
		i, _ := strconv.Atoi(text)
		i -= 1
		fmt.Println("Please enter the Skate 3 PID: ")

		pid, _ := reader.ReadString('\n')
		pid = CleanLineString(pid)
		id, _ := strconv.Atoi(pid)
		return i, id
	} else {
		log.Panic("[ERROR] Incorrect graphic option")
	}
	return 0, 0
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(`
	
██╗    ██╗      ██████╗ ██╗   ██╗███████╗    ███████╗██╗  ██╗ █████╗ ████████╗███████╗
██║    ██║     ██╔═══██╗██║   ██║██╔════╝    ██╔════╝██║ ██╔╝██╔══██╗╚══██╔══╝██╔════╝
██║    ██║     ██║   ██║██║   ██║█████╗      ███████╗█████╔╝ ███████║   ██║   █████╗  
██║    ██║     ██║   ██║╚██╗ ██╔╝██╔══╝      ╚════██║██╔═██╗ ██╔══██║   ██║   ██╔══╝  
██║    ███████╗╚██████╔╝ ╚████╔╝ ███████╗    ███████║██║  ██╗██║  ██║   ██║   ███████╗
╚═╝    ╚══════╝ ╚═════╝   ╚═══╝  ╚══════╝    ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   ╚══════╝
<3 by Hazz


Please drag & drop your PSG/DDS Graphic into here`)

	text, _ := reader.ReadString('\n')
	clean := strings.Replace(text, `"`, "", -1)
	clean = CleanLineString(clean)

	_, err := os.Stat(clean)
	if err != nil {
		log.Fatal("[ERROR] Graphic not valid")
	}
	switch {
	case strings.HasSuffix(clean, ".psg"):
		n, pid := Select(reader)
		h := importer.Attach(pid)
		switch n {
		case 1:
			importer.ModifyGraphic(h, clean, 0x70, importer.GraphicOne)
		case 2:
			importer.ModifyGraphic(h, clean, 0x70, importer.GraphicTwo)
		case 3:
			importer.ModifyGraphic(h, clean, 0x70, importer.GraphicThree)
		case 4:
			importer.ModifyGraphic(h, clean, 0x70, importer.GraphicFour)
		default:
			log.Panic("[ERROR] Incorrect graphic option")
		}

		//importer.HandlePSG(choice, clean)
	case strings.HasSuffix(clean, ".dds"):
		//choice := Select(reader)
	default:
		log.Fatal("[FATAL] File is not of DDS/PSG type")
	}
}