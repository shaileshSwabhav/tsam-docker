import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, Validators, FormArray, FormControl } from '@angular/forms';
import { GeneralService, IFeedbackQuestion, ISpecialization, IState } from 'src/app/service/general/general.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { FacultyService, IFacultyAssessment } from 'src/app/service/faculty/faculty.service';
import { BatchService, IBatchTiming } from 'src/app/service/batch/batch.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { NgbModal, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { IPermission } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { ActivatedRoute, Router } from '@angular/router';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { UrlConstant } from 'src/app/service/constant';
import { DatePipe } from '@angular/common';
import { IExperienceDTO } from 'src/app/service/talent/talent.service';
import { CourseService, ICourse, ICourseTechnicalAssessment } from 'src/app/service/course/course.service';
import { TechnologyService } from 'src/app/service/technology/technology.service';
import { DegreeService } from 'src/app/service/degree/degree.service';

@Component({
	selector: 'app-faculty-master',
	templateUrl: './faculty-master.component.html',
	styleUrls: ['./faculty-master.component.css']
})
export class FacultyMasterComponent implements OnInit {

	// form
	facultyForm: FormGroup;
	modelAction: () => void;
	formAction: () => void;
	facultySearchForm: FormGroup;

	// general types
	selectedCountry: any
	stateList: any[]
	countryList: any[]
	faculties: any[]
	selectedFaculty: any
	qualificationList: any[]
	specializationList: ISpecialization[]
	technologyList: any[]
	designationList: any[]
	collegeBranchList: any[]
	currentDate: Date

	// faculty
	facultyID: string

	// batch
	totalBatch: number
	batchList: any[];
	batchLimit: number
	batchOffset: number
	currentBatchPageNumber: number
	isBatchList: boolean;

	// modal
	modalHeader: string;
	modalButton: string;
	otherQualification: boolean[];
	otherDesignation: boolean[];
	otherDegree: boolean[];
	modalRef: NgbModalRef;
	searchFormValue: any;

	// spinner


	// number
	limit: number;
	offset: number;
	currentPage: number;
	totalFaculty: number;
	paginationString: string
	techLimit: number
	techOffset: number
	collegeLimit: number
	collegeOffset: number


	// boolean
	isViewClicked: boolean
	isAddClicked: boolean
	searched: boolean
	experienced: boolean;
	showSearch: boolean
	apiBtn: boolean
	isTechLoading: boolean
	isCollegeLoading: boolean

	// permission
	permission: IPermission
	loginID: string
	isAdmin: boolean

	// resume
	isResumeUploadedToServer: boolean
	isFileUploading: boolean
	docStatus: string
	displayedFileName: string

	indexOfCurrentWorkingExp: number;

	// map
	specializationMap: Map<string, ISpecialization[]> = new Map()

	// faculty assessment
	facultyAssessmentForm: FormGroup
	facultyAssementQuestions: IFeedbackQuestion[]
	facultyAssessment: IFacultyAssessment[]
	viewFacultyAssessment: boolean
	isFacultyAssessment: boolean
	averageAssessmentScore: number

	// course technical assessment
	technicalAssessments: ICourseTechnicalAssessment[]
	technicalAssessmentForm: FormGroup
	viewTechnicalAssessment: boolean
	isTechnicalAssessment: boolean
	isUpdateAssessmentClick: boolean

	// course
	courseList: any[]

	private readonly IGNORE_SEARCH_FIELD: string[] = ["limit", "offset"]

	@ViewChild("drawer") drawer: any
	@ViewChild("facultyModal") facultyModal: any
	@ViewChild("facultyAssessmentModal") facultyAssessmentModal: any
	@ViewChild("technicalAssessmentModal") technicalAssessmentModal: any
	@ViewChild("addBatchToFacultyModal") addBatchToFacultyModal: any
	@ViewChild("deleteConfirmationModal") deleteConfirmationModal: any

	constructor(
		private formBuilder: FormBuilder,
		private generalService: GeneralService,
		public utilService: UtilityService,
		private facultyService: FacultyService,
		private courseService: CourseService,
		private batchService: BatchService,
		private techService: TechnologyService,
		private degreeService: DegreeService,
		private fileOpsService: FileOperationService,
		private localService: LocalService,
		private urlConstant: UrlConstant,
		private spinnerService: SpinnerService,
		private modalService: NgbModal,
		private router: Router,
		private route: ActivatedRoute,
		private datePipe: DatePipe
	) {
		this.initializeVariables()
		this.createForms()
		this.getAllComponents();
	}

	initializeVariables(): void {
		this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.ADMIN_FACULTY)
		this.loginID = this.localService.getJsonValue("loginID")
		this.isAdmin = (this.localService.getJsonValue("roleName") == "Admin" ? true : false)

		this.spinnerService.loadingMessage = "Getting faculties"
		this.displayedFileName = "Select file"

		this.isViewClicked = false
		this.isAddClicked = false
		this.searched = false
		this.isBatchList = false
		this.experienced = false
		this.showSearch = false
		this.apiBtn = false
		this.isFileUploading = false
		this.isResumeUploadedToServer = false
		this.isTechLoading = false

		// faculty assessment
		this.viewFacultyAssessment = true
		this.isFacultyAssessment = true

		// course technical assessment
		this.viewTechnicalAssessment = false
		this.isTechnicalAssessment = true
		this.isUpdateAssessmentClick = false

		this.limit = 5;
		this.offset = 0;
		this.batchLimit = 5;
		this.batchOffset = 0;
		this.totalFaculty = 0
		this.totalBatch = 0
		this.indexOfCurrentWorkingExp = -1;
		this.averageAssessmentScore = 0
		this.techLimit = 10
		this.techOffset = 0

		this.currentDate = new Date()

