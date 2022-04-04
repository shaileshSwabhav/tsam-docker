import { Component, OnInit } from '@angular/core';
import { ISearchFilterField, TalentService } from 'src/app/service/talent/talent.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { FormBuilder, FormControl, FormGroup } from '@angular/forms';
import { IPermission } from 'src/app/service/menu/menu.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UrlConstant } from 'src/app/service/constant';
import { Router, ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-lifetime-value-report',
  templateUrl: './lifetime-value-report.component.html',
  styleUrls: ['./lifetime-value-report.component.css']
})
export class LifetimeValueReportComponent implements OnInit {

  // List.
  lifetimeValueReportList: any[]

  // Pagination.
  limit: number
  offset: number
  currentPage: number
  totalLifetimeValueEntries: number
  paginationStart: number
  paginationEnd: number

  // Total.
  totalLifetimeValue: number
  totalLifetimeValueInString: string



  // Search.
  isSearched: boolean
  lifetimeValueSearchForm: FormGroup
  searchFormValue: any
  searchFilterFieldList: ISearchFilterField[]

  // Permissions.
  permission: IPermission

  constructor(
    private talentService: TalentService,
    private spinnerService: SpinnerService,
    private formBuilder: FormBuilder,
    public utilityService: UtilityService,
    private urlConstant: UrlConstant,
    private router: Router,
    private route: ActivatedRoute,
  ) {
    this.initializeVariables()
    this.searchOrGetLifetimeValueReports()
  }

  // Initialize global variables.
  initializeVariables(): void {
    // Paginate.
    this.limit = 5
    this.offset = 0
    this.currentPage = 0

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Lifetime Value Report"

    // Total
    this.totalLifetimeValue = 0

    // Search.
    this.searchFormValue = {}
    this.isSearched = false
    this.searchFilterFieldList = []

    // Get permissions from menus using utilityService function.
    this.permission = this.utilityService.getPermission(this.urlConstant.TALENT_REPORT_LIFETIME)

    // Create forms.
    this.createLifetimeValueSearchForm()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Change page for pagination.
  changePage($event): void {

    this.offset = $event - 1
    this.currentPage = $event
    this.getAllLifetimeValueReports()
  }

  // Create new campus drive search form.
  createLifetimeValueSearchForm(): void {
    this.lifetimeValueSearchForm = this.formBuilder.group({
      firstName: new FormControl(null),
      lastName: new FormControl(null),
      email: new FormControl(null),
      upsellMinimum: new FormControl(null),
      upsellMaximum: new FormControl(null),
      placementMinimum: new FormControl(null),
      placementMaximum: new FormControl(null),
      knowledgeMinimum: new FormControl(null),
      knowledgeMaximum: new FormControl(null),
      teachingMinimum: new FormControl(null),
      teachingMaximum: new FormControl(null),
      totalMinimum: new FormControl(null),
      totalMaximum: new FormControl(null),
    })
  }

  // Get all lifetime value reports.
  getAllLifetimeValueReports(): void {
    this.spinnerService.loadingMessage = "Getting Lifetime Value Report"
    this.talentService.getLifetimeValueReports(this.limit, this.offset, this.searchFormValue).subscribe(response => {
      this.lifetimeValueReportList = response.body
      this.totalLifetimeValueEntries = parseInt(response.headers.get("X-Total-Count"))
      this.totalLifetimeValue = parseInt(response.headers.get("totalLifetimeValue"))
      this.totalLifetimeValueInString = this.formatLifetimeValueInIndianRupeeSystem(this.totalLifetimeValue)

    }, err => {
      console.error(err)

    }).add(() => {
      this.setPaginationString()
    })
  }

  // Reset search form and get all lifetime valuee entries.
  resetSearchAndGetAll() {
    this.searchFilterFieldList = []
    this.isSearched = false
    this.spinnerService.loadingMessage = "Getting Lifetime Value Report"

    this.lifetimeValueSearchForm.reset()
    this.searchFormValue = null
    this.createLifetimeValueSearchForm()
    this.changePage(1)
  }

  // Reset campus lifetime value search form.
  resetSearchForm(): void {
    this.searchFilterFieldList = []
    this.lifetimeValueSearchForm.reset()
  }

  //On clicking search lifetime values form button.
  searchLifetimeValues(): void {
    this.searchFormValue = { ...this.lifetimeValueSearchForm?.value }
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: this.searchFormValue,
    })
    for (let field in this.searchFormValue) {
      if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
        delete this.searchFormValue[field]
      } else {
        this.isSearched = true
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
    if (!this.isSearched) {
      return
    }
    this.spinnerService.loadingMessage = "Searching Lifetime Value Reports"
    this.changePage(1)
  }

  // Delete search criteria from lifetime value search form by search name.
  deleteSearchCriteria(searchName: string): void {
    this.lifetimeValueSearchForm.get(searchName).setValue(null)
    this.searchLifetimeValues()
  }

  // Format the number in indian rupee system.
  formatLifetimeValueInIndianRupeeSystem(lifetimeValue: number): string {
    var result = lifetimeValue.toString().split('.')
    var lastThree = result[0].substring(result[0].length - 3)
    var otherNumbers = result[0].substring(0, result[0].length - 3)
    if (otherNumbers != '')
      lastThree = ',' + lastThree
    var output = otherNumbers.replace(/\B(?=(\d{2})+(?!\d))/g, ",") + lastThree
    if (result.length > 1) {
      output += "." + result[1]
    }
    return output
  }

  // Get total lifetime value for each entry.
  totalLifetimeValueOfEachEntry(lifetimeValue: any): string {
    let total: number = 0
    if (lifetimeValue.upsell) {
      total = total + lifetimeValue.upsell
    }
    if (lifetimeValue.placement) {
      total = total + lifetimeValue.placement
    }
    if (lifetimeValue.knowledge) {
      total = total + lifetimeValue.knowledge
    }
    if (lifetimeValue.teaching) {
      total = total + lifetimeValue.teaching
    }
    return this.formatLifetimeValueInIndianRupeeSystem(total)
  }

  // Set total reports list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalLifetimeValueEntries < this.paginationEnd) {
      this.paginationEnd = this.totalLifetimeValueEntries
    }
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetLifetimeValueReports() {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilityService.isObjectEmpty(queryParams)) {
      this.getAllLifetimeValueReports()
      return
    }
    this.lifetimeValueSearchForm.patchValue(queryParams)
    this.searchLifetimeValues()
  }

}
