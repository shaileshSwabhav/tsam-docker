import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, FormControl, Validators } from '@angular/forms';
import { GeneralService } from 'src/app/service/general/general.service';
import { CompanyService, ICompanyEnquiry, ITechnology, IDomain, ICountry, IState, ICompanyEnquiryDTO, ICallRecordDTO, IPurpose, ICallRecord } from 'src/app/service/company/company.service';
import { Constant, Role, UrlConstant } from 'src/app/service/constant';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { ActivatedRoute, Router } from '@angular/router';
import { DatePipe } from '@angular/common';
import { ISearchSection } from 'src/app/service/talent/talent.service';
import { TechnologyService } from 'src/app/service/technology/technology.service';

@Component({
      selector: 'app-company-enquiry',
      templateUrl: './company-enquiry.component.html',
      styleUrls: ['./company-enquiry.component.css']
})
export class CompanyEnquiryComponent implements OnInit {

      //****************************COMPANY ENQUIRY*************************************** */

      // Components.
      domainList: IDomain[]
      technologyList: ITechnology[]
      salesPersonList: any[]
      stateList: IState[]
      stateSearchList: any[]
      countryList: ICountry[]
      enquirySourceList: any[]
      enquiryTypeList: any[]

      // technology component
      isTechLoading: boolean
      techLimit: number
      techOffset: number

      // Flags.
      isViewMode: boolean
      isOperationEnquiryUpdate: boolean

      // Company Enquiry.
      companyEnquiryForm: FormGroup
      companyEnquiryList: ICompanyEnquiryDTO[]
      selectedEnquiryID: string
      selectedEnquiriesList: any[]

      // Pagination.
      limit: number
      offset: number
      currentPage: number
      totalCompanyEnquiries: number
      paginationString: string

      // Modal.
      modalRef: any
      @ViewChild('companyEnquiryFormModal') companyEnquiryFormModal: any
      @ViewChild('deleteCompanyEnquiryModal') deleteCompanyEnquiryModal: any

      // Spinner.



      // Search.
      companyEnquirySearchForm: FormGroup
      searchFormValue: any
      isSearched: boolean
      showSearch: boolean
      searchSectionList: ISearchSection[]
      selectedSectionName: string

      // Permission.
      permission: IPermission
      roleName: string
      showForAdmin: boolean
      showForSalesPerson: boolean

      // Constants
      nilUUID: string

      //****************************CALL RECORDS*************************************** */
      // Components.
      purposeList: any[]
      outcomeList: any[]
      outcomeSearchList: any[]

      // Flags.
      showCallRecordForm: boolean
      isOperationCallRecordUpdate: boolean

      // Enquiry call record.
      callRecordList: ICallRecordDTO[]
      callRecordForm: FormGroup

      // Modal.
      @ViewChild('callRecordFormModal') callRecordFormModal: any

      //****************************ALLOCATION*************************************** */
      // Flags.
      showAllocateSalespersonToOneEnquiry: boolean
      showAllocateSalespersonToEnquiries: boolean
      multipleSelect: boolean

      // Modal.
      @ViewChild('allocateSalespersonModal') allocateSalespersonModal: any
      @ViewChild('drawer') drawer: any

      //constant
      private readonly IGNORE_SEARCH_FIELD: string[] = ["limit", "offset"]

      constructor(
            private formBuilder: FormBuilder,
            private companyService: CompanyService,
            private generalService: GeneralService,
            private techService: TechnologyService,
            private constant: Constant,
            private utilService: UtilityService,
            private localService: LocalService,
            private urlConstant: UrlConstant,
            private spinnerService: SpinnerService,
            private router: Router,
            private activatedRoute: ActivatedRoute,
            private role: Role,
            private datePipe: DatePipe,
            private route: ActivatedRoute,
            private modalService: NgbModal,

      ) {
            this.initializeVariables()
            this.getAllComponents()
      }


      get ongoingOperations() {
            return this.spinnerService.ongoingOperations
      }

      ngOnInit() {
      }

