package data

import "fmt"

type Repo struct {
	teams   []Team
	members []Member
}

func NewRepo() *Repo {
	teams := []Team{
		{1, "warwicks", nil},
		{2, "maveriks", nil},
	}
	members := []Member{
		{1, "najam awan", "nsa@gmail.com", 1},
		{2, "najaf awan", "najaf@gmail.com", 2},
	}
	teams[0].Leader = &members[0]
	teams[1].Leader = &members[1]

	return &Repo{teams: teams, members: members}
}

// Member defines the structure for an API product
type Member struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	TeamId int    `json:"teamID"`
}

// Team defines the structure for an API product
type Team struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Leader *Member `json:"leader"`
}

func (r *Repo) GetTeams() []Team {
	return r.teams
}

func (r *Repo) GetMembers() []Member {
	fmt.Println("data.members:", r.members)
	return r.members
}

func init() {
	fmt.Println("setup data")
}
