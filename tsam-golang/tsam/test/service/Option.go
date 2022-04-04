package service

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/repository"
	tst "github.com/techlabs/swabhav/tsam/models/test"
	"github.com/techlabs/swabhav/tsam/util"
)

// OptionService Provide method to Update, Delete, Add, Get Method For Option.
type OptionService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewOptionService creates a new instance of OptionService
func NewOptionService(db *gorm.DB, repository repository.Repository) *OptionService {
	return &OptionService{
		DB:         db,
		Repository: repository,
	}

}

// AddOption Add New Option to Database.
func (ser *OptionService) AddOption(option *tst.Option) error {

	// Add ID to Association if Not Present Or Return Error
	err := option.ValidateOption()
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	option.ID = util.GenerateUUID()

	// Create New Instace of UnitOfWork For Write Operation
	// Add Option to Database
	uow := repository.NewUnitOfWork(ser.DB, false)
	err = ser.Repository.Add(uow, option)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// AddOptions Add Multiple Option to Database.
func (ser *OptionService) AddOptions(Options *[]tst.Option, OptionsIDs *[]uuid.UUID) error {

	// Add individual Option To Database
	for _, Option := range *Options {

		// Add ID To Option
		Option.ID = util.GenerateUUID()

		// Add Option To Database
		err := ser.AddOption(&Option)
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
		*OptionsIDs = append(*OptionsIDs, Option.ID)
	}
	return nil
}

// UpdateOption Update Option data By Taking
// Option Struct & QueryProcessor
func (ser *OptionService) UpdateOption(option *tst.Option, queryProcessor ...repository.QueryProcessor) error {

	// ValidateOption function validate Option struct
	// for all required field
	err := option.ValidateOption()
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// GetOption Call for check Option with provided ID
	// Exist or not in database
	tempOption := tst.Option{}
	err = ser.GetOption(&tempOption, &option.ID)
	if err != nil {
		return err
	}

	// Create New UnitOfWork Instance For Update Option
	// By Repository
	uow := repository.NewUnitOfWork(ser.DB, false)
	err = ser.Repository.Update(uow, option)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteOption Delete Option & also there Options
// By Providing Option ID
func (ser *OptionService) DeleteOption(OptionID *uuid.UUID, queryProcessor ...repository.QueryProcessor) error {

	// Get Option From Database
	Option := tst.Option{}
	err := ser.GetOption(&Option, OptionID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Deleting redudant options in Option
	uow := repository.NewUnitOfWork(ser.DB, false)

	// Delete Option From Database
	err = ser.Repository.Delete(uow, &Option)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetOption Add New Option to Database.
func (ser *OptionService) GetOption(Option *tst.Option, OptionID *uuid.UUID,
	queryProcessor ...repository.QueryProcessor) error {
	// Get Option By ID From Database
	uow := repository.NewUnitOfWork(ser.DB, true)
	err := ser.Repository.Get(uow, *OptionID, Option, queryProcessor...)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// GetOptions Return All Option From Database.
func (ser *OptionService) GetOptions(Options *[]tst.Option, queryProcessor ...repository.QueryProcessor) error {
	uow := repository.NewUnitOfWork(ser.DB, true)

	// Get All Option Options From Database
	err := ser.Repository.GetAll(uow, Options, queryProcessor...)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
