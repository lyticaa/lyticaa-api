package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/lyticaa/lyticaa-api/internal/types"

	"bou.ke/monkey"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
)

const (
	jwToken      = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IlJrSkZRekF3UmtVeU16aEVSVVF6UkRjd05VTkNSakUyTlVJek1UaEVNakF5T0RJM01qYzBSUSJ9.eyJuaWNrbmFtZSI6ImJyZW50IiwibmFtZSI6ImJyZW50QGJtY2QuY28iLCJwaWN0dXJlIjoiaHR0cHM6Ly9zLmdyYXZhdGFyLmNvbS9hdmF0YXIvOGEyZjc2ZmIzNTE2ODg3YTVkMTE3ZDgyZGY1MDcwODM_cz00ODAmcj1wZyZkPWh0dHBzJTNBJTJGJTJGY2RuLmF1dGgwLmNvbSUyRmF2YXRhcnMlMkZici5wbmciLCJ1cGRhdGVkX2F0IjoiMjAxOS0xMi0xOVQxNTowNDo1Mi45ODNaIiwiaXNzIjoiaHR0cHM6Ly9kZXYtdWJvZmUyNTUuYXV0aDAuY29tLyIsInN1YiI6ImF1dGgwfDVkZTg5YWVhNWE2MTI4MGRlMWYxYmYyYiIsImF1ZCI6Ik5CV25zZTVXOXhDWUw3c0RaVzhldDJxOVlHV0M0MHBWIiwiaWF0IjoxNTc2ODE4NjMxLCJleHAiOjE1NzY4NTQ2MzEsImFjciI6Imh0dHA6Ly9zY2hlbWFzLm9wZW5pZC5uZXQvcGFwZS9wb2xpY2llcy8yMDA3LzA2L211bHRpLWZhY3RvciIsImFtciI6WyJtZmEiXX0.nLvYGO31_cpgYln1Zs4bsmu20TL5l5a59puzUuu1dUiJG2iF2UdJLEWl86WalDVGB6azmFwiAhLHljmxmHniczua0o2ood3HXFdP0IgPUnUF8_3az4ZdA5UxQzrW45OlS1zSO-WcVpgdwsrzWXwpLqCtu-kGK3jyB3DlysqDeIeK_DRtOH7avDh6YqkIoVmDu0VN7fq6UgzDdQfV63ZORvsEGEG1oufkAI__eEiUvRQ6Ce-kFVhTdjzmw-h82PTe7xVOQ20gZVEVwJQYm4EUMw-kZrym3ahT3dvnRaeZmE1DzmfSsLEjCHa1HWQ4KKqpTrQOP1RHh48MRz6WALkRVg"
	jwksResponse = `{"keys":[{"alg":"RS256",
							  "kty":"RSA",
							  "use":"sig",
							  "n":"6P2-6nzKB6_uhd1ejyhhWrqr4rxliNgiHvIZztDEv5XYAwRKDEWUGCSbcVZZUsAxgc20XOJGSU7TMmbVC_sFAE4998SvjihEBdtA7BolFuQ30CiJfiXYBSNDPSu3kqBmtYrwnKdObq7YdXnxFUdtVGx1IJ0FQs2np2iYVY1JNq47ozdzlTwq2fh5YVREdJZ7DzuIFm7HVhJckgcXTHX3O0jVPcHFMnlmZk094RO3drioxLDzaAejflVgaRtGDREI3oKyjbkEQKiQK0gsILZ08tj6rr2tNBRlT-jjseIeJ4lRO9Cr9mjaVl4k5h0M-hWaZUMOYXYmYiLtBcmpYhGMpw",
							  "e":"AQAB",
							  "kid":"RkJFQzAwRkUyMzhERUQzRDcwNUNCRjE2NUIzMThEMjAyODI3Mjc0RQ",
							  "x5t":"RkJFQzAwRkUyMzhERUQzRDcwNUNCRjE2NUIzMThEMjAyODI3Mjc0RQ",
							  "x5c":["MIIDBzCCAe+gAwIBAgIJdaeD3QGvtdKqMA0GCSqGSIb3DQEBCwUAMCExHzAdBgNVBAMTFmRldi11Ym9mZTI1NS5hdXRoMC5jb20wHhcNMTkxMjA1MDQ0NzA1WhcNMzMwODEzMDQ0NzA1WjAhMR8wHQYDVQQDExZkZXYtdWJvZmUyNTUuYXV0aDAuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA6P2+6nzKB6/uhd1ejyhhWrqr4rxliNgiHvIZztDEv5XYAwRKDEWUGCSbcVZZUsAxgc20XOJGSU7TMmbVC/sFAE4998SvjihEBdtA7BolFuQ30CiJfiXYBSNDPSu3kqBmtYrwnKdObq7YdXnxFUdtVGx1IJ0FQs2np2iYVY1JNq47ozdzlTwq2fh5YVREdJZ7DzuIFm7HVhJckgcXTHX3O0jVPcHFMnlmZk094RO3drioxLDzaAejflVgaRtGDREI3oKyjbkEQKiQK0gsILZ08tj6rr2tNBRlT+jjseIeJ4lRO9Cr9mjaVl4k5h0M+hWaZUMOYXYmYiLtBcmpYhGMpwIDAQABo0IwQDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBRguX78r4+iu/oh2N5iOUrXEg/bWjAOBgNVHQ8BAf8EBAMCAoQwDQYJKoZIhvcNAQELBQADggEBAH46r5i545ulJagXPDHVvPVargzgEffa78+8cL8RT4CFdPjGB6hWeOv4TvDlwx6s5Dm4hfDXOVaCWwaIs7vR4UEERiYHNfo1vFt5yycPU0iLDjhkdClP6IqLaFJDwDetWGzRyEZI/akkjZFsDhVNcrcbu8dFfg3Hk1koFTzKvFVC0dyYLo125RSUDO8snuqF1FlIaA2aLzWUGmY4EEVBp2AM29LsLejO3GWiVMg420HiQShJoHWvUL9Z6cuQBNBKxHW/Yy+qKXELwxViI5dlLp5KLVQq/S5Do044K+wuZOJxeQVVpU8V7jMWRgA9fyAsdV6gp3V7bO4zXC88B7g42OQ="]}]
					}`
)

