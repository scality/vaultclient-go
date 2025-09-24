package vaultclient

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

// Vault provides the API operation methods for making requests to Vault.
type Vault struct {
	endpoint   string
	region     string
	httpClient *http.Client
	creds      aws.CredentialsProvider
}

// Service information constants
const (
	ServiceName = "iam" // Name of service.
)

func New(accessKey, secretKey, sessionToken, endpoint, region string) *Vault {
	return &Vault{
		endpoint:   endpoint,
		region:     region,
		httpClient: &http.Client{Timeout: 60 * time.Second},
		creds: credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			sessionToken,
		),
	}
}

// makeAWSRequest is a helper function that handles AWS v4 signing and HTTP requests with retry logic
// This centralizes the common request logic used by all API operations
func (c *Vault) makeAWSRequest(ctx context.Context, action string, formData url.Values) (*http.Response, error) {
	formData.Set("Action", action)
	formData.Set("Version", "2010-05-08")

	requestBody := formData.Encode()

	maxAttempts := 3 // Use retry similar to default config here : https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/configure-retries-timeouts.html
	delay := 2 * time.Second

	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, strings.NewReader(requestBody))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Host", req.URL.Host)

		// Hash the request body payload for AWS signing
		hash := sha256.Sum256([]byte(requestBody))
		payloadHash := hex.EncodeToString(hash[:])

		// Sign the request using AWS SDK signer
		signer := v4.NewSigner()
		creds, err := c.creds.Retrieve(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve credentials: %w", err)
		}

		err = signer.SignHTTP(ctx, creds, req, payloadHash, ServiceName, c.region, time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to sign request: %w", err)
		}

		resp, err := c.httpClient.Do(req)
		if err == nil {
			// Check if we should retry based on status code
			if resp.StatusCode < 500 && resp.StatusCode != 429 {
				return resp, nil
			}
			// Server error or rate limit - close body and retry
			if closeErr := resp.Body.Close(); closeErr != nil {
				return nil, fmt.Errorf("failed to close response body: %w", closeErr)
			}
			lastErr = fmt.Errorf("server error: status %d", resp.StatusCode)
		} else {
			lastErr = err
		}

		// Don't sleep on the last attempt
		if attempt < maxAttempts-1 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
			delay *= 2
		}
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", maxAttempts, lastErr)
}

type APIError struct {
	StatusCode int
	Code       string
	Err        error
}

func (h *APIError) Error() string {
	return h.Err.Error()
}

func (h *APIError) Unwrap() error {
	return h.Err
}

func (h *APIError) HTTPStatusCode() int {
	return h.StatusCode
}

func (h *APIError) GetCode() string {
	return h.Code
}

func (vc *Vault) handleAWSResponse(resp *http.Response, result interface{}) (err error) {
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close response body: %w", closeErr)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var xmlError struct {
			Error struct {
				Code    string `xml:"Code"`
				Message string `xml:"Message"`
			} `xml:"Error"`
			RequestID string `xml:"RequestId"`
		}

		if err := xml.Unmarshal(respBody, &xmlError); err != nil {
			message := string(respBody)
			return &APIError{
				StatusCode: resp.StatusCode,
				Err: &types.ServiceFailureException{
					Message: &message,
				},
			}
		}

		switch xmlError.Error.Code {
		case (&types.EntityAlreadyExistsException{}).ErrorCode():
			return &APIError{
				StatusCode: resp.StatusCode,
				Code:       xmlError.Error.Code,
				Err: &types.EntityAlreadyExistsException{
					Message: &xmlError.Error.Message,
				},
			}
		case (&types.InvalidInputException{}).ErrorCode():
			return &APIError{
				StatusCode: resp.StatusCode,
				Code:       xmlError.Error.Code,
				Err: &types.InvalidInputException{
					Message: &xmlError.Error.Message,
				},
			}
		case (&types.NoSuchEntityException{}).ErrorCode():
			return &APIError{
				StatusCode: resp.StatusCode,
				Code:       xmlError.Error.Code,
				Err: &types.NoSuchEntityException{
					Message: &xmlError.Error.Message,
				},
			}
		case (&types.DeleteConflictException{}).ErrorCode():
			return &APIError{
				StatusCode: resp.StatusCode,
				Code:       xmlError.Error.Code,
				Err: &types.DeleteConflictException{
					Message: &xmlError.Error.Message,
				},
			}
		default:
			return &APIError{
				StatusCode: resp.StatusCode,
				Code:       xmlError.Error.Code,
				Err: &types.ServiceFailureException{
					Message: &xmlError.Error.Message,
				},
			}
		}
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}
