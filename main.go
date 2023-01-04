package main

import (
	"birdnest/printdata"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Report struct {
	XMLName xml.Name `xml:"report"`

	DeviceInformation struct {
		DeviceId         string `xml:"deviceId,attr"`
		ListenRange      string `xml:"listenRange"`
		DeviceStarted    string `xml:"deviceStarted"`
		UptimeSeconds    string `xml:"uptimeSeconds"`
		UpdateIntervalMs string `xml:"updateIntervalMs"`
	} `xml:"deviceInformation"`
	Capture struct {
		SnapshotTimestamp string `xml:"snapshotTimestamp,attr"`
		Drone             []struct {
			SerialNumber string `xml:"serialNumber"`
			Model        string `xml:"model"`
			Manufacturer string `xml:"manufacturer"`
			Mac          string `xml:"mac"`
			Ipv4         string `xml:"ipv4"`
			Ipv6         string `xml:"ipv6"`
			Firmware     string `xml:"firmware"`
			PositionY    string `xml:"positionY"`
			PositionX    string `xml:"positionX"`
			Altitude     string `xml:"altitude"`
		} `xml:"drone"`
	} `xml:"capture"`
}

func main() {

	file, err := os.Create("database.db")
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	tick := time.Tick(2 * time.Second)
	for range tick {
		log.Println("Update")
		readXmlApi(file)
	}

}

func readXmlApi(file *os.File) {
	// Send a GET request to the API endpoint
	resp, err := http.Get("https://assignments.reaktor.com/birdnest/drones")
	if err != nil {
		// Handle the error
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Handle the error
		return
	}
	var xmlData Report
	err = xml.Unmarshal(body, &xmlData)
	if err != nil {
		// Handle the error
		return
	}
	fmt.Println(xmlData.Capture.SnapshotTimestamp)
	for _, x := range xmlData.Capture.Drone {
		X, err := strconv.ParseFloat(x.PositionX, 64)
		if err != nil {
			// Handle the error
			fmt.Println(err)
			return
		}
		Y, err := strconv.ParseFloat(x.PositionY, 64)
		if err != nil {
			// Handle the error
			fmt.Println(err)
			return
		}

		// Round the float to the nearest int
		if X <= 35000 && X >= 15000 {
			//text := fmt.Sprintf("SeriaNumb %v\nLocation was X: %v Y: %v\n", x.SerialNumber, X, Y)
			printdata.NDZData(x.SerialNumber, X, Y, file)
		} else if Y <= 35000 && Y >= 15000 {
			//text := fmt.Sprintf("SeriaNumb %v\nLocation was X: %v Y: %v\n", x.SerialNumber, X, Y)
			printdata.NDZData(x.SerialNumber, X, Y, file)
		}
	}

}
