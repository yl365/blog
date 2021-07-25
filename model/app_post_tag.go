package model

//PostTag 文章标签
type PostTag struct {
	Id     int   `xorm:"INT(11) PK AUTOINCR comment('主键')" json:"id"` //主键
	PostId int   `xorm:"INT(11) UNIQUE(post_tag)" json:"post_id"`     //文章Id
	TagId  int   `xorm:"INT(11) UNIQUE(post_tag)" json:"tag_id"`      //标签Id
	Post   *Post `xorm:"-" swaggerignore:"true" json:"post"`          //文章
	Tag    *Tag  `xorm:"-" swaggerignore:"true" json:"tag"`           //标签
}

// PostTags 文章对应的标签
func PostTags(pid int) ([]PostTag, error) {
	mods := make([]PostTag, 0, 4)
	err := Db.Where("post_id = ? ", pid).Find(&mods)
	if err == nil {
		ids := make([]int, 0, len(mods))
		for idx := range mods {
			if !inOf(mods[idx].TagId, ids) {
				ids = append(ids, mods[idx].TagId)
			}
		}
		mapTag := TagIds(ids)
		if mapTag != nil {
			for idx := range mods {
				mods[idx].Tag = mapTag[mods[idx].TagId]
			}
		}
	}
	return mods, nil
}

// TagPostCount 通过标签查询文章分页总数
func TagPostCount(tid int) int {
	var count int
	Db.SQL(`SELECT count(post_id) as count FROM post_tag WHERE tag_id=?`, tid).Get(&count)
	return count
}

// TagPostPage 通过标签查询文章分页
func TagPostPage(tagId, pi, ps int) ([]PostTag, error) {
	mods := make([]PostTag, 0, ps)
	err := Db.SQL(`SELECT id,post_id,tag_id FROM post_tag WHERE tag_id=? limit ?,? `, tagId, (pi-1)*ps, ps).Find(&mods)
	if len(mods) > 0 {
		ids := make([]int, 0, len(mods))
		for idx := range mods {
			if !inOf(mods[idx].PostId, ids) {
				ids = append(ids, mods[idx].PostId)
			}
		}
		mapSet := PostIds(ids)
		for idx := range mods {
			mods[idx].Post = mapSet[mods[idx].PostId]
		}
	}
	return mods, err
}

// TagPostAdds 添加文章标签[]
func TagPostAdds(mod *[]PostTag) bool {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	affect, _ := sess.Insert(mod)
	if affect < 1 {
		sess.Rollback()
		return false
	}
	sess.Commit()
	return true
}

// TagPostDrop 删除标签对应的标签_文章
func TagPostDrop(tid int) bool {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if affect, err := sess.Where("tag_id = ?", tid).Delete(&PostTag{}); affect > 0 && err == nil {
		sess.Commit()
		Db.ClearCache(new(PostTag)) //清空缓存
		return true
	}
	sess.Rollback()
	return false
}

// PostTagDrops 删除文章对应的标签_文章 修改的时候 变更
func PostTagDrops(pid int, tids []int) bool {
	if len(tids) < 1 {
		return true
	}
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if affect, err := sess.Where("post_id = ?", pid).In("tag_id", tids).Delete(&PostTag{}); affect > 0 && err == nil {
		sess.Commit()
		Db.ClearCache(new(PostTag)) //清空缓存
		return true
	}
	sess.Rollback()
	return false
}

// PostTagDrop 删除文章对应的标签_文章删除文章的时候
func PostTagDrop(pid int) bool {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if affect, err := sess.Where("post_id = ?", pid).Delete(&PostTag{}); affect > 0 && err == nil {
		sess.Commit()
		Db.ClearCache(new(PostTag)) //清空缓存
		return true
	}
	sess.Rollback()
	return false
}
