import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, Validators, FormControl, FormArray } from '@angular/forms';
import { IState, ICompany, ICompanyBranch, ICountry, CompanyService, IDomain, ITechnology } from 'src/app/service/company/company.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { NgbModal, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { IPermission } from 'src/app/service/menu/menu.service';
import { ActivatedRoute, Router } from '@angular/router';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { Role, UrlConstant } from 'src/app/service/constant';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { TechnologyService } from 'src/app/service/technology/technology.service';
import { DegreeService } from 'src/app/service/degree/degree.service';
import { requirementDateValidator } from 'src/app/Validators/custom.validators';
import { CompanyRequirementService, IRequirement } from 'src/app/service/company/company-requirement/company-requirement.service';
import { UrlService } from 'src/app/service/url.service';

@Component({
      selector: 'app-company-master',
      templateUrl: './company-master.component.html',
      styleUrls: ['./company-master.component.css']
})
export class CompanyMasterComponent implements OnInit {

      // Company.
      company: ICompany;
      companyList: ICompany[]
      companyBranchList: ICompanyBranch[];
      companyForm: FormGroup
      addBranch: boolean
      selectedCompanyID: string
      isAddCompanyClicked: boolean
      isCompanyListFound: boolean
      isTechLoading: boolean

      // Required Componenets/Fields.
      allDomains: IDomain[];
      allTechnologies: ITechnology[];
      salesPeople: any[];
      allStates: IState[];
      stateList: IState[]
      allCountries: ICountry[];
      ratingList: number[];

      // Search.
      searchCompanyForm: FormGroup;
      searchFormValue: any;
      searched: boolean;
      showSearch: boolean;

      // Pagination.
      limit: number;
      offset: number;
      currentPage: number;
      totalCompanies: number;
      paginationString: string
      techLimit: number
      techOffset: number

      // Modal.
      modalHeader: string;
      modalButton: string;
      selectedCompanyBranch: ICompanyBranch;
      selectedCompany: ICompany;
      modalAction: () => void
      formHandler: () => void;
      modalRef: any;
      isViewClicked: boolean
      isUpdateClicked: boolean

      // credentials
      permission: IPermission

      // one pager
      isOnePagerUploadedToServer: boolean[]
      isOnePagerFileUploading: boolean[]
      onePagerDocStatus: string[]
      onePagerDisplayedFileName: string[]

      // terms and conditions
      isTermsUploadedToServer: boolean[]
      isTermsFileUploading: boolean[]
      termsDocStatus: string[]
      termsDisplayedFileName: string[]

      // logo
      logoDisplayedFileName: string
      logoDocStatus: string
      isLogoUploadedToServer: boolean
      isLogoFileUploading: boolean

      // route
      previousUrl: string

      //constant
      private readonly IGNORE_SEARCH_FIELD: string[] = ["limit", "offset"]

      // Modals
      @ViewChild('companyDetailModal') companyDetailModal: any
      @ViewChild('deleteConfirmationModal') deleteConfirmationModal: any
      isOfferLetterUploadedToServer: boolean;
      isOfferLetterFileUploading: boolean;
      requirementOfferDisplayedFileName: string;
      requirementOfferLetterDocStatus: string;
      minValueForMaxExp: number;
      isSalesPerson: boolean;
      isAdmin: boolean;
      showForAdmin: boolean;

      constructor(
            private formBuilder: FormBuilder,
            private companyService: CompanyService,
            private generalService: GeneralService,
            private techService: TechnologyService,
            private localService: LocalService,
            private utilService: UtilityService,
            private spinnerService: SpinnerService,
            private modalService: NgbModal,
            private router: Router,
            private activatedRoute: ActivatedRoute,
            private fileOps: FileOperationService,
            private urlConstant: UrlConstant,
            private degreeService: DegreeService,
            private companyRequirementService: CompanyRequirementService,
            private roleConstant: Role,
            private urlService: UrlService,


      ) {
            this.initializeVariables()
            this.editorConfig()
            this.createForms()
            this.getAllComponents()
      }

      initializeVariables(): void {
            this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.COMPANY_MASTER)
            this.isSalesPerson = this.localService.getJsonValue("roleName") == this.roleConstant.SALES_PERSON
            this.isAdmin = this.localService.getJsonValue("roleName") == this.roleConstant.ADMIN
            this.salesPeople = [];
            this.stateList = []

            this.limit = 5;
            this.offset = 0;
            this.techLimit = 10
            this.techOffset = 0

