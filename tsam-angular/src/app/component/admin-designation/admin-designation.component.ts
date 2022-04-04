import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, Validators, FormControl } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { DesignationService, IDesignation, ISearchFilterField } from 'src/app/service/designation/designation.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
      selector: 'app-admin-designation',
      templateUrl: './admin-designation.component.html',
      styleUrls: ['./admin-designation.component.css']
})
export class AdminDesignationComponent implements OnInit {

      // Flags.
      isSearched: boolean
      isOperationUpdate: boolean
      isViewMode: boolean

      // Designation.
      designationList: IDesignation[]
      designationForm: FormGroup

      // Pagination.
      limit: number
      currentPage: number
      totalDesignations: number
      offset: number
      paginationStart: number
      paginationEnd: number

      // Modal.
      modalRef: any
      @ViewChild('designationFormModal') designationFormModal: any
      @ViewChild('deleteDesignationModal') deleteDesignationModal: any

      // Spinner.



      // Search.
      designationSearchForm: FormGroup
      searchFormValue: any
      searchFilterFieldList: ISearchFilterField[]

      constructor(
            private formBuilder: FormBuilder,
            private designationService: DesignationService,
            private spinnerService: SpinnerService,
            private modalService: NgbModal,
            private router: Router,
            private route: ActivatedRoute,
            private utilService: UtilityService,
      ) {
            this.initializeVariables()
            this.searchOrGetDesignations()
      }


      get ongoingOperations() {
            return this.spinnerService.ongoingOperations
      }

      ngOnInit(): void { }

      // Initialize all global variables.
      initializeVariables() {
            // Components.
            this.designationList = [] as IDesignation[]

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
            this.createDesignationForm()
            this.createDesignationSearchForm()

            // Spinner.
            this.spinnerService.loadingMessage = "Getting Designations"


            // Search.
            this.searchFilterFieldList = []
            this.searchFormValue = {}
      }

      // =============================================================CREATE FORMS==========================================================================
      // Create designation form.
      createDesignationForm(): void {
            this.designationForm = this.formBuilder.group({
                  id: new FormControl(null),
                  position: new FormControl(null, [Validators.required, Validators.maxLength(200)]),
            })
      }

      // Create designation search form.
      createDesignationSearchForm(): void {
            this.designationSearchForm = this.formBuilder.group({
                  position: new FormControl(null)
            })
      }
      // =============================================================DESIGNATION CRUD FUNCTIONS==========================================================================
      // On clicking add new designation button.
      onAddNewDesignationClick(): void {
            this.isViewMode = false
            this.isOperationUpdate = false
            this.createDesignationForm()
            this.openModal(this.designationFormModal, 'sm')
      }

      // Add new designation.
      addDesignation(): void {
            this.spinnerService.loadingMessage = "Adding Designation"


            this.designationService.addDesignation(this.designationForm.value).subscribe((response: any) => {
                  this.modalRef.close()
                  this.getAllDesignations()
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

      // On clicking view designation button.
      onViewDesignationClick(designation: IDesignation): void {
            this.isViewMode = true
            this.createDesignationForm()
            this.designationForm.patchValue(designation)
            this.designationForm.disable()
            this.openModal(this.designationFormModal, 'sm')
      }

      // On cliking update form button in designation form.
      onUpdateDesignationClick(): void {
            this.isViewMode = false
            this.isOperationUpdate = true
            this.designationForm.enable()
      }

      // Update designation.
      updateDesignation(): void {
            this.spinnerService.loadingMessage = "Updating Designation"


            this.designationService.updateDesignation(this.designationForm.value).subscribe((response: any) => {
                  this.modalRef.close()
                  this.getAllDesignations()
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

      // On clicking delete designation button. 
      onDeleteDesignationClick(designationID: string): void {
            this.openModal(this.deleteDesignationModal, 'md').result.then(() => {
                  this.deleteDesignation(designationID)
            }, (err) => {
                  console.error(err)
                  return
            })
      }

      // Delete designation after confirmation from user.
      deleteDesignation(designationID: string): void {
            this.spinnerService.loadingMessage = "Deleting Designation"


            this.designationService.deleteDesignation(designationID).subscribe((response: any) => {
                  this.modalRef.close()
                  this.getAllDesignations()
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

      // =============================================================DESIGNATION SEARCH FUNCTIONS==========================================================================
      // Reset search form and renaviagte page.
      resetSearchAndGetAll(): void {
            this.searchFilterFieldList = []
            this.designationSearchForm.reset()
            this.searchFormValue = {}
            this.changePage(1)
            this.isSearched = false
      }

      // Reset search form.
      resetSearchForm(): void {
            this.searchFilterFieldList = []
            this.designationSearchForm.reset()
      }

      // Redirect to Same page.
      redirectToSamePage(): void {
            this.router.navigate([], {
                  relativeTo: this.route,
                  queryParams: this.searchFormValue,
            })
      }

      // Search designations.
      searchDesignations(): void {
            this.searchFormValue = { ...this.designationSearchForm?.value }
            this.searchFormValue.limit = this.limit
            this.searchFormValue.offset = this.offset
            for (let field in this.searchFormValue) {
                  if (field == "limit" || field == "offset") {
                        continue
                  }
                  if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
                        delete this.searchFormValue[field]
                  } else {
                        if (field)
                              this.isSearched = true
                  }
            }
            this.searchFilterFieldList = []
            for (var property in this.searchFormValue) {
                  if (property == "limit" || property == "offset") {
                        continue
                  }
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
            this.spinnerService.loadingMessage = "Searching Designations"
            this.changePage(1)
      }

      // Delete search criteria from designation search form by search name.
      deleteSearchCriteria(searchName: string): void {
            this.designationSearchForm.get(searchName).setValue(null)
            this.searchDesignations()
      }

      // ================================================OTHER FUNCTIONS FOR DESIGNATION===============================================
      // Page change.
      changePage(pageNumber: number): void {
            this.currentPage = pageNumber
            this.offset = this.currentPage - 1
            this.getAllDesignations()
      }

      // Checks the url's query params and decides to call whether to call get or search.
      searchOrGetDesignations() {
            let queryParams = this.route.snapshot.queryParams
            if (queryParams?.limit != null) {
                  this.limit = queryParams.limit
            } if (queryParams?.offset != null) {
                  this.offset = queryParams.offset
            }
            if (this.utilService.isObjectEmpty(queryParams)) {
                  this.getAllDesignations()
                  return
            }
            this.designationSearchForm.patchValue(queryParams)
            this.searchDesignations()
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

      // On clicking sumbit button in designation form.
      onFormSubmit(): void {
            if (this.designationForm.invalid) {
                  this.designationForm.markAllAsTouched()
                  return
            }
            if (this.isOperationUpdate) {
                  this.updateDesignation()
                  return
            }
            this.addDesignation()
      }

      // Set total talents list on current page.
      setPaginationString(): void {
            this.paginationStart = this.limit * this.offset + 1
            this.paginationEnd = +this.limit + this.limit * this.offset
            if (this.totalDesignations < this.paginationEnd) {
                  this.paginationEnd = this.totalDesignations
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
      // Get all designations.
      getAllDesignations() {
            this.spinnerService.loadingMessage = "Getting All Designations"


            this.searchFormValue.limit = this.limit
            this.searchFormValue.offset = this.offset
            this.redirectToSamePage()
            this.designationService.getAllDesignations(this.searchFormValue).subscribe((response) => {
                  this.totalDesignations = response.headers.get('X-Total-Count')
                  this.designationList = response.body
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