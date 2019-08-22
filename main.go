package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/Shopify/ejson/crypto"
	"github.com/atotto/clipboard"
)

const clipboardPreviewLength = 50

var stdin *bufio.Reader
var publicKeyBytes []byte

func init() {
	stdin = bufio.NewReader(os.Stdin)
}

func main() {
	if len(os.Args) != 2 {
		usageAndDie()
	}
	switch os.Args[1] {
	case "send":
		sendSecret()
	case "receive":
		receiveSecret()
	default:
		usageAndDie()
	}
}

func usageAndDie() {
	fmt.Fprintln(os.Stderr, "usage: secret-sender send|receive")
	os.Exit(1)
}

func sendSecret() {
	fmt.Println("Ask the receiving party to run `secret-sender receive` and send you the public key that it generates.")
	fmt.Println("Paste the public key here:")
	publicKeyBytes = readline()
	hexPublicKeyBytes, err := hex.DecodeString(string(publicKeyBytes))
	if err != nil {
		log.Fatal(err)
	}

	var receiverKP crypto.Keypair
	for i, b := range hexPublicKeyBytes {
		receiverKP.Public[i] = b
	}

	var ephemeralKP crypto.Keypair
	if err := ephemeralKP.Generate(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Copy your secret into your clipboard:")
	plaintextResult := make(chan []byte)
	go getClipboardContents(plaintextResult, "encrypt it")
	plaintext := <-plaintextResult

	encrypter := crypto.NewEncrypter(&ephemeralKP, receiverKP.Public)
	ciphertext, err := encrypter.Encrypt(plaintext)
	if err != nil {
		log.Fatal(err)
	}

	copyToClipboard(string(ciphertext))
	fmt.Println("This is the encrypted string. Paste it to the receiver. (We've already put it in your clipboard):")
	fmt.Println(string(ciphertext))
}

func receiveSecret() {
	var kp crypto.Keypair
	if err := kp.Generate(); err != nil {
		log.Fatal(err)
	}
	publicKey := kp.PublicString()
	publicKeyBytes = []byte(publicKey)
	copyToClipboard(publicKey)

	fmt.Println("Paste this key to the sender (we've already put it in your clipboard):")
	fmt.Println(publicKey)

	fmt.Println("They'll respond with a big encrypted-looking blob. Once you receive it, copy it to your clipboard:")
	ciphertextResult := make(chan []byte)
	go getClipboardContents(ciphertextResult, "decrypt it")
	ciphertext := <-ciphertextResult

	decrypter := &crypto.Decrypter{Keypair: &kp}
	plaintext, err := decrypter.Decrypt(ciphertext)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("This is the secret:")
	fmt.Println(string(plaintext))
}

func readline() []byte {
	line, _, err := stdin.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	return line
}

func getClipboardContents(textResult chan []byte, action string) {
	done := make(chan bool)

	go getClipboardContentsInner(done, textResult, action)
	go waitUntilEnter(done)
}

func getClipboardContentsInner(done chan bool, textResult chan []byte, action string) {
	var clipboardContents []byte

	for {
		select {
		case <-done:
			textResult <- clipboardContents
			return
		default:
			newClipboardContents := copyFromClipboard()

			if bytes.Compare(newClipboardContents, clipboardContents) != 0 && bytes.Compare(newClipboardContents, publicKeyBytes) != 0 {
				newClipboardPreview := strings.Join(strings.Fields(strings.TrimSpace(string(newClipboardContents))), " ")
				newClipboardPreviewLength := math.Min(float64(len(newClipboardPreview)), clipboardPreviewLength)

				fmt.Printf(
					"Got new clipboard contents: length=%d preview='%s'. Press enter/return to %s.\n",
					len(newClipboardContents),
					newClipboardPreview[0:int(newClipboardPreviewLength)],
					action,
				)

				clipboardContents = newClipboardContents
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func waitUntilEnter(done chan bool) {
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	done <- true
}

func copyToClipboard(text string) {
	if err := clipboard.WriteAll(text); err != nil {
		log.Fatal(err)
	}
}

func copyFromClipboard() []byte {
	text, err := clipboard.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return []byte(text)
}
