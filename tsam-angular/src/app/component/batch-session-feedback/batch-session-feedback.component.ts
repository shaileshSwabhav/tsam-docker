import { Location } from '@angular/common';
import { HttpParams } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { FormGroup, FormArray, FormControl, Validators, FormBuilder } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BatchService, IFacultyBatchSessionFeedback, ITalentBatchSessionFeebackDTO, ITalentBatchSessionFeedback } from 'src/app/service/batch/batch.service';
import { Role, UrlConstant } from 'src/app/service/constant';
import { FeedbackGroupService } from 'src/app/service/feedback-group/feedbak-group.service';
import { FeedbackService, IFeedbackQuestionGroup } from 'src/app/service/feedback/feedback.service';
import { IFeedbackOptions, IFeedbackQuestion } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-batch-session-feedback',
  templateUrl: './batch-session-feedback.component.html',
  styleUrls: ['./batch-session-feedback.component.css']
})
export class BatchSessionFeedbackComponent implements OnInit {

  // foreign keys
  batchID: string
  courseID: string
  sessionID: string
  loginID: string
  batchName: string
  sessionName: string

  // session feedback
  sessionFeedbackForm: FormGroup
  feedbackGroupForm: FormGroup
  talentSessionFeedbackForm: FormGroup

  // feedback
  facultyFeedbacks: IFacultyBatchSessionFeedback[]
  facultySessionFeedback: IFacultyBatchSessionFeedback
  talentFeedbacks: ITalentBatchSessionFeedback[]
  talentSessionFeedback: ITalentBatchSessionFeedback

  // feedback questions/options
  feedbackGroupQuestions: IFeedbackQuestionGroup[]
  facultyFeedbackQuestions: IFeedbackQuestionGroup[]
  talentFeedbackQuestions: IFeedbackQuestion[]
  feedbackOptions: IFeedbackOptions[]
  feedbackResponse: any[]

  // modal
  modalRef: any;

  // talent
  talentID: string
  totalTalents: number

  // faculty
  facultyID: string
  totalFaculty: number

  // boolean
  deleteTalentFeedback: boolean

  // spinner


  // access
  permission: IPermission
  isFaculty: boolean
  isAdmin: boolean
  isTalent: boolean
  isSalesperson: boolean

  talentFeedbackAddClick: boolean
  facultyFeedbackAddClick: boolean

  feedbackTo: string

  // disable button
  disableButton: boolean

  // navigated
  feedbackTalentID: string
  isNavigatedFromTalent: boolean

  // constants
  INITIAL_SCORE = 0
  MINIMUM_SCORE = 7
  MAX_SCORE = 10

  // score
  facultyAverageScore: number

  // spinner


  constructor(
    private formBuilder: FormBuilder,
    private route: ActivatedRoute,
    private batchService: BatchService,
    private localService: LocalService,
    private urlConstant: UrlConstant,
    public utilService: UtilityService,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private _location: Location,
    private feedbackService: FeedbackService,
    private role: Role,
    private feedbackGroupService: FeedbackGroupService
  ) {
    this.initializeVariables()
    this.getAllComponents()
  }

  initializeVariables(): void {
    if (this.localService.getJsonValue("roleName") == this.role.ADMIN || this.localService.getJsonValue("roleName") == this.role.SALES_PERSON) {
      this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_BATCH_MASTER_SESSION_FEEDBACK)
    }
    if (this.localService.getJsonValue("roleName") == this.role.FACULTY) {
      this.permission = this.utilService.getPermission(this.urlConstant.MY_BATCH_SESSION_FEEDBACK)
    }

    this.loginID = this.localService.getJsonValue("loginID")
    this.isFaculty = (this.localService.getJsonValue("roleName") == this.role.FACULTY ? true : false)
    this.isAdmin = (this.localService.getJsonValue("roleName") == this.role.ADMIN ? true : false)
    this.isTalent = (this.localService.getJsonValue("roleName") == this.role.TALENT ? true : false)
    this.isSalesperson = (this.localService.getJsonValue("roleName") == this.role.SALES_PERSON ? true : false)

    this.spinnerService.loadingMessage = "Getting feedback"

