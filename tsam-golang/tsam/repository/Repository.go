package repository

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
)

// Repository defines all methods to be present in repository.
type Repository interface {
	Get(uow *UnitOfWork, id uuid.UUID, out interface{}, queryProcessor ...QueryProcessor) error
	GetAll(uow *UnitOfWork, out interface{}, queryProcessor ...QueryProcessor) error
	GetRecord(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) error
	GetAllInOrder(uow *UnitOfWork, out, orderBy interface{}, queryProcessor ...QueryProcessor) error

	GetCount(uow *UnitOfWork, out, count interface{}, queryProcessors ...QueryProcessor) error
	GetCountUnscoped(uow *UnitOfWork, out, count interface{}, queryProcessors ...QueryProcessor) error

	// For tenants.
	GetForTenant(uow *UnitOfWork, tenantID, id uuid.UUID, out interface{}, queryProcessor ...QueryProcessor) error
	GetAllForTenant(uow *UnitOfWork, tenantID uuid.UUID, out interface{}, queryProcessor ...QueryProcessor) error
	GetRecordForTenant(uow *UnitOfWork, tenantID uuid.UUID, out interface{}, queryProcessors ...QueryProcessor) error
	GetAllInOrderForTenant(uow *UnitOfWork, tenantID uuid.UUID, out, orderBy interface{}, queryProcessor ...QueryProcessor) error

	GetCountForTenant(uow *UnitOfWork, tenantID uuid.UUID, out, count interface{}, queryProcessors ...QueryProcessor) error
	GetCountUnscopedForTenant(uow *UnitOfWork, tenantID uuid.UUID, out, count interface{}, queryProcessors ...QueryProcessor) error

	// Other CRUD operations.
	Add(uow *UnitOfWork, out interface{}) error
	Update(uow *UnitOfWork, out interface{}) error
	UpdateWithMap(uow *UnitOfWork, value, out interface{}, queryProcessors ...QueryProcessor) error
	BatchUpdate(uow *UnitOfWork, value, condition, out interface{}) error

	Save(uow *UnitOfWork, value interface{}) error
	NewDelete(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) error
	Delete(uow *UnitOfWork, out interface{}, where ...interface{}) error
	DeleteForTenant(uow *UnitOfWork, tenantID uuid.UUID, out interface{}, where ...interface{}) error

	RemoveAssociations(uow *UnitOfWork, out interface{}, associationName string, associations ...interface{}) error
	ReplaceAssociations(uow *UnitOfWork, out interface{}, associationName string, associations ...interface{}) error

	Scan(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) error

	// other
	NewPluckColumn(uow *UnitOfWork, columnName string, out interface{}, queryProcessors ...QueryProcessor) error
	SubQuery(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) (*gorm.SqlExpr, error)
}

// NewDelete deletes a record from table.
func (repository *GormRepository) NewDelete(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().Delete(out).Error
}

// BatchUpdate updates a group of entities.
func (repository *GormRepository) BatchUpdate(uow *UnitOfWork, value, condition, out interface{}) error {
	return uow.DB.Model(value).Where(condition).Updates(out).Error
}

// // GetColumn :- This is a beta function which might replace PluckColumn
// func (repository *GormRepository) GetColumn(uow *UnitOfWork, tablename string, columnNames []string,
// 	condition, out interface{}, queryProcessors ...QueryProcessor) error {
// 	db := uow.DB
// 	var err error
// 	if db, err = executeQueryProcessors(db, out, queryProcessors...); err != nil {
// 		return err
// 	}
// 	return db.Debug().Table(tablename).Select(columnNames).Where(condition).Find(out).Error
// }

// PluckColumn plucks column from Table
func PluckColumn(db *gorm.DB, tableName string, columnName string, out interface{}, queryProcessors ...QueryProcessor) error {
	var err error
	if db, err = executeQueryProcessors(db, out, queryProcessors...); err != nil {
		return err
	}
	return db.Debug().Table(tableName).Pluck(columnName, out).Error
}

// NewPluckColumn plucks column from Table
func (repository *GormRepository) NewPluckColumn(uow *UnitOfWork, columnName string, out interface{}, queryProcessors ...QueryProcessor) error {
	var err error
	db := uow.DB
	if db, err = executeQueryProcessors(db, out, queryProcessors...); err != nil {
		return err
	}
	return db.Debug().Pluck(columnName, out).Error
}

