import { DatePipe } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModalRef, NgbModalOptions, NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BatchSessionService } from 'src/app/service/batch-session/batch-session.service';
import { BatchTalentService } from 'src/app/service/batch-talent/batch-talent.service';
import { BatchTopicAssignmentService } from 'src/app/service/batch-topic-assignment/batch-topic-assignment.service';
import { BatchService } from 'src/app/service/batch/batch.service';
import { CourseService } from 'src/app/service/course/course.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { StorageService } from 'src/app/service/storage/storage.service';
import { ITalentReport, TalentDashboardService } from 'src/app/service/talent-dashboard/talent-dashboard.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-talent-dashboard',
  templateUrl: './talent-dashboard.component.html',
  styleUrls: ['./talent-dashboard.component.css']
})
export class TalentDashboardComponent implements OnInit {

  readonly BATCH_STAUS_ONGOING: string = "Ongoing"
  readonly BATCH_STAUS_FINISHED: string = "Finished"
  readonly BATCH_STAUS_UPCOMING: string = "Upcoming"

  // Talent.
  firstName: string
  talentID: string

  // Flags.
  isSessionPresent: boolean
  isSessionToday: boolean
  isFacultyFeddbackForTalentThisAndPreviousWeekPresent: boolean

  // Batch.
  selectedBatchID: string
  telegramBaseURL: string
  currentBatchStatus: string
  batchRelatedDetailsCount: number

  // Batch talent.
  batchTalentDetails: any

  // Modal.
  modalRef: any

  // Batch Session.
  batchSessionList: any[]
  @ViewChild('batchSessionPlanTemplate') batchSessionPlanTemplate: any

  // Faculty feedback for talent.
  facultyToTalentFeedbackLists: any
  facultyFeedbackForTalentThisWeek: number
  facultyFeedbackForTalentThisWeekInSttring: string
  facultyFeedbackForTalentPreviousWeek: number
  facultyFeedbackForTalentLeaderBoard: any[]

  // Talent feedback for faculty.
  @ViewChild('talentFeedbackToFacultyModal') talentFeedbackToFacultyModal: any
  @ViewChild('talentFeedback') talentFeedback: any
  talentFeedbackForFacultyList: any[]
  pendingTalentFeedbackForFacultyList: any[]
  pendingFeedbacksCount: number
  selectedBatchSessionFaculty: any
  selectedBatchSessionForTalentFeedbackID: string
  selectedBatchSessionForTalentFeedbackTopics: any[]
  selectedBatchSessionForTalentFeedbackDate: string
  selectedBatchSessionForTalentFeedbackSessionNumber: number

  // Batch session talent.
  batchSessionTalentList: any[]

  // Batch topic assignment.
  public readonly PENDING_STATUS = "Pending"
  public readonly SUBMITTED_STATUS = "Submitted"
  public readonly COMPLETED_STATUS = "Completed"
  batchTopicAssignmentList: any[]
  pendingBatchTopicAssignmentList: any[]
  submittedBatchTopicAssignmentList: any[]
  completedBatchTopicAssignmentList: any[]
  pendingAssignmentsCount: number
  submittedAssignmentsCount: number
  completedAssignmentsCount: number
  weeklyAssignmentAverageInString: string
  weeklyAssignmentAverage: number
  totalBatchAssignmentCount: number
  totalAssignmentsCompletedCount: number

  // ******************************************** NOT IN USE CURRENTLY ********************************************* 

  // Problem.
  // problemList: IProblemOfTheDayQuestion[]

  // Programming concept.
  conceptList: any[]

  // Leader.
  // leaderBoard: ILeaderBoard
  // performerList: IPerformer[]

  // Cousre.
  courseList: any[]
  batchesForTalentList: any[]

  // My courses.
  myCourseList: any[]
  onGoingCourseList: any[]
  finishedCourseList: any[]

  // Course time table.
  talentTimeTable: ITalentReport

  // Time table related flags.
  isTimetableVisible: boolean
  // isPiechartVisible: boolean

  // Time table label.
  workScheduleLabel: string

