package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"net/url"
)

func main() {

	// -h and --help flag
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
	fmt.Println(`Usage: luffuf -u <url> -w <wordlist> [other ffuf options]

	This wrapper:
	- Automatically injects stealth headers (unless you override them)
	- Adds Referer/Origin from the provided -u URL
	- Defaults to: -rate 5 -timeout 10 (unless overridden)
	- You can add any ffuf flag you want

	Examples:
	luffuf -u https://site.com/FUZZ -w paths.txt -mc 200
	`)
	os.Exit(0) // successful exit
}


	if len(os.Args) < 2 {
		fmt.Println("Usage: luffuf [ffuf arguments like -u, -w, etc.]")
		os.Exit(1) // error exit
	}

	userArgs := os.Args[1:]
	argsStr := strings.Join(userArgs, " ") // for easier contains checks

	// Extract base URL from -u argument
	baseURL := ""
	for i := 0; i < len(userArgs)-1; i++ {
		if userArgs[i] == "-u" && i+1 < len(userArgs) {
			parsedURL, err := url.Parse(userArgs[i+1])
			if err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
				baseURL = fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
			}
			break
		}
	}


	customArgs := []string{}

	// Conditionally inject headers
	if !strings.Contains(argsStr, "User-Agent:") {
		customArgs = append(customArgs, "-H", "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:108.0) Gecko/20100101 Firefox/108.0")
	}
	if !strings.Contains(argsStr, "Accept:") {
		customArgs = append(customArgs, "-H", "Accept: text/html")
	}
	if !strings.Contains(argsStr, "Accept-Language:") {
		customArgs = append(customArgs, "-H", "Accept-Language: en-US,en;q=0.5")
	}
	if !strings.Contains(argsStr, "Connection:") {
		customArgs = append(customArgs, "-H", "Connection: close")
	}

	// Only inject rate if user didn't specify
	if !strings.Contains(argsStr, "-rate") {
		customArgs = append(customArgs, "-rate", "5")
	}
	if !strings.Contains(argsStr, "-timeout") {
		customArgs = append(customArgs, "-timeout", "10")
	}

	if baseURL != "" {
		if !strings.Contains(argsStr, "Referer:") {
			customArgs = append(customArgs, "-H", fmt.Sprintf("Referer: %s", baseURL))
		}
		if !strings.Contains(argsStr, "Origin:") {
			customArgs = append(customArgs, "-H", fmt.Sprintf("Origin: %s", baseURL))
		}
	}

	// Combine injected args with user args
	fullArgs := append(customArgs, userArgs...)

	cmd := exec.Command("ffuf", fullArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running ffuf: %v\n", err)
	}
}
