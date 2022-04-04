import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, Validators, FormControl } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { DegreeService, IDegree, ISearchFilterField, ISpecialization, ISpecializationDTO } from 'src/app/service/degree/degree.service';
import { SpecializationService } from 'src/app/service/specialization/specialization.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { UrlConstant } from 'src/app/service/constant';

@Component({
	selector: 'app-admin-degree',
	templateUrl: './admin-degree.component.html',
	styleUrls: ['./admin-degree.component.css']
})
export class AdminDegreeComponent implements OnInit {

	// Components.
	degreeListing: IDegree[]

	// Flags.
	isDegreeSearched: boolean
	isSpecializationSearched: boolean
	isOperationUpdate: boolean
	isViewMode: boolean

	// Degree.
	degreeList: IDegree[]
	degreeForm: FormGroup

	// Specialization.
	specializationList: ISpecializationDTO[]
	specializationForm: FormGroup

	// Pagination for degree.
	limitDegree: number
	currentPageDegree: number
	totalDegrees: number
	offsetDegree: number
	paginationStart: number
	paginationEnd: number

	// Pagination for specialization.
	limitSpecialization: number
	currentPageSpecialization: number
	offsetSpecialization: number
	totalSpecializations: number

	// Modal.
	modalRef: any
	@ViewChild('degreeFormModal') degreeFormModal: any
	@ViewChild('deleteDegreeModal') deleteDegreeModal: any
	@ViewChild('speicalizationFormModal') specializationFormModal: any
	@ViewChild('deleteSpeicalizationModal') deleteSpeicalizationModal: any

	// Spinner.



	// Search for degree.
	degreeSearchForm: FormGroup
	degreeSearchFormValue: any
	searchDegreeFilterFieldList: ISearchFilterField[]

	// Search for specialization.
	specializationSearchForm: FormGroup
	specializationSearchFormValue: any
	searchSpecializationFilterFieldList: ISearchFilterField[]

	constructor(
		private formBuilder: FormBuilder,
		private degreeService: DegreeService,
		private specializationService: SpecializationService,
		private spinnerService: SpinnerService,
		private modalService: NgbModal,
		private router: Router,
		private route: ActivatedRoute,
		private utilService: UtilityService,
		private urlConstant: UrlConstant,
	) {
		this.initializeVariables()
		this.getAllComponents()
	}


	get ongoingOperations() {
		return this.spinnerService.ongoingOperations
	}

	ngOnInit(): void { }

	// Initialize all global variables.
	initializeVariables() {
		// Components.
		this.specializationList = [] as ISpecializationDTO[]
		this.degreeList = [] as IDegree[]
		this.degreeListing = [] as IDegree[]

		// Flags.
		this.isOperationUpdate = false
		this.isViewMode = false
		this.isDegreeSearched = false
		this.isSpecializationSearched = false

		// Degree pagination.
		this.limitDegree = 5
		this.offsetDegree = 0
		this.currentPageDegree = 0
		this.paginationStart = 0
		this.paginationEnd = 0

		// Specialization pagination.
		this.limitSpecialization = 5
		this.offsetSpecialization = 0
		this.currentPageSpecialization = 0

		// Degree search.
		this.searchDegreeFilterFieldList = []
		this.degreeSearchFormValue = {}

		// Specialization search.
		this.searchSpecializationFilterFieldList = []
		this.specializationSearchFormValue = {}

		// Initialize forms
		this.createDegreeForm()
		this.createSpecializationForm()
		this.createDegreeSearchForm()
		this.createSpecializationSearchForm()

		// Spinner.
		this.spinnerService.loadingMessage = "Getting degrees and specializations"

	}

	// =============================================================CREATE FORMS==========================================================================
	// Create degree form.
	createDegreeForm(): void {
		this.degreeForm = this.formBuilder.group({
			id: new FormControl(null),
			name: new FormControl(null, [Validators.required, Validators.maxLength(200)]),
		})
	}

	// Create specialization Form.
	createSpecializationForm(): void {
		this.specializationForm = this.formBuilder.group({
			id: new FormControl(null),
			branchName: new FormControl(null, [Validators.required, Validators.maxLength(200)]),
			degree: new FormControl(null, [Validators.required]),
		})
	}

	// Create specialization search form.
	createSpecializationSearchForm(): void {
		this.specializationSearchForm = this.formBuilder.group({
			branchName: new FormControl(null),
			degreeID: new FormControl(null),
		})
	}

	// Create degree search form.
	createDegreeSearchForm(): void {
		this.degreeSearchForm = this.formBuilder.group({
			name: new FormControl(null)
		})
	}

