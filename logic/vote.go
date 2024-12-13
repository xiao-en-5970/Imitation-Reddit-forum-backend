package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"go.uber.org/zap"
	"strconv"
)

// 1 vote + 432 score
// 86400/200 -> 200
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug(
		"VoteForPost",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction),
	)
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
