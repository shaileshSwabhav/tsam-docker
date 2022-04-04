import { Component, OnInit, ViewChild } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { Role, UrlConstant } from 'src/app/service/constant';
import { GeneralService, ISearchFilterField } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { ModuleService } from 'src/app/service/module/module.service';
import { ProgrammingConceptService } from 'src/app/service/programming-concept/programming-concept.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { IProgrammingConcept } from 'src/app/models/programming/concept';

@Component({
  selector: 'app-programming-concept',
  templateUrl: './programming-concept.component.html',
  styleUrls: ['./programming-concept.component.css']
})
export class ProgrammingConceptComponent implements OnInit {

  // Component.
  programmingConceptLevel: any[]
  moduleList: any[]

  // Flags.
  isSearched: boolean
  isOperationUpdate: boolean
  isViewMode: boolean

  // Programming concept.
  conceptList: IProgrammingConcept[]
  conceptForm: FormGroup
  selectedConcept: IProgrammingConcept

  // Pagination.
  limit: number
  currentPage: number
  totalConcepts: number
  offset: number
  paginationStart: number
  paginationEnd: number

  // Modal.
  modalRef: any
  @ViewChild('conceptFormModal') conceptFormModal: any
  @ViewChild('deleteConceptModal') deleteConceptModal: any

  // Search.
  conceptSearchForm: FormGroup
  searchFormValue: any
  searchFilterFieldList: ISearchFilterField[]

  // Access.
  permission: IPermission
  roleName: string

  // Params.
  operation: string

  constructor(
    public utilityService: UtilityService,
    private formBuilder: FormBuilder,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private router: Router,
    private route: ActivatedRoute,
    private programmingConceptService: ProgrammingConceptService,
    private generalService: GeneralService,
    public utilService: UtilityService,
    private urlConstant: UrlConstant,
    private localService: LocalService,
    private moduleService: ModuleService,
    private role: Role,
  ) {
    this.initializeVariables()
    this.getAllComponents()
  }

  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  ngAfterViewInit() {
    if (this.operation && this.operation == "add") {
      this.onAddNewConceptClick()
      this.operation = null
      return
    }
  }

