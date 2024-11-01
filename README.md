# rlpencoder

Rlpencoder is a tool to encode and decode RLP (Recursive Length Prefix) data. It allows you to save and restore RLP data (e.g. Ethereum blocks data) from a file.

## Features

- Decode RLP data to human-readable JSON format
- Encode JSON data back to RLP format
- Stream processing for handling large files

## Installation

```bash
go install github.com/dmitriy-b/rlp-to-blockchain/cmd/rlpencoder@latest
```

## Usage

### Decoding RLP to JSON

To decode an RLP file to JSON:

```bash
# Decode chain.rlp to output.json (default output file)
rlpencoder

# Decode chain.rlp to a custom output file
rlpencoder -output custom_output.json
```

### Encoding JSON to RLP

To encode JSON back to RLP format:

```bash
# Load JSON file and convert it to RLP format (saves to output.rlp)
rlpencoder -load input.json
```

## File Format

### Input RLP File

- The tool expects an RLP-encoded file named `chain.rlp` in the current directory when decoding
- The RLP data can contain multiple blocks or entries

### Output JSON File

- The JSON output preserves the structure of the RLP data
- All byte arrays are converted to hexadecimal strings for readability
- The output is pretty-printed with proper indentation

### Input JSON File

- When encoding back to RLP, the JSON structure should match the original format
- Hexadecimal strings are automatically converted back to byte arrays

## Examples

Example of decoded JSON output:

```json
[
  [
    "f9021ca0...",
    "a0...",
    [
      "0x123...",
      "0x456..."
    ]
  ]
]
```

## Building from Source

```bash
git clone https://github.com/yourusername/rlpencoder
cd rlpencoder
go build ./cmd/rlpencoder
```

## Requirements

- Go 1.22 or later
- The `github.com/ethereum/go-ethereum` package for RLP encoding/decoding
