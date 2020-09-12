package ivstream

import (
	"github.com/huynguyen-quoc/go/streams/schema/common"
	"time"
)

// GetPartitionID ...
func (message *InvestmentOpen) GetStreamInfo() common.StreamInfoEntity {
	if message.StreamMetadata == nil {
		return common.StreamInfoEntity{}
	}

	return convert(message.StreamMetadata)
}

func convert(info *common.StreamInfo) common.StreamInfoEntity {
	return common.StreamInfoEntity{
		StreamPartitionID: info.StreamPartitionID,
		StreamTime:        time.Unix(info.StreamTime, 0),
	}
}
