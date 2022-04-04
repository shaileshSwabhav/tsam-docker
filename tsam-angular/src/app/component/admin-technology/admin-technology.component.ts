import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, Validators, FormControl } from '@angular/forms';
import { ISearchFilterField, TechnologyService, ITechnology } from 'src/app/service/technology/technology.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { Router, ActivatedRoute } from '@angular/router';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { UrlConstant } from 'src/app/service/constant';

@Component({
      selector: 'app-admin-technology',
      templateUrl: './admin-technology.component.html',
      styleUrls: ['./admin-technology.component.css']
})
export class AdminTechnologyComponent implements OnInit {

      // Flags.
      isSearched: boolean
      isOperationUpdate: boolean
      isViewMode: boolean

      // Technology.
      technologyList: ITechnology[]
      technologyForm: FormGroup

      // Pagination.
      limit: number
      currentPage: number
      totalTechnologies: number
      offset: number
      paginationStart: number
      paginationEnd: number

      // Modal.
      modalRef: any
      @ViewChild('technologyFormModal') technologyFormModal: any
      @ViewChild('deleteTechnologyModal') deleteTechnologyModal: any

      // Spinner.



      // Search.
      technologySearchForm: FormGroup
      searchFormValue: any
      searchFilterFieldList: ISearchFilterField[]

      constructor(
            private formBuilder: FormBuilder,
            private technologyService: TechnologyService,
            private spinnerService: SpinnerService,
            private modalService: NgbModal,
            private router: Router,
            private route: ActivatedRoute,
            private utilService: UtilityService,
            private urlConstant: UrlConstant,
      ) {
            this.initializeVariables()
            this.searchOrGetTechnologies()
      }


      get ongoingOperations() {
            return this.spinnerService.ongoingOperations
      }

      ngOnInit(): void { }

      // Initialize all global variables.
      initializeVariables() {
            // Components.
            this.technologyList = [] as ITechnology[]

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
            this.createTechnologyForm()
            this.createTechnologySearchForm()

            // Spinner.
            this.spinnerService.loadingMessage = "Getting Technologies"


            // Search.
            this.searchFilterFieldList = []
            this.searchFormValue = {}
      }

      // =============================================================CREATE FORMS==========================================================================
      // Create technology form.
      createTechnologyForm(): void {
            this.technologyForm = this.formBuilder.group({
                  id: new FormControl(null),
                  language: new FormControl(null, [Validators.required, Validators.maxLength(50)]),
                  rating: new FormControl(null, [Validators.required, Validators.min(1), Validators.max(5)])
            })
      }

      // Create technology search form.
      createTechnologySearchForm(): void {
            this.technologySearchForm = this.formBuilder.group({
                  language: new FormControl(null),
                  rating: new FormControl(null, [Validators.min(1), Validators.max(5)])
            })
      }
      // =============================================================TECHNOLOGY CRUD FUNCTIONS==========================================================================
      // On clicking add new technology button.
      onAddNewTechnologyClick(): void {
            this.isViewMode = false
            this.isOperationUpdate = false
            this.createTechnologyForm()
            this.openModal(this.technologyFormModal, 'md')
      }

      // Add new technology.
      addTechnology(): void {
            this.spinnerService.loadingMessage = "Adding Technology"


            this.technologyService.addTechnology(this.technologyForm.value).subscribe((response: any) => {
                  this.modalRef.close()
                  this.getAllTechnologies()
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

      // On clicking view technology button.
      onViewTechnologyClick(technology: ITechnology): void {
            this.isViewMode = true
            this.createTechnologyForm()
            this.technologyForm.patchValue(technology)
            this.technologyForm.disable()
            this.openModal(this.technologyFormModal, 'md')
      }

      // On cliking update form button in technology form.
      onUpdateTechnologyClick(): void {
            this.isViewMode = false
            this.isOperationUpdate = true
            this.technologyForm.enable()
      }

      // Update technology.
      updateTechnology(): void {
            this.spinnerService.loadingMessage = "Updating Technology"


            this.technologyService.updateTechnology(this.technologyForm.value).subscribe((response: any) => {
                  this.modalRef.close()
                  this.getAllTechnologies()
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

      // On clicking delete technology button. 
      onDeleteTechnologyClick(technologyID: string): void {
            this.openModal(this.deleteTechnologyModal, 'md').result.then(() => {
                  this.deleteTechnology(technologyID)
            }, (err) => {
                  console.error(err)
                  return
            })
      }

      // Delete technology after confirmation from user.
      deleteTechnology(technologyID: string): void {
            this.spinnerService.loadingMessage = "Deleting Technology"


            this.technologyService.deleteTechnology(technologyID).subscribe((response: any) => {
                  this.modalRef.close()
                  this.getAllTechnologies()
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

      // =============================================================TECHNOLOGY SEARCH FUNCTIONS==========================================================================
      // Reset search form and renaviagte page.
      resetSearchAndGetAll(): void {
            this.searchFilterFieldList = []
            this.technologySearchForm.reset()
            this.searchFormValue = {}
            this.changePage(1)
            this.isSearched = false
            this.router.navigate([this.urlConstant.TRAINING_TECHNOLOGY])
      }

      // Reset search form.
      resetSearchForm(): void {
            this.searchFilterFieldList = []
            this.technologySearchForm.reset()
      }

      // Search technologies.
      searchTechnologies(): void {
            this.searchFormValue = { ...this.technologySearchForm?.value }
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
            this.spinnerService.loadingMessage = "Searching Technologies"
            this.changePage(1)
      }

      // Delete search criteria from technology search form by search name.
      deleteSearchCriteria(searchName: string): void {
            this.technologySearchForm.get(searchName).setValue(null)
            this.searchTechnologies()
      }

      // ================================================OTHER FUNCTIONS FOR TECHNOLOGY===============================================
      // Page change.
      changePage(pageNumber: number): void {
            this.currentPage = pageNumber
            this.offset = this.currentPage - 1
            this.getAllTechnologies()
      }

      // Checks the url's query params and decides to call whether to call get or search.
      searchOrGetTechnologies() {
            let queryParams = this.route.snapshot.queryParams
            if (this.utilService.isObjectEmpty(queryParams)) {
                  this.getAllTechnologies()
                  return
            }
            this.technologySearchForm.patchValue(queryParams)
            this.searchTechnologies()
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

      // On clicking sumbit button in technology form.
      onFormSubmit(): void {
            if (this.technologyForm.invalid) {
                  this.technologyForm.markAllAsTouched()
                  return
            }
            if (this.isOperationUpdate) {
                  this.updateTechnology()
                  return
            }
            this.addTechnology()
      }

      // Set total talents list on current page.
      setPaginationString(): void {
            this.paginationStart = this.limit * this.offset + 1
            this.paginationEnd = +this.limit + this.limit * this.offset
            if (this.totalTechnologies < this.paginationEnd) {
                  this.paginationEnd = this.totalTechnologies
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
      // Get all technologies.
      getAllTechnologies() {
            this.spinnerService.loadingMessage = "Getting All Technologies"


            this.searchFormValue.limit = this.limit
            this.searchFormValue.offset = this.offset
            this.technologyService.getAllTechnologies(this.limit, this.offset, this.searchFormValue).subscribe((response) => {
                  this.totalTechnologies = response.headers.get('X-Total-Count')
                  this.technologyList = response.body
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