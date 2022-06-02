package model

import "gorm.io/gorm"

// 用户点赞表
// idx_userid_videoid: 查找用户点赞列表，查找用户是否给某视频点了赞
type Favorite struct {
	gorm.Model
	UserID  int64 `gorm:"not null; index:idx_userid_videoid" json:"user_id"`
	VideoID int64 `gorm:"not null; index:idx_userid_videoid" json:"video_id"`
	Video   Video `gorm:"foreignkey:VideoID;association_foreignkey:ID"`
	User    User  `gorm:"foreignkey:UserID;association_foreignkey:ID"`
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
	if DB.Create(favorite).Error != nil {
		return 2
	}
	// 更新video表中的点赞数
	// if DB.Model(&Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count + 1")).Error != nil {
	// 	return 2
	// }
	return 0
}

// DeleteFavorite 删除点赞信息 返回值 0-删除成功 1-没有点赞 2-数据库错误
func DeleteFavorite(userID int64, videoID int64) int {
	if !IsFavorite(userID, videoID) {
		return 1
	}
	// 在favorite表中删除该记录
	if DB.Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&Favorite{}).Error != nil {
		return 2
	}
	// 更新video表中的点赞数
	// if DB.Model(&Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count - 1")).Error != nil {
	// 	return 2
	// }
	return 0
}

// GetFavoriteVideoList 获取用户的点赞视频
func GetFavoriteVideoList(userID int64) ([]*VideoAuthorUnion, error) {
	var favorites []*VideoAuthorUnion
	sql := `
			SELECT 
				t1.*, 
				t2.favorite_count 
			FROM 
				(
				SELECT 
					tv.id AS id,
					tu.id AS author_id, 
					tu.username AS author_name, 
					tv.play_url AS play_url, 
					tv.cover_url AS cover_url, 
					tv.title AS title 
				FROM 
					favorites tf 
					LEFT JOIN videos tv ON tf.video_id = tv.id 
					LEFT JOIN users tu ON tv.author_id = tu.id 
				WHERE 
					user_id = ?
				) t1 
				LEFT JOIN (
				SELECT 
					video_id, 
					count(*) AS favorite_count 
				FROM 
					favorites 
				GROUP BY 
					video_id
				) t2 ON t1.id = t2.video_id;
			`
	err := DB.Raw(sql, userID).Scan(&favorites).Error

	if err != nil {
		return nil, err
	}
	return favorites, nil
}

// IsFavorite 判断是否已经点赞
func IsFavorite(userID int64, videoID int64) bool {
	return DB.First(&Favorite{}, "user_id = ? and video_id = ?", userID, videoID).Error == nil
}
