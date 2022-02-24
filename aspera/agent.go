package aspera

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strings"

	"github.com/IBM/ibm-cos-sdk-go/aws/request"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	sdk "github.com/IBM/ibmcloud-cos-cli/aspera/transfersdk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultAddress = "127.0.0.1"
	defaultPort    = "55002"
)

type AccessKey struct {
	_      struct{} `type:"structure"`
	Id     *string  `type:"string"`
	Secret *string  `type:"string"`
}

type BucketAsperaInfo struct {
	AccessKey   *AccessKey `type:"structure"`
	ATSEndpoint *string    `type:"string"`
}

type Agent struct {
	sdk.TransferServiceClient
	s3     *s3.S3
	apikey string
}

func New(s3 *s3.S3, apikey string) (*Agent, error) {
	optInsecure := grpc.WithTransportCredentials(insecure.NewCredentials())
	target := fmt.Sprintf("%s:%s", defaultAddress, defaultPort)
	cc, err := grpc.Dial(target, optInsecure)
	if err != nil {
		return nil, err
	}
	client := sdk.NewTransferServiceClient(cc)

	return &Agent{client, s3, apikey}, nil
}

func SDKDir() string {
	if sdk_path, ok := os.LookupEnv("ASPERA_SDK_PATH"); ok {
		return sdk_path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return path.Join(home, ".aspera_sdk")
}

func TransferdBinPath() string {
	daemonName := "asperatransferd"
	if runtime.GOOS == "windows" {
		daemonName = "asperatransferd.exe"
	}
	return path.Join(SDKDir(), "bin", daemonName)
}

func (a *Agent) StartServer(ctx context.Context) error {
	if a.IsTransferdRunning() {
		return nil
	}
	transferd := TransferdBinPath()
	err := exec.CommandContext(ctx, transferd).Start()
	if err != nil {
		return fmt.Errorf("failed to start asperatransferd(%s): %v", transferd, err)
	}
	return nil
}

func (a *Agent) IsTransferdRunning() bool {
	if _, err := a.GetInfo(context.Background(), &sdk.InstanceInfoRequest{}); err != nil {
		return false
	}
	return true
}

func (a *Agent) GetBucketAspera(bucket string) (*BucketAsperaInfo, error) {
	opGetBucketAspera := &request.Operation{
		Name:       "GetBucketAspera",
		HTTPMethod: "GET",
		HTTPPath:   fmt.Sprintf("/%s?faspConnectionInfo", "backint"),
	}

	output := &BucketAsperaInfo{}
	req := a.s3.NewRequest(opGetBucketAspera, nil, output)
	if err := req.Send(); err != nil {
		return nil, err
	}
	return output, nil
}

func (a *Agent) GetICOSSpec(bucket string) *sdk.ICOSSpec {
	creds, _ := a.s3.Config.Credentials.Get()
	ICOSSpec := &sdk.ICOSSpec{
		ApiKey:               a.apikey,
		Bucket:               bucket,
		IbmServiceInstanceId: creds.ServiceInstanceID,
		IbmServiceEndpoint:   a.s3.Endpoint,
	}
	return ICOSSpec
}

func (a *Agent) GetAsperaTransferSpecV2(action string, bucket string, paths []*sdk.Path) (spec string, err error) {

	ICOSSpec := a.GetICOSSpec(bucket)
	direction := "recv"
	if action == "upload" {
		direction = "send"
	}

	var meta *BucketAsperaInfo
	if meta, err = a.GetBucketAspera(bucket); err != nil {
		return
	}

	transferSpec := &sdk.TransferSpecV2{
		SessionInitiation: &sdk.Initiation{
			Icos: ICOSSpec,
		},
		Direction: direction,
		Assets: &sdk.Assets{
			// The definition in the original proto is Path, not Paths
			// Reported to Aspera Team
			Paths:           paths,
			DestinationRoot: "/",
		},
		RemoteHost: *meta.ATSEndpoint,
		Title:      "IBMCloud COS CLI",
	}

	data, err := json.Marshal(transferSpec)
	if err != nil {
		return "", fmt.Errorf("unable to marchal transferspecv2: %s", err)
	}

	spec = string(data)
	return
}

func (a *Agent) GetAsperaTransferSpecV1(action string, bucket string, paths []*sdk.Path) (spec string, err error) {
	meta, err := a.GetBucketAspera(bucket)
	if err != nil {
		return "", fmt.Errorf("unable to get aspera metadata: %s", err)
	}
	creds, err := a.s3.Config.Credentials.Get()
	if err != nil {
		return "", fmt.Errorf("unable to get aws credentials: %s", err)
	}
	credentials := fmt.Sprintf(`{"type":"token","token":{ "delegated_refresh_token": "%s"}}`, creds.Token.AccessToken)

	j, err := json.Marshal(paths)
	if err != nil {
		return "", err
	}
	jsonPaths := string(j)

	// The type of `tags` in the original proto file is string
	// so I can't use TransferSpecV1 struct directly.
	// Reported to Aspera team
	data := fmt.Sprintf(`{
		"transfer_requests": [
		  {
			"transfer_request": {
			  "paths": %s,
			  "tags": {
				"aspera": {
				  "node": {
					"storage_credentials": %s
				  }
				}
			  }
			}
		  }
		]
	  }
	  `, jsonPaths, credentials)

	url := fmt.Sprintf("%s/files/%s_setup", *meta.ATSEndpoint, action)
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Aspera-Storage-Credentials", credentials)
	req.SetBasicAuth(*meta.AccessKey.Id, *meta.AccessKey.Secret)
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to %s: %v", url, err)
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP Error: %v : %s", res.StatusCode, res.Request.URL)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	type reqErr struct {
		Code    int    `json:"code,omitempty"`
		Reason  string `json:"reason,omitempty"`
		Message string `json:"user_message,omitempty"`
	}

	if matched, _ := regexp.Match(`"error":`, body); matched {
		var e reqErr
		if err := json.Unmarshal(body, &e); err != nil {
			return "", fmt.Errorf("request error: %d: %s", e.Code, e.Message)
		}
	}

	specs := map[string][]map[string]map[string]interface{}{
		"transfer_specs": {
			{"transfer_spec": map[string]interface{}{}},
		},
	}
	if err := json.Unmarshal(body, &specs); err != nil {
		return "", fmt.Errorf("failed to unmarshal transferspecs: %s: %s", err, string(body))
	}
	spec_map := specs["transfer_specs"][0]["transfer_spec"]
	j, err = json.Marshal(spec_map)
	if err != nil {
		return "", fmt.Errorf("failed to marshal spec: %s", err)
	}
	spec = string(j)
	return
}

