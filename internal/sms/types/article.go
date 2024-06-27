package types

// CanalLikeMsg canal解析like binlog消息.
type CanalLikeMsg struct {
	ID         string `json:"id"`
	BizID      string `json:"biz_id"`
	ObjID      string `json:"obj_id"`
	LikeNum    string `json:"like_num"`
	DislikeNum string `json:"dislike_num"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}
