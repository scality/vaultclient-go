package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/scality/vaultclient-go/vaultclient"
)

func main() {
	ctx := context.Background()

	vaultAccessKey := "D4IT2AWSB588GO5J9T00"
	vaultSecretKey := "UEEu8tYlsOGGrgf4DAiSZD6apVNPUWqRiPG0nTB6"
	vaultEndpoint := "http://localhost:8600"

	client := vaultclient.New(vaultAccessKey, vaultSecretKey, "", vaultEndpoint, "us-east-1")

	fmt.Println("Vault Client tests")

	// Example 1.1: Create an account
	fmt.Println("1. Creating a new account...")
	accountName := fmt.Sprintf("test-account-%d", time.Now().Unix())
	createInput := &vaultclient.CreateAccountInput{}
	createInput.SetName(accountName).
		SetEmail(fmt.Sprintf("%d@example.com", time.Now().Unix())).
		SetQuotaMax(100)

	createResult, err := client.CreateAccount(ctx, createInput)
	if err == nil {
		fmt.Printf("Account created successfully!\n")
		if createResult.Account != nil && createResult.Account.AccountData != nil {
			fmt.Printf("Account ID: %s\n", *createResult.Account.AccountData.ID)
			fmt.Printf("Account Name: %s\n", *createResult.Account.AccountData.Name)
			fmt.Printf("Account Email: %s\n", *createResult.Account.AccountData.Email)
		}
	}
	fmt.Println()

	// Example 1.2: Test error handling on duplicate account creation
	_, err = client.CreateAccount(ctx, createInput)
	if err != nil {
		var apiErr *vaultclient.APIError
		var entityExists *types.EntityAlreadyExistsException
		if errors.As(err, &apiErr) {
			fmt.Printf("StatusCode: %d\n", apiErr.HTTPStatusCode())
		}
		if errors.As(err, &entityExists) {
			fmt.Printf("entity exists: %s: %s\n", entityExists.ErrorCode(), *entityExists.Message)
		} else {
			fmt.Printf("unexpected error: %v\n", err)
		}
	}

	// Example 2: List accounts
	fmt.Println("2. Listing accounts...")
	listInput := &vaultclient.ListAccountsInput{}
	listInput.SetMaxItems(10)

	listResult, err := client.ListAccounts(ctx, listInput)
	if err != nil {
		log.Printf("Error listing accounts: %v", err)
	} else {
		fmt.Printf("Found %d accounts\n", len(listResult.Accounts))
		for i, account := range listResult.Accounts {
			if account.Name != nil && account.Email != nil {
				fmt.Printf("   %d. %s (%s)\n", i+1, *account.Name, *account.Email)
			}
		}
	}
	fmt.Println()

	// Example 3: Get account by name
	fmt.Println("3. Getting account by name...")
	getInput := &vaultclient.GetAccountInput{}
	getInput.SetName(accountName)

	getResult, err := client.GetAccount(ctx, getInput)
	if err != nil {
		log.Printf("Error getting account: %v", err)
	} else {
		fmt.Printf("Account found!\n")
		if getResult.Name != nil && getResult.Email != nil {
			fmt.Printf("Name: %s\n", *getResult.Name)
			fmt.Printf("Email: %s\n", *getResult.Email)
			if getResult.ID != nil {
				fmt.Printf("ID: %s\n", *getResult.ID)
			}
		}
	}
	fmt.Println()

	// Example 4: Generate access key for account
	fmt.Println("4. Generating access key...")
	genKeyInput := &vaultclient.GenerateAccountAccessKeyInput{}
	genKeyInput.SetAccountName(accountName)

	genKeyResult, err := client.GenerateAccountAccessKey(ctx, genKeyInput)
	if err != nil {
		log.Printf("Error generating access key: %v", err)
	} else {
		fmt.Printf("Access key generated!\n")
		if genKeyResult.GeneratedKey != nil {
			if genKeyResult.GeneratedKey.ID != nil {
				fmt.Printf("Access Key ID: %s\n", *genKeyResult.GeneratedKey.ID)
			}
			if genKeyResult.GeneratedKey.Value != nil {
				fmt.Printf("Secret Key: %s\n", *genKeyResult.GeneratedKey.Value)
			}
		}
	}
	fmt.Println()

	// Example 5: Delete account (optional - uncomment if you want to test)
	fmt.Println("5. Deleting account...")
	deleteInput := &vaultclient.DeleteAccountInput{}
	deleteInput.SetAccountName(accountName) // Use the same account name we've been working with

	_, err = client.DeleteAccount(ctx, deleteInput)
	if err != nil {
		log.Printf("Error deleting account: %v", err)
	} else {
		fmt.Printf("Account deleted successfully!\n")
	}
	fmt.Println()
}
