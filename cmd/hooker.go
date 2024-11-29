package main

type HookResponse map[string]interface{}

func (r HookResponse) Action() string {
	return r["action"].(string)
}

func (r HookResponse) TaskDescription() string {
	return r["data"].(map[string]interface{})["description"].(string)
}

func (r HookResponse) CreatedBy() string {
	return r["data"].(map[string]interface{})["owner"].(map[string]interface{})["username"].(string)
}

func (r HookResponse) PermaLink() string {
	return r["data"].(map[string]interface{})["permalink"].(string)
}

func (r HookResponse) AssignedToUserName() string {
	return r["data"].(map[string]interface{})["assigned_to"].(map[string]interface{})["username"].(string)
}

type NotificationData struct {
	Action      string
	Owner       string
	Description string
	PermaLink   string
}
