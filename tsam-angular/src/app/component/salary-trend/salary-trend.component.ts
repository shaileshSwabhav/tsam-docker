import { DatePipe } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { UrlConstant } from 'src/app/service/constant';
import { GeneralService, IDesignation, ISalaryTrend, ITechnologies } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { TechnologyService } from 'src/app/service/technology/technology.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-salary-trend',
  templateUrl: './salary-trend.component.html',
  styleUrls: ['./salary-trend.component.css']
})
export class SalaryTrendComponent implements OnInit {

  // components
  technologyList: ITechnologies[]
  desginationList: IDesignation[]
  companyRatingList: any[]

  // salary-trend
  salaryTrends: ISalaryTrend[]
  salaryTrendForm: FormGroup
  totalSalaryTrend: number
  isSalaryTrend: boolean

  // search
  salaryTrendSearchForm: FormGroup
  searchFormValue: any
  isSearched: boolean
  showSearch: boolean

  // flags
  isOperationUpdate: boolean
  isViewMode: boolean
  isTechLoading: boolean

  // pagination
  limit: number
  currentPage: number
  offset: number
  paginationString: string
  techLimit: number
  techOffset: number

  // modal
  modalRef: NgbModalRef
  @ViewChild('salaryTrendFormModal') salaryTrendFormModal: any
  @ViewChild('deleteConfirmationModal') deleteConfirmationModal: any
  @ViewChild('drawer') drawer: any

  // spinner



  // permission
  permission: IPermission

  constructor(
    private formBuilder: FormBuilder,
    private utilService: UtilityService,
    private localService: LocalService,
    private urlConstant: UrlConstant,
    private generalService: GeneralService,
    private techService: TechnologyService,
    private modalService: NgbModal,
    private spinnerService: SpinnerService,
    private router: Router,
    private route: ActivatedRoute,
    private datePipe: DatePipe
  ) {
    this.initializeVariables()
    this.getAllComponents()
  }

  initializeVariables(): void {
    this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.COMPANY_SALARY_TREND)

    this.limit = 5
    this.offset = 0
    this.totalSalaryTrend = 0
    this.techLimit = 10
    this.techOffset = 0

    this.isSalaryTrend = true
    this.isViewMode = false
    this.isOperationUpdate = false
    this.isSearched = false
    this.showSearch = false
    this.isTechLoading = false

    this.technologyList = []
    this.desginationList = []

    this.companyRatingList = []
    this.salaryTrends = []