func (a *Agent) Download(ctx context.Context, bucket string, key string, localPath string) (err error) {
	return a.DoTransfer(ctx, "download", bucket, key, localPath)
}

func (a *Agent) Upload(ctx context.Context, localPath string, bucket string, key string) (err error) {
	return a.DoTransfer(ctx, "upload", bucket, key, localPath)
}

func (a *Agent) DoTransfer(ctx context.Context, action string, bucket string, key string, localPath string) (err error) {
	rpcCtx := context.TODO()
	if err = a.StartServer(rpcCtx); err != nil {
		return
	}

	p := &sdk.Path{Source: key, Destination: localPath}
	if action == "upload" {
		p = &sdk.Path{Source: localPath, Destination: key}
	}

	transferSpec, err := a.GetAsperaTransferSpecV2(action, bucket, []*sdk.Path{p})
	if err != nil {
		return
	}

	req := &sdk.TransferRequest{
		TransferSpec: transferSpec,
		Config: &sdk.TransferConfig{
			Retry: &sdk.RetryStrategy{},
		},
		TransferType: sdk.TransferType_FILE_REGULAR,
	}
	transferResp, err := a.StartTransfer(ctx, req)
	if err != nil {
		return
	}
	transferId := transferResp.GetTransferId()

	// This only can cancel the transfer that already started
	// It can't cancel the transfer request (should also be a design flaw...)
	// TODO: Ask aspera team if they can support cancel transfer request
	go func() {
		<-ctx.Done()
		stop := &sdk.StopTransferRequest{TransferId: []string{transferId}}
		if _, err := a.StopTransfer(rpcCtx, stop); err != nil {
			log.Println("failed to stop transfer:", err)
		}
	}()

	stream, err := a.MonitorTransfers(rpcCtx, &sdk.RegistrationRequest{TransferId: []string{transferId}})
	if err != nil {
		return
	}

	// stream, err := a.StartTransferWithMonitor(ctx, req)

	var started bool
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch resp.Status {
		case sdk.TransferStatus_QUEUED:
			log.Printf("task %s queued", resp.TransferId)
		case sdk.TransferStatus_RUNNING:
			if !started && resp.TransferInfo.BytesTransferred == 0 {
				log.Printf("task %s started", resp.TransferId)
				started = true
			} else {
				log.Printf("transfered: %d", resp.TransferInfo.BytesTransferred)
			}
		case sdk.TransferStatus_FAILED, sdk.TransferStatus_CANCELED:
			log.Println("failed or cancelled")
			return fmt.Errorf("transfer %s: %s", resp.Status, resp.Error.GetDescription())
		case sdk.TransferStatus_COMPLETED:
			log.Println("finished")
			// MonitorTransfers doesn't works like StartTransferWithMonitor,
			// the response it returns doesn't emit EOF because of multiple transfers
			// so the loop will block infinitely even the transfer's finished.
			// I have to return here explicitly.
			return nil
		}
	}

	return
}