  // Initialize all global variables
  initializeVariables(): void {

    // Component.
    this.programmingConceptLevel = []

    // Programming concept.
    this.conceptList = [] as IProgrammingConcept[]

    // Flags.
    this.isOperationUpdate = false
    this.isViewMode = false
    this.isSearched = false

    // Pagination.
    this.limit = 5
    this.offset = 0
    this.currentPage = 0
    this.paginationStart = 0
    this.paginationEnd = 0

    // Initialize forms
    this.createConceptSearchForm()

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Programming Concepts"

    // Search.
    this.searchFilterFieldList = []

    // Access.
    this.roleName = this.localService.getJsonValue("roleName")
    if (this.roleName == this.role.ADMIN || this.roleName == this.role.SALES_PERSON) {
      this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.TRAINING_PROGRAMMING_CONCEPT)
    }
    if (this.roleName == this.role.FACULTY) {
      this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.BANK_PROGRAMMING_CONCEPT)
    }
  }

  // =============================================================CREATE FORMS==========================================================================

  // Create programming concept form.
  createConceptForm() {
    this.conceptForm = this.formBuilder.group({
      id: new FormControl(null),
      name: new FormControl(null, [Validators.required, Validators.maxLength(100)]),
      complexity: new FormControl(null, [Validators.required]),
      isModuleIndependent: new FormControl(false, [Validators.required]),
      description: new FormControl(null, [Validators.maxLength(2000)]),
    })
  }

  // Add new module to concept form.
  addModule(): void {
    this.conceptModuleControlArray.push(this.formBuilder.group({
      id: new FormControl(null),
      moduleID: new FormControl(null, [Validators.required]),
      programmingConceptID: new FormControl(null),
      level: new FormControl(null, [Validators.required, Validators.min(0), Validators.max(99)]),
    }))
  }

  // Get concept modules array from concept form.
  get conceptModuleControlArray(): FormArray {
    return this.conceptForm.get("modules") as FormArray
  }

  // Create programming concept search form.
  createConceptSearchForm() {
    this.conceptSearchForm = this.formBuilder.group({
      name: new FormControl(null)
    })
  }

  // =============================================================PROGRAMMING CONCEPT CRUD FUNCTIONS==========================================================================

  // On clicking add new programming concept button.
  onAddNewConceptClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.createConceptForm()
    this.openModal(this.conceptFormModal, 'lg')
  }

  // Add new programming concept.
  addConcept(): void {
    this.spinnerService.loadingMessage = "Adding Programming Concept"
    this.programmingConceptService.addConcept(this.conceptForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllConcepts()
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

  // On clicking view programming concept button.
  onViewConceptClick(concept: IProgrammingConcept): void {
    this.selectedConcept = concept
    this.isViewMode = true
    this.createConceptForm()
    // this.conceptForm.setControl("modules", this.formBuilder.array([]))
    // if (concept.modules && concept.modules.length > 0) {
    //   for (let i = 0; i < concept.modules.length; i++) {
    //     this.addModule()
    //   }
    // }
    this.conceptForm.patchValue(concept)
    this.conceptForm.disable()
    this.openModal(this.conceptFormModal, 'lg')
  }

  // On cliking update form button in programming concept form.
  onUpdateConceptClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = true
    this.conceptForm.enable()
  }

  // Update programming concept.
  updateConcept(): void {
    this.spinnerService.loadingMessage = "Updating Programming Concept"
    this.programmingConceptService.updateConcept(this.conceptForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllConcepts()
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

  // On clicking delete programming concept button. 
  onDeleteConceptClick(conceptID: string): void {
    this.openModal(this.deleteConceptModal, 'md').result.then(() => {
      this.deleteConcept(conceptID)
    }, (err) => {
      console.error(err)
      return
    })
  }

  // Delete programming concept after confirmation from user.
  deleteConcept(conceptID: string): void {
    this.spinnerService.loadingMessage = "Deleting Programming Concept"
    this.modalRef.close()
    this.programmingConceptService.deleteConcept(conceptID).subscribe((response: any) => {
      this.getAllConcepts()
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

  // =============================================================PROGRAMMING CONCEPT SEARCH FUNCTIONS==========================================================================

  // Reset search form and renaviagte page.
  resetSearchAndGetAll(): void {
    this.searchFilterFieldList = []
    this.conceptSearchForm.reset()
    this.searchFormValue = null
    this.changePage(1)
    this.isSearched = false
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: this.searchFormValue,
    })
  }

  // Reset search form.
  resetSearchForm(): void {
    this.searchFilterFieldList = []
    this.conceptSearchForm.reset()
  }

  // Search programming concepts.
  searchConcepts(): void {
    this.searchFormValue = { ...this.conceptSearchForm?.value }
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
    this.spinnerService.loadingMessage = "Searching Programming Concepts"
    this.changePage(1)
  }

  // Delete search criteria from programming concept search form by search name.
  deleteSearchCriteria(searchName: string): void {
    this.conceptSearchForm.get(searchName).setValue(null)
    this.searchConcepts()
  }

  // ================================================OTHER FUNCTIONS FOR PROGRAMMING CONCEPT===============================================

  // Page change.
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
    this.getAllConcepts()
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetConcepts() {
    let queryParams = this.route.snapshot.queryParams
    if (queryParams.operation) {
      this.operation = queryParams.operation
      // queryParams = {}
    }
    if (this.utilityService.isObjectEmpty(queryParams)) {
      this.getAllConcepts()
      return
    }
    this.conceptSearchForm.patchValue(queryParams)
    this.searchConcepts()
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

  // On clicking sumbit button in programming concept form.
  onFormSubmit(): void {
    if (this.conceptForm.invalid) {
      this.conceptForm.markAllAsTouched()
      return
    }
    if (this.isOperationUpdate) {
      this.updateConcept()
      return
    }
    this.addConcept()
  }

  // Get all invalid controls in talent form.
  public findInvalidControls(): any[] {
    const invalid = []
    const controls = this.conceptForm.controls
    for (const name in controls) {
      if (controls[name].invalid) {
        invalid.push(name)
      }
    }
    return invalid
  }

  // Set total feelings list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalConcepts < this.paginationEnd) {
      this.paginationEnd = this.totalConcepts
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

  // On changing value of is module independent field in concept form.
  onIsModuleIndependentChange(): void {

    // If is module independent true then delete all concept modules.
    if (this.conceptForm.get('isModuleIndependent').value == true) {
      this.conceptForm.setControl("modules", this.formBuilder.array([]))
    }

    // If is module independent false then add concept module.
    if (this.conceptForm.get('isModuleIndependent').value == false) {
      this.addModule()
    }
  }

  // Delete module from concept form.
  deleteModule(index: number) {
    if (confirm("Are you sure you want to delete module?")) {
      this.conceptForm.markAsDirty()
      this.conceptModuleControlArray.removeAt(index)
    }
  }

  // =============================================================GET FUNCTIONS==========================================================================

  // Get all components.
  getAllComponents(): void {
    this.getProgrammingConcepLevel()
    this.getModuleList()
    this.searchOrGetConcepts()
  }

  // Get all programming concepts by limit and offset.
  getAllConcepts() {
    this.spinnerService.loadingMessage = "Getting All Programming Concepts"
    if (!this.searchFormValue) {
      this.searchFormValue = {}
    }
    this.searchFormValue.limit = this.limit
    this.searchFormValue.offset = this.offset
    this.programmingConceptService.getAllConcepts(this.searchFormValue).subscribe((response) => {
      this.totalConcepts = +response.headers.get('X-Total-Count')
      this.conceptList = response.body
    }, error => {
      console.error(error);
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.setPaginationString()
    })
  }

  // Get programming concept complexity list.
  getProgrammingConcepLevel(): void {
    this.generalService.getGeneralTypeByType("programming_concept_level").subscribe((response: any) => {
      this.programmingConceptLevel = response
    }, (err: any) => {
      console.error(err)
    })
  }

  // Get module list.
  getModuleList(): void {
    let queryParams: any = {
      limit: -1,
      offset: 0
    }
    this.moduleService.getModule(queryParams).subscribe((response: any) => {
      this.moduleList = response.body
    }, (err: any) => {
      console.error(err)
    })
  }

}
