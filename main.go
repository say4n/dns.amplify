package main

func main() {
	query := "example.com"
	queryMessage := GenerateDNSMessage(query)

	err := PerformDNSRequest("dns.google:53", "localhost:2000", queryMessage)
	if err != nil {
		panic(err)
	}
}
