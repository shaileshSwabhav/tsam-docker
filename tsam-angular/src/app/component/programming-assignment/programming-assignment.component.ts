import { Component, OnInit, ViewChild } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { Role, UrlConstant } from 'src/app/service/constant';
import { CourseService, ICourseProgrammingAssignment } from 'src/app/service/course/course.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { ResourceService } from 'src/app/service/resource/resource.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { IProgrammingAssignment, ProgrammingAssignmentService } from 'src/app/service/programming-assignment/programming-assignment.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
	selector: 'app-programming-assignment',
	templateUrl: './programming-assignment.component.html',
	styleUrls: ['./programming-assignment.component.css']
})
export class ProgrammingAssignmentComponent implements OnInit {

	// components
	assignmentType: any[]
	sourceList: any[]
	programmingQuestionList: any[]
	complexityLevel: any[]
	resourcesList: any[]

	// programming-assignment
	programmingAssignments: IProgrammingAssignment[]
	totalAssignments: number

	// forms
	assignmentSearchForm: FormGroup
	assignmentForm: FormGroup

	// pagination
	limit: number
	offset: number
	currentPage: number
	paginationString: string

	// search
	showSearch: boolean
	isSearched: boolean
	searchFormValue: any

	// flags
	isViewMode: boolean
	isUpdateMode: boolean
	isAssignmentLoaded: boolean

	// access
	permission: IPermission
	isAdmin: boolean
	isFaculty: boolean

	// modal
	modalRef: any;
	@ViewChild('assignmentFormModal') assignmentFormModal: any
	@ViewChild('deleteModal') deleteModal: any
	@ViewChild('drawer') drawer: any
	@ViewChild('addAssignmentModal') addAssignmentModal: any
	@ViewChild('addCourseAssignmentModal') addCourseAssignmentModal: any

	//spinner



	// ck-editor
	ckConfig: any

	// redirection
	isNavigated: boolean

	// course-assignmnet
	courseAssignmentForm: FormGroup
	courseID: string
	courseName: string
	assignedCourseAssignments: string[]

