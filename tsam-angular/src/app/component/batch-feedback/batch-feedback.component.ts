import { Location } from '@angular/common';
import { HttpParams } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BatchService, IBatchFeedbackDTO, IFacultyBatchFeedback, ITalentBatchFeedback } from 'src/app/service/batch/batch.service';
import { Role, UrlConstant } from 'src/app/service/constant';
import { FeedbackGroupService } from 'src/app/service/feedback-group/feedbak-group.service';
import { FeedbackService, IFeedbackQuestionGroup, IFeedbackOptions, IFeedbackQuestion } from 'src/app/service/feedback/feedback.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-batch-feedback',
  templateUrl: './batch-feedback.component.html',
  styleUrls: ['./batch-feedback.component.css']
})
export class BatchFeedbackComponent implements OnInit {

  // foreign keys
  batchID: string
  loginID: string
  feedbackTalentID: string

  batchName: string

  // feedback form
  batchFeedbackForm: FormGroup
  feedbackGroupForm: FormGroup
  talentBatchFeedbackForm: FormGroup

  // feedback
  facultyFeedback: IFacultyBatchFeedback
  facultyFeedbacks: IFacultyBatchFeedback[]
  talentFeedback: ITalentBatchFeedback
  talentFeedbacks: ITalentBatchFeedback[]
  // feedbackTalents: ITalentFeedback[]
  // batchFeedback: ITalentBatchFeedback

  // feedback questions/options
  feedbackGroupQuestions: IFeedbackQuestionGroup[]
  talentFeedbackQuestions: IFeedbackQuestion[]
  facultyFeedbackQuestions: IFeedbackQuestionGroup[]
  feedbackOptions: IFeedbackOptions[]
  feedbackResponse: any[]

  // boolean
  deleteTalentFeedback: boolean

  // modal
  modalRef: any;

  // faculty
  facultyID: string

  // talent
  talentID: string
  totalStudents: number
  totalFaculty: number

  talentFeedbackAddClick: boolean
  facultyFeedbackAddClick: boolean

  // spinner


  // access
  permission: IPermission
  isFaculty: boolean
  isAdmin: boolean
  isTalent: boolean
  isSalesperson: boolean

  // disable button
  disableButton: boolean

  feedbackTo: string

  // navigated
  isNavigatedFromTalent: boolean

  // constant
  INITIAL_SCORE = 0
  MINIMUM_SCORE = 7
  MAX_SCORE = 10

