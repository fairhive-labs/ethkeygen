# Eth Key Generator

Eth Key Generator is a tool for generating and signing Ethereum keys and messages. It is designed to facilitate secure key management and message signing for Ethereum-based applications.

## Features

- **Key Generation:** Generate Ethereum keys easily.
- **Message Signing:** Sign messages using your Ethereum private key.
- **API Integration:** Use the provided API to generate signatures in bulk.

## Usage

```bash
curl -s -X PUT "https://ethkeygen.herokuapp.com/bulk-signature?n=2" \
  -H 'content-type: application/json' \
  -d '{"message":"your message"}' | jq .

