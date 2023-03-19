package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Fatal("Error reading input")
	}

	text := scanner.Text()
	handleInput(text)

}

func handleInput(text string) {
	emails := strings.Split(text, ",")

	c := make(chan string)

	for _, email := range emails {
		go checkEmail(strings.Trim(email, " "), c)
	}

	fmt.Print("\nProcessing...\n\n")

	for i := 0; i < len(emails); i++ {
		fmt.Println(<-c)
	}

	fmt.Println()
}

func checkEmail(email string, c chan string) {

	domain := email

	if parts := strings.Split(email, "@"); len(parts) > 1 {
		domain = parts[1]
	}

	mxRecords, err := net.LookupMX(domain)

	if err != nil || len(mxRecords) <= 0 {
		c <- fmt.Sprintf("%s: %s", email, "Is not valid ❌")
		return
	}

	txtRecords, err := net.LookupTXT(domain)

	if err != nil {
		c <- fmt.Sprintf("%s: %s", email, "Is not valid ❌ (Error fetching SPF record)")
		return
	}

	for i, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			break
		} else if i == len(txtRecords)-1 {
			c <- fmt.Sprintf("%s: %s", email, "Is not valid ❌ (No SPF record)")
			return
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		c <- fmt.Sprintf("%s: %s", email, "Is not valid ❌ (Error fetching DMARC record)")
		return
	}

	for i, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			break
		} else if i == len(dmarcRecords)-1 {
			c <- fmt.Sprintf("%s: %s", email, "Is not valid ❌ (No DMARC record)")
			return
		}
	}

	c <- fmt.Sprintf("%s: %s", email, "Is valid ✅")
}
