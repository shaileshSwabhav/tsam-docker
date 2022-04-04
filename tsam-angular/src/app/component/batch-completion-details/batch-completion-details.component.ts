import { DatePipe } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, ParamMap, Params, Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BatchSessionService } from 'src/app/service/batch-session/batch-session.service';
import { BatchTalentService } from 'src/app/service/batch-talent/batch-talent.service';
import { BatchTopicAssignmentService } from 'src/app/service/batch-topic-assignment/batch-topic-assignment.service';
import { BatchService, IBatchSessionFeedback, IBatchSessionsTalent, IFacultyBatchSessionFeedback, IMappedSession } from 'src/app/service/batch/batch.service';
import { UrlConstant } from 'src/app/service/constant';
import { IFeedbackOptions } from 'src/app/service/feedback/feedback.service';
import { GeneralService, IFeedbackQuestion } from 'src/app/service/general/general.service';
import { LocalService } from 'src/app/service/storage/local.service';

@Component({
  selector: 'app-batch-completion-details',
  templateUrl: './batch-completion-details.component.html',
  styleUrls: ['./batch-completion-details.component.css']
})
export class BatchCompletionDetailsComponent implements OnInit {


  batchName: string
  courseID: string
  sessionDate: string
  sessionName: string
  sessionNumber: number
  completedModules: any = []
  nextSessionModules: any = []
  isNextClicked: boolean
  moduleName: string
  topicName: string
  isTommorrowSessionCompleted: boolean
  isAttendanceComponent: boolean

  // session
  totalSessions: number
  batchModules: any[]
  pendingModules: any[]
  isPendingSessions: boolean

  sessionMarkAsDone: any[]

  // batch-session talent attendance
  sessionTalentAttendance: IBatchSessionsTalent[]
  sessionTalentAttendanceForm: FormGroup

  // spinner


  // batch
  batch: any
  batchID: string
  // batchName: string
  // sessionID: string
  batchSessionID: string
  batchTalents: any[]
  totalBatchModules: number
  talentLimit: number
  talentOffset: number
  talents: any[]
  totalTalent: number
  talentCount: number
  assignments: any = []

  // session attendance
  talentAttendance: IBatchSessionsTalent[]

  // session feedback
  sessionFeedbackForm: FormGroup
  feedbackForm: FormGroup
  talentSessionFeedbackForm: FormGroup
  isSingleTalentFeedback: boolean
  isSessionFeedbackLoaded: boolean
  isTalentsLoaded: boolean

  // feedback
  // facultyFeedbackQuestions: IFeedbackQuestionGroup[]
  loginID: string
  facultyName: string
  isFeedback: boolean
  facultyFeedbackQuestions: IFeedbackQuestion[]
  facultyID: string
  talentID: string
  feedbackTo: string
  facultySessionFeedback: IFacultyBatchSessionFeedback
  feedbackResponse: any[]
  talentFeedbacks: IMappedSession[]
  talentSessionFeedback: IBatchSessionFeedback[]
  // talentFeedbackComments: IFeedbackComments[]
  talentFeedbackKeywords: string[]
  tabToBeOpen: string


  //assignment
  selectedAssignments: any[]

  topicsID: string[] = []

  // constants
  INITIAL_SCORE = 0
  MINIMUM_SCORE = 7
  MAX_SCORE = 10
  FEEDBACK_TALENT_INDEX = 0
  nextSessionDate: string

  starHovered: any

  public readonly ATTENDANCEOPERATION: string = "attendance"
  public readonly COMPLETEDOPERATION: string = "completed"
  public readonly FEEDBACKOPERATION: string = "feedback"
  public readonly ASSIGNMENTOPERATION: string = "assignment"

  constructor(
    private formBuilder: FormBuilder,
    private urlConstant: UrlConstant,
    private router: Router,
    private route: ActivatedRoute,
    private datePipe: DatePipe,
    private batchSessionService: BatchSessionService,
    private spinnerService: SpinnerService,
    private batchService: BatchService,
    private batchTalentService: BatchTalentService,
    private activatedRoute: ActivatedRoute,
    private generalService: GeneralService,
    private localService: LocalService,
    private batchTopicAssignmentService: BatchTopicAssignmentService,
  ) {
    this.extractID()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
    this.initializeVariables()
    this.getAllComponents()
  }


