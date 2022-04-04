import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, Validators, FormControl } from '@angular/forms';
import { LocalService } from 'src/app/service/storage/local.service';
import { ActivatedRoute, Router } from '@angular/router';
import { IPermission } from 'src/app/service/menu/menu.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { DatePipe } from '@angular/common';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { UrlConstant } from 'src/app/service/constant';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { ISearchFilterField, ISeminarDTO, ISeminarTopic, ISeminarTopicDTO, IUpdateMultipleStudent, SeminarService } from 'src/app/service/college/seminar/seminar.service';
import { DegreeService } from 'src/app/service/degree/degree.service';

@Component({
	selector: 'app-seminar',
	templateUrl: './seminar.component.html',
	styleUrls: ['./seminar.component.css']
})
export class SeminarComponent implements OnInit {

	//****************************SEMINAR*************************************** */

	// Components.
	collegeBranchList: any[]
	salesPersonList: any[]
	speakerList: any[]

	// Flags.
	isViewMode: boolean
	isOperationUpdate: boolean

	// Seminar.
	seminarList: ISeminarDTO[]
	seminarForm: FormGroup
	selectedSeminarID: string

	// Pagination.
	limitSeminar: number
	offsetSeminar: number
	currentPageSeminar: number
	totalSeminars: number
	paginationStartSeminar: number
	paginationEndSeminar: number

	// Modal.
	modalRef: any
	@ViewChild('seminarFormModal') seminarFormModal: any
	@ViewChild('deleteSeminarModal') deleteSeminarModal: any

	// Spinner.



	// Search.
	seminarSearchForm: FormGroup
	isSearchedSeminar: boolean
	searchFormValueSeminar: any
	searchSeminarFilterFieldList: ISearchFilterField[]

	// Seminar retaled variables.
	currentYear: number
	activeTab: number
	resumeURL: string

	// Permission.
	permission: IPermission
	roleName: string
	showForSalesPersonLogin: boolean

	// Excel.
	arrayBuffer: any
	file: File
	excelData: any
	talentData: any[]
	errorMsg: string
	talentError: boolean
	upload: boolean
	excelFileName: string

	//****************************SEMINAR TOPICS*************************************** */
	// Components.
	seminarTopicSpeakerList: any[]

	// Flags.
	showTopicForm: boolean
	isOperationTopicUpdate: boolean

	// Topic.
	topicList: ISeminarTopicDTO[]
	topicForm: FormGroup

	// Modal.
	@ViewChild('topicFormModal') topicFormModal: any

	//****************************STUDENT*************************************** */
	// Components.
	stateList: any[]
	countryList: any[]
	qualificationList: any[]
	specializationList: any[]
	academicYearList: any[]
	selectedSeminarCollegeList: any[]

	// Flags.
	showRegisterStudentForm: boolean
	showSeminarTalentRegistrationFields: boolean
	showStudentsList: boolean

	// Student.
	studentList: any[]
	selectedStudentList: string[]
	studentForm: FormGroup

	// Pagination.
	limitStudent: number
	offsetStudent: number
	currentPageStudent: number
	totalStudents: number
	paginationStartStudent: number
	paginationEndStudent: number

	// Modal.
	@ViewChild('studentFormModal') studentFormModal: any

	// Search.
	isSearchedStudent: boolean
	studentSearchForm: FormGroup
	searchFormValueStudent: any
	searchStudentFilterFieldList: ISearchFilterField[]

	// Resume.
	isResumeUploadedToServer: boolean
	isFileUploading: boolean
	docStatus: string
	displayedFileName: string

	//****************************UPDATE MULTIPLE STUDENT*************************************** */
	// Flags.
	multipleSelect: boolean

	// Modal.
	@ViewChild('updateMultipleStudentModal') updateMultipleStudentModal: any

	// Forms.
	updateMultipleStudentForm: FormGroup
	updateMultipleStudentFormValue: IUpdateMultipleStudent

	constructor(
		private formBuilder: FormBuilder,
		public utilityService: UtilityService,
		private localService: LocalService,
		private spinnerService: SpinnerService,
		private generalService: GeneralService,
		private seminarService: SeminarService,
		private urlConstant: UrlConstant,
		private modalService: NgbModal,
		private fileOperationService: FileOperationService,
		private datePipe: DatePipe,
		private router: Router,
		private route: ActivatedRoute,
		private degreeService: DegreeService,
	) {
		this.initializeVariables()
		this.getAllComponents()
	}


	get ongoingOperations() {
		return this.spinnerService.ongoingOperations
	}

	ngOnInit(): void { }

