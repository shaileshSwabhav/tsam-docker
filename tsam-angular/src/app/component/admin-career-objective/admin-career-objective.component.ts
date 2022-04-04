import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, FormArray, FormControl, Validators } from '@angular/forms';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { Router, ActivatedRoute } from '@angular/router';
import { CareerObjectiveService, ICareerObjective, ISearchFilterField } from 'src/app/service/career-objective/career-objective.service';
import { UrlConstant } from 'src/app/service/constant';

@Component({
  selector: 'app-admin-career-objective',
  templateUrl: './admin-career-objective.component.html',
  styleUrls: ['./admin-career-objective.component.css']
})
export class AdminCareerObjectiveComponent implements OnInit {

  // Components.
  courseList: any[]

  // Flags.
  isSearched: boolean
  isOperationUpdate: boolean
  isViewMode: boolean
  showSearch: boolean

  // Career objective.
  careerObjectiveList: ICareerObjective[]
  careerObjectiveForm: FormGroup

  // Pagination.
  limit: number
  currentPage: number
  totalCareerObjectives: number
  offset: number
  paginationStart: number
  paginationEnd: number

  // Modal.
  modalRef: any
  @ViewChild('careerObjectiveFormModal') careerObjectiveFormModal: any
  @ViewChild('deleteCareerObjectiveModal') deleteCareerObjectiveModal: any

  // Spinner.



  // Search.
  careerObjectiveSearchForm: FormGroup
  searchFormValue: any
  searchFilterFieldList: ISearchFilterField[]

  constructor(
    public utilityService: UtilityService,
    private formBuilder: FormBuilder,
    private generalService: GeneralService,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private router: Router,
    private route: ActivatedRoute,
    private careerObjectiveService: CareerObjectiveService,
    private urlConstant: UrlConstant
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
    // Components.
    this.careerObjectiveList = [] as ICareerObjective[]
    this.courseList = [] as any[]

    // Flags.
    this.isOperationUpdate = false
    this.isViewMode = false
    this.isSearched = false
    this.showSearch = false

    // Pagination.
    this.limit = 5
    this.offset = 0
    this.currentPage = 0

    // Initialize forms
    this.createCareerObjectiveForm()
    this.createCareerObjectiveSearchForm()

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Career Objectives"


    // Search.
    this.searchFilterFieldList = []
    this.searchFormValue = {}
  }

