package types

import (
	"encoding/json"
)

type XrayTraceSegmentDocument struct {
	TraceId   string  `json:"trace_id"`
	Id        string  `json:"id"`
	StartTime float64 `json:"start_time"`
	EndTime   float64 `json:"end_time"`
	Name      string  `json:"name"`
}

func (r *XrayTraceSegmentDocument) ToString() string {
	json, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(json)
}
