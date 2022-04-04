import { DatePipe } from '@angular/common';
import { HttpParams } from '@angular/common/http';
import { ChangeDetectorRef, Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BatchService, IBatch, IBatchSession, IBatchSessionAssignment, IBatchSessionFeedback, IBatchTopicAssignment, IBatchSessionsTalent, IFacultyBatchSessionFeedback } from 'src/app/service/batch/batch.service';
import { UrlConstant, Role } from 'src/app/service/constant';
import { ICourseProgrammingAssignment } from 'src/app/service/course/course.service';
import { IFeedbackOptions, IFeedbackQuestion } from 'src/app/service/feedback/feedback.service';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { IProgrammingAssignment } from 'src/app/service/programming-assignment/programming-assignment.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { ITalentAssignmentScore, ITalentAssignmentSubmission } from 'src/app/service/talent/talent.service';
import { TechnologyService } from 'src/app/service/technology/technology.service';
import { UrlService } from 'src/app/service/url.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { BatchMasterComponent } from '../batch-master/batch-master.component';

@Component({
	selector: 'app-batch-details',
	templateUrl: './batch-details.component.html',
	styleUrls: ['./batch-details.component.css']
})
export class BatchDetailsComponent implements OnInit {

	// components
	ratinglist: any[];
	salesPeople: any[];
	technologies: any[];
	academicYears: any[];
	technologyList: any[];
	courseList: any[];
	courseSessionList: any[]
	daysList: any[]
	facultyList: any[];
	requirementList: any[]
	batchSessionList: any[]
	batchStatus: string[]
	batchObjectives: string[]
	programmingAssignmentList: IProgrammingAssignment[]
	courseAssignments: ICourseProgrammingAssignment[]


	//lock_status
	lockStatus: boolean
	// tech component
	isTechLoading: boolean
	techLimit: number
	techOffset: number

	// batch
	batchDetails: IBatch
	doesEligibilityExists: boolean
	batchForm: FormGroup

	//sidenav
	opened: boolean

	// brochure
	isBrochureUploadedToServer: boolean
	isFileUploading: boolean
	docStatus: string
	displayedFileName: string

	// image
	isLogoUploadedToServer: boolean
	isLogoFileUploading: boolean
	logoDocStatus: string
	logoDisplayedFileName: string

	// batch
	batch: any
	batchID: string
	batchName: string
	sessionName: string
	sessionDate: string
	sessionID: string

	// session
	isSessionManage: boolean
	isSessionAssign: boolean

	// course
	courseID: string
	courseSessions: IBatchSession[]

	// permission
	permission: IPermission
	loginID: string
	isFaculty: boolean
	isTalent: boolean
	isAdmin: boolean
	isSalesPerson: boolean;

	// spinner


	// modal
	modalRef: any;

	// spinner


	// sub-session icon list
	iconList: string[]

	// batch details menu
	batchDetailsMenu: IMenuDetails[]
	assignmentsMenu: IMenuDetails[]
	projectsMenu: IMenuDetails[]

	// batch session programming assignment
	batchSessionAssignment: IBatchSessionAssignment[]
	totalSessionAssignment: number
	assignedSessionAssignments: string[]
	sessionAssignmentList: IBatchTopicAssignment[]
	talentAssignmentScores: ITalentAssignmentScore[]
	scoreSearchForm: FormGroup
	sessionAssignmentTitle: string

	// batch session
	batchSession: any

	// response
	responseProgrammingAsignment: IBatchTopicAssignment[]
	talentAssignmentSubmissions: ITalentAssignmentSubmission[]

	// form
	sessionAssignmentForm: FormGroup
	assignedBatchSessionForm: FormGroup

	// flags
	isViewClicked: boolean
	isUpdateClicked: boolean
	isAttendanceClicked: boolean
	isFeedbackClicked: boolean
	isFeedbackAlreadySubmitted: boolean

	// deleteModalText
	deleteModalText: string

	// search in batch details
	searchBatchDetailsForm: FormGroup

	// tab
	batchDetailsTab: ITabDetails[]

	// batch-session talent attendance
	sessionTalentAttendance: IBatchSessionsTalent[]
	sessionTalentAttendanceForm: FormGroup

	// batch talents
	batchTalents: IMappedTalentDTO[]
	batchTalentLimit: number
	batchTalentOffset: number
	totalBatchTalents: number

	// feedback
	// facultyFeedbackQuestions: IFeedbackQuestionGroup[]
	facultyFeedbackQuestions: IFeedbackQuestion[]
	facultyID: string
	talentID: string
	feedbackTo: string
	facultySessionFeedback: IFacultyBatchSessionFeedback
	feedbackResponse: any[]
	talentFeedbacks: any = []
	talentSessionFeedback: IBatchSessionFeedback[]
	//batch feedbacks
	batchFeedbacks: IBatchSessionFeedback[]
	talentFeedbackComments: IFeedbackComments[]
	talentFeedbackKeywords: string[]

	// session attendance
	talentAttendance: IBatchSessionsTalent[]

	// session feedback
	sessionFeedbackForm: FormGroup
	feedbackForm: FormGroup
	talentSessionFeedbackForm: FormGroup
	isSingleTalentFeedback: boolean

	// faculty
	batchFaculty: any

	// assignment
	isAssignmentManage: boolean
	isAssignmentAssign: boolean

	// project
	isProjectAssignClicked: boolean
	isProjectManageClicked: boolean

	INITIAL_SCORE = 0
	MINIMUM_SCORE = 7
	MAX_SCORE = 10
	FEEDBACK_TALENT_INDEX = 0

	//ngb-rating
	hovered: number

	// modal
	@ViewChild("deleteModal") deleteModal: ElementRef
	@ViewChild("assignAssignmentToBatchSessionModal") assignAssignmentToBatchSessionModal: ElementRef
	@ViewChild("manageAssignmentUpdate") manageAssignmentUpdate: ElementRef
	@ViewChild("talentAttendanceFeedbackModal") talentAttendanceFeedbackModal: ElementRef
	@ViewChild("talentAttendanceModal") talentAttendanceModal: ElementRef
	@ViewChild("batchSessionFeedbackModal") batchSessionFeedbackModal: ElementRef
	@ViewChild("commentsModal") commentsModal: ElementRef

	// template ref
	@ViewChild("batchView") batchView: ElementRef
	@ViewChild("progressReport") progressReport: ElementRef
	@ViewChild("batchSessionPlanCreate") batchSessionPlanCreate: ElementRef
	@ViewChild("batchSessionPlanView") batchSessionPlanView: ElementRef
	@ViewChild("assignmentAssignView") assignmentAssignView: ElementRef
	@ViewChild("assignmentScoreView") assignmentScoreView: ElementRef
	@ViewChild("assignmentManageView") assignmentManageView: ElementRef
	@ViewChild("projectAssignView") projectAssignView: ElementRef
	@ViewChild("projectScoreView") projectScoreView: ElementRef
	@ViewChild("projectManageView") projectManageView: ElementRef
	@ViewChild("studentDetailsFeedbackView") studentDetailsFeedbackView: ElementRef
	@ViewChild("studentDetailsManage") studentDetailsManage: ElementRef
	@ViewChild("noView") noView: ElementRef
	@ViewChild("batchModules") batchModules: ElementRef
	@ViewChild("talentConceptTreeTemplate") talentConceptTreeTemplate: ElementRef

	activatedModulesTab: number;
	eligibilityID: string

	constructor(
		private formBuilder: FormBuilder,
		private urlConstant: UrlConstant,
		private batchService: BatchService,
		private generalService: GeneralService,
		private techService: TechnologyService,
		private utilService: UtilityService,
		private localService: LocalService,
		private urlService: UrlService,
		private fileOpsService: FileOperationService,
		private activatedRoute: ActivatedRoute,
		private router: Router,
		private datePipe: DatePipe,
		private modalService: NgbModal,
		private spinnerService: SpinnerService,
		private role: Role,
		private cdr: ChangeDetectorRef
	) { 
	}

	ngAfterViewInit() {
		this.extractID()
    this.createFeedbackTypeForm()
		this.initializeVariables()
		this.initializeTabs()
		this.createBatchForm()
		this.getAllComponents()
		this.createSearchForm()
		this.cdr.detectChanges()
	}

