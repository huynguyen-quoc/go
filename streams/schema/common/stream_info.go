package common

import (
	"time"

	"github.com/golang/protobuf/proto"

)

// StreamInfoEntity is the standard or base information used by the streaming code
type StreamInfoEntity struct {
	// The key by which the data is partitioned into shards
	StreamPartitionID int64
	// Time the data was added to the stream
	StreamTime time.Time
}

// FromPBBytes will convert from PB bytes to Info
func FromPBBytes(bytes []byte) StreamInfoEntity {
	info := &StreamInfo{}
	err := proto.Unmarshal(bytes, info)
	if err != nil {
		// swallow the error and return empty object
		return StreamInfoEntity{}
	}
	return FromPB(info)
}

// ToPBBytes will convert Info to PB bytes
func ToPBBytes(in StreamInfoEntity) []byte {
	info := ToPB(in)
	bytes, err := proto.Marshal(info)
	if err != nil {
		// swallow the error and return nil
		return nil
	}
	return bytes
}

// FromPB will convert from PB format to Info
func FromPB(in *StreamInfo) StreamInfoEntity {
	if in == nil {
		return StreamInfoEntity{}
	}
	return StreamInfoEntity{
		StreamPartitionID: in.StreamPartitionID,
		StreamTime:        time.Unix(in.StreamTime, 0),
	}
}

// ToPB will convert Info to from PB format
func ToPB(in StreamInfoEntity) *StreamInfo {
	return &StreamInfo{
		StreamPartitionID: in.StreamPartitionID,
		StreamTime:        in.StreamTime.Unix(),
	}
}
