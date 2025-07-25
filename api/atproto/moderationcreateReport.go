// Code generated by cmd/lexgen (see Makefile's lexgen); DO NOT EDIT.

package atproto

// schema: com.atproto.moderation.createReport

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bluesky-social/indigo/lex/util"
)

// ModerationCreateReport_Input is the input argument to a com.atproto.moderation.createReport call.
type ModerationCreateReport_Input struct {
	ModTool *ModerationCreateReport_ModTool `json:"modTool,omitempty" cborgen:"modTool,omitempty"`
	// reason: Additional context about the content and violation.
	Reason *string `json:"reason,omitempty" cborgen:"reason,omitempty"`
	// reasonType: Indicates the broad category of violation the report is for.
	ReasonType *string                               `json:"reasonType" cborgen:"reasonType"`
	Subject    *ModerationCreateReport_Input_Subject `json:"subject" cborgen:"subject"`
}

type ModerationCreateReport_Input_Subject struct {
	AdminDefs_RepoRef *AdminDefs_RepoRef
	RepoStrongRef     *RepoStrongRef
}

func (t *ModerationCreateReport_Input_Subject) MarshalJSON() ([]byte, error) {
	if t.AdminDefs_RepoRef != nil {
		t.AdminDefs_RepoRef.LexiconTypeID = "com.atproto.admin.defs#repoRef"
		return json.Marshal(t.AdminDefs_RepoRef)
	}
	if t.RepoStrongRef != nil {
		t.RepoStrongRef.LexiconTypeID = "com.atproto.repo.strongRef"
		return json.Marshal(t.RepoStrongRef)
	}
	return nil, fmt.Errorf("cannot marshal empty enum")
}
func (t *ModerationCreateReport_Input_Subject) UnmarshalJSON(b []byte) error {
	typ, err := util.TypeExtract(b)
	if err != nil {
		return err
	}

	switch typ {
	case "com.atproto.admin.defs#repoRef":
		t.AdminDefs_RepoRef = new(AdminDefs_RepoRef)
		return json.Unmarshal(b, t.AdminDefs_RepoRef)
	case "com.atproto.repo.strongRef":
		t.RepoStrongRef = new(RepoStrongRef)
		return json.Unmarshal(b, t.RepoStrongRef)

	default:
		return nil
	}
}

// ModerationCreateReport_ModTool is a "modTool" in the com.atproto.moderation.createReport schema.
//
// Moderation tool information for tracing the source of the action
type ModerationCreateReport_ModTool struct {
	// meta: Additional arbitrary metadata about the source
	Meta *interface{} `json:"meta,omitempty" cborgen:"meta,omitempty"`
	// name: Name/identifier of the source (e.g., 'bsky-app/android', 'bsky-web/chrome')
	Name string `json:"name" cborgen:"name"`
}

// ModerationCreateReport_Output is the output of a com.atproto.moderation.createReport call.
type ModerationCreateReport_Output struct {
	CreatedAt  string                                 `json:"createdAt" cborgen:"createdAt"`
	Id         int64                                  `json:"id" cborgen:"id"`
	Reason     *string                                `json:"reason,omitempty" cborgen:"reason,omitempty"`
	ReasonType *string                                `json:"reasonType" cborgen:"reasonType"`
	ReportedBy string                                 `json:"reportedBy" cborgen:"reportedBy"`
	Subject    *ModerationCreateReport_Output_Subject `json:"subject" cborgen:"subject"`
}

type ModerationCreateReport_Output_Subject struct {
	AdminDefs_RepoRef *AdminDefs_RepoRef
	RepoStrongRef     *RepoStrongRef
}

func (t *ModerationCreateReport_Output_Subject) MarshalJSON() ([]byte, error) {
	if t.AdminDefs_RepoRef != nil {
		t.AdminDefs_RepoRef.LexiconTypeID = "com.atproto.admin.defs#repoRef"
		return json.Marshal(t.AdminDefs_RepoRef)
	}
	if t.RepoStrongRef != nil {
		t.RepoStrongRef.LexiconTypeID = "com.atproto.repo.strongRef"
		return json.Marshal(t.RepoStrongRef)
	}
	return nil, fmt.Errorf("cannot marshal empty enum")
}
func (t *ModerationCreateReport_Output_Subject) UnmarshalJSON(b []byte) error {
	typ, err := util.TypeExtract(b)
	if err != nil {
		return err
	}

	switch typ {
	case "com.atproto.admin.defs#repoRef":
		t.AdminDefs_RepoRef = new(AdminDefs_RepoRef)
		return json.Unmarshal(b, t.AdminDefs_RepoRef)
	case "com.atproto.repo.strongRef":
		t.RepoStrongRef = new(RepoStrongRef)
		return json.Unmarshal(b, t.RepoStrongRef)

	default:
		return nil
	}
}

// ModerationCreateReport calls the XRPC method "com.atproto.moderation.createReport".
func ModerationCreateReport(ctx context.Context, c util.LexClient, input *ModerationCreateReport_Input) (*ModerationCreateReport_Output, error) {
	var out ModerationCreateReport_Output
	if err := c.LexDo(ctx, util.Procedure, "application/json", "com.atproto.moderation.createReport", nil, input, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