	// constants
	private readonly URL_REGEX = /^(http(s)?:\/\/)[\w.-]+(\.[\w\.-]+)+[\w\-\._~:?#[\]@!\/$&'\(\)\*\+,;=.]+$/
	private readonly IGNORE_SEARCH_FIELD: string[] = ["limit", "offset", "courseID", "courseName"]
	readonly READING = "Reading"
	readonly CODING = "Coding"

	constructor(
		private formBuilder: FormBuilder,
		public utilService: UtilityService,
		private generalService: GeneralService,
		private courseService: CourseService,
		private assignmentService: ProgrammingAssignmentService,
		private resourceService: ResourceService,
		private spinnerService: SpinnerService,
		private urlConstant: UrlConstant,
		private role: Role,
		private localService: LocalService,
		private router: Router,
		private route: ActivatedRoute,
		private modalService: NgbModal,
	) {
		this.initializeVariables()
		this.createSearchForm()
		this.editorConfig()
		this.getAllComponents()
	}

	editorConfig(): void {
		this.ckConfig = {
			extraPlugins: 'codeTag,kbdTag',
			removePlugins: "exportpdf",
			toolbar: [
				{ name: 'styles', items: ['Styles', 'Format'] },
				{
					name: 'basicstyles', groups: ['basicstyles', 'cleanup'],
					items: ['Bold', 'Italic', 'Underline', 'Strike', 'Subscript', 'Superscript', 'RemoveFormat', 'Code', 'Kbd']
				},
				{
					name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'],
					items: ['NumberedList', 'BulletedList', '-', 'Outdent', 'Indent', '-', 'Blockquote']
				},
				{ name: 'document', groups: ['mode', 'document', 'doctools'], items: ['Source'] },
			],
			toolbarGroups: [
				{ name: 'styles' },
				{ name: 'basicstyles', groups: ['basicstyles', 'cleanup'] },
				{ name: 'document', groups: ['mode', 'document', 'doctools'] },
				{ name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'] },
				{ name: 'links' },
			],
			removeButtons: "",
			language: "en",
			resize_enabled: false,
			width: "100%",
			height: "80%",
			forcePasteAsPlainText: false,
		}
	}

	initializeVariables(): void {
		this.isAdmin = (this.localService.getJsonValue("roleName") == this.role.ADMIN ? true : false)
		this.isFaculty = (this.localService.getJsonValue("roleName") == this.role.FACULTY ? true : false)

		if (this.isAdmin) {
			this.permission = this.utilService.getPermission(this.urlConstant.ADMIN_PROGRAMMING_ASSIGNMENT)
		}

		if (this.isFaculty) {
			this.permission = this.utilService.getPermission(this.urlConstant.FACULTY_PROGRAMMING_ASSIGNMENT)
		}

		this.limit = 5
		this.offset = 0
		this.currentPage = 1

		this.totalAssignments = 0

		this.isViewMode = false
		this.isUpdateMode = false
		this.isNavigated = false
		this.isAssignmentLoaded = true

		this.assignmentType = []
		this.sourceList = []
		this.programmingQuestionList = []
		this.complexityLevel = []
		this.assignedCourseAssignments = []
		this.programmingAssignments = []
		this.resourcesList = []

		this.spinnerService.loadingMessage = "Getting programming assignment"
	}

	getAllComponents(): void {
		this.getProgrammingAssignmentType()
		this.getProgrammingQuestionList()
		this.getProgrammingAssignmentSource()
		this.getComplexityLevel()
		this.getResourcesList()
		this.searchOrGetAssignments()
	}


	get ongoingOperations() {
		return this.spinnerService.ongoingOperations
	}

	ngOnInit(): void {
	}

	createSearchForm(): void {
		this.assignmentSearchForm = this.formBuilder.group({
			title: new FormControl(null),
			programmingAssignmentType: new FormControl(null),
			courseID: new FormControl(null),
			courseName: new FormControl(null),
			limit: new FormControl(this.limit),
			offset: new FormControl(this.offset),
			// currentPage: new FormControl(this.currentPage)
		})
	}

	createProgrammingAssignmentForm(): void {
		this.assignmentForm = this.formBuilder.group({
			id: new FormControl(null),
			title: new FormControl(null, [Validators.required]),
			programmingAssignmentType: new FormControl(null, [Validators.required]),
			timeHour: new FormControl(null, [Validators.required, Validators.min(0)]),
			timeMin: new FormControl(null, [Validators.required, Validators.min(0), Validators.max(59)]),
			timeRequired: new FormControl(null),
			complexityLevel: new FormControl(null, [Validators.required]),
			score: new FormControl(null, [Validators.required, Validators.min(1)]),
			programmingQuestion: new FormControl(null),
			programmingAssignmentSubTask: new FormArray([]),

			// taskDescription: new FormControl(null), //, [Validators.required]
			// additionalComments: new FormControl(null),
			// source: new FormControl(null),
			// sourceURL: new FormControl(null), //, [Validators.required, Validators.pattern(this.URL_REGEX)]
		})
	}

	get subTasks() {
		return this.assignmentForm.get("programmingAssignmentSubTask") as FormArray
	}

	createSubTaskForm(): void {
		this.subTasks.push(this.formBuilder.group({
			id: new FormControl(null),
			resourceID: new FormControl(null, [Validators.required]),
			description: new FormControl(null)
			// source: new FormControl(null),
			// sourceURL: new FormControl(null, [Validators.pattern(this.URL_REGEX)]),
		}))
	}

	deleteSubTask(index: number): void {
		this.subTasks.removeAt(index)
		this.assignmentForm.markAsDirty()
	}

	onAssignmentTypeChange(): void {
		// console.log(this.assignmentForm.get("programmingAssignmentType").value);

		if (this.assignmentForm.get("programmingAssignmentType").value == this.READING) {
			this.createSubTaskForm()
			this.assignmentForm.get("programmingQuestion").setValue(null)
			this.assignmentForm.get("programmingQuestion").clearValidators()
			this.utilService.updateValueAndValiditors(this.assignmentForm)
			return
		}

		// console.log(this.assignmentForm.get("programmingAssignmentType").value);

		if (this.assignmentForm.get("programmingAssignmentType").value == this.CODING) {
			this.subTasks.clear()
			this.assignmentForm.get("programmingQuestion").setValidators([Validators.required])
			this.assignmentForm.get("programmingAssignmentSubTask").setValue([])
			this.utilService.updateValueAndValiditors(this.assignmentForm)
			return
		}
	}

	searchOrGetAssignments(): void {
		let queryParams = this.route.snapshot.queryParams
		if (this.utilService.isObjectEmpty(queryParams)) {
			this.changePage(1)
			return
		}
		let courseID = this.route.snapshot.queryParamMap.get("courseID")
		let courseName = this.route.snapshot.queryParamMap.get("courseName")
		if (courseID && courseName) {
			this.courseID = courseID
			this.courseName = courseName
			this.isNavigated = true
			// this.assignmentSearchForm.patchValue(queryParams)
		}
		this.assignmentSearchForm.patchValue(queryParams)
		this.searchAssignment()
	}

	changePage(pageNumber: number): void {
		// this.currentPage = pageNumber;
		// this.offset = this.currentPage - 1;

		// this.assignmentSearchForm.get("limit").setValue(this.limit)
		this.assignmentSearchForm.get("offset").setValue(pageNumber - 1)

		this.limit = this.assignmentSearchForm.get("limit").value
		this.offset = this.assignmentSearchForm.get("offset").value

		this.searchAssignment();
	}

	searchAndCloseDrawer(): void {
		this.drawer.toggle()
		this.searchAssignment()
	}

	searchAssignment(): void {
		// console.log(this.searchBatchForm.value)
		this.addSearchToQueryParams()
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
		this.getAssignments()
	}

	addSearchToQueryParams(): void {
		this.searchFormValue = { ...this.assignmentSearchForm?.value }
		this.router.navigate([], {
			relativeTo: this.route,
			queryParams: this.searchFormValue,
		})
	}

	onAddClick(): void {
		this.isViewMode = false
		this.isUpdateMode = false
		this.createProgrammingAssignmentForm()
		this.openModal(this.assignmentFormModal, "xl")
	}

	onViewClick(assignment: IProgrammingAssignment): void {
		// console.log(assignment);
		this.isViewMode = true
		this.isUpdateMode = false
		this.createProgrammingAssignmentForm()

		if (assignment.programmingAssignmentSubTask) {
			for (let index = 0; index < assignment.programmingAssignmentSubTask.length; index++) {
				this.createSubTaskForm()
				assignment.programmingAssignmentSubTask[index].resourceID = assignment.programmingAssignmentSubTask[index].resource.id
			}
		}

		this.assignmentForm.patchValue(assignment)
		let hr: number = Math.floor(Number(assignment.timeRequired) / 60);
		let min = Number(assignment.timeRequired) % 60;
		this.assignmentForm.get('timeHour').setValue(hr);
		this.assignmentForm.get('timeMin').setValue(min);
		this.assignmentForm.disable()
		this.openModal(this.assignmentFormModal, "xl")
	}

	onUpdateClick(): void {
		this.isViewMode = false
		this.isUpdateMode = true
		this.assignmentForm.enable()
	}

	onDeleteClick(assignmentID: string): void {
		this.openModal(this.deleteModal, 'md').result.then(() => {
			this.deleteProgrammingAssignment(assignmentID)
		}, (err) => {
			console.error(err);
			return
		})
	}

	setPaginationString() {
		this.paginationString = ''
		let limit = this.assignmentSearchForm.get('limit').value
		let offset = this.assignmentSearchForm.get('offset').value

		let start: number = limit * offset + 1
		let end: number = +limit + limit * offset

		if (this.totalAssignments < end) {
			end = this.totalAssignments
		}

		if (this.totalAssignments == 0) {
			this.paginationString = ''
			return
		}

		this.paginationString = `${start} - ${end}`
	}

	getAssignments(): void {
		this.spinnerService.loadingMessage = "Getting assignments"

		this.isAssignmentLoaded = true
		this.programmingAssignments = []
		this.assignmentService.getProgrammingAssignments(this.searchFormValue).subscribe((response: any) => {
			this.programmingAssignments = response.body
			this.totalAssignments = response.headers.get("X-Total-Count")
			// console.log(this.programmingAssignments);
			this.assignedCourseAssignments = []
		}, (err: any) => {
			this.totalAssignments = 0
			console.error(err)
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err?.error?.error)
		}).add(() => {
			if (this.totalAssignments == 0) {
				this.isAssignmentLoaded = false
			}
			this.setPaginationString()

		})
	}

	addProgrammingAssignment(): void {
		this.spinnerService.loadingMessage = "Adding Assignment"

		this.assignmentService.addAssignment(this.assignmentForm.value).subscribe((response: any) => {
			this.modalRef.close()
			alert("Assignment successfully added")
			this.getAssignments()
		}, (err: any) => {
			console.error(err)
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err?.error?.error)
		}).add(() => {

		})
	}

