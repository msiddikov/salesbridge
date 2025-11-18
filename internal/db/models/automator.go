package models

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ---------- Enums ----------

type PublishState string

const (
	StateDraft    PublishState = "draft"
	StateActive   PublishState = "active"
	StateArchived PublishState = "archived"
)

type NodeKind string

const (
	KindTrigger    NodeKind = "trigger"
	KindCondition  NodeKind = "condition"
	KindAction     NodeKind = "action"
	KindControl    NodeKind = "control"
	KindCollection NodeKind = "collection"
)

type BackoffStrategy string

const (
	BackoffNone        BackoffStrategy = "none"
	BackoffFixed       BackoffStrategy = "fixed"
	BackoffExponential BackoffStrategy = "exponential"
)

type OnNodeError string

const (
	NodeErrorFailRun        OnNodeError = "fail_run"
	NodeErrorContinue       OnNodeError = "continue"
	NodeErrorRouteErrorPort OnNodeError = "route_error_port"
)

// ---------- Persistence Models (GORM) ----------

// Automation is the aggregate root. Nodes/Edges are children via AutomationID.
// The JSON "graph" field is a transient DTO (ignored by GORM) for API convenience.
type Automation struct {
	// Primary key (UUID/ULID as string). If you want DB-generated UUIDs,
	// switch to: ID uuid.UUID and use a BeforeCreate hook.
	ID string `json:"id" gorm:"type:uuid;primaryKey"`

	Name        string       `json:"name" gorm:"not null"`
	Description string       `json:"description,omitempty"`
	LocationId  string       `json:"locationId,omitempty" gorm:"index"`
	Location    Location     `json:"-" gorm:"foreignKey:LocationId;constraint:OnDelete:SET NULL"`
	State       PublishState `json:"state" gorm:"type:text;not null;default:'draft'"`

	CreatedAt time.Time `json:"createdAt" gorm:"not null"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"not null"`

	CreatorId uint `json:"-"`
	UpdaterId uint `json:"-"`
	Creator   User `json:"-"`
	Updater   User `json:"-"`

	// Ordered list of entry node IDs stored in JSONB (keeps order).
	Entry datatypes.JSON `json:"-" gorm:"type:jsonb;not null;default:'[]'::jsonb"`

	// Optional notes, stored as JSONB array of strings for simplicity.
	Notes datatypes.JSON `json:"-" gorm:"type:jsonb;not null;default:'[]'::jsonb"`

	// Children
	Nodes []Node `json:"-" gorm:"foreignKey:AutomationID;constraint:OnDelete:CASCADE"`
	Edges []Edge `json:"-" gorm:"foreignKey:AutomationID;constraint:OnDelete:CASCADE"`

	// -------- Transient/API-only --------
	Graph Graph `json:"graph" gorm:"-"`
}

// Node is persisted as a row; Config is JSONB.
type Node struct {
	ID           string         `json:"id" gorm:"type:uuid;primaryKey"`
	AutomationID string         `json:"-" gorm:"type:uuid;not null;index:idx_nodes_automation"`
	Automation   Automation     `json:"-" gorm:"foreignKey:AutomationID;constraint:OnDelete:CASCADE"`
	Type         string         `json:"type" gorm:"not null;index"`
	Kind         NodeKind       `json:"kind" gorm:"type:text;not null;index"`
	Name         string         `json:"name,omitempty"`
	Notes        string         `json:"notes"`
	PositionX    float64        `json:"positionX,omitempty" gorm:"type:double precision"`
	PositionY    float64        `json:"positionY,omitempty" gorm:"type:double precision"`
	Config       datatypes.JSON `json:"config,omitempty" gorm:"type:jsonb;not null;default:'{}'::jsonb"`
	CreatedAt    time.Time      `json:"-" gorm:"not null;default:now()"`
	UpdatedAt    time.Time      `json:"-" gorm:"not null;default:now()"`
}

