package connector

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

// To fetch 1000 results
// https://confluence.atlassian.com/bitbucketserverkb/how-to-apply-the-limit-filter-in-bitbucket-server-and-datacenter-rest-api-and-query-more-than-the-max-limit-of-1000-results-1142440445.html
const ITEMSPERPAGE = 1000

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
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_GROUP,
		},
	}
)
