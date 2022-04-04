package service

import (
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// ProgrammingQuestionTalentAnswerService provide method to update, delete, add, get method for programming question talent answer.
type ProgrammingQuestionTalentAnswerService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewProgrammingQuestionTalentAnswerService returns new instance of ProgrammingQuestionTalentAnswerService.
func NewProgrammingQuestionTalentAnswerService(db *gorm.DB, repository repository.Repository) *ProgrammingQuestionTalentAnswerService {
	return &ProgrammingQuestionTalentAnswerService{
		DB:         db,
		Repository: repository,
	}
}

// AddProgrammingQuestionTalentAnswer adds new programming question talent answer to database.
func (service *ProgrammingQuestionTalentAnswerService) AddProgrammingQuestionTalentAnswer(answer *programming.ProgrammingQuestionTalentAnswer) error {

	// Validate tenant id.
	err := service.doesTenantExist(answer.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate foreign keys id.
	err = service.doForeignKeysExistForAnswer(answer)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(answer.CreatedBy, answer.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	//  Creating unit of work.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add programming question talent answer to database.
	err = service.Repository.Add(uow, answer)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateProgrammingQuestionTalentAnswer updates programming question talent answer to database.
func (service *ProgrammingQuestionTalentAnswerService) UpdateProgrammingQuestionTalentAnswer(answer *programming.ProgrammingQuestionTalentAnswer) error {

	// Validate tenant ID.
	err := service.doesTenantExist(answer.TenantID)
	if err != nil {
		return err
	}

	// Validate foreign keys id.
	err = service.doForeignKeysExistForAnswer(answer)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate programming question talent answer ID.
	err = service.doesAnswerExist(answer.ID, answer.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(answer.UpdatedBy, answer.TenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update programming question talent answer.
	err = service.Repository.Update(uow, answer)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateProgrammingQuestionTalentAnswerScore updates programming question talent answer's score and isCorrect field to database.
func (service *ProgrammingQuestionTalentAnswerService) UpdateProgrammingQuestionTalentAnswerScore(answer *programming.ProgrammingQuestionTalentAnswerScore,
	tenantID, credentialID uuid.UUID) error {

	// Validate tenant ID.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Validate programming question talent answer ID.
	err = service.doesAnswerExist(answer.ID, tenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, tenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update isCorrect field of programming question talent answer.
	if err := service.Repository.UpdateWithMap(uow, &programming.ProgrammingQuestionTalentAnswer{}, map[string]interface{}{
		"isCorrect": answer.IsCorrect,
		"score":     answer.Score,
		"UpdatedBy": credentialID,
	},
		repository.Filter("id=?", answer.ID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Score could not be updated for the talent answer", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// DeleteProgrammingQuestionTalentAnswer deletes programming question talent answer from database.
func (service *ProgrammingQuestionTalentAnswerService) DeleteProgrammingQuestionTalentAnswer(answer *programming.ProgrammingQuestionTalentAnswer) error {
	credentialID := answer.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(answer.TenantID)
	if err != nil {
		return err
	}

	// Validate programming question talent answer ID.
	err = service.doesAnswerExist(answer.ID, answer.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, answer.TenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update programming question talent answer for updating deleted_by and deleted_at fields of programming question talent answer.
	if err := service.Repository.UpdateWithMap(uow, answer, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", answer.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent answer could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetProgrammingQuestionTalentAnswer returns one programming question talent answer.
func (service *ProgrammingQuestionTalentAnswerService) GetProgrammingQuestionTalentAnswer(answer *programming.ProgrammingQuestionTalentAnswerWithFullQuestionDTO,
	tenantID uuid.UUID) error {

	// Validate tenant ID.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get one programming question talent answer by id.
	err = service.Repository.GetForTenant(uow, tenantID, answer.ID, answer,
		repository.Filter("programming_question_talent_answers.`deleted_at` IS NULL"),
		repository.PreloadAssociations([]string{"ProgrammingQuestion", "ProgrammingQuestion.Solutions",
			"ProgrammingQuestion.Solutions.ProgrammingLanguage", "ProgrammingLanguage", "Talent"}),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "ProgrammingQuestion.TestCases",
			Queryprocessors: []repository.QueryProcessor{
				repository.Filter("programming_question_test_cases.`is_active` = 1"),
				repository.Filter("programming_question_test_cases.`is_hidden` = 0")},
		}),
	)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Check if there is an entry in programming_soltion_is_viewed for the programming question id and talent id.
	// Create bucket for solution is viewed.
	solutionIsViewed := programming.ProgrammingQuestionSolutionIsViewed{}

	// Get solution is viewed from database.
	err = service.Repository.GetRecordForTenant(uow, tenantID, &solutionIsViewed,
		repository.Filter("`programming_question_id`=? AND `talent_id`=?", answer.ProgrammingQuestionID, answer.TalentID))

	// Make solution is viewed true for question.
	answer.ProgrammingQuestion.SolutonIsViewed = true

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// If solution is viewed record not found then make it false.
			answer.ProgrammingQuestion.SolutonIsViewed = false
		} else {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// GetProgrammingQuestionTalentAnswers returns all programming question talent answers.
func (service *ProgrammingQuestionTalentAnswerService) GetProgrammingQuestionTalentAnswers(answers *[]programming.ProgrammingQuestionTalentAnswerDTO, tenantID uuid.UUID,
	form url.Values, limit, offset int, totalCount *int) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// If limit is -1 then make simple query.
	if limit == -1 {

		// Create query precessors, add search quesries.
		queryProcessors := service.addSearchQueries(form)

		// Add all necessary queries.
		queryProcessors = append(queryProcessors,
			repository.Select("programming_question_talent_answers.`created_at` as `date`, programming_question_talent_answers.*"),
			repository.Filter("`programming_question_option_id` IS NULL"),
			repository.Filter("programming_question_talent_answers.`deleted_at` IS NULL"),
			repository.Filter("programming_question_talent_answers.`tenant_id`=?", tenantID),
			repository.PreloadAssociations([]string{"ProgrammingQuestion", "ProgrammingLanguage", "Talent"}))

		// Get all programming question talent answers from database.
		err = service.Repository.GetAllInOrder(uow, answers, "programming_question_talent_answers.`is_correct` IS NULL DESC, programming_question_talent_answers.`created_at` DESC",
			queryProcessors...)
		if err != nil {
			uow.RollBack()
			return err
		}
	} else {

		// Create query precessors, add search quesries.
		queryProcessors := service.addSearchQueries(form)

		// Add all necessary queries for sub query one.
		queryProcessors = append(queryProcessors,
			repository.Table("programming_question_talent_answers"),
			repository.Select("programming_question_talent_answers.`created_at` as `date`, programming_question_talent_answers.*, "+
				"IF(programming_question_talent_answers.is_correct is null, 1, 0) as not_checked"),
			repository.Filter("`programming_question_option_id` IS NULL"),
			repository.Filter("programming_question_talent_answers.`deleted_at` IS NULL"),
			repository.Filter("programming_question_talent_answers.`tenant_id`=?", tenantID))

		// Create query expression for sub query one.
		subQueryOne, err := service.Repository.SubQuery(uow, programming.ProgrammingQuestionSolution{}, queryProcessors...)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

		// Get all programming question talent answers by limit and offset.
		if err := service.Repository.GetAll(uow, &answers,
			repository.RawQuery("SELECT COUNT(*) as total_answers, sum(not_checked) as total_not_checked, subone.* "+
				"from ? as subone GROUP BY programming_question_talent_answers.`programming_question_id`,  "+
				"programming_question_talent_answers.`programming_language_id`, programming_question_talent_answers.`talent_id` "+
				"ORDER BY subone.`is_correct` IS NULL DESC, subone.`created_at` DESC LIMIT ? OFFSET ?", subQueryOne, limit, (limit*offset)),
		); err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

		// Get count of all programming question talent answers.
		// Create query expression for sub query two.
		subQueryTwo, err := service.Repository.SubQuery(uow, programming.ProgrammingQuestionSolution{},
			repository.RawQuery("SELECT COUNT(*) as total_answers, sum(not_checked) as total_not_checked, subone.* "+
				"from ? as subone GROUP BY programming_question_talent_answers.`programming_question_id`,  "+
				"programming_question_talent_answers.`programming_language_id`, programming_question_talent_answers.`talent_id` "+
				"ORDER BY subone.`is_correct` IS NULL DESC, subone.`created_at` DESC", subQueryOne),
		)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

		// Create bucket for total count.
		totalCountModel := programming.TotalCount{}

		// Get total count of programming programming question talent answers from database.
		if err := service.Repository.GetAll(uow, &totalCountModel,
			repository.RawQuery("SELECT COUNT(*) as total_count FROM ? as subtwo", subQueryTwo),
		); err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

		// Set the totalCount variable.
		*totalCount = totalCountModel.TotalCount
	}

	for index := range *answers {

		// If limit is not -1 then preload programming question talent answers.
		if limit != -1 {

			// Preload all programming question talent answers.
			if err := service.Repository.Get(uow, (*answers)[index].ID, &(*answers)[index],
				repository.PreloadAssociations([]string{"ProgrammingQuestion", "ProgrammingLanguage", "Talent"}),
			); err != nil {
				uow.RollBack()
				log.NewLogger().Error(err.Error())
				return errors.NewValidationError("Record not found")
			}
		}

		// Create bucket for solution is viewed.
		solutionIsViewed := programming.ProgrammingQuestionSolutionIsViewed{}

		// Get solution is viewed from database.
		err = service.Repository.GetRecordForTenant(uow, tenantID, &solutionIsViewed,
			repository.Filter("`programming_question_id`=? AND `talent_id`=?", (*answers)[index].ProgrammingQuestion.ID,
				(*answers)[index].TalentID))

		// Make solution is viewed true for question of programming question talent answers.
		(*answers)[index].ProgrammingQuestion.SolutonIsViewed = true

		if err != nil {
			if err == gorm.ErrRecordNotFound {

				// If solution is viewed record not found then make it false for question of programming question talent answers.
				(*answers)[index].ProgrammingQuestion.SolutonIsViewed = false
			} else {
				uow.RollBack()
				return err
			}
		}
	}

	uow.Commit()
	return nil
}

// AddProgrammingQuestionSolutionIsViewed adds new programming question talent answer is views entry to database.
func (service *ProgrammingQuestionTalentAnswerService) AddProgrammingQuestionSolutionIsViewed(solutionIsViewed *programming.ProgrammingQuestionSolutionIsViewed) error {

	// Validate tenant id.
	err := service.doesTenantExist(solutionIsViewed.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate if foreign key ids exist for the solution is viewed.
	err = service.doForeignKeysExistForSolutionIsViewed(solutionIsViewed)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(solutionIsViewed.CreatedBy, solutionIsViewed.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	//  Creating unit of work.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add solution is viewed to database.
	err = service.Repository.Add(uow, solutionIsViewed)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// addSearchQueries adds search queries.
func (service *ProgrammingQuestionTalentAnswerService) addSearchQueries(requestForm url.Values) []repository.QueryProcessor {

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	var queryProcesors []repository.QueryProcessor

	// Programming type.
	if programmingType, ok := requestForm["programmingType"]; ok {
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN programming_questions ON programming_question_talent_answers.`programming_question_id` = programming_questions.`id`"),
			repository.Join("JOIN programming_question_types ON programming_questions.`programming_question_type_id` = programming_question_types.`id`"),
			repository.Filter("programming_question_types.`deleted_at` IS NULL"))

		util.AddToSlice("programming_question_types.`programming_type`", "=?", "AND", programmingType,
			&columnNames, &conditions, &operators, &values)
	}

	// Prgramming question id.
	if programmingQuestionID, ok := requestForm["programmingQuestionID"]; ok {
		util.AddToSlice("programming_question_talent_answers.`programming_question_id`", "=?", "AND", programmingQuestionID,
			&columnNames, &conditions, &operators, &values)
	}

	// Prgramming language id.
	if programmingLanguageID, ok := requestForm["programmingLanguageID"]; ok {
		util.AddToSlice("programming_question_talent_answers.`programming_language_id`", "=?", "AND", programmingLanguageID,
			&columnNames, &conditions, &operators, &values)
	}

	// Talent id.
	if talentID, ok := requestForm["talentID"]; ok {
		util.AddToSlice("programming_question_talent_answers.`talent_id`", "=?", "AND", talentID,
			&columnNames, &conditions, &operators, &values)
	}

	// Add all filters
	queryProcesors = append(queryProcesors,
		repository.FilterWithOperator(columnNames, conditions, operators, values))

	return queryProcesors
}

// doesTalentAndQuestionExist validates if the talent id and question id already exist for the programming question talent answer.
func (service *ProgrammingQuestionTalentAnswerService) doesTalentAndQuestionExist(answer *programming.ProgrammingQuestionTalentAnswer) error {

	// If programing language id is nil.
	if answer.ProgrammingLanguageID == nil {
		exists, err := repository.DoesRecordExistForTenant(service.DB, answer.TenantID, &programming.ProgrammingQuestionTalentAnswer{},
			repository.Filter("`talent_id`=? AND `programming_question_id`=? AND `id`!=?",
				answer.TalentID, answer.ProgrammingQuestionID, answer.ID))
		if err := util.HandleIfExistsError("You have already answered this question", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}
	// If programing language id is not nil.
	if answer.ProgrammingLanguageID != nil {
		exists, err := repository.DoesRecordExistForTenant(service.DB, answer.TenantID, &programming.ProgrammingQuestionTalentAnswer{},
			repository.Filter("`talent_id`=? AND `programming_question_id`=? AND `id`!=? AND `programming_language_id`=?",
				answer.TalentID, answer.ProgrammingQuestionID, answer.ID, answer.ProgrammingLanguageID))
		if err := util.HandleIfExistsError("You have already answered this question", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	return nil
}

// doForeignKeysExistForSolutionIsViewed validates if the foreign key ids exist for the solution is viewed.
func (service *ProgrammingQuestionTalentAnswerService) doForeignKeysExistForSolutionIsViewed(solutionIsViewed *programming.ProgrammingQuestionSolutionIsViewed) error {

	// Check for programming question id.
	exists, err := repository.DoesRecordExistForTenant(service.DB, solutionIsViewed.TenantID, programming.ProgrammingQuestion{},
		repository.Filter("`id` = ?", solutionIsViewed.ProgrammingQuestionID))
	if err := util.HandleError("Invalid programming question ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check for talent id.
	exists, err = repository.DoesRecordExistForTenant(service.DB, solutionIsViewed.TenantID, tal.Talent{},
		repository.Filter("`id` = ?", solutionIsViewed.TalentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doForeignKeysExistForAnswer validates if the foreign key ids exist for the programming question talent answer.
func (service *ProgrammingQuestionTalentAnswerService) doForeignKeysExistForAnswer(answer *programming.ProgrammingQuestionTalentAnswer) error {

	// Check for programming question id.
	exists, err := repository.DoesRecordExistForTenant(service.DB, answer.TenantID, programming.ProgrammingQuestion{},
		repository.Filter("`id` = ?", answer.ProgrammingQuestionID))
	if err := util.HandleError("Invalid programming question ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check for talent id.
	exists, err = repository.DoesRecordExistForTenant(service.DB, answer.TenantID, tal.Talent{},
		repository.Filter("`id` = ?", answer.TalentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check for programming question option id if it exists.
	if answer.ProgrammingQuestionOptionID != nil {
		exists, err = repository.DoesRecordExistForTenant(service.DB, answer.TenantID, programming.ProgrammingQuestionOption{},
			repository.Filter("`id` = ?", answer.ProgrammingQuestionOptionID))
		if err := util.HandleError("Invalid programming question option ID", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *ProgrammingQuestionTalentAnswerService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesAnswerExist validates if programming question talent answer exists or not in database.
func (service *ProgrammingQuestionTalentAnswerService) doesAnswerExist(answerID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingQuestionTalentAnswer{},
		repository.Filter("`id` = ?", answerID))
	if err := util.HandleError("Invalid programming question talent answer ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *ProgrammingQuestionTalentAnswerService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
