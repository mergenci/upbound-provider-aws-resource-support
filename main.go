// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/support"
	"github.com/aws/aws-sdk-go-v2/service/support/types"
)

func main() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Errorf("cannot load deafult configuration. make sure that AWS_PROFILE is set: %v", err))
	}

	supportClient := support.NewFromConfig(sdkConfig)
	language := "en"

	// CreateCase API call requires ServiceCode, which could be retrieved with DescribeServices call
	describeServicesOutput, err := supportClient.DescribeServices(context.TODO(), &support.DescribeServicesInput{
		Language: &language,
	})
	if err != nil {
		panic(fmt.Errorf("cannot describe services: %v", err))
	}

	sort.Slice(describeServicesOutput.Services, func(i, j int) bool {
		return strings.Compare(*describeServicesOutput.Services[i].Name, *describeServicesOutput.Services[j].Name) == -1
	})

	var theService *types.Service
	for i, service := range describeServicesOutput.Services {
		fmt.Printf("%-65q: %s\n", *service.Name, *service.Code)
		if *service.Code == "support-api" {
			theService = &describeServicesOutput.Services[i]
		}
	}

	fmt.Println()
	fmt.Println("The service for which SupportCase will be created (name: code):")
	fmt.Printf("%q: %s\n", *theService.Name, *theService.Code)
	fmt.Println()
	fmt.Println("Categories of the service (name: code):")
	for _, category := range theService.Categories {
		fmt.Printf("%-20q: %s\n", *category.Name, *category.Code)
	}

	createCaseInput := support.CreateCaseInput{
		CommunicationBody: aws.String("This support case is created for AWS SDK development purposes."),
		// IMPORTANT: Don't change subject in non-production. See https://docs.aws.amazon.com/awssupport/latest/user/about-support-api.html#endpoint
		Subject: aws.String("TEST CASE-Please ignore"),
		// AttachmentSetId:  new(string), // To be implemented later.
		CategoryCode:     theService.Categories[3].Code,
		CcEmailAddresses: []string{},
		IssueType:        aws.String("technical"), // "customer-service" or "technical" (default)
		Language:         &language,
		ServiceCode:      theService.Code,
		SeverityCode:     aws.String("low"), // https://docs.aws.amazon.com/awssupport/latest/APIReference/API_SeverityLevel.html
	}
	createCaseOutput, err := supportClient.CreateCase(context.TODO(), &createCaseInput)
	if err != nil {
		panic(fmt.Errorf("cannot create case: %v", err))
	}

	// displayId := "170550314501996"
	describeCasesInput := support.DescribeCasesInput{
		// AfterTime:             new(string),
		// BeforeTime:            new(string),
		CaseIdList: []string{*createCaseOutput.CaseId},
		// DisplayId: aws.String(displayId),
		// IncludeCommunications: new(bool),
		// IncludeResolvedCases:  false,
		// Language:              new(string),
		// MaxResults:            new(int32),
		// NextToken:             new(string),
	}
	describeCasesOutput, err := supportClient.DescribeCases(context.TODO(), &describeCasesInput)
	if err != nil {
		panic(fmt.Errorf("cannot describe case: %v", err))
	}
	fmt.Printf("%p", describeCasesOutput)

	resolveCaseOutput, err := supportClient.ResolveCase(context.TODO(), &support.ResolveCaseInput{
		CaseId: createCaseOutput.CaseId,
	})
	if err != nil {
		panic(fmt.Errorf("cannot resolve case: %v", err))
	}
	fmt.Printf("%p", resolveCaseOutput)
}
