import { DatePipe } from '@angular/common';
import { Component, Input, OnInit, Output, ViewChild, EventEmitter } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BatchService, IAhaMoment, IBatchModule, IBatchSession, IBatchSessionsTalent, IMappedSession, ITalentBatchDTO } from 'src/app/service/batch/batch.service';
import { Role, UrlConstant } from 'src/app/service/constant';
import { CourseService, ICourseSession } from 'src/app/service/course/course.service';
import { IFacultyCredentialList } from 'src/app/service/faculty/faculty.service';
import { GeneralService, IFeedbackQuestion, IFeeling, IFeelingLevel } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { ITalent } from 'src/app/service/talent/talent.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-batch-session',
  templateUrl: './batch-session.component.html',
  styleUrls: ['./batch-session.component.css']
})
export class BatchSessionComponent implements OnInit {

  // batch session
  batchSessionForm: FormGroup
  batchSessions: IMappedSession[]
  batchSessionID: string
  sessionList: any[]
  totalSessions: number
  isSessionPresent: boolean
  updateSession: IMappedSession
  session: ICourseSession
  assignedCourseSessions: string[]

  // batch
  batch: any
  batchID: string
  batchName: string
  sessionName: string
  sessionID: string

  // batch module
  batchModules: IBatchModule[]
  totalBatchModules: number

  // course
  courseID: string
  courseSessions: IBatchSession[]

  // faculty
  facultyList: IFacultyCredentialList[]

  // permission
  permission: IPermission
  loginID: string
  isFaculty: boolean

  // modal
  modalRef: any;

  sessionMap: any

  datesArray: any[]
  isCompletedArray: any[]

  // aha-moment
  ahaMoments: IAhaMoment[]
  ahaMomentForm: FormGroup
  ahaMomentResponseForm: FormGroup
  selectTalent: boolean
  selectFeeling: boolean
  selectResponse: boolean
  viewAhaMoment: boolean

  // talent
  talents: ITalent[]
  addSelectedTalents: any[]
  addDisabledTalents: any[]
  searchedTalentForm: FormGroup
  disableButton: boolean
  isTalentLoaded: boolean
  isTalentSearched: boolean
  totalTalent: number
  talentLimit: number
  talentOffset: number
  currentTalentPage: number;

  // feeling
  feelingForm: FormGroup
  feelingList: IFeeling[]
  feelingLevelList: IFeelingLevel[]
  feelingLevelDescription: string

  // feedback question form
  feedbackQuestionForm: FormGroup
  feedbackQuestions: IFeedbackQuestion[]

  // spinner


  // sub-session icon list
  iconList: string[]

  // batch-session talent attendance
  sessionTalentAttendance: IBatchSessionsTalent[]
  sessionTalentAttendanceForm: FormGroup

  // batch talents
  batchTalents: ITalentBatchDTO[]
  batchTalentLimit: number
  batchTalentOffset: number
  totalBatchTalents: number

  // Modal
  @ViewChild("deleteModal") deleteModal: any
  @ViewChild("ahamomentModal") ahamomentModal: any
  @ViewChild("batchSessionUpdateModal") batchSessionUpdateModal: any
  @ViewChild("batchSessionModal") batchSessionModal: any
  @ViewChild("assignAssignmentToBatchSessionModal") assignAssignmentToBatchSessionModal: any
  @ViewChild("manageAssignmentUpdate") manageAssignmentUpdate: any

  // Input from parent component
  @Input() isManage: boolean
  @Input() isAssign: boolean


  // session feedback
  @Output() assignAttendanceAndFeedbackHandler: EventEmitter<{ sessionID: string, isAttendanceGiven: boolean }>
  @Output() talentAttendanceHandler: EventEmitter<{ sessionID: string, isAttendanceGiven: boolean }>
  @Output() assignFeedbackHandler: EventEmitter<string>

  constructor(
    private formBuilder: FormBuilder,
    private urlConstant: UrlConstant,
    private courseService: CourseService,
    private batchService: BatchService,
    private generalService: GeneralService,
    private utilService: UtilityService,
    private localService: LocalService,
    private route: ActivatedRoute,
    private modalService: NgbModal,
    private spinnerService: SpinnerService,
    private role: Role,
    private datePipe: DatePipe
  ) {
    this.extractID()
    this.initializeVariables()
    this.getAllComponents()
  }

  initializeVariables(): void {
    this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_BATCH_MASTER_SESSION)

    this.loginID = this.localService.getJsonValue("loginID")
    this.isFaculty = (this.localService.getJsonValue("roleName") == this.role.FACULTY ? true : false)

    // batch-session
    this.sessionList = []
    this.batchSessions = []
    this.courseSessions = []
    this.datesArray = []
    this.isCompletedArray = []
    this.assignedCourseSessions = []


    this.totalSessions = 0
    this.totalBatchModules = 0

    this.spinnerService.loadingMessage = "Getting sessions"
    this.feelingLevelDescription = null

    this.isSessionPresent = false
    this.sessionMap = new Map()

    // aha-moment
    this.ahaMoments = []
    this.selectTalent = false
    this.selectFeeling = false
    this.selectResponse = false
    this.viewAhaMoment = false

    // talent
    this.addSelectedTalents = []
    this.addDisabledTalents = []
    this.talents = []
    this.disableButton = false
    this.totalTalent = 0
    this.isTalentLoaded = true
    this.isTalentSearched = false
    this.talentLimit = 5
    this.talentOffset = 0

    // feeling
    this.feelingLevelList = []
    this.feelingList = []

    // feedback question
    this.feedbackQuestions = []

    this.iconList = ["menu_book", "auto_stories", "import_contacts"]

    // batch session talent
    this.sessionTalentAttendance = []
    this.batchTalents = []
    this.batchTalentLimit = 10
    this.batchTalentOffset = 0
    this.totalBatchTalents = 0

