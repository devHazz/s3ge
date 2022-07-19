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

func Select(reader *bufio.Reader) (int) {
	fmt.Println("Which graphic would you like to assign this image to?")
	text, _ := reader.ReadString('\n')
	text = CleanLineString(text)
	if text == "1" || text == "2" || text == "3" || text == "4" {
		fmt.Printf("You have selected graphic: %s\n", text)
		i, _ := strconv.Atoi(text)
		return i - 1
	} else {
		log.Panic("[ERROR] Incorrect graphic option")
	}
	return 0
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
		choice := Select(reader)
		importer.HandlePSG(choice, clean)
	case strings.HasSuffix(clean, ".dds"):
		//choice := Select(reader)
	default:
		log.Fatal("[FATAL] File is not of DDS/PSG type")
	}
}