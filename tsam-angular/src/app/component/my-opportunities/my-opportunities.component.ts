import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { CompanyRequirementService } from 'src/app/service/company/company-requirement/company-requirement.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { TalentService } from 'src/app/service/talent/talent.service';

@Component({
  selector: 'app-my-opportunities',
  templateUrl: './my-opportunities.component.html',
  styleUrls: ['./my-opportunities.component.css']
})
export class MyOpportunitiesComponent implements OnInit {

  // Company requirement.
  myOppurtunityList: any[]

  // Pagination.
  limit: number
  currentPage: number
  offset: number
  paginationString: string
  totalMyOppurtunities: number

  // Spinner.



  // Flags.
  allJobsSelected: boolean
  appliedSelected: boolean
  eligibleSelected: boolean

  // Eligible talent.
  eligibleTalent: any
  eligibleDegreeIDs: any
  eligibleDesignationIDs: any
  eligibleTechnologyIDs: any
  eligibleExpTechnologyIDs: any
  eligibleModel: any

  constructor(
    private spinnerService: SpinnerService,
    private companyRequirementService: CompanyRequirementService,
    private talentService: TalentService,
    private localService: LocalService,
    private router: Router,
    private route: ActivatedRoute,
  ) {
    this.initializeVariables()
    this.getEligibleTalent()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize global variables.
  initializeVariables(): void {

    // Components.
    this.myOppurtunityList = []

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Opportunities"


    // Pagination.
    this.limit = 8
    this.offset = 0
    this.currentPage = 0

    // Flags
    this.allJobsSelected = false
    this.appliedSelected = false
    this.eligibleSelected = true

    // Eligible talent.
    this.eligibleDegreeIDs = []
    this.eligibleDesignationIDs = []
    this.eligibleTechnologyIDs = []
    this.eligibleExpTechnologyIDs = []
    this.eligibleModel = {}
  }

  // Get my oppurtunities.
  getMyOpportunities(): void {
    this.spinnerService.loadingMessage = "Getting Oppurtunities"


    let queryParams: any = {}
    if (this.allJobsSelected) {
      queryParams.isActive = "1"


      this.routeTo("all-jobs")
    }
    if (this.appliedSelected) {
      queryParams.isActive = "1"
      queryParams.talentID = this.localService.getJsonValue("loginID")


      this.routeTo("applied")
    }
    if (this.eligibleSelected) {
      queryParams = this.eligibleModel
      queryParams.isActive = "1"


      this.routeTo("eligible")
    }
    queryParams.limit = this.limit
    queryParams.offset = this.offset

    console.log(queryParams)
    this.companyRequirementService.getMyOppurtunities(queryParams).subscribe((response) => {
      this.totalMyOppurtunities = response.headers.get('X-Total-Count')
      this.myOppurtunityList = response.body
      this.formatCompanyRequirementsFields()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.setPaginationString()
    })
  }

  routeTo(type: string) {
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: {
        type: type
      },
    });
  }

  // Get eligible talent.
  getEligibleTalent(): void {
    this.spinnerService.loadingMessage = "Getting Oppurtunities"


    this.talentService.getEligibleTalent(this.localService.getJsonValue("loginID")).subscribe((response) => {
      this.eligibleTalent = response
      this.getEligibleCriteriaFromTalent()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.setPaginationString()
    })
  }

  // Get the eligibility criteria from talent.
  getEligibleCriteriaFromTalent(): void {
    // Get technology IDs.
    if (this.eligibleTalent.technologies && this.eligibleTalent.technologies?.length > 0) {
      for (let i = 0; i < this.eligibleTalent.technologies.length; i++) {
        this.eligibleTechnologyIDs.push(this.eligibleTalent.technologies[i].id)
      }
    }

    // Get degree IDs.
    if (this.eligibleTalent.academics && this.eligibleTalent.academics?.length > 0) {
      for (let i = 0; i < this.eligibleTalent.academics.length; i++) {
        this.eligibleDegreeIDs.push(this.eligibleTalent.academics[i].degree.id)
      }
    }

    // Get designation and experience technology IDs.
    if (this.eligibleTalent.experiences && this.eligibleTalent.experiences.length != 0) {
      for (let i = 0; i < this.eligibleTalent.experiences.length; i++) {
        this.eligibleDesignationIDs.push(this.eligibleTalent.experiences[i].designation.id)
        if (this.eligibleTalent.experiences[i]?.technologies && this.eligibleTalent.experiences[i].technologies?.length > 0) {
          for (let j = 0; j < this.eligibleTalent.experiences[i].technologies.length; j++) {
            if (this.eligibleExpTechnologyIDs.includes(this.eligibleTalent.experiences[i].technologies[j].id)) {
              continue
            }
            this.eligibleExpTechnologyIDs.push(this.eligibleTalent.experiences[i].technologies[j].id)
          }
        }
      }
    }

    // Make the eligibility critera params.
    if (!this.eligibleTalent.isExperience && this.eligibleTechnologyIDs?.length > 0) {
      this.eligibleModel.technologies = this.eligibleTechnologyIDs
      this.eligibleModel.isFresher = "true"
    }
    if (this.eligibleTalent.isExperience && this.eligibleExpTechnologyIDs?.length > 0) {
      this.eligibleModel.technologies = this.eligibleExpTechnologyIDs
      this.eligibleModel.isFresher = "false"
    }
    if (this.eligibleTalent.isExperience && this.eligibleDesignationIDs?.length > 0) {
      this.eligibleModel.designation = this.eligibleDesignationIDs
      this.eligibleModel.isFresher = "false"
    }
    if (this.eligibleDegreeIDs?.length > 0) {
      this.eligibleModel.qualifications = this.eligibleDegreeIDs
    }

    // Call get my oppurtunities based on the criteria.
    this.getMyOpportunities()
  }

  // On clicking all jobs button.
  onAllJobsButtonClick(): void {
    this.allJobsSelected = true
    this.appliedSelected = false
    this.eligibleSelected = false
    this.getMyOpportunities()
  }

  // On clicking applied button.
  onAppliedButtonClick(): void {
    this.allJobsSelected = false
    this.appliedSelected = true
    this.eligibleSelected = false
    this.getMyOpportunities()
  }

  // On clicking eligible button.
  onEligibleButtonClick(): void {
    this.allJobsSelected = false
    this.appliedSelected = false
    this.eligibleSelected = true
    // this.getMyOpportunities()
    this.getEligibleTalent()
  }

  // Format fields of company requirement through interface.
  formatCompanyRequirementsFields(): void {
    for (let requirement of this.myOppurtunityList) {
      // Format package offered in indian rupee system.
      requirement.packageOfferedInstring = this.formatPackageInLPA(requirement['minimumPackage'], requirement['maximumPackage'])
    }
  }

  // Format the package number in terms of LPA.
  formatPackageInLPA(min: number, max: number): string {
    if (!min && !max) {
      return
    }
    let minNumber: number = min / 100000
    let maxNumber: number = max / 100000
    let output: string = minNumber + " - " + maxNumber + " Lpa"
    return output
  }

  // Set total company requirements list on current page.
  setPaginationString() {
    this.paginationString = ''
    let start: number = this.limit * this.offset + 1
    let end: number = +this.limit + this.limit * this.offset
    if (this.totalMyOppurtunities < end) {
      end = this.totalMyOppurtunities
    }
    if (this.totalMyOppurtunities == 0) {
      this.paginationString = ''
      return
    }
    this.paginationString = `${start} - ${end} of ${this.totalMyOppurtunities}`
  }

  // On page change.
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
    this.getMyOpportunities()
  }

  // Redirect to company details page.
  redirectToCompanyDetails(companyRequirementID: string): void {
    this.router.navigate(['/my-opportunities/company-details'], {
      queryParams: {
        "companyRequirementID": companyRequirementID,
        "email": this.localService.getJsonValue("email"),
        "talentID": this.localService.getJsonValue("loginID"),
      }
    }).catch(err => {
      console.error(err)
    })
  }

}