// GormRepository will implement repository interface.
type GormRepository struct{}

//NewGormRepository returns new instance of GormRepository.
func NewGormRepository() *GormRepository {
	return &GormRepository{}
}

// Scan will fill the out interface with data(fields) based on the given QP conditions.
func (repository *GormRepository) Scan(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Scan(out).Error
}

// GetRecord returns a specific record from table with the given filter.
func (repository *GormRepository) GetRecord(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().First(out).Error
}

// GetRecordForTenant returns a specific record of the tenant from table with the given filter.
func (repository *GormRepository) GetRecordForTenant(uow *UnitOfWork, tenantID uuid.UUID, out interface{}, queryProcessors ...QueryProcessor) error {
	// #tenantID should be the first element in slice if "where" is appeneded in QP.
	queryProcessors = append([]QueryProcessor{Filter("tenant_id = ?", tenantID)}, queryProcessors...)
	return repository.GetRecord(uow, out, queryProcessors...)
}

// Get returns record from table by ID.
func (repository *GormRepository) Get(uow *UnitOfWork, id uuid.UUID, out interface{}, queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().First(out, "id = ?", id).Error
}

// GetForTenant returns record from table by ID of the speicfic tenant.
func (repository *GormRepository) GetForTenant(uow *UnitOfWork, tenantID, id uuid.UUID, out interface{}, queryProcessors ...QueryProcessor) error {
	// #tenantID should be the first element in slice if "where" is appeneded in QP.
	queryProcessors = append([]QueryProcessor{Filter("tenant_id = ?", tenantID)}, queryProcessors...)
	return repository.Get(uow, id, out, queryProcessors...)
}

// GetAll returns all records from the table.
func (repository *GormRepository) GetAll(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().Find(out).Error
}

// GetAllForTenant returns all records from the table for the specific tenant.
func (repository *GormRepository) GetAllForTenant(uow *UnitOfWork, tenantID uuid.UUID, out interface{}, queryProcessors ...QueryProcessor) error {
	// #tenantID should be the first element in slice if "where" is appeneded in QP.
	queryProcessors = append([]QueryProcessor{Filter("tenant_id = ?", tenantID)}, queryProcessors...)
	return repository.GetAll(uow, out, queryProcessors...)
}

// GetAllInOrder returns all records from table in specified order.
func (repository *GormRepository) GetAllInOrder(uow *UnitOfWork, out, orderBy interface{}, queryProcessors ...QueryProcessor) error {
	db := uow.DB

	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().Order(orderBy).Find(out).Error
}

// GetAllInOrderForTenant returns all records from table in specified order for the specific tenant.
func (repository *GormRepository) GetAllInOrderForTenant(uow *UnitOfWork, tenantID uuid.UUID, out, orderBy interface{}, queryProcessors ...QueryProcessor) error {
	// #tenantID should be the first element in slice if "where" is appeneded in QP.
	queryProcessors = append([]QueryProcessor{Filter("tenant_id = ?", tenantID)}, queryProcessors...)
	return repository.GetAllInOrder(uow, out, orderBy, queryProcessors...)
}

// Add adds record to table.
func (repository *GormRepository) Add(uow *UnitOfWork, out interface{}) error {
	return uow.DB.Create(out).Error
}

// Update updates the record in table.
func (repository *GormRepository) Update(uow *UnitOfWork, out interface{}) error {
	return uow.DB.Model(out).Update(out).Error
}

// Save updates the record in table. If value doesn't have primary key, new record will be inserted.
func (repository *GormRepository) Save(uow *UnitOfWork, value interface{}) error {
	return uow.DB.Save(value).Error
}

// UpdateWithMap updates the record in table using map.
// 	UpdateWithMap(uow, user{id="101"},map[string]interface{}{"name":"Ramesh"}
// It will filter by ID only if value has a primary key.
// 	Query: UPDATE users WHERE `id`="101" SET `name`="Ramesh";
func (repository *GormRepository) UpdateWithMap(uow *UnitOfWork, value, out interface{},
	queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().Model(value).Update(out).Error
}

// Delete deletes a record from table.
func (repository *GormRepository) Delete(uow *UnitOfWork, out interface{}, where ...interface{}) error {
	return uow.DB.Delete(out, where...).Error
}

