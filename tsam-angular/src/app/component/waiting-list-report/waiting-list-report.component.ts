import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { UrlConstant } from 'src/app/service/constant';
import { IPermission } from 'src/app/service/menu/menu.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { IWaitingListBatchDTO, IWaitingListCompanyBranchDTO, IWaitingListCourseDTO, IWaitingListRequirementDTO, IWaitingListTechnologyDTO, WaitingListReportService } from 'src/app/service/waiting-list-report/waiting-list-report.service';

@Component({
  selector: 'app-waiting-list-report',
  templateUrl: './waiting-list-report.component.html',
  styleUrls: ['./waiting-list-report.component.css']
})
export class WaitingListReportComponent implements OnInit {

  // Waiting list report.
  waitingListCompanyBranchReportList: IWaitingListCompanyBranchDTO[]
  waitingListCourseReportList: IWaitingListCourseDTO[]
  waitingListRequirementReportList: IWaitingListRequirementDTO[]
  waitingListBatchReportList: IWaitingListBatchDTO[]
  waitingListTechnologyReportList: IWaitingListTechnologyDTO[]




  // Permissions.
  permission: IPermission

  // Pagination.
  paginationStart: number
  paginationEnd: number

  // Pagination.
  limit: number
  offset: number
  currentPage: number
  totalEntries: number
  totalSubEntries: number

  // Waiting list realted variables.
  currentCompanyBranchID: string
  currentCourseID: string

