import { DatePipe } from '@angular/common';
import { ChangeDetectorRef, Component, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BatchService } from 'src/app/service/batch/batch.service';
import { StorageService } from 'src/app/service/storage/storage.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { CourseService, ICourseList } from 'src/app/service/course/course.service';
import { BatchTalentService } from 'src/app/service/batch-talent/batch-talent.service';
import { LocalService } from 'src/app/service/storage/local.service';

@Component({
  selector: 'app-talent-batch-details',
  templateUrl: './talent-batch-details.component.html',
  styleUrls: ['./talent-batch-details.component.css']
})
export class TalentBatchDetailsComponent implements OnInit {

  // Tab.
  batchDetailsTabList: any[]
  selectedTemplateName: any
  tab: string
  subTab: string
  lockStatus: boolean
  @ViewChild("batchDetailsTemplate") batchDetailsTemplate: any
  @ViewChild("talentFeedbackSubmitTemplate") talentFeedbackSubmitTemplate: any
  @ViewChild("talentFeedbackLeaderboardTemplate") talentFeedbackLeaderboardTemplate: any
  @ViewChild("talentBatchSessionPlanTemplate") talentBatchSessionPlanTemplate: any
  @ViewChild("talentBatchSessionPlanViewAllTemplate") talentBatchSessionPlanViewAllTemplate: any
  @ViewChild("talentAssignmentTemplate") talentAssignmentTemplate: any
  @ViewChild("talentProjectTemplate") talentProjectTemplate: any
  @ViewChild("talentConceptTreeTemplate") talentConceptTreeTemplate: any
  @ViewChild("talentPerformanceTemplate") talentPerformanceTemplate: any

  // Batch.
  batchDetailsList: any[]
  telegramBaseURL: string
  batchID: string
  batch: any
  batchesForTalentList: any[]

  // Batch topic assignment.
  batchTopicAssignmentID: string
  batchProjectID : string

  // Flags.
  isSidePanelOpen: boolean