	initializeVariables(): void {

		this.loginID = this.localService.getJsonValue("loginID")
		this.isSalesPerson = (this.localService.getJsonValue("roleName") == this.role.SALES_PERSON ? true : false)
		this.isFaculty = (this.localService.getJsonValue("roleName") == this.role.FACULTY ? true : false)
		this.isTalent = (this.localService.getJsonValue("roleName") == this.role.TALENT ? true : false)
		this.isAdmin = (this.localService.getJsonValue("roleName") == this.role.ADMIN ? true : false)

		if (this.isAdmin || this.isSalesPerson) {
			this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_BATCH_MASTER_SESSION_DETAILS)
		}
		if (this.isFaculty) {
			this.permission = this.utilService.getPermission(this.urlConstant.MY_BATCH_SESSION_DETAILS)
		}

		// tech components
		this.isTechLoading = false
		this.techLimit = 10
		this.techOffset = 0

		//lock status
		this.lockStatus = false

		this.FEEDBACK_TALENT_INDEX = 0

		// batch
		this.doesEligibilityExists = false
		this.eligibilityID = null

		// flags
		this.isViewClicked = false
		this.isUpdateClicked = false
		this.isAttendanceClicked = false
		this.isFeedbackClicked = false
		this.isFeedbackAlreadySubmitted = false
		this.isBrochureUploadedToServer = false
		// session
		this.isSessionManage = false
		this.isSessionAssign = false

		// assignment
		this.isAssignmentManage = true
		this.isAssignmentAssign = false

		// project
		this.isProjectAssignClicked = false
		this.isProjectManageClicked = false

		this.courseSessions = []
		this.programmingAssignmentList = []
		this.batchSessionAssignment = []
		this.courseAssignments = []
		this.assignedSessionAssignments = []
		this.sessionAssignmentList = []
		this.talentAssignmentScores = []
		this.talentAssignmentSubmissions = []
		this.talentAttendance = []
		this.talentFeedbacks = []
		this.talentSessionFeedback = []
		this.talentFeedbackKeywords = []
		this.talentFeedbackComments = []
		this.daysList = []

		this.totalSessionAssignment = 0

		this.spinnerService.loadingMessage = "Loading..."
		this.displayedFileName = "Select file"
		this.logoDisplayedFileName = "Select file"

		// batch session talent
		this.sessionTalentAttendance = []
		this.batchTalents = []
		this.batchTalentLimit = 10
		this.batchTalentOffset = 0
		this.totalBatchTalents = 0

		// response
		this.responseProgrammingAsignment = []

		// feedback
		this.facultyFeedbackQuestions = []
		this.isSingleTalentFeedback = false

		this.iconList = ["menu_book", "auto_stories", "import_contacts"]

