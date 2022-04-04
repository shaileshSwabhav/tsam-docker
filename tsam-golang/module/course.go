package module

import (
	"github.com/techlabs/swabhav/tsam"
	coursecon "github.com/techlabs/swabhav/tsam/course/controller"
	courseser "github.com/techlabs/swabhav/tsam/course/service"
	"github.com/techlabs/swabhav/tsam/repository"
)

func registerCourseRoutes(app *tsam.App, repository repository.Repository) {
	defer app.WG.Done()

	// Course Module
	courseService := courseser.NewCourseService(app.DB, repository)
	courseController := coursecon.NewCourseController(courseService, app.Log, app.Auth)
	courseSessionResourceService := courseser.NewModuleResourceService(app.DB, repository)
	courseSessionResourceController := coursecon.NewModuleResourceController(courseSessionResourceService, app.Log, app.Auth)
	courseAssessmentService := courseser.NewCourseTechnicalAssessmentService(app.DB, repository)
	courseAssessmentController := coursecon.NewAssessmentController(courseAssessmentService, app.Log, app.Auth)
	courseProgrammingAssignmentService := courseser.NewCourseTopicQuestionService(app.DB, repository)
	courseProgrammingAssignmentController := coursecon.NewCourseTopicQuestionController(courseProgrammingAssignmentService,
		app.Log, app.Auth)
	courseModuleService := courseser.NewCourseModuleService(app.DB, repository)
	courseModuleController := coursecon.NewCourseModuleController(courseModuleService, app.Log, app.Auth)
	moduleService := courseser.NewModuleService(app.DB, repository)
	moduleController := coursecon.NewModuleController(moduleService, app.Log, app.Auth)
	topicService := courseser.NewTopicService(app.DB, repository)
	topicController := coursecon.NewTopicController(topicService, app.Log, app.Auth)
	courseTopicConceptService := courseser.NewCourseTopicConceptService(app.DB, repository)
	courseTopicConceptController := coursecon.NewCourseTopicConceptController(courseTopicConceptService, app.Log, app.Auth)

	app.RegisterControllerRoutes([]tsam.Controller{
		courseController, moduleController, courseSessionResourceController, topicController,
		courseTopicConceptController, courseAssessmentController,
		courseProgrammingAssignmentController, courseModuleController,
	})

	// app.RouterRegister([]tsam.RouterSpecifier{
	// 	courseController, courseSessionController, courseSessionResourceController,
	// 	courseAssessmentController, courseProgrammingAssignmentController, courseModuleController,
	// })
}
