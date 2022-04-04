package service

import (
	"fmt"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// ProblemOfTheDayService provides methods to do different CRUD operations on problem of the day table.
type ProblemOfTheDayService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewProblemOfTheDayService returns a new instance Of ProblemOfTheDayService.
func NewProblemOfTheDayService(db *gorm.DB, repository repository.Repository) *ProblemOfTheDayService {
	return &ProblemOfTheDayService{
		DB:         db,
		Repository: repository,
	}
}

// GetProblemsOfTheDay gets problems of the day by date.
func (service *ProblemOfTheDayService) GetProblemsOfTheDay(tenantID uuid.UUID,
	questions *[]programming.QuestionProblemOfTheDayDTO) error {

	// ******************************* Automatic API Testing ****************************************

	// fmt.Println("Go Tickers Tutorial")
	// // `tick` every 1 second.
	// ticker := time.NewTicker(1 * time.Second)
	
	// // for every `tick` that our `ticker`
	// // emits, we print `tock`
	// for _ = range ticker.C {
	// 	fmt.Println("tock")
	// }


	// ******************************************************************************************************

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Print today's date and time.
	currentDate := time.Now().Format("2006-01-02")
	fmt.Println("*********************************Todays date and time*************************************")
	fmt.Println("Today's date: ", currentDate)
	fmt.Println("Today's date and time: ", time.Now())
	fmt.Println("****************************************************************************************")

	// Check if there are any problem of the day for today's date.
	err = service.Repository.GetAll(uow, questions,
		repository.Join("JOIN problem_of_the_day ON programming_questions.`id` = problem_of_the_day.`programming_question_id`"),
		repository.Filter("programming_questions.`tenant_id`=?", tenantID),
		repository.Filter("problem_of_the_day.`tenant_id`=?", tenantID),
		repository.Filter("date=?", currentDate),
	)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// fmt.Println("************************Records already in problem of the day:", questions)

	// If questions are already present for today's date then send them.
	if len(*questions) > 0 {

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
				repository.Filter("programming_question_talent_answers.`is_correct` = 1"),
				repository.GroupBy("`talent_id`"))
			if err != nil {
				uow.RollBack()
				return err
			}
			(*questions)[index].SolvedByCount = uint(totalCount)
		}

		return nil
	}

	// If no questions by today's date then add four questions by today's date to problem of the day table.
	// Get four questions whose entry is not in problem-of-the-day table.
	err = service.Repository.GetAll(uow, &questions,
		repository.Table("programming_questions"),
		repository.Join("LEFT JOIN problem_of_the_day ON programming_questions.`id` = problem_of_the_day.`programming_question_id`"),
		repository.Join("JOIN programming_question_types ON programming_question_types.`id` = programming_questions.`programming_question_type_id`"),
		repository.Filter("programming_questions.`tenant_id`=?", tenantID),
		repository.Filter("programming_question_types.`programming_type`=?", "Problem of the day"),
		repository.Filter("date IS NULL"),
		repository.Limit(4),
	)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// fmt.Println("************************Records taken for problem of the day:", problemOfTheDayQuestions)

	// Add 4 question IDs for today's date problem of the day.
	for _, qusetion := range *questions {

		// Create bucket for problem of the day.
		problemOfTheDay := programming.ProblemOfTheDay{}
		problemOfTheDay.ID = util.GenerateUUID()
		problemOfTheDay.Date = currentDate
		problemOfTheDay.TenantID = tenantID
		problemOfTheDay.ProgrammingQuestionID = qusetion.ID

		// Add problem of the day to database.
		err = service.Repository.Add(uow, problemOfTheDay)
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return err
		}
	}

	return nil
}

// GetProblemsOfThePreviousDays gets problems of the previous days by date.
func (service *ProblemOfTheDayService) GetProblemsOfThePreviousDays(tenantID uuid.UUID,
	questions *[]programming.QuestionProblemOfTheDayDTO, form url.Values) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if there are any problems of the day for previous days' date.
	err = service.Repository.GetAll(uow, questions,
		repository.Join("JOIN problem_of_the_day ON programming_questions.`id` = problem_of_the_day.`programming_question_id`"),
		repository.Filter("programming_questions.`tenant_id`=?", tenantID),
		repository.Filter("problem_of_the_day.`tenant_id`=?", tenantID),
		service.addSearchQueries(form),
	)
	if err != nil {
		log.NewLogger().Error(err.Error())
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
			repository.Filter("programming_questions.`id` = ?", (*questions)[index].ID))
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
			repository.Filter("programming_questions.`score` = programming_question_talent_answers.`score`"))
		if err != nil {
			uow.RollBack()
			return err
		}
		(*questions)[index].SolvedByCount = uint(totalCount)
	}

	return nil
}