// DeleteForTenant deletes a record from table for the specific tenant.
func (repository *GormRepository) DeleteForTenant(uow *UnitOfWork, tenantID uuid.UUID, out interface{}, where ...interface{}) error {
	return uow.DB.Where("tenant_id = ?", tenantID).Delete(out, where...).Error
}

// GetCount gives number of records in database.
func (repository *GormRepository) GetCount(uow *UnitOfWork, out, count interface{}, queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().Model(out).Count(count).Error
}

// GetCountForTenant gives number of records in database for the specific tenant.
func (repository *GormRepository) GetCountForTenant(uow *UnitOfWork, tenantID uuid.UUID, out, count interface{}, queryProcessors ...QueryProcessor) error {
	// #tenantID should be the first element in slice if "where" is appeneded in QP.
	queryProcessors = append([]QueryProcessor{Filter("tenant_id = ?", tenantID)}, queryProcessors...)
	return repository.GetCount(uow, out, count, queryProcessors...)
}

// ReplaceAssociations replaces associations from the given entity.
func (repository *GormRepository) ReplaceAssociations(uow *UnitOfWork, out interface{}, associationName string, associations ...interface{}) error {
	if err := uow.DB.Model(out).Association(associationName).Replace(associations...).Error; err != nil {
		return err
	}
	return nil
}

// RemoveAssociations removes associations from the given entity.
func (repository *GormRepository) RemoveAssociations(uow *UnitOfWork, out interface{}, associationName string, associations ...interface{}) error {
	if err := uow.DB.Model(out).Association(associationName).Delete(associations...).Error; err != nil {
		return err
	}
	return nil
}

// GetCountUnscoped (INCLUDES SOFT DELETED RECORDS) returns total count of specified entity after applying filters, if any.
func (repository *GormRepository) GetCountUnscoped(uow *UnitOfWork, out, count interface{}, queryProcessors ...QueryProcessor) error {
	db := uow.DB.Unscoped()
	var err error
	if db, err = executeQueryProcessors(db, out, queryProcessors...); err != nil {
		return err
	}
	return db.Model(out).Count(count).Error
}

// GetCountUnscopedForTenant (INCLUDES SOFT DELETED RECORDS) returns total count of specified entity
// for the specific tenant after applying filters, if any.
func (repository *GormRepository) GetCountUnscopedForTenant(uow *UnitOfWork, tenantID uuid.UUID, out, count interface{}, queryProcessors ...QueryProcessor) error {
	// #tenantID should be the first element in slice if "where" is appeneded in QP.
	queryProcessors = append([]QueryProcessor{Filter("tenant_id = ?", tenantID)}, queryProcessors...)
	return repository.GetCountUnscoped(uow, out, count, queryProcessors...)
}

// SubQuery returns query as a sub query.
func (repository *GormRepository) SubQuery(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) (*gorm.SqlExpr, error) {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return nil, err
	}
	return db.Debug().SubQuery(), nil
}

// ******************************** All GormRepository methods above this line ********************************

// OrderBy specifies order when retrieving records from database, set reorder to `true` to overwrite defined conditions
// 	Order("name DESC")
// 	Order("name DESC", true) // reorder
// 	Order(gorm.Expr("name = ? DESC", "first")) // sql expression
func OrderBy(value interface{}, reorder ...bool) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Order(value, reorder...)
		return db, nil
	}
}

// Select specify fields that you want to retrieve from database when querying, by default, will select all fields;
// When creating/updating, specify fields that you want to save to database.
func Select(query interface{}, args ...interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Select(query, args...)
		return db, nil
	}
}

// UnScoped will set the db to unscoped which will not fire the "where deleted_at IS NULL" automatically
// or other automatically generated conditions & queries.
func UnScoped() QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Unscoped()
		return db, nil
	}
}

// Raw allows you create a raw SQL query
func Raw(db *gorm.DB, sql string, out interface{}, values ...interface{}) error {
	return db.Debug().Raw(sql, values...).Scan(out).Error
}

// RawQuery uses raw sql as conditions, won't run it unless invoked by other methods
// 	Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&result)
func RawQuery(sql string, values ...interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Debug().Raw(sql, values...)
		return db, db.Error
	}
}

// Join specifies join conditions as query processors. (Use Find() or something similar to get results)
// 	Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "tsam@example.org")
func Join(query string, args ...interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Joins(query, args...)
		return db, nil
	}
}