func configEnv() {
	_ = os.Setenv("JWKS_URL", "http://dev-ubofe255.auth0.com/.well-known/jwks.json")
	_ = os.Setenv("JWT_AUD", "NBWnse5W9xCYL7sDZW8et2q9YGWC40pV")
	_ = os.Setenv("JWT_ISS", "https://dev-ubofe255.auth0.com/")
	_ = os.Setenv("AWS_S3_UPLOAD_BUCKET", "com.sellernomics.test.reports")
}

func httpClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return client, s.Close
}

func TestUsersUploadUrl(t *testing.T) {
	configEnv()

	// Monkey patch the time.
	wayback := time.Date(2019, time.December, 20, 5, 20, 0, 0, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	a := NewApi()

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(jwksResponse))
	})
	client, teardown := httpClient(h)
	defer teardown()

	a.Client = client

	backend := s3mem.New()
	faker := gofakes3.New(backend)
	ts := httptest.NewServer(faker.Server())
	defer ts.Close()

	config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials("AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", ""),
		Endpoint:         aws.String(ts.URL),
		Region:           aws.String("us-west-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	sess, _ := session.NewSession(config)
	s3Client := s3.New(sess)
	a.usersUploadUrl(s3Client)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/upload_url", nil)
	req.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", jwToken)},
	}

	resp := httptest.NewRecorder()
	a.Router.ServeHTTP(resp, req)

	if http.StatusOK != resp.Code {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, resp.Code)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var url types.Url

	if err := json.Unmarshal(body, &url); err != nil {
		t.Error("Unable to parse response body.")
	}

	if url.Url == "" {
		t.Error("Expected a signed URL but received an empty string.")
	}
}
