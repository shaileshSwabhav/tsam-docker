import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { UrlConstant } from 'src/app/service/constant';
import { IExperienceTechnologySummary, IPackageSummary, ITechnologyPackageSummary, ReportService } from 'src/app/service/report/report.service';
import { ITechnology } from 'src/app/service/technology/technology.service';

@Component({
  selector: 'app-package-summary',
  templateUrl: './package-summary.component.html',
  styleUrls: ['./package-summary.component.css']
})
export class PackageSummaryComponent implements OnInit {

  // packageSummary
  packageSummary: IPackageSummary[]
  experienceTechnologyPackageSummary: IExperienceTechnologySummary[]

  searchFormValue: any
  isReportLoaded: boolean
  isTechVisible: boolean

  // component
  technologyList: ITechnology[]



  // flags
  isLessThanThreeVisibile: boolean
  isThreeToFiveVisible: boolean
  isFiveToTenVisible: boolean
  isTenToFifteenVisible: boolean
  isGreaterThanFifteenVisible: boolean
  isTotalPackageVisible: boolean

  // Constants
  readonly LESSTHANTHREE = "LessThanThree"
  readonly THREETOFIVE = "ThreeToFive"
  readonly FIVETOTEN = "FiveToTen"
  readonly TENTOFIFTEEN = "TenToFifteen"
  readonly GREATERTHANFIFTEEN = "GreaterThanFifteen"
  readonly TOTAL = "Total"

  constructor(
    private spinnerService: SpinnerService,
    private reportService: ReportService,
    private urlConstant: UrlConstant,
    private router: Router,
  ) {
    this.initializeVariables()
  }


