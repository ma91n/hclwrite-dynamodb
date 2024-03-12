package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatalf("Usage: %s <filepath>\n", os.Args[0])
	}

	hclFilePath := os.Args[1]
	file, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Usage: %s <filepath>\n", os.Args[1])
	}

	tfFile, diags := hclwrite.ParseConfig(file, hclFilePath, hcl.Pos{Line: 1, Column: 1})
	if diags != nil && diags.HasErrors() {
		log.Fatalf("hclwrite parse: %s", diags)
	}
	if tfFile == nil {
		log.Fatalf("parse result is nil: %s", hclFilePath)
	}

	blocks := tfFile.Body().Blocks()
	referenceNames := make([]string, 0, len(blocks))
	for _, b := range blocks {
		if b.Type() != "resource" || b.Labels()[0] != "aws_dynamodb_table" {
			continue
		}
		referenceNames = append(referenceNames, b.Labels()[1])
	}

	fmt.Println(referenceNames)
}
