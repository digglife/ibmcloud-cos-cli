package aspera

import (
	"io"
	"log"

	sdk "github.com/IBM/ibmcloud-cos-cli/aspera/transfersdk"
	"gopkg.in/cheggaaa/pb.v1"
)

type Subscriber interface {
	Queued(resp *sdk.TransferResponse)
	Running(resp *sdk.TransferResponse)
	Done(resp *sdk.TransferResponse)
}

type DefaultSubscriber struct{}

func (b *DefaultSubscriber) Queued(resp *sdk.TransferResponse) {
	log.Printf("task %s queued", resp.TransferId)
}

func (b *DefaultSubscriber) Running(resp *sdk.TransferResponse) {
	log.Printf("transfered: %d", resp.TransferInfo.BytesTransferred)
}

func (b *DefaultSubscriber) Done(resp *sdk.TransferResponse) {
	log.Printf("task %s done", resp.TransferId)
}

type ProgressBarSubscriber struct {
	bar *pb.ProgressBar
}

func NewProgressBarSubscriber(total int, out io.Writer) *ProgressBarSubscriber {
	bar := pb.New(total).SetUnits(pb.U_BYTES)
	bar.Output = out
	return &ProgressBarSubscriber{bar: bar}
}

func (p *ProgressBarSubscriber) Queued(resp *sdk.TransferResponse) {
	p.bar.Prefix("Queued")
	p.bar.Start()
}

func (p *ProgressBarSubscriber) Running(resp *sdk.TransferResponse) {
	p.bar.Prefix("Running")
	p.bar.Set(int(resp.TransferInfo.BytesTransferred))
}

func (p *ProgressBarSubscriber) Done(resp *sdk.TransferResponse) {
	p.bar.Finish()
}
