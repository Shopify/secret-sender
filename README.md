# secret-sender

Send passwords and tokens manually over insecure channels

## Installation

```
$ brew tap shopify/shopify
$ brew install secret-sender
```

## How it works

`secret-sender` requires two users to run the program at the same time in cooperation, and paste messages to each other.
These messages are not secret so they can be sent over (e.g.) Slack.

1. The user receving the secret token or password will run `secret-sender receive`. This will copy their public key to
   the clipboard and display it in their shell.
2. The receiving user will message that public key to the sender.
3. The sender will run `secret-sender send`, paste the public key into their shell, and press Enter/Return.
4. The sender will then be prompted to copy their secret to their clipboard and press Enter/Return.
5. secret-sender will copy the encrypted string to the clipboard and display it in the shell.
6. The sender will message that encrypted string to the receiver.
7. The receiver copies that encrypted string into their clip and presses Enter/Return in secret-sender.
8. secret-sender displays the secret in the shell.

## Under the hood

Secret-sender uses NaCl Box cryptograpy, or curve25519xsalsa20poly1305.
The receiver generates an ephemeral keypair and sends the public portion to the sender, who encrypts the secret to that
key, before sending the ciphertext to the receiver. The receiver then recovers the plaintext and terminates, discarding
the private key.

Neither subcommand takes any arguments, but both ask the sender and receiver to copy their plaintext and ciphertext
(respectively) into their clipboard once the program is running.

Scripting with secret-sender is discouraged: instead, use [ejson](http://github.com/Shopify/ejson) directly.
