package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"go.uber.org/zap"
	"strconv"
)
import "bluebell/pkg/snow"

func CreatePost(p *models.Post) (err error) {
	p.ID = snow.GenID()
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	if err != nil {
		return err
	}
	return
}
func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {

	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}

	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
		return
	}
	// search community detail
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)

	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Int64("author_id", post.CommunityID), zap.Error(err))
		return
	}
	voteNum, err := redis.GetPostVoteDataByID(strconv.Itoa(int(pid)))
	if err != nil {
		zap.L().Error("GetPostVoteDataByID(post.CommunityID) failed", zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
		VoteNum:         voteNum,
	}

	return
}
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed", zap.Error(err))
		return
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}
		// search community detail
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Int64("author_id", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		zap.L().Error("redis.GetPostIDsInOrder() failed", zap.Error(err))
		return
	}
	//根据id去mysql数据库查询
	zap.L().Debug("ids:", zap.Any("ids", ids))
	if len(ids) == 0 {
		zap.L().Warn("len(ids)==0")
		return
	}
	posts, err := mysql.GetPostListByIDs(ids)
	zap.L().Debug("posts:", zap.Any("posts", posts))
	if err != nil {
		zap.L().Error("mysql.GetPostListByIDs() failed", zap.Error(err))
		return
	}
	// 提前查询每篇帖子的投票
	voteData, err := redis.GetPostVoteData(ids)
	zap.L().Debug("voteData", zap.Any("voteData", voteData))
	if err != nil {
		zap.L().Error("redis.GetPostVoteData() failed", zap.Error(err))
		return
	}
	// 将帖子的作者及分区信息查询出来填充到帖子中

	data = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}
		// search community detail
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Int64("author_id", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		zap.L().Error("redis.GetPostIDsInOrder() failed", zap.Error(err))
		return
	}
	//根据id去mysql数据库查询
	zap.L().Debug("ids:", zap.Any("ids", ids))
	if len(ids) == 0 {
		zap.L().Warn("len(ids)==0")
		return
	}
	posts, err := mysql.GetPostListByIDs(ids)
	zap.L().Debug("posts:", zap.Any("posts", posts))
	if err != nil {
		zap.L().Error("mysql.GetPostListByIDs() failed", zap.Error(err))
		return
	}
	// 提前查询每篇帖子的投票
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		zap.L().Error("redis.GetPostVoteData() failed", zap.Error(err))
		return
	}
	// 将帖子的作者及分区信息查询出来填充到帖子中

	data = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}
		// search community detail
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Int64("author_id", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// 将两个查询逻辑合二为一
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	zap.L().Debug("p", zap.Any("p", p))
	if p.CommunityID == 0 {
		data, err = GetPostList2(p)
	} else {
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew() failed", zap.Error(err))

	}
	return
}
