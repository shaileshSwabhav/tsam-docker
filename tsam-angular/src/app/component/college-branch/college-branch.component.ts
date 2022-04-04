import { Component, OnInit, ViewChild } from '@angular/core'
import { FormGroup, FormBuilder, Validators, FormArray, FormControl } from '@angular/forms'
import { ICollegeBranch, IState, ICountry, CollegeService, ICollege, ISalesperson } from 'src/app/service/college/college.service'
import { GeneralService } from 'src/app/service/general/general.service'
import { Constant, Role, UrlConstant } from 'src/app/service/constant'
import { SpinnerService } from 'src/app/service/spinner/spinner.service';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap'
import { ActivatedRoute, Router } from '@angular/router'
import { IPermission } from 'src/app/service/menu/menu.service'
import { UtilityService } from 'src/app/service/utility/utility.service'
import { LocalService } from 'src/app/service/storage/local.service'
import { Location } from '@angular/common'
import { CollegeMasterComponent } from '../college-master/college-master.component'
import { UrlService } from 'src/app/service/url.service'

@Component({
      selector: 'app-college-branch',
      templateUrl: './college-branch.component.html',
      styleUrls: ['./college-branch.component.css']
})
export class CollegeBranchComponent implements OnInit {

      // College.
      collegeBranch: ICollegeBranch
      collegeBranchList: any[]
      collegeBranchForm: FormGroup

      // Additional Components.
      salesPeople: any[]
      states: IState[]
      countries: ICountry[]
      ratingList: number[]
      allUniversities: any[]

      // Search.
      searchBranchForm: FormGroup
      searchFormValue: any
      isSearched: boolean
      showSearch: boolean

      // Pagination.
      limit: number
      offset: number
      currentPage: number
      totalCollegeBranches: number
      paginationString: string

      // Modal.
      selectedCollegeBranch: ICollegeBranch
      modalHeader: string
      modalButton: string
      modalAction: () => void
      formHandler: () => void
      modalRef: any

      // access
      permission: IPermission
      isViewClicked: boolean
      roleName: string
      isAdmin: boolean
      loginID: string

      //spinner


      // navigation
      isNavigatedFromColleges: boolean
      navigatedCollegeID: string

      // constants
      nilUUID: string

      option: string

      @ViewChild('branchDetailModal') branchDetailModal: any
      @ViewChild('deleteConfirmationModal') deleteConfirmationModal: any
      @ViewChild('drawer') drawer: any

      constructor(
            private formBuilder: FormBuilder,
            private collegeService: CollegeService,
            private generalService: GeneralService,
            private activatedRoute: ActivatedRoute,
            private constant: Constant,
            private roleConstant: Role,
            private urlConstant: UrlConstant,
            private spinnerService: SpinnerService,
            private modalService: NgbModal,
            private utilService: UtilityService,
            private localService: LocalService,
            private router: Router,
            private urlService: UrlService,
      ) {
            this.initializeVariables()
            this.createForms()
            this.getAllComponents()
      }

      initializeVariables(): void {

            this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'),
                  this.urlConstant.COLLEGE_BRANCH)

            this.roleName = this.localService.getJsonValue("roleName");
            this.loginID = this.localService.getJsonValue("loginID");

            this.option = ""
            this.salesPeople = []
            this.ratingList = []

            this.limit = 5
            this.offset = 0

            this.isSearched = false
            this.showSearch = false
            this.isNavigatedFromColleges = false
            this.isAdmin = false
            this.nilUUID = this.constant.NIL_UUID
            this.searchFormValue = {}

            this.spinnerService.loadingMessage = "Getting college branches"