  initializeVariables(): void {

    this.loginID = this.localService.getJsonValue("loginID")
    this.facultyName = this.localService.getJsonValue("firstName")

    this.isNextClicked = false
    this.isTommorrowSessionCompleted = false
    this.isAttendanceComponent = false
    this.isPendingSessions = false
    // this.isFeedback = false

    this.pendingModules = []
    this.batchModules = []
    this.talents = []
    this.completedModules = []

    this.totalSessions = 0
    this.totalBatchModules = 0
    this.totalTalent = 0
    this.sessionNumber = 4

    this.talentCount = 0
    this.sessionMarkAsDone = []
    this.assignments = []
    this.isSessionFeedbackLoaded = false
    this.isTalentsLoaded = false
    this.talentCount = 0
    this.selectedAssignments = []

  }

  getAllComponents() {
    this.processTabs()
    // this.getAllPendingModules()
    // this.getTalentsForBatch()
    // this.createSessionTalentAttendanceForm()
  }

  extractID(): void {
    this.activatedRoute.queryParamMap.subscribe(
      (params: any) => {
        this.batchID = params.get("batchID")
        this.sessionDate = params.get("currentDate")
        this.batchName = params.get("batchName")
        this.courseID = params.get("courseID")
        this.batchSessionID = params.get("batchSessionID")
        this.tabToBeOpen = params.get("operation")
        this.nextSessionDate = params.get("nextDate")
        this.topicsID = JSON.parse(params.get("topicsID"))
      }, (err: any) => {
        console.error(err);
      })
  }

  async processTabs(): Promise<void> {
    if (this.tabToBeOpen == this.ATTENDANCEOPERATION) {
      // console.log(this.ATTENDANCEOPERATION);
      this.getTalentsForBatch()
      this.isNextClicked = false
      this.isAttendanceComponent = true
      return
    }
    if (this.tabToBeOpen == this.FEEDBACKOPERATION) {
      this.isFeedback = true
      // this.talentCount = 0
      // this.facultyFeedbackQuestions = await this.getFeedbackQuestionsForFaculty("Faculty_Session_Feedback")
      // this.batchTalents = await this.getPresentTalents()
      // console.log(this.batchTalents);

      // this.addPresentTalentsToAttendanceForm()
      // this.processPresentTalents()
      // console.log(this.FEEDBACKOPERATION);
      return
    }
    if (this.tabToBeOpen == this.COMPLETEDOPERATION) {
      // console.log(this.COMPLETEDOPERATION);
      this.isAttendanceComponent = false
      this.isFeedback = false
      this.getBatchSessions()
      if (this.nextSessionDate) {
        this.getNextSessions()
      }
      return
    }
    if (this.tabToBeOpen == this.ASSIGNMENTOPERATION) {
      // console.log(this.ASSIGNMENTOPERATION);
      console.log(this.topicsID);

      this.isAttendanceComponent = false
      this.isFeedback = false
      this.getModuleTopicsAssignment()
      return
    }
  }
  markAllPendingSessionAsDone(topic: any): void {
    for (let index = 0; index < topic.subTopics.length; index++) {
      const element = topic.subTopics[index];
      this.markSessionAsDone(topic.id, element.id)
    }

  }

  markSessionAsDone(topicID: string, subTopicID: string): void {
    let session = {
      completedDate: this.sessionDate,
      subTopicID: subTopicID,
      topicID: topicID,
      isCompleted: true

    }
    if (!this.isSubTopicmarked(subTopicID)) {
      this.sessionMarkAsDone.push(session)
      // console.log(this.sessionMarkAsDone);
      return
    }
    // let index = this.sessionMarkAsDone.indexOf(subTopicID)
    let index = this.getIndexOFSessionPresentID(subTopicID)
    this.sessionMarkAsDone.splice(index, 1)

  }
  isSubTopicmarked(subTopicID: string): boolean {
    for (let index = 0; index < this.sessionMarkAsDone.length; index++) {
      let id = this.sessionMarkAsDone[index].subTopicID
      if (id == subTopicID) {
        return true
      }
    }
    return false
  }

