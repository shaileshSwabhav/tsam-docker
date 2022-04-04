package repository

import "github.com/jinzhu/gorm"

// UnitOfWork represent connection
type UnitOfWork struct {
	DB        *gorm.DB
	Committed bool
	Readonly  bool
}

// NewUnitOfWork Create New Instance Of UnitOfWork.
func NewUnitOfWork(db *gorm.DB, readonly bool) *UnitOfWork {
	commit := false
	if readonly {
		return &UnitOfWork{
			DB:        db.New(),
			Committed: commit,
			Readonly:  readonly,
		}
	}

	return &UnitOfWork{
		DB:        db.New().Begin(),
		Committed: commit,
		Readonly:  readonly,
	}
}

// Commit use to commit after a successful transaction.
func (uow *UnitOfWork) Commit() {
	if !uow.Readonly && !uow.Committed {
		uow.Committed = true
		uow.DB.Commit()
	}
}

// RollBack is used to rollback a transaction on failure.
func (uow *UnitOfWork) RollBack() {
	if !uow.Committed && !uow.Readonly {
		uow.DB.Rollback()
	}
}
