package model

import "gorm.io/gorm"

// 用户点赞表
// idx_userid_videoid: 查找用户点赞列表，查找用户是否给某视频点了赞
type Favorite struct {
	gorm.Model
	UserID  int64 `gorm:"not null; index:idx_userid_videoid" json:"user_id"`
	VideoID int64 `gorm:"not null; index:idx_userid_videoid" json:"video_id"`
}

// AddFavorite 添加点赞信息 返回值 0-点赞成功 1-已经点赞 2-数据库错误
func AddFavorite(userID int64, videoID int64) int {
	if IsFavorite(userID, videoID) {
		return 1
	}
	favorite := &Favorite{
		UserID:  userID,
		VideoID: videoID,
	}
	// 在favorite表中添加该记录
	if db.Create(favorite).Error != nil {
		return 2
	}
	// 更新video表中的点赞数
	if db.Model(&Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count + 1")).Error != nil {
		return 2
	}
	return 0
}

// DeleteFavorite 删除点赞信息 返回值 0-删除成功 1-没有点赞 2-数据库错误
func DeleteFavorite(userID int64, videoID int64) int {
	if !IsFavorite(userID, videoID) {
		return 1
	}
	// 在favorite表中删除该记录
	if db.Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&Favorite{}).Error != nil {
		return 2
	}
	// 更新video表中的点赞数
	if db.Model(&Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count - 1")).Error != nil {
		return 2
	}
	return 0
}

// GetFavoriteVideoList 获取用户的点赞视频
func GetFavoriteVideoList(userID int64) ([]*VideoAuthorUnion, error) {
	var favorites []*VideoAuthorUnion
	sql := "SELECT v.ID AS id, u.id AS author_id, u.username AS author_name, v.play_url AS play_url, v.cover_url AS cover_url, v.favorite_count AS favorite_count FROM users u LEFT JOIN favorites f ON u.id = f.user_id LEFT JOIN videos v  ON f.video_id = v.id WHERE user_id = ?"
	err := db.Raw(sql, userID).Scan(&favorites).Error

	if err != nil {
		return nil, err
	}
	return favorites, nil
}

// IsFavorite 判断是否已经点赞
func IsFavorite(userID int64, videoID int64) bool {
	return db.First(&Favorite{}, "user_id = ? and video_id = ?", userID, videoID).Error == nil
}
