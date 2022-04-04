import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, FormControl, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalRef, NgbModalOptions } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { ISearchFilterField } from 'src/app/service/general/general.service';
import { IProgrammingLanguage, ProgrammingLanguageService } from 'src/app/service/programming-language/programming-language.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { UrlConstant } from 'src/app/service/constant';

@Component({
  selector: 'app-programming-language',
  templateUrl: './programming-language.component.html',
  styleUrls: ['./programming-language.component.css']
})
export class ProgrammingLanguageComponent implements OnInit {

  // Flags.
  isSearched: boolean
  isOperationUpdate: boolean
  isViewMode: boolean

  // Programming language.
  languageList: IProgrammingLanguage[]
  languageForm: FormGroup

  // Pagination.
  limit: number
  currentPage: number
  totalLanguages: number
  offset: number
  paginationStart: number
  paginationEnd: number

  // Modal.
  modalRef: any
  @ViewChild('languageFormModal') languageFormModal: any
  @ViewChild('deleteLanguageModal') deleteLanguageModal: any

  // Spinner.



  // Search.
  languageSearchForm: FormGroup
  searchFormValue: any
  searchFilterFieldList: ISearchFilterField[]

  constructor(
    private formBuilder: FormBuilder,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private router: Router,
    private route: ActivatedRoute,
    private utilService: UtilityService,
    private languageService: ProgrammingLanguageService,
    private urlConstant: UrlConstant,
  ) {
    this.initializeVariables()
    this.searchOrGetLanguages()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize all global variables.
  initializeVariables() {

    // Components.
    this.languageList = [] as IProgrammingLanguage[]

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
    this.createLanguageForm()
    this.createLanguageSearchForm()

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Programming Languages"


    // Search.
    this.searchFilterFieldList = []
  }

  // =============================================================CREATE FORMS==========================================================================
  // Create programming language form.
  createLanguageForm(): void {
    this.languageForm = this.formBuilder.group({
      id: new FormControl(null),
      name: new FormControl(null, [Validators.required, Validators.maxLength(100)]),
      rating: new FormControl(null, [Validators.required, Validators.min(1), Validators.max(5)]),
      fileExtension: new FormControl(null, [Validators.maxLength(20)]),
    })
  }

  // Create programming language search form.
  createLanguageSearchForm(): void {
    this.languageSearchForm = this.formBuilder.group({
      name: new FormControl(null),
      rating: new FormControl(null, [Validators.min(1), Validators.max(5)])
    })
  }
  // =============================================================PROGRAMMING LANGUAGE CRUD FUNCTIONS==========================================================================
  // On clicking add new programming language button.
  onAddNewLanguageClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.createLanguageForm()
    this.openModal(this.languageFormModal, 'lg')
  }

  // Add new programming language.
  addLanguage(): void {
    this.spinnerService.loadingMessage = "Adding Programming Language"


    this.languageService.addProgrammingLanguage(this.languageForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllLanguages()
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

  // On clicking view programming language button.
  onViewLanguageClick(language: IProgrammingLanguage): void {
    this.isViewMode = true
    this.createLanguageForm()
    this.languageForm.patchValue(language)
    this.languageForm.disable()
    this.openModal(this.languageFormModal, 'lg')
  }

  // On cliking update form button in programming language form.
  onUpdateLanguageClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = true
    this.languageForm.enable()
  }

  // Update programming language.
  updateLanguage(): void {
    this.spinnerService.loadingMessage = "Updating Programming Language"


    this.languageService.updateProgrammingLanguage(this.languageForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllLanguages()
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

  // On clicking delete programming language button. 
  onDeleteLanguageClick(languageID: string): void {
    this.openModal(this.deleteLanguageModal, 'md').result.then(() => {
      this.deleteLanguage(languageID)
    }, (err) => {
      console.error(err)
      return
    })
  }

  // Delete programming language after confirmation from user.
  deleteLanguage(languageID: string): void {
    this.spinnerService.loadingMessage = "Deleting Programming Language"


    this.languageService.deleteProgrammingLanguage(languageID).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllLanguages()
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

  // =============================================================PROGRAMMING LANGUAGE SEARCH FUNCTIONS==========================================================================
  // Reset search form and renaviagte page.
  resetSearchAndGetAll(): void {
    this.searchFilterFieldList = []
    this.languageSearchForm.reset()
    this.searchFormValue = null
    this.changePage(1)
    this.isSearched = false
    this.router.navigate([this.urlConstant.TRAINING_INPUT_LANGUAGE])
  }

  // Reset search form.
  resetSearchForm(): void {
    this.searchFilterFieldList = []
    this.languageSearchForm.reset()
  }

  // Search programming languages.
  searchLanguages(): void {
    this.searchFormValue = { ...this.languageSearchForm?.value }
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
    this.spinnerService.loadingMessage = "Searching Programming Languages"
    this.changePage(1)
  }

  // Delete search criteria from programming language search form by search name.
  deleteSearchCriteria(searchName: string): void {
    this.languageSearchForm.get(searchName).setValue(null)
    this.searchLanguages()
  }

  // ================================================OTHER FUNCTIONS FOR PROGRAMMING LANGUAGE===============================================
  // Page change.
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
    this.getAllLanguages()
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetLanguages() {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getAllLanguages()
      return
    }
    this.languageSearchForm.patchValue(queryParams)
    this.searchLanguages()
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

  // On clicking sumbit button in programming language form.
  onFormSubmit(): void {
    if (this.languageForm.invalid) {
      this.languageForm.markAllAsTouched()
      return
    }
    if (this.isOperationUpdate) {
      this.updateLanguage()
      return
    }
    this.addLanguage()
  }

  // Set total talents list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalLanguages < this.paginationEnd) {
      this.paginationEnd = this.totalLanguages
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

  // Get all programming languages.
  getAllLanguages() {
    this.spinnerService.loadingMessage = "Getting Programming Languages"


    this.languageService.getAllProgrammingLanguages(this.limit, this.offset, this.searchFormValue).subscribe((response) => {
      this.totalLanguages = response.headers.get('X-Total-Count')
      this.languageList = response.body
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
