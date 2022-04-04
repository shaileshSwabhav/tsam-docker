import { DatePipe } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormControlName, FormGroup } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BatchService, IBatch } from 'src/app/service/batch/batch.service';
import { CourseService, ICourse } from 'src/app/service/course/course.service';
import { GeneralService, IRole } from 'src/app/service/general/general.service';
import { ICredentialLoginReports, ILoginReports, ReportService } from 'src/app/service/report/report.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { UrlConstant } from 'src/app/service/constant';

@Component({
  selector: 'app-login-report',
  templateUrl: './login-report.component.html',
  styleUrls: ['./login-report.component.css']
})
export class LoginReportComponent implements OnInit {

  // modal
  modalRef: any

  // loginReports
  loginReports: ILoginReports[]
  totalReports: number

  // credentialLoginReports
  credentialLoginReports: ICredentialLoginReports[]
  totalCredentialReports: number
  credentialID: string

  // components
  roleList: IRole[]
  batchList: IBatch[]
  courseList: ICourse[]

  // form
  loginReportForm: FormGroup
  loginReportSearchForm: FormGroup

  // credential-pagination
  credentialLimit: number;
  credentialCurrentPage: number;
  credentialOffset: number;
  credentialPaginationString: string

  // pagination
  limit: number;
  currentPage: number;
  offset: number;
  paginationString: string

  // search
  searchFormValue: any
  showSearch: boolean
  isSearched: boolean

  disableButton: boolean

  isCredentialReportLoaded: boolean
  isReportLoaded: boolean
  loginName: string

  // spinner

  isBatchLoading: boolean

  // index
  weekIndex: number
  isWeeklySearch: boolean

  @ViewChild('loginReportFormModal') loginReportFormModal: any
  @ViewChild('loginReportModal') loginReportModal: any

  constructor(
    private formBuilder: FormBuilder,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private reportService: ReportService,
    private batchService: BatchService,
    private coursService: CourseService,
    private generalService: GeneralService,
    private route: ActivatedRoute,
    private router: Router,
    private utilService: UtilityService,
    private datePipe: DatePipe,
    private urlConstant: UrlConstant,
  ) {
    this.initializeVariables()
    this.getAllComponents()
  }

  initializeVariables(): void {

    this.credentialLoginReports = []
    this.loginReports = []
    this.roleList = []

    this.isSearched = false
    // this.showSearch = false
    this.isCredentialReportLoaded = true
    this.isReportLoaded = true
    this.isWeeklySearch = false
    this.isBatchLoading = false

    this.credentialLimit = 5
    this.credentialOffset = 0
    this.totalCredentialReports = 0

    this.limit = 10
    this.offset = 0
    this.currentPage = 0
    this.totalReports = 0

    this.spinnerService.loadingMessage = "Getting login reports"
    this.credentialPaginationString = ""
    this.weekIndex = 0

    this.createLoginReportSearchForm()
  }

