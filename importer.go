package main

import (
	"bufio"
	"bytes"
	"converter"
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
		//i -= 1
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
		b, _ := os.ReadFile(clean)
		graphic := converter.ReadGraphicFromPSG(bytes.NewReader(b), b)
		switch n {
		case 1:
			err = importer.ModifyByPSG(graphic, h, b, importer.GraphicOne)
			if err != nil {
				log.Panic(err)
			}
		case 2:
			err = importer.ModifyByPSG(graphic, h, b, importer.GraphicTwo)
			if err != nil {
				log.Panic(err)
			}
		case 3:
			err = importer.ModifyByPSG(graphic, h, b, importer.GraphicThree)
			if err != nil {
				log.Panic(err)
			}
		case 4:
			err = importer.ModifyByPSG(graphic, h, b, importer.GraphicFour)
			if err != nil {
				log.Panic(err)
			}
		default:
			log.Panic("[ERROR] Incorrect graphic option")
		}

		//importer.HandlePSG(choice, clean)
	case strings.HasSuffix(clean, ".dds"):
		n, pid := Select(reader)
		h := importer.Attach(pid)
		buf, _ := os.ReadFile(clean)
		switch n {
		case 1:
			err = importer.ModifyByDDS(buf, h, importer.GraphicOne)
			if err != nil {
				log.Panic(err)
			}
		case 2:
			err = importer.ModifyByDDS(buf, h, importer.GraphicTwo)
			if err != nil {
				log.Panic(err)
			}
		case 3:
			err = importer.ModifyByDDS(buf, h, importer.GraphicThree)
			if err != nil {
				log.Panic(err)
			}
		case 4:
			err = importer.ModifyByDDS(buf, h, importer.GraphicFour)
			if err != nil {
				log.Panic(err)
			}
		default:
			log.Panic("[ERROR] Incorrect graphic option") 
		}
	default:
		log.Fatal("[FATAL] File is not of DDS/PSG type")
	}
}