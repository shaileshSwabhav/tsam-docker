import { Location } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { Role, UrlConstant } from 'src/app/service/constant';
import { CourseService, ICourseProgrammingAssignment, ICourseSession } from 'src/app/service/course/course.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { IProgrammingAssignment, ProgrammingAssignmentService } from 'src/app/service/programming-assignment/programming-assignment.service';
import { ResourceService } from 'src/app/service/resource/resource.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { LocalService } from 'src/app/service/storage/local.service';

@Component({
	selector: 'app-course-programming-assignment',
	templateUrl: './course-programming-assignment.component.html',
	styleUrls: ['./course-programming-assignment.component.css']
})
export class CourseProgrammingAssignmentComponent implements OnInit {

	// components
	assignmentType: any[]
	sourceList: any[]
	programmingQuestionList: any[]
	complexityLevel: any[]
	resourcesList: any[]
	courseSessionList: any[]
	courseTopicList: any[]

	// programming-assignment.
	programmingAssignments: IProgrammingAssignment[]
	courseAssignments: ICourseProgrammingAssignment[]
	courseSessionAssignments: ICourseSession[]
	totalCourseAssignment: number
	totalProgrammingAssignment: number
	assignedCourseAssignments: string[]
	assignmentTitle: string

	// forms.
	searchAssignmentForm: FormGroup
	courseAssignmentForm: FormGroup
	selectedCourseAssignmentForm: FormGroup
	programmingAssignmentForm: FormGroup
	assignmentForm: FormGroup
	programmingAssignmentSearchForm: FormGroup

	// course details.
	courseID: string
	courseName: string

	// spinner



	// pagination
	limit: number
	offset: number
	currentPage: number
	paginationString: string
	programmingAssignmentLimit: number
	programmingAssignmentOffset: number
	programmingAssignmentCurrentPage: number

	// search
	searchFormValue: any
	searchProgrammingAssignmentFormValue: any
	isSearched: boolean
	isProgrammingAssignmentSearched: boolean

	// modal
	modalRef: any
	@ViewChild("deleteModal") deleteModal: any
	@ViewChild("addAssignmentModal") addAssignmentModal: any
	@ViewChild("updateCourseAssignmentModal") updateCourseAssignmentModal: any
	@ViewChild("assignmentFormModal") assignmentFormModal: any
	@ViewChild("addCourseAssignmentModal") addCourseAssignmentModal: any
	@ViewChild('drawer') drawer: any

	// flags
	isAddMode: boolean
	isUpdateMode: boolean

	// access
	permission: IPermission

	// ck-editor
	ckConfig: any

	// constants
	private readonly IGNORE_SEARCH_FIELD: string[] = ["limit", "offset"]
	readonly READING = "Reading"
	readonly CODING = "Coding"

	constructor(
		private assignmentService: ProgrammingAssignmentService,
		private utilService: UtilityService,
		private generalService: GeneralService,
		private courseService: CourseService,
		private resourceService: ResourceService,
		private urlConstant: UrlConstant,
		private formBuilder: FormBuilder,
		private spinnerService: SpinnerService,
		private modalService: NgbModal,
		private router: Router,
		private route: ActivatedRoute,
		private location: Location,
		private localService: LocalService,
		private role: Role,
	) {
		this.extractCourseDetails()
		this.editorConfig()
		this.initializeVariables()
		this.createForms()
		this.getAllComponents()
	}


	get ongoingOperations() {
		return this.spinnerService.ongoingOperations
	}

	ngOnInit(): void {
	}

	backToPreviousPage(): void {
		this.location.back()
	}

	extractCourseDetails(): void {
		this.route.queryParamMap.subscribe(params => {
			this.courseID = params.get('courseID')
			this.courseName = params.get('courseName')
		})
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
		this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_COURSE_MASTER_PROGRAMMING_ASSIGNMENT)