	updateProgrammingAssignment(): void {
		this.spinnerService.loadingMessage = "Updating Assignment"

		this.assignmentService.updateAssignment(this.assignmentForm.value).subscribe((response: any) => {
			this.modalRef.close()
			alert("Assignment successfully updated")
			this.getAssignments()
		}, (err: any) => {
			console.error(err)
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err?.error?.error)
		}).add(() => {

		})
	}

	deleteProgrammingAssignment(assignmentID: string): void {
		this.spinnerService.loadingMessage = "Deleting Assignment"

		this.assignmentService.deleteAssignment(assignmentID).subscribe((response: any) => {
			this.modalRef.close()
			alert("Assignment successfully deleted")
			this.getAssignments()
		}, (err: any) => {
			console.error(err)
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err?.error?.error)
		}).add(() => {

		})
	}

	setProgrammingTimeRequired(): void {
		this.assignmentForm.get('timeRequired').setValue(this.assignmentForm.get('timeHour').value * 60 + this.assignmentForm.get('timeMin').value);
		console.log("time", this.assignmentForm.get('timeRequired').value)
	}

	validateAssignmentForm(): void {
		this.setProgrammingTimeRequired()
		if (this.assignmentForm.get('timeRequired').value == 0) {
			alert("Time is Required, value cann't be Zero")
			return
		}
		if (this.assignmentForm.invalid) {
			this.assignmentForm.markAllAsTouched()
			return
		}
		if (this.isUpdateMode) {
			this.updateProgrammingAssignment()
			return
		}
		this.addProgrammingAssignment()
	}

	// ============================================================= COURSE PROGRAMMING ASSIGNMENT =============================================================

	createCourseAssignmentForm(): void {
		this.courseAssignmentForm = this.formBuilder.group({
			courseAssignments: new FormArray([])
		})
	}

	get assignmentControlArray() {
		return this.courseAssignmentForm.get("courseAssignments") as FormArray
	}

	addAssignment(): void {
		this.assignmentControlArray.push(this.formBuilder.group({
			id: new FormControl(null),
			courseID: new FormControl(this.courseID),
			programmingAssignment: new FormControl(null),
			programmingAssignmentID: new FormControl(null),
			order: new FormControl(null, [Validators.min(1)]),
			isActive: new FormControl(null),
		}))
	}

	navigatedFromCourseAssignment(): void {
		this.assignedCourseAssignments = []
		// this.openModal(this.addAssignmentModal, "lg")
	}

	onAssignClick(): void {
		this.createCourseAssignmentForm()
		this.patchAssignmentToForm()
		this.openModal(this.addCourseAssignmentModal, "lg")
	}

	patchAssignmentToForm(): void {

		for (let index = 0; index < this.assignedCourseAssignments.length; index++) {
			this.addAssignment()
			this.programmingAssignments.find((obj) => {
				if (obj.id == this.assignedCourseAssignments[index]) {
					this.assignmentControlArray.at(index).get("programmingAssignment").setValue(obj)
					this.assignmentControlArray.at(index).get("programmingAssignmentID").setValue(obj.id)
					this.assignmentControlArray.at(index).get("order").setValidators([Validators.required])
					this.assignmentControlArray.at(index).get("isActive").setValidators([Validators.required])
					this.assignmentControlArray.at(index).get("order").setValue(index + 1)
					this.assignmentControlArray.at(index).get("isActive").setValue(true)
					this.courseAssignmentForm.markAsDirty()
				}
			})
		}
	}

	deleteExtraFormFields(assignment: ICourseProgrammingAssignment): void {
		for (let field in assignment) {
			if (assignment[field] === null || assignment[field] === "") {
				delete assignment[field];
			}
			if (field == "isMarked" || field == "programmingAssignment") {
				delete assignment[field];
			}
		}
	}

	addCourseProgrammingAssignment(assignment: ICourseProgrammingAssignment, errors?: string): void {
		this.spinnerService.loadingMessage = "Adding assignment to course"

		this.deleteExtraFormFields(assignment)
		this.courseService.addCourseProgrammingAssignment(assignment, this.courseID).subscribe((response: any) => {
			// console.log(response);
		}, (err: any) => {
			console.error(err);
			errors += err.error.error
		}).add(() => {

			if (this.ongoingOperations == 0) {
				if (errors.length == 0) {
					alert("Assignments successufully added to course")
					this.modalRef.close()
					this.redirectToCourseAssignment()
				} else {
					console.error(errors);
					alert(errors)
				}
			}
		})
	}

	updateCourseProgrammingAssignment(assignment: ICourseProgrammingAssignment, errors?: string): void {
		this.spinnerService.loadingMessage = "Updating assignment to course"

		this.deleteExtraFormFields(assignment)
		this.courseService.updateCourseProgrammingAssignment(assignment, this.courseID).subscribe((response: any) => {
			// console.log(response);
		}, (err: any) => {
			console.error(err);
			errors += err.error.error
		}).add(() => {

			if (this.ongoingOperations == 0) {
				if (errors.length == 0) {
					alert("Assignments successufully updated in course")
					this.modalRef.close()
					this.redirectToCourseAssignment()
				} else {
					console.error(errors);
					alert(errors)
				}
			}
		})
	}

	redirectToCourseAssignment(): void {
		this.router.navigate([this.urlConstant.TRAINING_COURSE_MASTER_PROGRAMMING_ASSIGNMENT], {
			queryParams: {
				"courseName": this.courseName,
				"courseID": this.courseID
			}
		}).catch(err => {
			console.error(err)

		});
	}

	checkAllAssignmentsAdded(): boolean {
		for (let index = 0; index < this.programmingAssignments.length; index++) {
			if (!this.assignedCourseAssignments.includes(this.programmingAssignments[index].id)) {
				return false
			}
		}
		return true
	}

	addAllAssignments(): void {
		let isChecked = this.checkAllAssignmentsAdded()

		for (let index = 0; index < this.programmingAssignments.length; index++) {
			if (isChecked) {
				this.removeAssignedProgrammingAssignment(this.programmingAssignments[index].id)
				// index--
				continue
			}
			if (!this.assignedCourseAssignments.includes(this.programmingAssignments[index].id)) {
				this.assignProgrammingAssignment(this.programmingAssignments[index].id)
			}
		}
		// console.log(this.assignedCourseAssignments);
	}

	checkAssignmentAdded(programmingAssignmentID: string): boolean {
		return this.assignedCourseAssignments.includes(programmingAssignmentID)
	}

	toggleAssignment(programmingAssignmentID: string) {
		if (this.assignedCourseAssignments.includes(programmingAssignmentID)) {
			this.removeAssignedProgrammingAssignment(programmingAssignmentID)
			return
		}
		// add
		if (!this.assignedCourseAssignments.includes(programmingAssignmentID)) {
			this.assignProgrammingAssignment(programmingAssignmentID)
		}
		// console.log(this.assignedCourseAssignments);
	}

	assignProgrammingAssignment(programmingAssignmentID: string): void {
		this.assignedCourseAssignments.push(programmingAssignmentID)
	}

	removeAssignedProgrammingAssignment(programmingAssignmentID: string): void {
		let index = this.assignedCourseAssignments.indexOf(programmingAssignmentID)
		this.assignedCourseAssignments.splice(index, 1)
	}

	validateCourseAssignment(): void {
		if (this.courseAssignmentForm.invalid) {
			this.courseAssignmentForm.markAllAsTouched()
			return
		}

		let assignment: ICourseProgrammingAssignment[] = this.courseAssignmentForm.value.courseAssignments
		// this.removeUnmarkedAssignments(assignment)
		// console.log(assignment);
		if (this.doesDuplicateOrderExist(assignment)) {
			alert("Multiple assignments cannot have same order")
			return
		}

		let errors: string = ""
		for (let index = 0; index < assignment.length; index++) {
			if (assignment[index].id) {
				this.updateCourseProgrammingAssignment(assignment[index], errors)
				continue
			}
			this.addCourseProgrammingAssignment(assignment[index], errors)
		}
	}

	removeUnmarkedAssignments(assignment: ICourseProgrammingAssignment[]): void {
		for (let index = 0; index < assignment.length; index++) {
			if (!assignment[index].isMarked) {
				assignment.splice(index, 1)
				index = 0
			}
		}
	}

	doesDuplicateOrderExist(assignments: ICourseProgrammingAssignment[]): boolean {
		let assignmentMap = new Map<number, number>()
		if (assignments.length > 1) {
			for (let index = 0; index < assignments.length; index++) {
				// console.log(assignmentMap);
				if (assignmentMap.has(assignments[index].order)) {
					return true
				}
				assignmentMap.set(assignments[index].order, 1)
				continue
			}
		}
		return false
	}


	openModal(modalContent: any, modalSize?: string): NgbModalRef {
		if (!modalSize) {
			modalSize = "lg"
		}
		let options: NgbModalOptions = {
			ariaLabelledBy: 'modal-basic-title', keyboard: false,
			backdrop: 'static', size: modalSize, centered: true,
		}
		this.modalRef = this.modalService.open(modalContent, options)
		return this.modalRef
	}

	resetSearchAndGetAll(): void {
		this.resetSearchForm()
		this.searchFormValue = null
		this.isSearched = false
		this.router.navigate([this.urlConstant.ADMIN_PROGRAMMING_ASSIGNMENT])
		this.changePage(1)
	}

	resetSearchForm(): void {
		this.limit = this.assignmentSearchForm.get("limit").value
		this.offset = this.assignmentSearchForm.get("offset").value
		// this.currentPage = this.assignmentSearchForm.get("currentPage").value

		this.assignmentSearchForm.reset({
			limit: this.limit,
			offset: this.offset,
			// currentPage: new FormControl(this.currentPage),
		})
	}

	getProgrammingAssignmentType(): void {
		this.generalService.getGeneralTypeByType("programming_assignment_type").subscribe((respond: any) => {
			this.assignmentType = respond;
		}, (err) => {
			console.error(err)
		})
	}

	getProgrammingAssignmentSource(): void {
		this.generalService.getGeneralTypeByType("programming_assignment_source").subscribe((respond: any) => {
			this.sourceList = respond;
		}, (err) => {
			console.error(err)
		})
	}

	getComplexityLevel(): void {
		this.generalService.getGeneralTypeByType("programming_question_level").subscribe((response: any) => {
			this.complexityLevel = response
		}, (err: any) => {
			console.error(err)
		})
	}

	getProgrammingQuestionList(): void {
		let queryParams: any = {
			"programmingType": "Assignment"
		}
		this.generalService.getProgrammingQuestionList(queryParams).subscribe((respond: any) => {
			this.programmingQuestionList = respond;
		}, (err) => {
			console.error(err)
		})
	}

	getResourcesList(): void {
		let queryParams: any = {
			"isBook": 1
		}
		this.resourcesList = []
		this.resourceService.getResourcesList(queryParams).subscribe((response: any) => {
			this.resourcesList = response.body
			console.log(this.resourcesList);
		}, (err: any) => {
			console.error(err);
		})
	}



}