// Edge connects two nodes within the same automation.
type Edge struct {
	ID           string `json:"id" gorm:"type:uuid;primaryKey"`
	AutomationID string `json:"-" gorm:"type:uuid;not null;index:idx_edges_automation"`

	FromNodeID string `json:"-" gorm:"type:uuid;not null;index"`
	FromPort   string `json:"-" gorm:"type:text"` // optional
	ToNodeID   string `json:"-" gorm:"type:uuid;not null;index"`

	// Helpful composite index to quickly fetch outgoing edges from a node+port.
	// (GORM will create single-column indexes from the tags above; for composite,
	// you can also add a manual index in a migration if desired.)
	// gorm:"index:idx_edges_from,priority:1" etc. kept simple here.

	CreatedAt time.Time `json:"-" gorm:"not null;default:now()"`
	UpdatedAt time.Time `json:"-" gorm:"not null;default:now()"`
}

// Optional: Enforce node/edge FK integrity at the DB level via separate migrations.
// With GORM, you can add constraints like this in AutoMigrate hooks or raw SQL:
//
//   ALTER TABLE edges
//     ADD CONSTRAINT fk_edges_fromnode
//       FOREIGN KEY ("from_node_id") REFERENCES nodes("id") ON DELETE CASCADE,
//     ADD CONSTRAINT fk_edges_tonode
//       FOREIGN KEY ("to_node_id")   REFERENCES nodes("id") ON DELETE CASCADE;
//
// Keeping them in SQL avoids circular migration issues.

// ---------- Public Graph DTO (API shape you already have) ----------

type Graph struct {
	Nodes []APINode `json:"nodes"`
	Edges []APIEdge `json:"edges"`
	Entry []string  `json:"entry"`
}

const DefaultNodeConfigEdge = "default"

type NodeConfig map[string]map[string]interface{}

func (c NodeConfig) hasEntries() bool {
	return c != nil && len(c) > 0
}

func (c NodeConfig) Clone() NodeConfig {
	if !c.hasEntries() {
		return nil
	}
	clone := make(NodeConfig, len(c))
	for edgeID, cfg := range c {
		if cfg == nil {
			clone[edgeID] = nil
			continue
		}
		copied := make(map[string]interface{}, len(cfg))
		for key, val := range cfg {
			copied[key] = val
		}
		clone[edgeID] = copied
	}
	return clone
}

func (c NodeConfig) EdgeConfig(edgeID string) map[string]interface{} {
	if !c.hasEntries() {
		return map[string]interface{}{}
	}
	if edgeID != "" {
		if cfg, ok := c[edgeID]; ok && cfg != nil {
			return cfg
		}
	}
	if cfg, ok := c[DefaultNodeConfigEdge]; ok && cfg != nil {
		return cfg
	}
	if cfg, ok := c[""]; ok && cfg != nil {
		return cfg
	}
	for _, cfg := range c {
		if cfg != nil {
			return cfg
		}
	}
	return map[string]interface{}{}
}

func (c *NodeConfig) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*c = nil
		return nil
	}
	type nodeConfigAlias map[string]map[string]interface{}
	var parsed nodeConfigAlias
	if err := json.Unmarshal(data, &parsed); err == nil {
		if len(parsed) == 0 {
			*c = nil
			return nil
		}
		*c = NodeConfig(parsed)
		return nil
	}
	var fallback map[string]interface{}
	if err := json.Unmarshal(data, &fallback); err != nil {
		return err
	}
	if len(fallback) == 0 {
		*c = nil
		return nil
	}
	*c = NodeConfig{
		DefaultNodeConfigEdge: fallback,
	}
	return nil
}

type APINode struct {
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	Kind      NodeKind   `json:"kind"`
	Name      string     `json:"name,omitempty"`
	Notes     string     `json:"notes,omitempty"`
	PositionX float64    `json:"positionX,omitempty"`
	PositionY float64    `json:"positionY,omitempty"`
	Config    NodeConfig `json:"config,omitempty"`
}

type APIEdge struct {
	ID         string `json:"id"`
	FromNodeId string `json:"fromNodeId"`
	FromPort   string `json:"fromPort,omitempty"`
	ToNodeId   string `json:"toNodeId"`
}
type AutomationRunStatus string

