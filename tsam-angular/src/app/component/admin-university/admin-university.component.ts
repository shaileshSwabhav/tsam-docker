import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { UrlConstant } from 'src/app/service/constant';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { ICountry, ISearchFilterField, IUniversity, UniveristyService } from 'src/app/service/university/univeristy.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-admin-university',
  templateUrl: './admin-university.component.html',
  styleUrls: ['./admin-university.component.css']
})
export class AdminUniversityComponent implements OnInit {

  // Components.
  countryList: ICountry[]
  uploadedUniversities: IUniversity[]

  // Flags.
  isOperationUpdate: boolean
  isUploadedSuccessfully: boolean
  isViewMode: boolean
  isSearched: boolean

  // University.
  universityList: IUniversity[]
  universityForm: FormGroup

  // Pagination.
  limit: number
  currentPage: number
  totalUniversities: number
  offset: number
  paginationStart: number
  paginationEnd: number

  // Modal.
  modalRef: any
  @ViewChild('universityFormModal') universityFormModal: any
  @ViewChild('deleteUniversityModal') deleteUniversityModal: any

  // Spinner.



  // Search.
  universitySearchForm: FormGroup
  searchFormValue: any
  searchFilterFieldList: ISearchFilterField[]

  // For modal nav tab.
  active: number

  // Constants.
  readonly UNIVERSITY_EXCEL_DEMO_LINK = this.urlConstant.UNIVERSITY_DEMO

  // beta
  readonly OPERATIONS = { 'add': "Add University", 'update': "Update University" }

  constructor(
    private formBuilder: FormBuilder,
    private universityService: UniveristyService,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private generalService: GeneralService,
    private fileOps: FileOperationService,
    private urlConstant: UrlConstant,
    private router: Router,
    private route: ActivatedRoute,
    private utilService: UtilityService,
  ) {
    this.initializeVariables()
    this.getAllComponents()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit() { }

  // Initialize all global variables.
  initializeVariables() {
    this.universityList = [] as IUniversity[]
    this.uploadedUniversities = [] as IUniversity[]
    this.countryList = [] as ICountry[]

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
    this.createUniversityForm()
    this.createUniversitySearchForm()

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Universities"


    // For nav tab.
    this.active = 1

    // Search.
    this.searchFilterFieldList = []
    this.searchFormValue = {}
  }

  // =============================================================CREATE FORMS==========================================================================
  // Create university form.
  createUniversityForm(): void {
    this.universityForm = this.formBuilder.group({
      id: new FormControl(null),
      universityName: new FormControl(null, [Validators.required, Validators.max(200)]),
      country: new FormControl(null, [Validators.required])
    })
  }

  // Create university search form.
  createUniversitySearchForm(): void {
    this.universitySearchForm = this.formBuilder.group({
      universityName: new FormControl(null, [Validators.max(200)]),
      countryID: new FormControl(null)
    })
  }

  // =============================================================UNIVERSITY CRUD FUNCTIONS==========================================================================
  // On clicking add new university button.
  onAddNewUniversityClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.isUploadedSuccessfully = false
    this.createUniversityForm()
    this.openModal(this.universityFormModal, 'md')
  }

