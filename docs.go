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
