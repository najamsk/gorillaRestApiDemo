// Package classification My Sample API New
//
// Documentation for Sample API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package main

import "gorilla/internal/data"

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body data.GenericError
}

// No content is returned by this API endpoint
// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// A list of teams
// swagger:response teamsResponse
type teamsResponseWrapper struct {
	// All current teams
	// in: body
	Body []data.Team
}

// A list of members
// swagger:response membersResponse
type membersResponseWrapper struct {
	// All current members
	// in: body
	Body []data.Member
}

// A member
// swagger:response memberResponse
type memberResponseWrapper struct {
	// Member Details
	// in: body
	Body data.Member
}

// swagger:parameters createMember updateMember
type memberParamsWrapper struct {
	// Member data structure to Create or Update.
	// in: body
	// required: true
	Body data.Member
}

// swagger:parameters delMember
type memberDeleteParamsWrapper struct {
	// ID of memeber
	// in: path
	ID string `json:"memid"`
}
