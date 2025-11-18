package automator

import (
	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

var (
	catalogFull = Catalog{
		Meta: templateMeta,
		Categories: []Category{
			ghlCategory,
			zenotiCategory,
			attributionCategory,
			othersCategory,
		},
	}

	templateMeta = CatalogMeta{
		Name:      "Salesbridge Automator",
		Version:   "1.0.0",
		Publisher: "Salesbridge",
	}

	// Attribution Category
	attributionCategory = Category{
		Id:   "attribution",
		Name: "Attribution",
		Nodes: []Node{
			attributionActionFindPerson,
			attributionActionCreatePerson,
			attributionActionUpdatePerson,
			attributionStageHit,
		},
	}

	// Actions
	attributionActionFindPerson = Node{
		Id:          "attribution.person.find",
		Title:       "Find Person",
		Description: "Finds a person in Attribution. with email or phone.",
		Type:        NodeTypeAction,
		Icon:        "ri:form",
		Color:       ColorAction,
		Ports: []NodePort{
			{
				Name:    "success",
				Payload: zenotiGuestNodeFields,
			},
			errorPort,
		},
		Fields: []NodeField{
			{Key: "email", Type: "string"},
			{Key: "phone", Type: "string"},
		},
	}
	attributionActionCreatePerson = Node{
		Id:          "attribution.person.create",
		Title:       "Create Person",
		Description: "Creates a new person in Attribution.",
		Type:        NodeTypeAction,
		Icon:        "ri:form",
		Color:       ColorAction,
		Ports: []NodePort{
			{
				Name:    "success",
				Payload: attributionPersonNodeFields,
			},
			errorPort,
		},
		Fields: attributionPersonNodeFields,
	}

	attributionActionUpdatePerson = Node{
		Id:          "attribution.person.update",
		Title:       "Update Person",
		Description: "Updates an existing person in Attribution.",
		Type:        NodeTypeAction,
		Icon:        "ri:form",
		Color:       ColorAction,
		Ports: []NodePort{
			{
				Name:    "success",
				Payload: attributionPersonNodeFields,
			},
			errorPort,
		},
		Fields: attributionPersonNodeFields,
	}

	attributionStageHit = Node{
		Id:          "attribution.stage.hit",
		Title:       "Stage Hit",
		Description: "Triggers when a person hits a specific stage in Attribution.",
		Type:        NodeTypeTrigger,
		Icon:        "ri:form",
		Color:       ColorAction,
		Ports: []NodePort{
			{
				Name:    "out",
				Payload: attributionStageHitNodeFields,
			},
		},
		Fields: attributionStageHitNodeFields,
	}
)

func GetCatalog(c *gin.Context) {
	c.Data(lvn.Res(200, getCatalogWithImplementedNodes(catalogFull), "OK"))
}
