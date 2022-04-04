import { formatDate } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, FormControl } from '@angular/forms';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { ICompanyBranch } from 'src/app/service/company/company.service';
import { UrlConstant } from 'src/app/service/constant';
import { ICourse } from 'src/app/service/course/course.service';
import { ITechnologies, GeneralService, ISearchFilterField } from 'src/app/service/general/general.service';
import { INextActionReport, ReportService } from 'src/app/service/report/report.service';
import { INextActionType, ISalesperson } from 'src/app/service/talent/talent.service';

@Component({
  selector: 'app-next-action-report',
  templateUrl: './next-action-report.component.html',
  styleUrls: ['./next-action-report.component.css']
})
export class NextActionReportComponent implements OnInit {

  // modal
  modalRef: any

  // nextActionReports
  nextActionReports: INextActionReport[]
  totalReports: number

  // components
  technologyList: ITechnologies[]
  courseList: ICourse[]
  companyList: ICompanyBranch[]
  userCredentialList: ISalesperson[]
  nextActionList: INextActionType[]

  // form
  nextActionReportForm: FormGroup
  nextActionReportSearchForm: FormGroup

  //pagination
  limit: number;
  currentPage: number;
  offset: number;
  paginationStart: number
  paginationEnd: number

  // search
  searchFormValue: any
  isSearched: boolean
  searchFilterFieldList: ISearchFilterField[]

  disableButton: boolean

  isReportLoaded: boolean
  isCoursesAvailable: boolean
  isCompanyAvailable: boolean
  isTechnologyAvailable: boolean

  // spinner


  constructor(
    private formBuilder: FormBuilder,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private generalService: GeneralService,
    private reportService: ReportService,
    private urlConstant: UrlConstant
  ) {
    this.initializeVariables()
    this.getAllComponents()
    this.createForm()
  }

  initializeVariables(): void {

    this.spinnerService.loadingMessage = "Getting next action report"

    this.limit = 5
    this.offset = 0
    this.totalReports = 0

    this.isSearched = false
    this.isCoursesAvailable = false
    this.isCompanyAvailable = false
    this.isTechnologyAvailable = false
    this.isReportLoaded = true
    this.disableButton = false

    this.courseList = []
    this.companyList = []
    this.userCredentialList = []
    this.nextActionList = []
    this.technologyList = []

    // Search.
    this.searchFilterFieldList = []
  }

  getAllComponents(): void {
    this.getCourseList()
    this.getCompanyList()
    this.getTechnologyList()
    this.getUserCredentialList()
    this.getNextActionTypeList()
    this.getNextActionReports()
  }

