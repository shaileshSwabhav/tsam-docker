import { DatePipe } from '@angular/common';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BatchSessionService } from 'src/app/service/batch-session/batch-session.service';
import { BatchService } from 'src/app/service/batch/batch.service';
import { AccessLevel, Role, UrlConstant } from 'src/app/service/constant';
import { GeneralService } from 'src/app/service/general/general.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { ProgrammingQuestionModalComponent } from '../programming-question-modal/programming-question-modal.component';
import { IModuleTopic } from 'src/app/models/course/module_topic';
import { IModule } from 'src/app/models/course/module';
import { ITopicProgrammingQuestion } from 'src/app/models/course/topic_programming_question';
import { IBatchTopicAssignment } from 'src/app/models/batch-topic-assignment';
import { ModuleService } from 'src/app/service/module/module.service';
import { BatchTopicAssignmentService } from 'src/app/service/batch-topic-assignment/batch-topic-assignment.service';

@Component({
  selector: 'app-batch-session-plan',
  templateUrl: './batch-session-plan.component.html',
  styleUrls: ['./batch-session-plan.component.css']
})
export class BatchSessionPlanComponent implements OnInit {

  // Components.
  dayList: any[]
  loginID: string

  // Batch.
  batchID: string
  batchRelatedDetailsCount: number
  batchName: string

  // Batch details.
  batchDetailsList: any[]

  // Flags.
  isCurrentSessionToday: boolean
  isSessionTaken: boolean
  isTalentPresent: boolean
  isUpdatePreRequisiteMode: boolean
  showPrerequisiteForm: boolean

  // Session plan.
  batchSessionList: any[]
  dateList: any[]
  currentDateListIndex: number
  currentDateListDateIndex: number
  currentBatchSessionDate: string
  toDayDateInDateForm: Date
  currentBatchSessionPlan: any
  currentBatchSession: any
  actualBatchSessionIndex: number
  actualBatchSessionDate: string
  currentBatchSessionPlanCurrentSubtopicCount: number
  currentBatchSessionPlanPendingSubtopicCount: number
  currentBatchSessionPlanCurrentAssignmentCountForFaculty: number
  currentBatchSessionPlanCurrentAssignmentCountForTalent: number
  currentBatchSessionPlanPendingAssignmentCountForFaculty: number
  currentBatchSessionPlanPendingAssignmentCountForTalent: number
  batchSessionsCounts: any
  public readonly MONTHS: any[] = ["January", "February", "March", "April", "May", "June", "July", "August",
    "September", "October", "November", "December"]
  public readonly ATTENDANCEOPERATION: string = "attendance"
  public readonly COMPLETEDOPERATION: string = "completed"
  public readonly FEEDBACKOPERATION: string = "feedback"
  public readonly ASSIGNMENTOPERATION: string = "assignment"
  @ViewChild('sessionCompletionStatus') sessionCompletionStatus: any

  // Programming question.
  selectedAssignment: any
  modalRef: any
  @ViewChild('programmingQuestionModal') programmingQuestionModal: any

  // Talent feedback for faculty.
  selectedBatchSessionFaculty: any
  talentFeedbackForFacultyList: any[]
  selectedBatchSessionForTalentFeedbackSessionNumber: number
  @ViewChild('talentFeedbackToFacultyModal') talentFeedbackToFacultyModal: any

  // Batch session talent.
  batchSessionTalentList: any[]

  // Talent.
  talentID: string

  // Values to parent.
  @Output() isUpdatedEmitter = new EventEmitter<any>()

  // PreRequisite.
  preRequisiteForm: FormGroup
  @ViewChild('preRequisiteModal') preRequisiteModal: any
  @ViewChild("ckeditorPreRequisite") ckeditorPreRequisite: any

  // Cke editor configuration.
  ckeEditorConfig: any

  // Access.
  access: any
  permission: IPermission

  // Completion details.
  topicsID: string[]

  // Programming question.
  programmingQuestionTopicList: IModuleTopic[]
  programmingQuestionSelectedTopicList: IModuleTopic[]
  programmigQuestionTopicSelectionForm: FormGroup
  @ViewChild('topicSelection') topicSelection: any