  isAllSubTopicMarked(topic: any): boolean {
    for (let index = 0; index < topic.subTopics.length; index++) {
      const element = topic.subTopics[index];
      if (!this.isSubTopicmarked(element.id)) {
        return false
      }
    }
    return true
  }

  getIndexOFSessionPresentID(subTopicID: any): number {
    for (let index = 0; index < this.sessionMarkAsDone.length; index++) {
      let id = this.sessionMarkAsDone[index].subTopicID
      // console.log(subTopicID);
      // console.log(id);


      if (id == subTopicID) {
        return index
      }
    }
    return -1
  }

  createSessionTalentAttendanceForm(): void {
    this.sessionTalentAttendanceForm = this.formBuilder.group({
      talentAttendance: new FormArray([])
    })
  }
  get talentAttendanceControlArray() {
    return this.sessionTalentAttendanceForm.get('talentAttendance') as FormArray
  }

  addTalentAttendanceToControlArray(): void {
    this.talentAttendanceControlArray.push(this.formBuilder.group({
      id: new FormControl(null),
      batchID: new FormControl(this.batchID),
      // batchSessionID: new FormControl(null, [Validators.required]),
      talentID: new FormControl(null, [Validators.required]),
      talent: new FormControl(null),
      isPresent: new FormControl(null, [Validators.required]),
      isFeedbackGiven: new FormControl(null),
      attendedDate: new FormControl(this.datePipe.transform(new Date(), "yyyy-MM-dd"), [Validators.required]),
      averageRating: new FormControl(null),
    }))
  }

  addTalentsToAttendanceForm(): void {
    this.createSessionTalentAttendanceForm()
    for (let index = 0; index < this.batchTalents?.length; index++) {
      this.addTalentAttendanceToControlArray()
      // console.log(this.searchBatchDetailsForm.value);
      // this.talentAttendanceControlArray.at(index).get("batchSessionID").setValue(this.searchBatchDetailsForm.get("sessionID").value)
      this.talentAttendanceControlArray.at(index).get("talentID").setValue(this.batchTalents[index].id)
      this.talentAttendanceControlArray.at(index).get("talent").setValue(this.batchTalents[index])
      this.talentAttendanceControlArray.at(index).get("isPresent").setValue(true)
      this.talentAttendanceControlArray.at(index).get("attendedDate").setValue(this.datePipe.transform(this.sessionDate, "yyyy-MM-dd"))
    }
    this.sessionTalentAttendanceForm.markAsDirty()
  }

  addPresentTalentsToAttendanceForm(): void {
    this.createSessionTalentAttendanceForm()
    for (let index = 0; index < this.batchTalents?.length; index++) {
      this.addTalentAttendanceToControlArray()
      // console.log(this.searchBatchDetailsForm.value);
      // this.talentAttendanceControlArray.at(index).get("batchSessionID").setValue(this.searchBatchDetailsForm.get("sessionID").value)
      console.log(this.talentAttendanceControlArray.at(index).get("isFeedbackGiven"));


      this.talentAttendanceControlArray.at(index).get("isFeedbackGiven").setValue(this.batchTalents[index].isFeedbackGiven)
      this.talentAttendanceControlArray.at(index).get("talentID").setValue(this.batchTalents[index].talent.id)
      this.talentAttendanceControlArray.at(index).get("talent").setValue(this.batchTalents[index].talent)
      this.talentAttendanceControlArray.at(index).get("isPresent").setValue(true)
      this.talentAttendanceControlArray.at(index).get("attendedDate").setValue(this.datePipe.transform(this.sessionDate, "yyyy-MM-dd"))
    }
    this.sessionTalentAttendanceForm.markAsDirty()
  }