	// =============================================================DEGREE CRUD FUNCTIONS==========================================================================
	// On clicking add new degree button.
	onAddNewDegreeClick(): void {
		this.isViewMode = false
		this.isOperationUpdate = false
		this.createDegreeForm()
		this.openModal(this.degreeFormModal, 'sm')
	}

	// Add new degree.
	addDegree(): void {
		this.spinnerService.loadingMessage = "Adding Degree"


		this.degreeService.addDegree(this.degreeForm.value).subscribe((response: any) => {
			this.modalRef.close()
			this.getAllDegrees()
			this.getDegreeList()
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

	// On clicking view degree button.
	onViewDegreeClick(degree: IDegree): void {
		this.isViewMode = true
		this.createDegreeForm()
		this.degreeForm.patchValue(degree)
		this.degreeForm.disable()
		this.openModal(this.degreeFormModal, 'sm')
	}

	// On cliking update form button in degree form.
	onUpdateDegreeClick(): void {
		this.isViewMode = false
		this.isOperationUpdate = true
		this.degreeForm.enable()
	}

	// Update degree.
	updateDegree(): void {
		this.spinnerService.loadingMessage = "Updating degree"


		this.degreeService.updateDegree(this.degreeForm.value).subscribe((response: any) => {
			this.modalRef.close()
			this.getAllDegrees()
			this.getDegreeList()
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

	// On clicking delete degree button. 
	onDeleteDegreeClick(degreeID: string): void {
		this.openModal(this.deleteDegreeModal, 'md').result.then(() => {
			this.deleteDegree(degreeID)
		}, (err) => {
			console.error(err)
			return
		})
	}

	// Delete degree after confirmation from user.
	deleteDegree(degreeID: string): void {
		this.spinnerService.loadingMessage = "Deleting degree"


		this.degreeService.deleteDegree(degreeID).subscribe((response: any) => {
			this.modalRef.close()
			this.getAllDegrees()
			this.getDegreeList()
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

	// =============================================================SPECIALIZATION CRUD FUNCTIONS==========================================================================
	// On clicking add new specialization button.
	onAddNewSpecializationClick(): void {
		this.isViewMode = false
		this.isOperationUpdate = false
		this.createSpecializationForm()
		this.openModal(this.specializationFormModal, 'md')
	}

	// Add new specialization.
	addSpecialization(): void {
		this.spinnerService.loadingMessage = "Adding Specialization"


		let specialization: ISpecialization = this.specializationForm.value
		this.patchIDFromObjects(specialization)
		this.specializationService.addSpecialization(this.specializationForm.value).subscribe((response: any) => {
			this.modalRef.close()
			this.getAllSpecializations()
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

	// On clicking view specialization button.
	onViewSpecializationClick(specialization: ISpecializationDTO): void {
		this.isViewMode = true
		this.createSpecializationForm()
		this.specializationForm.patchValue(specialization)
		console.log(this.specializationForm.value)
		this.specializationForm.disable()
		this.openModal(this.specializationFormModal, 'md')
	}

	// On cliking update form button in specialization form.
	onUpdateSpecializationClick(): void {
		this.isViewMode = false
		this.isOperationUpdate = true
		this.specializationForm.enable()
	}

	// Update specialization.
	updateSpecialization(): void {
		this.spinnerService.loadingMessage = "Updating Specialization"


		let specialization: ISpecialization = this.specializationForm.value
		this.patchIDFromObjects(specialization)
		this.specializationService.updateSpecialization(this.specializationForm.value).subscribe((response: any) => {
			this.modalRef.close()
			this.getAllSpecializations()
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

	// On clicking delete specialization button. 
	onDeleteSpecializationClick(specializationID: string): void {
		this.openModal(this.deleteSpeicalizationModal, 'md').result.then(() => {
			this.deleteSpeicalization(specializationID)
		}, (err) => {
			console.error(err)
			return
		})
	}

	// Delete specialization after confirmation from user.
	deleteSpeicalization(specializationID: string): void {
		this.spinnerService.loadingMessage = "Deleting Specialization"


		this.specializationService.deleteSpecialization(specializationID).subscribe((response: any) => {
			this.modalRef.close()
			this.getAllSpecializations()
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

	// =============================================================DEGREE SEARCH FUNCTIONS==========================================================================
	// Reset degree search form and renaviagte page.
	resetSearchAndGetAllForDegree(): void {
		this.searchDegreeFilterFieldList = []
		this.degreeSearchForm.reset()
		this.degreeSearchFormValue = {}
		this.changeDegreePage(1)
		this.isDegreeSearched = false
		this.router.navigate([this.urlConstant.TALENT_DEGREE])
	}

	// Reset degree search form.
	resetDegreeSearchForm(): void {
		this.searchDegreeFilterFieldList = []
		this.degreeSearchForm.reset()
	}

	// Search degrees.
	searchDegrees(): void {
		this.degreeSearchFormValue = { ...this.degreeSearchForm?.value }
		this.router.navigate([], {
			relativeTo: this.route,
			queryParams: this.degreeSearchFormValue,
		})
		for (let field in this.degreeSearchFormValue) {
			if (this.degreeSearchFormValue[field] === null || this.degreeSearchFormValue[field] === "") {
				delete this.degreeSearchFormValue[field]
			} else {
				this.isDegreeSearched = true
			}
		}
		this.searchDegreeFilterFieldList = []
		for (var property in this.degreeSearchFormValue) {
			let text: string = property
			let result: string = text.replace(/([A-Z])/g, " $1");
			let finalResult: string = result.charAt(0).toUpperCase() + result.slice(1);
			let valueArray: any[] = []
			if (Array.isArray(this.degreeSearchFormValue[property])) {
				valueArray = this.degreeSearchFormValue[property]
			}
			else {
				valueArray.push(this.degreeSearchFormValue[property])
			}
			this.searchDegreeFilterFieldList.push(
				{
					propertyName: property,
					propertyNameText: finalResult,
					valueList: valueArray
				})
		}
		if (this.searchDegreeFilterFieldList.length == 0) {
			this.resetSearchAndGetAllForDegree()
		}
		if (!this.isDegreeSearched) {
			return
		}
		this.spinnerService.loadingMessage = "Searching Degrees"
		this.changeDegreePage(1)
	}

	// Delete search criteria from degree search form by search name.
	deleteSearchDegreeCriteria(searchName: string): void {
		this.degreeSearchForm.get(searchName).setValue(null)
		this.searchDegrees()
	}

	// =============================================================SPECIALIZATION SEARCH FUNCTIONS==========================================================================
	// Reset specialization search form and renaviagte page.
	resetSearchAndGetAllForSpecialization(): void {
		this.searchSpecializationFilterFieldList = []
		this.specializationSearchForm.reset()
		this.specializationSearchFormValue = {}
		this.changeSpecializationPage(1)
		this.isSpecializationSearched = false
		this.router.navigate([this.urlConstant.TALENT_DEGREE])
	}

	// Reset specialization search form.
	resetSpecializationSearchForm(): void {
		this.searchSpecializationFilterFieldList = []
		this.specializationSearchForm.reset()
	}

	// Search specializations.
	searchSpecializations(): void {
		this.specializationSearchFormValue = { ...this.specializationSearchForm?.value }
		this.router.navigate([], {
			relativeTo: this.route,
			queryParams: this.specializationSearchFormValue,
		})
		for (let field in this.specializationSearchFormValue) {
			if (this.specializationSearchFormValue[field] === null || this.specializationSearchFormValue[field] === "") {
				delete this.specializationSearchFormValue[field]
			} else {
				this.isSpecializationSearched = true
			}
		}
		this.searchSpecializationFilterFieldList = []
		for (var property in this.specializationSearchFormValue) {
			let text: string = property
			let result: string = text.replace(/([A-Z])/g, " $1");
			let finalResult: string = result.charAt(0).toUpperCase() + result.slice(1);
			let valueArray: any[] = []
			if (Array.isArray(this.specializationSearchFormValue[property])) {
				valueArray = this.specializationSearchFormValue[property]
			}
			else {
				valueArray.push(this.specializationSearchFormValue[property])
			}
			this.searchSpecializationFilterFieldList.push(
				{
					propertyName: property,
					propertyNameText: finalResult,
					valueList: valueArray
				})
		}
		if (this.searchSpecializationFilterFieldList.length == 0) {
			this.resetSearchAndGetAllForSpecialization()
		}
		if (!this.isSpecializationSearched) {
			return
		}
		this.spinnerService.loadingMessage = "Searching Specializations"
		this.changeSpecializationPage(1)
	}

	// Delete search criteria from specialization search form by search name.
	deleteSearchSpecializationCriteria(searchName: string): void {
		this.specializationSearchForm.get(searchName).setValue(null)
		this.searchSpecializations()
	}

	// =======================================OTHER FUNCTIONS FOR DEGREE AND SPECIALZATION==========================================
	// Page change for degree.
	changeDegreePage(pageNumber: number): void {
		this.currentPageDegree = pageNumber
		this.offsetDegree = this.currentPageDegree - 1
		this.getAllDegrees()
	}

	// Page change for specialization.
	changeSpecializationPage(pageNumber: number): void {
		this.currentPageSpecialization = pageNumber
		this.offsetSpecialization = this.currentPageSpecialization - 1
		this.getAllSpecializations()
	}

	// Checks the url's query params and decides to call whether to call get or search.
	searchOrGetDegrees() {
		let queryParams = this.route.snapshot.queryParams
		if (this.utilService.isObjectEmpty(queryParams)) {
			this.getAllDegrees()
			return
		}
		this.degreeSearchForm.patchValue(queryParams)
		this.searchDegrees()
	}

	// Checks the url's query params and decides to call whether to call get or search.
	searchOrGetSpecializations() {
		let queryParams = this.route.snapshot.queryParams
		if (this.utilService.isObjectEmpty(queryParams)) {
			this.getAllSpecializations()
			return
		}
		this.specializationSearchForm.patchValue(queryParams)
		this.searchSpecializations()
	}

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

	// On clicking sumbit button in degree form.
	onDegreeFormSubmit(): void {
		if (this.degreeForm.invalid) {
			this.degreeForm.markAllAsTouched()
			return
		}
		if (this.isOperationUpdate) {
			this.updateDegree()
			return
		}
		this.addDegree()
	}

	// On clicking sumbit button in specialization form.
	onSpecializationFormSubmit(): void {
		if (this.specializationForm.invalid) {
			this.specializationForm.markAllAsTouched()
			return
		}
		if (this.isOperationUpdate) {
			this.updateSpecialization()
			return
		}
		this.addSpecialization()
	}

	// Set total list on current page.
	setPaginationString(limit: number, offset: number, total: number): void {
		this.paginationStart = limit * offset + 1
		this.paginationEnd = +limit + limit * offset
		if (total < this.paginationEnd) {
			this.paginationEnd = total
		}
	}

	// Compare for select option field.
	compareFn(optionOne: any, optionTwo: any): boolean {
		if (optionOne == null && optionTwo == null) {
			return true
		}
		if (optionTwo != undefined && optionOne != undefined) {
			return optionOne.id === optionTwo.id
		}
		return false
	}

	// On changing tab get the component lists.
	onTabChange(event: any) {
		if (event == 1) {
			this.getAllDegrees()
		}
		if (event == 2) {
			this.getAllSpecializations()
		}
	}

	// Extract ID from objects and delete objects before adding or updating.
	patchIDFromObjects(specialization: ISpecialization): void {
		if (this.specializationForm.get('degree').value) {
			specialization.degreeID = this.specializationForm.get('degree').value.id
			delete specialization['degree']
		}
	}

	// =============================================================GET FUNCTIONS==========================================================================
	// Get all components.
	getAllComponents(): void {
		// this.searchOrGetSpecializations()
		this.getDegreeList()
		this.searchOrGetDegrees()
	}

	// Get all degrees.
	getAllDegrees() {
		this.spinnerService.loadingMessage = "Getting All Degrees"


		this.degreeSearchFormValue.limit = this.limitDegree
		this.degreeSearchFormValue.offset = this.offsetDegree
		this.degreeService.getAllDegrees(this.degreeSearchFormValue).subscribe((response) => {
			this.totalDegrees = response.headers.get('X-Total-Count')
			this.degreeList = response.body
		}, error => {
			console.error(error)
			if (error.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
			}
		}).add(() => {
			this.setPaginationString(this.limitDegree, this.offsetDegree, this.totalDegrees)
		})
	}

	// Get all specialization list.
	getAllSpecializations() {
		this.spinnerService.loadingMessage = "Getting Specializations"


		this.specializationSearchFormValue.limit = this.limitSpecialization
		this.specializationSearchFormValue.offset = this.offsetSpecialization
		this.specializationService.getAllSpecializations(this.specializationSearchFormValue).subscribe((response) => {
			this.totalSpecializations = response.headers.get('X-Total-Count')
			this.specializationList = response.body
		}, error => {
			console.error(error)
			if (error.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
			}
		}).add(() => {
			this.setPaginationString(this.limitSpecialization, this.offsetSpecialization, this.totalSpecializations)
		})
	}

	// Get degree list.
	getDegreeList() {
		let queryParams: any = {
			limit: -1,
			offset: 0,
		}
		this.degreeService.getAllDegrees(queryParams).subscribe((response) => {
			this.degreeListing = response.body
		}, err => {
			console.error(err)
		})
	}
}
