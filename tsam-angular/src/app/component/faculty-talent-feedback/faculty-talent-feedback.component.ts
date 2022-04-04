import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { BatchTalentService } from 'src/app/service/batch-talent/batch-talent.service';
import { BatchService } from 'src/app/service/batch/batch.service';
import { GeneralService, IFeedbackOptions, IFeedbackQuestion } from 'src/app/service/general/general.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';
import { LocalService } from 'src/app/service/storage/local.service';

@Component({
  selector: 'app-faculty-talent-feedback',
  templateUrl: './faculty-talent-feedback.component.html',
  styleUrls: ['./faculty-talent-feedback.component.css']
})
export class FacultyTalentFeedbackComponent implements OnInit {
  loginID: string;
  facultyName: string;
  batchTalents: any[];
  facultyFeedbackQuestions: IFeedbackQuestion[]
  totalTalent: number;
  batchSessionID: string;
  batchID: string;
  feedbackForm: FormGroup
  talents: any[]

  feedbackTypeForm: FormGroup 

  @Output() sendToHome: EventEmitter<any> = new EventEmitter
  @Output() sendToNext: EventEmitter<any> = new EventEmitter
  
  readonly FACULTY_SESSION_FEEDBACK: string = 'Faculty_Session_Feedback'
  readonly FACULTY_SESSION_FEEDBACK_NON_TECH: string = 'Faculty_Session_Feedback_Non_Technical_Question'


  constructor(
    private formBuilder: FormBuilder,
    private localService: LocalService,
    private spinnerService: SpinnerService,
    private generalService: GeneralService,
    private batchTalentService: BatchTalentService,
    private activatedRoute: ActivatedRoute,
    private batchService: BatchService,
  ) {
    this.extractID()
    this.createFeedbackTypeForm()
  }

  ngOnInit(): void {
    this.initialize()
  }

  createFeedbackTypeForm(): void {
    this.feedbackTypeForm = this.formBuilder.group({
      feedbackType: new FormControl(this.FACULTY_SESSION_FEEDBACK)
    })
  }

  async onFeedbackTypeChange(): Promise<void> {
    console.log(this.feedbackTypeForm);
    this.facultyFeedbackQuestions = await this.getFeedbackQuestionsForFaculty(this.feedbackTypeForm.get("feedbackType").value)
    
    this.initFeedbackQuestionsForm()
  }

  async initialize(): Promise<void> {
    this.loginID = this.localService.getJsonValue("loginID")
    this.facultyName = this.localService.getJsonValue("firstName")
    this.batchTalents = []
    this.totalTalent = 0
    this.facultyFeedbackQuestions = []
    this.talents = []

    await this.getAllComponents()
    this.initFeedbackQuestionsForm()
  }

  async getAllComponents(): Promise<void> {
    let errors: string[] = []

    try {
      this.facultyFeedbackQuestions = await this.getFeedbackQuestionsForFaculty(this.feedbackTypeForm.get("feedbackType").value)
      this.batchTalents = await this.getPresentTalents()
      console.log(this.batchTalents);

    } catch (error) {
      errors.push(error)
    }

    if (errors.length > 0) {
      let errorString = ""
      for (let index = 0; index < errors.length; index++) {
        errorString += (index == 0 ? "" : "\n") + errors[index]
      }
      alert(errorString)
      return
    }
  }

  extractID(): void {
    this.activatedRoute.queryParamMap.subscribe(
      (params: any) => {
        this.batchID = params.get("batchID")
        this.batchSessionID = params.get("batchSessionID")
      }, (err: any) => {
        console.error(err);
      }).unsubscribe()
  }



  createFeedbackForm(): void {
    this.feedbackForm = this.formBuilder.group({
      feedbacks: new FormArray([])
    })
  }

  get feedbackArray(): FormArray {
    return this.feedbackForm.get("feedbacks") as FormArray
  }

  addFeedbackTalentsandQuestionsToForm(): void {
    this.feedbackArray.push(
      this.formBuilder.group({
        talentID: new FormControl(null, [Validators.required]),
        talentFirstName: new FormControl(null, Validators.required),
        talentLastName: new FormControl(null, Validators.required),
        feedbackQuestions: new FormArray([]),
      })
    )
  }

  getFeedbackQuestionsControlArray(index: number): FormArray {
    return this.feedbackArray.at(index).get("feedbackQuestions") as FormArray
  }

  addFeedbackQuestionsToArray(index: number): void {
    this.getFeedbackQuestionsControlArray(index).push(this.formBuilder.group({
      batchID: new FormControl(this.batchID),
      talentID: new FormControl(null, [Validators.required]),
      facultyID: new FormControl(this.loginID),
      questionID: new FormControl(null, [Validators.required]),
      optionID: new FormControl(null, [Validators.required]),
      option: new FormControl(null, [Validators.required]),
      answer: new FormControl(null, [Validators.required]),
    }))
  }


