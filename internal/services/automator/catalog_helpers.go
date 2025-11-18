package automator

import (
	"sync"
)

var (
	catalogLookupOnce sync.Once
	catalogLookup     map[string]Node
)

func getCatalogNode(nodeType string) (Node, bool) {
	catalogLookupOnce.Do(func() {
		catalogLookup = make(map[string]Node)
		for _, category := range catalogFull.Categories {
			for _, node := range category.Nodes {
				catalogLookup[node.Id] = node
			}
		}
	})

	node, ok := catalogLookup[nodeType]
	return node, ok
}

func errorPayload(err error, message string) map[string]map[string]interface{} {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	return map[string]map[string]interface{}{
		"error": {
			"message": message,
			"error":   errorMsg,
		},
	}
}

func successPayload(payload map[string]interface{}) map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"success": payload,
	}
}

func customPayload(port string, payload map[string]interface{}) map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		port: payload,
	}
}

var errorPort = NodePort{
	Name: "error",
	Payload: []NodeField{
		{Key: "message", Type: "string"},
		{Key: "error", Type: "string"},
	},
}

func successPort(payload []NodeField) NodePort {
	return NodePort{
		Name:    "success",
		Payload: payload,
	}
}

func customPort(portName string, payload []NodeField) NodePort {
	return NodePort{
		Name:    portName,
		Payload: payload,
	}
}

func getCatalogWithImplementedNodes(cat Catalog) Catalog {
	res := Catalog{
		Meta: cat.Meta,
	}
	for _, category := range cat.Categories {
		newCategory := Category{
			Id:   category.Id,
			Name: category.Name,
		}
		for _, node := range category.Nodes {
			if node.ExecFunc != nil || // for all nodes the exec func must be implemented
				node.Type == NodeTypeTrigger || // for trigger nodes the exec func is not required
				(node.Type == NodeTypeCollection && node.CollectorFunc != nil) { // for collection nodes the collector func must be implemented
				newCategory.Nodes = append(newCategory.Nodes, node)
			}
		}
		if len(newCategory.Nodes) > 0 {
			res.Categories = append(res.Categories, newCategory)
		}
	}
	return res

}