            this.searched = false
            this.showSearch = false
            this.isViewClicked = false
            this.isUpdateClicked = false
            this.isAddCompanyClicked = false
            this.addBranch = false
            this.isCompanyListFound = true
            this.isTechLoading = false

            // this.minRequiredDate = formatDate(new Date().setDate(new Date().getDate() + 10), 'yyyy/MM/dd', 'en');

            this.isOnePagerUploadedToServer = []
            this.isOnePagerFileUploading = []
            this.onePagerDisplayedFileName = []
            this.onePagerDocStatus = []

            this.isTermsUploadedToServer = []
            this.isTermsFileUploading = []
            this.termsDisplayedFileName = []
            this.termsDocStatus = []

            this.isLogoUploadedToServer = false
            this.isLogoFileUploading = false
            this.logoDisplayedFileName = "Select File"
            this.logoDocStatus = ""


            this.previousUrl = null
      }

      // gets all required components like all states,countries etc
      getAllComponents(): void {
            this.getAllCountries()
            this.getAllRatings()
            this.getAllStates()
            this.getAllSalesPeople()
            this.getAllDomains()
            this.getAllTechnologies()
            this.searchOrGetCompanies()
      }

      createForms(): void {
            this.createSearchCompanyForm()
            this.createAddCompanyForm()
      }


      get ongoingOperations() {
            return this.spinnerService.ongoingOperations
      }

      ngOnInit() {
      }

      // Handles pagination
      changePage(pageNumber: number): void {
            // $event will be the page number & offset will be 1 less than it.

            // this.offset = $event - 1
            // this.currentPage = $event
            // this.getAllCompanies()
            this.searchCompanyForm.get("offset").setValue(pageNumber - 1)

            this.limit = this.searchCompanyForm.get("limit").value
            this.offset = this.searchCompanyForm.get("offset").value
            this.searchCompany()
      }

      compareFn(c1: any, c2: any): boolean {
            return c1 && c2 ? c1.id === c2.id : c1 === c2;
      }

      // Compare two Object.
      compareFunc(ob1: any, ob2: any): any {
            if (ob1 == null && ob2 == null) {
                  return true
            }
            return ob1 && ob2 ? ob1.id == ob2.id : ob1 == ob2;
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
                  companyName: new FormControl(null, [Validators.maxLength(100)]),
                  limit: new FormControl(this.limit),
                  offset: new FormControl(this.offset),
                  // city: new FormControl(null, [Validators.maxLength(50), Validators.pattern("^[a-zA-Z]+([a-zA-Z ]?)+")]),
                  // companyRating: new FormControl(null),
                  // coordinatorID: new FormControl(null),
                  // state: new FormControl(null),
                  // hrHeadName: new FormControl(null, [Validators.maxLength(100), Validators.pattern("^[a-zA-Z]+([a-zA-Z. ]?)+")]),
                  // numberOfEmployees: new FormControl(null, [Validators.max(1000000), Validators.pattern("^[1-9][0-9]*$")]),
                  // salesPersonID: new FormControl(null),
                  // domains: new FormControl(null),
                  // technologies: new FormControl(null),
            })
      }

      createAddCompanyForm(): void {
            this.createCompanyForm()
            if (this.isAddCompanyClicked) {
                  this.addNewBranchInForm()
            }

      }

      // Create company form.
      createCompanyForm(): void {
            // const URLregex = "^(https?:\/\/)[a-z0-9-]+(\.[a-z0-9-]+)+(\/[a-z0-9-]+)*\/?$";
            const URLregex = /^(http(s)?:\/\/)[\w.-]+(\.[\w\.-]+)+[\w\-\._~:?#[\]@!\/$&'\(\)\*\+,;=.]+$/

            this.companyForm = this.formBuilder.group({
                  id: new FormControl(null),
                  code: new FormControl(null),
                  companyName: new FormControl(null,
                        [Validators.required, Validators.maxLength(200), Validators.pattern("^[a-zA-Z]+([a-zA-Z. ]?)+")]),
                  about: new FormControl(null, [Validators.maxLength(2000)]),
                  logo: new FormControl(null, [Validators.required]),
                  website: new FormControl(null,
                        [Validators.pattern(URLregex)]),
                  branches: this.formBuilder.array([])
            });
            // this.addNewBranchInForm();
      }


      // Gets companyBranches formBuilder array.
      get companyBranches() {
            return this.companyForm.get('branches') as FormArray
      }

      // Add New Branch.
      addNewBranchInForm(): void {
            this.companyBranches.push(this.formBuilder.group({
                  id: new FormControl(null),
                  branchName: new FormControl(null, [Validators.required, Validators.pattern("^[a-zA-Z]+([a-zA-Z. ]?)+")]),
                  code: new FormControl(null),
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
                  termsAndConditions: new FormControl(null),
            }))
            this.onePagerDisplayedFileName.push("Select file")
            this.termsDisplayedFileName.push("Select file")
            this.logoDisplayedFileName = "Select File"

            this.onePagerDocStatus.push("")
            this.termsDocStatus.push("")
            this.logoDocStatus = ""

            this.isOnePagerFileUploading.push(false)
            this.isTermsFileUploading.push(false)
            this.isLogoFileUploading = false

            this.isOnePagerUploadedToServer.push(false)
            this.isTermsUploadedToServer.push(false)
            this.isLogoUploadedToServer = false

      }

      getBranchesByCompany(companyID: string): void {
            // console.log("company id ->" + companyID)
            // this.router.navigate([this.urlConstant.COMPANY_BRANCH], {
            //       queryParams: {
            //             "companyID": companyID
            //       }
            // }).catch(err => {
            //       console.error(err)

            // })
            this.urlService.setUrlFrom(this.constructor.name, this.router.url)
            this.router.navigate([this.urlConstant.COMPANY_BRANCH], {
                  queryParams: {
                        "companyID": companyID
                  }
            }).catch(err => {
                  console.error(err)

            });

      }

      onAddBranchClick(companyID: string): void {
            this.addBranch = true
            this.isViewClicked = false
            this.isAddCompanyClicked = true
            this.openModal(this.companyDetailModal)
            this.setVariableForAddCompanyBranch()
            this.formHandler();
            this.selectedCompanyID = companyID
            // this.companyFormValidators()
            // console.log(this.companyForm.value);
      }

      companyFormValidators() {
            if (this.addBranch) {
                  this.companyForm.get('companyName').setValidators(null)
                  this.companyForm.get('website').setValidators(null)
                  this.companyForm.get('about').setValidators(null)
                  this.utilService.updateValueAndValiditors(this.companyForm)
                  return
            }
            // const URLregex = "^(https?:\/\/)[a-z0-9-]+(\.[a-z0-9-]+)+(\/[a-z0-9-]+)*\/?$";
            const URLregex = /^(http(s)?:\/\/)[\w.-]+(\.[\w\.-]+)+[\w\-\._~:?#[\]@!\/$&'\(\)\*\+,;=.]+$/

            this.companyForm.get('companyName').setValidators([Validators.required, Validators.maxLength(100), Validators.pattern("^[a-zA-Z]+([a-zA-Z. ]?)+")])
            this.companyForm.get('website').setValidators([Validators.pattern(URLregex)])
            this.companyForm.get('about').setValidators([Validators.maxLength(2000)])

            this.utilService.updateValueAndValiditors(this.companyForm)
      }

      setVariableForAddCompanyBranch(): void {
            this.updateVariable(this.createAddCompanyForm, "Add New Branch", "Add New Branch", this.addCompanyBranch)
      }

      onViewCompanyClick(branch: ICompany): void {
            this.createCompanyForm()
            this.isViewClicked = true
            this.isUpdateClicked = false
            this.isAddCompanyClicked = false
            this.addBranch = false
            this.updateVariable(this.createCompanyForm, "Company Details", "Company", this.addCompany)
            this.company = branch
            console.log(this.company);
            // this.updateForm()
            this.companyForm.patchValue(branch)
            this.companyForm.disable()
            this.openModal(this.companyDetailModal, "lg")
      }

      onUpdateCompanyClick(): void {
            // this.setVariableForUpdateCompany();
            this.modalHeader = "Update Company";
            this.modalButton = "Update";
            this.modalAction = this.updateCompany;
            this.companyForm.enable()
            this.isViewClicked = false
            this.isUpdateClicked = true
            this.isAddCompanyClicked = false
            if (this.company.logo) {
                  this.logoDisplayedFileName = `<a href=${this.company.logo} target="_blank">Company Logo</a>`
            }
      }

      onAddNewCompanyButtonClick(): void {
            this.addBranch = false
            this.isViewClicked = false
            this.isAddCompanyClicked = true
            this.openModal(this.companyDetailModal);
            this.setVariableForAddCompany();
            this.formHandler();
            this.companyFormValidators()
      }

      onDeleteCompanyClick(company: any) {
            this.selectedCompany = company
            this.openModal(this.deleteConfirmationModal, "md")
      }

      // Delete Company Branch from company Form.
      deleteBranchInForm(index: number): void {
            this.companyForm.markAsDirty();
            this.companyBranches.removeAt(index);
      }

      // Used to dismiss modal.
      dismissFormModal(modal: NgbModalRef) {

            if (this.isTermsFileUploading.includes(true) || this.isOnePagerFileUploading.includes(true)
                  || this.isLogoFileUploading) {
                  alert("Please wait till file is being uploaded")
                  return
            }
            if (this.isTermsUploadedToServer.includes(true) || this.isOnePagerUploadedToServer.includes(true)
                  || this.isLogoUploadedToServer) {
                  if (!confirm("Uploaded file will be deleted.\nAre you sure you want to close?")) {
                        return
                  }
            }
            modal.dismiss()
            this.isTermsUploadedToServer = []
            this.isOnePagerUploadedToServer = []
            this.termsDisplayedFileName = []
            this.onePagerDisplayedFileName = []
            this.termsDocStatus = []
            this.onePagerDocStatus = []

            this.logoDocStatus = ""
            this.logoDisplayedFileName = "Select File"
            this.isLogoFileUploading = false
            this.isLogoUploadedToServer = false
      }

      setVariableForAddCompany(): void {
            this.updateVariable(this.createAddCompanyForm, "Add New Company", "Add New Company", this.addCompany)
      }

      // setVariableForUpdateCompany(): void {
      //       this.updateVariable(this.updateForm, "Update Company", "Update Company", this.updateCompany)
      // }

      updateVariable(formaction: () => void, modalheader: string, modalbutton: string, modalaction: () => void): void {
            this.formHandler = formaction;
            this.modalHeader = modalheader;
            this.modalButton = modalbutton;
            this.modalAction = modalaction;
      }

      setPaginationString() {
            this.paginationString = ''
            let limit = this.searchCompanyForm.get('limit').value
            let offset = this.searchCompanyForm.get('offset').value

            let start: number = limit * offset + 1
            let end: number = +limit + limit * offset

            if (this.totalCompanies < end) {
                  end = this.totalCompanies
            }
            if (this.totalCompanies == 0) {
                  this.paginationString = ''
                  return
            }
            this.paginationString = `${start} - ${end}`
      }

      openModal(modalContent: any, modalSize?: string): void {

            this.isOnePagerFileUploading = []
            this.isTermsFileUploading = []
            this.isLogoFileUploading = false

            this.isOnePagerUploadedToServer = []
            this.isTermsUploadedToServer = []
            this.isLogoUploadedToServer = false

            this.onePagerDocStatus = []
            this.termsDocStatus = []
            this.logoDocStatus = ""

            this.onePagerDisplayedFileName = []
            this.termsDisplayedFileName = []
            this.logoDisplayedFileName = "Select File"

            if (modalSize == undefined) {
                  modalSize = "xl"
            }

            this.modalRef = this.modalService.open(modalContent, {
                  ariaLabelledBy: 'modal-basic-title',
                  keyboard: false,
                  backdrop: 'static',
                  size: modalSize
            });
            /*this.modalRef.result.subscribe((result) => {
            }, (reason) => {

            });*/
      }

      validateFormFields(): void {
            if (this.addBranch) {
                  this.companyForm.get("companyName").clearValidators()
                  this.companyForm.get("logo").clearValidators()
            } else {
                  this.companyForm.get("companyName").setValidators(
                        [Validators.required, Validators.maxLength(200), Validators.pattern("^[a-zA-Z]+([a-zA-Z. ]?)+")])
                  this.companyForm.get("logo").setValidators([Validators.required])

            }
            this.utilService.updateValueAndValiditors(this.companyForm)
      }

      // validate company form
      validateCompanyForm(): void {
            this.validateFormFields()
            // console.log(this.companyForm.controls);

            // for (let index = 0; index < this.companyForm.get('branches').value.length; index++) {
            // console.log(this.companyForm.get('branches').value[index]);
            if (this.isOnePagerFileUploading.includes(true) || this.isTermsFileUploading.includes(true)
                  || this.isLogoFileUploading) {
                  alert("Please wait while file is getting uploaded.")
                  return
            }
            // }

            if (this.companyForm.invalid) {
                  this.companyForm.markAllAsTouched();
                  return
            }

            if (this.addBranch) {
                  this.addCompanyBranch()
                  return
            }

            this.modalAction();
      }

      resetSearchAndGetAll(): void {
            this.searched = false
            this.resetSearchCompanyForm()
            this.searchFormValue = null
            this.router.navigate([this.urlConstant.COMPANY_MASTER])
            this.changePage(1)
      }

      resetSearchCompanyForm(): void {
            this.limit = this.searchCompanyForm.get("limit").value
            this.offset = this.searchCompanyForm.get("offset").value
            // this.currentPage = this.searchCompanyForm.get("currentPage").value

            this.searchCompanyForm.reset({
                  limit: this.limit,
                  offset: this.offset,
                  // currentPage: new FormControl(this.currentPage),
            })
      }
      searchCompany(): void {
            this.searchFormValue = { ...this.searchCompanyForm?.value }
            this.router.navigate([], {
                  relativeTo: this.activatedRoute,
                  queryParams: this.searchFormValue,
            })
            this.searchFormValue = { ...this.searchCompanyForm.value }
            let flag: boolean = true
            for (let field in this.searchFormValue) {
                  if (!this.searchFormValue[field]) {
                        delete this.searchFormValue[field];
                  } else {
                        if (!this.IGNORE_SEARCH_FIELD.includes(field)) {
                              this.searched = true
                        }
                        // this.searched = true
                        flag = false
                  }
            }
            // No API call on empty search.
            if (flag) {
                  return
            }
            // this.changePage(1)
            this.getAllCompanies()
      }

      searchOrGetCompanies(): void {
            let queryParams = this.activatedRoute.snapshot.queryParams
            if (this.utilService.isObjectEmpty(queryParams)) {
                  // this.getAllCompanies()
                  this.changePage(1)
                  return
            }
            this.searchCompanyForm.patchValue(queryParams)
            this.searchCompany()

            // this.activatedRoute.queryParams.subscribe((param) => {
            //       // console.log(param)
            //       this.searchFormValue = { ...param }
            //       for (let field in this.searchFormValue) {
            //             if (!this.searchFormValue[field]) {
            //                   delete this.searchFormValue[field];
            //             } else {
            //                   this.searched = true
            //             }
            //       }
            //       this.searchCompanyForm.patchValue(this.searchFormValue)
            //       this.changePage(1)
            // })
      }

      // =============================================================CRUD=============================================================

      getAllCompanies(): void {
            // getAllSearchedBranches calls same service and passes same arguments to it.
            // Calls search if a search has been made
            // if (this.searched) {
            //       this.getAllSearchedBranches()
            //       return
            // }
            console.log(this.searchFormValue);

            this.spinnerService.loadingMessage = "Getting companies";

            this.companyList = []
            this.totalCompanies = 0
            this.isCompanyListFound = true

            this.companyService.getAllCompanies(this.searchFormValue).subscribe(data => {
                  this.companyList = data.body;
                  // this.removeCompanyBranches()
                  // console.log(this.companyList);
                  // this.assignOtherFields()
                  this.totalCompanies = parseInt(data.headers.get('X-Total-Count'))
                  this.setPaginationString()
                  if (this.totalCompanies == 0) {
                        this.isCompanyListFound = false
                  }

            }, (error: any) => {
                  this.totalCompanies = 0
                  this.setPaginationString()
                  this.isCompanyListFound = false
                  console.error(error);

                  if (error.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                        return
                  }
                  alert(error.error.error)
            });
      }

      getAllSearchedBranches(): void {
            this.spinnerService.loadingMessage = "Searching companies"

            this.companyList = []
            this.totalCompanies = 0
            this.isCompanyListFound = true
            this.companyService.getAllCompanies(this.searchFormValue).subscribe(data => {
                  this.companyList = data.body;
                  // this.removeCompanyBranches()
                  // this.assignOtherFields()
                  this.totalCompanies = parseInt(data.headers.get('X-Total-Count'))
                  this.setPaginationString()
                  if (this.totalCompanies == 0) {
                        this.isCompanyListFound = false
                  }

            }, (error: any) => {
                  this.totalCompanies = 0
                  this.setPaginationString()
                  this.isCompanyListFound = false
                  console.error(error);

                  if (error.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                        return
                  }
                  alert(error.error.error)
            });
      }


      addCompany(): void {
            this.spinnerService.loadingMessage = "Adding company";

            this.companyService.addCompany(this.companyForm.value).subscribe((data: string) => {
                  this.modalRef.close();
                  this.getAllCompanies()
                  this.companyForm.reset()
                  alert("Company successfully added")
            }, (error: any) => {
                  console.error(error);

                  if (error.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                        return
                  }
                  alert(error.error.error)
            })
      }

      addCompanyBranch(): void {
            this.spinnerService.loadingMessage = "Adding company branch";

            let companyBranches = this.companyForm.get("branches")?.value
            console.log(companyBranches);
            console.log(this.selectedCompanyID);

            // for (let index = 0; index < this.companyBranches.value.length; index++) {
            //       this.companyBranches.value[index].mainBranch = false
            // }
            console.log(this.companyBranches.value);
            this.companyService.addCompanyBranches(this.companyBranches.value, this.selectedCompanyID).subscribe((data: string) => {
                  this.modalRef.close();
                  this.getAllCompanies()
                  this.companyForm.reset()
                  alert("Branch successfully added")
            }, (error: any) => {
                  console.error(error);

                  if (error.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                        return
                  }
                  alert(error.error?.error)
            })
      }

      // Update Company.
      updateCompany(): void {
            this.spinnerService.loadingMessage = "Updating company";

            // let company = this.companyForm.value
            // for (let branch of company.branches) {
            //       branch.companyID = company.id
            // }

            this.companyService.updateCompany(this.companyForm.value).subscribe((data: string) => {
                  this.modalRef.close();
                  this.getAllCompanies()
                  this.companyForm.reset()
                  // alert(data)
                  alert("Company successfully updated")
            }, (error: any) => {
                  console.error(error);

                  if (error.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                        return
                  }
                  alert(error.error.error)
            })

      }

      deleteCompany(): void {
            this.modalRef.close();
            this.spinnerService.loadingMessage = "Deleting company";

            this.companyService.deleteCompany(this.selectedCompany.id).subscribe((data: any) => {
                  this.getAllCompanies()
                  // alert(data)
                  alert("Company successfully deleted")
            }, (error: any) => {
                  console.error(error);

                  if (error.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                        return
                  }
                  alert(error.error.error)
            });
      }

      deleteCompanyBranch(): void {
            this.modalRef.close();
            this.spinnerService.loadingMessage = "Deleting company";

            this.companyService.deleteCompanyBranch(this.selectedCompanyBranch.companyID,
                  this.selectedCompanyBranch.id).subscribe((data: any) => {
                        this.getAllCompanies()
                        // alert(data)
                        alert("Branch successfully deleted")
                  }, (error: any) => {
                        console.error(error);

                        if (error.statusText.includes('Unknown')) {
                              alert("No connection to server. Check internet.")
                              return
                        }
                        alert(error.error.error)
                  });
      }

      // =======================================================Uploading One Pager and terms and condition================================================ //

      //On uplaoding one pager
      onOnePagerSelect(event: any, companyBranch: any, index: number) {
            this.onePagerDocStatus[index] = ""
            // companyBranch.controls['onePagerDocStatus'].setValue("true")
            let files = event.target.files
            if (files && files.length) {
                  let file = files[0]
                  let err = this.fileOps.isDocumentFileValid(file)
                  if (err != null) {
                        this.onePagerDocStatus[index] = `<p><span>&#10060;</span> ${err}</p>`
                        return
                  }
                  // Upload one pager if it is present.]
                  this.isOnePagerFileUploading[index] = true
                  // companyBranch.controls['isOnePagerUploading'].setValue(true)
                  this.fileOps.uploadOnePager(file).subscribe((data: any) => {

                        this.companyForm.markAsDirty()
                        companyBranch.patchValue({
                              onePager: data,
                              // isOnePagerUploading: false,
                              // onePagerDisplayedFileName: file.name,
                              // onePagerDocStatus: "<p><span class='green'>&#10003;</span> File uploaded.</p>"
                        })
                        // console.log(companyBranch);
                        this.onePagerDisplayedFileName[index] = file.name
                        this.isOnePagerFileUploading[index] = false
                        this.isOnePagerUploadedToServer[index] = true
                        this.onePagerDocStatus[index] = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
                  }, (error) => {
                        // companyBranch.patchValue({
                        //       isOnePagerUploading: false,
                        //       termsDocStatus: `<p><span>&#10060;</span> ${error}</p>`
                        // })
                        this.isOnePagerFileUploading[index] = false
                        this.onePagerDocStatus[index] = `<p><span>&#10060;</span> ${error}</p>`
                  })
            }
      }

      //On uplaoding terms and conditions
      onTermsAndConditionSelect(event: any, companyBranch: any, index: number) {
            this.termsDocStatus[index] = ""
            // companyBranch.controls['termsDocStatus'].setValue("true")
            let files = event.target.files
            if (files && files.length) {
                  let file = files[0]
                  let err = this.fileOps.isDocumentFileValid(file)
                  if (err != null) {
                        this.termsDocStatus[index] = `<p><span>&#10060;</span> ${err}</p>`
                        return
                  }
                  // Upload terms and condition if it is present.]
                  this.isTermsFileUploading[index] = true
                  // companyBranch.controls['isTermsAndConditionsUploading'].setValue(true)

                  this.fileOps.uploadTermsAndCondition(file).subscribe((data: any) => {
                        this.companyForm.markAsDirty()
                        companyBranch.patchValue({
                              termsAndConditions: data,
                              // isTermsAndConditionsUploading: false,
                              // termsDisplayedFileName: file.name,
                              // termsDocStatus: "<p><span class='green'>&#10003;</span> File uploaded.</p>"
                        })
                        // console.log(companyBranch);

                        this.termsDisplayedFileName[index] = file.name
                        this.isTermsFileUploading[index] = false
                        this.isTermsUploadedToServer[index] = true
                        this.termsDocStatus[index] = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
                  }, (error) => {
                        // companyBranch.patchValue({
                        //       isTermsAndConditionsUploading: false,
                        //       termsDocStatus: `<p><span>&#10060;</span> ${error}</p>`
                        // })
                        this.isTermsFileUploading[index] = false
                        this.termsDocStatus[index] = `<p><span>&#10060;</span> ${error}</p>`
                  })
            }
      }

      //On uplaoding logo
      onLogoSelect(event: any) {
            this.companyForm.get("logo").reset()
            this.logoDocStatus = ""
            let files = event.target.files
            if (files && files.length) {
                  let file = files[0]
                  this.logoDisplayedFileName = file.name
                  let err = this.fileOps.isImageFileValid(file)
                  if (err != null) {
                        this.logoDocStatus = `<p><span>&#10060;</span> ${err}</p>`
                        return
                  }
                  // Upload terms and condition if it is present.]
                  this.isLogoFileUploading = true
                  // companyBranch.controls['isTermsAndConditionsUploading'].setValue(true)

                  this.fileOps.uploadLogo(file,
                        this.fileOps.COMPANY_FOLDER + this.fileOps.LOGO_FOLDER)
                        .subscribe((data: any) => {
                              this.companyForm.markAsDirty()
                              this.companyForm.patchValue({
                                    logo: data,
                              })
                              // console.log(data);

                              this.logoDisplayedFileName = file.name
                              this.isLogoFileUploading = false
                              this.isLogoUploadedToServer = true
                              this.logoDocStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
                        }, (error) => {
                              this.isLogoFileUploading = false
                              this.logoDocStatus = `<p><span>&#10060;</span> ${error}</p>`
                        })
            }
      }

      // =======================================================COMPONENTS================================================ //

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

      getAllStates(): void {
            this.generalService.getStates().subscribe((data: IState[]) => {
                  this.allStates = data;
            }, (err) => {
                  console.error(err)
            })
      }

      getStateList(countryID: string): void {
            if (countryID) {
                  this.stateList = []
                  this.generalService.getStatesByCountryID(countryID).subscribe((data: IState[]) => {
                        this.stateList = data;
                  }, (err) => {
                        console.error(err)
                  })
            }
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


      // public findInvalidControls() {
      //       const invalid = []
      //       const controls = this.companyForm.controls
      //       for (const name in controls) {
      //             if (controls[name].invalid) {
      //                   invalid.push(name)
      //             }
      //       }
      //       return invalid
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

      // // On Update Button Click in Table.
      // updateCompanyForm(branch: ICompany):void {
      //       this.setVariableForUpdateCompany();
      //       this.company = branch
      //       this.updateForm();
      // }


      // getAllBranches():void {
      //       // Calls search if a search has been made
      //       // if (this.searched) {
      //       //       this.getAllSearchedBranches()
      //       //       return
      //       // }

      //       this.spinnerService.loadingMessage = "Getting companies";
      //       this.companyService.getAllCompanyBranches(this.limit, this.offset).subscribe(data => {
      //             // console.log(data.body);
      //             this.companyBranchList = data.body;
      //          
      //             this.assignOtherFields()
      //             this.totalCompanies = parseInt(data.headers.get('X-Total-Count'))
      //       }, (error) => {
      //             console.log(error);
      //          
      //             if (error.statusText.includes('Unknown')) {
      //                   alert("No connection to server. Check internet.")
      //             }
      //       });
      // }
      //=====================================================Enquiry====================================================================
      //open company requirement add modal 
      // terms and conditions
      isRequirementTermsUploadedToServer: boolean
      isRequirementTermsFileUploading: boolean
      requirementTermsDocStatus: string
      requirementermsDisplayedFileName: string
      isRequirementOnePagerUploadedToServer: boolean
      requirementOnePagerDisplayedFileName: string
      requirementOnePagerDocStatus: string
      requirementTermsDisplayedFileName: string

      openCompanyRequirementModal(companyRequirementModal: any, modalSize?: string): void {

            this.isRequirementTermsUploadedToServer = false
            this.requirementermsDisplayedFileName = "Select file"
            this.requirementTermsDocStatus = ""

            this.isRequirementOnePagerUploadedToServer = false
            this.requirementOnePagerDisplayedFileName = "Select file"
            this.requirementOnePagerDocStatus = ""

            this.isRequirementTermsUploadedToServer = false
            this.requirementTermsDisplayedFileName = "Select file"
            this.requirementTermsDocStatus = ""

            this.isOfferLetterUploadedToServer = false
            this.isOfferLetterFileUploading = false
            this.requirementOfferDisplayedFileName = "Select File"
            this.requirementOfferLetterDocStatus = ""

            this.openModal(companyRequirementModal, modalSize)
      }

      @ViewChild('companyRequirementFormModal') companyRequirementModal: any

      showExperienceFields: boolean
      companyRequirementForm: FormGroup;

      onAddRequirementClick(companyID: string): void {
            this.showExperienceFields = false
            this.isViewClicked = false
            if (this.isAdmin) {
                  this.showForAdmin = true
            }
            this.getAllRequirementComponents()
            this.createCompanyRequirementForm()
            this.companyRequirementForm.get('companyID').setValue(companyID)
            this.openCompanyRequirementModal(this.companyRequirementModal, "xl")
      }

      getAllRequirementComponents() {
            this.getDesignationList()
            this.getAllQualifications()
            // this.getAllCollegeNames()
            this.getAllUniversities()
            this.getAllPersonalityTypes()
            this.getTalentRatingList()
      }

      designationList: any[]
      getDesignationList(): void {
            this.generalService.getDesignations().subscribe((data: any[]) => {
                  this.designationList = data
            }, (err) => {
                  console.error(err)
            })
      }

      talentRatingList: any[]
      getTalentRatingList(): void {
            this.generalService.getGeneralTypeByType("talent_rating").subscribe((data: any[]) => {
                  this.talentRatingList = data;
            }, (err) => {
                  console.error(err)
            })
      }

      allPersonalityTypes: any[]
      getAllPersonalityTypes(): void {
            this.generalService.getGeneralTypeByType("personality_type").subscribe((data: any[]) => {
                  this.allPersonalityTypes = data;
            }, (err) => {
                  console.error(err)
            })
      }

      allQualifications: any[]
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

      allUniversities: any[]
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
                  designation: new FormControl(null, [Validators.required]),
                  designationID: new FormControl(null),

                  // Objects.
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

      get jobLocation() {
            return this.companyRequirementForm.get('jobLocation') as FormGroup
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
            this.requirementOfferLetterDocStatus = ""
            let files = event.target.files
            if (files && files.length) {
                  let file = files[0]
                  let err = this.fileOps.isDocumentFileValid(file)
                  if (err != null) {
                        this.requirementOfferLetterDocStatus = `<p><span>&#10060;</span> ${err}</p>`
                        return
                  }
                  // Upload sample offer letter.
                  this.isOfferLetterFileUploading = true
                  this.fileOps.uploadOfferLetter(file).subscribe((data: any) => {
                        this.companyRequirementForm.markAsDirty()
                        this.companyRequirementForm.patchValue({
                              sampleOfferLetter: data
                        })
                        this.requirementOfferDisplayedFileName = file.name
                        this.isOfferLetterFileUploading = false
                        this.isOfferLetterUploadedToServer = true
                        this.requirementOfferLetterDocStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
                  }, (error) => {
                        this.isOfferLetterFileUploading = false
                        this.requirementOfferLetterDocStatus = `<p><span>&#10060;</span> ${error}</p>`
                  })
            }
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
            }, (error) => {
                  console.error(error);

                  if (error.error) {
                        alert(error.error)
                        return
                  }
                  alert(error.error?.error)
            });
      }
}

