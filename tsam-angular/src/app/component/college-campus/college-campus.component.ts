import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, Validators, FormControl } from '@angular/forms';
import { IPermission } from 'src/app/service/menu/menu.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { GeneralService } from 'src/app/service/general/general.service';
import { CollegeCampusService, ICampusDrive, ISearchFilterField, IUpdateMultipleCandidate } from 'src/app/service/college/college-campus/college-campus.service';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { UrlConstant } from 'src/app/service/constant';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { DatePipe } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { DegreeService } from 'src/app/service/degree/degree.service';

@Component({
      selector: 'app-college-campus',
      templateUrl: './college-campus.component.html',
      styleUrls: ['./college-campus.component.css']
})
export class CollegeCampusComponent implements OnInit {

      //****************************CAMPUS DRIVE*************************************** */

      // Components.
      companyRequirementList: any[]
      salesPersonList: any[]
      collegeBranchList: any[]
      facultyList: any[]
      developerList: any[]

      // Flags.
      isViewMode: boolean
      isOperationUpdate: boolean

      // Campus Drive.
      campusDriveList: any[]
      campusDriveForm: FormGroup
      selectedCampusdriveID: string

      // Pagination.
      limitCampusDrive: number
      offsetCampusDrive: number
      currentPageCampusDrive: number
      totalCampusDrives: number
      paginationStartCampusDrive: number
      paginationEndCampusDrive: number

      // Modal.
      modalRef: any
      @ViewChild('campusDriveFormModal') campusDriveFormModal: any
      @ViewChild('deleteCampusDriveModal') deleteCampusDriveModal: any

      // Spinner.



      // Search.
      campusDriveSearchForm: FormGroup
      isSearchedCampusDrive: boolean
      searchFormValueCampusDrive: any
      searchCampusDriveFilterFieldList: ISearchFilterField[]

      // Campus drive retaled variables.
      currentYear: number
      activeTab: number
      resumeURL: string

      // Permission.
      permission: IPermission
      roleName: string
      showForSalesPersonLogin: boolean
      showForFacultyLogin: boolean
      showForDeveloperLogin: boolean

      // Excel.
      arrayBuffer: any
      file: File
      excelData: any
      talentData: any[]
      errorMsg: string
      talentError: boolean
      upload: boolean
      excelFileName: string

      //****************************CANDIDATE*************************************** */
      // Components.
      stateList: any[]
      countryList: any[]
      qualificationList: any[]
      specializationList: any[]
      academicYearList: any[]
      selectedCampusDriveCollegeList: any[]
      candidateResultList: any[]

      // Flags.
      showRegisterCandidateForm: boolean
      showCampusTalentRegistrationFields: boolean
      showCandidatesList: boolean

      // Candidate.
      candidateList: any[]
      selectedCandidateList: string[]
      candidateForm: FormGroup

      // Pagination.
      limitCandidate: number
      offsetCandidate: number
      currentPageCandidate: number
      totalCandidates: number
      paginationStartCandidate: number
      paginationEndCandidate: number

      // Modal.
      @ViewChild('candidateFormModal') candidateFormModal: any

      // Search.
      isSearchedCandidate: boolean
      candidateSearchForm: FormGroup
      searchFormValueCandidate: any
      searchCandidateFilterFieldList: ISearchFilterField[]

      // Resume.
      isResumeUploadedToServer: boolean
      isFileUploading: boolean
      docStatus: string
      displayedFileName: string

      //****************************UPDATE MULTIPLE CANDIDATE *************************************** */
      // Flags.
      multipleSelect: boolean

      // Modal.
      @ViewChild('updateMultipleCandidateModal') updateMultipleCandidateModal: any

      // Forms.
      updateMultipleCandidateForm: FormGroup
      updateMultipleCandidateFormValue: IUpdateMultipleCandidate

      constructor(
            private formBuilder: FormBuilder,
            public utilityService: UtilityService,
            private localService: LocalService,
            private spinnerService: SpinnerService,
            private generalService: GeneralService,
            private degreeService: DegreeService,
            private collegeCampusService: CollegeCampusService,
            private urlConstant: UrlConstant,
            private modalService: NgbModal,
            private fileOperationService: FileOperationService,
            private datePipe: DatePipe,
            private router: Router,
            private route: ActivatedRoute,
      ) {
            this.initializeVariables()
            this.getAllComponents()

      }


      get ongoingOperations() {
            return this.spinnerService.ongoingOperations
      }

      ngOnInit(): void { }

