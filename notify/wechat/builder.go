package wechat

import "strings"

type Builder interface {
	Build() Receivers
	SetToUser(toUser []string) Builder
	SetToParty(toParty string) Builder
	SetToTag(toTag string) Builder
	SetAgentId(agentId int) Builder
}

type Receivers struct {
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty"`
	ToTag   string `json:"totag"`
	AgentId int    `json:"agentid"`
}

type receiversBuilder struct {
	receivers Receivers
}

func NewReceiversBuilder() Builder {
	return &receiversBuilder{
		receivers: Receivers{},
	}
}

func (b *receiversBuilder) Build() Receivers {
	return b.receivers
}

func (b *receiversBuilder) SetToUser(toUsers []string) Builder {
	if len(toUsers) > 0 {
		b.receivers.ToUser = strings.Join(toUsers, "|")
	} else {
		b.receivers.ToUser = ""
	}
	return b
}

func (b *receiversBuilder) SetToParty(toParty string) Builder {
	b.receivers.ToParty = toParty
	return b
}

func (b *receiversBuilder) SetToTag(toTag string) Builder {
	b.receivers.ToTag = toTag
	return b
}

func (b *receiversBuilder) SetAgentId(agentId int) Builder {
	b.receivers.AgentId = agentId
	return b
}