  constructor(
    public utilService: UtilityService,
    private route: ActivatedRoute,
    private spinnerService: SpinnerService,
    private generalService: GeneralService,
    private datePipe: DatePipe,
    private role: Role,
    private urlConstant: UrlConstant,
    private accessLevel: AccessLevel,
    private localService: LocalService,
    private batchSessionService: BatchSessionService,
    private modalService: NgbModal,
    private router: Router,
    private batchService: BatchService,
    private formBuilder: FormBuilder,
    private moduleService: ModuleService,
    private btaService: BatchTopicAssignmentService
  ) {
    this.initializeVariables()
  }

  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize global variables.
  initializeVariables(): void {

    // Components.
    this.dayList = []

    // Batch details.
    this.batchDetailsList = [
      {
        fieldName: "Modules",
        fieldValue: 0,
        fieldImage: "assets/course/Modules.png"
      },
      {
        fieldName: "Topics",
        fieldValue: 0,
        fieldImage: "assets/course/topics.png"
      },
      {
        fieldName: "Assignments",
        fieldValue: 0,
        fieldImage: "assets/course/assignment.png"
      },
      {
        fieldName: "Sessions",
        fieldValue: 0,
        fieldImage: "assets/icon/colour-icons/sessions2.png"
      },
      {
        fieldName: "Hours",
        fieldValue: 0,
        fieldImage: "assets/icon/colour-icons/clock.png"
      },
      {
        fieldName: "Projects",
        fieldValue: 0,
        fieldImage: "assets/icon/colour-icons/project.png"
      },
    ]

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Batch Session Plan"

    // Session plan.
    this.batchSessionList = []
    this.dateList = []
    this.currentDateListIndex = 0
    this.currentDateListDateIndex = 0
    this.currentBatchSessionDate = null
    this.toDayDateInDateForm = new Date()
    this.currentBatchSessionPlanCurrentSubtopicCount = 0
    this.currentBatchSessionPlanPendingSubtopicCount = 0
    this.currentBatchSessionPlanCurrentAssignmentCountForFaculty = 0
    this.currentBatchSessionPlanCurrentAssignmentCountForTalent
    this.currentBatchSessionPlanPendingAssignmentCountForFaculty = 0
    this.currentBatchSessionPlanPendingAssignmentCountForTalent = 0

    // Flags.
    this.isCurrentSessionToday = false
    this.isSessionTaken = false
    this.isTalentPresent = false
    this.isUpdatePreRequisiteMode = false
    this.showPrerequisiteForm = false

    // Role.
    if (this.role.TALENT == this.localService.getJsonValue("roleName")) {
      this.access = this.accessLevel.ONLY_TALENT
    }

    if (this.role.FACULTY == this.localService.getJsonValue("roleName")) {
      this.access = this.accessLevel.ONLY_FACULTY
    }
    if (this.role.ADMIN == this.localService.getJsonValue("roleName") || this.role.SALES_PERSON == this.localService.getJsonValue("roleName")) {
      this.access = this.accessLevel.ADMIN_AND_SALESPERSON
    }

    // Access.
    if (this.access.isTalent) {
      this.permission = this.utilService.getPermission(this.urlConstant.MY_BATCHES)
    }
    if (!this.access.isTalent) {
      this.permission = this.utilService.getPermission(this.urlConstant.MY_BATCH_SESSION_DETAILS)
    }

    this.loginID = this.localService.getJsonValue("loginID")

    // Talent feedback for faculty.
    this.talentFeedbackForFacultyList = []

    // Batch session talent.
    this.batchSessionTalentList = []

    // Talent.
    this.talentID = this.localService.getJsonValue("loginID")

    // Cke editor configuration.
    this.ckeEditorCongiguration()

    // Completion details.
    this.topicsID = []

    // Programming question.
    this.programmingQuestionTopicList = []
    this.programmingQuestionSelectedTopicList = []

    // Batch.
    if (this.access.isTalent) {
      this.batchRelatedDetailsCount = 3
    }
    if (this.access.isFaculty) {
      this.batchRelatedDetailsCount = 1
    }
    this.route.queryParamMap.subscribe(params => {
      this.batchID = params.get("batchID")
      this.batchName = params.get("batchName")
      if (this.batchID) {
        this.getBatchSessionListAndDayList()
      }
    }, err => {
      console.error(err)
    })
  }

  // Format fields of batch session list. 
  formatBatchSessionList(): void {
    this.getYearMonthDateFromBatchSessionList()
    this.getCurrentBatchSessionDate()
    this.getCurrentSessionPlanRelatedDetails()
  }

  // Get years. months and dates form batch session list. 
  getYearMonthDateFromBatchSessionList(): void {

    this.dateList = []
    let isDateFound: boolean = false
    this.currentBatchSessionDate = null
    let todayDateInString: string = this.datePipe.transform(this.toDayDateInDateForm, 'yyyy-MM-dd')

    // Iterate batch session list.
    for (let i = 0; i < this.batchSessionList.length; i++) {
      let newDate: Date = new Date(this.batchSessionList[i].date)
      isDateFound = false

      // Set the last batch session with isSessionTaken as false to be the current batch session.
      if (this.currentBatchSessionDate == null && this.batchSessionList[i].isSessionTaken == false) {
        this.currentBatchSessionDate = this.batchSessionList[i].date
        this.actualBatchSessionIndex = i
        this.actualBatchSessionDate = this.batchSessionList[i].date

        // If batch session date is today's date then set isCurrentSessionToday to true.
        if (this.batchSessionList[i].date == todayDateInString) {
          this.isCurrentSessionToday = true
        }
      }

      // // Set the current batch session.
      // if (this.currentBatchSessionDate == null){

      //   // Calculate difference between todays date and batch session date.
      //   let currentBatchSessionDateInDateForm: Date = new Date(this.batchSessionList[i].date)
      //   let dateDifference: number = Math.floor((Date.UTC(currentBatchSessionDateInDateForm.getFullYear(), 
      //     currentBatchSessionDateInDateForm.getMonth(), currentBatchSessionDateInDateForm.getDate()) - 
      //     Date.UTC(this.toDayDateInDateForm.getFullYear(), this.toDayDateInDateForm.getMonth(), 
      //     this.toDayDateInDateForm.getDate()) ) /(1000 * 60 * 60 * 24))

      //   // If batch session date is todays date then show the batch session of today's date.
      //   if (this.currentBatchSessionDate == null && this.batchDetailsList[i].isSessionTaken == false){
      //     this.isCurrentSessionToday = true
      //     this.currentBatchSessionDate = this.batchSessionList[i].date
      //     this.actualBatchSessionIndex = i
      //   }

      //   // If batch session date is after todays date then show the next batch session.
      //   if (dateDifference > 0){
      //     this.isCurrentSessionToday = false
      //     this.currentBatchSessionDate = this.batchSessionList[i].date
      //     this.actualBatchSessionIndex = i
      //   }
      // }

      // Push the first date of batch session in the date list. 
      if (i == 0) {
        let dates: any[] = []
        dates.push(
          {
            date: newDate.getDate(),
            day: newDate.getDay(),
            isSessionTaken: this.batchSessionList[i].isSessionTaken,
            fromTime: this.utilService.convertTime(this.batchSessionList[i].moduleTiming.fromTime),
            toTime: this.utilService.convertTime(this.batchSessionList[i].moduleTiming.toTime)
          })
        this.dateList.push({
          year: newDate.getFullYear(),
          month: newDate.getMonth(),
          dates: dates
        })
        continue
      }

      // Iterate date list.
      for (let j = 0; j < this.dateList.length; j++) {

        // If year and month exists in date list then add only the date to the year and month entry.
        if (this.dateList[j].year == newDate.getFullYear() && this.dateList[j].month == newDate.getMonth()) {
          isDateFound = true
          this.dateList[j].dates.push(
            {
              date: newDate.getDate(),
              day: newDate.getDay(),
              isSessionTaken: this.batchSessionList[i].isSessionTaken,
              fromTime: this.utilService.convertTime(this.batchSessionList[i].moduleTiming.fromTime),
              toTime: this.utilService.convertTime(this.batchSessionList[i].moduleTiming.toTime)
            })
          break
        }
      }

      // If year and month combo does not exist in date list then add new entry in date list.
      if (!isDateFound) {
        let dates: any[] = []
        dates.push(
          {
            date: newDate.getDate(),
            day: newDate.getDay(),
            isSessionTaken: this.batchSessionList[i].isSessionTaken,
            fromTime: this.utilService.convertTime(this.batchSessionList[i].moduleTiming.fromTime),
            toTime: this.utilService.convertTime(this.batchSessionList[i].moduleTiming.toTime)
          })
        this.dateList.push({
          year: newDate.getFullYear(),
          month: newDate.getMonth(),
          dates: dates
        })
      }
    }

    // If current session date is not found then show the last batch session.
    if (this.currentBatchSessionDate == null) {
      this.currentBatchSessionDate = this.batchSessionList[this.batchSessionList.length - 1].date
    }
  }