      initializeVariables(): void {
            //****************************COMPANY ENQUIRY*************************************** */
            // Components.
            this.domainList = []
            this.technologyList = []
            this.salesPersonList = []
            this.stateList = []
            this.stateSearchList = []
            this.countryList = []
            this.enquirySourceList = []
            this.enquiryTypeList = []

            // Flags.
            this.isViewMode = false
            this.isOperationEnquiryUpdate = false
            this.isTechLoading = false

            // Company Enquiry.
            this.enquiryTypeList = [] as ICompanyEnquiryDTO[]
            this.selectedEnquiriesList = []

            // Paginate.
            this.limit = 5
            this.offset = 0
            this.currentPage = 0
            this.techLimit = 10
            this.techOffset = 0

            // Spinner.
            this.spinnerService.loadingMessage = "Getting company enquiries"


            // Search.
            this.isSearched = false
            this.showSearch = false

            // Permision.
            this.showForAdmin = false
            this.showForSalesPerson = true
            // Get permissions from menus using utilityService function.
            this.permission = this.utilService.getPermission(this.urlConstant.COMPANY_ENQUIRY)
            // Get role name for menu for calling their specific apis.
            this.roleName = this.localService.getJsonValue("roleName")
            // If admin is logged in then show its features.
            if (this.roleName?.toLocaleLowerCase() == this.role.ADMIN.toLocaleLowerCase()) {
                  this.showForAdmin = true
            }
            // Hide features for salesperson.
            if (this.roleName == this.role.SALES_PERSON) {
                  this.showForSalesPerson = false
            }

            // Constants.
            this.nilUUID = this.constant.NIL_UUID

            //****************************CALL RECORDS*************************************** */
            // Components.
            this.purposeList = []
            this.outcomeList = []
            this.outcomeSearchList = []

            // Flags.
            this.showCallRecordForm = false
            this.isOperationCallRecordUpdate = false

            // Enquiry call record.
            this.callRecordList = [] as ICallRecordDTO[]

            //****************************ALLOCATION*************************************** */
            // Flags.
            this.showAllocateSalespersonToOneEnquiry = false
            this.showAllocateSalespersonToEnquiries = false

            //****************************INITIALIZE FORMS*************************************** */
            this.createSearchCompanyEnquiryForm()

            this.setSearchSectionFields()
      }

      setSearchSectionFields(): void {
            this.searchSectionList = [
                  {
                        name: "Company",
                        isSelected: true
                  },
                  {
                        name: "Enquiry",
                        isSelected: false
                  },
                  {
                        name: "Location",
                        isSelected: false
                  },
                  {
                        name: "Calling",
                        isSelected: false
                  },
                  {
                        name: "Other",
                        isSelected: false
                  },
            ]
            this.selectedSectionName = "Company"
      }

      // On clicking search section name.
      onSearchSectionNameClick(sectionName: string): void {
            for (let i = 0; i < this.searchSectionList.length; i++) {
                  if (this.searchSectionList[i].name == sectionName) {
                        this.searchSectionList[i].isSelected = true
                        this.selectedSectionName = this.searchSectionList[i].name
                  } else {
                        this.searchSectionList[i].isSelected = false
                  }
            }
      }


      //*********************************************CREATE FORMS************************************************************
      // Create company enquiry form.
      createCompanyEnquiryForm(): void {
            this.companyEnquiryForm = this.formBuilder.group({
                  id: new FormControl(null),
                  code: new FormControl(null),
                  address: new FormControl(null, [Validators.required, Validators.pattern(/^[.0-9a-zA-Z\s,-\/]+$/)]),
                  pinCode: new FormControl(null, [Validators.required, Validators.pattern(/^[1-9][0-9]{5}$/)]),
                  country: new FormControl(null, [Validators.required]),
                  state: new FormControl(null, [Validators.required]),
                  city: new FormControl(null, [Validators.required, Validators.maxLength(50), Validators.pattern(/^[a-zA-Z]+([a-zA-Z ]?)+$/)]),
                  email: new FormControl(null, [Validators.pattern(/^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$/), Validators.maxLength(100)]),
                  hrName: new FormControl(null, [Validators.maxLength(100), Validators.pattern(/^[a-zA-Z ]*$/)]),
                  hrContact: new FormControl(null, [Validators.pattern(/^[6789]\d{9}$/)]),
                  founderName: new FormControl(null, [Validators.maxLength(100), Validators.pattern(/^[a-zA-Z ]*$/)]),
                  vacancy: new FormControl(null, [Validators.min(0), Validators.max(10000)]),
                  enquiryDate: new FormControl(null),
                  enquiryType: new FormControl(null),
                  enquirySource: new FormControl(null),
                  jobRole: new FormControl(null, [Validators.maxLength(100)]),
                  packageOffered: new FormControl(null, [Validators.min(0), Validators.max(9999999999)]),
                  subject: new FormControl(null, [Validators.required, Validators.maxLength(500)]),
                  message: new FormControl(null, [Validators.required, Validators.maxLength(4000)]),
                  companyName: new FormControl(null, [Validators.required, Validators.maxLength(200)]),
                  domains: new FormControl(null, [Validators.required]),
                  technologies: new FormControl(null, [Validators.required]),
                  website: new FormControl(null, [Validators.pattern('(https?://)?([\\da-z.-]+)\\.([a-z.]{2,6})[/\\w .-]*/?')]),
                  salesPerson: new FormControl(null),
                  companyBranch: new FormControl(null),
            })
      }

