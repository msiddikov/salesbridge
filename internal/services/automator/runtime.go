package automator

import "context"

func (rt *automationRuntime) nextNodes(ctx context.Context, parent *queuedNode, payloads map[string]map[string]interface{}) []*queuedNode {
	var nextNodes []*queuedNode
	if parent == nil {
		return nextNodes
	}

	for port := range payloads {
		nextTargets := rt.next(parent.node.ID, port)
		for _, target := range nextTargets {
			nextNode, ok := rt.nodes[target.nodeID]
			if !ok {
				continue
			}
			child := &queuedNode{
				node:           nextNode,
				incomingEdgeID: target.edgeID,
				incomingPort:   port,
				parent:         parent,
			}
			nextNodes = append(nextNodes, child)
		}
	}

	return nextNodes
}