  // Get the current date from date list.
  getCurrentBatchSessionDate(): void {
    let currentBatchSessionInDateForm: Date = new Date(this.currentBatchSessionDate)

    // Iterate date list.
    for (let i = 0; i < this.dateList.length; i++) {
      if (currentBatchSessionInDateForm.getFullYear() == this.dateList[i].year && currentBatchSessionInDateForm.getMonth() == this.dateList[i].month) {
        this.currentDateListIndex = i

        // Iterate date list dates.
        for (let j = 0; j < this.dateList[i].dates.length; j++) {
          if (this.dateList[i].dates[j].date == currentBatchSessionInDateForm.getDate()) {
            this.currentDateListDateIndex = j
          }
        }
      }
    }
  }

  // Get current session plan related details.
  getCurrentSessionPlanRelatedDetails(): void {
    if (this.access.isTalent) {
      this.batchRelatedDetailsCount = 3
      this.getCurrentBatchSessionPlan()
      this.getBatchSessionTalentList()
      this.getTalentFeedbackForFacultyList()
      return
    }
    if (!this.access.isTalent) {
      this.batchRelatedDetailsCount = 1
      this.getCurrentBatchSessionPlan()
      return
    }
  }

  // On clicking month year icon.
  onMonthYearClick(isNext: boolean): void {

    // If next month year is clicked.
    if (isNext) {
      this.currentDateListIndex = this.currentDateListIndex + 1
      this.currentDateListDateIndex = 0
    }

    // If previous month year is clicked.
    if (!isNext) {
      this.currentDateListIndex = this.currentDateListIndex - 1
      this.currentDateListDateIndex = 0
    }

    this.currentBatchSessionDate = this.createDateString()
    this.getCurrentSessionPlanRelatedDetails()
  }

  // On clicking date tab.
  onDateTabClick(event: any): void {
    this.currentDateListDateIndex = event.index
    this.currentBatchSessionDate = this.createDateString()
    this.getCurrentSessionPlanRelatedDetails()
  }

  // Create date string from date list.
  createDateString(): string {
    let monthInString: string
    let dateInString: string

    // If month is single then add 0 to it.
    if ((this.dateList[this.currentDateListIndex].month + 1) < 10) {
      monthInString = "0" + (this.dateList[this.currentDateListIndex].month + 1)
    }
    else {
      monthInString = (this.dateList[this.currentDateListIndex].month + 1).toString()
    }

    // If date is single then add 0 to it.
    if (this.dateList[this.currentDateListIndex].dates[this.currentDateListDateIndex].date < 10) {
      dateInString = "0" + this.dateList[this.currentDateListIndex].dates[this.currentDateListDateIndex].date
    }
    else {
      dateInString = (this.dateList[this.currentDateListIndex].dates[this.currentDateListDateIndex]).date.toString()
    }

    // Create date string.
    let date: string = this.dateList[this.currentDateListIndex].year + "-" + monthInString + "-" + dateInString
    return date
  }

  // Format current batch session plan.
  formatCurrentBatchSessionPlan(): void {
    this.checkCurrentSessionDateIsTodayDate()
    this.setShowStatusButtonForCurrentBatchSessionPlan()
    this.formatCurrentModuleList()
    this.formatPendingModuleList()
    if (this.access.isTalent) {
      this.getPreRequisite()
    }
  }

  // Check if current batch session date is todays date.
  checkCurrentSessionDateIsTodayDate(): void {
    let todayDateInString: string = this.datePipe.transform(this.toDayDateInDateForm, 'yyyy-MM-dd')
    if (todayDateInString == this.currentBatchSessionDate) {
      this.isCurrentSessionToday = true
      return
    }
    this.isCurrentSessionToday = false
  }

  // Set the show status button for current session plan only for the first session with is session taken as false.
  setShowStatusButtonForCurrentBatchSessionPlan(): void {
    // let actualBatchSessionDateInDateForm: Date = new Date(this.actualBatchSessionDate)
    // let currentBatchSessionDateInDateForm: Date = new Date(this.currentBatchSessionDate)
    // let dateDifference: number = Math.floor((Date.UTC(currentBatchSessionDateInDateForm.getFullYear(),
    //   currentBatchSessionDateInDateForm.getMonth(), currentBatchSessionDateInDateForm.getDate()) -
    //   Date.UTC(actualBatchSessionDateInDateForm.getFullYear(), actualBatchSessionDateInDateForm.getMonth(),
    //     actualBatchSessionDateInDateForm.getDate())) / (1000 * 60 * 60 * 24))
    // if (dateDifference < 1) {
    //   this.currentBatchSessionPlan.showStatusButton = true
    // }
    // if (this.actualBatchSessionDate == this.currentBatchSessionDate){
    //   this.currentBatchSessionPlan.showStatusButton = true
    // }
    let todayDateInDateForm: Date = new Date()
    let currentBatchSessionDateInDateForm: Date = new Date(this.currentBatchSessionDate)
    let dateDifference: number = Math.floor((Date.UTC(currentBatchSessionDateInDateForm.getFullYear(),
      currentBatchSessionDateInDateForm.getMonth(), currentBatchSessionDateInDateForm.getDate()) -
      Date.UTC(todayDateInDateForm.getFullYear(), todayDateInDateForm.getMonth(),
        todayDateInDateForm.getDate())) / (1000 * 60 * 60 * 24))
    if (dateDifference <= 1) {
      this.currentBatchSessionPlan.showStatusButton = true
    }
  }

