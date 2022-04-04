import { DatePipe } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { AdminService, IEvent } from 'src/app/service/admin/admin.service';
import { UrlConstant } from 'src/app/service/constant';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { GeneralService, ICountry, ISearchFilterField, IState } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-admin-event',
  templateUrl: './admin-event.component.html',
  styleUrls: ['./admin-event.component.css']
})
export class AdminEventComponent implements OnInit {

  // Flags.
  isSearched: boolean
  isViewClicked: boolean
  isUpdateClicked: boolean

  // components
  eventStatusList: any[]
  countryList: ICountry[]
  stateList: IState[]

  // events
  events: IEvent[]
  totalEvents: number

  // form
  eventForm: FormGroup
  eventSearchForm: FormGroup

  // Search.
  searchFormValue: any
  searchFilterFieldList: ISearchFilterField[]

  // Pagination.
  limit: number
  currentPage: number
  offset: number
  paginationString: string

  // Spinner.



  // access
  permission: IPermission

  // image upload
  docStatus: string
  displayedImageName: string
  isImageUploadedToServer: boolean
  isImageUploading

  // Modal.
  modalRef: any
  @ViewChild('eventFormModal') eventFormModal: any
  @ViewChild('deleteModal') deleteModal: any
  @ViewChild('drawer') drawer: any

  constructor(
    private formBuilder: FormBuilder,
    public utilService: UtilityService,
    private adminService: AdminService,
    private generalService: GeneralService,
    private fileOpsService: FileOperationService,
    private urlConstant: UrlConstant,
    private modalService: NgbModal,
    private spinnerService: SpinnerService,
    private datePipe: DatePipe,
    private route: ActivatedRoute,
    private router: Router,
  ) {
    this.initalizeVariables()
    this.createEventSearchForm()
    this.createEventForm()
    this.getAllComponents()
  }

  initalizeVariables(): void {
    this.permission = this.utilService.getPermission(this.urlConstant.COLLEGE_WORKSHOP)

    this.isSearched = false
    this.isViewClicked = false
    this.isUpdateClicked = false

    this.limit = 5
    this.offset = 0
    this.totalEvents = 0


    this.events = []
    this.eventStatusList = []

    this.resetImageUploadFields()
  }

  resetImageUploadFields(): void {
    this.docStatus = ""
    this.displayedImageName = "Select Image"
    this.isImageUploadedToServer = false
    this.isImageUploading = false
  }

  getAllComponents(): void {
    this.getEventStatusList()
    this.getCountryList()
    this.searchOrGetEvents()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  createEventForm(): void {
    this.eventForm = this.formBuilder.group({
      id: new FormControl(null),
      address: new FormControl(null, [Validators.required, Validators.pattern(/^[.0-9a-zA-Z\s,-\/]+$/)]),
      state: new FormControl(null, Validators.required),
      city: new FormControl(null, [Validators.required]),
      country: new FormControl(null, Validators.required),
      pinCode: new FormControl(null, [Validators.required, Validators.pattern(/^[0-9]{6}$/)]),
      title: new FormControl(null, [Validators.required]),
      description: new FormControl(null, [Validators.required]),
      entryFee: new FormControl(null, [Validators.required, Validators.min(0)]),
      fromDate: new FormControl(null, [Validators.required]),
      toDate: new FormControl(null, [Validators.required]),
      fromTime: new FormControl(null, [Validators.required]),
      toTime: new FormControl(null, [Validators.required]),
      totalHours: new FormControl(null, [Validators.required, Validators.min(1)]),
      lastRegistrationDate: new FormControl(null, [Validators.required]),
      eventStatus: new FormControl(null, [Validators.required]),
      eventMeetingLink: new FormControl(null),
      isOnline: new FormControl(true, [Validators.required]),
      isActive: new FormControl(true, [Validators.required]),
      eventImage: new FormControl(null),
    })
  }

  createEventSearchForm(): void {
    this.eventSearchForm = this.formBuilder.group({
      title: new FormControl(null),
      fromDate: new FormControl(null),
      toDate: new FormControl(null),
      lastRegistrationDate: new FormControl(null),
      isOnline: new FormControl(null),
      isActive: new FormControl(null),
    })
  }

  // Compare Obj1 and Obj2
  compareFn(objectOne: any, objectTwo: any): boolean {
    if (objectOne == null && objectTwo == null) {
      return true;
    }
    return objectOne && objectTwo ? objectOne.id === objectTwo.id : objectOne === objectTwo
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetEvents() {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getEvents()
      return
    }
    this.eventSearchForm.patchValue(queryParams)
    this.searchEvents()
  }

  searchAndCloseDrawer(): void {
    this.drawer.toggle()
    this.searchEvents()
  }

  searchEvents(): void {
    // console.log(this.searchBatchForm.value)
    this.searchFormValue = { ...this.eventSearchForm?.value }
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: this.searchFormValue,
    })
    let flag: boolean = true