  constructor(
    private spinnerService: SpinnerService,
    private router: Router,
    private courseService: CourseService,
    private localService: LocalService,
    private talentDashboardService: TalentDashboardService,
    public utilService: UtilityService,
    // private problemOfTheDayService: ProblemOfTheDayService,
    private batchTalentService: BatchTalentService,
    private batchService: BatchService,
    private batchSessionService: BatchSessionService,
    private datePipe: DatePipe,
    private modalService: NgbModal,
    private route: ActivatedRoute,
    private storageService: StorageService,
    private batchTopicAssignmentService: BatchTopicAssignmentService,
  ) {
    this.initializeVariables()
    this.getAllComponents()
  }

  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize global variables.
  initializeVariables(): void {

    // Talent.
    this.firstName = this.localService.getJsonValue("firstName")
    this.talentID = this.localService.getJsonValue("loginID")

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."

    // Course
    this.courseList = []
    this.batchesForTalentList = []

    // My courses.
    this.myCourseList = []
    this.onGoingCourseList = []
    this.finishedCourseList = []

    // Batch.
    this.telegramBaseURL = "https://telegram.me/"
    this.batchRelatedDetailsCount = 4

    // Batch session.
    this.batchSessionList = []

    // Faculty feedback for talent.
    this.facultyToTalentFeedbackLists = []

    // Flags.
    this.isSessionPresent = false
    this.isSessionToday = false
    this.isFacultyFeddbackForTalentThisAndPreviousWeekPresent = false

    // Time table related flags.
    this.isTimetableVisible = true
    // this.isPiechartVisible = false

    // Time table label.
    this.workScheduleLabel = "Work Schedule This Week"

    // Problem.
    // this.problemList = []

    // Programming concept.
    this.conceptList = []

    // Leader.
    // this.performerList = []

    // Faculty feedback for talent.
    this.facultyFeedbackForTalentThisWeek = 0
    this.facultyFeedbackForTalentPreviousWeek = 0
    this.facultyFeedbackForTalentLeaderBoard = []

    // Talent feedback for faculty.
    this.pendingTalentFeedbackForFacultyList = []
    this.talentFeedbackForFacultyList = []
    this.pendingFeedbacksCount = 0

    // Batch session talent.
    this.batchSessionTalentList = []

    // Batch topic assignment.
    this.pendingBatchTopicAssignmentList = []
    this.pendingAssignmentsCount = 0
    this.submittedBatchTopicAssignmentList = []
    this.submittedAssignmentsCount = 0
    this.completedBatchTopicAssignmentList = []
    this.completedAssignmentsCount = 0
    this.weeklyAssignmentAverage = 0
    this.weeklyAssignmentAverageInString = "-"
  }

  //*********************************************FORMAT FUNCTIONS************************************************************

  // Format fields of course list.
  formatCourseListFields(): void {
    for (let i = 0; i < this.courseList.length; i++) {

      // Give default logo to courses that dont have their own logo.
      if (this.courseList[i].logo == null) {
        this.courseList[i].logo = "assets/logo/java.png"
      }
    }
  }

  // Format the my courses field list.
  formatMyCourseList(): void {

    // Push the courses that have been enrolled in by the talent to the front of the array.
    let tempCourseList: any = []
    for (let i = 0; i < this.courseList.length; i++) {
      if (this.courseList[i].isEnrolled > 0) {
        tempCourseList.push(this.courseList[i]);
      }
    }
    for (let i = 0; i < this.courseList.length; i++) {
      if (this.courseList[i].isEnrolled == 0) {
        tempCourseList.push(this.courseList[i])
      }
    }
    this.courseList = tempCourseList

    // Separate the my course list into ongoing and finished course lists.
    for (let i = 0; i < this.myCourseList.length; i++) {
      if (this.myCourseList[i].batchStatus == "Ongoing") {
        this.onGoingCourseList.push(this.myCourseList[i])
      }
      if (this.myCourseList[i].batchStatus == "Finished") {
        this.finishedCourseList.push(this.myCourseList[i])
      }
    }

    // Convert the fees from number into indian rupee system string.
    for (let i = 0; i < this.onGoingCourseList.length; i++) {
      this.onGoingCourseList[i].feesInString = this.formatNumberInIndianRupeeSystem(this.onGoingCourseList[i].price)
    }
  }

