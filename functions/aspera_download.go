package functions

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	sdk "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/aspera"
	"github.com/IBM/ibmcloud-cos-cli/config/fields"
	"github.com/IBM/ibmcloud-cos-cli/config/flags"
	"github.com/IBM/ibmcloud-cos-cli/errors"
	"github.com/IBM/ibmcloud-cos-cli/utils"
	"github.com/urfave/cli"
)

func AsperaDownload(c *cli.Context) (err error) {

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

	// Monitor the file
	keepFile := false

	// Download location
	var dstPath string

	// In case of error removes incomplete downloads
	defer func() {
		if !keepFile && dstPath != "" {
			cosContext.Remove(dstPath)
		}
	}()

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

	dstPath = c.Args().First()

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

	if err = asp.Download(ctx, aws.StringValue(input.Bucket), aws.StringValue(input.Key), dstPath); err != nil {
		return
	}

	keepFile = true

	return
}