	//Initialize all global variables
	initializeVariables(): void {
		//****************************SEMINAR*************************************** */
		// Components.
		this.seminarList = []
		this.collegeBranchList = []
		this.salesPersonList = []
		this.speakerList = []

		// Flags.
		this.isViewMode = false
		this.isOperationUpdate = false

		// Pagination.
		this.limitSeminar = 5
		this.offsetSeminar = 0
		this.currentPageSeminar = 0

		// Spinner.
		this.spinnerService.loadingMessage = "Getting All Seminars"


		// Search.
		this.isSearchedSeminar = false
		this.searchFormValueSeminar = {}
		this.searchSeminarFilterFieldList = []

		// Seminar retaled variables.
		this.activeTab = 1
		this.currentYear = new Date().getFullYear()

		// Permission.
		this.showForSalesPersonLogin = true
		// Get permissions from menus using utilityService function.
		this.permission = this.utilityService.getPermission(this.urlConstant.COLLEGE_SEMINAR)
		// Get role name for menu for calling their specific apis.
		this.roleName = this.localService.getJsonValue("roleName")
		// Hide features for salesperson.
		if (this.roleName == "Salesperson") {
			this.showForSalesPersonLogin = false
		}

		//****************************TOPIC************************************** */
		// Components.
		this.seminarTopicSpeakerList = []

		// Flags.
		this.showTopicForm = false
		this.isOperationTopicUpdate = false

		// Topic.
		this.topicList = [] as ISeminarTopicDTO[]

		//****************************STUDENT*************************************** */
		// Components.
		this.stateList = []
		this.countryList = []
		this.qualificationList = []
		this.specializationList = []
		this.academicYearList = []
		this.selectedSeminarCollegeList = []

		// Flags.
		this.showRegisterStudentForm = false
		this.showSeminarTalentRegistrationFields = false
		this.showStudentsList = false

		// Student.
		this.selectedStudentList = []

		// Pagination.
		this.limitStudent = 5
		this.offsetStudent = 0
		this.currentPageStudent = 0

		// Search.
		this.isSearchedStudent = false
		this.searchFormValueStudent = {}
		this.searchStudentFilterFieldList = []

		// Resume.
		this.docStatus = ""
		this.displayedFileName = "Select file"
		this.isResumeUploadedToServer = false
		this.isFileUploading = false

		//****************************UPDATE MULTIPLE STUDENT*************************************** */
		// Flags.
		this.multipleSelect = false

		//****************************INITIALIZE FORMS*************************************** */
		this.createSeminarSearchForm()
	}

	//*********************************************CREATE FORMS************************************************************
	// Create new seminar form.
	createSeminarForm(): void {
		this.seminarForm = this.formBuilder.group({
			id: new FormControl(null),
			code: new FormControl({ value: null, disabled: true }),
			seminarName: new FormControl(null, [Validators.required]),
			description: new FormControl(null),
			location: new FormControl(null),
			collegeBranches: new FormControl(Array(), [Validators.required]),
			salesPeople: new FormControl(Array()),
			speakers: new FormControl(Array()),
			totalRegisteredStudents: new FormControl(null),
			totalVisitedStudents: new FormControl(null),
			seminarDate: new FormControl(null, [Validators.required]),
			fromTime: new FormControl(null, [Validators.required]),
			toTime: new FormControl(null, [Validators.required]),
			studentRegistrationLink: new FormControl(null),
			isActive: new FormControl(false),
		})
	}

	// Create new seminar search form.
	createSeminarSearchForm(): void {
		this.seminarSearchForm = this.formBuilder.group({
			seminarName: new FormControl(null),
			fromDate: new FormControl(null),
			toDate: new FormControl(null),
			collegeIDs: new FormControl(null),
			salesPersonIDs: new FormControl(null),
			speakerIDs: new FormControl(null),
		})
	}

	// Create new topic form.
	createTopicForm(): void {
		this.topicForm = this.formBuilder.group({
			id: new FormControl(),
			topicName: new FormControl(null, [Validators.required, Validators.maxLength(100)]),
			date: new FormControl(null, [Validators.required]),
			fromTime: new FormControl(null, [Validators.required]),
			toTime: new FormControl(null, [Validators.required]),
			description: new FormControl(null, [Validators.maxLength(500)]),
			speaker: new FormControl(null),
			seminarID: new FormControl(null)
		})
	}

