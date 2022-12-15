package data

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

const name = "handlers"

type Error string

func (e Error) Error() string {
	return string(e)
}

const ErrNotFound = Error("Not found")

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

func (r *Repo) GetTeams(ctx context.Context) []Team {
	_, span := otel.Tracer(name).Start(ctx, "repo/GetTeams")
	defer span.End()
	span.SetAttributes(attribute.String("request.n", "getTeamsAttr"))
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

func (r *Repo) UpdateMember(m Member) (*Member, error) {
	mem, err := r.FindMember(m.Id)
	if err != nil {
		return nil, err
	}
	r.db.Members[mem] = m
	return &m, nil
}

func (r *Repo) DeleteMember(m int) error {
	k, err := r.FindMember(m)
	if err != nil {
		return err
	}
	//if k is first item
	if k == 0 {
		r.db.Members = r.db.Members[1:]

	} else if k == len(r.db.Members)-1 {
		r.db.Members = r.db.Members[:len(r.db.Members)]
		//last item
	} else {
		left := r.db.Members[:k]
		right := r.db.Members[k+1:]
		final := append(left, right...)
		r.db.Members = final
	}
	return nil
}

func (r *Repo) FindMember(x int) (int, error) {
	for k, v := range r.db.Members {
		if v.Id == x {
			return k, nil
		}
	}
	return -1, ErrNotFound
}

func (r *Repo) CreateTeam(t Team) Team {
	r.db.Teams = append(r.db.Teams, t)
	t.Id = len(r.db.Teams)
	return t
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// func init() {
// 	fmt.Println("setup data")
// }
