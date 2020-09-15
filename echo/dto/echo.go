package dto

import v1 "github.com/huynguyen-quoc/go/echo/pb/v1"

type EchoRequest struct {
	Value string `json:"value"`
}

type EchoResponse struct {
	Value string `json:"value"`
}

func (e EchoRequest) ToPB() interface{} {
	return &v1.EchoRequest{
		Value: e.Value,
	}
}

func (e EchoRequest) fromPB(pb interface{}) interface{} {
	if pb == nil {
		return nil
	}

	pbData, ok := pb.(*v1.EchoRequest)
	if !ok {
		return nil
	}
	// int64
	pbValue := pbData.Value

	e.Value = pbValue
	return e
}


func (e EchoResponse) ToPB() interface{} {
	return &v1.EchoRequest{
		Value: e.Value,
	}
}

func (e EchoResponse) fromPB(pb interface{}) interface{} {
	if pb == nil {
		return nil
	}

	pbData, ok := pb.(*v1.EchoResponse)
	if !ok {
		return nil
	}
	// int64
	pbValue := pbData.Value

	e.Value = pbValue
	return e
}
