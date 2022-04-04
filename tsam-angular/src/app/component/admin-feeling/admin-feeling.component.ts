import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, FormControl, Validators, FormArray } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalRef, NgbModalOptions } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { FeelingService, IFeeling, ISearchFilterField } from 'src/app/service/feeling/feeling.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-admin-feeling',
  templateUrl: './admin-feeling.component.html',
  styleUrls: ['./admin-feeling.component.css']
})
export class AdminFeelingComponent implements OnInit {

  // Flags.
  isSearched: boolean
  isOperationUpdate: boolean
  isViewMode: boolean

  // Feeling.
  feelingList: IFeeling[]
  feelingForm: FormGroup

  // Pagination.
  limit: number
  currentPage: number
  totalFeelings: number
  offset: number
  paginationStart: number
  paginationEnd: number

  // Modal.
  modalRef: any
  @ViewChild('feelingFormModal') feelingFormModal: any
  @ViewChild('deleteFeelingModal') deleteFeelingModal: any

  // Spinner.



  // Search.
  feelingSearchForm: FormGroup
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
    private feelingService: FeelingService
  ) {
    this.initializeVariables()
    this.searchOrGetFeelings()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  //Initialize all global variables
  initializeVariables(): void {
    // Components.
    this.feelingList = [] as IFeeling[]

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
    this.createFeelingForm()
    this.createFeelingSearchForm()

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Feelings"


    // Search.
    this.searchFilterFieldList = []
    this.searchFormValue = {}
  }

  // =============================================================CREATE FORMS==========================================================================
  // Create feeling form.
  createFeelingForm() {
    this.feelingForm = this.formBuilder.group({
      id: new FormControl(null),
      feelingName: new FormControl(null, [Validators.required]),  //, Validators.pattern(/^[a-zA-Z ]*$/)
      feelingLevels: this.formBuilder.array([]),
    })
    this.addFeelingLevel()
  }

  // Add new feeling level to feeling form.
  addFeelingLevel() {
    this.feelingLevelControlArray.push(this.formBuilder.group({
      id: new FormControl(null),
      feelingID: new FormControl(null),
      levelNumber: new FormControl(null, [Validators.required, Validators.min(1)]),
      description: new FormControl(null, [Validators.required, Validators.maxLength(200)]),
    }))
  }

  // Create feeling search form.
  createFeelingSearchForm() {
    this.feelingSearchForm = this.formBuilder.group({
      feelingName: new FormControl(null)
    })
  }

  // =============================================================FEELLING CRUD FUNCTIONS==========================================================================
  // On clicking add new feeling button.
  onAddNewFeelingClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.createFeelingForm()
    this.openModal(this.feelingFormModal, 'xl')
  }

  // Add new feeling.
  addFeeling(): void {
    this.spinnerService.loadingMessage = "Adding Feeling"


    this.feelingService.addFeeling(this.feelingForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllFeelings()
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

  // On clicking view feeling button.
  onViewFeelingClick(feeling: IFeeling): void {
    this.isViewMode = true
    this.createFeelingForm()
    this.feelingForm.setControl("feelingLevels", this.formBuilder.array([]))
    if (feeling.feelingLevels && feeling.feelingLevels.length > 0) {
      for (let i = 0; i < feeling.feelingLevels.length; i++) {
        this.addFeelingLevel()
      }
    }
    this.feelingForm.patchValue(feeling)
    this.feelingForm.disable()
    this.openModal(this.feelingFormModal, 'xl')
  }

  // On cliking update form button in feeling form.
  onUpdateFeelingClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = true
    this.feelingForm.enable()
  }

  // Update feeling.
  updateFeeling(): void {
    this.spinnerService.loadingMessage = "Updating Feeling"


    this.feelingService.updateFeeling(this.feelingForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllFeelings()
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

  // On clicking delete feeling button. 
  onDeleteFeelingClick(feelingID: string): void {
    this.openModal(this.deleteFeelingModal, 'md').result.then(() => {
      this.deleteFeeling(feelingID)
    }, (err) => {
      console.error(err)
      return
    })
  }

  // Delete feeling after confirmation from user.
  deleteFeeling(feelingID: string): void {
    this.spinnerService.loadingMessage = "Deleting Feeling"


    this.feelingService.deleteFeeling(feelingID).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllFeelings()
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

  // =============================================================FEELING SEARCH FUNCTIONS==========================================================================
  // Reset search form and renaviagte page.
  resetSearchAndGetAll(): void {
    this.searchFilterFieldList = []
    this.feelingSearchForm.reset()
    this.searchFormValue = {}
    this.changePage(1)
    this.isSearched = false
    this.router.navigate(['/admin/models/feeling'])
  }

  // Reset search form.
  resetSearchForm(): void {
    this.searchFilterFieldList = []
    this.feelingSearchForm.reset()
  }

  // Search feelings.
  searchFeelings(): void {
    this.searchFormValue = { ...this.feelingSearchForm?.value }
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
    this.spinnerService.loadingMessage = "Searching Feelings"
    this.changePage(1)
  }

  // Delete search criteria from designation search form by search name.
  deleteSearchCriteria(searchName: string): void {
    this.feelingSearchForm.get(searchName).setValue(null)
    this.searchFeelings()
  }

  // ================================================OTHER FUNCTIONS FOR FEELING===============================================
  // Page change.
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
    this.getAllFeelings()
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetFeelings() {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilityService.isObjectEmpty(queryParams)) {
      this.getAllFeelings()
      return
    }
    this.feelingSearchForm.patchValue(queryParams)
    this.searchFeelings()
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

  // On clicking sumbit button in feeling form.
  onFormSubmit(): void {
    for (let i = 0; i < this.feelingLevelControlArray.length; i++) {
      if (this.checkLevelNumberAlreadyExists(this.feelingLevelControlArray.at(i).get('levelNumber'), i)) {
        return
      }
    }
    if (this.feelingForm.invalid) {
      this.feelingForm.markAllAsTouched()
      return
    }
    if (this.isOperationUpdate) {
      this.updateFeeling()
      return
    }
    this.addFeeling()
  }

  // Set total feelings list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalFeelings < this.paginationEnd) {
      this.paginationEnd = this.totalFeelings
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

  // Get feeling levels array from feeling form.
  get feelingLevelControlArray(): FormArray {
    return this.feelingForm.get("feelingLevels") as FormArray
  }

  // Delete feeling level from feeling form.
  deleteFeelingLevel(index: number) {
    if (confirm("Are you sure you want to delete feeling level?")) {
      this.feelingForm.markAsDirty()
      this.feelingLevelControlArray.removeAt(index)
    }
  }

  // Check if level number is already given.
  checkLevelNumberAlreadyExists(levelNumberControl: any, index: number): boolean {
    for (let i = 0; i < this.feelingLevelControlArray.length; i++) {
      if (levelNumberControl.value == this.feelingLevelControlArray.at(i).get('levelNumber').value && index != i
        && levelNumberControl.value != null && this.feelingLevelControlArray.at(i).get('levelNumber').value != null) {
        return true
      }
    }
    return false
  }

  // =============================================================GET FUNCTIONS==========================================================================
  // Get all feelings by limit and offset.
  getAllFeelings() {
    this.spinnerService.loadingMessage = "Getting All Feelings"


    this.searchFormValue.limit = this.limit
    this.searchFormValue.offset = this.offset
    this.feelingService.getAllFeelings(this.searchFormValue).subscribe((response) => {
      this.totalFeelings = response.headers.get('X-Total-Count')
      this.feelingList = response.body
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
