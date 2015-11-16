package model

import (
	"container/list"
	"database/sql"
	"log"
	"time"
)

import (
	"moetang.info/prod/Princess/repo"
)

var (
	_postUtil postUtil
)

func init() {
	_postUtil = postUtil{}
}

type postUtil struct {
}

func PostUtil() postUtil {
	return _postUtil
}

func (postUtil) GetIndexPosts() ([]*Post, error) {
	// record limit is 10
	rows, err := repo.Run().Query(`
SELECT id, title, abstract, content, status, posttime, lastupdatetime, userid FROM post where status = 0 order by posttime desc limit 10;
	`)
	if err != nil {
		return nil, err
	}
	defer func() {
		e := rows.Close()
		if e != nil {
			log.Println(e)
		}
	}()
	l := list.New()
	for rows.Next() {
		post, err := WrapPost(rows)
		if err != nil {
			return nil, err
		}
		l.PushBack(post)
	}
	if l.Len() > 0 {
		result := make([]*Post, l.Len())
		idx := 0
		for e := l.Front(); e != nil; e = e.Next() {
			result[idx] = e.Value.(*Post)
			idx++
		}
		return result, nil
	} else {
		return []*Post{}, nil
	}
}

func (postUtil) GetPostById(id int64) (*Post, bool, error) {
	// return post, isExist, error
	post := NewPost()
	post.Id = id
	err := post.FillPostById()
	if err == NO_SUCH_RECORD {
		return nil, false, nil
	} else if err != nil {
		return nil, false, nil
	}
	return post, true, nil
}

type Post struct {
	Id             int64
	Title          string
	Abstract       string
	Content        string
	Status         int64
	PostTime       int64
	LastUpdateTime int64
	PostUser       *User
}

func NewPost() *Post {
	return new(Post)
}

func WrapPost(r *sql.Rows) (*Post, error) {
	post := new(Post)
	var userId int64 = -1
	e := r.Scan(&post.Id, &post.Title, &post.Abstract, &post.Content, &post.Status, &post.PostTime, &post.LastUpdateTime, &userId)
	if e != nil {
		return nil, e
	}
	if userId <= 0 {
		return nil, NO_USER_RELATED
	}
	user := NewUser()
	user.Id = userId
	err := user.GetUserById()
	if err != nil {
		return nil, err
	}
	post.PostUser = user
	return post, nil
}

func (this *Post) Save() error {
	this.PostTime = time.Now().UnixNano()
	this.LastUpdateTime = this.PostTime
	_, err := repo.Run().Exec(`
INSERT INTO post(title, abstract, content, status, posttime, lastupdatetime, userid)
VALUES ($1, $2, $3, $4, $5, $6, $7);
	`, this.Title, this.Abstract, this.Content, this.Status, this.PostTime, this.LastUpdateTime, this.PostUser.Id)
	return err
}

func (this *Post) FillPostById() error {
	var userId int64 = -1
	r := repo.Run().QueryRow(`
SELECT title, abstract, content, status, posttime, lastupdatetime, userid FROM post where id = $1;
	`, this.Id)
	e := r.Scan(&this.Title, &this.Abstract, &this.Content, &this.Status, &this.PostTime, &this.LastUpdateTime, &userId)
	if e != nil {
		if e == sql.ErrNoRows {
			return NO_SUCH_RECORD
		} else {
			return e
		}
	}
	if userId <= 0 {
		return NO_USER_RELATED
	}
	user := NewUser()
	user.Id = userId
	err := user.GetUserById()
	if err != nil {
		return err
	}
	this.PostUser = user
	return nil
}
