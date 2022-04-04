package module

import (
	"github.com/techlabs/swabhav/tsam/callreport/callreportcntlr"
	"github.com/techlabs/swabhav/tsam/callreport/callreportsvc"
	reportCon "github.com/techlabs/swabhav/tsam/report/controller"
	reportSer "github.com/techlabs/swabhav/tsam/report/service"
	testCon "github.com/techlabs/swabhav/tsam/test/controller"
	testSer "github.com/techlabs/swabhav/tsam/test/service"

	companycon "github.com/techlabs/swabhav/tsam/company/controller"
	companyser "github.com/techlabs/swabhav/tsam/company/service"

	collegecon "github.com/techlabs/swabhav/tsam/college/controller"
	collegeser "github.com/techlabs/swabhav/tsam/college/service"

	"github.com/techlabs/swabhav/tsam"
	admincon "github.com/techlabs/swabhav/tsam/administration/controller"
	adminserv "github.com/techlabs/swabhav/tsam/administration/service"
	batchcon "github.com/techlabs/swabhav/tsam/batch/controller"
	batchser "github.com/techlabs/swabhav/tsam/batch/service"
	blogcon "github.com/techlabs/swabhav/tsam/blog/controller"
	blogser "github.com/techlabs/swabhav/tsam/blog/service"
	facultycon "github.com/techlabs/swabhav/tsam/faculty/controller"
	facultyser "github.com/techlabs/swabhav/tsam/faculty/service"
	generalcon "github.com/techlabs/swabhav/tsam/general/controller"
	generalser "github.com/techlabs/swabhav/tsam/general/service"
	programmingcon "github.com/techlabs/swabhav/tsam/programming/controller"
	programmingser "github.com/techlabs/swabhav/tsam/programming/service"
	"github.com/techlabs/swabhav/tsam/repository"
	resourcecon "github.com/techlabs/swabhav/tsam/resource/controller"
	resourceser "github.com/techlabs/swabhav/tsam/resource/service"
	tcon "github.com/techlabs/swabhav/tsam/talent/controller"
	tser "github.com/techlabs/swabhav/tsam/talent/service"

	androidUSerController "github.com/techlabs/swabhav/tsam/android/controller"
	androidUserService "github.com/techlabs/swabhav/tsam/android/service"

	communitycon "github.com/techlabs/swabhav/tsam/community/controller"
	communityser "github.com/techlabs/swabhav/tsam/community/service"

	notiController "github.com/techlabs/swabhav/tsam/notification_demo/controller"
	notiService "github.com/techlabs/swabhav/tsam/notification_demo/service"

	emailController "github.com/techlabs/swabhav/tsam/homePage/contact/controller"
	emailService "github.com/techlabs/swabhav/tsam/homePage/contact/service"
)