  createFeedbackForm(): void {
    this.feedbackForm = this.formBuilder.group({
      feedbacks: new FormArray([])
    })
  }

  get feedbackArray(): FormArray {
    return this.feedbackForm.get("feedbacks") as FormArray
  }


  addFeedbackQuestionsToArray(): void {
    // this.sessionFeedbackForm = this.formBuilder.group({
    this.feedbackArray.push(this.formBuilder.group({
      batchID: new FormControl(this.batchID),
      talentID: new FormControl(null, [Validators.required]),
      facultyID: new FormControl(this.loginID),
      questionID: new FormControl(null, [Validators.required]),
      optionID: new FormControl(null, [Validators.required]),
      option: new FormControl(null, [Validators.required]),
      answer: new FormControl(null, [Validators.required]),
    }))
  }


  createFeedbackQuestionsForm(talentID: string): void {
    console.log("createFeedbackQuestionsForm -> ", this.facultyFeedbackQuestions);

    this.createFeedbackForm()

    if (this.facultyFeedbackQuestions) {
      for (let index = 0; index < this.facultyFeedbackQuestions.length; index++) {
        this.addFeedbackQuestionsToArray()

        // initialize form
        this.feedbackArray.at(index).get('facultyID').setValue(this.loginID)
        this.feedbackArray.at(index).get("talentID").setValue(talentID)

        this.feedbackArray.at(index).get("questionID").setValue(this.facultyFeedbackQuestions[index].id)

        if (this.facultyFeedbackQuestions[index].hasOptions) {
          this.feedbackArray.at(index).get("optionID").setValidators([Validators.required])
          this.feedbackArray.at(index).get("option").setValidators([Validators.required])
          this.feedbackArray.at(index).get("answer").setValidators(null)
        } else {
          this.feedbackArray.at(index).get("optionID").setValidators(null)
          this.feedbackArray.at(index).get("option").setValidators(null)
          this.feedbackArray.at(index).get("answer").setValidators([Validators.required])
        }
      }
      this.isFeedback = true
    }

    // console.log("createFeedbackQuestionsForm -> ", this.isFeedback);
  }

  onFacultyFeedbackChange(feedbackQuestionControl: any, feedbackOptions: IFeedbackOptions[]) {

    let optionID = feedbackQuestionControl.get("optionID").value
    for (let index = 0; index < feedbackOptions.length; index++) {
      if (optionID == feedbackOptions[index].id) {
        feedbackQuestionControl.get('answer').setValue(feedbackOptions[index].value)
        feedbackQuestionControl.get('option').setValue(feedbackOptions[index])
      }
    }
  }
  onFacultyFeedbackInput(event: any, index: number, feedbackQuestionControl: any): void {

    let options: IFeedbackOptions[] = this.facultyFeedbackQuestions[index].options
    feedbackQuestionControl.get('answer').setValue(String(event))
    feedbackQuestionControl.get('optionID').setValue(null)
    feedbackQuestionControl.get('option').setValue(null)

    for (let index = 0; index < options.length; index++) {
      if (event == options[index].key) {

        feedbackQuestionControl.get('answer').setValue((options[index].key).toString())
        feedbackQuestionControl.get('optionID').setValue(options[index].id)
        feedbackQuestionControl.get('option').setValue(options[index])
      }
    }
  }

  onNextClick(): void {
    this.isNextClicked = !this.isNextClicked
  }

  async openAttendance(): Promise<void> {
    console.log(this.sessionMarkAsDone);
    try {
      for (let index = 0; index < this.sessionMarkAsDone.length; index++) {
        console.log(index);
        await this.markSessionsAsComplete(this.sessionMarkAsDone[index])
      }
    } catch (err) {
      console.error(err);
    }
    const queryParams: any = { operation: this.ATTENDANCEOPERATION };
    this.tabToBeOpen = this.ATTENDANCEOPERATION
    this.updateRoute(queryParams)
  }

