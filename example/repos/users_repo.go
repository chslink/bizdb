// Code generated by bizdb. DO NOT EDIT.
package repos

import (
	"errors"
	"fmt"
	"time"

	"github.com/chslink/bizdb"
	"github.com/chslink/bizdb/example/models"
)

var (
	ErrUsersNotFound = errors.New("users not found")
)

type UsersRepo struct {
	memDB     *bizdb.MemoryDB
	tableName string
	tx        *bizdb.Transaction
}

// NewUsersRepo 创建普通Repository
func NewUsersRepo(memDB *bizdb.MemoryDB) *UsersRepo {
	return &UsersRepo{
		memDB:     memDB,
		tableName: "users",
	}
}

// WithTx 创建带事务的Repository
func (r *UsersRepo) WithTx(tx *bizdb.Transaction) *UsersRepo {
	return &UsersRepo{
		memDB:     r.memDB,
		tableName: r.tableName,
		tx:        tx,
	}
}

type UsersQuery struct {
	Id *int64 `db:"id"`
	Name *string `db:"name"`
	Email *string `db:"email"`
	Age *int64 `db:"age"`
	CreateAt *time.Time `db:"create_at"`
	Limit  *int
	Offset *int
}

// GetById 根据主键获取
func (r *UsersRepo) GetById(Id int64) (*models.Users, error) {
	val, err := r.memDB.Get(r.tx, r.tableName, Id)
	if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, ErrUsersNotFound
	}

	model, ok := val.(*models.Users)
	if !ok {
		return nil, fmt.Errorf("type assertion failed")
	}

	return model.Copy(), nil
}


// GetByEmail 根据唯一索引查询
func (r *UsersRepo) GetByEmail(val string) (*models.Users, error) {
	query := UsersQuery{
		Email: &val,
	}

	results, err := r.Query(query)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrUsersNotFound
	}

	return results[0], nil
}


// Create 创建记录
func (r *UsersRepo) Create(model *models.Users) error {
	if r.tx == nil {
		return errors.New("write operation requires a transaction")
	}

	if existing, _ := r.memDB.Get(r.tx, r.tableName, model.Id); existing != nil {
		return fmt.Errorf("record with Id %v already exists", model.Id)
	}

	return r.memDB.Put(r.tx, r.tableName, model.Id, model.Copy())
}

// Update 更新记录
func (r *UsersRepo) Update(model *models.Users) error {
	if r.tx == nil {
		return errors.New("write operation requires a transaction")
	}

	if existing, _ := r.memDB.Get(r.tx, r.tableName, model.Id); existing == nil {
		return ErrUsersNotFound
	}

	return r.memDB.Put(r.tx, r.tableName, model.Id, model.Copy())
}

// Delete 删除记录
func (r *UsersRepo) Delete(Id int64) error {
	if r.tx == nil {
		return errors.New("delete operation requires a transaction")
	}

	if existing, _ := r.memDB.Get(r.tx, r.tableName, Id); existing == nil {
		return ErrUsersNotFound
	}

	return r.memDB.Delete(r.tx, r.tableName, Id)
}

// Range 遍历查询
func (r *UsersRepo) Range(f func(Id int64, val *models.Users) bool) error {
	return r.memDB.Range(r.tx, r.tableName, func(id, val any) bool {
		return f(id.(int64), val.(*models.Users))
	})
}

// Query 高级查询
func (r *UsersRepo) Query(q UsersQuery) ([]*models.Users, error) {
	var results []*models.Users
	var end int

	// 分页处理
	if q.Limit != nil || q.Offset != nil {
		offset := 0
		if q.Offset != nil {
			offset = *q.Offset
		}
		limit := 0
		if q.Limit != nil {
			limit = *q.Limit
		}
		end = offset + limit
	}

	err := r.Range(func(Id int64, val *models.Users) bool {
		model := val
		match := true

		
		if q.Id != nil && *q.Id != model.Id {
			match = false
		}
		
		if q.Name != nil && *q.Name != model.Name {
			match = false
		}
		
		if q.Email != nil && *q.Email != model.Email {
			match = false
		}
		
		if q.Age != nil && *q.Age != model.Age {
			match = false
		}
		
		if q.CreateAt != nil && *q.CreateAt != model.CreateAt {
			match = false
		}
		

		if match {
			results = append(results, model.Copy())
		}

		if end > 0 && len(results) >= end {
			return false
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	// 应用分页
	if q.Limit != nil || q.Offset != nil {
		offset := 0
		if q.Offset != nil {
			offset = *q.Offset
		}
		limit := len(results)
		if q.Limit != nil {
			limit = *q.Limit
		}

		end := offset + limit
		if end > len(results) {
			end = len(results)
		}

		if offset > len(results) {
			results = []*models.Users{}
		} else {
			results = results[offset:end]
		}
	}

	return results, nil
}
