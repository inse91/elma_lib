package e365_gateway

import "time"

const (
	StateDone = "done"
	StateExec = "exec"
)

type ProcCommon struct {
	ID                        string      `json:"__id"`
	CreatedAt                 time.Time   `json:"__createdAt"`
	CreatedBy                 string      `json:"__createdBy"`
	CurrentPerformers         interface{} `json:"__currentPerformers"`
	DebugUserId               interface{} `json:"__debugUserId"`
	FieldVisibilityConditions struct {
	} `json:"__fieldVisibilityConditions"`
	FlowHistory struct {
		Transitions struct {
		} `json:"transitions"`
		Items struct {
		} `json:"items"`
	} `json:"__flowHistory"`
	InformLists struct {
	} `json:"__informLists"`
	Logged        bool        `json:"__logged"`
	Name          string      `json:"__name"`
	NameTemplate  interface{} `json:"__nameTemplate"`
	NotifyOnStart bool        `json:"__notifyOnStart"`
	Parent        interface{} `json:"__parent"`
	Performers    struct {
	} `json:"__performers"`
	State       string   `json:"__state"`
	Subscribers []string `json:"__subscribers"`
	Tasks       struct {
	} `json:"__tasks"`
	Template struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
		Code      string `json:"code"`
		Version   int    `json:"version"`
	} `json:"__template"`
	TemplateId    string    `json:"__templateId"`
	UpdatedAt     time.Time `json:"__updatedAt"`
	UpdatedBy     string    `json:"__updatedBy"`
	CrossInstance struct {
		IsCrossInstance  bool   `json:"is_cross_instance"`
		CompanyInitiator string `json:"company_initiator"`
		ParentTask       string `json:"parent_task"`
		ParentInstance   string `json:"parent_instance"`
		ParentBranch     string `json:"parent_branch"`
		ParentPath       string `json:"parent_path"`
		Callback         string `json:"callback"`
	} `json:"cross_instance"`
	//Goods interface{} `json:"goods"`
}