// GetLeaderBoardPotd gets performers for problem ff the day by ranking and score.
func (service *ProblemOfTheDayService) GetLeaderBoardPotd(tenantID, talentID uuid.UUID,
	leaderBoard *general.LeaderBoradDTO, form url.Values) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create bucket for multiple perfomers.
	performers := []general.Performer{}

	// Create query precessors for sub query one.
	var queryProcessorsForSubQueryOne []repository.QueryProcessor
	queryProcessorsForSubQueryOne = append(queryProcessorsForSubQueryOne,
		repository.Select("MAX(programming_question_talent_answers.`score`) as total_score_question_wise, "+
			" programming_questions.`id` as programming_question_id, talents.`first_name` as first_name, "+
			"talents.`last_name` as last_name, talents.`image` as image, talents.`id` as talents_id "),
		repository.Table("programming_questions"),
		repository.Join("JOIN programming_question_types ON programming_questions.`programming_question_type_id` = programming_question_types.`id`"),
		repository.Join("JOIN programming_question_talent_answers ON programming_question_talent_answers.`programming_question_id` = programming_questions.`id`"),
		repository.Join("JOIN talents ON programming_question_talent_answers.`talent_id` = talents.`id`"),
		repository.Join("JOIN problem_of_the_day ON programming_questions.`id` = problem_of_the_day.`programming_question_id`"),
		service.addSearchQueries(form),
		repository.Filter("programming_questions.`tenant_id`=?", tenantID),
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("programming_question_types.`deleted_at` IS NULL"),
		repository.Filter("programming_questions.`deleted_at` IS NULL"),
		repository.Filter("programming_question_talent_answers.`deleted_at` IS NULL"),
		repository.Filter("programming_question_types.`programming_type` =?", "Problem of the day"),
		repository.GroupBy("programming_question_id, talents.`id`"))

	// Create query expression for sub query one.
	subQueryOne, err := service.Repository.SubQuery(uow, programming.ProgrammingQuestion{}, queryProcessorsForSubQueryOne...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Create query expression for query number two.
	subQueryTwo, err := service.Repository.SubQuery(uow, programming.ProgrammingQuestion{},
		repository.RawQuery("SELECT SUM(total_score_question_wise) as total_score, subone.* from ? as subone GROUP BY talents_id LIMIT 9", subQueryOne),
	)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Create query expression for query number three.
	subQueryThree, err := service.Repository.SubQuery(uow, programming.ProgrammingQuestion{},
		repository.RawQuery("SELECT subtwo.*, DENSE_RANK() OVER(ORDER BY  total_score desc) as `rank` from ? as subtwo", subQueryTwo),
	)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Get all performers.
	if err := service.Repository.GetAll(uow, &performers,
		repository.RawQuery("SELECT * FROM ? as subthree where talents_id != ? AND total_score > 0", subQueryThree, talentID)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Create bucket for performer.
	performer := general.Performer{}

	// Get self performer.
	isRecordNotFound := false
	if err := service.Repository.GetAll(uow, &performer,
		repository.RawQuery("SELECT * FROM ? as subthree where talents_id = ?", subQueryThree, talentID)); err != nil {
		if err == gorm.ErrRecordNotFound {
			isRecordNotFound = true
		} else {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Internal Server error")
		}
	}

	// If self performer record not found.
	if isRecordNotFound {

		// Create query precessors for talent sub query one.
		var queryProcessorsFortalentSubQueryOne []repository.QueryProcessor
		queryProcessorsFortalentSubQueryOne = append(queryProcessorsFortalentSubQueryOne,
			repository.Select("MAX(programming_question_talent_answers.`score`) as total_score_question_wise, "+
				" programming_questions.`id` as programming_question_id, talents.`first_name` as first_name, "+
				"talents.`last_name` as last_name, talents.`image` as image, talents.`id` as talents_id "),
			repository.Table("programming_questions"),
			repository.Join("JOIN programming_question_types ON programming_questions.`programming_question_type_id` = programming_question_types.`id`"),
			repository.Join("JOIN programming_question_talent_answers ON programming_question_talent_answers.`programming_question_id` = programming_questions.`id`"),
			repository.Join("JOIN talents ON programming_question_talent_answers.`talent_id` = talents.`id`"),
			repository.Join("JOIN problem_of_the_day ON programming_questions.`id` = problem_of_the_day.`programming_question_id`"),
			service.addSearchQueries(form),
			repository.Filter("programming_questions.`tenant_id`=?", tenantID),
			repository.Filter("talents.`deleted_at` IS NULL"),
			repository.Filter("programming_question_types.`deleted_at` IS NULL"),
			repository.Filter("programming_questions.`deleted_at` IS NULL"),
			repository.Filter("programming_question_talent_answers.`deleted_at` IS NULL"),
			repository.Filter("programming_question_types.`programming_type` =?", "Problem of the day"),
			repository.GroupBy("programming_question_id, talents.`id`"))

		// Create query expression for talent sub query one.
		talentSubQueryOne, err := service.Repository.SubQuery(uow, programming.ProgrammingQuestion{}, queryProcessorsForSubQueryOne...)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

		// Create query expression for talent query number two.
		talentSubQueryTwo, err := service.Repository.SubQuery(uow, programming.ProgrammingQuestion{},
			repository.RawQuery("SELECT SUM(total_score_question_wise) as total_score, subone.* from ? as subone GROUP BY talents_id", talentSubQueryOne),
		)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

		// Create query expression for talent query number three.
		talentSubQueryThree, err := service.Repository.SubQuery(uow, programming.ProgrammingQuestion{},
			repository.RawQuery("SELECT subtwo.*, DENSE_RANK() OVER(ORDER BY  total_score desc) as `ranking` from ? as subtwo", talentSubQueryTwo),
		)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

		// Get performer with least rank.
		if err := service.Repository.GetAll(uow, &performer,
			repository.RawQuery("SELECT MAX(IF(total_score = 0, `ranking`, `ranking` + 1)) as `rank` FROM ? as subthree", talentSubQueryThree)); err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

		// Get talent deatils.
		// Create bucket for talent.
		tempTalent := talent.Talent{}
		if err := service.Repository.GetForTenant(uow, tenantID, talentID, &tempTalent,
			repository.Select("first_name, last_name, image")); err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

		// Give talent details to performer.
		performer.FirstName = tempTalent.FirstName
		performer.LastName = tempTalent.LastName
		performer.Image = tempTalent.Image

		// Give performer to leader board.
		leaderBoard.SelfPerformer = performer

		// Give performers to leader board.
		leaderBoard.AllPerformers = performers
	}

	// If self performer record found.
	if !isRecordNotFound {
		// Give performer to leader board.
		leaderBoard.SelfPerformer = performer

		// Give performers to leader board.
		leaderBoard.AllPerformers = performers
	}

	return nil
}

// GetLeaderBoard gets performers for all questions by ranking and score.
func (service *ProblemOfTheDayService) GetLeaderBoard(tenantID, talentID uuid.UUID,
	leaderBoard *general.LeaderBoradDTO) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create bucket for multiple perfomers.
	performers := []general.Performer{}

	// Create query precessors for sub query one.
	var queryProcessorsForSubQueryOne []repository.QueryProcessor
	queryProcessorsForSubQueryOne = append(queryProcessorsForSubQueryOne,
		repository.Select("MAX(programming_question_talent_answers.`score`) as total_score_question_wise, "+
			" programming_questions.`id` as programming_question_id, talents.`first_name` as first_name, "+
			"talents.`last_name` as last_name, talents.`image` as image, talents.`id` as talents_id "),
		repository.Table("programming_questions"),
		repository.Join("JOIN programming_question_talent_answers ON programming_question_talent_answers.`programming_question_id` = programming_questions.`id`"),
		repository.Join("JOIN talents ON programming_question_talent_answers.`talent_id` = talents.`id`"),
		repository.Filter("programming_questions.`tenant_id`=?", tenantID),
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("programming_questions.`deleted_at` IS NULL"),
		repository.Filter("programming_question_talent_answers.`deleted_at` IS NULL"),
		repository.GroupBy("programming_question_id, talents.`id`"))

	// Create query expression for sub query one.
	subQueryOne, err := service.Repository.SubQuery(uow, programming.ProgrammingQuestion{}, queryProcessorsForSubQueryOne...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Create query expression for query number two.
	subQueryTwo, err := service.Repository.SubQuery(uow, programming.ProgrammingQuestion{},
		repository.RawQuery("SELECT SUM(total_score_question_wise) as total_score, subone.* from ? as subone GROUP BY talents_id LIMIT 9", subQueryOne),
	)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Create query expression for query number three.
	subQueryThree, err := service.Repository.SubQuery(uow, programming.ProgrammingQuestion{},
		repository.RawQuery("SELECT subtwo.*, DENSE_RANK() OVER(ORDER BY  total_score desc) as `rank` from ? as subtwo", subQueryTwo),
	)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Get all performers.
	if err := service.Repository.GetAll(uow, &performers,
		repository.RawQuery("SELECT * FROM ? as subthree where talents_id != ? AND total_score > 0", subQueryThree, talentID)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Create bucket for performer.
	performer := general.Performer{}

	// Get self performer.
	isRecordNotFound := false
	if err := service.Repository.GetAll(uow, &performer,
		repository.RawQuery("SELECT * FROM ? as subthree where talents_id = ?", subQueryThree, talentID)); err != nil {
		if err == gorm.ErrRecordNotFound {
			isRecordNotFound = true
		} else {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Internal Server error")
		}
	}

	// If self performer record not found.
	if isRecordNotFound {

		// Create query precessors for talent sub query one.
		var queryProcessorsFortalentSubQueryOne []repository.QueryProcessor
		queryProcessorsFortalentSubQueryOne = append(queryProcessorsFortalentSubQueryOne,
			repository.Select("MAX(programming_question_talent_answers.`score`) as total_score_question_wise, "+
				" programming_questions.`id` as programming_question_id, talents.`first_name` as first_name, "+
				"talents.`last_name` as last_name, talents.`image` as image, talents.`id` as talents_id "),
			repository.Table("programming_questions"),
			repository.Join("JOIN programming_question_talent_answers ON programming_question_talent_answers.`programming_question_id` = programming_questions.`id`"),
			repository.Join("JOIN talents ON programming_question_talent_answers.`talent_id` = talents.`id`"),
			repository.Filter("programming_questions.`tenant_id`=?", tenantID),
			repository.Filter("talents.`deleted_at` IS NULL"),
			repository.Filter("programming_questions.`deleted_at` IS NULL"),
			repository.Filter("programming_question_talent_answers.`deleted_at` IS NULL"),
			repository.GroupBy("programming_question_id, talents.`id`"))

		// Create query expression for talent sub query one.
		talentSubQueryOne, err := service.Repository.SubQuery(uow, programming.ProgrammingQuestion{}, queryProcessorsForSubQueryOne...)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

		// Create query expression for talent query number two.
		talentSubQueryTwo, err := service.Repository.SubQuery(uow, programming.ProgrammingQuestion{},
			repository.RawQuery("SELECT SUM(total_score_question_wise) as total_score, subone.* from ? as subone GROUP BY talents_id", talentSubQueryOne),
		)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

		// Create query expression for talent query number three.
		talentSubQueryThree, err := service.Repository.SubQuery(uow, programming.ProgrammingQuestion{},
			repository.RawQuery("SELECT subtwo.*, DENSE_RANK() OVER(ORDER BY  total_score desc) as `ranking` from ? as subtwo", talentSubQueryTwo),
		)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

		// Get performer with least rank.
		if err := service.Repository.GetAll(uow, &performer,
			repository.RawQuery("SELECT MAX(IF(total_score = 0, `ranking`, `ranking` + 1)) as `rank` FROM ? as subthree", talentSubQueryThree)); err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

		// Get talent deatils.
		// Create bucket for talent.
		tempTalent := talent.Talent{}
		if err := service.Repository.GetForTenant(uow, tenantID, talentID, &tempTalent,
			repository.Select("first_name, last_name, image")); err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

		// Give talent details to performer.
		performer.FirstName = tempTalent.FirstName
		performer.LastName = tempTalent.LastName
		performer.Image = tempTalent.Image

		// Give performer to leader board.
		leaderBoard.SelfPerformer = performer

		// Give performers to leader board.
		leaderBoard.AllPerformers = performers
	}

	// If self performer record found.
	if !isRecordNotFound {
		// Give performer to leader board.
		leaderBoard.SelfPerformer = performer

		// Give performers to leader board.
		leaderBoard.AllPerformers = performers
	}

	return nil
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *ProblemOfTheDayService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// addSearchQueries will append search queries from queryParams to queryProcessor
func (service *ProblemOfTheDayService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if searchDate, ok := requestForm["searchDate"]; ok {
		util.AddToSlice("`date`", "=?", "AND", searchDate, &columnNames, &conditions, &operators, &values)
	} else {
		// Get today's date.
		currentDate := time.Now().Format("2006-01-02")
		util.AddToSlice("`date`", "=?", "AND", currentDate, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}