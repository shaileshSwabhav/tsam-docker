import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { FormGroup, FormArray, FormControl, Validators, FormBuilder } from '@angular/forms';
import { IFeedbackOptions, IFeedbackQuestion } from 'src/app/service/general/general.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { FeedbackService } from 'src/app/service/feedback/feedback.service';
import { BatchService } from 'src/app/service/batch/batch.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { FacultyService } from 'src/app/service/faculty/faculty.service';

@Component({
  selector: 'app-talent-feedback-modal',
  templateUrl: './talent-feedback-modal.component.html',
  styleUrls: ['./talent-feedback-modal.component.css']
})
export class TalentFeedbackModalComponent implements OnInit {

  readonly MAX_SCORE = 10

  // Talent feedback for faculty.
  talentFeedbackToFaculytyForm: FormGroup
  talentTofacultyFeedbackQuestionList: IFeedbackQuestion[]
  // feedbackRelatedDetailsCount: number

  // Talent.
  talentID: string

  // Spinner.



  // Values to parent.
  @Output() isSuccessfulEmitter = new EventEmitter<any>()

  // Values from parent.
  @Input() batchSessionID: string
  @Input() topics: any[]
  @Input() batchID: string
  @Input() batchSessionDate: string
  @Input() sessionNumber: number
  @Input() feedbacks: any
  @Input() faculty: any

  constructor(
    private formBuilder: FormBuilder,
    private spinnerService: SpinnerService,
    private feedbackService: FeedbackService,
    private batchService: BatchService,
    private localService: LocalService,
    private facultyService: FacultyService,
  ) {
    this.initializeVariables()
    this.getFeedbackQuestionsForTalent("Talent_Session_Feedback")
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize global variables.
  initializeVariables(): void {

    // Talent
    this.talentID = this.localService.getJsonValue("loginID")

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."


    // Talent feedback for faculty.
    this.talentTofacultyFeedbackQuestionList = []
    // this.feedbackRelatedDetailsCount = 2
  }

  // =============================================================CREATE FORMS==========================================================================

  // Create all forms.
  createForms(): void {
    this.createTalentFeedbackToFacultyForm()

    // If feedback is not given.
    if (!this.feedbacks || this.feedbacks?.length == 0) {
      this.createFeedbackQuestionsForm()
    }

    // If feedback is given.
    if (this.feedbacks?.length > 0) {
      this.createSubmittedFeedbackQuestionsForm()
    }
  }

  // Create talent feedback to faculty form.
  createTalentFeedbackToFacultyForm(): void {
    this.talentFeedbackToFaculytyForm = this.formBuilder.group({
      feedbacks: new FormArray([])
    })
  }

  // Add feedback questions to form.
  addFeedbackQuestionsToArray(): void {
    this.feedbackArray.push(this.formBuilder.group({
      batchID: new FormControl(this.batchID),
      batchSessionID: new FormControl(this.batchSessionID),
      talentID: new FormControl(this.talentID),
      facultyID: new FormControl(this.faculty.id),
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

  // Set fields of form controls of form for submitted feedback.
  createSubmittedFeedbackQuestionsForm(): void {
    for (let i = 0; i < this.feedbacks?.length; i++) {
      let feedback: any = this.feedbacks[i]
      this.addFeedbackQuestionsToArray()
      this.feedbackArray.at(i).get("questionID").setValue(feedback.questionID)
      this.feedbackArray.at(i).get("optionID").setValue(feedback.optionID)
      this.feedbackArray.at(i).get("optionID").disable()
      this.feedbackArray.at(i).get("option").setValue(feedback.option)
      this.feedbackArray.at(i).get("option").disable()
      this.feedbackArray.at(i).get("answer").setValue(feedback.answer)
      this.feedbackArray.at(i).get("answer").disable()
    }
  }

  // Get feedbacks form control of talent feedback to faculty form.
  get feedbackArray(): FormArray {
    return this.talentFeedbackToFaculytyForm.get("feedbacks") as FormArray
  }

  // ============================================================= FORM FUNCTIONS ==========================================================================

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
    this.batchService.addTalentSessionFeedbacks(this.batchID, this.batchSessionID, talentFeedback.feedbacks).
      subscribe((response: any) => {
        // this.getListsRelatedToBatch()
        this.sendIsSuccessful(true)
        alert("Feedback successfully sent")
      }, (err: any) => {
        this.sendIsSuccessful(false)
        console.error(err)
        alert("Feedback could not be sent, please try again later")
      })
  }

  // Extract ID from objects and delete objects before adding or updating.
  patchIDFromObjects(talentFeedback: any): void {
    for (let i = 0; i < talentFeedback.feedbacks?.length; i++) {
      delete talentFeedback.feedbacks[i]['option']
    }
  }

  // ============================================================= OTHER FUNCTIONS ==========================================================================

  // Send is successful value to parent.
  sendIsSuccessful(isSuccessful: boolean): void {
    this.isSuccessfulEmitter.emit(isSuccessful)
  }

  // // Decrement batch related details count.
  // decrementBatchRelatedDetailsCount(): void{
  //   this.feedbackRelatedDetailsCount = this.feedbackRelatedDetailsCount - 1
  //   if (this.feedbackRelatedDetailsCount == 0){
  //     this.createForms()
  //   }
  // }

  // ============================================================= GET FUNCTIONS ==========================================================================

  // Get talent feedback questions for faculty.
  getFeedbackQuestionsForTalent(questionType: string): void {


    this.feedbackService.getFeedbackQuestionByType(questionType).subscribe((response: any) => {
      this.talentTofacultyFeedbackQuestionList = response.body
    }, (err: any) => {
      console.error(err)
    }).add(() => {
      this.createForms()
      // this.decrementBatchRelatedDetailsCount()


    })
  }
}
