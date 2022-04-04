import { DatePipe } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { url } from 'inspector';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { Role, UrlConstant } from 'src/app/service/constant';
import { IDepartmentDTO } from 'src/app/service/department/department.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { ICredential, ITargetCommunity, TargetCommunityService, ICollege, ICompany, ICourse, IFaculty, ITargetCommunityDTO, ITargetCommunityFunction, ISearchFilterField, ISearchSection } from 'src/app/service/target-community/target-community.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-target-community',
  templateUrl: './target-community.component.html',
  styleUrls: ['./target-community.component.css']
})
export class TargetCommunityComponent implements OnInit {

  // Components.
  targetTypeList: any[]
  departmentList: IDepartmentDTO[]
  // credentialList: ICredential[]
  // credentialSearchList: ICredential[]
  salesPersonList: any[]
  studentTypeList: any[]
  targetCommunityFunctionList: ITargetCommunityFunction[]
  targetCommunityFunctionSearchList: ITargetCommunityFunction[]
  collegeBranchList: ICollege[]
  companyBranchList: ICompany[]
  courseList: ICourse[]
  facultyList: IFaculty[]
  talentTypeList: any[]
  ratingNumberList: number[]

  // Flags.
  isOperationUpdate: boolean
  isViewMode: boolean
  showCollegeFieldInForm: boolean
  showCompanyFieldInForm: boolean
  showRatingFieldInForm: boolean
  areCollegesLoading: boolean
  areCompaniesLoading: boolean
  areFunctionsLoading: boolean
  // areCredentialsLoading: boolean
  // showCredentialField: boolean

  // Target community.
  targetCommunityList: ITargetCommunityDTO[]
  targetCommunityForm: FormGroup

  // Pagination.
  limit: number
  currentPage: number
  offset: number
  paginationStart: number
  paginationEnd: number
  totalTargetCommunities: number

  // Modal.
  modalRef: any
  @ViewChild('targetCommunityFormModal') targetCommunityFormModal: any
  @ViewChild('deleteTargetCommunityModal') deleteTargetCommunityModal: any

  // Spinner.



  // Search.
  isSearched: boolean
  targetCommunitySearchForm: FormGroup
  searchFormValue: any
  searchFilterFieldList: ISearchFilterField[]
  searchSectionList: ISearchSection[]
  selectedSectionName: string

  // Permission.
  permission: IPermission
  roleName: string

  constructor(
    private formBuilder: FormBuilder,
    public utilService: UtilityService,
    private generalService: GeneralService,
    private modalService: NgbModal,
    private spinnerService: SpinnerService,
    private targetCommunityService: TargetCommunityService,
    private router: Router,
    private route: ActivatedRoute,
    private datePipe: DatePipe,
    private urlConstant: UrlConstant,
    private localService: LocalService,
    private role: Role,
  ) {
    this.initializeVariables()
    this.getAllComponents()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize all global variables.
  initializeVariables(): void {

    // Components.
    this.targetTypeList = []
    this.targetCommunityList = [] as ITargetCommunityDTO[]
    this.departmentList = [] as IDepartmentDTO[]
    // this.credentialList = [] as ICredential[]
    // this.credentialSearchList = [] as ICredential[]
    this.salesPersonList = []
    this.studentTypeList = []
    this.targetCommunityFunctionList = [] as ITargetCommunityFunction[]
    this.targetCommunityFunctionSearchList = [] as ITargetCommunityFunction[]
    this.collegeBranchList = [] as ICollege[]
    this.companyBranchList = [] as ICompany[]
    this.courseList = [] as ICourse[]
    this.facultyList = [] as IFaculty[]
    this.talentTypeList = []
    this.ratingNumberList = []
    for (let i = 1; i <= 10; i++) {
      this.ratingNumberList.push(i)
    }

    // Flags.
    this.isOperationUpdate = false
    this.isViewMode = false
    this.showCollegeFieldInForm = false
    this.showCompanyFieldInForm = false
    this.showRatingFieldInForm = false
    this.areCollegesLoading = false
    this.areCompaniesLoading = false
    this.areFunctionsLoading = false
    // this.areCredentialsLoading = false
    // this.showCredentialField = false

    // Search.
    this.isSearched = false
    this.searchFilterFieldList = []
    this.searchSectionList = [
      {
        name: "Department",
        isSelected: true
      },
      {
        name: "Target",
        isSelected: false
      },
      {
        name: "Talent",
        isSelected: false
      },
      {
        name: "Valuation",
        isSelected: false
      },
      {
        name: "Others",
        isSelected: false
      }
    ]
    this.selectedSectionName = "Department"

    // Pagination.
    this.limit = 5
    this.offset = 0
    this.currentPage = 0

    // Initialize forms
    this.createTargetCommunitySearchForm()

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Target Communities"


    // Permision.
    // Get permissions from menus using utilityService function.
    this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.TALENT_TARGET_COMMUNITY)
    this.roleName = this.localService.getJsonValue("roleName")
  }