const (
	RunPending    AutomationRunStatus = "pending"
	RunRunning    AutomationRunStatus = "running"
	RunSuccess    AutomationRunStatus = "success"
	RunFailed     AutomationRunStatus = "failed"
	RunWithErrors AutomationRunStatus = "with_errors"
	RunCanceled   AutomationRunStatus = "canceled"
)

type AutomationRun struct {
	ID                 string  `gorm:"type:uuid;primaryKey"`
	AutomationID       string  `gorm:"type:uuid;not null;index"`
	BatchRunID         *string `json:"batchRunId,omitempty" gorm:"type:uuid;index"`
	LocationID         string  `gorm:"index"`
	TriggerType        string  `gorm:"not null"`
	TriggerPort        string
	TriggerPayload     map[string]interface{}   `json:"triggerPayload,omitempty" gorm:"-"`
	TriggerPayloadRaw  datatypes.JSON           `json:"-" gorm:"column:trigger_payload;type:jsonb"`
	ExecutionPath      []map[string]interface{} `json:"executionPath,omitempty" gorm:"-"`
	ExecutionPathRaw   datatypes.JSON           `json:"-" gorm:"column:execution_path;type:jsonb"` // ordered node IDs or [{nodeId,port}] to render the path quickly
	Status             AutomationRunStatus      `gorm:"type:text;not null"`
	ErrorMessage       string                   `gorm:"type:text"`
	RunNodesWithErrors int                      `gorm:"not null;default:0"`
	StartedAt          time.Time                `gorm:"not null"`
	CompletedAt        *time.Time
	RunNodes           []AutomationRunNode `gorm:"foreignKey:RunID"`
	CreatedAt          time.Time
}

type AutomationRunNode struct {
	ID                string `gorm:"type:uuid;primaryKey"`
	RunID             string `gorm:"type:uuid;index"`
	NodeID            string `gorm:"type:uuid;index"`
	NodeName          string
	NodeType          string
	Sequence          int                               `gorm:"index"` // to rebuild order
	InputFields       map[string]interface{}            `json:"inputFields,omitempty" gorm:"-"`
	InputFieldsRaw    datatypes.JSON                    `json:"-" gorm:"column:input_fields;type:jsonb"`
	OutputPayloads    map[string]map[string]interface{} `json:"outputPayloads,omitempty" gorm:"-"`
	OutputPayloadsRaw datatypes.JSON                    `json:"-" gorm:"column:output_payloads;type:jsonb"` // map of port name to payload
	ErrorMessage      string
	Status            AutomationRunStatus
	StartedAt         time.Time
	CompletedAt       *time.Time
}

func (r *AutomationRun) wrapJSONFields() error {
	if r == nil {
		return nil
	}

	if r.TriggerPayload != nil {
		raw, err := json.Marshal(r.TriggerPayload)
		if err != nil {
			return err
		}
		r.TriggerPayloadRaw = datatypes.JSON(raw)
	}

	if r.ExecutionPath != nil {
		raw, err := json.Marshal(r.ExecutionPath)
		if err != nil {
			return err
		}
		r.ExecutionPathRaw = datatypes.JSON(raw)
	}

	return nil
}

func (r *AutomationRun) unwrapJSONFields() error {
	if r == nil {
		return nil
	}

	if len(r.TriggerPayloadRaw) > 0 {
		var payload map[string]interface{}
		if err := json.Unmarshal(r.TriggerPayloadRaw, &payload); err != nil {
			return err
		}
		r.TriggerPayload = payload
	}

	if len(r.ExecutionPathRaw) > 0 {
		var path []map[string]interface{}
		if err := json.Unmarshal(r.ExecutionPathRaw, &path); err != nil {
			return err
		}
		r.ExecutionPath = path
	}

	return nil
}

func (r *AutomationRun) BeforeSave(tx *gorm.DB) (err error) {
	return r.wrapJSONFields()
}

func (r *AutomationRun) AfterFind(tx *gorm.DB) (err error) {
	return r.unwrapJSONFields()
}

