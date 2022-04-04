import { Component, Input, OnInit, Output, ViewChild, EventEmitter } from '@angular/core';
import { FormGroup, FormBuilder, FormControl, Validators, FormArray } from '@angular/forms';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { UrlConstant, Role } from 'src/app/service/constant';
import { GeneralService } from 'src/app/service/general/general.service';
import { IProgrammingQuestion, IProgrammingQuestionType, ProgrammingQuestionService } from 'src/app/service/programming-question/programming-question.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { ProgrammingConceptService } from 'src/app/service/programming-concept/programming-concept.service';
import { HttpParams } from '@angular/common/http';
import { IProgrammingConcept } from 'src/app/models/programming/concept';
import { environment } from 'src/environments/environment';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { timeValidator } from 'src/app/Validators/custom.validators';
import { IPermission } from 'src/app/service/menu/menu.service';

@Component({
  selector: 'app-programming-question-modal',
  templateUrl: './programming-question-modal.component.html',
  styleUrls: ['./programming-question-modal.component.css']
})
export class ProgrammingQuestionModalComponent implements OnInit {

  // Components.
  programmingQuestionTypeList: IProgrammingQuestionType[]
  programmingQuestionLevel: any[]
  programmingLanguageList: any[]
  programmingConceptList: IProgrammingConcept[]

  // Topic.
  @Input() topicID: string
  @Input() programmingQuestion: IProgrammingQuestion
  @Input() isViewClicked: boolean
  @Input() isUpdateClicked: boolean

  // Form.
  programmingQuestionForm: FormGroup

  // Cke editor configuration.
  ckeEditorCommonConfig: any
  ckeEditorExampleConfig: any
  ckeEditorQuestionConfig: any
  ckeEditorOptionConfig: any
  ckeEditorConstraintsConfig: any
  programmingQuestionImagesFolderPath: string
  imageUploadURL: string = environment.UPLOAD_API_PATH
  fileUploadLocation: string = environment.FILE_UPLOAD_LOACTION

  // Child elements.
  @ViewChild("programmingModal") programmingModal: any
  @ViewChild("ckeditorExample") ckeditorExample: any
  @ViewChild("ckeditorQuestion") ckeditorQuestion: any
  @ViewChild("ckeditorOption") ckeditorOption: any

  // Flags.
  showTestCasesInProgrammingForm: boolean
  hasAnyTalentAnswered: boolean

  // Child to parent.
  @Output() dismissModalEvent: EventEmitter<void> = new EventEmitter()
  @Output() addEvent: EventEmitter<string> = new EventEmitter<string>()
  @Output() questionUpdatedEvent: EventEmitter<IProgrammingQuestion> = new EventEmitter<IProgrammingQuestion>()

  // Programming Question Solution.
  solutionCount: number

  // Access.
  permission: IPermission
  roleName: string