  moduleList: IModule[] = []
  // Format the current modules list.
  formatCurrentModuleList(): void {
    this.currentBatchSessionPlanCurrentSubtopicCount = 0
    this.currentBatchSessionPlanCurrentAssignmentCountForFaculty = 0
    this.currentBatchSessionPlanCurrentAssignmentCountForTalent = 0

    // If current modules are present.
    if (this.currentBatchSessionPlan.module?.length > 0) {

      // Iterate current modules.
      for (let i = 0; i < this.currentBatchSessionPlan.module?.length; i++) {

        let module: any = this.currentBatchSessionPlan.module[i]
        let moduleSubTopicsCount: number = 0
        let moduleAssignmentsCount: number = 0
        this.moduleList.push(module)
        // Iterate module topics.
        for (let j = 0; j < module.moduleTopics?.length; j++) {
          let topic: any = module.moduleTopics[j]
          if (!this.topicsID.includes(topic?.id)) {
            this.topicsID.push(topic?.id)
          }
          this.programmingQuestionTopicList.push(topic)
          let topicAssignments: any[] = []

          // Iterate topic programming questions.
          for (let l = 0; l < topic.batchTopicAssignment?.length; l++) {

            // Count the number of topic programming questions for current batch session plan for faculty.
            this.currentBatchSessionPlanCurrentAssignmentCountForFaculty = this.currentBatchSessionPlanCurrentAssignmentCountForFaculty + 1

            // If batch topic is assigned then count the number of topic programming questions for current batch session plan for talent.
            if (topic.batchTopicAssignment[l].assignedDate != null) {
              this.currentBatchSessionPlanCurrentAssignmentCountForTalent = this.currentBatchSessionPlanCurrentAssignmentCountForTalent + 1
            }

            // If talent then push only those topic assignments that have assigned date.
            if (this.access.isTalent && topic.batchTopicAssignment[l].assignedDate != null) {
              topicAssignments.push(topic.batchTopicAssignment[l])
            }

            // If not talent then push all topic programming questions in topicAssignments.
            if (!this.access.isTalent) {
              topicAssignments.push(topic.batchTopicAssignment[l])
            }

            // // Calculate due date of .
            // let currentBatchSessionDateInDateForm: Date = new Date(this.currentBatchSessionDate)
            // currentBatchSessionDateInDateForm.setDate(currentBatchSessionDateInDateForm.getDate() + subTopic.batchTopicAssignment[l].days)
            // subTopic.batchTopicAssignment[l].dueDate = currentBatchSessionDateInDateForm

            // Give module name and topic name to assignment.
            topic.batchTopicAssignment[l].moduleName = module.moduleName
            topic.batchTopicAssignment[l].topicName = topic.topicName

            // Give difficulty to assignment.
            if (topic.batchTopicAssignment[l]?.programmingQuestion?.level == 1) {
              topic.batchTopicAssignment[l].programmingQuestion.difficulty = "Easy"
            }
            if (topic.batchTopicAssignment[l]?.programmingQuestion?.level == 2) {
              topic.batchTopicAssignment[l].programmingQuestion.difficulty = "Medium"
            }
            if (topic.batchTopicAssignment[l]?.programmingQuestion?.level == 3) {
              topic.batchTopicAssignment[l].programmingQuestion.difficulty = "Hard"
            }

            // Increment the module assignment count.
            moduleAssignmentsCount = moduleAssignmentsCount + 1
          }

          // Iterate sub topics.
          for (let k = 0; k < topic.subTopics?.length; k++) {
            let subTopic: any = topic.subTopics[k]

            // If batch session topic is completed.
            if (this.currentBatchSessionPlan.date == this.actualBatchSessionDate && subTopic.batchSessionTopic.isCompleted == true) {
              return
            }

            // Count the number of sub topics for current batch session plan.
            this.currentBatchSessionPlanCurrentSubtopicCount = this.currentBatchSessionPlanCurrentSubtopicCount + 1

            // Increment the sub topic count.
            moduleSubTopicsCount = moduleSubTopicsCount + 1
          }

          // Give all topicAssignments to topic.
          topic.topicAssignments = topicAssignments
        }

        // Give moduleSubTopicsCount to module.
        module.moduleSubTopicsCount = moduleSubTopicsCount

        // Give moduleAssignmentsCount to module.
        module.moduleAssignmentsCount = moduleAssignmentsCount
      }
    }
  }