  constructor(
    private waitingListReportService: WaitingListReportService,
    private spinnerService: SpinnerService,
    public utilityService: UtilityService,
    private router: Router,
    private urlConstant: UrlConstant,
  ) {
    this.initializeVariables()
    this.getWaitingListComapnyBranchReportList()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize global variables.
  initializeVariables(): void {

    // Waiting list report.
    this.waitingListCompanyBranchReportList = [] as IWaitingListCompanyBranchDTO[]
    this.waitingListCourseReportList = [] as IWaitingListCourseDTO[]
    this.waitingListRequirementReportList = [] as IWaitingListRequirementDTO[]
    this.waitingListBatchReportList = [] as IWaitingListBatchDTO[]
    this.waitingListTechnologyReportList = [] as IWaitingListTechnologyDTO[]

    // Permission.
    this.permission = this.utilityService.getPermission(this.urlConstant.TALENT_REPORT_WAITING_LIST)

    // Spinner.
    this.spinnerService.loadingMessage = "Getting waiting list report"


    // Paginate.
    this.limit = 10
    this.offset = 0
    this.currentPage = 0
    this.totalSubEntries = 0

    // Waiting list realted variables.
    this.currentCompanyBranchID = null
    this.currentCourseID = null
  }

  // Change page for pagination.
  changePage($event): void {

    this.offset = $event - 1
    this.currentPage = $event
    this.getWaitingListComapnyBranchReportList()
  }

  // Get waiting list report by company branch.
  getWaitingListComapnyBranchReportList(): void {
    this.spinnerService.loadingMessage = "Getting waiting list report"


    this.waitingListReportService.getWaitingListCompanyBranchReportList(this.limit, this.offset).subscribe(response => {
      this.waitingListCompanyBranchReportList = response.body
      this.totalEntries = parseInt(response.headers.get("X-Total-Count"))
    }, err => {
      console.error(err)
    }).add(() => {
      this.setPaginationString()


    })
  }

  // Get waiting list report by course.
  getWaitingListCourseReportList(): void {
    this.spinnerService.loadingMessage = "Getting waiting list report"


    this.waitingListReportService.getWaitingListCourseReportList(this.limit, this.offset).subscribe(response => {
      this.waitingListCourseReportList = response.body
      this.totalEntries = parseInt(response.headers.get("X-Total-Count"))

    }, err => {
      console.error(err)

    }).add(() => {
      this.setPaginationString()


    })
  }

  // Get waiting list report by technology.
  getWaitingListTechnologyReportList(): void {
    this.spinnerService.loadingMessage = "Getting waiting list report"


    this.waitingListReportService.getWaitingListTechnologyReportList(this.limit, this.offset).subscribe(response => {
      this.waitingListTechnologyReportList = response.body
      this.totalEntries = parseInt(response.headers.get("X-Total-Count"))

    }, err => {
      console.error(err)

    }).add(() => {
      this.setPaginationString()


    })
  }

  // Set total talents list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalEntries < this.paginationEnd) {
      this.paginationEnd = this.totalEntries
    }
  }

  // on Tab change in list.
  onTabChange(event: any) {
    this.limit = 10
    this.offset = 0
    this.currentPage = 0
    this.totalSubEntries = 0

    if (event == 1) {
      this.getWaitingListComapnyBranchReportList()
    }
    if (event == 2) {
      this.getWaitingListCourseReportList()
    }
    if (event == 3) {
      this.getWaitingListTechnologyReportList()
    }
  }

  // On clicking on any rows of waiting list by company branch list.
  onCompanyBranchRowClick(waitingListCompanyBranchReport: IWaitingListCompanyBranchDTO): void {
    for (let i = 0; i < this.waitingListCompanyBranchReportList.length; i++) {
      if (waitingListCompanyBranchReport.companyBranch.id == this.waitingListCompanyBranchReportList[i].companyBranch.id) {
        continue
      }
      this.waitingListCompanyBranchReportList[i]['isVisible'] = false
    }
    if (waitingListCompanyBranchReport['isVisible']) {
      waitingListCompanyBranchReport['isVisible'] = false
      return
    }
    waitingListCompanyBranchReport['isVisible'] = true
    if (this.currentCompanyBranchID == waitingListCompanyBranchReport.companyBranch.id) {
      this.totalSubEntries = this.waitingListRequirementReportList.length
      return
    }
    this.getRequirementList(waitingListCompanyBranchReport)
  }

  // Get waiting list by requirement report list.
  getRequirementList(waitingListCompanyBranchReport: IWaitingListCompanyBranchDTO): void {
    console.log("get req called")
    this.currentCompanyBranchID = waitingListCompanyBranchReport.companyBranch.id
    this.spinnerService.loadingMessage = "Getting waiting list report"


    this.waitingListReportService.getWaitingListRequirementReportList(waitingListCompanyBranchReport.companyBranch.id).subscribe(response => {
      this.waitingListRequirementReportList = response
      this.totalSubEntries = this.waitingListRequirementReportList.length

    }, err => {
      console.error(err)

    }).add(() => {
      this.setPaginationString()


    })
  }

  // On clicking on any rows of waiting list by course list.
  onCourseRowClick(waitingListCourseReport: IWaitingListCourseDTO): void {
    for (let i = 0; i < this.waitingListCourseReportList.length; i++) {
      if (waitingListCourseReport.course.id == this.waitingListCourseReportList[i].course.id) {
        continue
      }
      this.waitingListCourseReportList[i]['isVisible'] = false
    }
    if (waitingListCourseReport['isVisible']) {
      waitingListCourseReport['isVisible'] = false
      return
    }
    waitingListCourseReport['isVisible'] = true
    if (this.currentCourseID == waitingListCourseReport.course.id) {
      this.totalSubEntries = this.waitingListBatchReportList.length
      return
    }
    this.getBatchList(waitingListCourseReport)
  }

  // Get waiting list by batch report list.
  getBatchList(waitingListCourseReport: IWaitingListCourseDTO): void {
    this.currentCourseID = waitingListCourseReport.course.id
    this.spinnerService.loadingMessage = "Getting waiting list report"


    this.waitingListReportService.getWaitingListBatchReportList(waitingListCourseReport.course.id).subscribe(response => {
      this.waitingListBatchReportList = response
      this.totalSubEntries = this.waitingListBatchReportList.length

    }, err => {
      console.error(err)

    }).add(() => {
      this.setPaginationString()
    })
  }

  // Redirect to talent or enquiries page filtered by waiting list ids.
  redirectByWaitingList(ID: string, type: string, isTalent: boolean): void {
    let queryParams: any = {}
    if (type == "Company branch") {
      queryParams.waitingListCompanyBranchID = ID
    }
    if (type == "Requirement") {
      queryParams.waitingListRequirementID = ID
    }
    if (type == "Course") {
      queryParams.waitingListCourseID = ID
    }
    if (type == "Batch") {
      queryParams.waitingListBatchID = ID
    }
    if (type == "Technology") {
      queryParams.waitingListTechnologyID = ID
    }
    let url: string = ""
    if (isTalent) {
      url = this.urlConstant.TALENT_MASTER
    }
    else {
      url = this.urlConstant.TALENT_ENQUIRY
    }
    this.router.navigate([url], {
      queryParams
    }).catch(err => {
      console.error(err)

    })
  }

}
