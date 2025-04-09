package model

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// User 表示系统中的用户实体
type User struct {
	// 用户唯一标识符，自动生成且只能创建时设置
	ID int64 `gorm:"<-:create;primaryKey" json:"id"`
	// 用户名，唯一且只能创建时设置
	UserName string `gorm:"<-:create;unique;type:varchar(32)" json:"user_name" validate:"required,min=3,max=32"`
	// 用户密码，建议存储加密后的哈希值
	Password string `gorm:"not null;type:varchar(255)" json:"password" validate:"required,min=6,max=32"`
	// 用户昵称，可为空
	NickName sql.NullString `gorm:"type:varchar(32)" json:"nick_name" validate:"required,min=3,max=32"`
	// 用户手机号，唯一且可为空
	Phone sql.NullString `gorm:"unique;type:varchar(20)" json:"phone" validate:"required,len=11"`
	// 用户年龄，默认为零
	Age int32 `gorm:"default:0;not null" json:"age"`
	// 用户性别，0表示未知，1表示男，2表示女
	Gender int32 `gorm:"default:0;not null" json:"gender"`
	// 电子邮箱，唯一且可为空
	Email sql.NullString `gorm:"unique;type:varchar(100)" json:"email" validate:"required,email"`
	// 用户头像URL
	Avatar sql.NullString `gorm:"type:varchar(255)" json:"avatar"`
	// 用户状态，true表示正常，false表示禁用
	Status bool `gorm:"default:true;not null" json:"status"`
	// 是否为管理员
	IsAdmin bool `gorm:"default:false;not null" json:"is_admin"`
	// 用户等级，默认为 300，level3
	Level int32 `gorm:"default:300;not null" json:"level"`
	// 记录创建时间
	CreatedAt time.Time `gorm:"<-:create" json:"created_at"`
	// 更新时间
	UpdatedAt time.Time `gorm:"autoCreateTime;autoUpdateTime" json:"updated_at"`
	// 软删除时间戳
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	// 用户的粉丝数量（忽略）
	FollowerCount int64 `gorm:"-" json:"follower_count"`
	// 用户的关注数量（忽略）
	FollowCount int64 `gorm:"-" json:"follow_count"`
	// 用户是否被关注（忽略）
	IsFollowed bool `gorm:"-" json:"is_followed"`
	// 用户的文章数量（忽略）
	ArticleCount int64 `gorm:"-" json:"article_count"`
}

// CacheKeyByID 根据用户ID生成缓存键
func (u *User) CacheKeyByID(id int64) string {
	return fmt.Sprintf("User:id:%d", id)
}

// CacheKeyByName 根据用户名生成缓存键
func (u *User) CacheKeyByName(name string) string {
	return fmt.Sprintf("User:name:%s", name)
}
