package automator

const (
	ColorDefault    = "#6366F1"
	ColorTrigger    = "#8910b9ff"
	ColorAction     = "#3B82F6"
	ColorCollection = "#f6ae5cff"
	ColorOther      = "#0bd2f5ff"

	iconTrigger    = "ri:rocket-line"
	iconAction     = "ri:play-circle-line"
	iconCollection = "ri:stack-line"
	iconOther      = "ri:settings-4-line"
)

func designCatalog(c Catalog) Catalog {

	for _, category := range c.Categories {
		for i, node := range category.Nodes {
			// Set default icon based on node type if not set
			if node.Color == "" {
				switch node.Type {
				case NodeTypeTrigger:
					node.Color = ColorTrigger
				case NodeTypeAction:
					node.Color = ColorAction
				case NodeTypeCollection:
					node.Color = ColorCollection
				default:
					node.Color = ColorOther
				}
			}
			if node.Icon == "" {
				switch node.Type {
				case NodeTypeTrigger:
					node.Icon = iconTrigger
				case NodeTypeAction:
					node.Icon = iconAction
				case NodeTypeCollection:
					node.Icon = iconCollection
				default:
					node.Icon = iconOther
				}
			}
			category.Nodes[i] = node
		}
	}

	return c
}