            // Check whether the user is admin.
            if (this.roleName === "Admin") {
                  this.isAdmin = true;
            }
      }


      get ongoingOperations() {
            return this.spinnerService.ongoingOperations
      }

      ngOnInit() {
      }

      backToPreviousPage(): void {
            this.urlService.goBack(CollegeMasterComponent.name)
      }

      // gets all required components like all states,countries etc
      getAllComponents(): void {
            this.getAllCountries()
            this.getAllRatings()
            this.getAllSalesPeople()
            this.getAllUniversities()
            this.searchOrGetCollegeBranches()
      }

      createForms(): void {
            this.createCollegeBranchForm()
            this.createSearchBranchForm()
      }

      searchOrGetCollegeBranches(): void {
            let queryParams = this.activatedRoute.snapshot.queryParams
            if (this.utilService.isObjectEmpty(queryParams)) {
                  this.getAllBranches()
                  return
            }
            if (queryParams.collegeID) {
                  this.navigatedCollegeID = queryParams.collegeID
                  this.isNavigatedFromColleges = true
                  this.searchBranchForm.patchValue(queryParams)
                  this.searchBranches()
            }

      }


      // Handles pagination
      changePage($event: any): void {
            // $event will be the page number & offset will be 1 less than it.
            this.offset = $event - 1
            this.currentPage = $event
            this.getAllBranches()
      }

      getAllUniversities(): void {
            this.generalService.getAllUniversities().subscribe((data: any) => {
                  this.allUniversities = data
                  // console.log(this.allUniversities)
            }, (err) => {
                  console.error(err)
            })
      }

      getAllSalesPeople(): void {
            this.generalService.getSalesPersonList().subscribe((data: any) => {
                  this.salesPeople = data.body
            }, (err) => {
                  console.error(err)
            })
      }

      // return state list.
      getStateList(countryID: any) {
            this.generalService.getStatesByCountryID(countryID).subscribe((data: IState[]) => {
                  this.states = data
            }, (err) => {
                  console.error(err)
            })
      }

      getAllCountries(): void {
            this.generalService.getCountries().subscribe((data: ICountry[]) => {
                  this.countries = data
            }, (err) => {
                  console.error(err)
            })
      }

      getAllRatings(): void {
            this.ratingList = []
            this.generalService.getGeneralTypeByType("college_rating").subscribe((data: any[]) => {
                  for (let rating of data) {
                        this.ratingList.push(+rating.value)
                  }
                  this.ratingList = this.ratingList.sort()
            }, (err) => {
                  console.error(err)
            })
      }



      compareFn(c1: any, c2: any): boolean {
            return c1 && c2 ? c1.id === c2.id : c1 === c2
      }

      // Create search form.
      createSearchBranchForm(): void {

            this.searchBranchForm = this.formBuilder.group({
                  city: new FormControl(null, [Validators.maxLength(50), Validators.pattern("^[a-zA-Z]+([a-zA-Z ]?)+")]),
                  collegeRating: new FormControl(null),
                  allIndiaRanking: new FormControl(null, [Validators.max(100000), Validators.pattern("^[0-9]*$")]),
                  salesPersonID: new FormControl(null),
                  universityID: new FormControl(null),
                  collegeID: new FormControl(null),
                  stateID: new FormControl(null),
                  tpoName: new FormControl(null, [Validators.maxLength(70), Validators.pattern("^[a-zA-Z]+([a-zA-Z. ]?)+")]),
                  branchName: new FormControl(null, [Validators.maxLength(150)]),
            })
      }

      // Create college form.
      createCollegeBranchForm(): void {
            this.collegeBranchForm = this.formBuilder.group({
                  id: new FormControl(null),
                  collegeID: new FormControl(null, [Validators.required]),
                  branchName: new FormControl(null, [Validators.required, Validators.maxLength(150)]),
                  university: new FormControl(null, [Validators.required]),
                  code: new FormControl(null),
                  salesPerson: new FormControl(null),
                  tpoName: new FormControl(null, [Validators.maxLength(70), Validators.pattern("^[a-zA-Z ]*$")]),
                  tpoContact: new FormControl(null, [Validators.pattern(/^(?:(?:\+|0{0,2})91(\s*[\-]\s*)?|[0]?)?[6789]\d{9}$/)]),
                  tpoAlternateContact: new FormControl(null, [Validators.pattern(/^(?:(?:\+|0{0,2})91(\s*[\-]\s*)?|[0]?)?[6789]\d{9}$/)]),
                  tpoEmail: new FormControl(null, [Validators.pattern(/^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$/)]),
                  collegeRating: new FormControl(null),
                  allIndiaRanking: new FormControl(null, [Validators.max(100000), Validators.pattern("^[1-9][0-9]*$")]),
                  email: new FormControl(null, [Validators.pattern(/^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$/)]),
                  address: this.formBuilder.group({
                        state: new FormControl(null, [Validators.required]),
                        country: new FormControl(null, [Validators.required]),
                        address: new FormControl(null, [Validators.required, Validators.maxLength(100)]),
                        city: new FormControl(null, [Validators.required, Validators.maxLength(50), Validators.pattern(/^[a-zA-Z]+([a-zA-Z ]?)+$/)]),
                        pinCode: new FormControl(null, [Validators.required, Validators.pattern(/^[1-9][0-9]{5}$/)])
                  })
            })
      }

      deleteCollegeBranch(): void {
            this.modalRef.close()
            this.spinnerService.loadingMessage = "Deleting college branch"

            this.collegeService.deleteCollegeBranch(this.selectedCollegeBranch.collegeID,
                  this.selectedCollegeBranch.id).subscribe((data: any) => {
                        alert(data)
                        if (this.totalCollegeBranches == 1) {
                              this.collegeBranchList = []
                              this.totalCollegeBranches = 0

                              return
                        }
                        this.getAllBranches()
                  }, (error) => {
                        console.error(error)
                        if (error.statusText.includes('Unknown')) {
                              alert("No connection to server. Check internet.")
                              return
                        }
                        alert(error.statusText)
                  })
      }

      setPaginationString(): void {
            this.paginationString = ''
            let start: number = this.limit * this.offset + 1
            let end: number = +this.limit + this.limit * this.offset
            if (this.totalCollegeBranches < end) {
                  end = this.totalCollegeBranches
            }
            if (this.totalCollegeBranches == 0) {
                  this.paginationString = ''
                  return
            }
            this.paginationString = `${start} - ${end}`
      }

      closeSearchDrawer(): void {
            this.drawer.toggle()
            this.searchBranches()
      }

      searchBranches(): void {
            if (this.roleName == this.roleConstant.SALES_PERSON) {
                  this.searchBranchForm.get("salesPersonID").setValue(this.loginID)
            }

            this.searchFormValue = { ...this.searchBranchForm.value }
            this.router.navigate([], {
                  relativeTo: this.activatedRoute,
                  queryParams: this.searchFormValue,
            })

            let flag: boolean = true
            for (let field in this.searchFormValue) {
                  if (!this.searchFormValue[field]) {
                        delete this.searchFormValue[field]
                  } else {
                        if (field != "collegeID") {
                              this.isSearched = true
                        }
                        flag = false
                  }
            }

            // No API call on empty search.
            if (flag) {
                  return
            }
            this.changePage(1)
      }

      getAllBranches(): void {
            this.spinnerService.loadingMessage = "Getting college branches"


            // if (this.roleName === this.roleConstant.SALES_PERSON) {
            //       let loginID: string = this.localService.getJsonValue("loginID");
            //       this.searchFormValue["salesPersonID"] = loginID
            // }
            // if (this.isNavigatedFromColleges && !this.isSearched) {
            //       this.searchFormValue["collegeID"] = this.navigatedCollegeID
            // }
            this.collegeService.getAllCollegeBranches(this.limit, this.offset, this.searchFormValue).subscribe(data => {
                  this.collegeBranchList = data.body
                  this.totalCollegeBranches = parseInt(data.headers.get('X-Total-Count'))
            }, (error) => {
                  this.totalCollegeBranches = 0
                  console.error(error)
                  if (error.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                        return
                  }
                  alert(error.statusText)
            }).add(() => {
                  this.setPaginationString()

            })
      }

      // getAllBranchesByCollege(collegeID: string): void {
      //       this.searched = true
      //       this.spinnerService.loadingMessage = "Getting college branches"
      //       this.collegeService.getAllBranchesOfCollege(collegeID, this.searchFormValue).subscribe(data => {
      //             this.collegeBranchList = data.body
      //             console.log("Branches are :", data.body, collegeID);
      //             
      //             this.totalCollegeBranches = this.collegeBranchList.length
      //       }, (error) => {
      //             this.totalCollegeBranches = 0
      //             
      //             console.log(error)
      //             if (error.statusText.includes('Unknown')) {
      //                   alert("No connection to server. Check internet.")
      //             }
      //       })
      // }

      resetSearchAndGetAll(): void {
            this.isSearched = false
            this.showSearch = false
            this.searchBranchForm.reset()
            this.searchFormValue = {}

            if (this.roleName == this.roleConstant.SALES_PERSON) {
                  this.searchBranchForm.get("salesPersonID").setValue(this.loginID)
            }
            if (this.navigatedCollegeID) {
                  this.searchBranchForm.get("collegeID").setValue(this.navigatedCollegeID)

                  this.searchFormValue = { ...this.searchBranchForm?.value }
                  this.router.navigate([this.urlConstant.COLLEGE_BRANCH], {
                        relativeTo: this.activatedRoute,
                        queryParams: this.searchFormValue,
                  })
                  this.searchBranches()
                  return
            }

            this.router.navigate([this.urlConstant.COLLEGE_BRANCH])
            this.changePage(1)
      }

      resetSearchForm(): void {
            this.searchBranchForm.reset()
      }

      // Update College.
      updateCollegeBranch(): void {
            this.spinnerService.loadingMessage = "Updating college branch"

            this.collegeService.updateCollegeBranch(this.collegeBranchForm.value).subscribe((data: string) => {
                  this.modalRef.close()
                  this.getAllBranches()
                  this.collegeBranchForm.reset()
                  alert(data)
            }, (error) => {
                  console.error(error)

                  if (error.error) {
                        alert(error.error)
                        return
                  }
                  alert(error.statusText)
            })
      }

      onUpdateBranchClick() {
            this.collegeBranchForm.enable()
            this.setVariableForUpdateCollegeBranch()
            this.isViewClicked = false
      }

      onViewCollegeClick(branch: ICollegeBranch): void {
            this.getStateList(branch.address?.country?.id)
            // console.log(branch);
            this.modalHeader = "View College Branch"
            this.isViewClicked = true
            this.collegeBranch = branch
            this.updateForm()
            this.collegeBranchForm.disable()
            this.openModal(this.branchDetailModal)
      }

      onDeleteBranchClick(branch: ICollegeBranch): void {
            this.selectedCollegeBranch = branch
            this.openModal(this.deleteConfirmationModal, "md")
      }

      updateForm(): void {
            this.createCollegeBranchForm()
            // console.log(this.collegeBranch);
            this.collegeBranchForm.patchValue(this.collegeBranch)
            // console.log("form", this.collegeBranchForm);
      }

      setVariableForUpdateCollegeBranch(): void {
            this.updateVariable(this.updateForm, "Update College Branch", "Update College Branch", this.updateCollegeBranch)
      }

      updateVariable(formaction: () => void, modalheader: string, modalbutton: string, modalaction: () => void): void {
            this.formHandler = formaction
            this.modalHeader = modalheader
            this.modalButton = modalbutton
            this.modalAction = modalaction
      }

      // modal for add/update form
      openModal(branchDetailModal: any, modalSize?: string): void {

            if (modalSize == undefined) {
                  modalSize = "xl"
            }

            this.modalRef = this.modalService.open(branchDetailModal, {
                  ariaLabelledBy: 'modal-basic-title', keyboard: false,
                  backdrop: 'static', size: modalSize
            })
            /*this.modalRef.result.subscribe((result) => {
            }, (reason) => {
            })*/
      }

      //check form validation
      validate(): void {
            if (this.collegeBranchForm.invalid) {
                  this.collegeBranchForm.markAllAsTouched()
            } else {
                  this.modalAction()
            }
      }

}