		this.programmingAssignments = []
		this.courseAssignments = []
		this.courseSessionAssignments = []
		this.assignedCourseAssignments = []
		this.assignmentType = []
		this.sourceList = []
		this.programmingQuestionList = []
		this.complexityLevel = []
		this.resourcesList = []
		this.courseSessionList = []

		this.totalCourseAssignment = 0

		this.limit = 5
		this.offset = 0

		this.programmingAssignmentLimit = 5
		this.programmingAssignmentOffset = 0
		this.programmingAssignmentCurrentPage = 1

		this.isAddMode = false
		this.isUpdateMode = false
		this.isSearched = false
		this.isProgrammingAssignmentSearched = false
	}

	createForms(): void {
		this.createSearchForm()
		this.createAssignmentSearchForm()
	}

	getAllComponents(): void {
		this.getProgrammingAssignmentType()
		this.getProgrammingQuestionList()
		this.getProgrammingAssignmentSource()
		this.getComplexityLevel()
		this.getResourcesList()
		this.getCourseSessionList()
		this.getCourseSessionWiseProgrammingAssignment()
		// this.getCourseProgrammingAssignment()
	}

	createSearchForm(): void {
		this.searchAssignmentForm = this.formBuilder.group({
			courseID: new FormControl(this.courseID),
			limit: new FormControl(this.limit),
			offset: new FormControl(this.offset),
			name: new FormControl(null),
			isActive: new FormControl(null),
		})
	}

	createAssignmentSearchForm(): void {
		this.programmingAssignmentSearchForm = this.formBuilder.group({
			title: new FormControl(null),
			limit: new FormControl(this.programmingAssignmentLimit),
			offset: new FormControl(this.programmingAssignmentOffset),
		})
	}

	createCourseAssignmentForm(): void {
		this.courseAssignmentForm = this.formBuilder.group({
			courseAssignments: new FormArray([])
		})
	}

	createSelectedCourseAssignmentForm(): void {
		this.selectedCourseAssignmentForm = this.formBuilder.group({
			selectedCourseAssignments: new FormArray([])
		})
	}

	get assignmentControlArray() {
		return this.courseAssignmentForm.get("courseAssignments") as FormArray
	}

	get selectedAssignmentControlArray() {
		return this.selectedCourseAssignmentForm.get("selectedCourseAssignments") as FormArray
	}

	addAssignment(): void {
		this.assignmentControlArray.push(this.formBuilder.group({
			id: new FormControl(null),
			courseID: new FormControl(this.courseID),
			programmingAssignmentID: new FormControl(null),
			courseSessionID: new FormControl(null),
			order: new FormControl(null, [Validators.min(1)]),
			isActive: new FormControl(null),
			isMarked: new FormControl(false),
			showDetails: new FormControl(false),
		}))
	}

	addSelectedCourseSessionAssignment(): void {
		this.selectedAssignmentControlArray.push(this.formBuilder.group({
			id: new FormControl(null),
			courseID: new FormControl(this.courseID),
			programmingAssignmentID: new FormControl(null),
			courseSessionID: new FormControl(null),
			order: new FormControl(null, [Validators.min(1)]),
			isActive: new FormControl(null),
		}))
	}

	redirectToProgrammingAssignment(): void {
		this.router.navigate([this.urlConstant.ADMIN_PROGRAMMING_ASSIGNMENT], {
			queryParams: {
				"courseName": this.courseName,
				"courseID": this.courseID,
			}
		}).catch(err => {
			console.error(err)

		});
	}

	onAssignClick(): void {
		this.createCourseAssignmentForm()
		this.patchAssignmentToForm()
		this.openModal(this.addCourseAssignmentModal, "lg")
	}

	onAddClick(): void {
		this.isAddMode = true
		this.isUpdateMode = false
		this.isProgrammingAssignmentSearched = false
		this.searchProgrammingAssignmentFormValue = null

		this.programmingAssignmentLimit = 5
		this.programmingAssignmentOffset = 0
		this.assignedCourseAssignments = []
		// console.log(this.courseAssignments);

		this.createAssignmentSearchForm()
		this.createCourseAssignmentForm()
		this.createSelectedCourseAssignmentForm()

		this.changeAssignmentPage(1)
		this.openModal(this.addAssignmentModal, "xl")
	}

	onCourseAssignmentFormInput(assignmentControl: FormGroup): void {
		// console.log(assignmentControl.value);
		for (let index = 0; index < this.selectedAssignmentControlArray.controls.length; index++) {
			if (this.selectedAssignmentControlArray.at(index).get("programmingAssignmentID").value ==
				assignmentControl.get("programmingAssignmentID").value) {
				this.selectedAssignmentControlArray.at(index).patchValue(assignmentControl.value)
			}
		}
		// console.log(this.selectedAssignmentControlArray.value);
	}


	onUpdateClick(assignment: ICourseProgrammingAssignment): void {
		this.isAddMode = false
		this.isUpdateMode = true
		this.assignmentTitle = assignment.programmingAssignment.title
		this.createCourseAssignmentForm()
		this.addAssignment()

		// used tempAssignment array as patchValue reqiures an array
		let tempAssignment: ICourseProgrammingAssignment[] = []
		tempAssignment.push(assignment)
		this.courseAssignmentForm.get("courseAssignments").patchValue(tempAssignment)
		this.assignmentControlArray.at(0).get("programmingAssignmentID").setValue(assignment.programmingAssignment.id)
		this.assignmentControlArray.at(0).get("order").setValidators([Validators.required, Validators.min(1)])
		this.assignmentControlArray.at(0).get("courseSessionID").setValidators([Validators.required])
		this.assignmentControlArray.at(0).get("isActive").setValidators([Validators.required])
		this.openModal(this.updateCourseAssignmentModal, "md")
	}

	onStatusChangeClick(assignment: ICourseProgrammingAssignment): void {
		assignment.programmingAssignmentID = assignment.programmingAssignment.id
		assignment.isActive = !assignment.isActive
		this.updateCourseProgrammingAssignmentStatus(assignment)
	}

	createAddAssignmentComponents(): void {
		this.createCourseAssignmentForm()
		this.patchAssignmentToForm()
	}


	patchAssignmentToForm(): void {
		for (let index = 0; index < this.programmingAssignments.length; index++) {
			this.addAssignment()
			this.assignmentControlArray.at(index).get("isMarked").setValue(false)
			this.assignmentControlArray.at(index).get("programmingAssignmentID").setValue(this.programmingAssignments[index].id)
			this.assignmentControlArray.at(index).get("order").disable()
			this.assignmentControlArray.at(index).get("isActive").disable()
			this.assignmentControlArray.at(index).get("courseSessionID").disable()


			// currently selected assignment for adding it to course-session
			this.checkSelectedAssignments(index)

			// already assigned assignments to course-session
			this.assignCourseAssignmentDetails(index)
		}
	}

	checkSelectedAssignments(index: number): void {
		for (let j = 0; j < this.selectedAssignmentControlArray.controls.length; j++) {
			if (this.assignmentControlArray.at(index).get("programmingAssignmentID").value ==
				this.selectedAssignmentControlArray.at(j).get("programmingAssignmentID").value) {
				this.assignmentControlArray.at(index).get("isMarked").setValue(true)
				this.assignmentControlArray.at(index).get("courseSessionID").setValue(this.selectedAssignmentControlArray.at(j).get("courseSessionID").value)
				this.assignmentControlArray.at(index).get("order").setValue(this.selectedAssignmentControlArray.at(j).get("order").value)
				this.assignmentControlArray.at(index).get("isActive").setValue(this.selectedAssignmentControlArray.at(j).get("isActive").value)

				this.assignmentControlArray.at(index).get("order").enable()
				this.assignmentControlArray.at(index).get("isActive").enable()
				this.assignmentControlArray.at(index).get("courseSessionID").enable()
			}
		}
		// console.log(this.assignmentControlArray.value);
	}

	assignCourseAssignmentDetails(index: number): void {
		for (let j = 0; j < this.courseAssignments.length; j++) {
			if (this.courseAssignments[j].programmingAssignment?.id == this.programmingAssignments[index].id) {
				this.assignedCourseAssignments.push(this.programmingAssignments[index].id)

				this.assignmentControlArray.at(index).get("id").setValue(this.courseAssignments[j].id)
				this.assignmentControlArray.at(index).get("order").setValue(this.courseAssignments[j].order)
				this.assignmentControlArray.at(index).get("isActive").setValue(this.courseAssignments[j].isActive)
				this.assignmentControlArray.at(index).get("isMarked").setValue(true)
				this.assignmentControlArray.at(index).get("courseSessionID").setValue(this.courseAssignments[j].courseSessionID)

				this.assignmentControlArray.at(index).get("order").enable()
				this.assignmentControlArray.at(index).get("isActive").enable()

				// add course-session assignment to selected form
				this.addSelectedCourseSessionAssignment()
				let len = this.selectedAssignmentControlArray.controls.length
				this.selectedAssignmentControlArray.at(len - 1).patchValue(this.assignmentControlArray.at(index).value)

				this.assignmentControlArray.at(index).get("isMarked").disable()
			}
		}

		// console.log(this.selectedAssignmentControlArray.value);
	}


	onDeleteClick(courseAssignmentID: string) {
		this.openModal(this.deleteModal, "md").result.then(() => {
			this.deleteCourseProgrammingAssignment(courseAssignmentID)
		}, (err) => {
			console.error(err);
			return
		})
	}

	// setCourseAssignmentFlags(): void {
	// 	for (let index = 0; index < this.courseAssignments.length; index++) {
	// 		this.courseAssignments[index].showAssignmentDetails = false
	// 	}
	// }

	changePage(pageNumber: number): void {
		this.currentPage = pageNumber;
		this.offset = this.currentPage - 1;

		this.searchAssignmentForm.get("offset").setValue(this.offset)
		this.searchAssignmentForm.get("limit").setValue(this.limit)

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
				if (field != "limit" && field != "offset") {
					this.isSearched = true
				}
				flag = false
			}
		}

		// No API call on empty search.
		if (flag) {
			return
		}
		this.getCourseSessionWiseProgrammingAssignment()
	}

	addSearchToQueryParams(): void {
		this.searchFormValue = { ...this.searchAssignmentForm?.value }
		this.router.navigate([], {
			relativeTo: this.route,
			queryParams: this.searchFormValue,
		})
	}

	setPaginationString() {
		this.paginationString = ''
		let limit = this.searchAssignmentForm.get('limit').value
		let offset = this.searchAssignmentForm.get('offset').value


		let start: number = limit * offset + 1
		let end: number = +limit + limit * offset
		if (this.totalCourseAssignment < end) {
			end = this.totalCourseAssignment
		}
		if (this.totalCourseAssignment == 0) {
			this.paginationString = ''
			return
		}
		this.paginationString = `${start} - ${end}`
	}

	getCourseSessionWiseProgrammingAssignment(): void {
		this.courseAssignments = []
		this.totalCourseAssignment = 0
		this.spinnerService.loadingMessage = "Getting all course assignments"

		let queryParams: any = {
			"isGetAssignments": 1
		}
		this.courseService.getSessionsForCourse(this.courseID, queryParams).subscribe((response: any) => {
			// this.courseService.getCourseSessionWiseProgrammingAssignment(this.courseID).subscribe((response: any) => {
			this.courseSessionAssignments = response.body
			this.getCourseProgrammingAssignment()
			// console.log(this.courseSessionAssignments);
			// this.totalCourseAssignment = response.headers.get("X-Total-Count")
			// this.setCourseAssignmentFlags()
		}, (err: any) => {
			console.error(err)
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err.error.error)
		}).add(() => {
			this.setPaginationString()

		})
	}

	getCourseProgrammingAssignment(): void {
		this.courseAssignments = []
		this.totalCourseAssignment = 0

		for (let session of this.courseSessionAssignments) {
			if (session.courseProgrammingAssignment?.length > 0) {
				this.courseAssignments.push(...session.courseProgrammingAssignment)
			}
		}
		this.totalCourseAssignment = this.courseAssignments.length
	}

	deleteExtraFormFields(assignment: ICourseProgrammingAssignment): void {
		for (let field in assignment) {
			if (assignment[field] === null || assignment[field] === "") {
				delete assignment[field];
			}
			if (field == "isMarked" || field == "showDetails") {
				delete assignment[field];
			}
		}
	}

	addCourseProgrammingAssignment(assignment: ICourseProgrammingAssignment, errors: string[] = []): void {
		this.spinnerService.loadingMessage = "Adding assignment to course"

		this.deleteExtraFormFields(assignment)
		// console.log(assignment);
		this.courseService.addCourseProgrammingAssignment(assignment, this.courseID).subscribe((response: any) => {
			// console.log(response);
		}, (err: any) => {
			errors.push(err.error.error)
		}).add(() => {

			if (this.ongoingOperations == 0) {
				if (errors.length == 0) {
					alert("Assignments successufully added to course")
					this.modalRef.close()
					this.getCourseSessionWiseProgrammingAssignment()
				} else {
					console.error(errors);
					let errorString = ""
					for (let index = 0; index < errors.length; index++) {
						errorString += (index == 0 ? "" : "\n") + errors[index]
					}
					alert(errorString)
				}
			}
		})
	}

	updateCourseProgrammingAssignment(assignment: ICourseProgrammingAssignment, errors: string[] = []): void {
		this.spinnerService.loadingMessage = "Updating assignment to course"

		this.deleteExtraFormFields(assignment)
		this.courseService.updateCourseProgrammingAssignment(assignment, this.courseID).subscribe((response: any) => {
			// console.log(response);
		}, (err: any) => {
			errors.push(err.error.error)
		}).add(() => {

			if (this.ongoingOperations == 0) {
				if (errors.length == 0) {
					alert("Assignments successufully updated in course")
					this.modalRef.close()
					this.getCourseSessionWiseProgrammingAssignment()
				} else {
					console.error(errors);
					let errorString = ""
					for (let index = 0; index < errors.length; index++) {
						errorString += (index == 0 ? "" : "\n") + errors[index]
					}
					alert(errorString)
				}
			}
		})
	}

	updateCourseProgrammingAssignmentStatus(assignment: ICourseProgrammingAssignment): void {
		this.spinnerService.loadingMessage = "Updating assignment to course"

		this.deleteExtraFormFields(assignment)
		this.courseService.updateCourseProgrammingAssignment(assignment, this.courseID).subscribe((response: any) => {
			// console.log(response);
			alert("Assignments successufully updated in course")
			this.getCourseSessionWiseProgrammingAssignment()
		}, (err: any) => {
			assignment.isActive = !assignment.isActive
			console.error(err)
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err.error.error)
		}).add(() => {

		})
	}

	deleteCourseProgrammingAssignment(courseAssignmentID: string) {
		this.spinnerService.loadingMessage = "Deleting assignment from course"

		this.courseService.deleteCourseProgrammingAssignment(courseAssignmentID, this.courseID).subscribe((response: any) => {
			// console.log(response);
			this.getCourseSessionWiseProgrammingAssignment()
		}, (err: any) => {
			console.error(err)
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err.error.error)
		}).add(() => {

		})
	}

	toggleAssignment(assignmentControl: any) {
		if (this.assignedCourseAssignments.includes(assignmentControl.get("programmingAssignmentID").value)) {
			this.removeAssignedProgrammingAssignment(assignmentControl)
			return
		}
		// add
		if (!this.assignedCourseAssignments.includes(assignmentControl.get("programmingAssignmentID").value)) {
			this.assignProgrammingAssignment(assignmentControl)
		}
	}

	assignProgrammingAssignment(assignmentControl: FormGroup): void {
		this.assignedCourseAssignments.push(assignmentControl.get("programmingAssignmentID").value)

		assignmentControl.get("programmingAssignmentID").setValidators([Validators.required])
		assignmentControl.get("courseSessionID").setValidators([Validators.required])
		assignmentControl.get("order").setValidators([Validators.required, Validators.min(1)])
		assignmentControl.get("isActive").setValidators([Validators.required])

		assignmentControl.get("order").enable()
		assignmentControl.get("isActive").enable()
		assignmentControl.get("courseSessionID").enable()

		assignmentControl.get("isMarked").setValue(true)
		assignmentControl.get("isActive").setValue(true)

		this.utilService.updateValueAndValiditors(assignmentControl)
		this.assignToSelectedProgrammingAssignment(assignmentControl)
		// console.log(this.selectedCourseAssignmentForm.value);
	}

	assignToSelectedProgrammingAssignment(assignmentControl: FormGroup): void {
		let tempSelectedCourseAssignments = this.selectedCourseAssignmentForm.value.selectedCourseAssignments

		for (let index = 0; index < tempSelectedCourseAssignments.length; index++) {
			if (this.assignedCourseAssignments.includes(tempSelectedCourseAssignments.programmingAssignmentID)) {
				console.error("value present -> ", tempSelectedCourseAssignments);
				return
			}
		}

		this.addSelectedCourseSessionAssignment()
		let len = this.selectedAssignmentControlArray.controls.length
		this.selectedAssignmentControlArray.at(len - 1).patchValue(assignmentControl.value)
	}

	removeAssignedProgrammingAssignment(assignmentControl: any): void {
		let programmingAssignmentID = assignmentControl.get("programmingAssignmentID").value
		let index = this.assignedCourseAssignments.indexOf(assignmentControl.get("programmingAssignmentID").value)
		this.assignedCourseAssignments.splice(index, 1)

		assignmentControl.get("programmingAssignmentID").clearValidators()
		assignmentControl.get("order").clearValidators()
		assignmentControl.get("courseSessionID").clearValidators()
		assignmentControl.get("isActive").clearValidators()

		assignmentControl.get("order").disable()
		assignmentControl.get("isActive").disable()
		assignmentControl.get("courseSessionID").disable()

		assignmentControl.get("isActive").setValue(null)
		assignmentControl.get("order").setValue(null)
		assignmentControl.get("isMarked").setValue(false)
		assignmentControl.get("courseSessionID").setValue(null)

		this.utilService.updateValueAndValiditors(assignmentControl)
		this.removeSelectedProgrammingAssignment(programmingAssignmentID)
	}

	removeSelectedProgrammingAssignment(programmingAssignmentID: string): void {
		let tempSelectedCourseAssignments = this.selectedCourseAssignmentForm.value.selectedCourseAssignments

		for (let index = 0; index < this.selectedAssignmentControlArray.controls.length; index++) {
			if (programmingAssignmentID == this.selectedAssignmentControlArray.at(index).get("programmingAssignmentID").value) {
				console.error("value present -> ", tempSelectedCourseAssignments);
				this.selectedAssignmentControlArray.removeAt(index)
				return
			}
		}
	}

	validateCourseAssignment(): void {

		if (this.courseAssignmentForm.invalid) {
			this.courseAssignmentForm.markAllAsTouched()
			return
		}

		let errors: string[] = []

		let selectedAssignment: ICourseProgrammingAssignment[] = this.selectedCourseAssignmentForm.value.selectedCourseAssignments
		if (this.checkAssignmentOrder(selectedAssignment)) {
			alert("Multiple assignments cannot have same order.")
			return
		}

		for (let index = 0; index < selectedAssignment.length; index++) {
			if (selectedAssignment[index].id) {
				this.updateCourseProgrammingAssignment(selectedAssignment[index], errors)
				continue
			} else {
				this.addCourseProgrammingAssignment(selectedAssignment[index], errors)
			}
		}
	}

	checkAssignmentOrder(assignment: ICourseProgrammingAssignment[]): boolean {
		let assignmentMap = new Map<number, number>()

		for (let index = 0; index < assignment.length; index++) {
			if (assignmentMap.has(assignment[index].order)) {
				return true
			}
			assignmentMap.set(assignment[index].order, 1)
			continue
		}
		return false
	}

	validateSingleAssignment(): void {
		// console.log(this.courseAssignmentForm.controls);

		if (this.courseAssignmentForm.invalid) {
			this.courseAssignmentForm.markAllAsTouched()
			return
		}
		let errors: string[] = []
		this.updateCourseProgrammingAssignment(this.courseAssignmentForm.value.courseAssignments[0], errors)
	}

	openModal(modalContent: any, modalSize: string = "lg"): NgbModalRef {

		let options: NgbModalOptions = {
			ariaLabelledBy: 'modal-basic-title', keyboard: false,
			backdrop: 'static', size: modalSize, centered: true,
		}
		this.modalRef = this.modalService.open(modalContent, options)
		return this.modalRef
	}

	changeAssignmentPage(pageNumber: number): void {
		// this.programmingAssignmentCurrentPage = pageNumber
		// this.programmingAssignmentOffset = pageNumber - 1

		// this.programmingAssignmentSearchForm.get("limit").setValue(this.programmingAssignmentLimit)
		this.programmingAssignmentLimit = this.programmingAssignmentSearchForm.get("limit").value
		this.programmingAssignmentSearchForm.get("offset").setValue(pageNumber - 1)
		this.searchProgrammingAssignment()
	}

	searchProgrammingAssignment(): void {
		// console.log(this.programmingAssignmentSearchForm.value)
		this.searchProgrammingAssignmentFormValue = { ...this.programmingAssignmentSearchForm.value }
		let flag: boolean = true

		for (let field in this.searchProgrammingAssignmentFormValue) {
			if (this.searchProgrammingAssignmentFormValue[field] === null || this.searchProgrammingAssignmentFormValue[field] === "") {
				delete this.searchProgrammingAssignmentFormValue[field];
			} else {
				if (!this.IGNORE_SEARCH_FIELD.includes(field)) {
					this.isProgrammingAssignmentSearched = true
				}
				flag = false
			}
		}

		// No API call on empty search.
		if (flag) {
			return
		}
		this.getProgrammingAssignments()
	}


	getProgrammingAssignments(): void {
		this.spinnerService.loadingMessage = "Getting programming assignments"

		// console.log(this.searchProgrammingAssignmentFormValue);
		this.assignmentService.getProgrammingAssignments(this.searchProgrammingAssignmentFormValue).subscribe((response: any) => {
			this.programmingAssignments = response.body
			this.totalProgrammingAssignment = response.headers.get("X-Total-Count")
			// console.log(this.programmingAssignments);
			// console.log(this.assignedCourseAssignments);
			this.createAddAssignmentComponents()
		}, (err: any) => {
			console.error(err);
			alert(err.error?.error)
		}).add(() => {

		})
	}

	// ============================================================= PROGRAMMING ASSIGNMENT =============================================================

	createProgrammingAssignmentForm(): void {
		this.assignmentForm = this.formBuilder.group({
			id: new FormControl(null),
			title: new FormControl(null, [Validators.required]),
			programmingAssignmentType: new FormControl(null, [Validators.required]),
			timeRequired: new FormControl(null, [Validators.required]),
			complexityLevel: new FormControl(null, [Validators.required]),
			score: new FormControl(null, [Validators.required, Validators.min(1)]),
			programmingQuestion: new FormControl(null),
			programmingAssignmentSubTask: new FormArray([]),

			// taskDescription: new FormControl(null, [Validators.required]),
			// additionalComments: new FormControl(null),
			// sourceURL: new FormControl(null, [Validators.required, Validators.pattern(this.URL_REGEX)]),
			// source: new FormControl(null),
		})
	}

	get subTasks() {
		return this.assignmentForm.get("programmingAssignmentSubTask") as FormArray
	}

	createSubTaskForm(): void {
		this.subTasks.push(this.formBuilder.group({
			id: new FormControl(null),
			resourceID: new FormControl(null, [Validators.required]),
			description: new FormControl(null),
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

	onAddAssignmentClick(): void {
		this.createProgrammingAssignmentForm()
		this.openModal(this.assignmentFormModal, "xl")
	}

	addProgrammingAssignment(): void {
		this.spinnerService.loadingMessage = "Adding Assignment"

		this.assignmentService.addAssignment(this.assignmentForm.value).subscribe((response: any) => {
			this.modalRef.close()
			alert("Assignment successfully added")
			this.getProgrammingAssignments()
		}, (err: any) => {
			console.error(err)
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err.error)
		}).add(() => {

		})
	}

	validateAssignmentForm(): void {
		// console.log(this.assignmentForm.controls);

		if (this.assignmentForm.invalid) {
			this.assignmentForm.markAllAsTouched()
			return
		}
		this.addProgrammingAssignment()
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

	getResourcesList(): void {
		let queryParams: any = {
			"isBook": 1
		}
		this.resourcesList = []
		this.resourceService.getResourcesList(queryParams).subscribe((response: any) => {
			this.resourcesList = response.body
		}, (err: any) => {
			console.error(err);
		})
	}

	getCourseSessionList(): void {
		this.courseSessionList = []
		this.courseService.getCourseSessionList(this.courseID).subscribe((response: any) => {
			this.courseSessionList = response.body
		}, (err: any) => {
			console.error(err);
		})
	}

	getProgrammingQuestionList(): void {
		this.generalService.getProgrammingQuestionList().subscribe((respond: any) => {
			this.programmingQuestionList = respond;
		}, (err) => {
			console.error(err)
		})
	}

	resetSearchAndGetAll(): void {
		this.resetSearchForm()
		this.searchFormValue = null
		this.isSearched = false
		let url: string
		if (this.localService.getJsonValue("roleName") == this.role.ADMIN || this.localService.getJsonValue("roleName") == this.role.SALES_PERSON) {
			this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_COURSE_MASTER_PROGRAMMING_ASSIGNMENT)
		}
		if (this.localService.getJsonValue("roleName") == this.role.FACULTY) {
			this.permission = this.utilService.getPermission(this.urlConstant.BANK_COURSE_PROGRAMMING_ASSIGNMENT)
		}
		this.router.navigate([url])
		this.changePage(1)
	}

	resetSearchForm(): void {
		this.limit = this.searchAssignmentForm.get("limit").value
		this.offset = this.searchAssignmentForm.get("offset").value

		this.searchAssignmentForm.reset({
			courseID: this.courseID,
			limit: this.limit,
			offset: this.offset,
		})
	}

	resetAssignmentSearchAndGetAll(): void {
		this.programmingAssignmentLimit = this.programmingAssignmentSearchForm.get("limit").value
		this.programmingAssignmentOffset = this.programmingAssignmentSearchForm.get("offset").value
		this.createAssignmentSearchForm()
		this.searchProgrammingAssignmentFormValue = null
		this.isProgrammingAssignmentSearched = false
		this.searchProgrammingAssignment()
	}




	// Compare for select option field.
	compareFn(objectOne: any, objectTwo: any): boolean {
		if (objectOne == null && objectTwo == null) {
			return true
		}
		if (objectTwo != undefined && objectOne != undefined) {
			return objectOne.id === objectTwo
		}
		return false
	}


}