	// Create new student form.
	createStudentForm(): void {
		this.studentForm = this.formBuilder.group({
			firstName: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z]*$/)]),
			lastName: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z]*$/)]),
			email: new FormControl(null, [Validators.pattern(/^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$/), Validators.required]),
			contact: new FormControl(null, [Validators.pattern(/^[6789]\d{9}$/), Validators.required]),
			address: new FormControl(null, [Validators.pattern(/^[.0-9a-zA-Z\s,-\/]+$/)]),
			city: new FormControl(null, [Validators.pattern(/^[a-zA-Z]+([a-zA-Z ]?)+$/)]),
			state: new FormControl(null),
			country: new FormControl(null),
			pinCode: new FormControl(null, [Validators.pattern(/^[1-9][0-9]{5}$/)]),
			academicYear: new FormControl(null, [Validators.required]),
			isSwabhavTalent: new FormControl(null, [Validators.required]),
			college: new FormControl(null, [Validators.required, Validators.maxLength(200)]),
			percentage: new FormControl(null, [Validators.min(1), Validators.pattern(/(^100([.]0{1,2})?)$|(^\d{1,2}([.]\d{1,2})?)$/i), Validators.required]),
			passout: new FormControl(null, [Validators.min(1980), Validators.max(this.currentYear + 3), Validators.required]),
			degree: new FormControl(null, Validators.required),
			specialization: new FormControl(null, Validators.required),
			talentID: new FormControl(null),
			seminarID: new FormControl(null),
			seminarTalentRegistrationID: new FormControl(null),
			talentAcademicID: new FormControl(null),
			resume: new FormControl(null),
			hasVisited: new FormControl(false, [Validators.required])
		})
	}

	// Create new update multiple students form.
	createUpdateMultipleStudentForm(): void {
		this.updateMultipleStudentForm = this.formBuilder.group({
			hasVisited: new FormControl(null),
			seminarID: new FormControl(null),
			seminarTalentRegistrationIDs: new FormControl(null),
		})
	}

	// Create student search form.
	createStudentSearchForm(): void {
		this.studentSearchForm = this.formBuilder.group({
			firstName: new FormControl(null, [Validators.pattern(/^[a-zA-Z]*$/)]),
			lastName: new FormControl(null, [Validators.pattern(/^[a-zA-Z]*$/)]),
			email: new FormControl(null),
			contact: new FormControl(null),
			degreeID: new FormControl(null),
			specializationID: new FormControl(null),
			collegeID: new FormControl(null),
			academicYear: new FormControl(null),
			isSwabhavTalent: new FormControl(null),
			hasVisited: new FormControl(null),
			fromDate: new FormControl(null),
			toDate: new FormControl(null),
		})
	}

	//*********************************************CRUD FUNCTIONS FOR SEMINAR FORM************************************************************
	// On add new seminar button click.
	onAddNewSeminarButtonClick(): void {
		this.isViewMode = false
		this.isOperationUpdate = false
		this.createSeminarForm()
		this.enableSeminarForm()
		this.openModal(this.seminarFormModal, 'xl')
	}

	//On clicking add seminar in seminar form.
	addSeminar(): void {
		this.spinnerService.loadingMessage = "Adding Seminar"


		this.seminarService.addSeminar(this.seminarForm.value).subscribe((response: any) => {
			this.modalRef.close('success')
			this.getSeminars()
			alert(response)
		}, (error) => {
			console.error(error)
			if (error.error?.error) {
				alert(error.error?.error)
				return
			}
			alert(error.statusText)
		})
	}

	// On clicking view seminar button.
	onViewSeminarClick(seminar: ISeminarDTO): void {
		this.isViewMode = true
		this.createSeminarForm()
		// Seminar date.
		let seminarDate = seminar.seminarDate
		if (seminarDate) {
			seminar.seminarDate = this.datePipe.transform(seminarDate, 'yyyy-MM-dd')
		}
		this.seminarForm.patchValue(seminar)
		this.seminarForm.disable()
		this.openModal(this.seminarFormModal, 'xl')
	}

	// On cliking update form button in seminar form.
	onUpdateSeminarClick(): void {
		this.isViewMode = false
		this.isOperationUpdate = true
		this.enableSeminarForm()
	}

	// On clicking update seminar in seminar form.
	updateSeminar(): void {
		this.spinnerService.loadingMessage = "Updating Seminar"


		this.seminarService.updateSeminar(this.seminarForm.value).subscribe((response: any) => {
			this.modalRef.close()
			this.getSeminars()
			alert(response)
		}, (error) => {
			console.error(error)
			if (error.error?.error) {
				alert(error.error.error)
				return
			}
			alert("Check connection")
		})
	}

	// On clicking delete seminar button. 
	onDeleteSeminarClick(seminarID: string): void {
		this.openModal(this.deleteSeminarModal, 'md').result.then(() => {
			this.deleteSeminar(seminarID)
		}, (err) => {
			console.error(err)
			return
		})
	}

	// Delete seminar.
	deleteSeminar(seminarID: string): void {
		this.spinnerService.loadingMessage = "Deleting Seminar"


		this.seminarService.deleteSeminar(seminarID).subscribe((response: any) => {
			this.modalRef.close()
			this.getSeminars()
			alert(response)
		}, (error) => {
			console.error(error)
			if (error.error) {
				alert(error.error)
				return
			}
			alert(error.statusText)
		})
	}

	// ==================================================SEMINAR SEARCH FUNCTIONS==========================================================================
	// Reset seminar search form and renaviagte page.
	resetSearchAndGetAllForSeminar(): void {
		this.searchSeminarFilterFieldList = []
		this.seminarSearchForm.reset()
		this.searchFormValueSeminar = {}
		this.changePageForSeminar(1)
		this.isSearchedSeminar = false
		this.router.navigate([this.urlConstant.COLLEGE_SEMINAR])
	}

	// Reset seminar search form.
	resetSeminarSearchForm(): void {
		this.searchSeminarFilterFieldList = []
		this.seminarSearchForm.reset()
	}

	// Search seminars.
	searchSeminars(): void {
		this.searchFormValueSeminar = { ...this.seminarSearchForm?.value }
		this.router.navigate([], {
			relativeTo: this.route,
			queryParams: this.searchFormValueSeminar,
		})
		for (let field in this.searchFormValueSeminar) {
			if (this.searchFormValueSeminar[field] === null || this.searchFormValueSeminar[field] === "") {
				delete this.searchFormValueSeminar[field]
			} else {
				this.isSearchedSeminar = true
			}
		}
		this.searchSeminarFilterFieldList = []
		for (var property in this.searchFormValueSeminar) {
			let text: string = property
			let result: string = text.replace(/([A-Z])/g, " $1");
			let finalResult: string = result.charAt(0).toUpperCase() + result.slice(1);
			let valueArray: any[] = []
			if (Array.isArray(this.searchFormValueSeminar[property])) {
				valueArray = this.searchFormValueSeminar[property]
			}
			else {
				valueArray.push(this.searchFormValueSeminar[property])
			}
			this.searchSeminarFilterFieldList.push(
				{
					propertyName: property,
					propertyNameText: finalResult,
					valueList: valueArray
				})
		}
		if (this.searchSeminarFilterFieldList.length == 0) {
			this.resetSearchAndGetAllForSeminar()
		}
		if (!this.isSearchedSeminar) {
			return
		}
		this.spinnerService.loadingMessage = "Searching Seminars"
		this.changePageForSeminar(1)
	}

	// Delete search criteria from seminar search form by search name.
	deleteSearchSeminarCriteria(searchName: string): void {
		this.seminarSearchForm.get(searchName).setValue(null)
		this.searchSeminars()
	}

	//*********************************************OTHER FUNCTIONS FOR SEMINAR************************************************************
	// Enable the seminar form.
	enableSeminarForm(): void {
		this.seminarForm.enable()
		this.seminarForm.get('code').disable()
		this.seminarForm.get('totalRegisteredStudents').disable()
		this.seminarForm.get('totalVisitedStudents').disable()
	}

	// Change page for pagination for seminar.
	changePageForSeminar($event): void {
		this.offsetSeminar = $event - 1
		this.currentPageSeminar = $event
		this.getSeminars()
	}

	// On clicking sumbit button in seminar form.
	onSeminarFormSubmit(): void {
		if (this.seminarForm.invalid) {
			this.seminarForm.markAllAsTouched()
			return
		}
		if (this.isOperationUpdate) {
			this.updateSeminar()
			return
		}
		this.addSeminar()
	}

	// Redirect to talents page filtered by seminar id.
	redirectToTalentsForRegisteredStudents(seminarID: string): void {
		this.router.navigate([this.urlConstant.TALENT_MASTER], {
			queryParams: {
				"seminarID": seminarID,
				"hasVisited": "0"
			}
		}).catch(err => {
			console.error(err)

		})
	}

	// Redirect to talents page filtered by seminar id and has appeared fields.
	redirectToTalentsForVisitedStudents(seminar: any): void {
		this.router.navigate([this.urlConstant.TALENT_MASTER], {
			queryParams: {
				"seminarID": seminar.id,
				"hasVisited": "1"
			}
		}).catch(err => {
			console.error(err)

		})
	}

	// Checks the url's query params and decides to call whether to call get or search.
	searchOrGetSeminars() {
		let queryParams = this.route.snapshot.queryParams
		if (this.utilityService.isObjectEmpty(queryParams)) {
			this.getSeminars()
			return
		}
		this.seminarSearchForm.patchValue(queryParams)
		this.searchSeminars()
	}

	// Set total list on current page.
	setPaginationStringForSeminar(): void {
		this.paginationStartSeminar = this.limitSeminar * this.offsetSeminar + 1
		this.paginationEndSeminar = +this.limitSeminar + this.limitSeminar * this.offsetSeminar
		if (this.totalSeminars < this.paginationEndSeminar) {
			this.paginationEndSeminar = this.totalSeminars
		}
	}

	//*********************************************CRUD FUNCTIONS FOR STUDENT FORM************************************************************
	// On clicking register students button.
	getStudentsBySelectedSeminar(seminarID: any): void {
		this.selectedSeminarID = seminarID
		this.getSelectedSeminarCollegeList()
		this.createStudentSearchForm()
		this.initializeStudentVariables()
		this.getAllStudents()
		this.openModal(this.studentFormModal, 'xl')
	}

	// Get all students by seminar id.
	getAllStudents(): void {
		this.spinnerService.loadingMessage = "Getting All Students"


		this.seminarService.getStudentsBySeminar(this.selectedSeminarID,
			this.limitStudent, this.offsetStudent, this.searchFormValueStudent).subscribe(response => {
				this.studentList = response.body
				this.formatAddress()
				this.totalStudents = parseInt(response.headers.get("X-Total-Count"))
				this.getSeminars()
			}, err => {
				console.error(this.utilityService.getErrorString(err))
			}).add(() => {
				this.setPaginationStringForStudent()

			})
	}

	// On clicking add new student button.
	onAddNewStudentButtonClick(): void {
		this.createStudentForm()
		this.isViewMode = false
		this.isOperationUpdate = false
		this.showSeminarTalentRegistrationFields = false
		this.stateList = []
		this.studentForm.get('state').disable()
		this.specializationList = []
		this.studentForm.get('specialization').disable()
		this.toggleVisibilityOfCanditateFormAndList()
		this.activeTab = 1
	}

	// On clicking add student button on student form.
	addStudent(): void {
		this.spinnerService.loadingMessage = "Adding Student"


		let student: any = this.studentForm.value
		this.patchIDFromObjectsForStudent(student)
		this.seminarService.addStudent(this.studentForm.value, this.selectedSeminarID).subscribe((response: any) => {
			this.toggleVisibilityOfCanditateFormAndList()
			this.getAllStudents()
			alert(response)
		}, (error) => {
			console.error(error)
			if (typeof error.error == 'object' && error) {
				alert(this.utilityService.getErrorString(error))
				return
			}
			if (error.error == undefined) {
				alert('Student could not be added, try again')
				return
			}
			alert(error.statusText)
		})
	}

	// On clicking view student button.
	OnViewStudentButtonClick(index: number): void {
		this.createStudentForm()
		this.isViewMode = true
		this.addSeminarTalentRegsitrationFields()
		this.showSeminarTalentRegistrationFields = true
		this.toggleVisibilityOfCanditateFormAndList()
		if (this.studentList[index]?.registrationDate) {
			this.studentList[index].registrationDate = this.datePipe.transform(this.studentList[index]?.registrationDate, 'yyyy-MM-dd')
		}
		this.studentForm.patchValue(this.studentList[index])
		this.resumeURL = this.studentForm.get('resume').value
		// Resume.
		this.displayedFileName = "No resume uploaded"
		if (this.studentList[index].resume) {
			this.displayedFileName = `<a href=${this.studentList[index].resume} target="_blank">Resume present</a>`
		}
		this.stateList = []
		this.specializationList = []
		if (this.studentList[index].country != undefined) {
			this.getStateListByCountry(this.studentList[index].country)
		}
		if (this.studentList[index].degree != undefined) {
			this.getSpecializationListByDegree(this.studentList[index].degree)
		}
		this.studentForm.disable()
	}

	// On clicking update student button in view form.
	onUpdateStudentFormButtonClick(): void {
		this.isViewMode = false
		this.isOperationUpdate = true
		this.studentForm.enable()
		if (this.studentForm.get('country').value == null) {
			this.studentForm.get('state').disable()
		}
	}

	// Add seminar talent registration fields to form in update form.
	addSeminarTalentRegsitrationFields(): void {
		this.studentForm.addControl('registrationDate', this.formBuilder.control(null, [Validators.required]))
	}

	// On clicking update student in candodate form.
	updateStudent(): void {
		this.spinnerService.loadingMessage = "Updating Student"


		let student: any = this.studentForm.value
		this.patchIDFromObjectsForStudent(student)
		this.seminarService.updateStudent(this.studentForm.value, this.selectedSeminarID).subscribe((response: any) => {
			this.toggleVisibilityOfCanditateFormAndList()
			this.getAllStudents()
			alert(response)
		}, (error) => {
			console.error(error)
			if (typeof error.error == 'object' && error) {
				alert(this.utilityService.getErrorString(error))
				return
			}
			if (error.error == undefined) {
				alert('Student could not be added, try again')
				return
			}
			alert(error.statusText)
		})
	}

	// Delete student.
	deleteStudent(seminarTalentRegistrationID: string, resume: string): void {
		if (confirm("Are you sure you want to delete the student?")) {
			this.spinnerService.loadingMessage = "Deleting Student"


			this.seminarService.deleteStudent(seminarTalentRegistrationID, this.selectedSeminarID).subscribe((response: any) => {
				this.getAllStudents()
				this.getSeminars()
				this.fileOperationService.deleteUploadedFile(resume)
				alert(response)
			}, (error) => {
				console.error(error)
				if (typeof error.error == 'object' && error) {
					alert(this.utilityService.getErrorString(error))
					return
				}
				if (error.error == undefined) {
					alert('Student could not be added, try again')
					return
				}
				alert(error.statusText)
			})
		}
	}

	//Delete resume
	deleteResume(): void {
		this.fileOperationService.deleteUploadedFile().subscribe((data: any) => {
		}, (error) => {
			console.error(error)
		})
	}

	// ==================================================STUDENT SEARCH FUNCTIONS==========================================================================
	// Reset student search form and renaviagte page.
	resetSearchAndGetAllForStudent(): void {
		this.searchStudentFilterFieldList = []
		this.studentSearchForm.reset()
		this.searchFormValueStudent = {}
		this.changePageForStudent(1)
		this.isSearchedStudent = false
	}

	// Reset student search form.
	resetStudentSearchForm(): void {
		this.searchStudentFilterFieldList = []
		this.studentSearchForm.reset()
	}

	// Search students.
	searchStudents(): void {
		this.searchFormValueStudent = { ...this.studentSearchForm?.value }
		for (let field in this.searchFormValueStudent) {
			if (this.searchFormValueStudent[field] === null || this.searchFormValueStudent[field] === "") {
				delete this.searchFormValueStudent[field]
			} else {
				this.isSearchedStudent = true
			}
		}
		this.searchStudentFilterFieldList = []
		for (var property in this.searchFormValueStudent) {
			let text: string = property
			let result: string = text.replace(/([A-Z])/g, " $1");
			let finalResult: string = result.charAt(0).toUpperCase() + result.slice(1);
			let valueArray: any[] = []
			if (Array.isArray(this.searchFormValueStudent[property])) {
				valueArray = this.searchFormValueStudent[property]
			}
			else {
				valueArray.push(this.searchFormValueStudent[property])
			}
			this.searchStudentFilterFieldList.push(
				{
					propertyName: property,
					propertyNameText: finalResult,
					valueList: valueArray
				})
		}
		if (this.searchStudentFilterFieldList.length == 0) {
			this.resetSearchAndGetAllForStudent()
		}
		if (!this.isSearchedStudent) {
			return
		}
		this.spinnerService.loadingMessage = "Searching Students"
		this.changePageForStudent(1)
	}

	// Delete search criteria from student search form by search name.
	deleteSearchStudentCriteria(searchName: string): void {
		this.studentSearchForm.get(searchName).setValue(null)
		this.searchStudents()
	}

	//*********************************************OTHER FUNCTIONS FOR STUDENT************************************************************
	// Make all address fields into one comma separated string.
	formatAddress(): void {
		let addressArray: string[] = []
		for (let i = 0; i < this.studentList.length; i++) {
			addressArray = []
			if (this.studentList[i].address) {
				addressArray.push(this.studentList[i].address)
			}
			if (this.studentList[i].pinCode) {
				addressArray.push(this.studentList[i].pinCode)
			}
			if (this.studentList[i].city) {
				addressArray.push(this.studentList[i].city)
			}
			if (this.studentList[i].state) {
				addressArray.push(this.studentList[i].state.name)
			}
			if (this.studentList[i].country) {
				addressArray.push(this.studentList[i].country.name)
			}
			if (addressArray.length > 1) {
				this.studentList[i].fullAddress = addressArray.join(", ")
			}
			if (addressArray.length == 1) {
				this.studentList[i].fullAddress = addressArray[0]
			}
		}
	}

	// On clicking sumbit button in student form.
	onStudentFormSubmit(): void {
		if (this.studentForm.invalid) {
			this.studentForm.markAllAsTouched()
			return
		}
		if (this.isOperationUpdate) {
			this.updateStudent()
			return
		}
		this.addStudent()
	}

	// Toggle visibility of student list and student form.
	toggleVisibilityOfCanditateFormAndList(): void {
		if (this.showRegisterStudentForm) {
			this.showRegisterStudentForm = false
			this.showStudentsList = true
			return
		}
		this.showRegisterStudentForm = true
		this.showStudentsList = false
	}

	// Change page for pagination.
	changePageForStudent($event): void {
		this.offsetStudent = $event - 1
		this.currentPageStudent = $event
		this.getAllStudents()
	}

	// On closing the student form.
	onCloseStudentForm(): void {
		if (this.isFileUploading) {
			alert("Please wait till file is being uploaded")
			return
		}
		if (this.isResumeUploadedToServer) {
			if (!confirm("Uploaded resume will be deleted.\nAre you sure you want to close?")) {
				return
			}
			this.deleteResume()
		}
		this.isResumeUploadedToServer = false
		this.displayedFileName = "Select file"
		this.docStatus = ""
		this.toggleVisibilityOfCanditateFormAndList()
		this.isViewMode = false
		this.isOperationUpdate = false
	}

	// Initialize all variables related to student.
	initializeStudentVariables(): void {
		this.studentList = []
		this.spinnerService.loadingMessage = "Getting All Students"
		this.showRegisterStudentForm = false
		this.showStudentsList = true
		this.limitStudent = 5
		this.offsetStudent = 0
		this.currentPageStudent = 0
		this.selectedStudentList = []
		this.multipleSelect = false
		this.isViewMode = false
		this.isOperationUpdate = false
		this.isSearchedStudent = false
		this.searchFormValueStudent = {}
	}

	// On uplaoding resume.
	onResourceSelect(event: any): void {
		this.docStatus = ""
		let files = event.target.files
		if (files && files.length) {
			let file = files[0]

			// Upload resume if it is present.]
			this.isFileUploading = true
			this.fileOperationService.uploadResume(file).subscribe((data: any) => {
				this.studentForm.markAsDirty()
				this.studentForm.patchValue({
					resume: data
				})
				this.displayedFileName = file.name
				this.isFileUploading = false
				this.isResumeUploadedToServer = true
				this.docStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
			}, (error) => {
				this.isFileUploading = false
				this.docStatus = `<p><span>&#10060;</span> ${error}</p>`
			})
		}
	}

	// Will be called on [addTag]= "addCollegeToList" args will be passed automatically.
	addCollegeToList(option: any): Promise<any> {
		return new Promise((resolve) => {
			resolve(option)
		}
		)
	}

	// Extract ID from objects in student form.
	patchIDFromObjectsForStudent(Student: any): void {
		if (this.studentForm.get('degree').value) {
			Student.degreeID = this.studentForm.get('degree').value.id
			delete Student['degree']
		}
		if (this.studentForm.get('specialization').value) {
			Student.specializationID = this.studentForm.get('specialization').value.id
			delete Student['specialization']
		}
	}

	// Set total list on current page.
	setPaginationStringForStudent(): void {
		this.paginationStartSeminar = this.limitSeminar * this.offsetSeminar + 1
		this.paginationEndSeminar = +this.limitSeminar + this.limitSeminar * this.offsetSeminar
		if (this.totalSeminars < this.paginationEndSeminar) {
			this.paginationEndSeminar = this.totalSeminars
		}
	}

	//*********************************************OTHER FUNCTIONS************************************************************
	// Used to open modal.
	openModal(content: any, size?: string): NgbModalRef {
		if (!size) {
			size = 'lg'
		}
		let options: NgbModalOptions = {
			ariaLabelledBy: 'modal-basic-title', keyboard: false,
			backdrop: 'static', size: size
		}
		this.modalRef = this.modalService.open(content, options)
		return this.modalRef
	}

	//Compare for select option field.
	compareFn(optionOne: any, optionTwo: any): boolean {
		if (optionOne == null && optionTwo == null) {
			return true
		}
		if (optionTwo != undefined && optionOne != undefined) {
			return optionOne.id === optionTwo.id
		}
		return false
	}

	//*********************************************CRUD FUNCTIONS FOR TOPIC************************************************************
	// On clicking topic popup in seminar list.
	getTopicsForSelectedSeminar(seminarID: any): void {
		this.isOperationTopicUpdate = false
		this.topicList = []
		this.selectedSeminarID = seminarID
		this.setTopicSpeakerList()
		this.getAllTopics()
		this.showTopicForm = false
		this.openModal(this.topicFormModal, 'xl')
	}

	// Get all topics by enquiry id.
	getAllTopics(): void {
		this.spinnerService.loadingMessage = "Getting All Topics"


		this.seminarService.getTopicsBySeminar(this.selectedSeminarID).subscribe(response => {
			this.topicList = response
			this.formatDateTimeOfTopicList()
		}, err => {
			console.error(this.utilityService.getErrorString(err))
		})
	}

	// On clicking add new topic button.
	onAddNewTopicButtonClick(): void {
		this.isOperationTopicUpdate = false
		this.createTopicForm()
		this.showTopicForm = true
	}

	// Add topic.
	addTopic(): void {
		this.spinnerService.loadingMessage = "Adding Topic"


		let topic: ISeminarTopic = this.topicForm.value
		this.patchIDFromObjectsForTopic(topic)
		this.seminarService.addTopic(this.topicForm.value, this.selectedSeminarID).subscribe((response: any) => {
			this.showTopicForm = false
			this.getAllTopics()
			alert(response)
		}, (error) => {
			console.error(error)
			if (error.error?.error) {
				alert(error.error?.error)
				return
			}
			alert(error.statusText)
		})
	}

	// Update topic form.
	OnUpdateTopicButtonClick(index: number): void {
		this.isOperationTopicUpdate = true
		this.createTopicForm()
		this.showTopicForm = true
		this.topicForm.patchValue(this.topicList[index])
	}

	// Update topic.
	updateTopic(): void {
		this.spinnerService.loadingMessage = "Updating Topic"


		let topic: ISeminarTopic = this.topicForm.value
		this.patchIDFromObjectsForTopic(topic)
		this.seminarService.updateTopic(this.topicForm.value, this.selectedSeminarID).subscribe((response: any) => {
			this.showTopicForm = false
			this.getAllTopics()
			alert(response)
		}, (error) => {
			console.error(error)
			if (error.error?.error) {
				alert(error.error.error)
				return
			}
			alert("Check connection")
		})
	}

	// Delete topic.
	deleteTopic(topicID: string): void {
		if (confirm("Are you sure you want to delete the topic?")) {
			this.spinnerService.loadingMessage = "Deleting Topic"


			this.seminarService.deleteTopic(topicID, this.selectedSeminarID).subscribe((response: any) => {
				this.getAllTopics()
				alert(response)
			}, (error) => {
				console.error(error)
				if (error.error) {
					alert(error.error)
					return
				}
				alert(error.statusText)
			})
		}
	}

	// Validate topic form.
	validateTopicForm(): void {
		if (this.topicForm.invalid) {
			this.topicForm.markAllAsTouched()
			return
		}
		if (this.isOperationTopicUpdate) {
			this.updateTopic()
			return
		}
		this.addTopic()
	}

	// Format date time field of topic list by removing timestamp.
	formatDateTimeOfTopicList(): void {
		for (let i = 0; i < this.topicList?.length; i++) {
			let date = this.topicList[i].date
			if (date) {
				this.topicList[i].date = this.datePipe.transform(date, 'yyyy-MM-dd')
			}
		}
	}

	// Extract ID from objects in topic form.
	patchIDFromObjectsForTopic(topic: ISeminarTopic): void {
		if (this.topicForm.get('speaker').value) {
			topic.speakerID = this.topicForm.get('speaker').value.id
			delete topic['speaker']
		}
	}

	//*********************************************UPDATE MULTIPLE STUDENT FUNCTIONS************************************************************
	// On clicking update multiple students.
	OnClickingUpdateMultipleStudents(): void {
		this.createUpdateMultipleStudentForm()
		if (this.selectedStudentList.length == 0) {
			alert("Please select students to be updated")
			return
		}
		this.openModal(this.updateMultipleStudentModal, 'sm')
	}

	// Toggle visibility of multiple select checkbox.
	toggleMultipleSelect(): void {
		if (this.multipleSelect) {
			this.multipleSelect = false
			this.setSelectAllStudents(this.multipleSelect)
			return
		}
		this.multipleSelect = true
	}

	// Set isChecked field of all selected students.
	setSelectAllStudents(isSelectedAll: boolean): void {
		for (let i = 0; i < this.studentList.length; i++) {
			this.addStudentToList(isSelectedAll, this.studentList[i])
		}
	}

	// Check if all students in selected students are added in multiple select or not.
	checkStudentsAdded(): boolean {
		let count: number = 0

		for (let i = 0; i < this.studentList.length; i++) {
			if (this.selectedStudentList.includes(this.studentList[i].id))
				count = count + 1
		}
		return (count == this.studentList.length)
	}

	// Check if student is added in multiple select or not.
	checkStudentAdded(seminarTalentRegistrationID: string): boolean {
		return this.selectedStudentList.includes(seminarTalentRegistrationID)
	}

	// Takes a list called selectedStudentList and adds all the checked students to list, also does not contain duplicate values.
	addStudentToList(isChecked: boolean, student: any): void {
		if (isChecked) {
			if (!this.selectedStudentList.includes(student.seminarTalentRegistrationID)) {
				this.selectedStudentList.push(student.seminarTalentRegistrationID)
			}
			return
		}
		if (this.selectedStudentList.includes(student.seminarTalentRegistrationID)) {
			let index = this.selectedStudentList.indexOf(student.seminarTalentRegistrationID)
			this.selectedStudentList.splice(index, 1)
		}
	}

	// Update fields of multiple students.
	updateMultipleStudentFunction(): void {
		this.updateMultipleStudentFormValue = { ...this.updateMultipleStudentForm.value }
		// Check if all fields are null.
		let flag: boolean = true
		for (let field in this.updateMultipleStudentFormValue) {
			if (this.updateMultipleStudentFormValue[field] == null) {
				delete this.updateMultipleStudentFormValue[field]
			} else {
				flag = false
			}
		}
		// No API call on empty search.
		if (flag) {
			alert("Please select any one value to be updated")
			return
		}
		this.updateMultipleStudentFormValue.seminarTalentRegistrationIDs = this.selectedStudentList
		this.spinnerService.loadingMessage = "Updating Students"


		this.updateMultipleStudentFormValue.seminarID = this.selectedSeminarID
		this.seminarService.updateMultipleStudent(this.updateMultipleStudentFormValue).subscribe((response: any) => {
			this.getAllStudents()
			this.getSeminars()
			alert(response)
			this.modalRef.close('success')
			this.selectedStudentList = []
		}, (error) => {
			console.error(error)
			if (typeof error.error == 'object' && error) {
				alert(this.utilityService.getErrorString(error))
				return
			}
			if (error.error == undefined) {
				alert('Students could not be updated')
			}
			alert(error.statusText)
		})
	}

	//*********************************************GET FUNCTIONS************************************************************

	// Get all lists.
	getAllComponents(): void {
		this.getCountryList()
		this.getQualificationList()
		this.getSalesPersonList()
		this.getSpeakerList()
		this.getAcademicYear()
		this.getCollegeBranchList()
		this.searchOrGetSeminars()
	}

	// Get academic year list.
	getAcademicYear(): void {
		this.generalService.getGeneralTypeByType("academic_year").subscribe((response: any[]) => {
			this.academicYearList = response
		}, (err) => {
			console.error(this.utilityService.getErrorString(err))
		})
	}

	// Get all seminars by limit and offset.
	getSeminars(): void {
		this.spinnerService.loadingMessage = "Getting Seminars"


		this.searchFormValueSeminar.roleName = this.roleName
		this.searchFormValueSeminar.loginID = this.localService.getJsonValue("loginID")
		this.seminarService.getSeminars(this.limitSeminar, this.offsetSeminar, this.searchFormValueSeminar).subscribe(response => {
			this.seminarList = response.body
			this.totalSeminars = parseInt(response.headers.get("X-Total-Count"))
		}, err => {
			console.error(this.utilityService.getErrorString(err))
		}).add(() => {
			this.setPaginationStringForSeminar()
		})
	}

	// Get state list by country.
	getStateListByCountry(country: any): void {
		if (country == null) {
			this.stateList = []
			this.studentForm.get('state').setValue(null)
			this.studentForm.get('state').disable()
			return
		}
		if (country.name == "U.S.A." || country.name == "India") {
			this.studentForm.get('state').enable()
			this.generalService.getStatesByCountryID(country.id).subscribe((respond: any[]) => {
				this.stateList = respond
			}, (err) => {
				console.error(err)
			})
		}
		else {
			this.stateList = []
			this.studentForm.get('state').setValue(null)
			this.studentForm.get('state').disable()
		}
	}

	// Get specialization list by degree.
	getSpecializationListByDegree(degree: any): void {
		if (degree == null) {
			this.specializationList = []
			this.studentForm.get('specialization').setValue(null)
			this.studentForm.get('specialization').disable()
			return
		}
		this.studentForm.get('specialization').enable()
		this.generalService.getSpecializationByDegreeID(degree.id).subscribe((response: any) => {
			this.specializationList = response.body
		}, (err) => {
			console.error(err)
		})
	}

	// Get specialization list for students search.
	getSpecializationListForStudentSearch(degreeID: string): void {
		if (degreeID == null) {
			this.specializationList = []
			return
		}
		for (let i = 0; i < this.qualificationList.length; i++) {
			if (this.qualificationList[i].id == degreeID) {
				this.generalService.getSpecializationByDegreeID(degreeID).subscribe(response => {
					this.specializationList = response.body
				}, err => {
					console.error(this.utilityService.getErrorString(err))
				})
				return
			}
		}
	}

	// Get country list.
	getCountryList(): void {
		this.generalService.getCountries().subscribe((respond: any[]) => {
			this.countryList = respond
		}, (err) => {
			console.error(err.error)
		})
	}

	// Get qualification list.
	getQualificationList(): void {
		let queryParams: any = {
			limit: -1,
			offset: 0,
		}
		this.degreeService.getAllDegrees(queryParams).subscribe((respond: any) => {
			//respond.push(this.other)
			this.qualificationList = respond.body
		}, (err) => {
			console.error(this.utilityService.getErrorString(err))
		})
	}

	// Get college branch list.
	getCollegeBranchList(): void {
		this.generalService.getCollegeBranchList().subscribe(response => {
			this.collegeBranchList = response
		}, err => {
			console.error(this.utilityService.getErrorString(err))
		})
	}

	// Get selected seminar's college list.
	getSelectedSeminarCollegeList(): void {
		for (let i = 0; i < this.seminarList.length; i++) {
			if (this.selectedSeminarID == this.seminarList[i].id) {
				this.selectedSeminarCollegeList = this.seminarList[i].collegeBranches
			}
		}
	}

	// Get salesperson list.
	getSalesPersonList(): void {
		this.generalService.getSalesPersonList().subscribe(
			response => {
				this.salesPersonList = response.body
			}), err => {
				console.error(err)
			}
	}

	// Get speaker list.
	getSpeakerList(): void {
		this.generalService.getSpeakerList().subscribe(
			response => {
				this.speakerList = response
			}), err => {
				console.error(err)
			}
	}

	// Set the topic speaker list as the speaker list of seminar.
	setTopicSpeakerList(): void {
		for (let i = 0; i < this.seminarList.length; i++) {
			if (this.selectedSeminarID == this.seminarList[i].id) {
				this.seminarTopicSpeakerList = this.seminarList[i].speakers
			}
		}
	}
}
