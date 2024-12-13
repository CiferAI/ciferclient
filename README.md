# Data Minting

This repository provides code examples and instructions for minting data on the Cifer blockchain using metadata extracted from files. The process is simplified for users to focus on the essential steps to mint their data.

## **How to Mint Data**

### Steps:
1. **Connect your wallet** with your Cifer address.
2. **Define the data path** (e.g., `dataset_example.png`).
3. **Run the command:**
   ```bash
   go run main.go
   ```
<br>

---
<br>
The code in this repository automates the preparation, metadata generation, and submission to the blockchain, making it easy for users to mint data. While you only need to follow the simple steps outlined above, the code handles the detailed workflow, as described below:

### Summary of Automated Data Minting Workflow
1. **Collect metadata from the file** including size, format, and checksum.
2. **Save metadata as JSON** for reference.
3. **Send metadata to the blockchain** via a smart contract or supported API.
4. **Verify the transaction** and ensure the data is successfully minted on the blockchain.

---

### **1. Prepare Metadata**

Gather metadata from the file to be minted. The metadata typically includes:

- **File Name:** Name of the file
- **File Size:** File size in bytes
- **Image Dimensions (Width, Height):** Width and height of the image (if applicable)
- **File Format:** e.g., PNG, JPEG
- **Checksum:** SHA256 hash for file integrity verification

Example code for extracting metadata:
```go
metadata := Metadata{
    FileName: "dataset_example.png",
    Width:    800,
    Height:   600,
    Format:   "png",
    Checksum: "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd98a3457f0a65d2c10",
    FileSize: 15023,
}
```

---

### **2. Save Metadata as JSON**

The metadata is converted into a JSON format and saved as a file (e.g., metadata.json) for submission to the blockchain.

Example code:
```go
jsonData, _ := json.MarshalIndent(metadata, "", "  ")
os.WriteFile("metadata.json", jsonData, 0644)
```

Sample `metadata.json` output:
```json
{
  "file_name": "example.png",
  "width": 800,
  "height": 600,
  "format": "png",
  "checksum": "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd98a3457f0a65d2c10",
  "file_size": 15023
}
```

---

### **3. Prepare Code to Submit Metadata to Blockchain**

The metadata will be submitted to the blockchain using a smart contract that supports **mint data** commands, such as `MsgCreateMintdata` in the Cosmos SDK or a mint function in Ethereum.

Example code for Cosmos SDK:
```go
msg := &types.MsgCreateMintdata{
    Creator:        "cosmos1...", // Creator address
    CreatorAddress: "cosmos1...",
    Metadata:       string(jsonData), // Convert JSON to string
}
```

---

### **4. Broadcast the Transaction**

Use a blockchain client (e.g., **Cosmos SDK Client** or Web3 for Ethereum) to broadcast the transaction. Ensure the sender's address and metadata are accurate.

Example code for broadcasting a transaction in Cosmos SDK:
```go
txResp, err := client.BroadcastTx(ctx, account, msg)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Transaction Response:", txResp)
```

---

### **5. Verify Minted Data on Blockchain**

After broadcasting the transaction, verify the minted data on the blockchain to ensure successful submission.

Example code for querying minted data:
```go
queryResp, err := queryClient.MintdataAll(ctx, &types.QueryAllMintdataRequest{})
if err != nil {
    log.Fatal(err)
}

fmt.Println("Minted Data:", queryResp)
```

---

### **6. Check Transaction Status**

Use blockchain commands or APIs to check the transaction status, such as:
- Verify the transaction hash
- Check the block containing the transaction