func (n *AutomationRunNode) wrapJSONFields() error {
	if n == nil {
		return nil
	}

	if n.InputFields != nil {
		raw, err := json.Marshal(n.InputFields)
		if err != nil {
			return err
		}
		n.InputFieldsRaw = datatypes.JSON(raw)
	}

	if n.OutputPayloads != nil {
		raw, err := json.Marshal(n.OutputPayloads)
		if err != nil {
			return err
		}
		n.OutputPayloadsRaw = datatypes.JSON(raw)
	}

	return nil
}

func (n *AutomationRunNode) unwrapJSONFields() error {
	if n == nil {
		return nil
	}

	if len(n.InputFieldsRaw) > 0 {
		var fields map[string]interface{}
		if err := json.Unmarshal(n.InputFieldsRaw, &fields); err != nil {
			return err
		}
		n.InputFields = fields
	}

	if len(n.OutputPayloadsRaw) > 0 {
		var payloads map[string]map[string]interface{}
		if err := json.Unmarshal(n.OutputPayloadsRaw, &payloads); err != nil {
			return err
		}
		n.OutputPayloads = payloads
	}

	return nil
}

func (n *AutomationRunNode) BeforeSave(tx *gorm.DB) (err error) {
	return n.wrapJSONFields()
}

func (n *AutomationRunNode) AfterFind(tx *gorm.DB) (err error) {
	return n.unwrapJSONFields()
}

type BatchRunStatus string

const (
	BatchRunPending    BatchRunStatus = "pending"
	BatchRunRunning    BatchRunStatus = "running"
	BatchRunSuccess    BatchRunStatus = "success"
	BatchRunFailed     BatchRunStatus = "failed"
	BatchRunWithErrors BatchRunStatus = "with_errors"
	BatchRunCanceled   BatchRunStatus = "canceled"
)

