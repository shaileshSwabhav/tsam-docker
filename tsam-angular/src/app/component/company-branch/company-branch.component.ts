import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, Validators, FormControl } from '@angular/forms';
import { IState, ICompanyBranch, ICountry, CompanyService, IDomain, ITechnology } from 'src/app/service/company/company.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { formatDate, Location } from '@angular/common';
import { requirementDateValidator } from 'src/app/Validators/custom.validators';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { NgbModal, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { ActivatedRoute, Router } from '@angular/router';
import { Constant, Role, UrlConstant } from 'src/app/service/constant';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { CompanyRequirementService, IRequirement } from 'src/app/service/company/company-requirement/company-requirement.service';
import { TechnologyService } from 'src/app/service/technology/technology.service';
import { DegreeService } from 'src/app/service/degree/degree.service';
import { UrlService } from 'src/app/service/url.service';
import { CompanyMasterComponent } from '../company-master/company-master.component';


@Component({
  selector: 'app-company-branch',
  templateUrl: './company-branch.component.html',
  styleUrls: ['./company-branch.component.css']
})
export class CompanyBranchComponent implements OnInit {

  // Flags.
  showExperienceFields: boolean

  // Requirement form relatef=d variables.
  // Minimum value for maximum experience.
  minValueForMaxExp: number

  // Company.
  selectedBranch: any;
  companyBranchList: ICompanyBranch[];
  companyForm: FormGroup;
  companyBranchForm: FormGroup
  isCompanyBranchList: boolean

  // Required Componenets/Fields.
  allDomains: IDomain[];
  allTechnologies: ITechnology[];
  salesPeople: any[];
  allStates: IState[];
  stateList: IState[]
  allCountries: ICountry[];
  ratingList: number[];

  // tech components
  isTechLoading: boolean
  techLimit: number
  techOffset: number

  // Search.
  searchCompanyForm: FormGroup;
  searchFormValue: any;
  searched: boolean;
  showSearch: boolean;

  // constant
  nilUUID: string

  // Pagination.
  limit: number;
  offset: number;
  currentPage: number;
  totalCompanyBranches: number;
  paginationString: string

  // Modal.
  modalHeader: string;
  modalButton: string;
  selectedCompanyBranch: ICompanyBranch;
  modalAction: () => void
  formHandler: () => void;
  modalRef: any;
  isViewClicked: boolean

  // navigation
  navigatedCompanyID: string
  isNavigatedFromCompany: boolean

  // access
  permission: IPermission
  isSalesPerson: boolean
  isAdmin: boolean

  // one pager
  isOnePagerUploadedToServer: boolean
  isOnePagerFileUploading: boolean
  onePagerDocStatus: string
  onePagerDisplayedFileName: string

  // terms and conditions
  isTermsUploadedToServer: boolean
  isTermsFileUploading: boolean
  termsDocStatus: string
  termsDisplayedFileName: string

  // requirement terms and conditions
  isRequirementTermsUploadedToServer: boolean
  isRequirementTermsFileUploading: boolean
  requirementTermsDocStatus: string
  requirementTermsDisplayedFileName: string

  // url
  previousUrl: string

  // Modals
  @ViewChild('companyDetailModal') companyDetailModal: any
  @ViewChild('companyRequirementModal') companyRequirementModal: any
  @ViewChild('deleteConfirmationModal') deleteConfirmationModal: any
  @ViewChild('drawer') drawer: any

  //constant
  private readonly IGNORE_SEARCH_FIELD: string[] = ["limit", "offset"]

  constructor(
    private formBuilder: FormBuilder,
    private companyService: CompanyService,
    private companyRequirementService: CompanyRequirementService,
    private generalService: GeneralService,
    private techService: TechnologyService,
    private degreeService: DegreeService,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private activatedRoute: ActivatedRoute,
    private router: Router,
    private constant: Constant,
    private urlConstant: UrlConstant,
    private utilService: UtilityService,
    private localService: LocalService,
    private fileOps: FileOperationService,
    private roleConstant: Role,
    private location: Location,
    private urlService: UrlService,
  ) {
    this.initializeVariables()
    // this.getPreviousPageUrl()
    this.editorConfig()
    this.createForms()
    this.getAllComponents()
  }

  initializeVariables(): void {

    // Flags.
    this.showExperienceFields = false
    this.minValueForMaxExp = 1

    this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.COMPANY_BRANCH)

    this.isSalesPerson = this.localService.getJsonValue("roleName") == this.roleConstant.SALES_PERSON
    this.isAdmin = this.localService.getJsonValue("roleName") == this.roleConstant.ADMIN

    this.salesPeople = []
    this.allStates = []
    this.stateList = []

    this.limit = 5;
    this.offset = 0;
    this.techLimit = 10
    this.techOffset = 0

    this.searched = false
    this.showSearch = false
    this.isViewClicked = false
    this.isNavigatedFromCompany = false
    this.isCompanyBranchList = true
    this.isTechLoading = false

    // one-pager.
    this.isOnePagerUploadedToServer = false
    this.isOnePagerFileUploading = false
    this.onePagerDisplayedFileName = "Select file"

    // terms and conditions.
    this.isTermsUploadedToServer = false
    this.isTermsFileUploading = false
    this.termsDisplayedFileName = "Select file"

    // requirement terms and conditions.
    this.isRequirementTermsUploadedToServer = false
    this.isRequirementTermsFileUploading = false
    this.requirementTermsDisplayedFileName = "Select file"

    // sample offer letter.
    this.isOfferLetterUploadedToServer = false
    this.isOfferLetterFileUploading = false
    this.offerDisplayedFileName = "Select File"
    this.offerLetterDocStatus = ""


    this.searchFormValue = {}
    this.nilUUID = this.constant.NIL_UUID
    this.previousUrl = null

    this.minRequiredDate = formatDate(new Date().setDate(new Date().getDate() + 10), 'yyyy/MM/dd', 'en')

    // console.log(this.router.getCurrentNavigation())
  }

  createForms(): void {
    this.createSearchCompanyForm()
    this.createCompanyRequirementForm()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit() {
  }

  backToPreviousPage(): void {
    this.urlService.goBack(CompanyMasterComponent.name)
  }

  // Handles pagination
  changePage(pageNumber: number): void {
    // $event will be the page number & offset will be 1 less than it.

    // this.offset = $event - 1
    // this.currentPage = $event
    // this.getAllBranches()

    this.searchCompanyForm.get("offset").setValue(pageNumber - 1)

    this.limit = this.searchCompanyForm.get("limit").value
    this.offset = this.searchCompanyForm.get("offset").value

    // this.getAllBranches();
    this.searchBranches();
  }

  // gets all required components like all states,countries etc
  getAllComponents(): void {
    this.getAllCountries()
    this.getAllRatings()
    this.getAllStates()
    this.getAllSalesPeople()
    this.getAllDomains()
    this.getAllTechnologies()
    this.searchOrGetCompanyBranches()
  }

  // Compare two Object.
  compareFn(obj1: any, obj2: any): any {
    if (obj1 == null && obj2 == null) {
      return true
    }
    return obj1 && obj2 ? obj1.id == obj2.id : obj1 == obj2;
  }

  showSpecificStates(branch: any, state: IState): boolean {
    if (branch.get('country').value == null) {
      return false
    }
    if (branch.get('country').value.id === state.countryID) {
      return true
    }
    return false
  }

  // Create search form.
  createSearchCompanyForm(): void {
    this.searchCompanyForm = this.formBuilder.group({
      branchName: new FormControl(null, [Validators.maxLength(100)]),
      city: new FormControl(null, [Validators.maxLength(50), Validators.pattern("^[a-zA-Z]+([a-zA-Z ]?)+")]),
      companyRating: new FormControl(null),
      coordinatorID: new FormControl(null),
      companyID: new FormControl(null),
      state: new FormControl(null),
      hrHeadName: new FormControl(null, [Validators.maxLength(100), Validators.pattern("^[a-zA-Z]+([a-zA-Z. ]?)+")]),
      numberOfEmployees: new FormControl(null, [Validators.max(1000000), Validators.pattern("^[1-9][0-9]*$")]),
      salesPersonID: new FormControl(null),
      domains: new FormControl(null),
      technologies: new FormControl(null),
      limit: new FormControl(this.limit),
      offset: new FormControl(this.offset),
    })
  }

  // Add New Branch.
  createCompanyBranchForm(): void {
    this.companyBranchForm = this.formBuilder.group({
      id: new FormControl(null),
      branchName: new FormControl(null, [Validators.required, Validators.pattern("^[a-zA-Z]+([a-zA-Z. ]?)+")]),
      code: new FormControl(null),
      companyID: new FormControl(null),
      address: new FormControl(null, [Validators.required]),
      pinCode: new FormControl(null, [Validators.required, Validators.pattern("^[0-9]{6}$")]),
      country: new FormControl(null, [Validators.required]),
      state: new FormControl(null, [Validators.required]),
      city: new FormControl(null, [Validators.required, Validators.maxLength(50), Validators.pattern("^[a-zA-Z]+([a-zA-Z ]?)+")]),
      mainBranch: new FormControl(null, [Validators.required]),
      domains: new FormControl(null), //, [Validators.required]
      technologies: new FormControl(null, [Validators.required]),
      companyRating: new FormControl(null),
      hrHeadName: new FormControl(null), //, [Validators.required, Validators.maxLength(100), Validators.pattern("^[a-zA-Z]+([a-zA-Z. ]?)+")]
      hrHeadContact: new FormControl(null, [Validators.pattern(/^[6789]\d{9}$/)]),
      hrHeadEmail: new FormControl(null, [Validators.email]),
      unitHeadName: new FormControl(null), //, [Validators.required, Validators.maxLength(100), Validators.pattern("^[a-zA-Z]+([a-zA-Z. ]?)+")]
      unitHeadContact: new FormControl(null, [Validators.pattern(/^[6789]\d{9}$/)]),
      unitHeadEmail: new FormControl(null, [Validators.email]),
      technologyHeadName: new FormControl(null), //, [Validators.required, Validators.maxLength(100), Validators.pattern("^[a-zA-Z]+([a-zA-Z. ]?)+")]
      technologyHeadContact: new FormControl(null, [Validators.pattern(/^[6789]\d{9}$/)]),
      technologyHeadEmail: new FormControl(null, [Validators.email]),
      financeHeadName: new FormControl(null), //, [Validators.required, Validators.maxLength(100), Validators.pattern("^[a-zA-Z]+([a-zA-Z. ]?)+")]
      financeHeadContact: new FormControl(null, [Validators.pattern(/^[6789]\d{9}$/)]),
      financeHeadEmail: new FormControl(null, [Validators.email]),
      recruitmentHeadName: new FormControl(null), //, [Validators.required, Validators.maxLength(100), Validators.pattern("^[a-zA-Z]+([a-zA-Z. ]?)+")]
      recruitmentHeadContact: new FormControl(null, [Validators.pattern(/^[6789]\d{9}$/)]),
      recruitmentHeadEmail: new FormControl(null, [Validators.email]),
      numberOfEmployees: new FormControl(null, [Validators.max(1000000), Validators.pattern("^[1-9][0-9]*$")]),
      salesPerson: new FormControl(null, [Validators.required]),
      onePager: new FormControl(null),
      termsAndConditions: new FormControl(null)

    })
  }

  onViewBranchClick(branch: any): void {
    // console.log(branch)
    this.isViewClicked = true
    this.createCompanyBranchForm()
    this.getStateList(branch.country.id)
    this.selectedBranch = branch
    this.setVariableForViewCompanyBranch()
    this.companyBranchForm.disable()
    this.companyBranchForm.patchValue(branch)
    console.log(this.companyBranchForm.value);
    this.openCompanyDetailModal(this.companyDetailModal, "xl")
  }

  onUpdateBranchClick() {
    this.setVariableForUpdateCompanyBranch();
    this.isViewClicked = false
    this.companyBranchForm.enable()
    if (this.isSalesPerson) {
      this.companyBranchForm.get('salesPerson').disable()
    }
    if (this.companyBranchForm.get('termsAndConditions')?.value) {
      this.termsDisplayedFileName = `<a href=${this.companyBranchForm.get('termsAndConditions')?.value} 
                                                      target="_blank">Terms and conditions</a>`
    }
    if (this.companyBranchForm.get('onePager')?.value) {
      this.onePagerDisplayedFileName = `<a href=${this.companyBranchForm.get('onePager')?.value} 
                                                      target="_blank">One Pager</a>`
    }
  }

  onDeleteBranchClick(branch: ICompanyBranch): void {
    this.selectedCompanyBranch = branch
    this.openModal(this.deleteConfirmationModal, "md")
  }

  updateForm(): void {
    this.companyBranchForm.patchValue(this.selectedBranch)
  }

  setVariableForViewCompanyBranch(): void {
    this.updateVariable(this.createCompanyBranchForm, "Branch Details", "Branch", this.updateBranch)
  }

  setVariableForUpdateCompanyBranch(): void {
    this.updateVariable(this.updateForm, "Update Company", "Update Branch", this.updateBranch)
  }

  updateVariable(formaction: () => void, modalheader: string, modalbutton: string, modalaction: () => void): void {
    this.formHandler = formaction;
    this.modalHeader = modalheader;
    this.modalButton = modalbutton;
    this.modalAction = modalaction;
  }

  setPaginationString(): void {
    this.paginationString = ''
    let start: number = this.limit * this.offset + 1
    let end: number = +this.limit + this.limit * this.offset
    if (this.totalCompanyBranches < end) {
      end = this.totalCompanyBranches
    }
    if (this.totalCompanyBranches == 0) {
      this.paginationString = ''
      return
    }
    this.paginationString = `${start} - ${end}`
  }

  searchOrGetCompanyBranches(): void {
    let queryParams = this.activatedRoute.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.changePage(1)
      // this.getAllBranches()
      return
    }
    if (queryParams.companyID) {
      this.navigatedCompanyID = queryParams.companyID
      this.isNavigatedFromCompany = true
      // this.searchCompanyForm.patchValue(queryParams)
      // this.searchBranches()
      // return
    }
    this.searchCompanyForm.patchValue(queryParams)
    this.searchBranches()
  }

  // =========================================================GET AND SEARCH FUNCTIONS=====================================================

  closeSearchDrawer(): void {
    this.drawer.toggle()
    this.searchBranches()
  }

  searchBranches(): void {
    this.searchFormValue = { ...this.searchCompanyForm?.value }
    this.router.navigate([], {
      relativeTo: this.activatedRoute,
      queryParams: this.searchFormValue,
    })
    let flag: boolean = true
    for (let field in this.searchFormValue) {
      if (!this.searchFormValue[field]) {
        delete this.searchFormValue[field];
      } else {
        if (!this.IGNORE_SEARCH_FIELD.includes(field)) {
          this.searched = true
        }
        flag = false
      }
    }
    // No API call on empty search.
    if (flag) {
      return
    }
    // this.changePage(1)
    this.getAllBranches()
  }


  // gets all branches for admin login
  getAllBranches(): void {
    // Calls search if a search has been made
    if (this.searched) {
      this.getAllSearchedBranches()
      return
    }
    // calls salesperson branches
    if (this.isSalesPerson) {
      this.getAllBranchesForSalesPerson()
      return
    }
    // calls branches of company
    if (this.isNavigatedFromCompany) {
      this.getAllBranchesByCompany()
      return
    }
    this.companyBranchList = []
    this.totalCompanyBranches = 0
    this.isCompanyBranchList = true
    this.spinnerService.loadingMessage = "Getting companies";
    this.companyService.getAllCompanyBranches(this.searchFormValue).subscribe(data => {
      this.companyBranchList = data.body;
      // console.log(this.companyBranchList);
      // this.assignOtherFields()
      this.totalCompanyBranches = parseInt(data.headers.get('X-Total-Count'))
      this.setPaginationString()
      if (this.totalCompanyBranches == 0) {
        this.isCompanyBranchList = false
      }
    }, (error) => {
      this.totalCompanyBranches = 0
      this.setPaginationString()
      this.isCompanyBranchList = false
      console.error(error);
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(this.utilService.getErrorString(error));
    });
  }

  // search branches for all logins
  getAllSearchedBranches(): void {
    // console.log(this.searchFormValue);
    if (this.isSalesPerson) {
      this.searchFormValue.salesPersonID = this.localService.getJsonValue('loginID')
    }
    this.spinnerService.loadingMessage = "Searching branches"

    this.companyBranchList = []
    this.totalCompanyBranches = 0
    this.isCompanyBranchList = true
    this.companyService.getAllCompanyBranches(this.searchFormValue).subscribe(data => {
      this.companyBranchList = data.body;
      console.log(this.companyBranchList);
      // this.assignOtherFields()
      this.totalCompanyBranches = parseInt(data.headers.get('X-Total-Count'))
      this.setPaginationString()
      if (this.totalCompanyBranches == 0) {
        this.isCompanyBranchList = false
      }
    }, (error) => {
      this.totalCompanyBranches = 0
      this.setPaginationString()
      this.isCompanyBranchList = false
      console.error(error);
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(this.utilService.getErrorString(error));
    });
  }

  // gets all branches for salesperson login 
  getAllBranchesForSalesPerson(): void {
    this.spinnerService.loadingMessage = "Getting companies"

    // this.searchFormValue = {}
    // sales person navigating from company

    if (this.isNavigatedFromCompany) {
      this.searchFormValue.companyID = this.navigatedCompanyID
    }
    this.companyService.getAllBranchesForSalesPerson(this.localService.getJsonValue('loginID'),
      this.searchFormValue).
      subscribe((response: any) => {
        // console.log(response.body);

        this.companyBranchList = response.body
        this.totalCompanyBranches = parseInt(response.headers.get('X-Total-Count'))

      }, (error: any) => {

        if (error.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
          return
        }

        alert(this.utilService.getErrorString(error));
        console.error(error);
      })
  }

  // get all branches for the navigated company
  getAllBranchesByCompany(): void {
    // this.searched = true
    // calls navigated branch company for logged in salesperson as it gets called from constructor
    if (this.isSalesPerson) {
      this.getAllBranchesForSalesPerson()
      return
    }
    this.spinnerService.loadingMessage = "Getting company branches"
    this.companyService.getAllBranchesOfCompany(this.navigatedCompanyID, this.searchFormValue).subscribe((response: any) => {
      // console.log(response.body);
      this.companyBranchList = response.body
      this.setPaginationString()

      this.totalCompanyBranches = this.companyBranchList.length
    }, (error) => {

      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  resetSearchAndGetAll(): void {
    this.spinnerService.loadingMessage = "Getting companies";

    this.searched = false
    // this.isNavigatedFromCompany = false
    this.resetSearchCompanyForm()
    this.searchFormValue = null
    // changes -> after search form reset check if companyID exists and if it exist patch it to form and add it in queryparams.
    if (this.navigatedCompanyID) {
      this.searchCompanyForm.get("companyID").setValue(this.navigatedCompanyID)

      this.searchFormValue = { ...this.searchCompanyForm?.value }

      this.router.navigate([this.urlConstant.COMPANY_BRANCH], {
        relativeTo: this.activatedRoute,
        queryParams: this.searchFormValue,
      })
      this.searchBranches()
      return
    }
    this.router.navigate([this.urlConstant.COMPANY_BRANCH])
    this.changePage(1)

  }

  resetSearchCompanyForm(): void {
    this.limit = this.searchCompanyForm.get("limit").value
    this.offset = this.searchCompanyForm.get("offset").value
    // this.currentPage = this.companyBranchForm.get("currentPage").value

    this.searchCompanyForm.reset({
      limit: this.limit,
      offset: this.offset,
      // currentPage: new FormControl(this.currentPage),
    })
  }

  // assignOtherFields(): void {
  //       for (let branch of this.companyBranchList) {
  //             for (let salesPerson of this.salesPeople) {
  //                   if (branch.salesPersonID == salesPerson.id) {
  //                         branch.salesPersonName = salesPerson.name
  //                   }
  //             }
  //       }
  // }

  // =========================================================GET AND SEARCH FUNCTIONS END=====================================================


  // Update company branch.
  updateBranch(): void {
    this.spinnerService.loadingMessage = "Updating company";

    // console.log(this.companyBranchForm.value);
    this.companyService.updateCompanyBranch(this.companyBranchForm.value).subscribe((response: any) => {

      this.modalRef.close()
      this.companyBranchForm.reset()
      alert("Branch successfully updated")
      this.getAllBranches()
    }, (err) => {

      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(this.utilService.getErrorString(err));
      console.error(err.error.error)
    })

  }

  deleteCompanyBranch(): void {
    // console.log(this.selectedCompanyBranch.mainBranch);
    if (this.selectedCompanyBranch.mainBranch) {
      if (!confirm("You are deleting main branch of this company.\nAre you sure you want to continue?")) {
        return
      }
    }
    this.modalRef.close();
    this.spinnerService.loadingMessage = "Deleting company";

    this.companyService.deleteCompanyBranch(this.selectedCompanyBranch.companyID,
      this.selectedCompanyBranch.id).subscribe((data: any) => {
        // alert(data)
        alert("Branch successfully deleted")

        if (this.totalCompanyBranches == 1) {
          this.companyBranchList = []
          this.totalCompanyBranches = 0
          return
        }
        this.getAllBranches()
      }, (err) => {

        if (err.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
          return
        }
        alert(this.utilService.getErrorString(err));
        console.error(err.error.error)
      })
  }

  // =======================================================Uploading One Pager and terms and condition================================================ //

  //On uplaoding one pager
  onOnePagerSelect(event: any) {
    this.onePagerDocStatus = ""
    let files = event.target.files
    if (files && files.length) {
      let file = files[0]
      let err = this.fileOps.isDocumentFileValid(file)
      if (err != null) {
        this.onePagerDocStatus = `<p><span>&#10060;</span> ${err}</p>`
        return
      }
      // Upload one pager if it is present.]
      this.isOnePagerFileUploading = true
      this.fileOps.uploadOnePager(file).subscribe((data: any) => {
        // console.log(data);
        this.companyBranchForm.markAsDirty()
        this.companyBranchForm.patchValue({
          onePager: data
        })
        this.onePagerDisplayedFileName = file.name
        this.isOnePagerFileUploading = false
        this.isOnePagerUploadedToServer = true
        this.onePagerDocStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
      }, (error) => {
        this.isOnePagerFileUploading = false
        this.onePagerDocStatus = `<p><span>&#10060;</span> ${error}</p>`
      })
    }
  }

  //On uplaoding one pager
  onTermsAndConditionSelect(event: any) {
    this.termsDocStatus = ""
    let files = event.target.files
    if (files && files.length) {
      let file = files[0]
      let err = this.fileOps.isDocumentFileValid(file)
      if (err != null) {
        this.termsDocStatus = `<p><span>&#10060;</span> ${err}</p>`
        return
      }
      // Upload terms and condition if it is present.]
      this.isTermsFileUploading = true
      this.fileOps.uploadTermsAndCondition(file).subscribe((data: any) => {
        this.companyBranchForm.markAsDirty()
        this.companyBranchForm.patchValue({
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


  // =======================================================Company Requirement specific.================================================ //
  companyRequirementForm: FormGroup;

  minRequiredDate: any;

  // sample offer letter.
  isOfferLetterUploadedToServer: boolean
  isOfferLetterFileUploading: boolean
  offerLetterDocStatus: string
  offerDisplayedFileName: string

  designationList: any[];
  allQualifications: any[];
  allColleges: any[];
  allUniversities: any[];
  allPersonalityTypes: any[];
  talentRatingList: any[];

  // ck-editor
  ckConfig: any

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

  getAllRequirementComponents(): void {
    this.getDesignationList()
    this.getAllQualifications()
    // this.getAllCollegeNames()
    this.getAllUniversities()
    this.getAllPersonalityTypes()
    this.getTalentRatingList()
  }


  getDesignationList(): void {
    this.generalService.getDesignations().subscribe((data: any[]) => {
      this.designationList = data
    }, (err) => {
      console.error(err)
    })
  }

  getTalentRatingList(): void {
    this.generalService.getGeneralTypeByType("talent_rating").subscribe((data: any[]) => {
      this.talentRatingList = data;
    }, (err) => {
      console.error(err)
    })
  }

  getAllPersonalityTypes(): void {
    this.generalService.getGeneralTypeByType("personality_type").subscribe((data: any[]) => {
      this.allPersonalityTypes = data;
    }, (err) => {
      console.error(err)
    })
  }

  getAllQualifications(): void {
    let queryParams: any = {
      limit: -1,
      offset: 0,
    }
    this.degreeService.getAllDegrees(queryParams).subscribe((data: any) => {
      this.allQualifications = data.body
    }, (err) => {
      console.error(err)
    })
  }

  // getAllCollegeNames(): void {
  //       this.generalService.getAllCollegeNamesList().subscribe((data: any) => {
  //             this.allColleges = data.body
  //       }, (err) => {
  //             console.error(err)
  //       })
  // }

  getAllUniversities(): void {
    this.generalService.getAllUniversities().subscribe((data: any) => {
      this.allUniversities = data
    }, (err) => {
      console.error(err)
    })
  }

  // Create company form.
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
      designationID: new FormControl(null, [Validators.required]),

      // Objects.
      companyBranchID: new FormControl(null),
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

  onAddRequirementClick(branchID: string): void {
    this.showExperienceFields = false
    this.isViewClicked = false
    this.getAllRequirementComponents()
    this.createCompanyRequirementForm()
    this.companyRequirementForm.get('companyBranchID').setValue(branchID)
    this.openCompanyRequirementModal(this.companyRequirementModal, "xl")
  }

  // Used to dismiss modal.
  dismissFormModal(modal: NgbModalRef) {
    if (this.isTermsFileUploading || this.isOnePagerFileUploading) {
      alert("Please wait till file is being uploaded")
      return
    }
    if (this.isTermsUploadedToServer || this.isOnePagerUploadedToServer) {
      if (!confirm("Uploaded file will be deleted.\nAre you sure you want to close?")) {
        return
      }
    }
    modal.dismiss()
    this.isTermsUploadedToServer = false
    this.isOnePagerUploadedToServer = false
    this.termsDisplayedFileName = "Select file"
    this.onePagerDisplayedFileName = "Select file"
    this.termsDocStatus = ""
    this.onePagerDocStatus = ""
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

  get jobLocation() {
    return this.companyRequirementForm.get('jobLocation') as FormGroup
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

  addCompanyRequirement(): void {
    this.spinnerService.loadingMessage = "Adding company requirement";

    let companyRequirement: IRequirement = this.companyRequirementForm.value
    this.setFormFields(companyRequirement)
    console.log(companyRequirement);
    this.companyRequirementService.addCompanyRequirement(companyRequirement).subscribe((data: string) => {
      this.modalRef.close();
      this.companyRequirementForm.reset();

      // alert("requirement added with id:" + data)
      alert("Requirement successfully added")
    }, (error) => {
      console.error(error);

      if (error.error) {
        alert(error.error)
        return
      }
      alert(error.error?.error)
    });
  }

  // Extract id from objects and give it to company requirement.
  setFormFields(companyRequirement: IRequirement): void {
    if (this.companyRequirementForm.get('salesPerson')?.value) {
      companyRequirement.salesPersonID = this.companyRequirementForm.get('salesPerson').value.id
      delete companyRequirement['salesPerson']
    }
    if (this.companyRequirementForm.get('companyBranch')?.value) {
      companyRequirement.companyID = this.companyRequirementForm.get('companyBranch').value.id
      delete companyRequirement['companyBranch']
    }

    // designation is compulsory
    companyRequirement.designationID = this.companyRequirementForm.get("designation")?.value.id
    delete companyRequirement['designation']

    // companyRequirement.requiredBefore = new Date(this.companyRequirementForm.get("requiredBefore")?.value).toISOString()
    // companyRequirement.requiredFrom = new Date(this.companyRequirementForm.get("requiredFrom")?.value).toISOString()
  }


  //open company add/update modal
  openCompanyDetailModal(companyDetailModal: any, modalSize?: string): void {

    this.isTermsUploadedToServer = false
    this.termsDisplayedFileName = "Select file"
    this.termsDocStatus = ""

    this.isOnePagerUploadedToServer = false
    this.onePagerDisplayedFileName = "Select file"
    this.onePagerDocStatus = ""

    this.isRequirementTermsUploadedToServer = false
    this.requirementTermsDisplayedFileName = "Select file"
    this.requirementTermsDocStatus = ""

    this.openModal(companyDetailModal, modalSize)

  }

  //open company requirement add modal 
  openCompanyRequirementModal(companyRequirementModal: any, modalSize?: string): void {

    this.isTermsUploadedToServer = false
    this.termsDisplayedFileName = "Select file"
    this.termsDocStatus = ""

    this.isOnePagerUploadedToServer = false
    this.onePagerDisplayedFileName = "Select file"
    this.onePagerDocStatus = ""

    this.isRequirementTermsUploadedToServer = false
    this.requirementTermsDisplayedFileName = "Select file"
    this.requirementTermsDocStatus = ""

    this.isOfferLetterUploadedToServer = false
    this.isOfferLetterFileUploading = false
    this.offerDisplayedFileName = "Select File"
    this.offerLetterDocStatus = ""

    this.openModal(companyRequirementModal, modalSize)
  }

  openModal(modalContent: any, modalSize?: string): void {

    if (modalSize == undefined) {
      modalSize = "xl"
    }

    this.modalRef = this.modalService.open(modalContent, {
      ariaLabelledBy: 'modal-basic-title',
      keyboard: false,
      backdrop: 'static',
      size: modalSize
    });
    /*this.modalRef.result.then((result) => {
    }, (reason) => {

    });*/
  }

  // validate company form
  validateCompanyBranchForm(): void {

    if (this.isOnePagerFileUploading || this.isTermsFileUploading) {
      alert("Please wait while file is getting uploaded.")
      return
    }

    if (this.companyBranchForm.invalid) {
      this.companyBranchForm.markAllAsTouched();
      return
    }

    this.modalAction();
  }

  validateFormFields(): void {

    this.companyRequirementForm.get('maximumPackage').setValidators([Validators.required,
    Validators.min(this.companyRequirementForm.get('minimumPackage')?.value),
    Validators.max(100000000000)])

    this.companyRequirementForm.get('maximumPackage').updateValueAndValidity()
  }

  //validate company requirement form
  validateCompanyRequirementForm(): void {
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

  onOfferLetterClick(event: any) {
    this.offerLetterDocStatus = ""
    let files = event.target.files
    if (files && files.length) {
      let file = files[0]
      let err = this.fileOps.isDocumentFileValid(file)
      if (err != null) {
        this.offerLetterDocStatus = `<p><span>&#10060;</span> ${err}</p>`
        return
      }
      // Upload sample offer letter.
      this.isOfferLetterFileUploading = true
      this.fileOps.uploadOfferLetter(file).subscribe((data: any) => {
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


  // =====================================================Company Requirement Ends==================================================== //




  // =====================================================Components==================================================== //


  getAllDomains(): void {
    this.companyService.getAllDomains().subscribe((data: IDomain[]) => {
      this.allDomains = data
    }, (err) => {
      console.error(err)
    })
  }

  getAllTechnologies(event?: any): void {
    // this.generalService.getTechnologies().subscribe((data: ITechnology[]) => {
    //       this.allTechnologies = data
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
      this.allTechnologies = []
      this.allTechnologies = this.allTechnologies.concat(response.body)
    }, (err) => {
      console.error(err)
    }).add(() => {
      this.isTechLoading = false
    })
  }

  getAllSalesPeople(): void {
    this.generalService.getSalesPersonList().subscribe((data: any) => {
      this.salesPeople = data.body
    }, (err) => {
      console.error(err)
    })
  }

  getStateList(countryID: string): void {
    if (this.companyRequirementForm) {
      this.companyRequirementForm.get('jobLocation')?.get('state')?.reset()
      this.companyRequirementForm.get('jobLocation')?.get('state')?.disable()
    }
    if (this.companyBranchForm) {
      this.companyBranchForm.get('state')?.reset()
      this.companyBranchForm.get('state')?.disable()
    }

    this.stateList = [] as IState[]

    if (countryID) {
      this.stateList = []
      this.generalService.getStatesByCountryID(countryID).subscribe((data: IState[]) => {
        this.stateList = data
        if (this.stateList.length > 0) {
          if (!this.isViewClicked && this.companyRequirementForm) {
            this.companyRequirementForm.get('jobLocation')?.get('state')?.enable()
          }
          if (!this.isViewClicked && this.companyBranchForm) {
            this.companyBranchForm.get('state')?.enable()
          }
        }
      }, (err) => {
        console.error(err)
      })
    }
  }

  getAllStates(): void {
    this.generalService.getStates().subscribe((data: IState[]) => {
      this.allStates = data;
    }, (err) => {
      console.error(err)
    })
  }

  getAllCountries(): void {
    this.generalService.getCountries().subscribe((data: ICountry[]) => {
      this.allCountries = data.sort();
    }, (err) => {
      console.error(err)
    })
  }

  getAllRatings(): void {
    this.ratingList = [];
    this.generalService.getGeneralTypeByType("company_rating").subscribe((data: any[]) => {
      for (let rating of data) {
        this.ratingList.push(+rating.value)
      }
      this.ratingList = this.ratingList.sort()
    }, (err) => {
      console.error(err)
    })
  }


  // =====================================================Components End==================================================== //

  // addCompany():void {
  //       this.spinnerService.loadingMessage = "Adding company";
  //    
  //       this.companyService.addCompany(this.companyBranchForm.value).subscribe((data: string) => {
  //             this.modalRef.close();
  //             this.getAllBranches()
  //             this.companyBranchForm.reset()
  //             alert("company added with id:" + data)
  //       }, (error) => {
  //             console.log(error);
  //          
  //             if (error.error) {
  //                   alert(error.error)
  //                   return
  //             }
  //             alert(error.statusText)
  //       });
  // }



  // // On Update Button Click in Table.
  // updateCompanyForm(companyID: string, companyDetailModal: any):void {
  //       this.openCompanyDetailModal(companyDetailModal);
  //       this.spinnerService.loadingMessage = "Getting company details";
  //    
  //       this.setVariableForUpdateCompany();
  //       this.companyService.getCompanyByID(companyID).subscribe((data: ICompany) => {
  //             this.company = data[0]
  //          
  //             this.updateForm();
  //       }, (error) => {
  //             console.log(error);
  //          
  //             alert(error.error)
  //       })
  // }


  // setVariableForAddCompany():void {
  //       this.updateVariable(this.createCompanyForm, "Add New Company", "Add New Company", this.addCompany)
  // }


  // // Delete Company Branch from company Form.
  // deleteBranchInForm(index: number):void {
  //       this.companyForm.markAsDirty();
  //       this.companyBranches.removeAt(index);
  // }


  // // On Update Button Click in Table.
  // updateCompanyForm(branch: ICompany):void {
  //       this.setVariableForUpdateCompanyBranch();
  //       this.branch = branch
  //       this.updateForm();
  // }


  // onAddNewCompanyButtonClick(companyDetailModal: any):void {
  //       this.openCompanyDetailModal(companyDetailModal);
  //       this.setVariableForAddCompany();
  //       this.formHandler();
  // }


  // // Create company form.
  // createCompanyForm():void {
  //       const URLregex = '(https?://)?([\\da-z.-]+)\\.([a-z.]{2,6})[/\\w .-]*/?';
  //       this.companyForm = this.formBuilder.group({
  //             id: new FormControl(null),
  //             code: new FormControl(null),
  //             name: new FormControl(null, [Validators.required, Validators.maxLength(100)]),
  //             website: new FormControl(null,
  //                   [Validators.pattern(URLregex)]),
  //             branches: this.formBuilder.array([])
  //       });
  //       this.addNewBranchInForm();
  // }


  // Gets companyBranches formBuilder array.
  // get companyBranches() {
  //       return this.companyForm.get('branches') as FormArray
  // }


  // // Quill config.
  // quillModuleConfig = {
  //       // **** Comments include all option. Please do not remove them.
  //       toolbar: [
  //             ['bold', 'italic', 'underline', 'strike'],        // toggled buttons
  //             // ['blockquote', 'code-block'],

  //             // [{ 'header': 1 }, { 'header': 2 }],           // custom button values
  //             [{ 'list': 'ordered' }, { 'list': 'bullet' }],
  //             // [{ 'script': 'sub' }, { 'script': 'super' }],      // superscript/subscript
  //             [{ 'indent': '-1' }, { 'indent': '+1' }],          // outdent/indent
  //             // [{ 'direction': 'rtl' }],                         // text direction

  //             // [{ 'size': ['small', false, 'large', 'huge'] }],  // custom dropdown
  //             [{ 'header': [1, 2, 3, 4, 5, 6, false] }],

  //             // [{ 'color': [] }, { 'background': [] }],   // dropdown with defaults from theme
  //             // [{ 'font': [] }],
  //             [{ 'align': [] }],

  //             ['clean']                                 // remove formatting button
  //       ],

  // }

  // quillEditorStyle = {
  //       height: '100px',
  //       fontSize: 'large',
  //       backgroundColor: '#ffffff',
  //       // borderColor: '#FC6E26',
  //       // borderWidth: "1.5px",
  // }

  // quillCommentEditorStyle = {
  //       height: '80px',
  //       fontSize: 'large',
  //       backgroundColor: '#ffffff',
  //       // borderColor: '#FC6E26',
  //       // borderWidth: "1.5px",
  // }

}
