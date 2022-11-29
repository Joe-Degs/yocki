package api

import ylog "github.com/Joe-Degs/yocki/server/log"

type ProduceRequest struct {
	Record ylog.Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record ylog.Record `json:"record"`
}
