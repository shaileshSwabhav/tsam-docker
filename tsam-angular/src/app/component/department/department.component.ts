import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { DepartmentService, IDepartmentDTO, IRole, ISearchFilterField, ITargetCommunityFunctionDTO } from 'src/app/service/department/department.service';
import { GeneralService, IDepartment, ITargetCommunityFunction } from 'src/app/service/general/general.service';
import { TargetCommunityFunctionService } from 'src/app/service/target-community-function/target-community-function.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { UrlConstant } from 'src/app/service/constant';

@Component({
  selector: 'app-department',
  templateUrl: './department.component.html',
  styleUrls: ['./department.component.css']
})
export class DepartmentComponent implements OnInit {

  // Components.
  roleList: IRole[]
  departmentListing: IDepartment[]

  // Flags.
  isSearchedDepartment: boolean
  isSearchedFunction: boolean
  isOperationUpdate: boolean
  isViewMode: boolean

  // Department.
  departmentList: IDepartmentDTO[]
  departmentForm: FormGroup

  // Target Community Function.
  functionList: ITargetCommunityFunctionDTO[]
  functionForm: FormGroup

  // Pagination for department.
  limitDepartment: number
  currentPageDepartment: number
  totalDepartments: number
  offsetDepartment: number
  paginationStart: number
  paginationEnd: number

  // Pagination for target community function.
  limitFunction: number
  currentPageFunction: number
  totalFunctions: number
  offsetFunction: number

  // Modal.
  modalRef: any
  @ViewChild('departmentFormModal') departmentFormModal: any
  @ViewChild('deleteDepartmentModal') deleteDepartmentModal: any
  @ViewChild('functionFormModal') functionFormModal: any
  @ViewChild('deletefunctionModal') deletefunctionModal: any

  // Spinner.



  // Search for department.
  departmentSearchForm: FormGroup
  searchFormValueForDepartment: any
  searchDepartmentFilterFieldList: ISearchFilterField[]

  // Search for department.
  functionSearchForm: FormGroup
  searchFormValueForFunction: any
  searchFunctionFilterFieldList: ISearchFilterField[]

  constructor(
    private formBuilder: FormBuilder,
    private departmentService: DepartmentService,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private router: Router,
    private route: ActivatedRoute,
    private utilService: UtilityService,
    private generalService: GeneralService,
    private functionService: TargetCommunityFunctionService,
    private urlConstant: UrlConstant
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
    this.departmentList = [] as IDepartmentDTO[]
    this.roleList = [] as IRole[]
    this.departmentListing = [] as IDepartment[]
    this.functionList = [] as ITargetCommunityFunctionDTO[]

    // Flags.
    this.isOperationUpdate = false
    this.isViewMode = false
    this.isSearchedDepartment = false
    this.isSearchedFunction = false

    // Department Pagination.
    this.limitDepartment = 5
    this.offsetDepartment = 0
    this.currentPageDepartment = 0
    this.paginationStart = 0
    this.paginationEnd = 0

    // Target Community Function Pagination.
    this.limitFunction = 5
    this.offsetFunction = 0
    this.currentPageFunction = 0

    // Department search.
    this.searchDepartmentFilterFieldList = []
    this.searchFormValueForDepartment = {}

    // Function search.
    this.searchFunctionFilterFieldList = []
    this.searchFormValueForFunction = {}

    // Initialize forms
    this.createDepartmentForm()
    this.createDepartmentSearchForm()
    this.createFunctionForm()
    this.createFunctionSearchForm()

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Departments and Target Community Functions"

  }

  // =============================================================CREATE FORMS==========================================================================
  // Create department form.
  createDepartmentForm(): void {
    this.departmentForm = this.formBuilder.group({
      id: new FormControl(null),
      name: new FormControl(null, [Validators.required, Validators.maxLength(50)]),
      role: new FormControl(null, [Validators.required]),
    })
  }

  // Create department search form.
  createDepartmentSearchForm(): void {
    this.departmentSearchForm = this.formBuilder.group({
      name: new FormControl(null),
      roleID: new FormControl(null),
    })
  }

  // Create target community function form.
  createFunctionForm(): void {
    this.functionForm = this.formBuilder.group({
      id: new FormControl(null),
      functionName: new FormControl(null, [Validators.required, Validators.maxLength(100)]),
      department: new FormControl(null, [Validators.required]),
    })
  }

