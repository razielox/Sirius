package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

	sirius "github.com/0sm0s1z/Sirius-Scan/Engine/lib"
	siriusNmap "github.com/0sm0s1z/Sirius-Scan/Engine/lib/nmap"
	"github.com/lair-framework/go-nmap"
)


func main() {
	log.Println("asdf")

	var scanRequest sirius.ScanRequest
	scanRequest.ScanID = "1"

	VulnerabilityScanner(scanRequest)

}


// VulnerabilityScanner subscribes to the queue and listens for scan requests
// When a scan request is received, it will execute scans for each target up to the scan queue
func VulnerabilityScanner(scanRequest sirius.ScanRequest) {
	//Get target from scanRequest (last host in the list)
	//target := scanRequest.ScanReport.ScanResults[len(scanRequest.ScanReport.ScanResults)-1].IP

	//target := "10.0.50.20"

	//Execute Nmap Scan
	rawScanResults := "./nmap-tmp.xml"

	//log.Println(string(cmd))

	//Process Nmap Scan Results
	dat, err := os.ReadFile(rawScanResults)
	if err != nil {
		log.Println(err)
	}
	//Parse Nmap XML
	svdbHost, err := parseNmapXML(dat)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	cveList := processScanResults(dat)

	//Create a SVDBHost
	svdbHost.CVE = cveList

	//Clear Scan Results
	scanRequest.ScanReport.ScanResults = nil

	//Append SVDBHost to ScanReport
	scanRequest.ScanReport.ScanResults = append(scanRequest.ScanReport.ScanResults, *svdbHost)
	scanRequest.Command = "complete"

	//Send ScanReport to Queue
	log.Println(scanRequest)
	//sirius.SendToQueue(scanRequest, "scan")
	//sirius.SendToQueue(scanRequest, "scan-report")
}

// processScanResults processes the raw scan results and returns a list of CVEs
func processScanResults(dat []byte) []string {
	//Parse XML Using Lair Project's Nmap Parser
	var nmapResults []siriusNmap.CVE
	nmapResults = siriusNmap.ProcessReport(dat)

	//Create CVEList
	var cveList []string
	for _, cve := range nmapResults {
		newCVE := "CVE-" + cve.CVEID
		cveList = append(cveList, newCVE)
	}

	return cveList
}

func parseNmapXML(data []byte) (*sirius.SVDBHost, error) {
	var nmapRun nmap.NmapRun
	if err := xml.Unmarshal(data, &nmapRun); err != nil {
		return nil, fmt.Errorf("unable to unmarshal XML data: %v", err)
	}

	if len(nmapRun.Hosts) == 0 {
		return nil, fmt.Errorf("no hosts found in the nmap XML data")
	}

	host := nmapRun.Hosts[0]
	var ip string
	for _, address := range host.Addresses {
		if address.AddrType == "ipv4" || address.AddrType == "ipv6" {
			ip = address.Addr
			break
		}
	}

	var osName, osVersion string
	if len(host.Os.OsMatches) > 0 {
		osMatch := host.Os.OsMatches[0]
		osName = osMatch.Name
		osVersion = osMatch.OsClasses[0].OsGen
	}

	svdbHost := &sirius.SVDBHost{
		IP:        ip,
		//Hostname:  host.Hostnames[0].Name,
		OS:        osName,
		OSVersion: osVersion,
	}

	// Parse the services, CPEs, CVEs, and agent data as needed
	// and fill in the corresponding fields in the SVDBHost struct

	return svdbHost, nil
}