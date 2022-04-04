package service

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// ProgrammingQuestionService provides methods to do different CRUD operations on programming_questions table.
type ProgrammingQuestionService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// NewProgrammingQuestionService returns a new instance Of ProgrammingQuestionService.
func NewProgrammingQuestionService(db *gorm.DB, repository repository.Repository) *ProgrammingQuestionService {
	return &ProgrammingQuestionService{
		DB:          db,
		Repository:  repository,
		association: []string{},
	}
}

// AddProgrammingQuestion will add programming questions to the table
func (service *ProgrammingQuestionService) AddProgrammingQuestion(question *programming.ProgrammingQuestion,
	uows ...*repository.UnitOfWork) error {

	fmt.Println("------------------------------------------", question)
	// Check if all foreign keys exist.
	err := service.doesForeignKeyExist(question, question.CreatedBy)
	if err != nil {
		return err
	}

	// Give created by and tenant id to all options of the question.
	if question.Options != nil {
		for index := range question.Options {
			question.Options[index].TenantID = question.TenantID
			question.Options[index].CreatedBy = question.CreatedBy
		}
	}

	// Give created by and tenant id to all test cases of the question.
	if question.TestCases != nil {
		for index := range question.TestCases {
			question.TestCases[index].TenantID = question.TenantID
			question.TestCases[index].CreatedBy = question.CreatedBy
		}
	}

	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Add question.
	err = service.Repository.Add(uow, question)
	if err != nil {
		if length == 0 {
			uow.RollBack()
		}
		return err
	}

	if length == 0 {
		uow.Commit()
	}

	return nil
}

// AddProgrammingQuestions will add multiple programming questions to the table
func (service *ProgrammingQuestionService) AddProgrammingQuestions(questions *[]programming.ProgrammingQuestion, tenantID,
	credentialID uuid.UUID) error {

	uow := repository.NewUnitOfWork(service.DB, false)
	for _, question := range *questions {
		question.TenantID = tenantID
		question.CreatedBy = credentialID
		err := service.AddProgrammingQuestion(&question)
		if err != nil {
			uow.RollBack()
			return err
		}
	}
	uow.Commit()
	return nil
}

// UpdateProgrammingQuestion will update the specified programming question.
func (service *ProgrammingQuestionService) UpdateProgrammingQuestion(question *programming.ProgrammingQuestion) error {

	// if !question.IsActive {
	// 	return errors.NewValidationError("Inactive questions cannot be updated.")
	// }

	// Check if all foreign key exist.
	err := service.doesForeignKeyExist(question, question.UpdatedBy)
	if err != nil {
		return err
	}

	// Check is feedback question exist
	err = service.doesProgrammingQuestionExist(question.TenantID, question.ID)
	if err != nil {
		return err
	}

	// If the question has any test cases then make it nil.
	question.TestCases = nil

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// If question has options then delete all its test cases.
	if question.HasOptions {

		// Update programming question test case for updating deleted_by and deleted_at fields of programming question test case.
		if err := service.Repository.UpdateWithMap(uow, &programming.ProgrammingQuestionTestCase{}, map[string]interface{}{
			"DeletedBy": question.UpdatedBy,
			"DeletedAt": time.Now(),
		},
			repository.Filter("`tenant_id`=?", question.TenantID),
			repository.Filter("`programming_question_id`=?", question.ID),
		); err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError("Programming question could not be updated", http.StatusInternalServerError)
		}
	}

	// If question does not have options then delete all its options.
	if !question.HasOptions {
		question.Options = nil
	}

	// Update options of question.
	err = service.updateOptions(uow, question, question.UpdatedBy)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Create bucket for getting programming question already present in database.
	tempQuestion := programming.ProgrammingQuestion{}

	// Get programming question round for getting created_by field of programming question from database.
	if err := service.Repository.GetForTenant(uow, question.TenantID, question.ID, &tempQuestion); err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp programming question to programming question to be updated.
	question.CreatedBy = tempQuestion.CreatedBy

	// Update programming question associations.
	if err := service.updateProgrammingQuestionAssociations(uow, question); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Programming question could not be updated", http.StatusInternalServerError)
	}

	// Make programming question types nil so that it is not inserted again.
	question.ProgrammingQuestionTypes = nil

	// Make programming concepts nil so that it is not inserted again.
	question.ProgrammingConcepts = nil

	// Make programming langauges nil so that it is not inserted again.
	question.ProgrammingLanguages = nil

	// Update question in database.
	err = service.Repository.Save(uow, question)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateProgrammingQuestionIsActive updates question's isActive field to database.
