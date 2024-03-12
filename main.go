package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
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

	newFile := hclwrite.NewFile()
	newFile.Body().AppendUnstructuredTokens(hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte("// DO NOT EDIT, MADE BY hclwrite-dynamodb-generator\n"),
		},
	})
	for _, resourceName := range referenceNames {
		newFile.Body().AppendUnstructuredTokens(hclwrite.Tokens{
			{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")}, // 先頭行に改行を入れる
		})
		labels := []string{"aws_cloudwatch_metric_alarm", fmt.Sprintf("dynamodb_throttledrequests_%s", resourceName)}
		resource := newFile.Body().AppendNewBlock("resource", labels).Body()
		resource.SetAttributeRaw("alarm_name", hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenIdent,
				Bytes: []byte(`"${aws_dynamodb_table.myproduct_read.name}-throttledrequests"`),
			},
		})
		resource.SetAttributeValue("comparison_operator", cty.StringVal("GreaterThanOrEqualToThreshold"))
		resource.SetAttributeValue("datapoints_to_alarm", cty.StringVal("1"))
		resource.SetAttributeValue("evaluation_periods", cty.StringVal("1"))
		resource.SetAttributeValue("metric_name", cty.StringVal("ThrottledRequests"))
		resource.SetAttributeValue("namespace", cty.StringVal("AWS/DynamoDB"))
		resource.SetAttributeValue("period", cty.StringVal("60"))
		resource.SetAttributeValue("statistic", cty.StringVal("Maximum"))
		resource.SetAttributeValue("threshold", cty.StringVal("1"))
		resource.SetAttributeTraversal("alarm_actions", hcl.Traversal{
			hcl.TraverseRoot{Name: "aws_sns_topic"},
			hcl.TraverseAttr{Name: "myproduct_alert"},
			hcl.TraverseAttr{Name: "arn"},
		})
		dimensions := resource.AppendNewBlock("dimensions", nil).Body()
		dimensions.SetAttributeTraversal("TableName", hcl.Traversal{
			hcl.TraverseRoot{Name: "aws_dynamodb_table"},
			hcl.TraverseAttr{Name: "myproduct_read"},
			hcl.TraverseAttr{Name: "name"},
		})
	}

	out := hclwrite.Format(newFile.BuildTokens(nil).Bytes())
	_, _ = fmt.Fprint(os.Stdout, string(out))
}
