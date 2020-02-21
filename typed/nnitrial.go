package typed

import (
	"encoding/json"
	"fmt"
	"os"
)

type NNITrial struct {
	Id              string                   `json:"id"`
	Status          string                   `json:"status"`
	HyperParameters []string                 `json:"hyperParameters"`
	StartTime       int64                    `json:"startTime"`
	EndTime         int64                    `json:"endTime"`
	MetricData      []map[string]interface{} `json:"finalMetricData"`
}

type hyperParameter struct {
	ParameterId int                    `json:"parameter_id"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type trialResult struct {
	Sequence   int                      `json:"sequence"`
	Status     string                   `json:"status"`
	StartTime  int64                    `json:"startTime"`
	EndTime    int64                    `json:"endTime"`
	MetricData []map[string]interface{} `json:"finalMetricData"`
	Param      map[string]interface{}   `json:"params"`
}

func (t NNITrial) String() string {
	return fmt.Sprintf("ID:%s\nStatus:%s\nHyperParameters:%s\nStartTime:%v\nEndTime:%v\nMetricData:%v\n", t.Id,
		t.Status, t.HyperParameters[0], t.StartTime, t.EndTime, t.MetricData[0])
}

func (t NNITrial) getParameter() *hyperParameter {
	var params hyperParameter
	if len(t.HyperParameters) == 0 {
		return nil
	}
	err := json.Unmarshal([]byte(t.HyperParameters[0]), &params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unmashal Parameter fail:%v\n", err)
		return nil
	}
	return &params
}

func (t NNITrial) ToJSON() string {
	param := t.getParameter()
	temp := trialResult{
		Status:     t.Status,
		StartTime:  t.StartTime,
		EndTime:    t.EndTime,
		MetricData: t.MetricData,
	}
	if param != nil {
		temp.Sequence = param.ParameterId
		temp.Param = param.Parameters
	}

	res, err := json.Marshal(temp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Marshal error:%v\n", err)
		return ""
	}
	return string(res)
}
