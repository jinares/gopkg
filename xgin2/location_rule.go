package xgin2

type (
	Rule struct {
		IsRuning int               `json:"IsRuning,omitempty" yaml:"IsRuning"` //0:有效 1:无效
		Location string            `json:"Location",omitempty yaml:"Location"`
		Params   map[string]string `json:"Params,omitempty" yaml:"Params"` // 和 Location 一起进行匹配g

	}
)