    this.assignAttendanceAndFeedbackHandler = new EventEmitter()
    this.assignFeedbackHandler = new EventEmitter()
    this.talentAttendanceHandler = new EventEmitter()
  }

  getAllComponents(): void {
    // this.getSessionsForCourse()
    this.getSessionsForBatch()
    this.getBatchModules()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  extractID(): void {
    this.route.queryParamMap.subscribe(params => {
      this.batchID = params.get("batchID")
      this.courseID = params.get("courseID")
      this.batchName = params.get("batchName")
    }, err => {
      console.error(err);
    })
  }

  onAttendanceClick(session: IMappedSession): void {
    console.log(session);
    this.talentAttendanceHandler.emit({ sessionID: session.id, isAttendanceGiven: session.isAttendanceGiven })
  }

  // create a form for adding batch-sessions for form
  createBatchSessionForm() {
    this.batchSessionForm = this.formBuilder.group({
      sessions: this.formBuilder.array([], [Validators.required])
    })
  }

  get batchSessionControlArray() {
    return this.batchSessionForm.get('sessions') as FormArray
  }

  addBatchSessionsToForm() {
    this.batchSessionControlArray.push(this.formBuilder.group({
      id: new FormControl(null),
      batchID: new FormControl(this.batchID),
      courseSessionID: new FormControl(null),
      session: new FormControl(null),
      startDate: new FormControl(null, [Validators.required]),
      order: new FormControl(null, [Validators.required]),
      isCompleted: new FormControl(false, [Validators.required]),
      isMarked: new FormControl(false),
      showDetails: new FormControl(false),
      subSessions: new FormArray([]),
    }))
  }

  subSessionControlArray(index: number) {
    return this.batchSessionControlArray.at(index).get('subSessions') as FormArray
  }

  addSubsessionsToForm(index: number) {
    this.subSessionControlArray(index).push(this.formBuilder.group({
      id: new FormControl(null),
      batchID: new FormControl(this.batchID),
      courseSessionID: new FormControl(null),
      session: new FormControl(null),
      startDate: new FormControl(null),
      order: new FormControl(null),
      isCompleted: new FormControl(false),
      isMarked: new FormControl(false),
      showDetails: new FormControl(false),
      subSessions: new FormArray([]),
    }))
  }

  onAddBatchSessionClick(): void {
    this.getSessionsForCourse()
  }

  addCourseSessionsToForm(): void {
    this.totalSessions = this.batchSessions.length
    if (this.totalSessions == 0) {
      this.isSessionPresent = false
    }
    this.createBatchSessionForm()
    this.openModal(this.batchSessionModal, "xl", true)

    for (let index = 0; index < this.courseSessions.length; index++) {
      this.addBatchSessionsToForm()
      this.batchSessionControlArray.at(index).get("session").setValue(this.courseSessions[index])
      this.batchSessionControlArray.at(index).get("isMarked").setValue(false)
      this.batchSessionControlArray.at(index).get("showDetails").setValue(false)
      this.batchSessionControlArray.at(index).get("courseSessionID").setValue(this.courseSessions[index].id)
      this.batchSessionControlArray.at(index).get("order").setValue(this.courseSessions[index].order)

      this.batchSessionControlArray.at(index).get("order").disable()
      this.batchSessionControlArray.at(index).get("startDate").disable()
      this.batchSessionControlArray.at(index).get("isCompleted").disable()

      // check if session is added to batch-session
      this.markSubsessionInForm(this.courseSessions[index].subSessions, index)
      this.markAssignedSessionsInForm(index)
    }

    // console.log(this.batchSessionControlArray.value);
  }

  markSubsessionInForm(subSessions: IBatchSession[], index: number): void {
    for (let j = 0; j < subSessions.length; j++) {
      this.addSubsessionsToForm(index)

      this.subSessionControlArray(index).at(j).get("session").setValue(subSessions[j])
      this.subSessionControlArray(index).at(j).get("isMarked").setValue(false)
      this.subSessionControlArray(index).at(j).get("showDetails").setValue(false)
      this.subSessionControlArray(index).at(j).get("courseSessionID").setValue(subSessions[j].id)
    }
  }

  markAssignedSessionsInForm(index: number): void {
    this.batchSessions.find((session: IMappedSession) => {
      if (session.courseSessionID == this.courseSessions[index].id) {
        this.assignedCourseSessions.push(session.courseSessionID)
        this.batchSessionControlArray.at(index).get("id").setValue(session.id)
        this.batchSessionControlArray.at(index).get("startDate").setValue(
          this.datePipe.transform(new Date(session.startDate).toUTCString(), "yyyy-MM-ddTHH:mm", "GMT")
        )
        this.batchSessionControlArray.at(index).get("order").setValue(session.order)
        this.batchSessionControlArray.at(index).get("isCompleted").setValue(session.isCompleted)
        this.batchSessionControlArray.at(index).get("isMarked").setValue(true)
        this.batchSessionControlArray.at(index).get("showDetails").setValue(false)

        this.batchSessionControlArray.at(index).get("order").enable()
        this.batchSessionControlArray.at(index).get("order").setValue(session.order)
        this.batchSessionControlArray.at(index).get("startDate").enable()
        this.batchSessionControlArray.at(index).get("isCompleted").enable()

        if (session.session.subSessions.length > 0) {
          this.batchSessionControlArray.at(index).get("showDetails").setValue(true)
        }

        this.markAssignedSubSessions(session.session.subSessions, index)
        return
      }
    })
  }

  markAssignedSubSessions(subSessions: IBatchSession[], index: number): void {
    for (let j = 0; j < subSessions.length; j++) {
      this.assignedCourseSessions.push(subSessions[j].id)
      this.subSessionControlArray(index).at(j).get("session").setValue(subSessions[j])
      this.subSessionControlArray(index).at(j).get("courseSessionID").setValue(subSessions[j].id)
      this.subSessionControlArray(index).at(j).get("isMarked").setValue(true)
    }
  }

  toggleSubSessionWithForm(subSessionControl: any, index: number, j: number): void {
    // remove
    if (this.assignedCourseSessions.includes(subSessionControl.get("courseSessionID").value)) {
      this.removeSubSessionInFormArray(subSessionControl, index)
      return
    }

    // add
    if (!this.assignedCourseSessions.includes(subSessionControl.get("courseSessionID").value)) {
      this.addSubSessionToFormArray(subSessionControl, index)
    }
  }

  addSubSessionToFormArray(subSessionControl: any, index: number): void {
    subSessionControl.get("isMarked").setValue(true)
    this.assignedCourseSessions.push(subSessionControl.get("courseSessionID").value)
    this.batchSessionControlArray.at(index).get('isMarked').setValue(true)
  }

  removeSubSessionInFormArray(subSessionControl: any, index: number): void {
    let courseSessionIDIndex = this.assignedCourseSessions.indexOf(subSessionControl.get("courseSessionID").value)
    this.assignedCourseSessions.splice(courseSessionIDIndex, 1)
    subSessionControl.get('isMarked').setValue(false)

    let isMarkedFound = false
    for (let k = 0; k < this.subSessionControlArray(index).controls.length; k++) {
      if (this.subSessionControlArray(index)?.at(k)?.get("isMarked")?.value) {
        isMarkedFound = true
        break
      }
    }

    if (!isMarkedFound) {
      this.batchSessionControlArray.at(index).get("isMarked").setValue(false)

      this.batchSessionControlArray.at(index).get("order").disable()
      this.batchSessionControlArray.at(index).get("startDate").disable()
      this.batchSessionControlArray.at(index).get("isCompleted").disable()
    }

  }

  toggleSessionWithForm(sessionControl: any, index: number): void {
    // remove
    if (this.assignedCourseSessions.includes(sessionControl.get("courseSessionID").value)) {
      this.removeAssignedSession(sessionControl, index)
      return
    }
    // add
    if (!this.assignedCourseSessions.includes(sessionControl.get("courseSessionID").value)) {
      this.assignSession(sessionControl, index)
    }
  }

  removeAssignedSession(sessionControl: any, index: number): void {
    // console.log(sessionControl);
    let courseSessionIDIndex = this.assignedCourseSessions.indexOf(sessionControl.get("courseSessionID").value)
    this.assignedCourseSessions.splice(courseSessionIDIndex, 1)

    sessionControl.get("isMarked").setValue(false)

    // sessionControl.get("order").setValue(0)
    // sessionControl.get("startDate").setValue(null)

    sessionControl.get("order").disable()
    sessionControl.get("startDate").disable()
    sessionControl.get("isCompleted").disable()

    // remove all sub-sessions added for current session
    this.removeSubsesssionToForm(index)
  }

  removeSubsesssionToForm(index: number): void {
    // console.log(sessionControl);
    for (let j = 0; j < this.subSessionControlArray(index).length; j++) {
      this.subSessionControlArray(index).at(j).get("isMarked").setValue(false)
    }
  }

  assignSession(sessionControl: any, index: number): void {
    // console.log(sessionControl);

    this.assignedCourseSessions.push(sessionControl.get("courseSessionID").value)
    sessionControl.get("isMarked").setValue(true)
    sessionControl.get("showDetails").setValue(false)

    sessionControl.get("order").enable()
    sessionControl.get("startDate").enable()
    sessionControl.get("isCompleted").enable()

    // this.utilService.updateValueAndValiditors(sessionControl)

    // add all sub-sessions in form
    this.addSubsessionToForm(sessionControl, index)
  }

  addSubsessionToForm(sessionControl: any, index: number): void {
    if (sessionControl.get("subSessions").value.length > 0) {
      sessionControl.get("showDetails").setValue(true)
    }

    for (let j = 0; j < this.subSessionControlArray(index).controls.length; j++) {
      this.assignedCourseSessions.push(this.subSessionControlArray(index).at(j).get("courseSessionID").value)
      this.subSessionControlArray(index).at(j).get("isMarked").setValue(true)
    }
  }

  // called when single batch is being updated
  onBatchSessionUpdateClick(session: IMappedSession): void {
    // console.log(session);
    // this.getFacultyCredentialList()
    this.openModal(this.batchSessionUpdateModal, 'md')
    this.updateSession = session
    this.createBatchSessionForm()
    this.addBatchSessionsToForm()

    // used tempSession array as patchValue reqiures array
    let tempSession = []
    tempSession.push(session)
    this.batchSessionForm.get('sessions').patchValue(tempSession)
    this.batchSessionControlArray.at(0).get("startDate").setValue(
      this.datePipe.transform(new Date(session.startDate).toUTCString(), "yyyy-MM-ddTHH:mm", "GMT")
    )
  }

  updateSessionForm(): void {
    for (let index = 0; index < this.courseSessions.length; index++) {
      this.addBatchSessionsToForm()

      const isCompletedControl = this.batchSessionControlArray.at(index).get('isCompleted')
      if (this.isSessionPresent) {
        isCompletedControl.setValue(false)
        isCompletedControl.setValidators([Validators.required])
      } else {
        isCompletedControl.setValue(false)
        isCompletedControl.setValidators(null)
        isCompletedControl.markAsUntouched()
      }
      isCompletedControl.updateValueAndValidity()
    }
  }

  // sets the view for cards
  setBatchSessionViewField(): void {
    for (let index = 0; index < this.batchSessions.length; index++) {
      this.batchSessions[index].session.viewSubSessionClicked = false
      // this.batchSessions[index].session.cardColumn = "col-md-6 col-sm-12 d-flex"
    }
  }

  // sets the view for cards
  setCourseSessionViewField(): void {
    for (let index = 0; index < this.courseSessions.length; index++) {
      this.courseSessions[index].viewSubSessionClicked = false
    }
  }

  deleteSession(session: ICourseSession, batchSessionID: string): void {
    // this.session = session
    // this.batchSessionID = batchSessionID
    // this.openModal(this.deleteModal, 'md')
    this.openModal(this.deleteModal, 'md').result.then(() => {
      this.deleteSessionFromBatch(session.id, batchSessionID)
    }, (err) => {
      console.error(err);
      return
    })
  }

  addSubSessionToBatchSessionFormArray(): void {
    for (let index = 0; index < this.batchSessionControlArray.controls.length; index++) {
      for (let j = 0; j < this.subSessionControlArray(index).controls.length; j++) {
        let subSession = this.subSessionControlArray(index).at(j).value
        if (subSession.isMarked) {
          this.addBatchSessionsToForm()
          let lastIndex = this.batchSessionControlArray.controls.length - 1

          this.batchSessionControlArray.at(lastIndex).get("session").setValue(subSession.session)
          this.batchSessionControlArray.at(lastIndex).get("isMarked").setValue(subSession.isMarked)
          this.batchSessionControlArray.at(lastIndex).get("showDetails").setValue(false)
          this.batchSessionControlArray.at(lastIndex).get("courseSessionID").setValue(subSession.courseSessionID)
          this.batchSessionControlArray.at(lastIndex).get("order").setValue(null)
          this.batchSessionControlArray.at(lastIndex).get("startDate").setValue(null)
          this.batchSessionControlArray.at(lastIndex).get("isCompleted").setValue(null)

          this.batchSessionControlArray.at(lastIndex).get("order").clearValidators()
          this.batchSessionControlArray.at(lastIndex).get("startDate").clearValidators()
          this.batchSessionControlArray.at(lastIndex).get("isCompleted").clearValidators()

          this.utilService.updateValueAndValiditors(this.batchSessionControlArray.at(lastIndex) as FormGroup)
          // this.subSessionControlArray(index).removeAt(j)
        }
      }
      // (this.batchSessionControlArray.at(index) as FormGroup).removeControl("subSessions")
    }
  }

  validateBatchSession(): void {
    this.addSubSessionToBatchSessionFormArray()
    console.log(this.batchSessionForm.controls);

    if (this.batchSessionForm.invalid) {
      return
    }

    let batchSessions = this.batchSessionForm.value.sessions
    this.removeUnmarkedSessions(batchSessions)
    console.log(batchSessions);

    if (this.isSessionPresent) {
      this.updateSessionsForBatch(batchSessions)
      return
    }
    this.addSessionToBatch(batchSessions)
  }

  removeUnmarkedSessions(batchSessions: any[]): void {
    for (let index = 0; index < batchSessions.length; index++) {
      if (!batchSessions[index].isMarked) {
        batchSessions.splice(index, 1)
        index--
      }
    }
  }

  validateSingleSession() {
    if (this.batchSessionForm.invalid) {
      this.batchSessionControlArray.markAllAsTouched()
      return
    }
    this.updateSessionForBatch()
  }

  setSessionForUpdate(session: IMappedSession): void {
    if (confirm("Are you sure?")) {
      this.updateSession = session
      this.updateSession.isCompleted = true
      // console.log(session);
      this.updateSessionCompleteStatus()
      this.assignAttendanceAndFeedbackHandler.emit({ sessionID: session.id, isAttendanceGiven: session.isAttendanceGiven })
    }
  }

  assignMaterialIIconToBatchSubSessions(): void {
    for (let index = 0; index < this.batchSessions.length; index++) {
      if (this.batchSessions[index].session.subSessions) {
        this.getRandomIcon(this.batchSessions[index].session.subSessions)
      }
    }
  }

  assignMaterialIIconToCourseSubSessions(): void {
    for (let index = 0; index < this.courseSessions.length; index++) {
      if (this.courseSessions[index].subSessions) {
        this.getRandomIcon(this.courseSessions[index].subSessions)
      }
    }
  }

  getRandomIcon(subSessions: IBatchSession[]): void {
    for (let index = 0; index < subSessions.length; index++) {
      let iconIndex = Math.floor(Math.random() * this.iconList.length)
      subSessions[index].materialIcon = this.iconList[iconIndex]
    }
  }

  getBatchModules(): void {
    this.spinnerService.loadingMessage = "Getting modules"

    this.totalBatchModules = 0
    this.batchModules = []
    this.batchService.getOldBatchModules(this.batchID).subscribe((response: any) => {
      this.batchModules = response.body
      this.totalBatchModules = this.batchModules.length
      if (this.totalBatchModules > 0) {
        this.isSessionPresent = true
      }
      console.log(this.batchModules);

    }, (err: any) => {
      this.totalBatchModules = 0
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  getSessionsForBatch(): void {
    this.spinnerService.loadingMessage = "Getting sessions"
    // this.changeLoadingMessage.emit(this.spinnerService.loadingMessage)

    this.totalSessions = 0

    this.batchService.getSessionForBatch(this.batchID).subscribe((response: any) => {
      this.batchSessions = response.body
      console.log(this.batchSessions);
      this.totalSessions = this.batchSessions.length
      if (this.totalSessions > 0) {
        this.isSessionPresent = true
      }
      this.assignMaterialIIconToBatchSubSessions()
      this.setBatchSessionViewField()
    }, (err: any) => {
      this.totalSessions = 0
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  getSessionsForCourse(): void {
    this.spinnerService.loadingMessage = "Getting sessions"
    // this.changeLoadingMessage.emit(this.spinnerService.loadingMessage)

    this.courseService.getSessionsForCourse(this.courseID).subscribe((response: any) => {
      this.courseSessions = response.body
      // console.log(this.courseSessions);
      this.addCourseSessionsToForm()
    }, (err: any) => {
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  addSessionToBatch(batchSessions: any[]): void {
    // console.log(this.batchSessionForm.value.sessions);
    this.spinnerService.loadingMessage = "Adding Sessions"
    // this.changeLoadingMessage.emit(this.spinnerService.loadingMessage)

    this.batchService.addSessions(this.batchID, batchSessions).subscribe((response: any) => {
      // console.log(response);
      alert("Sessions successfully added for batch")
      this.modalRef.close()
      this.sessionList = []
      // this.resetArrayVariables()
      this.getSessionsForBatch()
      this.getBatchModules()
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
      console.error(err.error.error);
    })
  }

  updateSessionCompleteStatus(): void {
    this.spinnerService.loadingMessage = "Updating Session"
    // this.changeLoadingMessage.emit(this.spinnerService.loadingMessage)

    // this.updateSession.startDate = this.updateSession.startDate.substring(0, 16)
    this.updateSession.startDate = this.datePipe.transform(new Date(this.updateSession.startDate).toUTCString(), "yyyy-MM-ddTHH:mm", "GMT")
    // console.log(this.updateSession);
    this.batchService.updateSession(this.batchID, this.updateSession).subscribe((response: any) => {
      alert("Session successfully updated for batch")
      this.getSessionsForBatch()
      this.getBatchModules()
    }, (err: any) => {
      this.updateSession.isCompleted = false
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
      console.error(err.error.error);
    })
  }

  // update's single session
  updateSessionForBatch(): void {
    // console.log(this.batchSessionForm.value.sessions);
    this.spinnerService.loadingMessage = "Updating Session"
    // this.changeLoadingMessage.emit(this.spinnerService.loadingMessage)

    this.batchService.updateSession(this.batchID, this.batchSessionForm.value.sessions[0]).subscribe((response: any) => {
      this.modalRef.close()
      alert("Session successfully updated for batch")
      this.getSessionsForBatch()
      this.getBatchModules()
    }, (err: any) => {
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
      console.error(err.error.error);
    })
  }

  // update's multiple sessions
  updateSessionsForBatch(batchSessions: any[]): void {
    this.spinnerService.loadingMessage = "Updating Sessions"
    // this.changeLoadingMessage.emit(this.spinnerService.loadingMessage)

    // console.log(this.batchSessionForm.value.sessions);
    this.batchService.updateSessions(this.batchID, batchSessions).subscribe((response: any) => {
      // console.log(response);
      alert("Sessions successfully updated for batch")
      this.modalRef.close()
      this.sessionList = []
      // this.resetArrayVariables()
      this.getSessionsForBatch()
      this.getBatchModules()
    }, (err: any) => {
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
      console.error(err.error.error);
    })
  }

  deleteSessionFromBatch(sessionID: string, batchSessionID: string) {
    // console.log(this.session);
    this.spinnerService.loadingMessage = "Deleting session from batch"
    // this.changeLoadingMessage.emit(this.spinnerService.loadingMessage)

    this.batchService.deleteSessionFromBatch(this.batchID, sessionID, batchSessionID).subscribe((response: any) => {
      this.modalRef.close()
      alert("Session successfully deleted from batch")
      // this.resetArrayVariables()
      this.getSessionsForBatch()
      this.getBatchModules()
    }, (err: any) => {
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  // getFacultyCredentialList(): void {
  //   this.facultyService.getFacultyCredentialList().subscribe((response: any) => {
  //     this.facultyCredentialList = response.body
  //   }, (err: any) => {
  //     console.error(err);
  //   })
  // }

  // ============================================================= AHA MOMENT =============================================================

  createAhaMomentsForm(): void {
    this.ahaMomentForm = this.formBuilder.group({
      moments: this.formBuilder.array([])
    })
  }

  get ahaMomentsFormControlArray() {
    return this.ahaMomentForm.get('moments') as FormArray
  }

  addAhaMomentsToForm(): void {
    this.ahaMomentsFormControlArray.push(this.formBuilder.group({
      id: new FormControl(null),
      facultyID: new FormControl(this.loginID),
      talentID: new FormControl(null, [Validators.required]),
      feeling: new FormControl(null, [Validators.required]),
      feelingLevel: new FormControl(null, [Validators.required]),
      ahaMomentResponse: new FormControl(Array(), [Validators.required]),
    }))
  }

  onAhaMomentsClick(session: IMappedSession) {
    this.sessionID = session.id
    this.sessionName = session.session.name
    this.initializeAhaMomentVariables()
    this.openModal(this.ahamomentModal, 'xl')
  }

  initializeAhaMomentVariables(): void {
    this.addSelectedTalents = []
    this.addDisabledTalents = []
    this.talents = []
    this.ahaMoments = []
    this.getAllAhaMoments()
    this.changeTalentPage(1)
    this.selectTalent = false
    this.selectFeeling = false
    this.selectResponse = false
    this.viewAhaMoment = true
  }

  onAddAhaMomentClick(): void {
    this.getFeedbackQuestions()
    this.getFeelingList()
    this.createAhaMomentsForm()
    this.createFeelingForm()
    this.createFeedbackQuestionForm()
    this.viewAhaMoment = false
    this.selectTalent = true
    this.feelingLevelDescription = null
  }

  onBackButtonClick(): void {
    this.addSelectedTalents = []
    this.feelingForm.reset()
    this.feedbackQuestionForm.reset()
    this.ahaMomentForm.reset()
  }

  onPreviousButtonClick() {

    if (!this.selectTalent && this.selectFeeling && !this.selectResponse) {
      this.selectTalent = true
      this.selectFeeling = false
      this.selectResponse = false
      return
    }
    if (!this.selectTalent && !this.selectFeeling && this.selectResponse) {
      this.selectTalent = false
      this.selectFeeling = true
      this.selectResponse = false
      return
    }
  }

  onNextButtonClick(): void {
    // STEP - 1
    if (this.selectTalent && !this.selectFeeling && !this.selectResponse) {
      if (this.addSelectedTalents.length == 0) {
        alert("Please select atleast 1 student.")
        return
      }

      this.addTalentsToForm()
      this.selectTalent = false
      this.selectFeeling = true
      this.selectResponse = false
      return
    }

    // STEP - 2
    if (!this.selectTalent && this.selectFeeling && !this.selectResponse) {
      if (this.feelingForm.invalid) {
        this.feelingForm.markAllAsTouched()
        return
      }
      this.addFeelingsToForm()
      this.selectTalent = false
      this.selectFeeling = false
      this.selectResponse = true

      if (this.feedbackQuestionForm && this.feedbacksArray.controls.length == 0) {
        for (let index = 0; index < this.feedbackQuestions.length; index++) {
          this.addFeedbackQuestions()
          this.feedbacksArray.at(index).get("question").setValue(this.feedbackQuestions[index])
        }
      }
      return
    }
  }

  addTalentsToForm(): void {
    this.ahaMomentForm.reset()
    this.createAhaMomentsForm()

    for (let index = 0; index < this.addSelectedTalents.length; index++) {
      this.addAhaMomentsToForm()
      this.ahaMomentsFormControlArray.at(index).get("talentID").setValue(this.addSelectedTalents[index])
    }
    // console.log(this.ahaMomentsFormControlArray.controls)

  }

  addFeelingsToForm(): void {
    for (let index = 0; index < this.ahaMomentsFormControlArray.controls.length; index++) {
      this.ahaMomentsFormControlArray.at(index).get("feeling").setValue(this.feelingForm.get('feeling').value)
      this.ahaMomentsFormControlArray.at(index).get("feelingLevel").setValue(this.feelingForm.get('feelingLevel').value)
    }
  }

  addFeedbackQuestionsToForm(): void {
    for (let index = 0; index < this.ahaMomentsFormControlArray.controls.length; index++) {
      this.ahaMomentsFormControlArray.at(index).get("ahaMomentResponse").setValue(this.feedbackQuestionForm.value.feedbacks)
    }
  }


  // ******************************************* CRUD *******************************************

  addAhaMoment(): void {
    this.spinnerService.loadingMessage = "Adding aha moments"
    // this.changeLoadingMessage.emit(this.spinnerService.loadingMessage)

    this.batchService.addAhaMoment(this.batchID, this.sessionID, this.ahaMomentForm.value.moments).subscribe((response: any) => {
      // console.log(response)
      this.initializeAhaMomentVariables()
      this.ahaMomentForm.reset()
      this.feelingForm.reset()
      this.feedbackQuestionForm.reset()
      this.feelingLevelDescription = null
    }, (error: any) => {
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error.error)
      console.error(error)
    })
  }

  getAllAhaMoments(): void {
    this.spinnerService.loadingMessage = "Getting aha moments"
    // this.changeLoadingMessage.emit(this.spinnerService.loadingMessage)

    this.ahaMoments = []
    this.addDisabledTalents = []
    this.batchService.getAllAhaMoments(this.batchID, this.sessionID).subscribe((response: any) => {
      this.ahaMoments = response.body
      // console.log(this.ahaMoments);
      for (let index = 0; index < this.ahaMoments.length; index++) {
        this.addDisabledTalents.push(this.ahaMoments[index].talent.id)
      }
    }, (error: any) => {
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error.error)
      console.error(error)
    })
  }

  deleteTalentAhaMoment(ahaMomentID: string): void {
    // console.log(ahaMomentID);
    if (!confirm("Are you sure you want to delete talent's aha moment?")) {
      return
    }
    this.spinnerService.loadingMessage = "Deleting talents aha moment"
    // this.changeLoadingMessage.emit(this.spinnerService.loadingMessage)

    this.batchService.deleteAhaMoment(ahaMomentID).subscribe((response: any) => {
      // console.log(response);
      alert("Talent's aha moment successfully deleted")
      this.getAllAhaMoments()
    }, (error: any) => {
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error.error)
      console.error(error)
    })
  }

  validateAhaMoment(): void {

    this.addFeedbackQuestionsToForm()
    // console.log(this.feedbackQuestionForm.controls);

    if (this.feedbackQuestionForm.invalid) {
      this.feedbackQuestionForm.markAllAsTouched()
      return
    }

    // console.log(this.ahaMomentForm.controls)
    if (this.ahaMomentForm.invalid) {
      alert("Form is invalid.")
      return
    }

    this.addAhaMoment()
  }

  // ============================================================= FEELING =============================================================

  createFeelingForm(): void {
    this.feelingForm = this.formBuilder.group({
      feeling: new FormControl(null, [Validators.required]),
      feelingLevel: new FormControl(null, [Validators.required]),
    })
  }

  getFeelingList(): void {
    // this.spinnerService.loadingMessage = "Getting feelings list"

    this.feelingList = []
    this.generalService.getFeelingList().subscribe((response: any) => {
      this.feelingList = response.body
    }, (error: any) => {
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error.error)
      console.error(error)
    })
  }

  getFeelingLevelList(feelingID: string): void {
    this.spinnerService.loadingMessage = "Getting feelings list"

    this.feelingLevelList = []
    this.generalService.getFeelingLevelList(feelingID).subscribe((response: any) => {
      this.feelingLevelList = response.body
    }, (error: any) => {
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error.error)
      console.error(error)
    })
  }

  // ============================================================= AHA MOMENT QUESTIONS =============================================================

  createFeedbackQuestionForm(): void {
    this.feedbackQuestionForm = this.formBuilder.group({
      feedbacks: this.formBuilder.array([], [Validators.required])
    })
  }

  get feedbacksArray() {
    return this.feedbackQuestionForm.get('feedbacks') as FormArray
  }

  addFeedbackQuestions(): void {
    this.feedbacksArray.push(this.formBuilder.group({
      question: new FormControl(null, [Validators.required]),
      response: new FormControl(null, [Validators.required, Validators.maxLength(200)])
    }))
  }

  getFeedbackQuestions(): void {
    // this.spinnerService.loadingMessage = "Getting feedback questions"

    this.generalService.getFeedbackQuestionByType("Aha_Moment_Feedback").subscribe((response: any) => {
      this.feedbackQuestions = response.body
      // console.log(response.body);
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  // ============================================================= TALENTS =============================================================

  checkTalentAdded(talentID) {
    return this.addSelectedTalents.includes(talentID)
  }

  checkAllTalentAdded() {
    for (let index = 0; index < this.talents.length; index++) {
      if (!this.addSelectedTalents.includes(this.talents[index].id)) {
        return false
      }
    }
    return true
  }

  isCheckboxDisabled(talentID) {
    return this.addDisabledTalents.includes(talentID)
  }

  // takes a list called selectedTalent and adds all the checked talents to list, also does not contain duplicate values
  addTalentsToList(event, id: string) {
    if (event.target.checked) {
      this.addSelectedTalents.push(id)
    } else {
      if (this.addSelectedTalents.includes(id)) {
        let index = this.addSelectedTalents.indexOf(id)
        this.addSelectedTalents.splice(index, 1)
      }
    }
    // console.log(this.addSelectedTalents);
  }

  addAllTalentsToList(event) {
    if (event.target.checked) {
      for (let index = 0; index < this.talents.length; index++) {
        if (!this.addSelectedTalents.includes(this.talents[index].id)) {
          this.addSelectedTalents.push(this.talents[index].id)
        }
      }
      for (let index = 0; index < this.addDisabledTalents.length; index++) {
        if (this.addSelectedTalents.includes(this.addDisabledTalents[index])) {
          let indexOf = this.addSelectedTalents.indexOf(this.addDisabledTalents[index])
          this.addSelectedTalents.splice(indexOf, 1)
        }
      }
    } else {
      this.addSelectedTalents = []
    }

  }

  resetSelectedTalent() {
    this.addSelectedTalents = []
    // this.ahaMomentForm.reset()
  }

  // Returns all the talents.
  getTalentsForBatch(): void {
    // this.spinnerService.loadingMessage = "Getting Talents"

    this.totalTalent = 0
    this.isTalentLoaded = true
    this.talents = []
    this.batchService.getTalentsForBatch(this.batchID, this.talentLimit, this.talentOffset).subscribe((response: any) => {
      this.talents = response.body
      this.totalTalent = response.headers.get('X-Total-Count')
      if (this.totalTalent == 0) {
        this.isTalentLoaded = false
      }
      // console.log(this.talents)
    }, (err: any) => {
      this.totalTalent = 0
      this.isTalentLoaded = false
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
      console.error(err);
    })
  }

  // Handles pagination
  changeTalentPage($event: any): void {
    // $event will be the page number & offset will be 1 less than it.
    this.talentOffset = $event - 1
    this.currentTalentPage = $event
    this.getTalentsForBatch()
  }



  dismissModal(modal: NgbModalRef): void {
    modal.dismiss()
    this.sessionList = []
  }

  onSessionFeedbackClick(session: any): void {
    this.assignFeedbackHandler.emit(session.id)
  }

  openModal(modalContent: any, modalSize = "lg", withExtraOptions: boolean = false): NgbModalRef {
    let options: NgbModalOptions = {
      ariaLabelledBy: 'modal-basic-title', backdrop: 'static',
      size: modalSize, keyboard: false
    }

    if (withExtraOptions) {
      options.scrollable = true
      options.centered = true
    }
    this.modalRef = this.modalService.open(modalContent, options)
    return this.modalRef
  }

  // ============================================================= THE END =============================================================


  // resetArrayVariables(): void {
  //   this.datesArray = []
  //   this.isCompletedArray = []
  // }

  // firstly add the new batchsession form to formarray
  // it then checks if its id, startDate, isCompleted fields are present and sets those fields
  // createAndAddSessionToForm(session: IBatchSession, mappedSession?: IMappedSession): void {
  //   this.addBatchSessionsToForm()
  //   let formIndex = this.batchSessionControlArray.controls.length - 1
  //   if (mappedSession?.id) {
  //     this.batchSessionControlArray.at(formIndex).get('id').setValue(mappedSession.id)
  //   }
  //   if (mappedSession?.startDate) {
  //     this.batchSessionControlArray.at(formIndex).get('startDate').setValue(mappedSession.startDate)
  //   }
  //   this.batchSessionControlArray.at(formIndex).get('courseSessionID').setValue(session.id)

  //   // start date is required for main-session 
  //   // check if sessionID exist and if exist then make startDate and isCompleted field not required
  //   if (session.sessionID) {
  //     const startDateControl = this.batchSessionControlArray.at(formIndex).get('startDate')
  //     startDateControl.setValue(null)
  //     startDateControl.setValidators(null)
  //     startDateControl.updateValueAndValidity()

  //     const isCompletedControl = this.batchSessionControlArray.at(formIndex).get('isCompleted')
  //     isCompletedControl.setValue(null)
  //     isCompletedControl.setValidators(null)
  //     isCompletedControl.updateValueAndValidity()
  //   }
  // }

  // toggleSession will throw an alert for user confirmation and accordingly calls toggleSesionFromList func
  // toggleSession(event, session: IBatchSession, arrayIndex?: number): void {
  //   if (this.sessionList.includes(session.id)) {
  //     if (confirm("Are your sure want to remove this session?")) {
  //     session.viewSubSessionClicked = false
  //     this.toggleSessionFromList(event, session, arrayIndex)
  //       return
  //     }
  //     event.target.checked = true
  //     return
  //   }
  //   if (session?.subSessions?.length > 0) {
  //     session.viewSubSessionClicked = true
  //   }
  //   this.toggleSessionFromList(event, session, arrayIndex)
  // }

  // will add or remove the sessionID from sessionList array and also create or remove batchSession form from formarray
  // toggleSessionFromList(event, session: IBatchSession, arrayIndex?: number): void {

  //   // when btn checked is true add sessionID to sessionList and formarray
  //   if (event.target.checked) {
  //     this.addSessionsWhenChecked(session)
  //     return
  //   }

  //   // when btn checked is false remove sessionID from sessionList and formarray
  //   if (this.sessionList.includes(session.id)) {
  //     this.removeSessionsWhenUnchecked(session, arrayIndex)
  //   }
  // }

  // addSessionsWhenChecked(session: IBatchSession): void {
  //   this.sessionList.push(session.id)
  //   this.createAndAddSessionToForm(session)

  //   // add sub-sessions to list if they exist
  //   if (session.subSessions) {
  //     this.addSubSessionsToSessionList(session)
  //   }

  //   // add session if sub-session is add
  //   if (session.sessionID && !this.sessionList.includes(session.sessionID)) {
  //     this.addSessionWhenSubSessionAdded(session)
  //   }
  // }

  // addSubSessionsToSessionList(session: IBatchSession): void {
  //   for (let index = 0; index < session.subSessions.length; index++) {
  //     this.sessionList.push(session.subSessions[index].id)
  //     this.createAndAddSessionToForm(session.subSessions[index])
  //   }
  // }

  // addSessionWhenSubSessionAdded(session: IBatchSession): void {
  //   for (let i = 0; i < this.courseSessions.length; i++) {
  //     if (session.sessionID == this.courseSessions[i].id) {
  //       this.sessionList.push(this.courseSessions[i].id)
  //       this.createAndAddSessionToForm(this.courseSessions[i])
  //     }
  //   }
  // }

  // removeSessionsWhenUnchecked(session: IBatchSession, arrayIndex?: number): void {
  //   let index = this.sessionList.indexOf(session.id)
  //   this.sessionList.splice(index, 1)
  //   this.batchSessionControlArray.removeAt(index)
  //   // console.log(session);

  //   // if session.sessionID does not exist means it is a session and its date and isCompleted fields should be reset
  //   if (!session.sessionID) {
  //     if (arrayIndex) {
  //       this.datesArray[arrayIndex] = null
  //       this.isCompletedArray[arrayIndex] = false
  //     } else {
  //       this.datesArray[index] = null
  //       this.isCompletedArray[index] = false
  //     }
  //   }

  //   // sub-sessions for sessions are removed from sessionList if they were added
  //   if (session.subSessions) {
  //     this.removeSubSessionsFromSessionList(session)
  //   }
  // }

  // removeSubSessionsFromSessionList(session: IBatchSession): void {
  //   for (let i = 0; i < session.subSessions.length; i++) {
  //     if (this.sessionList.includes(session.subSessions[i].id)) {
  //       let index = this.sessionList.indexOf(session.subSessions[i].id)
  //       this.sessionList.splice(index, 1)
  //       this.batchSessionControlArray.removeAt(index)
  //     }
  //   }
  // }


  // checkSessionAdded(session: IBatchSession) {
  //   return this.sessionList.includes(session.id)
  // }

  // initalizes the isCompleted and dates array 
  // also checks all the id's of sessions that were previously added to batch
  // updateSessionList() {
  //   // used to get the index to store the date and isCompleted field
  //   let tempCourseSessionID = []

  //   // populate sessionMap
  //   tempCourseSessionID = this.initializeCourseSessionVariables()

  //   for (let index = 0; index < this.batchSessions.length; index++) {
  //     this.sessionMap.set(this.batchSessions[index].session.id, this.getBatchSessionMapCount(index) + 1)

  //     // checks if session has repeated itself which means session is added to the batch
  //     if (this.getBatchSessionMapCount(index) > 1) {
  //       this.sessionList.push(this.batchSessions[index].session.id)
  //       this.createAndAddSessionToForm(this.batchSessions[index].session, this.batchSessions[index])
  //       // firstly pop's dateArray as it contains null and is replaced by the start date of the session
  //       // pop is made to save the sessions date at the appropriate location and same for isCompletedArray
  //       let k = tempCourseSessionID.indexOf(this.batchSessions[index].session.id)
  //       this.spliceDateAndIsCompleted(index, k)
  //       this.setDate(this.batchSessions[index].session.id, k)
  //       this.setIsCompleted(this.batchSessions[index].session.id, k)

  //       // add sub-session to sessionList if they exist
  //       if (this.batchSessions[index].session.subSessions) {
  //         this.addSubSessionsToList(index)
  //       }
  //     }
  //   }
  // }

  // initializeCourseSessionVariables(): any[] {
  //   let tempCourseSessionID = []

  //   for (let index = 0; index < this.courseSessions.length; index++) {
  //     this.sessionMap.set(this.courseSessions[index].id, 1)
  //     tempCourseSessionID.push(this.courseSessions[index].id)
  //     this.datesArray.push(null)
  //     this.isCompletedArray.push(false)
  //   }

  //   return tempCourseSessionID
  // }

  // spliceDateAndIsCompleted(batchSessionIndex: number, spliceIndex:number): void {
  //   this.datesArray.splice(spliceIndex, 1)
  //   this.datesArray.splice(spliceIndex, 0, this.batchSessions[batchSessionIndex].startDate)
  //   this.isCompletedArray.splice(spliceIndex, 1)
  //   this.isCompletedArray.splice(spliceIndex, 0, this.batchSessions[batchSessionIndex].isCompleted)
  // }

  // getBatchSessionMapCount(elementIndex: number): number {
  //   return this.sessionMap.get(this.batchSessions[elementIndex].session.id)
  // }

  // addSubSessionsToList(batchSessionIndex: number): void {
  //   for (let j = 0; j < this.batchSessions[batchSessionIndex].session.subSessions.length; j++) {
  //     if (!this.sessionList.includes(this.batchSessions[batchSessionIndex].session.subSessions[j].id)) {
  //       this.sessionList.push(this.batchSessions[batchSessionIndex].session.subSessions[j].id)
  //       this.createAndAddSessionToForm(this.batchSessions[batchSessionIndex].session.subSessions[j])
  //     }
  //   }
  // }

  // sets the date field in the form array
  // setDate(sessionID: string, k: number): void {
  //   let index = this.sessionList.indexOf(sessionID)
  //   this.batchSessionControlArray.at(index).get('startDate').setValue(this.datesArray[k].slice(0, 16))
  //   this.datesArray[k] = (this.datesArray[k]).slice(0, 16)
  // }

  // sets the isCompleted field in the form array
  // setIsCompleted(id: any, k: number): void {
  //   let index = this.sessionList.indexOf(id)
  //   this.batchSessionControlArray.at(index).get('isCompleted').setValue(this.isCompletedArray[k])
  // }


}
