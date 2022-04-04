import { Component, OnInit, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { IFresherSummary, IAcademicTechnologySummary, ReportService, ITechnologyFresherSummary } from 'src/app/service/report/report.service';
import { ITechnology } from 'src/app/service/technology/technology.service';
import { UrlConstant } from 'src/app/service/constant';

@Component({
  selector: 'app-fresher-summary',
  templateUrl: './fresher-summary.component.html',
  styleUrls: ['./fresher-summary.component.css']
})
export class FresherSummaryComponent implements OnInit {

  // fresherSummary
  fresherSummary: IFresherSummary[]
  academicTechnologyFresherSummary: IAcademicTechnologySummary[]
  technologyFresherSummary: ITechnologyFresherSummary[]

  searchFormValue: any
  isReportLoaded: boolean

  // spinner


  // flags
  isOutstandingTechVisible: boolean
  isExcellentTechVisible: boolean
  isAverageTechVisible: boolean
  isUnrankedTechVisible: boolean
  isTotalVisible: boolean
  isTechVisible: boolean

  // Constants
  readonly OUTSTANDING = "Outstanding"
  readonly EXCELLENT = "Excellent"
  readonly AVERAGE = "Average"
  readonly UNRANKED = "Unranked"
  readonly PROFESSIONAL = "Professional"
  readonly JOBSEEKER = "Looking to switch"
  readonly REQUIREMENT = "Current Requirements"
  readonly TOTAL = "Total"

  // component
  technologyList: ITechnology[]

  @ViewChild('fresherSummaryFormModal') fresherSummaryFormModal: any

  constructor(
    private spinnerService: SpinnerService,
    private reportService: ReportService,
    private router: Router,
    private urlConstant: UrlConstant,
  ) {

    this.initializeVariables()
  }

  initializeVariables(): void {

    this.isReportLoaded = true
    this.isOutstandingTechVisible = false
    this.isExcellentTechVisible = false
    this.isAverageTechVisible = false
    this.isUnrankedTechVisible = false

    this.fresherSummary = []
    this.technologyFresherSummary = []
    this.technologyList = []

    this.searchFormValue = null

    this.spinnerService.loadingMessage = "Getting fresher summary"

    this.getAllFresherSummary()
    this.getSummaryTechnologyList()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }


  onOutstandingClick(): void {
    this.isOutstandingTechVisible = !this.isOutstandingTechVisible
    this.isAverageTechVisible = false
    this.isExcellentTechVisible = false
    this.isUnrankedTechVisible = false
    this.isTechVisible = false
    if (this.isOutstandingTechVisible) {
      this.getTechnologyFresherSummary(this.OUTSTANDING)
    }
  }

  onExcellentClick(): void {
    this.isExcellentTechVisible = !this.isExcellentTechVisible
    this.isAverageTechVisible = false
    this.isOutstandingTechVisible = false
    this.isUnrankedTechVisible = false
    this.isTotalVisible = false
    this.isTechVisible = false
    if (this.isExcellentTechVisible) {
      this.getTechnologyFresherSummary(this.EXCELLENT)
    }
  }

  onAverageClick(): void {
    this.isAverageTechVisible = !this.isAverageTechVisible
    this.isExcellentTechVisible = false
    this.isOutstandingTechVisible = false
    this.isUnrankedTechVisible = false
    this.isTotalVisible = false
    this.isTechVisible = false
    if (this.isAverageTechVisible) {
      this.getTechnologyFresherSummary(this.AVERAGE)
    }
  }

  onUnrankedClick(): void {
    this.isUnrankedTechVisible = !this.isUnrankedTechVisible
    this.isExcellentTechVisible = false
    this.isOutstandingTechVisible = false
    this.isAverageTechVisible = false
    this.isTotalVisible = false
    this.isTechVisible = false
    if (this.isUnrankedTechVisible) {
      this.getTechnologyFresherSummary(this.UNRANKED)
    }
  }

  onTotalClick(): void {
    this.isTotalVisible = !this.isTotalVisible
    this.isUnrankedTechVisible = false
    this.isExcellentTechVisible = false
    this.isOutstandingTechVisible = false
    this.isAverageTechVisible = false
    this.isTechVisible = false
    if (this.isTotalVisible) {
      this.getTechnologyFresherSummary()
    }
  }

  calculateAcademicYearTotal(summary: IFresherSummary): number {
    return summary.outstandingCount + summary.excellentCount + summary.averageCount + summary.unrankedCount
  }

  calculateFresherCount(talentType: string): number {
    let count: number = 0

    const conditionArray = [
      this.PROFESSIONAL, this.JOBSEEKER, this.REQUIREMENT
    ]

    for (let index = 0; index < this.fresherSummary.length; index++) {
      if (!conditionArray.includes(this.fresherSummary[index].columnName)) {
        if (talentType == this.OUTSTANDING) {
          count += this.fresherSummary[index].outstandingCount
          continue
        }
        if (talentType == this.EXCELLENT) {
          count += this.fresherSummary[index].excellentCount
          continue
        }
        if (talentType == this.AVERAGE) {
          count += this.fresherSummary[index].averageCount
          continue
        }
        if (talentType == this.UNRANKED) {
          count += this.fresherSummary[index].unrankedCount
          continue
        }
      }
    }
    return count
  }

  calculateTotalTalentsCount(talentType: string): number {
    let count: number = 0

    const conditionArray = [
      this.JOBSEEKER, this.REQUIREMENT
    ]

    for (let index = 0; index < this.fresherSummary.length; index++) {
      if (!conditionArray.includes(this.fresherSummary[index].columnName)) {
        if (talentType == this.OUTSTANDING) {
          count += this.fresherSummary[index].outstandingCount
          continue
        }
        if (talentType == this.EXCELLENT) {
          count += this.fresherSummary[index].excellentCount
          continue
        }
        if (talentType == this.AVERAGE) {
          count += this.fresherSummary[index].averageCount
          continue
        }
        if (talentType == this.UNRANKED) {
          count += this.fresherSummary[index].unrankedCount
          continue
        }
      }
    }
    return count
  }

  calculateTalent(): number {
    let count: number = 0

    const conditionArray = [
      this.JOBSEEKER, this.REQUIREMENT
    ]

    for (let index = 0; index < this.fresherSummary.length; index++) {
      if (!conditionArray.includes(this.fresherSummary[index].columnName)) {
        count += this.calculateAcademicYearTotal(this.fresherSummary[index])
      }
    }
    return count
  }

  calculateTotalFresherCount(): number {
    let count: number = 0

    count += this.calculateFresherCount(this.OUTSTANDING)
    count += this.calculateFresherCount(this.EXCELLENT)
    count += this.calculateFresherCount(this.AVERAGE)
    count += this.calculateFresherCount(this.UNRANKED)

    return count
  }

  calculateFresherTechnologyTotal(technology: ITechnology): number {
    let count: number = 0
    for (let index = 0; index < this.academicTechnologyFresherSummary.length; index++) {
      if (this.academicTechnologyFresherSummary[index].academic != null &&
        this.academicTechnologyFresherSummary[index].columnName != this.PROFESSIONAL) {

        count += this.academicTechnologyFresherSummary[index].technologySummary.find((val) => {
          return technology.language === val.technologyLangugage
        }).totalCount

      }
    }

    return count
  }

  calculateFresherOtherTechnologyTotal(): number {
    let count: number = 0
    for (let index = 0; index < this.academicTechnologyFresherSummary.length; index++) {
      if (this.academicTechnologyFresherSummary[index].academic != null &&
        this.academicTechnologyFresherSummary[index].columnName != this.PROFESSIONAL) {

        count += this.academicTechnologyFresherSummary[index].technologySummary.find((val) => {
          return "Other" === val.technologyLangugage
        }).totalCount

      }
    }

    return count
  }

  calculateTechnologyTotal(technology: ITechnology): number {
    let count: number = 0

    for (let index = 0; index < this.academicTechnologyFresherSummary.length; index++) {
      if (this.academicTechnologyFresherSummary[index].academic != null) {

        count += this.academicTechnologyFresherSummary[index].technologySummary.find((val) => {
          return technology.language === val.technologyLangugage
        }).totalCount
      }
    }

    return count
  }

  calculateOtherTechnologyTotal(): number {
    let count: number = 0
    for (let index = 0; index < this.academicTechnologyFresherSummary.length; index++) {
      if (this.academicTechnologyFresherSummary[index].academic != null) {

        count += this.academicTechnologyFresherSummary[index].technologySummary.find((val) => {
          return "Other" === val.technologyLangugage
        }).totalCount
      }
    }

    return count
  }

  getAllFresherSummary(): void {
    this.spinnerService.loadingMessage = "Getting fresher summary"

    this.fresherSummary = []
    this.reportService.getAllFresherSummary().subscribe((response: any) => {
      this.fresherSummary = response.body
    }, (err: any) => {
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  getTechnologyFresherSummary(talentType?: string): void {
    this.spinnerService.loadingMessage = "Getting summary"

    this.technologyFresherSummary = []
    this.academicTechnologyFresherSummary = []

    let queryParams: any = {}
    if (talentType) {
      queryParams = {
        talentType: talentType
      }
    }
    this.reportService.getTechnologyFresherSummary(queryParams).subscribe((response: any) => {
      this.academicTechnologyFresherSummary = response.body
      this.isTechVisible = true
      // this.technologyFresherSummary = response.body
      // console.log(response.body);
    }, (err: any) => {
      this.isTechVisible = false
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  // Redirect to talents page filtered by talent-type, columnName and technology.
  redirectToTalents(summary: any, isLookingForJob: string, talentType?: string, tech?: ITechnologyFresherSummary): void {
    // console.log(summary)
    // console.log(talentType)
    // console.log(tech)
    let isExperienced: string = '0'
    let fresherTechnology: string

    if (summary.columnName == this.PROFESSIONAL) {
      isExperienced = '1'
    }

    if (tech != null && tech.technology != null) {
      fresherTechnology = tech.technology.id
    }

    if (tech != null && tech.technologyLangugage == "Other") {
      fresherTechnology = "Other"
    }

    if (isLookingForJob == "1") {
      isExperienced = "1"
    }

    if (summary.columnName == this.REQUIREMENT) {
      this.router.navigate([this.urlConstant.COMPANY_REQUIREMENT], {
        queryParams: {
          "isFresherSummary": true,
          "talentRating": talentType,
          "technologies": fresherTechnology
        }
      })
      return
    }


    this.router.navigate([this.urlConstant.TALENT_MASTER], {
      queryParams: {
        "isFresherSummary": true,
        "talentType": talentType,
        "academicYear": summary.academic.key,
        "fresherTechnology": fresherTechnology,
        "isExperienced": isExperienced,
        "isLookingForJob": isLookingForJob
      }
    }).catch(err => {
      console.error(err)

    });
  }

  redirectTotalToTalents(talentType: string, isExperienced?: string, technology?: string): void {

    console.log(technology);

    this.router.navigate([this.urlConstant.TALENT_MASTER], {
      queryParams: {
        "isFresherSummary": true,
        "talentType": talentType,
        "fresherTechnology": technology,
        "isExperienced": isExperienced
      }
    }).catch(err => {
      console.error(err)

    });
  }

  redirect(isExperienced?: string): void {

    this.router.navigate([this.urlConstant.TALENT_MASTER], {
      queryParams: {
        "isFresherSummary": true,
        "isExperienced": isExperienced
      }
    }).catch(err => {
      console.error(err)

    });
  }


  getSummaryTechnologyList(): void {
    this.reportService.getSummaryTechnologyList().subscribe((response: any) => {
      this.technologyList = response.body
    }, (err: any) => {
      console.error(err);
    })
  }
}