// CreateRouterInstance will create and register the routes of all routers.
func CreateRouterInstance(app *tsam.App, repository repository.Repository) {
	// app.WG.Add(1)
	// go registerCommunityRoutes(app, repository)

	log := app.Log
	// General Module
	countryService := generalser.NewCountryService(app.DB, repository)
	countryController := generalcon.NewCountryController(countryService, log, app.Auth)
	stateService := generalser.NewStateService(app.DB, repository)
	stateController := generalcon.NewStateController(stateService, log, app.Auth)
	cityService := generalser.NewCityService(app.DB, repository)
	cityController := generalcon.NewCityController(cityService, log, app.Auth)
	technologyService := generalser.NewTechnologyService(app.DB, repository)
	technologyController := generalcon.NewTechnologyController(technologyService, log, app.Auth)
	degreeService := generalser.NewDegreeService(app.DB, repository)
	degreeController := generalcon.NewDegreeController(degreeService, log, app.Auth)
	GeneralTypeService := generalser.NewGeneralTypeService(app.DB, repository)
	GeneralTypeController := generalcon.NewGeneralTypeController(GeneralTypeService, log, app.Auth)
	roleService := generalser.NewRoleService(app.DB, repository)
	roleController := generalcon.NewRoleController(roleService)
	credentialService := generalser.NewCredentialService(app.DB, repository)
	credentialController := generalcon.NewCredentialController(credentialService, app.Auth)
	menuService := generalser.NewMenuService(app.DB, repository)
	menuController := generalcon.NewMenuController(menuService, log, app.Auth)
	tenantService := generalser.NewTenantService(app.DB, repository)
	tenantController := generalcon.NewTenantController(tenantService)
	designationService := generalser.NewDesignationService(app.DB, repository)
	designationController := generalcon.NewDesignationController(designationService, log, app.Auth)

	purposeService := generalser.NewPurposeService(app.DB, repository)
	purposeController := generalcon.NewPurposeController(purposeService, log, app.Auth)

	outcomeService := generalser.NewOutcomeService(app.DB, repository)
	outcomeController := generalcon.NewOutcomeController(outcomeService, log, app.Auth)

	sourceService := generalser.NewSourceService(app.DB, repository)
	sourceController := generalcon.NewSourceController(sourceService, log, app.Auth)

	examinationService := generalser.NewExaminationService(app.DB, repository)
	examinationController := generalcon.NewExaminationController(examinationService, log, app.Auth)

	userService := generalser.NewUserService(app.DB, repository)
	userController := generalcon.NewUserController(userService)

	dayService := generalser.NewDayService(app.DB, repository)
	dayController := generalcon.NewDayController(dayService)
	feedbackQuestionService := generalser.NewFeedbackQuestionService(app.DB, repository)
	feedbackQuestionController := generalcon.NewFeedbackQuestionController(feedbackQuestionService)
	feedbackOptionSerivce := generalser.NewFeedbackOptionService(app.DB, repository)
	feedbackOptionController := generalcon.NewFeedbackOptionController(feedbackOptionSerivce)

	departmentService := generalser.NewDepartmentService(app.DB, repository)
	departmentController := generalcon.NewDepartmentController(departmentService, log, app.Auth)
	projectService := generalser.NewSwabhavProjectService(app.DB, repository)
	projectController := generalcon.NewSwabhavProjectController(projectService, log, app.Auth)

	employeeService := generalser.NewEmployeeService(app.DB, repository)
	employeeController := generalcon.NewEmployeeController(employeeService)
	careerObjectiveService := generalser.NewCareerObjectiveService(app.DB, repository)
	careerObjectiveController := generalcon.NewCareerObjectiveController(careerObjectiveService, log, app.Auth)

	feelingService := generalser.NewFeelingService(app.DB, repository)
	feelingController := generalcon.NewFeelingController(feelingService, log, app.Auth)

	feedbackQuestionGroupService := generalser.NewFeedbackQuestionGroupService(app.DB, repository)
	feedbackQuestionGroupController := generalcon.NewFeedbackQuestionGroupController(feedbackQuestionGroupService)

	targetCommunityFunctionService := generalser.NewTargetCommunityFunctionService(app.DB, repository)
	targetCommunityFunctionController := generalcon.NewTargetCommunityFunctionController(targetCommunityFunctionService, log, app.Auth)

	salaryTrendService := generalser.NewSalaryTrendService(app.DB, repository)
	salaryTrendController := generalcon.NewSalaryTrendController(salaryTrendService)

	//android User
	androidUserSer := androidUserService.NewAndroidUserService(app.DB, repository)
	androidUSerControl := androidUSerController.NewAndroidUserController(androidUserSer, log, app.Auth)

	// programming
	programmingAssignmentService := programmingser.NewProgrammingAssignmentService(app.DB, repository)
	programmingAssignmentController := programmingcon.NewProgrammingAssignmentController(programmingAssignmentService,
		log, app.Auth)

	programmingQuestionService := programmingser.NewProgrammingQuestionService(app.DB, repository)
	programmingQuestionController := programmingcon.NewProgrammingQuestionController(programmingQuestionService)
	programmingQuestionTypeService := programmingser.NewProgrammingQuestionTypeService(app.DB, repository)
	programmingQuestionTypeController := programmingcon.NewProgrammingQuestionTypeController(programmingQuestionTypeService)
	problemOfTheDayService := programmingser.NewProblemOfTheDayService(app.DB, repository)
	problemOfTheDayController := programmingcon.NewProblemOfTheDayController(problemOfTheDayService)
	programmingSolutionService := programmingser.NewProgrammingQuestionTalentAnswerService(app.DB, repository)
	programmingSolutionController := programmingcon.NewProgrammingQuestionTalentAnswerController(programmingSolutionService)
	programmingConceptService := programmingser.NewProgrammingConceptService(app.DB, repository)
	programmingConceptController := programmingcon.NewProgrammingConceptController(programmingConceptService, log, app.Auth)
	moduleProgrammingConceptService := programmingser.NewModuleProgrammingConceptService(app.DB, repository)
	moduleProgrammingConceptController := programmingcon.NewModuleProgrammingConceptController(moduleProgrammingConceptService, log, app.Auth)
	programmingLanguageService := programmingser.NewProgrammingLanguageService(app.DB, repository)
	programmingLanguageController := programmingcon.NewProgrammingLanguageController(programmingLanguageService)
	programmingQuestionSolutionService := programmingser.NewProgrammingQuestionSolutionService(app.DB, repository)
	programmingQuestionSolutionController := programmingcon.NewProgrammingQuestionSolutionController(programmingQuestionSolutionService)
	programmingQuestionTestCaseService := programmingser.NewProgrammingQuestionTestCaseService(app.DB, repository)
	programmingQuestionTestCaseController := programmingcon.NewProgrammingQuestionTestCaseController(programmingQuestionTestCaseService)
	languageCompilerService := programmingser.NewLanguageCompilerService(app.DB, repository)
	languageCompilerController := programmingcon.NewLanguageCompilerController(languageCompilerService)
	programmingProjectService := programmingser.NewProgrammingProjectService(app.DB, repository)
	programmingProjectController := programmingcon.NewProgrammingProjectController(programmingProjectService, log, app.Auth)
	conceptDashboardService := programmingser.NewConceptDashboardService(app.DB, repository)
	conceptDashboardController := programmingcon.NewConceptDashboardController(conceptDashboardService, log, app.Auth)

	// blog
	blogTopicService := blogser.NewBlogTopicService(app.DB, repository)
	blogTopicController := blogcon.NewBlogTopicController(blogTopicService)
	blogService := blogser.NewBlogService(app.DB, repository)
	blogController := blogcon.NewBlogController(blogService)
	blogReplyService := blogser.NewBlogReplyService(app.DB, repository)
	blogReplyController := blogcon.NewBlogReplyController(blogReplyService)
	blogReactionService := blogser.NewBlogReactionService(app.DB, repository)
	blogReactionController := blogcon.NewBlogReactionController(blogReactionService)
	blogViewService := blogser.NewBlogViewService(app.DB, repository)
	blogViewController := blogcon.NewBlogViewController(blogViewService)

	// Timesheet
	// timesheetService := adminserv.TimesheetService(app.DB, repository)
	// timesheetController := admincon.TimesheetController(timesheetService)
	timesheetService := adminserv.NewTimesheetService(app.DB, repository)
	timesheetController := admincon.TimesheetController(timesheetService)

	// User
	salesPersonService := adminserv.NewSalesPersonService(app.DB, repository)
	salesPersonController := admincon.NewSalesPersonController(salesPersonService)
	salesPersonDashboardService := adminserv.NewSalesPersonDashboardService(app.DB, repository)
	salesPersonDashboardController := admincon.NewSalesPersonDashboardController(salesPersonDashboardService)
	adminDashboardService := adminserv.NewAdminDashboardService(app.DB, repository)
	adminDashboardController := admincon.NewAdminDashboardController(adminDashboardService)

	// Management
	managementService := adminserv.NewManagementService(app.DB, repository)
	managementController := admincon.NewManagementController(managementService, app.Auth)

	// Event
	eventService := adminserv.NewSwabhavEventService(app.DB, repository)
	eventController := admincon.NewSwabhavEventController(eventService)

	// Faculty Module
	facultyService := facultyser.NewFacultyService(app.DB, repository)
	facultyController := facultycon.NewFacultyController(facultyService, log, app.Auth)
	facultyDashboardService := facultyser.NewFacultyDashboardService(app.DB, repository)
	facultyDashboardController := facultycon.NewFacultyDashboardController(facultyDashboardService)
	facultyAssessmentService := facultyser.NewFacultyAssessmentService(app.DB, repository)
	facultyAssessmentController := facultycon.NewFacultyAssessmentController(facultyAssessmentService, log, app.Auth)

	// Talent Module
	enquiryService := tser.NewEnquiryService(app.DB, repository)
	enquiryController := tcon.NewEnquiryController(enquiryService)
	enquiryDashboardService := tser.NewEnquiryDashboardService(app.DB, repository)
	enquiryDashboardController := tcon.NewEnquiryDashboardController(enquiryDashboardService)
	talentService := tser.NewTalentService(app.DB, repository)
	talentController := tcon.NewTalentController(talentService)
	talentDashboardService := tser.NewTalentDashboardService(app.DB, repository)
	talentDashboardController := tcon.NewTalentDashboardController(talentDashboardService)
	talentCallRecordService := tser.NewTalentCallRecordService(app.DB, repository)
	talentCallRecordController := tcon.NewTalentCallRecordController(talentCallRecordService)
	specializationService := generalser.NewSpecializationService(app.DB, repository)
	specializationController := generalcon.NewSpecializationController(specializationService, log, app.Auth)
	enquiryCallRecordService := tser.NewEnquiryCallRecordService(app.DB, repository)
	enquiryCallRecordController := tcon.NewEnquiryCallRecordController(enquiryCallRecordService)
	lifetimeValueService := tser.NewTalentLifetimeValueService(app.DB, repository)
	lifetimeValueController := tcon.NewTalentLifetimeValueController(lifetimeValueService)
	interviewRoundService := tser.NewInterviewRoundService(app.DB, repository)
	interviewRoundController := tcon.NewInterviewRoundController(interviewRoundService)
	interviewScheduleService := tser.NewInterviewScheduleService(app.DB, repository)
	interviewScheduleController := tcon.NewInterviewShceduleController(interviewScheduleService)
	interviewService := tser.NewInterviewService(app.DB, repository)
	interviewController := tcon.NewInterviewController(interviewService)
	nextActionService := tser.NewTalentNextActionService(app.DB, repository)
	nextActionController := tcon.NewTalentNextActionController(nextActionService)
	nextActionTypeService := tser.NewNextActionTypeService(app.DB, repository)
	nextActionTypeController := tcon.NewNextActionTypeController(nextActionTypeService)
	careerPlanService := tser.NewCareerPlanService(app.DB, repository)
	careerPlanController := tcon.NewCareerPlanController(careerPlanService)
	waitingListService := tser.NewWaitingListService(app.DB, repository)
	waitingListController := tcon.NewWaitingListController(waitingListService)
	waitingListReportService := tser.NewWaitingListReportService(app.DB, repository)
	waitingListReportController := tcon.NewWaitingListReportController(waitingListReportService)
	proSummaryReportService := tser.NewProfessionalSummaryReportService(app.DB, repository)
	proSummaryReportController := tcon.NewProfessionalSummaryReportController(proSummaryReportService)
	talentRegistrationService := tser.NewTalentEventRegistrationService(app.DB, repository)
	talentRegistrationController := tcon.NewTalentEventRegistrationController(talentRegistrationService)
	talentAssignmentSubmissionService := tser.NewTalentAssignmentSubmissionService(app.DB, repository)
	talentAssignmentSubmissionController := tcon.NewTalentAssignmentSubmissionController(talentAssignmentSubmissionService,
		log, app.Auth)
	talentProjectSubmissionService := tser.NewTalentProjectSubmissionService(app.DB, repository)
	talentProjectSubmissionController := tcon.NewTalentProjectSubmissionController(talentProjectSubmissionService, log, app.Auth)
	// report
	loginReportService := reportSer.NewLoginReportService(app.DB, repository)
	loginReportController := reportCon.NewLoginReportController(loginReportService)
	fresherSummaryService := reportSer.NewFresherSummaryService(app.DB, repository)
	fresherSummaryController := reportCon.NewFresherSummaryController(fresherSummaryService)
	packageSummaryService := reportSer.NewPackageSummaryService(app.DB, repository)
	packageSummaryController := reportCon.NewPackageSummaryController(packageSummaryService)
	facultyReportService := reportSer.NewFacultyReportService(app.DB, repository)
	facultyReportController := reportCon.NewFacultyReportController(facultyReportService)
	talentReportService := reportSer.NewTalentReportService(app.DB, repository)
	talentReportController := reportCon.NewTalentReportController(talentReportService)

	// Course Module
	// courseService := courseser.NewCourseService(app.DB, repository)
	// courseController := coursecon.NewCourseController(courseService, log, app.Auth)
	// // courseSessionService := courseser.NewCourseSessionService(app.DB, repository)
	// // courseSessionController := coursecon.NewCourseSessionController(courseSessionService, log, app.Auth)
	// courseSessionResourceService := courseser.NewModuleResourceService(app.DB, repository)
	// courseSessionResourceController := coursecon.NewModuleResourceController(courseSessionResourceService, log, app.Auth)
	// courseAssessmentService := courseser.NewCourseTechnicalAssessmentService(app.DB, repository)
	// courseAssessmentController := coursecon.NewAssessmentController(courseAssessmentService, log, app.Auth)
	// courseProgrammingAssignmentService := courseser.NewCourseTopicQuestionService(app.DB, repository)
	// courseProgrammingAssignmentController := coursecon.NewCourseTopicQuestionController(courseProgrammingAssignmentService,
	// 	log, app.Auth)
	// courseModuleService := courseser.NewCourseModuleService(app.DB, repository)
	// courseModuleController := coursecon.NewCourseModuleController(courseModuleService, log, app.Auth)
	// moduleService := courseser.NewModuleService(app.DB, repository)
	// moduleController := coursecon.NewModuleController(moduleService, log, app.Auth)
	// topicService := courseser.NewTopicService(app.DB, repository)
	// topicController := coursecon.NewTopicController(topicService, log, app.Auth)
	// courseTopicConceptService := courseser.NewCourseTopicConceptService(app.DB, repository)
	// courseTopicConceptController := coursecon.NewCourseTopicConceptController(courseTopicConceptService, log, app.Auth)

	// College Module
	collegeService := collegeser.NewCollegeService(app.DB, repository)
	collegeController := collegecon.NewCollegeController(collegeService)
	collegeDriveService := collegeser.NewCollegeBranchService(app.DB, repository)
	collegeDriveController := collegecon.NewCollegeBranchController(collegeDriveService)
	collegeCampusService := collegeser.NewCampusDriveService(app.DB, repository)
	collegeCampusController := collegecon.NewCampusDriveController(collegeCampusService)
	collegeDashboardService := collegeser.NewCollegeDashboardService(app.DB, repository)
	collegeDashboardController := collegecon.NewCollegeDashboardController(collegeDashboardService)
	universityService := generalser.NewUniversityService(app.DB, repository)
	universityController := generalcon.NewUniversityController(universityService, log, app.Auth)
	candidateService := collegeser.NewCandidateService(app.DB, repository)
	candidateController := collegecon.NewCandidateController(candidateService)
	speakerService := collegeser.NewSpeakerService(app.DB, repository)
	speakerController := collegecon.NewSpeakerController(speakerService)
	seminarService := collegeser.NewSeminarService(app.DB, repository)
	seminarController := collegecon.NewSeminarController(seminarService)
	seminarTopicService := collegeser.NewSeminarTopicService(app.DB, repository)
	seminarTopicController := collegecon.NewSeminarTopicController(seminarTopicService)
	studentSeminarService := collegeser.NewStudentService(app.DB, repository)
	studentSeminaController := collegecon.NewStudentController(studentSeminarService)

	// Company Module
	companyService := companyser.NewCompanyService(app.DB, repository)
	companyController := companycon.NewCompanyController(companyService, log, app.Auth)
	companyBranchService := companyser.NewCompanyBranchService(app.DB, repository)
	companyBranchController := companycon.NewCompanyBranchController(companyBranchService, log, app.Auth)
	companyEnquiryService := companyser.NewCompanyEnquiryService(app.DB, repository)
	companyEnquiryController := companycon.NewCompanyEnquiryController(companyEnquiryService, log, app.Auth)
	companyReqirementService := companyser.NewCompanyRequirementService(app.DB, repository)
	companyRequirementController := companycon.NewCompanyRequirementController(companyReqirementService, log, app.Auth)
	domainService := companyser.NewDomainService(app.DB, repository)
	domainController := companycon.NewDomainController(domainService)
	companyDashboardService := companyser.NewCompanyDashboardService(app.DB, repository)
	companyDashboardController := companycon.NewCompanyDashboardController(companyDashboardService)
	companyEnquiryCallRecordService := companyser.NewEnquiryCallRecordService(app.DB, repository)
	companyEnquiryCallRecordController := companycon.NewEnquiryCallRecordController(companyEnquiryCallRecordService, log, app.Auth)

	// batch Module
	batchService := batchser.NewBatchService(app.DB, repository)
	batchController := batchcon.NewBatchController(batchService, log, app.Auth)
	courseBatchDashboardService := batchser.NewCourseBatchDashboardService(app.DB, repository)
	courseBatchDashboardController := batchcon.NewCourseBatchDashboardController(courseBatchDashboardService)
	oldBatchSessionService := batchser.OldBatchSessionService(app.DB, repository)
	oldBatchSessionController := batchcon.OldBatchSessionController(oldBatchSessionService)
	facultyBatchSessionFeedbackService := batchser.NewFacultySessionFeedbackService(app.DB, repository)
	facultyBatchSessionFeedbackController := batchcon.NewFacultySessionFeedbackController(facultyBatchSessionFeedbackService, log, app.Auth)
	facultyBatchFeedbackService := batchser.NewFacultyFeedbackService(app.DB, repository)
	facultyBatchFeedbackController := batchcon.NewFacultyFeedbackController(facultyBatchFeedbackService)
	talentBatchSessionFeedbackService := batchser.NewTalentSessionFeedbackService(app.DB, repository)
	talentBatchSessionFeedbackController := batchcon.NewTalentSessionFeedbackController(talentBatchSessionFeedbackService, log, app.Auth)
	talentBatchFeedbackService := batchser.NewTalentFeedbackService(app.DB, repository)
	talentBatchFeedbackController := batchcon.NewTalentFeedbackController(talentBatchFeedbackService)
	sessionProgrammingAssignmentService := batchser.NewBatchTopicAssignmentService(app.DB, repository)
	sessionProgrammingAssignmentController := batchcon.NewBatchTopicAssignmentController(sessionProgrammingAssignmentService,
		log, app.Auth)
	batchSessionsTalentService := batchser.NewSessionBatchTalentService(app.DB, repository)
	batchSessionsTalentController := batchcon.NewBatchSessionsTalentController(batchSessionsTalentService, log, app.Auth)
	batchTalentService := batchser.NewBatchTalentService(app.DB, repository)
	batchTalentController := batchcon.NewBatchTalentController(batchTalentService, log, app.Auth)
	// batchSessionPlanService := batchser.NewBatchSessionTopicService(app.DB, repository)
	// batchSessionPlanController := batchcon.NewBatchSessionTopicController(batchSessionPlanService, log, app.Auth)
	batchSessionService := batchser.NewSessionService(app.DB, repository)
	batchSessionController := batchcon.NewSessionController(batchSessionService, log, app.Auth)
	batchProjectService := batchser.NewProjectService(app.DB, repository)
	batchProjectController := batchcon.NewBatchProjectController(batchProjectService, log, app.Auth)
	prerequisiteService := batchser.NewPrerequisiteService(app.DB, repository)
	prerequisiteController := batchcon.NewPrerequisiteController(prerequisiteService, log, app.Auth)
	batchModuleService := batchser.NewModuleService(app.DB, repository)
	batchModuleController := batchcon.NewModuleController(batchModuleService, log, app.Auth)

	// Test Management
	questionService := testSer.NewQuestionService(app.DB, repository)
	questionController := testCon.NewQuestionController(questionService)
	optionService := testSer.NewOptionService(app.DB, repository)
	optionController := testCon.NewOptionController(optionService)

	// Calling reports
	callReportService := callreportsvc.New(app.DB, repository)
	callReportController := callreportcntlr.New(callReportService)
	nextActionReportService := reportSer.New(app.DB, repository)
	nextActionReportController := reportCon.New(nextActionReportService)

	// Aha moment
	ahaMomentService := batchser.NewAhaMomentService(app.DB, repository)
	ahaMomentController := batchcon.NewAhaMomentController(ahaMomentService)

	// Target Community
	targetCommunityService := adminserv.NewTargetCommunityService(app.DB, repository)
	targetCommunityController := admincon.NewTargetCommunityController(targetCommunityService)

	// resource
	resourceService := resourceser.NewResourceService(app.DB, repository)
	resourceController := resourcecon.NewResourceController(resourceService, log, app.Auth)
	resourceDownloadService := resourceser.NewResourceDownloadService(app.DB, repository)
	resourceDownloadController := resourcecon.NewResourceDownloadController(resourceDownloadService)
	resourceLikeService := resourceser.NewResourceLikeService(app.DB, repository)
	resourceLikeController := resourcecon.NewResourceLikeController(resourceLikeService)

	// community
	replyService := communityser.NewReplyService(app.DB, repository)
	replyController := communitycon.NewReplyController(replyService, log, app.Auth)
	// reactionService := communityser.NewReactionService(app.DB, repository)
	// reactionController := communitycon.NewReactionController(reactionService, log, app.Auth)
	// commentService := communityser.NewCommentService(app.DB, repository)
	// commentController := communitycon.NewCommentController(commentService)
	// notificationService := communityser.NewNotificationService(app.DB, repository)
	// notificationController := communitycon.NewNotificationController(notificationService)
	// notificationTypeService := communityser.NewNotificationTypeService(app.DB, repository)
	// notificationTypeController := communitycon.NewNotificationTypeController(notificationTypeService)
	channelService := communityser.NewChannelService(app.DB, repository)
	channelController := communitycon.NewChannelController(channelService, log, app.Auth)

	//notification_test
	notiService := notiService.NewNotificationService(app.DB, repository)
	notificationsController := notiController.NewAndroidUserController(notiService, log, app.Auth)

	//email
	emailService := emailService.NewContactInfoService(app.DB, repository)
	emailController := emailController.NewContactInfoController(emailService, log, app.Auth)

	//project_rating
	projectRatingParameterService := batchser.NewProjectRatingService(app.DB, repository)
	projectRatingParameterController := batchcon.NewPrjectRatingController(projectRatingParameterService, log, app.Auth)
	_ = channelController
	routerSpecifier := []tsam.RouterSpecifier{
		replyController, batchModuleController,
		salesPersonController, lifetimeValueController, careerObjectiveController, careerPlanController,
		collegeController, talentController, interviewRoundController, interviewScheduleController, interviewController,
		collegeDashboardController, facultyController, // batchSessionPlanController,
		companyController, companyEnquiryController, companyDashboardController, batchSessionController,
		companyEnquiryCallRecordController, batchProjectController,
		domainController, roleController, tenantController,
		menuController, purposeController, outcomeController, dayController, feedbackQuestionController, feedbackOptionController,
		questionController, optionController, oldBatchSessionController, facultyBatchSessionFeedbackController, facultyBatchFeedbackController,
		talentBatchSessionFeedbackController, talentBatchFeedbackController, prerequisiteController,
		courseBatchDashboardController, facultyDashboardController, talentDashboardController,
		enquiryDashboardController, salesPersonDashboardController, adminDashboardController,
		talentCallRecordController, enquiryCallRecordController, callReportController,
		cityController, userController, nextActionController, waitingListController,
		nextActionTypeController, targetCommunityController, departmentController, projectController,
		timesheetController, employeeController, nextActionReportController, feelingController, ahaMomentController,
		feedbackQuestionGroupController, facultyAssessmentController, managementController,
		targetCommunityFunctionController, resourceController, programmingProjectController, salaryTrendController,
		credentialController, collegeDriveController, sourceController, batchSessionsTalentController,
		enquiryController, countryController, stateController, technologyController, specializationController,
		universityController, GeneralTypeController, degreeController, designationController, collegeCampusController,
		examinationController, candidateController, companyBranchController, companyRequirementController,
		batchController, resourceDownloadController, resourceLikeController,
		waitingListReportController, speakerController, loginReportController, packageSummaryController,
		facultyReportController, seminarController, seminarTopicController, fresherSummaryController,
		proSummaryReportController, studentSeminaController, talentReportController, programmingQuestionController,
		programmingQuestionTypeController, problemOfTheDayController, eventController, talentRegistrationController,
		programmingSolutionController, programmingConceptController, programmingLanguageController,
		programmingQuestionSolutionController, programmingQuestionTestCaseController, programmingAssignmentController,
		sessionProgrammingAssignmentController, talentAssignmentSubmissionController,
		blogTopicController, blogController, blogReplyController, blogReactionController, blogViewController,
		batchTalentController, languageCompilerController, androidUSerControl,
		notificationsController, moduleProgrammingConceptController, talentProjectSubmissionController, emailController,
		projectRatingParameterController, conceptDashboardController}

	app.InitializeRouter(routerSpecifier)
	app.WG.Add(2)
	go registerCourseRoutes(app, repository)
	// #Niranjan new router.
	go registerCommunityRoutes(app, repository)
	app.WG.Wait()
}