      // Create company enquiry search form.
      createSearchCompanyEnquiryForm(): void {
            this.companyEnquirySearchForm = this.formBuilder.group({
                  companyName: new FormControl(null, [Validators.maxLength(100)]),
                  stateID: new FormControl(null),
                  countryID: new FormControl(null),
                  salesPersonID: new FormControl(null),
                  enquiryType: new FormControl(null),
                  enquirySource: new FormControl(null),
                  fromDate: new FormControl(null),
                  tillDate: new FormControl(null),
                  city: new FormControl(null, [Validators.maxLength(50)]),
                  domainIDs: new FormControl(null),
                  technologyIDs: new FormControl(null),
                  purposeID: new FormControl(null),
                  outcomeID: new FormControl(null),
                  limit: new FormControl(this.limit),
                  offset: new FormControl(this.offset),
            })
            this.disableSearchFormFields()
      }

      // Create new enquiry call record form.
      createCallRecordForm(): void {
            this.callRecordForm = this.formBuilder.group({
                  id: new FormControl(),
                  dateTime: new FormControl(null, [Validators.required]),
                  purpose: new FormControl(null, [Validators.required]),
                  outcome: new FormControl(null, [Validators.required]),
                  comment: new FormControl(null, [Validators.maxLength(500)]),
                  enquiryID: new FormControl(null)
            })
      }

      //*********************************************ADD FOR COMPANY ENQUIRY FUNCTIONS************************************************************
      // On add new company enquiry button click.
      onAddCompanyEnquiryClick(): void {
            this.isViewMode = false
            this.isOperationEnquiryUpdate = false
            this.createCompanyEnquiryForm()
            this.stateList = []
            this.companyEnquiryForm.get('state').disable()
            this.openModal(this.companyEnquiryFormModal, 'xl')
      }

