import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, FormArray, Validators, FormControl } from '@angular/forms';
import { GeneralService } from 'src/app/service/general/general.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { EnquiryService, ICallRecord, ICallRecordDTO, IDegree, IEnquiry, IEnquiryDTO, IEnquiryExcel, IExamination, IExperience, IExperienceDTO, IMastersAbroad, IScore, ISearchFilterField, ISearchSection, ISpecialization, IUniversity, IWaitingList, IWaitingListDTO } from 'src/app/service/talent/enquiry.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { IPermission } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { Role, UrlConstant } from 'src/app/service/constant';
import { DatePipe, Location } from '@angular/common';
import { BatchService } from 'src/app/service/batch/batch.service';
import { ActivatedRoute } from '@angular/router';
import { TechnologyService } from 'src/app/service/technology/technology.service';
import { DegreeService } from 'src/app/service/degree/degree.service';

@Component({
      selector: 'app-talent-enquiry',
      templateUrl: './talent-enquiry.component.html',
      styleUrls: ['./talent-enquiry.component.css']
})
export class TalentEnquiryComponent implements OnInit {

      //****************************ENQUIRY*************************************** */
      // Components.
      stateList: any[]
      countryList: any[]
      enquiryTypeList: any[]
      technologyList: any[]
      degreeList: any[]
      degreeIDSet: Set<string>
      designationList: any[]
      specializationList: any[]
      salesPersonList: any[]
      academicYearList: any[]
      collegeBranchList: any[]
      examinationList: any[]
      universityList: any[]
      sourceList: any[]
      yearOfMSList: number[]
      courseList: any[]

      // tech and college component
      isTechLoading: boolean
      techLimit: number
      techOffset: number
      isCollegeLoading: boolean
      collegeLimit: number
      collegeOffset: number

      // Flags.
      isViewMode: boolean
      isOperationEnquiryUpdate: boolean

      // Enquiry.
      enquiryList: IEnquiryDTO[]
      selectedEnquiryList: any[]
      selectedEnquiry: any
      selectEnquiryTemp: any
      enquiryForm: FormGroup

      // Experience.
      showExperiencesInForm: boolean
      showExperienceColumns: boolean

      // Academic.
      showAcademicsInForm: boolean

      // Pagination.
      limit: number
      offset: number
      totalEnquiries: number
      currentPage: number
      paginationString: string
      paginationStart: number
      paginationEnd: number

      // Modal.
      modalRef: any
      @ViewChild('enquiryFormModal') enquiryFormModal: any
      @ViewChild('deleteEnquiryModal') deleteEnquiryModal: any

      // Spinner.



      // Search.
      enquirySearchForm: FormGroup
      showSearch: boolean
      isSearched: boolean
      selectedSectionName: string
      searchFilterFieldList: ISearchFilterField[]
      searchSectionList: ISearchSection[]

      // Enquiry retaled numbers.
      indexOfCurrentWorkingExp: number
      currentYear: number

      // Masters abroad related type.
      greExam: IExamination
      gmatExam: IExamination
      toeflExam: IExamination
      ieltsExam: IExamination
      universityMap: Map<string, IUniversity[]> = new Map()
      areUniversitiesLoading: boolean
      showMastersAbroadInForm: boolean
      showMBADegreeRequiredError: boolean

      // Permission.
      permission: IPermission
      roleName: string
      showForAdmin: boolean
      showForSalesPerson: boolean

      // Resume.
      isResumeUploadedToServer: boolean
      isFileUploading: boolean
      docStatus: string
      displayedFileName: string

      // Excel.
      excelUploadedEnquiries: IEnquiryExcel[]
      reformedEnquiries: any[]
      isExcelUploaded: boolean
      addMultipleErrorList: string[]
      excelEnquiriesAdded: string
      showExcelProgress: boolean
      enquiriesExcelTotalCount: number
      enquiriesExcelAddedCount: number
      enquiriesExcelProcessedCount: number
      excelErrorList: string[]
      disableAddEnquiriesButton: boolean

      // Constant.
      readonly ENQUIRY_EXCEL_DEMO_LINK = this.urlConstant.ENQUIRY_BASIC_DEMO

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
      enquiryCallRecordList: ICallRecordDTO[]
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

      //****************************WAITING LIST*************************************** */
      // Components.
      companyBranchList: any[]
      activeRequirementByCompanyList: any[]
      activeBatchByCourseList: any[]
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

      // Waiting list.
      waitingList: IWaitingListDTO[]
      waitingListForm: FormGroup

      // Waiting lits realted variables.
      waitingListBatchID: string
      waitingListRequirementID: string
      waitingListCompanyBranchID: string
      waitingListCourseID: string
      waitingListTechnologyID: string

      // Modal.
      @ViewChild('waitingListFormModal') waitingListFormModal: any

      constructor(
            private formBuilder: FormBuilder,
            private generalService: GeneralService,
            private techService: TechnologyService,
            private degreeService: DegreeService,
            public utilityService: UtilityService,
            private enquiryService: EnquiryService,
            private spinnerService: SpinnerService,
            private localService: LocalService,
            private fileOperationService: FileOperationService,
            private modalService: NgbModal,
            private datePipe: DatePipe,
            private urlConstant: UrlConstant,
            private role: Role,
            private batchService: BatchService,
            private activatedRoute: ActivatedRoute,
            private _location: Location,
      ) {
            this.initializeVariables()
      }


      get ongoingOperations() {
            return this.spinnerService.ongoingOperations
      }

      ngOnInit(): void {

            this.createEnquiryForm()
            this.createCallRecordForm()
            this.createEnquirySearchForm()
            this.getAllComponents()
      }

      // Initialize all global variables.
      initializeVariables(): void {
            //****************************ENQUIRY*************************************** */
            // Components.
            this.stateList = []
            this.countryList = []
            this.enquiryTypeList = []
            this.technologyList = []
            this.degreeList = []
            this.specializationList = []
            this.designationList = []
            this.salesPersonList = []
            this.academicYearList = []
            this.sourceList = []
            this.collegeBranchList = []
            this.examinationList = []
            this.universityList = []
            this.yearOfMSList = []
            this.courseList = []
            this.degreeIDSet = new Set()

            // Flags.
            this.isViewMode = false
            this.isOperationEnquiryUpdate = false
            this.isTechLoading = false
            this.isCollegeLoading = false

            // Enquiry.
            this.enquiryList = []
            this.selectedEnquiryList = []
            this.selectedEnquiry = {}
            this.selectEnquiryTemp = {}

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
            this.techLimit = 10
            this.techOffset = 0
            this.collegeLimit = 10
            this.collegeOffset = 0

            // Spinner.
            this.spinnerService.loadingMessage = "Getting Talents"


            // Search.
            this.isSearched = false
            this.showSearch = false
            this.searchFilterFieldList = []
            this.searchSectionList = [
                  {
                        name: "Personal",
                        isSelected: true
                  },
                  {
                        name: "Enquiry",
                        isSelected: false
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
                        name: "Waiting List",
                        isSelected: false
                  },
                  {
                        name: "Others",
                        isSelected: false
                  }
            ]
            this.selectedSectionName = "Personal"

            // Enquiry retaled numbers.
            this.indexOfCurrentWorkingExp = -1
            this.currentYear = new Date().getFullYear()

            // Masters abroad.
            this.areUniversitiesLoading = false
            this.showMastersAbroadInForm = false
            this.showMBADegreeRequiredError = false

            // Permision.
            this.showForAdmin = false
            this.showForSalesPerson = true
            // Get permissions from menus using utilityService function.
            this.permission = this.utilityService.getPermission(this.urlConstant.TALENT_ENQUIRY)
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

            // Resume.
            this.docStatus = ""
            this.displayedFileName = "Select file"
            this.isResumeUploadedToServer = false
            this.isFileUploading = false

            // Excel.
            this.isExcelUploaded = false
            this.addMultipleErrorList = []
            this.enquiriesExcelAddedCount = 0
            this.enquiriesExcelTotalCount = 0
            this.enquiriesExcelProcessedCount = 0
            this.showExcelProgress = false
            this.excelErrorList = []
            this.disableAddEnquiriesButton = false

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
            this.enquiryCallRecordList = [] as ICallRecordDTO[]

            //****************************ALLOCATION*************************************** */
            // Flags.
            this.showAllocateSalespersonToOneEnquiry = false
            this.showAllocateSalespersonToEnquiries = false
            this.multipleSelect = false

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

            // Waiting list redirection variables.
            this.waitingListBatchID = null
            this.waitingListRequirementID = null
            this.waitingListCompanyBranchID = null
            this.waitingListCourseID = null
            this.waitingListTechnologyID = null

            // Waiting list.
            this.waitingList = [] as IWaitingListDTO[]
      }

