import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, FormControl, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalRef, NgbModalOptions } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { GeneralService } from 'src/app/service/general/general.service';
import { ISearchFilterField, ISpeaker, ISpeakerDTO, SpeakerService } from 'src/app/service/speaker/speaker.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-admin-speaker',
  templateUrl: './admin-speaker.component.html',
  styleUrls: ['./admin-speaker.component.css']
})
export class AdminSpeakerComponent implements OnInit {

  // Components.
  designationList: any[]

  // Flags.
  isSearched: boolean
  isOperationUpdate: boolean
  isViewMode: boolean

  // Speaker.
  speakerList: ISpeaker[]
  speakerForm: FormGroup

  // Pagination.
  limit: number
  currentPage: number
  totalSpeakers: number
  offset: number
  paginationStart: number
  paginationEnd: number

  // Modal.
  modalRef: any
  @ViewChild('speakerFormModal') speakerFormModal: any
  @ViewChild('deleteSpeakerModal') deleteSpeakerModal: any

  // Spinner.



  // Search.
  speakerSearchForm: FormGroup
  searchFormValue: any
  searchFilterFieldList: ISearchFilterField[]

  constructor(
    private formBuilder: FormBuilder,
    private speakerService: SpeakerService,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private router: Router,
    private route: ActivatedRoute,
    private utilService: UtilityService,
    private generalService: GeneralService,
  ) {
    this.initializeVariables()
    this.getDesignationList()
    this.searchOrGetSpeakers()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize all global variables.
  initializeVariables() {

    // Components.
    this.designationList = []

    // Speaker.
    this.speakerList = [] as ISpeaker[]

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
    this.createSpeakerForm()
    this.createSpeakerSearchForm()

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Speakers"


    // Search.
    this.searchFilterFieldList = []
  }

  // =============================================================CREATE FORMS==========================================================================
  // Create speaker form.
  createSpeakerForm(): void {
    this.speakerForm = this.formBuilder.group({
      id: new FormControl(null),
      firstName: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z]*$/), Validators.maxLength(50)]),
      lastName: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z]*$/), Validators.maxLength(50)]),
      company: new FormControl(null, [Validators.maxLength(200)]),
      experienceInYears: new FormControl(null, [Validators.max(60), Validators.min(0)]),
      designation: new FormControl(null),
    })
  }

  // Create speaker search form.
  createSpeakerSearchForm(): void {
    this.speakerSearchForm = this.formBuilder.group({
      firstName: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z]*$/), Validators.maxLength(50)]),
      lastName: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z]*$/), Validators.maxLength(50)]),
      company: new FormControl(null, [Validators.maxLength(200)]),
      minimumExperience: new FormControl(null, [Validators.max(60), Validators.min(0)]),
      maximumExperience: new FormControl(null, [Validators.max(60), Validators.min(0)]),
      designationID: new FormControl(null),
    })
  }
  // =============================================================SPEAKER CRUD FUNCTIONS==========================================================================
  // On clicking add new speaker button.
  onAddNewSpeakerClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.createSpeakerForm()
    this.openModal(this.speakerFormModal, 'xl')
  }

  // Add new speaker.
  addSpeaker(): void {
    this.spinnerService.loadingMessage = "Adding Speaker"


    let speaker: ISpeaker = this.speakerForm.value
    this.patchIDFromObjects(speaker)
    this.speakerService.addSpeaker(speaker).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllSpeakers()
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

  // On clicking view speaker button.
  onViewSpeakerClick(speaker: ISpeakerDTO): void {
    this.isViewMode = true
    this.createSpeakerForm()
    this.speakerForm.patchValue(speaker)
    this.speakerForm.disable()
    this.openModal(this.speakerFormModal, 'xl')
  }

  // On cliking update form button in speaker form.
  onUpdateSpeakerClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = true
    this.speakerForm.enable()
  }

  // Update speaker.
  updateSpeaker(): void {
    this.spinnerService.loadingMessage = "Updating Speaker"


    let speaker: ISpeaker = this.speakerForm.value
    this.patchIDFromObjects(speaker)
    this.speakerService.updateSpeaker(speaker).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllSpeakers()
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

  // On clicking delete speaker button. 
  onDeleteSpeakerClick(speakerID: string): void {
    this.openModal(this.deleteSpeakerModal, 'md').result.then(() => {
      this.deleteSpeaker(speakerID)
    }, (err) => {
      console.error(err)
      return
    })
  }

  // Delete speaker after confirmation from user.
  deleteSpeaker(speakerID: string): void {
    this.spinnerService.loadingMessage = "Deleting Speaker"


    this.speakerService.deleteSpeaker(speakerID).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllSpeakers()
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

  // =============================================================SPEAKER SEARCH FUNCTIONS==========================================================================
  // Reset search form and renaviagte page.
  resetSearchAndGetAll(): void {
    this.searchFilterFieldList = []
    this.speakerSearchForm.reset()
    this.searchFormValue = null
    this.changePage(1)
    this.isSearched = false
    this.router.navigate(['/admin/models/speaker'])
  }

  // Reset search form.
  resetSearchForm(): void {
    this.searchFilterFieldList = []
    this.speakerSearchForm.reset()
  }

  // Search speakers.
  searchSpeakers(): void {
    this.searchFormValue = { ...this.speakerSearchForm?.value }
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
    this.spinnerService.loadingMessage = "Searching Speakers"
    this.changePage(1)
  }

  // Delete search criteria from speaker search form by search name.
  deleteSearchCriteria(searchName: string): void {
    this.speakerSearchForm.get(searchName).setValue(null)
    this.searchSpeakers()
  }

  // ================================================OTHER FUNCTIONS FOR SPEAKER===============================================
  // Page change.
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
    this.getAllSpeakers()
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetSpeakers() {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getAllSpeakers()
      return
    }
    this.speakerSearchForm.patchValue(queryParams)
    this.searchSpeakers()
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

  // On clicking sumbit button in speaker form.
  onFormSubmit(): void {
    if (this.speakerForm.invalid) {
      this.speakerForm.markAllAsTouched()
      return
    }
    if (this.isOperationUpdate) {
      this.updateSpeaker()
      return
    }
    this.addSpeaker()
  }

  // Set total speakers list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalSpeakers < this.paginationEnd) {
      this.paginationEnd = this.totalSpeakers
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

  // Extract ID from objects and delete objects before adding or updating.
  patchIDFromObjects(speaker: ISpeaker): void {
    if (this.speakerForm.get('designation').value) {
      speaker.designationID = this.speakerForm.get('designation').value.id
      delete speaker['designation']
    }
  }

  // =============================================================GET FUNCTIONS==========================================================================
  // Get all speakers.
  getAllSpeakers() {
    this.spinnerService.loadingMessage = "Getting All Speakers"


    this.speakerService.getAllSpeakers(this.limit, this.offset, this.searchFormValue).subscribe((response) => {
      this.totalSpeakers = response.headers.get('X-Total-Count')
      this.speakerList = response.body
      console.log(this.speakerList)
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.setPaginationString()
    })
  }

  // Get designation list.
  getDesignationList(): void {
    this.generalService.getDesignations().subscribe((respond: any[]) => {
      this.designationList = respond
    }, (err) => {
      console.error(this.utilService.getErrorString(err))
    })
  }

}