    for (let field in this.searchFormValue) {
      if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
        delete this.searchFormValue[field];
      } else {
        this.isSearched = true
        flag = false
      }
    }

    // No API call on empty search.
    if (flag) {
      return
    }
    this.changePage(1)
  }

  // page change function
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber;
    this.offset = this.currentPage - 1;
    this.getEvents();
  }

  getEvents(): void {
    this.spinnerService.loadingMessage = "Getting events"

    this.totalEvents = 0
    this.events = []
    if (!this.isSearched) {
      this.searchFormValue = null
    }
    this.adminService.getEvents(this.limit, this.offset, this.searchFormValue).subscribe((response: any) => {
      this.events = response.body
      this.totalEvents = response.headers.get("X-Total-Count")
    }, (err: any) => {
      this.totalEvents = 0
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

  onAddClick(): void {
    this.isViewClicked = false
    this.isUpdateClicked = false
    this.createEventForm()
    this.openModal(this.eventFormModal, "xl")
  }

  convertDateToISOString(event: IEvent): void {
    event.toDate = new Date(event.toDate).toISOString()
    event.fromDate = new Date(event.fromDate).toISOString()
    event.lastRegistrationDate = new Date(event.lastRegistrationDate).toISOString()
  }

  addEvent(): void {
    this.spinnerService.loadingMessage = "Adding Event"


    let event: IEvent = this.eventForm.value
    this.convertDateToISOString(event)
    this.adminService.addEvent(event).subscribe((response: any) => {
      this.modalRef.close()
      alert("Event successfully added")
      this.getEvents()
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  onViewClick(event: IEvent): void {
    this.isViewClicked = true
    this.isUpdateClicked = false
    this.createEventForm()

    event.fromDate = this.datePipe.transform(event.fromDate, 'yyyy-MM-dd')
    event.toDate = this.datePipe.transform(event.toDate, 'yyyy-MM-dd')
    event.lastRegistrationDate = this.datePipe.transform(event.lastRegistrationDate, 'yyyy-MM-dd')

    this.eventForm.patchValue(event)
    this.eventForm.disable()
    this.openModal(this.eventFormModal, "xl")
  }

  onUpdateClick(): void {
    this.isViewClicked = false
    this.isUpdateClicked = true

    if (this.eventForm.get("eventImage").value) {
      this.displayedImageName = `<a href=${this.eventForm.get("eventImage").value} target="_blank">Event Image</a>`
    }
    this.eventForm.enable()
  }

  updateEvent(): void {
    this.spinnerService.loadingMessage = "Updating Event"


    let event: IEvent = this.eventForm.value
    this.convertDateToISOString(event)

    this.adminService.updateEvent(event).subscribe((response: any) => {
      this.modalRef.close()
      alert("Event successfully updated")
      this.getEvents()
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  onDeleteClick(eventID: string): void {
    this.openModal(this.deleteModal, 'md').result.then(() => {
      this.deleteEvent(eventID)
    }, (err) => {
      console.error(err);
      return
    })
  }

  deleteEvent(eventID: string): void {
    this.spinnerService.loadingMessage = "Deleting Event"

    this.adminService.deleteEvent(eventID).subscribe((response: any) => {
      this.modalRef.close()
      alert("Event successfully deleted")
      this.getEvents()
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  getEventStatusList(): void {
    this.generalService.getGeneralTypeByType("event_status").subscribe((response: any) => {
      this.eventStatusList = response
    }, (err: any) => {
      console.error(err)
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

  getStateList(countryid: any): void {
    this.eventForm.get('state')?.reset()
    this.eventForm.get('state')?.disable()
    this.stateList = [] as IState[]
    this.generalService.getStatesByCountryID(countryid).subscribe((respond: any[]) => {
      this.stateList = respond
      if (this.stateList.length > 0 && !this.isViewClicked) {
        this.eventForm.get('state')?.enable()
      }
    }, (err) => {
      console.error(this.utilService.getErrorString(err))
    })
  }

  //On uplaoding brochure
  onResourceSelect(event: any) {
    this.docStatus = ""
    let files = event.target.files
    if (files && files.length) {
      let file = files[0]
      let err = this.fileOpsService.isImageFileValid(file)
      if (err != null) {
        this.docStatus = `<p><span>&#10060;</span> ${err}</p>`
        return
      }
      // console.log(file)
      // Upload brochure if it is present.]
      this.isImageUploading = true
      this.fileOpsService.uploadFile(file, this.fileOpsService.EVENT_IMAGE).subscribe((data: any) => {
        this.eventForm.markAsDirty()
        this.eventForm.patchValue({
          eventImage: data
        })
        this.displayedImageName = file.name
        this.isImageUploading = false
        this.isImageUploadedToServer = true
        this.docStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
      }, (error) => {
        this.isImageUploading = false
        this.docStatus = `<p><span>&#10060;</span> ${error}</p>`
      })
    }
  }

  //validate add/update form
  validate(): void {
    if (this.isImageUploading) {
      alert("Please wait till file is being uploaded")
      return
    }

    // console.log(this.eventForm.controls);
    if (this.eventForm.invalid) {
      this.eventForm.markAllAsTouched();
      return
    }

    if (this.isUpdateClicked) {
      this.updateEvent()
      return
    }
    this.addEvent()
  }

  setPaginationString() {
    this.paginationString = ''
    let start: number = this.limit * this.offset + 1
    let end: number = +this.limit + this.limit * this.offset
    if (this.totalEvents < end) {
      end = this.totalEvents
    }
    if (this.totalEvents == 0) {
      this.paginationString = ''
      return
    }
    this.paginationString = `${start} - ${end}`
  }

  openModal(modalContent: any, modalSize?: string): NgbModalRef {
    this.resetImageUploadFields()
    if (!modalSize) {
      modalSize = 'lg'
    }
    let options: NgbModalOptions = {
      ariaLabelledBy: 'modal-basic-title', keyboard: false,
      backdrop: 'static', size: modalSize, centered: true,
    }
    this.modalRef = this.modalService.open(modalContent, options)
    return this.modalRef
  }

  // Used to dismiss modal.
  dismissFormModal(modal: NgbModalRef) {
    if (this.isImageUploading) {
      alert("Please wait till file is being uploaded")
      return
    }
    if (this.isImageUploadedToServer) {
      if (!confirm("Uploaded brochure will be deleted.\nAre you sure you want to close?")) {
        return
      }
    }
    modal.dismiss()
    this.resetImageUploadFields()
  }

  resetSearchAndGetAll(): void {
    this.resetSearchForm()
    this.searchFormValue = null
    this.changePage(1)
    this.isSearched = false
    this.router.navigate([this.urlConstant.COLLEGE_WORKSHOP])
  }

  // Reset search form.
  resetSearchForm(): void {
    this.eventSearchForm.reset()
  }
}