func (service *ProgrammingQuestionService) UpdateProgrammingQuestionIsActive(question *programming.ProgrammingQuestionIsActive,
	tenantID, credentialID uuid.UUID) error {

	// Validate tenant ID.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Validate question ID.
	err = service.doesProgrammingQuestionExist(tenantID, question.ID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(tenantID, credentialID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update isActive field of solution.
	if err := service.Repository.UpdateWithMap(uow, &programming.ProgrammingQuestion{}, map[string]interface{}{
		"IsActive":  question.IsActive,
		"UpdatedBy": credentialID,
	},
		repository.Filter("id=?", question.ID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Is active could not be updated for the programming question", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// DeleteProgrammingQuestion will delete the specified programming question.
func (service *ProgrammingQuestionService) DeleteProgrammingQuestion(question *programming.ProgrammingQuestion) error {

	// check if tenant exist.
	err := service.doesTenantExist(question.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist.
	err = service.doesCredentialExist(question.TenantID, question.DeletedBy)
	if err != nil {
		return err
	}

	// check is feedback question exist.
	err = service.doesProgrammingQuestionExist(question.TenantID, question.ID)
	if err != nil {
		return err
	}

	// Make ProgrammingQuestionTypes nil to avoid any updates or inserts.
	question.ProgrammingQuestionTypes = nil

	// Make programming concepts nil  to avoid any updates or inserts.
	question.ProgrammingConcepts = nil

	// Make programming languages nil  to avoid any updates or inserts.
	question.ProgrammingLanguages = nil

	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.UpdateWithMap(uow, new(programming.ProgrammingQuestion), map[string]interface{}{
		"DeletedBy": question.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id` = ?", question.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.deleteProgrammingOptions(uow, question.DeletedBy, question.ID)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetProgrammingQuestions will return all the programming question from the table.
func (service *ProgrammingQuestionService) GetProgrammingQuestions(questions *[]programming.ProgrammingQuestionDTO,
	form url.Values, tenantID uuid.UUID, limit, offset int, totalCount *int) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Create query precessors.
	queryProcessors := service.addSearchQueries(form)
	queryProcessors = append(queryProcessors,
		repository.Filter("programming_questions.`tenant_id`=?", tenantID),
		repository.PreloadAssociations([]string{"ProgrammingQuestionTypes", "ProgrammingConcepts", "ProgrammingLanguages"}),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "Options",
			Queryprocessors: []repository.QueryProcessor{
				repository.Filter("programming_question_options.`is_active` = 1"),
				repository.OrderBy("programming_question_options.`order`")},
		}),
		repository.GroupBy("programming_questions.`id`"),
		// repository.OrderBy("programming_questions.`is_active` DESC"),
		repository.Paginate(limit, offset, totalCount))

	err = service.Repository.GetAll(uow, questions, queryProcessors...)
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range *questions {

		// Get solution count.
		var totalCount int

		err = service.Repository.GetCount(uow, &programming.ProgrammingQuestion{}, &totalCount,
			repository.Join("JOIN programming_question_solutions on programming_questions.`id` = programming_question_solutions.`programming_question_id`"),
			repository.Filter("programming_questions.`deleted_at` IS NULL AND programming_question_solutions.`deleted_at` IS NULL"),
			repository.Filter("programming_questions.`tenant_id` = ?", tenantID),
			repository.Filter("programming_questions.`id` = ?", (*questions)[index].ID))
		if err != nil {
			uow.RollBack()
			return err
		}
		(*questions)[index].SolutionCount = uint16(totalCount)

		// Check if any talent has answered this question.
		var answeredCount int

		err = service.Repository.GetCount(uow, &programming.ProgrammingQuestionTalentAnswer{}, &answeredCount,
			repository.Filter("programming_question_talent_answers.`deleted_at` IS NULL"),
			repository.Filter("programming_question_talent_answers.`tenant_id` = ?", tenantID),
			repository.Filter("programming_question_talent_answers.`programming_question_id` = ?", (*questions)[index].ID))
		if err != nil {
			uow.RollBack()
			return err
		}

		(*questions)[index].HasAnyTalentAnswered = false

		if answeredCount > 0 {
			(*questions)[index].HasAnyTalentAnswered = true
		}
	}

	uow.Commit()
	return nil
}

// GetProgrammingQuestionList will return all the feedback question from the table.
func (service *ProgrammingQuestionService) GetProgrammingQuestionList(questions *[]list.ProgrammingQuestion,
	tenantID uuid.UUID, form url.Values) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create query precessors.
	quesryProcessors := service.addSearchQueries(form)
	quesryProcessors = append(quesryProcessors,
		repository.Filter("programming_questions.`tenant_id`=?", tenantID),
		repository.Filter("programming_questions.is_active = 1"))

	err = service.Repository.GetAllInOrder(uow, questions, "is_active DESC", quesryProcessors...)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetProgrammingQuestionsPractice will return all the questions in problem of the day format for practice.
func (service *ProgrammingQuestionService) GetProgrammingQuestionsPractice(questions *[]programming.QuestionProblemOfTheDayDTO,
	tenantID uuid.UUID, form url.Values, limit, offset int, totalCount *int) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create query precessors.
	quesryProcessors := service.addSearchQueries(form)
	quesryProcessors = append(quesryProcessors,
		repository.Filter("programming_questions.`tenant_id`=?", tenantID),
		repository.Filter("programming_questions.is_active = 1"),
		repository.Paginate(limit, offset, totalCount))

	err = service.Repository.GetAllInOrder(uow, questions, "`label` DESC", quesryProcessors...)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Get attempted by and solved by count.
	for index := range *questions {
		var totalCount int

		// Count attempted by.
		err = service.Repository.GetCount(uow, &programming.ProgrammingQuestion{}, &totalCount,
			repository.Join("JOIN programming_question_talent_answers on programming_questions.`id` = programming_question_talent_answers.`programming_question_id`"),
			repository.Filter("programming_questions.`deleted_at` IS NULL AND programming_question_talent_answers.`deleted_at` IS NULL"),
			repository.Filter("programming_questions.`tenant_id` = ?", tenantID),
			repository.Filter("programming_questions.`id` = ?", (*questions)[index].ID),
			repository.GroupBy("`talent_id`"))
		if err != nil {
			uow.RollBack()
			return err
		}
		(*questions)[index].AttemptedByCount = uint(totalCount)

		// Count solved by.
		err = service.Repository.GetCount(uow, &programming.ProgrammingQuestion{}, &totalCount,
			repository.Join("JOIN programming_question_talent_answers on programming_questions.`id` = programming_question_talent_answers.`programming_question_id`"),
			repository.Filter("programming_questions.`deleted_at` IS NULL AND programming_question_talent_answers.`deleted_at` IS NULL"),
			repository.Filter("programming_questions.`tenant_id` = ?", tenantID),
			repository.Filter("programming_questions.`id` = ?", (*questions)[index].ID),
			repository.Filter("programming_questions.`score` = programming_question_talent_answers.`score`"),
			repository.GroupBy("`talent_id`"))
		if err != nil {
			uow.RollBack()
			return err
		}
		(*questions)[index].SolvedByCount = uint(totalCount)
	}

	uow.Commit()
	return nil
}

// GetProgrammingQuestion returns one question by id.
func (service *ProgrammingQuestionService) GetProgrammingQuestion(question *programming.ProgrammingQuestionWithTalentAnswerDTO,
	tenantID uuid.UUID, form url.Values) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create query precessors.
	quesryProcessors := service.addSearchQueries(form)
	quesryProcessors = append(quesryProcessors,
		repository.Join("LEFT JOIN programming_question_talent_answers ON programming_questions.`id` = programming_question_talent_answers.`programming_question_id`"),
		repository.Select("IF(talent_id IS NULL, false, true) AS is_answered, talent_id, programming_question_option_id, answer, "+
			"programming_question_talent_answers.`programming_language_id`, programming_questions.*"),
		repository.Filter("programming_questions.`tenant_id`=?", tenantID),
		repository.Filter("programming_question_talent_answers.`deleted_at` IS NULL"),
		repository.OrderBy("programming_question_talent_answers.`created_at` DESC"),
		repository.PreloadAssociations([]string{"Solutions", "Solutions.ProgrammingLanguage",
			"ProgrammingLanguage", "ProgrammingQuestionTypes", "ProgrammingConcepts", "ProgrammingLanguages"}),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "Options",
			Queryprocessors: []repository.QueryProcessor{
				repository.Filter("programming_question_options.is_active = 1"),
				repository.OrderBy("programming_question_options.order")},
		}),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "TestCases",
			Queryprocessors: []repository.QueryProcessor{
				repository.Filter("programming_question_test_cases.`is_active` = 1"),
				repository.Filter("programming_question_test_cases.`is_hidden` = 0")},
		}),
	)

	// Get one question by id.
	hasTalentAnswered := true
	err = service.Repository.GetRecord(uow, question, quesryProcessors...)

	// If record not found then return question without adding search queries.
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			hasTalentAnswered = false
		} else {
			uow.RollBack()
			return err
		}
	}

	if !hasTalentAnswered {

		// If talent has not answered then return the question.
		err = service.Repository.GetRecord(uow, question,
			repository.Filter("programming_questions.`tenant_id`=?", tenantID),
			repository.PreloadAssociations([]string{"Solutions", "Solutions.ProgrammingLanguage",
				"ProgrammingQuestionTypes", "ProgrammingConcepts", "ProgrammingLanguages"}),
			repository.PreloadWithCustomCondition(repository.Preload{
				Schema: "Options",
				Queryprocessors: []repository.QueryProcessor{
					repository.Filter("programming_question_options.is_active = 1"),
					repository.OrderBy("programming_question_options.order")},
			}),
			repository.PreloadWithCustomCondition(repository.Preload{
				Schema: "TestCases",
				Queryprocessors: []repository.QueryProcessor{
					repository.Filter("programming_question_test_cases.`is_active` = 1"),
					repository.Filter("programming_question_test_cases.`is_hidden` = 0")},
			}),
		)

		if err != nil {
			uow.RollBack()
			return err
		}
	}

	// Create bucket for attempted by count and solved by count.
	var totalCount int

	// Count attempted by.
	err = service.Repository.GetCount(uow, &programming.ProgrammingQuestion{}, &totalCount,
		repository.Join("JOIN programming_question_talent_answers on programming_questions.`id` = programming_question_talent_answers.`programming_question_id`"),
		repository.Filter("programming_questions.`deleted_at` IS NULL AND programming_question_talent_answers.`deleted_at` IS NULL"),
		repository.Filter("programming_questions.`tenant_id` = ?", tenantID),
		repository.Filter("programming_questions.`id` = ?", question.ID),
		repository.GroupBy("`talent_id`"))
	if err != nil {
		uow.RollBack()
		return err
	}
	question.AttemptedByCount = uint(totalCount)

	// Count solved by.
	err = service.Repository.GetCount(uow, &programming.ProgrammingQuestion{}, &totalCount,
		repository.Join("JOIN programming_question_talent_answers on programming_questions.`id` = programming_question_talent_answers.`programming_question_id`"),
		repository.Filter("programming_questions.`deleted_at` IS NULL AND programming_question_talent_answers.`deleted_at` IS NULL"),
		repository.Filter("programming_questions.`tenant_id` = ?", tenantID),
		repository.Filter("programming_questions.`id` = ?", question.ID),
		repository.Filter("programming_question_talent_answers.`is_correct` = 1"),
		repository.GroupBy("`talent_id`"))
	if err != nil {
		uow.RollBack()
		return err
	}
	question.SolvedByCount = uint(totalCount)

	// Check if there is an entry in programming_soltion_is_viewed for the programming question id and talent id.
	// Create bucket for solution is viewed.
	solutionIsViewed := programming.ProgrammingQuestionSolutionIsViewed{}

	// Get talent id from from values.
	talentID := ""
	if _, ok := form["talentID"]; ok {
		talentID = form.Get("talentID")
	} else {
		uow.RollBack()
		return errors.NewValidationError("Talent id is not present")
	}

	// Get solution is viewed from database.
	err = service.Repository.GetRecordForTenant(uow, tenantID, &solutionIsViewed,
		repository.Filter("`programming_question_id`=? AND `talent_id`=?", question.ID, talentID))

	// Make solution is viewed true for question.
	question.SolutonIsViewed = true

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// If solution is viewed record not found then make it false.
			question.SolutonIsViewed = false
		} else {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// GetProgrammingQuestionForTopic will get all programming question for specified concepts.
func (service *ProgrammingQuestionService) GetProgrammingQuestionForTopic(tenantID uuid.UUID,
	programmingQuestions *[]programming.ProgrammingQuestionDTO, totalCount *int, parser *web.Parser) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	tempTopicQuestions := []programming.TopicProgrammingQuestionDTO{}

	queryProcessors := []repository.QueryProcessor{}
	queryProcessors = append(queryProcessors, service.addTopicSearchQueries(parser.Form)...)
	queryProcessors = append(queryProcessors, repository.Select("`id`, `programming_question_id`"),
		repository.Filter("`topic_programming_questions`.`tenant_id` = ?", tenantID))

	err = service.Repository.GetAll(uow, &tempTopicQuestions, queryProcessors...)
	if err != nil {
		uow.RollBack()
		return err
	}

	var programmingQuestionIDs []uuid.UUID

	for i := 0; i < len(tempTopicQuestions); i++ {
		programmingQuestionIDs = append(programmingQuestionIDs, tempTopicQuestions[i].ProgrammingQuestionID)
	}

	limit, offset := parser.ParseLimitAndOffset()

	queryProcessors = []repository.QueryProcessor{}
	queryProcessors = append(queryProcessors, service.addSearchQueries(parser.Form)...)

	if len(tempTopicQuestions) > 0 {
		queryProcessors = append(queryProcessors, repository.Filter("`programming_questions`.`id` NOT IN (?)", programmingQuestionIDs))
	}

	queryProcessors = append(queryProcessors, repository.Filter("programming_questions.`tenant_id` = ?", tenantID),
		repository.GroupBy("programming_questions.`id`"), repository.Paginate(limit, offset, totalCount))

	err = service.Repository.GetAll(uow, programmingQuestions, queryProcessors...)
	if err != nil {
		uow.RollBack()
		return err
	}

	// for i := 0; i < len(tempTopicQuestions); i++ {
	// 	for j := 0; j < len(*programmingQuestions); j++ {
	// 		if tempTopicQuestions[i].ProgrammingQuestionID == (*programmingQuestions)[j].ID {
	// 			(*programmingQuestions) = append((*programmingQuestions)[:j], (*programmingQuestions)[j+1:]...)
	// 		}
	// 	}
	// }

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// updateProgrammingQuestionAssociations updates programming question's associations.
func (service *ProgrammingQuestionService) updateProgrammingQuestionAssociations(uow *repository.UnitOfWork, question *programming.ProgrammingQuestion) error {

	// Replace programming question type.
	if err := service.Repository.ReplaceAssociations(uow, question, "ProgrammingQuestionTypes",
		question.ProgrammingQuestionTypes); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Replace programming concept.
	if err := service.Repository.ReplaceAssociations(uow, question, "ProgrammingConcepts",
		question.ProgrammingConcepts); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Replace programming language.
	if err := service.Repository.ReplaceAssociations(uow, question, "ProgrammingLanguages",
		question.ProgrammingLanguages); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	return nil
}

// updateOptions will check the new options with the options already in table.
func (service *ProgrammingQuestionService) updateOptions(uow *repository.UnitOfWork,
	question *programming.ProgrammingQuestion, credentialID uuid.UUID) error {

	if question.Options == nil {
		err := service.Repository.UpdateWithMap(uow, new(programming.ProgrammingQuestionOption), map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		}, repository.Filter("`programming_question_id`=?", question.ID))
		if err != nil {
			return err
		}
		return nil
	}

	programmingOptions := question.Options
	tempProgrammingOptions := []programming.ProgrammingQuestionOption{}

	err := service.Repository.GetAllForTenant(uow, question.TenantID, &tempProgrammingOptions,
		repository.Filter("`programming_question_id`=?", question.ID))
	if err != nil {
		return err
	}

	programmingOptionMap := make(map[uuid.UUID]uint)

	// initialize map.
	for _, tempProgrammingOption := range tempProgrammingOptions {
		programmingOptionMap[tempProgrammingOption.ID]++
	}

	for _, programmingOption := range programmingOptions {

		if util.IsUUIDValid(programmingOption.ID) {
			programmingOptionMap[programmingOption.ID]++

			if programmingOptionMap[programmingOption.ID] > 1 {
				programmingOption.UpdatedBy = credentialID
				err = service.Repository.Update(uow, &programmingOption)
				if err != nil {
					return err
				}
				programmingOptionMap[programmingOption.ID] = 0
			}
			continue
		}
		programmingOption.CreatedBy = credentialID
		programmingOption.TenantID = question.TenantID
		programmingOption.ProgrammingQuestionID = question.ID
		err = service.Repository.Add(uow, &programmingOption)
		if err != nil {
			return err
		}
	}

	for _, tempProgrammingOption := range tempProgrammingOptions {
		if programmingOptionMap[tempProgrammingOption.ID] == 1 {
			err = service.Repository.UpdateWithMap(uow, new(programming.ProgrammingQuestionOption), map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			}, repository.Filter("`id` = ?", tempProgrammingOption.ID))
			if err != nil {
				return err
			}
			programmingOptionMap[tempProgrammingOption.ID] = 0
		}
	}
	question.Options = nil
	return nil
}

// deleteProgrammingOptions will delete options for specified question.
func (service *ProgrammingQuestionService) deleteProgrammingOptions(uow *repository.UnitOfWork,
	credentialID, feedbackQuestionID uuid.UUID) error {

	err := service.Repository.UpdateWithMap(uow, new(programming.ProgrammingQuestionOption), map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`programming_question_id`=?", feedbackQuestionID))
	if err != nil {
		return err
	}
	return nil
}

// doesForeignKeyExist will check if all foreign keys are valid.
func (service *ProgrammingQuestionService) doesForeignKeyExist(question *programming.ProgrammingQuestion,
	credentialID uuid.UUID) error {

	// check if tenant exist.
	err := service.doesTenantExist(question.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist.
	err = service.doesCredentialExist(question.TenantID, credentialID)
	if err != nil {
		return err
	}

	// Check if question label is unique.
	err = service.doesProgrammingQuestionLabelExist(question.TenantID, question.ID, question.Label)
	if err != nil {
		return err
	}
	return nil
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *ProgrammingQuestionService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, new(general.Tenant),
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *ProgrammingQuestionService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, new(general.Credential),
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesProgrammingQuestionExist returns error if there is no programming question record for the given tenant.
func (service *ProgrammingQuestionService) doesProgrammingQuestionExist(tenantID, programmingQuestionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, new(programming.ProgrammingQuestion),
		repository.Filter("`id` = ?", programmingQuestionID))
	if err := util.HandleError("Invalid programming question ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesProgrammingQuestionLabelExist returns error if there is no programming question record for the given tenant.
func (service *ProgrammingQuestionService) doesProgrammingQuestionLabelExist(tenantID, programmingQuestionID uuid.UUID,
	programLabel string) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, new(programming.ProgrammingQuestion),
		repository.Filter("`id` != ? AND `label` = ?", programmingQuestionID, programLabel))
	if err := util.HandleIfExistsError("Question label already exist", exists, err); err != nil {
		return err
	}
	return nil
}

func (service *ProgrammingQuestionService) addTopicSearchQueries(requestForm url.Values) []repository.QueryProcessor {

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	var queryProcessors []repository.QueryProcessor

	if topicIDs, ok := requestForm["topicIDs"]; ok {
		util.AddToSlice("topic_programming_questions.`topic_id`", "IN (?)", "AND", topicIDs, &columnNames, &conditions, &operators, &values)
	}

	// queryProcessors = append(queryProcessors,
	// 	repository.Join("LEFT OUTER JOIN topic_programming_questions ON programming_questions.`id` = topic_programming_questions.`programming_question_id` AND "+
	// 		"topic_programming_questions.`tenant_id` = programming_questions.`tenant_id`"),
	// 	repository.Filter("topic_programming_questions.`deleted_at` IS NULL"))
	// if len(topicIDs) > 0 {
	// queryProcessors = append(queryProcessors, repository.Filter("topic_programming_questions.`topic_id` IS NULL OR "+
	// 	"topic_programming_questions.`topic_id` != ?", topicID))

	// queryProcessors = append(queryProcessors, repository.Filter("topic_programming_questions.`topic_id` IN (?)", topicIDs))
	// }

	queryProcessors = append(queryProcessors,
		repository.FilterWithOperator(columnNames, conditions, operators, values))

	return queryProcessors
}

// addSearchQueries will append search queries from queryParams to queryProcessor
func (service *ProgrammingQuestionService) addSearchQueries(requestForm url.Values) []repository.QueryProcessor {

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	var queryProcessors []repository.QueryProcessor
	// fmt.Println("****formvalue", requestForm["programmingConcept"])

	// Label.
	if _, ok := requestForm["label"]; ok {
		util.AddToSlice("programming_questions.`label`", "LIKE ?", "AND", "%"+requestForm.Get("label")+"%", &columnNames, &conditions, &operators, &values)
	}

	// Is active.
	if isActive, ok := requestForm["isActive"]; ok {
		util.AddToSlice("programming_questions.`is_active`", "=?", "AND", isActive, &columnNames, &conditions, &operators, &values)
	}

	// // Programming type id.
	// if _, ok := requestForm["programmingTypeID"]; ok {
	// 	util.AddToSlice("programming_question_type_id", "= ?", "AND", requestForm.Get("programmingTypeID"), &columnNames, &conditions, &operators, &values)
	// }

	// Talent id.
	if _, ok := requestForm["talentID"]; ok {
		util.AddToSlice("`talent_id`", "= ?", "AND", requestForm.Get("talentID"), &columnNames, &conditions, &operators, &values)
	}

	// Programming langauge id.
	if _, ok := requestForm["programmingLanguageID"]; ok {
		util.AddToSlice("programming_questions.`programming_language_id`", "= ?", "AND", requestForm.Get("programmingLanguageID"), &columnNames, &conditions, &operators, &values)
	}

	//Programming concept.
	if len(requestForm.Get("programmingConcept")) > 0 {
		queryProcessors = append(queryProcessors,
			repository.Join("INNER JOIN programming_questions_programming_concepts ON programming_questions.`id` = programming_questions_programming_concepts.`programming_question_id`"))
		util.AddToSlice("programming_questions_programming_concepts.`programming_concept_id`", "IN (?)", "AND",
			requestForm["programmingConcept"], &columnNames, &conditions, &operators, &values)
	}

	// Programming type.
	if programmingType, ok := requestForm["programmingType"]; ok {
		queryProcessors = append(queryProcessors,
			repository.Join("JOIN programming_questions_programming_question_types ON programming_questions.`id` = programming_questions_programming_question_types.`programming_question_id`"),
			repository.Join("JOIN programming_question_types ON programming_questions_programming_question_types.`programming_question_type_id` = programming_question_types.`id`"),
			repository.Filter("programming_question_types.`deleted_at` IS NULL"))
		util.AddToSlice("programming_question_types.`programming_type`", "=?", "AND", programmingType,
			&columnNames, &conditions, &operators, &values)
	}

	// // Programming concept.
	// if _, ok := requestForm["programmingConceptID"]; ok {
	// 	queryProcesors = append(queryProcesors,
	// 		repository.Join("JOIN programming_concepts_programming_questions ON programming_concepts_programming_questions.`programming_question_id` = programming_questions.`id`"),
	// 	)

	// 	util.AddToSlice("programming_concepts_programming_questions.`programming_concept_id`", "=?", "AND", requestForm.Get("programmingConceptID"),
	// 		&columnNames, &conditions, &operators, &values)
	// }

	// Add all filters
	queryProcessors = append(queryProcessors,
		repository.FilterWithOperator(columnNames, conditions, operators, values))

	return queryProcessors
}
