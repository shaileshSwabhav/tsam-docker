import { HttpParams } from '@angular/common/http';
import { Component, HostListener, OnInit } from '@angular/core';
import { FormBuilder, FormControl, FormGroup } from '@angular/forms';
import { Router, Event as NavigationEvent } from '@angular/router';
import { NgbModal, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { MenuData } from 'src/app/providers/menu/menu-data';
import { BatchService } from 'src/app/service/batch/batch.service';
import { DashboardService, IFeedbackKeyword, IGroupWiseKeywordName, IKeyword, ISessionGroupScore, ITalentFeedbackScore, ITalentSessionFeedbackScore } from 'src/app/service/dashboard/dashboard.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { IMenu } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  ;
  scrWidth: any;
  resp: boolean;
  salesPeople: any[]
  salesPerson: any
  enquirySource: any
  feebackKeyword: IGroupWiseKeywordName[]

  // batch
  isBatchDashboardVisible: boolean
  isBatchDashboardApiCalled: boolean
  batchPerformanceDetails: ITalentFeedbackScore[]
  sessionWiseTalentScore: ITalentSessionFeedbackScore
  batchTalentPerformanceDetails: ITalentFeedbackScore[]
  batchStatusDetails: ITalentFeedbackScore[]
  batchName: string

  // hide/show variables
  isTechnologyDashboardVisible: boolean
  isTechnologyDashboardApiCalled: boolean


  showCourses: boolean
  isCourseDashboardVisible: boolean
  isCourseDashboardApiCalled: boolean
  showOutstandingTalents: boolean
  showAverageTalents: boolean
  showGoodTalents: boolean
  isBatchTalentDetailsVisible: boolean
  isFeelingVisible: boolean

  // modal
  modalRef: any

  // components
  batchStatusList: any[]
  facultyList: any[]
  batchList: any[]
  sourceList: any[]

  // forms
  batchForm: FormGroup

  // search
  searchBatchFormValue: any

  // enquiry source
  enquirySourceDetails: any[]


  // feedback keyword score
  feedbackKeywordScore: IFeedbackKeyword[]

  template: string = `<div class="loading pointer"> <span>Loading</span> </div>`

  constructor(
    private formBuilder: FormBuilder,
    private dashboardService: DashboardService,
    private spinnerService: SpinnerService,
    private generalService: GeneralService,
    private batchService: BatchService,
    private util: UtilityService,
    private router: Router,
    private modalService: NgbModal,
  ) {
    this.getScreenSize();
    this.initializeVariables()
    this.createForms()
    this.getAllComponents()
  }

  initializeVariables() {
    this.spinnerService.loadingMessage = "Loading data.."
    this.batchName = ""

    this.resp = false
    this.salesPeople = []
    this.salesPerson = null
    this.enquirySource = null

    this.showCourses = false
    // technology
    this.isTechnologyDashboardVisible = false
    this.isTechnologyDashboardApiCalled = false
    // course
    this.isCourseDashboardVisible = false
    this.isCourseDashboardApiCalled = false
    // batch
    this.isBatchDashboardVisible = false
    this.isBatchDashboardApiCalled = false
    this.showOutstandingTalents = false
    this.isBatchTalentDetailsVisible = false
    this.isFeelingVisible = false

    this.batchList = []
    this.batchStatusList = []
    this.facultyList = []
    this.feedbackKeywordScore = []
    this.sourceList = []

    this.enquirySourceDetails = []

  }

  createForms(): void {
    this.createBatchForm()
  }

  getAllSalesPeople(): void {
    this.generalService.getSalesPersonList().subscribe((data: any) => {
      this.salesPeople = data.body
    }, (err) => {
      console.error(err)
    })
  }


  @HostListener('window:resize', ['$event'])
  getScreenSize(event?: any) {

    this.scrWidth = window.innerWidth;
    // console.log(this.scrWidth)

    if (this.scrWidth < 993) {
      this.resp = true
    }
    else if (this.scrWidth >= 993) {
      this.resp = false;
    }
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit() {

    this.getAllSalesPeople()
    this.getAdminDashboardDetails()
  }

  getAllComponents(): void {
    this.getBatchStatusList()
    this.getFacultyList()
    this.getSourceList()
    this.getBatchList(this.batchForm.get('batchStatus')?.value, this.batchForm.get('facultyID')?.value)
  }

  // Need to change 
  adminDashboard: any = {
    salesPeopleDetails: {},
    talentDetails: {},
    enquiryDetails: {},
    collegeDetails: {},
    companyDetails: {},
    courseDetails: {},
    facultyDetails: {},
    batchDetails: {},
    technologyDetails: {},
    // batch
    batchScore: {},
    batchStatusDetails: {},
    batchTalents: {},
  };

  getSalesPeopleDashboardDetails() {

    // let params = new HttpParams();
    // params = params.append('salesPersonID', this.salesPerson.id);

    let queryParams: any = {}
    if (this.salesPerson) {
      queryParams.salesPersonID = this.salesPerson.id
    }
    this.getTalentEnquirySourceDetails()

    this.dashboardService.getSalesPeopleDashboardDetails(queryParams).subscribe(data => {

      this.adminDashboard.salesPeopleDetails = data.body;

    }, err => {
      console.error(err)

    })
  }

  getTalentDashboardDetails() {

    this.dashboardService.getTalentDashboardDetails().subscribe(data => {

      this.adminDashboard.talentDetails = data.body;

    }, err => {
      console.error(err)

    })
  }

  getTalentEnquiryDashboardDetails() {


    let queryParams: any = {}
    if (this.enquirySource) {
      queryParams.sourceID = this.enquirySource.id
    }

    this.dashboardService.getTalentEnquiryDashboardDetails(queryParams).subscribe(data => {

      this.adminDashboard.enquiryDetails = data.body;

    }, err => {
      console.error(err)

    })
  }

  getTalentEnquirySourceDetails() {
    let queryParams: any = {}
    if (this.salesPerson) {
      queryParams.salesPersonID = this.salesPerson.id
    }


    this.dashboardService.getTalentEnquirySource(queryParams).subscribe(data => {

      this.enquirySourceDetails = data.body

    }, err => {
      console.error(err)

    })
  }

  getCollegeDashboardDetails() {

    this.dashboardService.getCollegeDashboardDetails().subscribe(data => {

      this.adminDashboard.collegeDetails = data.body;

    }, err => {
      console.error(err)

    })
  }

  getFacultyDashboardDetails() {

    this.dashboardService.getFacultyDashboardDetails().subscribe(data => {

      this.adminDashboard.facultyDetails = data.body;

    }, err => {
      console.error(err)

    })
  }

  getCompanyDashboardDetails() {

    this.dashboardService.getCompanyDashboardDetails().subscribe(data => {

      this.adminDashboard.companyDetails = data.body;

    }, err => {
      console.error(err)

    })
  }

  getCourseDashboardDetails() {
    this.spinnerService.loadingMessage = "Loading Courses..."

    this.dashboardService.getCourseDashboardDetails().subscribe(data => {
      this.adminDashboard.courseDetails = data.body;
      this.isCourseDashboardApiCalled = true

    }, err => {
      console.error(err)

    })
  }

  getBatchDashboardDetails() {

    this.dashboardService.getBatchDashboardDetails().subscribe(data => {
      this.adminDashboard.batchDetails = data.body;
      // console.log(this.adminDashboard.batchDetails);
      this.adminDashboard.batchDetails.live.isVisible = false
      this.adminDashboard.batchDetails.completed.isVisible = false
      this.adminDashboard.batchDetails.upcoming.isVisible = false

    }, err => {
      console.error(err)

    })
  }

  getTechnologyDashboardDetails() {
    this.spinnerService.loadingMessage = "Loading Technology Data..."

    this.dashboardService.getTechnologyDashboardDetails().subscribe(data => {

      this.adminDashboard.technologyDetails = data.body;
      this.isTechnologyDashboardApiCalled = true

    }, err => {
      console.error(err)

    })
  }

  showTechnologyDashboard() {
    this.isTechnologyDashboardVisible = !this.isTechnologyDashboardVisible
    if (!this.isTechnologyDashboardApiCalled) {
      this.getTechnologyDashboardDetails()
    }
  }

  showCourseDashboard() {
    this.isCourseDashboardVisible = !this.isCourseDashboardVisible
    if (!this.isCourseDashboardApiCalled) {
      this.getCourseDashboardDetails()
    }
  }

  showCoursesData() {
    this.showCourses = true
  }

  getAdminDashboardDetails() {
    this.getAdminDashboardData()
    this.getSalesPeopleDashboardDetails()
    this.getFacultyDashboardDetails()
    this.getTalentDashboardDetails()
    this.getTalentEnquiryDashboardDetails()
    this.getCollegeDashboardDetails()
    this.getCompanyDashboardDetails()
    this.getBatchDashboardDetails()
  }

  getAdminDashboardData() {
    // this.spinnerService.loadingMessage="Getting data.."
    this.dashboardService.getAllDashboardDetails().subscribe(data => {

      this.adminDashboard.technologyDetails.totalCount = data.body?.totalTechnologies
      this.adminDashboard.courseDetails.totalCourses = data.body?.totalCourses
      // 
    }, err => {
      console.error(err)

      if (err.status == 400) {
        alert(err.error);
        this.util.rediectTo("login")
      }
    })
  }

  getBatchPerformance() {
    this.spinnerService.loadingMessage = "Loading Batch Data..."

    this.dashboardService.getBatchPerformance(this.searchBatchFormValue).subscribe((response: any) => {
      this.adminDashboard.batchScore = response.body
      this.feebackKeyword = response.body.keywordNames
      this.isBatchDashboardApiCalled = true
      // console.log(this.adminDashboard.batchScore);

    }, (err: any) => {
      console.error(err)

    })
  }

  assignFeedbackKeywordFields(): void {
    for (let index = 0; index < this.feebackKeyword.length; index++) {
      this.feebackKeyword[index].showGroupDetails = false
    }
  }

  createBatchForm(): void {
    this.batchForm = this.formBuilder.group({
      batchStatus: new FormControl("Ongoing"),
      batchID: new FormControl(null),
      facultyID: new FormControl(null)
    })
  }

  showBatchDashboard(): void {
    this.isBatchDashboardVisible = !this.isBatchDashboardVisible
    if (!this.isBatchDashboardApiCalled) {
      this.getSearchBatchPerformance()
    }
  }

  getSearchBatchPerformance(): void {
    this.showOutstandingTalents = false
    this.searchBatchFormValue = { ...this.batchForm.value }
    let flag: boolean = true

    for (let field in this.searchBatchFormValue) {
      if (!this.searchBatchFormValue[field]) {
        delete this.searchBatchFormValue[field]
      } else {
        flag = false
      }
    }
    // No API call on empty search.
    if (flag) {
      return
    }
    this.getBatchPerformance()
  }


  getOutstandingTalents(): void {
    this.showOutstandingTalents = !this.showOutstandingTalents
    this.showGoodTalents = false
    this.showAverageTalents = false
    this.batchPerformanceDetails = this.adminDashboard.batchScore.outstandingTalent
  }

  getGoodTalents(): void {
    this.showGoodTalents = !this.showGoodTalents
    this.showOutstandingTalents = false
    this.showAverageTalents = false
    this.batchPerformanceDetails = this.adminDashboard.batchScore.goodTalent
  }

  getAverageTalents(): void {
    this.showAverageTalents = !this.showAverageTalents
    this.showGoodTalents = false
    this.showOutstandingTalents = false
    this.batchPerformanceDetails = this.adminDashboard.batchScore.averageTalent
  }

  getTalentScore(performance: ITalentFeedbackScore, modalContent: any) {

    if (performance.score == 0) {
      alert("Talent has no feedbacks")
      return
    }
    this.spinnerService.loadingMessage = "Getting session wise score"

    this.sessionWiseTalentScore = null
    this.dashboardService.getSessionWiseTalentScore(performance.talentID, performance.batchID).subscribe((response: any) => {
      this.sessionWiseTalentScore = response.body
      this.assignFeedbackKeywordFields()
      this.isFeelingVisible = false
      // console.log(this.sessionWiseTalentScore)
      this.openModal(modalContent, "xl")

    }, (err: any) => {
      console.error(err)

    })
  }

  assignFeedbackScoreFields(): void {
    for (let index = 0; index < this.sessionWiseTalentScore.keywordNames.length; index++) {
      this.sessionWiseTalentScore.keywordNames[index].showGroupDetails = false
    }
  }

  calculateAvgGroupScore(groupScore: ISessionGroupScore): number {
    let score = 0
    for (let index = 0; index < groupScore.feedbackScore.length; index++) {
      score += groupScore.feedbackScore[index].keywordScore
    }
    return score / groupScore.feedbackScore.length
  }

  getBatchDetails(batchStatus: string, totalBatches: number): void {

    if (totalBatches == 0) {
      alert("There are no " + batchStatus.toLowerCase() + " batches")
      return
    }

    if (batchStatus == "Ongoing" && !this.adminDashboard.batchDetails.live.isVisible) {
      return
    }

    if (batchStatus == "Finished" && !this.adminDashboard.batchDetails.completed.isVisible) {
      return
    }

    if (batchStatus == "Upcoming" && !this.adminDashboard.batchDetails.upcoming.isVisible) {
      return
    }

    // let params = new HttpParams()
    let queryParams: any = {}
    if (batchStatus) {
      // params = params.append("batchStatus", batchStatus)
      queryParams.batchStatus = batchStatus
    }
    this.spinnerService.loadingMessage = "Getting batch details"

    this.dashboardService.getBatchDetails(queryParams).subscribe((response: any) => {
      this.adminDashboard.batchStatusDetails = response.body

      if (batchStatus == "Ongoing") {
        this.adminDashboard.batchDetails.live.batchStatusDetails = response.body
        for (let index = 0; index < this.adminDashboard.batchDetails.live.batchStatusDetails.length; index++) {
          this.adminDashboard.batchDetails.live.batchStatusDetails[index].isVisible = false
        }
      }

      if (batchStatus == "Upcoming") {
        this.adminDashboard.batchDetails.upcoming.batchStatusDetails = response.body
        for (let index = 0; index < this.adminDashboard.batchDetails.upcoming.batchStatusDetails.length; index++) {
          this.adminDashboard.batchDetails.upcoming.batchStatusDetails[index].isVisible = false
        }
      }

      if (batchStatus == "Finished") {
        this.adminDashboard.batchDetails.completed.batchStatusDetails = response.body
        for (let index = 0; index < this.adminDashboard.batchDetails.completed.batchStatusDetails.length; index++) {
          this.adminDashboard.batchDetails.completed.batchStatusDetails[index].isVisible = false
        }
      }


    }, (err: any) => {
      console.error(err)

    })
  }

  getBatchTalent(batchID: string, batchName: string, modalContent: any): void {

    this.batchName = batchName
    this.batchTalentPerformanceDetails = []
    this.spinnerService.loadingMessage = "Getting batch talents"

    this.dashboardService.getBatchTalents(batchID).subscribe((response: any) => {
      // console.log(response.body);
      this.batchTalentPerformanceDetails = response.body
      this.isBatchTalentDetailsVisible = true
      this.assignFeedbackKeywordFields()
      this.openModal(modalContent, "xl")

    }, (err: any) => {
      console.error(err)

    })
  }

  openModal(modalContent: any, modalSize?: string): void {
    if (modalSize == undefined) {
      modalSize = 'lg'
    }
    this.modalRef = this.modalService.open(modalContent,
      { ariaLabelledBy: 'modal-basic-title', backdrop: 'static', size: modalSize, keyboard: false }
    );
    /*this.modalRef.result.subscribe((result) => {
    }, (reason) => {

    });*/
  }

  goTo(path: String) {
    this.router.navigateByUrl("/" + path)
  }

  // =============================================================COMPONENTS=============================================================


  // Get All Batch Status
  getBatchStatusList(): void {
    this.generalService.getGeneralTypeByType("batch_status").subscribe((respond: any) => {
      // console.log(respond);
      this.batchStatusList = respond;
    }, (err) => {
      console.error(err)
    })
  }

  getBatchList(batchStatus: string, facultyID: string): void {
    // let batchStatus = this.batchForm.controls["batchStatus"].value
    // console.log(batchStatus);
    let queryParams: any = {
      is_Active: '1'
    }

    // let batchParams = new HttpParams()

    if (batchStatus) {
      // batchParams = batchParams.append("batchStatus", batchStatus)
      queryParams.batchStatus = batchStatus
    }
    if (facultyID) {
      // batchParams = batchParams.append("facultyID", facultyID)
      queryParams.facultyID = facultyID
    }
    // batchParams = batchParams.append("is_Active", '1')

    this.batchService.getBatchList(queryParams).subscribe((response: any) => {
      this.batchList = response.body
      this.batchForm.controls["batchID"].setValue(null)
    }, (err) => {
      console.error(err)
    })
  }

  // Get All Faculty List
  getFacultyList(): void {
    this.generalService.getFacultyList().subscribe((data: any) => {
      this.facultyList = data.body
    }, (err) => {
      console.error(err);
    })
  }

  getSourceList(): void {
    this.generalService.getSources().subscribe(data => {
      this.sourceList = data
    }, (err) => {
      console.error(err);
    })
  }


  toggleGroupVisibility(group: IGroupWiseKeywordName): void {
    group.showGroupDetails = !group.showGroupDetails
  }

  toggleFeelingVisibility(): void {
    this.isFeelingVisible = !this.isFeelingVisible
  }

  dismissBatchTalentModal(modal: NgbModalRef): void {
    this.assignFeedbackKeywordFields()
    modal.close()
  }

}