  getAllComponents(): void {
    this.getBatchList()
    this.getCourseList()
    this.getRoleList()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  createLoginReportForm(): void {
    this.loginReportForm = this.formBuilder.group({
      loginName: new FormControl(null),
      loginTime: new FormControl(null),
      logoutTime: new FormControl(null),
      totalHours: new FormControl(null),
      roleName: new FormControl(null),
    })
  }

  createLoginReportSearchForm(): void {
    this.loginReportSearchForm = this.formBuilder.group({
      loginName: new FormControl(null),
      roleID: new FormControl(null),
      role: new FormControl(null),
      fromDate: new FormControl(null),
      toDate: new FormControl(null),
      isActive: new FormControl("1"),
      batchID: new FormControl(null),
      courseID: new FormControl(null),
      // duration: new FormControl(null),
      // firstName: new FormControl(null),
      // lastName: new FormControl(null),
    })
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetReports() {
    let queryParams: any = this.route.snapshot.queryParams

    if (this.utilService.isObjectEmpty(queryParams)) {
      this.loginReportSearchForm.get("role").setValue(this.roleList[0])
      this.getWeeklyReport()
      return
    }

    this.loginReportSearchForm.patchValue(queryParams)

    this.loginReportSearchForm.get("role").setValue(
      this.roleList.find((val) => {
        console.log(val.roleName);
        return val.id == queryParams.roleID
      })
    )

    this.searchLoginReport(false)
  }

  onLoginCountClick(report: ILoginReports): void {
    this.credentialID = report.credentialID
    this.loginName = report.loginName
    this.openModal(this.loginReportModal, "lg")
    this.getCredentialLoginReports()
  }

  changeCredentialPage($event: any): void {

    this.credentialOffset = $event - 1
    this.credentialCurrentPage = $event;
    this.getCredentialLoginReports()
  }

  changePage($event: any): void {
    this.offset = $event - 1
    this.currentPage = $event;
    this.getLoginReports()
  }

  resetSearchForm(): void {
    this.loginReportSearchForm.reset({
      isActive: "1",
      role: this.roleList[0]
    })
  }

  resetSearchAndGetAll(): void {
    this.loginReportSearchForm.reset({
      isActive: "1",
      role: this.roleList[0]
    })
    this.searchFormValue = null
    this.weekIndex = null
    this.isSearched = false
    this.isWeeklySearch = false
    // this.showSearch = false
    this.router.navigate([this.urlConstant.ADMIN_LOGIN_REPORT])
    this.searchLoginReport(false, true)
  }

  openModal(modalContent: any, modalSize?: string): void {
    this.resetCredentialLoginReportVariables()

    if (modalSize == undefined) {
      modalSize = 'lg'
    }
    this.modalRef = this.modalService.open(modalContent, {
      ariaLabelledBy: 'modal-basic-title',
      backdrop: 'static', size: modalSize,
      keyboard: false, centered: true,
    });
  }

  resetCredentialLoginReportVariables(): void {
    this.credentialLoginReports = []
    this.credentialLimit = 5
    this.credentialOffset = 0
    this.credentialCurrentPage = 0
    this.totalCredentialReports = 0
    this.isCredentialReportLoaded = true
  }

  setCredentialPaginationString(): void {
    this.credentialPaginationString = ''
    let start: number = this.credentialLimit * this.credentialOffset + 1
    let end: number = +this.credentialLimit + this.credentialLimit * this.credentialOffset
    if (this.totalCredentialReports < end) {
      end = this.totalCredentialReports
    }
    if (this.totalCredentialReports == 0) {
      this.credentialPaginationString = ''
      return
    }
    this.credentialPaginationString = `${start} - ${end} of ${this.totalCredentialReports}`
  }

  setPaginationString(): void {
    this.paginationString = ''
    let start: number = this.limit * this.offset + 1
    let end: number = +this.limit + this.limit * this.offset
    if (this.totalReports < end) {
      end = this.totalReports
    }
    if (this.totalReports == 0) {
      this.paginationString = ''
      return
    }
    this.paginationString = `${start} - ${end} of ${this.totalReports}`
  }

  searchLoginReport(isWeeklySearch?: boolean, isReset?: boolean): void {
    console.log("searchLoginReport");

    if (!isWeeklySearch) {
      this.isWeeklySearch = false
      this.weekIndex = null
    }

    if (this.loginReportSearchForm.get("role").value) {
      this.loginReportSearchForm.get("roleID").setValue(this.loginReportSearchForm.get("role").value.id)
    }

    if (this.loginReportSearchForm.get("role")?.value?.roleName != "Talent") {
      this.loginReportSearchForm.get("batchID").setValue(null)
      this.loginReportSearchForm.get("courseID").setValue(null)
    }

    this.searchFormValue = { ...this.loginReportSearchForm.value }

    let flag: boolean = true
    for (let field in this.searchFormValue) {
      console.log("field -> ", field);
      if (!this.searchFormValue[field] || field == "role") {
        delete this.searchFormValue[field];
      } else {
        if (!isReset) {
          this.isSearched = true
        }
        flag = false
      }
    }

    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: this.searchFormValue,
    })

    console.log("searchFormValue -> ", this.searchFormValue);
    // No API call on empty search.
    if (flag) {

      return
    }
    this.changePage(1)
  }

  getWeeklyReport(): void {
    if (this.weekIndex === null) {
      // this.isWeeklySearch = false
      this.getLoginReports()
    }
    this.isWeeklySearch = true
    let fromDate: Date = this.getMonday()
    let toDate: Date = new Date(fromDate)
    toDate.setDate(toDate.getDate() + 6)

    this.loginReportSearchForm.get('toDate').setValue(this.datePipe.transform(toDate, 'yyyy-MM-dd'))
    this.loginReportSearchForm.get('fromDate').setValue(this.datePipe.transform(fromDate, 'yyyy-MM-dd'))
    this.loginReportSearchForm.get("role").setValue(this.roleList[0])

    this.searchLoginReport(true)
  }

