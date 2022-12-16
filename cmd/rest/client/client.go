package main

import (
	"gorilla/client"
	"gorilla/client/members"
	"gorilla/client/teams"

	"go.uber.org/zap"
)

type remoteClient struct {
	cfg    *client.TransportConfig
	client *client.GorillaAPI
	l      *zap.Logger
}

func NewRemoteClient(cfg *client.TransportConfig, client *client.GorillaAPI, logger *zap.Logger) *remoteClient {
	return &remoteClient{
		cfg:    cfg,
		client: client,
		l:      logger,
	}
}

// GetAllTeams makes remote clal and prints all the teams
func (r *remoteClient) GetAllTeams() {
	params := teams.NewListTeamsParams()
	t, err := r.client.Teams.ListTeams(params)
	if err != nil {
		r.l.Error("listing teams gets error:", zap.Error(err))
	}
	teams := t.GetPayload()
	for _, v := range teams {
		r.l.Info("team:",
			zap.Int64("id", v.ID),
			zap.String("name", v.Name),
			zap.String("leader name", v.Leader.Name),
		)
	}
}

// GetAllTeams makes remote clal and prints all the teams
func (r *remoteClient) GetAllMembers() {
	params := members.NewListMembersParams()
	t, err := r.client.Members.ListMembers(params)
	if err != nil {
		r.l.Error("listing teams gets error:", zap.Error(err))
	}
	rows := t.GetPayload()
	for _, v := range rows {
		r.l.Info("member:",
			zap.Int64("id", v.ID),
			zap.String("name", v.Name),
		)
	}
}

// GetAllTeams makes remote clal and prints all the teams
func (r *remoteClient) DeleteMember() {
	params := members.NewDelMemberParams()
	params.ID = "11"
	_, err := r.client.Members.DelMember(params)
	if err != nil {
		r.l.Error("delete member error:", zap.Error(err))
	}
	// r.l.Info("delete result code:", zap.Bool("code", t.IsCode(http.StatusNotFound)))

}

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// c := client.Default
	tc := client.DefaultTransportConfig().WithHost("localhost:8000")
	cl := client.NewHTTPClientWithConfig(nil, tc)

	rc := NewRemoteClient(tc, cl, logger)

	// rc.GetAllTeams()
	// rc.GetAllMembers()
	rc.DeleteMember()
}
