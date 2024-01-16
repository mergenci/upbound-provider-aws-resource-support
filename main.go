// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/support"
)

func main() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Errorf("cannot load deafult configuration. make sure that AWS_PROFILE is set: %v", err))
	}

	supportClient := support.NewFromConfig(sdkConfig)
	language := "en"

	serviceCodes := make([]string, 100)
	describeServicesoutput, err := supportClient.DescribeServices(context.TODO(), &support.DescribeServicesInput{
		Language:        &language,
		ServiceCodeList: serviceCodes,
	})
	if err != nil {
		panic(fmt.Errorf("cannot describe services: %v", err))
	}
	fmt.Printf("%v", describeServicesoutput)

	createCaseInput := support.CreateCaseInput{
		CommunicationBody: aws.String("This support case is created for AWS SDK development purposes."),
		// IMPORTANT: Don't change subject in non-production. See https://docs.aws.amazon.com/awssupport/latest/user/about-support-api.html#endpoint
		Subject: aws.String("TEST CASE-Please ignore"),
		// AttachmentSetId:  new(string), // To be implemented later.
		CategoryCode:     new(string),
		CcEmailAddresses: []string{},
		IssueType:        aws.String("technical"), // "customer-service" or "technical" (default)
		Language:         &language,
		ServiceCode:      aws.String("10"),  // TODO: This value should be replaced with a proper one returned by DescribeServices call above.
		SeverityCode:     aws.String("low"), // https://docs.aws.amazon.com/awssupport/latest/APIReference/API_SeverityLevel.html
	}
	createCaseOutput, err := supportClient.CreateCase(context.TODO(), &createCaseInput)
	if err != nil {
		panic(fmt.Errorf("cannot create case: %v", err))
	}
	fmt.Printf("%v", createCaseOutput)
}