  // Format the pending modules list.
  formatPendingModuleList(): void {
    this.currentBatchSessionPlanPendingSubtopicCount = 0
    this.currentBatchSessionPlanPendingAssignmentCountForFaculty = 0
    this.currentBatchSessionPlanPendingAssignmentCountForTalent = 0
    let currentBatchSessionDateInDateForm: Date = new Date(this.currentBatchSessionDate)

    // Get diffrence between todays date and current batch session date.
    let dateDifference: number = Math.floor((Date.UTC(currentBatchSessionDateInDateForm.getFullYear(),
      currentBatchSessionDateInDateForm.getMonth(), currentBatchSessionDateInDateForm.getDate()) -
      Date.UTC(this.toDayDateInDateForm.getFullYear(), this.toDayDateInDateForm.getMonth(),
        this.toDayDateInDateForm.getDate())) / (1000 * 60 * 60 * 24))

    // If current batch session date is after todays date then dont show pending assignments and sub topics.
    if (dateDifference > 0 && this.currentBatchSessionDate != this.actualBatchSessionDate) {
      return
    }

    // If pending modules are present.
    if (this.currentBatchSessionPlan.pendingModule?.length > 0) {

      // Iterate pending modules.
      for (let i = 0; i < this.currentBatchSessionPlan.pendingModule?.length; i++) {
        let module: any = this.currentBatchSessionPlan.pendingModule[i]
        let moduleSubTopicsCount: number = 0
        let moduleAssignmentsCount: number = 0
        this.moduleList.push(module)
        // Iterate module topics.
        for (let j = 0; j < module.moduleTopics?.length; j++) {
          let topic: any = module.moduleTopics[j]
          if (!this.topicsID.includes(topic?.id)) {
            this.topicsID.push(topic?.id)
          }
          this.programmingQuestionTopicList.push(topic)
          let topicAssignments: any[] = []

          // Iterate topic programming questions.
          for (let l = 0; l < topic.batchTopicAssignment?.length; l++) {

            // Count the  number of topic programming questions for current batch session plan.
            this.currentBatchSessionPlanPendingAssignmentCountForFaculty = this.currentBatchSessionPlanPendingAssignmentCountForFaculty + 1

            // If batch topic is assigned then count the number of topic programming questions for current batch session plan for talent.
            if (this.access.isTalent && topic.batchTopicAssignment[l].assignedDate != null) {
              this.currentBatchSessionPlanPendingAssignmentCountForTalent = this.currentBatchSessionPlanPendingAssignmentCountForTalent + 1
            }

            // If talent then push only those topic assignments that have assigned date.
            if (this.access.isTalent && topic.batchTopicAssignment[l].assignedDate != null) {
              topicAssignments.push(topic.batchTopicAssignment[l])
            }

            // If not talent then push all topic programming questions in topicAssignments.
            if (!this.access.isTalent) {
              topicAssignments.push(topic.batchTopicAssignment[l])
            }

            // Calculate due date of .
            // let currentBatchSessionDateInDateForm: Date = new Date(this.currentBatchSessionDate)
            // currentBatchSessionDateInDateForm.setDate(currentBatchSessionDateInDateForm.getDate() + subTopic.batchTopicAssignment[l].days)
            // subTopic.batchTopicAssignment[l].dueDate = currentBatchSessionDateInDateForm

            // Give module name and topic name to assignment.
            topic.batchTopicAssignment[l].moduleName = module.moduleName
            topic.batchTopicAssignment[l].topicName = topic.topicName

            // Give difficulty to assignment.
            if (topic.batchTopicAssignment[l]?.programmingQuestion?.level == 1) {
              topic.batchTopicAssignment[l].programmingQuestion.difficulty = "Easy"
            }
            if (topic.batchTopicAssignment[l]?.programmingQuestion?.level == 2) {
              topic.batchTopicAssignment[l].programmingQuestion.difficulty = "Medium"
            }
            if (topic.batchTopicAssignment[l]?.programmingQuestion?.level == 3) {
              topic.batchTopicAssignment[l].programmingQuestion.difficulty = "Hard"
            }

            // Increment the module assignment count.
            moduleAssignmentsCount = moduleAssignmentsCount + 1
          }

          // Iterate sub topics.
          for (let k = 0; k < topic.subTopics?.length; k++) {
            let subTopic: any = topic.subTopics[k]

            // If batch session topic is completed.
            if (this.currentBatchSessionPlan.date == this.actualBatchSessionDate && subTopic.batchSessionTopic.isCompleted == true) {
              continue
            }

            // Count the  number of sub topics for pending batch session plan.
            this.currentBatchSessionPlanPendingSubtopicCount = this.currentBatchSessionPlanPendingSubtopicCount + 1

            // Increment the sub topic count.
            moduleSubTopicsCount = moduleSubTopicsCount + 1
          }

          // Give all topicAssignments to topic.
          topic.topicAssignments = topicAssignments
        }

        // Give moduleSubTopicsCount to module.
        module.moduleSubTopicsCount = moduleSubTopicsCount

        // Give moduleAssignmentsCount to module.
        module.moduleAssignmentsCount = moduleAssignmentsCount
      }
    }
  }

  // On clicking view assignment button.
  onViewAssignment(assignment: any): void {
    this.selectedAssignment = assignment
    this.openModal(this.programmingQuestionModal, 'lg')
  }

  // Format batch details list
  formatBatchDetailsList(): void {
    this.batchDetailsList[0].fieldValue = this.batchSessionsCounts.moduleCount
    this.batchDetailsList[1].fieldValue = this.batchSessionsCounts.topicCount
    this.batchDetailsList[2].fieldValue = this.batchSessionsCounts.assignmentCount
    this.batchDetailsList[3].fieldValue = this.batchSessionsCounts.sessionCount
    this.batchDetailsList[4].fieldValue = Math.floor(this.batchSessionsCounts.totalBatchHours / 60)
    this.batchDetailsList[5].fieldValue = this.batchSessionsCounts.projectCount
  }

  // Used to open modal.
  openModal(content: any, size?: string, options: NgbModalOptions = {
    ariaLabelledBy: 'modal-basic-title', keyboard: false,
    backdrop: 'static', size: size
  }): NgbModalRef {
    if (!size) {
      options.size = 'lg'
    }
    this.modalRef = this.modalService.open(content, options)
    return this.modalRef
  }


  dismissModal(reason?: string) {
    this.modalRef.dismiss(reason)
  }