  async openFeedback(): Promise<void> {
    try {
      let response = await this.markSessionAttendance()
      alert(response)
    } catch (error) {
      console.log(error);
      alert(error)
      this.sendToSessionPlan()
      return
    }
    const queryParams: Params = { operation: this.FEEDBACKOPERATION };
    this.tabToBeOpen = this.FEEDBACKOPERATION;
    this.updateRoute(queryParams)
  }

  toggleTommorrow(): void {
    this.isTommorrowSessionCompleted = !this.isTommorrowSessionCompleted
  }


  submitFeedback(): void {
    // console.log(this.feedbackArray.valid);
    let errors: string[] = []
    try {
      // this.talentCount++;
      // if (this.feedbackArray.valid) {
      // if (this.talentCount == this.talents.length) {
      // this.addFacultyFeedback(errors)
      const queryParams: Params = { operation: this.ASSIGNMENTOPERATION };
      this.tabToBeOpen = this.ASSIGNMENTOPERATION
      this.updateRoute(queryParams)
      // return
      // }
      // this.addFacultyFeedback(errors)
      // this.createFeedbackQuestionsForm(this.talents[this.talentCount]?.id)
    }
    catch (error) {
      console.log(error);

    }

  }


  sendToSessionPlan(): void {
    let queryParamMap = {
      batchID: this.batchID,
      courseID: this.courseID,
      batchName: this.batchName,
    }
    this.router.navigate(['/my/batch/session/details'], {
      relativeTo: this.route,
      queryParams: queryParamMap,
    })
  }

  processPresentTalents(): void {
    this.talents = []

    for (let index = 0; index < this.talentAttendanceControlArray.controls.length; index++) {
      if (this.talentAttendanceControlArray.at(index).get("isPresent").value && !this.talentAttendanceControlArray.at(index).get("isFeedbackGiven").value) {
        // this.talents.push(this.talentAttendanceControlArray.at(index).get("talentID").value)
        this.talents.push(this.talentAttendanceControlArray.at(index).get("talent").value)

      }
    }

    console.log("processPresentTalents -> ", this.talents);

    this.talentCount = 0
    this.createFeedbackQuestionsForm(this.talents[this.talentCount]?.id)
  }


  // *********************************Assignment********************
  isAssignmentsValid(): boolean {
    if (this.selectedAssignments.length < 1) {
      return false
    }

    for (let index = 0; index < this.selectedAssignments.length; index++) {
      if (this.selectedAssignments[index].dueDate == null) {
        return false
      }
    }
    return true
  }

  async addAssignments(): Promise<void> {
    let errors: string[] = []

    for (let index = 0; index < this.selectedAssignments.length; index++) {
      let topicID = this.selectedAssignments[index].topicID
      let topicAssignmentID = this.selectedAssignments[index].id
      this.selectedAssignments[index].ProgrammingQuestionID = this.selectedAssignments[index].programmingQuestion.id
      this.selectedAssignments[index].batchSessionID = this.batchSessionID
      // console.log(topicID, topicAssignmentID, this.selectedAssignments[index]);
      try {

        await this.updateTopicAssignments(topicID, topicAssignmentID, this.selectedAssignments[index])
      } catch (error) {
        errors.push(error)
        return
      }
    }
    if (errors.length > 0) {
      let errorString = ""
      for (let index = 0; index < errors.length; index++) {
        errorString += (index == 0 ? "" : "\n") + errors[index]
      }
      alert(errorString)
      return
    }
    alert('assignments assigned successfully')
    this.sendToSessionPlan()
  }

  updateAssignment(topic: any): void {
    // console.log(topic);

    if (!this.isAssignmentMarked(topic.id)) {
      this.selectedAssignments.push(topic)
      // console.log(this.selectedAssignments);
      return
    }
    // let index = this.selectedAssignments.indexOf(subTopicID)
    let index = this.getIndexOfAssignmentMarkedID(topic.id)
    this.selectedAssignments.splice(index, 1)
    topic.dueDate = null

  }
  isAssignmentMarked(subTopicID: string): boolean {
    // console.log("isAssignment->", this.selectedAssignments);

    for (let index = 0; index < this.selectedAssignments.length; index++) {
      let id = this.selectedAssignments[index].id
      if (id == subTopicID) {
        return true
      }
    }
    return false
  }


