package ivstream

import (
	"fmt"
	"github.com/huynguyen-quoc/go/streams/core"
	"github.com/huynguyen-quoc/go/streams/kafka"
	"time"

	"github.com/huynguyen-quoc/go/streams/schema/common"
)

type InvestmentOpenEntity struct {
	Message        string             `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	InvestmentId   int64              `protobuf:"varint,2,opt,name=investmentId,proto3" json:"investmentId,omitempty"`
	UserId         int64              `protobuf:"varint,3,opt,name=userId,proto3" json:"userId,omitempty"`
	StreamMetadata *common.StreamInfo `protobuf:"bytes,1235,opt,name=StreamMetadata,proto3" json:"StreamMetadata,omitempty"`
}


// GetPartitionKey ...
func (entity *InvestmentOpenEntity) GetPartitionKey() interface{} {
	return entity.InvestmentId
}

// GetStreamInfo ...
func (entity *InvestmentOpenEntity) GetStreamInfo() common.StreamInfoEntity {
	return convert(entity.StreamMetadata)
}

// ToProtoBuf ...
func (entity *InvestmentOpenEntity) ToProtoBuf(isUpdate bool) core.Message {
	out := &InvestmentOpen{
		Message:      entity.Message,
		InvestmentId: entity.InvestmentId,
		UserId:       entity.UserId,
	}

	out.StreamMetadata = &common.StreamInfo{
		StreamPartitionID: entity.GetPartitionID(entity.GetPartitionKey()),
		StreamTime:        time.Now().UnixNano(),
	}

	return out
}

// FromPB ...
func (entity *InvestmentOpenEntity) FromPB(pb core.Message) kafka.Entity {
	in := pb.(*InvestmentOpen)
	out := &InvestmentOpenEntity{
		Message:      in.Message,
		InvestmentId: in.InvestmentId,
		UserId:       in.UserId,
	}

	// custom stream metadata used by the Coban team
	out.StreamMetadata = in.StreamMetadata

	return out
}

// GetPartitionID is responsible for returning a valid partitionId
func (entity *InvestmentOpenEntity) GetPartitionID(id interface{}) int64 {
	idInt64, ok := id.(int64)
	if !ok {
		fmt.Println("Error. `id` is not of int64 type")
	}

	// Ensure this stays positive
	return int64(idInt64 % 1000000000)
}

// GetMessage ...
func (entity *InvestmentOpenEntity) GetMessage() core.Message {
	return &InvestmentOpen{
		InvestmentId: entity.InvestmentId,
		Message: entity.Message,
		UserId: entity.UserId,
	}
}
