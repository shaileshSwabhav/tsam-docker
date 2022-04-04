import { Component, OnInit } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { IFacultyBatchSessionFeedback, IMappedSession, IBatchSessionFeedback, BatchService } from 'src/app/service/batch/batch.service';
import { IFeedbackQuestion, IFeedbackOptions, GeneralService } from 'src/app/service/general/general.service';
import { LocalService } from 'src/app/service/storage/local.service';

@Component({
  selector: 'app-batch-session-completion-feedback',
  templateUrl: './batch-session-completion-feedback.component.html',
  styleUrls: ['./batch-session-completion-feedback.component.css']
})
export class BatchSessionCompletionFeedbackComponent implements OnInit {


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

  // constants
  INITIAL_SCORE = 0
  MINIMUM_SCORE = 7
  MAX_SCORE = 10
  FEEDBACK_TALENT_INDEX = 0
  nextSessionDate: string

  // session feedback
  sessionFeedbackForm: FormGroup
  feedbackForm: FormGroup
  talentSessionFeedbackForm: FormGroup
  isSingleTalentFeedback: boolean
  isSessionFeedbackLoaded: boolean
  isTalentsLoaded: boolean

  batchID: string
  batchSessionID: string


  talent: any
  sessionDate: string

  // spinner


  constructor(
    private spinnerService: SpinnerService,
    private formBuilder: FormBuilder,
    private generalService: GeneralService,
    private localService: LocalService,
    private batchService: BatchService,
  ) { }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
    this.initializeVariables()
    this.getAllComponents()

  }

  initializeVariables() {
    this.loginID = this.localService.getJsonValue("loginID")
    this.facultyName = this.localService.getJsonValue("firstName")
  }
  getAllComponents() {
    this.getFeedbackQuestionsForFaculty("Faculty_Session_Feedback")
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


  onFacultyFeedbackInput(event: any, index: number, feedbackQuestionControl: any): void {
    // console.log(feedbackQuestionControl)

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


  getFeedbackQuestionsForFaculty(questionType: string): void {

    this.generalService.getFeedbackQuestionByType(questionType).subscribe((response: any) => {
      // console.log(response.body);

      this.facultyFeedbackQuestions = response.body
    }, (err: any) => {
      console.error(err)
    }).add(() => {


    })
  }



  submitFeedback(): void {
    // console.log(this.feedbackArray.valid);
    let errors: string[] = []

    // this.talentCount++;
    if (this.feedbackArray.valid) {
      // if (this.talentCount == this.talents.length) {
      //   this.addFacultyFeedback(errors)
      //   const queryParams: Params = { operation: this.ASSIGNMENTOPERATION };
      //   this.tabToBeOpen = this.ASSIGNMENTOPERATION
      //   this.updateRoute(queryParams)
      //   return
      // }
      this.addFacultyFeedback(errors)
      // this.createFeedbackQuestionsForm(this.tale)
    }

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
}
