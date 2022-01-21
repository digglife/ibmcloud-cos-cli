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
	"runtime"
	"strings"

	"github.com/IBM/ibm-cos-sdk-go/aws/request"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/aspera/transfersdk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultAddress string = "127.0.0.1"
	defaultPort    string = "55002"
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
	transfersdk.TransferServiceClient
	s3 *s3.S3
}

func New(s3 *s3.S3) (*Agent, error) {
	optInsecure := grpc.WithTransportCredentials(insecure.NewCredentials())
	target := fmt.Sprintf("%s:%s", defaultAddress, defaultPort)
	cc, err := grpc.Dial(target, optInsecure)
	if err != nil {
		return nil, err
	}
	client := transfersdk.NewTransferServiceClient(cc)

	return &Agent{client, s3}, nil
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
	return exec.CommandContext(ctx, TransferdBinPath()).Start()
}

func (a *Agent) IsTransferdRunning() bool {
	if _, err := a.GetInfo(context.Background(), &transfersdk.InstanceInfoRequest{}); err != nil {
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

func (a *Agent) GetICOSSpec(bucket string) *transfersdk.ICOSSpec {
	ICOSSpec := &transfersdk.ICOSSpec{
		ApiKey:               "",
		Bucket:               "",
		IbmServiceInstanceId: "",
		IbmServiceEndpoint:   "",
	}
	return ICOSSpec
}

type Path struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func (a *Agent) GetAsperaTransferSpec(action string, bucket string, paths []Path) (string, error) {
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
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	specs := map[string][]map[string]map[string]interface{}{
		"transfer_specs": {
			{"transfer_spec": map[string]interface{}{}},
		},
	}
	if err := json.Unmarshal(body, &specs); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %s\n %s", err, string(body))
	}
	spec_map := specs["transfer_specs"][0]["transfer_spec"]
	spec, err := json.Marshal(spec_map)
	if err != nil {
		return "", fmt.Errorf("failed to marshal spec: %s", err)
	}
	fmt.Println(string(spec))
	return string(spec), nil
}

func (a *Agent) COSDownload(bucket string, key string, localPath string) error {
	return nil
}

func (a *Agent) COSUpload(localPath string, bucket string, key string) error {
	return nil
}

func (a *Agent) DoCOSTransfer(ctx context.Context, bucket string, action string, key string, localPath string) error {
	err := a.StartServer(context.TODO())
	if err != nil {
		return err
	}

	p := Path{key, localPath}
	//transferType := transfersdk.TransferType_FILE_TO_STREAM_DOWNLOAD
	if action == "upload" {
		p = Path{localPath, key}
		//transferType = transfersdk.TransferType_STREAM_TO_FILE_UPLOAD
	}

	transferSpec, err := a.GetAsperaTransferSpec(action, bucket, []Path{p})
	if err != nil {
		return err
	}

	req := &transfersdk.TransferRequest{
		TransferSpec: transferSpec,
		Config: &transfersdk.TransferConfig{
			Retry: &transfersdk.RetryStrategy{},
		},
		TransferType: transfersdk.TransferType_FILE_REGULAR,
	}
	steam, err := a.StartTransferWithMonitor(ctx, req)
	if err != nil {
		return err
	}

	started := false
	for {
		resp, err := steam.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch resp.Status {
		case transfersdk.TransferStatus_RUNNING:
			if !started && resp.SessionInfo.BytesTransferred != 0 {
				log.Println("Begin transfer")
				started = true
			} else if started {
				log.Printf("Transfered: %d", resp.SessionInfo.BytesTransferred)
			}
		case transfersdk.TransferStatus_FAILED, transfersdk.TransferStatus_CANCELED:
			log.Println("Failed")
		case transfersdk.TransferStatus_COMPLETED:
			log.Println("Finished")
		default:
			log.Println("Unknown")
		}
	}

	return nil

}