  constructor(
    private spinnerService: SpinnerService,
    private batchService: BatchService,
    private route: ActivatedRoute,
    private datePipe: DatePipe,
    private cdr: ChangeDetectorRef,
    private storageService: StorageService,
    private router: Router,
    private utilityService: UtilityService,
    private batchTalentService: BatchTalentService,
    private localService: LocalService,
  ) {
    this.initializeVariables()
  }

  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
    this.getBatchesForTalent()
  }

  ngAfterViewInit() {
    this.initializeTabs()
    this.cdr.detectChanges()
  }

  // Initialize global variables.
  initializeVariables(): void {

    // Tabs.
    this.lockStatus = false

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Batch Details..."

    // Batch.
    this.batchesForTalentList = []

    // Query params.
    this.telegramBaseURL = "https://telegram.me/"
    this.route.queryParamMap.subscribe(params => {
      this.batchID = params.get("batchID")
      this.tab = params.get("tab")
      this.subTab = params.get("subTab")
      this.batchTopicAssignmentID = params.get("batchTopicAssignmentID")
      this.batchProjectID= params.get("batchProjectID")
      if (!this.batchID) {
        this.batchID = this.storageService.getItem("selectedTalentBatchID")
      }
      this.getBatchDetails()
      this.redirectToSamePage()
    }, err => {
      console.error(err)
    })

    // Flags.
    this.isSidePanelOpen = false
  }

  // Initialize the batch details tabs.
  initializeTabs(): void {
    this.batchDetailsTabList = [
      {
        tabName: "My Class", subTabs: [
          { tabName: "Batch Details", isActive: false, isRedirect: false, url: null, isVisible: true, templateName: this.batchDetailsTemplate },
          // { tabName: "Brochure", isActive: false, isRedirect: true, url: null, isVisible: true },
        ]
      },
      {
        tabName: "Sessions", subTabs: [
          { tabName: "Session Plan", isActive: false, isRedirect: false, url: null, isVisible: true, templateName: this.talentBatchSessionPlanTemplate },
          { tabName: "View All", isActive: false, isRedirect: false, url: null, isVisible: true, templateName: this.talentBatchSessionPlanViewAllTemplate },
        ]
      },
      {
        tabName: "Assignments", subTabs: [
          { tabName: "Submit", isActive: false, isRedirect: false, url: null, isVisible: true, templateName: this.talentAssignmentTemplate },
        ]
      },
      {
        tabName: "Projects", subTabs: [
          { tabName: "Submit", isActive: false, isRedirect: false, url: null, isVisible: true, templateName: this.talentProjectTemplate },
        ]
      },
      {
        tabName: "Feedback", subTabs: [
          { tabName: "Submit", isActive: false, isRedirect: false, url: null, isVisible: true, templateName: this.talentFeedbackSubmitTemplate },
          { tabName: "Leaderboard", isActive: false, isRedirect: false, url: null, isVisible: true, templateName: this.talentFeedbackLeaderboardTemplate },
        ]
      },
      {
        tabName: "My Progress Report", subTabs: [
          {
            tabName: "Concept Tree", isActive: false, isRedirect: false, url: "my-batches/concept-tree", isVisible: true, templateName: this.talentConceptTreeTemplate,
            queryParams: { batchID: this.batchID }
          },
          { tabName: "My Performance", isActive: false, isRedirect: false, url: null, isVisible: true, templateName: this.talentPerformanceTemplate },
        ]
      }
    ]

    // If no tab then set it as default.
    if (!this.tab) {
      this.tab = "My Class"
      this.subTab = "Batch Details"
    }
    this.onSubTabClick(this.tab, this.subTab)
  }

  //********************************************* BATCH DETAILS FUNCTIONS ************************************************************

  // Format the batch details list.
  formatBatchDetailsList(): void {
    this.batchDetailsList = [
      {
        fieldName: "Course Name",
        fieldValue: this.batch.course.name
      },
      {
        fieldName: "Batch Name",
        fieldValue: this.batch.batchName
      },
      {
        fieldName: "Total Students",
        fieldValue: this.batch.totalIntake
      },
      {
        fieldName: "Starting Date",
        fieldValue: this.datePipe.transform(this.batch.startDate, 'EEEE, MMMM d, y')
      },
      {
        fieldName: "End Date",
        fieldValue: this.datePipe.transform(this.batch.endDate, 'EEEE, MMMM d, y')
      },
      {
        fieldName: "Total Hours",
        fieldValue: Math.floor(this.batch.totalHours / 60)
      },
    ]

    for (let i = 0; i < this.batch.faculty?.length; i++) {
      if (this.batch.faculty != null) {
        this.batchDetailsList.push({
          fieldName: "Mentor Name",
          fieldValue: this.batch.faculty[i].firstName + " " + this.batch.faculty[i].lastName
        },
          {
            fieldName: "Email Id",
            fieldValue: this.batch.faculty[i].email
          },
          {
            fieldName: null,
            fieldValue: null
          })
      }
    }

    // Format batch timings start time and batch timing end time in AM and PM.
    for (let i = 0; i < this.batch.batchTimings.length; i++) {
      this.batch.batchTimings[i].fromTime = this.utilityService.convertTime(this.batch.batchTimings[i].fromTime)
      this.batch.batchTimings[i].toTime = this.utilityService.convertTime(this.batch.batchTimings[i].toTime)
    }

    // Sort batch timings by order of days.
    this.batch.batchTimings.sort(this.sortBatchTimings)

    // Give redirection to tabs with url. 
    // for (let i = 0; i < this.batchDetailsTabList.length; i++) {
    //   if (this.batchDetailsTabList[i].subTabs) {
    //     for (let j = 0; j < this.batchDetailsTabList[i].subTabs.length; j++) {

    //       // If brochure is not present then hide the brochure tab.
    //       if (this.batchDetailsTabList[i].subTabs[j].tabName == "Brochure" && this.batchDetailsTabList[i].tabName == "My Class" &&
    //         this.batch.brochure == null) {
    //         this.batchDetailsTabList[i].subTabs[j].isVisible = false
    //         break
    //       }

    //       // If brochure is present then set the brochure tab's url.
    //       if (this.batchDetailsTabList[i].subTabs[j].tabName == "Brochure" && this.batchDetailsTabList[i].tabName == "My Class" &&
    //         this.batch.brochure != null) {
    //         this.batchDetailsTabList[i].subTabs[j].url = this.batch.brochure
    //         break
    //       }
    //     }
    //   }
    // }
  }

  //********************************************* FORMAT FUNCTIONS ************************************************************

  // On clicking sub tab.
  onSubTabClick(tabName: string, subTabName: string, url?: string): void {

    // If tab is for redirection.
    for (let i = 0; i < this.batchDetailsTabList.length; i++) {
      if (this.batchDetailsTabList[i].subTabs) {
        for (let j = 0; j < this.batchDetailsTabList[i].subTabs.length; j++) {
          if (this.batchDetailsTabList[i].subTabs[j].tabName == subTabName && this.batchDetailsTabList[i].tabName == tabName) {
            if (this.batchDetailsTabList[i].subTabs[j].isRedirect) {
              this.redirectToExternalLink(url)
              return
            }
          }
        }
      }
    }

    // If tab is not for redirection.
    for (let i = 0; i < this.batchDetailsTabList.length; i++) {
      if (this.batchDetailsTabList[i].subTabs) {
        for (let j = 0; j < this.batchDetailsTabList[i].subTabs.length; j++) {
          if (this.batchDetailsTabList[i].subTabs[j].tabName == subTabName && this.batchDetailsTabList[i].tabName == tabName) {
            this.tab = tabName
            this.subTab = subTabName
            this.batchDetailsTabList[i].subTabs[j].isActive = true
            if (this.batchDetailsTabList[i].subTabs[j].url != null) {
              this.redirectToOtherPage(this.batchDetailsTabList[i].subTabs[j].url, this.batchDetailsTabList[i].subTabs[j].queryParams)
            }
            if (this.batchDetailsTabList[i].subTabs[j].url == null) {
              this.selectedTemplateName = this.batchDetailsTabList[i].subTabs[j].templateName
              this.redirectToSamePage()
            }
            continue
          }
          this.batchDetailsTabList[i].subTabs[j].isActive = false
        }
      }
    }
  }

  // Format fields of batches for talent list.
  formatBatchesForTalentList(): void{
    let tempBatchesForTalentList: any[] = []
    for (let i = 0; i < this.batchesForTalentList.length; i++){
      if (this.batchesForTalentList[i].course){

        // Push those batches for talent whose course exists.
        tempBatchesForTalentList.push(this.batchesForTalentList[i])
      }
    }
    this.batchesForTalentList = tempBatchesForTalentList
  }

  //********************************************* REDIRECT FUNCTIONS ************************************************************

  // Redirect to telegram link of faculty.
  redirectToTelegramLink(): void {
    window.open(this.telegramBaseURL + this.batch.batchTelegramLink, "_blank")
  }

  // Redirect to meet link of batch.
  redirectToMeetLink(): void {
    window.open(this.batch.batchMeetLink, "_blank")
  }

  // Redirect to external link.
  redirectToExternalLink(url): void {
    window.open(url, "_blank")
  }

  // Redirect to same page.
  redirectToSamePage(): void {
    let queryParams: any = {
      "batchID": this.batchID,
      "tab": this.tab,
      "subTab": this.subTab,
    }
    if (this.tab == "Assignments" && this.subTab == "Submit") {
      queryParams.batchTopicAssignmentID = this.batchTopicAssignmentID
    }
    if (this.tab == "Projects" && this.subTab == "Submit") {
      queryParams.batchProjectID = this.batchProjectID
    }
    this.router.navigate(['/my-batches'], { queryParams }).catch(err => {
      console.error(err)
    })
  }

  // Redirect to same page but only changing the batch id.
  redirectOnBatchForTalentChange(): void {
    let queryParams: any = {
      batchID: this.batchID 
    }
    this.router.navigate([],
      { relativeTo: this.route,
        queryParams: queryParams,
        queryParamsHandling: "merge"
      } 
    ).catch(err => {
      console.error(err)
    })
  }

  // Redirect to other page.
  redirectToOtherPage(url: string, queryParams: any): void {
    this.router.navigate([url], { queryParams }).catch(err => {
      console.error(err)
    })
  }

  //********************************************* OTHER FUNCTIONS ************************************************************

  // On changing batches for talent.
  onBathcesForTalentChange(): void{
    this.redirectOnBatchForTalentChange()
    if (this.tab == "My Class" && this.subTab == "Batch Detail"){
      this.getBatchDetails()
    }
  } 

  // Sort batch timings.
  sortBatchTimings(a, b) {
    if (a.day.order < b.day.order) {
      return -1
    }
    if (a.day.order > b.day.order) {
      return 1
    }
    return 0
  }

  //*********************************************GET FUNCTIONS************************************************************

  // Get course list.
  getBatchDetails(): void {
    this.spinnerService.loadingMessage = "Getting Batch Details..."
    this.batchService.getBatchDetails(this.batchID).subscribe((response) => {
      this.batch = response
      this.formatBatchDetailsList()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get batches for talent.
  getBatchesForTalent(): void {
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."
    this.batchTalentService.getBatchesForTalent(this.localService.getJsonValue("loginID")).subscribe((response) => {
      this.batchesForTalentList = response.body
      this.formatBatchesForTalentList()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

}