    this.createSearchForm()
  }

  getAllComponents(): void {
    this.getTechnologyList()
    this.getDesignationList()
    this.getCompanyRatingList()
    this.searchOrGetSalaryTrend()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  createSearchForm(): void {
    this.salaryTrendSearchForm = this.formBuilder.group({
      fromDate: new FormControl(null),
      toDate: new FormControl(null),
      technology: new FormControl(null),
      designation: new FormControl(null),
      minimumExperience: new FormControl(null, [Validators.min(0)]),
      maximumExperience: new FormControl(null, [Validators.min(1)]),
      companyRating: new FormControl(null),
    })
  }

  createSalaryTrendForm(): void {
    this.salaryTrendForm = this.formBuilder.group({
      id: new FormControl(null),
      date: new FormControl(null, [Validators.required]),
      companyRating: new FormControl(null, [Validators.required]),
      minimumExperience: new FormControl(null, [Validators.required, Validators.min(0)]),
      maximumExperience: new FormControl(null, [Validators.required]),
      minimumSalary: new FormControl(null, [Validators.required, Validators.min(1)]),
      maximumSalary: new FormControl(null, [Validators.required]),
      medianSalary: new FormControl(null, [Validators.required, Validators.min(1)]),
      technology: new FormControl(null, [Validators.required]),
      designation: new FormControl(null, [Validators.required])
    })
  }

  compareFn(optionOne: any, optionTwo: any): boolean {
    if (optionOne == null && optionTwo == null) {
      return true
    }
    return optionOne == optionTwo
  }

  onAddClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.createSalaryTrendForm()
    this.openModal(this.salaryTrendFormModal, "xl")
  }

  onViewClick(salaryTrend: ISalaryTrend): void {
    this.isViewMode = true
    this.isOperationUpdate = false
    this.createSalaryTrendForm()
    this.salaryTrendForm.disable()
    salaryTrend.date = this.datePipe.transform(salaryTrend.date, 'yyyy-MM-dd')
    this.salaryTrendForm.patchValue(salaryTrend)
    this.openModal(this.salaryTrendFormModal, "xl")
  }

  onUpdateClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = true
    // this.setMaximumExperienceRequired()
    // this.setMaximumSalaryRequired()
    this.salaryTrendForm.enable()
  }

  onDeleteClick(salaryTrendID: string): void {
    this.openModal(this.deleteConfirmationModal, 'md').result.then(() => {
      this.deleteSalaryTrend(salaryTrendID)
    }, (err) => {
      console.error(err);
      return
    })
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetSalaryTrend() {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getAllSalaryTrend()
      return
    }
    this.salaryTrendSearchForm.patchValue(queryParams)
    this.searchSalaryTrend()
  }

  searchAndCloseDrawer() {
    this.drawer.toggle()
    this.searchSalaryTrend()
  }

  searchSalaryTrend(): void {
    this.searchFormValue = { ...this.salaryTrendSearchForm?.getRawValue() }
    let flag: boolean = true

    for (let field in this.searchFormValue) {
      if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
        delete this.searchFormValue[field]
      } else {
        this.isSearched = true
        this.showSearch = true
        flag = false
      }
    }

    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: this.searchFormValue,
    })

    // No API call on empty search.
    if (flag) {
      return
    }
    this.changePage(1)
  }

  // page change function
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber;
    this.offset = this.currentPage - 1;
    this.getAllSalaryTrend();
  }

  // setMaximumExperienceRequired(): void {


  //   this.salaryTrendForm.get('maximumExperience').updateValueAndValidity()
  // }

  // setMaximumSalaryRequired(): void {


  //   this.salaryTrendForm.get('maximumSalary').updateValueAndValidity()
  // }


  validateFormFields(): void {
    this.salaryTrendForm.get('maximumExperience').
      setValidators([Validators.min(+(this.salaryTrendForm.get('minimumExperience')?.value) + 1)])

    this.salaryTrendForm.get('maximumSalary').
      setValidators([Validators.min(+(this.salaryTrendForm.get('minimumSalary')?.value))])

    this.salaryTrendForm.get('medianSalary').
      setValidators([Validators.min(+(this.salaryTrendForm.get('minimumSalary')?.value)),
      Validators.max(+(this.salaryTrendForm.get('maximumSalary')?.value))])

    this.utilService.updateValueAndValiditors(this.salaryTrendForm)
  }

  onSubmit(): void {
    console.log(this.salaryTrendForm.controls);

    this.validateFormFields()

    if (this.salaryTrendForm.invalid) {
      this.salaryTrendForm.markAllAsTouched();
      return
    }
    if (this.isOperationUpdate) {
      this.updateSalaryTrend()
      return
    }
    this.addSalaryTrend()
  }

  setPaginationString() {
    this.paginationString = ''
    let start: number = this.limit * this.offset + 1
    let end: number = +this.limit + this.limit * this.offset
    if (this.totalSalaryTrend < end) {
      end = this.totalSalaryTrend
    }
    if (this.totalSalaryTrend == 0) {
      this.paginationString = ''
      return
    }
    this.paginationString = `${start}-${end}`
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

  resetSearchAndGetAll(): void {
    this.salaryTrendSearchForm.reset()
    this.searchFormValue = null
    this.changePage(1)
    this.isSearched = false
    this.showSearch = false
    this.router.navigate([this.urlConstant.COMPANY_SALARY_TREND])
  }


  resetSearchForm(): void {
    this.salaryTrendSearchForm.reset()
  }

  // =============================================================== CRUD ===============================================================   

  getAllSalaryTrend(): void {

    if (this.isSearched) {
      this.getSearchSalaryTrend()
      return
    }

    this.spinnerService.loadingMessage = "Getting salary trends"


    this.totalSalaryTrend = 0
    this.salaryTrends = []
    this.isSalaryTrend = true
    this.generalService.getSalaryTrend(this.limit, this.offset).subscribe((response: any) => {
      this.salaryTrends = response.body
      this.totalSalaryTrend = response.headers.get('X-Total-Count')
    }, (error) => {
      this.totalSalaryTrend = 0
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error?.error)
    }).add(() => {
      if (this.totalSalaryTrend == 0) {
        this.isSalaryTrend = false
      }
      this.setPaginationString()
    })
  }

  getSearchSalaryTrend(): void {
    this.spinnerService.loadingMessage = "Searching salary trends"


    this.totalSalaryTrend = 0
    this.salaryTrends = []
    this.isSalaryTrend = true
    this.generalService.getSalaryTrend(this.limit, this.offset, this.searchFormValue).subscribe((response: any) => {
      this.salaryTrends = response.body
      this.totalSalaryTrend = response.headers.get('X-Total-Count')
    }, (error) => {
      this.totalSalaryTrend = 0
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error?.error)
    }).add(() => {


      if (this.totalSalaryTrend == 0) {
        this.isSalaryTrend = false
      }
      this.setPaginationString()
    })
  }

  addSalaryTrend(): void {
    this.spinnerService.loadingMessage = "Adding salary trend"


    this.salaryTrendForm.get('companyRating').setValue(+(this.salaryTrendForm.get('companyRating')?.value))
    this.generalService.addSalaryTrend(this.salaryTrendForm.value).subscribe((response: any) => {
      console.log(response)
      this.modalRef.close()
      alert("Salary trend successfully added")
      this.getAllSalaryTrend()
    }, (error) => {
      console.error(error);
      if (error.error?.error) {
        alert(error.error?.error);
        return;
      }
      alert(error.statusText);
    })
  }

  updateSalaryTrend(): void {
    this.spinnerService.loadingMessage = "Updating salary trend"


    this.salaryTrendForm.get('companyRating').setValue(+(this.salaryTrendForm.get('companyRating')?.value))
    this.generalService.updateSalaryTrend(this.salaryTrendForm.value).subscribe((response: any) => {
      console.log(response)
      this.modalRef.close()
      alert("Salary trend successfully updated")
      this.getAllSalaryTrend()
    }, (error) => {
      console.error(error)
      if (error.error?.error) {
        alert(error.error?.error);
        return
      }
      alert(error.statusText)
    })
  }

  deleteSalaryTrend(salaryTrendID: string): void {
    this.spinnerService.loadingMessage = "Deleting salary trend"


    this.generalService.deleteSalaryTrend(salaryTrendID).subscribe((response: any) => {
      console.log(response)
      this.modalRef.close()
      alert("Salary trend successfully deleted")
      this.getAllSalaryTrend()
    }, (error) => {
      console.error(error);
      if (error.error?.error) {
        alert(error.error?.error);
        return;
      }
      alert(error.statusText);
    })
  }

  // =============================================================== COMPONENTS ===============================================================   

  // Get Technology List.
  getTechnologyList(event?: any): void {
    // this.generalService.getTechnologies().subscribe((respond: any[]) => {
    //   this.technologyList = respond;
    // }, (err) => {
    //   console.error(this.utilService.getErrorString(err))
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

  // Get Designation List.
  getDesignationList(): void {
    this.generalService.getDesignations().subscribe((respond: any[]) => {
      this.desginationList = respond;
    }, (err) => {
      console.error(this.utilService.getErrorString(err))
    })
  }

  getCompanyRatingList(): void {
    this.companyRatingList = [];
    this.generalService.getGeneralTypeByType("company_rating").subscribe((data: any[]) => {
      // this.companyRatingList = data
      for (let rating of data) {
        this.companyRatingList.push(+rating.value)
      }
      this.companyRatingList = this.companyRatingList.sort()
    }, (err) => {
      console.error(err)
    })
  }

}
