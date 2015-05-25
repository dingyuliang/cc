package api

import (
	"github.com/jxwr/cc/topo"
)

type Response struct {
	Errno  int         `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Body   interface{} `json:"body"`
}

func MakeResponse(errno int, msg string, body interface{}) Response {
	return Response{errno, msg, body}
}

func MakeSuccessResponse(body interface{}) Response {
	return Response{0, "OK", body}
}

func MakeFailureResponse(msg string) Response {
	return Response{777, msg, nil}
}

type MapResp map[string]interface{}

type RegionSnapshotParams struct {
	Region      string            `json:"region"`
	PostTime    int64             `json:"posttime"`
	Nodes       []*topo.Node      `json:"nodes"`
	FailureInfo *topo.FailureInfo `json:"failure_info"`
}

type MigrateParams struct {
	SourceId string   `json:"source_id"`
	TargetId string   `json:"target_id"`
	Ranges   []string `json:"ranges"`
}

type ToggleModeParams struct {
	Action string `json:"action"`
	Perm   string `json:"perm"`
	NodeId string `json:"node_id"`
}

type MakeReplicaSetParams struct {
	NodeIds []string `json:"node_ids"`
}

type RebalanceParams struct {
	Method       string   `json:"method"`
	TargetIds    []string `json:"target_ids"`
	ShowPlanOnly bool     `json:"show_plan_only"`
}

type MeetNodeParams struct {
	NodeId string `json:"node_id"`
}

type SetAsMasterParams struct {
	NodeId string `json:"node_id"`
}

type ForgetAndResetNodeParams struct {
	NodeId string `json:"node_id"`
}

type ReplicateParams struct {
	ChildId  string `json:"child_id"`
	ParentId string `json:"parent_id"`
}

type FailoverTakeoverParams struct {
	NodeId string `json:"node_id"`
}
