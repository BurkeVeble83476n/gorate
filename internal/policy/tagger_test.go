package policy

import (
	"testing"
)

func makeTaggerPolicy(name, tags string) Policy {
	return Policy{Name: name, Tags: tags}
}

func TestAddTag_NewTag(t *testing.T) {
	p := makeTaggerPolicy("p1", "")
	AddTag(&p, "critical")
	if p.Tags != "critical" {
		t.Errorf("expected 'critical', got %q", p.Tags)
	}
}

func TestAddTag_DuplicateTag(t *testing.T) {
	p := makeTaggerPolicy("p1", "critical")
	AddTag(&p, "critical")
	if p.Tags != "critical" {
		t.Errorf("expected no duplicate, got %q", p.Tags)
	}
}

func TestAddTag_MultipleUnique(t *testing.T) {
	p := makeTaggerPolicy("p1", "")
	AddTag(&p, "critical")
	AddTag(&p, "prod")
	if p.Tags != "critical,prod" {
		t.Errorf("expected 'critical,prod', got %q", p.Tags)
	}
}

func TestRemoveTag_Existing(t *testing.T) {
	p := makeTaggerPolicy("p1", "critical,prod")
	RemoveTag(&p, "critical")
	if p.Tags != "prod" {
		t.Errorf("expected 'prod', got %q", p.Tags)
	}
}

func TestRemoveTag_NonExisting(t *testing.T) {
	p := makeTaggerPolicy("p1", "prod")
	RemoveTag(&p, "critical")
	if p.Tags != "prod" {
		t.Errorf("expected 'prod' unchanged, got %q", p.Tags)
	}
}

func TestGetTags_Empty(t *testing.T) {
	p := makeTaggerPolicy("p1", "")
	tags := GetTags(&p)
	if len(tags) != 0 {
		t.Errorf("expected no tags, got %v", tags)
	}
}

func TestGetTags_Multiple(t *testing.T) {
	p := makeTaggerPolicy("p1", "a,b,c")
	tags := GetTags(&p)
	if len(tags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(tags))
	}
}

func TestHasTag_True(t *testing.T) {
	p := makeTaggerPolicy("p1", "alpha,beta")
	if !HasTag(&p, "beta") {
		t.Error("expected HasTag to return true for 'beta'")
	}
}

func TestHasTag_False(t *testing.T) {
	p := makeTaggerPolicy("p1", "alpha")
	if HasTag(&p, "gamma") {
		t.Error("expected HasTag to return false for 'gamma'")
	}
}

func TestFilterByTag_ReturnsMatching(t *testing.T) {
	policies := []Policy{
		makeTaggerPolicy("p1", "prod,critical"),
		makeTaggerPolicy("p2", "dev"),
		makeTaggerPolicy("p3", "prod"),
	}
	result := FilterByTag(policies, "prod")
	if len(result) != 2 {
		t.Errorf("expected 2 policies, got %d", len(result))
	}
}

func TestFilterByTag_NoMatch(t *testing.T) {
	policies := []Policy{
		makeTaggerPolicy("p1", "dev"),
	}
	result := FilterByTag(policies, "prod")
	if len(result) != 0 {
		t.Errorf("expected 0 policies, got %d", len(result))
	}
}
