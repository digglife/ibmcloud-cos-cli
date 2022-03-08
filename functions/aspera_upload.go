package functions

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	sdk "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/common/file_helpers"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/aspera"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/errors"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

func AsperaUpload(c *cli.Context) (err error) {

	// check the number of arguments
	if c.NArg() > 1 {
		err = &errors.CommandError{
			CLIContext: c,
			Cause:      errors.InvalidNArg,
		}
		return
	}

	APIKeyEnv := sdk.EnvAPIKey.Get()
	if APIKeyEnv == "" {
		err = fmt.Errorf("missing IBMCLOUD_API_KEY Environment Variable")
		return
	}

	// Load COS Context
	var cosContext *utils.CosContext
	if cosContext, err = GetCosContext(c); err != nil {
		return
	}

	if !file_helpers.FileExists(aspera.TransferdBinPath()) {
		cosContext.UI.Warn("Aspera Transferd binary not found. Downloading...")
		if err = DownloadSDK(cosContext); err != nil {
			return
		}
	}

	// Build GetObjectInput
	input := new(s3.GetObjectInput)

	// Required parameters for GetObjectInput
	mandatory := map[string]string{
		fields.Bucket: flags.Bucket,
		fields.Key:    flags.Key,
	}

	//
	// Optional parameters for GetObjectInput
	options := map[string]string{}

	//
	// Check through user inputs for validation
	if err = MapToSDKInput(c, input, mandatory, options); err != nil {
		return
	}

	srcPath := c.Args().First()

	if _, err = os.Stat(srcPath); os.IsNotExist(err) {
		err = fmt.Errorf("%s: %s", err, srcPath)
		return
	}

	var region string
	if region, err = cosContext.GetCurrentRegion(""); err != nil {
		return
	}

	var serviceEndpoint string
	if serviceEndpoint, err = cosContext.GetServiceEndpoint(); err != nil {
		return
	}

	cfg := new(aws.Config).WithRegion(region).WithEndpoint(serviceEndpoint)
	sess := cosContext.Session.Copy(cfg)
	client := s3.New(sess)
	asp, _ := aspera.New(client, APIKeyEnv)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	var size int64
	size, err = GetLocalTotalSize(srcPath)
	transferInput := &aspera.TransferInput{
		Bucket: aws.StringValue(input.Bucket),
		Key:    aws.StringValue(input.Key),
		Path:   srcPath,
		Sub:    aspera.NewProgressBarSubscriber(size, cosContext.UI.Writer()),
	}

	if err = asp.Upload(ctx, transferInput); err != nil {
		return
	}

	return
}

func GetLocalTotalSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
