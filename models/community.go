package models

import "time"

type Community struct {
	ID   string `json:"id" db:"community_id"`
	Name string `json:"name" db:"community_name"`
}
type CommunityDetail struct {
	ID           string    `json:"id" db:"community_id"`
	Name         string    `json:"name" db:"community_name"`
	Introduction string    `json:"introduction,omitempty" db:"introduction"`
	CreateTime   time.Time `json:"time" db:"create_time"`
}
