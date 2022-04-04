import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, FormControl, Validators, FormArray } from '@angular/forms';
import { GeneralService, ICountry, ISalaryTrend, IState } from 'src/app/service/general/general.service';
import { CompanyService, ICompanyRequirement, IDomain } from 'src/app/service/company/company.service';
import { Constant, Role, UrlConstant } from 'src/app/service/constant';
import { ActivatedRoute, Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { IPermission } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { DatePipe, Location } from '@angular/common';
import { ISearchFilterField, ISearchSection, TalentService } from 'src/app/service/talent/talent.service';
import { CompanyRequirementService, IRequirement, IRequirementDTO, ISearchTalentParams } from 'src/app/service/company/company-requirement/company-requirement.service';
import { requirementDateValidator } from 'src/app/Validators/custom.validators';
import { DummyDataService } from 'src/app/service/dummydata/dummy-data.service';
import { ITechnology, TechnologyService } from 'src/app/service/technology/technology.service';
import { DegreeService } from 'src/app/service/degree/degree.service';


@Component({
      selector: 'app-company-master-requirement',
      templateUrl: './company-master-requirement.component.html',
      styleUrls: ['./company-master-requirement.component.css']
})

export class CompanyMasterRequirementComponent implements OnInit {

      // Components.
      technologyList: ITechnology[]
      salesPersonList: any[]
      stateList: any[]
      countryList: any[]
      designationList: any[]
      degreeList: any[]
      personalityTypeList: any[]
      talentRatingList: any[]
      activeRequirementList: any[]
      activeBatchList: any[]
      isTechLoading: boolean
      companyList: any[]

      // Required Componenets/Fields.
      allDomains: IDomain[];
      // allTechnologies: ITechnology[];
      // salesPeople: any[];
      // allStates: IState[];
      // allCountries: ICountry[];
      // ratingList: number[];

      // Flags.
      isOperationUpdate: boolean
      isViewMode: boolean
      showExperienceFields: boolean

      // Company requirement.
      companyRequirement: IRequirementDTO[]
      companyRequirementForm: FormGroup
      selectedRequirement: any
      isRequrimentLoaded: boolean

      // Pagination.
      limit: number
      currentPage: number
      offset: number
      paginationString: string
      totalCompanyRequirements: number
      techLimit: number
      techOffset: number


      // Modal.
      modalRef: any
      @ViewChild('companyRequirementFormModal') companyRequirementFormModal: any
      @ViewChild('companyRequirementModal') companyRequirementModal: any
      @ViewChild('deleteCompanyRequirementModal') deleteCompanyRequirementModal: any
      @ViewChild('closeCompanyRequirementModal') closeCompanyRequirementModal: any
      @ViewChild('allocateSalespersonModal') allocateSalespersonModal: any
      @ViewChild('searchTalentsModal') searchTalentsModal: any
      @ViewChild('showApplicantsModal') showApplicantsModal: any
      @ViewChild('addRequirementRatingModal') addRequirementRatingModal: any
      @ViewChild('drawer') drawer: any

      // Spinner.



      // Search.
      isSearched: boolean
      showSearch: boolean
      searchFormValue: any
      searchParams: ISearchTalentParams
      companyRequirementSearchForm: FormGroup
      selectedSectionName: string
      searchSectionList: ISearchSection[]

      // Permission.
      permission: IPermission
      roleName: string
      showForAdmin: boolean
      showForSalesPerson: boolean
      showForCompany: boolean

      // Terms and conditions.
      isTermsUploadedToServer: boolean
      isTermsFileUploading: boolean
      termsDocStatus: string
      termsDisplayedFileName: string

      // sample offer letter.
      isOfferLetterUploadedToServer: boolean
      isOfferLetterFileUploading: boolean
      offerLetterDocStatus: string
      offerDisplayedFileName: string

      // Company requirement relatd numbers.
      minValueForMaxExp: number

      // Constants
      nilUUID: string

      // Waiting list.
      twoWaitingLists: any
      selecetdrRequirementForTransfer: any
      selectedBatchForTransfer: any
      showTransferButton: boolean
      showTransferForm: boolean
      showCompanyReqField: boolean
      showBatchField: boolean

      // requirement rating
      requirementRatingForm: FormGroup
      isRatingAddClick: boolean

      // feedback questions
      feedbackQuestions: any[]
      fixedRatingQuestion: any[]

      // salary-trend
      salaryTrend: ISalaryTrend

      // navigation
      isNavigated: boolean

      // quill editor
      // quillModuleConfig: any
      // quillEditorStyle: any
      // quillCommentEditorStyle: any

      // ck-editor
      ckConfig: any

      //constant
      private readonly IGNORE_SEARCH_FIELD: string[] = ["limit", "offset"]

      // requirement terms and conditions
      isRequirementTermsUploadedToServer: boolean
      isRequirementTermsFileUploading: boolean
      requirementTermsDocStatus: string
      requirementTermsDisplayedFileName: string
      isSalesPerson: boolean;
      isAdmin: boolean;

      constructor(
            private formBuilder: FormBuilder,
            private companyRequirementService: CompanyRequirementService,
            private generalService: GeneralService,
            private companyService: CompanyService,
            private techService: TechnologyService,
            private localService: LocalService,
            private utilService: UtilityService,
            private talentService: TalentService,
            private degreeService: DegreeService,
            private fileOperationService: FileOperationService,
            private dummyDataService: DummyDataService,
            private constant: Constant,
            private urlConstant: UrlConstant,
            private spinnerService: SpinnerService,
            private modalService: NgbModal,
            private datePipe: DatePipe,
            private role: Role,
            private router: Router,
            private route: ActivatedRoute,
            private location: Location,
            private fileOps: FileOperationService,
            private roleConstant: Role,


      ) {
            this.initializeVariables()
            this.getAllComponents()
      }


      get ongoingOperations() {
            return this.spinnerService.ongoingOperations
      }

      ngOnInit() { }

      initializeVariables(): void {
            this.isSalesPerson = this.localService.getJsonValue("roleName") == this.roleConstant.SALES_PERSON
            this.isAdmin = this.localService.getJsonValue("roleName") == this.roleConstant.ADMIN
            // Components.
            this.technologyList = []
            this.salesPersonList = []
            this.stateList = []
            this.countryList = []
            this.designationList = []
            this.degreeList = []
            this.personalityTypeList = []
            this.talentRatingList = []
            this.feedbackQuestions = []

            // Flags.
            this.isOperationUpdate = false
            this.isViewMode = false
            this.showExperienceFields = false
            this.isRatingAddClick = false
            this.isRequrimentLoaded = true
            this.isTechLoading = false

            // Company requirement.
            this.selectedRequirement = {}

            // Search.
            this.isSearched = false
            this.showSearch = false
            this.searchFormValue = {}

            // Pagination.
            this.limit = 5
            this.offset = 0
            this.currentPage = 0
            this.techLimit = 10
            this.techOffset = 0

            // Permision.
            this.showForAdmin = false
            this.showForSalesPerson = true
            this.showForCompany = true

            // Get permissions from menus using utilityService function.
            this.permission = this.utilService.getPermission(this.urlConstant.COMPANY_REQUIREMENT)

            // Get role name for menu for calling their specific apis.
            this.roleName = this.localService.getJsonValue("roleName")

            // If admin is logged in then show its features.
            if (this.roleName == this.role.ADMIN) {
                  this.showForAdmin = true
            }
            // Hide features for salesperson.
            if (this.roleName == this.role.SALES_PERSON) {
                  this.showForSalesPerson = false
            }
            // Hide features for faculty.
            if (this.roleName == this.role.COMPANY) {
                  this.showForCompany = false
            }

            // Terms and conditions.
            this.isTermsUploadedToServer = false
            this.isTermsFileUploading = false
            this.termsDisplayedFileName = "Select file"
            this.termsDocStatus = ""

            // sample offer letter
            this.isOfferLetterUploadedToServer = false
            this.isOfferLetterFileUploading = false
            this.offerDisplayedFileName = "Select File"
            this.offerLetterDocStatus = ""

            // Constants.
            this.nilUUID = this.constant.NIL_UUID

            // Waiting list.
            this.twoWaitingLists = {}

            // Company requirement related numbers.
            this.minValueForMaxExp = 1

            // Initialize forms
            this.createCompanyRequirementSearchForm()

            this.setSearchSectionFields()
            this.editorConfig()

            // Spinner.
            this.spinnerService.loadingMessage = "Getting Company Requirements"
            // requirement terms and conditions.
            this.isRequirementTermsUploadedToServer = false
            this.isRequirementTermsFileUploading = false
            this.requirementTermsDisplayedFileName = "Select file"
      }

      setSearchSectionFields(): void {
            this.searchSectionList = [
                  {
                        name: "Requirement",
                        isSelected: true
                  },
                  {
                        name: "Talent",
                        isSelected: false
                  },
                  {
                        name: "Location",
                        isSelected: false
                  },
                  {
                        name: "Other",
                        isSelected: false
                  },
            ]
            this.selectedSectionName = "Requirement"
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

      editorConfig(): void {
            this.ckConfig = {
                  extraPlugins: 'codeTag',
                  removePlugins: "exportpdf",
                  toolbar: [
                        { name: 'styles', items: ['Styles', 'Format'] },
                        {
                              name: 'basicstyles', groups: ['basicstyles', 'cleanup'],
                              items: ['Bold', 'Italic', 'Underline', 'Strike', 'Subscript', 'Superscript', 'RemoveFormat', 'Code']
                        },
                        {
                              name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'],
                              items: ['NumberedList', 'BulletedList', '-', 'Outdent', 'Indent', '-', 'Blockquote']
                        },
                        { name: 'document', groups: ['mode', 'document', 'doctools'], items: ['Source'] },
                  ],
                  toolbarGroups: [
                        { name: 'styles' },
                        { name: 'basicstyles', groups: ['basicstyles', 'cleanup'] },
                        { name: 'document', groups: ['mode', 'document', 'doctools'] },
                        { name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'] },
                        { name: 'links' },
                        // { name: 'insert' },
                  ],
                  removeButtons: "",
                  resize_enabled: false,
                  language: "en",
                  width: "100%",
                  height: "80%",
                  forcePasteAsPlainText: false,
            }
      }

      // quillEditorConfiguration(): void {
      //       // Quill config.
      //       this.quillModuleConfig = {
      //             // **** Comments include all option. Please do not remove them.
      //             toolbar: [
      //                   ['bold', 'italic', 'underline', 'strike'],        // toggled buttons
      //                   [{ 'list': 'ordered' }, { 'list': 'bullet' }],
      //                   [{ 'indent': '-1' }, { 'indent': '+1' }],          // outdent/indent
      //                   [{ 'header': [1, 2, 3, 4, 5, 6, false] }],
      //                   [{ 'align': [] }],
      //                   ['clean']                                 // remove formatting button
      //             ],
      //       }

      //       this.quillEditorStyle = {
      //             height: '100px',
      //             fontSize: 'large',
      //             backgroundColor: '#ffffff',
      //       }

      //       this.quillCommentEditorStyle = {
      //             height: '80px',
      //             fontSize: 'large',
      //             backgroundColor: '#ffffff'
      //       }

      // }

      // =============================================================CREATE FORMS==========================================================================
      // Create company form.
      createCompanyEnquiryRequirementForm(): void {
            this.companyRequirementForm = this.formBuilder.group({
                  // Address.
                  id: new FormControl(null),
                  jobLocation: new FormGroup({
                        address: new FormControl(null, [Validators.required]),
                        pinCode: new FormControl(null, [Validators.required, Validators.pattern("^[0-9]{6}$")]),
                        country: new FormControl(null, [Validators.required]),
                        state: new FormControl(null, [Validators.required]),
                        city: new FormControl(null, [Validators.required, Validators.maxLength(50),
                        Validators.pattern("^[a-zA-Z]+([a-zA-Z ]?)+")]),
                  }),

                  // Multiple.
                  qualifications: new FormControl(null),
                  universities: new FormControl(null),
                  technologies: new FormControl(null, [Validators.required]),
                  selectedTalents: new FormControl(null),
                  designation: new FormControl(null, [Validators.required]),
                  designationID: new FormControl(null),

                  // Objects.
                  company: new FormControl(null),
                  companyID: new FormControl(null),
                  salesPersonID: new FormControl(null),

                  // Single.
                  code: new FormControl(null),
                  talentRating: new FormControl(null),
                  personalityType: new FormControl(null),
                  isExperience: new FormControl(false, [Validators.required]),
                  // jobRole: new FormControl(null, [Validators.required]),
                  jobDescription: new FormControl(null, [Validators.required]),
                  jobRequirement: new FormControl(null, [Validators.required]),
                  jobType: new FormControl(null, [Validators.required]),
                  // packageOffered: new FormControl(null, [Validators.required, Validators.min(100000), Validators.max(1000000000)]),
                  minimumPackage: new FormControl(null, [Validators.required, Validators.min(100000)]),
                  maximumPackage: new FormControl(null, [Validators.required, Validators.max(100000000000)]),
                  requiredBefore: new FormControl(null, [Validators.required, requirementDateValidator]),
                  requiredFrom: new FormControl(null, [Validators.required, requirementDateValidator]),
                  vacancy: new FormControl(null, [Validators.required, Validators.min(1), Validators.max(10000)]),
                  comment: new FormControl(null, [Validators.maxLength(4000)]),
                  termsAndConditions: new FormControl(null),
                  sampleOfferLetter: new FormControl(null),
                  isActive: new FormControl(true),

                  // Undecided fields.
                  // colleges: new FormControl(null),
                  // marksCriteria: new FormControl(null, [Validators.max(100)]),
            });
      }


      // Create company requirement form.
      createCompanyRequirementForm(): void {
            this.companyRequirementForm = this.formBuilder.group({
                  // Address.
                  id: new FormControl(null),
                  jobLocation: new FormGroup({
                        address: new FormControl(null, [Validators.required]),
                        pinCode: new FormControl(null, [Validators.required, Validators.pattern("^[0-9]{6}$")]),
                        country: new FormControl(null, [Validators.required]),
                        state: new FormControl(null, [Validators.required]),
                        city: new FormControl(null, [Validators.required, Validators.maxLength(50),
                        Validators.pattern("^[a-zA-Z]+([a-zA-Z ]?)+")]),
                  }),

                  // Multiple.
                  qualifications: new FormControl(null),
                  universities: new FormControl(null),
                  technologies: new FormControl(null, [Validators.required]),
                  selectedTalents: new FormControl(null),

                  // Objects.
                  company: new FormControl(null),
                  salesPerson: new FormControl(null),
                  designation: new FormControl(null, [Validators.required]),
                  designationID: new FormControl(null),

                  // Single.
                  code: new FormControl(null),
                  talentRating: new FormControl(null),
                  personalityType: new FormControl(null),
                  isExperience: new FormControl(false, [Validators.required]),
                  // jobRole: new FormControl(null, [Validators.required]),
                  jobDescription: new FormControl(null, [Validators.required]),
                  jobRequirement: new FormControl(null, [Validators.required]),
                  jobType: new FormControl(null, [Validators.required]),
                  minimumPackage: new FormControl(null, [Validators.required, Validators.min(100000)]),
                  maximumPackage: new FormControl(null, [Validators.required, Validators.max(100000000000)]),
                  // packageOffered: new FormControl(null, , Validators.max(100000000000)]),
                  requiredBefore: new FormControl(null, [Validators.required, requirementDateValidator]),
                  requiredFrom: new FormControl(null, [Validators.required, requirementDateValidator]),
                  vacancy: new FormControl(null, [Validators.required, Validators.min(1), Validators.max(10000)]),
                  comment: new FormControl(null, [Validators.maxLength(4000)]),
                  termsAndConditions: new FormControl(null),
                  sampleOfferLetter: new FormControl(null),
                  // minimumExperience: new FormControl(null, [Validators.min(0), Validators.max(30)]),
                  // maximumExperience: new FormControl(null, [Validators.min(0), Validators.max(30)]),

                  // Criteria
                  increment: new FormControl(null),
                  weeklyHoliday: new FormControl(null),
                  qualification: new FormControl(null),
                  bondPeriod: new FormControl(null),
                  talentLocation: new FormControl(null),
                  workShift: new FormControl(null),
                  companyType: new FormControl(null),
                  joiningPeriod: new FormControl(null),
                  genderPreference: new FormControl(null),
                  // criteria10: new FormControl(null),

                  // Rating
                  rating: new FormControl(null)

                  // Undecided fields.
                  // colleges: new FormControl(null),
                  // marksCriteria: new FormControl(null, [Validators.max(100)]),
            })
      }


      // Create company requirement search form.
      createCompanyRequirementSearchForm(): void {
            this.companyRequirementSearchForm = this.formBuilder.group({
                  countryID: new FormControl(null),
                  stateID: new FormControl({ value: null, disabled: true }),
                  city: new FormControl(null, [Validators.maxLength(50), Validators.pattern("^[a-zA-Z]+([a-zA-Z ]?)+")]),
                  salesPersonID: new FormControl(null),
                  technologies: new FormControl(null),
                  qualifications: new FormControl(null),
                  designation: new FormControl(null),
                  requirementFromDate: new FormControl(null),
                  requirementTillDate: new FormControl(null),
                  minimumExperience: new FormControl(null, [Validators.min(0), Validators.max(30)]),
                  maximumExperience: new FormControl(null, [Validators.min(0), Validators.max(30)]),
                  // jobRole: new FormControl(null),
                  personalityType: new FormControl(null),
                  minimumPackage: new FormControl(null, [Validators.min(100000), Validators.max(1000000000)]),
                  maximumPackage: new FormControl(null, [Validators.min(100000), Validators.max(1000000000)]),
                  isActive: new FormControl(null),
                  talentRating: new FormControl(null),
                  limit: new FormControl(this.limit),
                  offset: new FormControl(this.offset),
            })
      }

      // Nagivate to previous url.
      backToPreviousPage(): void {
            this.location.back()
      }

      // =============================================================COMPANY REQUIREMENT CRUD FUNCTIONS==========================================================================
      // On clicking view company requirement button.
      onViewCompanyRequirementClick(companyRequirement?: IRequirementDTO): void {

            this.isViewMode = true
            this.showExperienceFields = false
            this.isRatingAddClick = false
            this.createCompanyRequirementForm()

            console.log("companyRequirement.requiredFrom -> ", companyRequirement.requiredFrom);
            console.log("companyRequirement.requiredBefore -> ", companyRequirement.requiredBefore);

            // Format dates.
            // companyRequirement.requiredFrom = this.datePipe.transform(companyRequirement.requiredFrom, 'yyyy-MM-dd')
            // companyRequirement.requiredBefore = this.datePipe.transform(companyRequirement.requiredBefore, 'yyyy-MM-dd')

            // Experience.
            if (companyRequirement.minimumExperience) {
                  this.companyRequirementForm.get('isExperience').setValue(true)
                  this.onIsExperiencedValueChange(true)
            }

            // State.
            if (companyRequirement.jobLocation.country != undefined) {
                  this.getStateListByCountry(companyRequirement.jobLocation.country)
            }

            this.companyRequirementForm.patchValue(companyRequirement)
            this.companyRequirementForm.get("requiredFrom").setValue(
                  this.datePipe.transform(new Date(companyRequirement.requiredFrom).toUTCString(), 'yyyy-MM-dd', "GMT")
            )
            this.companyRequirementForm.get("requiredBefore").setValue(
                  this.datePipe.transform(new Date(companyRequirement.requiredBefore).toUTCString(), 'yyyy-MM-dd', "GMT")
            )
            this.companyRequirementForm.disable()
            this.openModal(this.companyRequirementFormModal, 'xl')
      }

      // On cliking update form button in company requirement form.
      onUpdateCompanyRequirementClick(): void {
            this.isViewMode = false
            this.isOperationUpdate = true
            this.isRatingAddClick = false
            if (this.companyRequirementForm.get('termsAndConditions')?.value) {
                  this.termsDisplayedFileName = `<a href=${this.companyRequirementForm.get('termsAndConditions')?.value} 
                                                      target="_blank">Terms and conditions</a>`
            }
            if (this.companyRequirementForm.get('sampleOfferLetter')?.value) {
                  this.offerDisplayedFileName = `<a href=${this.companyRequirementForm.get('sampleOfferLetter')?.value} 
                                                      target="_blank">Sample offer letter</a>`
            }
            this.enableCompanyRequirementForm()
      }

      onRequirementAddClick(): void {
            if (this.isAdmin) {
                  this.showForAdmin = true
            }
            this.createCompanyEnquiryRequirementForm()
            this.openModal(this.companyRequirementModal, 'xl')

      }
      // Update company requirement.
      updateCompanyRequirement(): void {
            this.spinnerService.loadingMessage = "Updating Company Requirement"
            if (this.isRatingAddClick) {
                  this.spinnerService.loadingMessage = "Adding Requirement Rating"
            }


            let companyRequirement: IRequirement = this.companyRequirementForm.value
            this.setFormFields(companyRequirement)
            // console.log(companyRequirement);
            this.companyRequirementService.updateCompanyRequirement(companyRequirement).subscribe((response: any) => {
                  this.modalRef.close()
                  this.getAllRequirements()
                  if (this.isRatingAddClick) {
                        alert("Rating successfully added.")
                        this.isRatingAddClick = false
                  } else {
                        alert(response)
                  }
            }, (error) => {
                  console.error(error)
                  if (error.error?.error) {
                        alert(error.error.error)
                        return
                  }
                  alert("Check connection")
            })
      }

      // On clicking close company requirement button. 
      onCloseCompanyRequirementClick(companyRequirementID: string): void {
            this.openModal(this.closeCompanyRequirementModal, 'md').result.then(() => {
                  this.closeCompanyRequirement(companyRequirementID)
            }, (err) => {
                  console.error(err)
                  return
            })
      }

      // On clicking close company requirement.
      closeCompanyRequirement(companyRequirementID?: string): void {
            this.spinnerService.loadingMessage = "Closing Company Requirement"


            this.companyRequirementService.closeCompanyRequirement(companyRequirementID).subscribe((response: any) => {
                  this.modalRef.close()
                  this.getAllRequirements()
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

      // On clicking delete company requirement button. 
      onDeleteCompanyRequirementClick(companyRequirementID: string): void {
            this.openModal(this.deleteCompanyRequirementModal, 'md').result.then(() => {
                  this.deleteCompanyRequirement(companyRequirementID)
            }, (err) => {
                  console.error(err)
                  return
            })
      }

      // Delete company requirement after confirmation from user.
      deleteCompanyRequirement(companyRequirementID: string): void {
            this.spinnerService.loadingMessage = "Deleting Company Requirement"


            this.companyRequirementService.deleteCompanyRequirement(companyRequirementID).subscribe((response: any) => {
                  this.modalRef.close()
                  this.getAllRequirements()
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

      // =============================================================COMPANY REQUIREMENT SEARCH FUNCTIONS==========================================================================
      // Reset search form and renaviagte page.
      resetSearchAndGetAll(): void {
            this.companyRequirementSearchForm.reset()
            this.searchFormValue = null
            this.changePage(1)
            this.isSearched = false
            this.showSearch = false
            this.router.navigate([this.urlConstant.COMPANY_REQUIREMENT])
      }

      // Reset search form.
      resetSearchForm(): void {
            this.limit = this.companyRequirementSearchForm.get("limit").value
            this.offset = this.companyRequirementSearchForm.get("offset").value
            this.companyRequirementSearchForm.reset({
                  limit: this.limit,
                  offset: this.offset,
            })
      }

      searchAndCloseDrawer(): void {
            this.drawer.toggle()
            this.searchCompanyRequirements()
      }

      // Search company requirements.
      searchCompanyRequirements(): void {
            this.searchFormValue = { ...this.companyRequirementSearchForm?.value }
            if (this.searchFormValue.departmentID) {
                  this.searchFormValue.departmentID = this.searchFormValue.departmentID.id
            }
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
            this.spinnerService.loadingMessage = "Searching Company Requirements"
            // this.changePage(1)
            this.getAllRequirements()
      }

      // ================================================OTHER FUNCTIONS FOR COMPANY REQUIREMENTS===============================================
      // On page change.
      changePage(pageNumber: number): void {

            this.companyRequirementSearchForm.get("offset").setValue(pageNumber - 1)

            this.limit = this.companyRequirementSearchForm.get("limit").value
            this.offset = this.companyRequirementSearchForm.get("offset").value
            this.searchCompanyRequirements()
            // this.currentPage = pageNumber
            // this.offset = this.currentPage - 1
            // this.getAllRequirements()
      }

      // Checks the url's query params and decides whether to call get or search.
      searchOrGetCompanyRequirements(): void {
            let queryParams = this.route.snapshot.queryParams
            if (this.utilService.isObjectEmpty(queryParams)) {
                  this.getAllRequirements()
                  return
            }
            let isFresherSummary = this.route.snapshot.queryParamMap.get("isFresherSummary")
            if (isFresherSummary) {
                  this.getFresherSummaryQueryparams()
                  return
            }
            this.companyRequirementSearchForm.patchValue(queryParams)
            this.searchCompanyRequirements()
      }

      // Navigated from fresher summary report.
      getFresherSummaryQueryparams() {
            let talentRating = this.route.snapshot.queryParamMap.get("talentRating") // -> outstanding, excellent, average
            let technologies = this.route.snapshot.queryParamMap.get("technologies") // -> techID
            this.isNavigated = true
            this.searchFormValue = {}
            if (technologies) {
                  this.searchFormValue.technologies = technologies
                  // this.companyRequirementSearchForm.get('technologies').setValue(technologies)
            }
            if (talentRating) {
                  this.searchFormValue.talentRating = talentRating
                  // this.companyRequirementSearchForm.get('talentRating').setValue(talentRating)
            }
            this.getAllRequirements()
            return
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

      get jobLocation() {
            return this.companyRequirementForm.get('jobLocation') as FormGroup
      }

      // Set total company requirements list on current page.
      setPaginationString() {
            this.paginationString = ''

            let limit = this.companyRequirementSearchForm.get('limit').value
            let offset = this.companyRequirementSearchForm.get('offset').value

            let start: number = limit * offset + 1
            let end: number = +limit + limit * offset

            if (this.totalCompanyRequirements < end) {
                  end = this.totalCompanyRequirements
            }
            if (this.totalCompanyRequirements == 0) {
                  this.paginationString = ''
                  return
            }
            this.paginationString = `${start} - ${end}`
      }

      // Used to open modal.
      openModal(content: any, size?: string): NgbModalRef {
            this.isTermsUploadedToServer = false
            this.termsDisplayedFileName = "Select file"
            this.termsDocStatus = ""

            this.isOfferLetterUploadedToServer = false
            this.isOfferLetterFileUploading = false
            this.offerDisplayedFileName = "Select File"
            this.offerLetterDocStatus = ""

            if (!size) {
                  size = 'lg'
            }

            let options: NgbModalOptions = {
                  ariaLabelledBy: 'modal-basic-title', keyboard: false,
                  backdrop: 'static', size: size, centered: true
            }
            this.modalRef = this.modalService.open(content, options)
            return this.modalRef
      }

      // On uplaoding one pager.
      onTermsAndConditionSelect(event: any) {
            this.termsDocStatus = ""
            let files = event.target.files
            if (files && files.length) {
                  let file = files[0]
                  let err = this.fileOperationService.isDocumentFileValid(file)
                  if (err != null) {
                        this.termsDocStatus = `<p><span>&#10060;</span> ${err}</p>`
                        return
                  }
                  // Upload terms and condition if it is present.]
                  this.isTermsFileUploading = true
                  this.fileOperationService.uploadTermsAndCondition(file).subscribe((data: any) => {
                        this.companyRequirementForm.markAsDirty()
                        this.companyRequirementForm.patchValue({
                              termsAndConditions: data
                        })
                        this.termsDisplayedFileName = file.name
                        this.isTermsFileUploading = false
                        this.isTermsUploadedToServer = true
                        this.termsDocStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
                  }, (error) => {
                        this.isTermsFileUploading = false
                        this.termsDocStatus = `<p><span>&#10060;</span> ${error}</p>`
                  })
            }
      }

      onOfferLetterClick(event: any) {
            this.offerLetterDocStatus = ""
            let files = event.target.files
            if (files && files.length) {
                  let file = files[0]
                  let err = this.fileOperationService.isDocumentFileValid(file)
                  if (err != null) {
                        this.offerLetterDocStatus = `<p><span>&#10060;</span> ${err}</p>`
                        return
                  }
                  // Upload sample offer letter.
                  this.isOfferLetterFileUploading = true
                  this.fileOperationService.uploadOfferLetter(file).subscribe((data: any) => {
                        this.companyRequirementForm.markAsDirty()
                        this.companyRequirementForm.patchValue({
                              sampleOfferLetter: data
                        })
                        this.offerDisplayedFileName = file.name
                        this.isOfferLetterFileUploading = false
                        this.isOfferLetterUploadedToServer = true
                        this.offerLetterDocStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
                  }, (error) => {
                        this.isOfferLetterFileUploading = false
                        this.offerLetterDocStatus = `<p><span>&#10060;</span> ${error}</p>`
                  })
            }
      }

      // To chcek if company requirement's required before date is coming near.
      checkUrgency(requiredBefore: any, limit: number): boolean {
            let limitDate = new Date(new Date().setDate(new Date().getDate() + limit))
            let limitDay = limitDate.getDate()
            let limitMonth = limitDate.getMonth()
            let limitYear = limitDate.getFullYear()

            let requirementDate = new Date(requiredBefore)
            let requirementDay = requirementDate.getDate()
            let requirementMonth = requirementDate.getMonth()
            let requirementYear = requirementDate.getFullYear()

            if (requirementYear < limitYear || (requirementYear == limitYear && requirementMonth < limitMonth)) {
                  return true
            }
            if (requirementYear > limitYear || (requirementYear == limitYear && requirementMonth > limitMonth)) {
                  return false
            }
            if (requirementDay > limitDay) {
                  return false
            }
            return true
      }

      // Format fields of company requirement through interface.
      formatCompanyRequirementsFields(): void {
            for (let requirement of this.companyRequirement) {
                  // Assign Urgency level to requirement
                  requirement.isUrgent = this.checkUrgency(requirement.requiredBefore, 3)

                  // Assign talent rating in number to requirement.
                  for (let i = 0; i < this.talentRatingList.length; i++) {
                        if (this.talentRatingList[i].key == requirement.talentRating) {
                              requirement.talentRatingValue = this.talentRatingList[i].value
                        }
                  }

                  //       // Format package offered in indian rupee system.
                  //       requirement.packageOfferedInstring = this.formatNumberInIndianRupeeSystem(requirement.packageOffered)
            }
      }

      // Extract id from objects and give it to company requirement.
      setFormFields(companyRequirement: IRequirement): void {
            if (this.companyRequirementForm.get('salesPerson')?.value) {
                  companyRequirement.salesPersonID = this.companyRequirementForm.get('salesPerson').value.id
                  delete companyRequirement['salesPerson']
            }
            if (this.companyRequirementForm.get('company')?.value) {
                  // console.log( " this.companyRequirementForm.get('companyBranch').value.id",this.companyRequirementForm.get('companyBranch').value);

                  companyRequirement.companyID = this.companyRequirementForm.get('company').value.id
                  delete companyRequirement['company']
            }

            // designation is compulsory
            console.log("this.companyRequirementForm.get('designation')?.value", this.companyRequirementForm.get("designation")?.value);

            companyRequirement.designationID = this.companyRequirementForm.get("designation")?.value.id
            delete companyRequirement['designation']

            companyRequirement.requiredBefore = new Date(this.companyRequirementForm.get("requiredBefore")?.value).toISOString()
            companyRequirement.requiredFrom = new Date(this.companyRequirementForm.get("requiredFrom")?.value).toISOString()
      }

      // Used to dismiss modal.
      dismissFormModal(modal: NgbModalRef) {
            if (this.isTermsFileUploading || this.isOfferLetterFileUploading) {
                  alert("Please wait till file is being uploaded")
                  return
            }
            if (this.isTermsUploadedToServer || this.isOfferLetterUploadedToServer) {
                  if (!confirm("Uploaded file will be deleted.\nAre you sure you want to close?")) {
                        return
                  }
                  // delete terms and contion *******************************************************
                  // this.deleteResume()
            }
            modal.dismiss()
            this.isTermsUploadedToServer = false
            this.termsDisplayedFileName = "Select file"
            this.termsDocStatus = ""
      }

      validateFormFields(): void {
            this.companyRequirementForm.get('maximumPackage').setValidators([Validators.required,
            Validators.min(this.companyRequirementForm.get('minimumPackage')?.value + 1),
            Validators.max(100000000000)])

            this.companyRequirementForm.updateValueAndValidity()
      }

      // On clicking sumbit button in company requirement form.
      validateRequirementForm(): void {
            this.validateFormFields()

            // console.log(this.companyRequirementForm.controls)

            if (this.isTermsFileUploading || this.isOfferLetterFileUploading) {
                  alert("Please wait till file is being uploaded")
                  return
            }

            if (this.companyRequirementForm.get('rating')?.value) {
                  this.calculateAverageRating()
            }

            if (this.companyRequirementForm.invalid) {
                  this.companyRequirementForm.markAllAsTouched()
                  return
            }
            this.updateCompanyRequirement()
      }

      // On 'is experience' control value change in company requirement form.
      onIsExperiencedValueChange(isExperienced: boolean): void {
            if (isExperienced) {
                  this.companyRequirementForm.addControl('minimumExperience',
                        this.formBuilder.control(null, [Validators.required, Validators.min(1), Validators.max(30)]))
                  this.companyRequirementForm.addControl('maximumExperience',
                        this.formBuilder.control(null, [Validators.min(1), Validators.max(30)]))
                  this.showExperienceFields = true
                  return
            }
            this.companyRequirementForm.removeControl('minimumExperience')
            this.companyRequirementForm.removeControl('maximumExperience')
            this.showExperienceFields = false
      }

      // Set the minimum value of maximum experience of company requirement.
      setMaximumExperienceMinimumValue(value: number): void {
            if (value == 0) {
                  this.minValueForMaxExp = 1
                  return
            }
            this.minValueForMaxExp = value
            this.companyRequirementForm.get('maximumExperience').
                  setValidators([Validators.min(this.minValueForMaxExp)])
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

      // Enable the company requirement form.
      enableCompanyRequirementForm(): void {
            this.companyRequirementForm.enable()
            this.companyRequirementForm.get("jobLocation").get("state").enable()
            this.companyRequirementForm.get('code').disable()
      }

      // ================================================REDIRECTION FUNCTIONS FOR COMPANY REQUIREMENTS===============================================

      // Redirect to talents page filtered by waiting lists' company requirement id.
      redirectToTalentsForWaitingList(requirementID: string): void {
            this.router.navigate([this.urlConstant.TALENT_MASTER], {
                  queryParams: {
                        "waitingListRequirementID": requirementID
                  }
            }).catch(err => {
                  console.error(err)
            })
      }

      // Redirect to enquiries page filtered by waiting lists' company requirement id.
      redirectToEnquiriesForWaitingList(requirementID: string): void {
            this.router.navigate([this.urlConstant.TALENT_ENQUIRY], {
                  queryParams: {
                        "waitingListRequirementID": requirementID
                  }
            }).catch(err => {
                  console.error(err)
            })
      }

      // ================================================TALENT SEARCH FUNCTIONS FOR COMPANY REQUIREMENTS===============================================

      // View the search critera for talents.
      viewSearchRequirement(requirement: ICompanyRequirement): void {
            let requirementTechnologies = []
            let requirementQualifications = []

            if (requirement.technologies) {
                  for (let index = 0; index < requirement.technologies.length; index++) {
                        requirementTechnologies.push(requirement.technologies[index].id)
                  }
            }

            if (requirement.qualifications) {
                  for (let index = 0; index < requirement.qualifications.length; index++) {
                        requirementQualifications.push(requirement.qualifications[index].id)
                  }
            }

            this.searchParams = {
                  qualifications: requirementQualifications,
                  personalityType: requirement.personalityType,
                  minimumExperience: requirement.minimumExperience,
                  maximumExperience: requirement.maximumExperience,
                  isExperience: false,
                  talentType: requirement.talentRating,
                  technologies: requirementTechnologies,
            }

            if (requirement.minimumExperience) {
                  this.searchParams.isExperience = true
            }

            if (!this.searchParams.talentType) {
                  this.deleteTalentType()
            }
            if (!this.searchParams.isExperience) {
                  delete this.searchParams.minimumExperience
                  delete this.searchParams.maximumExperience
            }
            if (this.searchParams.technologies.length == 0) {
                  delete this.searchParams.technologies
            }
            if (this.searchParams.qualifications.length == 0) {
                  delete this.searchParams.qualifications
            }

            this.searchParams.requirementSearch = requirement.id
            this.selectedRequirement = requirement
            this.openModal(this.searchTalentsModal, 'lg')
      }

      // Delete miminum experience from search.
      deleteMinimumExperience(): void {
            delete this.searchParams.minimumExperience
      }

      // Delete maxinum experience from search.
      deleteMaximumExperience(): void {
            delete this.searchParams.maximumExperience
      }

      // Delete is experience from search.
      deleteIsExperience(): void {
            delete this.searchParams.isExperience
            if (this.searchParams.minimumExperience) {
                  delete this.searchParams.minimumExperience
            }
            if (this.searchParams.maximumExperience) {
                  delete this.searchParams.maximumExperience
            }
      }

      // Delete personlaity type from search.
      deletePersonalityType(): void {
            delete this.searchParams.personalityType
      }

      // Delete talent type from search.
      deleteTalentType(): void {
            delete this.searchParams.talentType
      }

      // Delete technologies from search.
      deleteTechnologies(id: string): void {
            let index = this.searchParams.technologies.indexOf(id)
            this.searchParams.technologies.splice(index, 1)

            if (this.searchParams.technologies.length == 0) {
                  delete this.searchParams.technologies
            }
      }

      // Delete qualifications from search.
      deleteQualification(id: string): void {
            let index = this.searchParams.qualifications.indexOf(id)
            this.searchParams.qualifications.splice(index, 1)

            if (this.searchParams.qualifications.length == 0) {
                  delete this.searchParams.qualifications
            }
      }

      // Redirect to talents based on search criteria.
      searchTalents(): void {

            this.router.navigate([this.urlConstant.TALENT_MASTER], {
                  queryParams: this.searchParams
            }).catch(err => {
                  console.error(err)
            })
      }

      // Redirect to talents to show selected talents for requirement.
      viewSelectedTalents(requirementID: string): void {
            this.spinnerService.loadingMessage = "Getting selected talents"

            this.router.navigate([this.urlConstant.TALENT_MASTER], {
                  queryParams: {
                        "requirementID": requirementID
                  }
            }).catch(err => {
                  console.error(err)
            })
      }

      // ================================================APPLICANT FUNCTIONS FOR COMPANY REQUIREMENTS===============================================

      // On clicking show applicants button in batch list.
      onShowApplicantsButtonClick(requirement: any): void {
            this.showTransferButton = false
            this.showTransferForm = false
            this.showCompanyReqField = false
            this.showBatchField = false
            this.selecetdrRequirementForTransfer = null
            this.selectedBatchForTransfer = null
            let modalSize: string
            if (!requirement.isActive) {
                  this.showTransferButton = true
                  modalSize = 'md'
            } else {
                  modalSize = 'sm'
            }
            this.getWaitingListByRequirementID(requirement)
            this.openModal(this.showApplicantsModal, modalSize)
            this.selectedRequirement = requirement
      }

      // Get waiting list by requirement.
      getWaitingListByRequirementID(requirement: any): void {
            this.spinnerService.loadingMessage = "Getting Applicants"


            let queryParams: any = {
                  companyRequirementID: requirement.id
            }
            if (requirement.isActive) {
                  queryParams.isActive = "1"
            }
            else {
                  queryParams.isActive = "0"
            }
            this.talentService.getTwoWaitingLists(queryParams).subscribe((response) => {
                  this.twoWaitingLists = response
                  if (!this.selectedRequirement.isActive && ((this.twoWaitingLists.talentWaitingList?.length != 0) ||
                        (this.twoWaitingLists.talentWaitingList?.length != 0))) {
                        this.showTransferButton = true
                  }
            }, (err) => {
                  console.error(err)
            })
      }

      // On clicking transfer button of applicants modal.
      onTrasnferButtonClick(): void {
            this.showTransferButton = false
            this.showTransferForm = true
      }

      // On changing value of waiting for control in waiting list form.
      onWaitingForChange(waitingFor: string): void {
            if (waitingFor == "Requirement") {
                  this.showCompanyReqField = true
                  this.showBatchField = false
                  this.selecetdrRequirementForTransfer = null
                  this.selectedBatchForTransfer = null
                  this.getActiveRequirementList()
            }
            if (waitingFor == "Batch") {
                  this.showCompanyReqField = false
                  this.showBatchField = true
                  this.selecetdrRequirementForTransfer = null
                  this.selectedBatchForTransfer = null
                  this.getActiveBatchList()
            }
      }

      // On clicking sunmit button of transfer form in applicants modal.
      onSubmitTransferButtonClick(): void {

            // If waiting for is not selected.
            if (!this.showCompanyReqField && !this.showBatchField) {
                  alert("Please select waiting for field")
                  return
            }

            // If compnay is selected.
            if (this.showCompanyReqField) {
                  if (!this.selecetdrRequirementForTransfer) {
                        alert("Please select a requirement ID")
                        return
                  }
                  let waitingLists: any[] = []
                  for (let i = 0; i < this.twoWaitingLists.talentWaitingList.length; i++) {
                        waitingLists.push(this.twoWaitingLists.talentWaitingList[i])
                  }
                  for (let i = 0; i < this.twoWaitingLists.enquiryWaitingList.length; i++) {
                        waitingLists.push(this.twoWaitingLists.enquiryWaitingList[i])
                  }
                  let updateWaitingList: any = {
                        companyID: this.selecetdrRequirementForTransfer.companyBranchID,
                        requirementID: this.selecetdrRequirementForTransfer.id,
                        batchID: null,
                        courseID: null,
                        waitingLists: waitingLists
                  }
                  this.spinnerService.loadingMessage = "Transfering Applicants"


                  this.talentService.transferWaitingList(updateWaitingList).subscribe((response) => {
                        alert("Some talents' or enquiries' waiting list entries may not be transfered since they already have the requirement assigned to them")
                        this.getWaitingListByRequirementID(this.selectedRequirement)
                        this.getAllRequirements()
                        this.selecetdrRequirementForTransfer = null
                        this.selectedBatchForTransfer = null
                        this.showTransferForm = false
                        return
                  }, (err) => {
                        console.error(err)
                  }).add(() => {

                  })
            }

            // If batch is selected.
            if (this.showBatchField) {
                  if (!this.selectedBatchForTransfer) {
                        alert("Please select a batch ID")
                        return
                  }
                  let waitingLists: any[] = []
                  for (let i = 0; i < this.twoWaitingLists.talentWaitingList.length; i++) {
                        waitingLists.push(this.twoWaitingLists.talentWaitingList[i])
                  }
                  for (let i = 0; i < this.twoWaitingLists.enquiryWaitingList.length; i++) {
                        waitingLists.push(this.twoWaitingLists.enquiryWaitingList[i])
                  }
                  let updateWaitingList: any = {
                        companyID: null,
                        requirementID: null,
                        batchID: this.selectedBatchForTransfer.id,
                        courseID: this.selectedBatchForTransfer.courseID,
                        waitingLists: waitingLists
                  }
                  this.spinnerService.loadingMessage = "Transfering Applicants"


                  this.talentService.transferWaitingList(updateWaitingList).subscribe((response) => {
                        alert("Some talents' or enquiries' waiting list entries may not be transfered since they already have the requirement assigned to them")
                        this.getWaitingListByRequirementID(this.selectedRequirement)
                        this.getAllRequirements()
                        this.selecetdrRequirementForTransfer = null
                        this.selectedBatchForTransfer = null
                        this.showTransferForm = false
                        return
                  }, (err) => {
                        console.error(err)
                  }).add(() => {

                  })
            }
      }

      // On clicking close transfer form button in applicants modal.
      onCloseTrasferFormButtonClick(): void {
            this.selecetdrRequirementForTransfer = null
            this.selectedBatchForTransfer = null
            this.showTransferButton = true
            this.showTransferForm = false
            this.showCompanyReqField = false
            this.showBatchField = false
      }

      // =============================================================ALLOCATION FUNCTIONS==========================================================================

      // On clicking allocate salesperson to requirement button.
      onAllocateSalesPersonToRequirementClick(requirement: any): void {
            this.selectedRequirement = requirement
            this.openModal(this.allocateSalespersonModal, 'sm')
      }

      // Allocate salesperson to requirement.
      allocateSalesPersonToRequirement(salesPersonID: string, requirementID: string): void {
            if (salesPersonID == "null") {
                  alert("Please select sales person")
                  return
            }
            let requirementIDsToBeUpdated = []
            requirementIDsToBeUpdated.push({
                  "requirementID": requirementID
            })
            this.spinnerService.loadingMessage = "Salesperson is getting allocated to company requirement"


            // console.log(salesPersonID)
            // console.log(requirementIDsToBeUpdated)
            this.companyRequirementService.allocateSalesPersonToCompanyRequirement(requirementIDsToBeUpdated, salesPersonID).subscribe((response: any) => {
                  this.getAllRequirements()
                  alert(response)
                  this.modalRef.close('success')
            }, (error) => {
                  console.error(error)
                  if (typeof error.error == 'object' && error) {
                        alert(this.utilService.getErrorString(error))
                        return
                  }
                  if (error.error == undefined) {
                        alert('Sales person could not be allocated to company requireemnt')
                  }
                  alert(error.statusText)
            })
      }

      // =============================================================GET FUNCTIONS==========================================================================
      // Gets all components.
      getAllComponents(): void {
            this.getSalesPersonList()
            this.getTalentRatingList()
            this.getCountryList()
            this.getTechnologyList()
            this.getDesignationList()
            this.getDegreeList()
            this.getPersonalityTypeList()
            this.getRatingFeedbackQuestions()
            this.searchOrGetCompanyRequirements()
            this.getCompanyList()
      }

      // Get salesperson list.
      getSalesPersonList(): void {
            this.generalService.getSalesPersonList().subscribe((data: any) => {
                  this.salesPersonList = data.body
            }, (err) => {
                  console.error(err)
            })
      }

      // Get talent rating list.
      getTalentRatingList(): void {
            this.generalService.getGeneralTypeByType("talent_rating").subscribe((data: any[]) => {
                  this.talentRatingList = data
            }, (err) => {
                  console.error(err)
            })
      }

      // Get country list.
      getCountryList(): void {
            this.generalService.getCountries().subscribe((data: any[]) => {
                  this.countryList = data.sort()
            }, (err) => {
                  console.error(err)
            })
      }

      // Get technology list.
      getTechnologyList(event?: any): void {
            // this.generalService.getTechnologies().subscribe((data: any[]) => {
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


      // Get designation list.
      getDesignationList(): void {
            this.generalService.getDesignations().subscribe((data: any[]) => {
                  this.designationList = data
                  console.log(this.designationList);

            }, (err) => {
                  console.error(err)
            })
      }

      // Get degree list.
      getDegreeList(): void {
            let queryParams: any = {
                  limit: -1,
                  offset: 0,
            }
            this.degreeService.getAllDegrees(queryParams).subscribe((data: any) => {
                  this.degreeList = data.body
            }, (err) => {
                  console.error(err)
            })
      }

      // Get personality type list.
      getPersonalityTypeList(): void {
            this.generalService.getGeneralTypeByType("personality_type").subscribe((data: any[]) => {
                  this.personalityTypeList = data
            }, (err) => {
                  console.error(err)
            })
      }

      // Get state list by country in requirement form.
      getStateListByCountry(country: any): void {
            if (country == null) {
                  this.stateList = []
                  this.jobLocation.get('state').setValue(null)
                  this.jobLocation.get('state').disable()
                  return
            }
            this.generalService.getStatesByCountryID(country.id).subscribe((respond: any[]) => {
                  this.stateList = respond
                  if (this.stateList.length == 0) {
                        this.jobLocation.get('state').setValue(null)
                        this.jobLocation.get('state').disable()
                        return
                  }
                  if (this.isViewMode) {
                        this.jobLocation.get('state').disable()
                        return
                  }
                  this.jobLocation.get('state').enable()
            }, (err) => {
                  console.error(err)
            })
      }

      // Get state list by country in search form.
      getStateListByCountryID(countryID: string): void {
            if (countryID == null) {
                  this.stateList = []
                  this.companyRequirementSearchForm.get('stateID').setValue(null)
                  this.companyRequirementSearchForm.get('stateID').disable()
                  return
            }
            this.generalService.getStatesByCountryID(countryID).subscribe((respond: any[]) => {
                  this.stateList = respond
                  if (this.stateList.length == 0) {
                        this.companyRequirementSearchForm.get('stateID').setValue(null)
                        this.companyRequirementSearchForm.get('stateID').disable()
                        return
                  }
                  this.companyRequirementSearchForm.get('stateID').enable()
            }, (err) => {
                  console.error(err)
            })
      }

      // Get active company requirement list.
      getActiveRequirementList(): void {
            let queryParams: any = {
                  isActive: "1"
            }
            this.generalService.getRequirementList(queryParams).subscribe((response: any) => {
                  this.activeRequirementList = response.body
            }, (err: any) => {
                  console.error(err)
            })
      }

      // Get active batch list.
      getActiveBatchList(): void {
            let queryParams: any = {
                  batchStatus: ["Ongoing", "Upcoming"],
                  isActive: "1"
            }
            this.generalService.getBatchList(queryParams).subscribe((response: any) => {
                  this.activeBatchList = response
            }, (err: any) => {
                  console.error(err)
            })
      }

      getRatingFeedbackQuestions(): void {
            this.dummyDataService.getRatingFeedbackQuestions().then((response: any) => {
                  this.feedbackQuestions = response
            }).catch((err) => {
                  console.error(err);
            })
      }
      // Get company list.
      getCompanyList(): void {
            let querparams = {
                  limit: -1,
                  offset: 0
            }
            this.companyService.getAllCompanies(querparams).subscribe((response: any) => {
                  this.companyList = response.body
                  console.log(this.companyList);

            }, (err: any) => {
                  console.error(err)
            })
      }

      // Get all company requirements.
      getAllRequirements(): void {
            this.spinnerService.loadingMessage = "Getting Company Requirements"


            this.isRequrimentLoaded = true
            this.companyRequirement = []
            this.totalCompanyRequirements = 0
            if (!this.showForSalesPerson) {
                  if (this.searchFormValue == null) {
                        this.searchFormValue = {}
                  }
                  this.searchFormValue.salesPersonID = this.localService.getJsonValue("loginID")
            }
            if (!this.showForCompany) {
                  if (this.searchFormValue == null) {
                        this.searchFormValue = {}
                  }
                  this.searchFormValue.companyID = this.localService.getJsonValue("loginID")
            }
            this.companyRequirementService.getCompanyRequirements(this.searchFormValue).subscribe((response) => {
                  this.totalCompanyRequirements = response.headers.get('X-Total-Count')
                  this.companyRequirement = response.body
                  this.formatCompanyRequirementsFields()
            }, error => {
                  this.companyRequirement = []
                  console.error(error)
                  if (error.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                  }
            }).add(() => {
                  if (this.companyRequirement.length == 0) {
                        this.isRequrimentLoaded = false
                  }
                  this.setPaginationString()
            })
      }

      // ================================================ ADD REQUIREMENT RATING ===============================================

      // createFeedbackGroupForm(): void {
      //       this.requirementRatingForm = this.formBuilder.group({
      //             ratings: new FormArray([])
      //       })
      // }

      // get ratingArray(): FormArray {
      //       return this.requirementRatingForm.get("ratings") as FormArray
      // }

      // addRequirementRating(): void {
      //       this.ratingArray.push(this.formBuilder.group({
      //             requirementID: new FormControl(null, [Validators.required]),
      //             questionID: new FormControl(null, [Validators.required]),
      //             optionID: new FormControl(null, [Validators.required]),
      //             answer: new FormControl(null, [Validators.required]),
      //             requirement: new FormControl(null, [Validators.required]),
      //             question: new FormControl(null, [Validators.required]),
      //             option: new FormControl(null, [Validators.required]),
      //       }))
      // }

      onRatingViewClick(requirement: IRequirement): void {
            // console.log(requirement);
            this.selectedRequirement = requirement
            this.salaryTrend = requirement.salaryTrend

            this.createRequirementRatingForm(requirement)
            this.companyRequirementForm.disable()
            this.isRatingAddClick = false
            if (requirement.rating == null) {
                  this.isRatingAddClick = true
                  this.companyRequirementForm.enable()
            }
            this.openModal(this.addRequirementRatingModal, "lg")
      }

      onRatingAddClick(): void {
            this.isRatingAddClick = true
            this.companyRequirementForm.enable()
      }

      createRequirementRatingForm(requirement: IRequirement): void {
            this.createCompanyRequirementForm()
            this.setCompulsoryFieldsForRating()

            // Experience.
            this.companyRequirementForm.get('isExperience').setValue(requirement.minimumExperience != null)
            this.onIsExperiencedValueChange(requirement.minimumExperience != null)
            requirement.requiredBefore = this.datePipe.transform(requirement.requiredBefore, 'yyyy-MM-dd')
            requirement.requiredFrom = this.datePipe.transform(requirement.requiredFrom, 'yyyy-MM-dd')
            this.companyRequirementForm.patchValue(requirement)
      }

      setCompulsoryFieldsForRating(): void {

            this.companyRequirementForm.get("requiredBefore").setValidators([Validators.required])
            this.companyRequirementForm.get("requiredFrom").setValidators([Validators.required])

            for (let index = 0; index < this.feedbackQuestions.length; index++) {
                  this.companyRequirementForm.get(this.feedbackQuestions[index].columnName).setValidators([Validators.required])
            }

            this.companyRequirementForm.get("rating").setValidators([Validators.required])
            this.utilService.updateValueAndValiditors(this.companyRequirementForm)
      }

      validateRating(): void {
            this.companyRequirementForm.get("designationID").setValue(this.companyRequirementForm.get("designation")?.value.id)
            // console.log(this.companyRequirementForm.controls);

            this.calculateAverageRating()
            if (this.companyRequirementForm.invalid) {
                  this.companyRequirementForm.markAllAsTouched()
                  return
            }
            this.isRatingAddClick = true
            this.updateCompanyRequirement()
      }

      calculateAverageRating(): void {
            let averageScore: number = 0.0
            let totalScore: number = 0.0

            // Calculate maxscore.
            for (let index = 0; index < this.feedbackQuestions.length; index++) {
                  totalScore += this.feedbackQuestions[index].maxScore
            }

            for (let index = 0; index < this.feedbackQuestions.length; index++) {
                  for (let j = 0; j < this.feedbackQuestions[index].options.length; j++) {
                        if (this.feedbackQuestions[index].options[j].value ===
                              this.companyRequirementForm.get(this.feedbackQuestions[index].columnName).value) {
                              averageScore += this.feedbackQuestions[index].options[j].key
                        }
                  }
            }

            // ============================= PACKAGE SCORE =============================

            totalScore += 25 // maxscore for package is 25.
            let avgPackage = (this.companyRequirementForm.get("minimumPackage")?.value +
                  this.companyRequirementForm.get("maximumPackage")?.value) / 2

            averageScore += this.getPackageScore(avgPackage)

            // ============================= TECH SCORE =============================

            let averageTechScore: number = 0
            let totalTechScore: number = 0
            if (this.companyRequirementForm.get("technologies")?.value) {
                  let technologies: ITechnology[] = this.companyRequirementForm.get("technologies")?.value
                  for (let tech of technologies) {
                        if (tech.rating) {
                              totalTechScore += 5 // max rating for technology is 5.
                              averageTechScore += tech.rating
                        }
                  }
            }

            if (averageTechScore > 0 && totalTechScore > 0) {
                  totalScore += 15 // max score for tech is 15
                  averageScore += (averageTechScore * 15) / totalTechScore
            }

            // ============================= TALENT PERSONALTIY SCORE =============================
            if (this.companyRequirementForm.get("personalityType")?.value != null) {
                  totalScore += 3 // max score for personality type is 3
                  if (this.companyRequirementForm.get("personalityType")?.value == "Leader") {
                        averageScore += 1
                  } else if (this.companyRequirementForm.get("personalityType")?.value == "Introvert") {
                        averageScore += 2
                  } else if (this.companyRequirementForm.get("personalityType")?.value == "Extrovert") {
                        averageScore += 2
                  }
            } else {
                  totalScore += 3
                  averageScore += 3
            }

            // ============================= TALENT RATING SCORE =============================

            if (this.companyRequirementForm.get("talentRating")?.value != null) {
                  totalScore += 4 // max score for talent rating is 4
                  if (this.companyRequirementForm.get("talentRating")?.value >= 5) {
                        averageScore += 1
                  } else if (this.companyRequirementForm.get("talentRating")?.value <= 2) {
                        averageScore += 2
                  } else {
                        averageScore += 3
                  }
            } else {
                  totalScore += 4
                  averageScore += 4
            }

            // ============================= TALENT EXPERIENCE SCORE =============================

            if (this.companyRequirementForm.get("minimumExperience")?.value != null) {
                  let avgExp: number = 0
                  avgExp += this.companyRequirementForm.get("minimumExperience")?.value
                  if (this.companyRequirementForm.get("maximumExperience")?.value != null) {
                        avgExp += this.companyRequirementForm.get("maximumExperience")?.value
                        avgExp = avgExp / 2
                  }
                  averageScore += this.getExperienceScore(avgExp)
            } else {
                  averageScore += 10
            }
            totalScore += 10 // max score for experience is 10.

            console.log("Average Score -> ", averageScore, " totalScore -> ", totalScore);
            averageScore = (averageScore * 10.0) / totalScore
            this.companyRequirementForm.get('rating').setValue(averageScore)
      }

      getPackageScore(avgPackage: number): number {

            if (avgPackage <= 300000) {
                  return 5
            }

            if (avgPackage > 300000 && avgPackage <= 500000) {
                  return 10
            }

            if (avgPackage > 500000 && avgPackage <= 1000000) {
                  return 15
            }

            if (avgPackage > 1000000 && avgPackage <= 1500000) {
                  return 20
            }

            if (avgPackage > 1500000) {
                  return 25
            }

            return 0
      }

      getExperienceScore(avgExp: number): number {

            if (avgExp > 1 && avgExp <= 3) {
                  return 8
            }

            if (avgExp > 3 && avgExp <= 6) {
                  return 6
            }

            if (avgExp > 6 && avgExp <= 8) {
                  return 4
            }

            if (avgExp > 8 && avgExp <= 10) {
                  return 2
            }

            if (avgExp > 10) {
                  return 1
            }

            return 0
      }



      // Used to dismiss modal.
      dismissRequirementFormModal(modal: NgbModalRef) {
            if (this.isRequirementTermsFileUploading || this.isOfferLetterFileUploading) {
                  alert("Please wait till file is being uploaded")
                  return
            }
            if (this.isRequirementTermsUploadedToServer || this.isOfferLetterUploadedToServer) {
                  if (!confirm("Uploaded file will be deleted.\nAre you sure you want to close?")) {
                        return
                  }
            }
            modal.dismiss()
            this.isRequirementTermsUploadedToServer = false
            this.requirementTermsDisplayedFileName = "Select file"
            this.requirementTermsDocStatus = ""
      }

      //On uplaoding one pager
      onRequirementTermsAndConditionSelect(event: any) {
            this.requirementTermsDocStatus = ""
            let files = event.target.files
            if (files && files.length) {
                  let file = files[0]
                  let err = this.fileOps.isDocumentFileValid(file)
                  if (err != null) {
                        this.requirementTermsDocStatus = `<p><span>&#10060;</span> ${err}</p>`
                        return
                  }
                  // Upload terms and condition if it is present.]
                  this.isRequirementTermsFileUploading = true
                  this.fileOps.uploadTermsAndCondition(file).subscribe((data: any) => {
                        this.companyRequirementForm.markAsDirty()
                        this.companyRequirementForm.patchValue({
                              termsAndConditions: data
                        })
                        this.requirementTermsDisplayedFileName = file.name
                        this.isRequirementTermsFileUploading = false
                        this.isRequirementTermsUploadedToServer = true
                        this.requirementTermsDocStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
                  }, (error) => {
                        this.isRequirementTermsFileUploading = false
                        this.requirementTermsDocStatus = `<p><span>&#10060;</span> ${error}</p>`
                  })
            }
      }

      addCompanyRequirement(): void {
            this.spinnerService.loadingMessage = "Adding company requirement";

            let companyRequirement: IRequirement = this.companyRequirementForm.value
            if (this.isSalesPerson) {
                  companyRequirement.salesPersonID = this.localService.getJsonValue("loginID")
            }
            this.setFormFields(companyRequirement)
            console.log(companyRequirement);
            this.companyRequirementService.addCompanyRequirement(companyRequirement).subscribe((data: string) => {
                  this.modalRef.close();
                  this.companyRequirementForm.reset();

                  // alert("requirement added with id:" + data)
                  alert("Requirement successfully added")
                  this.searchOrGetCompanyRequirements()
            }, (error) => {
                  console.error(error);

                  if (error.error) {
                        alert(error.error?.error)
                        return
                  }
                  alert(error.error?.error)
            });
      }

      //validate company requirement form
      validateCompanyRequirementForm(): void {
            console.log(this.companyRequirementForm);

            this.validateFormFields()

            // console.log(this.companyRequirementForm.controls)

            if (this.isRequirementTermsFileUploading || this.isOfferLetterFileUploading) {
                  alert("Please wait while file is getting uploaded.")
                  return
            }

            if (this.isRequirementTermsFileUploading || this.isOfferLetterFileUploading) {
                  alert("Please wait while file is getting uploaded.")
                  return
            }

            if (this.companyRequirementForm.invalid) {
                  this.companyRequirementForm.markAllAsTouched();
                  return
            }

            this.addCompanyRequirement();
            console.log(this.companyRequirementForm);

      }
}