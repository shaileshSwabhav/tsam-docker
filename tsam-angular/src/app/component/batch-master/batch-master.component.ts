import { Component, OnInit, ViewChild } from '@angular/core';
import { Validators, FormBuilder, FormGroup, FormArray, FormControl, AbstractControl } from '@angular/forms';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { BatchService, IBatch, IBatchModule, IBatchTiming, IMappedTalent } from 'src/app/service/batch/batch.service';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { LocalService } from 'src/app/service/storage/local.service';
import { ActivatedRoute, Router } from '@angular/router';
import { Role, UrlConstant } from 'src/app/service/constant';
import { DatePipe } from '@angular/common';
import { TechnologyService } from 'src/app/service/technology/technology.service';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { TalentService } from 'src/app/service/talent/talent.service';
import { CourseModuleService, ICourseModule } from 'src/app/service/course-module/course-module.service';
import { UrlService } from 'src/app/service/url.service';
import { AdminService } from 'src/app/service/admin/admin.service';

@Component({
	selector: 'app-batch-master',
	templateUrl: './batch-master.component.html',
	styleUrls: ['./batch-master.component.css']
})
export class BatchMasterComponent implements OnInit {

	// listing
	ratinglist: any[];
	salesPeople: any[];
	technologies: any[];
	academicYears: any[];
	technologyList: any[];
	courseList: any[];
	courseSessionList: any[]
	daysList: any[]
	facultyList: any[];
	isActive: string[]
	batchStatus: string[]
	batchObjectives: string[]
	requirementList: any[]
	activeBatchList: any[]

	// tech component
	isTechLoading: boolean
	techLimit: number
	techOffset: number

	// form
	batchForm: FormGroup;
	searchBatchForm: FormGroup;
	searchedTalentForm: FormGroup

	// pagination
	limit: number;
	offset: number;
	currentPage: number;
	totalTalent: number;
	talentLimit: number;
	talentOffset: number;
	currentTalentPage: number;
	totalSelectedTalent: number
	selectedTalentLimit: number
	selectedTalentOffset: number
	currentSelectedTalentPage: number
	paginationString: string
	talentPaginationString: string
	selectedTalentPaginationString: string

	modalHandler: () => void;
	formHandler: (index?: number) => void;

	// modal handlers
	modalHeadingLabel: string;
	modalButtonLabel: string;
	isViewClicked: boolean;
	modalRef: NgbModalRef;
	entity: string
	waitingListMessage: string
	isAddClicked: boolean

	// batches
	totalBatch: number;
	batches: any[];
	selectedBatch: any;
	batchID: string
	supervisorFacultyCount: number
	isHeadFcaulty: boolean
	isViewAllBatches: boolean

	// talent
	student: any[];
	talentID: string
	toggleTalentSelect: boolean
	talents: [];
	addSelectedTalents: any[]
	selectedTalent: any[] = []
	isTalentAssignedOrDeleted: boolean
	totalBatchTalents: number
	eligibilityID: string

	// access
	permission: IPermission
	loginID: string

	// booleans
	isSearched: boolean
	doesEligibilityExists: boolean
	disableButton: boolean
	searchedTalent: boolean
	showSearch: boolean

	isSalesPerson: boolean
	isFaculty: boolean
	isTalent: boolean
	isAdmin: boolean

	//spinner

	selecetdrRequirementForTransfer: any
	selectedBatchForTransfer: any

	// search
	searchFormValue: any

	// booleans
	isBatchLoaded: boolean
	isTalentLoaded: boolean
	isSelectedTalentLoaded: boolean
	showTransferButton: boolean
	showTransferForm: boolean
	showCompanyReqField: boolean
	showBatchField: boolean

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

	// spinner


	// Waiting list.
	twoWaitingLists: any

	// course modules
	courseModules: ICourseModule[]
	totalCourseModules: number
	moduleIDList: string[]

	// batch modules
	batchModules: IBatchModule[]
	totalBatchModules: number
	batchModuleForm: FormGroup

	// module modal
	activeModuleTab: number

	//totalDaysSelected
	totalDaysSelected: number

	moduleTimings: any

	@ViewChild('batchModal') batchModal: any
	@ViewChild('talentModal') talentModal: any
	@ViewChild('waitingListModal') waitingListModal: any
	@ViewChild('deleteModal') deleteModal: any
	@ViewChild('module') moduleModal: any
	@ViewChild('drawer') drawer: any
	@ViewChild('dropdownButton') dropdownButton: any

	private readonly IGNORE_SEARCH_FIELD: string[] = ["limit", "offset", "facultyID", 'salesPersonID']

	constructor(
		private formBuilder: FormBuilder,
		public utilService: UtilityService,
		private generalService: GeneralService,
		private techService: TechnologyService,
		private talentService: TalentService,
		private batchService: BatchService,
		private courseModuleService: CourseModuleService,
		private fileOperationService: FileOperationService,
		private urlConstant: UrlConstant,
		private modalService: NgbModal,
		private spinnerService: SpinnerService,
		private localService: LocalService,
		private router: Router,
		private route: ActivatedRoute,
		private role: Role,
		private datePipe: DatePipe,
		private urlService: UrlService,
		private adminService: AdminService
	) {
		this.initializeVariables()
		this.createForm()
		this.getAllComponents()
	}

	initializeVariables(): void {
		this.isSalesPerson = (this.localService.getJsonValue("roleName") == this.role.SALES_PERSON ? true : false)
		this.isFaculty = (this.localService.getJsonValue("roleName") == this.role.FACULTY ? true : false)
		this.isTalent = (this.localService.getJsonValue("roleName") == this.role.TALENT ? true : false)
		this.isAdmin = (this.localService.getJsonValue("roleName") == this.role.ADMIN ? true : false)
		this.loginID = this.localService.getJsonValue("loginID")

		if (this.isAdmin || this.isSalesPerson) {
			this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_BATCH_MASTER)
		}

		if (this.isFaculty || this.isSalesPerson) {
			this.permission = this.utilService.getPermission(this.urlConstant.MY_BATCH)
		}

		this.salesPeople = []
		this.facultyList = []
		this.courseList = []
		this.selectedBatch = []
		this.addSelectedTalents = []
		this.activeBatchList = []
		this.courseModules = []
		this.batchModules = []
		this.moduleIDList = []
		this.daysList = []

		this.isViewClicked = false;
		this.toggleTalentSelect = false;
		this.isSearched = false
		this.showSearch = false
		this.isTalentAssignedOrDeleted = false
		this.disableButton = false
		this.doesEligibilityExists = false
		this.isBatchLoaded = true
		this.isTalentLoaded = true
		this.isSelectedTalentLoaded = true
		this.isBrochureUploadedToServer = false
		this.isFileUploading = false
		this.isAddClicked = false
		this.showTransferButton = false
		this.showTransferForm = false
		this.showCompanyReqField = false
		this.showBatchField = false
		this.isTechLoading = false
		this.eligibilityID = null


		this.limit = 5
		this.offset = 0
		this.talentLimit = 5
		this.talentOffset = 0
		this.totalTalent = 0
		this.selectedTalentLimit = 5
		this.selectedTalentOffset = 0
		this.totalBatch = 0
		this.totalBatchTalents = 0
		this.techLimit = 10
		this.techOffset = 0
		this.totalBatchModules = 0
		this.totalCourseModules = 0
		this.activeModuleTab = 1
		this.supervisorFacultyCount = 0
		this.isHeadFcaulty = false
		this.isViewAllBatches = false

		this.totalDaysSelected = 0

		this.spinnerService.loadingMessage = "Loading batches"
		this.displayedFileName = "Select file"

		this.searchFormValue = {}
		this.twoWaitingLists = {}

