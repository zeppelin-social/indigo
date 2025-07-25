// Code generated by cmd/lexgen (see Makefile's lexgen); DO NOT EDIT.

package bsky

// schema: app.bsky.unspecced.initAgeAssurance

import (
	"context"

	"github.com/bluesky-social/indigo/lex/util"
)

// UnspeccedInitAgeAssurance_Input is the input argument to a app.bsky.unspecced.initAgeAssurance call.
type UnspeccedInitAgeAssurance_Input struct {
	// countryCode: An ISO 3166-1 alpha-2 code of the user's location.
	CountryCode string `json:"countryCode" cborgen:"countryCode"`
	// email: The user's email address to receive assurance instructions.
	Email string `json:"email" cborgen:"email"`
	// language: The user's preferred language for communication during the assurance process.
	Language string `json:"language" cborgen:"language"`
}

// UnspeccedInitAgeAssurance calls the XRPC method "app.bsky.unspecced.initAgeAssurance".
func UnspeccedInitAgeAssurance(ctx context.Context, c util.LexClient, input *UnspeccedInitAgeAssurance_Input) (*UnspeccedDefs_AgeAssuranceState, error) {
	var out UnspeccedDefs_AgeAssuranceState
	if err := c.LexDo(ctx, util.Procedure, "application/json", "app.bsky.unspecced.initAgeAssurance", nil, input, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