  // =============================================================CREATE FORMS==========================================================================
  // Create career objective form.
  createCareerObjectiveForm() {
    this.careerObjectiveForm = this.formBuilder.group({
      id: new FormControl(null),
      name: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z ]*$/), Validators.maxLength(100)]),
      courses: this.formBuilder.array([]),
    })
    this.addCourse()
  }

  // Add new career objective course to career objective form.
  addCourse() {
    this.careerObjectiveCourseControlArray.push(this.formBuilder.group({
      id: new FormControl(null),
      careerObjectiveID: new FormControl(null),
      courseID: new FormControl(null, [Validators.required]),
      order: new FormControl(null, [Validators.required, Validators.min(1)]),
      technicalAspect: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z ,]*$/)]),
    }))
  }

  // Create career objective search form.
  createCareerObjectiveSearchForm() {
    this.careerObjectiveSearchForm = this.formBuilder.group({
      name: new FormControl(null)
    })
  }

  // =============================================================CAREER OBJECTIVE CRUD FUNCTIONS==========================================================================
  // On clicking add new career objective button.
  onAddNewCareerObjectiveClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.createCareerObjectiveForm()
    this.openModal(this.careerObjectiveFormModal, 'xl')
  }

  // Add new career objective.
  addCareerObjevtive(): void {
    this.spinnerService.loadingMessage = "Adding Career Objective"


    this.careerObjectiveService.addCareerObjective(this.careerObjectiveForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllCareerObjectives()
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

  // On clicking view career objective button.
  onViewCareerObjevtiveClick(careerObjective: ICareerObjective): void {
    this.isViewMode = true
    this.createCareerObjectiveForm()
    this.careerObjectiveForm.setControl("courses", this.formBuilder.array([]))
    if (careerObjective.courses && careerObjective.courses.length > 0) {
      for (let i = 0; i < careerObjective.courses.length; i++) {
        this.addCourse()
      }
    }
    this.careerObjectiveForm.patchValue(careerObjective)
    this.careerObjectiveForm.disable()
    this.openModal(this.careerObjectiveFormModal, 'xl')
  }

  // On cliking update form button in career objective form.
  onUpdateCareerObjevtiveClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = true
    this.careerObjectiveForm.enable()
  }

  // Update career objective.
  updateCareerObjective(): void {
    this.spinnerService.loadingMessage = "Updating Career Objective"


    this.careerObjectiveService.updateCareerObjective(this.careerObjectiveForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllCareerObjectives()
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

  // On clicking delete career objective button. 
  onDeleteCareerObjectiveClick(careerObjectiveID: string): void {
    this.openModal(this.deleteCareerObjectiveModal, 'md').result.then(() => {
      this.deleteCareerObjective(careerObjectiveID)
    }, (err) => {
      console.error(err)
      return
    })
  }

  // Delete career objective after confirmation from user.
  deleteCareerObjective(careerObjectiveID: string): void {
    this.spinnerService.loadingMessage = "Deleting Career Objective"

    this.modalRef.close()

    this.careerObjectiveService.deleteCareerObjective(careerObjectiveID).subscribe((response: any) => {
      this.getAllCareerObjectives()
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

  // =============================================================CAREER OBJECTIVE SEARCH FUNCTIONS==========================================================================
  // Reset search form and renaviagte page.
  resetSearchAndGetAll(): void {
    this.searchFilterFieldList = []
    this.careerObjectiveSearchForm.reset()
    this.searchFormValue = {}
    this.changePage(1)
    this.isSearched = false
    this.showSearch = false
    this.router.navigate([this.urlConstant.TALENT_CAREER_OBJECTVE])
  }

  // Reset search form.
  resetSearchForm(): void {
    this.searchFilterFieldList = []
    this.careerObjectiveSearchForm.reset()
  }

  // Search career objectives.
  searchCareerObjectives(): void {
    this.searchFormValue = { ...this.careerObjectiveSearchForm?.value }
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: this.searchFormValue,
    })
    for (let field in this.searchFormValue) {
      if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
        delete this.searchFormValue[field]
      } else {
        this.isSearched = true
      }
    }
    this.searchFilterFieldList = []
    for (var property in this.searchFormValue) {
      let text: string = property
      let result: string = text.replace(/([A-Z])/g, " $1");
      let finalResult: string = result.charAt(0).toUpperCase() + result.slice(1);
      let valueArray: any[] = []
      if (Array.isArray(this.searchFormValue[property])) {
        valueArray = this.searchFormValue[property]
      }
      else {
        valueArray.push(this.searchFormValue[property])
      }
      this.searchFilterFieldList.push(
        {
          propertyName: property,
          propertyNameText: finalResult,
          valueList: valueArray
        })
    }
    if (this.searchFilterFieldList.length == 0) {
      this.resetSearchAndGetAll()
    }
    if (!this.isSearched) {
      return
    }
    this.spinnerService.loadingMessage = "Searching Career Objectives"
    this.changePage(1)
  }

  // Delete search criteria from designation search form by search name.
  deleteSearchCriteria(searchName: string): void {
    this.careerObjectiveSearchForm.get(searchName).setValue(null)
    this.searchCareerObjectives()
  }

  // ================================================OTHER FUNCTIONS FOR CAREER OBJECTIVE===============================================
  // Page change.
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
    this.getAllCareerObjectives()
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetCareerObjectives() {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilityService.isObjectEmpty(queryParams)) {
      this.getAllCareerObjectives()
      return
    }
    this.careerObjectiveSearchForm.patchValue(queryParams)
    this.searchCareerObjectives()
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

  // On clicking sumbit button in career objective form.
  onFormSubmit(): void {
    for (let i = 0; i < this.careerObjectiveCourseControlArray.length; i++) {
      if (this.checkOrderAlreadyExists(this.careerObjectiveCourseControlArray.at(i).get('order'), i)) {
        return
      }
    }
    if (this.careerObjectiveForm.invalid) {
      this.careerObjectiveForm.markAllAsTouched()
      return
    }
    if (this.isOperationUpdate) {
      this.updateCareerObjective()
      return
    }
    this.addCareerObjevtive()
  }

  // Set total talents list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalCareerObjectives < this.paginationEnd) {
      this.paginationEnd = this.totalCareerObjectives
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

  // Get career objective courses array from careeer objective form.
  get careerObjectiveCourseControlArray(): FormArray {
    return this.careerObjectiveForm.get("courses") as FormArray
  }

  // Delete course from career objective form.
  deleteCourse(index: number) {
    if (confirm("Are you sure you want to delete course?")) {
      this.careerObjectiveForm.markAsDirty()
      this.careerObjectiveCourseControlArray.removeAt(index)
    }
  }

  // Check if order number is already given.
  checkOrderAlreadyExists(orderControl: any, index: number): boolean {
    for (let i = 0; i < this.careerObjectiveCourseControlArray.length; i++) {
      if (orderControl.value == this.careerObjectiveCourseControlArray.at(i).get('order').value && index != i
        && orderControl.value != null && this.careerObjectiveCourseControlArray.at(i).get('order').value != null) {
        return true
      }
    }
    return false
  }

  // Show course name bu course id.
  showCourseByID(courseID: string): void {
    if (this.courseList != null) {
      for (let i = 0; i < this.courseList.length; i++) {
        if (courseID == this.courseList[i].id) {
          return this.courseList[i].name
        }
      }
    }
  }

  // =============================================================GET FUNCTIONS==========================================================================
  // Get all components.
  getAllComponents(): void {
    this.searchOrGetCareerObjectives()
    this.getCourseList()
  }

  // Get course list.
  getCourseList(): void {
    this.generalService.getCourseList().subscribe(response => {
      this.courseList = response.body
    }
    ), err => {
      console.error(err)
    }
  }

  // Get all career objectives by limit and offset.
  getAllCareerObjectives() {
    this.spinnerService.loadingMessage = "Getting All Career Objectives"


    this.searchFormValue.limit = this.limit
    this.searchFormValue.offset = this.offset
    this.careerObjectiveService.getAllCareerObjectives(this.searchFormValue).subscribe((response) => {
      this.totalCareerObjectives = response.headers.get('X-Total-Count')
      this.careerObjectiveList = response.body
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.setPaginationString()
    })
  }

}
