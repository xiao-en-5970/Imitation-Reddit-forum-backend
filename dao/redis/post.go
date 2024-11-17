package redis

import (
	"bluebell/models"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	//按分数从大到小查询
	return client.ZRevRange(key, start, end).Result()
}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	return getIDsFromKey(key, p.Page, p.Size)
}

func GetPostVoteDataByID(id string) (data int64, err error) {

	key := getRedisKey(KeyPostVotedZSetPF + id)
	data, err = client.ZCount(key, "1", "1").Result()
	if err != nil {
		zap.L().Error("GetPostVoteDataByID", zap.String("key", key), zap.Error(err))
		return
	}
	return
}

// 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPF + id)
	//	// 查找key中分数是1的元素的数量->统计每篇帖子的赞成票的数量
	//	v := client.ZCount(key, "1", "1").Val()
	//	data = append(data,v)
	//}

	//使用pipeline一次发送多条命令，减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(ids))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}

	return
}

// 按社区查询ids
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//使用zinterstore 把分区的帖子set与帖子分数的zset 生成一个新的zset
	//针对新的zset 按之前的逻辑取数据
	//社区key
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}

	communityKey := strconv.Itoa(int(p.CommunityID))
	//利用缓存key减少zinterstore执行次数
	key := orderKey + communityKey
	if client.Exists(key).Val() < 1 {
		//不存在则需要计算
		pipeline := client.Pipeline()
		zap.L().Debug("key", zap.String("key", key))
		zap.L().Debug("communityKey", zap.String("communityKey", communityKey))
		zap.L().Debug("orderKey", zap.String("orderKey", orderKey))
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "Max",
		}, getRedisKey(KeyCommunitySetPF+communityKey), orderKey)
		pipeline.Expire(key, 120*time.Second)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	return getIDsFromKey(key, p.Page, p.Size)

}