  // Add new university.
  addUniversity(): void {
    this.spinnerService.loadingMessage = "Adding University"


    this.universityService.addUniversity(this.universityForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllUniversities()
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

  // On clicking view university button.
  onViewUniversityClick(university: IUniversity): void {
    this.active = 1
    this.isViewMode = true
    this.createUniversityForm()
    this.universityForm.patchValue(university)
    this.universityForm.disable()
    this.openModal(this.universityFormModal, 'md')
  }

  // On cliking update form button in university form.
  onUpdateUniversityClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = true
    this.isUploadedSuccessfully = false
    this.active = 1
    this.universityForm.enable()
  }

  // Update university.
  updateUniversity(): void {
    this.spinnerService.loadingMessage = "Updating University"


    this.universityService.updateUniversity(this.universityForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllUniversities()
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

  // On clicking delete university button. 
  onDeleteUniversityClick(universityID: string): void {
    this.openModal(this.deleteUniversityModal, 'md').result.then(() => {
      this.deleteUniversity(universityID)
    }, (err) => {
      console.error(err)
      return
    })
  }

  // Delete university after confirmation from user.
  deleteUniversity(universityID: string): void {
    this.spinnerService.loadingMessage = "Deleting University"

    this.modalRef.close()

    this.universityService.deleteUniversity(universityID).subscribe((response: any) => {
      this.getAllUniversities()
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

  // =============================================================UNIVERSITY SEARCH FUNCTIONS==========================================================================
  // Reset search form and renaviagte page.
  resetSearchAndGetAll(): void {
    this.searchFilterFieldList = []
    this.universitySearchForm.reset()
    this.searchFormValue = {}
    this.changePage(1)
    this.isSearched = false
    this.router.navigate([this.urlConstant.TALENT_UNIVERSITY])
  }

  // Reset search form.
  resetSearchForm(): void {
    this.searchFilterFieldList = []
    this.universitySearchForm.reset()
  }

  // Search universities.
  searchUniversities(): void {
    this.searchFormValue = { ...this.universitySearchForm?.value }
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
    this.spinnerService.loadingMessage = "Searching Universities"
    this.changePage(1)
  }

  // Delete search criteria from university search form by search name.
  deleteSearchCriteria(searchName: string): void {
    this.universitySearchForm.get(searchName).setValue(null)
    this.searchUniversities()
  }

  // ================================================OTHER FUNCTIONS FOR UNIVERSITY===============================================
  // Page change.
  changePage(pageNumber: number): void {

    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
    this.getAllUniversities()
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetUniversities() {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getAllUniversities()
      return
    }
    this.universitySearchForm.patchValue(queryParams)
    this.searchUniversities()
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

  // On clicking sumbit button in university form.
  onFormSubmit(): void {
    if (this.universityForm.invalid) {
      this.universityForm.markAllAsTouched()
      return
    }
    if (this.isOperationUpdate) {
      this.updateUniversity()
      return
    }
    this.addUniversity()
  }

  // Set total universities list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalUniversities < this.paginationEnd) {
      this.paginationEnd = this.totalUniversities
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

  // =============================================================GET FUNCTIONS==========================================================================
  // Get all components.
  getAllComponents(): void {
    this.searchOrGetUniversities()
    this.getCountryList()
  }

  // Get country list.
  getCountryList(): void {
    this.generalService.getCountries().subscribe((data: any) => {
      this.countryList = data
    }, (err) => {
      console.error(err)
    })
  }

  // Get all universities.
  getAllUniversities(): void {
    this.spinnerService.loadingMessage = "Getting Universities"


    this.searchFormValue.limit = this.limit
    this.searchFormValue.offset = this.offset
    this.universityService.getAllUniversities(this.searchFormValue).
      subscribe((data) => {
        this.totalUniversities = data.headers.get('X-Total-Count')
        this.universityList = data.body

        // // temporary as not sure about showing country.
        // this.universityList.forEach(university => {
        //   this.allCountries.forEach(country => {
        //     if (country.id == university.countryID) {
        //       university.country.name = country.name
        //     }
        //   })
        // })

      },
        (error) => {
          console.error(error)
          if (error.statusText.includes('Unknown')) {
            alert("No connection to server. Check internet.")
          }
        }).add(() => {
          this.setPaginationString()
        })
  }

  // =============================================================EXCEL UPLOAD FUNCTIONS==========================================================================
  // Use service.
  onFileChange(event: any, excelFile: any): void {
    const files = event.target.files
    if (files.length === 0) {
      return
    }
    if (files.length !== 1) {
      alert('only 1 file should be uploaded')
      return
    }
    const file = files[0]
    this.spinnerService.loadingMessage = "Uploading excel"

    this.fileOps.uploadExcel(file).subscribe((uploadedUniversities: IUniversity[]) => {
      if (this.validateUniversities(uploadedUniversities)) {
        this.uploadedUniversities = uploadedUniversities
        this.isUploadedSuccessfully = true
      } else {
        excelFile.value = ""
      }

    }, (err) => {

      alert(err)
      excelFile.value = ""
    })
  }

  // Validate universities in file.
  validateUniversities(universities: IUniversity[]): boolean {
    if (!universities || universities.length == 0) {
      alert("No universities")
      return false
    }
    for (let i = 0; i < universities.length; i++) {
      if (!universities[i].universityName || universities[i].universityName === "") {
        alert(`university name on row ${i + 2} not specified`)
        return false
      } else if (!universities[i].countryName || universities[i].countryName === "") {
        alert(`country name of university ${universities[i].universityName} is not specified`)
        return false
      }
      let country = { name: universities[i].countryName }
      universities[i].country = country
    }
    return true
  }

  onTabChange(event: any) {
    console.log(event)
  }

  // Add multiple universities.
  addMultipleUniversities() {
    this.spinnerService.loadingMessage = "Adding all universities"

    this.universityService.addMultipleUniversities(this.uploadedUniversities).subscribe((response: any) => {

      this.modalRef.close()
      this.getAllUniversities()
      alert(response)
    }, (error) => {
      console.error(error)

      alert(error.error?.error)
    })
  }
}
/* Notes

Excel:-
     Output format
            header?: "A"|number|string[];

            // Override worksheet range
            range?: any;

            // Include or omit blank lines in the output
            blankrows?: boolean;

            // Default value for null/undefined values
            defval?: any;

            // if true, return raw data; if false, return formatted text
            raw?: boolean;

            // if true, return raw numbers; if false, return formatted numbers
            rawNumbers?: boolean;




Tab change
pages: string[] = ["tab-1", "tab-2"];
  disableEnd: boolean = true;

  actualPage = 1

  goToEnd(tabSet): void {
    this.disableEnd = false;
    var that = this
    setTimeout(function () {
      tabSet.select(that.pages[1]);
    }, 1);
  }

  tabChange(evt) {
    console.log(evt)
  }
  */