      //*********************************************CREATE FORMS************************************************************
      // Create new enquiry form.
      createEnquiryForm(): void {
            this.enquiryForm = this.formBuilder.group({
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
                  enquirySource: new FormControl(null),
                  address: new FormControl(null, [Validators.pattern(/^[.0-9a-zA-Z\s,-\/]+$/)]),
                  city: new FormControl(null, [Validators.pattern(/^[a-zA-Z]+([a-zA-Z ]?)+$/)]),
                  state: new FormControl(null),
                  country: new FormControl(null),
                  pinCode: new FormControl(null, [Validators.pattern(/^[1-9][0-9]{5}$/)]),
                  technologies: new FormControl(Array()),
                  courses: new FormControl(null),
                  enquiryDate: new FormControl(null, [Validators.required]),
                  enquiryType: new FormControl(null, [Validators.required]),
                  academics: this.formBuilder.array([]),
                  isExperience: new FormControl(false),
                  isAcademic: new FormControl(false),
                  experiences: this.formBuilder.array([]),
                  salesPerson: new FormControl(null),
                  resume: new FormControl(null),
                  additionalDetails: new FormControl(null, [Validators.maxLength(1000)]),
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
            })
      }

      // Add new experience to enquiry form.
      addExperience(): void {
            this.enquiryExperienceControlArray.push(this.formBuilder.group({
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

      // Add new college to enquiry form.
      addAcademic(): void {
            this.enquiryAcademicControlArray.push(this.formBuilder.group({
                  id: new FormControl(null),
                  degree: new FormControl(null, [Validators.required]),
                  specialization: new FormControl(null, [Validators.required]),
                  college: new FormControl(null, [Validators.required, Validators.maxLength(200)]),
                  percentage: new FormControl(null, [Validators.required, Validators.min(1),
                  Validators.pattern(/(^100([.]0{1,2})?)$|(^\d{1,2}([.]\d{1,2})?)$/i), Validators.required]),
                  passout: new FormControl(null, [Validators.required, Validators.min(1980),
                  Validators.max(this.currentYear + 3)])
            }))
      }

      // Add new matsers abroad to enquiry form.
      addMastersAbroad(): void {
            this.enquiryForm.addControl('mastersAbroad', this.formBuilder.group({
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
                  ielts: new FormControl(null, [Validators.min(0), Validators.pattern(/^(10|\d)(\.5)?$/),
                  Validators.max(this.ieltsExam?.totalMarks)]),
                  talentID: new FormControl(null),
                  enquiryID: new FormControl(null)
            }))
      }

      // Create new enquiry search form.
      createEnquirySearchForm(): void {
            this.enquirySearchForm = this.formBuilder.group({
                  firstName: new FormControl(null),
                  lastName: new FormControl(null),
                  contact: new FormControl(null),
                  email: new FormControl(null),
                  college: new FormControl(null),
                  isExperience: new FormControl(null),
                  isMastersAbroad: new FormControl(null),
                  city: new FormControl(null),
                  technologies: new FormControl(null),
                  passout: new FormControl(null),
                  enquiryType: new FormControl(null),
                  degrees: new FormControl(null),
                  designations: new FormControl(null),
                  experienceTechnologies: new FormControl(null),
                  totalExperience: new FormControl(null),
                  salesPersonIDs: new FormControl(null),
                  academicYears: new FormControl(null),
                  callRecordPurposeID: new FormControl(null),
                  callRecordOutcomeID: new FormControl(null),
                  minimumExperience: new FormControl(null, [Validators.min(0), Validators.max(30)]),
                  maximumExperience: new FormControl(null, [Validators.min(0), Validators.max(30)]),
                  searchAllEnquiries: new FormControl(null),
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
                  enquiryFromDate: new FormControl(null),
                  enquiryToDate: new FormControl(null),
                  isLastThirtyDays: new FormControl(null),
                  enquirySource: new FormControl(null),
            })
      }

      // Create new call record form.
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

      //*********************************************ADD FOR ENQUIRY FUNCTIONS************************************************************
      // On clicking add new enquiry button.
      onAddNewEnquiryClick(): void {
            this.indexOfCurrentWorkingExp = -1
            this.isViewMode = false
            this.enableEnquiryForm()
            this.showExperiencesInForm = false
            this.showAcademicsInForm = false
            this.showMastersAbroadInForm = false
            this.showMBADegreeRequiredError = false
            this.isOperationEnquiryUpdate = false
            this.showExcelProgress = false
            this.createEnquiryForm()
            this.stateList = []
            this.enquiryForm.get('state').disable()
            this.universityList = []
            this.openModal(this.enquiryFormModal, 'xl')
      }

      // Handle creation or updation of masters abroad of enquiry.
      assignMastersAbroad(addOrUpdate: string): IEnquiry {
            let scoresArray: IScore[] = []
            let enquiry: IEnquiry = this.enquiryForm.value
            let mastersAbroad: IMastersAbroad
            if (enquiry.mastersAbroad) {
                  if (this.enquiryForm.get('mastersAbroad').get('gre').value != null) {
                        let score: IScore = {
                              id: this.enquiryForm.get('mastersAbroad').get('greID').value,
                              marksObtained: this.enquiryForm.get('mastersAbroad').get('gre').value,
                              examinationID: this.greExam.id,
                              mastersAbroadID: this.enquiryForm.get('mastersAbroad').get('id').value
                        }
                        scoresArray.push(score)
                  }
                  if (this.enquiryForm.get('mastersAbroad').get('gmat').value != null) {
                        let score: IScore = {
                              id: this.enquiryForm.get('mastersAbroad').get('gmatID').value,
                              marksObtained: this.enquiryForm.get('mastersAbroad').get('gmat').value,
                              examinationID: this.gmatExam.id,
                              mastersAbroadID: this.enquiryForm.get('mastersAbroad').get('id').value
                        }
                        scoresArray.push(score)
                  }
                  if (this.enquiryForm.get('mastersAbroad').get('toefl').value != null) {
                        let score: IScore = {
                              id: this.enquiryForm.get('mastersAbroad').get('toeflID').value,
                              marksObtained: this.enquiryForm.get('mastersAbroad').get('toefl').value,
                              examinationID: this.toeflExam.id,
                              mastersAbroadID: this.enquiryForm.get('mastersAbroad').get('id').value
                        }
                        scoresArray.push(score)
                  }
                  if (this.enquiryForm.get('mastersAbroad').get('ielts').value != null) {
                        let score: IScore = {
                              id: this.enquiryForm.get('mastersAbroad').get('ieltsID').value,
                              marksObtained: this.enquiryForm.get('mastersAbroad').get('ielts').value,
                              examinationID: this.ieltsExam.id,
                              mastersAbroadID: this.enquiryForm.get('mastersAbroad').get('id').value
                        }
                        scoresArray.push(score)
                  }
                  mastersAbroad = {
                        scores: scoresArray,
                        degreeID: this.enquiryForm.get('mastersAbroad').get('degree').value.id,
                        countries: enquiry.mastersAbroad.countries,
                        universities: enquiry.mastersAbroad.universities,
                        yearOfMS: enquiry.mastersAbroad.yearOfMS,
                        enquiryID: enquiry.id
                  }
                  if (addOrUpdate == "update") {
                        if (this.selectEnquiryTemp.mastersAbroad) {
                              mastersAbroad.id = this.selectEnquiryTemp.mastersAbroad.id
                        }
                  }
                  enquiry.mastersAbroad = mastersAbroad
                  return enquiry
            }
            enquiry.mastersAbroad = null
            return enquiry
      }

      // Add new enquiry.
      addEnquiry(): void {
            let enquiry: IEnquiry = this.assignMastersAbroad("add")
            this.patchIDFromObjectsForEnquiry(enquiry)
            this.spinnerService.loadingMessage = "Adding Enquiry"


            this.enquiryService.addEnquiry(enquiry).subscribe((response: any) => {
                  this.modalRef.close('success')
                  this.getEnquiries()
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

      //*********************************************UPDATE AND VIEW FOR ENQUIRY FUNCTIONS************************************************************
      // On clicking view enquiry button.
      onViewEnquiryClick(id: any): void {
            this.indexOfCurrentWorkingExp = -1
            this.isViewMode = true
            this.showExperiencesInForm = false
            this.showAcademicsInForm = false
            this.showMastersAbroadInForm = false
            this.showMBADegreeRequiredError = false
            this.selectedEnquiry = null
            this.createEnquiryForm()
            this.assignSelectedEnquiry(id)

            // Resume.
            this.displayedFileName = "No resume uploaded"
            if (this.selectedEnquiry.resume) {
                  this.displayedFileName = `<a href=${this.selectedEnquiry.resume} target="_blank">Resume present</a>`
            }
            // Enquiry date.
            let enquiryDate = this.selectedEnquiry.enquiryDate
            if (enquiryDate) {
                  this.selectedEnquiry.enquiryDate = this.datePipe.transform(enquiryDate, 'yyyy-MM-dd')
            }

            // Expereinces.
            this.enquiryForm.setControl("experiences", this.formBuilder.array([]))
            if (this.selectedEnquiry.experiences && this.selectedEnquiry.experiences.length > 0) {
                  for (let i = 0; i < this.selectedEnquiry.experiences.length; i++) {
                        if (this.selectedEnquiry.experiences[i].toDate == null) {
                              this.indexOfCurrentWorkingExp = i
                        }
                        this.calculateYearsOfExperience(this.selectedEnquiry.experiences[i])
                        let fromDate = this.selectedEnquiry.experiences[i]?.fromDate
                        if (fromDate) {
                              this.selectedEnquiry.experiences[i].fromDate = this.datePipe.transform(fromDate, 'yyyy-MM')
                        }
                        let toDate = this.selectedEnquiry.experiences[i]?.toDate
                        if (toDate) {
                              this.selectedEnquiry.experiences[i].toDate = this.datePipe.transform(toDate, 'yyyy-MM')
                        }
                        this.calculateYearsOfExperience(this.selectedEnquiry.experiences[i])
                        this.addExperience()
                  }
                  this.showExperiencesInForm = true
            }

            if (this.indexOfCurrentWorkingExp != -1) {
                  this.enquiryExperienceControlArray.at(this.indexOfCurrentWorkingExp).get('isCurrentWorking').setValue('true')
                  let control = this.enquiryExperienceControlArray.controls[this.indexOfCurrentWorkingExp]
                  if (control instanceof FormGroup) {
                        control.removeControl('toDate')
                  }
            }

            // Academics.
            this.enquiryForm.setControl("academics", this.formBuilder.array([]))
            if (this.selectedEnquiry.academics && this.selectedEnquiry.academics.length > 0) {
                  for (let i = 0; i < this.selectedEnquiry.academics.length; i++) {
                        this.addAcademic()
                        this.getSpecializationListByDegreeID(this.selectedEnquiry.academics[i].degree.id)
                  }
                  this.showAcademicsInForm = true
                  this.enquiryForm.get("isAcademic").setValue(true)
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

            this.selectEnquiryTemp = {}

            for (var k in this.selectedEnquiry) {
                  this.selectEnquiryTemp[k] = this.selectedEnquiry[k]
            }

            if (this.selectedEnquiry.mastersAbroad) {
                  this.addMastersAbroad()
                  this.showMastersAbroadInForm = true
                  this.enquiryForm.get("isMastersAbroad").setValue(true)
                  for (let i = 0; i < this.selectedEnquiry.mastersAbroad.scores.length; i++) {
                        if (this.selectedEnquiry.mastersAbroad.scores[i].examination.name == "GRE") {
                              greScore = this.selectedEnquiry.mastersAbroad.scores[i].marksObtained
                              greID = this.selectedEnquiry.mastersAbroad.scores[i].id
                        }
                        if (this.selectedEnquiry.mastersAbroad.scores[i].examination.name == "GMAT") {
                              gmatScore = this.selectedEnquiry.mastersAbroad.scores[i].marksObtained
                              gmatID = this.selectedEnquiry.mastersAbroad.scores[i].id
                        }
                        if (this.selectedEnquiry.mastersAbroad.scores[i].examination.name == "TOEFL") {
                              toeflScore = this.selectedEnquiry.mastersAbroad.scores[i].marksObtained
                              toeflID = this.selectedEnquiry.mastersAbroad.scores[i].id
                        }
                        if (this.selectedEnquiry.mastersAbroad.scores[i].examination.name == "IELTS") {
                              ieltsScore = this.selectedEnquiry.mastersAbroad.scores[i].marksObtained
                              ieltsID = this.selectedEnquiry.mastersAbroad.scores[i].id
                        }
                  }

                  let tempMastersAbroad: any = {
                        id: this.selectedEnquiry.mastersAbroad.id,
                        degree: this.selectedEnquiry.mastersAbroad.degree,
                        countries: this.selectedEnquiry.mastersAbroad.countries,
                        universities: this.selectedEnquiry.mastersAbroad.universities,
                        yearOfMS: this.selectedEnquiry.mastersAbroad.yearOfMS,
                        gre: greScore,
                        greID: greID,
                        gmat: gmatScore,
                        gmatID: gmatID,
                        toefl: toeflScore,
                        toeflID: toeflID,
                        ielts: ieltsScore,
                        ieltsID: ieltsID
                  }
                  this.selectEnquiryTemp.mastersAbroad = tempMastersAbroad
            }
            else {
                  this.selectEnquiryTemp.mastersAbroad = []
            }

            // State.
            if (this.selectedEnquiry.country != undefined) {
                  this.getStateListByCountryID(this.selectedEnquiry.country)
            }

            // Populate enquiry form.
            this.enquiryForm.patchValue(this.selectEnquiryTemp)

            // Disable talent form.
            this.enquiryForm.disable()

            // Open talent form modal.
            this.openModal(this.enquiryFormModal, 'xl')
      }

      // On clicking update enquiry button.
      onUpdateEnquiryClick(): void {
            this.isViewMode = false
            this.isOperationEnquiryUpdate = true
            this.enableEnquiryForm()
            if (this.enquiryForm.get('state').value == null) {
                  this.enquiryForm.get('state').disable()
            }
      }

      // Update enquiry.
      updateEnquiry(): void {
            let enquiry: IEnquiry = this.assignMastersAbroad("update")
            this.patchIDFromObjectsForEnquiry(enquiry)
            this.spinnerService.loadingMessage = "Updating Enquiry"


            this.enquiryService.updateEnquiry(enquiry).subscribe((response: any) => {
                  this.modalRef.close('success')
                  this.getEnquiries()
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

      // Convert enquiry to talent.
      convertEnquiryToTalent(enquiry: IEnquiryDTO): void {
            //************WHEN SOME COMPULSARY FIELDS WERE EMPTY WHEN COMING FROM ENQUIRY FORM********** */
            // if (enquiry.academics != null || enquiry.academics?.length != 0){
            //       let isSpecializationPresent: boolean = true
            //       let isCollegePresent: boolean = true
            //       let isPercentagePresent: boolean = true
            //       let fieldNamesArray:string[] = []
            //       let fieldNamesString:string 
            //       for (let i = 0; i < enquiry.academics.length; i++){
            //             if (enquiry.academics[i].specialization == null){
            //                   isSpecializationPresent = false
            //             }
            //             if (enquiry.academics[i].percentage == null){
            //                   isPercentagePresent = false
            //             }
            //             if (enquiry.academics[i].college == null){
            //                   isCollegePresent = false
            //             }
            //       }

            //       if (!isSpecializationPresent){
            //             fieldNamesArray.push("Specialization")
            //       }
            //       if (!isPercentagePresent){
            //             fieldNamesArray.push("Percentage")
            //       }
            //       if (!isCollegePresent){
            //             fieldNamesArray.push("College")
            //       }
            //       if (!enquiry.academicYear){
            //             fieldNamesArray.push("Academic Year")
            //       }
            //       if (!enquiry.enquiryType){
            //             fieldNamesArray.push("Enquiry Type")
            //       }
            //       if (fieldNamesArray.length != 0){
            //             fieldNamesString = fieldNamesArray.join(", ")
            //             alert("Please fill " + fieldNamesString + " fields before converting to talent")
            //             return
            //       }
            // }
            if (!confirm("This action will convert the enquiry to talent.\nDo you wish to continue?")) {
                  return
            }
            this.spinnerService.loadingMessage = "Converting to talent"


            this.enquiryService.convertEnquiryToTalent(enquiry.id).subscribe((response: any) => {
                  this.getEnquiries()
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

      //*********************************************DELETE FOR ENQUIRY FUNCTIONS************************************************************
      // On clicking delete enquiry button. 
      onDeleteEnquiryClick(enquiryID: string): void {
            this.assignSelectedEnquiry(enquiryID)
            this.openModal(this.deleteEnquiryModal, 'md').result.then(() => {
                  this.deleteEnquiry()
            }, (err) => {
                  console.error(err)
                  return
            })
      }

      // Delete enquiry.
      deleteEnquiry(): void {
            this.spinnerService.loadingMessage = "Deleting Enquiry"


            this.enquiryService.deleteEnquiry(this.selectedEnquiry.id).subscribe((response: any) => {
                  this.modalRef.close()
                  this.getEnquiries()
                  this.fileOperationService.deleteUploadedFile(this.selectedEnquiry.resume)
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

      // Delete resume.
      deleteResume(): void {
            this.fileOperationService.deleteUploadedFile().subscribe((data: any) => {
            }, (error) => {
                  console.error(error)
            })
      }

      //*********************************************FUNCTIONS FOR ENQUIRY FORM************************************************************
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

      // Validate enquiry form.
      validateEnquiryForm(): void {
            if (this.enquiryForm.invalid) {
                  // Left purposely.
                  //console.log("invalid controls:", this.findInvalidControls())
                  this.enquiryForm.markAllAsTouched()
                  return
            }
            if (this.isOperationEnquiryUpdate) {
                  this.updateEnquiry()
                  return
            }
            this.addEnquiry()
      }

      // Change page for pagination.
      changePage($event): void {
            this.offset = $event - 1
            this.currentPage = $event
            this.getEnquiries()
      }

      // Set total designations list on current page.
      setPaginationString(): void {
            this.paginationStart = this.limit * this.offset + 1
            this.paginationEnd = +this.limit + this.limit * this.offset
            if (this.totalEnquiries < this.paginationEnd) {
                  this.paginationEnd = this.totalEnquiries
            }
      }

      // On clicking currently working or not.
      isWorkingClicked(isWorking: string, index: number, expereince: any): void {
            if (isWorking == "false") {
                  expereince.addControl('toDate', this.formBuilder.control(null, [Validators.required]))
                  this.indexOfCurrentWorkingExp = -1
                  return
            }
            expereince.removeControl('toDate')
            this.indexOfCurrentWorkingExp = index
      }

      // Assign selected enquiry.
      assignSelectedEnquiry(id: string): void {
            let length = this.enquiryList.length
            for (let index = 0; index < length; index++) {
                  if (this.enquiryList[index].id == id) {
                        this.selectedEnquiry = this.enquiryList[index]
                  }
            }
      }

      // Delete experience from enquiry form.
      deleteExperience(index: number): void {
            if (confirm("Are you sure you want to delete experience?")) {
                  if (this.indexOfCurrentWorkingExp == index) {
                        this.indexOfCurrentWorkingExp = -1
                  }
                  if (this.indexOfCurrentWorkingExp > index) {
                        this.indexOfCurrentWorkingExp = this.indexOfCurrentWorkingExp - 1
                  }
                  this.enquiryForm.markAsDirty()
                  this.enquiryExperienceControlArray.removeAt(index)
                  //this.otherDesignation[index] = false
            }
      }

      // Delete college from enquiry form.
      deleteAcademic(index: number): void {
            if (confirm("Are you sure you want to delete academic?")) {
                  this.enquiryForm.markAsDirty()
                  this.enquiryAcademicControlArray.removeAt(index)
            }
      }

      // Get experiences array from enquiry form.
      get enquiryExperienceControlArray(): FormArray {
            return this.enquiryForm.get("experiences") as FormArray
      }

      // Get academics array from enquiry form.
      get enquiryAcademicControlArray(): FormArray {
            return this.enquiryForm.get("academics") as FormArray
      }

      // On clicking isExperience checkbox.
      toggleExperienceControls(event: any): void {
            if (this.showExperiencesInForm) {
                  if (confirm("Are you sure you want to delete all experience details?")) {
                        this.enquiryForm.setControl('experiences', this.formBuilder.array([]))
                        this.enquiryForm.get("isExperience").setValue(false)
                        this.showExperiencesInForm = false
                        return
                  }
                  event.target.checked = true
                  return
            }
            if (this.enquiryForm.get("academicYear").value != 5) {
                  event.target.checked = false
                  alert("Experience details can only be added when academic year is selected as completed")
                  return
            }
            this.addExperience()
            this.showExperiencesInForm = true
            this.enquiryForm.get("isExperience").setValue(true)
      }

      // On clicking isAcademic checkbox.
      toggleAcademicControls(event: any): void {
            if (this.showAcademicsInForm) {
                  if (confirm("Are you sure you want to delete all academic details?")) {
                        this.showAcademicsInForm = false
                        this.enquiryForm.setControl('academics', this.formBuilder.array([]))
                        this.enquiryForm.get("isAcademic").setValue(false)
                        return
                  }
                  event.target.checked = true
                  return
            }
            this.showAcademicsInForm = true
            this.enquiryForm.get("isAcademic").setValue(true)
            this.addAcademic()
      }

      // Used to open modal.
      openModal(content: any, size?: string): NgbModalRef {
            this.isResumeUploadedToServer = false
            this.disableAddEnquiriesButton = false
            this.displayedFileName = "Select file"
            this.docStatus = ""
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

      // Calculate total years of experience for one experience.
      calculateYearsOfExperience(experience: IExperienceDTO): void {
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

      // Calculate total years of experience for talent.
      calculateTotalYearsOfExperinceOfTalent(enquiry: IEnquiryDTO): string {
            let totalYears: number = 0
            enquiry.totalYearsOfExperience = totalYears
            if ((enquiry.experiences) && (enquiry.experiences.length != 0)) {
                  for (let i = 0; i < enquiry.experiences.length; i++) {
                        this.calculateYearsOfExperience(enquiry.experiences[i])
                        totalYears = totalYears + enquiry.experiences[i].yearsOfExperienceInNumber
                  }
                  enquiry.totalYearsOfExperience = totalYears
                  let numberOfYears: number = Math.floor(totalYears)
                  let numberOfMonths: number = Math.round((totalYears % 1) * 12)
                  if (numberOfYears == 0 && numberOfMonths == 0) {
                        enquiry.totalYearsOfExperience = 0
                        enquiry.totalYearsOfExperienceInString = "Fresher"
                        return
                  }
                  if (isNaN(numberOfYears) || isNaN(numberOfMonths)) {
                        enquiry.totalYearsOfExperience = 0
                        enquiry.totalYearsOfExperienceInString = "Fresher"
                        return
                  }
                  enquiry.totalYearsOfExperienceInString = numberOfYears + "." + numberOfMonths + " Year(s)"
                  return
            }
            enquiry.totalYearsOfExperience = 0
            enquiry.totalYearsOfExperienceInString = "Fresher"
      }

      // Get all invalid controls in enquiry form.
      public findInvalidControls(): any {
            const invalid = []
            const controls = this.enquiryForm.controls
            for (const name in controls) {
                  if (controls[name].invalid) {
                        invalid.push(name)
                  }
            }
            return invalid
      }

      // Used to dismiss modal.
      dismissFormModal(modal: NgbModalRef): void {
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
            modal.dismiss()
            this.isResumeUploadedToServer = false
            this.disableAddEnquiriesButton = false
            this.displayedFileName = "Select file"
            this.docStatus = ""
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
                        this.enquiryForm.markAsDirty()
                        this.enquiryForm.patchValue({
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

      // On clicking isMastersAbroad checkbox.
      toggleShowMSInUsFormControls(event: any): void {
            if (this.showMastersAbroadInForm) {
                  if (confirm("Are you sure you want to delete all MS abroad details?")) {
                        this.showMastersAbroadInForm = false
                        this.removeMastersAbroad()
                        this.enquiryForm.get("isMastersAbroad").setValue(false)
                        return
                  }
                  event.target.checked = true
                  return
            }
            this.showMastersAbroadInForm = true
            this.enquiryForm.get("isMastersAbroad").setValue(true)
            this.addMastersAbroad()
            this.enquiryForm.get('mastersAbroad').get('universities').disable()
      }

      // On changing degree in masters abroad, add or remove gmat score field.
      onDegreeChange(degree: IDegree): void {
            if (degree && degree.name == "M.B.A.") {
                  this.enquiryForm.get('mastersAbroad').get('gmat').setValidators([Validators.required, Validators.min(0), Validators.max(this.gmatExam?.totalMarks)])
                  this.enquiryForm.get('mastersAbroad').get('gmat').updateValueAndValidity()
                  this.showMBADegreeRequiredError = true
                  return
            }
            this.enquiryForm.get('mastersAbroad').get('gmat').setValidators([Validators.min(0), Validators.max(this.gmatExam?.totalMarks)])
            this.enquiryForm.get('mastersAbroad').get('gmat').updateValueAndValidity()
            this.showMBADegreeRequiredError = false
      }

      // Toggle visibility of experience columns.
      toggleShowExperienceColumns(): void {
            this.showExperienceColumns = !this.showExperienceColumns
      }

      // Calculates current company name of enquiry.
      calculateCurrentCompanyName(enquiry: IEnquiryDTO): void {
            if (enquiry.experiences && enquiry.experiences?.length != 0 && enquiry.experiences[enquiry.experiences.length - 1]?.toDate == null) {
                  enquiry.currentCompanyName = enquiry.experiences[enquiry.experiences.length - 1]?.company
                  return
            }
            enquiry.currentCompanyName = "Not Currently Working"
      }

      // Calculate the last designation of enquiry.
      calculateDesignation(enquiry: IEnquiryDTO): void {
            if (enquiry.experiences && enquiry.experiences?.length != 0) {
                  enquiry.lastDesignation = enquiry.experiences[enquiry.experiences.length - 1]?.designation?.position
                  return
            }
            enquiry.lastDesignation = "No experience"
      }

      // Calculate the package of current company of enquiry.
      calculateCurrentCompanyPackage(enquiry: IEnquiryDTO): void {
            if (enquiry.experiences && enquiry.experiences?.length != 0 && enquiry.experiences[enquiry.experiences.length - 1]?.toDate == null
                  && enquiry.experiences[enquiry.experiences.length - 1]?.package != null) {
                  enquiry.currentCompanyPackage = enquiry.experiences[enquiry.experiences.length - 1]?.package / 100000 + " Lpa"
                  return
            }
            enquiry.currentCompanyPackage = "Package not mentioned"
      }

      // Get the technologies of all the experiences of talent.
      calculateAllExperiencesTechnologies(enquiry: IEnquiryDTO): void {
            enquiry.allExperiencesTechnologies = []
            if (enquiry.experiences && enquiry.experiences.length != 0) {
                  for (let i = 0; i < enquiry.experiences.length; i++) {
                        if (enquiry.experiences[i]?.technologies && enquiry.experiences[i].technologies?.length > 0) {
                              for (let j = 0; j < enquiry.experiences[i].technologies.length; j++) {
                                    if (enquiry.allExperiencesTechnologies.includes(enquiry.experiences[i].technologies[j].language)) {
                                          continue
                                    }
                                    enquiry.allExperiencesTechnologies.push(enquiry.experiences[i].technologies[j].language)
                              }
                        }
                  }
            }
      }

      // Calculate the enquiry's first college name.
      calculateFirstCollegeName(enquiry: IEnquiryDTO): void {
            if (enquiry.academics && enquiry.academics?.length != 0 && enquiry.academics[0]?.college != null) {
                  enquiry.firstCollegeName = enquiry.academics[0]?.college
                  return
            }
            enquiry.firstCollegeName = "Not Mentioned"
      }

      // Enable enquiry form.
      enableEnquiryForm(): void {
            this.enquiryForm.enable()
            this.enquiryForm.get('code').disable()
      }

      // On clearing all countries in masters abroad countries.
      onClearCountries(): void {
            this.universityList = []
            // let arr = this.talentForm.get('mastersAbroad') as FormArray
            // Clear all universities. ************** -n
            this.enquiryForm.get('mastersAbroad')?.get('universities').reset()
            this.enquiryForm.get('mastersAbroad').get('universities').disable()
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
            let selectedUniversities: IUniversity[] = this.enquiryForm.get('mastersAbroad').get('universities').value
            // Change it to remove specific universities *********** -n
            this.enquiryForm.get('mastersAbroad').get('universities').reset()
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
            this.enquiryForm.get('mastersAbroad').get('universities').enable()
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

      // Remove masters abroad form enquiry form.
      removeMastersAbroad(): void {
            this.enquiryForm.removeControl('mastersAbroad')
            this.universityList = []
      }

      // Will be called on [addTag]= "addCollegeToList" args will be passed automatically.
      addCollegeToList(option: any): Promise<any> {
            return new Promise((resolve) => {
                  resolve(option)
            }
            )
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

      // Extract ID from objects in enquiry form.
      patchIDFromObjectsForEnquiry(enquiry: IEnquiry): void {
            if (this.enquiryForm.get('enquirySource').value) {
                  enquiry.sourceID = this.enquiryForm.get('enquirySource').value.id
                  delete enquiry['enquirySource']
            }
            if (this.enquiryForm.get('salesPerson').value) {
                  enquiry.salesPersonID = this.enquiryForm.get('salesPerson').value.id
                  delete enquiry['salesPerson']
            }
            for (let i = 0; i < this.enquiryExperienceControlArray.length; i++) {
                  if (this.enquiryExperienceControlArray.at(i).get('designation').value) {
                        enquiry.experiences[i].designationID = this.enquiryExperienceControlArray.at(i).get('designation').value.id
                        delete enquiry.experiences['designation']
                  }
            }
            for (let i = 0; i < this.enquiryAcademicControlArray.length; i++) {
                  if (this.enquiryAcademicControlArray.at(i).get('degree').value) {
                        enquiry.academics[i].degreeID = this.enquiryAcademicControlArray.at(i).get('degree').value.id
                        delete enquiry.academics['degree']
                  }
                  if (this.enquiryAcademicControlArray.at(i).get('specialization').value) {
                        enquiry.academics[i].specializationID = this.enquiryAcademicControlArray.at(i).get('specialization').value.id
                        delete enquiry.academics['specialization']
                  }
            }
      }

      // Calculate some fields of all enquiries.
      calculateFieldsOfEnquiry(enquiries: IEnquiryDTO[]): void {
            for (let i = 0; i < enquiries.length; i++) {
                  this.calculateCurrentCompanyName(enquiries[i])
                  this.calculateCurrentCompanyPackage(enquiries[i])
                  this.calculateAllExperiencesTechnologies(enquiries[i])
                  this.calculateDesignation(enquiries[i])
                  this.calculateFirstCollegeName(enquiries[i])
                  this.calculateTotalYearsOfExperinceOfTalent(enquiries[i])
                  if (enquiries[i].expectedCTC && enquiries[i].expectedCTC > 0) {
                        enquiries[i].expectedCTCInString = this.formatNumberInIndianRupeeSystem(enquiries[i].expectedCTC)
                  }
            }
      }

      // Nagivate to previous url.
      backToPreviousPage(): void {
            this._location.back()
      }

      //*********************************************CRUD FUNCTIONS FOR ENQUIRY CALL RECORDS************************************************************
      // On clicking call records popup in enquiry list.
      getCallRecordsForSelectedEnquiry(enquiryID: any): void {
            this.isOperationCallRecordUpdate = false
            this.enquiryCallRecordList = []
            this.assignSelectedEnquiry(enquiryID)
            this.getAllCallRecords()
            this.showCallRecordForm = false
            this.openModal(this.callRecordFormModal, 'xl')
      }

      // Get all call records by enquiry id.
      getAllCallRecords(): void {
            this.spinnerService.loadingMessage = "Getting All Call Records"


            this.enquiryService.getCallRecordsByEnquiry(this.selectedEnquiry.id).subscribe(response => {
                  this.enquiryCallRecordList = response
                  this.formatDateTimeOfEnquiryCallRecords()
                  this.getEnquiries()
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
            this.enquiryService.addCallRecord(this.callRecordForm.value, this.selectedEnquiry.id).subscribe((response: any) => {
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

      // Update call record form.
      OnUpdateCallRecordButtonClick(index: number): void {
            this.isOperationCallRecordUpdate = true
            this.showPurposeJobFields = false
            this.showOutcomeTargetDateField = false
            this.createCallRecordForm()
            this.showCallRecordForm = true
            if (this.enquiryCallRecordList[index].purpose) {
                  this.getOutcomesByPurpose(this.enquiryCallRecordList[index].purpose)
            }
            if (this.enquiryCallRecordList[index].targetDate) {
                  this.addTargetDateInCallRecordForm()
            }
            this.callRecordForm.patchValue(this.enquiryCallRecordList[index])
      }

      // Update call record.
      updateCallrecord(): void {
            this.spinnerService.loadingMessage = "Updating Call Record"


            let callRecord: ICallRecord = this.callRecordForm.value
            this.patchIDFromObjectsForCallRecord(callRecord)
            this.enquiryService.updateCallRecord(this.callRecordForm.value, this.selectedEnquiry.id).subscribe((response: any) => {
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
      deleteCallRecord(enquiryCallRecordID: string): void {
            if (confirm("Are you sure you want to delete the call record?")) {
                  this.spinnerService.loadingMessage = "Deleting Call record"


                  this.enquiryService.deleteCallRecord(enquiryCallRecordID, this.selectedEnquiry.id).subscribe((response: any) => {
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
      togglePurposeJobFields(purpose: any): void {
            if (purpose != null && purpose.purpose == "Placement") {
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
      toggleOutcomeTargetDateField(outcome: any): void {
            if (outcome != null && outcome.outcome == "Interested, Follow Up") {
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

      // Format date time field of enquiry call records by removing timestamp.
      formatDateTimeOfEnquiryCallRecords(): void {
            for (let i = 0; i < this.enquiryCallRecordList?.length; i++) {
                  let dateTime = this.enquiryCallRecordList[i].dateTime
                  if (dateTime) {
                        this.enquiryCallRecordList[i].dateTime = this.datePipe.transform(dateTime, 'yyyy-MM-ddTHH:mm:ss')
                  }
                  let targetDate = this.enquiryCallRecordList[i].targetDate
                  if (targetDate) {
                        this.enquiryCallRecordList[i].targetDate = this.datePipe.transform(targetDate, 'yyyy-MM-dd')
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

      //*********************************************FUNCTIONS FOR ENQUIRY SEARCH FORM************************************************************
      // Reset search form and get all enquiries.
      resetSearchAndGetAll(): void {
            this.searchFilterFieldList = []
            this.spinnerService.loadingMessage = "Getting All Enquiries"
            this.isSearched = false
            this.showSearch = false
            this.createEnquirySearchForm()
            this.changePage(1)
      }

      // Reset searhc form.
      resetSearchForm(): void {
            this.searchFilterFieldList = []
            this.outcomeList = []
            this.enquirySearchForm.reset()
      }

      // On clicking search button.
      OnSerachEquiriesButtonClick(): void {
            this.spinnerService.loadingMessage = "Getting Searched Enquiries"
            this.isSearched = true
            this.changePage(1)
      }

      // Get all searched enquiries.
      getAllSearchedEnquiries(): void {
            let roleNameAndLogin: any = {
                  roleName: this.localService.getJsonValue("roleName"),
                  loginID: this.localService.getJsonValue("loginID")
            }
            this.spinnerService.loadingMessage = "Getting All Searched Enquiries"


            let data = this.enquirySearchForm.value
            this.utilityService.deleteNullValuePropertyFromObject(data)
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
            this.enquiryService.getAllSearchedEnquiries(data, this.limit, this.offset, roleNameAndLogin).subscribe(response => {
                  this.enquiryList = response.body
                  this.calculateFieldsOfEnquiry(this.enquiryList)
                  this.totalEnquiries = parseInt(response.headers.get('X-Total-Count'))
            }, (error) => {
                  console.error(error)
                  if (error.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                  }
            })
      }

      // Set the limit of maximum experiences required.
      setMaximumExperienceRequired(value: number): void {
            this.enquiryForm.get('maximumExperience').
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

      // Delete search criteria from enquiry search form by search name.
      deleteSearchCriteria(searchName: string): void {
            this.enquirySearchForm.get(searchName).setValue(null)
            this.getAllSearchedEnquiries()
      }

      //*********************************************ALLOCATION FUNCTIONS************************************************************
      // On clicking allocate one enquiry to batch button click.
      onAllocateSalespersonToOneEnquiryClick(enquiryID: any): void {
            this.showAllocateSalespersonToEnquiries = false
            this.showAllocateSalespersonToOneEnquiry = true
            this.assignSelectedEnquiry(enquiryID)
            this.openModal(this.allocateSalespersonModal, 'sm')
      }

      // On clicking allocate multiple enquiries to batch button click.
      onAllocateSalespersonToEnquiriesClick(): void {
            this.showAllocateSalespersonToEnquiries = true
            this.showAllocateSalespersonToOneEnquiry = false
            this.openModal(this.allocateSalespersonModal, 'sm')
      }

      // Toggle visibility og multiple select checkbox.
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
            for (let i = 0; i < this.enquiryList.length; i++) {
                  this.addEnquiryToList(isSelectedAll, this.enquiryList[i])
            }
      }

      // Check if all enquiries in selected enquiries are added in multiple select or not.
      checkEnquiriesAdded(): boolean {
            let count: number = 0
            for (let i = 0; i < this.enquiryList.length; i++) {
                  if (this.selectedEnquiryList.includes(this.enquiryList[i].id))
                        count = count + 1
            }
            return (count == this.enquiryList.length)
      }

      // Check if enquiry is added in multiple select or not.
      checkEnquiryAdded(enquiryID): boolean {
            return this.selectedEnquiryList.includes(enquiryID)
      }

      // Takes a list called selected enquiries and adds all the checked enquiries to list, also does not contain duplicate values.
      addEnquiryToList(isChecked: boolean, enquiry: IEnquiryDTO): void {
            if (isChecked) {
                  if (!this.selectedEnquiryList.includes(enquiry.id)) {
                        this.selectedEnquiryList.push(enquiry.id)
                  }
                  return
            }
            if (this.selectedEnquiryList.includes(enquiry.id)) {
                  let index = this.selectedEnquiryList.indexOf(enquiry.id)
                  this.selectedEnquiryList.splice(index, 1)
            }
      }

      // Allocate salesperson to enquiry(s).
      allocateSalesPersonToEnquiries(salesPersonID: string, enquiryID?: string): void {
            if (salesPersonID == "null") {
                  alert("Please select sales person")
                  return
            }
            let enquiryIDsToBeUpdated = []
            if (!enquiryID) {
                  for (let index = 0; index < this.selectedEnquiryList.length; index++) {
                        enquiryIDsToBeUpdated.push({
                              "enquiryID": this.selectedEnquiryList[index]
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


            this.enquiryService.updateEnquiriesSalesPerson(enquiryIDsToBeUpdated, salesPersonID).subscribe((response: any) => {
                  this.getEnquiries()
                  alert(response)
                  this.modalRef.close('success')
                  this.selectedEnquiryList = []
            }, (error) => {
                  console.error(error)
                  if (typeof error.error == 'object' && error) {
                        alert(this.utilityService.getErrorString(error))
                        return
                  }
                  if (error.error == undefined) {
                        alert('Sales person could not be allocated to enquiries')
                  }
                  alert(error.statusText)
            })
      }

      //*********************************************CRUD FUNCTIONS FOR WAITING LIST************************************************************
      // On clicking get waiting list by enquiry id.
      getWaitingListForSelectedEnquiry(enquiryID: any): void {
            this.isOperationWaitingListUpdate = false
            this.waitingList = []
            this.spinnerService.loadingMessage = "Getting waiting list"
            this.assignSelectedEnquiry(enquiryID)
            this.getWaitingList()
            this.showWaitingListForm = false
            this.openModal(this.waitingListFormModal, 'xl')
      }

      // Get waiting list by enquiry id.
      getWaitingList(): void {
            this.spinnerService.loadingMessage = "Getting Waiting List"


            this.enquiryService.getWaitingListByEnquiry(this.selectedEnquiry.id).subscribe(response => {
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


            this.waitingListForm.get('email').setValue(this.selectedEnquiry.email)
            this.waitingListForm.get('isActive').setValue(true)
            let waitingList: IWaitingList = this.waitingListForm.value
            this.patchIDFromObjectsForWaitingList(waitingList)
            waitingList.enquiryID = this.selectedEnquiry.id
            this.enquiryService.addWaitingList(waitingList).subscribe((response: any) => {
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
            this.enquiryService.updateWaitingList(waitingList).subscribe((response: any) => {
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


                  this.enquiryService.deleteWaitingList(waitingListID).subscribe((response: any) => {
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

      //*********************************************GET FUNCTIONS************************************************************
      // Get all lists.
      getAllComponents(): void {
            this.getPurposeList()
            this.getCountryList()
            this.getEnquiryTypeList()
            this.getDesignnationList()
            this.getTechnologyList()
            this.getSalesPersonList()
            this.getDegreeList()
            this.getAcademicYearList()
            this.getSourceList()
            this.getCollegeBranchList()
            this.getExaminationList()
            this.getYearOfMSList()
            this.getCourseList()
            this.getCompanyBranchList()
            this.getAllBatchList()
            this.getAllRequirementList()
            this.getQueryParams()
      }

      // Get designation list.
      getDesignnationList(): void {
            this.generalService.getDesignations().subscribe((respond: any[]) => {
                  //respond.push(this.other)
                  this.designationList = respond
            }, (err) => {
                  console.error(err)
            })
      }

      // Get purpose list.
      getPurposeList(): void {
            this.generalService.getPurposeListByType("talent_enquiry").subscribe(
                  data => {
                        this.purposeList = data
                  }
            ), err => {
                  console.error(err)
            }
      }

      // Get outcome list by purpose object.
      getOutcomesByPurpose(purpose: any): void {
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

      // Get outcome list by purpose id.
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

      // Get salesPerson list.
      getSalesPersonList(): void {
            this.generalService.getSalesPersonList().subscribe((data: any) => {
                  this.salesPersonList = data.body
            }, (err) => {
                  console.error(err)
            })
      }

      // Get enquiry type list.
      getEnquiryTypeList(): void {
            this.generalService.getGeneralTypeByType("talent_enquiry_type").subscribe((respond: any[]) => {
                  this.enquiryTypeList = respond
            }, (err) => {
                  console.error(err)
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

      // Get state list by country ID.
      getStateListByCountryID(country: any): void {
            if (country == null) {
                  this.stateList = []
                  this.enquiryForm.get('state').setValue(null)
                  this.enquiryForm.get('state').disable()
                  return
            }
            if (country.name == "U.S.A." || country.name == "India") {
                  this.enquiryForm.get('state').enable()
                  this.generalService.getStatesByCountryID(country.id).subscribe((respond: any[]) => {
                        this.stateList = respond
                  }, (err) => {
                        console.error(err)
                  })
            }
            else {
                  this.enquiryForm.get('state').setValue(null)
                  this.enquiryForm.get('state').disable()
            }
      }

      // Get degree list.
      getDegreeList(): void {
            let queryParams: any = {
                  limit: -1,
                  offset: 0,
            }
            this.degreeService.getAllDegrees(queryParams).subscribe((respond: any) => {
                  this.degreeList = respond.body
            }, (err) => {
                  console.error(err)
            })
      }

      // Get specialization list.
      getSpecializationListByDegreeID(degreeID: string, index?: number, shouldClear?: boolean): void {
            if (index != undefined) {
                  this.enquiryAcademicControlArray.at(index).get('specialization').setValue(null)
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
            })
      }

      // Get academic year list.
      getAcademicYearList(): void {
            this.generalService.getGeneralTypeByType("academic_year").subscribe((respond: any[]) => {
                  this.academicYearList = respond
            }, (err) => {
                  console.error(this.utilityService.getErrorString(err))
            })
      }

      // Get technology list.
      getTechnologyList(event?: any): void {
            // this.generalService.getTechnologies().subscribe((respond: any[]) => {
            //       this.technologyList = respond
            // })
            let queryParams: any = {}
            if (event && event?.term != "") {
                  queryParams.language = event?.term
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

      // Get all enquiries by limit and offset.
      getEnquiries(): void {
            let roleNameAndLogin: any = {
                  roleName: this.localService.getJsonValue("roleName"),
                  loginID: this.localService.getJsonValue("loginID")
            }
            if (this.isNavigatedFromWaitingList) {
                  this.getEnquiriesForWaitingList()
                  return
            }
            if (this.isSearched) {
                  this.getAllSearchedEnquiries()
                  return
            }
            this.spinnerService.loadingMessage = "Getting all enquiries"


            this.enquiryService.getEnquiries(this.limit, this.offset, roleNameAndLogin).subscribe(res => {
                  this.enquiryList = res.body
                  this.calculateFieldsOfEnquiry(this.enquiryList)
                  this.totalEnquiries = parseInt(res.headers.get("X-Total-Count"))
            }, err => {
                  console.error(this.utilityService.getErrorString(err))
            }).add(() => {
                  this.setPaginationString()
            })
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
                  // console.log("getCollegeBranchList -> ", response);
                  this.collegeBranchList = []
                  this.collegeBranchList = this.collegeBranchList.concat(response)
            }, (err) => {
                  console.error(err)
            }).add(() => {
                  this.isCollegeLoading = false
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

      getYearOfMSList(): void {
            this.yearOfMSList.push(this.currentYear)
            this.yearOfMSList.push((this.currentYear + 1))
            this.yearOfMSList.push((this.currentYear + 2))
      }

      // Get course list.
      getCourseList(): void {
            this.generalService.getCourseList().subscribe(
                  response => {
                        this.courseList = response.body
                  }
            ), err => {
                  console.error(err)
            }
      }

      // Get company branch list.
      getCompanyBranchList(): void {
            this.generalService.getCompanyBranchList().subscribe((response: any) => {
                  this.companyBranchList = response.body
            }, (err: any) => {
                  console.error(err)
            })
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

      // Get query params from url and get enquiries accordingly.
      getQueryParams(): void {
            if (this.activatedRoute.snapshot.queryParamMap.keys != []) {
                  this.enquiryList = []
                  // Navigated from waiting list by company branch id.
                  let waitingListCompanyBranchID = this.activatedRoute.snapshot.queryParamMap.get("waitingListCompanyBranchID")
                  if (waitingListCompanyBranchID) {
                        this.isNavigatedFromWaitingList = true
                        this.waitingListCompanyBranchID = waitingListCompanyBranchID
                        this.getEnquiriesForWaitingList()
                        return
                  }

                  // Navigated from waiting list by course id.
                  let waitingListCourseID = this.activatedRoute.snapshot.queryParamMap.get("waitingListCourseID")
                  if (waitingListCourseID) {
                        this.isNavigatedFromWaitingList = true
                        this.waitingListCourseID = waitingListCourseID
                        this.getEnquiriesForWaitingList()
                        return
                  }

                  // Navigated from waiting list by batch id.
                  let waitingListBatchID = this.activatedRoute.snapshot.queryParamMap.get("waitingListBatchID")
                  if (waitingListBatchID) {
                        this.isNavigatedFromWaitingList = true
                        this.waitingListBatchID = waitingListBatchID
                        this.getEnquiriesForWaitingList()
                        return
                  }

                  // Navigated from waiting list by requirement id.
                  let waitingListRequirementID = this.activatedRoute.snapshot.queryParamMap.get("waitingListRequirementID")
                  if (waitingListRequirementID) {
                        this.isNavigatedFromWaitingList = true
                        this.waitingListRequirementID = waitingListRequirementID
                        this.getEnquiriesForWaitingList()
                        return
                  }

                  // Navigated from waiting list by requirement id.
                  let waitingListTechnologyID = this.activatedRoute.snapshot.queryParamMap.get("waitingListTechnologyID")
                  if (waitingListTechnologyID) {
                        this.isNavigatedFromWaitingList = true
                        this.waitingListTechnologyID = waitingListTechnologyID
                        this.getEnquiriesForWaitingList()
                        return
                  }

            }
            this.changePage(1)
      }

      // Get enquiries for waiting list.
      getEnquiriesForWaitingList(): void {
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
            this.spinnerService.loadingMessage = "Getting Enquiries"


            this.enquiryService.getEnquiriesByWaitingList(this.limit, this.offset, queryParams).
                  subscribe((response: any) => {
                        this.enquiryList = response.body
                        this.totalEnquiries = response.headers.get('X-Total-Count')
                        this.calculateFieldsOfEnquiry(this.enquiryList)
                  },
                        (error) => {
                              console.error(error)
                              if (error.error) {
                                    alert(error.error)
                                    return
                              }
                              alert(error.error.error)
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
            this.docStatus = ""
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

            this.fileOperationService.uploadExcel(file).subscribe((uploadedEnquiries: any) => {
                  if (this.validateEnquiries(uploadedEnquiries)) {
                        this.excelUploadedEnquiries = uploadedEnquiries
                        this.isExcelUploaded = true
                        this.docStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
                        this.reformEnquiriesFromExcel()
                  } else {
                        excelFile.value = ""
                        this.docStatus = ""
                        this.isExcelUploaded = false
                  }
            }, (err) => {
                  alert(err)
                  excelFile.value = ""
            }).add(() => {

            })
      }

      // Validate all enquiries from excel file upload.
      validateEnquiries(uploadedEnquiries: IEnquiryExcel[]): boolean {

            if (!uploadedEnquiries || uploadedEnquiries.length == 0) {
                  alert("No Uploaded Enquiries")
                  return false
            }

            for (let index = 0; index < uploadedEnquiries.length; index++) {
                  if (!uploadedEnquiries[index].firstName || uploadedEnquiries[index].firstName == "") {
                        alert(`firstname on row ${index + 2} is not specified`)
                        return false
                  }
                  if (!uploadedEnquiries[index].lastName || uploadedEnquiries[index].lastName == "") {
                        alert(`lastname on row ${index + 2} is not specified`)
                        return false
                  }
                  if (!uploadedEnquiries[index].email || uploadedEnquiries[index].email == "") {
                        alert(`email on row ${index + 2} is not specified`)
                        return false
                  }
                  if (!uploadedEnquiries[index].contact || uploadedEnquiries[index].contact == "") {
                        alert(`contact on row ${index + 2} is not specified`)
                        return false
                  }
                  if (!uploadedEnquiries[index].academicYear || uploadedEnquiries[index].academicYear == 0) {
                        alert(`academic year on row ${index + 2} is not specified`)
                        return false
                  }

                  // Check for academic details.
                  if (uploadedEnquiries[index].degreeName || uploadedEnquiries[index].specializationName ||
                        uploadedEnquiries[index].collegeName || uploadedEnquiries[index].percentage ||
                        uploadedEnquiries[index].yearOfPassout) {

                        if (!uploadedEnquiries[index].degreeName || uploadedEnquiries[index].degreeName == "") {
                              alert(`degree  on row ${index + 2} is not specified`)
                              return false
                        }

                        if (!uploadedEnquiries[index].specializationName || uploadedEnquiries[index].specializationName == "") {
                              alert(`specialization  on row ${index + 2} is not specified`)
                              return false
                        }

                        if (!uploadedEnquiries[index].collegeName || uploadedEnquiries[index].collegeName == "") {
                              alert(`college  on row ${index + 2} is not specified`)
                              return false
                        }

                        if (!uploadedEnquiries[index].percentage || uploadedEnquiries[index].percentage == 0) {
                              alert(`percentage  on row ${index + 2} is not specified`)
                              return false
                        }

                        if (!uploadedEnquiries[index].yearOfPassout || uploadedEnquiries[index].yearOfPassout == 0) {
                              alert(`passout  on row ${index + 2} is not specified`)
                              return false
                        }
                  }
            }
            return true
      }

      // Create enquiries json to be sent to API.
      reformEnquiriesFromExcel(): void {
            this.reformedEnquiries = []
            for (let index = 0; index < this.excelUploadedEnquiries.length; index++) {

                  let enquiry: any = {
                        firstName: this.excelUploadedEnquiries[index].firstName,
                        lastName: this.excelUploadedEnquiries[index].lastName,
                        email: this.excelUploadedEnquiries[index].email,
                        contact: this.excelUploadedEnquiries[index].contact.toString(),
                        academicYear: this.excelUploadedEnquiries[index].academicYear,
                        city: this.excelUploadedEnquiries[index].city,
                        address: this.excelUploadedEnquiries[index].address,
                        pinCode: this.excelUploadedEnquiries[index].pinCode,
                        countryName: this.excelUploadedEnquiries[index].countryName,
                        stateName: this.excelUploadedEnquiries[index].stateName,
                        academics: [{
                              degreeName: this.excelUploadedEnquiries[index].degreeName,
                              specializationName: this.excelUploadedEnquiries[index].specializationName,
                              collegeName: this.excelUploadedEnquiries[index].collegeName,
                              percentage: this.excelUploadedEnquiries[index].percentage,
                              yearOfPassout: this.excelUploadedEnquiries[index].yearOfPassout,
                        }],
                  }
                  this.reformedEnquiries.push(enquiry)
                  this.enquiriesExcelTotalCount = this.reformedEnquiries.length
            }
            this.disableAddEnquiriesButton = false
            // console.log(this.reformedTalents)
      }

      // Add multiple enquiries from excel.
      addMultipleEnquiries(): void {

            this.showExcelProgress = true
            this.disableAddEnquiriesButton = true

            // Call add enquiry from excel API for every enquiry.
            for (let i = 0; i < this.enquiriesExcelTotalCount; i++) {
                  this.enquiryService.addEnquiryFromExcel(this.reformedEnquiries[i]).subscribe((response: any) => {
                        this.enquiriesExcelAddedCount = this.enquiriesExcelAddedCount + 1
                  }, (error) => {
                        if (error?.error?.error) {
                              this.excelErrorList.push(this.reformedEnquiries[i].email + " : " + error?.error?.error)
                              return
                        }
                        this.excelErrorList.push(this.reformedEnquiries[i].email + " : " + error.error)
                  }).add(() => {
                        this.enquiriesExcelProcessedCount = this.enquiriesExcelProcessedCount + 1
                        if (this.enquiriesExcelProcessedCount == this.enquiriesExcelTotalCount) {
                              this.getEnquiries()
                        }
                  })
            }
      }

      // On cancelling uploaded excel file.
      onCalcelUploadedExcelFile(): void {
            this.docStatus = ""
            this.isExcelUploaded = false
            this.addMultipleErrorList = []
      }

      // // Add multiple enquiries from excel.
      // addMultipleEnquiries(): void {
      //       this.spinnerService.loadingMessage = "Adding enquiries"
      //       
      //       
      //       this.enquiryService.addEnquiries(this.reformedEnquiries).subscribe((response: any) => {
      //             // console.log(response)
      //             this.modalRef.close('success')
      //             alert(response)
      //             this.getEnquiries()
      //       }, (error) => {
      //             this.docStatus = ""
      //             this.isExcelUploaded = false
      //             console.error(error)
      //             if (error.error?.error) {
      //                   alert(error.error?.error)
      //                   return
      //             }
      //             this.addMultipleErrorList = []
      //             if (error?.error?.errorList?.length > 0) {
      //                   this.excelEnquiriesAdded = error?.error?.message
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

