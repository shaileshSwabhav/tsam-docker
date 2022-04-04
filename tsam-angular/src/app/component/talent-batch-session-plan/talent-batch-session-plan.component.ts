import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { BatchService, IMappedSession } from 'src/app/service/batch/batch.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { TalentDashboardService } from 'src/app/service/talent-dashboard/talent-dashboard.service';
import { ActivatedRoute } from '@angular/router';
import { LocalService } from 'src/app/service/storage/local.service';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { NgbModalRef, NgbModalOptions, NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { IFeedbackOptions, IFeedbackQuestion } from 'src/app/service/general/general.service';
import { FeedbackService } from 'src/app/service/feedback/feedback.service';

@Component({
  selector: 'app-talent-batch-session-plan',
  templateUrl: './talent-batch-session-plan.component.html',
  styleUrls: ['./talent-batch-session-plan.component.css']
})
export class TalentBatchSessionPlanComponent implements OnInit {

  readonly BATCH_STAUS_ONGOING: string = "Ongoing"
  readonly MAX_SCORE = 10

  // Batch Session.
  batchID: string
  existingSessions: number
  batchSessionList: IMappedSession[]
  @Input() batchStatus: string
  currentSessionOrder: number

  // Spinner.



  // Flags.
  isSessionPresent: boolean

  // Talent.
  talentID: string

  // Talent feedback for faculty.
  selectedBatchSessionForTalentFeedback: any
  @ViewChild('talentFeedbackToFacultyModal') talentFeedbackToFacultyModal: any
  talentFeedbackToFaculytyForm: FormGroup
  selectedBatchFaculty: any
  talentTofacultyFeedbackQuestionList: IFeedbackQuestion[]

  // Modal.
  modalRef: any

  constructor(
    private spinnerService: SpinnerService,
    private talentDashboardService: TalentDashboardService,
    private route: ActivatedRoute,
    private localService: LocalService,
    private formBuilder: FormBuilder,
    private batchService: BatchService,
    private modalService: NgbModal,
    private feedbackService: FeedbackService,
  ) {
    this.initializeVariables()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
    if (this.batchID && this.batchStatus) {
      this.getSessionsForBatch()
      this.getFeedbackQuestionsForTalent("Talent_Session_Feedback")
    }
  }

  // Initialize global variables.
  initializeVariables(): void {

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."


    // Batch session.
    this.batchSessionList = []
    this.currentSessionOrder = 0

    // Flags.
    this.isSessionPresent = false

    // Batch id.
    this.batchID = this.route.snapshot.queryParamMap.get("batchID")

    // Talent.
    this.talentID = this.localService.getJsonValue("loginID")

    // Talent feedback for faculty.
    this.talentTofacultyFeedbackQuestionList = []
  }

  //********************************************* TALENT FEEDACK FOR FACULTY FUNCTIONS ************************************************************

  // On clicking talent feedback to faculty.
  onTalentFeeddbackToFacultyClick(batchSession: any): void {
    this.selectedBatchSessionForTalentFeedback = batchSession
    this.selectedBatchFaculty = batchSession.faculty
    this.createTalentFeedbackToFaculytyForm()
    this.createFeedbackQuestionsForm()
    this.openModal(this.talentFeedbackToFacultyModal, 'lg')
  }

  // Create talent feedback to faculty form.
  createTalentFeedbackToFaculytyForm(): void {
    this.talentFeedbackToFaculytyForm = this.formBuilder.group({
      feedbacks: new FormArray([])
    })
  }

  // Get feedbacks form control of talent feedback to faculty form.
  get feedbackArray(): FormArray {
    return this.talentFeedbackToFaculytyForm.get("feedbacks") as FormArray
  }

  // Add feedback questions to form.
  addFeedbackQuestionsToArray(): void {
    this.feedbackArray.push(this.formBuilder.group({
      batchID: new FormControl(this.batchID),
      batchSessionID: new FormControl(this.selectedBatchSessionForTalentFeedback.id),
      talentID: new FormControl(this.talentID),
      facultyID: new FormControl(this.selectedBatchFaculty.id),
      questionID: new FormControl(null, [Validators.required]),
      optionID: new FormControl(null, [Validators.required]),
      option: new FormControl(null, [Validators.required]),
      answer: new FormControl(null, [Validators.required]),
    }))
  }

  // Set fields of form controls of form.
  createFeedbackQuestionsForm(): void {
    if (this.talentTofacultyFeedbackQuestionList) {
      for (let index = 0; index < this.talentTofacultyFeedbackQuestionList.length; index++) {
        this.addFeedbackQuestionsToArray()
        this.feedbackArray.at(index).get("questionID").setValue(this.talentTofacultyFeedbackQuestionList[index].id)
        if (this.talentTofacultyFeedbackQuestionList[index].hasOptions) {
          this.feedbackArray.at(index).get("optionID").setValidators([Validators.required])
          this.feedbackArray.at(index).get("option").setValidators([Validators.required])
          this.feedbackArray.at(index).get("answer").setValidators(null)
        } else {
          this.feedbackArray.at(index).get("optionID").setValidators(null)
          this.feedbackArray.at(index).get("option").setValidators(null)
          this.feedbackArray.at(index).get("answer").setValidators([Validators.required])
        }
      }
    }
  }

  // On changing option for feedback.
  onTalentFeedbackForFacultyChange(feedbackQuestionControl: any, feedbackOptions: IFeedbackOptions[]) {
    let optionID = feedbackQuestionControl.get("optionID").value
    for (let index = 0; index < feedbackOptions.length; index++) {
      if (optionID == feedbackOptions[index].id) {
        feedbackQuestionControl.get('answer').setValue(feedbackOptions[index].value)
        feedbackQuestionControl.get('option').setValue(feedbackOptions[index])
      }
    }
  }

  // On giving input to feedback.
  onTalentFeedbackForFacultyInput(event: any, index: number, feedbackQuestionControl: any): void {
    let options: IFeedbackOptions[] = this.talentTofacultyFeedbackQuestionList[index].options
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

  // Validate talent feedback before adding.
  validateTalentFeedback(): void {
    if (this.talentFeedbackToFaculytyForm.invalid) {
      this.talentFeedbackToFaculytyForm.markAllAsTouched()
      return
    }
    this.addTalentFeedback()
  }

  // Add talent feedback.
  addTalentFeedback(): void {
    this.spinnerService.loadingMessage = "Sending feedback"


    let talentFeedback: any = this.talentFeedbackToFaculytyForm.value
    this.patchIDFromObjects(talentFeedback)
    this.batchService.addTalentSessionFeedbacks(this.batchID, this.selectedBatchSessionForTalentFeedback.id,
      talentFeedback.feedbacks).
      subscribe((response: any) => {
        this.modalRef.close()
        this.getSessionsForBatch()
        alert("Feedback successfully sent")
      }, (err: any) => {
        console.error(err)
        alert("Feedback could not be sent, please try again later")
      })
  }

  // Extract ID from objects and delete objects before adding or updating.
  patchIDFromObjects(talentFeedback: any): void {
    for (let i = 0; i < talentFeedback.feedbacks.length; i++) {
      delete talentFeedback.feedbacks[i]['option']
    }
  }

  //********************************************* FORMAT FUNCTIONS ************************************************************

  // Set the view sub session clicked as false initially.
  setBatchSessionViewField(): void {
    for (let index = 0; index < this.batchSessionList.length; index++) {
      this.batchSessionList[index].session.viewSubSessionClicked = false
    }
  }

  // Format batch sesssion list.
  formatBatchSessions(): void {

    if (this.batchStatus == this.BATCH_STAUS_ONGOING) {

      // Set the first incomplete session's order as the current session order.
      for (let i = 0; i < this.batchSessionList.length; i++) {
        if (!this.batchSessionList[i].isCompleted) {
          this.currentSessionOrder = this.batchSessionList[i].order
          break
        }
      }
    }
  }

  //********************************************* OTHER FUNCTIONS ************************************************************

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

  // Get all sessions for batch.
  getSessionsForBatch(): void {
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."
    this.existingSessions = 0


    this.talentDashboardService.getBatchSessionWithBatchTalentDetails(this.batchID, this.talentID).subscribe((response: any) => {
      this.batchSessionList = response
      this.setBatchSessionViewField()
      this.formatBatchSessions()
      this.existingSessions = this.batchSessionList.length
      if (this.existingSessions > 0) {
        this.isSessionPresent = true
      }
    }, (err: any) => {
      this.existingSessions = 0
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  // Get talent feedback questions for faculty.
  getFeedbackQuestionsForTalent(questionType: string): void {


    this.feedbackService.getFeedbackQuestionByType(questionType).subscribe((response: any) => {
      this.talentTofacultyFeedbackQuestionList = response.body
    }, (err: any) => {
      console.error(err)
    })
  }

}