		this.selecetdrRequirementForTransfer = null
		this.selectedBatchForTransfer = null

	}

	getAllComponents(): void {
		this.getTechnology()
		this.getAcademicYear()
		this.getStudentRatingList()
		this.getAllSalesPeople()
		this.getBatchStatus()
		this.getBatchObjectives()
		this.getCourseList()
		this.getFacultyList()
		this.getRequirementList()
		this.searchOrGetBatches()
		this.getDaysList()
		this.getSupervisorCount()
		// this.getBatches();
	}

	createForm(): void {
		this.createSearchBatchForm()
		this.createBatchForm()
	}


	get ongoingOperations() {
		return this.spinnerService.ongoingOperations
	}

	ngOnInit() {
	}

	createSearchBatchForm(): void {
		this.searchBatchForm = this.formBuilder.group({
			// batchID: new FormControl(null),
			batchName: new FormControl(null),
			batchStatus: new FormControl(null),
			batchObjective: new FormControl(null),
			isActive: new FormControl(null),
			startDate: new FormControl(null),
			estimatedEndDate: new FormControl(null),
			courseID: new FormControl(null),
			salesPersonID: new FormControl(null),
			limit: new FormControl(this.limit),
			offset: new FormControl(this.offset),
			facultyID: new FormControl(null),
		})
		this.searchBatchForm.markAsDirty()
	}

	createBatchForm(): void {
		this.doesEligibilityExists = false

		this.batchForm = this.formBuilder.group({
			id: new FormControl(),
			batchName: new FormControl(null),
			code: new FormControl(null),
			course: new FormControl(null),
			startDate: new FormControl(null),
			// endDate: new FormControl(null),
			estimatedEndDate: new FormControl(null),
			totalStudents: new FormControl(null),
			totalIntake: new FormControl(null),
			batchStatus: new FormControl("Upcoming"),
			batchObjective: new FormControl(null),
			isActive: new FormControl(true),
			isB2B: new FormControl(null),
			brochure: new FormControl(null),
			meetLink: new FormControl(null),
			telegramLink: new FormControl(null),
			eligibility: this.formBuilder.group({
				id: [],
				technologies: [null],
				studentRating: [null],
				experience: [null],
				academicYear: [null]
			}),
			salesPerson: new FormControl(null),
			requirement: new FormControl(null),
			batchTimings: this.formBuilder.array([]),
			isApplyToAllSessions: new FormControl(false),
			// logo: new FormControl(null),
			// faculty: new FormControl(null),
			// sessions: new FormControl(Array()),
		})
		this.addValidatorsToBatchForm()
	}

	addValidatorsToBatchForm(): void {
		if (this.isAdmin) {
			this.addAdminValidators()

			return
		}
		if (this.isAddClicked) {
			this.addBatchTimingForm()
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
		// this.batchForm.get("requirement").setValidators([Validators.required])
		// this.batchForm.get("brochure").setValidators([Validators.required])

		this.utilService.updateValueAndValiditors(this.batchForm)
	}

	// addFacultyValidators(): void {
	// 	this.batchForm.get("startDate").setValidators([Validators.required])
	// 	this.batchForm.get("batchStatus").setValidators([Validators.required])
	// 	this.batchForm.get("faculty").setValidators([Validators.required])

	// 	for (let index = 0; index < this.facultyList.length; index++) {
	// 		if (this.isFaculty && this.loginID == this.facultyList[index].id) {
	// 			this.batchForm.get('faculty').setValue(this.facultyList[index])
	// 		}
	// 	}

	// 	this.utilService.updateValueAndValiditors(this.batchForm)
	// }

	get batchTiming() {
		return this.batchForm.get('batchTimings') as FormArray
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

	createBatchTimingForm(day: any): FormGroup {

		return this.formBuilder.group({
			id: new FormControl(),
			day: new FormControl(day, [Validators.required]),
			fromTime: new FormControl(null, [Validators.required]),
			toTime: new FormControl(null, [Validators.required])
		})

		// this.batchTiming.push(this.formBuilder.group({
		// 	id: new FormControl(),
		// 	day: new FormControl(null, [Validators.required]),
		// 	fromTime: new FormControl(
		// 		(this.batchTiming.at(0)?.get('fromTime')?.value != null) ? this.batchTiming.at(0)?.get('fromTime')?.value : null,
		// 		[Validators.required]),
		// 	toTime: new FormControl(
		// 		(this.batchTiming.at(0)?.get('toTime')?.value != null) ? this.batchTiming.at(0)?.get('toTime')?.value : null,
		// 		[Validators.required])
		// }))
	}

	// delete batch timing from batch form.
	deleteBatchTimingInForm(index: number): void {
		this.batchTiming.removeAt(index)
		this.batchForm.markAsDirty()
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
		this.batchTiming.insert(index + 1, this.createBatchTimingForm(day))

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

	// Delete batch timing from batch form.
	deleteBatchTimingInBatchForm(index: number): void {
		if (confirm("Are you sure you want to delete the batch timing?")) {
			this.batchTiming.removeAt(index)

			// If removed batch timing is first index, then make is apply to all sessions false and enable all other batch timings.
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
		this.batchForm.markAsDirty();

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

	studentRatingCheck() {
		const studentRatingControl = this.batchForm.get('eligibility.studentRating')
		if (this.doesEligibilityExists) {
			studentRatingControl.setValidators([Validators.required])
		} else {
			studentRatingControl.setValue(null)
			studentRatingControl.setValidators(null)
			studentRatingControl.markAsUntouched()
		}
		studentRatingControl.updateValueAndValidity()
	}

	experienceCheck() {
		const experienceControl = this.batchForm.get('eligibility.experience')
		if (this.doesEligibilityExists) {
			experienceControl.setValidators([Validators.required])
		} else {
			experienceControl.setValue(null)
			experienceControl.setValidators(null)
			experienceControl.markAsUntouched()
		}
		experienceControl.updateValueAndValidity()
	}

	academicYearCheck() {
		const academicYearControl = this.batchForm.get('eligibility.academicYear')
		if (this.doesEligibilityExists) {
			academicYearControl.setValidators([Validators.required])
		} else {
			academicYearControl.setValue(null)
			academicYearControl.setValidators(null)
			academicYearControl.markAsUntouched()
		}
		academicYearControl.updateValueAndValidity()
	}

	createSearchTalentForm() {
		this.searchedTalentForm = this.formBuilder.group({
			firstName: new FormControl(null),
			lastName: new FormControl(null),
			contact: new FormControl(null),
			email: new FormControl(null),
		});
	}

	updateBatchForm(): void {

		if (this.selectedBatch.eligibility != null) {
			this.doesEligibilityExists = true
		} else {
			this.selectedBatch.eligibility = {}
			this.doesEligibilityExists = false
		}

		// Sort batch timings by order of days.
		this.selectedBatch.batchTimings.sort(this.sortBatchTimings)
		console.log("selectedBatchbatchTimings", this.selectedBatch.batchTimings)

		// Add batch timings to batch timing form.
		for (let i = 0; i < this.selectedBatch.batchTimings.length; i++) {
			this.addBatchTimingForm()
		}
		// console.log(this.selectedBatch);
		this.batchForm.patchValue(this.selectedBatch)
		this.getDaysList()

		// make eligibility null if eligibility does not exist
		// if (Object.keys(this.selectedBatch.eligibility).length == 0) {
		// 	this.selectedBatch.eligibility = null
		// }

		if (!this.doesEligibilityExists) {
			this.selectedBatch.eligibility = null
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

	onBatchAddClick(): void {
		this.isAddClicked = true
		this.isViewClicked = false
		this.updateModalVariable("Add a new batch", "Add Batch", this.createBatchForm, this.addBatch)
		this.formHandler();
		this.openModal(this.batchModal, "xl")
	}

	onBatchViewClick(event: any, batch: any): void {
		event.stopPropagation();
		this.isViewClicked = true
		this.isAddClicked = false
		// console.log(batch)
		this.updateModalVariable("Batch details", "Batch", this.createBatchForm, this.addBatch)
		this.formHandler();

		batch.startDate = this.datePipe.transform(batch.startDate, 'yyyy-MM-dd')
		batch.estimatedEndDate = this.datePipe.transform(batch.estimatedEndDate, 'yyyy-MM-dd')

		this.selectedBatch = batch
		this.updateBatchForm()
		this.batchForm.disable()
		this.openModal(this.batchModal, "xl")
	}

	onBatchUpdateClick(event: any, batch: any): void {
		event.stopPropagation();

		this.isViewClicked = false
		this.isAddClicked = false

		this.updateModalVariable("Update Batch details", "Update Batch", this.createBatchForm, this.updateBatch)
		this.formHandler();

		batch.startDate = this.datePipe.transform(batch.startDate, 'yyyy-MM-dd')
		batch.estimatedEndDate = this.datePipe.transform(batch.estimatedEndDate, 'yyyy-MM-dd')

		this.selectedBatch = batch
		this.updateBatchForm()

		if (this.selectedBatch.brochure) {
			this.displayedFileName = `<a href=${this.selectedBatch.brochure} target="_blank">Batch Brochure</a>`
		}
		if (this.selectedBatch.logo) {
			this.logoDisplayedFileName = `<a href=${this.selectedBatch.logo} target="_blank">Batch Logo</a>`
		}

		// this.modalHeadingLabel = "Update Batch Details"
		// this.modalButtonLabel = "Update Batch"
		// this.modalHandler = this.updateBatch
		this.openModal(this.batchModal, "xl")
	}

	// Update Modal Variable 
	updateModalVariable(modalHeadingLabel: string, modalButtonLabel: string, formHandler: any, modalHandler: any): void {
		this.modalHeadingLabel = modalHeadingLabel;
		this.modalButtonLabel = modalButtonLabel;
		this.formHandler = formHandler;
		this.modalHandler = modalHandler;
	}

	// Compare two Object.
	compareFn(ob1: any, ob2: any): any {

		if (ob1 == null && ob2 == null) {
			return true
		}
		return ob1 && ob2 ? ob1.id == ob2.id : ob1 == ob2;
	}

	assignSelectedBatch(event: any, batch: any): void {
		event.stopPropagation();
		this.selectedBatch = batch;
		this.waitingListMessage = "This will also mean that all the related waiting list entries in talent and enquiry will be inactive."
		this.entity = "batch"
		this.toggleTalentSelect = true;
		// this.modalHandler = this.deleteBatch
		this.openModal(this.deleteModal, 'md').result.then(() => {
			this.deleteBatch(batch)
		}, (err) => {
			console.error(err);
			return
		})
	}


	deleteEligibilityIfEmpty(batch: any): void {
		if (Object.keys(batch.eligibility).length === 0 && batch.eligibility.constructor === Object) {
			delete batch.eligibility
		}
	}

	setPaginationString() {
		this.paginationString = ''
		let limit = this.searchBatchForm.get('limit').value
		let offset = this.searchBatchForm.get('offset').value

		let start: number = limit * offset + 1
		let end: number = +limit + limit * offset
		// let start: number = this.limit * this.offset + 1
		// let end: number = +this.limit + this.limit * this.offset
		if (this.totalBatch < end) {
			end = this.totalBatch
		}
		if (this.totalBatch == 0) {
			this.paginationString = ''
			return
		}
		this.paginationString = `${start} - ${end}`
	}

	setTalentsPaginationString() {
		this.talentPaginationString = ''
		let start: number = this.talentLimit * this.talentOffset + 1
		let end: number = +this.talentLimit + this.talentLimit * this.talentOffset
		if (this.totalTalent < end) {
			end = this.totalTalent
		}
		if (this.totalTalent == 0) {
			this.talentPaginationString = ''
			return
		}
		this.talentPaginationString = `${start} - ${end} of ${this.totalTalent}`
	}

	setSelectedTalentsPaginationString() {
		this.selectedTalentPaginationString = ''
		let start: number = this.selectedTalentLimit * this.selectedTalentOffset + 1
		let end: number = +this.selectedTalentLimit + this.selectedTalentLimit * this.selectedTalentOffset
		if (this.totalSelectedTalent < end) {
			end = this.totalSelectedTalent
		}
		if (this.totalSelectedTalent == 0) {
			this.selectedTalentPaginationString = ''
			return
		}
		this.selectedTalentPaginationString = `${start} - ${end} of ${this.totalSelectedTalent}`
	}

	// Handles pagination
	changePage(pageNumber: number): void {
		this.searchBatchForm.get("offset").setValue(pageNumber - 1)
		this.searchBatch()
	}

	// Handles pagination
	changeTalentPage($event: any): void {
		// $event will be the page number & offset will be 1 less than it.
		this.talentOffset = $event - 1
		this.currentTalentPage = $event
		this.getAllTalents()
	}

	// Handles pagination
	changeSelectedTalentPage($event: any): void {
		// $event will be the page number & offset will be 1 less than it.
		this.selectedTalentOffset = $event - 1
		this.currentSelectedTalentPage = $event
		this.getSelectedTalents(this.batchID)
	}


	toggleButtons() {
		this.disableButton = false
		this.doesEligibilityExists = false
	}

	checkTalentAdded(talentID) {
		return this.addSelectedTalents.includes(talentID)
	}

	assignTalentToBatch(batchID: string) {
		this.addSelectedTalents = []
		this.isTalentAssignedOrDeleted = false
		this.createSearchTalentForm()
		this.openModal(this.talentModal, "xl")
		this.getSelectedTalents(batchID)
		this.getAllTalents()
	}

	getTalentStatus(talentID: string): boolean {
		for (let index = 0; index < this.selectedTalent.length; index++) {
			if (this.selectedTalent[index].id == talentID) {
				return true
			}
		}
		return false
	}

	// takes a list called selectedTalent and adds all the checked talents to list, also does not contain duplicate values
	addTalentsToList(event, id: string) {
		if (event.target.checked) {
			//   console.log(id);
			this.addSelectedTalents.push(id)
		} else {
			if (this.addSelectedTalents.includes(id)) {
				let index = this.addSelectedTalents.indexOf(id)
				this.addSelectedTalents.splice(index, 1)
			}
		}
		// console.log(this.addSelectedTalents);
	}

	assignTalentToDelete(talentID: string) {
		// this.talentID = talentID
		this.waitingListMessage = ""
		this.entity = "talent from the batch"
		this.openModal(this.deleteModal, 'md').result.then(() => {
			this.deleteTalentFromBatch(talentID)
		}, (err) => {
			console.error(err);
			return
		})
	}

	calculateBatchTiming(batchTimings: IBatchTiming[]): string {
		let batchSchedule: string = ""
		let order: number

		if (batchTimings.length == 0) {
			return "-"
		}


		batchSchedule += this.getBatchScheduleDays(batchTimings) + "\n"
		batchSchedule += " " + this.addTimingToBatchSchedule(batchTimings)
		return batchSchedule
	}

	getBatchScheduleDays(batchTimings: IBatchTiming[]): string {
		let batchSchedule: string = ""
		let endDay: string
		let order: number
		let isWeekly: boolean = true

		batchSchedule += batchTimings[0]?.day?.day
		order = batchTimings[0]?.day?.order

		if (batchTimings.length > 1) {
			batchSchedule = batchSchedule.substring(0, 3)
			for (let index = 1; index < batchTimings.length; index++) {
				endDay = batchTimings[index].day.day.substring(0, 3)
				if (order + 1 == batchTimings[index].day.order) {
					order = batchTimings[index].day.order
					isWeekly = (!isWeekly ? isWeekly : true)
					continue
				} else {
					batchSchedule += ", " + endDay
					isWeekly = false
				}
			}
			if (isWeekly) {
				batchSchedule += " to " + endDay
			}
		}
		return batchSchedule
	}

	addTimingToBatchSchedule(batchTimings: IBatchTiming[]): string {
		return this.utilService.parseTime(batchTimings[0]?.fromTime) +
			" - " + this.utilService.parseTime(batchTimings[0]?.toTime)
	}

	resetSearchAndGetAll(): void {
		// this.searchBatchForm.reset()
		this.resetSearchForm()
		// this.searchFormValue = {}
		this.isSearched = false
		this.showSearch = false
		// this.router.navigate([this.urlConstant.BATCH_MASTER])
		this.changePage(1)
	}

	resetSearchForm(): void {
		this.limit = this.searchBatchForm.get("limit").value
		this.offset = this.searchBatchForm.get("offset").value
		this.searchBatchForm.reset({
			limit: this.limit,
			offset: this.offset,
			facultyID: this.isFaculty ? this.loginID : null,
			salesPersonID: this.isSalesPerson ? this.loginID : null,
		})
	}

	redirectToBatchSession(batch: IBatch): void {
		this.router.navigate(['session'], {
			relativeTo: this.route,
			queryParams: {
				"batchID": batch.id,
				"courseID": batch.course.id,
				"batchName": batch.batchName
			}
		}).catch(err => {
			console.error(err)

		});
	}


	redirectToBatchDetails(batch: IBatch): void {
		this.urlService.setUrlFrom(this.constructor.name, this.router.url)
		if (batch.course.id != null) {
			this.router.navigate(["session/details"], {
				relativeTo: this.route,
				queryParams: {
					"batchID": batch.id,
					"courseID": batch.course.id,
					"batchName": batch.batchName
				}
			}).catch(err => {
				console.error(err)

			});
		} else {
			alert("Course ID Doesn't Exit")
		}
	}


	redirectToBatchFeedback(batch: IBatch): void {
		this.router.navigate(["feedback"], {
			relativeTo: this.route,
			queryParams: {
				"batchID": batch.id,
				"batchName": batch.batchName
			}
		}).catch(err => {
			console.error(err)

		});
	}

	redirectToTalentMaster(batchID: string): void {
		this.router.navigate(['/talent/master'], {
			queryParams: {
				"batchID": batchID
			}
		}).catch(err => {
			console.error(err)

		});
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

	openModal(modalContent: any, modalSize = "lg"): NgbModalRef {

		this.resetBrochureUploadFields()
		this.resetLogoUploadFields()

		let options: NgbModalOptions = {
			ariaLabelledBy: 'modal-basic-title', keyboard: false,
			backdrop: 'static', size: modalSize
		}
		this.modalRef = this.modalService.open(modalContent, options)
		return this.modalRef

		/*this.modalRef.result.subscribe((result) => {
		}, (reason) => {

		});*/
	}

	// Used to dismiss modal.
	dismissFormModal(modal: NgbModalRef) {
		if (this.isFileUploading || this.isLogoFileUploading) {
			alert("Please wait till file is being uploaded")
			return
		}
		if (this.isBrochureUploadedToServer || this.isLogoUploadedToServer) {
			if (!confirm("Uploaded brochure will be deleted.\nAre you sure you want to close?")) {
				return
			}
		}
		modal.dismiss()
		this.resetBrochureUploadFields()
		this.resetLogoUploadFields()
	}

	//validate add/update form
	validate(): void {
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
			if (confirm("Making batch inactive or its status as 'Finished' will also mean that all the related waiting list entries in talent and enquiry will be inactive. Do you want to go ahead??")) {
				this.modalHandler();
				return
			}
			return
		}

		this.modalHandler();
	}

	resetSelectedTalent() {
		this.addSelectedTalents = []
		if (this.isTalentAssignedOrDeleted) {
			this.getBatches()
		}
	}

	resetBatchForm() {
		this.doesEligibilityExists = false
		this.batchForm.reset()
	}

	searchAndCloseDrawer(): void {
		this.drawer.toggle()

		this.searchBatchForm.get("limit").setValue(5)
		this.searchBatchForm.get("offset").setValue(0)

		this.searchBatch()
	}

	searchBatch(): void {
		if (this.isFaculty) {
			this.searchBatchForm.get("facultyID").setValue(this.loginID)
		}

		if (this.isSalesPerson) {
			this.searchBatchForm.get("salesPersonID").setValue(this.loginID)
		}

		this.searchFormValue = { ...this.searchBatchForm?.value }
		this.router.navigate([], {
			relativeTo: this.route,
			queryParams: this.searchFormValue,
		})

		let flag: boolean = true

		for (let field in this.searchFormValue) {
			if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
				delete this.searchFormValue[field];
			} else {
				flag = false
				if (this.isFaculty && !this.IGNORE_SEARCH_FIELD.includes(field)) {
					this.isSearched = true
				}
				if (this.isSalesPerson && !this.IGNORE_SEARCH_FIELD.includes(field)) {
					this.isSearched = true
				}
				if (this.isAdmin && !this.IGNORE_SEARCH_FIELD.includes(field)) {
					this.isSearched = true
				}
			}
		}

		// No API call on empty search.
		if (flag) {
			return
		}
		this.getBatches()
	}

	searchOrGetBatches(): void {
		let queryParams = this.route.snapshot.queryParams
		if (!this.utilService.isObjectEmpty(queryParams)) {
			this.searchBatchForm.patchValue(queryParams)
		}
		this.searchBatch()
	}

	// =============================================================CRUD=============================================================

	// Get All Batches
	getBatches(): void {
		this.searchFormValue.isViewAllBatches = this.isViewAllBatches ? 1 : 0
		if (this.isTalent) {
			this.getTalentBatches()
			return
		}

		this.spinnerService.loadingMessage = "Loading batches";

		this.batches = []
		this.totalBatch = 0
		this.totalBatchTalents = 0
		this.disableButton = true
		this.isBatchLoaded = true

		// console.log(this.searchFormValue);

		this.batchService.getBatches(this.searchFormValue).subscribe((respond: any) => {
			this.batches = respond.body;
			this.totalBatch = parseInt(respond.headers.get('X-Total-Count'))
			this.totalBatchTalents = parseInt(respond.headers.get("Total-Batch-Talents"))
			// console.log(this.batches)
			if (this.totalBatch == 0) {
				this.isBatchLoaded = false
			}
		}, (error: any) => {
			this.totalBatch = 0
			this.totalBatchTalents = 0
			this.isBatchLoaded = false
			if (error.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			console.error(this.utilService.getErrorString(error))
			alert(this.utilService.getErrorString(error))
		}).add(() => {

			this.setPaginationString()
			this.disableButton = false
		})
	}

	formatBatchTimings(): void {
		for (let index = 0; index < this.batches.length; index++) {
			for (let j = 0; j < this.batches[index].batchTimings.length; j++) {
				// Remove seconds form batch timings.
				this.batches[index].batchTimings[j].fromTime = this.batches[index].batchTimings[j].fromTime.replace(/:[^:]*$/, '')
				this.batches[index].batchTimings[j].toTime = this.batches[index].batchTimings[j].toTime.replace(/:[^:]*$/, '')
			}
		}
	}

	// Get All Batches
	getSearchBatches(): void {

		this.spinnerService.loadingMessage = "Loading batches";

		this.batches = []
		this.totalBatch = 0
		this.totalBatchTalents = 0
		this.isBatchLoaded = true
		this.disableButton = true

		// console.log(this.searchFormValue);
		this.batchService.getBatches(this.searchFormValue).subscribe((respond: any) => {
			this.batches = respond.body;
			this.totalBatch = parseInt(respond.headers.get('X-Total-Count'))
			this.totalBatchTalents = parseInt(respond.headers.get("Total-Batch-Talents"))
			if (this.totalBatch == 0) {
				this.isBatchLoaded = false
			}
		}, (error) => {
			this.totalBatch = 0
			this.totalBatchTalents = 0
			this.isBatchLoaded = false
			if (error.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			console.error(this.utilService.getErrorString(error))
		}).add(() => {

			this.setPaginationString()
			this.disableButton = false
		})
	}

	getTalentBatches(): void {
		this.disableButton = true
		this.totalBatch = 0
		this.isBatchLoaded = true
		this.batches = []
		this.spinnerService.loadingMessage = "Loading batches";

		// console.log(this.localService.getJsonValue('loginID'));
		this.talentService.getBatchListOfOneTalent(this.loginID).
			subscribe((response: any) => {
				// console.log(response);
				this.batches = response
				this.totalBatch = this.batches.length
				if (this.totalBatch == 0) {
					this.isBatchLoaded = false
				}
				// this.setPaginationString()
				// 
				this.disableButton = false
			}, (error: any) => {
				this.totalBatch = 0
				// this.setPaginationString()

				this.isBatchLoaded = false
				this.disableButton = false
				if (error.statusText.includes('Unknown')) {
					alert("No connection to server. Check internet.")
					return
				}
				console.error(this.utilService.getErrorString(error))
			}).add(() => {

				this.setPaginationString()
			})
	}

	// Extract id from objects and give it to ompany requirement.
	setFormFields(batch: any): void {
		if (this.batchForm.get("startDate")?.value) {
			batch.startDate = new Date(this.batchForm.get("startDate")?.value).toISOString()
		}
	}

	addBatch(): void {
		this.disableButton = true
		let batch = this.batchForm.value
		this.utilService.deleteNullValueIDFromObject(batch)
		this.utilService.deleteNullValuePropertyFromObject(batch.eligibility)
		this.setFormFields(batch)
		if (batch.eligibility) {
			this.deleteEligibilityIfEmpty(batch)
		}
		this.spinnerService.loadingMessage = "Adding batch"

		batch.startDate = this.datePipe.transform(batch.startDate, "yyyy-MM-dd")
		// console.log(batch);
		this.batchService.addBatch(batch).subscribe((respond: any) => {
			// console.log(respond)
			this.modalRef.close();
			this.getBatches()
			alert("New Batch Added Successfully")
			this.batchForm.reset()
		}, (err) => {
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(this.utilService.getErrorString(err));
		}).add(() => {

			this.disableButton = false
		})
	}

	// Update Batch
	updateBatch(): void {
		this.batchForm.enable()
		this.disableButton = true
		let batch = this.batchForm.value
		this.utilService.deleteNullValueIDFromObject(batch)
		this.utilService.deleteNullValuePropertyFromObject(batch.eligibility)
		// this.setFormFields(batch)
		if (batch.eligibility) {
			this.deleteEligibilityIfEmpty(batch);
		}
		batch.startDate = this.datePipe.transform(batch.startDate, "yyyy-MM-dd")
		// console.log(batch);
		this.spinnerService.loadingMessage = "Updating batch";

		this.batchService.updateBatch(batch).subscribe((respond: any) => {
			this.eligibilityID = null
			// console.log(this.searched);
			this.modalRef.close()
			this.getBatches()
			alert("Batch Updated Successfully");
			this.batchForm.reset()
		}, (err) => {
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(this.utilService.getErrorString(err));
			console.error(err)
		}).add(() => {

			this.disableButton = false
		})
	}

	// Delete Batch
	deleteBatch(batch: any): void {
		this.disableButton = true
		this.spinnerService.loadingMessage = "Deleting batch";

		this.batchService.deleteBatch(batch.id).subscribe((respond: any) => {
			// 
			alert("Batch Deleted Successfully");
			this.modalRef.close();
			// console.log(this.offset);
			if (this.batches.length == 1) {
				this.changePage(this.offset)
			}
			this.getBatches();
		}, (err) => {

			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(this.utilService.getErrorString(err))
		}).add(() => {

			this.disableButton = false
		})
	}

	// returns all the selected talents for the batch
	getSelectedTalents(batchID: string) {
		// this.apiBtn = true
		this.disableButton = true
		this.batchID = batchID
		this.selectedTalent = []
		this.totalSelectedTalent = 0
		this.isSelectedTalentLoaded = true
		this.spinnerService.loadingMessage = "Getting Talents"

		this.batchService.getTalentsForBatch(batchID, this.selectedTalentLimit, this.selectedTalentOffset).
			subscribe((response: any) => {
				this.selectedTalent = response.body
				this.totalSelectedTalent = response.headers.get('X-Total-Count')
				this.disableButton = false
				if (this.totalSelectedTalent == 0) {
					this.isSelectedTalentLoaded = false
				}
				// this.setSelectedTalentsPaginationString()
				// 
				//     console.log(this.selectedTalent)
			}, (err: any) => {
				this.totalSelectedTalent = 0
				this.disableButton = false
				this.isSelectedTalentLoaded = false
				// this.setSelectedTalentsPaginationString()
				// 
				if (err.statusText.includes('Unknown')) {
					alert("No connection to server. Check internet.")
					return
				}
				// console.log("Error in talent for batch")
				console.error(err)
			}).add(() => {

				this.setSelectedTalentsPaginationString()
			})

	}

	// returns all the talents
	getAllTalents(): void {
		let roleNameAndLogin: any = {
			roleName: this.localService.getJsonValue("roleName"),
			loginID: this.loginID
		};

		if (this.searchedTalent) {
			this.getAllSearchedTalents()
			return
		}
		this.disableButton = true
		this.totalTalent = 0
		this.isTalentLoaded = true
		this.talents = []
		this.spinnerService.loadingMessage = "Getting Talents"

		this.talentService.getTalents(this.talentLimit, this.talentOffset, roleNameAndLogin).subscribe((response: any) => {

			this.talents = response.body
			this.totalTalent = response.headers.get('X-Total-Count')
			if (this.totalTalent == 0) {
				this.isTalentLoaded = false
			}
			// this.setTalentsPaginationString()
			// console.log(this.talents)
			// this.checkTalentAdded(this.talents)
			this.disableButton = false
		}, err => {
			this.totalTalent = 0
			this.isTalentLoaded = false
			// this.setTalentsPaginationString()

			this.disableButton = false
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err.error.error)
			console.error("Error" + err);
		}).add(() => {

			this.setTalentsPaginationString()
		})
	}

	getAllSearchedTalents() {
		let roleNameAndLogin: any = {
			roleName: this.localService.getJsonValue("roleName"),
			loginID: this.loginID
		};
		if (!this.searchedTalent) {
			this.searchedTalent = true
			this.offset = 0
			this.currentTalentPage = 1
		}
		this.disableButton = true
		this.talents = []
		this.totalTalent = 0
		this.isTalentLoaded = true
		let data = this.searchedTalentForm.value
		this.utilService.deleteNullValuePropertyFromObject(data)
		// console.log(data)
		this.spinnerService.loadingMessage = "Getting Searched Talents";

		this.talentService.getAllSearchedTalents(data, this.talentLimit, this.talentOffset, roleNameAndLogin)
			.subscribe(response => {
				this.disableButton = false
				this.talents = response.body
				this.totalTalent = parseInt(response.headers.get('X-Total-Count'))
				if (this.totalTalent == 0) {
					this.isTalentLoaded = false
				}
				// 
			}, (error) => {
				this.totalTalent = 0
				this.disableButton = false
				this.isTalentLoaded = false
				// 
				if (error.statusText.includes('Unknown')) {
					alert("No connection to server. Check internet.")
					return
				}
				console.error(error)
			}).add(() => {

				this.setTalentsPaginationString()
			})
	}

	// ********************************Add/Delete talent to batch********************************

	// will add multiple talents to batch
	addTalentsToBatch(): void {
		//create a object for mappedTalents
		let selectedTalentObject: IMappedTalent[] = []
		this.addSelectedTalents.forEach(element => {
			selectedTalentObject.push({
				"talentID": element
			})
		});
		this.disableButton = true
		this.spinnerService.loadingMessage = "Adding talent to batch"

		this.batchService.addTalentsToBatch(selectedTalentObject, this.batchID).subscribe((response: any) => {
			// console.log(response)
			this.disableButton = false
			this.addSelectedTalents = []
			this.isTalentAssignedOrDeleted = true
			alert("Talents successfully added to batch")

			this.getAllTalents()
			this.getSelectedTalents(this.batchID)
		}, err => {
			this.disableButton = false

			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			// this.addSelectedTalents = []
			console.error(err.error);
			alert(err.error.error)
		}).add(() => {

		})

	}

	deleteTalentFromBatch(talentID: string): void {
		this.disableButton = true
		this.spinnerService.loadingMessage = "Deleting talent from batch"

		this.batchService.deleteTalentFromBatch(talentID, this.batchID).subscribe(() => {
			this.modalRef.close()
			this.getAllTalents()
			this.disableButton = false
			this.isTalentAssignedOrDeleted = true
			this.getSelectedTalents(this.batchID)
			alert("Talent successfully deleted from the batch")
			// 
		}, err => {
			// 
			this.disableButton = false
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			console.error(err.error);
			alert(err.error)
		}).add(() => {

		})
	}

	//On uplaoding brochure
	onResourceSelect(event: any) {
		this.docStatus = ""
		let files = event.target.files
		if (files && files.length) {
			let file = files[0]
			let err = this.fileOperationService.isDocumentFileValid(file)
			if (err != null) {
				this.docStatus = `<p><span>&#10060;</span> ${err}</p>`
				return
			}
			// console.log(file)
			// Upload brochure if it is present.]
			this.isFileUploading = true
			this.fileOperationService.uploadBrochure(file,
				this.fileOperationService.BATCH_FOLDER + this.fileOperationService.BROCHURE_FOLDER)
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
			let err = this.fileOperationService.isImageFileValid(file)
			if (err != null) {
				this.logoDocStatus = `<p><span>&#10060;</span> ${err}</p>`
				return
			}
			// console.log(file)
			// Upload brochure if it is present.]
			this.isLogoFileUploading = true
			this.fileOperationService.uploadLogo(file,
				this.fileOperationService.BATCH_FOLDER + this.fileOperationService.LOGO_FOLDER).subscribe((data: any) => {
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


	// ********************************Wating List related functions********************************

	// On clicking show applicants button in batch list.
	onShowApplicantsButtonClick(event: any, batch: any): void {
		event.stopPropagation();
		this.showTransferButton = false
		this.showTransferForm = false
		this.showCompanyReqField = false
		this.showBatchField = false
		this.selecetdrRequirementForTransfer = null
		this.selectedBatchForTransfer = null
		let modalSize: string
		if (!batch.isActive || batch.batchStatus == "Finished") {
			this.showTransferButton = true
			modalSize = 'md'
		} else {
			modalSize = 'sm'
		}
		this.getWaitingListByBatchID(batch)
		this.openModal(this.waitingListModal, modalSize)
		this.selectedBatch = batch
	}

	// Get waiting list by batch.
	getWaitingListByBatchID(batch: any): void {
		this.spinnerService.loadingMessage = "Getting Applicants"

		let queryParams: any = {
			batchID: batch.id
		}
		if (batch.isActive && batch.status != "Finished") {
			queryParams.isActive = "1"
		}
		else {
			queryParams.isActive = "0"
		}
		this.talentService.getTwoWaitingLists(queryParams).subscribe((response) => {
			this.twoWaitingLists = response

			if ((!this.selectedBatch.isActive || batch.batchStatus == "Finished") && ((this.twoWaitingLists.talentWaitingList?.length != 0) ||
				(this.twoWaitingLists.talentWaitingList?.length != 0))) {
				this.showTransferButton = true
				return
			}
			this.showTransferButton = false
		}, (err) => {
			console.error(err)

		})
	}

	// On clicking transfer button of applicants modal.
	onTrasnferButtonClick(): void {
		this.showTransferButton = false
		this.showTransferForm = true
	}

	// On changing value of waiting for control in waiting list form.
	onWaitingForChange(waitingFor: string): void {
		if (waitingFor == "Requirement") {
			this.showCompanyReqField = true
			this.showBatchField = false
			this.selecetdrRequirementForTransfer = null
			this.selectedBatchForTransfer = null
			this.getActiveRequirementList()
		}
		if (waitingFor == "Batch") {
			this.showCompanyReqField = false
			this.showBatchField = true
			this.selecetdrRequirementForTransfer = null
			this.selectedBatchForTransfer = null
			this.getActiveBatchList()
		}
	}

	// On clicking sunmit button of transfer form in applicants modal.
	onSubmitTransferButtonClick(): void {

		// If waiting for is not selected.
		if (!this.showCompanyReqField && !this.showBatchField) {
			alert("Please select waiting for field")
			return
		}

		// If compnay is selected.
		if (this.showCompanyReqField) {
			if (!this.selecetdrRequirementForTransfer) {
				alert("Please select a requirement ID")
				return
			}
			let waitingLists: any[] = []
			for (let i = 0; i < this.twoWaitingLists.talentWaitingList.length; i++) {
				waitingLists.push(this.twoWaitingLists.talentWaitingList[i])
			}
			for (let i = 0; i < this.twoWaitingLists.enquiryWaitingList.length; i++) {
				waitingLists.push(this.twoWaitingLists.enquiryWaitingList[i])
			}
			let updateWaitingList: any = {
				companyBranchID: this.selecetdrRequirementForTransfer.companyBranchID,
				requirementID: this.selecetdrRequirementForTransfer.id,
				batchID: null,
				courseID: null,
				waitingLists: waitingLists
			}
			this.spinnerService.loadingMessage = "Transfering Applicants"

			this.talentService.transferWaitingList(updateWaitingList).subscribe((response) => {

				alert("Some talents' or enquiries' waiting list entries may not be transfered since they already have the batch assigned to them")
				this.getWaitingListByBatchID(this.selectedBatch)
				this.getBatches()
				this.selecetdrRequirementForTransfer = null
				this.selectedBatchForTransfer = null
				this.showTransferForm = false
				return
			}, (err) => {
				console.error(err)

			})
		}

		// If batch is selected.
		if (this.showBatchField) {
			if (!this.selectedBatchForTransfer) {
				alert("Please select a batch ID")
				return
			}
			let waitingLists: string[] = []
			for (let i = 0; i < this.twoWaitingLists.talentWaitingList.length; i++) {
				waitingLists.push(this.twoWaitingLists.talentWaitingList[i])
			}
			for (let i = 0; i < this.twoWaitingLists.enquiryWaitingList.length; i++) {
				waitingLists.push(this.twoWaitingLists.enquiryWaitingList[i])
			}
			let updateWaitingList: any = {
				companyBranchID: null,
				requirementID: null,
				batchID: this.selectedBatchForTransfer.id,
				courseID: this.selectedBatchForTransfer.courseID,
				waitingLists: waitingLists
			}
			this.spinnerService.loadingMessage = "Transfering Applicants"

			this.talentService.transferWaitingList(updateWaitingList).subscribe((response) => {

				alert("Some talents' or enquiries' waiting list entries may not be transfered since they already have the batch assigned to them")
				this.getWaitingListByBatchID(this.selectedBatch)
				this.getBatches()
				this.selecetdrRequirementForTransfer = null
				this.selectedBatchForTransfer = null
				this.showTransferForm = false
				return
			}, (err) => {
				console.error(err)

			})
		}
	}

	// On clicking close transfer form button in applicants modal.
	onCloseTrasferFormButtonClick(): void {
		this.selecetdrRequirementForTransfer = null
		this.selectedBatchForTransfer = null
		this.showTransferButton = true
		this.showTransferForm = false
		this.showCompanyReqField = false
		this.showBatchField = false
	}

	// Redirect to talents page filtered by waiting lists' batch id.
	redirectToTalentsForWaitingList(batchID: string): void {
		this.router.navigate(['/talent/master'], {
			queryParams: {
				"waitingListBatchID": batchID
			}
		}).catch(err => {
			console.error(err)

		});
	}

	// Redirect to enquiries page filtered by waiting lists' batch id.
	redirectToEnquiriesForWaitingList(batchID: string): void {
		this.router.navigate(['/talent/enquiry'], {
			queryParams: {
				"waitingListBatchID": batchID
			}
		}).catch(err => {
			console.error(err)

		});
	}

	// ====================================================== MODULES ======================================================

	createBatchModuleForm(): void {
		this.batchModuleForm = this.formBuilder.group({
			batchModules: new FormArray([])
		})
	}

	get batchModuleControlArray(): FormArray {
		return this.batchModuleForm.get("batchModules") as FormArray
	}

	addBatchModulesToForm(): void {
		this.batchModuleControlArray.push(this.formBuilder.group({
			id: new FormControl(null),
			batchID: new FormControl(this.selectedBatch.id),
			moduleID: new FormControl(null, [Validators.required]),
			facultyID: new FormControl(null, [Validators.required]),
			order: new FormControl(null, [Validators.required, Validators.min(1)]),
			startDate: new FormControl(null),
			completedDate: new FormControl(null),
			isCompleted: new FormControl(false),
			isMarked: new FormControl(false),
			isModuleTimingMarked: new FormControl(false),
			moduleTimings: new FormArray([]),
			isApplyToAllSessions: new FormControl(false),

		}))
	}

	getModuleTimingControlArray(index: number): FormArray {
		return this.batchModuleControlArray.at(index).get("moduleTimings") as FormArray
	}

	addModuleTimingToForm(index: number): void {
		this.getModuleTimingControlArray(index).push(this.formBuilder.group({
			id: new FormControl(null),
			dayID: new FormControl(null),//required
			day: new FormControl(null),
			moduleID: new FormControl(null),
			facultyID: new FormControl(null),
			isMarked: new FormControl(null),
			fromTime: new FormControl(null),//required
			toTime: new FormControl(null),//required
		}))
	}

	// On checking apply it to all sessions.
	onApplyToAllModulesClick(index: number): void {
		this.batchModuleControlArray.at(index).get('isApplyToAllSessions')
			.setValue(!this.batchModuleControlArray.at(index).get('isApplyToAllSessions').value)

		// If checkbox is unselected then enable all the batch timings.
		if (this.batchModuleControlArray.at(index).get('isApplyToAllSessions').value == false) {
			// 	for (let i = 1; i < this.getModuleTimingControlArray(index).length; i++) {
			// 		this.getModuleTimingControlArray(index).at(i).get("toTime").enable()
			// 		this.getModuleTimingControlArray(index).at(i).get("fromTime").enable()
			// 	}
			return
		}

		var element = -1
		for (let k = 0; k < this.getModuleTimingControlArray(index).length; k++) {
			if (this.getModuleTimingControlArray(index).at(k).get("isMarked").value && element == -1) {
				element = k;
			}
		}
		// console.log(element);
		// console.log(this.getModuleTimingControlArray(index).at(element).get("fromTime").value);
		// console.log(this.getModuleTimingControlArray(index).at(element).get("toTime").value);

		// If first session day timings are empty then give alert.
		if (this.getModuleTimingControlArray(index).at(element).get("fromTime").value == null
			|| this.getModuleTimingControlArray(index).at(element).get("toTime").value == null) {
			this.batchModuleControlArray.at(index).get('isApplyToAllSessions').setValue(false)

			alert("Please fill in the from time and to time of the first session day")

			return
		}


		// If start time and end time are same.
		if (this.getModuleTimingControlArray(index).at(element).get("fromTime").value == this.getModuleTimingControlArray(index).at(element).get("toTime").value) {
			this.batchModuleControlArray.at(index).get('isApplyToAllSessions').setValue(false)
			alert("Start time and end time cannot be same")
			return
		}

		// If there are more than one batch timings then set all the other batch timings as the value for the first 
		// batch timing.
		// let moduleTimingsControl = (this.batchModuleControlArray.at(index).get("moduleTimings") as FormArray).controls
		if (this.getModuleTimingControlArray(index).length > 1) {
			for (let i = 0; i < this.getModuleTimingControlArray(index).length; i++) {
				if (i == element) {
					// console.log(i);
					continue
				}
				this.getModuleTimingControlArray(index).at(i).get('fromTime').setValue(this.getModuleTimingControlArray(index).at(element).get('fromTime').value)
				this.getModuleTimingControlArray(index).at(i).get('toTime').setValue(this.getModuleTimingControlArray(index).at(element).get('toTime').value)
				// this.getModuleTimingControlArray(index).at(i).get('fromTime').disable()
				// this.getModuleTimingControlArray(index).at(i).get('toTime').disable()

			}
		}

	}

	toggleModuleTimingValidators(i: number): void {
		this.batchModuleControlArray.at(i).get("isModuleTimingMarked").
			setValue(!this.batchModuleControlArray.at(i).get("isModuleTimingMarked").value)

		if (this.batchModuleControlArray.at(i).get("isModuleTimingMarked").value) {
			this.addModuleTimingsValidators(this.batchModuleControlArray.at(i).get("moduleTimings"))
			return
		}

		this.removeModuleTimingsValidators(this.batchModuleControlArray.at(i).get("moduleTimings"))
	}

	addModuleTimingsValidators(formControl: AbstractControl): void {

		let moduleTimingsControl = (formControl as FormArray).controls
		for (let index = 0; index < moduleTimingsControl.length; index++) {
			const element = moduleTimingsControl[index];

			if (element.get("isMarked").value) {
				element.get("dayID").setValidators(Validators.required)
				element.get("fromTime").setValidators(Validators.required)
				element.get("toTime").setValidators(Validators.required)
				this.utilService.updateValueAndValiditors(element as FormGroup)
			}
		}

	}

	getTotalDays(i: number): number {
		this.totalDaysSelected = 0
		let moduleTimingsControl = (this.batchModuleControlArray.at(i).get("moduleTimings") as FormArray).controls
		for (let index = 0; index < moduleTimingsControl.length; index++) {
			const element = moduleTimingsControl[index];
			// console.log(element);
			if (element.get("isMarked").value) {
				// console.log(this.totalDaysSelected);
				this.totalDaysSelected++
				element.get("dayID").setValidators(Validators.required)
				element.get("fromTime").setValidators(Validators.required)
				element.get("toTime").setValidators(Validators.required)
				this.utilService.updateValueAndValiditors(element as FormGroup)
			}
		}
		return this.totalDaysSelected
	}

	removeModuleTimingsValidators(formControl: AbstractControl): void {
		let moduleTimingsControl = (formControl as FormArray).controls
		for (let index = 0; index < moduleTimingsControl.length; index++) {
			const element = moduleTimingsControl[index];

			element.get("isMarked").setValue(false)
			element.get("dayID").setValidators(null)
			element.get("fromTime").setValidators(null)
			element.get("toTime").setValidators(null)
			this.utilService.updateValueAndValiditors(element as FormGroup)
		}
		this.totalDaysSelected = 0
	}

	onModuleTimeChange(i: number, j: number, moduleTimingForm: FormGroup): void {
		moduleTimingForm.get('isMarked').setValue(!moduleTimingForm.get('isMarked').value)

		if (j > 0 && moduleTimingForm.get('isMarked').value) {
			let fromTime = this.getModuleTimingControlArray(i).at(0).get("fromTime").value
			let toTime = this.getModuleTimingControlArray(i).at(0).get("toTime").value

			// console.log("fromtime -> ", fromTime, " toTime -> ", toTime);
			moduleTimingForm.get("fromTime").setValue(fromTime)
			moduleTimingForm.get("toTime").setValue(toTime)
		}
	}

	deleteModuleTiming(i: number, j: number): void {
		this.getModuleTimingControlArray(i).at(j).get("isMarked").setValue(false)
	}

	async onModuleClick(event: any, batch: IBatch): Promise<void> {
		// try {
		this.urlService.setUrlFrom(this.constructor.name, this.router.url)
		if (batch.course.id != null) {
			this.router.navigate(["session/details"], {
				relativeTo: this.route,
				queryParams: {
					"batchID": batch.id,
					"courseID": batch.course.id,
					"batchName": batch.batchName,
					"moduleTab": "true",
				}
			}).catch(err => {
				console.error(err)

			});
		} else {
			alert("Course ID Doesn't Exit")
		}
		// 	event.stopPropagation()
		// 	this.selectedBatch = batch
		// 	console.log("selectedBatch", this.selectedBatch)

		// 	this.createBatchModuleForm()

		// 	this.activeModuleTab = 1
		// 	this.moduleIDList = []

		// 	if (this.selectedBatch.course) {
		// 		this.courseModules = await this.getCourseModules(this.selectedBatch.course?.id)
		// 	}
		// 	console.log(this.courseModules);

		// 	this.batchModules = await this.getBatchModules(this.selectedBatch.id)
		// 	console.log("batchModules", this.batchModules)

		// 	this.handleBatchModuleForm(true)
		// 	this.openModal(this.moduleModal, "xl")
		// } catch (error) {
		// 	console.error(error);
		// }
	}

	handleBatchModuleForm(isCourseAdded: boolean): void {
		if (isCourseAdded) {
			this.addCourseModuleToForm()
		}

		this.addBatchModuleIDToList()
		this.disableBatchModulesInForm()
	}

	async onTabChange(activeModuleTab: number): Promise<void> {
		this.activeModuleTab = activeModuleTab
		// console.log(this.activeModuleTab);
		try {
			switch (this.activeModuleTab) {
				case 1:
					this.courseModules = await this.getCourseModules(this.selectedBatch.course?.id)
					this.batchModules = await this.getBatchModules(this.selectedBatch.id)
					console.log(" === batchModules -> ", this.batchModules);

					this.handleBatchModuleForm(true)
					break
				case 2:
					this.batchModules = await this.getBatchModules(this.selectedBatch.id)
					this.handleBatchModuleForm(false)
					break
				default:
					this.activeModuleTab = 1
					this.courseModules = await this.getCourseModules(this.selectedBatch.course?.id)
					this.handleBatchModuleForm(true)
					break
			}
		} catch (error) {
			console.error(error);
		}
	}

	toggleModule(batchModuleForm: FormGroup): void {
		// console.log(batchModuleForm);
		batchModuleForm.disable()
		this.removeFormValueAndValidators(batchModuleForm)
		batchModuleForm.get("isMarked").setValue(!batchModuleForm.get("isMarked").value)

		if (batchModuleForm.get("isMarked").value) {
			this.setModuleFormValidators(batchModuleForm)
			batchModuleForm.enable()
			return
		}
	}

	removeFormValueAndValidators(batchModuleForm: FormGroup): void {
		batchModuleForm.get("facultyID").clearValidators()
		batchModuleForm.get("moduleID").clearValidators()
		batchModuleForm.get("order").clearValidators()
		batchModuleForm.get("startDate").clearValidators()

		batchModuleForm.get("facultyID").setValue(null)
		batchModuleForm.get("order").setValue(null)
		batchModuleForm.get("startDate").setValue(null)

		this.utilService.updateValueAndValiditors(batchModuleForm)
	}

	setModuleFormValidators(batchModuleForm: FormGroup): void {
		batchModuleForm.get("facultyID").setValidators([Validators.required])
		batchModuleForm.get("moduleID").setValidators([Validators.required])
		batchModuleForm.get("order").setValidators([Validators.required, Validators.min(1)])
		// batchModuleForm.get("startDate").setValidators([Validators.required])

		this.utilService.updateValueAndValiditors(batchModuleForm)
	}

	async getCourseModules(courseID: string): Promise<ICourseModule[]> {
		try {
			return await new Promise<ICourseModule[]>((resolve, reject) => {
				this.spinnerService.loadingMessage = "Getting course modules";
				;
				this.courseModules = [];
				this.totalCourseModules = 0;
				let queryParams: any = {
					isActive: "1",
					limit: -1,
					offset: 0,
				};
				this.courseModuleService.getCourseModules(courseID, queryParams).subscribe((response) => {
					// this.courseModules = response.body
					resolve(response.body);
					this.totalCourseModules = parseInt(response.headers.get("X-Total-Count"));
				}, (err: any) => {
					console.error(err);
					if (err.statusText.includes('Unknown')) {
						reject("No connection to server. Check internet.");
						return;
					}
					reject(err.error.error);
				});
			});
		} finally {
			;
		}
	}

	disableBatchModulesInForm(): void {
		let batchModuleForm: IBatchModule[] = this.batchModuleForm.value.batchModules
		for (let index = 0; index < this.batchModules?.length; index++) {
			for (let j = 0; j < batchModuleForm?.length; j++) {
				if (this.batchModules[index].module.id == batchModuleForm[j].moduleID) {
					this.batchModuleControlArray.at(j).get("id").setValue(this.batchModules[index].id)
					this.batchModuleControlArray.at(j).get("facultyID").setValue(this.batchModules[index].faculty?.id)
					this.batchModuleControlArray.at(j).get("isCompleted").setValue(this.batchModules[index].isCompleted)
					this.batchModuleControlArray.at(j).get("startDate").
						setValue(this.datePipe.transform(this.batchModules[index].startDate, "yyyy-MM-dd"))
					if (this.batchModules[index].estimatedEndDate != null) {
						this.batchModuleControlArray.at(j).get("estimatedEndDate").
							setValue(this.datePipe.transform(this.batchModules[index].estimatedEndDate, "yyyy-MM-dd"))
					}
					this.batchModuleControlArray.at(j).get("order").setValue(this.batchModules[index].order)
					this.batchModuleControlArray.at(j).get("isMarked").setValue(true)
					this.batchModuleControlArray.at(j).get("isModuleTimingMarked").setValue(false)

					if (this.batchModules[index]?.moduleTimings?.length > 0) {
						this.batchModuleControlArray.at(j).get("isModuleTimingMarked").setValue(true)
					}

					for (let k = 0; k < this.batchModules[index]?.moduleTimings?.length; k++) {
						const element = this.batchModules[index]?.moduleTimings[k];

						this.getModuleTimingControlArray(j).at(element.day.order - 1).get("id").setValue(element.id)
						// this.getModuleTimingControlArray(j).at(element.day.order-1).get("batchID").setValue(element.batchID)
						this.getModuleTimingControlArray(j).at(element.day.order - 1).get("moduleID").setValue(element.moduleID)
						// this.getModuleTimingControlArray(j).at(element.day.order-1).get("facultyID").setValue(element.facultyID)
						this.getModuleTimingControlArray(j).at(element.day.order - 1).get("fromTime").setValue(element.fromTime)
						this.getModuleTimingControlArray(j).at(element.day.order - 1).get("toTime").setValue(element.toTime)
						this.getModuleTimingControlArray(j).at(element.day.order - 1).get("isMarked").setValue(true)
					}
					this.batchModuleControlArray.at(j).enable()
				}
			}
		}
	}

	addCourseModuleToForm(): void {
		this.createBatchModuleForm()
		for (let index = 0; index < this.courseModules.length; index++) {
			this.addBatchModulesToForm()
			let len = this.batchModuleControlArray.controls?.length
			this.batchModuleControlArray.at(len - 1).get("moduleID").setValue(this.courseModules[index].module?.id)
			this.addModuleTiming(index, this.courseModules[index].module?.id)
			this.batchModuleControlArray.at(len - 1).disable()
			// this.batchModuleControlArray.at(len - 1).get("order").setValue(this.courseModules[index].order)
		}
	}

	addModuleTiming(j: number, moduleID: string): void {
		for (let index = 0; index < this.daysList.length; index++) {
			this.addModuleTimingToForm(j)
			let len = this.getModuleTimingControlArray(j).controls?.length
			this.getModuleTimingControlArray(j).at(len - 1).get("dayID").setValue(this.daysList[index].id)
			this.getModuleTimingControlArray(j).at(len - 1).get("moduleID").setValue(moduleID)
			this.getModuleTimingControlArray(j).at(len - 1).get("day").setValue(this.daysList[index])
			this.getModuleTimingControlArray(j).at(len - 1).get("isMarked").setValue(false)
		}
	}

	async getBatchModules(batchID: string): Promise<IBatchModule[]> {
		try {
			return await new Promise<IBatchModule[]>((resolve, reject) => {
				this.spinnerService.loadingMessage = "Getting batch modules";
				;
				this.batchModules = [];
				this.totalBatchModules = 0;
				let queryParams: any = {
					limit: -1,
					offset: 0
				}
				this.batchService.getBatchModules(batchID, queryParams).subscribe((response) => {
					resolve(response.body);
					this.totalBatchModules = parseInt(response.headers.get("X-Total-Count"));
				}, (err: any) => {
					console.error(err);
					if (err.statusText.includes('Unknown')) {
						reject("No connection to server. Check internet.");
						return;
					}
					reject(err.error.error);
				});
			});
		} finally {
			;
		}
	}

	addBatchModuleIDToList(): void {
		for (let index = 0; index < this.batchModules?.length; index++) {
			this.moduleIDList.push(this.batchModules[index].module.id)
		}
	}

	async addBatchModule(batchModule: IBatchModule): Promise<void> {
		try {
			return await new Promise<any>((resolve, reject) => {

				this.spinnerService.loadingMessage = "Updating batch modules"

				this.batchService.addBatchModule(this.selectedBatch.id, batchModule).subscribe((response) => {
					// console.log(response);
					resolve(response)
				}, (err: any) => {
					// console.error(err)
					reject(err)
				})
			})
		} finally {

		}
	}

	async updateBatchModule(batchModule: IBatchModule): Promise<void> {
		try {
			return await new Promise<any>((resolve, reject) => {

				this.spinnerService.loadingMessage = "Updating batch modules"

				this.batchService.updateBatchModule(this.selectedBatch.id, batchModule).subscribe((response) => {
					// console.log(response);
					resolve(response)
				}, (err: any) => {
					// console.error(err)
					reject(err)
				})
			})

		} finally {

		}
	}

	validateBatchModule(): void {
		console.log(this.batchModuleForm.controls);

		if (!this.batchModuleForm.valid) {
			this.utilService.markGroupDirty(this.batchModuleForm)
			return
		}
		this.submitBatchModules()
	}

	async submitBatchModules(): Promise<void> {
		let batchModules: IBatchModule[] = this.batchModuleForm.value.batchModules
		let errors: string[] = []
		// console.log(batchModules);

		this.setModuleTimings(batchModules)

		for (let index = 0; index < batchModules?.length; index++) {
			console.log(" module -> ", batchModules[index])
			try {
				if (batchModules[index]?.id) {
					await this.updateBatchModule(batchModules[index])
					continue
				}
				await this.addBatchModule(batchModules[index])
			} catch (err) {
				errors.push(err?.error?.error)
				console.error(err);
			}
		}

		if (errors.length != 0) {
			let errorString = ""
			for (let index = 0; index < errors.length; index++) {
				errorString += (index == 0 ? "" : "\n") + errors[index]
			}
			alert(errorString)
			return
		}

		alert("Modules successufully saved")
		this.onTabChange(2)
		this.searchBatch()

	}

	setModuleTimings(batchModules: IBatchModule[]): void {
		for (let index = 0; index < batchModules.length; index++) {
			for (let j = 0; j < batchModules[index].moduleTimings?.length; j++) {
				if (!batchModules[index].moduleTimings[j].isMarked) {
					batchModules[index].moduleTimings.splice(j, 1)
					j = j - 1
					continue
				}
				batchModules[index].moduleTimings[j].facultyID = batchModules[index].facultyID
			}
		}
	}

	onDeleteModuleClick(batchModuleID: string): void {
		this.waitingListMessage = null
		this.entity = null

		this.openModal(this.deleteModal, "md").result.then(() => {
			this.deleteBatchModule(batchModuleID)
		}, (err) => {
			console.error(err);
			return
		})
	}

	deleteBatchModule(batchModuleID: string): void {
		this.spinnerService.loadingMessage = "Deleting module"

		this.batchService.deleteBatchModule(this.selectedBatch.id, batchModuleID).subscribe((response: any) => {
			// console.log(response);
			alert("Batch module successfully deleted")
			this.onTabChange(2)
		}, (err: any) => {
			console.error(err);
			console.error(err);
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.");
				return;
			}
			alert(err.error?.error)
		}).add(() => {

		})
	}

	// ====================================================== CRUD END ======================================================

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
			// console.log(respond);
			this.batchStatus = respond;
		}, (err) => {
			console.error(err.error.error)
		})
	}

	// Get All Batch objectives
	getBatchObjectives(): void {
		this.generalService.getGeneralTypeByType("batch_objective").subscribe((respond: any) => {
			// console.log(respond);
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
	getTechnology(event?: any): void {
		let queryParams: any = {}
		if (event && event?.term != "") {
			queryParams.language = event.term
		}
		this.technologies = []
		this.isTechLoading = true
		this.techService.getAllTechnologies(this.techLimit, this.techOffset, queryParams).subscribe((response) => {
			// console.log("getTechnology -> ", response);
			this.technologies = response.body
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
			// console.log(respond);
			this.academicYears = respond;
		}, (err) => {
			console.error(err.error.error)
		})
	}

	// Get All Sales Persons
	getAllSalesPeople(): void {
		this.generalService.getSalesPersonList().subscribe((data: any) => {
			this.salesPeople = this.salesPeople.concat(data.body)
		}, (err) => {
			console.error(err)
		})
	}

	// Get All Faculty List
	getFacultyList(): void {
		this.generalService.getFacultyList().subscribe((data: any) => {
			this.facultyList = this.facultyList.concat(data.body)
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

	// Get active company requirement list.
	getActiveRequirementList(): void {
		let queryParams: any = {
			isActive: "1"
		}
		this.generalService.getRequirementList(queryParams).subscribe((response: any) => {
			this.requirementList = response.body
		}, (err: any) => {
			console.error(err);
		})
	}

	// Get active batch list.
	getActiveBatchList(): void {
		let queryParams: any = {
			batchStatus: ["Ongoing", "Upcoming"],
			isActive: "1"
		}
		this.generalService.getBatchList(queryParams).subscribe((response: any) => {
			this.activeBatchList = response
		}, (err: any) => {
			console.error(err);
		})
	}

	// Get count of faculty supervisors.
	getSupervisorCount(): void {
		let queryParams: any = {
			roleName: this.localService.getJsonValue("roleName"),
	  	}
		this.adminService.getSupervisorCount(queryParams).subscribe(response => {
			this.supervisorFacultyCount = response.body.totalCount
			if (this.supervisorFacultyCount == 0) {
				this.isHeadFcaulty = true
			}
			if (this.supervisorFacultyCount > 0) {
				this.isHeadFcaulty = false
			}
		}, err => {
			console.error(this.utilService.getErrorString(err))
		})
	}

	// Toggle between my and all talents.
	toggleMyAllTalents(): void {
		this.isViewAllBatches = !this.isViewAllBatches
		if (this.isSearched) {
			this.getAllSearchedTalents()
			return
		}
		this.getBatches()
	}

}