		// spinner

	}

	// constants
	readonly BATCH_DETAILS_INDEX: number = 0
	readonly SESSION_PLAN_INDEX: number = 1
	readonly MODULES_INDEX: number = 2
	readonly ASSIGNMENT_INDEX: number = 3
	readonly PROJECT_INDEX: number = 4
	readonly STUDENT_DETAILS_INDEX: number = 5

	readonly FACULTY_CONCEPT_TREE_URL = "/my/batch/session/details/concept-tree"
	readonly TALENT_CONCEPT_TREE_URL = "/my-batches/concept-trees"
	readonly ADMIN_CONCEPT_TREE_URL = "/training/batch/master/session/details/concept-tree"


	initializeTabs(): void {
		this.batchDetailsTab = [
			{
				tabName: "Batch", subTab: [
					{ tabName: "Details", isActive: false, isRedirect: false, templateName: this.batchView },
					{ tabName: "Progress Report", isActive: false, isRedirect: false, templateName: this.progressReport },
					{ tabName: "Feedback", isActive: false, isRedirect: false, templateName: this.studentDetailsFeedbackView },
				]
			},
			{
				tabName: "Session Plan", subTab: [
					{ tabName: "Manage", isActive: false, isRedirect: false, templateName: this.batchSessionPlanCreate },
					{ tabName: "View", isActive: false, isRedirect: false, templateName: this.batchSessionPlanView },
				]
			},
			{
				tabName: "Modules", subTab: [
					{
						tabName: "Concept Tree", isActive: false, isRedirect: true, fn: () => {
							let url: string

							url = this.isFaculty ? this.FACULTY_CONCEPT_TREE_URL :
								(this.isTalent ? this.TALENT_CONCEPT_TREE_URL : this.ADMIN_CONCEPT_TREE_URL)

							console.log("url -> ", url);

							this.router.navigate([url], {
								queryParams: {
									"batchID": this.batchID,
									"courseID": this.courseID,
									"batchName": this.batchDetails?.batchName
								}
							}).catch(err => {
								console.error(err)
							})

						}
					},
					{ tabName: "Manage", isActive: false, isRedirect: false, templateName: this.batchModules },
					{ tabName: "Assign New", isActive: false, isRedirect: false, templateName: this.batchModules },
				]
			},
			{
				tabName: "Assignments", subTab: [
					{ tabName: "Scores", isActive: false, isRedirect: false, templateName: this.assignmentScoreView },
					{ tabName: "Manage", isActive: false, isRedirect: false, templateName: this.assignmentManageView },
					{ tabName: "Assign New", isActive: false, isRedirect: false, templateName: this.assignmentAssignView },
				]
			},
			{
				tabName: "Projects", subTab: [
					{ tabName: "Scores", isActive: false, isRedirect: false, templateName: this.projectScoreView },
					{ tabName: "Manage", isActive: false, isRedirect: false, templateName: this.projectManageView },
					{ tabName: "Assign New", isActive: false, isRedirect: false, templateName: this.projectAssignView },
				]
			},
			{
				tabName: "Student Details", subTab: [
					{
						tabName: "Student info", isActive: false, isRedirect: true, fn: () => {
							let url: string
							// if (this.isFaculty) {
							// 	url = this.urlConstant.MY_TALENT
							// }

							// if (!this.isFaculty) {
							// 	url = this.urlConstant.TALENT_MASTER
							// }

							url = this.isFaculty ? this.urlConstant.MY_TALENT : this.urlConstant.TALENT_MASTER

							this.router.navigate([url], {
								queryParams: {
									"batchID": this.batchID
								}
							}).catch(err => {
								console.error(err)
							})
						}
					},
					{ tabName: "Manage", isActive: false, isRedirect: false, templateName: this.studentDetailsManage },
				]
			}
		]
	}

	resetSessionFlags(): void {
		this.isSessionManage = false
		this.isSessionAssign = false
	}

	resetAssignmentFlags(): void {
		this.isAssignmentAssign = false
		this.isAssignmentManage = false
	}

	resetProjectFlags(): void {
		this.isProjectAssignClicked = false
		this.isProjectManageClicked = false
	}

	onBatchDetailTabClick(subMenu: IMenuDetails, menuIndex?: number): void {
		this.resetSessionFlags()
		this.resetAssignmentFlags()
		this.resetProjectFlags()

		for (let index = 0; index < this.batchDetailsTab.length; index++) {
			this.batchDetailsTab[index].subTab.find((tab: IMenuDetails) => {
				tab.isActive = false
				if (tab.tabName == subMenu.tabName && menuIndex == index) {
					tab.isActive = true
					this.searchBatchDetailsForm.get("tab").setValue(index)
					return
				}
			})
		}

		this.searchBatchDetailsForm.get("subTab").setValue(subMenu.tabName)
		this.checkTabClick(subMenu, menuIndex)
		this.mergeSearchToQueryParams()
	}

	checkTabClick(subMenu: IMenuDetails, menuIndex: number): void {
		switch (Number(menuIndex)) {
			case this.BATCH_DETAILS_INDEX:
				this.onBatchDetailsClick(subMenu)
				break
			case this.SESSION_PLAN_INDEX:
				break
			case this.MODULES_INDEX:
				this.onModulesDetailsClick(subMenu)
				break
			case this.ASSIGNMENT_INDEX:
				this.onAssignmentsClick(subMenu)
				break
			case this.PROJECT_INDEX:
				this.onProjectClick(subMenu)
				break
			case this.STUDENT_DETAILS_INDEX:
				this.onStudentDetailsClick(subMenu)
				break
		}

	}

	onBatchDetailsClick(menu: IMenuDetails): void {
		switch (menu.tabName) {
			case "Details":
				this.onBatchClick()
				break;
			case "Feedback":
				this.isAttendanceClicked = false
				this.isFeedbackClicked = false

				// this.getBatchSessionDetail()
				this.getBatchTalentDetails()
				this.getAllTalentBatchSessionFeedback()
				break;
			default:
				break;
		}
	}

	onModulesDetailsClick(menu: IMenuDetails): void {

		switch (menu.tabName) {
			case "Assign New":
				this.activatedModulesTab = 1
				break;
			case "Manage":
				this.activatedModulesTab = 2
				break;
			// case "Concept Tree":
			// 	this.redirectToConceptTree()
			// break;
			default:
				break;
		}
		// console.log("activatedModulesTab", this.activatedModulesTab);

	}

	onAssignmentsClick(menu: IMenuDetails): void {
		switch (menu.tabName) {
			case "Assign New":
				this.isAssignmentAssign = true
				break;
			case "Manage":
				this.isAssignmentManage = true
				break;
		}
	}

	onProjectClick(menu: IMenuDetails): void {
		switch (menu.tabName) {
			case "Assign New":
				this.isProjectAssignClicked = true
				break;
			case "Manage":
				this.isProjectManageClicked = true
				break;
		}
	}

	onStudentDetailsClick(menu: IMenuDetails): void {
		switch (menu.tabName) {
			case "Manage List":
				break;
		}
	}

	createSearchForm(): void {
		this.searchBatchDetailsForm = this.formBuilder.group({
			batchID: new FormControl(this.batchID),
			courseID: new FormControl(this.courseID),
			batchName: new FormControl(this.batchName),
			sessionID: new FormControl(null),
			tab: new FormControl(null),
			subTab: new FormControl(null),
		})
		this.getTab()
		this.mergeSearchToQueryParams()
	}

	getTab(): void {
		this.activatedRoute.queryParams.subscribe(queryParams => {
			console.log("queryParams", queryParams)
			this.searchBatchDetailsForm.patchValue(queryParams)
			console.log(this.searchBatchDetailsForm.value);

			if (queryParams.tab && queryParams.subTab) {

				for (let index = 0; index < this.batchDetailsTab[queryParams.tab].subTab.length; index++) {
					if (this.batchDetailsTab[queryParams.tab].subTab[index].tabName === queryParams.subTab) {
						this.onBatchDetailTabClick(this.batchDetailsTab[queryParams.tab].subTab[index], queryParams.tab)
						return
					}
				}
			}

			if (queryParams.moduleTab) {
				this.batchDetailsTab[2].subTab[2].isActive = true

				// this.onBatchDetailTabClick({ tabName: "Manage", isActive: true }, this.SESSION_PLAN_INDEX)
				this.onBatchDetailTabClick(this.batchDetailsTab[2].subTab[2], 2)
				return

			}

			this.batchDetailsTab[this.SESSION_PLAN_INDEX].subTab[0].isActive = true

			// this.onBatchDetailTabClick({ tabName: "Manage", isActive: true }, this.SESSION_PLAN_INDEX)
			this.onBatchDetailTabClick(this.batchDetailsTab[this.SESSION_PLAN_INDEX].subTab[0], this.SESSION_PLAN_INDEX)

		}, err => {
			console.error(err);
		})

		// let queryParams = this.activatedRoute.snapshot.queryParams		
	}

  readonly FACULTY_SESSION_FEEDBACK: string = 'Faculty_Session_Feedback'
  readonly FACULTY_SESSION_FEEDBACK_NON_TECH: string = 'Faculty_Session_Feedback_Non_Technical_Question'

	async getAllComponents(): Promise<void> {
		this.getBatchDetails()
		this.getBatchObjectives()
		this.getTechnology()
		this.getAcademicYear()
		this.getStudentRatingList()
		this.getAllSalesPeople()
		this.getBatchStatus()
		this.getCourseList()
		this.getFacultyList()
		this.getDaysList()
		this.getRequirementList()
    this.facultyFeedbackQuestions = await this.getFeedbackQuestionsForFaculty(this.feedbackTypeForm.get("feedbackType").value)
		// this.getProgrammingAssignments()
	}

	// Testing
	backToPreviousPage(): void {
		this.urlService.goBack(BatchMasterComponent.name)
	}


	get ongoingOperations() {
		return this.spinnerService.ongoingOperations
	}

	ngOnInit(): void {
	}

	extractID(): void {
		this.activatedRoute.queryParamMap.subscribe(
			(params: any) => {
				this.batchID = params.get("batchID")
				this.courseID = params.get("courseID")
				this.batchName = params.get("batchName")
			}, (err: any) => {
				console.error(err);
			})
	}

	onBatchClick(): void {
		this.getBatchDetails()
	}

	// ============================================================= BATCH =============================================================

	createBatchForm(): void {
		this.doesEligibilityExists = false

		this.batchForm = this.formBuilder.group({
			id: new FormControl(null),
			batchName: new FormControl(null, [Validators.required]),
			code: new FormControl(null),
			course: new FormControl(null, [Validators.required]),
			startDate: new FormControl(null),
			estimatedEndDate: new FormControl(null),
			totalStudents: new FormControl(null),
			totalIntake: new FormControl(null),
			batchStatus: new FormControl("Upcoming"),
			batchObjective: new FormControl(null),
			isActive: new FormControl(true),
			isB2B: new FormControl(null),
			brochure: new FormControl(null),
			logo: new FormControl(null),
			meetLink: new FormControl(null),
			telegramLink: new FormControl(null),
			eligibility: this.formBuilder.group({
				id: new FormControl(null),
				technologies: new FormControl(null),
				studentRating: new FormControl(null),
				experience: new FormControl(null),
				academicYear: new FormControl(null),
			}),
			// faculty: new FormControl(null, [Validators.required]),
			salesPerson: new FormControl(null, [Validators.required]),
			requirement: new FormControl(null),
			batchTimings: this.formBuilder.array([]),
			isApplyToAllSessions: new FormControl(false),
		})
		this.batchForm.disable()
		this.addValidatorsToBatchForm()
	}

	addValidatorsToBatchForm(): void {
		if (this.isAdmin) {
			this.addAdminValidators()

			return
		}

		// this.addFacultyValidators()
	}

	addAdminValidators(): void {
		this.batchForm.get("batchName").setValidators([Validators.required])
		this.batchForm.get("batchObjective").setValidators([Validators.required])
		this.batchForm.get("course").setValidators([Validators.required])
		this.batchForm.get("totalIntake").setValidators([Validators.required])
		this.batchForm.get("salesPerson").setValidators([Validators.required])
		this.batchForm.get("isB2B").setValidators([Validators.required])

		this.batchForm.get("startDate").disable()
		// this.batchForm.get("requirement").setValidators([Validators.required])
		// this.batchForm.get("brochure").setValidators([Validators.required])

		this.utilService.updateValueAndValiditors(this.batchForm)
	}

	get batchTiming() {
		return this.batchForm.get('batchTimings') as FormArray
	}

	// createBatchTimingForm(day: any): void {
	// 	this.batchTiming.push(this.formBuilder.group({
	// 		id: new FormControl(),
	// 		day: new FormControl(null, [Validators.required]),
	// 		fromTime: new FormControl(
	// 			(this.batchTiming.at(0)?.get('fromTime')?.value != null) ? this.batchTiming.at(0)?.get('fromTime')?.value : null,
	// 			[Validators.required]),
	// 		toTime: new FormControl(
	// 			(this.batchTiming.at(0)?.get('toTime')?.value != null) ? this.batchTiming.at(0)?.get('toTime')?.value : null,
	// 			[Validators.required])
	// 	}))
	// }

	// delete batch timing from batch form.
	deleteBatchTimingInBatchForm(index: number): void {
		if (confirm("Are you sure you want to delete the batch timing?")) {
			this.batchTiming.removeAt(index)
			if (index == 0 && this.batchTiming.length > 1) {
				for (let i = 0; i < this.batchTiming.length; i++) {
					this.batchTiming.at(i).get('fromTime').enable()
					this.batchTiming.at(i).get('toTime').enable()
				}
			}
			this.batchForm.markAsDirty()
			this.formatDayListFields()
		}
	}

	setFromTime(index: number): void {
		if (index == 0 && this.batchTiming.at(index)?.get('fromTime')?.value) {
			for (let index = 0; index < this.batchTiming.controls.length; index++) {
				if (!this.batchTiming.at(index)?.get('fromTime')?.value) {
					this.batchTiming.at(index)?.get('fromTime')?.setValue(this.batchTiming.at(0)?.get('fromTime')?.value)
				}
			}
		}
	}

	setToTime(index: number): void {
		if (index == 0 && this.batchTiming.at(index)?.get('toTime')?.value) {
			for (let index = 0; index < this.batchTiming.controls.length; index++) {
				if (!this.batchTiming.at(index)?.get('toTime')?.value) {
					this.batchTiming.at(index)?.get('toTime')?.setValue(this.batchTiming.at(0)?.get('toTime')?.value)
				}
			}
		}
	}

	eligilibilityChecked(event) {
		event.target.checked = true
		if (this.doesEligibilityExists) {
			if (this.batchForm.get('eligibility').valid) {
				if (confirm("This action will delete the batch eligibility. Are you sure you want to delete batch eligibility?")) {
					event.target.checked = false
				} else {
					return
				}
			}
			this.eligibilityID = this.batchForm.get('eligibility.id').value
			// console.log("this.eligibilityID",this.eligibilityID)
			this.batchForm.get('eligibility.id').setValue(null)
			this.doesEligibilityExists = false
			this.setEligibility()
			return;
		}
		this.doesEligibilityExists = true;
		this.setEligibility()
	}

	setEligibility() {
		this.technologyCheck()
		this.batchForm.markAsDirty()

		// this.studentRatingCheck()
		// this.experienceCheck()
		// this.academicYearCheck()
		// console.log(this.batchForm.get('eligibility'));
	}

	technologyCheck() {
		const technologyControl = this.batchForm.get('eligibility.technologies')
		if (this.doesEligibilityExists) {
			technologyControl.setValidators([Validators.required])
			if (this.eligibilityID != null) {
				this.batchForm.get('eligibility.id').setValue(this.eligibilityID)
			}
		} else {
			technologyControl.setValue(null)
			technologyControl.setValidators(null)
			technologyControl.markAsUntouched()
		}

		technologyControl.updateValueAndValidity()
	}

	// studentRatingCheck() {
	// 	const studentRatingControl = this.batchForm.get('eligibility.studentRating')
	// 	if (this.doesEligibilityExists) {
	// 		studentRatingControl.setValidators([Validators.required])
	// 	} else {
	// 		studentRatingControl.setValue(null)
	// 		studentRatingControl.setValidators(null)
	// 		studentRatingControl.markAsUntouched()
	// 	}
	// 	studentRatingControl.updateValueAndValidity()
	// }

	// experienceCheck() {
	// 	const experienceControl = this.batchForm.get('eligibility.experience')
	// 	if (this.doesEligibilityExists) {
	// 		experienceControl.setValidators([Validators.required])
	// 	} else {
	// 		experienceControl.setValue(null)
	// 		experienceControl.setValidators(null)
	// 		experienceControl.markAsUntouched()
	// 	}
	// 	experienceControl.updateValueAndValidity()
	// }

	// academicYearCheck() {
	// 	const academicYearControl = this.batchForm.get('eligibility.academicYear')
	// 	if (this.doesEligibilityExists) {
	// 		academicYearControl.setValidators([Validators.required])
	// 	} else {
	// 		academicYearControl.setValue(null)
	// 		academicYearControl.setValidators(null)
	// 		academicYearControl.markAsUntouched()
	// 	}
	// 	academicYearControl.updateValueAndValidity()
	// }

	getBatchDetails(): void {
		this.spinnerService.loadingMessage = "Loading..."

		// let queryParams: any = {
		//   batchID: this.batchID
		// }
		this.batchService.getBatch(this.batchID).subscribe((response: any) => {
			this.batchDetails = response.body
			// console.log(this.batchDetails);
			this.courseID = this.batchDetails?.course?.id
			this.batchFaculty = this.batchDetails?.faculty
			this.searchBatchDetailsForm.get('courseID').setValue(this.courseID)
			this.mergeSearchToQueryParams()
			this.editBatchForm()

		}, (err: any) => {
			console.error(err);

			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err.error.error)
		})
	}

	mergeSearchToQueryParams(): void {
		this.router.navigate([], {
			relativeTo: this.activatedRoute,
			queryParams: this.searchBatchDetailsForm.value,
			queryParamsHandling: 'merge'
		});
	}

	editBatchForm(): void {
		this.isViewClicked = true
		this.isUpdateClicked = false

		this.resetBrochureUploadFields()
		this.resetLogoUploadFields()

		// console.log(this.batchDetails);

		if (this.batchDetails?.brochure) {
			this.displayedFileName = `<a href=${this.batchDetails?.brochure} target="_blank">Batch Brochure</a>`
		}
		if (this.batchDetails?.logo) {
			this.logoDisplayedFileName = `<a href=${this.batchDetails?.logo} target="_blank">Batch Logo</a>`
		}

		this.createBatchForm()
		this.batchForm.disable()

		this.batchDetails.startDate = this.datePipe.transform(this.batchDetails?.startDate, 'yyyy-MM-dd')
		this.batchDetails.estimatedEndDate = this.datePipe.transform(this.batchDetails?.estimatedEndDate, 'yyyy-MM-dd')

		this.patchBatchForm()
		// console.log(this.batchForm.controls);
	}

	patchBatchForm(): void {
		this.checkBatchEligibility();

		// Sort batch timings by order of days.
		this.batchDetails.batchTimings.sort(this.sortBatchTimings)

		// Add batch timings to batch timing form.
		for (let i = 0; i < this.batchDetails.batchTimings.length; i++) {
			this.addBatchTimingForm()
		}
		// console.log("batchDetails",this.batchDetails);
		this.batchForm.patchValue(this.batchDetails)
		this.getDaysList()

		// make eligibility null if eligibility does not exist
		// if (Object.keys(this.selectedBatch.eligibility).length == 0) {
		// 	this.selectedBatch.eligibility = null
		// }

		// if (!this.doesEligibilityExists) {
		// 	this.selectedBatch.eligibility = null
		// }

		// (this.batchForm.controls["batchTimings"] as FormArray).clear()
		// for (let index = 0; index < this.batchDetails?.batchTimings?.length; index++) {
		// 	this.createBatchTimingForm()
		// }

		// // console.log(this.batchDetails)
		// this.batchForm.patchValue(this.batchDetails)
		// this.disableForm(this.batchForm)
		this.batchForm.disable()

		this.editBatchEligibility()
	}

	checkBatchEligibility(): void {
		if (this.batchDetails.eligibility != null) {
			this.doesEligibilityExists = true
		} else {
			this.batchDetails.eligibility = {}
			this.doesEligibilityExists = false
		}
	}

	// make eligibility null if eligibility does not exist
	editBatchEligibility(): void {
		if (Object.keys(this.batchDetails.eligibility).length == 0) {
			this.batchDetails.eligibility = null
		}
	}

	onBatchUpdateClick(): void {
		this.isViewClicked = false
		this.isUpdateClicked = true

		if (this.batchDetails.brochure) {
			this.displayedFileName = `<a href=${this.batchDetails.brochure} target="_blank">Batch Brochure</a>`
		}
		if (this.batchDetails.logo) {
			this.logoDisplayedFileName = `<a href=${this.batchDetails.logo} target="_blank">Batch Logo</a>`
		}
		this.batchForm.enable()
	}

	updateBatch(): void {
		this.batchForm.enable()
		let batch = this.batchForm.value
		this.utilService.deleteNullValueIDFromObject(batch)
		this.utilService.deleteNullValuePropertyFromObject(batch.eligibility)
		batch.startDate = this.datePipe.transform(batch.startDate, "yyyy-MM-dd")
		// this.setFormFields(batch)
		if (batch.eligibility) {
			this.deleteEligibilityIfEmpty(batch);
		}
		// console.log(batch);
		this.spinnerService.loadingMessage = "Updating batch";

		this.batchService.updateBatch(batch).subscribe((respond: any) => {
			this.eligibilityID = null
			this.createBatchForm()
			this.batchForm.disable()
			alert("Batch Updated Successfully");

			this.getBatchDetails()
		}, (err) => {

			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(this.utilService.getErrorString(err));
			console.error(err.error.error)
		})
	}

	// Extract id from objects and give it to ompany requirement.
	setFormFields(batch: any): void {
		batch.startDate = new Date(this.batchForm.get("startDate")?.value).toISOString()
		batch.estimatedEndDate = new Date(this.batchForm.get("estimatedEndDate")?.value).toISOString()
	}

	deleteEligibilityIfEmpty(batch: any): void {
		if (Object.keys(batch.eligibility).length === 0 && batch.eligibility.constructor === Object) {
			delete batch.eligibility
		}
	}

	//validate update form
	validateBatchDetails(): void {
		if (this.isFileUploading || this.isLogoFileUploading) {
			alert("Please wait till file is being uploaded")
			return
		}
		// console.log(this.batchForm.controls);

		if (this.batchForm.invalid) {
			this.batchForm.markAllAsTouched();
			return
		}
		if (this.batchForm.get('isActive').value == false || this.batchForm.get('batchStatus').value == "Finished") {
			if (confirm("Making batch inactive or its status as 'Finished' will also mean that all the related waiting list entries in talent"
				+ "and enquiry will be inactive. Do you want to go ahead??")) {
				this.updateBatch();
				return
			}
			return
		}
		this.updateBatch();
	}

	resetBrochureUploadFields(): void {
		this.isBrochureUploadedToServer = false
		this.isFileUploading = false
		this.displayedFileName = "Select file"
		this.docStatus = ""
	}

	resetLogoUploadFields(): void {
		this.isLogoUploadedToServer = false
		this.isLogoFileUploading = false
		this.logoDisplayedFileName = "Select file"
		this.logoDocStatus = ""
	}

	// Compare two Object.
	compareFn(ob1: any, ob2: any): any {
		if (ob1 == null && ob2 == null) {
			return true
		}
		return ob1 && ob2 ? ob1.id == ob2.id : ob1 == ob2;
	}

	//On uplaoding brochure
	onResourceSelect(event: any) {
		this.docStatus = ""
		let files = event.target.files
		if (files && files.length) {
			let file = files[0]
			let err = this.fileOpsService.isDocumentFileValid(file)
			if (err != null) {
				this.docStatus = `<p><span>&#10060;</span> ${err}</p>`
				return
			}
			// console.log(file)
			// Upload brochure if it is present.]
			this.isFileUploading = true
			this.fileOpsService.uploadBrochure(file,
				this.fileOpsService.BATCH_FOLDER + this.fileOpsService.BROCHURE_FOLDER)
				.subscribe((data: any) => {
					this.batchForm.markAsDirty()
					this.batchForm.patchValue({
						brochure: data
					})
					this.displayedFileName = file.name
					this.isFileUploading = false
					this.isBrochureUploadedToServer = true
					this.docStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
				}, (error) => {
					this.isFileUploading = false
					this.docStatus = `<p><span>&#10060;</span> ${error}</p>`
				})
		}
	}

	// On uplaoding logo
	onImageSelect(event: any) {
		this.logoDocStatus = ""
		let files = event.target.files
		if (files && files.length) {
			let file = files[0]
			let err = this.fileOpsService.isImageFileValid(file)
			if (err != null) {
				this.logoDocStatus = `<p><span>&#10060;</span> ${err}</p>`
				return
			}
			// console.log(file)
			// Upload brochure if it is present.]
			this.isLogoFileUploading = true
			this.fileOpsService.uploadLogo(file,
				this.fileOpsService.BATCH_FOLDER + this.fileOpsService.LOGO_FOLDER).subscribe((data: any) => {
					this.batchForm.markAsDirty()
					this.batchForm.patchValue({
						logo: data
					})
					this.logoDisplayedFileName = file.name
					this.isLogoFileUploading = false
					this.isLogoUploadedToServer = true
					this.logoDocStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
				}, (error) => {
					this.isLogoFileUploading = false
					this.logoDocStatus = `<p><span>&#10060;</span> ${error}</p>`
				})
		}
	}


	getAllTalentBatchSessionFeedback(): void {
		this.spinnerService.loadingMessage = "Getting all feedback"

		this.batchService.getBatchTalentsFeedback(this.batchID).subscribe((response: any) => {
			this.talentFeedbacks = response.body
			// console.log(this.talentFeedbacks);
			this.seperateTalentSessionFeedback()
			this.getTalentFeedbackKeywords()
		},
			(err: any) => {
				console.error(err)
				if (err.statusText.includes('Unknown')) {
					alert("No connection to server. Check internet.")
					return
				}
				alert(err.error.error)
			}).add(() => {

			})
	}

	seperateTalentSessionFeedback(): void {
		this.talentSessionFeedback = []
		for (let index = 0; index < this.talentFeedbacks.length; index++) {
			// for (let j = 0; j < this.talentFeedbacks[index]?.batchSessionFeedback.length; j++) {
			this.talentSessionFeedback.push(...this.talentFeedbacks[index]?.batchSessionFeedback)
			// }
		}
		// console.log(this.talentSessionFeedback);
	}

	getTalentFeedbackKeywords(): void {
		this.talentFeedbackKeywords = []
		for (let index = 0; index < this.talentFeedbacks.length; index++) {
			if (this.talentFeedbacks[index]?.batchSessionFeedback &&
				this.talentFeedbacks[index]?.batchSessionFeedback?.length > 0) {
				for (let question of this.talentFeedbacks[index]?.batchSessionFeedback) {
					// console.log(question);

					if (question.question.keyword && question.question.keyword != "" &&
						!this.talentFeedbackKeywords.includes(question.question.keyword)) {
						this.talentFeedbackKeywords.push(question.question.keyword)
					}
				}
				// console.log(this.talentFeedbackKeywords);
				return
			}
		}
	}

	calculateKeywordAvgScore(keyword: string, sessionID: string): string {

		let averageScore: number = 0
		let totalMaxScore: number = 0

		for (let index = 0; index < this.talentSessionFeedback.length; index++) {
			if (this.talentSessionFeedback[index].question.keyword == keyword && this.talentSessionFeedback[index]?.option &&
				this.talentSessionFeedback[index].batchSessionID == sessionID) {
				totalMaxScore += this.talentSessionFeedback[index]?.question.maxScore
				averageScore += this.talentSessionFeedback[index]?.option?.key
			}
		}

		if (averageScore > 0) {
			return ((averageScore * 10 / totalMaxScore)) + "/10"
		}
		return "-"
	}

	onCommentClick(sessionID: string): void {
		// console.log(sessionID);
		let feedbackCommentsMap = new Map<string, string[]>()

		for (let index = 0; index < this.talentSessionFeedback.length; index++) {
			if (this.talentSessionFeedback[index].batchSessionID == sessionID) {
				if (!this.talentSessionFeedback[index]?.question?.hasOptions)
					if (feedbackCommentsMap.get(this.talentSessionFeedback[index]?.question?.question)) {
						let question = this.talentSessionFeedback[index]?.question?.question
						let answer = this.talentSessionFeedback[index]?.answer

						let mappedAnswers = feedbackCommentsMap.get(question)
						mappedAnswers.push(answer)

						feedbackCommentsMap.set(question, mappedAnswers)
					} else {
						feedbackCommentsMap.set(this.talentSessionFeedback[index]?.question?.question,
							[this.talentSessionFeedback[index]?.answer])
					}

				this.talentFeedbackComments = []

				for (let [key, value] of feedbackCommentsMap) {
					this.talentFeedbackComments.push({
						question: key, answers: value
					})
				}
			}
		}

		this.openModal(this.commentsModal, "lg")
	}

	getBatchTalentDetails(isTalentAttendanceHandler: boolean = false, isOpenModal: boolean = false): void {
		this.spinnerService.loadingMessage = "Getting talents"

		this.batchTalents = []
		this.totalBatchTalents = 0

		let queryParams: any = {
			"isActive": 1,
		}

		this.batchService.getBatchTalentDetails(this.batchID, queryParams).subscribe((response: any) => {
			this.batchTalents = response.body
			// console.log(this.batchTalents);
			this.totalBatchTalents = this.batchTalents.length

			if (this.isFeedbackClicked) {
				this.FEEDBACK_TALENT_INDEX = -1
				this.getTalentFeedbackToGiven()
			}


		}, (err: any) => {
			this.totalBatchTalents = 0
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err.error.error)
			console.error(err);
		}).add(() => {

		})
	}


	getBatchSessionDetail(): void {
		let queryParams: any = {
			batchSessionID: this.searchBatchDetailsForm.get("sessionID").value
		}

		this.batchService.getBatchSessionList(this.batchID, queryParams).subscribe((response: any) => {
			this.batchSession = response.body[0]
			// console.log(this.batchSession);
		}, (err: any) => {
			console.error(err);
		})
	}

	onBatchTalentClick(talentDetail: IMappedTalentDTO): void {
		talentDetail.showDetails = !talentDetail.showDetails

		if (!talentDetail.showDetails) {
			return
		}

		this.spinnerService.loadingMessage = "Getting talent session details"

		this.batchService.getTalentTopicDetails(this.batchID, talentDetail.talent.id).subscribe((response: any) => {
			talentDetail.sessionDetails = response.body
			// console.log(talentDetail.sessionDetails);
		}, (err: any) => {
			talentDetail.showDetails = false
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err.error.error)
			console.error(err);
		}).add(() => {

		})
	}



	// ================================================= FEEDBACK =================================================

  feedbackTypeForm: FormGroup 
	mappedTalent: IMappedTalentDTO


  createFeedbackTypeForm(): void {
    this.feedbackTypeForm = this.formBuilder.group({
      feedbackType: new FormControl(this.FACULTY_SESSION_FEEDBACK)
    })
  }

  async onFeedbackTypeChange(): Promise<void> {
    console.log(this.feedbackTypeForm);
    this.facultyFeedbackQuestions = await this.getFeedbackQuestionsForFaculty(this.feedbackTypeForm.get("feedbackType").value)
    
    this.createFeedbackComponents(this.mappedTalent)
  }

	onFeedbackAlreadyGivenClick(details: IMappedTalentDTO, batchSession: any): void {
		console.log("batchSession -> ", batchSession);
		this.batchSession = batchSession
		this.getTalentFeedbackDetails(batchSession.id, details.talent.id)
		this.feedbackTo = details.talent.firstName + " " + details.talent.lastName
		this.sessionDate = batchSession.date
		this.isFeedbackClicked = true
		this.isFeedbackAlreadySubmitted = true
		this.isAttendanceClicked = false
	}

	createSubmittedFeedbackQuestionsForm(submittedFeedback: any): void {
		// if (this.facultyFeedbackQuestions) {
		// console.log(submittedFeedback);

		for (let index = 0; index < submittedFeedback.length; index++) {
			this.addFeedbackQuestionsToArray()
			// let index = 0
			// initialize form
			if (this.isFaculty) {
				this.feedbackArray.at(index).get('facultyID').setValue(submittedFeedback[index].facultyID)
				this.feedbackArray.at(index).get("talentID").setValue(submittedFeedback[index].talentID)
			}
			// if (this.isTalent) {
			// 	this.feedbackArray.at(index).get('facultyID').setValue(this.facultyID)
			// 	this.feedbackArray.at(index).get('talentID').setValue(this.loginID)
			// }
			// if (this.isAdmin) {
			// 	this.feedbackArray.at(index).get('facultyID').setValue(this.facultyID)
			// 	this.feedbackArray.at(index).get('talentID').setValue(details.talent.id)
			// }

			this.feedbackArray.at(index).get("questionID").setValue(submittedFeedback[index].questionID)

			if (submittedFeedback?.question?.hasOptions) {
				this.feedbackArray.at(index).get("optionID").setValue(submittedFeedback[index]?.optionID)
				this.feedbackArray.at(index).get("option").setValue(submittedFeedback[index]?.option)
				this.feedbackArray.at(index).get("answer").setValue(submittedFeedback[index]?.answer)
				this.feedbackArray.at(index).get("hover").setValue(null)
			} else {
				this.feedbackArray.at(index).get("optionID").setValue(submittedFeedback[index]?.optionID)
				this.feedbackArray.at(index).get("option").setValue(submittedFeedback[index]?.option)
				this.feedbackArray.at(index).get("answer").setValue(submittedFeedback[index]?.answer)
				this.feedbackArray.at(index).get("hover").setValue(null)
			}
		}
		// console.log(this.feedbackArray.controls);

	}
	// 	// console.log(this.feedbackGroupArray.controls);
	// }

	getTalentFeedbackDetails(sessionID: string, talentID: string): void {
		let queryParams = new HttpParams()
		queryParams = queryParams.append("talentID", talentID)
		this.batchService.getAllTalentBatchSessionFeedback(this.batchID, sessionID, queryParams).subscribe((response: any) => {
			// console.log(response.body);
			// this.facultyFeedbackQuestions = response.body
			this.createFeedbackForm()
			this.createSubmittedFeedbackQuestionsForm(response.body)
			// this.addFacultyFeedbackToForm(details)
			this.openModal(this.batchSessionFeedbackModal, "lg")
		}, (err: any) => {
			console.error(err)
		})
	}

  async getFeedbackQuestionsForFaculty(questionType: string): Promise<any> {
    try {
      return await new Promise<any>((resolve, reject) => {
        this.spinnerService.loadingMessage = "Getting Feedback Questions"
        // this.spinnerService.startSpinner()
        this.generalService.getFeedbackQuestionByType(questionType).subscribe((response: any) => {
          resolve(response.body)
        }, (err: any) => {
          reject(err)
        })
      })
    } finally {
      // this.spinnerService.stopSpinner()
    }
  }

	createFeedbackComponents(details: IMappedTalentDTO): void {
		this.createFeedbackForm()
		this.addFacultyFeedbackToForm(details)
	}

	onFeedbackClick(details: IMappedTalentDTO, batchSession: any): void {
		this.sessionID = batchSession.id
		this.batchSession = batchSession
		this.mappedTalent = details
		this.sessionDate = batchSession.date

		this.isSingleTalentFeedback = true
		this.isFeedbackAlreadySubmitted = false
		this.createFeedbackForm()
		this.addFacultyFeedbackToForm(details)
		this.openModal(this.batchSessionFeedbackModal, "xl")
	}

	addFacultyFeedbackToForm(details: IMappedTalentDTO): void {
		this.feedbackResponse = []
		this.feedbackTo = details.talent.firstName + " " + details.talent.lastName

		if (this.isAdmin) {
			// this.facultyID = sessionFeedback.faculty.id
			this.facultyID = this.batchFaculty.id
			this.talentID = details.talent.id
		}
		this.createFeedbackQuestionsForm(details)
		this.isFeedbackClicked = true
		this.isAttendanceClicked = false
		// this.openModal(this.batchSessionFeedbackModal, "lg")
	}

	createFeedbackForm(): void {
		this.feedbackForm = this.formBuilder.group({
			feedbacks: new FormArray([])
		})
	}

	get feedbackArray(): FormArray {
		return this.feedbackForm.get("feedbacks") as FormArray
	}

	addFeedbackQuestionsToArray(): void {
		// this.sessionFeedbackForm = this.formBuilder.group({
		this.feedbackArray.push(this.formBuilder.group({
			batchID: new FormControl(this.batchID),
			talentID: new FormControl(null, [Validators.required]),
			facultyID: new FormControl(this.loginID),
			questionID: new FormControl(null, [Validators.required]),
			optionID: new FormControl(null, [Validators.required]),
			option: new FormControl(null, [Validators.required]),
			answer: new FormControl(null, [Validators.required]),
			hover: new FormControl(null, [Validators.required]),
		}))
	}

	createFeedbackQuestionsForm(details: IMappedTalentDTO): void {
		if (this.facultyFeedbackQuestions) {
			for (let index = 0; index < this.facultyFeedbackQuestions.length; index++) {
				this.addFeedbackQuestionsToArray()

				// initialize form
				if (this.isFaculty) {
					this.feedbackArray.at(index).get('facultyID').setValue(this.loginID)
					this.feedbackArray.at(index).get("talentID").setValue(details.talent.id)
				}
				if (this.isTalent) {
					this.feedbackArray.at(index).get('facultyID').setValue(this.facultyID)
					this.feedbackArray.at(index).get('talentID').setValue(this.loginID)
				}
				if (this.isAdmin) {
					this.feedbackArray.at(index).get('facultyID').setValue(this.facultyID)
					this.feedbackArray.at(index).get('talentID').setValue(details.talent.id)
				}

				this.feedbackArray.at(index).get("questionID").setValue(this.facultyFeedbackQuestions[index].id)

				if (this.facultyFeedbackQuestions[index].hasOptions) {
					this.feedbackArray.at(index).get("optionID").setValidators([Validators.required])
					this.feedbackArray.at(index).get("option").setValidators([Validators.required])
					this.feedbackArray.at(index).get("answer").setValidators(null)
					this.feedbackArray.at(index).get("hover").setValidators(null)
				} else {
					this.feedbackArray.at(index).get("optionID").setValidators(null)
					this.feedbackArray.at(index).get("option").setValidators(null)
					this.feedbackArray.at(index).get("answer").setValidators([Validators.required])
					this.feedbackArray.at(index).get("hover").setValidators(null)
				}
			}
		}
		// console.log(this.feedbackGroupArray.controls);
	}

	onFacultyFeedbackChange(feedbackQuestionControl: any, feedbackOptions: IFeedbackOptions[]) {
		let optionID = feedbackQuestionControl.get("optionID").value
		for (let index = 0; index < feedbackOptions.length; index++) {
			if (optionID == feedbackOptions[index].id) {
				feedbackQuestionControl.get('answer').setValue(feedbackOptions[index].value)
				feedbackQuestionControl.get('option').setValue(feedbackOptions[index])
			}
		}
	}

	onFacutyFeedbackInput(event: any, index: number, feedbackQuestionControl: any): void {
		// console.log(feedbackQuestionControl)

		let options: IFeedbackOptions[] = this.facultyFeedbackQuestions[index].options
		feedbackQuestionControl.get('answer').setValue(String(event))
		feedbackQuestionControl.get('optionID').setValue(null)
		feedbackQuestionControl.get('option').setValue(null)

		for (let index = 0; index < options.length; index++) {
			if (event == options[index].key) {
				feedbackQuestionControl.get('answer').setValue((options[index].key).toString())
				feedbackQuestionControl.get('optionID').setValue(options[index].id)
				feedbackQuestionControl.get('option').setValue(options[index])
			}
		}
	}

	addFacultyFeedback(errors: string[]): void {
		// console.log(this.feedbackForm.value.feedbacks);
		// console.log(this.sessionID);

		this.spinnerService.loadingMessage = "Adding feedback"

		this.batchService.addFacultySessionFeedbacks(this.batchID, this.sessionID, this.feedbackForm.value.feedbacks).
			subscribe((response: any) => {
				// console.log(response)
				this.feedbackResponse = []
			}, (err: any) => {
				console.error(err)
				if (err.statusText.includes('Unknown')) {
					alert("No connection to server. Check internet.")
					return
				}
				errors.push(err.error.error)
			}).add(() => {

			})
	}

	validateFacultyFeedback(): void {
		// console.log(this.feedbackForm.controls);
		// console.log(this.feedbackForm.value);

		if (this.feedbackForm.invalid) {
			this.feedbackForm.markAllAsTouched()
			return
		}

		let errors: string[] = []
		this.addFacultyFeedback(errors)

		if (this.isSingleTalentFeedback) {
			this.getBatchTalentDetails()
			this.getAllTalentBatchSessionFeedback()
			this.modalRef.close()
			return
		}

		this.getTalentFeedbackToGiven()

		if (this.FEEDBACK_TALENT_INDEX == this.talentAttendance.length) {
			if (errors.length == 0) {
				this.onBatchDetailTabClick({ tabName: "View", isActive: true }, this.SESSION_PLAN_INDEX)
				this.modalRef.close()
				return
			}
			if (errors.length > 0) {
				let errorString = ""
				for (let index = 0; index < errors.length; index++) {
					errorString += (index == 0 ? "" : "\n") + errors[index]
				}
				alert(errorString)
			}
		}
	}

	getTalentFeedbackToGiven(): void {
		let isTalentFound = false

		for (let index = ++this.FEEDBACK_TALENT_INDEX; index < this.talentAttendance.length; index++) {
			if (this.talentAttendance[index].isPresent) {
				this.batchTalents.find((val: IMappedTalentDTO) => {
					if (val.talent.id == this.talentAttendance[index].talent.id) {
						this.createFeedbackComponents(val)
						this.FEEDBACK_TALENT_INDEX = index
						isTalentFound = true
						return
					}
				})
				if (isTalentFound) {
					return
				}
			}
		}
	}

	// ============================================================= COMPONENTS =============================================================

	// Get All Student Rating
	getStudentRatingList(): void {
		this.generalService.getGeneralTypeByType("talent_rating").subscribe((respond: any) => {
			// console.log(respond);
			this.ratinglist = respond;
		}, (err) => {
			console.error(err.error.error)
		})
	}

	// Get All Batch Status
	getBatchStatus(): void {
		this.generalService.getGeneralTypeByType("batch_status").subscribe((respond: any) => {
			this.batchStatus = respond;
		}, (err) => {
			console.error(err.error.error)
		})
	}

	// Get All Batch objectives
	getBatchObjectives(): void {
		this.generalService.getGeneralTypeByType("batch_objective").subscribe((respond: any) => {
			this.batchObjectives = respond;
		}, (err) => {
			console.error(err.error.error)
		})
	}

	timeout: NodeJS.Timeout

	searchTechnology(event?: any): void {
		if (this.timeout) {
			clearTimeout(this.timeout)
		}
		this.timeout = setTimeout(() => {
			this.spinnerService.isDisabled = true
			this.getTechnology(event)
		}, 500)
	}

	// Get All Technology List
	// Change the function call from search & use timeout in that #niranjan 
	getTechnology(event?: any): void {
		// console.log("get tech called")
		let queryParams: any = {}
		if (event && event?.term != "") {
			queryParams.language = event.term
		}
		this.isTechLoading = true
		this.techService.getAllTechnologies(this.techLimit, this.techOffset, queryParams).subscribe((response) => {
			// console.log("getTechnology -> ", response);
			this.technologies = response.body
			// this.technologies = this.technologies.concat(response.body)
		}, (err) => {
			console.error(err)
		}).add(() => {
			this.isTechLoading = false
			this.spinnerService.isDisabled = false
		})
	}

	// Get Academic Year List.
	getAcademicYear(): void {
		this.generalService.getGeneralTypeByType("academic_year").subscribe((respond: any) => {
			this.academicYears = respond;
		}, (err) => {
			console.error(err.error.error)
		})
	}

	// Get All Sales Persons
	getAllSalesPeople(): void {
		this.generalService.getSalesPersonList().subscribe((data: any) => {
			this.salesPeople = data.body
		}, (err) => {
			console.error(err)
		})
	}

	// Get All Faculty List
	getFacultyList(): void {
		this.generalService.getFacultyList().subscribe((data: any) => {
			this.facultyList = data.body
		}, (err) => {
			console.error(err);
		})

	}

	// Get All Course List
	getCourseList(): void {
		this.generalService.getCourseList().subscribe((data: any) => {
			this.courseList = data.body
		}, (err) => {
			console.error(err)
		})
	}

	getDaysList(): void {
		this.generalService.getDaysList().subscribe((response: any) => {
			this.daysList = response.body
			this.formatDayListFields()
		}, (err) => {
			console.error(err);
		})
	}

	getRequirementList(): void {
		this.generalService.getRequirementList().subscribe((response: any) => {
			this.requirementList = response.body
		}, (err: any) => {
			console.error(err);
		})
	}

	geBatchSessionList(): void {
		this.batchService.getBatchSessionList(this.batchID).subscribe((response: any) => {
			this.batchSessionList = response.body
		}, (err: any) => {
			console.error(err);
		})
	}

	getBatchSessionTalents(batchSessionID: string): void {
		this.batchService.getBatchSessionTalents(this.batchID, batchSessionID).subscribe((response: any) => {
			this.talentAttendance = response.body
			// console.log(this.talentAttendance);
		}, (err: any) => {
			console.error(err)
		})


	}

	// ============================================================= OTHER =============================================================

	redirectToBatchFeedback(): void {
		let url: string
		if (this.localService.getJsonValue("roleName") == this.role.ADMIN || this.localService.getJsonValue("roleName") == this.role.SALES_PERSON) {
			url = this.urlConstant.TRAINING_BATCH_MASTER_FEEDBACK
		}
		if (this.localService.getJsonValue("roleName") == this.role.FACULTY) {
			url = this.urlConstant.MY_BATCH_FEEDBACK
		}
		this.router.navigate([url], {
			queryParams: {
				"batchID": this.batchID,
				"batchName": this.batchDetails?.batchName
			}
		}).catch(err => {
			console.error(err)

		});
	}

	redirectToTalentMaster(): void {
		let url: string
		if (this.isFaculty) {
			url = this.urlConstant.MY_TALENT
		}

		if (!this.isFaculty) {
			url = this.urlConstant.TALENT_MASTER
		}

		// console.log(" === url -> ", url);
		this.router.navigate([url], {
			queryParams: {
				"batchID": this.batchID
			}
		}).catch(err => {
			console.error(err)
		});
	}

	// redirectToConceptTree(): void {
	// 	let url: string

	// 	if (this.isFaculty) {
	// 		url = this.FACULTY_CONCEPT_TREE_URL
	// 	}

	// 	if (this.isTalent) {
	// 		url = this.TALENT_CONCEPT_TREE_URL
	// 	}

	// 	if (this.isAdmin || this.isSalesPerson) {
	// 		url = this.ADMIN_CONCEPT_TREE_URL
	// 	}

	// 	// this.router.navigate([url], {
	// 	// 	queryParams: {
	// 	// 		"batchID": this.batchID,
	// 	// 		"courseID": this.courseID,
	// 	// 		"batchName": this.batchDetails?.batchName
	// 	// 	}
	// 	// }).catch(err => {
	// 	// 	console.error(err)
	// 	// });
	// }

	openModal(modalContent: any, modalSize = "xl", withExtraOptions?: boolean): NgbModalRef {
		this.resetBrochureUploadFields()
		this.resetLogoUploadFields()

		let options: NgbModalOptions = {
			ariaLabelledBy: 'modal-basic-title', backdrop: 'static',
			size: modalSize, keyboard: false
		}

		if (withExtraOptions) {
			options.scrollable = true
			options.centered = true
		}
		this.modalRef = this.modalService.open(modalContent, options)
		return this.modalRef
	}



	dismissModal(modal: NgbModalRef): void {
		modal.dismiss()
	}

	closeModal(modal: NgbModalRef): void {
		modal.close()
	}

	// Format the fields of day list.
	formatDayListFields(): void {

		for (let i = 0; i < this.daysList.length; i++) {
			this.daysList[i].isSelected = false
			for (let j = 0; j < this.batchTiming.controls.length; j++) {
				if (this.daysList[i].id == this.batchTiming.at(j).get('day').value.id) {
					this.daysList[i].isSelected = true
					continue
				}
			}
		}
	}

	createBatchTimingForm(index: number, day: any): FormGroup {

		return this.formBuilder.group({
			id: new FormControl(),
			day: new FormControl(day, [Validators.required]),
			fromTime: new FormControl(null,
				[Validators.required]),
			toTime: new FormControl(null,
				[Validators.required])
		})


	}

	// On clicking day button.
	onDayClick(day: any): void {

		// If day is selected then unselect it and remove it from batch timings form.
		if (day.isSelected) {
			if (confirm("Are you sure you want to delete the batch timing?")) {
				for (let i = 0; i < this.batchTiming.controls.length; i++) {
					if (day.id == this.batchTiming.at(i).get('day').value.id) {
						this.batchTiming.removeAt(i)

						// If removed batch timing is first index, then make is apply to all sessions false and enable all other batch timings.
						if (i == 0 && this.batchTiming.length > 1) {
							this.batchForm.get('isApplyToAllSessions').setValue(false)
							for (let j = 0; j < this.batchTiming.length; j++) {
								this.batchTiming.at(j).get('fromTime').enable()
								this.batchTiming.at(j).get('toTime').enable()
							}
						}
						break
					}
				}
				this.batchForm.markAsDirty()
				day.isSelected = false
			}
			return
		}

		// If day is unselected then select it and add it to batch timings form.
		day.isSelected = true
		let index: number = -1
		for (let i = 0; i < this.batchTiming.controls.length; i++) {
			if (day.order > this.batchTiming.at(i).value.day.order) {
				index = i
			}
		}
		this.batchTiming.insert(index + 1, this.createBatchTimingForm(index + 1, day))

		// If apply to all is true and new batch timing is inserted then set the value and disable it.
		// If added batch timing is not in first place then set its time as the first batch timing's value and disable it.
		if (index != -1 && this.batchTiming.length > 1 && this.batchForm.get('isApplyToAllSessions').value) {
			this.batchTiming.at(index + 1).get('fromTime').setValue(this.batchTiming.at(0).get('fromTime').value)
			this.batchTiming.at(index + 1).get('toTime').setValue(this.batchTiming.at(0).get('toTime').value)
			this.batchTiming.at(index + 1).get('fromTime').disable()
			this.batchTiming.at(index + 1).get('toTime').disable()
		}

		// If added batch timing is in first place then set first batch timing's value as the second batch timing's value
		// and disable second batch timing.
		if (index == -1 && this.batchTiming.length > 1 && this.batchForm.get('isApplyToAllSessions').value) {
			this.batchTiming.at(0).get('fromTime').setValue(this.batchTiming.at(1).get('fromTime').value)
			this.batchTiming.at(0).get('toTime').setValue(this.batchTiming.at(1).get('toTime').value)
			this.batchTiming.at(1).get('fromTime').disable()
			this.batchTiming.at(1).get('toTime').disable()
		}
	}

	// Sort batch timings.
	sortBatchTimings(a, b) {
		if (a.day.order < b.day.order) {
			return -1
		}
		if (a.day.order > b.day.order) {
			return 1
		}
		return 0
	}

	// Add batch timing form to batch form.
	addBatchTimingForm(): void {
		this.batchTiming.push(this.formBuilder.group({
			id: new FormControl(),
			day: new FormControl(null, [Validators.required]),
			fromTime: new FormControl(null, [Validators.required]),
			toTime: new FormControl(null, [Validators.required])
		}))
	}

	// On changing value of batch time.
	onBatchTimeChange(index: number, time: any): void {
		// console.log("time -> ", time);

		if (index == 0 && this.batchTiming.length > 1 && this.batchForm.get('isApplyToAllSessions').value) {
			for (let i = 1; i < this.batchTiming.length; i++) {
				this.batchTiming.at(i).get('fromTime').setValue(this.batchTiming.at(0).get('fromTime').value)
				this.batchTiming.at(i).get('toTime').setValue(this.batchTiming.at(0).get('toTime').value)
			}
		}
	}

	// On checking apply it to all sessions.
	onApplyToAllSessionsClick(): void {

		// If checkbox is unselected then enable all the batch timings.
		if (this.batchForm.get('isApplyToAllSessions').value == false) {
			for (let i = 1; i < this.batchTiming.length; i++) {
				this.batchTiming.at(i).get('fromTime').enable()
				this.batchTiming.at(i).get('toTime').enable()
			}
			return
		}

		// If first session day timings are empty then give alert.
		if (this.batchTiming.at(0).value.fromTime == null || this.batchTiming.at(0).value.toTime == null) {
			alert("Please fill in the from time and to time of the first session day")
			this.batchForm.get('isApplyToAllSessions').setValue(false)
			return
		}

		// If start time and end time are same.
		if (this.batchTiming.at(0).value.fromTime == this.batchTiming.at(0).value.toTime) {
			alert("Start time and end time cannot be same")
			this.batchForm.get('isApplyToAllSessions').setValue(false)
			return
		}

		// If there are more than one batch timings then set all the other batch timings as the value for the first 
		// batch timing.


		if (this.batchTiming.length > 1) {
			for (let i = 1; i < this.batchTiming.length; i++) {
				this.batchTiming.at(i).get('fromTime').setValue(this.batchTiming.at(0).get('fromTime').value)
				this.batchTiming.at(i).get('toTime').setValue(this.batchTiming.at(0).get('toTime').value)
				this.batchTiming.at(i).get('fromTime').disable()
				this.batchTiming.at(i).get('toTime').disable()
			}
		}
	}


}

interface IMenuDetails {
	tabName: string
	isActive: boolean
	templateName?: any
	isRedirect?: boolean
	fn?: () => void
}

interface ITabDetails {
	tabName: string
	subTab: IMenuDetails[]
	fn?: (subMenu: IMenuDetails) => void
}

interface IMappedTalentDTO {
	id?: string
	batch: any
	talent: any
	dateOfJoining: string
	sessionsAttendedCount: number
	totalSessionsCount: number
	totalHours: number
	attendedHours: number
	averageRating: number
	totalFeedbacksGiven: number
	showDetails: boolean
	sessionDetails: any[]
}

interface IFeedbackComments {
	question: string
	answers: string[]
}