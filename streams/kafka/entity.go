package kafka

import (
	"github.com/huynguyen-quoc/go/streams/core"
	"github.com/huynguyen-quoc/go/streams/schema/common"
)

//go:generate mockery --inpackage --case underscore --name Entity
type Entity interface {
	core.WriterDTO
	FromPB(pb core.Message) Entity
	GetStreamInfo() common.StreamInfoEntity
	GetPartitionID(id interface{}) int64
	GetPartitionKey() interface{}
	GetMessage() core.Message
}