      // Add new company enquiry.
      addCompanyEnquiry(): void {
            this.spinnerService.loadingMessage = "Adding Company Enquiry"


            let enquiry: ICompanyEnquiry = this.companyEnquiryForm.value
            this.patchIDFromObjectsForCompanyEnquiry(enquiry)
            this.companyService.addCompanyEnquiry(enquiry).subscribe((response: any) => {
                  this.modalRef.close('success')
                  this.getAllEnquiries()
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

      //*********************************************UPDATE AND VIEW FOR COMPANY ENQUIRY FUNCTIONS************************************************************
      // On clicking view company enquiry button.
      onViewCompanyEnquiryClick(enquiry: ICompanyEnquiryDTO): void {
            this.isViewMode = true
            this.createCompanyEnquiryForm()

            // State.
            if (enquiry.country != undefined) {
                  this.getStateListByCountry(enquiry.country)
            }

            this.companyEnquiryForm.patchValue(enquiry)
            this.companyEnquiryForm.disable()
            this.openModal(this.companyEnquiryFormModal, 'xl')
      }

      // On cliking update form button in company enquiry form.
      onUpdateCompanyEnquiryClick(): void {
            this.isViewMode = false
            this.isOperationEnquiryUpdate = true
            this.enableCompanyEnquiryForm()
      }

      // Update company enquiry.
      updateCompanyEnquiry(): void {
            this.spinnerService.loadingMessage = "Updating Company Enquiry"


            let enquiry: ICompanyEnquiry = this.companyEnquiryForm.value
            this.patchIDFromObjectsForCompanyEnquiry(enquiry)
            this.companyService.updateCompanyEnquiry(this.companyEnquiryForm.value).subscribe((response: any) => {
                  this.modalRef.close()
                  this.getAllEnquiries()
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

      //*********************************************DELETE FOR COMPANY ENQUIRY FUNCTIONS************************************************************

      // On clicking delete company enquiry button. 
      onDeleteCompanyEnquiryClick(enquiryID: string): void {
            this.openModal(this.deleteCompanyEnquiryModal, 'md').result.then(() => {
                  this.deleteCompanyEnquiry(enquiryID)
            }, (err) => {
                  console.error(err)
                  return
            })
      }

      // Delete company enquiry after confirmation from user.
      deleteCompanyEnquiry(enquiryID: string): void {
            this.spinnerService.loadingMessage = "Deleting Company Enquiry"

            this.modalRef.close()

            this.companyService.deleteCompanyEnquiry(enquiryID).subscribe((response: any) => {
                  this.getAllEnquiries()
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

      // =============================================================COMPANY ENQUIRY SEARCH FUNCTIONS==========================================================================
      // Reset search form and renaviagte page.
      resetSearchAndGetAll(): void {
            this.companyEnquirySearchForm.reset()
            this.searchFormValue = null
            this.changePage(1)
            this.isSearched = false
            this.showSearch = false
            this.router.navigate([this.urlConstant.COMPANY_ENQUIRY])
      }

      // Reset search form.
      resetSearchForm(): void {
            this.limit = this.companyEnquirySearchForm.get("limit").value
            this.offset = this.companyEnquirySearchForm.get("offset").value
            this.companyEnquirySearchForm.reset({
                  limit: this.limit,
                  offset: this.offset,
            })
      }

      searchAndCloseDrawer(): void {
            this.drawer.toggle()
            this.searchCompanyEnquiries()
      }

      // Search company enquiries.
      searchCompanyEnquiries(): void {
            this.searchFormValue = { ...this.companyEnquirySearchForm?.value }
            this.router.navigate([], {
                  relativeTo: this.route,
                  queryParams: this.searchFormValue,
            })
            let flag: boolean = true
            for (let field in this.searchFormValue) {
                  if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
                        delete this.searchFormValue[field]
                  } else {
                        if (!this.IGNORE_SEARCH_FIELD.includes(field)) {
                              this.isSearched = true
                        }
                        flag = false
                  }
            }
            // if (!this.isSearched) {
            //       return
            // }
            if (flag) {
                  return
            }
            this.spinnerService.loadingMessage = "Searching Company Enquiries"
            // this.changePage(1)
            this.getAllEnquiries()
      }

      // Disable dependent field controls on search form creation.
      disableSearchFormFields(): void {
            this.companyEnquirySearchForm.get('stateID').disable()
            this.companyEnquirySearchForm.get('outcomeID').disable()
      }

      // On country value change in search form.
      onCountryChangeInSearchForm(countryID: string): void {
            if (countryID) {
                  console.log("will do something")
            }
      }

      // ================================================OTHER FUNCTIONS FOR COMPNAY ENQUIRY===============================================
      // Page change.
      changePage(pageNumber: number): void {
            // 
            // this.currentPage = pageNumber
            // this.offset = this.currentPage - 1
            this.companyEnquirySearchForm.get("offset").setValue(pageNumber - 1)

            this.limit = this.companyEnquirySearchForm.get("limit").value
            this.offset = this.companyEnquirySearchForm.get("offset").value
            // this.searchCompany()
            this.searchCompanyEnquiries()
      }

      // Checks the url's query params and decides whether to call get or search.
      searchOrGetCompanyEnquiries(): void {
            let queryParams = this.route.snapshot.queryParams
            if (this.utilService.isObjectEmpty(queryParams)) {
                  this.getAllEnquiries()
                  return
            }
            this.companyEnquirySearchForm.patchValue(queryParams)
            this.searchCompanyEnquiries()
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

      // On clicking sumbit button in company enquiry form.
      onFormSubmit(): void {

            console.log(this.companyEnquiryForm.controls);

            if (this.companyEnquiryForm.invalid) {
                  this.companyEnquiryForm.markAllAsTouched()
                  return
            }

            this.companyEnquiryForm.get("enquiryDate")?.setValue(
                  this.datePipe.transform(this.companyEnquiryForm.get("enquiryDate")?.value, 'yyyy-MM-dd'))

            if (this.isOperationEnquiryUpdate) {
                  this.updateCompanyEnquiry()
                  return
            }
            this.addCompanyEnquiry()
      }

      // Set total company enquiries list on current page.
      setPaginationString(): void {
            this.paginationString = ''
            let limit = this.companyEnquirySearchForm.get('limit').value
            let offset = this.companyEnquirySearchForm.get('offset').value
            let start: number = limit * offset + 1
            let end: number = +limit + limit * offset
            if (this.totalCompanyEnquiries < end) {
                  end = this.totalCompanyEnquiries
            }
            if (this.totalCompanyEnquiries == 0) {
                  this.paginationString = ''
                  return
            }
            this.paginationString = `${start}-${end}`
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

      // Format the fees in indian rupee system.
      formatFeesInIndianRupeeSystem(fees: number): string {
            var result = fees.toString().split('.')
            var lastThree = result[0].substring(result[0].length - 3)
            var otherNumbers = result[0].substring(0, result[0].length - 3)
            if (otherNumbers != '')
                  lastThree = ',' + lastThree
            var output = otherNumbers.replace(/\B(?=(\d{2})+(?!\d))/g, ",") + lastThree
            if (result.length > 1) {
                  output += "." + result[1]
            }
            return output
      }

      // Extract id from objects and give it to company enquiry.
      patchIDFromObjectsForCompanyEnquiry(enquiry: ICompanyEnquiry): void {
            if (this.companyEnquiryForm.get('salesPerson').value) {
                  enquiry.salesPersonID = this.companyEnquiryForm.get('salesPerson').value.id
                  delete enquiry['salesPerson']
            }
      }

      // Enable the company enquiry form.
      enableCompanyEnquiryForm(): void {
            this.companyEnquiryForm.enable()
            this.companyEnquiryForm.get('code').disable()
      }

      //*********************************************CRUD FUNCTIONS FOR COMPANY ENQUIRY CALL RECORDS************************************************************
      // On clicking get all call records by enquiry id.
      getCallRecordsForSelectedEnquiry(enquiryID: any): void {
            this.isOperationCallRecordUpdate = false
            this.callRecordList = []
            this.selectedEnquiryID = enquiryID
            this.getAllCallRecords()
            this.showCallRecordForm = false
            this.openModal(this.callRecordFormModal, 'xl')
      }

      // Get all call records by enquiry id.
      getAllCallRecords(): void {
            this.spinnerService.loadingMessage = "Getting all call records"


            this.companyService.getCallRecordsByEnquiry(this.selectedEnquiryID).subscribe(response => {
                  this.callRecordList = response
                  this.formatDateTimeOfEnquiryCallRecords()
            }, err => {
                  console.error(this.utilService.getErrorString(err))
            })
      }

      // On clicking add new call record button.
      onAddNewCallRecordButtonClick(): void {
            this.isOperationCallRecordUpdate = false
            this.outcomeList = []
            this.createCallRecordForm()
            this.showCallRecordForm = true
            this.callRecordForm.get('outcome').disable()
      }

      // Add call record.
      addCallRecord(): void {
            this.spinnerService.loadingMessage = "Adding Call Record"


            let callRecord: ICallRecord = this.callRecordForm.value
            this.patchIDFromObjectsForCallRecord(callRecord)
            this.companyService.addCallRecord(this.callRecordForm.value, this.selectedEnquiryID).subscribe((response: any) => {
                  this.showCallRecordForm = false
                  this.getAllCallRecords()
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

      // On clicking update call record button.
      OnUpdateCallRecordButtonClick(index: number): void {
            this.isOperationCallRecordUpdate = true
            this.createCallRecordForm()
            this.showCallRecordForm = true
            if (this.callRecordList[index].purpose) {
                  this.getOutcomesByPurpose(this.callRecordList[index].purpose)
            }
            this.callRecordForm.patchValue(this.callRecordList[index])
      }

      // Update call record.
      updateCallrecord(): void {
            this.spinnerService.loadingMessage = "Updating Call Record"


            let callRecord: ICallRecord = this.callRecordForm.value
            this.patchIDFromObjectsForCallRecord(callRecord)
            this.companyService.updateCallRecord(this.callRecordForm.value, this.selectedEnquiryID).subscribe((response: any) => {
                  this.showCallRecordForm = false
                  this.getAllCallRecords()
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

      // Delete call record.
      deleteCallRecord(callRecordID: string): void {
            if (confirm("Are you sure you want to delete the call record?")) {
                  this.spinnerService.loadingMessage = "Deleting Call Record"


                  this.companyService.deleteCallRecord(callRecordID, this.selectedEnquiryID).subscribe((response: any) => {
                        this.getAllCallRecords()
                        alert(response)
                  }, (error) => {
                        console.error(error)
                        if (error.error) {
                              alert(error.error)
                              return
                        }
                        alert(error.statusText)
                  }).add(() => {

                  })
            }
      }

      // Validate call record form.
      validateCallRecordForm(): void {
            if (this.callRecordForm.invalid) {
                  this.callRecordForm.markAllAsTouched()
                  return
            }
            if (this.isOperationCallRecordUpdate) {
                  this.updateCallrecord()
                  return
            }
            this.addCallRecord()
      }

      // Format date time field of enquiry call records by removing timestamp.
      formatDateTimeOfEnquiryCallRecords(): void {
            for (let i = 0; i < this.callRecordList?.length; i++) {
                  let dateTime = this.callRecordList[i].dateTime
                  if (dateTime) {
                        this.callRecordList[i].dateTime = this.datePipe.transform(dateTime, 'yyyy-MM-ddTHH:mm:ss')
                  }
            }
      }

      // Extract ID from objects in call record form.
      patchIDFromObjectsForCallRecord(callRecord: ICallRecord): void {
            if (this.callRecordForm.get('purpose').value) {
                  callRecord.purposeID = this.callRecordForm.get('purpose').value.id
                  delete callRecord['purpose']
            }
            if (this.callRecordForm.get('outcome').value) {
                  callRecord.outcomeID = this.callRecordForm.get('outcome').value.id
                  delete callRecord['outcome']
            }
      }

      //*********************************************ALLOCATION FUNCTIONS************************************************************
      // On clicking allocate salesperson to one enquiry button click.
      onAllocateSalespersonToOneEnquiryClick(enquiryID: any): void {
            this.showAllocateSalespersonToEnquiries = false
            this.showAllocateSalespersonToOneEnquiry = true
            this.selectedEnquiryID = enquiryID
            this.openModal(this.allocateSalespersonModal, 'sm')
      }

      // On clicking allocate salesperson to multiple enquiries button click.
      onAllocateSalespersonToEnquiriesClick(): void {
            this.showAllocateSalespersonToEnquiries = true
            this.showAllocateSalespersonToOneEnquiry = false
            this.openModal(this.allocateSalespersonModal, 'sm')
      }

      // Toggle visibility of multiple select checkbox.
      toggleMultipleSelect(): void {
            if (this.multipleSelect) {
                  this.multipleSelect = false
                  this.setSelectAllEnquiries(this.multipleSelect)
                  return
            }
            this.multipleSelect = true
      }

      // Set isChecked field of all selected enquiries.
      setSelectAllEnquiries(isSelectedAll: boolean): void {
            for (let i = 0; i < this.companyEnquiryList.length; i++) {
                  this.addenquiryToList(isSelectedAll, this.companyEnquiryList[i])
            }
      }

      // Check if all enquiries in selected enquiries are added in multiple select or not.
      checkEnquiriesAdded(): boolean {
            let count: number = 0

            for (let i = 0; i < this.companyEnquiryList.length; i++) {
                  if (this.selectedEnquiriesList.includes(this.companyEnquiryList[i].id))
                        count = count + 1
            }
            return (count == this.companyEnquiryList.length)
      }

      // Check if enquiry is added in multiple select or not.
      checkEnquiryAdded(enquiryID): boolean {
            return this.selectedEnquiriesList.includes(enquiryID)
      }

      // Takes a list called selectedEnquiriesList and adds all the checked enquiries to list, also does not contain duplicate values.
      addenquiryToList(isChecked: boolean, enquiry: ICompanyEnquiryDTO): void {
            if (isChecked) {
                  if (!this.selectedEnquiriesList.includes(enquiry.id)) {
                        this.selectedEnquiriesList.push(enquiry.id)
                  }
                  return
            }
            if (this.selectedEnquiriesList.includes(enquiry.id)) {
                  let index = this.selectedEnquiriesList.indexOf(enquiry.id)
                  this.selectedEnquiriesList.splice(index, 1)
            }
      }

      // Allocate salesperson to enquiries(s).
      allocateSalesPersonToEnquiries(salesPersonID: string, enquiryID?: string): void {
            if (salesPersonID == "null") {
                  alert("Please select sales person")
                  return
            }

            let enquiryIDsToBeUpdated = []

            if (!enquiryID) {
                  for (let index = 0; index < this.selectedEnquiriesList.length; index++) {
                        enquiryIDsToBeUpdated.push({
                              "enquiryID": this.selectedEnquiriesList[index]
                        })
                  }
                  this.spinnerService.loadingMessage = "Salesperson is getting allocated to enquiries"
            }
            else {
                  enquiryIDsToBeUpdated.push({
                        "enquiryID": enquiryID
                  })
                  this.spinnerService.loadingMessage = "Salesperson is getting allocated to enquiry"
            }
            if (enquiryIDsToBeUpdated.length == 0) {
                  alert("Please select enquiries")
                  return
            }


            this.companyService.updateCompanyEnquirysSalesPerson(enquiryIDsToBeUpdated, salesPersonID).subscribe((response: any) => {
                  this.getAllEnquiries()
                  alert(response)
                  this.modalRef.close('success')
                  this.selectedEnquiriesList = []
            }, (error) => {
                  console.error(error)
                  if (typeof error.error == 'object' && error) {
                        alert(this.utilService.getErrorString(error))
                        return
                  }
                  if (error.error == undefined) {
                        alert('Sales person could not be allocated to enquiries')
                  }
                  alert(error.statusText)
            })
      }

      // =============================================================GET FUNCTIONS==========================================================================
      // Get all components.
      getAllComponents(): void {
            this.getSalesPeronList()
            this.getCountryList()
            this.getDomainList()
            this.getTechnologyList()
            this.getCompanyEnquirySourceList()
            this.getCompanyEnquiryTypeList()
            this.searchOrGetCompanyEnquiries()
            this.getPurposes()
      }

      // Get salesperson list.
      getSalesPeronList(): void {
            this.generalService.getSalesPersonList().subscribe((data: any) => {
                  this.salesPersonList = data.body
            }, (err) => {
                  console.error(err)
            })
      }

      // Get country list.
      getCountryList(): void {
            this.generalService.getCountries().subscribe((data: ICountry[]) => {
                  this.countryList = data.sort()
            }, (err) => {
                  console.error(err)
            })
      }

      // Get state list by country.
      getStateListByCountry(country: any): void {
            if (country == null) {
                  this.stateList = []
                  this.companyEnquiryForm.get('state').setValue(null)
                  this.companyEnquiryForm.get('state').disable()
                  return
            }
            if (country.name == "U.S.A." || country.name == "India") {
                  this.companyEnquiryForm.get('state').enable()
                  this.generalService.getStatesByCountryID(country.id).subscribe((respond: any[]) => {
                        this.stateList = respond
                  }, (err) => {
                        console.error(err)
                  })
            }
            else {
                  this.stateList = []
                  this.companyEnquiryForm.get('state').setValue(null)
                  this.companyEnquiryForm.get('state').disable()
            }
      }

      // Get state list by country id for enquiry search form.
      getStateListByCountryByID(countryID: any): void {
            if (countryID == null) {
                  this.stateSearchList = []
                  this.companyEnquirySearchForm.get('stateID').setValue(null)
                  this.companyEnquirySearchForm.get('stateID').disable()
                  return
            }
            let countryName: string = ""
            for (let i = 0; i < this.countryList.length; i++) {
                  if (this.countryList[i].id == countryID) {
                        countryName = this.countryList[i].name
                  }
            }
            if (countryName == "U.S.A." || countryName == "India") {
                  this.companyEnquirySearchForm.get('stateID').enable()
                  this.generalService.getStatesByCountryID(countryID).subscribe((respond: any[]) => {
                        this.stateSearchList = respond
                  }, (err) => {
                        console.error(err)
                  })
            }
            else {
                  this.stateSearchList = []
                  this.companyEnquirySearchForm.get('stateID').setValue(null)
                  this.companyEnquirySearchForm.get('stateID').disable()
            }
      }

      // Get domain list.
      getDomainList(): void {
            this.companyService.getAllDomains().subscribe((data: IDomain[]) => {
                  this.domainList = data
            }, (err) => {
                  console.error(err)
            })
      }

      // Get technology list.
      getTechnologyList(event?: any): void {
            // this.generalService.getTechnologies().subscribe((data: ITechnology[]) => {
            //       this.technologyList = data
            // }, (err) => {
            //       console.error(err)
            // })
            let queryParams: any = {}
            if (event && event?.term != "") {
                  queryParams.language = event.term
            }
            this.isTechLoading = true
            this.techService.getAllTechnologies(this.techLimit, this.techOffset, queryParams).subscribe((response) => {
                  // console.log("getTechnology -> ", response);
                  this.technologyList = []
                  this.technologyList = this.technologyList.concat(response.body)
            }, (err) => {
                  console.error(err)
            }).add(() => {
                  this.isTechLoading = false
            })
      }

      // Get enquiry sources list.
      getCompanyEnquirySourceList(): void {
            this.generalService.getGeneralTypeByType("company_enquiry_source").subscribe((respond: any[]) => {
                  this.enquirySourceList = respond
            }, (err) => {
                  console.error(err)
            })
      }

      // Get enquiry types list.
      getCompanyEnquiryTypeList(): void {
            this.generalService.getGeneralTypeByType("company_enquiry_type").subscribe((respond: any[]) => {
                  this.enquiryTypeList = respond
                  console.log(this.enquiryTypeList)
            }, (err) => {
                  console.error(err)
            })
      }

      // Get company enquiry list.
      getAllEnquiries(): void {
            this.spinnerService.loadingMessage = "Getting All Company Enquiries"



            // If not searched.
            if (this.searchFormValue == null) {
                  this.searchFormValue = {
                        roleName: this.localService.getJsonValue("roleName"),
                        loginID: this.localService.getJsonValue("loginID")
                  }
            }
            // If searched.
            else {
                  this.searchFormValue.roleName = this.localService.getJsonValue("roleName")
                  this.searchFormValue.loginID = this.localService.getJsonValue("loginID")
            }
            console.log("searchFormValue", this.searchFormValue);

            this.companyService.getAllCompanyEnquiries(this.searchFormValue).subscribe((response) => {
                  this.totalCompanyEnquiries = response.headers.get('X-Total-Count')
                  this.companyEnquiryList = response.body
            }, error => {
                  console.error(error)
                  if (error.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                  }
            }).add(() => {

                  this.setPaginationString()
            })
      }

      // Get purpose list.
      getPurposes(): void {
            this.generalService.getPurposeListByType("company_enquiry").subscribe(
                  data => {
                        this.purposeList = data
                  }
            ), err => {
                  console.error(err)
            }
      }

      // Get outcome list by purpose for call record form.
      getOutcomesByPurpose(purpose: IPurpose): void {
            this.callRecordForm.get('outcome').setValue(null)
            this.outcomeList = []
            if (!purpose) {
                  this.callRecordForm.get('outcome').setValue(null)
                  this.outcomeList = []
                  this.callRecordForm.get('outcome').disable()
                  return
            }
            this.generalService.getOutcomeListByPurpose(purpose.id).subscribe((data: any) => {
                  this.outcomeList = data
                  this.callRecordForm.get('outcome').enable()
            }, (err) => {
                  console.error(err)
            })
      }

      // Get outcome list by purpose id for enquiry search form.
      getOutcomesByPurposeID(purposeID: any): void {
            this.companyEnquirySearchForm.get('outcomeID').setValue(null)
            this.outcomeSearchList = []
            if (!purposeID) {
                  this.companyEnquirySearchForm.get('outcomeID').setValue(null)
                  this.outcomeList = []
                  this.companyEnquirySearchForm.get('outcomeID').disable()
                  return
            }
            this.generalService.getOutcomeListByPurpose(purposeID).subscribe((data: any) => {
                  this.outcomeSearchList = data
                  this.companyEnquirySearchForm.get('outcomeID').enable()
            }, (err) => {
                  console.error(err)
            })
      }
}