      //Initialize all global variables
      initializeVariables(): void {
            //****************************CAMPUS DRIVE*************************************** */
            // Components.
            this.companyRequirementList = []
            this.salesPersonList = []
            this.collegeBranchList = []
            this.facultyList = []
            this.developerList = []

            // Flags.
            this.isViewMode = false
            this.isOperationUpdate = false

            // Pagination.
            this.limitCampusDrive = 5
            this.offsetCampusDrive = 0
            this.currentPageCampusDrive = 0
            this.paginationStartCampusDrive = 0
            this.paginationEndCampusDrive = 0

            // Spinner.
            this.spinnerService.loadingMessage = "Getting All Campus Drives"


            // Search.
            this.isSearchedCampusDrive = false
            this.searchFormValueCampusDrive = {}
            this.searchCampusDriveFilterFieldList = []

            // Campus drive retaled variables.
            this.activeTab = 1
            this.currentYear = new Date().getFullYear()

            // Permission.
            this.showForSalesPersonLogin = true
            this.showForFacultyLogin = true
            this.showForDeveloperLogin = true
            // Get permissions from menus using utilityService function.
            this.permission = this.utilityService.getPermission(this.urlConstant.COLLEGE_CAMPUS_DRIVE)
            // Get role name for menu for calling their specific apis.
            this.roleName = this.localService.getJsonValue("roleName")
            // Hide features for salesperson.
            if (this.roleName == "Salesperson") {
                  this.showForSalesPersonLogin = false
            }
            // Hide features for salesperson.
            if (this.roleName == "Faculty") {
                  this.showForFacultyLogin = false
            }
            // Hide features for salesperson.
            if (this.roleName == "developer") {
                  this.showForDeveloperLogin = false
            }

            // Excel.
            this.talentData = []
            this.errorMsg = ""
            this.talentError = false
            this.upload = false
            this.excelFileName = ""

            //****************************CANDIDATE*************************************** */
            // Components.
            this.stateList = []
            this.countryList = []
            this.qualificationList = []
            this.specializationList = []
            this.academicYearList = []
            this.selectedCampusDriveCollegeList = []
            this.candidateResultList = []

            // Flags.
            this.showRegisterCandidateForm = false
            this.showCampusTalentRegistrationFields = false
            this.showCandidatesList = false

            // Candidate.
            this.selectedCandidateList = []

            // Pagination.
            this.limitCandidate = 5
            this.offsetCandidate = 0
            this.currentPageCandidate = 0
            this.paginationStartCandidate = 0
            this.paginationEndCandidate = 0

            // Search.
            this.isSearchedCandidate = false
            this.searchFormValueCandidate = {}
            this.searchCandidateFilterFieldList = []

            // Resume.
            this.docStatus = ""
            this.displayedFileName = "Select file"
            this.isResumeUploadedToServer = false
            this.isFileUploading = false

            //****************************UPDATE MULTIPLE CANDIDATE *************************************** */
            // Flags.
            this.multipleSelect = false

            //****************************INITIALIZE FORMS*************************************** */
            this.createCampusDriveSearchForm()
      }

      //*********************************************CREATE FORMS************************************************************
      // Create new campus drive form.
      createCampusDriveForm(): void {
            this.campusDriveForm = this.formBuilder.group({
                  id: new FormControl(null),
                  code: new FormControl({ value: null, disabled: true }),
                  campusName: new FormControl(null, [Validators.required]),
                  description: new FormControl(null),
                  location: new FormControl(null),
                  salesPeople: new FormControl(Array()),
                  faculties: new FormControl(Array()),
                  developers: new FormControl(Array()),
                  collegeBranches: new FormControl(Array(), [Validators.required]),
                  totalRegisteredCandidates: new FormControl(null),
                  totalAppearedCandidates: new FormControl(null),
                  totalRequirements: new FormControl(null),
                  campusDate: new FormControl(null, [Validators.required]),
                  studentRegistrationLink: new FormControl(null),
                  companyRequirements: new FormControl(Array(), [Validators.required]),
                  cancelled: new FormControl(false),
            })
      }

