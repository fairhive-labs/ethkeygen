# Eth Key Generator

Eth Key Generator is a tool for generating and signing Ethereum keys and messages. It is designed to facilitate secure key management and message signing for Ethereum-based applications.

## Features

- **Key Generation:** Generate Ethereum keys easily.
- **Message Signing:** Sign messages using your Ethereum private key.
- **API Integration:** Use the provided API to generate signatures in bulk.

## Usage

### Generates Keys

Open your browser and navigate to [ethkeygen.herokuapp.com](https://ethkeygen.herokuapp.com)

By default, this will generate 10 Ethereum keys. If you want to generate a different number of keys, specify the `l` query parameter. For example, to generate 20 keys, use the following URL [https://ethkeygen.herokuapp.com/?l=20](https://ethkeygen.herokuapp.com/?l=20)

Replace `20` with the desired number of keys, and the service will return the requested number of Ethereum keys ;)

### Sign the same message using multiple keys (API)

```bash
curl -s -X PUT "https://ethkeygen.herokuapp.com/bulk-signature?n=2" \
  -H 'content-type: application/json' \
  -d '{"message":"your message"}' | jq .
```

## Docker Usage

### Build Locally

To build the Docker image locally:

```bash
docker build -t ethkeygen .
```

### Pull from GitHub Package Registry

Alternatively, pull the pre-built image from the GitHub package registry:

```bash
docker pull ghcr.io/fairhive-labs/ethkeygen
docker tag ghcr.io/fairhive-labs/ethkeygen ethkeygen
```

### Run the Container

To run the container in iterative mode (with auto-destruct, and exposing port 8080 of the container on port 8081 on your machine), use the following command:

```bash
docker run -it --rm -p 8081:8080 ethkeygen
```

Open your browser and enter the following URL: [http://localhost:8081](http://localhost:8081)

### Stop the Container

To stop the running container:

```bash
docker stop <container_id>
```

### Remove the Container

To remove the container:

```bash
docker rm <container_id>
```

### Remove the Image

To delete the Docker image:

```bash
docker rmi ethkeygen
```
