package aspera

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials"
	"github.com/IBM/ibm-cos-sdk-go/aws/request"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/awstesting"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/aspera/mocks"
	sdk "github.com/IBM/ibmcloud-cos-cli/aspera/transfersdk"
	"github.com/stretchr/testify/assert"
)

func TestIsTransferdRunning(t *testing.T) {
	mockClient := new(mocks.TransferServiceClient)
	mockClient.On("GetInfo", context.Background(), &sdk.InstanceInfoRequest{}).Return(&sdk.InstanceInfoResponse{}, nil)

	dummySession, _ := session.NewSession(new(aws.Config))
	s3srv := s3.New(dummySession)
	agent := Agent{mockClient, s3srv, "apikey"}

	assert.True(t, agent.IsTransferdRunning())
}

func TestStartSever(t *testing.T) {

}

func TestSDKDir(t *testing.T) {
	os.Setenv("ASPERA_SDK_PATH", "/path/to/sdk")
	assert.Equal(t, "/path/to/sdk", SDKDir())
	os.Unsetenv("ASPERA_SDK_PATH")

	if home, err := os.UserHomeDir(); err == nil {
		assert.Equal(t, filepath.Join(home, ".aspera_sdk"), SDKDir())
	} else {
		assert.Equal(t, ".aspera_sdk", SDKDir())
	}
}

func TestTransferdBinPath(t *testing.T) {
	os.Setenv("ASPERA_SDK_PATH", "/path/to/sdk")
	if runtime.GOOS == "windows" {
		assert.Equal(t, "/path/to/sdk/bin/asperatransferd.exe", TransferdBinPath())
	} else {
		assert.Equal(t, "/path/to/sdk/bin/asperatransferd", TransferdBinPath())
	}
}

func TestGetBucketAspera(t *testing.T) {
	res := http.Response{
		StatusCode: 200,
		Body: body(`{
			"AccessKey":
				{
					"Id": "id",
					"Secret":"secret"
				},
			"ATSEndpoint": "https://zshengli.aspera.io:443"
		}`),
	}

	s := awstesting.NewClient(&aws.Config{})
	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.PushBack(unmarshal)
	s.Handlers.Send.Clear()
	s.Handlers.Send.PushBack(func(r *request.Request) {
		r.HTTPResponse = &res
	})

	transferdClient := new(mocks.TransferServiceClient)
	s3svc := &s3.S3{Client: s}
	agent := Agent{transferdClient, s3svc, "apikey"}

	info, _ := agent.GetBucketAspera("a-bucket")
	assert.Equal(t, info, &BucketAsperaInfo{
		AccessKey:   &AccessKey{Id: aws.String("id"), Secret: aws.String("secret")},
		ATSEndpoint: aws.String("https://zshengli.aspera.io:443"),
	})

}
func TestGetAsperaTransferSpecV2(t *testing.T) {
	c := credentials.NewCredentials(&stubProvider{
		creds: credentials.Value{
			ServiceInstanceID: "SVCINTID",
		},
		expired: true,
	})

	s := awstesting.NewClient(&aws.Config{
		Credentials: c,
		Endpoint:    aws.String("a-endpoint"),
	})
	s3svc := &s3.S3{Client: s}

	transferdClient := new(mocks.TransferServiceClient)
	a := Agent{transferdClient, s3svc, "apikey"}

	AsperaEndpoint := `{
		"AccessKey":
			{
				"Id": "id",
				"Secret":"secret"
			},
		"ATSEndpoint": "https://zshengli.aspera.io:443"
	}`

	res := []*http.Response{
		{
			StatusCode: 200,
			Body:       body(AsperaEndpoint),
		},
		{
			StatusCode: 200,
			Body:       body(AsperaEndpoint),
		},
	}

	reqNum := 0
	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.PushBack(unmarshal)
	s.Handlers.Send.Clear()
	s.Handlers.Send.PushBack(func(r *request.Request) {
		r.HTTPResponse = res[reqNum]
		reqNum++
	})

	sendPaths := []*sdk.Path{{Source: "/local/path/to/file", Destination: "key"}}
	sendSpec, _ := a.GetAsperaTransferSpecV2("upload", "a-bucket", sendPaths)

	expectedSendSpec := `{"session_initiation":{"icos":{"api_key":"apikey","bucket":"a-bucket","ibm_service_instance_id":"SVCINTID","ibm_service_endpoint":"a-endpoint"}},"assets":{"destination_root":"/","paths":[{"source":"/local/path/to/file","destination":"key"}]},"direction":"send","remote_host":"https://zshengli.aspera.io:443","title":"IBMCloud COS CLI"}`

	assert.Equal(t, expectedSendSpec, sendSpec)

	recvPaths := []*sdk.Path{{Source: "key", Destination: "/local/path/to/file"}}
	recvSpec, _ := a.GetAsperaTransferSpecV2("download", "a-bucket", recvPaths)

	expectedRecvSpec := `{"session_initiation":{"icos":{"api_key":"apikey","bucket":"a-bucket","ibm_service_instance_id":"SVCINTID","ibm_service_endpoint":"a-endpoint"}},"assets":{"destination_root":"/","paths":[{"source":"key","destination":"/local/path/to/file"}]},"direction":"recv","remote_host":"https://zshengli.aspera.io:443","title":"IBMCloud COS CLI"}`

	assert.Equal(t, expectedRecvSpec, recvSpec)
}

type stubProvider struct {
	creds   credentials.Value
	expired bool
	err     error
}

func (s *stubProvider) Retrieve() (credentials.Value, error) {
	s.expired = false
	s.creds.ProviderName = "stubProvider"
	return s.creds, s.err
}
func (s *stubProvider) IsExpired() bool {
	return s.expired
}

// Refer to ibm-cos-sdk-go/aws/request/request_test.go
func unmarshal(req *request.Request) {
	defer req.HTTPResponse.Body.Close()
	if req.Data != nil {
		json.NewDecoder(req.HTTPResponse.Body).Decode(req.Data)
	}
}

func body(str string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader([]byte(str)))
}