  getIndexOfAssignmentMarkedID(subTopicID: any): number {
    for (let index = 0; index < this.selectedAssignments.length; index++) {
      let id = this.selectedAssignments[index].id
      // console.log(subTopicID);
      // console.log(id);

      if (id == subTopicID) {
        return index
      }
    }
    return -1
  }


  async markSessionsAsComplete(subTopic: any): Promise<void> {
    try {
      return await new Promise<void>((resolve, reject) => {
        this.spinnerService.loadingMessage = "Updating sessions"

        console.log(subTopic);

        this.batchSessionService.markSubTopicAsComplete(this.batchID, subTopic.subTopicID, subTopic).subscribe((response: any) => {
          console.log(response);
          // this.()
          // if (this.ongoingOperations == 0) {
          //   this.tabToBeOpen = this.ATTENDANCEOPERATION
          //   this.updateRoute(queryParams)
          // }
          resolve(response)
        }, (err: any) => {
          // this.()
          // console.error(err);
          if (err.statusText.includes('Unknown')) {
            alert("No connection to server. Check internet.")
            return
          }
          alert(err.error.error)
          reject()
        })
      })
    } finally {

    }

  }

  updateRoute(queryParams: Params): void {

    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: queryParams,
      queryParamsHandling: "merge"
    })
    // console.log(queryParams);
    // this.extractID()
    this.processTabs()

  }

  addFacultyFeedback(errors: string[]): void {
    // console.log(this.feedbackForm.value.feedbacks);
    // console.log(this.batchSessionID);


    this.spinnerService.loadingMessage = "Adding feedback"

    this.batchService.addFacultySessionFeedbacks(this.batchID, this.batchSessionID, this.feedbackForm.value.feedbacks).
      subscribe((response: any) => {
        // console.log(response)
        this.feedbackResponse = []
      }, (err: any) => {
        // console.error(err)
        if (err.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
          return
        }
        errors.push(err.error.error)
      }).add(() => {


      })
  }

  async getPresentTalents(): Promise<any> {
    try {
      return await new Promise<any>((resolve, reject) => {
        this.spinnerService.loadingMessage = "Getting Talents"

        this.totalTalent = 0
        this.batchTalents = []
        let queryParams: any = {
          isPresent: 1,
          batchSessionID: this.batchSessionID
        }
        this.batchTalentService.getBatchSessionTalents(this.batchID, this.batchSessionID, queryParams).subscribe((response: any) => {

          // this.batchTalents = response.body
          this.totalTalent = response.headers.get('X-Total-Count')

          // console.log(this.batchTalents);
          resolve(response.body)
        }, (err: any) => {
          this.totalTalent = 0
          if (err.statusText.includes('Unknown')) {
            alert("No connection to server. Check internet.")
            return
          }
          alert(err.error.error)
          console.error(err);
          reject(err)
        })
      })
    } finally {

    }

  }

  async markSessionAttendance(): Promise<void> {
    try {
      return await new Promise<any>((resolve, reject) => {

        this.spinnerService.loadingMessage = "Marking Attendance"

        // console.log(this.talentAttendanceControlArray.value);
        this.batchSessionService.addSessionPlanAttendance(this.batchID, this.batchSessionID,
          this.talentAttendanceControlArray.value).subscribe((response) => {
            // alert(response)
            // if (this.ongoingOperations <= 0) {

            // console.log(this.talents[this.talentCount].id);
            // this.createFeedbackQuestionsForm(this.talents[this.talentCount].id)
            // }
            resolve(response)
          }, (err: any) => {
            console.error(err);
            if (err.statusText.includes('Unknown')) {
              alert("No connection to server. Check internet.")
              return
            }
            // alert(err.error.error)
            reject(err.error.error)
          })
      }
      )
    } finally {

    }

  }

  getNextSessions(): void {
    this.spinnerService.loadingMessage = "Getting sessions"
    // this.totalSessions = 0

    this.nextSessionModules = []
    let queryParams: any = {
      sessionDate: this.nextSessionDate,
      sessionTopicCompletedDate: this.nextSessionDate,
      facultyID: this.loginID
    }
    this.batchSessionService.getBatchSessionPlan(this.batchID, queryParams).subscribe((response: any) => {

      // console.log(response);
      this.nextSessionDate = response.date
      this.nextSessionModules = response.module
    }, (err: any) => {
      this.totalBatchModules = 0
      // console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })


  }

  async getFeedbackQuestionsForFaculty(questionType: string): Promise<any> {
    try {
      return await new Promise<any>((resolve, reject) => {

        this.generalService.getFeedbackQuestionByType(questionType).subscribe((response: any) => {
          // console.log(response.body);
          // this.processTabs()
          // this.facultyFeedbackQuestions = response.body
          resolve(response.body)
        }, (err: any) => {
          // console.error(err)
          reject(err)
        })
      })
    } finally {

    }

  }

  getBatchSessions(): void {
    this.spinnerService.loadingMessage = "Getting sessions"
    this.totalSessions = 0

    this.completedModules = []
    this.pendingModules = []
    // let queryParams: any = {
    //   date: this.sessionDate
    // }
    let queryParams: any = {
      sessionDate: this.sessionDate,
      sessionTopicCompletedDate: this.sessionDate,
      facultyID: this.loginID,
      pendingCompletedDate: this.sessionDate,
    }
    // console.log(queryParams);
    this.batchSessionService.getBatchSessionPlan(this.batchID, queryParams).subscribe((response: any) => {
      console.log(response);
      this.sessionDate = response.date
      this.pendingModules = response.pendingModule
      this.completedModules = response.module
      this.totalSessions = this.pendingModules.length
      if (this.totalSessions > 0) {
        this.isPendingSessions = true
      }
    }, (err: any) => {
      this.totalBatchModules = 0
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    }).add(() => {


    })

  }

  // Returns all the talents.
  getTalentsForBatch(): void {

    this.totalTalent = 0
    this.batchTalents = []
    this.batchTalentService.getBatchTalentList(this.batchID).subscribe((response: any) => {
      // console.log(response.body);

      this.batchTalents = response.body
      this.totalTalent = response.headers.get('X-Total-Count')
      this.addTalentsToAttendanceForm()
      this.isTalentsLoaded = true
    }, (err: any) => {
      this.totalTalent = 0
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
      console.error(err);
    }).add(() => {


    })
  }


  getModuleTopicsAssignment(): void {
    this.spinnerService.loadingMessage = "Getting Assignments"

    let queryParams: any = {
      completedDate: this.sessionDate,
      isCompleted: 1,
      topicID: this.topicsID,
      assignedDate: 0,
    }
    this.batchTopicAssignmentService.getAllTopicAssignments(this.batchID, queryParams).subscribe((response) => {
      // console.log(response.body);
      this.assignments = response.body
      console.log(this.assignments);
    }, (err: any) => {
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
      console.error(err);
    }).add(() => {


    })

  }

  async updateTopicAssignments(topicID: string, topicAssignmentID: string, assignment: any): Promise<void> {
    try {
      return await new Promise<any>((resolve, reject) => {
        this.spinnerService.loadingMessage = "updating Assignments"

        this.batchTopicAssignmentService.updateTopicAssignment(this.batchID, topicID, topicAssignmentID, assignment).subscribe((response) => {
          // console.log(response);
          resolve(response)
        }, (err: any) => {
          if (err.statusText.includes('Unknown')) {
            alert("No connection to server. Check internet.")
            return
          }
          // alert(err.error.error)
          reject(err.error.error)
          console.error(err);
        })
      }
      )
    } finally {

    }
  }
}



// interface IBatchTalentDTO {
//   id?: string
//   batch: any
//   talent: any
//   dateOfJoining: string
//   sessionsAttendedCount: number
//   totalSessionsCount: number
//   totalHours: number
//   attendedHours: number
//   averageRating: number
//   totalFeedbacksGiven: number
//   showDetails: boolean
//   sessionDetails: any[]
// }