      // Create new candidate form.
      createCandidateForm(): void {
            this.candidateForm = this.formBuilder.group({
                  firstName: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z]*$/)]),
                  lastName: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z]*$/)]),
                  email: new FormControl(null, [Validators.pattern(/^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$/), Validators.required]),
                  contact: new FormControl(null, [Validators.pattern(/^[6789]\d{9}$/), Validators.required]),
                  address: new FormControl(null, [Validators.pattern(/^[.0-9a-zA-Z\s,-\/]+$/)]),
                  city: new FormControl(null, [Validators.pattern(/^[a-zA-Z]+([a-zA-Z ]?)+$/)]),
                  state: new FormControl(null),
                  country: new FormControl(null),
                  pinCode: new FormControl(null, [Validators.pattern(/^[1-9][0-9]{5}$/)]),
                  academicYear: new FormControl(null, [Validators.required]),
                  isSwabhavTalent: new FormControl(null, [Validators.required]),
                  college: new FormControl(null, [Validators.required, Validators.maxLength(200)]),
                  percentage: new FormControl(null, [Validators.min(1), Validators.pattern(/(^100([.]0{1,2})?)$|(^\d{1,2}([.]\d{1,2})?)$/i), Validators.required]),
                  passout: new FormControl(null, [Validators.min(1980), Validators.max(this.currentYear + 3), Validators.required]),
                  degree: new FormControl(null, Validators.required),
                  specialization: new FormControl(null, Validators.required),
                  talentID: new FormControl(null),
                  campusDriveID: new FormControl(null),
                  campusTalentRegistrationID: new FormControl(null),
                  talentAcademicID: new FormControl(null),
                  resume: new FormControl(null),
                  isTestLinkSent: new FormControl(false, [Validators.required]),
                  hasAttempted: new FormControl(false, [Validators.required])
            })
      }

      // Create new campus drive search form.
      createCampusDriveSearchForm(): void {
            this.campusDriveSearchForm = this.formBuilder.group({
                  campusName: new FormControl(null),
                  fromDate: new FormControl(null),
                  toDate: new FormControl(null),
                  collegeIDs: new FormControl(null),
                  salesPersonIDs: new FormControl(null),
                  facultyIDs: new FormControl(null),
                  developerIDs: new FormControl(null),
                  companyRequirementIDs: new FormControl(null),
            })
      }

      // Create new update multiple candidates form.
      createUpdateMultipleCandidateForm(): void {
            this.updateMultipleCandidateForm = this.formBuilder.group({
                  isTestLinkSent: new FormControl(null),
                  hasAttempted: new FormControl(null),
                  result: new FormControl(null),
                  campusDriveID: new FormControl(null),
                  campusTalentRegistrationIDs: new FormControl(null),
            })
      }

      // Create candidate search form.
      createCandidateSearchForm(): void {
            this.candidateSearchForm = this.formBuilder.group({
                  firstName: new FormControl(null, [Validators.pattern(/^[a-zA-Z]*$/)]),
                  lastName: new FormControl(null, [Validators.pattern(/^[a-zA-Z]*$/)]),
                  email: new FormControl(null),
                  contact: new FormControl(null),
                  degreeID: new FormControl(null),
                  specializationID: new FormControl(null),
                  collegeID: new FormControl(null),
                  academicYear: new FormControl(null),
                  isSwabhavTalent: new FormControl(null),
                  isTestLinkSent: new FormControl(null),
                  hasAttempted: new FormControl(null),
                  result: new FormControl(null),
                  fromDate: new FormControl(null),
                  toDate: new FormControl(null),
            })
      }

      //*********************************************CRUD FUNCTIONS FOR CAMPUS DRIVE FORM************************************************************
      // On add new campus drive button click.
      onAddNewCampusDriveButtonClick(): void {
            this.isViewMode = false
            this.isOperationUpdate = false
            this.createCampusDriveForm()
            this.enableCampusDriveForm()
            this.openModal(this.campusDriveFormModal, 'xl')
      }

      //On clicking add campus drive in campus drive form.
      addCampusDrive(): void {
            this.spinnerService.loadingMessage = "Adding Campus Drive"


            this.collegeCampusService.addCampusDrive(this.campusDriveForm.value).subscribe((response: any) => {
                  this.modalRef.close('success')
                  this.getCampusDrives()
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

      // On clicking view campus drive button.
      onViewCampusDriveClick(campusDrive: ICampusDrive): void {
            this.isViewMode = true
            this.createCampusDriveForm()
            // Campus drive date.
            let campusDate = campusDrive.campusDate
            if (campusDate) {
                  campusDrive.campusDate = this.datePipe.transform(campusDate, 'yyyy-MM-ddTHH:mm:ss')
            }
            this.campusDriveForm.patchValue(campusDrive)
            this.campusDriveForm.disable()
            this.openModal(this.campusDriveFormModal, 'xl')
      }

      // On cliking update form button in campus drive form.
      onUpdateCampusDriveClick(): void {
            this.isViewMode = false
            this.isOperationUpdate = true
            this.enableCampusDriveForm()
      }

      // On clicking update campus drive in campus drive form.
      updateCampusDrive(): void {
            this.spinnerService.loadingMessage = "Updating Campus Drive"


            this.collegeCampusService.updateCampusDrive(this.campusDriveForm.value).subscribe((response: any) => {
                  this.modalRef.close()
                  this.getCampusDrives()
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

      // On clicking delete campus drive button. 
      onDeleteCampusDriveClick(campusDriveID: string): void {
            this.openModal(this.deleteCampusDriveModal, 'md').result.then(() => {
                  this.deleteCampusDrive(campusDriveID)
            }, (err) => {
                  console.error(err)
                  return
            })
      }

      // Delete campus drive.
      deleteCampusDrive(campusDriveID: string): void {
            this.spinnerService.loadingMessage = "Deleting Campus Drive"


            this.collegeCampusService.deleteCampusDrive(campusDriveID).subscribe((response: any) => {
                  this.modalRef.close()
                  this.getCampusDrives()
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

      // ==================================================CAMPUS DRIVE SEARCH FUNCTIONS==========================================================================
      // Reset campus drive search form and renaviagte page.
      resetSearchAndGetAllForCampusDrive(): void {
            this.searchCampusDriveFilterFieldList = []
            this.campusDriveSearchForm.reset()
            this.searchFormValueCampusDrive = {}
            this.changePageForCampusDrive(1)
            this.isSearchedCampusDrive = false
            this.router.navigate([this.urlConstant.COLLEGE_CAMPUS_DRIVE])
      }

      // Reset campus drive search form.
      resetCampusDriveSearchForm(): void {
            this.searchCampusDriveFilterFieldList = []
            this.campusDriveSearchForm.reset()
      }

      // Search campus drives.
      searchCampusDrives(): void {
            this.searchFormValueCampusDrive = { ...this.campusDriveSearchForm?.value }
            this.router.navigate([], {
                  relativeTo: this.route,
                  queryParams: this.searchFormValueCampusDrive,
            })
            for (let field in this.searchFormValueCampusDrive) {
                  if (this.searchFormValueCampusDrive[field] === null || this.searchFormValueCampusDrive[field] === "") {
                        delete this.searchFormValueCampusDrive[field]
                  } else {
                        this.isSearchedCampusDrive = true
                  }
            }
            this.searchCampusDriveFilterFieldList = []
            for (var property in this.searchFormValueCampusDrive) {
                  let text: string = property
                  let result: string = text.replace(/([A-Z])/g, " $1");
                  let finalResult: string = result.charAt(0).toUpperCase() + result.slice(1);
                  let valueArray: any[] = []
                  if (Array.isArray(this.searchFormValueCampusDrive[property])) {
                        valueArray = this.searchFormValueCampusDrive[property]
                  }
                  else {
                        valueArray.push(this.searchFormValueCampusDrive[property])
                  }
                  this.searchCampusDriveFilterFieldList.push(
                        {
                              propertyName: property,
                              propertyNameText: finalResult,
                              valueList: valueArray
                        })
            }
            if (this.searchCampusDriveFilterFieldList.length == 0) {
                  this.resetSearchAndGetAllForCampusDrive()
            }
            if (!this.isSearchedCampusDrive) {
                  return
            }
            this.spinnerService.loadingMessage = "Searching Campus Drives"
            this.changePageForCampusDrive(1)
      }

      // Delete search criteria from campus drive search form by search name.
      deleteSearchCampusDriveCriteria(searchName: string): void {
            this.campusDriveSearchForm.get(searchName).setValue(null)
            this.searchCampusDrives()
      }

      //*********************************************OTHER FUNCTIONS FOR CAMPUS DRIVE************************************************************
      // Enable the campus drive form.
      enableCampusDriveForm(): void {
            this.campusDriveForm.enable()
            this.campusDriveForm.get('code').disable()
            this.campusDriveForm.get('totalRequirements').disable()
            this.campusDriveForm.get('totalRegisteredCandidates').disable()
            this.campusDriveForm.get('totalAppearedCandidates').disable()
      }

      // Change page for pagination for campus drive.
      changePageForCampusDrive($event): void {
            this.offsetCampusDrive = $event - 1
            this.currentPageCampusDrive = $event
            this.getCampusDrives()
      }

      // On clicking sumbit button in campus drive form.
      onCampusDriveFormSubmit(): void {
            if (this.campusDriveForm.invalid) {
                  this.campusDriveForm.markAllAsTouched()
                  return
            }
            if (this.isOperationUpdate) {
                  this.updateCampusDrive()
                  return
            }
            this.addCampusDrive()
      }

      // Redirect to talents page filtered by campus drive id.
      redirectToTalentsForRegisteredCandidates(campusDriveID: string): void {
            this.router.navigate([this.urlConstant.TALENT_MASTER], {
                  queryParams: {
                        "campusDriveID": campusDriveID,
                        "hasAppeared": "0"
                  }
            }).catch(err => {
                  console.error(err)

            });
      }

      // Redirect to talents page filtered by campus drive id and has appeared fields.
      redirectToTalentsForAppearedCandidates(campusDrive: any): void {
            this.router.navigate([this.urlConstant.TALENT_MASTER], {
                  queryParams: {
                        "campusDriveID": campusDrive.id,
                        "hasAppeared": "1"
                  }
            }).catch(err => {
                  console.error(err)

            });
      }

      // Checks the url's query params and decides to call whether to call get or search.
      searchOrGetCampusDrives() {
            let queryParams = this.route.snapshot.queryParams
            if (this.utilityService.isObjectEmpty(queryParams)) {
                  this.getCampusDrives()
                  return
            }
            this.campusDriveSearchForm.patchValue(queryParams)
            this.searchCampusDrives()
      }

      // Set total list on current page.
      setPaginationStringForCampusDrive(): void {
            this.paginationStartCampusDrive = this.limitCampusDrive * this.offsetCampusDrive + 1
            this.paginationEndCampusDrive = +this.limitCampusDrive + this.limitCampusDrive * this.offsetCampusDrive
            if (this.totalCampusDrives < this.paginationEndCampusDrive) {
                  this.paginationEndCampusDrive = this.totalCampusDrives
            }
      }

      //*********************************************CRUD FUNCTIONS FOR CANDIDATE FORM************************************************************
      // On clicking register students button.
      getCandidatesBySelectedCampusDrive(campusDriveID: any): void {
            this.selectedCampusdriveID = campusDriveID
            this.getSelectedCampusDriveCollegeList()
            this.createCandidateSearchForm()
            this.initializeCandidateVariables()
            this.getAllCandidates()
            this.openModal(this.candidateFormModal, 'xl')
      }

      // Get all candidates by campus drive id.
      getAllCandidates(): void {
            this.spinnerService.loadingMessage = "Getting All Candidates"


            this.collegeCampusService.getCandidatesByCampusDrive(this.selectedCampusdriveID,
                  this.limitCandidate, this.offsetCandidate, this.searchFormValueCandidate).subscribe(response => {
                        this.candidateList = response.body
                        this.formatAddress()
                        this.totalCandidates = parseInt(response.headers.get("X-Total-Count"))
                        this.getCampusDrives()
                  }, err => {
                        console.error(this.utilityService.getErrorString(err))
                  }).add(() => {

                        this.setPaginationStringForCandidate()
                  })
      }

      // On clicking add new candidate button.
      onAddNewCandidateButtonClick(): void {
            this.createCandidateForm()
            this.isViewMode = false
            this.isOperationUpdate = false
            this.showCampusTalentRegistrationFields = false
            this.stateList = []
            this.candidateForm.get('state').disable()
            this.specializationList = []
            this.candidateForm.get('specialization').disable()
            this.toggleVisibilityOfCanditateFormAndList()
            this.activeTab = 1
      }

      // On clicking add candidate button on candidate form.
      addCandidate(): void {
            this.spinnerService.loadingMessage = "Adding Candidate"


            let candidate: any = this.candidateForm.value
            this.patchIDFromObjectsForCandidate(candidate)
            this.collegeCampusService.addCandidate(this.candidateForm.value, this.selectedCampusdriveID).subscribe((response: any) => {
                  this.toggleVisibilityOfCanditateFormAndList()
                  this.getAllCandidates()
                  alert(response)
            }, (error) => {
                  console.error(error)
                  if (typeof error.error == 'object' && error) {
                        alert(this.utilityService.getErrorString(error))
                        return
                  }
                  if (error.error == undefined) {
                        alert('Candidate could not be added, try again')
                        return
                  }
                  alert(error.statusText)
            })
      }

      // On clicking view candidate button.
      OnViewCandidateButtonClick(index: number): void {
            this.createCandidateForm()
            this.isViewMode = true
            this.addCampusTalentRegsitrationFields()
            this.showCampusTalentRegistrationFields = true
            this.toggleVisibilityOfCanditateFormAndList()
            if (this.candidateList[index]?.registrationDate) {
                  this.candidateList[index].registrationDate = this.datePipe.transform(this.candidateList[index]?.registrationDate, 'yyyy-MM-dd')
            }
            this.candidateForm.patchValue(this.candidateList[index])
            this.resumeURL = this.candidateForm.get('resume').value
            // Resume.
            this.displayedFileName = "No resume uploaded"
            if (this.candidateList[index].resume) {
                  this.displayedFileName = `<a href=${this.candidateList[index].resume} target="_blank">Resume present</a>`
            }
            this.stateList = []
            this.specializationList = []
            if (this.candidateList[index].country != undefined) {
                  this.getStateListByCountry(this.candidateList[index].country)
            }
            if (this.candidateList[index].degree != undefined) {
                  this.getSpecializationListByDegree(this.candidateList[index].degree)
            }
            this.candidateForm.disable()
      }

      // On clicking update candidate button in view form.
      onUpdateCandidateFormButtonClick(): void {
            this.isViewMode = false
            this.isOperationUpdate = true
            this.candidateForm.enable()
            if (this.candidateForm.get('country').value == null) {
                  this.candidateForm.get('state').disable()
            }
      }

      // Add campus talent registration fields to form in update form.
      addCampusTalentRegsitrationFields(): void {
            this.candidateForm.addControl('registrationDate', this.formBuilder.control(null, [Validators.required]))
            this.candidateForm.addControl('result', this.formBuilder.control(null))
      }

      // On clicking update candidate in candodate form.
      updateCandidate(): void {
            this.spinnerService.loadingMessage = "Updating Candidate"


            let candidate: any = this.candidateForm.value
            this.patchIDFromObjectsForCandidate(candidate)
            this.collegeCampusService.updateCandidate(this.candidateForm.value, this.selectedCampusdriveID).subscribe((response: any) => {
                  this.toggleVisibilityOfCanditateFormAndList()
                  this.getAllCandidates()
                  alert(response)
            }, (error) => {
                  console.error(error)
                  if (typeof error.error == 'object' && error) {
                        alert(this.utilityService.getErrorString(error))
                        return
                  }
                  if (error.error == undefined) {
                        alert('Candidate could not be added, try again')
                        return
                  }
                  alert(error.statusText)
            })
      }

      // Delete candidate.
      deleteCandidate(campusTalentRegistrationID: string, resume: string): void {
            if (confirm("Are you sure you want to delete the candidate?")) {
                  this.spinnerService.loadingMessage = "Deleting Candidate"


                  this.collegeCampusService.deleteCandidate(campusTalentRegistrationID, this.selectedCampusdriveID).subscribe((response: any) => {
                        this.getAllCandidates()
                        this.getCampusDrives()
                        this.fileOperationService.deleteUploadedFile(resume)
                        alert(response)
                  }, (error) => {
                        console.error(error)
                        if (typeof error.error == 'object' && error) {
                              alert(this.utilityService.getErrorString(error))
                              return
                        }
                        if (error.error == undefined) {
                              alert('Candidate could not be added, try again')
                              return
                        }
                        alert(error.statusText)
                  }).add(() => {

                  })
            }
      }

      // Delete resume
      deleteResume(): void {
            this.fileOperationService.deleteUploadedFile().subscribe((data: any) => {
            }, (error) => {
                  console.error(error)
            })
      }

      // ==================================================CANDIDATE SEARCH FUNCTIONS==========================================================================
      // Reset candidate search form and renaviagte page.
      resetSearchAndGetAllForCandidate(): void {
            this.searchCandidateFilterFieldList = []
            this.candidateSearchForm.reset()
            this.searchFormValueCandidate = {}
            this.changePageForCandidate(1)
            this.isSearchedCandidate = false
      }

      // Reset candidate search form.
      resetCandidateSearchForm(): void {
            this.searchCandidateFilterFieldList = []
            this.candidateSearchForm.reset()
      }

      // Search candidates.
      searchCandidates(): void {
            this.searchFormValueCandidate = { ...this.candidateSearchForm?.value }
            for (let field in this.searchFormValueCandidate) {
                  if (this.searchFormValueCandidate[field] === null || this.searchFormValueCandidate[field] === "") {
                        delete this.searchFormValueCandidate[field]
                  } else {
                        this.isSearchedCandidate = true
                  }
            }
            this.searchCandidateFilterFieldList = []
            for (var property in this.searchFormValueCandidate) {
                  let text: string = property
                  let result: string = text.replace(/([A-Z])/g, " $1");
                  let finalResult: string = result.charAt(0).toUpperCase() + result.slice(1);
                  let valueArray: any[] = []
                  if (Array.isArray(this.searchFormValueCandidate[property])) {
                        valueArray = this.searchFormValueCandidate[property]
                  }
                  else {
                        valueArray.push(this.searchFormValueCandidate[property])
                  }
                  this.searchCandidateFilterFieldList.push(
                        {
                              propertyName: property,
                              propertyNameText: finalResult,
                              valueList: valueArray
                        })
            }
            if (this.searchCandidateFilterFieldList.length == 0) {
                  this.resetSearchAndGetAllForCandidate()
            }
            if (!this.isSearchedCandidate) {
                  return
            }
            this.spinnerService.loadingMessage = "Searching Candidates"
            this.changePageForCandidate(1)
      }

      // Delete search criteria from candidate search form by search name.
      deleteSearchCandidateCriteria(searchName: string): void {
            this.candidateSearchForm.get(searchName).setValue(null)
            this.searchCandidates()
      }

      //*********************************************OTHER FUNCTIONS FOR CANDIDATE************************************************************
      // Make all address fields into one comma separated string.
      formatAddress(): void {
            let addressArray: string[] = []
            for (let i = 0; i < this.candidateList.length; i++) {
                  addressArray = []
                  if (this.candidateList[i].address) {
                        addressArray.push(this.candidateList[i].address)
                  }
                  if (this.candidateList[i].pinCode) {
                        addressArray.push(this.candidateList[i].pinCode)
                  }
                  if (this.candidateList[i].city) {
                        addressArray.push(this.candidateList[i].city)
                  }
                  if (this.candidateList[i].state) {
                        addressArray.push(this.candidateList[i].state.name)
                  }
                  if (this.candidateList[i].country) {
                        addressArray.push(this.candidateList[i].country.name)
                  }
                  if (addressArray.length > 1) {
                        this.candidateList[i].fullAddress = addressArray.join(", ")
                  }
                  if (addressArray.length == 1) {
                        this.candidateList[i].fullAddress = addressArray[0]
                  }
            }
      }

      // On clicking sumbit button in candidate form.
      onCandidateFormSubmit(): void {
            if (this.candidateForm.invalid) {
                  this.candidateForm.markAllAsTouched()
                  return
            }
            if (this.isOperationUpdate) {
                  this.updateCandidate()
                  return
            }
            this.addCandidate()
      }

      // Toggle visibility of candidate list and candidate form.
      toggleVisibilityOfCanditateFormAndList(): void {
            if (this.showRegisterCandidateForm) {
                  this.showRegisterCandidateForm = false
                  this.showCandidatesList = true
                  return
            }
            this.showRegisterCandidateForm = true
            this.showCandidatesList = false
      }

      // Change page for pagination.
      changePageForCandidate($event): void {
            this.offsetCandidate = $event - 1
            this.currentPageCandidate = $event
            this.getAllCandidates()
      }

      // On closing the candidate form.
      onCloseCandidateForm(): void {
            if (this.isFileUploading) {
                  alert("Please wait till file is being uploaded")
                  return
            }
            if (this.isResumeUploadedToServer) {
                  if (!confirm("Uploaded resume will be deleted.\nAre you sure you want to close?")) {
                        return
                  }
                  this.deleteResume()
            }
            this.isResumeUploadedToServer = false
            this.displayedFileName = "Select file"
            this.docStatus = ""
            this.toggleVisibilityOfCanditateFormAndList()
            this.isViewMode = false
            this.isOperationUpdate = false
      }

      // Initialize all variables related to candidate.
      initializeCandidateVariables(): void {
            this.candidateList = []
            this.spinnerService.loadingMessage = "Getting All Candidates"
            this.showRegisterCandidateForm = false
            this.showCandidatesList = true
            this.limitCandidate = 5
            this.offsetCandidate = 0
            this.currentPageCandidate = 0
            this.selectedCandidateList = []
            this.multipleSelect = false
            this.isViewMode = false
            this.isOperationUpdate = false
            this.isSearchedCandidate = false
            this.searchFormValueCandidate = {}
      }

      // On uplaoding resume.
      onResourceSelect(event: any): void {
            this.docStatus = ""
            let files = event.target.files
            if (files && files.length) {
                  let file = files[0]

                  // Upload resume if it is present.]
                  this.isFileUploading = true
                  this.fileOperationService.uploadResume(file).subscribe((data: any) => {
                        this.candidateForm.markAsDirty()
                        this.candidateForm.patchValue({
                              resume: data
                        })
                        this.displayedFileName = file.name
                        this.isFileUploading = false
                        this.isResumeUploadedToServer = true
                        this.docStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
                  }, (error) => {
                        this.isFileUploading = false
                        this.docStatus = `<p><span>&#10060;</span> ${error}</p>`
                  })
            }
      }

      // Will be called on [addTag]= "addCollegeToList" args will be passed automatically.
      addCollegeToList(option: any): Promise<any> {
            return new Promise((resolve) => {
                  resolve(option)
            }
            )
      }

      // Extract ID from objects in candidate form.
      patchIDFromObjectsForCandidate(candidate: any): void {
            if (this.candidateForm.get('degree').value) {
                  candidate.degreeID = this.candidateForm.get('degree').value.id
                  delete candidate['degree']
            }
            if (this.candidateForm.get('specialization').value) {
                  candidate.specializationID = this.candidateForm.get('specialization').value.id
                  delete candidate['specialization']
            }
      }

      // Set total list on current page.
      setPaginationStringForCandidate(): void {
            this.paginationStartCandidate = this.limitCandidate * this.offsetCandidate + 1
            this.paginationEndCandidate = +this.limitCandidate + this.limitCandidate * this.offsetCandidate
            if (this.totalCandidates < this.paginationEndCandidate) {
                  this.paginationEndCandidate = this.totalCandidates
            }
      }

      //*********************************************OTHER FUNCTIONS************************************************************
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

      //Compare for select option field.
      compareFn(optionOne: any, optionTwo: any): boolean {
            if (optionOne == null && optionTwo == null) {
                  return true
            }
            if (optionTwo != undefined && optionOne != undefined) {
                  return optionOne.id === optionTwo.id
            }
            return false
      }

      //*********************************************UPDATE MULTIPLE CANDIDATE FUNCTIONS************************************************************
      // On clicking update multiple candidates.
      OnClickingUpdateMultipleCandidates(): void {
            this.createUpdateMultipleCandidateForm()
            if (this.selectedCandidateList.length == 0) {
                  alert("Please select candidates to be updated")
                  return
            }
            this.openModal(this.updateMultipleCandidateModal, 'md')
      }

      // Toggle visibility of multiple select checkbox.
      toggleMultipleSelect(): void {
            if (this.multipleSelect) {
                  this.multipleSelect = false
                  this.setSelectAllCandidates(this.multipleSelect)
                  return
            }
            this.multipleSelect = true
      }

      // Set isChecked field of all selected candidates.
      setSelectAllCandidates(isSelectedAll: boolean): void {
            for (let i = 0; i < this.candidateList.length; i++) {
                  this.addCandidateToList(isSelectedAll, this.candidateList[i])
            }
      }

      // Check if all candidates in selected candidates are added in multiple select or not.
      checkCandidatesAdded(): boolean {
            let count: number = 0

            for (let i = 0; i < this.candidateList.length; i++) {
                  if (this.selectedCandidateList.includes(this.candidateList[i].id))
                        count = count + 1
            }
            return (count == this.candidateList.length)
      }

      // Check if candidate is added in multiple select or not.
      checkCandidateAdded(campusTalentRegistrationID: string): boolean {
            return this.selectedCandidateList.includes(campusTalentRegistrationID)
      }

      // Takes a list called selectedCandidatesList and adds all the checked candidates to list, also does not contain duplicate values.
      addCandidateToList(isChecked: boolean, candidate: any): void {
            if (isChecked) {
                  if (!this.selectedCandidateList.includes(candidate.campusTalentRegistrationID)) {
                        this.selectedCandidateList.push(candidate.campusTalentRegistrationID)
                  }
                  return
            }
            if (this.selectedCandidateList.includes(candidate.campusTalentRegistrationID)) {
                  let index = this.selectedCandidateList.indexOf(candidate.campusTalentRegistrationID)
                  this.selectedCandidateList.splice(index, 1)
            }
      }

      // Update fields of multiple candidates.
      updateMultipleCandidateFunction(): void {
            this.updateMultipleCandidateFormValue = { ...this.updateMultipleCandidateForm.value }
            // Check if all fields are null.
            let flag: boolean = true
            for (let field in this.updateMultipleCandidateFormValue) {
                  if (this.updateMultipleCandidateFormValue[field] == null) {
                        delete this.updateMultipleCandidateFormValue[field]
                  } else {
                        flag = false
                  }
            }
            // No API call on empty search.
            if (flag) {
                  alert("Please select any one value to be updated")
                  return
            }
            this.updateMultipleCandidateFormValue.campusTalentRegistrationIDs = this.selectedCandidateList
            this.spinnerService.loadingMessage = "Updating Candidates"


            this.updateMultipleCandidateFormValue.campusDriveID = this.selectedCampusdriveID
            this.collegeCampusService.updateMultipleCandidate(this.updateMultipleCandidateFormValue).subscribe((response: any) => {
                  this.getAllCandidates()
                  this.getCampusDrives()
                  alert(response)
                  this.modalRef.close('success')
                  this.selectedCandidateList = []
            }, (error) => {
                  console.error(error)
                  if (typeof error.error == 'object' && error) {
                        alert(this.utilityService.getErrorString(error))
                        return
                  }
                  if (error.error == undefined) {
                        alert('Candidates could not be updated')
                  }
                  alert(error.statusText)
            })
      }

      //*********************************************GET FUNCTIONS************************************************************
      // Get all lists.
      getAllComponents(): void {
            this.getCountryList()
            this.getQualificationList()
            this.getSalesPersonList()
            this.getAcademicYear()
            this.getCollegeBranchList()
            this.getCompanyRequirementList()
            this.getFacultyList()
            this.getCandidateResultList()
            this.getDeveloperList()
            this.searchOrGetCampusDrives()
      }

      // Get salesperson list.
      getSalesPersonList(): void {
            this.generalService.getSalesPersonList().subscribe(
                  response => {
                        this.salesPersonList = response.body
                  }), err => {
                        console.error(err)
                  }
      }

      // Get developer list.
      getDeveloperList(): void {
            this.generalService.getEmployeeList().subscribe(
                  response => {
                        this.developerList = response
                  }), err => {
                        console.error(err)
                  }
      }

      // Get academic year list.
      getAcademicYear(): void {
            this.generalService.getGeneralTypeByType("academic_year").subscribe((response: any[]) => {
                  this.academicYearList = response
            }, (err) => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get candidate result list.
      getCandidateResultList(): void {
            this.generalService.getGeneralTypeByType("candidate_result").subscribe((respond: any[]) => {
                  this.candidateResultList = respond
            }, (err) => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get all campus drives by limit and offset.
      getCampusDrives(): void {
            this.spinnerService.loadingMessage = "Getting Campus Drives"


            this.searchFormValueCampusDrive.roleName = this.roleName
            this.searchFormValueCampusDrive.loginID = this.localService.getJsonValue("loginID")
            this.collegeCampusService.getCampusDrives(this.limitCampusDrive, this.offsetCampusDrive, this.searchFormValueCampusDrive).subscribe(response => {
                  this.campusDriveList = response.body
                  this.totalCampusDrives = parseInt(response.headers.get("X-Total-Count"))
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            }).add(() => {

                  this.setPaginationStringForCampusDrive()
            })
      }

      // Get state list by country.
      getStateListByCountry(country: any): void {
            if (country == null) {
                  this.stateList = []
                  this.candidateForm.get('state').setValue(null)
                  this.candidateForm.get('state').disable()
                  return
            }
            if (country.name == "U.S.A." || country.name == "India") {
                  this.candidateForm.get('state').enable()
                  this.generalService.getStatesByCountryID(country.id).subscribe((respond: any[]) => {
                        this.stateList = respond
                  }, (err) => {
                        console.error(err)
                  })
            }
            else {
                  this.stateList = []
                  this.candidateForm.get('state').setValue(null)
                  this.candidateForm.get('state').disable()
            }
      }

      // Get specialization list by degree.
      getSpecializationListByDegree(degree: any): void {
            if (degree == null) {
                  this.specializationList = []
                  this.candidateForm.get('specialization').setValue(null)
                  this.candidateForm.get('specialization').disable()
                  return
            }
            this.candidateForm.get('specialization').enable()
            this.generalService.getSpecializationByDegreeID(degree.id).subscribe((response: any) => {
                  this.specializationList = response.body
            }, (err) => {
                  console.error(err)
            })
      }

      // Get specialization list for candidate search.
      getSpecializationListForCandidateSearch(degreeID: string): void {
            if (degreeID == null) {
                  this.specializationList = []
                  return
            }
            for (let i = 0; i < this.qualificationList.length; i++) {
                  if (this.qualificationList[i].id == degreeID) {
                        this.generalService.getSpecializationByDegreeID(degreeID).subscribe(response => {
                              this.specializationList = response.body
                        }, err => {
                              console.error(this.utilityService.getErrorString(err))
                        })
                        return
                  }
            }
      }

      // Get country list.
      getCountryList(): void {
            this.generalService.getCountries().subscribe((respond: any[]) => {
                  this.countryList = respond
            }, (err) => {
                  console.error(err.error)
            })
      }

      // Get qualification list.
      getQualificationList(): void {
            let queryParams: any = {
                  limit: -1,
                  offset: 0,
            }
            this.degreeService.getAllDegrees(queryParams).subscribe((respond: any) => {
                  //respond.push(this.other)
                  this.qualificationList = respond.body
            }, (err) => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get college branch list.
      getCollegeBranchList(): void {
            this.generalService.getCollegeBranchList().subscribe(response => {
                  this.collegeBranchList = response
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get active company requirement list.
      getCompanyRequirementList(): void {
            let queryParams: any = {
                  isActive: "1"
            }
            this.generalService.getRequirementList(queryParams).subscribe(response => {
                  this.companyRequirementList = response.body
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get faculty list.
      getFacultyList(): void {
            this.generalService.getFacultyList().subscribe(response => {
                  this.facultyList = response.body
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get selected campus drive's college list.
      getSelectedCampusDriveCollegeList(): void {
            for (let i = 0; i < this.campusDriveList.length; i++) {
                  if (this.selectedCampusdriveID == this.campusDriveList[i].id) {
                        this.selectedCampusDriveCollegeList = this.campusDriveList[i].collegeBranches
                  }
            }
      }

      /***************************************************************************************************************************** */

      // copyToClipboard
      copyToClipBoard() {
            document.execCommand("www.swabhavtechlabs.com/register?collegeid=20971FH&&campusid=87876NHG8757&&instance=98765jkjgjh")
      }

}
