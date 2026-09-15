// models/thunder_cdp_model.go (100行以下 - SPEC-PRINCIPLE-001)
package models

// ActiveThunderTask は迅雷 CDP (taskBaseMap + DOM) から抽出された詳細タスク情報
type ActiveThunderTask struct {
	FileName       string `json:"file_name"`
	Target         string `json:"target"`
	Status         string `json:"status"`
	ProgressText   string `json:"progress_text"`
	SpeedText      string `json:"speed_text"`
	DetailText     string `json:"detail_text"`
	TaskID         int64  `json:"task_id"`
	TaskType       int    `json:"task_type"`
	TaskStatusCode int    `json:"task_status_code"`
	ErrorCode      int    `json:"error_code"`
	URL            string `json:"url"`
	SavePath       string `json:"save_path"`
	FileSize       int64  `json:"file_size"`
	DownloadSize   int64  `json:"download_size"`
	DownloadSpeed  int64  `json:"download_speed"`
	VIPSpeed       int64  `json:"vip_speed"`
	CreateTime     int64  `json:"create_time"`
	CompletionTime int64  `json:"completion_time"`
	SrcTotal       int    `json:"src_total"`
	SrcUsing       int    `json:"src_using"`
	Peers          int    `json:"peers"`
	CID            string `json:"cid"`
	GCID           string `json:"gcid"`
	Origin         string `json:"origin"`
	RawJSON        string `json:"raw_json"`
	RawText        string `json:"raw_text"`
}

// ThunderTaskInputItem は迅雷登録用の URL と正規化ファイル名
type ThunderTaskInputItem struct {
	URL      string `json:"url"`
	FileName string `json:"file_name"`
}
