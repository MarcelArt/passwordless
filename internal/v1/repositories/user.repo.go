package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MarcelArt/passwordless/internal/common"
	"github.com/MarcelArt/passwordless/internal/entities"
	"github.com/MarcelArt/passwordless/internal/v1/models"
	"github.com/gin-gonic/gin"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type IUserRepo interface {
	common.IBaseCrudRepo[entities.User, models.UserInput, models.UserPage]
	GetByUsernameOrEmail(c context.Context, usernameOrEmail string) (entities.User, error)
	GetPermissions(id any) ([]string, error)
	GetRoles(id any) ([]models.UserRole, error)
}

type UserRepo struct {
	db        *gorm.DB
	pageQuery string
}

var _ IUserRepo = &UserRepo{}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
		pageQuery: `
			SELECT 
				u.id as id, 
				u.username as username, 
				u.email as email,
				json_agg(r."name") as roles
			FROM users u  
			left join user_roles ur on u.id = ur.user_id and ur.deleted_at isnull
			left join roles r on ur.role_id = r.id  and r.deleted_at isnull
			where u.deleted_at isnull
			group by
				u.id,
				u.username,
				u.email
		`,
	}
}

func (r *UserRepo) Create(c context.Context, input models.UserInput) (uint, error) {
	user, err := common.Cast[entities.User](input)
	if err != nil {
		return 0, fmt.Errorf("cannot cast input: %w", err)
	}
	user.Password = input.Password

	err = gorm.G[entities.User](r.db).Create(c, &user)

	return user.ID, err
}

func (r *UserRepo) Read(c *gin.Context) (paginate.Page, []models.UserPage) {
	users := make([]models.UserPage, 0)

	stmt := r.db.Raw(r.pageQuery)

	pg := paginate.New()

	page := pg.With(stmt).Request(c.Request).Response(&users)

	return page, users
}

func (r *UserRepo) Update(c context.Context, id any, input models.UserInput) error {
	user, err := common.Cast[entities.User](input)
	if err != nil {
		return fmt.Errorf("cannot cast input: %w", err)
	}
	user.Password = input.Password

	_, err = gorm.G[entities.User](r.db).Where("id = ?", id).Updates(c, user)

	return err
}

func (r *UserRepo) Delete(c context.Context, id any) error {
	_, err := gorm.G[entities.User](r.db).Where("id = ?", id).Delete(c)

	return err
}

func (r *UserRepo) GetByID(c context.Context, id any) (entities.User, error) {
	var user entities.User

	user, err := gorm.G[entities.User](r.db).Where("id = ?", id).First(c)

	return user, err
}

func (r *UserRepo) GetByUsernameOrEmail(c context.Context, usernameOrEmail string) (entities.User, error) {
	return gorm.G[entities.User](r.db).Where("username = $1 or email = $1", usernameOrEmail).First(c)
}

func (r *UserRepo) GetPermissions(id any) ([]string, error) {
	var permissionsJSON string
	permissions := make([]string, 0)

	query := `
		SELECT 
			jsonb_agg(DISTINCT t.permissions) as permissions
		FROM (
			SELECT jsonb_array_elements_text(r.permissions ) AS permissions
			FROM roles r 
			left join user_roles ur on r.id = ur.role_id 
			where r.deleted_at isnull
			and ur.user_id = ?
		) t;
	`

	if err := r.db.Raw(query, id).Scan(&permissionsJSON).Error; err != nil {
		return nil, fmt.Errorf("failed retrieving permissions: %w", err)
	}

	err := json.Unmarshal([]byte(permissionsJSON), &permissions)

	return permissions, err
}

func (r *UserRepo) GetRoles(id any) ([]models.UserRole, error) {
	roles := make([]models.UserRole, 0)

	query := `
		select 
			r.id as id,
			r."name" as name,
			r.description as description
		from user_roles ur 
		join roles r on ur.role_id = r.id and r.deleted_at isnull
		where ur.deleted_at isnull
		and ur.user_id = ?
	`

	err := r.db.Raw(query, id).Scan(&roles).Error

	return roles, err
}
