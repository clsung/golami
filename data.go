package golami

import "encoding/json"

// jsonInt clones int, to work properly with "status" in "desc_obj" have both int/string type
type jsonInt int

func (i jsonInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(i))
}

func (i *jsonInt) UnmarshalJSON(data []byte) error {
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}

	var tmp int
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	*i = jsonInt(tmp)
	return nil
}

// {"data":{"nli":[{"desc_obj":{"result":"離春節還有42天","status":0},"type":"date"}]},"status":"ok"}
// {"data":{"nli":[{"desc_obj":{"result":"主人，請問你想查哪裡的天氣呢？","type":"weather","status":"0"},"type":"question"}]},"status":"ok"}
// {"data":{"seg":"再 過 幾 天 就是 過 年 "},"status":"ok"}

type DescObj struct {
	Result string
	Type   string `json:"type,omitempty"`
	Status jsonInt
}

type NLI struct {
	DescObj DescObj `json:"desc_obj"`
	Type    string
}

type ASR struct {
	Result       string
	SpeechStatus int `json:"speech_status"`
	Final        bool
	Status       int
}

type NLUResult struct {
	NLI []NLI  `json:"nli,omitempty"`
	SEG string `json:"seg,omitempty"`
	ASR ASR    `json:"asr,omitempty"`
}

type Result struct {
	Data   NLUResult
	Status string
}