  // Format the number in indian rupee system.
  formatNumberInIndianRupeeSystem(number: number): string {
    if (!number) {
      return
    }
    var result = number.toString().split('.')
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

  // Format fields of programming concept list.
  formatConceptList(): void {
    for (let i = 0; i < this.conceptList.length; i++) {

      // Give default logo to concept that has no logo.
      if (!this.conceptList[i].logo) {
        this.conceptList[i].logo = "assets/logo/default-programming-concept-logo.png"
      }
      this.conceptList[i].description = this.conceptList[i].description.substr(0, 30)
    }
  }

  // // Format leader board.
  // formatLeaderBoard(): void {

  //   // Put all performers in one list.
  //   this.performerList = this.leaderBoard.allPerformers

  //   // If self performer exists then push it in performer list.
  //   this.performerList.splice(0, 0, this.leaderBoard.selfPerformer)

  //   // Give default image to talents that dont have an image.
  //   for (let i = 0; i < this.performerList.length; i++) {
  //     if (!this.performerList[i].image) {
  //       this.performerList[i].image = "assets/images/default-profile-image.png"
  //     }
  //   }
  // }

  // Format the fields of faculty feedback to talent list.
  formatFacultyFeedbackToTalentList(): void {
    this.isFacultyFeddbackForTalentThisAndPreviousWeekPresent = false
    this.facultyFeedbackForTalentThisWeek = 0

    // If there is no feedback for this week and pevious week.
    if (this.facultyToTalentFeedbackLists.thisWeekFeedbacks.length == 0 && this.facultyToTalentFeedbackLists.previousWeekFeedbacks.length == 0) {
      this.facultyToTalentFeedbackLists = null
      return
    }

    // If there is feedback for this week and previous week.
    if (this.facultyToTalentFeedbackLists.thisWeekFeedbacks.length > 0 && this.facultyToTalentFeedbackLists.previousWeekFeedbacks.length > 0) {
      this.isFacultyFeddbackForTalentThisAndPreviousWeekPresent = true
      for (let i = 0; i < this.facultyToTalentFeedbackLists.thisWeekFeedbacks.length; i++) {
        this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].difference = this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].answer - this.facultyToTalentFeedbackLists.previousWeekFeedbacks[i].answer
        this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].previousWeekAnswer = this.facultyToTalentFeedbackLists.previousWeekFeedbacks[i].answer
        this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].percent = this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].difference / this.facultyToTalentFeedbackLists.previousWeekFeedbacks[i].answer * 100
        this.facultyFeedbackForTalentThisWeek = this.facultyFeedbackForTalentThisWeek + this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].answer
        this.facultyFeedbackForTalentPreviousWeek = this.facultyFeedbackForTalentPreviousWeek + this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].previousWeekAnswer
        if (this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].answer % 1 != 0) {
          this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].answer = this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].answer.toFixed(1)
        }
      }
    }

    // If there is no feedback for this week and there is feedback for pevious week.
    if (this.facultyToTalentFeedbackLists.thisWeekFeedbacks.length == 0 && this.facultyToTalentFeedbackLists.previousWeekFeedbacks.length > 0) {
      this.facultyToTalentFeedbackLists.thisWeekFeedbacks = this.facultyToTalentFeedbackLists.previousWeekFeedbacks
      for (let i = 0; i < this.facultyToTalentFeedbackLists.previousWeekFeedbacks.length; i++) {
        this.facultyFeedbackForTalentThisWeek = this.facultyFeedbackForTalentThisWeek + this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].answer
        if (this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].answer % 1 != 0) {
          this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].answer = this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].answer.toFixed(1)
        }
      }
    }

    // If there is no feedback for previous week and there is feedback for this week.
    if (this.facultyToTalentFeedbackLists.thisWeekFeedbacks.length > 0 && this.facultyToTalentFeedbackLists.previousWeekFeedbacks.length == 0) {
      for (let i = 0; i < this.facultyToTalentFeedbackLists.thisWeekFeedbacks.length; i++) {
        this.facultyFeedbackForTalentThisWeek = this.facultyFeedbackForTalentThisWeek + this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].answer
        if (this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].answer % 1 != 0) {
          this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].answer = this.facultyToTalentFeedbackLists.thisWeekFeedbacks[i].answer.toFixed(1)
        }
      }
    }
    this.facultyFeedbackForTalentThisWeek = this.facultyFeedbackForTalentThisWeek / this.facultyToTalentFeedbackLists.thisWeekFeedbacks.length
    if (this.facultyFeedbackForTalentThisWeek % 1 != 0) {
      this.facultyFeedbackForTalentThisWeekInSttring = this.facultyFeedbackForTalentThisWeek.toFixed(1)
    }
    else {
      this.facultyFeedbackForTalentThisWeekInSttring = this.facultyFeedbackForTalentThisWeek.toFixed(0)
    }
    this.facultyFeedbackForTalentPreviousWeek = this.facultyFeedbackForTalentPreviousWeek / this.facultyToTalentFeedbackLists.thisWeekFeedbacks.length
  }

  // Format the faculty feedback for talent ledaer board fields.
  formatFacultyFeedbackForTalentLeaderBoard(): void {
    let rank: number = 1
    this.facultyFeedbackForTalentLeaderBoard.sort(this.sortLeaderBoardByRating)
    for (let i = 0; i < this.facultyFeedbackForTalentLeaderBoard.length; i++) {

      // Give rank.
      if (i == 0){
        this.facultyFeedbackForTalentLeaderBoard[i].rank = rank
      }
      if (i > 0){
        if (this.facultyFeedbackForTalentLeaderBoard[i-1].rating != this.facultyFeedbackForTalentLeaderBoard[i].rating){
          rank = rank + 1
        }
        this.facultyFeedbackForTalentLeaderBoard[i].rank = rank
      }

      // Set the images.
      if (i > 2){
        continue
      }
      this.facultyFeedbackForTalentLeaderBoard[i].imageURL = "assets/icon/talent-dashboard/badge.png"
      if (i == 0) {
        this.facultyFeedbackForTalentLeaderBoard[i].imageURL = "assets/icon/talent-dashboard/first.png"
      }
      if (i == 1) {
        this.facultyFeedbackForTalentLeaderBoard[i].imageURL = "assets/icon/talent-dashboard/second.png"
      }
      if (i == 2) {
        this.facultyFeedbackForTalentLeaderBoard[i].imageURL = "assets/icon/talent-dashboard/third.png"
      }
    }
  }

  // Sort leader borad by rating.
  sortLeaderBoardByRating(a, b) {
    if (a.rating > b.rating) {
      return -1
    }
    if (a.rating < b.rating) {
      return 1
    }
    return 0
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

  // Redirect to problem details page.
  redirectToProblemDetails(problemID): void {
    this.router.navigate(['/problems-of-the-day/problem-details'], {
      queryParams: {
        "problemID": problemID,
        "problemType": "Problem of the day",
      }
    }).catch(err => {
      console.error(err)
    })
  }

  // Redirect to problem of the day page.
  redirectToProblemOfTheDay(): void {
    this.router.navigate(['/problems-of-the-day'], {
    }).catch(err => {
      console.error(err)
    })
  }

  // Redirect to practice page.
  redirectToPractice(): void {
    this.router.navigate(['/practice'], {
    }).catch(err => {
      console.error(err)
    })
  }

  // Redirect to my-batches page.
  redirectToMyCourses(): void {
    this.router.navigate(['/my-batches'], {
    }).catch(err => {
      console.error(err)
    })
  }

  // Redirect to course-details page.
  redirectToCourseDetails(courseID: string): void {
    this.router.navigate(['course-details'], {
      queryParams: {
        "courseID": courseID
      }
    }).catch(err => {
      console.error(err)
    })
  }

  // Redirect to meet link of batch.
  redirectToMeetLink(): void {
    window.open(this.batchTalentDetails.batchMeetLink, "_blank");
  }

  // Redirect to telegram link of faculty.
  redirectToTelegramLink(): void {
    window.open(this.telegramBaseURL + this.batchTalentDetails.batchTelegramLink, "_blank");
  }

  // Redirect to talent batch details.
  redirectToTalentBatchDetails(tab?: string, subTab?: string): void {
    this.router.navigate(['/my-batches'], {
      queryParams: {
        "batchID": this.selectedBatchID,
        "tab": tab,
        "subTab": subTab,
      }
    }).catch(err => {
      console.error(err)
    })
  }

  // Redirect to batch topic assignment.
  redirectTalentbatchDetails(batchTopicAssignmentID: string, tab?: string, subTab?: string): void {
    this.router.navigate(['/my-batches'], {
      queryParams: {
        "batchID": this.selectedBatchID,
        "tab": tab,
        "subTab": subTab,
        "batchTopicAssignmentID": batchTopicAssignmentID,
      }
    }).catch(err => {
      console.error(err)
    })
  }

  // Redirect to Same page.
  redirectToSamePage(): void {
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: {
        "batchID": this.selectedBatchID,
      },
    })
  }

  //********************************************* BATCH TALENT FUNCTIONS ************************************************************

  // Format batch talent details.
  formatBatchTalentDetails(): void {

    // Calculate course completed percentage.
    this.batchTalentDetails.courseCompletedPercentage = this.batchTalentDetails.totalSessionsCompleted / this.batchTalentDetails.totalSessionsCount * 100

    // Format batch timings start time and batch timing end time in AM and PM.
    for (let i = 0; i < this.batchTalentDetails.batchTimings.length; i++) {
      this.batchTalentDetails.batchTimings[i].fromTime = this.utilService.convertTime(this.batchTalentDetails.batchTimings[i].fromTime)
      this.batchTalentDetails.batchTimings[i].toTime = this.utilService.convertTime(this.batchTalentDetails.batchTimings[i].toTime)
    }

    // Convert total completed hours and total hours from minutes to hours.
    this.batchTalentDetails.totalCompletedHours = Math.floor(this.batchTalentDetails.totalCompletedHours / 60)
    this.batchTalentDetails.totalHours = Math.floor(this.batchTalentDetails.totalHours / 60)
  }

  //********************************************* BATCH SESSION FUNCTIONS ************************************************************

  // Set the selected batch whose details will be diplayed on dashboard.
  setSelectedBatch(): void {
    let ongoingBatchList: any[] = []
    let upcomingBatchList: any[] = []
    let finishedBatchList: any[] = []

    for (let i = 0; i < this.batchesForTalentList.length; i++) {

      // Collect all ongoing batches.
      if (this.batchesForTalentList[i].batchStatus == this.BATCH_STAUS_ONGOING) {
        ongoingBatchList.push(this.batchesForTalentList[i])
      }

      // Collect all upcoming batches.
      if (this.batchesForTalentList[i].batchStatus == this.BATCH_STAUS_UPCOMING) {
        upcomingBatchList.push(this.batchesForTalentList[i])
      }

      // Collect all finished batches.
      finishedBatchList.push(this.batchesForTalentList[i])
    }

    // If ongoing batches is present then set it as selected batch.
    if (ongoingBatchList.length > 0) {
      this.currentBatchStatus = this.BATCH_STAUS_ONGOING
      this.setSelectedBatchDetails(ongoingBatchList)
      return
    }

    // If upcoming batches is present then set it as selected batch.
    if (upcomingBatchList.length > 0) {
      this.currentBatchStatus = this.BATCH_STAUS_UPCOMING
      this.setSelectedBatchDetails(upcomingBatchList)
      return
    }

    // If finished batches is present then set it as selected batch.
    if (finishedBatchList.length > 0) {
      this.currentBatchStatus = this.BATCH_STAUS_FINISHED
      this.setSelectedBatchDetails(finishedBatchList)
    }
  }

  // Set selected batch details.
  setSelectedBatchDetails(batchList: any): void{
    this.selectedBatchID = batchList[0].id
    this.storageService.setItem("selectedTalentBatchID", this.selectedBatchID)
    this.getListsRelatedToBatch()
    this.redirectToSamePage()
  }

  // Get all components related to batch.
  getListsRelatedToBatch(): void {
    this.batchRelatedDetailsCount = 4
    this.facultyToTalentFeedbackLists = null
    this.facultyFeedbackForTalentLeaderBoard = []
    if (this.currentBatchStatus == this.BATCH_STAUS_ONGOING) {
      this.getFacultyFeedbackForTalentWeekWise()
      this.getFacultyFeedbackForTalentLeaderBoard()
    }
    // if (this.currentBatchStatus == this.BATCH_STAUS_FINISHED){
    //   this.getFacultyFeedbackForTalent()
    // // }
    this.getBatchTalentDetails()
    this.getBatchSessionList()
    this.getTalentFeedbackForFacultyList()
    this.getBatchSessionTalentList()
    this.getAllBatchTopicAssignmentsForTalent()
  }

  // On changing batches for talent.
  onBathcesForTalentChange(): void{
    
    // Set the current Batch status as status of the batch.
    for (let i = 0; i < this.batchesForTalentList.length; i++) {
      if (this.selectedBatchID == this.batchesForTalentList[i].batchID) {
        this.currentBatchStatus = this.batchesForTalentList[i].batchStatus
      }
    }
    this.redirectToSamePage()
    this.storageService.setItem("selectedTalentBatchID", this.selectedBatchID)
    this.getListsRelatedToBatch()
  } 

  // Decrement batch related details count.
  decrementBatchRelatedDetailsCount(): void {
    this.batchRelatedDetailsCount = this.batchRelatedDetailsCount - 1
    if (this.batchRelatedDetailsCount == 0) {
      this.formatBatchSessionDateList()
    }
  }

  // Format batch session date list.
  formatBatchSessionDateList(): void {
    this.pendingFeedbacksCount = 0
    this.pendingTalentFeedbackForFacultyList = []

    // Give attendance of talent to all batch session date list.
    for (let i = 0; i < this.batchSessionList.length; i++) {
      this.batchSessionList[i].sessionNumber = (i + 1)

      // Set the topic and sub topics for the batch session list.
      let topics: any[] = []
      for (let j = 0; j < this.batchSessionList[i].batchSessionTopics.length; j++) {
        if (j == 0) {
          topics.push(this.batchSessionList[i].batchSessionTopics[j].topic)
          topics[0].subTopics = []
          topics[0].subTopics.push(this.batchSessionList[i].batchSessionTopics[j].subTopic)
          continue
        }
        if (j > 0 && this.batchSessionList[i].batchSessionTopics[j].topic.id == topics[topics.length - 1].id) {
          topics[topics.length - 1].subTopics.push(this.batchSessionList[i].batchSessionTopics[j].subTopic)
        }
        if (j > 0 && this.batchSessionList[i].batchSessionTopics[j].topic.id != topics[topics.length - 1].id) {
          topics.push(this.batchSessionList[i].batchSessionTopics[j].topic)
          topics[topics.length - 1].subTopics = []
          topics[topics.length - 1].subTopics.push(this.batchSessionList[i].batchSessionTopics[j].subTopic)
        }
      }
      this.batchSessionList[i].topics = topics

      // Iterate batch session list.
      for (let j = 0; j < this.batchSessionTalentList.length; j++) {
        if (this.batchSessionList[i].id == this.batchSessionTalentList[j].batchSessionID
          && this.batchSessionTalentList[j].isPresent == true) {
          this.batchSessionList[i].isPresent = true
        }
        if (this.batchSessionList[i].id == this.batchSessionTalentList[j].batchSessionID
          && this.batchSessionTalentList[j].isPresent == false) {
          this.batchSessionList[i].isPresent = false
        }
      }
    }

    // If session is ongoing then show current session.
    if (this.currentBatchStatus == this.BATCH_STAUS_ONGOING) {

      // Iterate batch session date list.
      for (let i = 0; i < this.batchSessionList.length; i++) {

        // Iterate talent feedback for faculty list.
        for (let j = 0; j < this.talentFeedbackForFacultyList.length; j++) {

          // If batch session is complete, attendance is present and feedback is not given then show 
          // option for feedback.
          if (this.batchSessionList[i].id == this.talentFeedbackForFacultyList[j].batchSessionID &&
            this.talentFeedbackForFacultyList[j].sessionFeedbacks?.length == 0 &&
            this.batchSessionList[i].isSessionTaken == true && this.batchSessionList[i].isPresent == true) {
            this.pendingFeedbacksCount = this.pendingFeedbacksCount + 1

            // Only get till 5 pending feedbacks.
            if (this.pendingFeedbacksCount > 5) {
              continue
            }
            this.talentFeedbackForFacultyList[j].sessionNumber = this.batchSessionList[i].sessionNumber
            this.talentFeedbackForFacultyList[j].topics = this.batchSessionList[i].topics
            this.pendingTalentFeedbackForFacultyList.push(this.talentFeedbackForFacultyList[j])
          }
        }
      }
    }
  }

  //********************************************* TALENT FEEDACK FOR FACULTY FUNCTIONS ************************************************************

  // On clicking talent feedback to faculty.
  onTalentFeeddbackToFacultyClick(pendingFeedback: any): void {
    this.selectedBatchSessionForTalentFeedbackID = pendingFeedback.batchSessionID
    this.selectedBatchSessionForTalentFeedbackTopics = pendingFeedback.topics
    this.selectedBatchSessionForTalentFeedbackDate = pendingFeedback.date
    this.selectedBatchSessionForTalentFeedbackSessionNumber = pendingFeedback.sessionNumber
    for (let i = 0; i < this.batchSessionList.length; i++) {
      if (this.batchSessionList[i].date == pendingFeedback.date) {
        this.selectedBatchSessionFaculty = this.batchSessionList[i].faculty
      }
    }
    this.openModal(this.talentFeedbackToFacultyModal, 'lg')
  }

  // Receive is feedback add successful from child.
  receiveIsFeedbackAddSuccessful(isSuccessful: any): void {

    // If modal is closed.
    if (isSuccessful == null) {
      this.modalRef.close()
    }

    // If add is successful.
    if (isSuccessful == true) {
      this.modalRef.close()
      this.getListsRelatedToBatch()
      this.batchSessionPlanTemplate.getCurrentSessionPlanRelatedDetails()
    }
  }

  // If current batch session is updated.
  isCurrentBatchSessionPlanUpdated(): void {
    this.getListsRelatedToBatch()
  }

  //********************************************* BATCH TOPIC ASSIGNMENT FUNCTIONS ************************************************************

  // Format the fields of batch topic assignment list.
  formatBatchTopicAssignmentFields(): void {
    this.pendingAssignmentsCount = 0
    this.completedAssignmentsCount = 0
    this.submittedAssignmentsCount = 0
    this.pendingBatchTopicAssignmentList = []
    this.completedBatchTopicAssignmentList = []
    this.submittedBatchTopicAssignmentList  = []

    // Set the total and completed assignments.
    this.totalBatchAssignmentCount = this.batchTopicAssignmentList.length
    this.totalAssignmentsCompletedCount = 0

    for (let i = 0; i < this.batchTopicAssignmentList.length; i++) {

      // Format submitted on date.
      if (this.batchTopicAssignmentList[i].submissions?.length > 0) {
        this.batchTopicAssignmentList[i].submissions[0].submittedOn = this.datePipe.transform(this.batchTopicAssignmentList[i].submissions[0]?.submittedOn, 'EEE, MMM d, y')
      }

      // Format faculty remarks.
      if (this.batchTopicAssignmentList[i].submissions?.length > 0 && this.batchTopicAssignmentList[i].submissions[0]?.facultyRemarks != null) {
        this.batchTopicAssignmentList[i].facultyRemarks = this.batchTopicAssignmentList[i].submissions[0]?.facultyRemarks
      }

      if (this.batchTopicAssignmentList[i].submissions?.length == 0 || this.batchTopicAssignmentList[i].submissions[0]?.facultyRemarks == null) {
        this.batchTopicAssignmentList[i].facultyRemarks = null
      }

      // Format is Checked field.
      if (this.batchTopicAssignmentList[i].submissions?.length > 0 && this.batchTopicAssignmentList[i].submissions[0]?.isChecked == true) {
        this.batchTopicAssignmentList[i].isChecked = true
      }

      // Set completed status.
      if (this.batchTopicAssignmentList[i].submissions?.length > 0 && this.batchTopicAssignmentList[i].submissions[0].isAccepted) {
        this.batchTopicAssignmentList[i].status = this.COMPLETED_STATUS
        this.totalAssignmentsCompletedCount = this.totalAssignmentsCompletedCount + 1
        continue
      }

      // Set submitted status.
      if (this.batchTopicAssignmentList[i].submissions?.length > 0 && !this.batchTopicAssignmentList[i].submissions[0].isChecked) {
        this.batchTopicAssignmentList[i].status = this.SUBMITTED_STATUS
        continue
      }

      // Set pending status.
      this.batchTopicAssignmentList[i].status = this.PENDING_STATUS
    }

    // Separate the batch topic assignment list by its status.
    for (let i = 0; i < this.batchTopicAssignmentList.length; i++) {

      // Pending.
      if (this.batchTopicAssignmentList[i].status == this.PENDING_STATUS) {
        this.pendingAssignmentsCount = this.pendingAssignmentsCount + 1
        if (this.pendingAssignmentsCount <= 5) {
          this.pendingBatchTopicAssignmentList.push(this.batchTopicAssignmentList[i])
        }
      }

      // Submitted.
      if (this.batchTopicAssignmentList[i].status == this.SUBMITTED_STATUS) {
        this.submittedAssignmentsCount = this.submittedAssignmentsCount + 1
        if (this.submittedAssignmentsCount <= 5) {
          this.submittedBatchTopicAssignmentList.push(this.batchTopicAssignmentList[i])
        }
      }

      // Completed.
      if (this.batchTopicAssignmentList[i].status == this.COMPLETED_STATUS) {
        this.completedAssignmentsCount = this.completedAssignmentsCount + 1
        if (this.completedAssignmentsCount <= 5) {
          this.completedBatchTopicAssignmentList.push(this.batchTopicAssignmentList[i])
        }
      }
    }
  }

  // Format the weekly assignment scores.
  formatWeeklyAssignmentScores(): void {
    let todayDateInDateForm = new Date()
    let firstDayOfWeekInDateForm: Date = new Date(todayDateInDateForm)
    let lastDayOfWeekInDateForm: Date = new Date(todayDateInDateForm)
    let weeklyAssignmentScore: number = 0
    let weeklyAssignmentTotal: number = 0

    // Get range for last 7 days.
    firstDayOfWeekInDateForm.setDate(firstDayOfWeekInDateForm.getDate() - 1)
    lastDayOfWeekInDateForm.setDate(lastDayOfWeekInDateForm.getDate() - 7)

    for (let i = 0; i < this.batchTopicAssignmentList.length; i++) {
      let dueDate = new Date(this.batchTopicAssignmentList[i].dueDate)

      // If date is equal or within the first and last dates of the week and assignment not submitted.
      if ((dueDate <= firstDayOfWeekInDateForm && dueDate >= lastDayOfWeekInDateForm)
        && (this.batchTopicAssignmentList[i].submissions == null || (this.batchTopicAssignmentList[i].submissions != null
          && this.batchTopicAssignmentList[i].submissions?.length == 0))) {
        weeklyAssignmentTotal = weeklyAssignmentTotal + this.batchTopicAssignmentList[i].programmingQuestion.score
      }

      // If date is equal or within the first and last dates of the week and assignment is submitted.
      if ((dueDate <= firstDayOfWeekInDateForm && dueDate >= lastDayOfWeekInDateForm)
        && (this.batchTopicAssignmentList[i].submissions != null && this.batchTopicAssignmentList[i].submissions?.length > 0)) {
        weeklyAssignmentScore = weeklyAssignmentScore + this.batchTopicAssignmentList[i].submissions[0].score
        weeklyAssignmentTotal = weeklyAssignmentTotal + this.batchTopicAssignmentList[i].programmingQuestion.score
      }
    }
    if (weeklyAssignmentTotal > 0) {
      this.weeklyAssignmentAverage = (weeklyAssignmentScore / weeklyAssignmentTotal) * 10
      this.weeklyAssignmentAverageInString = this.weeklyAssignmentAverage.toFixed(1)
    }
  }

  //********************************************* OTHER FUNCTIONS ************************************************************

  // Compare for select option field.
  compareFn(optionOne: any, optionTwo: any): boolean {
    if (optionOne == null && optionTwo == null) {
      return true
    }
    if (optionTwo != undefined && optionOne != undefined) {
      return optionOne.id === optionTwo.id
    }
    return false
  }

  // Used to open modal.
  openModal(content: any, size?: string): NgbModalRef {
    if (!size) {
      size = 'lg'
    }
    let options: NgbModalOptions = {
      ariaLabelledBy: 'modal-basic-title', keyboard: false,
      backdrop: 'static', size: size
    }
    this.modalRef = this.modalService.open(content, options)
    return this.modalRef
  }

  //*********************************************GET FUNCTIONS************************************************************

  // Get all components.
  getAllComponents(): void {
    this.getBatchesForTalent()
    // this.getCourseList()
    // this.getMyCourseList()
    // this.getTalentReport()
    // this.getProblemsOfTheDay()
    // this.getParentProgrammingConceptList()
    // this.getLeaderBoard()
  }

  // Get batches for talent.
  getBatchesForTalent(): void {
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."
    this.batchTalentService.getBatchesForTalent(this.talentID).subscribe((response) => {
      this.batchesForTalentList = response.body
      this.formatBatchesForTalentList()
      this.setSelectedBatch()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get all session dates for batch.
  getBatchSessionList(): void {
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."
    this.batchSessionService.getBatchSessionWithTopicNameList(this.selectedBatchID).subscribe((response: any) => {
      this.batchSessionList = response.body
    }, (err: any) => {
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    }).add(() => {
      this.decrementBatchRelatedDetailsCount()
    })
  }

  // Get batch talent details.
  getBatchTalentDetails(): void {
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."
    this.batchTalentService.getOneBatchTalentDetails(this.selectedBatchID, this.talentID).subscribe((response) => {
      this.batchTalentDetails = response
      this.formatBatchTalentDetails()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.decrementBatchRelatedDetailsCount()
    })
  }

  // Get faculty feedback for talent for current week and previous week.
  getFacultyFeedbackForTalentWeekWise(): void {
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."
    this.talentDashboardService.getFacultyFeedbackForTalentWeekWiseDashboard(this.talentID, this.selectedBatchID).subscribe((response: any) => {
      this.facultyToTalentFeedbackLists = response
      this.formatFacultyFeedbackToTalentList()
    }, (err: any) => {
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  // Get faculty feedback for talent.
  getFacultyFeedbackForTalent(): void {
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."
    this.talentDashboardService.getFacultyFeedbackForTalentDashboard(this.talentID, this.selectedBatchID).subscribe((response: any) => {
      this.facultyToTalentFeedbackLists = {}
      this.facultyToTalentFeedbackLists.previousWeekFeedbacks = []
      this.facultyToTalentFeedbackLists.thisWeekFeedbacks = response
      this.formatFacultyFeedbackToTalentList()
    }, (err: any) => {
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  // Get faculty feedback for talent leader board.
  getFacultyFeedbackForTalentLeaderBoard(): void {
    let queryParams: any = {
      limit: 5
    }
    this.talentDashboardService.getFacultyFeedbackForTalentLeaderBoard(this.selectedBatchID, queryParams).subscribe((response: any) => {
      this.facultyFeedbackForTalentLeaderBoard = response
      this.formatFacultyFeedbackForTalentLeaderBoard()
    }, (err: any) => {
      console.error(err)
    })
  }

  // Get talent feedback for faculty list.
  getTalentFeedbackForFacultyList(): void {
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."
    this.batchService.getTalentBatchFeedback(this.selectedBatchID, this.talentID).subscribe((response) => {
      this.talentFeedbackForFacultyList = response.body
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.decrementBatchRelatedDetailsCount()
    })
  }

  // Get batch session talent list.
  getBatchSessionTalentList(): void {
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."
    this.batchService.getAllBatchSessionTalentsForTalent(this.selectedBatchID, this.talentID).subscribe((response) => {
      this.batchSessionTalentList = response.body
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.decrementBatchRelatedDetailsCount()
    })
  }

  // Get all batch topic assignments for talent.
  getAllBatchTopicAssignmentsForTalent(): void {
    this.batchTopicAssignmentList = []
    this.spinnerService.loadingMessage = "Getting All Assignments..."
    this.batchTopicAssignmentService.getAllBatchTopicAssignmentsForTalent(this.selectedBatchID, this.talentID).subscribe((response) => {
      this.batchTopicAssignmentList = response.body
      this.formatBatchTopicAssignmentFields()
      this.formatWeeklyAssignmentScores()
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err?.error?.error)
    })
  }

  //  ******************************************** NOT IN USE CURRENTLY ********************************************* 

  // Get course list.
  getCourseList(): void {
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."
    this.courseService.getCourseMinimumDetails(this.talentID).subscribe((response) => {
      this.courseList = response
      this.formatCourseListFields()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get all courses.
  getMyCourseList(): void {
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."
    this.courseService.getMyCoursesByTalent(this.talentID).subscribe((response) => {
      this.myCourseList = response
      this.formatMyCourseList()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get talent time table.
  getTalentReport(): void {
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."
    let queryParams: any = {
      talentID: this.talentID,
      date: this.getMonday().toISOString()
    }
    this.talentDashboardService.getTalentReport(this.talentID, queryParams).subscribe((response: any) => {
      this.talentTimeTable = response
    }, (err: any) => {
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    }).add(() => {
      this.isTimetableVisible = true
    })
  }

  // // Get problem of the day questions.
  // getProblemsOfTheDay(): void {
  //   this.spinnerService.loadingMessage = "Getting Dashboard Ready..."


  //   this.problemOfTheDayService.getProblemsOfTheDay().subscribe((response) => {
  //     this.problemList = response
  //   }, error => {
  //     console.error(error)
  //     if (error.statusText.includes('Unknown')) {
  //       alert("No connection to server. Check internet.")
  //     }
  //   })
  // }

  // // Get parent programming concept list.
  // getParentProgrammingConceptList(): void {
  //   this.spinnerService.loadingMessage = "Getting Dashboard Ready..."
  //   
  //   
  //   this.programmingConceptService.getAllParentConcepts().subscribe((response) => {
  //     this.conceptList = response
  //     this.formatConceptList()
  //   }, error => {
  //     console.error(error)
  //     if (error.statusText.includes('Unknown')) {
  //       alert("No connection to server. Check internet.")
  //     }
  //   })
  // }

  // // Get leader board for all questions..
  // getLeaderBoard(): void {
  //   this.spinnerService.loadingMessage = "Getting Dashboard Ready..."


  //   this.problemOfTheDayService.getLeaderBoard(this.talentID).subscribe((response) => {
  //     this.leaderBoard = response
  //     this.formatLeaderBoard()
  //   }, error => {
  //     console.error(error)
  //     if (error.statusText.includes('Unknown')) {
  //       alert("No connection to server. Check internet.")
  //     }
  //   })
  // }

  /**
   * getMonday
   * @returns date depending on the weekIndex value, 0 is for this week,
   *  -7 is for previous week and 7 for next week.
   */
  getMonday(): Date {
    let date = new Date();
    let day = date.getDay();
    let finalDate = new Date();

    let offset = 1 - day;
    if (offset > 0) {
      offset = -6;
    }

    finalDate.setDate(new Date().getDate() + offset);
    return finalDate
  }
}

