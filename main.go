package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	_ "golang.org/x/image/bmp"  // Supports BMP
	_ "golang.org/x/image/webp" // Supports WebP

	// Importing the general purpose Cosmos blockchain client
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"

	// Importing the types package of your blog blockchain
	"cifer/x/cifer/types"
)

type Metadata struct {
	FileName string `json:"file_name"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Format   string `json:"format"`
	Checksum string `json:"checksum"`
	FileSize int64  `json:"file_size"`
}

// GetImageMetadata retrieves image metadata
func GetImageMetadata(filePath string) (*Metadata, error) {
	// Open image file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Decode image to retrieve size and format
	img, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %v", err)
	}

	// Retrieve file size
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("error getting file info: %v", err)
	}

	// Calculate checksum
	checksum, err := GetFileChecksum(filePath)
	if err != nil {
		return nil, fmt.Errorf("error generating checksum: %v", err)
	}

	// Combine metadata
	metadata := &Metadata{
		FileName: fileInfo.Name(),
		Width:    img.Width,
		Height:   img.Height,
		Format:   format,
		Checksum: checksum,
		FileSize: fileInfo.Size(),
	}

	return metadata, nil
}

// GetFileChecksum generates SHA256 checksum of file
func GetFileChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("error resetting file pointer: %v", err)
	}

	if _, err := file.Read(hash.Sum(nil)); err != nil {
		return "", fmt.Errorf("error calculating checksum: %v", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func main() {
	ctx := context.Background()
	addressPrefix := "cosmos"

	// Create a Cosmos client instance
	client, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix))
	if err != nil {
		log.Fatal(err)
	}

	// Account `alice` was initialized during `ignite chain serve`
	accountName := "alice"

	// Get account from the keyring
	account, err := client.Account(accountName)
	if err != nil {
		log.Fatal(err)
	}

	addr, err := account.Address(addressPrefix)
	if err != nil {
		log.Fatal(err)
	}

	// File Metadata
	filePath := "dataset_example.png"

	// GET Metadata
	metadata, err := GetImageMetadata(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print Metadata
	fmt.Println("Metadata:")
	fmt.Printf("File Name: %s\n", metadata.FileName)
	fmt.Printf("Width: %d px\n", metadata.Width)
	fmt.Printf("Height: %d px\n", metadata.Height)
	fmt.Printf("Format: %s\n", metadata.Format)
	fmt.Printf("Checksum (SHA256): %s\n", metadata.Checksum)
	fmt.Printf("File Size: %d bytes\n", metadata.FileSize)

	// Metadata for Blockchain
	jsonData, _ := json.Marshal(metadata)
	fmt.Println(string(jsonData))

	msg := &types.MsgCreateMintdata{
		Creator:        addr,
		CreatorAddress: addr,
		Metadata:       string(jsonData),
	}

	// save to JSON
	fileName := "metadata.json"
	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing JSON file: %v\n", err)
		return
	}

	fmt.Printf("JSON file %s created successfully\n", fileName)

	// Broadcast a transaction from account `alice` with the message
	// to create a post store response in txResp
	txResp, err := client.BroadcastTx(ctx, account, msg)
	if err != nil {
		log.Fatal(err)
	}

	// Print response from broadcasting a transaction
	fmt.Print("MsgCreatePost:\n\n")
	fmt.Println(txResp)

	// Instantiate a query client for your `blog` blockchain
	queryClient := types.NewQueryClient(client.Context())

	// Query the blockchain using the client's `PostAll` method
	// to get all posts store all posts in queryResp
	queryResp, err := queryClient.MintdataAll(ctx, &types.QueryAllMintdataRequest{})
	if err != nil {
		log.Fatal(err)
	}

	// Print response from querying all the posts
	fmt.Print("\n\nAll posts:\n\n")
	fmt.Println(queryResp)
}
