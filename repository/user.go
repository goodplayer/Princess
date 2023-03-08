package repository

type User struct {
	UserId string `gorm:"column:user_id;primaryKey"`

	UserName  string `gorm:"column:user_name;not null;unique"`
	AliasName string `gorm:"column:alias_name;not null"`
	Password  string `gorm:"column:password;not null"`

	UserStatus int `gorm:"column:user_status;not null;default:0"`

	CreatedAt int64 `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64 `gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt int64 `gorm:"column:deleted_at"`
}

func (u *User) TableName() string {
	return "princess_user"
}

func (d *Db) SaveUser(user *User) (*User, error) {
	result := d.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (d *Db) UpdateUser(user *User) (*User, error) {
	result := d.db.Save(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (d *Db) DeleteUser(user *User) error {
	result := d.db.Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *Db) CheckUserExistsByEmail(email string) (bool, error) {
	tx := d.db.Where("user_name = ?", email).Find(&User{})
	if tx.Error != nil {
		return false, tx.Error
	}
	return tx.RowsAffected > 0, nil
}
