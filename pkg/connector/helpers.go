package connector

import (
	"fmt"
	"strings"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	MembershipEntitlementIDTemplate = "membership:%s"
	V1GrantIDTemplate               = "grant:%s:%s"
)

func annotationsForUserResourceType() annotations.Annotations {
	annos := annotations.Annotations{}
	annos.Update(&v2.SkipEntitlementsAndGrants{})
	return annos
}

// Populate entitlement options for a 1password resource.
func PopulateOptions(displayName, permission, resource string) []ent.EntitlementOption {
	options := []ent.EntitlementOption{
		ent.WithGrantableTo(resourceTypeUser),
		ent.WithDescription(fmt.Sprintf("%s of Jenkins %s %s", permission, displayName, resource)),
		ent.WithDisplayName(fmt.Sprintf("%s %s %s", displayName, resource, permission)),
	}
	return options
}

func titleCase(s string) string {
	titleCaser := cases.Title(language.English)

	return titleCaser.String(s)
}

func ParseEntitlementID(id string) (*v2.ResourceId, []string, error) {
	parts := strings.Split(id, ":")
	// Need to be at least 3 parts type:entitlement_id:slug
	if len(parts) < 3 || len(parts) > 3 {
		return nil, nil, fmt.Errorf("jenkins-connector: invalid resource id")
	}

	resourceId := &v2.ResourceId{
		ResourceType: parts[0],
		Resource:     strings.Join(parts[1:len(parts)-1], ":"),
	}

	return resourceId, parts, nil
}

func ParseGrantID(id string) (*v2.ResourceId, []string, error) {
	parts := strings.Split(id, ":")
	// Need to be at least 5 parts type:grant_id:slug:resource_id:resource_type
	if len(parts) < 5 || len(parts) > 5 {
		return nil, nil, fmt.Errorf("jenkins-connector: invalid resource id")
	}

	resourceId := &v2.ResourceId{
		ResourceType: parts[0],
		Resource:     strings.Join(parts[1:len(parts)-1], ":"),
	}

	return resourceId, parts, nil
}

func V1MembershipEntitlementID(resourceID string) string {
	return fmt.Sprintf(MembershipEntitlementIDTemplate, resourceID)
}

func V1GrantID(entitlementID string, userID string) string {
	return fmt.Sprintf(V1GrantIDTemplate, entitlementID, userID)
}
