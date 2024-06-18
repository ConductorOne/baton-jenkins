package connector

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

var (
	resourceTypeUser = &v2.ResourceType{
		Id:          "user",
		DisplayName: "User",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_USER,
		},
	}
	resourceTypeNode = &v2.ResourceType{
		Id:          "node",
		DisplayName: "Node",
		// Traits: []v2.ResourceType_Trait{
		// 	v2.ResourceType_TRAIT_GROUP,
		// },
	}
	resourceTypeJob = &v2.ResourceType{
		Id:          "job",
		DisplayName: "Job",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_GROUP,
		},
	}
)