  initializeVariables(): void {

    this.isReportLoaded = true
    this.isTechVisible = false
    this.isLessThanThreeVisibile = false
    this.isThreeToFiveVisible = false
    this.isFiveToTenVisible = false
    this.isTenToFifteenVisible = false
    this.isTotalPackageVisible = false
    this.isGreaterThanFifteenVisible = false

    this.packageSummary = []
    this.experienceTechnologyPackageSummary = []
    this.technologyList = []

    this.searchFormValue = null

    this.spinnerService.loadingMessage = "Getting package summary"

    this.getAllPackageSummary()
    this.getSummaryTechnologyList()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  calculateTotalExperience(packageSummary: IPackageSummary): number {
    return packageSummary.lessThanThree + packageSummary.threeToFive
      + packageSummary.fiveToTen + packageSummary.tenToFifteen + packageSummary.greaterThanFifteen
  }

  calculatePackageTotal(packageType: string): number {
    let count: number = 0
    for (let index = 0; index < this.packageSummary.length; index++) {
      if (this.LESSTHANTHREE == packageType) {
        count += this.packageSummary[index].lessThanThree
      }
      if (this.THREETOFIVE == packageType) {
        count += this.packageSummary[index].threeToFive
      }
      if (this.FIVETOTEN == packageType) {
        count += this.packageSummary[index].fiveToTen
      }
      if (this.TENTOFIFTEEN == packageType) {
        count += this.packageSummary[index].tenToFifteen
      }
      if (this.GREATERTHANFIFTEEN == packageType) {
        count += this.packageSummary[index].greaterThanFifteen
      }
    }
    return count
  }

  calculatePackageTechTotal(technology: ITechnology): number {
    let count: number = 0
    for (let index = 0; index < this.experienceTechnologyPackageSummary.length; index++) {
      count += this.experienceTechnologyPackageSummary[index].technologySummary.find((val) => {
        return technology.language === val.techLanguage
      }).totalCount
    }
    return count
  }

  // calculateOtherTechTotal(): number {
  //   let count: number = 0
  //   for (let index = 0; index < this.experienceTechnologyPackageSummary.length; index++) {
  //     count += this.experienceTechnologyPackageSummary[index].technologySummary.find((val) => {
  //         return "Other" === val.techLanguage
  //       }).totalCount
  //   }
  //   return count
  // }

  calculateTotal(): number {
    let count: number = 0

    count += this.calculatePackageTotal(this.LESSTHANTHREE)
    count += this.calculatePackageTotal(this.THREETOFIVE)
    count += this.calculatePackageTotal(this.FIVETOTEN)
    count += this.calculatePackageTotal(this.TENTOFIFTEEN)
    count += this.calculatePackageTotal(this.GREATERTHANFIFTEEN)

    return count
  }

  onLessThanThreeClick(): void {
    this.isLessThanThreeVisibile = !this.isLessThanThreeVisibile
    this.isThreeToFiveVisible = false
    this.isFiveToTenVisible = false
    this.isTenToFifteenVisible = false
    this.isGreaterThanFifteenVisible = false
    this.isTotalPackageVisible = false
    this.isTechVisible = false
    if (this.isLessThanThreeVisibile) {
      this.getTechnologyPackageSummary(this.LESSTHANTHREE)
    }
  }

  onThreeToFiveClick(): void {
    this.isThreeToFiveVisible = !this.isThreeToFiveVisible
    this.isLessThanThreeVisibile = false
    this.isFiveToTenVisible = false
    this.isTenToFifteenVisible = false
    this.isGreaterThanFifteenVisible = false
    this.isTotalPackageVisible = false
    this.isTechVisible = false
    if (this.isThreeToFiveVisible) {
      this.getTechnologyPackageSummary(this.THREETOFIVE)
    }
  }

  onFiveToTenClick(): void {
    this.isFiveToTenVisible = !this.isFiveToTenVisible
    this.isLessThanThreeVisibile = false
    this.isThreeToFiveVisible = false
    this.isTenToFifteenVisible = false
    this.isGreaterThanFifteenVisible = false
    this.isTotalPackageVisible = false
    this.isTechVisible = false
    if (this.isFiveToTenVisible) {
      this.getTechnologyPackageSummary(this.FIVETOTEN)
    }
  }

  onTenToFifteenClick(): void {
    this.isTenToFifteenVisible = !this.isTenToFifteenVisible
    this.isLessThanThreeVisibile = false
    this.isThreeToFiveVisible = false
    this.isFiveToTenVisible = false
    this.isTotalPackageVisible = false
    this.isTechVisible = false
    this.isGreaterThanFifteenVisible = false
    if (this.isTenToFifteenVisible) {
      this.getTechnologyPackageSummary(this.TENTOFIFTEEN)
    }
  }

  onGreaterThanFifteenClick(): void {
    this.isGreaterThanFifteenVisible = !this.isGreaterThanFifteenVisible
    this.isLessThanThreeVisibile = false
    this.isThreeToFiveVisible = false
    this.isFiveToTenVisible = false
    this.isTotalPackageVisible = false
    this.isTenToFifteenVisible = false
    this.isTechVisible = false
    if (this.isGreaterThanFifteenVisible) {
      this.getTechnologyPackageSummary(this.GREATERTHANFIFTEEN)
    }
  }

  onTotalPackageClick(): void {
    this.isTotalPackageVisible = !this.isTotalPackageVisible
    this.isLessThanThreeVisibile = false
    this.isThreeToFiveVisible = false
    this.isFiveToTenVisible = false
    this.isTenToFifteenVisible = false
    this.isTechVisible = false
    if (this.isTotalPackageVisible) {
      this.getTechnologyPackageSummary()
    }
  }

  getAllPackageSummary(): void {
    this.spinnerService.loadingMessage = "Getting package summary"

    this.packageSummary = []
    this.reportService.getAllPackageSummary().subscribe((response: any) => {
      this.packageSummary = response.body
      // console.log(response.body);
    }, (err: any) => {
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  getTechnologyPackageSummary(packageType?: string): void {
    this.spinnerService.loadingMessage = "Getting summary"

    this.experienceTechnologyPackageSummary = []

    let queryParams: any = {}
    if (packageType) {
      queryParams = {
        packageType: packageType
      }
    }
    this.reportService.getTechnologyPackageSummary(queryParams).subscribe((response: any) => {
      this.experienceTechnologyPackageSummary = response.body
      this.isTechVisible = true
      // console.log(response.body);
    }, (err: any) => {
      console.error(err);
      this.isTechVisible = false
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  //  ============================================ REDIRECTS ============================================

  redirectToTalents(summary: any, packageType?: string, tech?: ITechnologyPackageSummary): void {
    // console.log(summary);
    // console.log(packageType);
    // console.log(tech);
    let packageExperience: string
    let packageTechnology: string = null

    if (tech && tech?.technology) {
      packageTechnology = tech?.technology.id
    }

    if (tech && tech?.technology == null) {
      packageTechnology = tech.techLanguage
    }

    switch (summary.experience) {
      case "0 to 3 Years":
        packageExperience = "0-3"
        break
      case "3 to 6 Years":
        packageExperience = "3-6"
        break
      case "6 to 8 Years":
        packageExperience = "6-8"
        break
      case "8 to 10 Years":
        packageExperience = "8-10"
        break
      case "10+ Years":
        packageExperience = "10+"
        break
    }

    // console.log(packageExperience);

    this.router.navigate([this.urlConstant.TALENT_MASTER], {
      queryParams: {
        "isPackageSummary": true,
        "packageExperience": packageExperience,
        "packageType": packageType,
        "packageTechnology": packageTechnology,
      }
    }).catch(err => {
      console.error(err)

    })
  }

  redirectTotalToTalents(packageType?: string, tech?: string): void {
    // console.log(packageType);
    // console.log(tech);

    this.router.navigate([this.urlConstant.TALENT_MASTER], {
      queryParams: {
        "isPackageSummary": true,
        "packageType": packageType,
        "packageTechnology": tech,
      }
    }).catch(err => {
      console.error(err)

    })
  }



  //  ============================================ COMPONENTS ============================================

  getSummaryTechnologyList(): void {
    this.reportService.getSummaryTechnologyList().subscribe((response: any) => {
      this.technologyList = response.body
    }, (err: any) => {
      console.error(err);
    })
  }
}
