package mysql

import (
	"bluebell/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"
)

// 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlstr := `insert into post(
                 post_id,title,content,author_id,community_id)
				values(?,?,?,?,?)`
	_, err = db.Exec(sqlstr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	if err != nil {
		return
	}
	return
}

// GetPostById 根据id查询单个帖子数据
func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlstr := `select 
    	post_id,title,content,author_id,community_id,create_time 
		from post 
		where post_id = ?`
	err = db.Get(post, sqlstr, pid)
	return
}

func GetUserById(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlstr := `select user_id ,user.username from user where user_id = ?`
	err = db.Get(user, sqlstr, uid)
	return
}

// 查询用户列表函数
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlstr := `select 
    	post_id,title,content,author_id,community_id,create_time 
		from post 
		order by create_time
		DESC
		limit ?,?`
	posts = make([]*models.Post, 0, size)
	zap.L().Debug(fmt.Sprintf("page:", page, "size:", size))
	err = db.Select(&posts, sqlstr, (page-1)*size, size)
	if err != nil {
		zap.L().Error("GetPostList failed", zap.Error(err))
		return
	}
	return
}

// 根据给定打id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlstr := `select post_id,title,content,author_id,community_id,create_time
				from post
				where post_id in (?)
				order by FIND_IN_SET(post_id,?)`
	query, args, err := sqlx.In(sqlstr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}