// Model specifies the model you would like to run db operations on
// 	// update all users's name to `hello`
// 	db.Model(&User{}).Update("name", "hello")
// 	// if user's primary key is non-blank, will use it as condition, then will only update the user's name to `hello`
// 	db.Model(&user).Update("name", "hello")
func Model(value interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Debug().Model(value)
		return db, nil
	}
}

// GroupBy returns QueryProcessor (used to Group Result Set)
func GroupBy(groupstr ...string) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		for _, entity := range groupstr {
			db = db.Group(entity)
		}
		return db, nil
	}
}

// Having use with GroupBy
func Having(conditions string, values ...interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Having(conditions, values...)
		return db, nil
	}
}

// Distinct adds to query for distinct result set.
func Distinct(columns ...string) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		length := len(columns)
		if length == 0 {
			return db, nil
		}
		str := "DISTINCT "
		for index, column := range columns {
			if length-1 != index {
				str += column + ","
			} else {
				str += column
			}
		}
		db = db.Select(str)
		return db, nil
	}
}

// Table specifies the table you would like to run db operations
func Table(tableName string) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Table(tableName)
		return db, nil
	}
}

// PreloadAssociations preloads data from the specified table.
// 	PreloadAssociations([]string{"Orders", "Customers"})
func PreloadAssociations(preloadAssociations []string) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		if preloadAssociations != nil {
			for _, association := range preloadAssociations {
				db = db.Debug().Preload(association)
			}
		}
		return db, nil
	}
}

// PreloadWithCondition preloads associations with given conditions.
// 	PreloadWithCondition(map[string][]interface{}{"Orders":[]interface{}{"state NOT IN (?)","cancelled"}})
func PreloadWithCondition(preloadAssociations map[string][]interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		if preloadAssociations != nil {
			for table, condition := range preloadAssociations {
				db = db.Debug().Preload(table, condition...)
			}
		}
		return db, nil
	}
}

// PreloadWithCustomCondition preloads associations with queryProcessor.
// 	'Cant use maps as they dont maintain the order of queries'
func PreloadWithCustomCondition(preloadAssociations ...Preload) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		// closureIndex is a separate index maintained for looping inside the anonymous func.
		var closureIndex uint8
		for _, association := range preloadAssociations {
			db = db.Preload(association.Schema, func(db *gorm.DB) *gorm.DB {
				db, err := executeQueryProcessors(db, out, preloadAssociations[closureIndex].Queryprocessors...)
				if err != nil {
					db.Error = err
				}
				closureIndex++
				return db
			})
			if db.Error != nil {
				return db, db.Error
			}
		}
		return db, nil
	}
}

// executeQueryProcessors executes all queryProcessor func.
func executeQueryProcessors(db *gorm.DB, out interface{}, queryProcessors ...QueryProcessor) (*gorm.DB, error) {
	var err error
	for _, query := range queryProcessors {
		if query != nil {
			db, err = query(db, out)
			if err != nil {
				return db, err
			}
		}
	}
	return db, nil
}

// Filter will filter the results based on condition.
// 	Filter("name= ?","Ramesh")
// Query : WHERE `name`= "Ramesh"
func Filter(condition string, args ...interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Debug().Where(condition, args...)
		return db, nil
	}
}

// Limit sets limit and returns as query processor.
func Limit(limit interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Limit(limit)
		return db, nil
	}
}

// FilterWithOperator adds multiple condition with operator.
// FilterWithOperator("`sales_person_id`", "IS NULL", "AND", nil) ===> Pass nil in value for NULL checks
// 	FilterWithOperator([]string{"name","age"},[]string{"LIKE ?",">"},[]string{"AND"},[]interface{"ajay",18}])
// Query: `name` LIKE "ajay" AND `age` > 18
func FilterWithOperator(columnNames []string, conditions []string, operators []string, values []interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {

		if len(columnNames) != len(conditions) && len(conditions) != len(values) {
			return db, nil
		}

		if len(conditions) == 1 {
			if values[0] == nil {
				db = db.Where(fmt.Sprintf("%v %v", columnNames[0], conditions[0]))
				return db, nil
			}
			db = db.Where(fmt.Sprintf("%v %v", columnNames[0], conditions[0]), values[0])
			return db, nil
		}
		if len(columnNames)-1 != len(operators) {
			return db, nil
		}

		str := ""
		nums := []int{}
		for index := 0; index < len(columnNames); index++ {
			if values[index] == nil {
				nums = append(nums, index)
			}
			if index == len(columnNames)-1 {
				str = fmt.Sprintf("%v%v %v", str, columnNames[index], conditions[index])
			} else {
				str = fmt.Sprintf("%v%v %v %v ", str, columnNames[index], conditions[index], operators[index])
			}
		}
		for ind, num := range nums {
			values = append(values[:num], values[num+1:]...)
			for i := ind; i < len(nums); i++ {
				// This is done to adjust indexes because we sliced.
				nums[i] = nums[i] - 1
			}
		}
		db = db.Where(str, values...)
		return db, nil
	}
}

