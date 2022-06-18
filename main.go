package main

func main() {
	query := "fb.me"
	queryMessage := GenerateDNSMessage(query)

	err := PerformDNSRequest("1.0.0.1:53", "", queryMessage)
	if err != nil {
		panic(err)
	}
}
