package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "setup" {
		fmt.Println("1. Go to github.com → Settings → Password and authentication")
		fmt.Println("2. Click 'Enable two-factor authentication'")
		fmt.Println("3. Click 'setup key' under the QR code to get the secret")
		fmt.Println("4. Run: totp <secret>")
		fmt.Println("5. Enter the 6-digit code on GitHub to verify")
		fmt.Println("6. Save the secret somewhere safe (e.g. ~/.totp)")
		fmt.Println()
		fmt.Println("To use permanently, save secret to ~/.totp and run without args:")
		fmt.Println("  echo 'YOURSECRETHERE' > ~/.totp && chmod 600 ~/.totp")
		return
	}

	secret := ""
	if len(os.Args) >= 2 {
		secret = os.Args[1]
	}
	if secret == "" {
		secret = readSecretFile()
	}
	if secret == "" {
		fmt.Fprintf(os.Stderr, "Usage: totp <secret>\n")
		fmt.Fprintf(os.Stderr, "       totp setup    — show instructions\n")
		fmt.Fprintf(os.Stderr, "Or save secret to ~/.totp and run: totp\n")
		os.Exit(1)
	}

	code, remaining := generateTOTP(secret)
	fmt.Printf("%06d  (%ds remaining)\n", code, remaining)
}

func readSecretFile() string {
	home, _ := os.UserHomeDir()
	data, err := os.ReadFile(home + "/.totp")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func generateTOTP(secret string) (uint32, int) {
	// Decode base32 secret (GitHub gives base32-encoded keys)
	secret = strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(secret), " ", ""))
	// Pad to multiple of 8 for base32
	for len(secret)%8 != 0 {
		secret += "="
	}
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid secret (must be base32): %v\n", err)
		os.Exit(1)
	}

	// Time step: 30 seconds
	now := time.Now().Unix()
	counter := uint64(math.Floor(float64(now) / 30))
	remaining := 30 - int(now%30)

	// HMAC-SHA1(key, counter)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)
	mac := hmac.New(sha1.New, key)
	mac.Write(buf)
	hash := mac.Sum(nil)

	// Dynamic truncation (RFC 4226)
	offset := hash[len(hash)-1] & 0x0F
	code := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7FFFFFFF
	code = code % 1000000

	return code, remaining
}