  initFeedbackQuestionsForm(): void {
    this.createFeedbackForm()
    for (let index = 0; index < this.batchTalents?.length; index++) {
      const element = this.batchTalents[index]?.talent;
      if (!this.batchTalents[index]?.isFeedbackGiven) {
        this.talents.push(element)
      }

      this.createFeedbackQuestionsForm(element, index)
    }
    console.log(this.talents);

  }
  createFeedbackQuestionsForm(talent: any, index: number): void {
    this.addFeedbackTalentsandQuestionsToForm()

    this.feedbackArray.at(index).get("talentID").setValue(talent?.id)
    this.feedbackArray.at(index).get("talentFirstName").setValue(talent?.firstName)
    this.feedbackArray.at(index).get("talentLastName").setValue(talent?.lastName)
    this.createFeedbackQuestionsFormArray(talent.id, index)

  }
  createFeedbackQuestionsFormArray(talentID: string, j: number): void {

    if (this.facultyFeedbackQuestions) {
      for (let index = 0; index < this.facultyFeedbackQuestions.length; index++) {
        // initialize form
        this.addFeedbackQuestionsToArray(j)

        this.getFeedbackQuestionsControlArray(j).at(index).get('facultyID').setValue(this.loginID)
        this.getFeedbackQuestionsControlArray(j).at(index).get("talentID").setValue(talentID)
        this.getFeedbackQuestionsControlArray(j).at(index).get("questionID").setValue(this.facultyFeedbackQuestions[index].id)

        if (this.facultyFeedbackQuestions[index].hasOptions) {
          this.getFeedbackQuestionsControlArray(j).at(index).get("optionID").setValidators([Validators.required])
          this.getFeedbackQuestionsControlArray(j).at(index).get("option").setValidators([Validators.required])
          this.getFeedbackQuestionsControlArray(j).at(index).get("answer").setValidators(null)
        } else {
          this.getFeedbackQuestionsControlArray(j).at(index).get("optionID").setValidators(null)
          this.getFeedbackQuestionsControlArray(j).at(index).get("option").setValidators(null)
          this.getFeedbackQuestionsControlArray(j).at(index).get("answer").setValidators([Validators.required])
        }

      }
    }
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
  async validateTalentsFeedbacks(): Promise<void> {
    let errors: string[] = []
    for (let index = 0; index < this.feedbackForm.value.feedbacks.length; index++) {
      const element = this.feedbackForm.value.feedbacks[index].feedbackQuestions;
      try {
        await this.addFacultyFeedback(element)
      } catch (error) {
        errors.push(error.error.error)
      }
    }
    if (errors.length > 0) {
      let errorString = ""
      for (let index = 0; index < errors.length; index++) {
        errorString += (index == 0 ? "" : "\n") + errors[index]
      }
      alert(errorString)
      this.navigateHome()
      return
    }
    alert('feedback submitted successfully')
    this.navigateToNext()
  }

  async addFacultyFeedback(feedbacks: any): Promise<void> {
    try {
      this.spinnerService.loadingMessage = "Adding feedback"
      this.spinnerService.isDisabled = true
      return await new Promise<any>((resolve, reject) => {
        this.batchService.addFacultySessionFeedbacks(this.batchID, this.batchSessionID, feedbacks).
          subscribe((response: any) => {
            resolve(response)
          }, (err: any) => {
            if (err.statusText.includes('Unknown')) {
              alert("No connection to server. Check internet.")
              return
            }
            reject(err)
          })
      })
    } finally {
      this.spinnerService.isDisabled = false
      // this.spinnerService.stopSpinner()
    }
  }

  async getFeedbackQuestionsForFaculty(questionType: string): Promise<any> {
    try {
      return await new Promise<any>((resolve, reject) => {
        this.spinnerService.loadingMessage = "Getting Feedback Questions"
        // this.spinnerService.startSpinner()
        this.generalService.getFeedbackQuestionByType(questionType).subscribe((response: any) => {
          resolve(response.body)
        }, (err: any) => {
          reject(err)
        })
      })
    } finally {
      // this.spinnerService.stopSpinner()
    }
  }

  navigateHome(): void {
    this.sendToHome.emit()
  }

  navigateToNext(): void {
    this.sendToNext.emit()
  }

  async getPresentTalents(): Promise<any> {
    try {
      return await new Promise<any>((resolve, reject) => {
        this.spinnerService.loadingMessage = "Getting Talents"
        // this.spinnerService.startSpinner()
        this.totalTalent = 0
        this.batchTalents = []
        let queryParams: any = {
          isPresent: 1,
          batchSessionID: this.batchSessionID
        }
        this.batchTalentService.getBatchSessionTalents(this.batchID, this.batchSessionID, queryParams).subscribe((response: any) => {
          this.totalTalent = response.headers.get('X-Total-Count')
          resolve(response.body)
        }, (err: any) => {
          this.totalTalent = 0
          if (err.statusText.includes('Unknown')) {
            alert("No connection to server. Check internet.")
            return
          }
          reject(err)
        })
      })
    } finally {
      // this.spinnerService.stopSpinner()
    }
  }

}
