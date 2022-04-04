import { ChangeDetectorRef, Component, NgZone, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, FormArray, FormControl, Validators } from '@angular/forms';
import { UtilityService } from 'src/app/service/utility/utility.service';

import { IExamination, IMastersAbroad, IScore, ITalent, TalentService, ILifetimeValue, IUniversity, INextAction, ITalentDTO, IExperienceDTO, ICallRecordDTO, IDegree, IPurpose, IOutcome, ICallRecord, ISpecialization, INextActionDTO, IWaitingList, IWaitingListDTO, IExperience, ITalentExcel, IAcademic, ISearchSection, ISearchFilterField, ITalentDownloadExcel } from 'src/app/service/talent/talent.service'
import { GeneralService } from 'src/app/service/general/general.service';
import { ActivatedRoute, Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { NgbModal, NgbModalRef, NgbModalOptions } from '@ng-bootstrap/ng-bootstrap';
import { DatePipe, Location } from '@angular/common';
import { IPermission } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { Role, UrlConstant } from 'src/app/service/constant';
import { BatchService, IMappedSession } from 'src/app/service/batch/batch.service';
import { CompanyRequirementService } from 'src/app/service/company/company-requirement/company-requirement.service';
import { TechnologyService } from 'src/app/service/technology/technology.service';
import { CareerObjectiveService } from 'src/app/service/career-objective/career-objective.service';
import { DegreeService } from 'src/app/service/degree/degree.service';
import { AdminService } from 'src/app/service/admin/admin.service';

@Component({
      selector: 'app-talent-master',
      templateUrl: './talent-master.component.html',
      styleUrls: ['./talent-master.component.css']
})
export class TalentMasterComponent implements OnInit {

      //****************************TALENT*************************************** */
      // Components.
      stateList: any[]
      countryList: any[]
      talentTypeList: any[]
      technologyList: any[]
      degreeList: any[]
      degreeIDSet: Set<string>
      specializationList: any[]
      personalityTypeList: any[]
      designationList: any[]
      salesPersonList: any[]
      academicYearList: any[]
      sourceList: any[]
      collegeBranchList: any[]
      examinationList: any[]
      universityList: IUniversity[]
      facultyList: any[]
      courseList: any[]
      yearOfMSList: number[]

      // Tech and college ng select.
      isCollegeLoading: boolean
      isTechLoading: boolean
      techLimit: number
      techOffset: number
      collegeLimit: number
      collegeOffset: number

      // Flags.
      isViewMode: boolean
      isOperationTalentUpdate: boolean

      // Talent.
      talents: ITalentDTO[]
      selectedTalents: ITalentDTO[]
      selectedTalent: any
      selectTalentTemp: any
      selectedTalentsList: any[]
      talentForm: FormGroup
      selectedTalentID: string

      // Experience.
      showExperiencesInForm: boolean
      showExperienceColumns: boolean

      // Academic.
      showAcademicsInForm: boolean

      // Pagination.
      limit: number
      offset: number
      totalTalents: number
      currentPage: number
      paginationStart: number
      paginationEnd: number

      // Modal.
      modalRef: any
      @ViewChild('talentFormModal') talentFormModal: any
      @ViewChild('deleteTalentModal') deleteTalentModal: any

      // Spinner.



      // Search.
      talentSearchForm: FormGroup
      isSearched: boolean
      selectedSectionName: string
      searchFilterFieldList: ISearchFilterField[]
      searchSectionList: ISearchSection[]

      // Talent retaled numbers.
      indexOfCurrentWorkingExp: number
      currentYear: number

      // Masters abroad related type.
      greExam: IExamination
      gmatExam: IExamination
      toeflExam: IExamination
      ieltsExam: IExamination
      // *************initialize and declare better -n*********************
      universityMap: Map<string, IUniversity[]> = new Map()
      areUniversitiesLoading: boolean
      showMastersAbroadInForm: boolean
      showMBADegreeRequiredError: boolean

      // Permission.
      permission: IPermission
      roleName: string
      showForAdmin: boolean
      showForFaculty: boolean
      showForSalesPerson: boolean

      // Resume.
      isResumeUploadedToServer: boolean
      isResumeUploading: boolean
      resumeDocStatus: string
      resumeDisplayedFileName: string

      // Image.
      isImageUploadedToServer: boolean
      isImageUploading: boolean
      imageDocStatus: string
      imageDisplayedFileName: string

      // Constants.
      readonly MAX_VALUE = 999999999999

      // Excel.
      excelUploadedTalents: ITalentExcel[]
      reformedTalents: any[]
      isExcelUploaded: boolean
      addMultipleErrorList: string[]
      excelErrorList: string[]
      excelTalentsAdded: string
      talentsExcelTotalCount: number
      talentsExcelAddedCount: number
      talentsExcelProcessedCount: number
      showExcelProgress: boolean
      disableAddTalentsButton: boolean
      excelDownloadTalents: ITalentDownloadExcel[]
      excelDownloadOffsetTotal: number
      excelDownloadOffsetCount: number
      public readonly EXCE_DOWNLOAD_LIMIT_5: number = 5
      public readonly EXCE_DOWNLOAD_LIMIT_50: number = 50

      // Constant.
      readonly TALENT_EXCEL_DEMO_LINK = this.urlConstant.TALENT_BASIC_DEMO

      // arrayBuffer: any
      // file: File
      // excelData: any
      // talentData: any[]
      // errorMsg: string
      // talentError: boolean
      // upload: boolean
      // excelFileName: string

      //****************************CALL RECORDS*************************************** */
      // Components.
      purposeList: any[]
      outcomeList: any[]

      // Flags.
      showCallRecordForm: boolean
      showPurposeJobFields: boolean
      showOutcomeTargetDateField: boolean
      isOperationCallRecordUpdate: boolean

      // Talent call record.
      talentCallRecordList: ICallRecordDTO[]
      callRecordForm: FormGroup

      // Modal.
      @ViewChild('callRecordFormModal') callRecordFormModal: any

      //****************************LIFETIME VALUE*************************************** */
      // Lifetime value.
      selectedLifetimeValue: ILifetimeValue
      isOperationLifetimeValueUpdate: boolean
      lifetimeValueForm: FormGroup

      // Flags.
      showLifetimeValue: boolean
      showLifetimeForm: boolean

      // Lifetime value related numbers.
      totalLifetimeValue: number

      // Modal.
      @ViewChild('lifetimeValueFormModal') lifetimeValueFormModal: any

      //****************************NEXT ACTION*************************************** */
      // Components.
      companyList: any[]
      nextActionTypeList: any[]

      // Next action.      
      nextActionForm: FormGroup
      nextActionList: INextActionDTO[]

      // Flags.
      isOperationNextActionUpdate: boolean
      showNextActionForm: boolean
      showNextActionTechnologies: boolean
      showNextActionTargetDate: boolean
      showNextActionStipend: boolean

      // Next action related numbers.
      totalNextActions: number

      //Modal.
      @ViewChild('nextActionFormModal') nextActionFormModal: any

      //****************************CAREER PLAN*************************************** */
      // Components.      
      careerObjectiveList: any[]

      // Career Plan.
      careerPlanList: any[]
      careerPlanForm: FormGroup
      careerPlanMap: Map<string, any[]> = new Map()
      selectedCareerPlan: any

      // Flags.
      showAddCareerPlanForm: boolean
      showUpdateCareerPlanForm: boolean

      // Career plan realted numbers.
      currentRatingUpdate: number

      // Modal
      @ViewChild('careerPlanFromModal') careerPlanFromModal: any

      //****************************BATCH DETAILS*************************************** */
      // Flags.
      showAdditionalBatchDetails: boolean
      isNavigatedFromBatch: boolean
      batchesFound: boolean
      showBatchSession: boolean

      // Batch.
      batchID: string
      batchSessions: IMappedSession[]
      batchesOfOneTalentList: any[]
      batchList: any[]
      supervisorCount: number
      isHeadFcaulty: boolean
      isHeadSalesPerson: boolean
      isViewAllTalents: boolean

      // Modal
      @ViewChild('showBatchesModal') showBatchesModal: any

      //****************************COMPANY REQUIREMENT*************************************** */
      // Requirement.
      requirementID: string

      // Flags.
      isNavigated: boolean
      // Make this false to show Allocate to batch button.
      showAllocateCompany: boolean
      requirementSearched: boolean
      showRequirementTalents: boolean

      //****************************ALLOCATION*************************************** */
      // Flags.
      showAllocateSalespersonToOneTalent: boolean
      showAllocateSalespersonToTalents: boolean
      showAllocateOneTalentToBatch: boolean
      showAllocateTalentsToBatch: boolean
      multipleSelect: boolean

      // Modal.
      @ViewChild('allocateToBatchModal') allocateToBatchModal: any
      @ViewChild('allocateSalespersonModal') allocateSalespersonModal: any

      //****************************WAITING LIST*************************************** */
      // Components.
      companyBranchList: any[]
      activeRequirementByCompanyList: any[]
      activeBatchByCourseList: any[]
      requiremenList: any[]
      allBatchList: any[]
      allRequirementList: any[]

      // Flags.
      showWaitingListForm: boolean
      isOperationWaitingListUpdate: boolean
      showCompanyBranchList: boolean
      showRequirementList: boolean
      showCourseList: boolean
      showBatchList: boolean
      isNavigatedFromWaitingList: boolean

      // Waiting list redirection variables.
      waitingListBatchID: string
      waitingListRequirementID: string
      waitingListCompanyBranchID: string
      waitingListCourseID: string
      waitingListTechnologyID: string

      // Waiting list.
      waitingList: IWaitingListDTO[]
      waitingListForm: FormGroup

      // Modal.
      @ViewChild('waitingListFormModal') waitingListFormModal: any

      //****************************CAMPUS DRIVE*************************************** */
      // Flags.
      isNavigatedFromCampusDrive: boolean

      // Campus drive redirection variables.
      campusDriveID: string
      hasAppeared: string

      //****************************SEMINAR*************************************** */
      // Flags.
      isNavigatedFromSeminar: boolean

      // Seminar redirection variables.
      seminarID: string
      hasVisited: string

      //****************************PROFESSIONAL SUMMARY REPORT************************************** */
      // Flags.
      isNavigatedFromProSummaryReport: boolean

      // Professional summary report redirection variables.
      companyName: string
      category: string
      isCompany: string

      //****************************PROFESSIONAL SUMMARY REPORT BY TECHNOLOGY COUNT************************************** */
      // Flags.
      isNavigatedFromProSummaryReportTechCount: boolean

      // Professional summary report redirection variables.
      companyNameTechCount: string
      categoryTechCount: string
      isCompanyTechCount: string
      technologyIDTechCount: string

      //****************************FRESHER SUMMARY REPORT************************************** */

      // Fresher summary report redirection variables.
      isNavigatedFromFresherSummaryReport: boolean
      academicYear: string
      talentType: string
      isExperienced: string
      isLookingForJob: string
      fresherTechnology: string

      //****************************PACKAGE SUMMARY REPORT************************************** */

      // Package summary report redirection variables.
      isNavigatedFromPackageSummaryReport: boolean
      packageType: string
      packageTechnology: string
      packageExperience: string

      constructor(
            public utilityService: UtilityService,
            public careerObjectiveService: CareerObjectiveService,
            private formBuilder: FormBuilder,
            private talentService: TalentService,
            private generalService: GeneralService,
            private degreeService: DegreeService,
            private techService: TechnologyService,
            private activatedRoute: ActivatedRoute,
            private companyRequirementService: CompanyRequirementService,
            private fileOps: FileOperationService,
            private router: Router,
            private spinnerService: SpinnerService,
            private modalService: NgbModal,
            private localService: LocalService,
            private batchService: BatchService,
            private urlConstant: UrlConstant,
            private role: Role,
            private datePipe: DatePipe,
            private _location: Location,
            private zone: NgZone,
            private changeDetectorRef: ChangeDetectorRef,
            private adminService: AdminService
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
            //****************************TALENT*************************************** */
            // Components.
            this.stateList = []
            this.countryList = []
            this.talentTypeList = []
            this.technologyList = []
            this.degreeList = []
            this.specializationList = []
            this.personalityTypeList = []
            this.designationList = []
            this.salesPersonList = []
            this.academicYearList = []
            this.sourceList = []
            this.collegeBranchList = []
            this.examinationList = []
            this.universityList = []
            this.facultyList = []
            this.courseList = []
            this.yearOfMSList = []
            this.degreeIDSet = new Set()

            // Flags.
            this.isViewMode = false
            this.isOperationTalentUpdate = false

            // tech & college component
            this.isCollegeLoading = false
            this.isTechLoading = false
            this.techOffset = 0
            this.collegeOffset = 0
            this.techLimit = 10
            this.collegeLimit = 10

            // Talent.
            this.selectTalentTemp = {}
            this.selectedTalents = []
            this.selectedTalent = {}
            this.selectedTalentsList = []

            // Experience
            this.showExperiencesInForm = false
            this.showExperienceColumns = false

            // Academic.
            this.showAcademicsInForm = false

            // Paginate.
            this.limit = 5
            this.offset = 0
            this.currentPage = 0
            this.paginationStart = 0
            this.paginationEnd = 0

            // Spinner.
            this.spinnerService.loadingMessage = "Getting Talents"


            // Search.
            this.isSearched = false
            this.searchFilterFieldList = []
            this.searchSectionList = [
                  {
                        name: "Personal",
                        isSelected: true
                  },
                  {
                        name: "Academics",
                        isSelected: false
                  },
                  {
                        name: "Experiences",
                        isSelected: false
                  },
                  {
                        name: "Calling",
                        isSelected: false
                  },
                  {
                        name: "Courses",
                        isSelected: false
                  },
                  {
                        name: "Waiting List",
                        isSelected: false
                  },
                  {
                        name: "Others",
                        isSelected: false
                  }
            ]
            this.selectedSectionName = "Personal"

            // Talent retaled numbers.
            this.indexOfCurrentWorkingExp = -1
            this.currentYear = new Date().getFullYear()

            // Masters abroad.
            this.areUniversitiesLoading = false
            this.showMastersAbroadInForm = false
            this.showMBADegreeRequiredError = false

            // Permision.
            this.showForAdmin = false
            this.showForFaculty = true
            this.showForSalesPerson = true
            this.roleName = this.localService.getJsonValue("roleName")
            // Get role name for menu for calling their specific apis.
            // If admin is logged in then show its features.
            if (this.roleName == this.role.ADMIN) {
                  this.showForAdmin = true
                  this.permission = this.utilityService.getPermission(this.urlConstant.TALENT_MASTER)
            }
            // Hide features for salesperson.
            if (this.roleName == this.role.SALES_PERSON) {
                  this.showForSalesPerson = false
                  this.permission = this.utilityService.getPermission(this.urlConstant.TALENT_MASTER)
            }
            // Hide features for faculty.
            if (this.roleName == this.role.FACULTY) {
                  this.showForFaculty = false
                  this.permission = this.utilityService.getPermission(this.urlConstant.MY_TALENT)
            }

            // Resume.
            this.resumeDocStatus = ""
            this.resumeDisplayedFileName = "Select file"
            this.isResumeUploadedToServer = false
            this.isResumeUploading = false

            // Image.
            this.imageDocStatus = ""
            this.imageDisplayedFileName = "Select file"
            this.isImageUploadedToServer = false
            this.isImageUploading = false

            // Excel.
            this.isExcelUploaded = false
            this.addMultipleErrorList = []
            this.excelErrorList = []
            this.talentsExcelAddedCount = 0
            this.talentsExcelTotalCount = 0
            this.talentsExcelProcessedCount = 0
            this.showExcelProgress = false
            this.disableAddTalentsButton = false
            this.excelDownloadTalents = []
            this.excelDownloadOffsetCount = 0
            this.excelDownloadOffsetTotal = 0
            // this.talentData = []
            // this.errorMsg = ""
            // this.talentError = false
            // this.upload = false
            // this.excelFileName = ""

            //****************************CALL RECORDS*************************************** */
            // Components.
            this.purposeList = []
            this.outcomeList = []

            // Flags.
            this.showCallRecordForm = false
            this.showPurposeJobFields = false
            this.showOutcomeTargetDateField = false
            this.isOperationCallRecordUpdate = false

            // Talent call record.
            this.talentCallRecordList = [] as ICallRecordDTO[]

            //****************************LIFETIME VALUE*************************************** */
            // Flags.
            this.showLifetimeValue = false
            this.showLifetimeForm = false
            this.isOperationLifetimeValueUpdate = false

            // Lifetime value related numbers.
            this.totalLifetimeValue = 0

            //****************************NEXT ACTION*************************************** */
            // Components.
            this.companyList = []
            this.nextActionTypeList = [] as INextActionDTO[]

            // Next action.      
            this.nextActionList = []

            // Flags.
            this.isOperationNextActionUpdate = false
            this.showNextActionForm = false
            this.showNextActionTechnologies = false
            this.showNextActionTargetDate = false
            this.showNextActionStipend = false

            // Next action realted numbers.
            this.totalNextActions = 0

            //****************************CAREER PLAN*************************************** */
            // Components.      
            this.careerObjectiveList = []

            // Career Plan.
            this.careerPlanList = []

            // Flags.
            this.showAddCareerPlanForm = false
            this.showUpdateCareerPlanForm = false

            // Career plan realted numbers.
            this.currentRatingUpdate = 0

            //****************************BATCH DETAILS*************************************** */
            // Flags.
            this.showAdditionalBatchDetails = true
            this.isNavigatedFromBatch = false
            this.batchesFound = true
            this.showBatchSession = false
            this.showAllocateOneTalentToBatch = false
            this.showAllocateTalentsToBatch = false
            this.isViewAllTalents = false

            // Batch.
            this.batchSessions = []
            this.batchesOfOneTalentList = []
            this.batchList = []
            this.supervisorCount = 0
            this.isHeadFcaulty = false
            this.isHeadSalesPerson = false

            //****************************COMPANY REQUIREMRNT*************************************** */
            // Flags.
            this.isNavigated = false
            this.showAllocateCompany = false
            this.requirementSearched = false
            this.showRequirementTalents = false

            //****************************ALLOCATION*************************************** */
            // Flags.
            this.showAllocateSalespersonToOneTalent = false
            this.showAllocateSalespersonToTalents = false

            //****************************WAITING LIST*************************************** */
            // Components.
            this.companyBranchList = []
            this.activeRequirementByCompanyList = []
            this.activeBatchByCourseList = []
            this.allBatchList = []
            this.allRequirementList = []

            // Flags.
            this.showWaitingListForm = false
            this.isOperationWaitingListUpdate = false
            this.showCompanyBranchList = false
            this.showRequirementList = false
            this.showCourseList = false
            this.showBatchList = false
            this.isNavigatedFromWaitingList = false

            // Waiting list.
            this.waitingList = [] as IWaitingListDTO[]

            // Waiting list redirection variables.
            this.waitingListBatchID = null
            this.waitingListRequirementID = null
            this.waitingListCompanyBranchID = null
            this.waitingListCourseID = null
            this.waitingListTechnologyID = null

            //****************************CAMPUS DRIVE*************************************** */
            // Flags.
            this.isNavigatedFromCampusDrive = false

            //****************************SEMINAR*************************************** */
            // Flags.
            this.isNavigatedFromSeminar = false

            //****************************PROFESSIONAL SUMMARY REPORT************************************** */
            // Flags.
            this.isNavigatedFromProSummaryReport = false

            // Professional summary report redirection variables.
            this.companyName = null
            this.category = null
            this.isCompany = null

            //****************************PROFESSIONAL SUMMARY REPORT BY TECHNOLOGY COUNT************************************** */
            // Flags.
            this.isNavigatedFromProSummaryReportTechCount = false

            // Professional summary report redirection variables.
            this.companyNameTechCount = null
            this.categoryTechCount = null
            this.isCompanyTechCount = null
            this.technologyIDTechCount = null

            //****************************FRESHER SUMMARY REPORT************************************** */
            // Flags.
            this.isNavigatedFromFresherSummaryReport = false

            // Fresher summary report redirection variables.
            this.academicYear = null
            this.talentType = null
            this.fresherTechnology = null
            this.isLookingForJob = "0"

            //****************************PACKAGE SUMMARY REPORT************************************** */
            // Flags.
            this.isNavigatedFromPackageSummaryReport = false

            // Package summary report redirection variables.
            this.packageType = null
            this.packageTechnology = null
            this.packageExperience = null

            //****************************INITIALIZE FORMS*************************************** */
            this.createSearchTalentForm()
      }

      // Nagivate to previous url.
      backToPreviousPage(): void {
            this._location.back()
      }

      // Get query params from url and get talents accordingly.
      getQueryParams(): void {
            if (this.activatedRoute.snapshot.queryParamMap.keys != []) {
                  this.talents = []

                  // Navigated from requirement search.
                  let requirementSearch = this.activatedRoute.snapshot.queryParamMap.get("requirementSearch")
                  if (requirementSearch) {
                        this.multipleSelect = true
                        this.showAllocateCompany = true
                        this.requirementSearched = true
                        this.isNavigated = true
                        this.requirementID = requirementSearch
                        let param = this.activatedRoute.snapshot.queryParams
                        let searchFormValue = { ...param }
                        for (let field in searchFormValue) {
                              if (!searchFormValue[field]) {
                                    delete searchFormValue[field]
                              } else {
                                    this.isSearched = true
                              }
                        }
                        delete searchFormValue["requirementSearch"]
                        if (searchFormValue["isExperience"]) {
                              searchFormValue["isExperience"] = searchFormValue["isExperience"] == "true" ? true : false
                        }
                        if (searchFormValue["talentType"]) {
                              searchFormValue["talentType"] = Number(searchFormValue["talentType"])
                        }
                        if (searchFormValue["minimumExperience"]) {
                              searchFormValue["minimumExperience"] = Number(searchFormValue["minimumExperience"])
                        }
                        if (searchFormValue["maximumExperience"]) {
                              searchFormValue["maximumExperience"] = Number(searchFormValue["maximumExperience"])
                        }
                        this.talentSearchForm.patchValue(searchFormValue)
                  }

                  // Navigated from company requirement by id.
                  let requirementID = this.activatedRoute.snapshot.queryParamMap.get("requirementID")
                  if (requirementID) {
                        this.requirementID = requirementID
                        this.isNavigated = true
                        this.showRequirementTalents = true
                        this.getAllTalentsInRequirement()
                        return
                  }

                  // Navigated from batch by id.
                  let batchID = this.activatedRoute.snapshot.queryParamMap.get("batchID")
                  if (batchID) {
                        this.isNavigatedFromBatch = true
                        this.isNavigated = true
                        this.batchID = batchID
                        this.getSelectedTalents()
                        return
                  }

                  // Navigated from waiting list.
                  let waitingListCompanyBranchID = this.activatedRoute.snapshot.queryParamMap.get("waitingListCompanyBranchID")
                  let waitingListCourseID = this.activatedRoute.snapshot.queryParamMap.get("waitingListCourseID")
                  let waitingListBatchID = this.activatedRoute.snapshot.queryParamMap.get("waitingListBatchID")
                  let waitingListRequirementID = this.activatedRoute.snapshot.queryParamMap.get("waitingListRequirementID")
                  let waitingListTechnologyID = this.activatedRoute.snapshot.queryParamMap.get("waitingListTechnologyID")
                  if (waitingListCompanyBranchID || waitingListCourseID || waitingListBatchID || waitingListRequirementID ||
                        waitingListTechnologyID) {
                        this.isNavigatedFromWaitingList = true
                        this.isNavigated = true
                        this.waitingListCompanyBranchID = waitingListCompanyBranchID
                        this.waitingListCourseID = waitingListCourseID
                        this.waitingListBatchID = waitingListBatchID
                        this.waitingListRequirementID = waitingListRequirementID
                        this.waitingListTechnologyID = waitingListTechnologyID
                        this.getTalentsForWaitingList()
                        return
                  }

                  // Navigated from campus drive.
                  let campusDriveID = this.activatedRoute.snapshot.queryParamMap.get("campusDriveID")
                  let hasAppeared = this.activatedRoute.snapshot.queryParamMap.get("hasAppeared")
                  if (campusDriveID) {
                        this.isNavigatedFromCampusDrive = true
                        this.isNavigated = true
                        this.campusDriveID = campusDriveID
                        this.hasAppeared = hasAppeared
                        this.getTalentsByCampusdDrive()
                        return
                  }

                  // Navigated from seminar.
                  let seminarID = this.activatedRoute.snapshot.queryParamMap.get("seminarID")
                  let hasVisited = this.activatedRoute.snapshot.queryParamMap.get("hasVisited")
                  if (seminarID) {
                        this.isNavigatedFromSeminar = true
                        this.isNavigated = true
                        this.seminarID = seminarID
                        this.hasVisited = hasVisited
                        this.getTalentsBySeminar()
                        return
                  }

                  // Navigated from professional summary report.
                  let companyName = this.activatedRoute.snapshot.queryParamMap.get("companyName")
                  let category = this.activatedRoute.snapshot.queryParamMap.get("category")
                  let isCompany = this.activatedRoute.snapshot.queryParamMap.get("isCompany")
                  if (companyName && category && isCompany) {
                        this.isNavigatedFromProSummaryReport = true
                        this.isNavigated = true
                        this.companyName = companyName
                        this.category = category
                        this.isCompany = isCompany
                        this.getTalentsForProSummaryReport()
                        return
                  }

                  // Navigated from professional summary report by technology count.
                  let companyNameTechCount = this.activatedRoute.snapshot.queryParamMap.get("companyNameTechCount")
                  let categoryTechCount = this.activatedRoute.snapshot.queryParamMap.get("categoryTechCount")
                  let isCompanyTechCount = this.activatedRoute.snapshot.queryParamMap.get("isCompanyTechCount")
                  let technologyIDTechCount = this.activatedRoute.snapshot.queryParamMap.get("technologyIDTechCount")
                  if (companyNameTechCount && categoryTechCount && isCompanyTechCount && technologyIDTechCount) {
                        this.isNavigatedFromProSummaryReportTechCount = true
                        this.isNavigated = true
                        this.companyNameTechCount = companyNameTechCount
                        this.categoryTechCount = categoryTechCount
                        this.isCompanyTechCount = isCompanyTechCount
                        this.technologyIDTechCount = technologyIDTechCount
                        this.getTalentsForProSummaryReportByTechnologyCount()
                        return
                  }

                  // Navigated from fresher summary report.
                  let isFresherSummary = this.activatedRoute.snapshot.queryParamMap.get("isFresherSummary")
                  if (isFresherSummary) {
                        let academicYear = this.activatedRoute.snapshot.queryParamMap.get("academicYear") // -> 1, 2, 3, 4, 5
                        let talentType = this.activatedRoute.snapshot.queryParamMap.get("talentType") // -> outstanding, excellent, average
                        let isExperienced = this.activatedRoute.snapshot.queryParamMap.get("isExperienced") // -> "0", "1"
                        let isLookingForJob = this.activatedRoute.snapshot.queryParamMap.get("isLookingForJob") // -> "0", "1"
                        let fresherTechnology = this.activatedRoute.snapshot.queryParamMap.get("fresherTechnology") // -> techID
                        this.isNavigatedFromFresherSummaryReport = true
                        this.isNavigated = true
                        if (fresherTechnology) {
                              this.fresherTechnology = fresherTechnology
                        }
                        if (talentType) {
                              this.talentType = talentType
                        }
                        if (academicYear) {
                              this.academicYear = academicYear
                        }
                        if (isExperienced) {
                              this.isExperienced = isExperienced
                        }
                        if (isLookingForJob) {
                              this.isLookingForJob = isLookingForJob
                        }
                        this.getTalentsForFresherSummaryReport()
                        return
                  }

                  // Navigated from package summary report.
                  let isPackageSummary = this.activatedRoute.snapshot.queryParamMap.get("isPackageSummary")
                  if (isPackageSummary) {
                        let packageType = this.activatedRoute.snapshot.queryParamMap.get("packageType")
                        let packageExperience = this.activatedRoute.snapshot.queryParamMap.get("packageExperience")
                        let packageTechnology = this.activatedRoute.snapshot.queryParamMap.get("packageTechnology") // -> techID
                        this.isNavigatedFromPackageSummaryReport = true
                        this.isNavigated = true
                        if (packageType) {
                              this.packageType = packageType
                        }
                        if (packageTechnology) {
                              this.packageTechnology = packageTechnology
                        }
                        if (packageExperience) {
                              this.packageExperience = packageExperience
                        }
                        this.getTalentsForPackageSummaryReport()
                        return
                  }
            }
            this.changePage(1)
      }

      //*********************************************CREATE FORMS************************************************************
      // Create talent form.
      createTalentForm(): void {
            this.talentForm = this.formBuilder.group({
                  id: new FormControl(null),
                  code: new FormControl({ value: null, disabled: true }),
                  firstName: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z]*$/), Validators.maxLength(50)]),
                  lastName: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z]*$/), Validators.maxLength(50)]),
                  email: new FormControl(null, [Validators.required, Validators.pattern(/^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$/),
                  Validators.maxLength(100)]),
                  alternateEmail: new FormControl(null, [Validators.pattern(/^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$/),
                  Validators.maxLength(100)]),
                  contact: new FormControl(null, [Validators.required, Validators.pattern(/^[6789]\d{9}$/)]),
                  alternateContact: new FormControl(null, [Validators.pattern(/^[6789]\d{9}$/)]),
                  address: new FormControl(null, [Validators.pattern(/^[.0-9a-zA-Z\s,-\/]+$/)]),
                  city: new FormControl(null, [Validators.pattern(/^[a-zA-Z]+([a-zA-Z ]?)+$/), Validators.maxLength(50)]),
                  state: new FormControl(null),
                  country: new FormControl(null),
                  pinCode: new FormControl(null, [Validators.pattern(/^[1-9][0-9]{5}$/)]),
                  talentType: new FormControl(null),
                  personalityType: new FormControl(null),
                  talentSource: new FormControl(null),
                  isActive: new FormControl(true),
                  technologies: new FormControl(Array()),
                  academics: this.formBuilder.array([]),
                  isExperience: new FormControl(false),
                  isAcademic: new FormControl(false),
                  experiences: this.formBuilder.array([]),
                  salesPerson: new FormControl(null),
                  resume: new FormControl(null),
                  image: new FormControl(null),
                  facebookUrl: new FormControl(null, [Validators.maxLength(200),
                  Validators.pattern(/^(?:https?:\/\/)?(?:www\.)?facebook\.com\/.(?:(?:\w)*#!\/)?(?:pages\/)?(?:[\w\-]*\/)*([\w\-\.]*)$/)]),
                  instagramUrl: new FormControl(null, [Validators.maxLength(200),
                  Validators.pattern(/(?:(?:http|https):\/\/)?(?:www.)?(?:instagram.com|instagr.am)\/([A-Za-z0-9-_]+)/)]),
                  githubUrl: new FormControl(null, [Validators.maxLength(200),
                  Validators.pattern(/^(http(s?):\/\/)?(www\.)?github\.([a-z])+\/([A-Za-z0-9]{1,})+\/?$/)]),
                  linkedInUrl: new FormControl(null, [Validators.maxLength(200),
                  Validators.pattern(/(https?)?:?(\/\/)?(([w]{3}||\w\w)\.)?linkedin.com(\w+:{0,1}\w*@)?(\S+)(:([0-9])+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?/)]),
                  academicYear: new FormControl(null, [Validators.required]),
                  isMastersAbroad: new FormControl(false),
                  isSwabhavTalent: new FormControl(null, [Validators.required]),
                  experienceInMonths: new FormControl(null),
                  courses: new FormControl(null)
            })
      }

      // Add new matsers abroad to talent form.
      addMastersAbroad(): void {
            this.talentForm.addControl('mastersAbroad', this.formBuilder.group({
                  id: new FormControl(null),
                  degree: new FormControl(null, [Validators.required]),
                  countries: new FormControl(null, [Validators.required]),
                  universities: new FormControl(Array(), [Validators.required]),
                  yearOfMS: new FormControl(null),
                  greID: new FormControl(null),
                  gre: new FormControl(null, [Validators.required, Validators.min(0), Validators.max(this.greExam?.totalMarks)]),
                  gmatID: new FormControl(null),
                  gmat: new FormControl(null, [Validators.min(0), Validators.max(this.gmatExam?.totalMarks)]),
                  toeflID: new FormControl(null),
                  toefl: new FormControl(null, [Validators.min(0), Validators.max(this.toeflExam?.totalMarks)]),
                  ieltsID: new FormControl(null),
                  ielts: new FormControl(null, [Validators.min(0), Validators.pattern(/^(10|\d)(\.5)?$/), Validators.max(this.ieltsExam?.totalMarks)]),
                  talentID: new FormControl(null),
                  enquiryID: new FormControl(null)
            }))
      }

      // Add new experience to talent form.
      addExperience(): void {
            this.talentExperienceControlArray.push(this.formBuilder.group({
                  id: new FormControl(null),
                  company: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z ]*$/), Validators.maxLength(255)]),
                  designation: new FormControl(null, [Validators.required]),
                  technologies: new FormControl(null),
                  isCurrentWorking: new FormControl("false"),
                  fromDate: new FormControl(null, [Validators.required]),
                  toDate: new FormControl(null, [Validators.required]),
                  package: new FormControl(null, [Validators.min(100000), Validators.max(100000000)]),
                  yearsOfExperience: new FormControl(null),
            }))
      }

      // Add new academic to talent form.
      addAcademic(): void {
            this.talentAcademicControlArray.push(this.formBuilder.group({
                  id: new FormControl(null),
                  degree: new FormControl(null, [Validators.required]),
                  specialization: new FormControl(null, [Validators.required]),
                  college: new FormControl(null, [Validators.required, Validators.maxLength(200)]), //, Validators.pattern(/^[a-zA-Z ]*$/)
                  percentage: new FormControl(null, [Validators.min(1), Validators.pattern(/(^100([.]0{1,2})?)$|(^\d{1,2}([.]\d{1,2})?)$/i), Validators.required]),
                  passout: new FormControl(null, [Validators.required, Validators.min(1980), Validators.max(this.currentYear + 3)])
            }))
      }

      // Create talent search form.
      createSearchTalentForm(): void {
            this.talentSearchForm = this.formBuilder.group({
                  email: new FormControl(null),
                  college: new FormControl(null),
                  firstName: new FormControl(null),
                  lastName: new FormControl(null),
                  isExperience: new FormControl(null),
                  technologies: new FormControl(null),
                  passout: new FormControl(null),
                  talentType: new FormControl(null),
                  isActive: new FormControl(true),
                  personalityType: new FormControl(null),
                  experienceTechnologies: new FormControl(null),
                  totalExperience: new FormControl(null),
                  marksCriteria: new FormControl(null),
                  facultyID: new FormControl(null),
                  batchID: new FormControl(null),
                  courseID: new FormControl(null),
                  degrees: new FormControl(null),
                  designations: new FormControl(null),
                  companyName: new FormControl(null),
                  salesPersonIDs: new FormControl(null),
                  academicYears: new FormControl(null),
                  callRecordPurposeID: new FormControl(null),
                  callRecordOutcomeID: new FormControl(null),
                  lifetimeValue: new FormControl(null),
                  minimumExperience: new FormControl(null, [Validators.min(0), Validators.max(30)]),
                  maximumExperience: new FormControl(null, [Validators.min(0), Validators.max(30)]),
                  nextActionTypeID: new FormControl(null),
                  isSwabhavTalent: new FormControl(null),
                  isMastersAbroad: new FormControl(null),
                  searchAllTalents: new FormControl(null),
                  city: new FormControl(null),
                  countryID: new FormControl(null),
                  yearOfMS: new FormControl(null),
                  waitingFor: new FormControl(null),
                  waitingForCompanyBranchID: new FormControl(null),
                  waitingForRequirementID: new FormControl(null),
                  waitingForCourseID: new FormControl(null),
                  waitingForBatchID: new FormControl(null),
                  waitingForIsActive: new FormControl(null),
                  waitingForFromDate: new FormControl(null),
                  waitingForToDate: new FormControl(null),
            })
      }

      // Create new talent call record form.
      createCallRecordForm(): void {
            this.callRecordForm = this.formBuilder.group({
                  id: new FormControl(),
                  dateTime: new FormControl(null, [Validators.required]),
                  purpose: new FormControl(null, [Validators.required]),
                  outcome: new FormControl(null, [Validators.required]),
                  comment: new FormControl(null, [Validators.maxLength(500)]),
                  talentID: new FormControl(null)
            })
      }

      // Create new talent lifetime value form.
      createLifetimeValueForm(): void {
            this.lifetimeValueForm = this.formBuilder.group({
                  id: new FormControl(),
                  upsell: new FormControl(null, [Validators.min(0), Validators.max(this.MAX_VALUE)]),
                  placement: new FormControl(null, [Validators.min(0), Validators.max(this.MAX_VALUE)]),
                  knowledge: new FormControl(null, [Validators.min(0), Validators.max(this.MAX_VALUE)]),
                  teaching: new FormControl(null, [Validators.min(0), Validators.max(this.MAX_VALUE)]),
            })
      }

      // Create next action form.
      createNextActionForm(): void {
            this.nextActionForm = this.formBuilder.group({
                  id: new FormControl(null),
                  talentID: new FormControl(null, [Validators.required]),
                  comment: new FormControl(null),
                  actionType: new FormControl(null, [Validators.required]),
            })
      }

      // Create new career plan form.
      createCareerPlanForm(): void {
            this.careerPlanForm = this.formBuilder.group({
                  careerObjective: new FormControl(null, [Validators.required]),
                  careerPlans: this.formBuilder.array([]),
            })
      }

      // Add career plan to career plans array in career plan form.
      addCareerPlan(): void {
            this.talentCareerPlansControlArray.push(this.formBuilder.group({
                  id: new FormControl(null),
                  careerObjectiveID: new FormControl(null),
                  careerObjectivesCoursesID: new FormControl(null),
                  facultyID: new FormControl(null),
                  talentID: new FormControl(null),
                  currentRating: new FormControl(null, [Validators.required, Validators.min(1), Validators.max(10)]),
                  technicalAspect: new FormControl(null),
            }))
      }

      // Create new waiting list form.
      createWaitingListForm(): void {
            this.waitingListForm = this.formBuilder.group({
                  id: new FormControl(),
                  waitingFor: new FormControl(null, [Validators.required]),
                  talentID: new FormControl(null),
                  enquiryID: new FormControl(null),
                  email: new FormControl(null),
                  isActive: new FormControl(null),
            })
      }

      //*********************************************ADD FOR TALENT FUNCTIONS************************************************************
      // On add new talent button click.
      onAddNewTalentClick(): void {
            this.indexOfCurrentWorkingExp = -1
            this.addMultipleErrorList = []
            this.isViewMode = false
            this.showExperiencesInForm = false
            this.showAcademicsInForm = false
            this.showMastersAbroadInForm = false
            this.showMBADegreeRequiredError = false
            this.isOperationTalentUpdate = false
            this.showExcelProgress = false
            this.createTalentForm()
            this.stateList = []
            this.talentForm.get('state').disable()
            this.universityList = []
            this.openModal(this.talentFormModal, 'xl')
      }

      // Handle creation or updation of masters abroad of talent.
      assignMastersAbroad(addOrUpdate: string): ITalent {
            let scoresArray: IScore[] = []
            let talent: ITalent = this.talentForm.value
            let mastersAbroad: IMastersAbroad
            if (talent.mastersAbroad) {
                  if (this.talentForm.get('mastersAbroad').get('gre').value != null) {
                        let score: IScore = {
                              id: this.talentForm.get('mastersAbroad').get('greID').value,
                              marksObtained: this.talentForm.get('mastersAbroad').get('gre').value,
                              examinationID: this.greExam.id,
                              mastersAbroadID: this.talentForm.get('mastersAbroad').get('id').value
                        }
                        scoresArray.push(score)
                  }
                  if (this.talentForm.get('mastersAbroad').get('gmat').value != null) {
                        let score: IScore = {
                              id: this.talentForm.get('mastersAbroad').get('gmatID').value,
                              marksObtained: this.talentForm.get('mastersAbroad').get('gmat').value,
                              examinationID: this.gmatExam.id,
                              mastersAbroadID: this.talentForm.get('mastersAbroad').get('id').value
                        }
                        scoresArray.push(score)
                  }
                  if (this.talentForm.get('mastersAbroad').get('toefl').value != null) {
                        let score: IScore = {
                              id: this.talentForm.get('mastersAbroad').get('toeflID').value,
                              marksObtained: this.talentForm.get('mastersAbroad').get('toefl').value,
                              examinationID: this.toeflExam.id,
                              mastersAbroadID: this.talentForm.get('mastersAbroad').get('id').value
                        }
                        scoresArray.push(score)
                  }
                  if (this.talentForm.get('mastersAbroad').get('ielts').value != null) {
                        let score: IScore = {
                              id: this.talentForm.get('mastersAbroad').get('ieltsID').value,
                              marksObtained: this.talentForm.get('mastersAbroad').get('ielts').value,
                              examinationID: this.ieltsExam.id,
                              mastersAbroadID: this.talentForm.get('mastersAbroad').get('id').value
                        }
                        scoresArray.push(score)
                  }
                  mastersAbroad = {
                        scores: scoresArray,
                        degreeID: this.talentForm.get('mastersAbroad').get('degree').value.id,
                        countries: talent.mastersAbroad.countries,
                        universities: talent.mastersAbroad.universities,
                        yearOfMS: talent.mastersAbroad.yearOfMS,
                        talentID: talent.id
                  }
                  if (addOrUpdate == "update") {
                        if (this.selectTalentTemp.mastersAbroad) {
                              mastersAbroad.id = this.selectTalentTemp.mastersAbroad.id
                        }
                  }
                  talent.mastersAbroad = mastersAbroad
                  return talent
            }
            talent.mastersAbroad = null
            return talent
      }

      // Add new talent.
      addTalent(): void {
            let talent: ITalent = this.assignMastersAbroad("add")
            this.patchIDFromObjectsForTalent(talent)
            this.spinnerService.loadingMessage = "Adding Talent"


            this.talentService.addTalent(talent).subscribe((response: any) => {
                  this.modalRef.close('success')
                  this.getTalents()
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

      //*********************************************UPDATE AND VIEW FOR TALENT FUNCTIONS************************************************************
      // On clicking view talent button.
      onViewTalentClick(id: string): void {
            this.indexOfCurrentWorkingExp = -1
            this.isViewMode = true
            this.showExperiencesInForm = false
            this.showAcademicsInForm = false
            this.showMastersAbroadInForm = false
            this.showMBADegreeRequiredError = false
            this.selectedTalent = null
            this.createTalentForm()
            this.assignSelectedTalent(id)

            // Resume.
            this.resumeDisplayedFileName = "No resume uploaded"
            if (this.selectedTalent.resume) {
                  this.resumeDisplayedFileName = `<a href=${this.selectedTalent.resume} target="_blank">Resume present</a>`
            }

            // Image.
            this.imageDisplayedFileName = "No image uploaded"
            if (this.selectedTalent.image) {
                  this.imageDisplayedFileName = `<a href=${this.selectedTalent.image} target="_blank">Image present</a>`
            }

            // Expereinces.
            this.talentForm.setControl("experiences", this.formBuilder.array([]))
            if (this.selectedTalent.experiences && this.selectedTalent.experiences.length > 0) {
                  for (let i = 0; i < this.selectedTalent.experiences.length; i++) {
                        if (this.selectedTalent.experiences[i].toDate == null) {
                              this.indexOfCurrentWorkingExp = i
                        }
                        this.calculateYearsOfExperienceForViewMode(this.selectedTalent.experiences[i])
                        let fromDate = this.selectedTalent.experiences[i]?.fromDate
                        if (fromDate) {
                              this.selectedTalent.experiences[i].fromDate = this.datePipe.transform(fromDate, 'yyyy-MM')
                        }
                        let toDate = this.selectedTalent.experiences[i]?.toDate
                        if (toDate) {
                              this.selectedTalent.experiences[i].toDate = this.datePipe.transform(toDate, 'yyyy-MM')
                        }
                        this.addExperience()
                  }
                  this.showExperiencesInForm = true
            }

            if (this.indexOfCurrentWorkingExp != -1) {
                  this.talentExperienceControlArray.at(this.indexOfCurrentWorkingExp).get('isCurrentWorking').setValue('true')
                  let control = this.talentExperienceControlArray.controls[this.indexOfCurrentWorkingExp]
                  if (control instanceof FormGroup) {
                        control.removeControl('toDate')
                  }
            }

            // Academics.
            this.talentForm.setControl("academics", this.formBuilder.array([]))
            if (this.selectedTalent.academics && this.selectedTalent.academics.length > 0) {
                  for (let i = 0; i < this.selectedTalent.academics.length; i++) {
                        this.addAcademic()
                        this.getSpecializationListByDegreeID(this.selectedTalent.academics[i].degree.id)
                  }
                  this.showAcademicsInForm = true
                  this.talentForm.get("isAcademic").setValue(true)
            }

            // Masters Aborad.
            let greScore: number = null
            let greID: string = null
            let gmatScore: any = null
            let gmatID: string = null
            let toeflScore: any = null
            let toeflID: string = null
            let ieltsScore: any = null
            let ieltsID: string = null

            this.selectTalentTemp = {}

            for (var k in this.selectedTalent) {
                  this.selectTalentTemp[k] = this.selectedTalent[k]
            }

            if (this.selectedTalent.mastersAbroad) {
                  this.addMastersAbroad()
                  this.showMastersAbroadInForm = true
                  this.talentForm.get("isMastersAbroad").setValue(true)
                  for (let i = 0; i < this.selectedTalent.mastersAbroad.scores.length; i++) {
                        if (this.selectedTalent.mastersAbroad.scores[i].examination.name == "GRE") {
                              greScore = this.selectedTalent.mastersAbroad.scores[i].marksObtained
                              greID = this.selectedTalent.mastersAbroad.scores[i].id
                        }
                        if (this.selectedTalent.mastersAbroad.scores[i].examination.name == "GMAT") {
                              gmatScore = this.selectedTalent.mastersAbroad.scores[i].marksObtained
                              gmatID = this.selectedTalent.mastersAbroad.scores[i].id
                        }
                        if (this.selectedTalent.mastersAbroad.scores[i].examination.name == "TOEFL") {
                              toeflScore = this.selectedTalent.mastersAbroad.scores[i].marksObtained
                              toeflID = this.selectedTalent.mastersAbroad.scores[i].id
                        }
                        if (this.selectedTalent.mastersAbroad.scores[i].examination.name == "IELTS") {
                              ieltsScore = this.selectedTalent.mastersAbroad.scores[i].marksObtained
                              ieltsID = this.selectedTalent.mastersAbroad.scores[i].id
                        }
                  }

                  let tempMastersAbroad: any = {
                        id: this.selectedTalent.mastersAbroad.id,
                        degree: this.selectedTalent.mastersAbroad.degree,
                        countries: this.selectedTalent.mastersAbroad.countries,
                        universities: this.selectedTalent.mastersAbroad.universities,
                        yearOfMS: this.selectedTalent.mastersAbroad.yearOfMS,
                        gre: greScore,
                        greID: greID,
                        gmat: gmatScore,
                        gmatID: gmatID,
                        toefl: toeflScore,
                        toeflID: toeflID,
                        ielts: ieltsScore,
                        ieltsID: ieltsID
                  }
                  this.selectTalentTemp.mastersAbroad = tempMastersAbroad
            }
            else {
                  this.selectTalentTemp.mastersAbroad = null
            }

            // State.
            if (this.selectedTalent.country != undefined) {
                  this.getStateListByCountry(this.selectedTalent.country)
            }

            // Populate talent form.
            this.talentForm.patchValue(this.selectTalentTemp)

            // Disable talent form.
            this.talentForm.disable()

            // Open talent form modal.
            this.openModal(this.talentFormModal, 'xl')
      }

      // On clicking update talent button.
      onUpdateTalentClick(): void {
            this.isViewMode = false
            this.isOperationTalentUpdate = true
            this.enableTalentForm()
            if (this.talentForm.get('state').value == null) {
                  this.talentForm.get('state').disable()
            }
      }

      // Update Talent.
      updateTalent(): void {
            let talent: ITalent = this.assignMastersAbroad("update")
            this.patchIDFromObjectsForTalent(talent)
            this.spinnerService.loadingMessage = "Updating Talent"


            this.talentService.updateTalent(talent).subscribe((response: any) => {
                  this.modalRef.close('success')
                  this.getTalents()
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

      //*********************************************DELETE FOR TALENT FUNCTIONS************************************************************
      // On clicking delete talent button. 
      onDeleteTalentClick(talentID: string): void {
            this.assignSelectedTalent(talentID)
            this.openModal(this.deleteTalentModal, 'md').result.then(() => {
                  this.deleteTalent()
            }, (err) => {
                  console.error(err)
                  return
            })
      }

      // Delete talent.
      deleteTalent(): void {
            this.spinnerService.loadingMessage = "Deleting Talent"
            this.modalRef.close()

            this.talentService.deleteTalent(this.selectedTalent.id).subscribe((response: any) => {
                  this.getTalents()
                  this.fileOps.deleteUploadedFile(this.selectedTalent.resume)
                  this.fileOps.deleteUploadedFile(this.selectedTalent.image)
                  alert(response)
            }, (error) => {
                  console.error(error)

                  if (error.error) {
                        if (error.error.error) {
                              alert(error.error.error)
                              return
                        }
                        alert(error.error)
                        return
                  }
                  alert(error.statusText)
            })
      }

      // Delete resume.
      deleteResume(): void {
            this.fileOps.deleteUploadedFile().subscribe((data: any) => {
            }, (error) => {
                  console.error(error)
            })
      }

      // Delete image.
      deleteImage(): void {
            this.fileOps.deleteUploadedFile().subscribe((data: any) => {
            }, (error) => {
                  console.error(error)
            })
      }

      //*********************************************FUNCTIONS FOR TALENT FORM************************************************************
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

      // Delete masters abroad from talent form.
      removeMastersAbroad(): void {
            this.talentForm.removeControl('mastersAbroad')
            this.universityList = []
      }

      // On clicking currently working or not add or remove toDate field to experiences arrat form.
      isWorkingClicked(isWorking: string, index: number, expereince: any): void {
            if (isWorking == "false") {
                  expereince.addControl('toDate', this.formBuilder.control(null, [Validators.required]))
                  this.indexOfCurrentWorkingExp = -1
                  return
            }
            expereince.removeControl('toDate')
            this.indexOfCurrentWorkingExp = index
      }

      // Delete college from talent form.
      deleteAcademic(index: number): void {
            if (confirm("Are you sure you want to delete academic?")) {
                  this.talentForm.markAsDirty()
                  this.talentAcademicControlArray.removeAt(index)
            }
      }

      // Delete experience from talent form.
      deleteExperience(index: number): void {
            if (confirm("Are you sure you want to delete experience?")) {
                  if (this.indexOfCurrentWorkingExp == index) {
                        this.indexOfCurrentWorkingExp = -1
                  }
                  if (this.indexOfCurrentWorkingExp > index) {
                        this.indexOfCurrentWorkingExp = this.indexOfCurrentWorkingExp - 1
                  }
                  this.talentForm.markAsDirty()
                  this.talentExperienceControlArray.removeAt(index)
                  //this.otherDesignation[index] = false
            }
      }

      // Get experiences array from talent form.
      get talentExperienceControlArray(): FormArray {
            return this.talentForm.get("experiences") as FormArray
      }

      // Get academics array from talent form.
      get talentAcademicControlArray(): FormArray {
            return this.talentForm.get("academics") as FormArray
      }

      // Get career plans array from career plan form.
      get talentCareerPlansControlArray(): FormArray {
            return this.careerPlanForm.get("careerPlans") as FormArray
      }

      // Assign selected talent.
      assignSelectedTalent(id: string): void {
            let length = this.talents.length
            for (let index = 0; index < length; index++) {
                  if (this.talents[index].id == id) {
                        this.selectedTalent = this.talents[index]
                  }
            }
      }

      // Change page for pagination.
      changePage($event): void {
            this.offset = $event - 1
            this.currentPage = $event
            this.getTalents()
      }

      // Set total talents list on current page.
      setPaginationString(): void {
            this.paginationStart = this.limit * this.offset + 1
            this.paginationEnd = +this.limit + this.limit * this.offset
            if (this.totalTalents < this.paginationEnd) {
                  this.paginationEnd = this.totalTalents
            }
      }

      // On clicking isExperience checkbox, add or remove experiences array.
      toggleExperienceControls(event: any): void {
            if (this.showExperiencesInForm) {
                  if (confirm("Are you sure you want to delete all experience details?")) {
                        this.talentForm.setControl('experiences', this.formBuilder.array([]))
                        this.talentForm.get("isExperience").setValue(false)
                        this.showExperiencesInForm = false
                        return
                  }
                  event.target.checked = true
                  return
            }
            if (this.talentForm.get("academicYear").value != 5) {
                  event.target.checked = false
                  alert("Experience details can only be added when academic year is selected as completed")
                  return
            }
            this.addExperience()
            this.showExperiencesInForm = true
            this.talentForm.get("isExperience").setValue(true)
      }

      // On clicking isAcademic checkbox, add or remove academics array.
      toggleAcademicControls(event: any): void {
            if (this.showAcademicsInForm) {
                  if (confirm("Are you sure you want to delete all academic details?")) {
                        this.showAcademicsInForm = false
                        this.talentForm.setControl('academics', this.formBuilder.array([]))
                        this.talentForm.get("isAcademic").setValue(false)
                        return
                  }
                  event.target.checked = true
                  return
            }
            this.showAcademicsInForm = true
            this.talentForm.get("isAcademic").setValue(true)
            this.addAcademic()
      }

      // On uplaoding resume.
      onResumeSelect(event: any): void {
            this.resumeDocStatus = ""
            let files = event.target.files
            if (files && files.length) {
                  let file = files[0]

                  // Upload resume if it is present.]
                  this.isResumeUploading = true
                  this.fileOps.uploadResume(file).subscribe((data: any) => {
                        this.talentForm.markAsDirty()
                        this.talentForm.patchValue({
                              resume: data
                        })
                        this.resumeDisplayedFileName = file.name
                        this.isResumeUploading = false
                        this.isResumeUploadedToServer = true
                        this.resumeDocStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
                  }, (error) => {
                        this.isResumeUploading = false
                        this.resumeDocStatus = `<p><span>&#10060;</span> ${error}</p>`
                  })
            }
      }

      // On uplaoding image.
      onImageSelect(event: any): void {
            this.imageDocStatus = ""
            let files = event.target.files
            if (files && files.length) {
                  let file = files[0]

                  // Upload image if it is present.
                  this.isImageUploading = true
                  this.fileOps.uploadTalentImage(file).subscribe((response: any) => {
                        this.talentForm.markAsDirty()
                        this.talentForm.patchValue({
                              image: response
                        })
                        this.imageDisplayedFileName = file.name
                        this.isImageUploading = false
                        this.isImageUploadedToServer = true
                        this.imageDocStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
                  }, (error) => {
                        this.isImageUploading = false
                        this.imageDocStatus = `<p><span>&#10060;</span> ${error}</p>`
                  })
            }
      }

      // Show specialization according to degree id.
      showSpecificSpecializations(academic: any, specialization: ISpecialization): boolean {
            if (!academic.get('degree').value) {
                  return false
            }
            if (academic.get('degree').value.id === specialization.degreeID) {
                  return true
            }
            return false
      }

      // Used to dismiss modal.
      dismissFormModal(modal: NgbModalRef): void {
            if (this.isResumeUploading || this.isImageUploading) {
                  alert("Please wait till file is being uploaded")
                  return
            }
            if (this.isResumeUploadedToServer) {
                  if (!confirm("Uploaded resume will be deleted.\nAre you sure you want to close?")) {
                        return
                  }
                  this.deleteResume()
            }
            if (this.isImageUploadedToServer) {
                  if (!confirm("Uploaded image will be deleted.\nAre you sure you want to close?")) {
                        return
                  }
                  this.deleteImage()
            }
            modal.dismiss()
            this.isResumeUploadedToServer = false
            this.disableAddTalentsButton = false
            this.isImageUploadedToServer = false
            this.resumeDisplayedFileName = "Select file"
            this.imageDisplayedFileName = "Select file"
            this.resumeDocStatus = ""
            this.imageDocStatus = ""
      }

      // Used to open modal.
      openModal(content: any, size?: string): NgbModalRef {
            this.disableAddTalentsButton = false
            this.isExcelUploaded = false
            this.isResumeUploadedToServer = false
            this.isImageUploadedToServer = false
            this.resumeDisplayedFileName = "Select file"
            this.imageDisplayedFileName = "Select file"
            this.resumeDocStatus = ""
            this.imageDocStatus = ""
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

      // Validate talent form.
      validateTalentForm(): void {
            if (this.talentForm.invalid) {
                  this.talentForm.markAllAsTouched()
                  return
            }
            if (this.isOperationTalentUpdate) {
                  this.updateTalent()
                  return
            }
            this.addTalent()
      }

      // Get all invalid controls in talent form.
      public findInvalidControls(): any[] {
            const invalid = []
            const controls = this.talentForm.controls
            for (const name in controls) {
                  if (controls[name].invalid) {
                        invalid.push(name)
                  }
            }
            return invalid
      }

      // Calculate total years of experience for one experience for view mode.
      calculateYearsOfExperienceForViewMode(experience: IExperienceDTO): void {
            let toDate
            let fromDate
            if (experience.fromDate != null && experience.fromDate != "") {
                  if (experience.toDate == null || experience.toDate == "") {
                        toDate = new Date()
                  }
                  else {
                        toDate = new Date(experience.toDate)
                  }
                  fromDate = new Date(experience.fromDate)
                  let monthDiff: number = toDate.getMonth() - fromDate.getMonth() + (12 * (toDate.getFullYear() - fromDate.getFullYear()))
                  let numberOfYears: number = Math.floor(monthDiff / 12)
                  let numberOfMonths: number = Math.round(((monthDiff / 12) % 1) * 12)
                  experience.yearsOfExperience = numberOfYears + "." + numberOfMonths + " Year(s)"
                  experience.yearsOfExperienceInNumber = +(monthDiff / 12).toFixed(2)
            }
      }

      // Calculate total years of experience for talent in string.
      calculateTotalYearsOfExperinceOfTalentInString(talent: ITalentDTO): string {
            let totalYears: number = 0
            talent.totalYearsOfExperience = totalYears
            if ((talent.experiences) && (talent.experiences.length != 0)) {
                  for (let i = 0; i < talent.experiences.length; i++) {
                        this.calculateYearsOfExperienceForViewMode(talent.experiences[i])
                        totalYears = totalYears + talent.experiences[i].yearsOfExperienceInNumber
                  }
                  talent.totalYearsOfExperience = totalYears
                  let numberOfYears: number = Math.floor(totalYears)
                  let numberOfMonths: number = Math.round((totalYears % 1) * 12)
                  if (numberOfYears == 0 && numberOfMonths == 0) {
                        talent.totalYearsOfExperience = 0
                        talent.totalYearsOfExperienceInString = "Fresher"
                        return
                  }
                  if (isNaN(numberOfYears) || isNaN(numberOfMonths)) {
                        talent.totalYearsOfExperience = 0
                        talent.totalYearsOfExperienceInString = "Fresher"
                        return
                  }
                  talent.totalYearsOfExperienceInString = numberOfYears + "." + numberOfMonths + " Year(s)"
                  return
            }
            talent.totalYearsOfExperience = 0
            talent.totalYearsOfExperienceInString = "Fresher"
      }

      // Calculates the years of experience for each experience before add or update of talent.
      calculateYearsOfExperience(experience: IExperience): number {
            let toDate
            let fromDate
            if (experience.fromDate != null && experience.fromDate != "") {
                  if (experience.toDate == null || experience.toDate == "") {
                        toDate = new Date()
                  }
                  else {
                        toDate = new Date(experience.toDate)
                  }
                  fromDate = new Date(experience.fromDate)
                  let monthDiff: number = toDate.getMonth() - fromDate.getMonth() + (12 * (toDate.getFullYear() - fromDate.getFullYear()))
                  return monthDiff
            }
      }

      // Calculate total years of experience for talent in number before add or update of talent.
      calculateTotalYearsOfExperinceOfTalent(talent: ITalent): string {
            let monthDiff: number = 0
            if ((talent.experiences) && (talent.experiences.length != 0)) {
                  for (let i = 0; i < talent.experiences.length; i++) {
                        monthDiff = this.calculateYearsOfExperience(talent.experiences[i]) + monthDiff
                  }
                  talent.experienceInMonths = monthDiff
                  return
            }
            talent.experienceInMonths = monthDiff
            return
      }

      // On clicking isMastersAbroad checkbox, add or remove masters aborad.
      toggleShowMSInUsFormControls(event: any): void {
            if (this.showMastersAbroadInForm) {
                  if (confirm("Are you sure you want to delete all MS abroad details?")) {
                        this.showMastersAbroadInForm = false
                        this.removeMastersAbroad()
                        this.talentForm.get("isMastersAbroad").setValue(false)
                        return
                  }
                  event.target.checked = true
                  return
            }
            this.showMastersAbroadInForm = true
            this.talentForm.get("isMastersAbroad").setValue(true)
            this.addMastersAbroad()
            this.talentForm.get('mastersAbroad').get('universities').disable()
      }

      // On changing degree in masters abroad, add or remove gmat score field.
      onDegreeChange(degree: IDegree): void {
            if (degree && degree.name == "M.B.A.") {
                  this.talentForm.get('mastersAbroad').get('gmat').setValidators([Validators.required, Validators.min(0), Validators.max(this.gmatExam?.totalMarks)])
                  this.talentForm.get('mastersAbroad').get('gmat').updateValueAndValidity()
                  this.showMBADegreeRequiredError = true
                  return
            }
            this.talentForm.get('mastersAbroad').get('gmat').setValidators([Validators.min(0), Validators.max(this.gmatExam?.totalMarks)])
            this.talentForm.get('mastersAbroad').get('gmat').updateValueAndValidity()
            this.showMBADegreeRequiredError = false
      }

      // Toggle visibility of experience columns.
      toggleShowExperienceColumns(): void {
            this.showExperienceColumns = !this.showExperienceColumns
            // this.getTalents()
      }

      // Calculates current company name of talent.
      calculateCurrentCompanyName(talent: ITalentDTO): void {
            if (talent.experiences && talent.experiences.length != 0 && talent.experiences[talent.experiences.length - 1]?.toDate == null) {
                  talent.currentCompanyName = talent.experiences[talent.experiences.length - 1]?.company
                  return
            }
            talent.currentCompanyName = "Not Currently Working"
      }

      // Calculate the last designation of talent.
      calculateDesignation(talent: ITalentDTO): void {
            if (talent.experiences && talent.experiences.length != 0) {
                  talent.lastDesignation = talent.experiences[talent.experiences.length - 1]?.designation?.position
                  return
            }
            talent.lastDesignation = "No experience"
      }

      // Calculate the package of current company of talent.
      calculateCurrentCompanyPackage(talent: ITalentDTO): void {
            if (talent.experiences && talent.experiences.length != 0 && talent.experiences[talent.experiences.length - 1]?.toDate == null
                  && talent.experiences[talent.experiences.length - 1]?.package != null) {
                  talent.currentCompanyPackage = talent.experiences[talent.experiences.length - 1]?.package / 100000 + " Lpa"
                  return
            }
            talent.currentCompanyPackage = "Package not mentioned"
      }

      // Get the technologies of all the experiences of talent.
      calculateAllExperiencesTechnologies(talent: ITalentDTO): void {
            talent.allExperiencesTechnologies = []
            if (talent.experiences && talent.experiences.length != 0) {
                  for (let i = 0; i < talent.experiences.length; i++) {
                        if (talent.experiences[i]?.technologies && talent.experiences[i].technologies?.length > 0) {
                              for (let j = 0; j < talent.experiences[i].technologies.length; j++) {
                                    if (talent.allExperiencesTechnologies.includes(talent.experiences[i].technologies[j].language)) {
                                          continue
                                    }
                                    talent.allExperiencesTechnologies.push(talent.experiences[i].technologies[j].language)
                              }
                        }
                  }
            }
      }

      // Calculate the talent's first college name.
      calculateFirstCollegeName(talent: ITalentDTO): void {
            if (talent.academics && talent.academics.length != 0 && talent.academics[0]?.college != null) {
                  talent.firstCollegeName = talent.academics[0]?.college
                  return
            }
            talent.firstCollegeName = "Not Mentioned"
      }

      // Enable the talent form.
      enableTalentForm(): void {
            this.talentForm.enable()
            this.talentForm.get('code').disable()
      }

      // Navigate to interviews page for the talent.
      getInterviewSchedulesByTalent(talentID: string, talentFirstName: string, talentLastName: string): void {
            this.router.navigate([this.urlConstant.INTERVIEW_SCHEDULE], {
                  queryParams: {
                        "talentID": talentID,
                        "talentName": talentFirstName + " " + talentLastName
                  }
            }).catch(err => {
                  console.error(err)
            })
      }

      // Sort the talents in ascending or desecending order by given property of talent.
      sortTalents(property: string): void {
            this.talents.sort(this.dynamicSort(property))
      }

      // Sort the talents in ascending ot desecending order by given multiple properties of talent.
      sortTalentsMultipleProperty(...properties: string[]): void {
            this.talents.sort(this.dynamicSortMultiple(...properties))
      }

      // Sort method for sorting by one property.
      dynamicSort(property) {
            var sortOrder = 1
            if (property[0] === "-") {
                  sortOrder = -1
                  property = property.substr(1)
            }
            return function (a, b) {
                  /* next line works with strings and numbers, 
                   * and you may want to customize it to your needs
                   */
                  var result = (a[property] < b[property]) ? -1 : (a[property] > b[property]) ? 1 : 0
                  return result * sortOrder
            }
      }

      // Sort method for sorting by multiple property.
      dynamicSortMultiple(...properties: string[]) {
            /*
             * save the arguments object as it will be overwritten
             * note that arguments object is an array-like object
             * consisting of the names of the properties to sort by
             */
            var props = properties
            return function (obj1, obj2) {
                  var i = 0, result = 0, numberOfProperties = props.length
                  /* try getting a different result from 0 (equal)
                   * as long as we have extra properties to compare
                   */
                  while (result === 0 && i < numberOfProperties) {
                        result = this.dynamicSort(props[i])(obj1, obj2)
                        i++
                  }
                  return result
            }
      }

      // Calculate some fields of all talents.
      calculateFieldsOfTalent(talents: ITalentDTO[]): void {
            for (let i = 0; i < talents.length; i++) {
                  this.calculateCurrentCompanyName(talents[i])
                  this.calculateCurrentCompanyPackage(talents[i])
                  this.calculateAllExperiencesTechnologies(talents[i])
                  this.calculateDesignation(talents[i])
                  this.calculateFirstCollegeName(talents[i])
                  this.calculateTotalYearsOfExperinceOfTalentInString(talents[i])
                  if (talents[i].expectedCTC && talents[i].expectedCTC > 0) {
                        let ctc: number = talents[i].expectedCTC / 100000
                        talents[i].expectedCTCInString = ctc + " Lpa"
                  }
            }
      }

      // On clearing all countries in masters abroad countries.
      onClearCountries(): void {
            this.universityList = []
            // let arr = this.talentForm.get('mastersAbroad') as FormArray
            // Clear all universities. ************** -n
            this.talentForm.get('mastersAbroad')?.get('universities').reset()
            this.talentForm.get('mastersAbroad').get('universities').disable()
            // Change it to remove specific universities
            for (let universities of Array.from(this.universityMap.values())) {
                  for (let university of universities) {
                        university.isVisible = false
                  }
            }
      }

      // On removing one country in masters abroad countries.
      onRemoveCountry(id: any): void {
            this.universityList = []

            // let arr = this.talentForm.get('mastersAbroad') as FormArray
            let selectedUniversities: IUniversity[] = this.talentForm.get('mastersAbroad').get('universities').value
            // Change it to remove specific universities *********** -n
            this.talentForm.get('mastersAbroad').get('universities').reset()
            // // Removing universities of the removed countries from form
            // for (let [index, university] of selectedUniversities.entries()) {
            //       if (university.country.id == id) {
            //       }
            // }

            for (let [countryID, universities] of this.universityMap) {
                  if (countryID == id) {
                        for (let university of universities) {
                              university.isVisible = false
                        }
                        continue
                  }
                  for (let university of universities) {
                        if (!university.isVisible) {
                              break
                        }
                        this.universityList = this.universityList.concat(university)
                  }
            }
      }

      // On adding one country in masters abroad countries.
      onAddCountry(countryID: string): void {
            this.talentForm.get('mastersAbroad').get('universities').enable()
            if (this.universityMap.has(countryID)) {
                  let allUniversities = this.universityMap.get(countryID)
                  for (let university of allUniversities) {
                        university.isVisible = true
                        this.universityList = this.universityList.concat(university)
                  }
                  return
            }
            this.areUniversitiesLoading = true
            this.generalService.getUniversityByCountryID(countryID).subscribe((response: IUniversity[]) => {
                  for (let i = 0; i < response.length; i++) {
                        response[i].isVisible = true
                        response[i].countryName = response[i].country?.name
                  }
                  this.areUniversitiesLoading = false
                  this.universityMap.set(countryID, response)
                  this.universityList = this.universityList.concat(response)
            }, err => {
                  this.areUniversitiesLoading = false
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Will be called on [addTag]= "addCollegeToList" args will be passed automatically.
      addCollegeToList(option: any): Promise<any> {
            return new Promise((resolve) => {
                  resolve(option)
            }
            )
      }

      // Extract ID from objects in talent form.
      patchIDFromObjectsForTalent(talent: ITalent): void {
            if (this.talentForm.get('talentSource').value) {
                  talent.sourceID = this.talentForm.get('talentSource').value.id
                  delete talent['talentSource']
            }
            if (this.talentForm.get('salesPerson').value) {
                  talent.salesPersonID = this.talentForm.get('salesPerson').value.id
                  delete talent['salesPerson']
            }
            for (let i = 0; i < this.talentExperienceControlArray.length; i++) {
                  if (this.talentExperienceControlArray.at(i).get('designation').value) {
                        talent.experiences[i].designationID = this.talentExperienceControlArray.at(i).get('designation').value.id
                        delete talent.experiences['designation']
                  }
            }
            for (let i = 0; i < this.talentAcademicControlArray.length; i++) {
                  if (this.talentAcademicControlArray.at(i).get('degree').value) {
                        talent.academics[i].degreeID = this.talentAcademicControlArray.at(i).get('degree').value.id
                        delete talent.academics['degree']
                  }
                  if (this.talentAcademicControlArray.at(i).get('specialization').value) {
                        talent.academics[i].specializationID = this.talentAcademicControlArray.at(i).get('specialization').value.id
                        delete talent.academics['specialization']
                  }
            }
            this.calculateTotalYearsOfExperinceOfTalent(talent)
      }

      // Toggle between my and all talents.
      toggleMyAllTalents(): void{
            this.isViewAllTalents = !this.isViewAllTalents
            if (this.isHeadSalesPerson && this.isViewAllTalents){
                  this.roleName = "" 
            }
            if (this.isHeadSalesPerson && !this.isViewAllTalents){
                  this.roleName =  this.localService.getJsonValue("roleName")
            }
            if (this.isSearched) {
                  this.getAllSearchedTalents()
                  return
            }
            this.getTalents()
      }

      //*********************************************CRUD FUNCTIONS FOR TALENT CALL RECORDS************************************************************
      // On clicking get all call records by talent id.
      getCallRecordsForSelectedTalent(talentID: any): void {
            this.isOperationCallRecordUpdate = false
            this.talentCallRecordList = []
            this.spinnerService.loadingMessage = "Getting all call records"
            this.assignSelectedTalent(talentID)
            this.getAllCallRecords()
            this.showCallRecordForm = false
            this.openModal(this.callRecordFormModal, 'xl')
      }

      // Get all call records by talent id.
      getAllCallRecords(): void {
            this.spinnerService.loadingMessage = "Getting All Call Records"


            this.talentService.getCallRecordsByTalent(this.selectedTalent.id).subscribe(response => {
                  this.talentCallRecordList = response
                  this.formatDateTimeOfTalentCallRecords()
                  this.getTalents()
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // On clicking add new call record button.
      onAddNewCallRecordButtonClick(): void {
            this.isOperationCallRecordUpdate = false
            this.outcomeList = []
            this.showPurposeJobFields = false
            this.showOutcomeTargetDateField = false
            this.createCallRecordForm()
            this.showCallRecordForm = true
            this.callRecordForm.get('outcome').disable()
      }

      // Add call record.
      addCallRecord(): void {
            this.spinnerService.loadingMessage = "Adding Call Record"


            let callRecord: ICallRecord = this.callRecordForm.value
            this.patchIDFromObjectsForCallRecord(callRecord)
            this.talentService.addCallRecord(this.callRecordForm.value, this.selectedTalent.id).subscribe((response: any) => {
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
            this.showPurposeJobFields = false
            this.showOutcomeTargetDateField = false
            this.createCallRecordForm()
            this.showCallRecordForm = true
            if (this.talentCallRecordList[index].purpose) {
                  this.getOutcomesByPurpose(this.talentCallRecordList[index].purpose)
            }
            if (this.talentCallRecordList[index].targetDate) {
                  this.addTargetDateInCallRecordForm()
            }
            this.callRecordForm.patchValue(this.talentCallRecordList[index])
      }

      // Update call record.
      updateCallrecord(): void {
            this.spinnerService.loadingMessage = "Updating Call Record"


            let callRecord: ICallRecord = this.callRecordForm.value
            this.patchIDFromObjectsForCallRecord(callRecord)
            this.talentService.updateCallRecord(this.callRecordForm.value, this.selectedTalent.id).subscribe((response: any) => {
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
      deleteCallRecord(talentCallRecordID: string): void {
            if (confirm("Are you sure you want to delete the call record?")) {
                  this.spinnerService.loadingMessage = "Deleting Call Record"


                  this.talentService.deleteCallRecord(talentCallRecordID, this.selectedTalent.id).subscribe((response: any) => {
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

      // Toggle visibility of expected CTC and notice period fields of call record.
      togglePurposeJobFields(purpose: IPurpose): void {
            if (purpose && purpose.purpose == "Placement") {
                  this.showPurposeJobFields = true
                  this.callRecordForm.addControl('expectedCTC', this.formBuilder.control(null, [Validators.min(100000), Validators.max(100000000)]))
                  this.callRecordForm.addControl('noticePeriod', this.formBuilder.control(null, [Validators.min(0), Validators.max(9)]))
                  return
            }
            this.showPurposeJobFields = false
            this.callRecordForm.removeControl('expectedCTC')
            this.callRecordForm.removeControl('noticePeriod')
      }

      // Toggle visibility of target field of call record.
      toggleOutcomeTargetDateField(outcome: IOutcome): void {
            if (outcome && outcome.outcome == "Follow Up") {
                  this.addTargetDateInCallRecordForm()
                  return
            }
            this.removetargetDateFromCallRecordForm()
      }

      // Add target field to call record.
      addTargetDateInCallRecordForm(): void {
            this.showOutcomeTargetDateField = true
            this.callRecordForm.addControl('targetDate', this.formBuilder.control(null, [Validators.required]))
      }

      // Remove target field from call record.
      removetargetDateFromCallRecordForm(): void {
            this.showOutcomeTargetDateField = false
            if (this.callRecordForm.contains('targetDate')) {
                  this.callRecordForm.removeControl('targetDate')
            }
      }

      // Format date time field of talent call records by removing timestamp.
      formatDateTimeOfTalentCallRecords(): void {
            for (let i = 0; i < this.talentCallRecordList?.length; i++) {
                  let dateTime = this.talentCallRecordList[i].dateTime
                  if (dateTime) {
                        this.talentCallRecordList[i].dateTime = this.datePipe.transform(dateTime, 'yyyy-MM-ddTHH:mm:ss')
                  }
                  let targetDate = this.talentCallRecordList[i].targetDate
                  if (targetDate) {
                        this.talentCallRecordList[i].targetDate = this.datePipe.transform(targetDate, 'yyyy-MM-dd')
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

      //*********************************************CRUD FUNCTIONS FOR TALENT LIFETIME VALUE************************************************************
      // On clicking lifetime value button in talents table.
      getLifetimeValueForSelectedTalent(talentID: string, lifetimeValue: number): void {
            if (!this.showForFaculty) {
                  return
            }
            this.selectedLifetimeValue = { upsell: null, teaching: null, knowledge: null, placement: null }
            this.spinnerService.loadingMessage = "Getting lifetime value details"
            this.showLifetimeForm = false
            this.showLifetimeValue = true
            this.assignSelectedTalent(talentID)
            this.createLifetimeValueForm()
            if (lifetimeValue == null) {
                  this.isOperationLifetimeValueUpdate = false
                  this.createLifetimeValueForm()
                  this.openModal(this.lifetimeValueFormModal, 'xl')
                  return
            }
            this.isOperationLifetimeValueUpdate = true
            this.openModal(this.lifetimeValueFormModal, 'xl')
            this.getLifetimeValue()
      }

      // Get lifetime value details by talent id.
      getLifetimeValue(): void {
            this.spinnerService.loadingMessage = "Getting Lifetime Value"


            this.talentService.getLifetimeValue(this.selectedTalent.id).subscribe(response => {
                  this.selectedLifetimeValue = response
                  this.getTalents()
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // On clicking update or add lifetime value button click in life time value details.
      onUpdateLifetimeValueButtonClick(): void {
            this.showLifetimeValue = false
            this.showLifetimeForm = true
            this.createLifetimeValueForm()
            if (this.isOperationLifetimeValueUpdate) {
                  this.lifetimeValueForm.patchValue(this.selectedLifetimeValue)
            }
            return
      }

      // On clicking update or add lifetime value button click in life time value form.
      onUpdateOrAddButtonClickInForm(): void {
            if (!this.isOperationLifetimeValueUpdate) {
                  this.addLifetimeValue()
            }
            else {
                  this.updateLifetimeValue()
            }
      }

      // Add lifetime value.
      addLifetimeValue(): void {
            this.spinnerService.loadingMessage = "Adding Lifetime Value"


            this.talentService.addLifetimeValue(this.lifetimeValueForm.value, this.selectedTalent.id).subscribe((response: any) => {
                  this.showLifetimeForm = false
                  this.showLifetimeValue = true
                  this.isOperationLifetimeValueUpdate = true
                  this.getLifetimeValue()
                  alert(response.body)
            }, (error) => {
                  console.error(error)
                  if (error.error?.error) {
                        alert(error.error?.error)
                        return
                  }
                  alert(error.statusText)
            })
      }

      // Update lifetime value.
      updateLifetimeValue(): void {
            this.spinnerService.loadingMessage = "Updating Lifetime Value"


            this.talentService.updateLifetimeValue(this.lifetimeValueForm.value, this.selectedTalent.id).subscribe((response: any) => {
                  this.showLifetimeForm = false
                  this.showLifetimeValue = true
                  this.getLifetimeValue()
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

      // Delete lifetime value.
      deleteLifetimeValue(): void {
            if (confirm("Are you sure you want to delete the lifetime value details?")) {
                  this.spinnerService.loadingMessage = "Deleting Lifetime value"


                  this.talentService.deleteLifetimeValue(this.selectedLifetimeValue.id, this.selectedTalent.id).subscribe((response: any) => {
                        alert(response)
                        this.showLifetimeForm = false
                        this.showLifetimeValue = true
                        this.isOperationLifetimeValueUpdate = false
                        this.selectedLifetimeValue = { upsell: null, teaching: null, knowledge: null, placement: null }
                        this.getTalents()
                  }, (error) => {
                        console.error(error)
                        if (error.error?.error) {
                              alert(error.error?.error)
                              return
                        }
                        alert(error.statusText)
                  }).add(() => {

                  })
            }
      }

      // Validate lifetime value form.
      validateLifetimeValueForm(): void {
            if (this.lifetimeValueForm.invalid) {
                  this.lifetimeValueForm.markAllAsTouched()
                  return
            }
            if (this.lifetimeValueForm.get('upsell').value == null && this.lifetimeValueForm.get('placement').value == null &&
                  this.lifetimeValueForm.get('knowledge').value == null && this.lifetimeValueForm.get('teaching').value == null) {
                  alert("Please fill atleast one value field!")
                  return
            }
            if (this.lifetimeValueForm.get('upsell').value + this.lifetimeValueForm.get('placement').value + this.lifetimeValueForm.get('knowledge').value +
                  this.lifetimeValueForm.get('teaching').value > this.MAX_VALUE) {
                  alert("Total lifetime value cannot be above " + this.MAX_VALUE)
                  return
            }
            this.onUpdateOrAddButtonClickInForm()
      }

      // Calculate the total lifetime value to be displayed on the modal popup.
      getTotalLifetimeValue(): string {
            let total: number = 0
            if (this.selectedLifetimeValue) {
                  if (this.selectedLifetimeValue.upsell) {
                        total = total + this.selectedLifetimeValue.upsell
                  }
                  if (this.selectedLifetimeValue.placement) {
                        total = total + this.selectedLifetimeValue.placement
                  }
                  if (this.selectedLifetimeValue.knowledge) {
                        total = total + this.selectedLifetimeValue.knowledge
                  }
                  if (this.selectedLifetimeValue.teaching) {
                        total = total + this.selectedLifetimeValue.teaching
                  }
            }
            return this.formatNumberInIndianRupeeSystem(total)
      }

      // Format the number in indian rupee system.
      formatNumberInIndianRupeeSystem(number: number): string {
            if (!number) {
                  return
            }
            var result = number.toString().split('.')
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

      // Format the total lifetime value to be shown in talent table.
      formatLifetimeValueInShortForm(lifetimeValue: number): string {
            if (lifetimeValue < 100000) {
                  return "" + lifetimeValue
            }
            if ((lifetimeValue / 100000) < 100) {
                  return (lifetimeValue / 100000).toFixed(2) + "lakh(s)"
            }
            return (lifetimeValue / 10000000).toFixed(2) + "crore(s)"
      }

      // Show all courses taken by takent in lifetime value form.
      showAllCoursesTaken(talent: any): string {
            if (talent.courses != null && talent.courses?.length != 0) {
                  let coursesArray = []
                  for (let i = 0; i < talent.courses.length; i++) {
                        coursesArray.push(talent.courses[i].name)
                  }
                  let coursesTaken: string = coursesArray.join(", ")
                  return coursesTaken
            }
            return "No course taken"
      }

      //*********************************************FUNCTIONS FOR TALENT SEARCH FORM************************************************************
      // Reset search form and get all talents.
      resetSearchAndGetAll(): void {
            this.searchFilterFieldList = []
            this.spinnerService.loadingMessage = "Getting all talents"
            this.isSearched = false
            this.createSearchTalentForm()
            this.changePage(1)
      }

      // Reset search form.
      resetSearchForm(): void {
            this.searchFilterFieldList = []
            this.outcomeList = []
            this.talentSearchForm.reset()
            this.talentSearchForm.get('isActive').setValue(true)
      }

      // On clicking search button
      OnSearchTalentsButtonClick(): void {
            this.spinnerService.loadingMessage = "Getting Searched Talents"
            this.isSearched = true
            this.changePage(1)
      }

      // Get all searched talents.
      getAllSearchedTalents(): void {
            let roleNameAndLogin: any = {
                  roleName: this.roleName,
                  loginID: this.localService.getJsonValue("loginID"),
                  isViewAllBatches: this.isViewAllTalents?1:0
            }
            let data = this.talentSearchForm.value
            this.utilityService.deleteNullValuePropertyFromObject(data)
            this.spinnerService.loadingMessage = "Getting All Searched Talents"


            this.searchFilterFieldList = []
            for (var property in data) {
                  let text: string = property
                  let result: string = text.replace(/([A-Z])/g, " $1")
                  let finalResult: string = result.charAt(0).toUpperCase() + result.slice(1)
                  let valueArray: any[] = []
                  if (Array.isArray(data[property])) {
                        valueArray = data[property]
                  }
                  else {
                        valueArray.push(data[property])
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
            this.talentService.getAllSearchedTalents(data, this.limit, this.offset, roleNameAndLogin).subscribe(response => {
                  this.talents = response.body
                  this.calculateFieldsOfTalent(this.talents)
                  this.totalTalents = parseInt(response.headers.get('X-Total-Count'))
                  this.totalLifetimeValue = parseInt(response.headers.get("totalLifetimeValue"))
                  // Now we keep the checkbox as selected if the talent is selected by setting 
                  // the isAllocatedForRequirement to true.
                  // if (this.selectedTalents.length > 0) {
                  //       this.selectedTalents.forEach(selectedTalent => {
                  //             this.talents.forEach(talent => {
                  //                   if (talent.id == selectedTalent.id) {
                  //                         talent.isChecked = true
                  //                   }
                  //             })
                  //       })
                  // }
            }, (error) => {
                  console.error(error)
                  if (error.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                  }
            }).add(() => {

                  this.setPaginationString()
            })
      }

      // Set the limit of maximum experiences required.
      setMaximumExperienceRequired(value: number): void {
            this.talentSearchForm.get('maximumExperience').
                  setValidators([Validators.max(30), Validators.min(value)])
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

      // Delete search criteria from talent search form by search name.
      deleteSearchCriteria(searchName: string): void {
            this.talentSearchForm.get(searchName).setValue(null)
            this.getAllSearchedTalents()
      }

      //*********************************************COMPANY REQUIREMENT RELATED FUNCTIONS************************************************************
      // ************ Added by niranjan *********************
      talentSearchByRequirements(): void {
            // this.addSelectedTalents = []
            // this.createSearchTalentForm()

            // this.talentSearchForm.patchValue(this.companyData.talentSearchParams)
            this.getAllSearchedTalents()
      }

      addTalentToCompanyRequirement(): void {
            let requirementTalents = []
            this.selectedTalentsList.forEach(talent => {
                  requirementTalents.push({
                        "requirementID": this.requirementID,
                        "talentID": talent
                  })
            })
            if (requirementTalents.length == 0) {
                  alert("Please select talents")
                  return
            }
            this.spinnerService.loadingMessage = "Adding talents to company requirement"


            this.companyRequirementService.addTalentsToRequirement(this.requirementID, requirementTalents).subscribe((response: any) => {
                  alert(response)
                  this.router.navigateByUrl("/company/requirement")
            }, (error: any) => {
                  console.error(error)
                  if (error.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                  }
                  if (error.error) {
                        alert(error.error.error)
                        return
                  }
                  alert(error.error)
            })
      }

      //*********************************************BATCH RELATED FUNCTIONS************************************************************

      // Get talents for batch.
      getSelectedTalents(): void {
            this.spinnerService.loadingMessage = "Getting Talents"


            this.batchService.getTalentsForBatch(this.batchID, this.limit, this.offset).
                  subscribe((response: any) => {
                        this.talents = response.body
                        this.totalTalents = response.headers.get('X-Total-Count')
                        this.totalLifetimeValue = parseInt(response.headers.get("totalLifetimeValue"))
                        this.calculateFieldsOfTalent(this.talents)
                  },
                        (error) => {
                              console.error(error)
                              if (error.error) {
                                    alert(error.error)
                                    return
                              }
                              alert(error.error.error)
                        }).add(() => {


                              this.setPaginationString()
                        })

      }

      // Redirect to batch feedback page.
      redirectToBatchFeedback(batch: any): void {
            this.router.navigate(["/batch/master/feedback"], {
                  queryParams: {
                        "batchID": batch.id,
                        "batchName": batch.batchName,
                        "talentID": this.selectedTalentID
                  }
            }).catch(err => {
                  console.error(err)
            })
      }

      // Redirect to batch session page.
      redirectToBatchSession(batch: any, session: any): void {
            this.router.navigate(['/batch/master/session/feedback'], {
                  queryParams: {
                        "batchID": batch.id,
                        "batchName": batch.batchName,
                        "courseID": batch.course.id,
                        "sessionID": session.id,
                        "sessionName": session.session.name,
                        "talentID": this.selectedTalentID
                  }
            }).catch(err => {
                  console.error(err)
            })
      }

      // Redirect to batch details page.
      redirectToBatchDetails(batch: any): void {
            let url: string
            if (!this.showForFaculty){
                  url = "/my/batch/session/details"
            }
            if (this.showForFaculty){
                  url = "/training/batch/master/session/details"
            }
            this.router.navigate([url], {
                  queryParams: {
                        "batchID": batch.id,
                        "batchName": batch.batchName,
                        "courseID": batch.course?.id,
                        "tab": 1,
                        "subTab": "Manage"
                  }
            }).catch(err => {
                  console.error(err)
            })
      }

      // Not tested
      deleteTalentFromBatch(batchID: string): void {
            if (confirm("Are you sure you want to delete talent from this batch?")) {
                  this.spinnerService.loadingMessage = "Deleting talent from batch"


                  this.batchService.deleteTalentFromBatch(this.selectedTalentID, batchID).subscribe((response: any) => {
                        alert("Talent successfully deleted from the batch")
                        this.getBatchListOfOneTalent(this.selectedTalentID)
                  }, err => {
                        if (err.statusText.includes('Unknown')) {
                              alert("No connection to server. Check internet.")
                              return
                        }
                        console.error(err.error)
                        alert(err.error)
                  }).add(() => {

                  })
            }
      }

      //*********************************************ALLOCATION FUNCTIONS************************************************************
      // On clicking allocate one talent to batch button click.
      onAllocateOneTalentToBatchClick(talentID: any): void {
            this.showAllocateTalentsToBatch = false
            this.showAllocateOneTalentToBatch = true
            this.assignSelectedTalent(talentID)
            this.openModal(this.allocateToBatchModal, 'sm')
      }

      // On clicking allocate multiple talents to batch button click.
      onAllocateTalentsToBatchClick(): void {
            this.showAllocateTalentsToBatch = true
            this.showAllocateOneTalentToBatch = false
            this.openModal(this.allocateToBatchModal, 'sm')
      }

      // On clicking allocate salesperson to one talent button click.
      onAllocateSalespersonToOneTalentClick(talentID: any): void {
            this.showAllocateSalespersonToTalents = false
            this.showAllocateSalespersonToOneTalent = true
            this.assignSelectedTalent(talentID)
            this.openModal(this.allocateSalespersonModal, 'sm')
      }

      // On clicking allocate salesperson to multiple talents button click.
      onAllocateSalespersonToTalentsClick(): void {
            this.showAllocateSalespersonToTalents = true
            this.showAllocateSalespersonToOneTalent = false
            this.openModal(this.allocateSalespersonModal, 'sm')
      }

      // Toggle visibility of multiple select checkbox.
      toggleMultipleSelect(): void {
            if (this.multipleSelect) {
                  this.multipleSelect = false
                  this.setSelectAllTalents(this.multipleSelect)
                  return
            }
            this.multipleSelect = true
      }

      // Set isChecked field of all selected talents.
      setSelectAllTalents(isSelectedAll: boolean): void {
            for (let i = 0; i < this.talents.length; i++) {
                  this.addTalentToList(isSelectedAll, this.talents[i])
            }
      }

      // Check if all talents in selected talents are added in multiple select or not.
      checkTalentsAdded(): boolean {
            let count: number = 0

            for (let i = 0; i < this.talents.length; i++) {
                  if (this.selectedTalentsList.includes(this.talents[i].id))
                        count = count + 1
            }
            return (count == this.talents.length)
      }

      // Check if talent is added in multiple select or not.
      checkTalentAdded(talentID): boolean {
            return this.selectedTalentsList.includes(talentID)
      }

      // Takes a list called selectedTalent and adds all the checked talents to list, also does not contain duplicate values.
      addTalentToList(isChecked: boolean, talent: ITalentDTO): void {
            if (isChecked) {
                  if (!this.selectedTalentsList.includes(talent.id)) {
                        this.selectedTalentsList.push(talent.id)
                  }
                  return
            }
            if (this.selectedTalentsList.includes(talent.id)) {
                  let index = this.selectedTalentsList.indexOf(talent.id)
                  this.selectedTalentsList.splice(index, 1)
            }
      }

      // Allocate salesperson to talent(s).
      allocateSalesPersonToTalents(salesPersonID: string, talentID?: string): void {
            if (salesPersonID == "null") {
                  alert("Please select sales person")
                  return
            }

            let talentIDsToBeUpdated = []

            if (!talentID) {
                  for (let index = 0; index < this.selectedTalentsList.length; index++) {
                        talentIDsToBeUpdated.push({
                              "talentID": this.selectedTalentsList[index]
                        })
                  }
                  this.spinnerService.loadingMessage = "Salesperson is getting allocated to talents"
            }
            else {
                  talentIDsToBeUpdated.push({
                        "talentID": talentID
                  })
                  this.spinnerService.loadingMessage = "Salesperson is getting allocated to talent"
            }
            if (talentIDsToBeUpdated.length == 0) {
                  alert("Please select talents")
                  return
            }


            this.talentService.updateTalentsSalesPerson(talentIDsToBeUpdated, salesPersonID).subscribe((response: any) => {
                  this.getTalents()
                  alert(response)
                  this.modalRef.close('success')
                  this.selectedTalentsList = []
            }, (error) => {
                  console.error(error)
                  if (typeof error.error == 'object' && error) {
                        alert(this.utilityService.getErrorString(error))
                        return
                  }
                  if (error.error == undefined) {
                        alert('Sales person could not be allocated to talents')
                  }
                  alert(error.statusText)
            })
      }

      // Allocate talent(s) to batch.
      allocateTalentToBatch(batchID: string, talentID?: string): void {
            if (batchID == "null") {
                  alert("Please select batch")
                  return
            }
            let talentIDsToBeUpdated = []
            if (!talentID) {
                  for (let index = 0; index < this.selectedTalentsList.length; index++) {
                        talentIDsToBeUpdated.push({
                              "talentID": this.selectedTalentsList[index]
                        })
                  }
                  this.spinnerService.loadingMessage = "Talents are getting allocated to batch"
            }
            else {
                  talentIDsToBeUpdated.push({
                        "talentID": talentID
                  })
                  this.spinnerService.loadingMessage = "Talent is getting allocated to batch"

            }
            if (talentIDsToBeUpdated.length == 0) {
                  alert("Please select talents")
                  return
            }



            this.batchService.addTalentsToBatch(talentIDsToBeUpdated, batchID).subscribe(() => {
                  this.getTalents()
                  alert("Talent allocated to batch successfully")
                  this.modalRef.close('success')
                  this.selectedTalentsList = []
            }, (error) => {
                  console.error(error)
                  if (typeof error.error == 'object' && error) {
                        alert(this.utilityService.getErrorString(error))
                        return
                  }
                  if (error.error == undefined) {
                        alert('Talent could not be allocated, try again')
                  }
                  alert(error.statusText)
            })
      }

      //*********************************************NEXT ACTION FUNCTIONS************************************************************
      // On selecting action type.       
      onNextActionTypeChange(nextAction: any): void {
            if (!nextAction) {
                  return
            }
            this.createNextActionForm()
            this.showNextActionForm = true
            this.showNextActionTechnologies = false
            this.showNextActionTargetDate = false
            this.showNextActionStipend = false
            this.nextActionForm.get('talentID').setValue(this.selectedTalent.id)
            this.nextActionForm.get('actionType').setValue(nextAction)

            if (this.nextActionForm.get('actionType').value.type === "Course") {
                  this.onCourseSelect()
            }
            if (this.nextActionForm.get('actionType').value.type === "Blog") {
                  this.onBlogSelect()
            }
            if (this.nextActionForm.get('actionType').value.type === "Placement") {
                  this.onPlacementSelect()
            }
            if (this.nextActionForm.get('actionType').value.type === "Internship") {
                  this.onInternshipSelect()
            }
            if (this.nextActionForm.get('actionType').value.type == "Referral") {
                  this.onReferralSelect()
            }
            if (this.nextActionForm.get('actionType').value.type === "Teaching Assistant") {
                  this.onTeachingSelect()
            }
            if (this.nextActionForm.get('actionType').value.type !== "Course"
                  && this.nextActionForm.get('actionType').value.type !== "Referral") {
                  return
            }
      }

      // On selecting course action type.
      onCourseSelect(): void {
            this.nextActionForm.addControl('courses', this.formBuilder.control(null, [Validators.required]))
            this.nextActionForm.addControl('targetDate', this.formBuilder.control(null, [Validators.required]))
            this.showNextActionTargetDate = true
      }

      // On selecting blog action type.
      onBlogSelect(): void {
            this.nextActionForm.addControl('technologies', this.formBuilder.control(null, [Validators.required]))
            this.showNextActionTechnologies = true
            this.nextActionForm.addControl('targetDate', this.formBuilder.control(null, [Validators.required]))
            this.showNextActionTargetDate = true
      }

      // On selecting internship action type.
      onInternshipSelect(): void {
            this.nextActionForm.addControl('technologies', this.formBuilder.control(null, [Validators.required]))
            this.showNextActionTechnologies = true
            this.nextActionForm.addControl('stipend', this.formBuilder.control(null, [Validators.required]))
            this.showNextActionStipend = true
            this.nextActionForm.addControl('fromDate', this.formBuilder.control(null, [Validators.required]))
            this.nextActionForm.addControl('toDate', this.formBuilder.control(null, [Validators.required]))
            this.nextActionForm.addControl('targetDate', this.formBuilder.control(null, [Validators.required]))
            this.showNextActionTargetDate = true
      }

      // On selecting referral action type.
      onReferralSelect(): void {
            this.nextActionForm.addControl('referralCount', this.formBuilder.control(null, [Validators.required]))
      }

      // On selecting placement action type.
      onPlacementSelect(): void {
            this.nextActionForm.addControl('companies', this.formBuilder.control(null, [Validators.required]))
            this.nextActionForm.addControl('technologies', this.formBuilder.control(null, [Validators.required]))
            this.showNextActionTechnologies = true
            this.nextActionForm.addControl('stipend', this.formBuilder.control(null, [Validators.required]))
            this.showNextActionStipend = true
            this.nextActionForm.addControl('targetDate', this.formBuilder.control(null, [Validators.required]))
            this.showNextActionTargetDate = true
      }

      // On selecting teaching action type.
      onTeachingSelect(): void {
            this.nextActionForm.addControl('technologies', this.formBuilder.control(null, [Validators.required]))
            this.showNextActionTechnologies = true
            this.nextActionForm.addControl('stipend', this.formBuilder.control(null, [Validators.required]))
            this.showNextActionStipend = true
            this.nextActionForm.addControl('targetDate', this.formBuilder.control(null, [Validators.required]))
            this.showNextActionTargetDate = true
      }

      // On clicking next action button.
      onNextActionClick(talentID: string): void {
            this.isOperationNextActionUpdate = false
            this.showNextActionForm = false
            this.assignSelectedTalent(talentID)
            this.getTalentNextAction(talentID)
            this.openModal(this.nextActionFormModal, 'xl')
      }

      // On clicking update next action.
      onUpdateNextActionClick(nextAction: INextActionDTO): void {
            this.isOperationNextActionUpdate = true
            this.onNextActionTypeChange(nextAction.actionType)
            this.nextActionForm.patchValue(nextAction)
      }

      // On clicking add new next action.
      onAddNewNextActionClick(): void {
            this.isOperationNextActionUpdate = false
            this.createNextActionForm()
            this.showNextActionForm = true
            this.showNextActionTechnologies = false
            this.showNextActionTargetDate = false
            this.showNextActionStipend = false
      }

      // Validate next action form.
      validateNextActionForm(): void {
            if (this.nextActionForm.invalid) {
                  this.nextActionForm.markAllAsTouched()
                  return
            }
            if (this.isOperationNextActionUpdate) {
                  this.updateNextAction()
                  return
            }
            this.addNextAction()
      }

      // On clicking delete next action.
      onDeleteNextActionClick(nextActionID: string): void {
            if (confirm("Are you sure you want to delete this next action?")) {
                  this.spinnerService.loadingMessage = "Deleting next action"


                  this.talentService.deleteTalentNextAction(nextActionID, this.selectedTalent.id).subscribe((response: any) => {
                        this.getTalentNextAction(this.selectedTalent.id)
                        alert(response)
                  }, (error: any) => {
                        console.error(error)
                        if (error.error?.error) {
                              alert(error.error?.error)
                              return
                        }
                        alert(error.statusText)
                  }).add(() => {

                  })
            }
      }

      // Get next action for talent.
      getTalentNextAction(talentID: string): void {
            this.spinnerService.loadingMessage = "Getting next action"


            this.totalNextActions = 1
            this.nextActionList = []
            this.talentService.getAllTalentNextActions(talentID).subscribe((response: any) => {
                  this.nextActionList = response.body
                  this.totalNextActions = this.nextActionList.length
                  // #Niranjan
                  this.formatDatesOfTalentNextActions()
            }, (error: any) => {
                  console.error(error)
                  if (error.error?.error) {
                        alert(error.error?.error)
                        return
                  }
                  alert(error.statusText)
            })
      }

      // Adding new next action for talent.
      addNextAction(): void {
            this.spinnerService.loadingMessage = "Adding next action"


            let nextAction: INextAction = this.nextActionForm.value
            this.patchIDFromObjectsForNextAction(nextAction)
            this.talentService.addTalentNextAction(this.nextActionForm.value).subscribe((response: any) => {
                  this.getTalentNextAction(this.selectedTalent.id)
                  this.showNextActionForm = false
                  alert(response)
            }, (error: any) => {
                  console.error(error)
                  if (error.error?.error) {
                        alert(error.error?.error)
                        return
                  }
                  alert(error.statusText)
            })
      }

      // Updating next action for talent.
      updateNextAction(): void {
            this.spinnerService.loadingMessage = "Update next action"


            let nextAction: INextAction = this.nextActionForm.value
            this.patchIDFromObjectsForNextAction(nextAction)
            this.talentService.updateTalentNextAction(this.nextActionForm.value).subscribe((response: any) => {
                  this.getTalentNextAction(this.selectedTalent.id)
                  this.showNextActionForm = false
                  alert(response)
            }, (error: any) => {
                  console.error(error)
                  if (error.error?.error) {
                        alert(error.error?.error)
                        return
                  }
                  alert(error.statusText)
            })
      }

      // Format date fields field of talent next actions by removing timestamp.
      formatDatesOfTalentNextActions(): void {
            for (let index in this.nextActionList) {
                  let fromDate = this.nextActionList[index].fromDate
                  let toDate = this.nextActionList[index].toDate
                  let targetDate = this.nextActionList[index].targetDate
                  if (!fromDate && !toDate && !targetDate) {
                        continue
                  }
                  if (fromDate) {
                        this.nextActionList[index].fromDate = this.datePipe.transform(fromDate, 'yyyy-MM-dd')
                  }
                  if (toDate) {
                        this.nextActionList[index].toDate = this.datePipe.transform(toDate, 'yyyy-MM-dd')
                  }
                  if (targetDate) {
                        this.nextActionList[index].targetDate = this.datePipe.transform(targetDate, 'yyyy-MM-dd')
                  }
            }
      }

      // Extract ID from objects in next action form.
      patchIDFromObjectsForNextAction(nextAction: INextAction): void {
            if (this.nextActionForm.get('actionType').value) {
                  nextAction.actionTypeID = this.nextActionForm.get('actionType').value.id
                  delete nextAction['actionType']
            }
      }

      //*********************************************CRUD FUNCTIONS FOR TALENT CAREER PLAN************************************************************
      // On clicking get all career plans by talent id.
      getCareerPlansForSelectedTalent(talentID: any): void {
            this.careerPlanList = []
            this.spinnerService.loadingMessage = "Getting all career plans"
            this.assignSelectedTalent(talentID)
            this.getAllCareerPlans()
            this.showAddCareerPlanForm = false
            this.showUpdateCareerPlanForm = false
            this.selectedCareerPlan = {}
            this.openModal(this.careerPlanFromModal, 'xl')
      }

      // Get all career plans by talent id.
      getAllCareerPlans(): void {
            this.spinnerService.loadingMessage = "Getting Career Plans"


            this.talentService.getCareerPlansByTalent(this.selectedTalent.id).subscribe(response => {
                  this.careerPlanList = response
                  this.setFieldsOfAllCareerPlans()
                  this.makeCareerPlanMap(this.careerPlanList)
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // On clicking add new career plan button.
      onAddNewCareerPlanButtonClick(): void {
            this.createCareerPlanForm()
            this.showAddCareerPlanForm = true
            this.showUpdateCareerPlanForm = false
      }

      // On changing the career objective in form.
      onCareerObjectiveChange(): void {
            if (this.careerPlanForm.get('careerObjective').value == null || this.checkCareerObjectiveExists(this.careerPlanForm.get('careerObjective'))) {
                  this.careerPlanForm.setControl('careerPlans', this.formBuilder.array([]))
                  return
            }
            let loginID: string = this.localService.getJsonValue("loginID")
            this.careerPlanForm.setControl('careerPlans', this.formBuilder.array([]))
            for (let i = 0; i < this.careerPlanForm.get('careerObjective').value.courses.length; i++) {
                  this.addCareerPlan()
                  let careerOjbectiveCourse: any = this.careerPlanForm.get('careerObjective').value.courses[i]

                  this.talentCareerPlansControlArray.at(i).get('technicalAspect').setValue(careerOjbectiveCourse.technicalAspect)
                  this.talentCareerPlansControlArray.at(i).get('careerObjectiveID').setValue(this.careerPlanForm.get('careerObjective').value.id)
                  this.talentCareerPlansControlArray.at(i).get('careerObjectivesCoursesID').setValue(careerOjbectiveCourse.id)
                  this.talentCareerPlansControlArray.at(i).get('facultyID').setValue(loginID)
            }
      }

      // Add career plan.
      addCareerPlanToTalent(): void {
            this.spinnerService.loadingMessage = "Adding Career Plan"


            this.talentService.addCareerPlan(this.careerPlanForm.get('careerPlans').value, this.selectedTalent.id).subscribe((response: any) => {
                  this.showAddCareerPlanForm = false
                  this.getAllCareerPlans()
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

      // On clicking update career plan button.
      OnUpdateCareerPlanButtonClick(careerPlan: any): void {
            this.selectedCareerPlan = careerPlan
            this.showAddCareerPlanForm = false
            this.showUpdateCareerPlanForm = true
            this.currentRatingUpdate = this.selectedCareerPlan.currentRating
      }

      // Update career plan.
      updateCareerPlan(): void {
            if (!this.currentRatingUpdate) {
                  alert("Please enter the rating")
                  return
            }
            if (this.currentRatingUpdate == 0) {
                  alert("Rating cannot be 0")
                  return
            }
            if (this.currentRatingUpdate < 1) {
                  alert("Rating cannot be less than 1")
                  return
            }
            if (this.currentRatingUpdate > 10) {
                  alert("Rating cannot be more than 10")
                  return
            }
            this.selectedCareerPlan.currentRating = this.currentRatingUpdate
            this.selectedCareerPlan.facultyID = this.localService.getJsonValue("loginID")
            this.spinnerService.loadingMessage = "Updating Career Plan Rating"


            this.talentService.updateCareerPlan(this.selectedCareerPlan).subscribe((response: any) => {
                  this.showUpdateCareerPlanForm = false
                  this.getAllCareerPlans()
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

      // Delete career plan.
      deleteCareerPlan(careerObjectiveID: string): void {
            if (confirm("Are you sure you want to delete the career plan?")) {
                  this.spinnerService.loadingMessage = "Deleting Career Plan"


                  this.talentService.deleteCareerPlan(careerObjectiveID, this.selectedTalent.id).subscribe((response: any) => {
                        this.getAllCareerPlans()
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

      // Validate career plan form.
      validateCareerPlanForm(): void {
            if (this.checkCareerObjectiveExists(this.careerPlanForm.get('careerObjective'))) {
                  return
            }
            if (this.careerPlanForm.invalid) {
                  this.careerPlanForm.markAllAsTouched()
                  return
            }
            this.addCareerPlanToTalent()
      }

      // Show the career objective name.
      showCareerObejctiveByID(id: string): string {
            for (let i = 0; i < this.careerObjectiveList.length; i++) {
                  if (this.careerObjectiveList[i].id == id) {
                        return this.careerObjectiveList[i].name
                  }
            }
      }

      // Show the faculty name.
      showFacultyByID(id: string): string {
            for (let i = 0; i < this.facultyList.length; i++) {
                  if (this.facultyList[i].id == id) {
                        return this.facultyList[i].firstName + " " + this.facultyList[i].lastName
                  }
            }
      }

      // Set the fields of all career plans.
      setFieldsOfAllCareerPlans(): void {
            let currentCareerObjectiveCourses: any = []
            let courseID: string = ""
            for (let j = 0; j < this.careerPlanList.length; j++) {
                  for (let i = 0; i < this.careerObjectiveList.length; i++) {
                        if (this.careerObjectiveList[i].id == this.careerPlanList[j].careerObjectiveID) {
                              currentCareerObjectiveCourses = this.careerObjectiveList[i].courses
                              break
                        }
                  }
                  for (let i = 0; i < currentCareerObjectiveCourses.length; i++) {
                        if (currentCareerObjectiveCourses[i].id == this.careerPlanList[j].careerObjectivesCoursesID) {
                              this.careerPlanList[j].technicalAspect = currentCareerObjectiveCourses[i].technicalAspect
                              this.careerPlanList[j].order = currentCareerObjectiveCourses[i].order
                              courseID = currentCareerObjectiveCourses[i].courseID
                        }
                  }
                  for (let i = 0; i < this.courseList.length; i++) {
                        if (this.courseList[i].id == courseID) {
                              this.careerPlanList[j].courseName = this.courseList[i].name
                        }
                  }
            }
      }

      // Make map of career objective id and its related career plans.
      makeCareerPlanMap(careerPlanList: any[]): void {
            let tempCareerPlanArray: any[] = []
            this.careerPlanMap = new Map()
            for (let i = 0; i < this.careerPlanList.length; i++) {
                  if (!this.careerPlanMap.get(careerPlanList[i].careerObjectiveID)) {
                        this.careerPlanMap.set(careerPlanList[i].careerObjectiveID, [])
                  }
                  tempCareerPlanArray = this.careerPlanMap.get(careerPlanList[i].careerObjectiveID)
                  tempCareerPlanArray.push(careerPlanList[i])
                  this.careerPlanMap.set(careerPlanList[i].careerObjectiveID, tempCareerPlanArray)
            }

            //sort career plans by order
            this.careerPlanMap.forEach((value: any[], key: string) => {
                  value.sort((a, b) => {
                        return a.order - b.order
                  })
            })
      }

      // Check if career objective already exists for the talent.
      checkCareerObjectiveExists(careerObjective: any): boolean {
            if (careerObjective && careerObjective.value) {
                  for (let i = 0; i < this.careerPlanList.length; i++) {
                        if (careerObjective.value.id == this.careerPlanList[i].careerObjectiveID) {
                              return true
                        }
                  }
                  return false
            }
            return false
      }

      //*********************************************CRUD FUNCTIONS FOR WAITING LIST************************************************************
      // On clicking get waiting list by talent id.
      getWaitingListForSelectedTalent(talentID: any): void {
            this.isOperationWaitingListUpdate = false
            this.waitingList = []
            this.spinnerService.loadingMessage = "Getting waiting list"
            this.assignSelectedTalent(talentID)
            this.getWaitingList()
            this.showWaitingListForm = false
            this.openModal(this.waitingListFormModal, 'xl')
      }

      // Get waiting list by talent id.
      getWaitingList(): void {
            this.spinnerService.loadingMessage = "Getting Waiting List"


            this.talentService.getWaitingListByTalent(this.selectedTalent.id).subscribe(response => {
                  this.waitingList = response
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // On clicking add new waiting list button.
      onAddNewWaitingListButtonClick(): void {
            this.isOperationWaitingListUpdate = false
            this.activeRequirementByCompanyList = []
            this.activeBatchByCourseList = []
            this.createWaitingListForm()
            this.showWaitingListForm = true
            this.showCompanyBranchList = false
            this.showRequirementList = false
            this.showCourseList = false
            this.showBatchList = false
      }

      // Add waiting list.
      addWaitingList(): void {
            this.spinnerService.loadingMessage = "Adding Waiting List"


            this.waitingListForm.get('email').setValue(this.selectedTalent.email)
            this.waitingListForm.get('isActive').setValue(true)
            let waitingList: IWaitingList = this.waitingListForm.value
            this.patchIDFromObjectsForWaitingList(waitingList)
            waitingList.talentID = this.selectedTalent.id
            this.talentService.addWaitingList(waitingList).subscribe((response: any) => {
                  this.showWaitingListForm = false
                  this.getWaitingList()
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

      // On clicking update waiting list button.
      OnUpdateWaitingListButtonClick(index: number): void {
            this.activeRequirementByCompanyList = []
            this.activeBatchByCourseList = []
            this.showWaitingListForm = true
            this.showCompanyBranchList = false
            this.showRequirementList = false
            this.showCourseList = false
            this.showBatchList = false
            this.isOperationWaitingListUpdate = true
            this.createWaitingListForm()
            this.showWaitingListForm = true
            if (this.waitingList[index].companyBranch) {
                  this.waitingListForm.get('waitingFor').setValue("Company")
                  this.waitingListForm.addControl('companyBranch', this.formBuilder.control(null, [Validators.required]))
                  this.showCompanyBranchList = true
                  if (this.waitingList[index].companyRequirement) {
                        this.getActiveRequirementByCompanyList(this.waitingList[index].companyBranch.id, index)
                  }
                  else {
                        this.waitingListForm.patchValue(this.waitingList[index])
                  }
                  return
            }
            if (this.waitingList[index].course) {
                  this.waitingListForm.get('waitingFor').setValue("Course")
                  this.waitingListForm.addControl('course', this.formBuilder.control(null, [Validators.required]))
                  this.showCourseList = true
                  if (this.waitingList[index].batch) {
                        this.getActiveBatchByCourseList(this.waitingList[index].course.id, index)
                  }
                  else {
                        this.waitingListForm.patchValue(this.waitingList[index])
                  }
            }
      }

      // Update waiting list.
      updateWaitingList(): void {
            this.spinnerService.loadingMessage = "Updating Waiting List"


            let waitingList: IWaitingList = this.waitingListForm.value
            this.patchIDFromObjectsForWaitingList(waitingList)
            this.talentService.updateWaitingList(waitingList).subscribe((response: any) => {
                  this.showWaitingListForm = false
                  this.getWaitingList()
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

      // Delete waiting list.
      deleteWaitingList(waitingListID: string): void {
            if (confirm("Are you sure you want to delete the waiting list?")) {
                  this.spinnerService.loadingMessage = "Deleting Waiting List"


                  this.talentService.deleteWaitingList(waitingListID).subscribe((response: any) => {
                        this.getWaitingList()
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

      // Validate waiting list form.
      validateWaitingListForm(): void {
            if (this.waitingListForm.invalid) {
                  this.waitingListForm.markAllAsTouched()
                  return
            }
            if (this.isOperationWaitingListUpdate) {
                  this.updateWaitingList()
                  return
            }
            this.addWaitingList()
      }

      // Extract ID from objects in waiting list form.
      patchIDFromObjectsForWaitingList(waitingList: IWaitingList): void {
            if (this.waitingListForm.contains('companyBranch') && this.waitingListForm.get('companyBranch').value) {
                  waitingList.companyBranchID = this.waitingListForm.get('companyBranch').value.id
                  delete waitingList['companyBranch']
            }
            if (this.waitingListForm.contains('companyRequirement') && this.waitingListForm.get('companyRequirement').value) {
                  waitingList.companyRequirementID = this.waitingListForm.get('companyRequirement').value.id
                  delete waitingList['companyRequirement']
            }
            if (this.waitingListForm.contains('course') && this.waitingListForm.get('course').value) {
                  waitingList.courseID = this.waitingListForm.get('course').value.id
                  delete waitingList['course']
            }
            if (this.waitingListForm.contains('batch') && this.waitingListForm.get('batch').value) {
                  waitingList.batchID = this.waitingListForm.get('batch').value.id
                  delete waitingList['batch']
            }
            if (this.waitingListForm.get('waitingFor').value) {
                  delete waitingList['waitingFor']
            }
            for (let i = 0; i < this.sourceList.length; i++) {
                  if (this.sourceList[i].name == "tt") {
                        waitingList.sourceID = this.sourceList[i].id
                  }
            }
      }

      // On changing value of waiting for control in waiting list form.
      onWaitingForChange(waitingFor: string): void {
            if (waitingFor == "Company") {
                  this.onCompanySelectionForWaitingList()
            }
            if (waitingFor == "Course") {
                  this.onCourseSelectionForWaitingList()
            }
      }

      // On selecting 'company' for 'waiting for' control in waiting list form.
      onCompanySelectionForWaitingList(): void {
            // Remove course realted controls and initialize course related variables.
            this.showCourseList = false
            this.showBatchList = false
            this.waitingListForm.removeControl('course')
            this.waitingListForm.removeControl('batch')
            this.activeBatchByCourseList = []

            // Add compnay branch related controls.
            this.waitingListForm.addControl('companyBranch', this.formBuilder.control(null, [Validators.required]))
            this.showCompanyBranchList = true
      }

      // On selecting 'course' for 'waiting for' control in waiting list form.
      onCourseSelectionForWaitingList(): void {
            // Remove company branch realted controls and initialize company branch related variables.
            this.showCompanyBranchList = false
            this.showRequirementList = false
            this.waitingListForm.removeControl('companyBranch')
            this.waitingListForm.removeControl('companyRequirement')
            this.activeRequirementByCompanyList = []

            // Add course related controls.
            this.waitingListForm.addControl('course', this.formBuilder.control(null, [Validators.required]))
            this.showCourseList = true
      }

      // On changing value of company branch control in waiting list form. 
      onCompanyBranchChangeForWaitingList(companyBranch: any): void {
            if (!companyBranch) {
                  this.showRequirementList = false
                  this.waitingListForm.removeControl('companyRequirement')
                  return
            }
            this.getActiveRequirementByCompanyList(companyBranch.id)
      }

      // On changing value of course control in waiting list form. 
      onCourseChangeForWaitingList(course: any): void {
            if (!course) {
                  this.showBatchList = false
                  this.waitingListForm.removeControl('batch')
                  return
            }
            this.getActiveBatchByCourseList(course.id)
      }

      // Get talents for waiting list.
      getTalentsForWaitingList(): void {
            let queryParams: any = {}
            if (this.waitingListCompanyBranchID) {
                  queryParams.companyBranchID = this.waitingListCompanyBranchID
            }
            if (this.waitingListRequirementID) {
                  queryParams.requirementID = this.waitingListRequirementID
            }
            if (this.waitingListCourseID) {
                  queryParams.courseID = this.waitingListCourseID
            }
            if (this.waitingListBatchID) {
                  queryParams.batchID = this.waitingListBatchID
            }
            if (this.waitingListTechnologyID) {
                  queryParams.technologyID = this.waitingListTechnologyID
            }
            this.spinnerService.loadingMessage = "Getting Talents"


            this.talentService.getTalentsByWaitingList(this.limit, this.offset, queryParams).
                  subscribe((response: any) => {
                        this.talents = response.body
                        this.totalTalents = response.headers.get('X-Total-Count')
                        this.totalLifetimeValue = parseInt(response.headers.get("totalLifetimeValue"))
                        this.calculateFieldsOfTalent(this.talents)
                  },
                        (error) => {
                              console.error(error)
                              if (error.error) {
                                    alert(error.error)
                                    return
                              }
                              alert(error.error.error)
                        }).add(() => {


                              this.setPaginationString()
                        })
      }

      //*********************************************CAMPUS DRIVE LIST TALENTS***********************************************************
      // Get talents by campus drive id.
      getTalentsByCampusdDrive(): void {
            this.spinnerService.loadingMessage = "Getting Talents"


            let queryParams: any = {
                  "hasAppeared": this.hasAppeared
            }
            this.talentService.getTalentsByCampusDrive(this.campusDriveID, this.limit, this.offset, queryParams).
                  subscribe((response: any) => {
                        this.talents = response.body
                        this.totalTalents = response.headers.get('X-Total-Count')
                        this.totalLifetimeValue = parseInt(response.headers.get("totalLifetimeValue"))
                        this.calculateFieldsOfTalent(this.talents)
                  },
                        (error) => {
                              console.error(error)
                              if (error.error) {
                                    alert(error.error)
                                    return
                              }
                              alert(error.error.error)
                        }).add(() => {


                              this.setPaginationString()
                        })
      }

      // Get talents for excel download.
      getExcelDownloadTalents(limit: number): void {
            this.spinnerService.loadingMessage = "Downloading Talents"


            this.talentService.getExcelDownloadTalents(limit, this.excelDownloadOffsetCount).
                  subscribe((response: any) => {
                        this.excelDownloadTalents = response.body
                        this.formatExcelDownloadTalents(this.excelDownloadTalents)
                        this.excelDownloadOffsetCount = this.excelDownloadOffsetCount + 1
                        if (this.excelDownloadOffsetCount <= (this.excelDownloadOffsetTotal + 1)) {
                              if (this.excelDownloadOffsetCount == 1) {
                                    this.fileOps.createExcelWorkSheet(this.excelDownloadTalents)
                                    this.getExcelDownloadTalents(limit)
                                    return
                              }
                              if (this.excelDownloadOffsetCount > 1) {
                                    this.fileOps.appendToExcelWorkSheet(this.excelDownloadTalents)
                                    this.getExcelDownloadTalents(limit)
                                    return
                              }
                        }
                        if (this.excelDownloadOffsetCount > (this.excelDownloadOffsetTotal + 1)) {
                              this.downloadExcelFile()
                              return
                        }
                  },
                        (error) => {
                              console.error(error)
                              this.getExcelDownloadTalents(limit)
                              // if (error.error) {
                              //       alert(error.error)
                              //       return
                              // }
                              alert("Could not download excel.. try again later")
                        }).add(() => {


                              this.setPaginationString()
                        })
      }

      // Get searched talents for excel download.
      getExcelDownloadSearchedTalents(limit: number): void {
            let data = this.talentSearchForm.value
            this.utilityService.deleteNullValuePropertyFromObject(data)
            this.spinnerService.loadingMessage = "Downloading Talents"


            this.talentService.getSearchedExcelDownloadTalents(data, limit, this.excelDownloadOffsetCount).
                  subscribe((response: any) => {
                        this.excelDownloadTalents = response.body
                        this.formatExcelDownloadTalents(this.excelDownloadTalents)
                        this.excelDownloadOffsetCount = this.excelDownloadOffsetCount + 1
                        if (this.excelDownloadOffsetCount <= (this.excelDownloadOffsetTotal + 1)) {
                              if (this.excelDownloadOffsetCount == 1) {
                                    this.fileOps.createExcelWorkSheet(this.excelDownloadTalents)
                                    this.getExcelDownloadSearchedTalents(limit)
                                    return
                              }
                              if (this.excelDownloadOffsetCount > 1) {
                                    this.fileOps.appendToExcelWorkSheet(this.excelDownloadTalents)
                                    this.getExcelDownloadSearchedTalents(limit)
                                    return
                              }
                        }
                        if (this.excelDownloadOffsetCount > (this.excelDownloadOffsetTotal + 1)) {
                              this.downloadExcelFile()
                              return
                        }
                  },
                        (error) => {
                              console.error(error)
                              this.getExcelDownloadSearchedTalents(limit)
                              // if (error.error) {
                              //       alert(error.error)
                              //       return
                              // }
                              alert("Could not download excel.. try again later")
                        }).add(() => {


                              this.setPaginationString()
                        })
      }

      //*********************************************SEMINAR LIST TALENTS***********************************************************
      // Get talents by seminar id.
      getTalentsBySeminar(): void {
            this.spinnerService.loadingMessage = "Getting Talents"


            let queryParams: any = {
                  "hasVisited": this.hasVisited
            }
            this.talentService.getTalentsBySeminar(this.seminarID, this.limit, this.offset, queryParams).
                  subscribe((response: any) => {
                        this.talents = response.body
                        this.totalTalents = response.headers.get('X-Total-Count')
                        this.totalLifetimeValue = parseInt(response.headers.get("totalLifetimeValue"))
                        this.calculateFieldsOfTalent(this.talents)
                  },
                        (error) => {
                              console.error(error)
                              if (error.error) {
                                    alert(error.error)
                                    return
                              }
                              alert(error.error.error)
                        }).add(() => {


                              this.setPaginationString()
                        })
      }

      //*********************************************PROFESSIONAL SUMMARY REPORT TALENTS***********************************************************
      // Get talents for professional summary report.
      getTalentsForProSummaryReport(): void {
            this.spinnerService.loadingMessage = "Getting Talents"


            let queryParams: any = {
                  "companyName": this.companyName,
                  "category": this.category,
                  "isCompany": this.isCompany
            }
            this.talentService.getTalentsForProSummaryReport(this.limit, this.offset, queryParams).
                  subscribe((response: any) => {
                        this.talents = response.body
                        this.totalTalents = response.headers.get('X-Total-Count')
                        this.totalLifetimeValue = parseInt(response.headers.get("totalLifetimeValue"))
                        this.calculateFieldsOfTalent(this.talents)
                  },
                        (error) => {
                              console.error(error)
                              if (error.error) {
                                    alert(error.error)
                                    return
                              }
                              alert(error.error.error)
                        }).add(() => {


                              this.setPaginationString()
                        })
      }

      // Get talents for professional summary report by technology count.
      getTalentsForProSummaryReportByTechnologyCount(): void {
            this.spinnerService.loadingMessage = "Getting Talents"


            let queryParams: any = {
                  "companyName": this.companyNameTechCount,
                  "category": this.categoryTechCount,
                  "isCompany": this.isCompanyTechCount,
                  "technologyID": this.technologyIDTechCount
            }
            this.talentService.getTalentsForProSummaryReportTechnologyCount(this.limit, this.offset, queryParams).
                  subscribe((response: any) => {
                        this.talents = response.body
                        this.totalTalents = response.headers.get('X-Total-Count')
                        this.totalLifetimeValue = parseInt(response.headers.get("totalLifetimeValue"))
                        this.calculateFieldsOfTalent(this.talents)
                  },
                        (error) => {
                              console.error(error)
                              if (error.error) {
                                    alert(error.error)
                                    return
                              }
                              alert(error.error.error)
                        }).add(() => {


                              this.setPaginationString()
                        })
      }

      //*********************************************FRESHER SUMMARY REPORT TALENTS***********************************************************
      // Get talents for fresher summary report.
      getTalentsForFresherSummaryReport(): void {
            this.spinnerService.loadingMessage = "Getting Talents"


            let queryParams: any = {}

            if (this.academicYear) {
                  queryParams.academicYear = this.academicYear
            }

            if (this.isExperienced) {
                  queryParams.isExperienced = this.isExperienced
            }

            if (this.talentType) {
                  queryParams.talentType = this.talentType
            }

            if (this.fresherTechnology) {
                  queryParams.fresherTechnology = this.fresherTechnology
            }
            if (this.isLookingForJob) {
                  queryParams.isLookingForJob = this.isLookingForJob
            }
            this.talentService.getTalentsForFresherSummaryReport(this.limit, this.offset, queryParams).
                  subscribe((response: any) => {
                        this.talents = response.body
                        this.totalTalents = response.headers.get('X-Total-Count')
                        this.totalLifetimeValue = parseInt(response.headers.get("totalLifetimeValue"))
                        this.calculateFieldsOfTalent(this.talents)
                  },
                        (error) => {
                              console.error(error)
                              if (error.error) {
                                    alert(error.error)
                                    return
                              }
                              alert(error.error.error)
                        }).add(() => {


                              this.setPaginationString()
                        })
      }

      //*********************************************PACKAGE SUMMARY REPORT TALENTS***********************************************************
      // Get talents for package summary report.
      getTalentsForPackageSummaryReport(): void {
            this.spinnerService.loadingMessage = "Getting Talents"


            let queryParams: any = {}

            if (this.packageType) {
                  queryParams.packageType = this.packageType
            }

            if (this.packageTechnology) {
                  queryParams.packageTechnology = this.packageTechnology
            }

            if (this.packageExperience) {
                  queryParams.packageExperience = this.packageExperience
            }

            this.talentService.getTalentsForPackageSummaryReport(this.limit, this.offset, queryParams).
                  subscribe((response: any) => {
                        this.talents = response.body
                        this.totalTalents = response.headers.get('X-Total-Count')
                        this.totalLifetimeValue = parseInt(response.headers.get("totalLifetimeValue"))
                        this.calculateFieldsOfTalent(this.talents)
                  },
                        (error) => {
                              console.error(error)
                              if (error.error) {
                                    alert(error.error)
                                    return
                              }
                              alert(error.error.error)
                        }).add(() => {


                              this.setPaginationString()
                        })
      }

      //*********************************************GET FUNCTIONS************************************************************
      // Get all lists.
      getAllComponents(): void {
            this.getTalentType()
            this.getCountryList()
            this.getTechnologyList()
            this.getDesignationList()
            this.getPersonalityTypes()
            this.getDegreeList()
            this.getSalesPersonList()
            this.getBatchList()
            this.getAcademicYear()
            this.getPurposes()
            this.getSourceList()
            this.getCollegeBranchList()
            this.getExaminationList()
            this.getFacultyList()
            this.getCourseList()
            this.getNextActionTypeList()
            this.getCompanyList()
            this.getCareerObjectiveList()
            this.getYearOfMSList()
            this.getCompanyBranchList()
            this.getAllBatchList()
            this.getAllRequirementList()
            this.getSupervisorCount()
            this.getQueryParams()
      }

      // Get batch list.
      getBatchList(): void {
            let queryParams: any = {
                  is_Active: '1'
            }
            this.batchService.getBatchList(queryParams).subscribe(data => {
                  this.batchList = data.body
            }), err => {
                  console.error(err)
            }
      }

      // Get faculty list.
      getFacultyList(): void {
            this.generalService.getFacultyList().subscribe((response: any) => {
                  this.facultyList = response.body
            }, err => {
                  console.error(err)
            })
      }

      // Get course list.
      getCourseList(): void {
            this.generalService.getCourseList().subscribe((response: any) => {
                  this.courseList = response.body
            }, err => {
                  console.error(err)
            })
      }

      // Get company list.
      getCompanyList(): void {
            this.generalService.getCompanyBranchList().subscribe((response: any) => {
                  this.companyList = response.body
            }, (err: any) => {
                  console.error(err)
            })
      }

      // Get company branch list.
      getCompanyBranchList(): void {
            this.generalService.getCompanyBranchList().subscribe((response: any) => {
                  this.companyBranchList = response.body
            }, (err: any) => {
                  console.error(err)
            })
      }

      //Get salesperson list.
      getSalesPersonList(): void {
            this.generalService.getSalesPersonList().subscribe(
                  data => {
                        this.salesPersonList = data.body
                  }
            ), err => {
                  console.error(err)
            }
      }

      // Get academic year list.
      getAcademicYear(): void {
            this.generalService.getGeneralTypeByType("academic_year").subscribe((respond: any[]) => {
                  this.academicYearList = respond
            }, (err) => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get purpose list.
      getPurposes(): void {
            this.generalService.getPurposeListByType("talent").subscribe(
                  data => {
                        this.purposeList = data
                  }
            ), err => {
                  console.error(err)
            }
      }

      // Get outcome list by purpose for call record form.
      getOutcomesByPurpose(purpose: IPurpose): void {
            this.togglePurposeJobFields(purpose)
            this.removetargetDateFromCallRecordForm()
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

      // Get outcome list by purpose id for talent search form.
      getOutcomesByPurposeID(purposeID: any): void {
            if (purposeID == null) {
                  this.outcomeList = []
                  return
            }
            this.generalService.getOutcomeListByPurpose(purposeID).subscribe((data: any) => {
                  this.outcomeList = data
            }, (err) => {
                  console.error(err)
            })
      }

      // Get technology list.
      getTechnologyList(event?: any): void {
            // this.generalService.getTechnologies().subscribe((respond: any[]) => {
            //       this.technologyList = respond
            // }, (err) => {
            //       console.error(this.utilityService.getErrorString(err))
            // })
            let queryParams: any = {}
            if (event && event?.term != "") {
                  queryParams.language = event.term
            }
            this.isTechLoading = true
            this.techService.getAllTechnologies(this.techLimit, this.techOffset, queryParams).subscribe((response) => {
                  this.technologyList = []
                  this.technologyList = this.technologyList.concat(response.body)
            }, (err) => {
                  console.error(err)
            }).add(() => {
                  this.isTechLoading = false
            })
      }

      // Get all talents by limit and offset.
      getTalents(): void {
            let roleNameAndLogin: any = {
                  roleName: this.roleName,
                  loginID: this.localService.getJsonValue("loginID"),
                  isViewAllBatches: this.isViewAllTalents?1:0
            }
            console.log(roleNameAndLogin)
            if (this.isNavigatedFromBatch) {
                  this.getSelectedTalents()
                  return
            }
            if (this.isNavigatedFromWaitingList) {
                  this.getTalentsForWaitingList()
                  return
            }
            if (this.isNavigatedFromCampusDrive) {
                  this.getTalentsByCampusdDrive()
                  return
            }
            if (this.isNavigatedFromSeminar) {
                  this.getTalentsBySeminar()
                  return
            }
            if (this.showRequirementTalents) {
                  this.getAllTalentsInRequirement()
                  return
            }
            if (this.isSearched || this.requirementSearched) {
                  this.getAllSearchedTalents()
                  return
            }
            if (this.isNavigatedFromProSummaryReport) {
                  this.getTalentsForProSummaryReport()
                  return
            }
            if (this.isNavigatedFromProSummaryReportTechCount) {
                  this.getTalentsForProSummaryReportByTechnologyCount()
                  return
            }
            if (this.isNavigatedFromFresherSummaryReport) {
                  this.getTalentsForFresherSummaryReport()
                  return
            }

            if (this.isNavigatedFromPackageSummaryReport) {
                  this.getTalentsForPackageSummaryReport()
                  return
            }

            this.spinnerService.loadingMessage = "Getting All talents"


            this.talentService.getTalents(this.limit, this.offset, roleNameAndLogin).subscribe(res => {
                  this.talents = res.body
                  this.calculateFieldsOfTalent(this.talents)
                  this.totalTalents = parseInt(res.headers.get("X-Total-Count"))
                  this.totalLifetimeValue = parseInt(res.headers.get("totalLifetimeValue"))
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            }).add(() => {

                  this.setPaginationString()
            })
      }

      // Get next action type list.
      getNextActionTypeList(): void {
            this.generalService.getNextActionTypeList().subscribe((response: any) => {
                  this.nextActionTypeList = response.body
            }, (err: any) => {
                  console.error(err)
            })
      }

      // Get batch list of one talent.
      getBatchListOfOneTalent(talentID: string): void {
            this.batchesFound = true
            this.showAdditionalBatchDetails = true
            this.batchesOfOneTalentList = []
            this.spinnerService.loadingMessage = "Getting Batches"
            this.assignSelectedTalent(talentID)
            this.talentService.getBatchListOfOneTalent(this.selectedTalent.id).subscribe(response => {
                  this.batchesOfOneTalentList = response
                  if (this.batchesOfOneTalentList.length == 0) {
                        this.batchesFound = false
                  }
                  for (let index = 0; index < this.batchesOfOneTalentList.length; index++) {
                        this.batchesOfOneTalentList[index].isVisible = false
                  }
            }, (err: any) => {
                  this.batchesFound = false
                  console.error(err)
                  if (err.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                        return
                  }
                  alert(err.error.error)
            })
            this.openModal(this.showBatchesModal, 'xl')
      }

      // Get session list for one batch.
      getSessionsForBatch(batch: any): void {
            if (batch.isVisible) {
                  batch.isVisible = false
                  this.showAdditionalBatchDetails = true
                  return
            }
            let queryParams: any = {
                  roleName: this.roleName,
                  loginID: this.localService.getJsonValue("loginID")
            }

            // this.showAdditionalBatchDetails = false
            this.spinnerService.loadingMessage = "Getting sessions"


            this.batchSessions = []
            this.batchService.getSessionForBatch(batch.id, queryParams).subscribe((response: any) => {
                  this.batchSessions = response.body
                  batch.isVisible = true
                  if (this.batchSessions.length > 0) {
                        this.showAdditionalBatchDetails = false
                  }
                  if (this.batchSessions.length == 0) {
                        this.showAdditionalBatchDetails = true
                        alert("No sessions found")
                  }
            }, (err: any) => {
                  console.error(err)
                  if (err.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                        return
                  }
                  alert(err.error.error)
            })
      }

      // Get talent list by company requirement.
      getAllTalentsInRequirement(): void {
            this.spinnerService.loadingMessage = "Getting All Talents"


            this.talentService.getAllTalentsByRequirement(this.requirementID, this.limit, this.offset).subscribe(res => {
                  this.talents = res.body
                  this.totalTalents = parseInt(res.headers.get("X-Total-Count"))
            }, err => {
                  console.error(err)
            }).add(() => {

                  this.setPaginationString()
            })
      }

      // Get talent type list.
      getTalentType(): void {
            this.generalService.getGeneralTypeByType("talent_rating").subscribe((respond: any[]) => {
                  this.talentTypeList = respond
            }, (err) => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get personality type list.
      getPersonalityTypes(): void {
            this.generalService.getGeneralTypeByType("personality_type").subscribe((respond: any[]) => {
                  this.personalityTypeList = respond
            }, (err) => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get state list by country.
      getStateListByCountry(country: any): void {
            if (country == null) {
                  this.stateList = []
                  this.talentForm.get('state').setValue(null)
                  this.talentForm.get('state').disable()
                  return
            }
            if (country.name == "U.S.A." || country.name == "India") {
                  this.talentForm.get('state').enable()
                  this.generalService.getStatesByCountryID(country.id).subscribe((respond: any[]) => {
                        this.stateList = respond
                  }, (err) => {
                        console.error(err)
                  })
            }
            else {
                  this.stateList = []
                  this.talentForm.get('state').setValue(null)
                  this.talentForm.get('state').disable()
            }
      }

      // Get specialization list.
      getSpecializationListByDegreeID(degreeID: string, index?: number, shouldClear?: boolean): void {
            if (index != undefined) {
                  this.talentAcademicControlArray.at(index).get('specialization').setValue(null)
            }
            if (this.degreeIDSet.has(degreeID)) {
                  return
            }
            this.degreeIDSet.add(degreeID)
            if (degreeID) {
                  this.generalService.getSpecializationByDegreeID(degreeID).subscribe((response: any) => {
                        this.specializationList = this.specializationList.concat(response.body)
                  }, (err) => {
                        console.error(err)
                  })
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

      // Get degree list.
      getDegreeList(): void {
            let queryParams: any = {
                  limit: -1,
                  offset: 0,
            }
            this.degreeService.getAllDegrees(queryParams).subscribe((respond: any) => {
                  //respond.push(this.other)
                  this.degreeList = respond.body
            }, (err) => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get designation list.
      getDesignationList(): void {
            this.generalService.getDesignations().subscribe((respond: any[]) => {
                  //respond.push(this.other)
                  this.designationList = respond
            }, (err) => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get source list.
      getSourceList(): void {
            this.generalService.getSources().subscribe(response => {
                  this.sourceList = response
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get college branch list.
      getCollegeBranchList(event?: any): void {
            // this.generalService.getCollegeBranchList().subscribe(response => {
            //       this.collegeBranchList = response
            // }, err => {
            //       console.error(this.utilityService.getErrorString(err))
            // })

            let queryParams: any = {}
            if (event && event?.term != "") {
                  queryParams.branchName = event.term
            }
            this.isCollegeLoading = true
            this.generalService.getCollegeBranchListWithLimit(this.collegeLimit, this.collegeOffset, queryParams).subscribe((response) => {
                  this.collegeBranchList = []
                  this.collegeBranchList = this.collegeBranchList.concat(response)
            }, (err) => {
                  console.error(err)
            }).add(() => {
                  this.isCollegeLoading = false
            })
      }

      // Get career objective list.
      getCareerObjectiveList(): void {
            let queryParamas: any = {
                  limit: -1,
                  offset: 0
            }
            this.careerObjectiveService.getAllCareerObjectives(queryParamas).subscribe(response => {
                  this.careerObjectiveList = response.body
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get examination list.
      getExaminationList(): void {
            this.generalService.getExaminationList().subscribe(response => {
                  this.examinationList = response
                  this.setExaminationsTotalMarks()
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get all batch list.
      getAllBatchList(): void {
            this.generalService.getBatchList().subscribe(response => {
                  this.allBatchList = response
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get all requirement list.
      getAllRequirementList(): void {
            this.generalService.getRequirementList().subscribe(response => {
                  this.allRequirementList = response.body
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get count of faculty supervisors.
      getSupervisorCount(): void {
            let queryParams: any = {
                  roleName: this.roleName,
            }
            this.adminService.getSupervisorCount(queryParams).subscribe(response => {
                  this.supervisorCount = response.body.totalCount
                  console.log(this.supervisorCount)

                  // For faculty.
                  if (this.roleName == this.role.FACULTY && this.supervisorCount == 0){
                        this.isHeadFcaulty = true
                  }
                  if (this.roleName == this.role.FACULTY && this.supervisorCount > 0){
                        this.isHeadFcaulty = false
                  }

                  // For salesperson.
                  if (this.roleName == this.role.SALES_PERSON && this.supervisorCount == 0){
                        this.isHeadSalesPerson = true
                  }
                  if (this.roleName == this.role.SALES_PERSON && this.supervisorCount > 0){
                        this.isHeadSalesPerson = false
                  }
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Make list of years for masters abroad.
      getYearOfMSList(): void {
            this.yearOfMSList.push(this.currentYear)
            this.yearOfMSList.push((this.currentYear + 1))
            this.yearOfMSList.push((this.currentYear + 2))
      }

      // Get active company requirement list by company branch id.
      getActiveRequirementByCompanyList(companyBranchID: string, index?: number): void {
            let queryParams: any = {
                  companyBranchID: companyBranchID,
                  isActive: "1"
            }
            if (this.waitingListForm.contains('companyRequirement')) {
                  this.waitingListForm.get('companyRequirement').setValue(null)
                  this.activeRequirementByCompanyList = []
            }
            this.generalService.getRequirementList(queryParams).subscribe(response => {
                  this.activeRequirementByCompanyList = response.body
                  if (this.activeRequirementByCompanyList?.length == 0) {
                        this.showRequirementList = false
                        this.waitingListForm.removeControl('companyRequirement')
                  }
                  if (this.activeRequirementByCompanyList?.length != 0 && !this.waitingListForm.contains('companyRequirement')) {
                        this.waitingListForm.addControl('companyRequirement', this.formBuilder.control(null))
                        this.showRequirementList = true
                  }
                  if (index != undefined) {
                        this.waitingListForm.patchValue(this.waitingList[index])
                  }
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get upcoming or live batch list by course id.
      getActiveBatchByCourseList(courseID: string, index?: number): void {
            let queryParams: any = {
                  courseID: courseID,
                  batchStatus: ["Ongoing", "Upcoming"],
                  isActive: "1"
            }
            if (this.waitingListForm.contains('batch')) {
                  this.waitingListForm.get('batch').setValue(null)
                  this.activeBatchByCourseList = []
            }
            this.batchService.getBatchList(queryParams).subscribe(response => {
                  this.activeBatchByCourseList = response.body
                  if (this.activeBatchByCourseList?.length == 0) {
                        this.showBatchList = false
                        this.waitingListForm.removeControl('batch')
                  }
                  if (this.activeBatchByCourseList?.length != 0 && !this.waitingListForm.contains('batch')) {
                        this.waitingListForm.addControl('batch', this.formBuilder.control(null, [Validators.required]))
                        this.showBatchList = true
                  }
                  if (index != undefined) {
                        this.waitingListForm.patchValue(this.waitingList[index])
                  }
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Set examination's total marks.
      setExaminationsTotalMarks(): void {
            for (let i = 0; i < this.examinationList.length; i++) {
                  if (this.examinationList[i].name == "GRE") {
                        this.greExam = this.examinationList[i]
                  }
                  if (this.examinationList[i].name == "GMAT") {
                        this.gmatExam = this.examinationList[i]
                  }
                  if (this.examinationList[i].name == "TOEFL") {
                        this.toeflExam = this.examinationList[i]
                  }
                  if (this.examinationList[i].name == "IELTS") {
                        this.ieltsExam = this.examinationList[i]
                  }
            }
      }

      //*********************************************EXCEL UPLOAD************************************************************

      // On uploading excel file.
      onExcelUpload(event: any, excelFile: any) {
            this.resumeDocStatus = ""
            this.addMultipleErrorList = []
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

            this.fileOps.uploadExcel(file).subscribe((uploadedTalents: any) => {
                  if (this.validateTalents(uploadedTalents)) {
                        this.excelUploadedTalents = uploadedTalents
                        this.isExcelUploaded = true
                        this.resumeDocStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
                        this.reformTalentsFromExcel()
                  } else {
                        excelFile.value = ""
                        this.resumeDocStatus = ""
                        this.isExcelUploaded = false
                  }
            }, (err) => {
                  alert(err)
                  excelFile.value = ""
            }).add(() => {

            })
      }

      // On downloading excel file.
      onExcelDownload(limit: number) {
            if (confirm("Are you sure you want to download excel of all talents?")) {
                  this.excelDownloadOffsetTotal = (Math.floor(this.totalTalents / limit))
                  this.excelDownloadOffsetCount = 0
                  this.fileOps.createExcelWorkBook()
                  if (this.isSearched) {
                        this.getExcelDownloadSearchedTalents(limit)
                        return
                  }
                  this.getExcelDownloadTalents(limit)
            }
      }

      // Format the excel download talents.
      formatExcelDownloadTalents(excelDownloadTalentList: ITalentDownloadExcel[]): any[] {
            for (let i = 0; i < excelDownloadTalentList.length; i++) {

                  // Clean all data.
                  excelDownloadTalentList[i]["academicYear"] = this.utilityService.getValueByKey(excelDownloadTalentList[i].academicYear, this.academicYearList)
                  excelDownloadTalentList['academicYear']
                  excelDownloadTalentList[i]["talentType"] = this.utilityService.getValueByKey(excelDownloadTalentList[i].talentType, this.talentTypeList)
                  excelDownloadTalentList['talentType']
                  if (excelDownloadTalentList[i]["fromYear"] != null) {
                        excelDownloadTalentList[i]["fromYear"] = this.datePipe.transform(excelDownloadTalentList[i]["fromYear"], 'd MMM, y')
                  }
                  if (excelDownloadTalentList[i]["toYear"] != null) {
                        excelDownloadTalentList[i]["toYear"] = this.datePipe.transform(excelDownloadTalentList[i]["toYear"], 'd MMM, y')
                  }
                  if (excelDownloadTalentList[i]["toYear"] == null) {
                        excelDownloadTalentList[i]["toYear"] = "Working"
                  }
                  if (excelDownloadTalentList[i]["isActive"] == true) {
                        excelDownloadTalentList[i]["isActive"] = "Yes"
                  }
                  if (excelDownloadTalentList[i]["isActive"] == false) {
                        excelDownloadTalentList[i]["isActive"] = "No"
                  }
                  if (excelDownloadTalentList[i]["isSwabhavTalent"] == true) {
                        excelDownloadTalentList[i]["isSwabhavTalent"] = "Yes"
                  }
                  if (excelDownloadTalentList[i]["isSwabhavTalent"] == false) {
                        excelDownloadTalentList[i]["isSwabhavTalent"] = "No"
                  }
                  if (excelDownloadTalentList[i]["isSwabhavTalent"] == false) {
                        excelDownloadTalentList[i]["isSwabhavTalent"] = "No"
                  }
                  if (excelDownloadTalentList[i]["totalYearOfExp"] == null) {
                        excelDownloadTalentList[i]["totalYearOfExp"] = 0
                  }
                  excelDownloadTalentList[i]["totalYearOfExp"] = Math.floor(excelDownloadTalentList[i]["totalYearOfExp"] / 12)

                  // Change the key names of all fields.
                  for (var key in excelDownloadTalentList[i]) {
                        if (excelDownloadTalentList[i].hasOwnProperty(key)) {
                              let tempKey: string = key
                              tempKey = tempKey.replace(/([A-Z])/g, " $1")
                              let convertedTempKey = tempKey.charAt(0).toUpperCase() + tempKey.slice(1)
                              excelDownloadTalentList[i][convertedTempKey] = excelDownloadTalentList[i][key]
                              delete excelDownloadTalentList[i][key]
                        }
                  }
            }
            return excelDownloadTalentList
      }

      // Download excel file.
      downloadExcelFile(): void {
            this.fileOps.saveAsExcelFile("talent-excel-download")
            alert("Excel downloaded successfully")
      }

      // Validate all talents from excel file upload.
      validateTalents(uploadedTalents: ITalentExcel[]): boolean {

            if (!uploadedTalents || uploadedTalents.length == 0) {
                  alert("No Uploaded Talents")
                  return false
            }

            for (let index = 0; index < uploadedTalents.length; index++) {
                  if (!uploadedTalents[index].firstName || uploadedTalents[index].firstName == "") {
                        alert(`firstname on row ${index + 2} is not specified`)
                        return false
                  }
                  if (!uploadedTalents[index].lastName || uploadedTalents[index].lastName == "") {
                        alert(`lastname on row ${index + 2} is not specified`)
                        return false
                  }
                  if (!uploadedTalents[index].email || uploadedTalents[index].email == "") {
                        alert(`email on row ${index + 2} is not specified`)
                        return false
                  }
                  if (!uploadedTalents[index].contact || uploadedTalents[index].contact == "") {
                        alert(`contact on row ${index + 2} is not specified`)
                        return false
                  }
                  if (!uploadedTalents[index].academicYear || uploadedTalents[index].academicYear == 0) {
                        alert(`academic year on row ${index + 2} is not specified`)
                        return false
                  }
                  if (uploadedTalents[index].isSwabhavTalent == null) {
                        alert(`isSwabhavTalent on row ${index + 2} is not specified`)
                        return false
                  }

                  // Check for academic details.
                  if (uploadedTalents[index].degreeName || uploadedTalents[index].specializationName ||
                        uploadedTalents[index].collegeName || uploadedTalents[index].percentage ||
                        uploadedTalents[index].yearOfPassout) {

                        if (!uploadedTalents[index].degreeName || uploadedTalents[index].degreeName == "") {
                              alert(`degree  on row ${index + 2} is not specified`)
                              return false
                        }

                        if (!uploadedTalents[index].specializationName || uploadedTalents[index].specializationName == "") {
                              alert(`specialization  on row ${index + 2} is not specified`)
                              return false
                        }

                        if (!uploadedTalents[index].collegeName || uploadedTalents[index].collegeName == "") {
                              alert(`college  on row ${index + 2} is not specified`)
                              return false
                        }

                        if (!uploadedTalents[index].percentage || uploadedTalents[index].percentage == 0) {
                              alert(`percentage  on row ${index + 2} is not specified`)
                              return false
                        }

                        if (!uploadedTalents[index].yearOfPassout || uploadedTalents[index].yearOfPassout == 0) {
                              alert(`passout  on row ${index + 2} is not specified`)
                              return false
                        }
                  }
            }
            return true
      }

      // Create talent json to be sent to API.
      reformTalentsFromExcel(): void {
            this.reformedTalents = []
            for (let index = 0; index < this.excelUploadedTalents.length; index++) {

                  let talent: any = {
                        firstName: this.excelUploadedTalents[index].firstName,
                        lastName: this.excelUploadedTalents[index].lastName,
                        email: this.excelUploadedTalents[index].email,
                        contact: this.excelUploadedTalents[index].contact.toString(),
                        academicYear: this.excelUploadedTalents[index].academicYear,
                        isSwabhavTalent: (this.excelUploadedTalents[index].isSwabhavTalent === "yes" ? true : false),
                        isActive: true,
                        city: this.excelUploadedTalents[index].city,
                        address: this.excelUploadedTalents[index].address,
                        pinCode: this.excelUploadedTalents[index].pinCode,
                        countryName: this.excelUploadedTalents[index].countryName,
                        stateName: this.excelUploadedTalents[index].stateName,
                        academics: [{
                              degreeName: this.excelUploadedTalents[index].degreeName,
                              specializationName: this.excelUploadedTalents[index].specializationName,
                              collegeName: this.excelUploadedTalents[index].collegeName,
                              percentage: this.excelUploadedTalents[index].percentage,
                              yearOfPassout: this.excelUploadedTalents[index].yearOfPassout,
                        }],
                  }
                  this.reformedTalents.push(talent)
                  this.talentsExcelTotalCount = this.reformedTalents.length
            }
            this.disableAddTalentsButton = false
      }

      // Add multiple talents from excel.
      addMultipleTalents(): void {

            this.showExcelProgress = true
            this.disableAddTalentsButton = true

            // Call add talent from excel API for every talent.
            for (let i = 0; i < this.talentsExcelTotalCount; i++) {
                  this.talentService.addTalentFromExcel(this.reformedTalents[i]).subscribe((response: any) => {
                        this.talentsExcelAddedCount = this.talentsExcelAddedCount + 1
                  }, (error) => {
                        if (error?.error?.error) {
                              this.excelErrorList.push(this.reformedTalents[i].email + " : " + error?.error?.error)
                              return
                        }
                        this.excelErrorList.push(this.reformedTalents[i].email + " : " + error.error)
                  }).add(() => {
                        this.talentsExcelProcessedCount = this.talentsExcelProcessedCount + 1
                        if (this.talentsExcelProcessedCount == this.talentsExcelTotalCount) {
                              this.getTalents()
                        }
                  })
            }
      }

      // On cancelling uploaded excel file.
      onCalcelUploadedExcelFile(): void {
            this.resumeDocStatus = ""
            this.isExcelUploaded = false
            this.addMultipleErrorList = []
      }

      // // Add multiple talents from excel (not being used now).
      // addMultipleTalents(): void {
      //       this.spinnerService.loadingMessage = "Adding talents"
      //       
      //       
      //       this.talentService.addTalents(this.reformedTalents).subscribe((response: any) => {
      //             this.modalRef.close('success')
      //             alert(response)
      //             this.getTalents()
      //       }, (error) => {
      //             this.resumeDocStatus = ""
      //             this.isExcelUploaded = false
      //             console.error(error)
      //             if (error.error?.error) {
      //                   alert(error.error?.error)
      //                   return
      //             }
      //             this.addMultipleErrorList = []
      //             if (error?.error?.errorList?.length > 0) {
      //                   this.excelTalentsAdded = error?.error?.message
      //                   error?.error?.errorList.forEach((err: any) => {
      //                         this.addMultipleErrorList.push(err.errorKey)
      //                   })
      //                   // for (let index = 0; index < error?.error?.length; index++) {
      //                   //       this.addMultipleErrorList.push(error.error[index]?.ErrorKey)
      //                   // }
      //                   alert("There are errors in your uploaded file.")
      //             } else {
      //                   alert(error.statusText)
      //             }
      //       })
      // }
}