  // =============================================================CREATE FORMS==========================================================================
  // Create target community form.  
  createTargetCommunityForm(): void {
    this.targetCommunityForm = this.formBuilder.group({
      id: new FormControl(null),
      department: new FormControl(null, [Validators.required]),
      // credential: new FormControl(null, [Validators.required]),
      salesPerson: new FormControl(null, [Validators.required]),
      function: new FormControl(null, [Validators.required]),
      targetType: new FormControl(null, [Validators.required]),
      studentType: new FormControl(null, [Validators.required]),
      courses: new FormControl(null, [Validators.required]),
      numberOfBatches: new FormControl(null, [Validators.min(0), Validators.max(255)]),
      targetStudentCount: new FormControl(null, [Validators.min(0), Validators.max(2147483647)]),
      faculty: new FormControl(null, [Validators.required]),
      hours: new FormControl(null, [Validators.min(0), Validators.max(65535)]),
      fees: new FormControl(null, [Validators.min(0), Validators.max(1000000), Validators.pattern(/^[0-9]+(\.[0-9]{1,2})?$/)]),
      targetStartDate: new FormControl(null, [Validators.required]),
      targetEndDate: new FormControl(null, [Validators.required]),
      isTargetAchieved: new FormControl(false, [Validators.required]),
      requiredTalentRating: new FormControl(null, [Validators.required]),
      talentType: new FormControl(null, [Validators.required]),
      minExperienceYears: new FormControl(null, [Validators.min(0), Validators.max(99)]),
      maxExperienceYears: new FormControl(null, [Validators.min(0), Validators.max(99)]),
      salary: new FormControl(null, [Validators.min(0), Validators.max(9999999999), Validators.pattern(/^[0-9]+(\.[0-9]{1,2})?$/)]),
      upSell: new FormControl(null, [Validators.min(0), Validators.max(1000)]),
      crossSell: new FormControl(null, [Validators.min(0), Validators.max(1000)]),
      referral: new FormControl(null, [Validators.min(0), Validators.max(1000)]),
      action: new FormControl(null, [Validators.maxLength(1000)]),
    })
  }

  // Create target community search form.
  createTargetCommunitySearchForm(): void {
    this.targetCommunitySearchForm = this.formBuilder.group({
      targetType: new FormControl(null),
      studentType: new FormControl(null),
      departmentID: new FormControl(null),
      // credentialID: new FormControl(null),
      salesPersonID: new FormControl(null),
      functionID: new FormControl(null),
      facultyID: new FormControl(null),
      numberOfBatches: new FormControl(null, [Validators.min(0)]),
      targetStudentCount: new FormControl(null, [Validators.min(0)]),
      hours: new FormControl(null, [Validators.min(0)]),
      fees: new FormControl(null, [Validators.min(0)]),
      isTargetAchieved: new FormControl(null),
      startDateFromDate: new FormControl(null),
      startDateToDate: new FormControl(null),
      endDateFromDate: new FormControl(null),
      endDateToDate: new FormControl(null),
      collegeIDs: new FormControl(null),
      companyIDs: new FormControl(null),
      courseIDs: new FormControl(null),
      talentType: new FormControl(null),
      minExperienceYears: new FormControl(null, [Validators.min(0)]),
      maxExperienceYears: new FormControl(null, [Validators.min(0)]),
      salary: new FormControl(null, [Validators.min(0), Validators.pattern(/^[0-9]+(\.[0-9]{1,2})?$/)]),
      upSell: new FormControl(null, [Validators.min(0)]),
      crossSell: new FormControl(null, [Validators.min(0)]),
      referral: new FormControl(null, [Validators.min(0)]),
      requiredTalentRating: new FormControl(null),
      rating: new FormControl(null),
    })
    this.disableSearchFormFields()
  }