  createForm(): void {
    this.createNextActionSearchForm()
    this.createNextActionForm()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  createNextActionSearchForm(): void {
    this.nextActionReportSearchForm = this.formBuilder.group({
      stipend: new FormControl(null),
      fromDate: new FormControl(null),
      toDate: new FormControl(null),
      targetDate: new FormControl(null),
      loginID: new FormControl(null),
      actionType: new FormControl(null),
      courses: new FormControl(null),
      companies: new FormControl(null),
      technologies: new FormControl(null),
    })
  }

  createNextActionForm(): void {
    this.nextActionReportForm = this.formBuilder.group({
      id: new FormControl(null),
      loginID: new FormControl(null),
      loginName: new FormControl(null),
      talent: new FormGroup({
        contact: new FormControl(null),
        email: new FormControl(null),
        firstName: new FormControl(null),
        lastName: new FormControl(null)
      }),
      talentID: new FormControl(null),
      stipend: new FormControl(null),
      referralCount: new FormControl(null),
      fromDate: new FormControl(null),
      toDate: new FormControl(null),
      targetDate: new FormControl(null),
      comment: new FormControl(null),
      actionType: new FormControl(null),
      courses: new FormControl(Array()),
      companies: new FormControl(Array()),
      technologies: new FormControl(Array()),
    })
  }

  onViewReportClick(report: INextActionReport, modalContent: any): void {
    this.createNextActionForm()
    this.isCoursesAvailable = false
    this.isCompanyAvailable = false
    this.isTechnologyAvailable = false

    // console.log(report);

    this.nextActionReportForm.disable()

    if (report.fromDate) {
      report.fromDate = formatDate(report.fromDate, 'mediumDate', 'en_US')
    }
    if (report.toDate) {
      report.toDate = formatDate(report.toDate, 'mediumDate', 'en_US')
    }
    if (report.targetDate) {
      report.targetDate = formatDate(report.targetDate, 'mediumDate', 'en_US')
    }
    if (report.courses) {
      this.isCoursesAvailable = true
    }
    if (report.companies) {
      this.isCompanyAvailable = true
    }
    if (report.technologies) {
      this.isTechnologyAvailable = true
    }

    this.nextActionReportForm.patchValue(report)
    this.openModal(modalContent, "xl")
  }

  // Set total reports list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalReports < this.paginationEnd) {
      this.paginationEnd = this.totalReports
    }
  }

  searchedNextActionReport(): void {
    this.spinnerService.loadingMessage = "Searching report";

    this.searchFormValue = { ...this.nextActionReportSearchForm.value }
    let flag: boolean = true
    for (let field in this.searchFormValue) {
      if (!this.searchFormValue[field]) {
        delete this.searchFormValue[field];
      } else {
        this.isSearched = true
        flag = false
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
    // No API call on empty search.
    if (flag) {

      return
    }
    this.isSearched = true
    this.changePage(1)
  }

  // Delete search criteria from employee search form by search name.
  deleteSearchCriteria(searchName: string): void {
    this.nextActionReportSearchForm.get(searchName).setValue(null)
    this.searchedNextActionReport()
  }

  changePage($event: any): void {

    this.offset = $event - 1
    this.currentPage = $event;
    this.getNextActionReports()
  }

  resetSearchForm(): void {
    this.searchFilterFieldList = []
    this.nextActionReportSearchForm.reset()
  }

  resetSearchAndGetAll(): void {
    this.searchFilterFieldList = []
    this.nextActionReportSearchForm.reset()
    this.isSearched = false
    this.changePage(1)
  }

  openModal(modalContent: any, modalSize?: string): void {
    if (modalSize == undefined) {
      modalSize = 'lg'
    }
    this.modalRef = this.modalService.open(modalContent,
      {
        ariaLabelledBy: 'modal-basic-title',
        backdrop: 'static', size: modalSize,
        keyboard: false
      }
    );
    /*this.modalRef.result.subscribe((result) => {
    }, (reason) => {

    });*/
  }

  // =============================================================CRUD=============================================================

  getNextActionReports(): void {
    if (this.isSearched) {
      this.getSearchNextActionReports()
      return
    }
    this.spinnerService.loadingMessage = "Getting next action reports"

    this.nextActionReports = []
    this.totalReports = 0
    this.isReportLoaded = true
    this.disableButton = true
    this.reportService.getCallingReports(this.limit, this.offset).subscribe((response: any) => {
      this.nextActionReports = response.body
      this.totalReports = response.headers.get("X-Total-Count")
      if (this.totalReports == 0) {
        this.isReportLoaded = false
      }
      this.setPaginationString()
      this.disableButton = false

    }, (error) => {
      this.totalReports = 0
      this.disableButton = false
      this.isReportLoaded = false
      this.setPaginationString()

      console.error(error);
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error.error)
    })
  }

  getSearchNextActionReports(): void {
    this.spinnerService.loadingMessage = "Searching next action reports"

    this.nextActionReports = []
    this.totalReports = 0
    this.disableButton = true
    this.isReportLoaded = true
    this.reportService.getCallingReports(this.limit, this.offset, this.searchFormValue).subscribe((response: any) => {
      this.nextActionReports = response.body
      this.totalReports = response.headers.get("X-Total-Count")
      if (this.totalReports == 0) {
        this.isReportLoaded = false
      }
      this.setPaginationString()
      this.disableButton = false

    }, (error) => {
      this.totalReports = 0
      this.isReportLoaded = false
      this.disableButton = false
      this.setPaginationString()

      console.error(error);
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error.error)
    })
  }

  // =============================================================COMPONENTS=============================================================

  // Get Technology List.
  getTechnologyList(): void {
    this.generalService.getTechnologies().subscribe((respond: any[]) => {
      this.technologyList = respond;
    }, (err) => {
      console.error(err)
    })
  }

  getCourseList(): void {
    this.generalService.getCourseList().subscribe((response: any) => {
      this.courseList = response.body
    }, err => {
      console.error(err)
    })
  }

  getCompanyList(): void {
    this.generalService.getCompanyBranchList().subscribe((response: any) => {
      this.companyList = response.body
    }, (err: any) => {
      console.error(err);
    })
  }

  getUserCredentialList() {
    this.generalService.getUserCredentialList().subscribe(data => {
      this.userCredentialList = data;
    }, err => {
      console.error(err)
    })
  }

  getNextActionTypeList(): void {
    this.generalService.getNextActionTypeList().subscribe((response: any) => {
      this.nextActionList = response.body
    }, (err: any) => {
      console.error(err);
    })
  }

}