// DoesRecordExist returns true if the record exists.
// 	If ID is to be checked then populate it in the model
func DoesRecordExist(db *gorm.DB, out interface{}, queryProcessors ...QueryProcessor) (bool, error) {
	count := 0
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return false, err
	}

	if err := db.Debug().Model(out).Count(&count).Error; err != nil {
		return false, err
	}
	if count <= 0 {
		return false, nil
	}
	return true, nil
}

// DoesRecordExistForTenant returns true if the record exists for the specific tenant.
// 	If ID is to be checked then populate it in the model
func DoesRecordExistForTenant(db *gorm.DB, tenantID uuid.UUID, out interface{}, queryProcessors ...QueryProcessor) (bool, error) {
	if tenantID == uuid.Nil {
		return false, errors.NewValidationError("DoesRecordExistForTenant: Invalid tenant ID")
	}
	count := 0
	// Below comment would make the tenant check before all query processor (Uncomment only if needed in future)
	// queryProcessors = append([]QueryProcessor{Filter("tenant_id = ?", tenantID)},queryProcessors... )
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return false, err
	}
	if err := db.Debug().Model(out).Where("tenant_id = ?", tenantID).Count(&count).Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// Paging will restrict the output of query with limit and offset and
// 	sets X-Total-Count header with the appropriate total value.
func Paging(w http.ResponseWriter, r *http.Request) QueryProcessor {
	queryparam := mux.Vars(r)
	limitparam := queryparam["limit"]
	offsetparam := queryparam["offset"]
	var err error
	limit := -1
	if len(limitparam) > 0 {
		limit, err = strconv.Atoi(limitparam)
		if err != nil {
			limit = -1
		}
	}
	offset := 0
	if len(offsetparam) > 0 {
		offset, err = strconv.Atoi(offsetparam)
		if err != nil {
			offset = 0
		}
	}
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		if out != nil {
			var totalRecords int
			if err := db.Model(out).Count(&totalRecords).Error; err != nil {
				return db, err
			}
			// Expose headers will let us see the headers
			w.Header().Add("Access-Control-Expose-Headers", "X-Total-Count")
			w.Header().Set("X-Total-Count", strconv.Itoa(totalRecords))
		}

		if limit != -1 {
			db = db.Limit(limit)
		}

		if offset > 0 {
			db = db.Offset(limit * offset)
		}
		return db, nil
	}
}

// Paginate will restrict the output of query with limit and offset & fill totalCount with total records.
func Paginate(limit, offset int, totalCount *int) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		if out != nil {
			if totalCount != nil {
				if err := db.Model(out).Count(totalCount).Error; err != nil {
					return db, err
				}
			}
		}

		if limit != -1 {
			db = db.Limit(limit)
		}

		if offset > 0 {
			db = db.Offset(limit * offset)
		}
		return db, nil
	}
}

// PaginateWithoutModel will restrict the output of query with limit and offset & fill totalCount with total records.
// 	Count does not use model, So the table name has to be specificied for using this function.
func PaginateWithoutModel(limit, offset int, totalCount *int) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		if out != nil {
			if totalCount != nil {
				if err := db.Count(totalCount).Error; err != nil {
					return db, err
				}
			}
		}

		if limit != -1 {
			db = db.Limit(limit)
		}

		if offset > 0 {
			db = db.Offset(limit * offset)
		}
		return db, nil
	}
}

// AddForeignKey adds foreign key to the given scope. e.g:
//     db.Model(&User{}).AddForeignKey("city_id", "cities(id)", "RESTRICT", "RESTRICT")
//     db.Model(&User{}).AddForeignKey("city_id", "cities(id)", "CASCADE", "CASCADE")
func AddForeignKey(db *gorm.DB, field string, dest string, onDelete string, onUpdate string) {
	// scope := s.NewScope(s.Value)
}