type AutomationBatchRun struct {
	ID                string                 `gorm:"type:uuid;primaryKey"`
	Notes             string                 `gorm:"type:text"`
	AutomationID      string                 `gorm:"type:uuid;not null;index"`
	NodeID            string                 `gorm:"type:uuid;not null"`
	LocationID        string                 `gorm:"index"`
	Status            BatchRunStatus         `gorm:"type:text;not null"`
	ConfigSnapshot    map[string]interface{} `json:"configSnapshot,omitempty" gorm:"-"`
	ConfigSnapshotRaw datatypes.JSON         `json:"-" gorm:"column:config_snapshot;type:jsonb"`
	PageFrom          int                    `gorm:"not null"`
	PageTo            *int
	CurrentPage       int
	TotalPages        *int
	ItemsProcessed    int
	TotalItems        *int
	ProgressPct       float64 `gorm:"type:double precision"`
	ErrorMessage      string  `gorm:"type:text"`
	RunsWithErrors    int     `gorm:"not null;default:0"`
	StartedAt         *time.Time
	CompletedAt       *time.Time
	CreatorID         uint
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (b *AutomationBatchRun) wrapJSONFields() error {
	if b == nil {
		return nil
	}

	if b.ConfigSnapshot != nil {
		raw, err := json.Marshal(b.ConfigSnapshot)
		if err != nil {
			return err
		}
		b.ConfigSnapshotRaw = datatypes.JSON(raw)
	}

	return nil
}

func (b *AutomationBatchRun) unwrapJSONFields() error {
	if b == nil {
		return nil
	}

	if len(b.ConfigSnapshotRaw) > 0 {
		var cfg map[string]interface{}
		if err := json.Unmarshal(b.ConfigSnapshotRaw, &cfg); err != nil {
			return err
		}
		b.ConfigSnapshot = cfg
	}

	return nil
}

func (b *AutomationBatchRun) BeforeSave(tx *gorm.DB) (err error) {
	return b.wrapJSONFields()
}

func (b *AutomationBatchRun) AfterFind(tx *gorm.DB) (err error) {
	return b.unwrapJSONFields()
}

// ---------- Hooks / Helpers ----------

// AfterFind hydrates Automation.Graph from persisted Nodes/Edges/Entry.
func (a *Automation) AfterFind(tx *gorm.DB) (err error) {
	// Load children if not preloaded
	if len(a.Nodes) == 0 {
		if err = tx.Model(a).Association("Nodes").Find(&a.Nodes); err != nil {
			return err
		}
	}
	if len(a.Edges) == 0 {
		if err = tx.Model(a).Association("Edges").Find(&a.Edges); err != nil {
			return err
		}
	}

	// Build DTO nodes
	apiNodes := make([]APINode, 0, len(a.Nodes))
	for _, n := range a.Nodes {
		var cfg NodeConfig

		if len(n.Config) > 0 {

			_ = json.Unmarshal(n.Config, &cfg)
		}

		apiNodes = append(apiNodes, APINode{
			ID:        n.ID,
			Type:      n.Type,
			Kind:      n.Kind,
			Name:      n.Name,
			Notes:     n.Notes,
			PositionX: n.PositionX,
			PositionY: n.PositionY,
			Config:    cfg,
		})
	}

	// Build DTO edges
	apiEdges := make([]APIEdge, 0, len(a.Edges))
	for _, e := range a.Edges {
		apiEdges = append(apiEdges, APIEdge{
			ID:         e.ID,
			FromNodeId: e.FromNodeID,
			FromPort:   e.FromPort,
			ToNodeId:   e.ToNodeID,
		})
	}

	// Entry
	var entry []string
	if len(a.Entry) > 0 {
		_ = json.Unmarshal(a.Entry, &entry)
	}

	a.Graph = Graph{
		Nodes: apiNodes,
		Edges: apiEdges,
		Entry: entry,
	}
	return nil
}

// BeforeSave flattens Automation.Graph back into Node/Edge rows and JSONB fields.
func (a *Automation) BeforeSave(tx *gorm.DB) (err error) {
	// Entry []string -> JSONB
	if a.Graph.Entry != nil {
		if b, err := json.Marshal(a.Graph.Entry); err == nil {
			a.Entry = datatypes.JSON(b)
		} else {
			return err
		}
	}

	return nil
}

// AfterSave persists child Node and Edge rows after the Automation has been
// saved to the DB. Doing this in AfterSave ensures the Automation row exists
// and avoids FK violations when nodes reference the parent automation ID.
func (a *Automation) AfterSave(tx *gorm.DB) (err error) {
	// Replace nodes
	if err := tx.Where("automation_id = ?", a.ID).Delete(&Node{}).Error; err != nil {
		return err
	}

	// Nodes DTO -> rows
	if a.Graph.Nodes != nil {
		nodes := make([]Node, 0, len(a.Graph.Nodes))
		for _, n := range a.Graph.Nodes {
			var cfgBytes []byte
			if n.Config != nil {
				if b, err := json.Marshal(n.Config); err == nil {
					cfgBytes = b
				} else {
					return err
				}
			} else {
				cfgBytes = []byte("{}")
			}
			nodes = append(nodes, Node{
				ID:           n.ID,
				AutomationID: a.ID,
				Type:         n.Type,
				Kind:         n.Kind,
				Name:         n.Name,
				Notes:        n.Notes,
				PositionX:    n.PositionX,
				PositionY:    n.PositionY,
				Config:       datatypes.JSON(cfgBytes),
			})
		}
		a.Nodes = nodes
	}
	if len(a.Nodes) > 0 {
		if err := tx.Create(&a.Nodes).Error; err != nil {
			return err
		}
	}

	// Replace edges
	if err := tx.Where("automation_id = ?", a.ID).Delete(&Edge{}).Error; err != nil {
		return err
	}

	// Edges DTO -> rows
	if a.Graph.Edges != nil {
		edges := make([]Edge, 0, len(a.Graph.Edges))
		for _, e := range a.Graph.Edges {
			edges = append(edges, Edge{
				ID:           e.ID,
				AutomationID: a.ID,
				FromNodeID:   e.FromNodeId,
				FromPort:     e.FromPort,
				ToNodeID:     e.ToNodeId,
			})
		}
		a.Edges = edges
	}

	if len(a.Edges) > 0 {
		if err := tx.Create(&a.Edges).Error; err != nil {
			return err
		}
	}

	return nil
}

// ---------- Convenience: AutoMigrate ----------

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Automation{},
		&Node{},
		&Edge{},
		&AutomationBatchRun{},
		&AutomationRun{},
		&AutomationRunNode{},
	)
}
