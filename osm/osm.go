package osm

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/qedus/osmpbf"
)

type Address struct {
	Street   string
	City     string
	Postcode string
	Hash     string
}

func hasAddressTags(tags map[string]string) bool {
	hasCity := false
	hasStreet := false
	hasPostcode := false

	for key := range tags {
		switch key {
		case "addr:city":
			hasCity = true
		case "addr:street":
			hasStreet = true
		case "addr:postcode":
			hasPostcode = true
		}
	}
	return hasCity && hasStreet && hasPostcode
}

func addressFromOSMTags(tags map[string]string) Address {
	return Address{
		Street:   tags["addr:street"],
		City:     tags["addr:city"],
		Postcode: tags["addr:postcode"],
		Hash:     GetHash(tags["addr:street"] + tags["addr:city"] + tags["addr:postcode"]),
	}
}

func ParseFromUrl(url string, result chan<- Address) {
	if !strings.HasPrefix(url, "https://") {
		fmt.Printf("Rejected non-HTTPS URL: %s\n", url)
		return
	}

	client := &http.Client{Timeout: 30 * time.Minute}
	response, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
		return
	}
	defer response.Body.Close()

	d := osmpbf.NewDecoder(response.Body)
	err = d.Start(4) // 4 worker goroutines
	if err != nil {
		fmt.Printf("Error starting decoder: %v\n", err)
		return
	}

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("Error decoding: %v\n", err)
			break
		} else {
			switch obj := v.(type) {
			case *osmpbf.Node:
				if hasAddressTags(obj.Tags) {
					result <- addressFromOSMTags(obj.Tags)
				}
			case *osmpbf.Way:
				if hasAddressTags(obj.Tags) {
					result <- addressFromOSMTags(obj.Tags)
				}
			}
		}
	}
}
