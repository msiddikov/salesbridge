package automator

import "context"

func (rt *automationRuntime) nextNodes(ctx context.Context, parent *queuedNode, payloads map[string]map[string]interface{}) []*queuedNode {
	var nextNodes []*queuedNode
	if parent == nil {
		return nextNodes
	}

	for port := range payloads {
		nextID, edgeID := rt.next(parent.node.ID, port)
		if nextID == "" {
			continue
		}
		nextNode, ok := rt.nodes[nextID]
		if !ok {
			continue
		}
		child := &queuedNode{
			node:           nextNode,
			incomingEdgeID: edgeID,
			incomingPort:   port,
			parent:         parent,
		}
		nextNodes = append(nextNodes, child)
	}

	return nextNodes
}
