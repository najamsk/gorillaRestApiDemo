package data

import "fmt"

type Repo struct {
	// teams   []Team
	// members []Member
	db *DataStore
}
type DataStore struct {
	Teams   []Team
	Members []Member
}

func NewDataStore() *DataStore {
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

	return &DataStore{Teams: teams, Members: members}
}

func NewRepo(db *DataStore) *Repo {
	// teams := []Team{
	// 	{1, "warwicks", nil},
	// 	{2, "maveriks", nil},
	// }
	// members := []Member{
	// 	{1, "najam awan", "nsa@gmail.com", 1},
	// 	{2, "najaf awan", "najaf@gmail.com", 2},
	// }
	// teams[0].Leader = &members[0]
	// teams[1].Leader = &members[1]

	return &Repo{db: db}
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
	return r.db.Teams
}

func (r *Repo) GetMembers() []Member {
	fmt.Println("data.members:", r.db.Members)
	return r.db.Members
}

func (r *Repo) CreateMember(m Member) Member {
	m.Id = len(r.db.Members) + 1
	fmt.Println(m.Id)
	r.db.Members = append(r.db.Members, m)
	return m
}
func (r *Repo) CreateTeam(t Team) Team {
	r.db.Teams = append(r.db.Teams, t)
	t.Id = len(r.db.Teams)
	return t
}

func init() {
	fmt.Println("setup data")
}
