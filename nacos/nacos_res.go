package nacos

type ConfigRes struct {
	Group       string      `json:"group"`
	AppName     string      `json:"appName"`
	DataId      string      `json:"dataId"`
	Tenant      string      `json:"tenant"`
	ContentType string      `json:"type"`
	ConfigTags  interface{} `json:"configTags"`
	Content     string      `json:"content"`
}