  // faculty average score
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
			this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_BATCH_MASTER)
		}

		if (this.localService.getJsonValue("roleName") == this.role.FACULTY) {
			this.permission = this.utilService.getPermission(this.urlConstant.MY_BATCH)

		}
    this.loginID = this.localService.getJsonValue("loginID")

    this.isFaculty = (this.localService.getJsonValue("roleName") == this.role.FACULTY ? true : false)
    this.isAdmin = (this.localService.getJsonValue("roleName") == this.role.ADMIN ? true : false)
    this.isTalent = (this.localService.getJsonValue("roleName") == this.role.TALENT ? true : false)
    this.isSalesperson = (this.localService.getJsonValue("roleName") == this.role.SALES_PERSON ? true : false)

    this.spinnerService.loadingMessage = "Getting feedback"

    // this.feedbackTalents = []
    this.facultyFeedbacks = []
    this.talentFeedbacks = []
    this.feedbackGroupQuestions = []
    this.talentFeedbackQuestions = []
    this.facultyFeedbackQuestions = []
    this.feedbackResponse = []


    this.totalStudents = 0
    this.totalFaculty = 0

    this.deleteTalentFeedback = false
    this.disableButton = false
    this.talentFeedbackAddClick = false
    this.facultyFeedbackAddClick = false
    this.isNavigatedFromTalent = false
  }

  getAllComponents(): void {
    this.extractID()

    if (this.isTalent) {
      this.getFeedbackQuestionsForTalent("Talent_Batch_Feedback")
    }
    if (this.isAdmin || this.isSalesperson || this.isFaculty) {
      this.getFeedbackQuestionsForFaculty("Faculty_Batch_Feedback")
      this.getFeedbackQuestionsForTalent("Talent_Batch_Feedback")
    }
    this.getAllTalentBatchFeedback()
    this.getAllFacultyBatchFeedback()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  extractID(): void {
    this.route.queryParamMap.subscribe(
      params => {
        this.batchID = params.get("batchID")
        this.batchName = params.get("batchName")
        this.feedbackTalentID = params.get("talentID")

        if (this.feedbackTalentID) {
          this.isNavigatedFromTalent = true
        }

      })
  }

  backToPreviousPage(): void {
    this._location.back()
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

  createFeedbackForm(): void {
    this.batchFeedbackForm = this.formBuilder.group({
      batchID: new FormControl(this.batchID),
      talentID: new FormControl(null, [Validators.required]),
      facultyID: new FormControl(this.loginID),
      questionID: new FormControl(null, [Validators.required]),
      optionID: new FormControl(null, [Validators.required]),
      answer: new FormControl(null, [Validators.required, Validators.min(1)]),
      isVisible: new FormControl(false),
    })
  }

  onFacultyFeedbackAddClick(feedback: IFacultyBatchFeedback, feedbackModal: any): void {
    feedback.showFeedback = !feedback.showFeedback
    this.facultyFeedbackAddClick = true
    this.feedbackResponse = []
    this.facultyFeedback = feedback
    this.feedbackGroupQuestions = this.facultyFeedbackQuestions
    this.feedbackTo = feedback.talent.firstName + " " + feedback.talent.lastName
    // console.log(this.facultyFeedback);
    if (this.isAdmin) {
      this.talentID = feedback.talent.id
      this.facultyID = feedback.faculty.id
    }
    this.createFeedbackGroupForm()
    this.createFeedbackQuestionsGroupForm()
    this.openModal(feedbackModal, "lg")
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

    // console.log(this.feedbackGroupArray.controls);
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
            feedbackQuestionControl.at(j).get('answer').setValue((feedbackQuestionGroup.answer).toString())
            feedbackQuestionControl.at(j).get('isVisible').setValue(false)
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
      // feedbackQuestionControl.at(j).get('answer').setValue(this.INITIAL_SCORE)

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
        this.createFeedbackForm()

        this.batchFeedbackForm.get('questionID').setValue(feedbackQuestions[index].id)
        // this.batchFeedbackForm.get('question').setValue(feedbackQuestions[index].question)
        this.batchFeedbackForm.get('isVisible').setValue(true)

        if (feedbackQuestions[index].hasOptions) {
          this.batchFeedbackForm.get('optionID').setValidators([Validators.required])
          this.batchFeedbackForm.get('answer').setValidators(null)
          if (answer) {
            for (let j = 0; j < feedbackQuestions[index].options.length; j++) {
              if (feedbackQuestions[index].options[j].key == answer) {
                this.batchFeedbackForm.get('optionID').setValue(feedbackQuestions[index].options[j].id)
                this.batchFeedbackForm.get('answer').setValue(answer.toString())
                this.batchFeedbackForm.get('isVisible').setValue(false)
              }
            }
          }
        } else {
          this.batchFeedbackForm.get('optionID').setValidators(null)
          this.batchFeedbackForm.get('answer').setValidators([Validators.required])
          this.batchFeedbackForm.get('isVisible').setValue(true)
        }

        if (this.isFaculty) {
          this.batchFeedbackForm.get('facultyID').setValue(this.loginID)
          this.batchFeedbackForm.get('talentID').setValue(this.facultyFeedback.talent.id)
        }
        if (this.isTalent) {
          this.batchFeedbackForm.get('facultyID').setValue(this.talentFeedback.faculty.id)
          this.batchFeedbackForm.get('talentID').setValue(this.loginID)
        }
        if (this.isAdmin) {
          this.batchFeedbackForm.get('facultyID').setValue(this.facultyID)
          this.batchFeedbackForm.get('talentID').setValue(this.talentID)
        }

        (this.feedbackGroupArray.at(feedbackIndex).get("feedbackQuestions") as FormArray).push(this.batchFeedbackForm)
      }
    }
  }

  onFacultyFeedbackInput(event: any, feedbackQuestionControl: any, options: IFeedbackOptions[]): void {
    // console.log(feedbackQuestionControl)

    // console.log("event -> ", event.target.value)
    feedbackQuestionControl.get('answer').setValue(event.target.value)
    feedbackQuestionControl.get('optionID').setValue(null)

    for (let index = 0; index < options.length; index++) {
      if (event.target.value == options[index].key) {
        // console.log(options[index])
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
          if (value == feedbackQuestions[index].options[j]?.key) {
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

  onFacultyFeedbackDeleteClick(feedback: IFacultyBatchFeedback, deleteModal: any): void {
    feedback.showFeedback = !feedback.showFeedback
    this.facultyFeedback = feedback
    this.deleteTalentFeedback = false
    this.openModal(deleteModal, "md")
  }

  validateFacultyFeedback() {
    // console.log(this.feedbackGroupArray.controls)

    if (this.feedbackGroupArray.invalid) {
      this.feedbackGroupArray.markAllAsTouched()
      return
    }
    this.feedbackResponse = []

    for (let index = 0; index < this.feedbackGroupArray.controls.length; index++) {
      this.feedbackResponse = this.feedbackResponse.concat(this.feedbackGroupArray.at(index).get('feedbackQuestions').value)
    }
    // console.log(this.feedbackResponse);

    if (this.facultyFeedbackAddClick) {
      this.addFacultyFeedback()
      this.facultyFeedbackAddClick = false
      return
    }
  }

  // ============================================================= TALENT FEEDBACK =============================================================

  createSessionFeedbackForm(): void {
    this.talentBatchFeedbackForm = this.formBuilder.group({
      feedbacks: this.formBuilder.array([])
    })
  }

  get talentBatchFeedbacksArray() {
    return this.talentBatchFeedbackForm.get('feedbacks') as FormArray
  }

  addBatchFeedback(): void {
    this.talentBatchFeedbacksArray.push(this.formBuilder.group({
      batchID: new FormControl(this.batchID),
      talentID: new FormControl(null, [Validators.required]),
      facultyID: new FormControl(this.loginID),
      questionID: new FormControl(null, [Validators.required]),
      optionID: new FormControl(null, [Validators.required]),
      answer: new FormControl(null, [Validators.required])
    }))
  }

  onTalentFeedbackAddClick(feedback: ITalentBatchFeedback, feedbackModal: any,
    talentFeedback?: IBatchFeedbackDTO): void {

    this.talentFeedbackAddClick = true
    this.feedbackResponse = []
    this.talentFeedback = feedback
    // this.feedbackGroupQuestions = this.talentFeedbackQuestions
    this.feedbackTo = feedback.faculty.firstName + " " + feedback.faculty.lastName
    // console.log(this.facultyFeedback);
    if (this.isTalent) {
      feedback.showFeedback = !feedback.showFeedback
    }
    if (this.isAdmin) {
      this.toggleFeedback(talentFeedback)
      this.facultyID = feedback.faculty.id
      this.talentID = talentFeedback.talent.id
    }
    this.createSessionFeedbackForm()
    this.createBatchSessionForm()
    this.openModal(feedbackModal, "lg")
  }

  createBatchSessionForm(): void {
    if (this.talentFeedbackQuestions) {
      for (let index = 0; index < this.talentFeedbackQuestions.length; index++) {
        this.addBatchFeedback()

        this.talentBatchFeedbacksArray.at(index).get('questionID').setValue(this.talentFeedbackQuestions[index].id)

        if (this.talentFeedbackQuestions[index].hasOptions) {
          this.talentBatchFeedbacksArray.at(index).get('optionID').setValidators([Validators.required])
          this.talentBatchFeedbacksArray.at(index).get('answer').setValidators(null)
        } else {
          this.talentBatchFeedbacksArray.at(index).get('optionID').setValidators(null)
          this.talentBatchFeedbacksArray.at(index).get('answer').setValidators([Validators.required])
        }

        if (this.isFaculty) {
          this.talentBatchFeedbacksArray.at(index).get('facultyID').setValue(this.loginID)
          this.talentBatchFeedbacksArray.at(index).get('talentID').setValue(this.facultyFeedback.talent.id)
        }
        if (this.isTalent) {
          this.talentBatchFeedbacksArray.at(index).get('facultyID').setValue(this.talentFeedback.faculty.id)
          this.talentBatchFeedbacksArray.at(index).get('talentID').setValue(this.loginID)
        }
        if (this.isAdmin) {
          this.talentBatchFeedbacksArray.at(index).get('facultyID').setValue(this.facultyID)
          this.talentBatchFeedbacksArray.at(index).get('talentID').setValue(this.talentID)
        }
      }
    }
  }

  onTalentFeedbackInput(event: any, batchFeedbackControl: any, feedbackOptions: IFeedbackOptions[]) {
    batchFeedbackControl.get('answer').setValue(event.target.value)
    batchFeedbackControl.get('optionID').setValue(null)

    for (let index = 0; index < feedbackOptions.length; index++) {
      if (event.target.value == feedbackOptions[index].key) {
        // batchFeedbackControl.get('answer').setValue((feedbackOptions[index].key).toString())
        batchFeedbackControl.get('optionID').setValue(feedbackOptions[index].id)
      }
    }
    // console.log(this.talentBatchFeedbackForm.controls)
  }

  onTalentFeedbackChange(batchFeedbackControl: any, feedbackOptions: IFeedbackOptions[]) {
    let optionID = batchFeedbackControl.get("optionID").value
    for (let index = 0; index < feedbackOptions.length; index++) {
      if (optionID == feedbackOptions[index].id) {
        batchFeedbackControl.get('answer').setValue(feedbackOptions[index].value)
      }
    }
  }

  onTalentFeedbackDeleteClick(feedback: ITalentBatchFeedback,
    talentFeedback: IBatchFeedbackDTO, deleteModal: any): void {
    // this.talentFeedback = feedback
    if (this.isAdmin) {
      this.toggleFeedback(talentFeedback)
    }
    this.facultyID = feedback.faculty.id
    this.talentID = talentFeedback.talent.id
    this.deleteTalentFeedback = true
    this.openModal(deleteModal, "md")
  }

  validateTalentFeedback(): void {

    // console.log(this.talentBatchFeedbackForm.controls)
    // console.log(this.talentBatchFeedbackForm.value)

    if (this.talentBatchFeedbackForm.invalid) {
      this.talentBatchFeedbackForm.markAllAsTouched()
      return
    }

    if (this.talentFeedbackAddClick) {
      this.addTalentFeedback()
      this.talentFeedbackAddClick = false
      return
    }
  }


  toggleFeedback(feedback: IBatchFeedbackDTO): void {
    feedback.showFeedback = !feedback.showFeedback
  }

  // =============================================================CRUD=============================================================

  // admin and faculty login
  getAllTalentBatchFeedback(): void {
    if (this.isTalent) {
      this.getTalentBatchFeeback()
      return
    }

    let params = new HttpParams()
    if (this.isNavigatedFromTalent) {
      params = params.append("talentID", this.feedbackTalentID)
    }

    this.totalFaculty = 1
    this.disableButton = true
    this.spinnerService.loadingMessage = "Getting feedback"


    this.batchService.getAllTalentBatchFeedback(this.batchID, params).subscribe((response: any) => {
      this.disableButton = false
      this.talentFeedbacks = response.body
      this.totalFaculty = this.talentFeedbacks.length
      // console.log(this.talentFeedbacks);
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

  // talent login
  getTalentBatchFeeback(): void {
    this.disableButton = true
    this.spinnerService.loadingMessage = "Getting feedback"

    this.totalFaculty = 1

    this.batchService.getTalentBatchFeedback(this.batchID, this.loginID).subscribe((response: any) => {
      this.disableButton = false
      this.talentFeedbacks = response.body
      // console.log(this.talentFeedbacks);
      this.totalFaculty = this.talentFeedbacks.length
      this.calculateTalentScore()
      // this.getFacultyTalentBatchFeedback()
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

  // admin login
  getAllFacultyBatchFeedback(): void {

    if (this.isFaculty) {
      this.getFacultyBatchFeedback()
      // this.getAllTalentBatchFeedback()
      return
    }
    if (this.isTalent) {
      return
    }
    this.totalStudents = 1
    this.disableButton = true
    this.spinnerService.loadingMessage = "Getting feedback"



    let params = new HttpParams()
    if (this.isNavigatedFromTalent) {
      params = params.append("talentID", this.feedbackTalentID)
    }

    this.batchService.getAllFacultyBatchFeedback(this.batchID, params).subscribe((response: any) => {
      this.disableButton = false
      this.facultyFeedbacks = response.body
      this.totalStudents = this.facultyFeedbacks.length
      this.assignFacultyFeedbackFields()
      // console.log(this.facultyFeedbacks);
    },
      (err: any) => {
        this.disableButton = false
        this.totalStudents = 0
        console.error(err)
        if (err.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
          return
        }
        alert(err.error.error)
      })
  }

  getFacultyBatchFeedback(): void {
    this.disableButton = true
    this.spinnerService.loadingMessage = "Getting feedback"
    this.totalStudents = 1



    let params = new HttpParams()
    if (this.isNavigatedFromTalent) {
      params = params.append("talentID", this.feedbackTalentID)
    }

    this.batchService.getFacultyBatchFeedback(this.batchID, this.loginID, params).subscribe((response: any) => {
      this.disableButton = false
      this.facultyFeedbacks = response.body
      // console.log(this.facultyFeedbacks);
      this.totalStudents = this.facultyFeedbacks.length
      this.assignFacultyFeedbackFields()
      // this.assignBatchFeedback()
    },
      (err: any) => {
        this.totalStudents = 0
        this.disableButton = false
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
    this.spinnerService.loadingMessage = "Adding feedback"


    this.batchService.addFacultyBatchFeedbacks(this.batchID, this.feedbackResponse).
      subscribe((response: any) => {
        // console.log(response)
        this.disableButton = false
        this.modalRef.close()
        this.getAllFacultyBatchFeedback()
        // this.getAllFeedbackForSession()
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
    // console.log(this.talentBatchFeedbackForm.value);
    this.disableButton = true
    this.spinnerService.loadingMessage = "Adding feedback"


    this.batchService.addTalentBatchFeedbacks(this.batchID, this.talentBatchFeedbackForm.value.feedbacks).
      subscribe((response: any) => {
        // console.log(response.body)
        this.disableButton = false
        this.modalRef.close()
        this.getAllTalentBatchFeedback()
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

  deleteFacultyBatchFeedback(): void {
    if (this.deleteTalentFeedback) {
      this.deleteTalentBatchFeedback()
      return
    }
    this.disableButton = true
    // console.log(this.facultyFeedback);
    this.spinnerService.loadingMessage = "Deleting feedback"


    this.batchService.deleteFacultyBatchFeedback(this.batchID, this.facultyFeedback.talent.id).subscribe((response: any) => {
      // console.log(response);
      this.disableButton = false
      alert("Feedback successfully deleted")
      this.getAllFacultyBatchFeedback()
      // this.getAllFeedbackForSession()
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

  deleteTalentBatchFeedback(): void {
    this.disableButton = true
    this.spinnerService.loadingMessage = "Deleting feedback"


    this.batchService.deleteTalentBatchFeedback(this.batchID, this.facultyID, this.talentID).
      subscribe((response: any) => {
        // console.log(response);
        this.disableButton = false
        alert("Feedback successfully deleted")
        this.getAllTalentBatchFeedback()
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

  getFeedbackQuestionsForFaculty(questionType: string): void {
    this.disableButton = true


    // this.generalService.getFeedbackQuestionByType(questionType).subscribe((response: any) => {
    this.feedbackGroupService.getFeedbackQuestionGroupByType(questionType).subscribe((response: any) => {
      this.disableButton = false
      this.facultyFeedbackQuestions = response.body
      // console.log(response.body);
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

  getFeedbackQuestionsForTalent(questionType: string): void {
    this.disableButton = true


    this.feedbackService.getFeedbackQuestionByType(questionType).subscribe((response: any) => {
      // this.feedbackGroupSerivce.getFeedbackQuestionGroupByType(questionType).subscribe((response: any) => {
      this.disableButton = false
      this.talentFeedbackQuestions = response.body
      // console.log(response.body);
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

  // ===============================================================================================================================


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

  calculateTalentScore(): void {
    for (let index = 0; index < this.talentFeedbacks.length; index++) {
      let averageScore = 0
      let noOfFeedbacks = 0
      this.talentFeedbacks[index].showFeedback = true
      this.talentFeedbacks[index].averageScore = 0
      for (let j = 0; j < this.talentFeedbacks[index].feedbacks?.length; j++) {
        this.talentFeedbacks[index].feedbacks[j].showFeedback = false
        if (this.talentFeedbacks[index].feedbacks[j].batchFeedbacks?.length > 0) {
          noOfFeedbacks++
        }
        this.talentFeedbacks[index].feedbacks[j].averageScore = this.calculateTalentFeedbackScore(this.talentFeedbacks[index].feedbacks[j])
        averageScore += this.talentFeedbacks[index].feedbacks[j].averageScore
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
  calculateTalentFeedbackScore(talentFeedback: IBatchFeedbackDTO): number {
    let averageScore = 0
    let totalScore = 0

    if (talentFeedback.batchFeedbacks.length == 0) {
      return 0
    }

    for (let index = 0; index < talentFeedback.batchFeedbacks.length; index++) {
      if (talentFeedback.batchFeedbacks[index].option && talentFeedback.batchFeedbacks[index].question) {
        averageScore += talentFeedback.batchFeedbacks[index].option?.key
        totalScore += talentFeedback.batchFeedbacks[index].question.maxScore
      }
    }

    averageScore = (averageScore * 10) / totalScore
    return averageScore
  }

  assignFacultyFeedbackFields(): void {
    for (let index = 0; index < this.facultyFeedbacks.length; index++) {
      this.facultyFeedbacks[index].showFeedback = false
      this.facultyFeedbacks[index].averageScore = this.calculateFacultyFeedbackScore(this.facultyFeedbacks[index])

    }
  }

  calculateFacultyFeedbackScore(facultyFeeback: IFacultyBatchFeedback): number {
    let averageScore = 0
    let totalScore = 0

    if (facultyFeeback.batchFeedbacks?.length == 0) {
      return 0
    }

    for (let index = 0; index < facultyFeeback.batchFeedbacks?.length; index++) {
      if (facultyFeeback.batchFeedbacks[index].option && facultyFeeback.batchFeedbacks[index].question) {
        averageScore += facultyFeeback.batchFeedbacks[index].option?.key
        totalScore += facultyFeeback.batchFeedbacks[index].question.maxScore
      }
    }

    averageScore = (averageScore * 10) / totalScore
    return averageScore
  }

  counter() {
    return new Array(this.MAX_SCORE + 1);
  }
}

