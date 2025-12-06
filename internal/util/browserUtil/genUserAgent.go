package browserUtil

import (
	"fmt"
	"math/rand"
)

// GenerateRandomUserAgent generates a random Chrome user agent string
func GenerateRandomUserAgent() string {
	osType := rand.Intn(3) // 0: Windows, 1: macOS, 2: Linux
	var osString string

	switch osType {
	case 0: // Windows
		windowsVersions := []string{"10.0", "11.0"}
		winVer := windowsVersions[rand.Intn(len(windowsVersions))]
		osString = fmt.Sprintf("Windows NT %s; Win64; x64", winVer)
	case 1: // macOS
		macVersions := []string{"10_15_7", "11_6", "12_6", "13_5", "14_0"}
		macVer := macVersions[rand.Intn(len(macVersions))]
		osString = fmt.Sprintf("Macintosh; Intel Mac OS X %s", macVer)
	case 2: // Linux
		osString = "X11; Linux x86_64"
	}

	// Generate random Chrome version
	// Major version between 115 and 125
	major := rand.Intn(11) + 115
	// Build number between 0 and 5000
	build := rand.Intn(5001)
	// Patch number between 0 and 200
	patch := rand.Intn(201)

	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.%d.%d Safari/537.36",
		osString, major, build, patch)
}