  // =============================================================TARGET COMMUNITY CRUD FUNCTIONS==========================================================================
  // On clicking add new target community button.
  onAddNewTargetCommunityClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.showCollegeFieldInForm = false
    this.showCompanyFieldInForm = false
    // this.showCredentialField = false
    this.showRatingFieldInForm = false
    this.createTargetCommunityForm()
    // this.targetCommunityForm.get('credential').disable()
    this.targetCommunityForm.get('function').disable()
    this.openModal(this.targetCommunityFormModal, 'xl')
  }

  // Add new target community.
  addTargetCommunity(): void {
    this.spinnerService.loadingMessage = "Adding Target Community"


    let targetCommunity: ITargetCommunity = this.targetCommunityForm.value
    this.patchIDFromObjects(targetCommunity)
    this.targetCommunityService.addTargetCommunity(targetCommunity).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllTargetCommunities()
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

  // On clicking view target community button.
  onViewTargetCommunityClick(targetCommunity: ITargetCommunityDTO): void {
    this.isViewMode = true
    this.createTargetCommunityForm()
    this.showCollegeFieldInForm = false
    this.showCompanyFieldInForm = false
    this.showRatingFieldInForm = false

    // Format dates.
    targetCommunity.targetStartDate = this.datePipe.transform(targetCommunity.targetStartDate, 'yyyy-MM-dd')
    targetCommunity.targetEndDate = this.datePipe.transform(targetCommunity.targetEndDate, 'yyyy-MM-dd')

    // Target type is college.
    if (targetCommunity.colleges && targetCommunity.colleges?.length > 0) {
      this.targetCommunityForm.addControl('colleges', this.formBuilder.control(null, [Validators.required]))
      this.targetCommunityForm.removeControl('companies')
      this.showCollegeFieldInForm = true
      this.showCompanyFieldInForm = false
      this.targetCommunityForm.addControl('rating', this.formBuilder.control(null, [Validators.required]))
      this.showRatingFieldInForm = true
    }

    // Target type is company.
    if (targetCommunity.companies && targetCommunity.companies?.length > 0) {
      this.targetCommunityForm.addControl('companies', this.formBuilder.control(null, [Validators.required]))
      this.targetCommunityForm.removeControl('colleges')
      this.showCollegeFieldInForm = false
      this.showCompanyFieldInForm = true
      this.targetCommunityForm.addControl('rating', this.formBuilder.control(null, [Validators.required]))
      this.showRatingFieldInForm = true
    }

    // Get function list.
    this.getFunctionListByRoleInForm(targetCommunity.department)

    // Get credential lit.
    // this.getCredentialListByDepartmentRoleInForm(targetCommunity.department)

    this.targetCommunityForm.patchValue(targetCommunity)
    this.targetCommunityForm.disable()
    this.openModal(this.targetCommunityFormModal, 'xl')
  }

  // On cliking update form button in target community form.
  onUpdateTargetCommunityClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = true
    this.targetCommunityForm.enable()
  }

  // Update target community.
  updateTargetCommunity(): void {
    this.spinnerService.loadingMessage = "Updating Target Community"


    let targetCommunity: ITargetCommunity = this.targetCommunityForm.value
    this.patchIDFromObjects(targetCommunity)
    this.targetCommunityService.updateTargetCommunity(targetCommunity).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllTargetCommunities()
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

  // On cicking target achieved button in list.
  onAchievedButtonClick(targetCommunityID: string): void {
    if (confirm("Are you sure you want to set the target community as achieved?")) {
      this.spinnerService.loadingMessage = "Updating Target Community"


      this.targetCommunityService.updateTargetCommunityIsTargetAchieved({ "targetCommunityID": targetCommunityID }).subscribe((response: any) => {
        this.getAllTargetCommunities()
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
  }

  // On clicking delete target community button. 
  onDeleteTargetCommunityClick(targetCommunityID: string): void {
    this.openModal(this.deleteTargetCommunityModal, 'md').result.then(() => {
      this.deleteTargetCommunity(targetCommunityID)
    }, (err) => {
      console.error(err)
      return
    })
  }

  // Delete target community after confirmation from user.
  deleteTargetCommunity(targetCommunityID: string): void {
    this.spinnerService.loadingMessage = "Deleting Target Community"


    this.targetCommunityService.deleteTargetCommunity(targetCommunityID).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllComponents()
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

  // =============================================================TARGET COMMUNITY SEARCH FUNCTIONS==========================================================================
  // Reset search form and renaviagte page.
  resetSearchAndGetAll(): void {
    this.searchFilterFieldList = []
    this.targetCommunitySearchForm.reset()
    this.searchFormValue = null
    this.changePage(1)
    this.isSearched = false
    this.router.navigate([this.urlConstant.TALENT_TARGET_COMMUNITY])
  }

  // Reset search form.
  resetSearchForm(): void {
    this.searchFilterFieldList = []
    this.targetCommunitySearchForm.reset()
  }

  // Search target communities.
  searchTargetCommunities(): void {
    this.searchFormValue = { ...this.targetCommunitySearchForm?.value }
    if (this.searchFormValue.departmentID) {
      this.searchFormValue.departmentID = this.searchFormValue.departmentID.id
    }
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
    this.spinnerService.loadingMessage = "Searching Target Communities"
    this.changePage(1)
  }

  // Disable dependent field controls on search form creation.
  disableSearchFormFields(): void {
    this.targetCommunitySearchForm.get('functionID').disable()
    // this.targetCommunitySearchForm.get('credentialID').disable()
  }

  // On department value change in search form.
  onDepartmentChangeInSearchForm(department: any): void {
    // this.getCredentialListByDepartmentRoleInSearchForm(department)
    this.getFunctionListByRoleInSearchForm(department)
  }

  // On clicking search section name.
  onSearchSectionNameClick(sectionName: string): void {
    for (let i = 0; i < this.searchSectionList.length; i++) {
      if (this.searchSectionList[i].name == sectionName) {
        this.searchSectionList[i].isSelected = true
        this.selectedSectionName = this.searchSectionList[i].name
      }
      else {
        this.searchSectionList[i].isSelected = false
      }
    }
  }

  // Delete search criteria from target community search form by search name.
  deleteSearchCriteria(searchName: string): void {
    this.targetCommunitySearchForm.get(searchName).setValue(null)
    this.searchTargetCommunities()
  }


  // ================================================OTHER FUNCTIONS FOR TARGET COMMUNITIES===============================================
  // Page change.
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
    this.getAllTargetCommunities()
  }

  // Checks the url's query params and decides whether to call get or search.
  searchOrGetTargetCommunities(): void {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getAllTargetCommunities()
      return
    }
    this.targetCommunitySearchForm.patchValue(queryParams)
    this.searchTargetCommunities()
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

  // On clicking sumbit button in target community form.
  onFormSubmit(): void {
    if (this.targetCommunityForm.invalid) {
      this.targetCommunityForm.markAllAsTouched()
      return
    }
    if (this.isOperationUpdate) {
      this.updateTargetCommunity()
      return
    }
    this.addTargetCommunity()
  }

  // Set total target community list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalTargetCommunities < this.paginationEnd) {
      this.paginationEnd = this.totalTargetCommunities
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

  // On department value change in form.
  onDepartmentChangeInForm(department: any): void {
    // this.getCredentialListByDepartmentRoleInForm(department)
    this.getFunctionListByRoleInForm(department)
  }

  // On target type value change in form.
  onTargetTypeChange(targetType: string): void {
    if (!targetType) {
      this.showCollegeFieldInForm = false
      this.showCompanyFieldInForm = false
      return
    }
    // Target type is company.
    if (targetType == "Company") {
      this.targetCommunityForm.addControl('companies', this.formBuilder.control(null, [Validators.required]))
      this.targetCommunityForm.removeControl('colleges')
      this.showCollegeFieldInForm = false
      this.showCompanyFieldInForm = true
      if (!this.targetCommunityForm.contains('rating')) {
        this.targetCommunityForm.addControl('rating', this.formBuilder.control(null, [Validators.required]))
        this.showRatingFieldInForm = true
        return
      }
      this.targetCommunityForm.get('rating').setValue(null)
      return
    }
    // Target type is college.
    if (targetType == "College") {
      this.targetCommunityForm.addControl('colleges', this.formBuilder.control(null, [Validators.required]))
      this.targetCommunityForm.removeControl('companies')
      this.showCollegeFieldInForm = true
      this.showCompanyFieldInForm = false
      if (!this.targetCommunityForm.contains('rating')) {
        this.targetCommunityForm.addControl('rating', this.formBuilder.control(null, [Validators.required]))
        this.showRatingFieldInForm = true
        return
      }
      this.targetCommunityForm.get('rating').setValue(null)
      return
    }

    // Target type is retail.
    if (targetType == "Retail") {
      this.targetCommunityForm.removeControl('colleges')
      this.targetCommunityForm.removeControl('companies')
      this.showCollegeFieldInForm = false
      this.showCompanyFieldInForm = false
      this.targetCommunityForm.removeControl('rating')
      this.showRatingFieldInForm = false
    }
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

  // Format all fields of target community.
  formatTargetCommunityFields(): void {
    for (let i = 0; i < this.targetCommunityList.length; i++) {
      if (this.targetCommunityList[i].fees) {
        this.targetCommunityList[i].feesInString = this.formatFeesInIndianRupeeSystem(this.targetCommunityList[i].fees)
      }
      if (this.targetCommunityList[i].salary) {
        this.targetCommunityList[i].salaryInString = this.formatFeesInIndianRupeeSystem(this.targetCommunityList[i].salary)
      }
    }
  }

  // Extract id from objects and give it to target community.
  patchIDFromObjects(targetCommunity: ITargetCommunity): void {
    if (this.targetCommunityForm.get('department').value) {
      targetCommunity.departmentID = this.targetCommunityForm.get('department').value.id
      delete targetCommunity['department']
    }
    // if (this.targetCommunityForm.get('credential').value) {
    //   targetCommunity.credentialID = this.targetCommunityForm.get('credential').value.id
    //   delete targetCommunity['credential']
    // }
    if (this.targetCommunityForm.get('function').value) {
      targetCommunity.functionID = this.targetCommunityForm.get('function').value.id
      delete targetCommunity['function']
    }
    if (this.targetCommunityForm.get('faculty').value) {
      targetCommunity.facultyID = this.targetCommunityForm.get('faculty').value.id
      delete targetCommunity['faculty']
    }
    if (this.targetCommunityForm.get('salesPerson').value) {
      targetCommunity.salesPersonID = this.targetCommunityForm.get('salesPerson').value.id
      delete targetCommunity['salesPerson']
    }
    // targetCommunity.requiredTalentRating = +targetCommunity['requiredTalentRating']
    // delete targetCommunity['requiredTalentRating']

    // if (this.targetCommunityForm.contains('rating')){
    //   targetCommunity.rating = +targetCommunity['rating']
    //   delete targetCommunity['rating']
    // }
  }

  // =============================================================GET FUNCTIONS==========================================================================
  // Get all components.
  getAllComponents(): void {
    this.getTargetTypeList()
    this.getDepartmentList()
    this.getStudentTypeList()
    this.getCourseList()
    this.searchOrGetTargetCommunities()
    this.getFacultyList()
    this.getSalesPersonList()
    this.getCollegeBranchList()
    this.getCompanyBranchList()
    this.getTalentTypeList()
  }

  // Get target type list.
  getTargetTypeList(): void {
    this.generalService.getGeneralTypeByType("target_type").subscribe((response: any) => {
      this.targetTypeList = response
    }, (err: any) => {
      console.error(err)
    })
  }

  // Get department list.
  getDepartmentList(): void {
    let queryParams: any = {
      roleNames: ["Salesperson", "Faculty"]
    }
    this.generalService.getDepartmentList(queryParams).subscribe((response: any) => {
      this.departmentList = response
    }, (err: any) => {
      console.error(err)
    })
  }

  // Get student type list.
  getStudentTypeList(): void {
    this.generalService.getGeneralTypeByType("target_community_student_type").subscribe((response: any) => {
      this.studentTypeList = response
    }, (err: any) => {
      console.error(err)
    })
  }

  // // Get credential list by department role name in target community form.
  // getCredentialListByDepartmentRoleInForm(department: any): void {
  //   if (!department) {
  //     this.showCredentialField = false
  //     this.targetCommunityForm.get('credential').setValue(null)
  //     this.credentialList = []
  //     this.targetCommunityForm.get('credential').disable()
  //     return
  //   }
  //   if (!this.isViewMode) {
  //     this.targetCommunityForm.get('credential').enable()
  //   }
  //   this.areCredentialsLoading = true
  //   this.showCredentialField = true
  //   this.targetCommunityForm.get('credential').setValue(null)
  //   this.credentialList = []
  //   let queryParams: any = {
  //     roleNames: [department.role.roleName]
  //   }
  //   this.generalService.getCredentialListByRole(queryParams).subscribe((response: any) => {
  //     this.credentialList = response
  //     this.areCredentialsLoading = false
  //   }, (err: any) => {
  //     this.areCredentialsLoading = false
  //     console.error(err)
  //   })
  // }

  // // Get credential list by department role name in target community search form.
  // getCredentialListByDepartmentRoleInSearchForm(department: any): void {
  //   if (!department) {
  //     this.targetCommunitySearchForm.get('credentialID').setValue(null)
  //     this.credentialSearchList = []
  //     this.targetCommunitySearchForm.get('credentialID').disable()
  //     return
  //   }
  //   this.targetCommunitySearchForm.get('credentialID').enable()
  //   this.areCredentialsLoading = true
  //   this.targetCommunitySearchForm.get('credentialID').setValue(null)
  //   this.credentialSearchList = []
  //   let queryParams: any = {
  //     roleNames: [department.role.roleName]
  //   }
  //   this.generalService.getCredentialListByRole(queryParams).subscribe((response: any) => {
  //     this.credentialSearchList = response
  //     this.areCredentialsLoading = false
  //   }, (err: any) => {
  //     this.areCredentialsLoading = false
  //     console.error(err)
  //   })
  // }

  // Get target community function list by department id in target community form.
  getFunctionListByRoleInForm(department: any): void {
    if (!department) {
      this.targetCommunityForm.get('function').setValue(null)
      this.targetCommunityFunctionList = []
      this.targetCommunityForm.get('function').disable()
      return
    }
    if (!this.isViewMode) {
      this.targetCommunityForm.get('function').enable()
    }
    this.areFunctionsLoading = true
    this.targetCommunityForm.get('function').setValue(null)
    this.targetCommunityFunctionList = []
    this.generalService.getTargetCommunityFunctionByDepartemnt(department.id).subscribe((response: any) => {
      this.targetCommunityFunctionList = response
      this.areFunctionsLoading = false
    }, (err: any) => {
      this.areFunctionsLoading = false
      console.error(err)
    })
  }

  // Get target community function list by department id in target community search form.
  getFunctionListByRoleInSearchForm(department: any): void {
    if (!department) {
      this.targetCommunitySearchForm.get('functionID').setValue(null)
      this.targetCommunityFunctionSearchList = []
      this.targetCommunitySearchForm.get('functionID').disable()
      return
    }
    this.targetCommunitySearchForm.get('functionID').enable()
    this.areFunctionsLoading = true
    this.targetCommunitySearchForm.get('functionID').setValue(null)
    this.targetCommunityFunctionSearchList = []
    this.generalService.getTargetCommunityFunctionByDepartemnt(department.id).subscribe((response: any) => {
      this.targetCommunityFunctionSearchList = response
      this.areFunctionsLoading = false
    }, (err: any) => {
      this.areFunctionsLoading = false
      console.error(err)
    })
  }

  // Get all target communities.
  getAllTargetCommunities(): void {
    this.spinnerService.loadingMessage = "Getting All Target Communities"


    this.targetCommunityService.getAllTargetCommunities(this.limit, this.offset, this.searchFormValue).subscribe((response) => {
      this.totalTargetCommunities = response.headers.get('X-Total-Count')
      this.targetCommunityList = response.body
      this.formatTargetCommunityFields()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.setPaginationString()
    })
  }

  // Get college branch list.
  getCollegeBranchList(): void {
    this.areCollegesLoading = true
    this.generalService.getCollegeBranchList().subscribe((response: any) => {
      this.collegeBranchList = response
      this.areCollegesLoading = false
    }, (err: any) => {
      this.areCollegesLoading = false
      console.error(err)
    })
  }

  // Get company branch list.
  getCompanyBranchList(): void {
    this.areCompaniesLoading = true
    this.generalService.getCompanyBranchList().subscribe((response: any) => {
      this.companyBranchList = response.body
      this.areCompaniesLoading = false
    }, (err: any) => {
      this.areCompaniesLoading = false
      console.error(err)
    })
  }

  // Get course list.
  getCourseList(): void {
    this.generalService.getCourseList().subscribe((response: any) => {
      this.courseList = response.body
    }, (err: any) => {
      console.error(err)
    })
  }

  // Get faculty list.
  getFacultyList(): void {
    this.generalService.getFacultyList().subscribe((response: any) => {
      this.facultyList = response.body
    }, (err: any) => {
      console.error(err)
    })
  }

  // Get salesPerson list.
  getSalesPersonList(): void {
    this.generalService.getSalesPersonList().subscribe((response: any) => {
      this.salesPersonList = response.body
    }, (err: any) => {
      console.error(err)
    })
  }

  // Get talent type list.
  getTalentTypeList(): void {
    this.generalService.getGeneralTypeByType("talent_rating").subscribe((respond: any[]) => {
      this.talentTypeList = respond
    }, (err) => {
      console.error(this.utilService.getErrorString(err))
    })
  }

}