  constructor(
    private formBuilder: FormBuilder,
    public utilService: UtilityService,
    private questionService: ProgrammingQuestionService,
    private generalService: GeneralService,
    private spinnerService: SpinnerService,
    private programmingConceptService: ProgrammingConceptService,
    private fileOps: FileOperationService,
    private localService: LocalService,
    private role: Role,
    private urlConstant: UrlConstant,
  ) {
    this.initializeVariables()
    this.getAllComponents()
    this.createProgrammingQuestionForm()
  }

  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
    if (this.isViewClicked){
      this.solutionCount = this.programmingQuestion.solutionCount
      this.hasAnyTalentAnswered = this.programmingQuestion.hasAnyTalentAnswered
      this.onViewProgrammingQuestionClick()
    }
    this.getProgrammingConceptsList()
   }

  // Initialize global variables.
  initializeVariables(): void {

    // Access.
    this.roleName = this.localService.getJsonValue("roleName")
    if (this.roleName == this.role.ADMIN || this.roleName == this.role.SALES_PERSON) {
      this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.TRAINING_PROGRAMMING_QUESTION)
    }
    if (this.roleName == this.role.FACULTY) {
      this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.BANK_PROGRAMMING_QUESTION)
    }

    // Components.
    this.programmingQuestionTypeList = []
    this.programmingQuestionLevel = []
    this.programmingLanguageList = []
    this.programmingConceptList  = []

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Programming Questions"

    // Flags.
    this.showTestCasesInProgrammingForm = false

    // Cke editor configuration.
    this.programmingQuestionImagesFolderPath = this.fileOps.PROGRAMMING_QUESTION_IMAGES
    this.ckeEditorCommonCongiguration()
    this.ckeEditorQuestionCongiguration()
    this.ckeEditorOptionCongiguration()
    this.ckeEditorConstraintsCongiguration()
    this.ckeEditorExampleCongiguration()
  }

  // =============================================================CREATE FORMS==========================================================================

  // Create programming question form.
  createProgrammingQuestionForm(): void {
    this.programmingQuestionForm = this.formBuilder.group({
      id: new FormControl(null),
      programmingQuestionTypes: new FormControl(Array()),
      label: new FormControl(null, [Validators.required]),
      question: new FormControl(null, [Validators.required, Validators.maxLength(1000)]),
      hasOptions: new FormControl(false),
      hasTestCases: new FormControl(false),
      isActive: new FormControl(true, [Validators.required]),
      isLanguageSpecific: new FormControl(false, [Validators.required]),
      level: new FormControl(null, [Validators.required]),
      score: new FormControl(null, [Validators.required, Validators.min(1)]),
      timeHour: new FormControl(null, [Validators.required, Validators.min(0), Validators.max(23)]),
      timeMin: new FormControl(null, [Validators.required, Validators.min(0), Validators.max(59),]),
      timeRequired: new FormControl(null),
      options: this.formBuilder.array([]),
      testCases: this.formBuilder.array([]),
      // programmingLanguages: new FormControl(Array()),
      example: new FormControl(null, [Validators.required, Validators.maxLength(2000)]),
      constraints: new FormControl(null, [Validators.required, Validators.maxLength(500)]),
      comment: new FormControl(null, [Validators.maxLength(500)]),
      programmingConcept: new FormControl(Array()),
    }, { validator: timeValidator("timeHour", "timeMin") })
  }

  // Add programming option to options of programming question form.
  addProgrammingOptionsToForm(): void {
    this.programmingOptions.push(this.formBuilder.group({
      id: new FormControl(null),
      programmingQuestionID: new FormControl(null),
      option: new FormControl(null, [Validators.required, Validators.maxLength(100)]),
      order: new FormControl(null, [Validators.required, Validators.min(1)]),
      isCorrect: new FormControl(false, [Validators.required]),
      isActive: new FormControl(true, [Validators.required]),
    }))
  }

  // Add programming test case to test cases of programming question form.
  addTestCasesToProgrammingForm(): void {
    this.programmingTestCases.push(this.formBuilder.group({
      id: new FormControl(null),
      programmingQuestionID: new FormControl(null),
      input: new FormControl(null, [Validators.required, Validators.maxLength(500)]),
      output: new FormControl(null, [Validators.required, Validators.maxLength(500)]),
      explanation: new FormControl(null, [Validators.maxLength(500)]),
      isHidden: new FormControl(false, [Validators.required]),
      isActive: new FormControl(true, [Validators.required]),
    }))
  }

  // =============================================================PROGRAMMING QUESTION CRUD FUNCTIONS==========================================================================

  // Add programming question.
  addProgrammingQuestion(): void {
    this.spinnerService.loadingMessage = "Adding Programming Question"
    // let programmingQuestion: IProgrammingQuestion = this.programmingQuestionForm.value
    // this.patchIDFromObject(programmingQuestion)
    this.questionService.addProgrammingQuestion(this.programmingQuestionForm.value).subscribe((response) => {
      alert("Programming question successfully added")
      this.addEvent.emit(response.body)
    }, (error) => {
      console.error(error);
      if (error.error?.error) {
        alert(error.error?.error)
        return;
      }
      alert(error.statusText)
    })
  }

  // On clicking view programming question button.
  onViewProgrammingQuestionClick(): void {
    console.log(this.programmingQuestion)
    if (this.programmingQuestion.isLanguageSpecific) {
      this.programmingQuestionForm.addControl("programmingLanguages", this.formBuilder.control(null, [Validators.required]))
    }
    this.patchForm(this.programmingQuestion)
    let hr: number = Math.floor(this.programmingQuestion.timeRequired / 60)
    let min = this.programmingQuestion.timeRequired % 60
    this.programmingQuestionForm.get('timeHour').setValue(hr)
    this.programmingQuestionForm.get('timeMin').setValue(min)
    this.programmingQuestionForm.disable()
  }

  // On clicking update programming question button.
  onUpdateProgrammingQuestionClick(): void {
    this.isViewClicked = false
    this.isUpdateClicked = true
    this.programmingQuestionForm.enable()
  }

  // Update programming question by sending programming question to parent component.
  updateProgrammingQuestion(): void {
    let programmingQuestion: IProgrammingQuestion = this.programmingQuestionForm.value
    this.questionUpdatedEvent.emit(programmingQuestion)
  }

  // =======================================OTHER FUNCTIONS FOR PROGRAMMING QUESTION==========================================

  // Getter method for options of programming question form.
  get programmingOptions() {
    return this.programmingQuestionForm.get('options') as FormArray
  }

  // Getter method for test cases of programming question form.
  get programmingTestCases() {
    return this.programmingQuestionForm.get('testCases') as FormArray
  }

  // Delete single option of programming question form. 
  deleteProgrammingQuestionOption(index: number): void {
    if (confirm("Are you sure you want to delete this option ?")) {
      this.programmingOptions.removeAt(index)
      this.programmingQuestionForm.markAsDirty()
    }
  }

  // Delete single test case of programming question form. 
  deleteTestCaseFromProgrammingForm(index: number): void {
    if (confirm("Are you sure you want to delete this test case ?")) {
      this.programmingTestCases.removeAt(index)
      this.programmingQuestionForm.markAsDirty()
    }
  }

  // Patch value from programming question to programming question form.
  patchForm(question: IProgrammingQuestion): void {

    // Add options to question form.
    if (!question.hasOptions) {
      question.options = []
      this.programmingQuestionForm.setControl("options", this.formBuilder.array([]))
    } else {
      this.programmingQuestionForm.setControl("options", this.formBuilder.array([]))
      for (let index = 0; index < question.options.length; index++) {
        this.addProgrammingOptionsToForm()
      }
    }
    this.programmingQuestionForm.patchValue(question)
  }

  // On changing the value of has options in programming question form.
  onHasOptionChange(): void {

    // If question has solutions then dont let it have options.
    if (this.programmingQuestionForm.controls["hasOptions"]?.value && this.solutionCount > 0) {
      alert("Programming question already has solutions")
      this.programmingQuestionForm.controls["hasOptions"].setValue(false)
      return
    }

    // If has options field is changed to true.
    if (this.programmingQuestionForm.controls["hasOptions"]?.value) {
      this.setValidatorsForOptionsAvailable()
      this.addProgrammingOptionsToForm()
      return
    }

    // If has options field is changed to false.
    this.setValidatorsForOptionsNotAvailable()
    this.programmingOptions.clear()
  }

  // Set form controls of question form when has options field is true.
  setValidatorsForOptionsAvailable(): void {

    this.programmingQuestionForm.controls["example"].setValue(null)
    this.programmingQuestionForm.controls["example"].setValidators(null)

    this.programmingQuestionForm.controls["constraints"].setValue(null)
    this.programmingQuestionForm.controls["constraints"].setValidators(null)

    this.programmingQuestionForm.controls["comment"].setValue(null)

    this.programmingQuestionForm.controls["options"].setValue([])
    this.programmingQuestionForm.controls["options"].setValidators([Validators.required])

    this.programmingQuestionForm.setControl("testCases", this.formBuilder.array([]))
    this.programmingQuestionForm.controls["hasTestCases"].setValue(false)

    this.utilService.updateValueAndValiditors(this.programmingQuestionForm)

  }

  // Set form controls of question form when has options field is false.
  setValidatorsForOptionsNotAvailable(): void {

    this.programmingQuestionForm.controls["example"].setValue(null)
    this.programmingQuestionForm.controls["example"].setValidators([Validators.required, Validators.maxLength(500)])

    this.programmingQuestionForm.controls["constraints"].setValue(null)
    this.programmingQuestionForm.controls["constraints"].setValidators([Validators.required, Validators.maxLength(500)])

    this.programmingQuestionForm.controls["comment"].setValue(null)
    this.programmingQuestionForm.controls["comment"].setValidators([Validators.maxLength(500)])

    this.programmingQuestionForm.setControl("options", this.formBuilder.array([]))

    this.programmingQuestionForm.controls["testCases"].setValue([])

    this.utilService.updateValueAndValiditors(this.programmingQuestionForm)
  }

  // // Extract id from objects before adding or updating programming question.
  // patchIDFromObject(programmingQuestion: any): void {
  //   if (this.programmingQuestionForm.get('programmingQuestionType').value) {
  //     programmingQuestion.programmingQuestionTypeID = this.programmingQuestionForm.get('programmingQuestionType').value.id
  //     delete programmingQuestion['programmingQuestionType']
  //   }
  // }

   // Sets the time required field form hours and mins to minutes.
   setProgrammingTimeRequired(): void {
    this.programmingQuestionForm.get('timeRequired').setValue(this.programmingQuestionForm.get('timeHour').value * 60 + this.programmingQuestionForm.get('timeMin').value);
  }

  // Validate add/update programming question form.
  validate(): void {
    this.setProgrammingTimeRequired()
    if (this.programmingQuestionForm.get('timeRequired').value <= 0) {
      // alert("Time is required, value cann't be zero")
      const err = { inValid: true }
      this.programmingQuestionForm.setErrors(err)
    }
    if (this.programmingQuestionForm.invalid) {
      this.programmingQuestionForm.markAllAsTouched()
      return
    }
    if (this.isUpdateClicked) {
      this.updateProgrammingQuestion()
      return
    }
    this.addProgrammingQuestion()
  }

  // Get all invalid controls in programming question form.
  public findInvalidControls(): any[] {
    const invalid = []
    const controls = this.programmingQuestionForm.controls
    for (const name in controls) {
      if (controls[name].invalid) {
        invalid.push(name)
      }
    }
    return invalid
  }

  // On clicking hasTestCases checkbox in programming question form.
  toggleTestCaseControls(event: any): void {
    if (this.showTestCasesInProgrammingForm) {
      if (confirm("Are you sure you want to delete all test cases?")) {
        this.showTestCasesInProgrammingForm = false
        this.programmingQuestionForm.setControl('testCases', this.formBuilder.array([]))
        this.programmingQuestionForm.get("hasTestCases").setValue(false)
        return
      }
      event.target.checked = true
      return
    }
    this.showTestCasesInProgrammingForm = true
    this.programmingQuestionForm.get("hasTestCases").setValue(true)
    this.programmingQuestionForm.get("testCases").setValue([])
    this.addTestCasesToProgrammingForm()
  }

  // On changing value of is language specific field in qusetion form.
  onIsLanguageSpecificChange(): void {

    // If is language specific true then add programming languages control.
    if (this.programmingQuestionForm.get('isLanguageSpecific').value == true) {
      this.programmingQuestionForm.addControl("programmingLanguages", this.formBuilder.control(null, [Validators.required]))
    }

    // If is language specific false then remove concept module.
    if (this.programmingQuestionForm.get('isLanguageSpecific').value == false) {
      this.programmingQuestionForm.removeControl('programmingLanguages')
    }
  }

  // Initialize time in hours and minutes.
  initializeTimeInHoursAndMinutes(): void {
    if (this.programmingQuestionForm.get('timeMin').value == null) {
      this.programmingQuestionForm.get('timeMin').setValue(0)
    }
    if (this.programmingQuestionForm.get('timeHour').value == null) {
      this.programmingQuestionForm.get('timeHour').setValue(0)
    }

    // const hour = control.get('timeHour');
    // const min = control.get('timeMin');
    // if (hour.value==null||min.value==null || hour.value==0|| min.value==0) {
    //     return {toTimeValidator: false}
    // }

    // console.log(hour.value > min.value);

    // return hour.value*60+ min.value>5 ? { toTimeValidator: true } : { toTimeValidator: false };
    // return this.programmingQuestionForm.get('timeHour').value*60+ this.programmingQuestionForm.get('timeMin').value;
  }

  // Common cke editor congiguration.
  ckeEditorCommonCongiguration(): void {
    this.ckeEditorCommonConfig = {
      extraPlugins: 'codeTag,kbdTag',
      removePlugins: "exportpdf",
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
        { name: 'document', groups: ['mode', 'document', 'doctools'], items: ['Source'] },
      ],
      toolbarGroups: [
        { name: 'styles' },
        { name: 'basicstyles', groups: ['basicstyles', 'cleanup'] },
        { name: 'document', groups: ['mode', 'document', 'doctools'] },
        { name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'] },
      ],
      removeButtons: "",
      language: 'en',
      forcePasteAsPlainText: false,
    }
  }

  // Configure the cke editor for question.
  ckeEditorQuestionCongiguration(): void {

    // Give common config to cke editor question.
    this.ckeEditorQuestionConfig = {}
    this.ckeEditorQuestionConfig = { ...this.ckeEditorCommonConfig }
    this.ckeEditorQuestionConfig.toolbar = [...this.ckeEditorCommonConfig.toolbar]
    this.ckeEditorQuestionConfig.toolbarGroups = [...this.ckeEditorCommonConfig.toolbarGroups]

    // Add extra options for cke editor question.
    this.ckeEditorQuestionConfig.toolbar.push(...[
      { name: 'insert', items: ['SImage'] },
      { name: 'links', items: ['Link', 'Unlink'] },
    ])

    this.ckeEditorQuestionConfig.toolbarGroups.push(...[
      { name: 'links' },
      { name: 'insert' },
    ])

    this.ckeEditorQuestionConfig.folderPath = this.programmingQuestionImagesFolderPath,
    this.ckeEditorQuestionConfig.imageUploadURL = this.imageUploadURL,
    this.ckeEditorQuestionConfig.fileUploadLocation = this.fileUploadLocation,
    this.ckeEditorQuestionConfig.allowedContent = true,
    this.ckeEditorQuestionConfig.extraAllowedContent = 'img'
    this.ckeEditorQuestionConfig.extraPlugins = 'codeTag,kbdTag,simage'
  }

  // Configure the cke editor for option.
  ckeEditorOptionCongiguration(): void {

    // Give common config to cke editor option.
    this.ckeEditorOptionConfig = {}
    this.ckeEditorOptionConfig = { ...this.ckeEditorCommonConfig }
  }

  // Configure the cke editor for constraints.
  ckeEditorConstraintsCongiguration(): void {
    
    // Give common config to cke editor constraints.
    this.ckeEditorConstraintsConfig = {}
    this.ckeEditorConstraintsConfig = { ...this.ckeEditorCommonConfig }
  }

  // Configure the cke editor for example.
  ckeEditorExampleCongiguration(): void {

    // Give common config to cke editor example.
    this.ckeEditorExampleConfig = {}
    this.ckeEditorExampleConfig = { ...this.ckeEditorCommonConfig }
    this.ckeEditorExampleConfig.toolbar = [...this.ckeEditorCommonConfig.toolbar]
    this.ckeEditorExampleConfig.toolbarGroups = [...this.ckeEditorCommonConfig.toolbarGroups]

    // Add extra options for cke editor example.
    this.ckeEditorExampleConfig.toolbar.push(...[
      { name: 'insert', items: ['SImage'] },
      { name: 'links', items: ['Link', 'Unlink'] },
    ])

    this.ckeEditorExampleConfig.toolbarGroups.push(...[
      { name: 'links' },
      { name: 'insert' },
    ])

    this.ckeEditorExampleConfig.folderPath = this.programmingQuestionImagesFolderPath,
    this.ckeEditorExampleConfig.imageUploadURL = this.imageUploadURL,
    this.ckeEditorExampleConfig.fileUploadLocation = this.fileUploadLocation,
    this.ckeEditorExampleConfig.allowedContent = true,
    this.ckeEditorExampleConfig.extraAllowedContent = 'img'
    this.ckeEditorExampleConfig.extraPlugins = 'codeTag,kbdTag,simage'
  }

  // Dismiss programming question modal.
  dismissModal(): void {
    this.dismissModalEvent.emit()
  }

  // =======================================GET FUNCTIONS========================================================

  // Get all components.
  getAllComponents(): void {
    this.getProgrammingQuestionTypeList()
    this.getProgrammingQuestionLevel()
    this.getProgrammingLanguageList()
  }

  // Get programming question type list.
  getProgrammingQuestionTypeList(): void {
    this.generalService.getProgrammingQuestionTypeList().subscribe((response: any) => {
      this.programmingQuestionTypeList = response
    }, (err: any) => {
      console.error(err)
    })
  }

  // Get programming question level list.
  getProgrammingQuestionLevel(): void {
    this.generalService.getGeneralTypeByType("programming_question_level").subscribe((response: any) => {
      this.programmingQuestionLevel = response
    }, (err: any) => {
      console.error(err)
    })
  }

  // Get programming language list.
  getProgrammingLanguageList(): void {
    this.generalService.getProgrammingLanguageList().subscribe((response: any) => {
      this.programmingLanguageList = response
    }, (err: any) => {
      console.error(err)
    })
  }

  // Get programming concept list.
  getProgrammingConceptsList(): void {
    let params: HttpParams = new HttpParams()
    params = params.append("limit", "-1")
    params = params.append("offset", "0")
    if (this.topicID) {
      console.log("wkdkwmlwfm")
      params = params.append("topicID", this.topicID)
    }
    this.programmingConceptService.getAllConcepts(params).subscribe((response) => {
      this.programmingConceptList = response.body
    }, (err: any) => {
      console.error(err)
    })
  }

}
