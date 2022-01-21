package functions

import (
	"context"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibmcloud-cos-cli/aspera"
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

	// Load COS Context
	var cosContext *utils.CosContext
	if cosContext, err = GetCosContext(c); err != nil {
		return
	}

	// Download location
	var dstPath string

	// Build GetObjectInput
	input := new(s3.GetObjectInput)

	// Required parameters for GetObjectInput
	// mandatory := map[string]string{
	// 	fields.Bucket: flags.Bucket,
	// 	fields.Key:    flags.Key,
	// }

	//
	// Check through user inputs for validation
	// if err = MapToSDKInput(c, input, mandatory, options); err != nil {
	// 	return
	// }

	//
	// Validate Download Location
	var file utils.WriteCloser
	if dstPath, file, err = getAndValidateDownloadPath(cosContext, c.Args().First(),
		aws.StringValue(input.Key)); err != nil || file == nil {
		return
	}

	var region string
	if region, err = cosContext.GetCurrentRegion(""); err != nil {
		return
	}
	serviceEndpoint, err := cosContext.GetServiceEndpoint()
	if err != nil {
		return
	}
	cfg := new(aws.Config).WithRegion(region).WithEndpoint(serviceEndpoint)
	sess := cosContext.Session.Copy(cfg)
	client := s3.New(sess)
	asp, _ := aspera.New(client)

	ctx := context.Background()
	return asp.DoCOSTransfer(ctx, flags.Bucket, "download", flags.Key, dstPath)
	// e := asp.IsTransferdRunning()
	// return fmt.Errorf("Running: %v; Dest: %s", e, dstPath)
	// Render DownloadOutput
	// output := &render.DownloadOutput{
	// 	TotalBytes: totalBytes,
	// }
	// Output the successful message
	// err = cosContext.GetDisplay(c.String(flags.Output), c.Bool(flags.JSON)).Display(input, output, nil)

	// Return
}
