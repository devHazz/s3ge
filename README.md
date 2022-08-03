# Welcome to the project!

## About
This program gives you the ability to import custom graphics into the Skate 3 RPCS3 Process but also allows you to pull any existing graphics on a `SKATER.P` file.

## How-to use
### Importer
Either make sure you have a DDS file (with mipmaps generated) or a PSG File. Files must be of 256x128 size. <br>
Run the importer.exe, drag & drop in the file and run through the prompts. <br>
Need help with retrieving the Skate 3 PID? <br> 
`Open Task Manager > Find the rpcs3.exe in the Details Tab > Locate the PID next to the process name`

### Converter
Make sure your skater has 4 exisiting graphics, otherwise this will NOT work! <br>
Make sure that your `SKATER.P` has been copied to the files directory, then execute `go run main.go`

## Requirements
 - Windows ONLY
 - Go Version 1.18.x
 - PS3 `SKATER.P` File (Can be retrieved through RPCS3)
 - A skater with 4 existing graphics (Won't be a requirement in future builds)
