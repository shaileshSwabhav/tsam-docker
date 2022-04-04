import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, Validators, FormControl } from '@angular/forms';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { CourseService, ICourse } from 'src/app/service/course/course.service';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { NgbModal, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { IPermission } from 'src/app/service/menu/menu.service';
import { Role, UrlConstant } from 'src/app/service/constant';
import { ActivatedRoute, Router } from '@angular/router';
import { ITechnology, TechnologyService } from 'src/app/service/technology/technology.service';
import { Observable, Subject } from 'rxjs';
import { LocalService } from 'src/app/service/storage/local.service';
import { SafeUrl } from '@angular/platform-browser';


@Component({
	selector: 'app-course-master',
	templateUrl: './course-master.component.html',
	styleUrls: ['./course-master.component.css']
})
export class CourseMasterComponent implements OnInit {

	// modal
	modalHeader: string;
	modalButton: string;
	isCourseUpdate: boolean;
	modalHandler: () => void;
	resourceHandler: () => void;
	formHandler: (index?: number) => void;
	modalRef: any;

	isViewClicked: boolean
	isCoursesLoaded: boolean

	// course
	courseForm: FormGroup;
	doesEligibilityExists: boolean
	courseSearchForm: FormGroup;
	selectedCourse: any;
	courseID: string

	// resource
	resourceButtonLabel: string;
	resourseType: string;

	// component
	technologies: ITechnology[];
	totalCourse: number;
	courses: any[];
	academicYears: any[];
	courseType: any[];
	courseLevel: any[];
	courseStatus: any[];
	ratinglist: any[];

	// search
	showSearch: boolean
	isSearched: boolean
	searchFormValue: any

	// access
	permission: IPermission
	isFaculty: boolean

	//spinner


	// pagination
	limit: number;
	offset: number;
	currentPage: number;
	paginationString: string

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

	isTechLoading: boolean
	technologies$: Observable<any>;
	technologyInput$: Subject<string>;
	minLengthTerm: number;
	techLimit: number
	techOffset: number

	private readonly IGNORE_SEARCH_FIELD: string[] = ["limit", "offset"]

	@ViewChild('courseSearchModal') courseSearchModal: any
	@ViewChild('courseModal') courseModal: any
	@ViewChild('drawer') drawer: any
	@ViewChild('deleteConfirmationModal') deleteConfirmationModal: any

	// cke editor.
	@ViewChild("ckeEditorPreRequisites") ckeditorPreRequisites: any

	// Cke editor configuration.
	ckeEditorConfig: any

	// Params.
	operation: string

	constructor(
		private formBuilder: FormBuilder,
		public utilService: UtilityService,
		private generalService: GeneralService,
		private courseService: CourseService,
		private techService: TechnologyService,
		private localService: LocalService,
		private spinnerService: SpinnerService,
		private urlConstant: UrlConstant,
		private role: Role,
		private fileOperationService: FileOperationService,
		private router: Router,
		private route: ActivatedRoute,
		private modalService: NgbModal
	) {
		this.initializeVariables()
		this.createForms()
		this.getAllComponents()
	}

	isAdmin: boolean = false

	initializeVariables(): void {
		this.isAdmin = (this.localService.getJsonValue("roleName") == this.role.ADMIN ? true : false)
		this.isFaculty = (this.localService.getJsonValue("roleName") == this.role.FACULTY ? true : false)
		if (this.isFaculty){
			this.permission = this.utilService.getPermission(this.urlConstant.BANK_COURSE)
		}
		if (!this.isFaculty){
			this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_COURSE_MASTER)
		}

		this.limit = 5
		this.offset = 0
		this.currentPage = 1

		this.isViewClicked = false
		this.doesEligibilityExists = false
		this.showSearch = false
		this.isBrochureUploadedToServer = false
		this.isCoursesLoaded = true
		this.isSearched = false

		this.spinnerService.loadingMessage = "Getting courses"
		this.displayedFileName = "Select file"

		this.isTechLoading = false
		this.technologyInput$ = new Subject<string>();
		this.minLengthTerm = 1
		this.techLimit = 10
		this.techOffset = 0

		// Cke editor configuration.
		this.ckeEditorCongiguration()
	}


	// Load Essential Data.
	getAllComponents(): void {
		this.getTechnology()
		// this.getTechAsync()
		this.getCourseType()
		this.getCourseLevel()
		this.getCourseStatus()
		this.getAcademicYear()
		this.getStudentRatingList()
		this.searchOrGetCourses()
	}

	createForms(): void {
		this.createCourseSearchForm()
		this.createCourseForm()
	}


	get ongoingOperations() {
		return this.spinnerService.ongoingOperations
	}

	ngOnInit() {
	}

	ngAfterViewInit() {
		console.log("=== ngAfterViewInit ===");
		
		if (this.operation && this.operation == "add"){
		  this.onCourseAddButtonClick()
		  this.operation = null
		  return
		}
	}

	createCourseSearchForm(): void {
		this.courseSearchForm = this.formBuilder.group({
			createdAt: new FormControl(null, [Validators.pattern(/^((19[8-9]\d|20[0-3]\d)\-\d{1,2}\-\d{1,2})/)]),
			technologies: new FormControl(null),
			courseName: new FormControl(null),
			courseType: new FormControl(null),
			courseLevel: new FormControl(null),
			technology: new FormControl(null),
			limit: new FormControl(this.limit),
			offset: new FormControl(this.offset),
			// currentPage: new FormControl(this.currentPage),
		})
	}

	// Return Current Date
	getDate(): any {
		let d = new Date();
		return d.getDate() + "/" + d.getMonth() + "/" + d.getFullYear();
	}

	// Create Cource Form Object.
	createCourseForm(): void {
		this.courseForm = this.formBuilder.group({
			id: new FormControl(),
			code: new FormControl(null),
			name: new FormControl(null, [Validators.required]),
			technologies: [Array(), Validators.required],
			courseType: new FormControl(null, [Validators.required]),
			courseLevel: new FormControl(null, [Validators.required]),
			description: new FormControl(null, [Validators.maxLength(2000)]),
			preRequisites: new FormControl(null, [Validators.maxLength(2000)]),
			eligibility: this.formBuilder.group({
				id: new FormControl(),
				technologies: [null],
				studentRating: new FormControl(null),
				experience: new FormControl(null),
				academicYear: new FormControl(null)
			}),
			price: new FormControl(null, [Validators.required]),
			durationInMonths: new FormControl(null, [Validators.required]),
			totalHours: new FormControl(null, [Validators.required]),
			totalSessions: new FormControl(null),
			brochure: new FormControl(null),
			logo: new FormControl(null),
			// sessions: this.formBuilder.array([]),
			// discount: new FormControl(null),
		});

		if (this.isFaculty) {
			this.courseForm.get("price").clearValidators()
			this.courseForm.updateValueAndValidity()
		}
	}

	eligilibilityChecked(event) {
		event.target.checked = true
		if (this.doesEligibilityExists) {
			if (this.courseForm.get('eligibility').valid) {
				if (confirm("This action will delete the batch eligibility. Are you sure you want to delete batch eligibility?")) {
					this.courseForm.get('eligibility.id').setValue(null)
					event.target.checked = false
				} else {
					return
				}
			}
			this.courseForm.get('eligibility.id').setValue(null)
			this.doesEligibilityExists = false
			this.setEligibility()
			return;
		}
		this.doesEligibilityExists = true;
		this.setEligibility()
	}

	setEligibility() {
		this.technologyCheck()
		this.studentRatingCheck()
		this.experienceCheck()
		this.academicYearCheck()
		this.courseForm.markAsDirty();
		// console.log(this.courseForm.get('eligibility'));
	}

	technologyCheck() {
		const technologyControl = this.courseForm.get('eligibility.technologies')
		if (this.doesEligibilityExists) {
			technologyControl.setValidators([Validators.required])
		} else {
			technologyControl.setValue(null)
			technologyControl.setValidators(null)
			technologyControl.markAsUntouched()
		}
		technologyControl.updateValueAndValidity()
	}

	studentRatingCheck() {
		const studentRatingControl = this.courseForm.get('eligibility.studentRating')
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
		const experienceControl = this.courseForm.get('eligibility.experience')
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
		const academicYearControl = this.courseForm.get('eligibility.academicYear')
		if (this.doesEligibilityExists) {
			academicYearControl.setValidators([Validators.required])
		} else {
			academicYearControl.setValue(null)
			academicYearControl.setValidators(null)
			academicYearControl.markAsUntouched()
		}
		academicYearControl.updateValueAndValidity()
	}


	updateCourseForm(course: any): void {
		this.createCourseForm();
		if (course.eligibility != null) {
			this.doesEligibilityExists = true
		} else {
			course.eligibility = {}
			this.doesEligibilityExists = false
		}
		this.courseForm.patchValue(course);

		// make eligibility null if eligibility does not exist
		if (Object.keys(course.eligibility).length == 0) {
			course.eligibility = null
		}
	}

	setPaginationString() {
		this.paginationString = ""
		let limit = this.courseSearchForm.get("limit").value
		let offset = this.courseSearchForm.get("offset").value

		let start: number = limit * offset + 1
		let end: number = +limit + limit * offset
		if (this.totalCourse < end) {
			end = this.totalCourse
		}
		if (this.totalCourse == 0) {
			this.paginationString = ""
			return
		}
		this.paginationString = `${start} - ${end}`
	}

	deleteEligibilityIfEmpty(course: ICourse): void {
		if (Object.keys(course.eligibility).length === 0 && course.eligibility.constructor === Object) {
			delete course.eligibility
		}
	}


	// On Course Add Button Click
	onCourseAddButtonClick(): void {
		this.selectedCourse = null
		this.isViewClicked = false
		this.doesEligibilityExists = false
		this.isFileUploading = false
		this.isBrochureUploadedToServer = false
		this.updateCommonVariable(this.createCourseForm, this.addCourse, "Add New Course", "Add Course");
		this.createCourseForm();
		this.openModal(this.courseModal)
	}

	onCourseViewButtonClick(event: any, course: ICourse) {
		event.stopPropagation()

		// console.log(course);
		this.isViewClicked = true
		this.updateCommonVariable(this.createCourseForm, this.addCourse, "View Course", "View Course");
		this.selectedCourse = course
		this.updateCourseForm(course)
		this.courseForm.disable()
		this.displayedFileName = `<a href=${course.brochure} target="_blank">Select file</a>`
		this.openModal(this.courseModal)
	}

	// On Update Course Button Click.
	onCourseUpdateButtonClick(): void {
		this.isViewClicked = false
		this.updateCommonVariable(this.updateCourseForm, this.updateCourse, "Update Course", "Update Course");
		if (this.selectedCourse.brochure) {
			this.displayedFileName = `<a class='custom-file-label' href=${this.selectedCourse.brochure} target="_blank">Select file</a>`
		}
		if (this.selectedCourse.logo) {
			this.logoDisplayedFileName = `<a class='custom-file-label' href=${this.selectedCourse.logo} target="_blank">Select file</a>`
		}
		this.courseForm.enable()
	}

	onCourseDeleteClick(event: any, courseID: string): void {
		event.stopPropagation()
		this.courseID = courseID
		this.openModal(this.deleteConfirmationModal, "md")
	}

	// Compare two Object.
	compareFn(ob1: any, ob2: any): boolean {
		if (ob1 == null && ob2 == null) {
			return true
		}
		if (ob1 != undefined && ob2 != undefined) {
			if (ob1.id == ob2.id) {
				return true;
			}
			return false;
		}
		return false
	}

	// Update Common Variable On Update, Add Button Click.
	updateCommonVariable(formhandler: any, actionhandler: any, headername: string, buttonname: string): void {
		this.formHandler = formhandler;
		this.modalHandler = actionhandler;
		this.modalHeader = headername;
		this.modalButton = buttonname;
	}

	// redirectToCourseSession(course: ICourse): void {
	// 	this.router.navigate(['session'], {
	// 		relativeTo: this.route,
	// 		queryParams: {
	// 			"courseName": course.name,
	// 			"courseID": course.id
	// 		}
	// 	}).catch(err => {
	// 		console.error(err)

	// 	});
	// }

	// redirectToProgrammingAssignment(course: ICourse): void {
	// 	this.router.navigate(['programming-assignment'], {
	// 		relativeTo: this.route,
	// 		queryParams: {
	// 			"courseName": course.name,
	// 			"courseID": course.id
	// 		}
	// 	}).catch(err => {
	// 		console.error(err)

	// 	});
	// }

	// redirectToCourseModule(course: ICourse): void {
	// 	this.router.navigate(['module'], {
	// 		relativeTo: this.route,
	// 		queryParams: {
	// 			"courseName": course.name,
	// 			"courseID": course.id
	// 		}
	// 	}).catch(err => {
	// 		console.error(err)

	// 	});
	// }

	// Handles pagination
	changePage(pageNumber: number): void {
		this.courseSearchForm.get("offset").setValue(pageNumber - 1)

		this.limit = this.courseSearchForm.get("limit").value
		this.offset = this.courseSearchForm.get("offset").value
		this.searchCourse()
	}

	searchOrGetCourses(): void {
		console.log("=== searchOrGetCourses ===");

		let queryParams = this.route.snapshot.queryParams
		if (queryParams.operation){
			this.operation = queryParams.operation
			// queryParams = {}
		}
		if (this.utilService.isObjectEmpty(queryParams)) {
			this.changePage(1)
			return
		}
		this.courseSearchForm.patchValue(queryParams)
		// let limit = this.route.snapshot.queryParamMap.get("limit")
		// let offset = this.route.snapshot.queryParamMap.get("offset")
		// if (limit && offset) {
		// 	this.changePage(this.offset - 1)
		// 	return
		// }
		this.searchCourse()
	}

	searchAndCloseDrawer(): void {
		this.drawer.toggle()
		this.searchCourse()
	}

	searchCourse(): void {
		this.searchFormValue = { ...this.courseSearchForm?.value }
		this.router.navigate([], {
			relativeTo: this.route,
			queryParams: this.searchFormValue,
		})
		let flag: boolean = true

		for (let field in this.searchFormValue) {
			if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
				delete this.searchFormValue[field];
			} else {
				if (!this.IGNORE_SEARCH_FIELD.includes(field)) {
					this.isSearched = true
				}
				flag = false
			}
		}

		// No API call on empty search.
		if (flag) {
			return
		}

		this.getAllCourse()
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
			if (this.isBrochureUploadedToServer) {
				this.deleteUploadedFile(this.uploadedBrochureURL)
			}
			if (this.isLogoUploadedToServer) {
				this.deleteUploadedFile(this.uploadedLogoURL)
			}
		}
		modal.dismiss()
		this.resetBrochureUploadFields()
		this.resetBrochureUploadFields()
		this.spinnerService.isDisabled = false
		this.isTechLoading = false
	}

	deleteUploadedFile(fileURL: string) {
		this.spinnerService.loadingMessage = "Deleting..."
		this.fileOperationService.deleteUploadedFile(fileURL).subscribe((data: any) => {
			console.log(data)
		}, (err) => {
			console.error(err)
		})
	}

	resetBrochureUploadFields(): void {
		this.uploadedBrochureURL = ""
		this.isBrochureUploadedToServer = false
		this.isFileUploading = false
		this.displayedFileName = "Select file"
		this.docStatus = ""
	}

	resetLogoUploadFields(): void {
		this.uploadedLogoURL = ""
		this.isLogoUploadedToServer = false
		this.isLogoFileUploading = false
		this.logoDisplayedFileName = "Select file"
		this.logoDocStatus = ""
	}


	// open modal for add/update course form
	openModal(courseModal: any, modalSize?: string): void {

		this.resetBrochureUploadFields()
		this.resetLogoUploadFields()

		if (modalSize == undefined) {
			modalSize = "xl"
		}

		this.modalRef = this.modalService.open(courseModal,
			{
				ariaLabelledBy: 'modal-basic-title',
				backdrop: 'static', size: modalSize,
				keyboard: false
			});
		/*this.modalRef.result.subscribe((result) => {
		}, (reason) => {

		});*/
	}

	// Get all invalid controls in talent form.
	public findInvalidControls(): any[] {
		const invalid = []
		const controls = this.courseForm.controls
		for (const name in controls) {
			if (controls[name].invalid) {
				invalid.push(name)
			}
		}
		return invalid
	}

	//validate add/update course form
	validate(): void {
		// console.log(this.findInvalidControls())
		if (this.isFileUploading || this.isLogoFileUploading) {
			alert("Please wait till file is being uploaded")
			return
		}

		if (this.courseForm.invalid) {
			this.courseForm.markAllAsTouched()
			return
		}

		this.modalHandler();
	}

	onAddFilterClick(): void {
		this.openModal(this.courseSearchModal)
	}

	resetSearchCourseForm(): void {
		this.limit = this.courseSearchForm.get("limit").value
		this.offset = this.courseSearchForm.get("offset").value
		// this.currentPage = this.courseSearchForm.get("currentPage").value

		this.courseSearchForm.reset({
			limit: this.limit,
			offset: this.offset,
			// currentPage: new FormControl(this.currentPage),
		})

		console.log(this.courseSearchForm.controls);
	}

	resetSearchAndGetAll(): void {
		// console.log("resetSearchAndGetAll")
		this.resetSearchCourseForm()
		this.isSearched = false
		this.showSearch = false
		// this.router.navigate([this.urlConstant.COURSE_MASTER])
		this.searchCourse()
	}

	// =============================================================CRUD STARTS=============================================================


	// Get All Course List
	getAllCourse(): void {
		if (this.isSearched) {
			this.getSearchedCourse()
			return
		}
		this.spinnerService.loadingMessage = "Getting courses"
			;
		this.courses = [];
		this.totalCourse = 0
		this.isCoursesLoaded = true
		console.log(this.searchFormValue);

		this.courseService.getCourses(this.searchFormValue).subscribe((res: any) => {
			this.courses = res.body
			this.totalCourse = parseInt(res.headers.get("X-Total-Count"))
			this.setPaginationString()
			if (this.totalCourse == 0) {
				this.isCoursesLoaded = false
			}
			// console.log(this.courses);
			;
		}, err => {
			console.error(err)

			this.totalCourse = 0
			this.setPaginationString()
			this.isCoursesLoaded = false
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err.error.error)
		})
	}

	// Get Searched Course List
	getSearchedCourse(): void {
		this.spinnerService.loadingMessage = "Searching courses";
		;
		this.totalCourse = 0
		this.courses = [];
		this.isCoursesLoaded = true
		// console.log(this.searchFormValue)
		this.courseService.getCourses(this.searchFormValue).subscribe((res: any) => {
			this.courses = res.body
			this.totalCourse = parseInt(res.headers.get("X-Total-Count"))
			this.setPaginationString()
			if (this.totalCourse == 0) {
				this.isCoursesLoaded = false
			}
			;
		}, err => {
			console.error(err)

			this.totalCourse = 0
			this.isCoursesLoaded = false
			this.setPaginationString()
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err.error.error)
		})
	}

	// Add Course.
	addCourse(): void {
		this.spinnerService.loadingMessage = "Adding course";
		;
		let data = this.courseForm.value;
		this.utilService.deleteNullValueIDFromObject(data)
		this.utilService.deleteNullValuePropertyFromObject(data.eligibility)
		if (data.eligibility) {
			this.deleteEligibilityIfEmpty(data)
		}
		console.log(data);

		this.courseService.addCourse(data).subscribe(respond => {
			this.modalRef.close()

			this.getAllCourse()
			alert("Course successfully added")
		}, (err) => {
			console.error(err)

			this.totalCourse = 0
			this.setPaginationString()
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err?.error?.error)
		})
		// console.log(this.courseForm);
	}

	// Update Course
	updateCourse(): void {
		this.spinnerService.loadingMessage = "Updating course";
		;
		let data = this.courseForm.value;
		this.utilService.deleteNullValueIDFromObject(data)
		this.utilService.deleteNullValuePropertyFromObject(data.eligibility)
		if (data.eligibility) {
			this.deleteEligibilityIfEmpty(data)
		}
		// console.log(data);
		this.courseService.updateCourse(data).subscribe(respond => {
			this.modalRef.close()

			this.getAllCourse()
			alert("Course successfully updated")
		}, (err) => {
			console.error(err)

			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err?.error?.error)
		});
		// console.log(this.courseForm.value);
	}

	// Delete Course
	deleteCourse(): void {
		this.spinnerService.loadingMessage = "Deleting course";
		this.modalRef.close();
		;
		this.courseService.deleteCourse(this.courseID).subscribe(respond => {
			alert("Course successfully deleted")

			this.getAllCourse();
		}, (err) => {
			console.error(err)

			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err?.error?.error)
		});
	}

	uploadedBrochureURL: string
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
			// Upload brochure if it is present.]
			this.isFileUploading = true
			this.spinnerService.isDisabled = true
			this.fileOperationService.uploadBrochure(file,
				this.fileOperationService.COURSE_FOLDER + this.fileOperationService.BROCHURE_FOLDER)
				.subscribe((data: any) => {
					this.courseForm.markAsDirty()
					this.courseForm.patchValue({
						brochure: data
					})
					this.displayedFileName = file.name
					this.isFileUploading = false
					this.uploadedBrochureURL = data
					// isBrochureUploadedToServer can be replaced by url link entirely #nj
					this.isBrochureUploadedToServer = true
					this.docStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
				}, (error) => {
					this.isFileUploading = false
					this.docStatus = `<p><span>&#10060;</span> ${error}</p>`
				}).add(() => this.spinnerService.isDisabled = false)
		}
	}


	uploadedLogoURL: string
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
			// Upload brochure if it is present.
			this.isLogoFileUploading = true
			this.spinnerService.loadingMessage = "Uploading logo..."
			this.fileOperationService.uploadLogo(file,
				this.fileOperationService.COURSE_FOLDER + this.fileOperationService.LOGO_FOLDER).subscribe((data: any) => {
					this.courseForm.markAsDirty()
					this.courseForm.patchValue({
						logo: data
					})
					this.logoDisplayedFileName = file.name
					this.isLogoFileUploading = false
					this.uploadedLogoURL = data
					// isLogoUploadedToServer can be replaced by url link entirely #nj
					this.isLogoUploadedToServer = true
					this.logoDocStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
				}, (error) => {
					this.isLogoFileUploading = false
					this.logoDocStatus = `<p><span>&#10060;</span> ${error}</p>`
				})
		}
	}

	// Redirect to course details page.
	redirectToCourseDetails(courseID: string): void {
		console.log(courseID)

		this.router.navigate([this.isAdmin ? 
				this.urlConstant.TRAINING_COURSE_MASTER_DETAILS : this.urlConstant.BANK_COURSE_DETAILS], {
			queryParams: {
				"courseID": courseID,
			}
		}).catch(err => {
			console.error(err)
		})
	}

	// Redirect to external link.
	redirectToExternalLink(url): void {
		window.open(url, "_blank")
	}

	// cke editor congiguration.
	ckeEditorCongiguration(): void {
		this.ckeEditorConfig = {
			extraPlugins: 'codeTag',
			removePlugins: "exportpdf",
			// stylesSet: 'new_styles',
			toolbar: [
				{ name: 'styles', items: ['Styles', 'Format'] },
				{
					name: 'basicstyles', groups: ['basicstyles', 'cleanup'],
					items: ['Bold', 'Italic', 'Underline', 'Strike', 'Subscript', 'Superscript', 'RemoveFormat', 'Code']
				},
				{
					name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'],
					items: ['NumberedList', 'BulletedList', '-', 'Outdent', 'Indent', '-', 'Blockquote']
				},
				// { name: 'links', items: [ 'Link', 'Unlink' ] }, //, 'Anchor' // Link
				// { name: 'insert', items: [ 'Image' ] }, //, 'Table', 'HorizontalRule' // Image
				{ name: 'document', groups: ['mode', 'document', 'doctools'], items: ['Source'] },
				// { name: 'clipboard', groups: [ 'clipboard', 'undo' ], items: [ 'Cut', 'Copy', 'Paste', 'PasteText', 'PasteFromWord', '-', 'Undo', 'Redo' ] },
				// { name: 'editing', groups: [ 'find', 'selection', 'spellchecker' ], items: [ 'Scayt' ] },
				// { name: 'tools', items: [ 'Maximize' ] },
				// { name: 'others', items: [ '-' ] },
				// { name: 'about', items: [ 'About' ] }
			],
			toolbarGroups: [
				{ name: 'styles' },
				{ name: 'basicstyles', groups: ['basicstyles', 'cleanup'] },
				{ name: 'document', groups: ['mode', 'document', 'doctools'] },
				{ name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'] },
				// { name: 'links' }, // Link
				// { name: 'insert' }, // Image
				// '/',
				// { name: 'colors' },
				// { name: 'clipboard', groups: [ 'clipboard', 'undo' ] },
				// { name: 'editing', groups: [ 'find', 'selection', 'spellchecker' ] },
				// { name: 'forms' },
				// { name: 'tools' },
				// { name: 'others' },
			],
			removeButtons: "",
			language: 'en',
			resize_enabled: false,
			width: "100%", height: "150px",
			forcePasteAsPlainText: false,
		}
	}

	// =============================================================COMPONENTS=============================================================


	// Get All Student Rating
	getStudentRatingList(): void {
		this.generalService.getGeneralTypeByType("talent_rating").subscribe((respond: any[]) => {
			this.ratinglist = respond;
		}, (err) => {
			console.error(err)
		})
	}

	// Get Course Type List 
	getCourseType(): void {
		this.generalService.getGeneralTypeByType("course_type").subscribe((respond: any) => {
			this.courseType = respond;
		}, (err) => {
			console.error(err)
		})
	}

	// Get Course Level List 
	getCourseLevel(): void {
		this.generalService.getGeneralTypeByType("course_level").subscribe((respond: any) => {
			this.courseLevel = respond;
		}, (err) => {
			console.error(err)
		})
	}

	// Get Academic Year List.
	getAcademicYear(): void {
		this.generalService.getGeneralTypeByType("academic_year").subscribe((respond: any[]) => {
			this.academicYears = respond;
		}, (err) => {
			console.error(err)
		})
	}

	// Get Course Status List.
	getCourseStatus(): void {
		this.generalService.getGeneralTypeByType("course_status").subscribe((respond: any[]) => {
			this.courseStatus = respond;
		}, (err) => {
			console.error(err)
		});
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
		}, (err) => {
			console.error(err)
		}).add(() => {
			this.isTechLoading = false
			this.spinnerService.isDisabled = false
		})
	}

	// getTechAsync(): void {
	//       this.technologies$ = concat(
	//             of([]), // default items
	//             this.technologyInput$.pipe(
	//                   startWith(''),
	//                   distinctUntilChanged(),
	//                   debounceTime(200),
	//                   tap(() => this.isTechLoading = true),
	//                   switchMap(term => {
	//                         let queryParams: any = {}
	//                         if (term != "") {
	//                               queryParams.language = term
	//                         }
	//                         return this.techService.getAllTechnologies(this.techLimit, this.techOffset, queryParams).pipe(
	//                               map(res => {
	//                                     console.log("technologies$ -> ", res);
	//                                     return res.body
	//                               }),
	//                               catchError(() => of([])), // empty list on error
	//                               tap(() => {
	//                                     this.isTechLoading = false
	//                               })
	//                         )
	//                   })
	//             )
	//       )
	// }

}