  // On clicking skip session button.
  onSkipSessionClick(): void {
    this.modalRef.close()
    this.spinnerService.loadingMessage = "Updating session plan"


    let queryParams: any
    if (this.access.isFaculty) {
      queryParams = {
        facultyID: this.loginID
      }
    }
    this.batchSessionService.skipPendingSession(this.batchID, queryParams).subscribe((response: any) => {
      alert("Session plan successfully updated")
      this.getBatchSessionList()
    }, (err: any) => {
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.statusText)
    })
  }


  // =================Niranjan====================



  // ====================== Niranjan end ============

  // On clicking any of the bacth session buttons.
  onBatchSessionButtonClick(operation: string): void {
    if (this.modalRef) {
      this.modalRef.close()
    }

    // Go to next date from current batch session date.
    let nextDate: string = null

    // If batch session is not the last then give next batch session's date.
    if (this.actualBatchSessionIndex < (this.batchSessionList.length - 1)) {
      nextDate = this.batchSessionList[this.actualBatchSessionIndex + 1].date
    }

    let queryParamMap = {
      operation: operation,
      batchID: this.batchID,
      batchSessionID: this.currentBatchSessionPlan.id,
      currentDate: this.currentBatchSessionDate,
      nextDate: nextDate,
      batchName: this.batchName,
      topicsID: JSON.stringify(this.topicsID),
    }

    let url: string
    if (this.access.isFaculty) {
      url = "/my/batch/session/completion"
    }
    if (this.access.isAdmin || this.access.isSalesPerson) {
      url = "/training/batch/master/session/completion"
    }
    this.router.navigate([url], {
      relativeTo: this.route,
      queryParams: queryParamMap,
    })
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

    // Give attendance of talent to all batch session date list.
    for (let i = 0; i < this.batchSessionList.length; i++) {

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

      this.batchSessionList[i].isPresent = false
      this.batchSessionList[i].isFeedbackGiven = false
      this.batchSessionList[i].sessionNumber = (i + 1)
      this.batchSessionList[i].feedbacks = []
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

    // Iterate batch session date list.
    for (let i = 0; i < this.batchSessionList.length; i++) {

      // If talent was present for that batch session.
      if (this.batchSessionList[i].isPresent) {

        // Iterate talent feedback for faculty list.
        for (let j = 0; j < this.talentFeedbackForFacultyList.length; j++) {

          if (this.batchSessionList[i].id == this.talentFeedbackForFacultyList[j].batchSessionID) {

            // Give current batch session feedbacks.
            this.batchSessionList[i].feedbacks = this.talentFeedbackForFacultyList[j].sessionFeedbacks

            // If talent feedback for faculty is present for that batch session then set isFeedbackGiven to true.
            if (this.talentFeedbackForFacultyList[j].sessionFeedbacks?.length > 0) {
              this.batchSessionList[i].isFeedbackGiven = true
            }
          }
        }
      }
    }

    // Set the current batch session by matching date with the currnt batch session date.
    for (let i = 0; i < this.batchSessionList.length; i++) {
      if (this.batchSessionList[i].date == this.currentBatchSessionDate) {
        this.currentBatchSession = this.batchSessionList[i]
      }
    }
  }

  //********************************************* TALENT FEEDACK FOR FACULTY FUNCTIONS ************************************************************

  // On clicking talent feedback to faculty.
  onTalentFeeddbackToFacultyClick(): void {
    this.selectedBatchSessionForTalentFeedbackSessionNumber = this.currentBatchSession.sessionNumber
    for (let i = 0; i < this.batchSessionList.length; i++) {
      if (this.batchSessionList[i].date == this.currentBatchSessionDate) {
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
      this.currentBatchSession.isFeedbackGiven = true
      this.getOneTalentOneBatchSessionFeedback()
      this.isUpdatedEmitter.emit()
    }
  }

  //********************************************* PROGRAMMING QUESTION RELATED FUNCTIONS ************************************************************

  // Create topic selection form.
  createTopicSelectionForm(): void {
    this.programmigQuestionTopicSelectionForm = this.formBuilder.group({
      moduleID: new FormControl(this.moduleList[0]?.id, [Validators.required]),
      // selectedTopic or topicList? possible bug with multiple modules(test) #niranjan 
      topicID: new FormControl(this.programmingQuestionTopicList[0]?.id, [Validators.required]),
      isPublished: new FormControl(true),
      dueDate: new FormControl(null)
    })
  }

  // On changing module selection.
  changeModule(e: any): void {
    this.programmingQuestionSelectedTopicList = e.target.moduleTopics
  }

  // On clicking add a new assignment button.
  onAddNewAssignmentClick() {
    this.createTopicSelectionForm()
    this.openModal(this.topicSelection, 'lg')
  }

  // Check if topic selection form is valid or not.
  isTopicSelectionFormValid(): boolean {
    if (this.programmigQuestionTopicSelectionForm.get('isPublished').value) {
      this.programmigQuestionTopicSelectionForm.get('dueDate').setValidators([Validators.required])
    } else {
      this.programmigQuestionTopicSelectionForm.get('dueDate').clearValidators()
    }
    this.utilService.updateValueAndValiditors(this.programmigQuestionTopicSelectionForm)
    if (this.programmigQuestionTopicSelectionForm.invalid) {
      this.programmigQuestionTopicSelectionForm.markAllAsTouched()
      return false
    }
    return true
  }

  // Add programming question to question bank.
  addProgrammingQuestion(topicID: string): Promise<string> {
    return new Promise<string>((resolve, reject) => {
      this.openModal(ProgrammingQuestionModalComponent, 'xl')
      let componentInstance = <ProgrammingQuestionModalComponent>this.modalRef.componentInstance
      componentInstance.topicID = topicID
      componentInstance.dismissModalEvent.subscribe(() => {
        this.modalRef.dismiss()
      })
      componentInstance.addEvent.subscribe((id: string) => {
        resolve(id)
        alert("Question added to bank")
      }, () => {
        reject("Question not added")
      })
    })
  }

  // Add programming qusetion to module topic.
  addTopicProgrammingQuestion(topicID: string, topicProgrammingQuesiton: ITopicProgrammingQuestion):
    Promise<void> {
    return new Promise<void>((resolve, reject) => {
      this.moduleService.addTopicProgrammingQuestion(topicID, topicProgrammingQuesiton).subscribe(() => {
        resolve()
        alert("Question added to topic")
      }, () => {
        reject("Question not added to topic.")
      })
    })
  }

  // Add programming question to batch topic.
  addBatchTopicAssignment(topicID: string, bta: IBatchTopicAssignment):
    Promise<void> {
    return new Promise<void>((resolve, reject) => {
      this.btaService.addBatchTopicAssignment(this.batchID, topicID, bta).subscribe(() => {
        resolve()
      }, () => {
        reject("Question not assigned to batch.")
      })
    })
  }

  // On clicking add new question button. 
  async onAddQuestionClick() {
    if (!this.isTopicSelectionFormValid()) {
      return
    }
    this.modalRef.close()
    // // to get from query params
    // let params = this.route.snapshot.queryParamMap

    let moduleID = this.programmigQuestionTopicSelectionForm.get('moduleID').value
    let topicID = this.programmigQuestionTopicSelectionForm.get('topicID').value
    let batchSessionID = this.currentBatchSessionPlan.id
    let dueDate = this.programmigQuestionTopicSelectionForm.get('dueDate').value

    try {
      let questionID = await this.addProgrammingQuestion(topicID)

      let topicProgrammingQuesiton: ITopicProgrammingQuestion = {
        isActive: true,
        programmingQuestionID: questionID
      }
      await this.addTopicProgrammingQuestion(topicID, topicProgrammingQuesiton)
      let bta: IBatchTopicAssignment = {
        batchID: this.batchID,
        programmingQuestionID: questionID,
        moduleID: moduleID,
        topicID: topicID,
        dueDate: dueDate,
        batchSessionID: batchSessionID,
      }
      console.log({ bta })
      await this.addBatchTopicAssignment(topicID, bta)
      alert("Operation complete.")
      this.modalRef.close()

    } catch (err) {
      alert("Some error occurred")
      console.error("catch err (p:undefined):", err)
    }
  }

  // Redirect to add new assignment tab in batch details page.
  redirectToAssignNewAssignment() {
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: {
        "tab": "3",
        "subTab": "Assign New"
      },
      queryParamsHandling: 'merge'
    }).catch((err: any) => {
      console.error(err)
    })
  }

  //********************************************* PREREQUISITES FUNCTIONS ************************************************************

  // On clicking add pre requisite button.
  onPreRequisiteClick(): void {
    this.showPrerequisiteForm = false
    this.isUpdatePreRequisiteMode = false
    this.openModal(this.preRequisiteModal, "lg")
    this.getPreRequisite()
  }

  // Create pre requisite form.
  createPreRequisiteForm(): void {
    this.preRequisiteForm = this.formBuilder.group({
      id: new FormControl(null),
      prerequisite: new FormControl(null, [Validators.required, Validators.maxLength(1000)]),
    })
  }

  // On clicking add new pre requisite button.
  onAddNewPreRequisiteButtonClick(): void {
    this.isUpdatePreRequisiteMode = false
    this.createPreRequisiteForm()
    this.showPrerequisiteForm = true
  }

  // On clicking update new pre requisite button.
  onUpdatePreRequisiteButtonClick(): void {
    this.isUpdatePreRequisiteMode = true
    this.createPreRequisiteForm()
    this.preRequisiteForm.get('prerequisite').setValue(this.currentBatchSessionPlan.prerequisiteList[0].prerequisite)
    this.showPrerequisiteForm = true
  }

  // Add pre requisite.
  addPreRequisite(): void {
    this.spinnerService.loadingMessage = "Adding Pre-Requisite"
    this.batchSessionService.addPreRequisite(this.batchID, this.currentBatchSessionPlan.id,
      this.preRequisiteForm.value).subscribe((response: any) => {
        this.showPrerequisiteForm = false
        this.getPreRequisite()
        alert(response.body)
      }, (error) => {
        if (error.error?.error) {
          alert(error.error?.error)
          return
        }
        alert(error.statusText)
      })
  }

  // Validate pre requisite form.
  validatePreRequisiteForm(): void {
    if (this.preRequisiteForm.invalid) {
      this.preRequisiteForm.markAllAsTouched()
      return
    }
    if (this.isUpdatePreRequisiteMode) {
      this.updatePreRequisite()
      return
    }
    this.addPreRequisite()
  }

  // Get pre requisite.
  getPreRequisite(): void {
    if (this.access.isTalent) {
      this.spinnerService.loadingMessage = "Getting Batch session Plan"
    }
    if (!this.access.isTalent) {
      this.spinnerService.loadingMessage = "Getting Pre-Requisite"
    }
    let quaryParams: any = {
      batchSessionID: this.currentBatchSessionPlan.id
    }
    this.batchSessionService.getPreRequisite(this.batchID, quaryParams).subscribe((response: any) => {
      this.currentBatchSessionPlan.prerequisiteList = response.body
    }, (error) => {
      if (error.error?.error) {
        alert(error.error?.error)
        return
      }
      alert(error.statusText)
    })
  }

  // Update pre requisite
  updatePreRequisite(): void {
    this.spinnerService.loadingMessage = "Updating Pre-Requisite"
    this.batchSessionService.updatePreRequisite(this.batchID, this.currentBatchSessionPlan.id,
      this.currentBatchSessionPlan.prerequisiteList[0].id, this.preRequisiteForm.value).subscribe((response: any) => {
        this.showPrerequisiteForm = false
        this.getPreRequisite()
        alert(response.body)
      }, (error) => {
        if (error.error?.error) {
          alert(error.error?.error)
          return
        }
        alert(error.statusText)
      })
  }

  // On clicking delete pre requisite.
  onDeletePreRequisite(): void {
    if (confirm("Are you sure you want to delete the pre requisite ?")) {
      this.deletePreRequisite()
    }
  }

  // Delete pre requisite
  deletePreRequisite(): void {
    this.spinnerService.loadingMessage = "Delete Pre-Requisite"
    this.batchSessionService.deletePreRequisite(this.batchID, this.currentBatchSessionPlan.id,
      this.currentBatchSessionPlan.prerequisiteList[0].id).subscribe((response: any) => {
        this.showPrerequisiteForm = false
        this.getPreRequisite()
        alert(response.body)
      }, (error) => {
        if (error.error?.error) {
          alert(error.error?.error)
          return
        }
        alert(error.statusText)
      })
  }

  // cke editor congiguration.
  ckeEditorCongiguration(): void {
    this.ckeEditorConfig = {
      extraPlugins: 'codeTag,kbdTag',
      removePlugins: "exportpdf",
      // stylesSet: 'new_styles',
      toolbar: [
        { name: 'styles', items: ['Styles', 'Format'] },
        {
          name: 'basicstyles', groups: ['basicstyles', 'cleanup'],
          items: ['Bold', 'Italic', 'Underline', 'Strike', 'Subscript', 'Superscript', 'RemoveFormat', 'Code']
        },
        {
          name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'],
          items: ['NumberedList', 'BulletedList', '-', 'Outdent', 'Indent', '-', 'Blockquote']
        },
        // { name: 'links', items: [ 'Link', 'Unlink' ] }, //, 'Anchor' // Link
        // { name: 'insert', items: [ 'Image' ] }, //, 'Table', 'HorizontalRule' // Image
        { name: 'document', groups: ['mode', 'document', 'doctools'], items: ['Source'] },
        // { name: 'clipboard', groups: [ 'clipboard', 'undo' ], items: [ 'Cut', 'Copy', 'Paste', 'PasteText', 'PasteFromWord', '-', 'Undo', 'Redo' ] },
        // { name: 'editing', groups: [ 'find', 'selection', 'spellchecker' ], items: [ 'Scayt' ] },
        // { name: 'tools', items: [ 'Maximize' ] },
        // { name: 'others', items: [ '-' ] },
        // { name: 'about', items: [ 'About' ] }
      ],
      toolbarGroups: [
        { name: 'styles' },
        { name: 'basicstyles', groups: ['basicstyles', 'cleanup'] },
        { name: 'document', groups: ['mode', 'document', 'doctools'] },
        { name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'] },
        // { name: 'links' }, // Link
        // { name: 'insert' }, // Image
        // '/',
        // { name: 'colors' },
        // { name: 'clipboard', groups: [ 'clipboard', 'undo' ] },
        // { name: 'editing', groups: [ 'find', 'selection', 'spellchecker' ] },
        // { name: 'forms' },
        // { name: 'tools' },
        // { name: 'others' },
      ],
      removeButtons: "",
      language: 'en',
      resize_enabled: false,
      width: "100%", height: "150px",
      forcePasteAsPlainText: false,
    }
  }

  //********************************************* GET FUNCTIONS ************************************************************

  // Get batch session list and day list by calling asynchronously.
  async getBatchSessionListAndDayList(): Promise<void> {
    try {

      // Get day list.
      if (this.dayList.length == 0) {
        this.dayList = await this.getDayList()
      }

      // Get batch session list.
      this.batchSessionList = await this.getBatchSessionList()
      if (this.batchSessionList.length == 0) {
        this.currentBatchSessionPlan = null
        this.dateList = []
        return
      }
      this.getBatchSessionsCounts()
      this.formatBatchSessionList()
    } catch (error) {
      console.error(error)
    }
  }

  // Get day list.
  async getDayList(): Promise<any[]> {
    try {
      return await new Promise<any[]>((resolve, reject) => {
        this.spinnerService.loadingMessage = "Getting Batch session Plan"
        this.generalService.getDaysList().subscribe((response) => {
          resolve(response.body)
        }, (err: any) => {
          console.error(err)
          if (err.statusText.includes('Unknown')) {
            reject("No connection to server. Check internet.");
            return
          }
          reject(err.error.error)
        })
      })
    } finally {
    }
  }

  // Get batch session list.
  async getBatchSessionList(): Promise<any[]> {
    try {
      return await new Promise<any[]>((resolve, reject) => {
        this.spinnerService.loadingMessage = "Getting Batch session Plan"
        let queryParams: any
        if (this.access.isFaculty) {
          queryParams = {
            facultyID: this.loginID
          }
        }
        this.batchSessionService.getBatchSessionWithTopicNameList(this.batchID, queryParams).subscribe((response) => {
          resolve(response.body)
        }, (err: any) => {
          console.error(err)
          if (err.statusText.includes('Unknown')) {
            reject("No connection to server. Check internet.")
            return;
          }
          reject(err.error.error)
        })
      })
    } finally {
    }
  }

  // Get current batch session plan by date.
  getCurrentBatchSessionPlan(): void {
    this.spinnerService.loadingMessage = "Getting Batch session Plan"
    let queryParams: any = {
      sessionTopicCompletedDate: this.currentBatchSessionDate,
      sessionDate: this.currentBatchSessionDate,
      pendingCompletedDate: this.currentBatchSessionDate,
    }

    // If faculty login then show only their session plan.
    if (this.access.isFaculty) {
      queryParams.facultyID = this.loginID
    }
    this.batchSessionService.getBatchSessionPlan(this.batchID, queryParams).subscribe((response) => {
      this.currentBatchSessionPlan = response
      this.formatCurrentBatchSessionPlan()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.decrementBatchRelatedDetailsCount()
    })
  }

  // Get batch session counts.
  getBatchSessionsCounts(): void {
    this.spinnerService.loadingMessage = "Getting Batch session Plan"
    let queryParams: any = {
      facultyID: this.access.isFaculty ? this.loginID : null
    }
    this.batchSessionService.getBatchSessionsCounts(this.batchID, queryParams).subscribe((response) => {
      this.batchSessionsCounts = response
      this.formatBatchDetailsList()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get talent feedback for faculty list.
  getTalentFeedbackForFacultyList(): void {
    this.spinnerService.loadingMessage = "Getting Batch session Plan"
    this.batchService.getTalentBatchFeedback(this.batchID, this.talentID).subscribe((response) => {
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
    this.spinnerService.loadingMessage = "Getting Batch session Plan"
    this.batchService.getAllBatchSessionTalentsForTalent(this.batchID, this.talentID).subscribe((response) => {
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

  // Get talent feedback for one batch session.
  getOneTalentOneBatchSessionFeedback(): void {
    this.spinnerService.loadingMessage = "Getting Batch session Plan"
    this.batchService.getOneTalentOneBatchSessionFeedback(this.batchID, this.talentID, this.currentBatchSession.id).subscribe((response) => {
      this.currentBatchSession.feedbacks = response.body
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }


}
