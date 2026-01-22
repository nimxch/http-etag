// package main

// import (
// 	"bytes"
// 	"crypto/rand"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"math/big"
// 	"net/http"
// 	"os"
// 	"time"
// )

// const (
// 	kb = 1024

// 	upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
// 	lower = "abcdefghijklmnopqrstuvwxyz"
// 	digit = "0123456789"

// 	// Add special chars here if you want: "!@#$%^&*()-_=+[]{}..."
// 	all = upper + lower + digit
// )

// func randInt(n int) (int, error) {
// 	x, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
// 	if err != nil {
// 		return 0, err
// 	}
// 	return int(x.Int64()), nil
// }

// func randChar(set string) (byte, error) {
// 	i, err := randInt(len(set))
// 	if err != nil {
// 		return 0, err
// 	}
// 	return set[i], nil
// }

// func shuffleBytes(b []byte) error {
// 	// Fisher–Yates shuffle with crypto/rand
// 	for i := len(b) - 1; i > 0; i-- {
// 		j, err := randInt(i + 1)
// 		if err != nil {
// 			return err
// 		}
// 		b[i], b[j] = b[j], b[i]
// 	}
// 	return nil
// }

// func GeneratePasswordKB(sizeKB int) (string, error) {
// 	if sizeKB <= 0 {
// 		return "", fmt.Errorf("sizeKB must be > 0")
// 	}

// 	total := sizeKB * kb
// 	if total < 8 {
// 		total = 8
// 	}

// 	pw := make([]byte, 0, total)

// 	fmt.Printf("Generating password: %d KB (%d bytes)\n", sizeKB, total)

// 	// Enforce rule by construction
// 	c, err := randChar(upper)
// 	if err != nil {
// 		return "", err
// 	}
// 	pw = append(pw, c)

// 	c, err = randChar(lower)
// 	if err != nil {
// 		return "", err
// 	}
// 	pw = append(pw, c)

// 	c, err = randChar(digit)
// 	if err != nil {
// 		return "", err
// 	}
// 	pw = append(pw, c)

// 	// Fill remaining
// 	// Progress update granularity: every 64 KB
// 	progressStep := 64 * kb
// 	nextPrint := progressStep

// 	for len(pw) < total {
// 		c, err := randChar(all)
// 		if err != nil {
// 			return "", err
// 		}
// 		pw = append(pw, c)

// 		if len(pw) >= nextPrint || len(pw) == total {
// 			percent := (len(pw) * 100) / total
// 			fmt.Printf("\rProgress: %3d%% (%d / %d bytes)", percent, len(pw), total)
// 			nextPrint += progressStep
// 		}
// 	}

// 	// Shuffle so first 3 chars aren't predictable type ordering
// 	if err := shuffleBytes(pw); err != nil {
// 		return "", err
// 	}

// 	fmt.Println("\nDone.")
// 	return string(pw), nil
// }

// type ResetPasswordRequest struct {
// 	UserName        string `json:"userName"`
// 	NewPassword     string `json:"newPassword"`
// 	ResetPasswordId string `json:"resetPasswordId"`
// }

// func Reset() {
// 	passwordSizeKB := 1024 // 1MB password string

// 	pw, err := GeneratePasswordKB(passwordSizeKB)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Save
// 	err = os.WriteFile("password.txt", []byte(pw), 0600)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("Saved password.txt (size=%d bytes)\n", len(pw))

// 	userName := "<email@email.emai>"
// 	resetPasswordId := "TEST_TEST_TEST"

// 	apiURL := "https://<Your Reset Password URL>"

// 	// Paste the cookie header string (from browser/curl) ONLY if required.
// 	// Otherwise leave empty.
// 	cookieHeader := "" // e.g. "G_ENABLED_IDPS=google; iqnextcares-_auuid=...; ..."

// 	// === Call API ===
// 	req := ResetPasswordRequest{
// 		UserName:        userName,
// 		NewPassword:     pw,
// 		ResetPasswordId: resetPasswordId,
// 	}

// 	status, respBody, err := ResetPassword(apiURL, req, cookieHeader)
// 	if err != nil {
// 		log.Fatalf("reset-password call failed: %v", err)
// 	}

// 	fmt.Println("HTTP Status:", status)
// 	fmt.Println(respBody)
// }

// func ResetPassword(apiURL string, req ResetPasswordRequest, cookieHeader string) (int, []byte, error) {
// 	bodyBytes, err := json.Marshal(req)
// 	if err != nil {
// 		return 0, nil, err
// 	}

// 	httpReq, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(bodyBytes))
// 	if err != nil {
// 		return 0, nil, err
// 	}

// 	// minimum necessary headers
// 	httpReq.Header.Set("Accept", "application/json")
// 	httpReq.Header.Set("Content-Type", "application/json")
// 	httpReq.Header.Set("Origin", "https://faraday.iqnext.io")
// 	httpReq.Header.Set("Referer", "https://faraday.iqnext.io/")

// 	// If endpoint requires cookie auth/session, pass it
// 	// NOTE: This is sensitive. Don't hardcode in real code.
// 	if cookieHeader != "" {
// 		httpReq.Header.Set("Cookie", cookieHeader)
// 	}

// 	client := &http.Client{
// 		Timeout: 30 * time.Second,
// 	}

// 	resp, err := client.Do(httpReq)
// 	if err != nil {
// 		return 0, nil, err
// 	}
// 	defer resp.Body.Close()

// 	respBody, _ := io.ReadAll(resp.Body)
// 	return resp.StatusCode, respBody, nil
// }
