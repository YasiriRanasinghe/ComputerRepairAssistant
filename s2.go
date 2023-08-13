package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {

	banner := `

	██████╗ ███████╗██████╗ ███████╗ ██████╗ ███╗   ██╗ █████╗ ██╗          ██████╗ ██████╗ ███╗   ███╗██████╗ ██╗   ██╗████████╗███████╗██████╗     
	██╔══██╗██╔════╝██╔══██╗██╔════╝██╔═══██╗████╗  ██║██╔══██╗██║         ██╔════╝██╔═══██╗████╗ ████║██╔══██╗██║   ██║╚══██╔══╝██╔════╝██╔══██╗    
	██████╔╝█████╗  ██████╔╝███████╗██║   ██║██╔██╗ ██║███████║██║         ██║     ██║   ██║██╔████╔██║██████╔╝██║   ██║   ██║   █████╗  ██████╔╝    
	██╔═══╝ ██╔══╝  ██╔══██╗╚════██║██║   ██║██║╚██╗██║██╔══██║██║         ██║     ██║   ██║██║╚██╔╝██║██╔═══╝ ██║   ██║   ██║   ██╔══╝  ██╔══██╗    
	██║     ███████╗██║  ██║███████║╚██████╔╝██║ ╚████║██║  ██║███████╗    ╚██████╗╚██████╔╝██║ ╚═╝ ██║██║     ╚██████╔╝   ██║   ███████╗██║  ██║    
	╚═╝     ╚══════╝╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═╝╚══════╝     ╚═════╝ ╚═════╝ ╚═╝     ╚═╝╚═╝      ╚═════╝    ╚═╝   ╚══════╝╚═╝  ╚═╝    
																																					 
	██████╗ ███████╗██████╗  █████╗ ██╗██████╗      █████╗ ███████╗███████╗██╗███████╗████████╗ █████╗ ███╗   ██╗████████╗                           
	██╔══██╗██╔════╝██╔══██╗██╔══██╗██║██╔══██╗    ██╔══██╗██╔════╝██╔════╝██║██╔════╝╚══██╔══╝██╔══██╗████╗  ██║╚══██╔══╝                           
	██████╔╝█████╗  ██████╔╝███████║██║██████╔╝    ███████║███████╗███████╗██║███████╗   ██║   ███████║██╔██╗ ██║   ██║                              
	██╔══██╗██╔══╝  ██╔═══╝ ██╔══██║██║██╔══██╗    ██╔══██║╚════██║╚════██║██║╚════██║   ██║   ██╔══██║██║╚██╗██║   ██║                              
	██║  ██║███████╗██║     ██║  ██║██║██║  ██║    ██║  ██║███████║███████║██║███████║   ██║   ██║  ██║██║ ╚████║   ██║                              
	╚═╝  ╚═╝╚══════╝╚═╝     ╚═╝  ╚═╝╚═╝╚═╝  ╚═╝    ╚═╝  ╚═╝╚══════╝╚══════╝╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═══╝   ╚═╝                              
																																					 

  SLITT POC Project | V0.1                                                                   
`

	fmt.Println(banner)
	time.Sleep(3 * time.Second)
	fmt.Println("Gathering Data ...")

	// Run systeminfo command to get hardware and BIOS information
	cmd := exec.Command("systeminfo")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to execute systeminfo command:", err)
		return
	}

	// Parse the output to extract relevant information
	hardwareSpecs := parseSystemInfoOutput(string(output))

	// Get pending software updates
	pendingUpdates := getPendingUpdates()

	// Create a new file to save the extracted information
	file, err := os.Create("hardware_specs.txt")
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	// Write the extracted information to the file
	file.WriteString(hardwareSpecs)

	// Write pending software updates to the file
	file.WriteString("\n## Pending Software Windows Updates\n")
	for _, update := range pendingUpdates {
		file.WriteString(update + "\n")
	}

	fmt.Println("Hardware specs, BIOS information, and pending software updates extracted and saved to hardware_specs.txt.")
	time.Sleep(3 * time.Second)
}

func parseSystemInfoOutput(output string) string {

	// Get serial number
	SerialNumber := getSerialNumber()

	var hardwareSpecs strings.Builder

	// Extract Manufacturer, Model, and Serial Number from the output
	manufacturer := extractValue(output, "System Manufacturer:")
	model := extractValue(output, "System Model:")
	biosVersion := extractValue(output, "BIOS Version:")
	osname := extractValue(output, "OS Name:")
	osversion := extractValue(output, "OS Version:")

	// Append the extracted information to the hardwareSpecs string
	hardwareSpecs.WriteString("## Hardware Specs\n")
	hardwareSpecs.WriteString("Manufacturer: " + manufacturer + "\n")

	for _, update := range SerialNumber {
		hardwareSpecs.WriteString("Serial Number: " + update + "\n")
	}

	//hardwareSpecs.WriteString("Serial Number: " + SerialNumber + "\n")
	hardwareSpecs.WriteString("Model: " + model + "\n\n")

	hardwareSpecs.WriteString("## BIOS Information\n")
	hardwareSpecs.WriteString("BIOS Version: " + biosVersion + "\n\n")

	hardwareSpecs.WriteString("## OS Specs\n")
	hardwareSpecs.WriteString("OS Name: " + osname + "\n")
	hardwareSpecs.WriteString("OS Version: " + osversion + "\n\n")

	return hardwareSpecs.String()
}

func getPendingUpdates() []string {
	cmd := exec.Command("powershell", "-Command", "Get-WmiObject -Class Win32_QuickFixEngineering | Select-Object -ExpandProperty HotFixID")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to execute PowerShell command:", err)
		return nil
	}

	updates := strings.Split(strings.TrimSpace(string(output)), "\r\n")
	return updates
}

func getSerialNumber() []string {
	cmd := exec.Command("powershell", "-Command", "(Get-WmiObject -Class Win32_BIOS).SerialNumber")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to execute PowerShell command:", err)
		return nil
	}

	updates := strings.Split(strings.TrimSpace(string(output)), "\r\n")
	return updates
}

func extractValue(output, key string) string {
	startIndex := strings.Index(output, key)
	if startIndex == -1 {
		return ""
	}

	startIndex += len(key)
	endIndex := strings.Index(output[startIndex:], "\n")
	if endIndex == -1 {
		return ""
	}

	return strings.TrimSpace(output[startIndex : startIndex+endIndex])
}
