package golami

// {"data":{"nli":[{"desc_obj":{"result":"離春節還有42天","status":0},"type":"date"}]},"status":"ok"}
// {"data":{"seg":"再 過 幾 天 就是 過 年 "},"status":"ok"}

type DescObj struct {
	Result string
	Status int
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