		this.stateList = []
		this.specializationList = []
		this.facultyAssessment = []
		this.technicalAssessments = []
		this.collegeBranchList = []
		// console.log(this.currentDate.getFullYear() - 1)
	}

	// Get Required List.
	getAllComponents(): void {
		this.getCountryList()
		this.getDesignationList()
		this.getQualificationList()
		this.getTechnologyList()
		this.getSpecializationList()
		this.getCollegeBranchList()
		this.getFeedbackQuestionsForFaculty()
		this.searchOrGetFaculties()
		// this.getFaculties()
	}

	createForms(): void {
		this.createFacultySearchForm()
	}


	get ongoingOperations() {
		return this.spinnerService.ongoingOperations
	}

	ngOnInit() {
		// this.setDefaultForOtherField();
	}

	// Create Faculty Form Object.
	createFacultyForm(): void {
		this.facultyForm = this.formBuilder.group({
			id: new FormControl(),
			code: new FormControl(null),
			firstName: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z]+$/)]),
			lastName: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z]+$/)]),
			contact: new FormControl(null, [Validators.required, Validators.pattern(/^[6789]\d{9}$/)]),
			dateOfBirth: new FormControl(null),
			dateOfJoining: new FormControl(null),
			email: new FormControl(null, [Validators.required, Validators.pattern(/^([a-zA-Z0-9._%+-]+@[a-zA-Z0-9]+\.[a-zA-z]{2,3})/)]),
			technologies: new FormControl(null),
			resume: new FormControl(null),
			isActive: new FormControl(true, [Validators.required]),
			isFullTime: new FormControl(null, [Validators.required]),
			address: new FormControl(null, [Validators.required, Validators.pattern(/^[.0-9a-zA-Z\s,-\/]+$/)]),
			state: new FormControl(null, [Validators.required]),
			city: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z\s*]+[a-zA-Z]+$/)]),
			country: new FormControl(null, Validators.required),
			pinCode: new FormControl(null, [Validators.required, Validators.pattern(/^[0-9]{6}$/)]),
			telegramID: new FormControl(null),
			academics: this.formBuilder.array([]),
			experiences: this.formBuilder.array([]),
			// isExperience: [false],
			// contact: [null, [Validators.required, Validators.pattern(/^(?:(?:\+|0{0,2})91(\s*[\-]\s*)?|[0]?)?[789]\d{9}$/)]],
		});
		this.addFacultyAcademic();
		this.facultyForm.get('state')?.disable()
	}

	// Add Faculty Academic Detail.
	addFacultyAcademic(): void {
		this.facultyAcademicControlsArray.push(this.formBuilder.group({
			id: new FormControl(),
			degree: new FormControl(null, [Validators.required]),
			specialization: new FormControl(null, [Validators.required]),
			college: new FormControl(null, [Validators.required]), //, Validators.pattern(/^[a-zA-Z\s*]+[a-zA-Z]+$/)
			percentage: new FormControl(null, [Validators.required, Validators.pattern(/^(\d{1,2}(\.\d{1,2})?)$/)]),
			passout: new FormControl(null, [Validators.required, Validators.pattern(/^(19[8-9]\d|20[0-4]\d|2020)$/),
			Validators.max(this.currentDate.getFullYear())]),
			specializationList: new FormControl(null),
		}))
		// this.otherDegree.push(false)
	}

	// Add Faculty Experience Details.
	addFacultyExperience(): void {
		this.facultyExperienceControlsArray.push(this.formBuilder.group({
			id: new FormControl(),
			companyName: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z\s*]+[a-zA-Z]+$/)]),
			designation: new FormControl(null, [Validators.required]),
			technologies: new FormControl(null, [Validators.required]),
			isCurrentlyWorking: new FormControl(false),
			yearsOfExperience: new FormControl(null),
			fromDate: new FormControl(null),
			toDate: new FormControl(null),
			package: new FormControl(null, [Validators.min(100000), Validators.max(100000000)]),
		}));
		// , { validators: fromAndToDateValidator }
		// this.otherDesignation.push(false);
	}

	// Faculty Academic Array Controls Return
	get facultyAcademicControlsArray() {
		return this.facultyForm.get('academics') as FormArray
	}

	// Faculty Experience Array Controls Return
	get facultyExperienceControlsArray() {
		return this.facultyForm.get('experiences') as FormArray
	}

	createFacultySearchForm(): void {
		this.facultySearchForm = this.formBuilder.group({
			firstName: new FormControl(null, [Validators.pattern(/^[a-zA-Z ]+$/)]),
			lastName: new FormControl(null, [Validators.pattern(/^[a-zA-Z ]+$/)]),
			contact: new FormControl(null),
			city: new FormControl(null),
			email: new FormControl(null),
			technologies: new FormControl(null),
			isActive: new FormControl(null),
			limit: new FormControl(this.limit),
			offset: new FormControl(this.offset)
		});
	}

	changePage(pageNumber: number): void {

		// this.offset = $event - 1
		// this.currentPage = $event;

		this.facultySearchForm.get("offset").setValue(pageNumber - 1)

		this.limit = this.facultySearchForm.get("limit").value
		this.offset = this.facultySearchForm.get("offset").value

		console.log(this.facultySearchForm.value);


		this.searchFaculty()
	}

	// Compare Ob1 and Ob2
	compareFn(ob1: any, ob2: any): boolean {
		if (ob1 == null && ob2 == null) {
			return true;
		}
		return ob1 && ob2 ? ob1.id === ob2.id : ob1 === ob2
	}

	// On Experience Checkbox Click
	experienceChecked(event): void {
		event.target.checked = true
		let formValid = false
		if (this.experienced) {
			// console.log(this.facultyExperienceControlsArray);
			for (let i = 0; i < this.facultyExperienceControlsArray.length; i++) {
				if (this.facultyExperienceControlsArray.at(i).valid) {
					formValid = true
					break
				}
			}

			if (formValid) {
				if (confirm("This will delete all the experiences for this faculty. Are you sure you want to delete it?")) {
					event.target.checked = false
					this.facultyForm.markAsDirty()
				} else {
					return
				}
			}
			this.facultyForm.setControl('experiences', this.formBuilder.array([]))
			this.experienced = false
			return
			// this.facultyForm.get("isExperience").setValue(false);
		}
		this.experienced = true;
		// this.facultyForm.get("isExperience").setValue(true);
		this.addFacultyExperience();
	}

	// Delete Experience.
	deleteExperience(index: number): void {
		// console.log(this.facultyExperienceControlsArray.at(index));

		if (this.facultyExperienceControlsArray.at(index).valid) {
			if (confirm("This will delete the faculty experience. Are you sure you want to delete it?")) {
				if (this.indexOfCurrentWorkingExp == index) {
					this.indexOfCurrentWorkingExp = -1;
				}
				if (this.indexOfCurrentWorkingExp > index) {
					this.indexOfCurrentWorkingExp = this.indexOfCurrentWorkingExp - 1;
				}
				this.facultyExperienceControlsArray.removeAt(index)
				this.facultyForm.markAsDirty()
				// return
			}
			return
		}
		this.facultyExperienceControlsArray.removeAt(index);
		this.facultyForm.markAsDirty()
		return
	}

	// Delete Academics
	deleteAcademic(index: number): void {
		// console.log(this.facultyAcademicControlsArray.at(index));
		if (this.facultyAcademicControlsArray.at(index).valid) {
			if (confirm("This will delete the faculty academic. Are you sure you want to delete it?")) {
				this.facultyAcademicControlsArray.removeAt(index)
				this.facultyForm.markAsDirty()
				// return
			}
			return
		}
		this.facultyAcademicControlsArray.removeAt(index);
		this.facultyForm.markAsDirty()
		return
	}

	//Show specialization according to degree id
	showSpecificSpecializations(academic: any, specialization: any): boolean {

		if (academic.get('degree').value == null) {
			return false
		}
		if (academic.get('degree').value.id === specialization.degreeID) {
			return true
		}
		return false
	}

	doesSpecializationExist(academic: any): boolean {
		for (let index = 0; index < this.specializationList.length; index++) {
			if (academic.get('degree').value.id === this.specializationList[index].degreeID) {
				return false
			}
		}
		return true
	}

	setMaxToDate(value: any, index: number): void {
		// console.log(value);
		// this.facultyExperienceControlsArray.at(index).get('toDate').setValidators([Validators.required, Validators.min(value)])
	}


	// Update Faculty Form.
	updateFacultyForm(): void {
		this.createFacultyForm();
		this.displayedFileName = "No resume uploaded"
		if (this.selectedFaculty.resume) {
			this.displayedFileName = `<a href=${this.selectedFaculty.resume} target="_blank">Resume present</a>`
		}
		this.facultyForm.setControl("academics", this.formBuilder.array([]))
		this.facultyForm.setControl("experiences", this.formBuilder.array([]));
		// this.facultyForm.get("isExperience").setValue(false);

		if (this.selectedFaculty.experiences && this.selectedFaculty.experiences.length > 0) {
			for (let i = 0; i < this.selectedFaculty.experiences.length; i++) {
				if (this.selectedFaculty.experiences[i].toDate == null) {
					this.indexOfCurrentWorkingExp = i
				}
				// if (this.selectedFaculty.experiences[i].fromDate != null) {
				//       this.selectedFaculty.experiences[i].fromDate = this.selectedFaculty.experiences[i].fromDate.slice(0, 10)
				// }
				// if (this.selectedFaculty.experiences[i].toDate != null) {
				//       this.selectedFaculty.experiences[i].toDate = this.selectedFaculty.experiences[i].toDate.slice(0, 10)
				// }
				this.calculateYearsOfExperience(this.selectedFaculty.experiences[i]);
				this.addFacultyExperience();
			}
			this.experienced = true;
		}

		for (let ind = 0; ind < this.selectedFaculty.academics.length; ind++) {
			this.addFacultyAcademic()
		}
		if (this.selectedFaculty.dateOfBirth != null) {
			this.selectedFaculty.dateOfBirth = this.selectedFaculty.dateOfBirth.slice(0, 10)
		}
		if (this.selectedFaculty.dateOfJoining != null) {
			this.selectedFaculty.dateOfJoining = this.selectedFaculty.dateOfJoining.slice(0, 10)
		}
		this.getStateList(this.selectedFaculty.country.id);
		// for (let i = 0; i < this.selectedFaculty.academics.length; i++) {
		//       this.getSpecializationListForDegree(this.selectedFaculty.academics[i].degree.id)
		// }

		// console.log(this.selectedFaculty);
		this.facultyForm.patchValue(this.selectedFaculty);

		if (this.indexOfCurrentWorkingExp != -1) {
			this.facultyExperienceControlsArray.at(this.indexOfCurrentWorkingExp).get('isCurrentlyWorking').setValue('true')
			let control = this.facultyExperienceControlsArray.controls[this.indexOfCurrentWorkingExp]
			if (control instanceof FormGroup) {
				control.removeControl('toDate')
			}
		}
	}

	// On Add New Faculty Click
	onAddNewFacultyClick(): void {
		this.isAddClicked = true
		this.isViewClicked = false
		this.indexOfCurrentWorkingExp = -1;
		this.experienced = false;
		this.updatemodelParamAccordingToAction("Add New Faculty", "Add Faculty", this.addFaculty, this.createFacultyForm);
		this.createFacultyForm();
		this.openModal(this.facultyModal, "xl")
	}

	// On Update Faculty Click
	onUpdateFacultyClick(): void {
		this.facultyForm.enable()
		this.updatemodelParamAccordingToAction("Update Faculty", "Update Faculty", this.updateFaculty, this.updateFacultyForm)
		this.isAddClicked = false
		this.isViewClicked = false
		this.displayedFileName = `<a href=${this.selectedFaculty.resume} target="_blank">${this.selectedFaculty.firstName}
                 ${this.selectedFaculty.lastName} </a>`
	}

	onViewFacultyClick(faculty: any) {
		this.indexOfCurrentWorkingExp = -1
		this.selectedFaculty = faculty
		this.experienced = false
		this.isViewClicked = true
		this.isAddClicked = false
		this.createFacultyForm()
		this.updateFacultyForm()
		this.facultyForm.disable()
		this.updatemodelParamAccordingToAction("Faculty Detail", "Faculty", this.addFaculty, this.createFacultyForm)
		this.openModal(this.facultyModal, 'xl')
	}

	//On clicking currently working or not
	isWorkingClicked(isWorking: string, index: number, expereince: any): void {
		if (isWorking == "false") {
			expereince.addControl('toDate', this.formBuilder.control(null, [Validators.required]));
			this.indexOfCurrentWorkingExp = -1;
			return
		}
		expereince.removeControl('toDate');
		this.indexOfCurrentWorkingExp = index
	}

	//Calculate total years of experience for one experience
	calculateYearsOfExperience(experience: IExperienceDTO): void {
		let toDate;
		let fromDate;
		if (experience.fromDate != null && experience.fromDate != "") {
			if (experience.toDate == null || experience.toDate == "") {
				toDate = new Date();
			}
			else {
				toDate = new Date(experience.toDate)
			}
			fromDate = new Date(experience.fromDate)
			if (fromDate > toDate) {
				experience.yearsOfExperience = "-"
				experience.yearsOfExperienceInNumber = 0
				return
			}
			let monthDiff: number = toDate.getMonth() - fromDate.getMonth() + (12 * (toDate.getFullYear() - fromDate.getFullYear()))
			let numberOfYears: number = Math.floor(monthDiff / 12)
			let numberOfMonths: number = Math.round(((monthDiff / 12) % 1) * 12)
			experience.yearsOfExperience = numberOfYears + "." + numberOfMonths + " Year(s)"
			experience.yearsOfExperienceInNumber = +(monthDiff / 12).toFixed(2)
		}
	}

	// Calculate total years of experience for faculty.
	calculateTotalYearsOfExperince(faculty: any): string {
		let totalYears: number = 0;
		faculty.totalYearsOfExperience = totalYears
		if ((faculty.experiences) && (faculty.experiences.length != 0)) {
			for (let i = 0; i < faculty.experiences.length; i++) {
				this.calculateYearsOfExperience(faculty.experiences[i]);
				totalYears = totalYears + faculty.experiences[i].yearsOfExperienceInNumber
			}
			faculty.totalYearsOfExperience = totalYears
			let numberOfYears: number = Math.floor(totalYears)
			let numberOfMonths: number = Math.round((totalYears % 1) * 12)
			if (numberOfYears == 0 && numberOfMonths == 0) {
				faculty.totalYearsOfExperience = 0
				faculty.totalYearsOfExperienceInString = "-"
				return
			}
			if (isNaN(numberOfYears) || isNaN(numberOfMonths)) {
				faculty.totalYearsOfExperience = 0
				faculty.totalYearsOfExperienceInString = "-"
				return
			}
			faculty.totalYearsOfExperienceInString = numberOfYears + "." + numberOfMonths + " Year(s)"
			return
		}
		faculty.totalYearsOfExperience = 0
		faculty.totalYearsOfExperienceInString = "-"
	}

	setPaginationString() {
		this.paginationString = ''
		let start: number = this.limit * this.offset + 1
		let end: number = +this.limit + this.limit * this.offset
		if (this.totalFaculty < end) {
			end = this.totalFaculty
		}
		if (this.totalFaculty == 0) {
			this.paginationString = ''
			return
		}
		this.paginationString = `${start} - ${end}`
	}

	onGetBatchListClick(facultyID: string) {
		this.facultyID = facultyID
		this.openModal(this.addBatchToFacultyModal, "xl")
		this.getBatchList()
		// this.changeBatchPage(1)
	}

	changeBatchPage($event: any): void {
		;
		this.batchOffset = $event - 1
		this.currentBatchPageNumber = $event;
		this.getBatchList()
	}

	resetSearchAndGetAll(): void {
		this.searched = false
		this.showSearch = false
		this.resetSearchFacultyForm()
		this.router.navigate([this.urlConstant.ADMIN_FACULTY])
		this.changePage(1)
	}

	resetSearchFacultyForm(): void {
		this.limit = this.facultySearchForm.get("limit").value
		this.offset = this.facultySearchForm.get("offset").value

		this.facultySearchForm.reset({
			limit: this.limit,
			offset: this.offset,
		})
	}

	// Get Faculty Index
	getFacultyIndexByID(id: string): any {
		for (let index = 0; index < this.faculties.length; index++) {
			if (this.faculties[index].id == id) {
				return index;
			}
		}
		return undefined;
	}

	// Assign to SelectFaculty
	assignSelectedFaculty(id: string): void {
		let index = this.getFacultyIndexByID(id);
		if (index != undefined) {
			this.selectedFaculty = this.faculties[index];
		}
		this.openModal(this.deleteConfirmationModal, "md")
	}

	// Update Param According to Model Action Change. like Update Click, View Click, Add New Click.
	updatemodelParamAccordingToAction(modalHeader: string, modalButton: string, modelAction: any, formAction: any): void {
		this.modalHeader = modalHeader;
		this.modalButton = modalButton;
		this.modelAction = modelAction;
		this.formAction = formAction;
	}

	// Reset Form Action.
	resetFacultyForm(): void {
		this.facultyForm.reset();
	}

	searchAndCloseDrawer(): void {
		this.drawer.toggle()
		this.searchFaculty()
	}

	searchFaculty(): void {
		this.searchFormValue = { ...this.facultySearchForm?.value }
		this.router.navigate([], {
			relativeTo: this.route,
			queryParams: this.searchFormValue,
		})

		this.apiBtn = true
		let flag: boolean = true
		for (let field in this.searchFormValue) {
			if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
				delete this.searchFormValue[field];
			} else {
				if (!this.IGNORE_SEARCH_FIELD.includes(field)) {
					this.searched = true
				}
				flag = false
			}
		}
		// No API call on empty search.
		if (flag) {
			this.apiBtn = false
			return
		}

		this.apiBtn = false
		this.getFaculties()
	}

	searchOrGetFaculties(): void {
		let queryParams = this.route.snapshot.queryParams

		if (this.utilService.isObjectEmpty(queryParams)) {
			this.changePage(1)
			return
		}
		this.facultySearchForm.patchValue(queryParams)
		this.searchFaculty()
	}

	calculateDates(): void {
		for (let index in this.faculties) {
			let dateOfBirth = this.faculties[index].dateOfBirth
			let dateOfJoining = this.faculties[index].dateOfJoining
			// if (!dateOfBirth && !dateOfJoining) {
			//       continue
			// }
			if (dateOfBirth) {
				this.faculties[index].dateOfBirth = this.datePipe.transform(dateOfBirth, 'yyyy-MM-dd')
			}
			if (dateOfJoining) {
				this.faculties[index].dateOfJoining = this.datePipe.transform(dateOfJoining, 'yyyy-MM-dd')
			}
			if (this.faculties[index].experiences) {
				for (let j in this.faculties[index].experiences) {
					let fromDate = this.faculties[index].experiences[j].fromDate
					let toDate = this.faculties[index].experiences[j].toDate

					if (!fromDate && !toDate) {
						continue
					}
					if (fromDate) {
						this.faculties[index].experiences[j].fromDate = this.datePipe.transform(fromDate, 'yyyy-MM')
					}
					if (toDate) {
						this.faculties[index].experiences[j].toDate = this.datePipe.transform(toDate, 'yyyy-MM')
					}
				}
			}

		}
	}

	calculateExperience(): any {
		for (let index = 0; index < this.faculties.length; index++) {
			this.calculateTotalYearsOfExperince(this.faculties[index])
		}
	}

	// ==========================================================CRUD START==========================================================

	// Get All Faculty List
	getFaculties(): void {
		if (this.searched) {
			this.getSearchedFaculty()
			return
		}
		this.spinnerService.loadingMessage = "Getting faculties";
		;
		this.faculties = []
		this.totalFaculty = 1

		this.facultyService.getAllFaculty(this.searchFormValue).subscribe(res => {
			this.faculties = res.body;
			this.totalFaculty = Number(res.headers.get('X-Total-Count'))
			this.setPaginationString()
			this.calculateDates()
			this.calculateExperience()

			// console.log(this.faculties)
		}, (err: any) => {
			// this.disableButton = false
			this.totalFaculty = 0
			this.setPaginationString()

			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			// this.addSelectedTalents = []
			console.error(err);
			alert(err?.error?.error)
		})
	}


	// Get All Faculty List
	getSearchedFaculty(): void {
		this.spinnerService.loadingMessage = "Searching faculty";
		;
		this.totalFaculty = 0
		this.faculties = []

		this.facultyService.getAllFaculty(this.searchFormValue).subscribe(res => {
			this.faculties = res.body;
			this.totalFaculty = Number(res.headers.get('X-Total-Count'))
			this.setPaginationString()
			this.calculateDates()
			this.calculateExperience()
			// console.log(this.faculties)

		}, err => {
			// this.disableButton = false

			this.totalFaculty = 0
			this.setPaginationString()
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			// this.addSelectedTalents = []
			console.error(err);
			alert(err?.error?.error)
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
		let isWeekly: boolean = true
		let order: number


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

	getBatchList(): void {
		this.batchList = []
		this.spinnerService.loadingMessage = "Getting batch details";
		;
		this.isBatchList = true;
		this.totalBatch = 0
		let queryParams: any = {
			facultyID: this.facultyID
		}

		this.batchService.getBatches(queryParams).subscribe(res => {
			this.batchList = res.body;
			this.totalBatch = parseInt(res.headers.get('X-Total-Count'))
			// console.log(this.batchList);
			if (this.totalBatch == 0) {
				this.isBatchList = false
			}
			;
		}, err => {
			// this.disableButton = false
			this.isBatchList = false
			this.totalBatch = 0

			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			// this.addSelectedTalents = []
			console.error(err);
			alert(err?.error?.error)
		})
	}


	// Add New Faculty.
	addFaculty(): void {
		this.spinnerService.loadingMessage = "Adding faculty"
		let faculty = this.facultyForm.value;
		this.utilService.deleteNullValueIDFromObject(faculty);
		for (let index = 0; index < faculty.academics?.length; index++) {
			delete faculty.academics[index].specializationList
		}
		console.log(faculty);
		;
		this.apiBtn = true
		this.facultyService.addFaculty(faculty).subscribe(respond => {
			this.apiBtn = false
			this.modalRef.close()
			this.changePage(1);

			alert("Faculty Added Successfully");
		}, err => {
			// this.disableButton = false

			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			// this.addSelectedTalents = []
			console.error(err);
			alert(err?.error?.error)
		})
	}

	// Update Faculty.
	updateFaculty(): void {
		this.apiBtn = true
		this.spinnerService.loadingMessage = "Updating faculty";

		let faculty = this.facultyForm.value;
		this.utilService.deleteNullValueIDFromObject(faculty);
		// console.log(faculty);

		this.facultyService.updateFaculty(faculty).subscribe(respond => {
			this.apiBtn = false
			this.modalRef.close();
			alert("Faculty Updated Successfully")
			this.getFaculties();

		}, err => {
			// this.disableButton = false

			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			// this.addSelectedTalents = []
			console.error(err);
			alert(err?.error?.error)
		})
	}


	// Delete Faculty.
	deleteFaculty(): void {
		this.apiBtn = true
		this.spinnerService.loadingMessage = "Deleting faculty";
		this.modalRef.close();
		;
		if (this.batchList && this.batchList.length > 0) {
			alert("You Can't Delete Faculty there are some batches taken by this faculty.")
				;
			return
		}

		this.facultyService.deleteFaculty(this.selectedFaculty.id).subscribe(respond => {
			this.apiBtn = false
			alert("Faculty Deleted Successfully")
			if (this.faculties.length == 1) {
				this.changePage(this.offset)
			}
			this.getFaculties()

		}, err => {
			// this.disableButton = false

			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			// this.addSelectedTalents = []
			console.error(err);
			alert(err?.error?.error)
		})
	}

	//On uplaoding resume
	onResourceSelect(event: any) {
		this.docStatus = ""
		let files = event.target.files
		if (files && files.length) {
			let file = files[0]

			// Upload resume if it is present.]
			this.isFileUploading = true
			this.spinnerService.isDisabled = true
			this.fileOpsService.uploadResume(file).subscribe((data: any) => {
				this.facultyForm.markAsDirty()
				this.facultyForm.patchValue({
					resume: data
				})
				this.displayedFileName = file.name
				this.isFileUploading = false
				this.isResumeUploadedToServer = true
				this.docStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
			}, (error) => {
				this.isFileUploading = false
				this.docStatus = `<p><span>&#10060;</span> ${error}</p>`
			}).add(() => {
				this.spinnerService.isDisabled = false
			})
		}
	}

	// ==========================================================CRUD END==========================================================


	timeout: NodeJS.Timeout

	searchTechnology(event?: any): void {
		if (this.timeout) {
			clearTimeout(this.timeout)
		}
		this.timeout = setTimeout(() => {
			this.spinnerService.isDisabled = true
			this.getTechnologyList(event)
		}, 500)
	}


	// Get Technology List.
	getTechnologyList(event?: any): void {

		// console.log("get tech called")
		let queryParams: any = {}
		if (event && event?.term != "") {
			queryParams.language = event.term
		}
		this.technologyList = []
		this.isTechLoading = true
		this.techService.getAllTechnologies(this.techLimit, this.techOffset, queryParams).subscribe((response) => {
			// console.log("getTechnology -> ", response);
			this.technologyList = response.body
			// this.technologies = this.technologies.concat(response.body)
		}, (err) => {
			console.error(err)
		}).add(() => {
			this.isTechLoading = false
			this.spinnerService.isDisabled = false
		})
	}

	// Get specialization list
	getSpecializationList() {
		this.generalService.getSpecializations().subscribe(response => {
			this.specializationList = response.body;

		}, err => {
			console.error(this.utilService.getErrorString(err))
		})
	}

	// Return State List.
	getStateList(countryid: any): void {
		this.facultyForm.get('state')?.reset()
		this.facultyForm.get('state')?.disable()
		this.stateList = [] as IState[]
		this.spinnerService.isDisabled = true
		this.generalService.getStatesByCountryID(countryid).subscribe((respond: any[]) => {
			this.stateList = respond;
			if (this.stateList.length > 0 && !this.isViewClicked) {
				this.facultyForm.get('state')?.enable()
			}
		}, (err) => {
			console.error(this.utilService.getErrorString(err))
		}).add(() => {
			this.spinnerService.isDisabled = false
		})
	}

	// Get Country List
	getCountryList(): void {
		this.generalService.getCountries().subscribe((respond: any[]) => {
			this.countryList = respond;
		}, (err) => {
			console.error(this.utilService.getErrorString(err))
		})
	}

	// Get Qualification List.
	getQualificationList(): void {
		let queryParams: any = {
			limit: -1,
			offset: 0
		}
		this.degreeService.getAllDegrees(queryParams).subscribe((respond: any) => {
			this.qualificationList = respond.body;
		}, (err) => {
			console.error(this.utilService.getErrorString(err))
		})
	}

	// Get Designation List.
	getDesignationList(): void {
		this.generalService.getDesignations().subscribe((respond: any[]) => {
			// respond.push(this.other);
			this.designationList = respond;
		}, (err) => {
			console.error(this.utilService.getErrorString(err))
		})
	}

	getSpecializationListForDegree(academicForm: FormGroup): void {
		let degreeID = academicForm.get('degree')?.value?.id
		academicForm.get('specializationList')?.reset()
		academicForm.get('specialization')?.reset()

		console.log(academicForm.get('degree')?.value);

		if (degreeID) {
			this.spinnerService.isDisabled = true
			this.generalService.getSpecializationByDegreeID(degreeID).subscribe((response) => {
				academicForm.get('specializationList').setValue(response.body)
				// this.specializationList = this.specializationList.concat(response.body)
				// this.specializationMap.set(degreeID, response.body)
			}, err => {
				console.error(this.utilService.getErrorString(err))
			}).add(() => {
				this.spinnerService.isDisabled = false
			})
		}
	}

	// Get college branch list.
	getCollegeBranchList(event?: any): void {
		let queryParams: any = {}
		if (event && event?.term != "") {
			queryParams.branchName = event.term
		}
		this.isCollegeLoading = true
		this.generalService.getCollegeBranchListWithLimit(this.collegeLimit, this.collegeOffset, queryParams).subscribe((response) => {
			// console.log("getCollegeBranchList -> ", response);
			this.collegeBranchList = []
			this.collegeBranchList = this.collegeBranchList.concat(response)
		}, (err) => {
			console.error(err)
		}).add(() => {
			this.isCollegeLoading = false
		})
	}

	// Will be called on [addTag]= "addCollegeToList" args will be passed automatically.
	addCollegeToList(option: any): Promise<any> {
		return new Promise((resolve) => {
			resolve(option)
		})
	}

	// // On Qualification Change
	// qualificationChange(index: number):void {
	//       if (this.compareFn(this.facultyAcademicControlsArray.value[index].degree, this.other)) {
	//             if (this.otherQualification) {
	//                   this.otherQualification[index] = true
	//                   return;
	//             }
	//       }
	//       this.otherQualification[index] = false;
	// }


	// // On Designation Change
	// designationChange(index: number):void {
	//       if (this.compareFn(this.facultyExperienceControlsArray.value[index].designation, this.other)) {
	//             this.otherDesignation[index] = true
	//             return
	//       }
	//       this.otherDesignation[index] = false
	// }

	// Set Degree By Other Selection
	setDegree(param: any, index: number): void {
		this.facultyAcademicControlsArray.at(index).get("degree").setValue({ name: param })
	}

	// Set Designation By Other Selection
	setDesignation(param: any, index: number): void {
		this.facultyExperienceControlsArray.at(index).get("designation").setValue({ position: param })
	}

	openModal(modalContent: any, modalSize?: string): void {
		// For resume
		this.isResumeUploadedToServer = false
		this.displayedFileName = "Select file"
		this.docStatus = ""
		if (modalSize == undefined) {
			modalSize = 'lg'
		}
		this.modalRef = this.modalService.open(modalContent,
			{ ariaLabelledBy: 'modal-basic-title', backdrop: 'static', size: modalSize, keyboard: false }
		);
	}

	//Delete resume
	deleteResume() {
		this.fileOpsService.deleteUploadedFile().subscribe((data: any) => {
		}, (error) => {
			console.error(error);
		})
	}

	// Used to dismiss modal.
	dismissFormModal(modal: NgbModalRef) {
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
		modal.dismiss()
		this.isResumeUploadedToServer = false
		this.displayedFileName = "Select file"
		this.docStatus = ""
	}

	//validate add/update form
	validate(): void {
		console.log(this.facultyForm.controls);

		if (this.facultyForm.invalid) {
			this.facultyForm.markAllAsTouched();
		}
		else {
			this.modelAction()
		}
	}

	// ============================================================= FACULTY ASSESSMENT =============================================================

	createFacultyAssessmentForm(): void {
		this.facultyAssessmentForm = this.formBuilder.group({
			feedbacks: this.formBuilder.array([])
		})
	}

	get facultyAssessmentArray() {
		return this.facultyAssessmentForm.get('feedbacks') as FormArray
	}

	addFacultyQuestion(): void {
		this.facultyAssessmentArray.push(this.formBuilder.group({
			id: new FormControl(null),
			credentialID: new FormControl(null),
			facultyID: new FormControl(null, [Validators.required]),
			questionID: new FormControl(null, [Validators.required]),
			optionID: new FormControl(null, [Validators.required]),
			answer: new FormControl(null, [Validators.required]),
			groupID: new FormControl(null, [Validators.required])
			// feedbackQuestionGroup: new FormControl(null, [Validators.required])
		}))
	}

	onViewFacultyAssessmentClick(faculty: any) {
		this.selectedFaculty = faculty
		this.getFacultyAssessment()
		this.viewFacultyAssessment = true
		this.openModal(this.facultyAssessmentModal, "lg")
	}

	onAddFacultyAssessmentClick(): void {
		this.createFacultyAssessmentForm()
		for (let index = 0; index < this.facultyAssementQuestions.length; index++) {

			this.addFacultyQuestion()
			this.facultyAssessmentArray.at(index).get("facultyID").setValue(this.selectedFaculty.id)
			this.facultyAssessmentArray.at(index).get("credentialID").setValue(this.loginID)
			this.facultyAssessmentArray.at(index).get("questionID").setValue(this.facultyAssementQuestions[index].id)
			this.facultyAssessmentArray.at(index).get("groupID").setValue(this.facultyAssementQuestions[index].feedbackQuestionGroup.id)
			// this.facultyAssessmentArray.at(index).get("feedbackQuestionGroup").
			//       setValue(this.facultyAssementQuestions[index].feedbackQuestionGroup)

			if (this.facultyAssementQuestions[index].hasOptions) {
				this.facultyAssessmentArray.at(index).get('optionID').setValidators([Validators.required])
				this.facultyAssessmentArray.at(index).get('answer').setValidators(null)
			} else {
				this.facultyAssessmentArray.at(index).get('optionID').setValidators(null)
				this.facultyAssessmentArray.at(index).get('answer').setValidators([Validators.required])
			}
			this.facultyAssessmentArray.at(index).get('optionID').updateValueAndValidity()
			this.facultyAssessmentArray.at(index).get('answer').updateValueAndValidity()
		}
		this.viewFacultyAssessment = false
	}


	getFeedbackQuestionsForFaculty(): void {

		this.facultyAssementQuestions = []
		this.generalService.getFeedbackQuestionByType("Faculty_Assessment").subscribe((response: any) => {
			this.facultyAssementQuestions = response.body

			// console.log(this.facultyAssementQuestions)
		}, (err: any) => {
			console.error(err)
		})
	}

	getFacultyAssessment(): void {
		this.spinnerService.loadingMessage = "Getting faculty assessment"

		this.facultyAssessment = []
		this.isFacultyAssessment = true
		this.facultyService.getFacultyAssessment(this.selectedFaculty.id).subscribe((response: any) => {
			this.facultyAssessment = response.body
			this.calculateFacultyAssessmentScore()
			// if (this.facultyAssessment.length == 0) {
			//       this.isFacultyAssessment = false
			// }
			console.log(this.facultyAssessment)
		}, (err: any) => {
			// this.isFacultyAssessment = false
			console.error(err)
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err?.error?.error)
		}).add(() => {

			if (this.facultyAssessment.length == 0) {
				this.isFacultyAssessment = false
			}
		})
	}

	addFacultyAssessment(): void {
		this.spinnerService.loadingMessage = "Adding faculty assessment"

		this.facultyService.addFacultyAssessment(this.facultyAssessmentForm.value.feedbacks).subscribe((response: any) => {
			console.log(response)

			this.viewFacultyAssessment = true
			this.getFaculties()
			this.getFacultyAssessment()
		}, (err: any) => {
			console.error(err)

			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err?.error?.error)
		})
	}

	validateFacultyAssessment(): void {
		// console.log(this.facultyAssessmentForm.controls)

		if (this.facultyAssessmentForm.invalid) {
			this.facultyAssessmentForm.markAllAsTouched()
			return
		}
		this.addFacultyAssessment()
	}

	calculateFacultyAssessmentScore(): void {
		let averageScore = 0
		let totalScore = 0

		if (this.facultyAssessment.length == 0) {
			this.averageAssessmentScore = 0
		}

		for (let index = 0; index < this.facultyAssessment?.length; index++) {
			if (this.facultyAssessment[index].option && this.facultyAssessment[index].question) {
				averageScore += this.facultyAssessment[index].option?.key
				totalScore += this.facultyAssessment[index].question.maxScore
			}
		}

		averageScore = (averageScore * 10) / totalScore
		this.averageAssessmentScore = averageScore
	}

	// ============================================================= COURSE TECHNICAL ASSESSMENT =============================================================

	createTechnicalAssessmentForm(): void {
		this.technicalAssessmentForm = this.formBuilder.group({
			assessments: new FormArray([], Validators.required)
		})
	}

	get technicalAssessmentArray() {
		return this.technicalAssessmentForm.get("assessments") as FormArray
	}

	addTechnicalAssessmentToForm(): void {
		this.technicalAssessmentArray.push(this.formBuilder.group({
			id: new FormControl(null),
			courseID: new FormControl(null, [Validators.required]),
			course: new FormControl(null),
			faculty: new FormControl(null),
			rating: new FormControl(null, [Validators.min(1), Validators.max(10)])
		}))
	}

	deleteTechnicalAssessmentFromArray(index: number): void {
		this.technicalAssessmentArray.removeAt(index)
	}

	onTechnicalAssessmentClick(faculty: any): void {
		this.selectedFaculty = faculty
		this.viewTechnicalAssessment = true
		this.getCourses()
		this.getTechnicalAssessmentForFaculty()
		// this.createTechnicalAssessmentForm()
		this.openModal(this.technicalAssessmentModal, "lg")
	}

	onAddTechnicalAssessmentClick(): void {
		this.isUpdateAssessmentClick = false
		this.viewTechnicalAssessment = false
		// this.createTechnicalAssessmentForm()
	}

	onUpdateTechnicalAssessmentClick(technicalAssessment: ICourseTechnicalAssessment): void {
		this.isUpdateAssessmentClick = true
		this.viewTechnicalAssessment = false
		this.createTechnicalAssessmentForm()
		this.addTechnicalAssessmentToForm()
		technicalAssessment.courseID = technicalAssessment.course.id
		let tempTechnicalAssessmentArray = [technicalAssessment]
		this.technicalAssessmentForm.get("assessments").patchValue(tempTechnicalAssessmentArray)
		console.log(this.technicalAssessmentForm.value);
	}

	addTechnicalAssessments(): void {
		this.createTechnicalAssessmentForm()
		for (let index = 0; index < this.courseList.length; index++) {
			this.addTechnicalAssessmentToForm()
			this.technicalAssessmentArray.at(index).get("course").setValue(this.courseList[index])
			this.technicalAssessmentArray.at(index).get("courseID").setValue(this.courseList[index].id)
			this.technicalAssessmentArray.at(index).get("faculty").setValue(this.selectedFaculty)
		}
	}

	getCourses(): void {

		this.courseList = []
		this.generalService.getCourseList().subscribe((response: any) => {
			this.courseList = response.body
			this.addTechnicalAssessments()

		}, (error: any) => {
			console.error(error)

		})
	}

	getTechnicalAssessmentForFaculty(): void {
		this.spinnerService.loadingMessage = "Getting course technical assessment"

		this.technicalAssessments = []
		this.isTechnicalAssessment = true
		this.courseService.getCourseTechnicalAssessmentForFaculty(this.selectedFaculty.id).subscribe((response: any) => {
			this.technicalAssessments = response.body
			if (this.technicalAssessments.length == 0) {
				this.isTechnicalAssessment = false
			}
			this.viewTechnicalAssessment = true

		}, (err: any) => {
			console.error(err)
			this.isTechnicalAssessment = false

			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err?.error?.error)
		})
	}

	addTechnicalAssessment(): void {
		this.spinnerService.loadingMessage = "Adding course technical assessment"

		this.courseService.addCourseTechnicalAssessmnets(this.technicalAssessmentForm.value.assessments,
			this.selectedFaculty.id).subscribe((response: any) => {
				console.log(response);
				this.getFaculties()
				this.getTechnicalAssessmentForFaculty()

			}, (err: any) => {
				console.error(err)

				if (err.statusText.includes('Unknown')) {
					alert("No connection to server. Check internet.")
					return
				}
				alert(err?.error?.error)
			})
	}

	updateTechnicalAssessment(): void {
		this.spinnerService.loadingMessage = "Updating course technical assessment"

		this.courseService.updateCourseTechnicalAssessmnet(this.technicalAssessmentForm.value.assessments[0],
			this.selectedFaculty.id).subscribe((response: any) => {
				console.log(response);
				this.getTechnicalAssessmentForFaculty()

			}, (err: any) => {
				console.error(err)

				if (err.statusText.includes('Unknown')) {
					alert("No connection to server. Check internet.")
					return
				}
				alert(err?.error?.error)
			})
	}

	deleteTechnicalAssessment(assessmentID: string): void {
		if (confirm("Are you sure you want to delete technical assessment")) {
			this.spinnerService.loadingMessage = "Updating course technical assessment"

			this.courseService.deleteCourseTechnicalAssessmnet(assessmentID).
				subscribe((response: any) => {
					console.log(response)
					this.getTechnicalAssessmentForFaculty()

				}, (err: any) => {
					console.error(err)

					if (err.statusText.includes('Unknown')) {
						alert("No connection to server. Check internet.")
						return
					}
					alert(err?.error?.error)
				})
		}
	}

	validateTechnicalAssessment(): void {
		// console.log(this.technicalAssessmentForm.controls)

		if (this.technicalAssessmentForm.invalid) {
			this.technicalAssessmentForm.markAllAsTouched()
			return
		}

		if (this.isUpdateAssessmentClick) {
			this.updateTechnicalAssessment()
			return
		}

		this.addTechnicalAssessment()
	}
}