  /**
  * getMonday
  * @returns date depending on the weekIndex value, 0 is for this week,
  *  1 is for previous week and so on.
  */
  getMonday(): Date {
    let date = new Date();
    let day = date.getDay()
    let prevMonday = new Date()
    let index = this.weekIndex * 7
    // if (day == 1) {
    //   prevMonday.setDate(date.getDate() - (7 + index))
    //   return prevMonday
    // }
    let subDate = ((day + 6) % 7) + index
    prevMonday.setDate(date.getDate() - Math.abs(subDate))
    return prevMonday
  }

  // ============================================================= GET OPERATION =============================================================

  getLoginReports(): void {

    if (this.isSearched) {
      this.getSearchLoginReports()
      return
    }

    this.spinnerService.loadingMessage = "Getting login reports"

    this.loginReports = []
    this.totalReports = 0
    this.isReportLoaded = true

    this.reportService.getLoginReports(this.limit, this.offset, this.searchFormValue).subscribe((response: any) => {
      this.loginReports = response.body
      this.totalReports = response.headers.get("X-Total-Count")
      console.log(response.body)
    }, (error) => {
      this.totalReports = 0
      console.error(error);
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error.error)
    }).add(() => {
      this.setPaginationString()
      this.disableButton = false
      if (this.totalReports == 0) {
        this.isReportLoaded = false
      }

    })
  }

  getSearchLoginReports(): void {
    this.spinnerService.loadingMessage = "Searching login reports"

    this.loginReports = []
    this.totalReports = 0
    this.isReportLoaded = true
    this.reportService.getLoginReports(this.limit, this.offset, this.searchFormValue).subscribe((response: any) => {
      this.loginReports = response.body
      this.totalReports = response.headers.get("X-Total-Count")
      console.log(response.body)
    }, (error) => {
      this.totalReports = 0
      console.error(error);
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error.error)
    }).add(() => {
      this.setPaginationString()
      this.disableButton = false
      if (this.totalReports == 0) {
        this.isReportLoaded = false
      }

    })
  }

  getCredentialLoginReports(): void {
    this.spinnerService.loadingMessage = "Getting credential login reports"

    this.credentialLoginReports = []
    this.totalCredentialReports = 0
    this.isCredentialReportLoaded = true

    this.reportService.getCredentialLoginReports(this.credentialID, this.credentialLimit, this.credentialOffset,
      this.searchFormValue).subscribe((response: any) => {
        this.credentialLoginReports = response.body
        this.totalCredentialReports = response.headers.get("X-Total-Count")
        // console.log(response.body)
      }, (error) => {
        this.totalCredentialReports = 0
        console.error(error);
        if (error.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
          return
        }
        alert(error.error.error)
      }).add(() => {
        this.setCredentialPaginationString()
        this.disableButton = false
        if (this.totalCredentialReports == 0) {
          this.isCredentialReportLoaded = false
        }

      })
  }

  // ============================================================= COMPONENTS =============================================================

  getRoleList(): void {
    this.generalService.getAllRoles().subscribe((response: any) => {
      this.roleList = response
      this.searchOrGetReports()
    }, (err: any) => {
      console.error(err)
    })
  }

  // Get batch list.
  getBatchList(courseID?: string): void {
    let queryParams: any = {
      is_Active: '1'
    }

    if (courseID) {
      queryParams.courseID = courseID
    }
    this.isBatchLoading = true

    this.batchService.getBatchList(queryParams).subscribe(data => {
      this.batchList = data.body
    }, err => {
      console.error(err)
    }).add(() => {
      this.isBatchLoading = false
    })
  }

  // Get course list.
  getCourseList(): void {
    this.generalService.getCourseList().subscribe(response => {
      this.courseList = response.body
    }, err => {
      console.error(err)
    })
  }

  // calculateTotalHours(loginTime: string, logoutTime: string) {
  //   let diff = Math.abs(new Date(logoutTime).getTime() - new Date(loginTime).getTime())
  //   console.log(Math.ceil(diff / (1000 * 3600 * 24)))
  //   return Math.ceil(diff / (1000 * 3600 * 24))
  // }

  // onViewReportClick(report: ICredentialLoginReports): void {
  //   this.createLoginReportForm()
  //   // console.log(report);

  //   this.loginReportForm.disable()

  //   if (report.loginTime) {
  //     report.loginTime = formatDate(report.loginTime, 'medium', 'en_US')
  //   }
  //   if (report.logoutTime) {
  //     report.logoutTime = formatDate(report.logoutTime, 'medium', 'en_US')
  //   }

  //   report.totalHours = report.totalHours.substring(0, 8)

  //   this.loginReportForm.patchValue(report)
  //   this.openModal(this.loginReportFormModal, "lg")
  // }

}
