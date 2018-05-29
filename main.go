package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/Shopify/ejson/crypto"
	"github.com/atotto/clipboard"
)

var stdin *bufio.Reader

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

func sendSecret() {
	fmt.Println("Ask the receiving party to run `secret-sender receive` and send you the public key that it generates.")
	fmt.Println("Paste the public key here:")
	pk := readline()
	bytes, err := hex.DecodeString(string(pk))
	if err != nil {
		log.Fatal(err)
	}
	var receiverKP crypto.Keypair
	for i, b := range bytes {
		receiverKP.Public[i] = b
	}

	var ephemeralKP crypto.Keypair
	if err := ephemeralKP.Generate(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Paste your secret here:")
	plaintext := readline()

	encrypter := crypto.NewEncrypter(&ephemeralKP, receiverKP.Public)
	ciphertext, err := encrypter.Encrypt(plaintext)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("This is the encrypted string. Paste it to the receiver. (We've already put it in your clipboard):")
	pbcopy(string(ciphertext))
	fmt.Println(string(ciphertext))
}

func receiveSecret() {
	var kp crypto.Keypair
	if err := kp.Generate(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Paste this key to the sender (we've already put it in your clipboard):")
	pbcopy(kp.PublicString())
	fmt.Println(kp.PublicString())
	fmt.Println("They'll respond with a big encrypted-looking blob. Paste it here, then press return:")

	ciphertext := readline()

	decrypter := &crypto.Decrypter{Keypair: &kp}
	plaintext, err := decrypter.Decrypt(ciphertext)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("This is the secret:")
	fmt.Println(string(plaintext))
}

func usageAndDie() {
	fmt.Fprintln(os.Stderr, "usage: secret-sender send|receive")
	os.Exit(1)
}

func readline() []byte {
	line, _, err := stdin.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	return line
}

func pbcopy(text string) {
	if err := clipboard.WriteAll(text); err != nil {
		log.Fatal(err)
	}
}
