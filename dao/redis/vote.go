package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"strconv"
	"time"
)

const (
	oneWeekInSeconds = 7 * 24 * 60 * 60
	scorePerVote     = 432
)

var (
	ErrVoteTimeExpire = errors.New("vote time expire")
	ErrVoteRepeated   = errors.New("vote repeated")
)

func CreatePost(postID, communityID int64) (err error) {

	pipeline := client.TxPipeline()
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)
	_, err = pipeline.Exec()
	if err != nil {
		zap.L().Error("pipeline exec error", zap.Error(err))
	}
	return
}

func VoteForPost(userID, postID string, value float64) error {
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()

	zap.L().Debug("time.now", zap.Float64("now", float64(time.Now().Unix())))
	zap.L().Debug("postime", zap.Float64("post", postTime))
	zap.L().Debug("time.now-posttime", zap.Float64("now-post", float64(time.Now().Unix())-postTime))
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	zap.L().Debug("redis.ZScore", zap.Float64("ov", ov))
	zap.L().Debug("value", zap.Float64("value", value))
	diff := value - ov
	if diff == 0 {
		return ErrVoteRepeated
	}
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), diff*scorePerVote, postID)

	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), userID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
