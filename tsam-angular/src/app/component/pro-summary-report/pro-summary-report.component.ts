import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { UrlConstant } from 'src/app/service/constant';
import { IPermission } from 'src/app/service/menu/menu.service';
import { ICompanyTechnologyTalent, IProfessionalSummaryReport, IProfessionalSummaryReportCounts, ProSummaryReportService } from 'src/app/service/pro-summary-report/pro-summary-report.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { url } from 'inspector';

@Component({
  selector: 'app-pro-summary-report',
  templateUrl: './pro-summary-report.component.html',
  styleUrls: ['./pro-summary-report.component.css']
})
export class ProSummaryReportComponent implements OnInit {

  // Flags.
  showFirstCategoryColumns: boolean
  showSecondCategoryColumns: boolean
  showThirdCategoryColumns: boolean
  showFourthCategoryColumns: boolean

  // Professional summary report.
  proSummaryReportList: IProfessionalSummaryReport[]
  proSummaryReportCounts: IProfessionalSummaryReportCounts
  companyTechCountsFirst: ICompanyTechnologyTalent[]


  // Permissions.
  permission: IPermission

  // Pagination.
  paginationString: string

  // Pagination.
  limit: number
  offset: number
  currentPage: number
  totalEntries: number

  constructor(
    private proSummaryReportService: ProSummaryReportService,
    private spinnerService: SpinnerService,
    public utilityService: UtilityService,
    private router: Router,
    private urlConstant: UrlConstant,
  ) {
    this.initializeVariables()
    this.getProfessionalSummaryReportList()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize global variables.
  initializeVariables(): void {

    // Flags.
    this.showFirstCategoryColumns = false
    this.showSecondCategoryColumns = false
    this.showThirdCategoryColumns = false
    this.showFourthCategoryColumns = false

    // Pro summary report.
    this.proSummaryReportList = [] as IProfessionalSummaryReport[]
    this.companyTechCountsFirst = [] as ICompanyTechnologyTalent[]
    this.proSummaryReportCounts = { firstCountTotal: 0, secondCountTotal: 0, thirdCountTotal: 0, fourthCountTotal: 0 }

    // Permission.
    this.permission = this.utilityService.getPermission(this.urlConstant.TALENT_REPORT_PRO_SUMMARY)

    // Spinner.
    this.spinnerService.loadingMessage = "Getting professional summary report"


    // Paginate.
    this.limit = 50
    this.offset = 0
    this.currentPage = 0
  }

  // Change page for pagination.
  changePage($event): void {

    this.offset = $event - 1
    this.currentPage = $event
    this.getProfessionalSummaryReportList()
  }

  // Get professional summary report.
  getProfessionalSummaryReportList(): void {
    this.spinnerService.loadingMessage = "Getting Professionals report"


    this.proSummaryReportService.getProfessionalSummaryReportList(this.limit, this.offset).subscribe(response => {
      this.proSummaryReportList = response.body
      this.totalEntries = parseInt(response.headers.get("X-Total-Count"))
      this.proSummaryReportCounts.firstCountTotal = parseInt(response.headers.get("fisrtTotalCount"))
      this.proSummaryReportCounts.secondCountTotal = parseInt(response.headers.get("secondTotalCount"))
      this.proSummaryReportCounts.thirdCountTotal = parseInt(response.headers.get("thirdTotalCount"))
      this.proSummaryReportCounts.fourthCountTotal = parseInt(response.headers.get("fourthTotalCount"))
    }, err => {
      console.error(err)
    }).add(() => {
      this.setPaginationString()
    })
  }

  // Get talent counts by technology and company name.
  getProfessionalSummaryReportByTechCount(category: string): void {
    this.spinnerService.loadingMessage = "Getting Talent Counts By Technology"


    let queryParams: any = {
      "category": category
    }
    this.proSummaryReportService.getProfessionalSummaryReportByTechCount(queryParams).subscribe(response => {
      this.companyTechCountsFirst = response
      this.assignTechCountToCompanyName(category)
      if (category == "first") {
        this.showFirstCategoryColumns = true
        return
      }
      if (category == "second") {
        this.showSecondCategoryColumns = true
        return
      }
      if (category == "third") {
        this.showThirdCategoryColumns = true
        return
      }
      if (category == "fourth") {
        this.showFourthCategoryColumns = true
        return
      }
    }, err => {
      console.error(err)
    }).add(() => {
      this.setPaginationString()
    })
  }

  // Display total entries on current page.
  setPaginationString(): void {
    this.paginationString = ''
    let start: number = this.limit * this.offset + 1
    let end: number = +this.limit + this.limit * this.offset
    if (this.totalEntries < end) {
      end = this.totalEntries
    }
    if (this.totalEntries == 0) {
      this.paginationString = ''
      return
    }
    this.paginationString = `${start}-${end} of ${this.totalEntries}`
  }

  // Redirect to talents page filtered by company name and category.
  redirectToTalentsForProSummaryReport(companyName: string, category: string, isCompany: boolean): void {
    this.router.navigate([this.urlConstant.TALENT_MASTER], {
      queryParams: {
        "companyName": companyName,
        "category": category,
        "isCompany": isCompany
      }
    }).catch(err => {
      console.error(err)

    })
  }

  // Redirect to talents page filtered by company name, category and technology id.
  redirectToTalentsForProSummaryReportByTechnologyCount(companyName: string, category: string, isCompany: boolean, technologyId: string): void {
    this.router.navigate([this.urlConstant.TALENT_MASTER], {
      queryParams: {
        "companyNameTechCount": companyName,
        "categoryTechCount": category,
        "isCompanyTechCount": isCompany,
        "technologyIDTechCount": technologyId
      }
    }).catch(err => {
      console.error(err)

    })
  }

  // Toggle technology name and talent count columns for each expereince colomns.
  toggleTechCount(category: string): void {
    if (category == "first") {
      if (this.showFirstCategoryColumns) {
        this.showFirstCategoryColumns = false
        return
      }
    }
    if (category == "second") {
      if (this.showSecondCategoryColumns) {
        this.showSecondCategoryColumns = false
        return
      }
    }
    if (category == "third") {
      if (this.showThirdCategoryColumns) {
        this.showThirdCategoryColumns = false
        return
      }
    }
    if (category == "fourth") {
      if (this.showFourthCategoryColumns) {
        this.showFourthCategoryColumns = false
        return
      }
    }
    this.getProfessionalSummaryReportByTechCount(category)
    return
  }

  // Assighn the talent counts by technology name for each company name.
  assignTechCountToCompanyName(category: string): void {
    for (let i = 0; i < this.proSummaryReportList.length; i++) {
      this.proSummaryReportList[i][category + 'CatergoryTechCount'] = []
    }
    for (let i = 0; i < this.proSummaryReportList.length; i++) {
      for (let j = 0; j < this.companyTechCountsFirst.length; j++) {
        if (this.companyTechCountsFirst[j].company == this.proSummaryReportList[i].company) {
          this.proSummaryReportList[i][category + 'CatergoryTechCount'].push(this.companyTechCountsFirst[j])
        }
      }
    }
  }

}
