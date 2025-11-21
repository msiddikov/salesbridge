package automator

import (
	"client-runaway-zenoti/internal/db/models"
	"context"
)

type (
	Catalog struct {
		Meta       CatalogMeta
		Categories []Category
	}
	CatalogMeta struct {
		Name      string
		Version   string
		Publisher string
	}

	Category struct {
		Id    string
		Name  string
		Icon  string
		Color string
		Nodes []Node
	}

	Node struct {
		Id            string
		ExecFunc      func(context.Context, map[string]interface{}, models.Location) map[string]map[string]interface{} `json:"-"`
		CollectorFunc func(context.Context, map[string]interface{}, models.Location) (collectionResult, error)         `json:"-"`
		Title         string
		Description   string
		Type          NodeType
		Kind          string
		Icon          string
		Color         string
		Fields        []NodeField
		Ports         []NodePort
	}

	NodePort struct {
		Name    string
		Payload []NodeField
	}

	NodeField struct {
		Key           string
		Label         string
		Type          string
		SelectOptions []string
		Required      bool
	}

	NodeType string
)

const (
	// NodeTypes
	NodeTypeTrigger    NodeType = "trigger"
	NodeTypeAction     NodeType = "action"
	NodeTypeCollection NodeType = "collection"
)

var (
	attributionPersonNodeFields = []NodeField{
		{Key: "personId", Type: "string"},
		{Key: "first_name", Type: "string"},
		{Key: "last_name", Type: "string"},
		{Key: "email", Type: "string"},
		{Key: "phone", Type: "string"},
	}

	attributionStageHitNodeFields = []NodeField{
		{Key: "personId", Type: "string"},
		{Key: "first_name", Type: "string"},
		{Key: "last_name", Type: "string"},
		{Key: "email", Type: "string"},
		{Key: "phone", Type: "string"},
	}
)
