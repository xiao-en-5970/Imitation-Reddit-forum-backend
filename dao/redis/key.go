package redis

const (
	Prefix             = "bluebell:"   //项目key前缀
	KeyPostTimeZSet    = "post:time"   //zset;帖子及发帖时间
	KeyPostScoreZSet   = "post:score"  //zset;帖子及投票分数
	KeyPostVotedZSetPF = "post:voted:" //zset;记录用户及投票类型;参数是post id
	KeyCommunitySetPF  = "community:"  // set;保存每个分区下帖子的id
)

func getRedisKey(key string) string {
	return Prefix + key
}
