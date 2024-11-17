package mysql

import (
	"bluebell/models"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"
	if err = db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}
func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlstr := `select 
    		community_id, community_name,introduction,create_time 
			from community 
			where community_id = ?`
	if err = db.Get(community, sqlstr, id); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no id in db")
			err = ErrorInvalidID
		}
	}
	zap.L().Debug(fmt.Sprintf("community:%v", community))
	return
}