    this.facultyFeedbacks = []
    this.feedbackGroupQuestions = []
    this.facultyFeedbackQuestions = []
    this.talentFeedbackQuestions = []


    this.totalTalents = 0
    this.totalFaculty = 0
    this.facultyAverageScore = 0

    this.deleteTalentFeedback = false
    this.disableButton = false
    this.talentFeedbackAddClick = false
    this.facultyFeedbackAddClick = false
    this.isNavigatedFromTalent = false
  }

  getAllComponents(): void {
    this.extractID()

    // if (this.isFaculty) {
    //   this.getFeedbackQuestionsForFaculty("Faculty_Session_Feedback")
    // }
    if (this.isTalent) {
      this.getFeedbackQuestionsForTalent("Talent_Session_Feedback")
    }
    if (this.isAdmin || this.isSalesperson || this.isFaculty) {
      this.getFeedbackQuestionsForFaculty("Faculty_Session_Feedback")
      this.getFeedbackQuestionsForTalent("Talent_Session_Feedback")
    }
    this.getAllFacultyBatchSessionFeedback()
    this.getAllTalentBatchSessionFeedback()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  extractID(): void {
    this.route.queryParamMap.subscribe(
      params => {
        this.batchID = params.get('batchID')
        this.courseID = params.get('courseID')
        this.sessionID = params.get('sessionID')
        this.sessionName = params.get('sessionName')
        this.batchName = params.get('batchName')
        this.feedbackTalentID = params.get("talentID")

        if (this.feedbackTalentID) {
          this.isNavigatedFromTalent = true
        }

      })
  }

  // ============================================================= FACULTY FEEDBACK =============================================================
  createFeedbackGroupForm(): void {
    this.feedbackGroupForm = this.formBuilder.group({
      feedbackGroups: new FormArray([])
    })
  }

  get feedbackGroupArray(): FormArray {
    return this.feedbackGroupForm.get("feedbackGroups") as FormArray
  }

  getFeedbackQuestions(index: number): FormArray {
    return this.feedbackGroupArray.at(index).get('feedbackQuestions') as FormArray
  }

  addFeedbackGroup(): void {
    this.feedbackGroupArray.push(this.formBuilder.group({
      id: new FormControl(null),
      groupName: new FormControl(null),
      groupDescription: new FormControl(null),
      answer: new FormControl(null, [Validators.required, Validators.min(1)]),
      feedbackQuestions: new FormArray([]),
    }))
  }

  createFeedbackSessionForm(): void {
    this.sessionFeedbackForm = this.formBuilder.group({
      batchID: new FormControl(this.batchID),
      talentID: new FormControl(null, [Validators.required]),
      facultyID: new FormControl(this.loginID),
      questionID: new FormControl(null, [Validators.required]),
      optionID: new FormControl(null, [Validators.required]),
      answer: new FormControl(null, [Validators.required]),
      isVisible: new FormControl(false),
    })
  }

  createFeedbackQuestionsGroupForm(): void {
    if (this.feedbackGroupQuestions) {
      for (let index = 0; index < this.feedbackGroupQuestions.length; index++) {
        this.addFeedbackGroup()
        this.feedbackGroupArray.at(index).get("groupName").setValue(this.feedbackGroupQuestions[index].groupName)
        this.feedbackGroupArray.at(index).get("groupDescription").setValue(this.feedbackGroupQuestions[index].groupDescription)
        this.feedbackGroupQuestions[index].answer = this.INITIAL_SCORE
        this.createFeedbackQuestionForm(this.feedbackGroupQuestions[index].feedbackQuestions,
          index, this.feedbackGroupQuestions[index].answer)
      }
    }
    // console.log(this.feedbackGroupArray.controls);
  }


  onFacultyFeedbackAddClick(sessionFeedback: IFacultyBatchSessionFeedback, feedbackModal: any): void {
    sessionFeedback.showFeedback = !sessionFeedback.showFeedback
    this.facultyFeedbackAddClick = true
    this.facultySessionFeedback = sessionFeedback
    this.feedbackGroupQuestions = this.facultyFeedbackQuestions
    this.feedbackResponse = []
    // console.log(this.facultySessionFeedback);
    this.feedbackTo = sessionFeedback.talent.firstName + " " + sessionFeedback.talent.lastName
    if (this.isAdmin) {
      this.facultyID = sessionFeedback.faculty.id
      this.talentID = sessionFeedback.talent.id
    }
    this.createFeedbackGroupForm()
    this.createFeedbackQuestionsGroupForm()
    this.openModal(feedbackModal, "lg")
  }

  addFeedbackQuestionsForGroup(event: any, index: number) {
    this.feedbackGroupQuestions[index].answer = event.target.value
    this.feedbackGroupArray.at(index).get("answer").setValue(event.target.value)
    // score is greater than 7
    if (event.target.value > this.MINIMUM_SCORE) {
      this.setFeedbackQuestionForm(this.feedbackGroupQuestions[index], this.feedbackGroupArray.at(index).get('feedbackQuestions'))
      return
    }
    // score is less than 7
    this.resetFeedbackQuestionForm(this.feedbackGroupQuestions[index], this.feedbackGroupArray.at(index).get('feedbackQuestions'))
  }

  setFeedbackQuestionForm(feedbackQuestionGroup: IFeedbackQuestionGroup, feedbackQuestionControl: any): void {
    for (let j = 0; j < feedbackQuestionControl.controls.length; j++) {
      feedbackQuestionControl.at(j).get('isVisible').setValue(true)
      feedbackQuestionControl.at(j).get('answer').setValue((this.INITIAL_SCORE).toString())
      feedbackQuestionControl.at(j).get('optionID').setValue(null)

      if (feedbackQuestionGroup.answer > this.INITIAL_SCORE && feedbackQuestionGroup.feedbackQuestions[j].hasOptions) {
        for (let k = 0; k < feedbackQuestionGroup.feedbackQuestions[j].options.length; k++) {

          if (feedbackQuestionGroup.feedbackQuestions[j].options[k].key == feedbackQuestionGroup.answer) {
            feedbackQuestionControl.at(j).get('optionID').setValue(feedbackQuestionGroup.feedbackQuestions[j].options[k].id)
            feedbackQuestionControl.at(j).get('isVisible').setValue(false)
            feedbackQuestionControl.at(j).get('answer').setValue((feedbackQuestionGroup.answer).toString())
          }
        }
      }
    }
    // console.log(feedbackQuestionControl);
  }

  resetFeedbackQuestionForm(feedbackQuestionGroup: IFeedbackQuestionGroup, feedbackQuestionControl: any): void {
    for (let j = 0; j < feedbackQuestionControl.controls.length; j++) {
      feedbackQuestionControl.at(j).get('isVisible').setValue(true)
      // feedbackQuestionControl.at(j).get('optionID').setValue(null)
      // feedbackQuestionControl.at(j).get('answer').setValue((this.INITIAL_SCORE).toString())

      if (feedbackQuestionGroup.feedbackQuestions[j].hasOptions) {

        for (let k = 0; k < feedbackQuestionGroup.feedbackQuestions[j].options.length; k++) {
          if (feedbackQuestionGroup.feedbackQuestions[j].options[k].key == feedbackQuestionGroup.answer) {
            feedbackQuestionControl.at(j).get('optionID').setValue(feedbackQuestionGroup.feedbackQuestions[j].options[k].id)
            feedbackQuestionControl.at(j).get('answer').setValue(feedbackQuestionGroup.feedbackQuestions[j].options[k].value)
          }
        }
      }
    }
    // console.log(feedbackQuestionControl);
  }

  createFeedbackQuestionForm(feedbackQuestions: IFeedbackQuestion[],
    feedbackIndex: number, answer?: number): void {
    if (feedbackQuestions) {
      for (let index = 0; index < feedbackQuestions.length; index++) {
        this.createFeedbackSessionForm()

        this.sessionFeedbackForm.get('questionID').setValue(feedbackQuestions[index].id)
        this.sessionFeedbackForm.get('isVisible').setValue(true)

        if (feedbackQuestions[index].hasOptions) {
          this.sessionFeedbackForm.get('optionID').setValidators([Validators.required])
          this.sessionFeedbackForm.get('answer').setValidators(null)
          if (answer) {
            for (let j = 0; j < feedbackQuestions[index].options.length; j++) {
              if (feedbackQuestions[index].options[j].key == answer) {
                this.sessionFeedbackForm.get('optionID').setValue(feedbackQuestions[index].options[j].id)
                this.sessionFeedbackForm.get('answer').setValue(answer.toString())
                this.sessionFeedbackForm.get('isVisible').setValue(false)
              }
            }
          }
        } else {
          this.sessionFeedbackForm.get('optionID').setValidators(null)
          this.sessionFeedbackForm.get('answer').setValidators([Validators.required])
          this.sessionFeedbackForm.get('isVisible').setValue(true)
        }

        if (this.isFaculty) {
          this.sessionFeedbackForm.get('facultyID').setValue(this.loginID)
          this.sessionFeedbackForm.get('talentID').setValue(this.facultySessionFeedback.talent.id)
        }
        if (this.isTalent) {
          this.sessionFeedbackForm.get('facultyID').setValue(this.talentSessionFeedback.faculty.id)
          this.sessionFeedbackForm.get('talentID').setValue(this.loginID)
        }
        if (this.isAdmin) {
          this.sessionFeedbackForm.get('facultyID').setValue(this.facultyID)
          this.sessionFeedbackForm.get('talentID').setValue(this.talentID)
        }

        (this.feedbackGroupArray.at(feedbackIndex).get("feedbackQuestions") as FormArray).push(this.sessionFeedbackForm)
      }
    }
  }

  onFacutyFeedbackInput(event: any, feedbackQuestionControl: any, options: IFeedbackOptions[]): void {
    // console.log(feedbackQuestionControl)

    feedbackQuestionControl.get('answer').setValue(event.target.value)
    feedbackQuestionControl.get('optionID').setValue(null)

    for (let index = 0; index < options.length; index++) {
      if (event.target.value == options[index].key) {
        // feedbackQuestionControl.get('answer').setValue((options[index].key).toString())
        feedbackQuestionControl.get('optionID').setValue(options[index].id)
      }
    }
    // console.log(feedbackQuestionControl)
  }

  onFacultyFeedbackChange(feedbackQuestionControl: any, feedbackOptions: IFeedbackOptions[]) {
    let optionID = feedbackQuestionControl.get("optionID").value
    for (let index = 0; index < feedbackOptions.length; index++) {
      if (optionID == feedbackOptions[index].id) {
        feedbackQuestionControl.get('answer').setValue(feedbackOptions[index].value)
      }
    }
  }

  isFeedbackOptionsAvailable(feedbackQuestions: IFeedbackQuestion[], value: number): boolean {

    let isOptionAvailabe = false
    for (let index = 0; index < feedbackQuestions.length; index++) {
      if (!feedbackQuestions[index].hasOptions) {
        return true
      } else {
        isOptionAvailabe = false
        for (let j = 0; j < feedbackQuestions[index].options.length; j++) {
          if (value == feedbackQuestions[index].options[j].key) {
            isOptionAvailabe = true
            break
          }
        }
        if (!isOptionAvailabe) {
          return true
        }
      }
    }
    return false
  }

  onFacultyFeedbackDeleteClick(sessionFeedback: IFacultyBatchSessionFeedback, deleteModal: any): void {
    sessionFeedback.showFeedback = !sessionFeedback.showFeedback
    this.facultySessionFeedback = sessionFeedback
    this.deleteTalentFeedback = false
    this.openModal(deleteModal, "md")
  }

  validateFacultyFeedback() {
    // console.log(this.feedbackGroupForm.controls);
    // console.log(this.feedbackGroupForm.value);

    if (this.feedbackGroupForm.invalid) {
      this.feedbackGroupForm.markAllAsTouched()
      return
    }
    this.feedbackResponse = []

    for (let index = 0; index < this.feedbackGroupArray.controls.length; index++) {
      this.feedbackResponse = this.feedbackResponse.concat(this.feedbackGroupArray.at(index).get('feedbackQuestions').value)
    }
    // console.log(this.feedbackResponse)

    if (this.facultyFeedbackAddClick) {
      this.addFacultyFeedback()
      this.facultyFeedbackAddClick = false
      return
    }

  }


  // ============================================================= TALENT FEEDBACK =============================================================

  createSessionFeedbackForm(): void {
    this.talentSessionFeedbackForm = this.formBuilder.group({
      feedbacks: this.formBuilder.array([])
    })
  }

  get talentSessionFeedbacksArray() {
    return this.talentSessionFeedbackForm.get('feedbacks') as FormArray
  }

  addBatchFeedback(): void {
    this.talentSessionFeedbacksArray.push(this.formBuilder.group({
      batchID: new FormControl(this.batchID),
      talentID: new FormControl(null, [Validators.required]),
      facultyID: new FormControl(this.loginID),
      questionID: new FormControl(null, [Validators.required]),
      optionID: new FormControl(null, [Validators.required]),
      answer: new FormControl(null, [Validators.required])
    }))
  }

  onTalentFeedbackAddClick(sessionFeedback: ITalentBatchSessionFeedback, feedbackModal: any,
    talentFeedback?: ITalentBatchSessionFeebackDTO): void {
    this.talentFeedbackAddClick = true
    this.talentSessionFeedback = sessionFeedback
    // this.feedbackGroupQuestions = this.talentFeedbackQuestions
    this.feedbackResponse = []
    // console.log(this.talentSessionFeedback);
    this.feedbackTo = sessionFeedback.faculty.firstName + " " + sessionFeedback.faculty.lastName
    if (this.isTalent) {
      sessionFeedback.showFeedback = !sessionFeedback.showFeedback
    }
    if (this.isAdmin) {
      this.facultyID = sessionFeedback.faculty.id
      this.talentID = talentFeedback.talent.id
      this.toggleFeedback(talentFeedback)
    }

    this.createSessionFeedbackForm()
    this.createBatchSessionForm()
    this.openModal(feedbackModal, "lg")
  }

  createBatchSessionForm(): void {
    if (this.talentFeedbackQuestions) {
      for (let index = 0; index < this.talentFeedbackQuestions.length; index++) {
        this.addBatchFeedback()

        this.talentSessionFeedbacksArray.at(index).get('questionID').setValue(this.talentFeedbackQuestions[index].id)

        if (this.talentFeedbackQuestions[index].hasOptions) {
          this.talentSessionFeedbacksArray.at(index).get('optionID').setValidators([Validators.required])
          this.talentSessionFeedbacksArray.at(index).get('answer').setValidators(null)
        } else {
          this.talentSessionFeedbacksArray.at(index).get('optionID').setValidators(null)
          this.talentSessionFeedbacksArray.at(index).get('answer').setValidators([Validators.required])
        }

        if (this.isFaculty) {
          this.talentSessionFeedbacksArray.at(index).get('facultyID').setValue(this.loginID)
          this.talentSessionFeedbacksArray.at(index).get('talentID').setValue(this.facultySessionFeedback.talent.id)
        }
        if (this.isTalent) {
          this.talentSessionFeedbacksArray.at(index).get('facultyID').setValue(this.talentSessionFeedback.faculty.id)
          this.talentSessionFeedbacksArray.at(index).get('talentID').setValue(this.loginID)
        }
        if (this.isAdmin) {
          this.talentSessionFeedbacksArray.at(index).get('facultyID').setValue(this.facultyID)
          this.talentSessionFeedbacksArray.at(index).get('talentID').setValue(this.talentID)
        }
      }
    }
  }

  onTalentFeedbackInput(event: any, sessionFeedbackControl: any, feedbackOptions: IFeedbackOptions[]) {
    sessionFeedbackControl.get('answer').setValue(event.target.value)
    sessionFeedbackControl.get('optionID').setValue(null)

    for (let index = 0; index < feedbackOptions.length; index++) {
      if (event.target.value == feedbackOptions[index].key) {
        // sessionFeedbackControl.get('answer').setValue((feedbackOptions[index].key).toString())
        sessionFeedbackControl.get('optionID').setValue(feedbackOptions[index].id)
      }
    }
    // console.log(this.talentSessionFeedbackForm.controls)
  }

  onTalentFeedbackChange(sessionFeedbackControl: any, feedbackOptions: IFeedbackOptions[]) {
    let optionID = sessionFeedbackControl.get("optionID").value
    for (let index = 0; index < feedbackOptions.length; index++) {
      if (optionID == feedbackOptions[index].id) {
        sessionFeedbackControl.get('answer').setValue(feedbackOptions[index].value)
      }
    }
  }

  onTalentFeedbackDeleteClick(sessionFeedback: ITalentBatchSessionFeedback,
    talentFeedback: ITalentBatchSessionFeebackDTO, deleteModal: any): void {
    // this.talentFeedback = feedback
    if (this.isAdmin) {
      this.toggleFeedback(talentFeedback)
    }
    this.facultyID = sessionFeedback.faculty.id
    this.talentID = talentFeedback.talent.id
    this.deleteTalentFeedback = true
    this.openModal(deleteModal, "md")
  }

  validateTalentFeedback(): void {

    // console.log(this.talentSessionFeedbackForm.controls)
    // console.log(this.talentSessionFeedbackForm.value)

    if (this.talentSessionFeedbackForm.invalid) {
      this.talentSessionFeedbackForm.markAllAsTouched()
      return
    }

    if (this.talentFeedbackAddClick) {
      this.addTalentFeedback()
      this.talentFeedbackAddClick = false
      return
    }
  }

  toggleFeedback(feedback: ITalentBatchSessionFeebackDTO): void {
    feedback.showFeedback = !feedback.showFeedback
  }

  // =============================================================CRUD=============================================================

  getAllFacultyBatchSessionFeedback(): void {

    if (this.isFaculty) {
      this.getFacultyBatchSessionFeedback()
      // this.getAllTalentBatchSessionFeedback()
      return
    }
    if (this.isTalent) {
      return
    }

    let params = new HttpParams()
    if (this.isNavigatedFromTalent) {
      params = params.append("talentID", this.feedbackTalentID)
    }

    this.totalTalents = 1
    this.disableButton = true
    this.spinnerService.loadingMessage = "Getting feedback"


    this.batchService.getAllFacultyBatchSessionFeedback(this.batchID, this.sessionID, params).
      subscribe((response: any) => {
        this.disableButton = false
        this.facultyFeedbacks = response.body
        this.totalTalents = this.facultyFeedbacks.length
        // console.log(this.facultyFeedbacks);
        this.calculateFacultyScore()
        // this.assignSessionFeedback()
      },
        (err: any) => {
          this.disableButton = false
          this.totalTalents = 0
          console.error(err)
          if (err.statusText.includes('Unknown')) {
            alert("No connection to server. Check internet.")
            return
          }
          alert(err.error.error)
        })
  }

  getFacultyBatchSessionFeedback(): void {
    this.disableButton = true
    this.totalTalents = 1
    this.spinnerService.loadingMessage = "Getting feedback"



    let params = new HttpParams()
    if (this.isNavigatedFromTalent) {
      params = params.append("talentID", this.feedbackTalentID)
    }

    this.batchService.getFacultyBatchSessionFeedback(this.batchID, this.sessionID, this.loginID, params).
      subscribe((response: any) => {
        this.disableButton = false
        this.facultyFeedbacks = response.body
        // console.log(this.facultyFeedbacks);
        this.totalTalents = this.facultyFeedbacks.length
        this.calculateFacultyScore()
        // this.assignSessionFeedback()
      },
        (err: any) => {
          this.disableButton = false
          this.totalTalents = 0
          console.error(err)
          if (err.statusText.includes('Unknown')) {
            alert("No connection to server. Check internet.")
            return
          }
          alert(err.error.error)
        })
  }

  getAllTalentBatchSessionFeedback(): void {

    if (this.isTalent) {
      this.getTalentBatchSessionFeedback()
      return
    }
    this.totalFaculty = 1
    this.disableButton = true
    this.spinnerService.loadingMessage = "Getting all feedback"



    let params = new HttpParams()
    if (this.isNavigatedFromTalent) {
      params = params.append("talentID", this.feedbackTalentID)
    }

    this.batchService.getTalentBatchSessionFeedback(this.batchID, this.sessionID, params).subscribe((response: any) => {
      this.disableButton = false
      this.talentFeedbacks = response.body
      // console.log(this.talentFeedbacks);
      this.totalFaculty = this.talentFeedbacks.length
      this.calculateTalentScore()
    },
      (err: any) => {
        this.disableButton = false
        this.totalFaculty = 0
        console.error(err)
        if (err.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
          return
        }
        alert(err.error.error)
      })
  }


  getTalentBatchSessionFeedback(): void {
    this.disableButton = true
    this.totalFaculty = 1
    this.spinnerService.loadingMessage = "Getting all feedback"


    this.batchService.getSpecifiedTalentBatchSessionFeedback(this.batchID, this.sessionID, this.loginID).subscribe((response: any) => {
      // console.log(response.body);
      this.disableButton = false
      this.talentFeedbacks = response.body
      this.totalFaculty = this.talentFeedbacks.length
      this.calculateTalentScore()
      // this.getFacultyTalentBatchSessionFeedback()
    },
      (err: any) => {
        this.disableButton = false
        this.totalFaculty = 0
        console.error(err)
        if (err.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
          return
        }
        alert(err.error.error)
      })
  }

  addFacultyFeedback(): void {
    // console.log(this.sessionFeedbackForm.value);
    this.disableButton = true
    this.spinnerService.loadingMessage = "Adding session feedback"

    // console.log(sessionID)
    this.batchService.addFacultySessionFeedbacks(this.batchID, this.sessionID, this.feedbackResponse).
      subscribe((response: any) => {
        // console.log(response)
        this.disableButton = false
        this.modalRef.close()
        this.getAllFacultyBatchSessionFeedback()
        alert("Feedback successfully added")
        this.feedbackResponse = []
      }, (err: any) => {
        this.disableButton = false
        console.error(err)
        if (err.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
          return
        }
        alert(err.error.error)
      })
  }

  addTalentFeedback(): void {
    // console.log(this.talentSessionFeedbackForm.value);
    this.disableButton = true
    this.spinnerService.loadingMessage = "Adding session feedback"


    this.batchService.addTalentSessionFeedbacks(this.batchID, this.sessionID, this.talentSessionFeedbackForm.value.feedbacks).
      subscribe((response: any) => {
        // console.log(response)
        this.disableButton = false
        this.modalRef.close()
        this.getAllTalentBatchSessionFeedback()
        alert("Feedback successfully added")
      }, (err: any) => {
        this.disableButton = false
        console.error(err)
        if (err.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
          return
        }
        alert(err.error.error)
      })
  }

  deleteFacultyBatchSessionFeedback(): void {
    if (this.deleteTalentFeedback) {
      this.deleteTalentBatchSessionFeedback()
      return
    }
    this.disableButton = true
    this.spinnerService.loadingMessage = "Deleting session feedback of talent"


    this.batchService.deleteFacultyBatchSessionFeedback(this.batchID, this.sessionID, this.facultySessionFeedback.talent.id).subscribe((response: any) => {
      // console.log(response);
      this.disableButton = false
      alert("Feedback successfully deleted")
      this.getAllFacultyBatchSessionFeedback()
      this.modalRef.close()
    }, (err: any) => {
      this.disableButton = false
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  deleteTalentBatchSessionFeedback(): void {
    this.disableButton = true
    this.spinnerService.loadingMessage = "Deleting session feedback of faculty"


    this.batchService.deleteTalentBatchSessionFeedback(this.batchID, this.sessionID, this.talentID).
      subscribe((response: any) => {
        this.disableButton = false
        // console.log(response)
        alert("Feedback successfully deleted")
        this.getAllTalentBatchSessionFeedback()
        this.modalRef.close()
      }, (err: any) => {
        this.disableButton = false
        console.error(err)
        if (err.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
          return
        }
        alert(err.error.error)
      })
  }

  // =============================================================COMPONENTS=============================================================

  getFeedbackQuestionsForFaculty(questionType: string): void {
    this.disableButton = true


    // this.generalService.getFeedbackQuestionByType(questionType).subscribe((response: any) => {
    this.feedbackGroupService.getFeedbackQuestionGroupByType(questionType).subscribe((response: any) => {
      this.facultyFeedbackQuestions = response.body
      this.disableButton = false
      // console.log(response.body);
    }, (err: any) => {
      this.disableButton = false
      console.error(err)
    })
  }

  getFeedbackQuestionsForTalent(questionType: string): void {
    this.disableButton = true


    this.feedbackService.getFeedbackQuestionByType(questionType).subscribe((response: any) => {
      // this.feedbackGroupSerivce.getFeedbackQuestionGroupByType(questionType).subscribe((response: any) => {
      this.talentFeedbackQuestions = response.body
      this.disableButton = false
      // console.log(response.body);
    }, (err: any) => {
      this.disableButton = false
      console.error(err)
    })
  }



  openModal(contentModal: any, modalSize?: string): void {
    if (modalSize == undefined || modalSize == "") {
      modalSize = "xl"
    }
    this.modalRef = this.modalService.open(contentModal, {
      ariaLabelledBy: 'modal-basic-title',
      backdrop: 'static',
      size: modalSize
    })
  }

  backToPreviousPage(): void {
    this._location.back()
  }

  calculateTalentScore(): void {
    for (let index = 0; index < this.talentFeedbacks.length; index++) {
      let averageScore = 0
      let noOfFeedbacks = 0
      this.talentFeedbacks[index].showFeedback = true
      this.talentFeedbacks[index].averageScore = 0

      for (let j = 0; j < this.talentFeedbacks[index].feedbacks?.length; j++) {
        this.talentFeedbacks[index].feedbacks[j].showFeedback = false
        if (this.talentFeedbacks[index].feedbacks[j].sessionFeedbacks?.length > 0) {
          noOfFeedbacks++
          this.talentFeedbacks[index].feedbacks[j].averageScore = this.calculateTalentFeedbackScore(this.talentFeedbacks[index].feedbacks[j])
          averageScore += this.talentFeedbacks[index].feedbacks[j].averageScore
        }
      }

      if (noOfFeedbacks == 0) {
        this.talentFeedbacks[index].averageScore = 0
        return
      }
      this.talentFeedbacks[index].averageScore = averageScore / noOfFeedbacks
      this.facultyAverageScore = this.talentFeedbacks[index].averageScore
    }
  }

  // will returns avergaeScore of every talent's feedback to faculty
  calculateTalentFeedbackScore(talentFeedback: ITalentBatchSessionFeedback): number {
    let averageScore = 0
    let totalScore = 0

    if (talentFeedback.sessionFeedbacks.length == 0) {
      return 0
    }

    for (let index = 0; index < talentFeedback.sessionFeedbacks.length; index++) {
      if (talentFeedback.sessionFeedbacks[index].option && talentFeedback.sessionFeedbacks[index].question) {
        averageScore += talentFeedback.sessionFeedbacks[index].option?.key
        totalScore += talentFeedback.sessionFeedbacks[index].question.maxScore
      }
    }

    averageScore = (averageScore * 10) / totalScore
    return averageScore
  }

  calculateFacultyScore(): void {
    for (let index = 0; index < this.facultyFeedbacks.length; index++) {
      this.facultyFeedbacks[index].showFeedback = false
      this.facultyFeedbacks[index].averageScore = this.calculateFacultyFeedbackScore(this.facultyFeedbacks[index])
    }
  }

  calculateFacultyFeedbackScore(facultyFeeback: IFacultyBatchSessionFeedback): number {
    let averageScore = 0
    let totalScore = 0

    if (facultyFeeback.sessionFeedbacks?.length == 0) {
      return 0
    }

    for (let index = 0; index < facultyFeeback.sessionFeedbacks?.length; index++) {
      if (facultyFeeback.sessionFeedbacks[index].option && facultyFeeback.sessionFeedbacks[index].question) {
        averageScore += facultyFeeback.sessionFeedbacks[index].option?.key
        totalScore += facultyFeeback.sessionFeedbacks[index].question.maxScore
      }
    }

    averageScore = (averageScore * 10) / totalScore
    return averageScore
  }

  counter() {
    return new Array(this.MAX_SCORE + 1);
  }

}
