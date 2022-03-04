package aspera

import (
	"log"

	sdk "github.com/IBM/ibmcloud-cos-cli/aspera/transfersdk"
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