  // Create target community function search form.
  createFunctionSearchForm(): void {
    this.functionSearchForm = this.formBuilder.group({
      functionName: new FormControl(null),
      departmentID: new FormControl(null),
    })
  }
  // =============================================================DEPARTMENT CRUD FUNCTIONS==========================================================================
  // On clicking add new department button.
  onAddNewDepartmentClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.createDepartmentForm()
    this.openModal(this.departmentFormModal, 'md')
  }

  // Add new department.
  addDepartment(): void {
    this.spinnerService.loadingMessage = "Adding Department"


    let department: IDepartment = this.departmentForm.value
    this.patchIDFromObjectsForDepartment(department)
    this.departmentService.addDepartment(this.departmentForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllDepartments()
      this.getDepartmentList()
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

  // On clicking view department button.
  onViewDepartmentClick(department: IDepartmentDTO): void {
    this.isViewMode = true
    this.createDepartmentForm()
    this.departmentForm.patchValue(department)
    this.departmentForm.disable()
    this.openModal(this.departmentFormModal, 'md')
  }

  // On cliking update form button in department form.
  onUpdateDepartmentClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = true
    this.departmentForm.enable()
  }

  // Update department.
  updateDepartment(): void {
    this.spinnerService.loadingMessage = "Updating Department"


    let department: IDepartment = this.departmentForm.value
    this.patchIDFromObjectsForDepartment(department)
    this.departmentService.updateDepartment(this.departmentForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllDepartments()
      this.getDepartmentList()
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

  // On clicking delete department button. 
  onDeleteDepartmentClick(departmentID: string): void {
    this.openModal(this.deleteDepartmentModal, 'md').result.then(() => {
      this.deleteDepartment(departmentID)
    }, (err) => {
      console.error(err)
      return
    })
  }

  // Delete department after confirmation from user.
  deleteDepartment(departmentID: string): void {
    this.spinnerService.loadingMessage = "Deleting Department"

    this.modalRef.close()

    this.departmentService.deleteDepartment(departmentID).subscribe((response: any) => {
      this.getAllDepartments()
      this.getDepartmentList()
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

  // =============================================================TARGET COMMUNITY FUNCTION CRUD FUNCTIONS==========================================================================
  // On clicking add new target community function button.
  onAddNewFunctionClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.createFunctionForm()
    this.openModal(this.functionFormModal, 'md')
  }

  // Add new target community function.
  addTFunction(): void {
    this.spinnerService.loadingMessage = "Adding Target Community Function"


    let targetCommunityFunction: ITargetCommunityFunction = this.functionForm.value
    this.patchIDFromObjectsForFunction(targetCommunityFunction)
    this.functionService.addTargetCommunityFunction(targetCommunityFunction).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllFunctions()
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

  // On clicking view target community function button.
  onViewFunctionClick(targetCommunityFunction: ITargetCommunityFunctionDTO): void {
    this.isViewMode = true
    this.createFunctionForm()
    this.functionForm.patchValue(targetCommunityFunction)
    this.functionForm.disable()
    this.openModal(this.functionFormModal, 'md')
  }

  // On cliking update form button in target community function form.
  onUpdateFunctionClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = true
    this.functionForm.enable()
  }

  // Update target community function.
  updateFunction(): void {
    this.spinnerService.loadingMessage = "Updating Target Community Function"


    let targetCommunityFunction: ITargetCommunityFunction = this.functionForm.value
    this.patchIDFromObjectsForFunction(targetCommunityFunction)
    this.functionService.updateTargetCommunityFunction(targetCommunityFunction).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllFunctions()
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

  // On clicking delete target community function button. 
  onDeleteFunctionClick(targetCommunityFunctionID: string): void {
    this.openModal(this.deletefunctionModal, 'md').result.then(() => {
      this.deleteFunctionForm(targetCommunityFunctionID)
    }, (err) => {
      console.error(err)
      return
    })
  }

  // Delete target community function after confirmation from user.
  deleteFunctionForm(targetCommunityFunctionID: string): void {
    this.spinnerService.loadingMessage = "Deleting Target Community Function"

    this.modalRef.close()

    this.functionService.deleteTargetCommunityFunction(targetCommunityFunctionID).subscribe((response: any) => {
      this.getAllFunctions()
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

  // =============================================================DEPARTMENT SEARCH FUNCTIONS==========================================================================
  // Reset search form and renaviagte page.
  resetSearchAndGetAllForDepartment(): void {
    this.searchDepartmentFilterFieldList = []
    this.departmentSearchForm.reset()
    this.searchFormValueForDepartment = {}
    this.changePageForDepartment(1)
    this.isSearchedDepartment = false
    this.router.navigate([this.urlConstant.ADMIN_DEPARTMENT])
  }

  // Reset search form.
  resetDepartmentSearchForm(): void {
    this.searchDepartmentFilterFieldList = []
    this.departmentSearchForm.reset()
  }

  // Search departments.
  searchDepartments(): void {
    this.searchFormValueForDepartment = { ...this.departmentSearchForm?.value }
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: this.searchFormValueForDepartment,
    })
    for (let field in this.searchFormValueForDepartment) {
      if (this.searchFormValueForDepartment[field] === null || this.searchFormValueForDepartment[field] === "") {
        delete this.searchFormValueForDepartment[field]
      } else {
        this.isSearchedDepartment = true
      }
    }
    this.searchDepartmentFilterFieldList = []
    for (var property in this.searchFormValueForDepartment) {
      let text: string = property
      let result: string = text.replace(/([A-Z])/g, " $1");
      let finalResult: string = result.charAt(0).toUpperCase() + result.slice(1);
      let valueArray: any[] = []
      if (Array.isArray(this.searchFormValueForDepartment[property])) {
        valueArray = this.searchFormValueForDepartment[property]
      }
      else {
        valueArray.push(this.searchFormValueForDepartment[property])
      }
      this.searchDepartmentFilterFieldList.push(
        {
          propertyName: property,
          propertyNameText: finalResult,
          valueList: valueArray
        })
    }
    if (this.searchDepartmentFilterFieldList.length == 0) {
      this.resetSearchAndGetAllForDepartment()
    }
    if (!this.isSearchedDepartment) {
      return
    }
    this.spinnerService.loadingMessage = "Searching Departments"
    this.changePageForDepartment(1)
  }

  // Delete search criteria from department search form by search name.
  deleteSearchDepartmentCriteria(searchName: string): void {
    this.departmentSearchForm.get(searchName).setValue(null)
    this.searchDepartments()
  }

  // =============================================================TARGET COMMUNITY FUNCTION SEARCH FUNCTIONS==========================================================================
  // Reset target community function search form and renaviagte page.
  resetSearchAndGetAllForFunction(): void {
    this.searchFunctionFilterFieldList = []
    this.functionSearchForm.reset()
    this.searchFormValueForFunction = {}
    this.changeFunctionPage(1)
    this.isSearchedFunction = false
    this.router.navigate([this.urlConstant.ADMIN_DEPARTMENT])
  }

  // Reset target community function search form.
  resetFunctionSearchForm(): void {
    this.searchFunctionFilterFieldList = []
    this.functionSearchForm.reset()
  }

  // Search target community functions.
  searchFunctions(): void {
    this.searchFormValueForFunction = { ...this.functionSearchForm?.value }
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: this.searchFormValueForFunction,
    })
    for (let field in this.searchFormValueForFunction) {
      if (this.searchFormValueForFunction[field] === null || this.searchFormValueForFunction[field] === "") {
        delete this.searchFormValueForFunction[field]
      } else {
        this.isSearchedFunction = true
      }
    }
    this.searchFunctionFilterFieldList = []
    for (var property in this.searchFormValueForFunction) {
      let text: string = property
      let result: string = text.replace(/([A-Z])/g, " $1");
      let finalResult: string = result.charAt(0).toUpperCase() + result.slice(1);
      let valueArray: any[] = []
      if (Array.isArray(this.searchFormValueForFunction[property])) {
        valueArray = this.searchFormValueForFunction[property]
      }
      else {
        valueArray.push(this.searchFormValueForFunction[property])
      }
      this.searchFunctionFilterFieldList.push(
        {
          propertyName: property,
          propertyNameText: finalResult,
          valueList: valueArray
        })
    }
    if (this.searchFunctionFilterFieldList.length == 0) {
      this.resetSearchAndGetAllForDepartment()
    }
    if (!this.isSearchedFunction) {
      return
    }
    this.spinnerService.loadingMessage = "Searching Target Community Functions"
    this.changeFunctionPage(1)
  }

  // Delete search criteria from function search form by search name.
  deleteSearchFunctionCriteria(searchName: string): void {
    this.functionSearchForm.get(searchName).setValue(null)
    this.searchFunctions()
  }

  // ================================================OTHER FUNCTIONS FOR DEPARTMENT===============================================
  // Page change.
  changePageForDepartment(pageNumber: number): void {

    this.currentPageDepartment = pageNumber
    this.offsetDepartment = this.currentPageDepartment - 1
    this.getAllDepartments()
  }

  // Page change for target community function.
  changeFunctionPage(pageNumber: number): void {

    this.currentPageFunction = pageNumber
    this.offsetFunction = this.currentPageFunction - 1
    this.getAllFunctions()
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetDepartments() {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getAllDepartments()
      return
    }
    this.departmentSearchForm.patchValue(queryParams)
    this.searchDepartments()
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetFunctions() {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getAllFunctions()
      return
    }
    this.functionSearchForm.patchValue(queryParams)
    this.searchFunctions()
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

  // On clicking sumbit button in department form.
  onFormSubmitForDepartment(): void {
    if (this.departmentForm.invalid) {
      this.departmentForm.markAllAsTouched()
      return
    }
    if (this.isOperationUpdate) {
      this.updateDepartment()
      return
    }
    this.addDepartment()
  }

  // On clicking sumbit button in target community function form.
  onFormSubmitForFunction(): void {
    if (this.functionForm.invalid) {
      this.functionForm.markAllAsTouched()
      return
    }
    if (this.isOperationUpdate) {
      this.updateFunction()
      return
    }
    this.addTFunction()
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

  onTabChange(event: any) {
    if (event == 1) {
      this.getAllDepartments()
    }
    if (event == 2) {
      this.getAllFunctions()
    }
  }

  patchIDFromObjectsForDepartment(department: IDepartment): void {
    if (this.departmentForm.get('role').value) {
      department.roleID = this.departmentForm.get('role').value.id
      delete department['role']
    }
  }

  patchIDFromObjectsForFunction(targetCommunityFunction: ITargetCommunityFunction): void {
    if (this.functionForm.get('department').value) {
      targetCommunityFunction.departmentID = this.functionForm.get('department').value.id
      delete targetCommunityFunction['department']
    }
  }

  // =============================================================GET FUNCTIONS==========================================================================
  // Get all components.
  getAllComponents() {
    this.getRoleList()
    this.getDepartmentList()
    this.searchOrGetDepartments()
    // this.searchOrGetTargetCommunityFunctions()
  }

  // Get all departments.
  getAllDepartments() {
    this.spinnerService.loadingMessage = "Getting All Departments"


    this.searchFormValueForDepartment.limit = this.limitDepartment
    this.searchFormValueForDepartment.offset = this.offsetDepartment
    this.departmentService.getAllDepartments(this.searchFormValueForDepartment).subscribe((response) => {
      this.totalDepartments = response.headers.get('X-Total-Count')
      this.departmentList = response.body
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.setPaginationString(this.limitDepartment, this.offsetDepartment, this.totalDepartments)
    })
  }

  // Get role list.
  getRoleList() {
    this.generalService.getAllRoles().subscribe((response) => {
      this.roleList = response
    }, err => {
      console.error(err)
    })
  }

  // Get department list.
  getDepartmentList() {
    this.generalService.getDepartmentList().subscribe((response) => {
      this.departmentListing = response
    }, err => {
      console.error(err)
    })
  }

  // Get all target community functions.
  getAllFunctions() {
    this.spinnerService.loadingMessage = "Getting All Target Community Functions"


    this.searchFormValueForFunction.limit = this.limitFunction
    this.searchFormValueForFunction.offset = this.offsetFunction
    this.functionService.getAllTargetCommunityFunctions(this.searchFormValueForFunction).subscribe((response) => {
      this.totalFunctions = response.headers.get('X-Total-Count')
      this.functionList = response.body
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.setPaginationString(this.limitFunction, this.offsetFunction, this.totalFunctions)
    })
  }
}
