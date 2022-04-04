import { DatePipe } from '@angular/common';
import { ChangeDetectorRef, Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BatchService } from 'src/app/service/batch/batch.service';
import { Role } from 'src/app/service/constant';
import { GeneralService } from 'src/app/service/general/general.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { CdkDragDrop, moveItemInArray } from '@angular/cdk/drag-drop';
import { BatchSessionService } from 'src/app/service/batch-session/batch-session.service';
import { BatchTopicAssignmentService } from 'src/app/service/batch-topic-assignment/batch-topic-assignment.service';

@Component({
  selector: 'app-batch-topic',
  templateUrl: './batch-topic.component.html',
  styleUrls: ['./batch-topic.component.css']
})
export class BatchTopicComponent implements OnInit {

  // Components.
  batchStatusList: any[]
  facultyList: any[]
  dayList: any[]
  loginID: string

  // Batch.
  batch: any
  batchForm: FormGroup
  batchID: string

  // Batch timing.
  batchModuleTimingsForm: FormGroup

  // Batch details.
  batchDetailsList: any[]

  // Flags.
  isFaculty: boolean

  // Batch module.
  previousBatchModule: any[]
  batchModuleList: any[]
  selectedBatchModuleList: any[]
  totalModules: number
  batchModuleTabList: any[]
  selectedCourseModuleTabList: any[]
  selectedBatchModule: any
  selectedBatchModuleIndex: number
  batchUpdateErrorCount: number
  batchUpdateCount: number

  // Session plan steps.
  sessionPlanStepsList: any[]
  selectedTemplateName: any
  currentStep: number
  batchSessionPlanEntryList: any[]
  @ViewChild("stepOne") stepOne: any
  @ViewChild("stepTwo") stepTwo: any
  @ViewChild("stepThree") stepThree: any
  @ViewChild("stepFour") stepFour: any
  @ViewChild("stepFive") stepFive: any
  @ViewChild("stepSix") stepSix: any
  @ViewChild("tempSessionPlan") sessionPlan: any

  // Batch topic assignment.
  batchTopicAssignmentEntryList: any[]
  batchTopicAssignmentAddCount: number

  // Session plan.
  batchSessionList: any[]

  // Faculty.
  facultyID: string

  constructor(
    private cdr: ChangeDetectorRef,
    private formBuilder: FormBuilder,
    public utilService: UtilityService,
    private route: ActivatedRoute,
    private spinnerService: SpinnerService,
    private batchService: BatchService,
    private generalService: GeneralService,
    private datePipe: DatePipe,
    private role: Role,
    private localService: LocalService,
    private batchSessionService: BatchSessionService,
    private batchTopicAssignmentService: BatchTopicAssignmentService,
  ) {
    this.initializeVariables()
    this.getAllComponents()
  }

  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  ngAfterViewInit() {
    this.initializeTabs()
    this.cdr.detectChanges()
  }

  // Initialize global variables.
  initializeVariables(): void {

    // Components.
    this.batchStatusList = []
    this.facultyList = []
    this.dayList = []
    this.loginID = this.localService.getJsonValue("loginID")

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Batch Session Plan"

    // Session plan.
    this.batchSessionList = []

    // Batch module.
    this.previousBatchModule = []
    this.batchModuleList = []
    this.selectedBatchModuleList = []
    this.batchModuleTabList = []
    this.selectedCourseModuleTabList = []
    this.selectedBatchModuleIndex = 0
    this.batchUpdateErrorCount = 0
    this.batchUpdateCount = 0

    // Session plan steps.
    this.batchSessionPlanEntryList = []

    // Batch topic assignment.
    this.batchTopicAssignmentEntryList = []
    this.batchTopicAssignmentAddCount = 0

    // Faculty.
    this.facultyID = this.localService.getJsonValue("loginID")

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

    // Flags.
    this.isFaculty = false

    // Role.
    if (this.role.FACULTY == this.localService.getJsonValue("roleName")) {
      this.isFaculty = true
    }

    // Create forms.
    this.createBatchForm()
    this.createBatchModuleTimingsForm()
  }

  // Initialize the session plan steps.
  initializeTabs(): void {

    // Session plan steps.
    this.sessionPlanStepsList = [
      { templateName: this.stepOne, step: 1 },
      { templateName: this.stepTwo, step: 2 },
      { templateName: this.stepThree, step: 3 },
      { templateName: this.stepFour, step: 4 },
      { templateName: this.stepFive, step: 5 },
      { templateName: this.stepSix, step: 6 },
      { templateName: this.sessionPlan, step: 7 },
    ]
    this.currentStep = 1

    // If batch id exists in params then get batch by id.
    this.route.queryParamMap.subscribe(params => {
      this.batchID = params.get("batchID")
      if (this.batchID) {
        this.getBatch()
      }
    }, err => {
      console.error(err)
    })
  }

  //********************************************* STEP ONE FUNCTIONS ************************************************************

  // On clicking create session plan button.
  onCreateSessionPlanClick(): void {
    this.currentStep = this.currentStep + 1
    this.showCurrentStepTemplate()
    this.getBatchModuleList()
  }

  //********************************************* STEP TWO FUNCTIONS ************************************************************

  // Format fields of module list.
  formatBatchModuleList(): void {

    // Count the number of sub topics, concepts and approx time in all topics.
    // Iterate all course modules.
    for (let i = 0; i < this.batchModuleList.length; i++) {
      let module: any = this.batchModuleList[i].module
      module.conceptCount = 0
      module.subTopicCount = 0
      module.approxTime = 0

      // Remove seconds form batch module timings.
      for (let j = 0; j < this.batchModuleList[i].moduleTimings?.length; j++) {
        this.batchModuleList[i].moduleTimings[j].fromTime = this.batchModuleList[i].moduleTimings[j].fromTime.replace(/:[^:]*$/, '')
        this.batchModuleList[i].moduleTimings[j].toTime = this.batchModuleList[i].moduleTimings[j].toTime.replace(/:[^:]*$/, '')
      }

      // Format start and end date.
      this.batchModuleList[i].startDate = this.datePipe.transform(this.batchModuleList[i].startDate, 'yyyy-MM-dd')
      this.batchModuleList[i].estimatedEndDate = this.datePipe.transform(this.batchModuleList[i].estimatedEndDate, 'yyyy-MM-dd')

      // Iterate all topics.
      for (let j = 0; j < module?.moduleTopics?.length; j++) {
        module.subTopicCount = module.subTopicCount + module.moduleTopics[j].subTopics.length
        module.approxTime = module.approxTime + module.moduleTopics[j].totalTime
        let topic: any = module?.moduleTopics[j]
        topic.isTopicClicked = false
        topic.isSelected = true
        topic.moduleID = module.id
        module.conceptCount = module.conceptCount + topic.topicProgrammingConcept.length

        // Iterate all sub topics.
        for (let k = 0; k < topic?.subTopics.length; k++) {
          // module.conceptCount = module.conceptCount + topic.subTopics[k].topicProgrammingConcept.length
          topic.subTopics[k].isSelected = true
          topic.subTopics[k].isSubTopicClicked = false
          topic.subTopics[k].moduleID = module.id
        }
      }
    }

    // If the first batch module is not the first in order for the batch, then get the previous batch module.
    if (this.batchModuleList[0].order > 1) {
      this.getBatchModule()
    }
  }

  // Create batch form.
  createBatchForm(): void {
    this.batchForm = this.formBuilder.group({
      id: new FormControl(),
      batchName: new FormControl(null),
      code: new FormControl(null),
      course: new FormControl(null),
      totalStudents: new FormControl(null),
      totalIntake: new FormControl(null),
      batchStatus: new FormControl("Upcoming"),
      batchObjective: new FormControl(null),
      isActive: new FormControl(true),
      isB2B: new FormControl(null),
      brochure: new FormControl(null),
      meetLink: new FormControl(null),
      startDate: new FormControl(null),
      eligibility: this.formBuilder.group({
        id: [],
        technologies: [null],
        studentRating: [null],
        experience: [null],
        academicYear: [null]
      }),
      faculty: new FormControl(null),
      batchTimings: new FormControl(null),
      salesPerson: new FormControl(null),
      requirement: new FormControl(null),
    })
    this.addFacultyValidators()
  }

  // Add validators.
  addFacultyValidators(): void {
    this.batchForm.get("batchStatus").setValidators([Validators.required])
    this.utilService.updateValueAndValiditors(this.batchForm)
  }

  // Create batch module timings form.
  createBatchModuleTimingsForm(): void {
    this.batchModuleTimingsForm = this.formBuilder.group({
      batchModuleTimings: new FormArray([]),
      isApplyToAllSessions: new FormControl(false),
      startDate: new FormControl(null),
    })
  }

  get batchModuleTimings(): FormArray {
    return this.batchModuleTimingsForm.get("batchModuleTimings") as FormArray
  }

  // Insert batch module timing form.
  insertBatchModuleTimingForm(day: any): any {
    return this.formBuilder.group({
      id: new FormControl(null),
      batchID: new FormControl(this.batchID),
      batchModuleID: new FormControl(this.selectedBatchModule.id),
      day: new FormControl(day),
      moduleID: new FormControl(this.selectedBatchModule.module.id),
      facultyID: new FormControl(this.facultyID),
      fromTime: new FormControl(null, [Validators.required]),
      toTime: new FormControl(null, [Validators.required]),
    })
  }

  // Add batch module timing form to batch module timing form.
  addBatchModuleTimingForm(): void {
    this.batchModuleTimings.push(this.formBuilder.group({
      id: new FormControl(null),
      batchID: new FormControl(this.batchID),
      batchModuleID: new FormControl(this.selectedBatchModule.id),
      day: new FormControl(null),
      moduleID: new FormControl(this.selectedBatchModule.module.id),
      facultyID: new FormControl(this.facultyID),
      fromTime: new FormControl(null, [Validators.required]),
      toTime: new FormControl(null, [Validators.required]),
    }))
  }

  // Update batch form.
  updateBatchForm(): void {
    this.createBatchForm()

    // Set empty eligibility.
    if (this.batch.eligibility == null) {
      this.batch.eligibility = {}
    }

    // Format start date.
    this.batch.startDate = this.datePipe.transform(this.batch.startDate, 'yyyy-MM-dd')
    this.batch.estimatedEndDate = this.datePipe.transform(this.batch.estimatedEndDate, 'yyyy-MM-dd')

    this.batchForm.patchValue(this.batch)
  }

  // Update batch module timings form.
  updateBatchModuleTimingsForm(): void {
    this.createBatchModuleTimingsForm()

    // Sort batch timings by order of days.
    this.batch.batchTimings.sort(this.sortBatchTimings)

    // Add batch module timings to batch module timing form.
    for (let i = 0; i < this.selectedBatchModule.moduleTimings.length; i++) {
      this.addBatchModuleTimingForm()
    }
    this.batchModuleTimings.patchValue(this.selectedBatchModule.moduleTimings)
    this.batchModuleTimingsForm.get('startDate').setValue(this.selectedBatchModule.startDate)

    // If it is the first batch module then make its start date compulsary.
    if (this.selectedBatchModule.isFirst) {
      this.batchModuleTimingsForm.get('startDate').setValidators([Validators.required])
      this.batchModuleTimingsForm.get('startDate').updateValueAndValidity()
    }

    this.getDayList()
  }

  // Sort batch timings.
  sortBatchTimings(a, b) {
    if (a.day.order < b.day.order) {
      return -1
    }
    if (a.day.order > b.day.order) {
      return 1
    }
    return 0
  }

  // Format the fields of day list.
  formatDayListFields(): void {
    for (let i = 0; i < this.dayList.length; i++) {
      this.dayList[i].isSelected = false
      for (let j = 0; j < this.batchModuleTimings.controls.length; j++) {
        if (this.dayList[i].id == this.batchModuleTimings.at(j).get('day').value.id) {
          this.dayList[i].isSelected = true
          continue
        }
      }
    }
  }

  // On clicking day button.
  onDayClick(day: any): void {

    // If day is selected then unselect it and remove it from batch timings form.
    if (day.isSelected) {
      if (confirm("Are you sure you want to delete the batch timing?")) {
        for (let i = 0; i < this.batchModuleTimings.controls.length; i++) {
          if (day.id == this.batchModuleTimings.at(i).get('day').value.id) {
            this.batchModuleTimings.removeAt(i)

            // If removed batch timing is first index, then make is apply to all sessions false and enable all other batch timings.
            if (i == 0 && this.batchModuleTimings.length > 1) {
              this.batchModuleTimingsForm.get('isApplyToAllSessions').setValue(false)
              for (let j = 0; j < this.batchModuleTimings.length; j++) {
                this.batchModuleTimings.at(j).get('fromTime').enable()
                this.batchModuleTimings.at(j).get('toTime').enable()
              }
            }
            break
          }
        }
        this.batchModuleTimingsForm.markAsDirty()
        day.isSelected = false
      }
      return
    }

    // If day is unselected then select it and add it to batch timings form.
    day.isSelected = true
    let index: number = -1
    for (let i = 0; i < this.batchModuleTimings.controls.length; i++) {
      if (day.order > this.batchModuleTimings.at(i).value.day.order) {
        index = i
      }
    }
    this.batchModuleTimings.insert(index + 1, this.insertBatchModuleTimingForm(day))

    // If apply to all is true and new batch timing is inserted then set the value and disable it.
    // If added batch timing is not in first place then set its time as the first batch timing's value and disable it.
    if (index != -1 && this.batchModuleTimings.length > 1 && this.batchModuleTimingsForm.get('isApplyToAllSessions').value) {
      this.batchModuleTimings.at(index + 1).get('fromTime').setValue(this.batchModuleTimings.at(0).get('fromTime').value)
      this.batchModuleTimings.at(index + 1).get('toTime').setValue(this.batchModuleTimings.at(0).get('toTime').value)
      this.batchModuleTimings.at(index + 1).get('fromTime').disable()
      this.batchModuleTimings.at(index + 1).get('toTime').disable()
    }

    // If added batch timing is in first place then set first batch timing's value as the second batch timing's value
    // and disable second batch timing.
    if (index == -1 && this.batchModuleTimings.length > 1 && this.batchModuleTimingsForm.get('isApplyToAllSessions').value) {
      this.batchModuleTimings.at(0).get('fromTime').setValue(this.batchModuleTimings.at(1).get('fromTime').value)
      this.batchModuleTimings.at(0).get('toTime').setValue(this.batchModuleTimings.at(1).get('toTime').value)
      this.batchModuleTimings.at(1).get('fromTime').disable()
      this.batchModuleTimings.at(1).get('toTime').disable()
    }
  }

  // Delete batch timing from batch form.
  deleteBatchTimingInBatchForm(index: number): void {
    if (confirm("Are you sure you want to delete the batch timing?")) {
      this.batchModuleTimings.removeAt(index)

      // If removed batch timing is first index, then make is apply to all sessions false and enable all other batch timings.
      if (index == 0 && this.batchModuleTimings.length > 1) {
        this.batchModuleTimingsForm.get('isApplyToAllSessions').setValue(false)
        for (let i = 0; i < this.batchModuleTimings.length; i++) {
          this.batchModuleTimings.at(i).get('fromTime').enable()
          this.batchModuleTimings.at(i).get('toTime').enable()
        }
      }
      this.batchModuleTimingsForm.markAsDirty()
      this.formatDayListFields()
    }
  }

  // On checking apply it to all sessions.
  onApplyToAllSessionsClick(): void {

    // If checkbox is unselected then enable all the batch timings.
    if (this.batchModuleTimingsForm.get('isApplyToAllSessions').value == false) {
      for (let i = 1; i < this.batchModuleTimings.length; i++) {
        this.batchModuleTimings.at(i).get('fromTime').enable()
        this.batchModuleTimings.at(i).get('toTime').enable()
      }
      return
    }

    // If first session day timings are empty then give alert.
    if (this.batchModuleTimings.at(0).value.fromTime == null || this.batchModuleTimings.at(0).value.toTime == null
      || this.batchModuleTimings.at(0).value.fromTime == "" || this.batchModuleTimings.at(0).value.toTime == "") {
      alert("Please fill in the from time and to time of the first session day")
      this.batchModuleTimingsForm.get('isApplyToAllSessions').setValue(false)
      return
    }

    // If start time and end time are same.
    if (this.batchModuleTimings.at(0).value.fromTime == this.batchModuleTimings.at(0).value.toTime) {
      alert("Start time and end time cannot be same")
      this.batchModuleTimingsForm.get('isApplyToAllSessions').setValue(false)
      return
    }

    // If there are more than one batch timings then set all the other batch timings as the value for the first 
    // batch timing.
    if (this.batchModuleTimings.length > 1) {
      for (let i = 1; i < this.batchModuleTimings.length; i++) {
        this.batchModuleTimings.at(i).get('fromTime').setValue(this.batchModuleTimings.at(0).get('fromTime').value)
        this.batchModuleTimings.at(i).get('toTime').setValue(this.batchModuleTimings.at(0).get('toTime').value)
        this.batchModuleTimings.at(i).get('fromTime').disable()
        this.batchModuleTimings.at(i).get('toTime').disable()
      }
    }
  }

  // On changing value of batch time.
  onBatchTimeChange(index: number): void {
    if (index == 0 && this.batchModuleTimings.length > 1 && this.batchModuleTimingsForm.get('isApplyToAllSessions').value) {
      for (let i = 1; i < this.batchModuleTimings.length; i++) {
        this.batchModuleTimings.at(i).get('fromTime').setValue(this.batchModuleTimings.at(0).get('fromTime').value)
        this.batchModuleTimings.at(i).get('toTime').setValue(this.batchModuleTimings.at(0).get('toTime').value)
      }
    }
  }

  // Check if batch module timings form is valid.
  isBatchModuleTimingsFormValid(): boolean {

    // If batch module timings is empty.
    if (this.batchModuleTimings.length == 0) {
      return false
    }

    // If from time and to time are same for any batch timing then return.
    for (let i = 0; i < this.batchModuleTimings.length; i++) {
      if (this.batchModuleTimings.at(i).get('fromTime').value == this.batchModuleTimings.at(i).get('toTime').value) {
        this.batchModuleTimingsForm.markAllAsTouched()
        return false
      }
    }

    // If batch module timings is invalid.
    if (this.batchModuleTimingsForm.invalid) {
      this.batchModuleTimingsForm.markAllAsTouched()
      return false
    }

    return true
  }

  // On clicking batch module tab for setting batch module timings.
  onBatchModuleTabClickForBatchModuleTimings(event: any): void {

    // Set the batch timings for the current module.
    if (this.selectedBatchModule) {
      this.batchModuleTimingsForm.enable()
      for (let i = 0; i < this.batchModuleList.length; i++) {
        if (this.selectedBatchModule.id == this.batchModuleList[i].id) {
          this.batchModuleList[i].moduleTimings = this.batchModuleTimings.value
          this.batchModuleList[i].startDate = this.batchModuleTimingsForm.get('startDate').value
        }
      }
    }

    // Check if the module batch module is valid, if not then set its incomplete field to false.
    if (this.batchModuleTabList.length > 0 && this.selectedBatchModule) {
      if (this.selectedBatchModule.moduleTimings?.length == 0 || ((this.selectedBatchModule.startDate == null || this.selectedBatchModule.startDate == "")
        && this.selectedBatchModule.isFirst) || this.checkBatchModuleTimingsFromTimeToTomeIsSame(this.selectedBatchModule)) {
        this.setIncompleteFieldForbBatchModuleTabList(this.selectedBatchModule.module.moduleName, false)
      }
      else {
        this.setIncompleteFieldForbBatchModuleTabList(this.selectedBatchModule.module.moduleName, true)
      }
    }

    // Change the selected batch module to the new index selected.
    this.selectedBatchModule = this.batchModuleList[event.index]

    // Only give start date compulsary validation for the first batch module.
    if (event.index == 0) {
      this.selectedBatchModule.isFirst = true
    }

    this.updateBatchModuleTimingsForm()
  }

  // Validate all batch related form.
  validateAllBatchRelatedForms(): void {

    let errorString: string = "Modules incomplete: "
    let moduleNameList: string[] = []

    // Validate batch module timings form for the current selected batch module.
    if (this.selectedBatchModule && !this.isBatchModuleTimingsFormValid()) {
      return
    }

    // Set the batch timings for the current module if batch timings form is valid.
    if (this.selectedBatchModule) {
      this.batchModuleTimingsForm.enable()
      for (let i = 0; i < this.batchModuleList.length; i++) {
        if (this.selectedBatchModule.id == this.batchModuleList[i].id) {
          this.batchModuleList[i].moduleTimings = this.batchModuleTimings.value
          this.batchModuleList[i].startDate = this.batchModuleTimingsForm.get('startDate').value
        }
      }
    }

    // Check if any batch module does not have batch timings, start date or has same from time and to time.
    for (let i = 0; i < this.batchModuleList.length; i++) {
      if (this.batchModuleList[i].moduleTimings?.length == 0 || ((this.batchModuleList[i].startDate == null || this.batchModuleList[i].startDate == "")
        && this.batchModuleList[i].isFirst) || this.checkBatchModuleTimingsFromTimeToTomeIsSame(this.batchModuleList[i])) {
        this.setIncompleteFieldForbBatchModuleTabList(this.batchModuleList[i].module.moduleName, false)
        moduleNameList.push(this.batchModuleList[i].module.moduleName)
      }
      else {
        this.setIncompleteFieldForbBatchModuleTabList(this.selectedBatchModule.module.moduleName, true)
      }
    }

    // Select the first incomplete batch module.
    this.setSelectedBatchModuleIndex()

    // If any batch module is incomplete then give alert and return.
    if (moduleNameList.length > 0) {
      errorString = errorString + moduleNameList.join(", ")
      alert(errorString)
      return
    }

    // If batch form is invalid then return.
    if (this.batchForm.invalid) {
      this.batchForm.markAllAsTouched()
      return
    }

    // Set the number of updates to be done.
    this.batchUpdateCount = this.batchModuleList.length

    // Update batch.
    this.updateBatch()
  }

  // Set the incomplete field for batch module tab list.
  setIncompleteFieldForbBatchModuleTabList(moduleName: string, isComplete: boolean): void {
    for (let i = 0; i < this.batchModuleTabList.length; i++) {
      if (this.batchModuleTabList[i].moduleName == moduleName) {
        this.batchModuleTabList[i].isComplete = isComplete
      }
    }
  }

  // Set the selected batch module index.
  setSelectedBatchModuleIndex(): void {
    for (let j = 0; j < this.batchModuleTabList.length; j++) {
      if (!this.batchModuleTabList[j].isComplete) {
        this.selectedBatchModuleIndex = j
        return
      }
    }
  }

  // Check if module's batch timings' from time and to time are same or not.
  checkBatchModuleTimingsFromTimeToTomeIsSame(batchModule: any): boolean {
    for (let i = 0; i < batchModule.moduleTimings.length; i++) {
      if (batchModule.moduleTimings[i].fromTime == batchModule.moduleTimings[i].toTime) {
        return true
      }
    }
    return false
  }

  // Extract ID from objects and delete objects before adding or updating.
  patchIDFromObjects(batchModule: any): void {
    batchModule.facultyID = batchModule.faculty.id
    batchModule.moduleID = batchModule.module.id
    if (!batchModule.isFirst && batchModule.startDate == "") {
      batchModule.startDate = null
    }
    for (let i = 0; i < batchModule.moduleTimings?.length; i++) {
      batchModule.moduleTimings[i].dayID = batchModule.moduleTimings[i].day.id
    }
  }

  // Update batch.
  updateBatch(): void {
    this.batchForm.enable()
    this.setFormFields()
    let batch = this.batchForm.value
    this.utilService.deleteNullValuePropertyFromObject(batch.eligibility)
    if (batch.eligibility) {
      this.deleteEligibilityIfEmpty()
    }
    this.spinnerService.loadingMessage = "Updating batch"
    this.batchService.updateBatch(batch).subscribe((respond: any) => {
      this.updateBatchModuleList()
    }, (err) => {
      if (err.statusText.includes('Unknown')) {
        return
      }
      console.error(err)
    })
  }

  // Update batch module.
  updateBatchModule(batchModule: any): void {
    this.spinnerService.loadingMessage = "Updating Batch"
    this.patchIDFromObjects(batchModule)
    this.batchService.updateBatchModule(this.batchID, batchModule).subscribe((respond: any) => {
    }, (err) => {
      if (err.statusText.includes('Unknown')) {
        return
      }
      this.batchUpdateErrorCount = this.batchUpdateErrorCount + 1
      console.error(err)
    }).add(() => {
      this.batchUpdateCount = this.batchUpdateCount - 1
      if (this.batchUpdateCount == 0) {
        this.showBatchUpdateSucessOrError()
      }
    })
  }

  // Show if batch and batch modules are successfully updated or not.
  showBatchUpdateSucessOrError(): void {

    // If error count is more than 0 then show failure alert.
    if (this.batchUpdateErrorCount > 0) {
      alert("Batch could not be updated")
      return
    }
    this.getBatch()
    this.getBatchModuleList()
    this.goToStepThree()
  }

  // Update batch module list.
  updateBatchModuleList(): void {
    for (let i = 0; i < this.batchModuleList.length; i++) {
      this.updateBatchModule(this.batchModuleList[i])
    }
  }

  // Delete eligibility from batch form if not present.
  deleteEligibilityIfEmpty(): void {
    if (Object.keys(this.batch?.eligibility)?.length === 0 && this.batch?.eligibility?.constructor === Object) {
      delete this.batch.eligibility
    }
  }

  // Set dates in proper format for batch.
  setFormFields(): void {
    if (this.batchForm.get("startDate")?.value) {
      this.batch.startDate = this.datePipe.transform(this.batch.startDate, 'yyyy-MM-dd')
    }
    if (this.batchForm.get("estimatedEndDate")?.value) {
      this.batch.estimatedEndDate = this.datePipe.transform(this.batch.estimatedEndDate, 'yyyy-MM-dd')
    }
  }

  // Go tp step three.
  goToStepThree(): void {
    this.currentStep = this.currentStep + 1
    this.showCurrentStepTemplate()
  }

  // On dragging and dropping modules table rows.
  onDropModulesTableRow(event: CdkDragDrop<string[]>) {
    moveItemInArray(this.batchModuleList, event.previousIndex, event.currentIndex)
    this.batchModuleList.forEach((module, index) => {
      module.order = index + 1
    })
  }

  //********************************************* STEP THREE FUNCTIONS ************************************************************

  // Go to step four.
  goToStepFour(): void {
    this.batchModuleTabList = this.formatBatchModuleTabList(this.batchModuleTabList, this.batchModuleList, true)
    this.currentStep = this.currentStep + 1
    this.showCurrentStepTemplate()
  }

  //********************************************* STEP FOUR FUNCTIONS ************************************************************

  // Format the batch modules tab list.
  formatBatchModuleTabList(tempBatchModuleTabList: any[], tempBatchModuleList: any[], isBatchModuleToBeSelected: boolean): any[] {

    // Create course module tab list.
    tempBatchModuleTabList = []
    for (let i = 0; i < tempBatchModuleList.length; i++) {
      tempBatchModuleTabList.push(
        {
          moduleName: tempBatchModuleList[i].module.moduleName,
          module: tempBatchModuleList[i].module,
        }
      )

      if (tempBatchModuleList[i].module.logo) {
        tempBatchModuleTabList[i].imageURL = tempBatchModuleList[i].module.logo
      }
      else {
        tempBatchModuleTabList[i].imageURL = "assets/icon/grey-icons/Score.png"
      }

      tempBatchModuleTabList[i].isComplete = true
      if (tempBatchModuleList[i].moduleTimings.length == 0 || tempBatchModuleList[i].startDate == null) {
        tempBatchModuleTabList[i].isComplete = false
      }
    }

    if (tempBatchModuleList.length > 0 && isBatchModuleToBeSelected) {
      this.selectedBatchModule = tempBatchModuleList[0].module
    }

    return tempBatchModuleTabList
  }

  // On clicking module tab.
  onModuleTabClick(event: any): void {

    // Till step 4 batch module list.
    if (this.currentStep <= 4) {
      for (let i = 0; i < this.batchModuleList.length; i++) {
        if (i == event.index) {
          this.selectedBatchModule = this.batchModuleList[i].module
        }
      }
    }

    // After step 4 batch module list.
    if (this.currentStep > 4) {
      for (let i = 0; i < this.batchModuleList.length; i++) {
        if (i == event.index) {
          this.selectedBatchModule = this.selectedBatchModuleList[i].module
        }
      }
    }
  }

  // On changing isSelected of topic.
  onTopicSelectedChange(topic: any, subTopic: any): void {
    let selectedSubTopicsCount: number = 0
    for (let i = 0; i < topic.subTopics?.length; i++) {

      // If topic's checkbox then make all sub topics isSelected same as topic's isSelected. 
      if (subTopic == null) {
        topic.subTopics[i].isSelected = topic.isSelected
      }

      // If sub topic's checkbox then count the number of checked sub topics.
      if (subTopic != null && topic.subTopics[i].isSelected) {
        selectedSubTopicsCount = selectedSubTopicsCount + 1
      }
    }

    // If sub topic's checkbox and no checkbox selected then make isSelected of topic false.
    if (subTopic != null && selectedSubTopicsCount == 0) {
      topic.isSelected = false
    }

    // If sub topic's checkbox and all checkboxes are selected then make isSelected of topic true.
    if (subTopic != null && selectedSubTopicsCount <= topic.subTopics?.length && selectedSubTopicsCount > 0) {
      topic.isSelected = true
    }
  }

  // Create selected course module list with all the selected topics and subtopics.
  createSelectedCourseModuleList(): boolean {
    let unselectedModuleList: string[] = []
    this.selectedBatchModuleList = []
    let tempBatchModuleList: any[] = JSON.parse(JSON.stringify(this.batchModuleList))

    // Iterate all modules.
    for (let i = 0; i < tempBatchModuleList.length; i++) {
      let module: any = tempBatchModuleList[i].module
      let selectedCourseModule: any
      let selectedTopics: any[] = []

      // Iterate all topics.
      for (let j = 0; j < module.moduleTopics?.length; j++) {
        let topic: any = module.moduleTopics[j]
        topic.isTopicClicked = false

        // // If sub topic is selected then give order to selected sub topics.
        // if (topic.isSelected){
        //   selectedTopics.push(topic)
        //   continue
        // }

        // If topic is not selected.
        if (topic.isSelected) {
          let selectedSubTopics: any[] = []

          // Iterate all the sub topics.
          for (let k = 0; k < topic.subTopics?.length; k++) {

            // If sub topic is selected then push it to selected sub topics.
            if (topic.subTopics[k].isSelected) {
              selectedSubTopics.push(topic.subTopics[k])
            }
          }

          // If sub topics are there then push the topic.
          if (selectedSubTopics.length > 0) {
            selectedTopics.push(topic)

            // Give the selected sub topics to topic.
            selectedTopics[selectedTopics.length - 1].subTopics = selectedSubTopics
          }
        }
      }

      // If selected topics is present then add the module to selected course modules.
      if (selectedTopics.length > 0) {
        selectedCourseModule = tempBatchModuleList[i]
        selectedCourseModule.module.moduleTopics = selectedTopics
        this.selectedBatchModuleList.push(selectedCourseModule)
      }

      // If no topic selected then give alert.
      if (selectedTopics.length == 0) {
        unselectedModuleList.push(module.moduleName)
        selectedCourseModule = tempBatchModuleList[i]
        selectedCourseModule.module.moduleTopics = selectedTopics
        this.selectedBatchModuleList.push(selectedCourseModule)
      }
    }

    // If there are unselected modules then give alert.
    if (unselectedModuleList.length > 0) {
      alert("Please select atleast one sub topic from the modules :" + unselectedModuleList.join(", "))
      return false
    }
    return true
  }

  // Go to step three from step four.
  goToStepThreeFromStepFour(): void {
    this.selectedBatchModule = this.batchModuleList[0]
    this.goBackOneStep()
  }

  // Go to step five.
  goToStepFive(): void {
    if (this.createSelectedCourseModuleList()) {
      this.selectedCourseModuleTabList = this.formatBatchModuleTabList(this.selectedCourseModuleTabList, this.selectedBatchModuleList, true)
      this.currentStep = this.currentStep + 1
      this.showCurrentStepTemplate()
    }
  }

  //********************************************* STEP FIVE FUNCTIONS ************************************************************

  // On dragging and dropping topics table rows.
  onDropTopicsTableRow(event: CdkDragDrop<string[]>, tempTopics: any[]) {
    moveItemInArray(tempTopics, event.previousIndex, event.currentIndex)
    tempTopics.forEach((topic, index) => {
      topic.order = index + 1
    })
  }

  // Set the order for all the sub topics in selected course module list.
  SetOrderAndCreateSessionPlanList(): void {

    // let subTopicOrder: number  = 0
    this.batchSessionPlanEntryList = []

    // Iterate all modules.
    for (let i = 0; i < this.selectedBatchModuleList.length; i++) {
      let module: any = this.selectedBatchModuleList[i].module
      let subTopicOrder: number = 0
      // let seletedModuleSubTopics: any[] = []

      // Iterate all topics.
      for (let j = 0; j < module.moduleTopics?.length; j++) {
        let topic: any = module.moduleTopics[j]
        topic.isAllQuestionsSelected = true

        // Iterate all topic programming questions and set their is selected as false.
        for (let l = 0; l < topic.topicProgrammingQuestions?.length; l++) {
          topic.topicProgrammingQuestions[l].isSelected = true
          if (topic.topicProgrammingQuestions[l]?.programmingQuestion?.level == 1) {
            topic.topicProgrammingQuestions[l].programmingQuestion.difficulty = "Easy"
          }
          if (topic.topicProgrammingQuestions[l]?.programmingQuestion?.level == 2) {
            topic.topicProgrammingQuestions[l].programmingQuestion.difficulty = "Medium"
          }
          if (topic.topicProgrammingQuestions[l]?.programmingQuestion?.level == 3) {
            topic.topicProgrammingQuestions[l].programmingQuestion.difficulty = "Hard"
          }
          // topic.topicProgrammingQuestions[l].days = 1
        }

        // Iterate all sub topics and give order.
        for (let k = 0; k < topic.subTopics?.length; k++) {
          let subTopic: any = topic.subTopics[k]
          subTopicOrder = subTopicOrder + 1
          // subTopic.parentTopicName = topic.topicName
          // subTopic.isAllQuestionsSelected = true

          // // Iterate all topic programming questions and set their is selected as false.
          // for (let l = 0; l < subTopic.topicProgrammingQuestions?.length; l++){
          //   subTopic.topicProgrammingQuestions[l].isSelected = true
          //   if (subTopic.topicProgrammingQuestions[l]?.programmingQuestion?.level == 1){
          //     subTopic.topicProgrammingQuestions[l].programmingQuestion.difficulty = "Easy"
          //   }
          //   if (subTopic.topicProgrammingQuestions[l]?.programmingQuestion?.level == 2){
          //     subTopic.topicProgrammingQuestions[l].programmingQuestion.difficulty = "Medium"
          //   }
          //   if (subTopic.topicProgrammingQuestions[l]?.programmingQuestion?.level == 3){
          //     subTopic.topicProgrammingQuestions[l].programmingQuestion.difficulty = "Hard"
          //   }
          //   subTopic.topicProgrammingQuestions[l].days = 1
          // }

          // Push all subtopics for a module in an array.
          // seletedModuleSubTopics.push(subTopic)

          this.batchSessionPlanEntryList.push({
            batchID: this.batchID,
            moduleID: subTopic.moduleID,
            topicID: subTopic.topicID,
            subTopicID: subTopic.id,
            order: subTopicOrder,
            totalTime: subTopic.totalTime,
            isCompleted: false,
            // // Temp!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
            // name: subTopic.topicName
            // // Temp!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
          })
        }
      }

      // Give the selected sub topics to selected course module.
      // this.selectedCourseModuleList[i].module.subTopics = seletedModuleSubTopics
    }
  }

  // Go to step six.
  goToStepSix(): void {
    this.SetOrderAndCreateSessionPlanList()
    this.selectedBatchModule = this.selectedBatchModuleList[0].module
    this.currentStep = this.currentStep + 1
    this.showCurrentStepTemplate()
  }

  // Got to step four from step five.
  goToStepFourFromStepFive(): void {
    this.batchModuleTabList = this.formatBatchModuleTabList(this.batchModuleTabList, this.batchModuleList, true)
    this.goBackOneStep()
  }

  //********************************************* STEP SIX FUNCTIONS ************************************************************

  // On changing isSelected of topic programming question.
  onProgrammingQuestionSelectedChange(topic: any, topicProgrammingQuestionIndex: number): void {
    let selectedProgrammingQuestionsCount: number = 0
    let topicProgrammingQuestion: any

    // If topic checkbox then change its value.
    if (topicProgrammingQuestionIndex == -1) {
      topic.isAllQuestionsSelected = !topic.isAllQuestionsSelected
    }

    // If programming question checkbox then change its value.
    if (topicProgrammingQuestionIndex > -1) {
      topicProgrammingQuestion = topic.topicProgrammingQuestions[topicProgrammingQuestionIndex]
      topicProgrammingQuestion.isSelected = !topicProgrammingQuestion.isSelected
    }

    for (let i = 0; i < topic.topicProgrammingQuestions?.length; i++) {

      // If topic's checkbox then make all programming questions isSelected same as topic's isAllQuestionsSelected. 
      if (topicProgrammingQuestionIndex == -1) {
        topic.topicProgrammingQuestions[i].isSelected = topic.isAllQuestionsSelected
      }

      // If programming question's checkbox then count the number of checked programming questions.
      if (topicProgrammingQuestionIndex > -1 && topic.topicProgrammingQuestions[i].isSelected) {
        selectedProgrammingQuestionsCount = selectedProgrammingQuestionsCount + 1
      }
    }

    // If programming question's checkbox and no checkbox selected then make isAllQuestionsSelected of topic false.
    if (topicProgrammingQuestionIndex > -1 && selectedProgrammingQuestionsCount == 0) {
      topic.isAllQuestionsSelected = false
    }

    // If programming question's checkbox and all checkboxes are selected then make isAllQuestionsSelected of topic true.
    if (topicProgrammingQuestionIndex > -1 && selectedProgrammingQuestionsCount <= topic.topicProgrammingQuestions?.length &&
      selectedProgrammingQuestionsCount > 0) {
      topic.isAllQuestionsSelected = true
    }
  }

  // Create session plan.
  createSessionPlan(): void {
    this.createBatchTopicAssignmentEntryList()
    this.addBatchSessionPlan()
  }

  // Create batch topic assignment list. 
  createBatchTopicAssignmentEntryList(): void {
    this.batchTopicAssignmentEntryList = []

    // Iterate all selected course modules.
    for (let i = 0; i < this.selectedBatchModuleList.length; i++) {
      let module: any = this.selectedBatchModuleList[i].module

      // Iterate all topics.
      for (let j = 0; j < module.moduleTopics?.length; j++) {
        let topic: any = module.moduleTopics[j]

        // Iterate all programming questions.
        for (let k = 0; k < topic.topicProgrammingQuestions.length; k++) {

          // Create batch topic assignment.
          if (topic.topicProgrammingQuestions[k].isSelected) {

            // If days of topicProgrammingQuestionsis less than 0 than give alert
            // if (subTopic.topicProgrammingQuestions[k].days < 1){
            //   alert("Required days required for assignment cannot be lesser than 1")
            //   return
            // }

            // Create batch topic assignment entry. 
            this.batchTopicAssignmentEntryList.push({
              batchID: this.batchID,
              topicID: topic.id,
              programmingQuestionID: topic.topicProgrammingQuestions[k].programmingQuestion.id,
              // days: subTopic.topicProgrammingQuestions[k].days,
              moduleID: module.id
              // // Temp!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
              // programmingQuestionName: subTopic.topicProgrammingQuestions[k].programmingQuestion.label
              // // Temp!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
            })
          }
        }
      }
    }
    this.batchTopicAssignmentAddCount = this.batchTopicAssignmentEntryList.length
  }

  // Add all batch topic assignment entries.
  addBatchTopicAssignmentEntryList(): void {
    for (let i = 0; i < this.batchTopicAssignmentEntryList.length; i++) {
      this.addBatchTopicAssignment(this.batchTopicAssignmentEntryList[i])
    }
  }

  // Add batch session plan.
  addBatchSessionPlan(): void {
    this.spinnerService.loadingMessage = "Creating batch session plan"
    let queryParams: any
    if (this.isFaculty) {
      queryParams = {
        facultyID: this.loginID
      }
    }
    this.batchSessionService.addBatchSessionPlan(this.batchID, this.batchSessionPlanEntryList, queryParams).subscribe((respond: any) => {

      // If there are no batch topic assignments then show batch session plan.
      if (this.batchTopicAssignmentEntryList.length == 0) {
        this.getBatchSessionList()
        alert("Session plan successfully created")
      }

      // If there are batch topic assignments then add them.
      if (this.batchTopicAssignmentEntryList.length > 0) {
        this.addBatchTopicAssignmentEntryList()
      }
    }, (err) => {
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(this.utilService.getErrorString(err))
      console.error(err)
    })
  }

  // Add batch topic assignment.
  addBatchTopicAssignment(batchTopicAssignment: any): void {
    this.spinnerService.loadingMessage = "Creating batch session plan"
    this.batchTopicAssignmentService.addBatchTopicAssignment(this.batchID, batchTopicAssignment.topicID, batchTopicAssignment).subscribe((respond: any) => {
    }, (err) => {
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(this.utilService.getErrorString(err))
      console.error(err)
    }).add(() => {
      this.decrementBatchTopicAssignmentAddCount()
    })
  }

  // Decrement batchTopicAssignmentAddCount and get session plan when it is 0.
  decrementBatchTopicAssignmentAddCount(): void {
    this.batchTopicAssignmentAddCount = this.batchTopicAssignmentAddCount - 1
    if (this.batchTopicAssignmentAddCount == 0) {
      this.getBatchSessionList()
      alert("Session plan successfully created")
    }
  }

  // Got to step five from step six.
  goToStepFiveFromStepSix(): void {
    this.selectedBatchModule = this.selectedCourseModuleTabList[0].module
    this.goBackOneStep()
  }

  //********************************************* OTHER FUNCTIONS ************************************************************

  // Selecting the template to be shown.
  showCurrentStepTemplate(): void {
    this.sessionPlanStepsList[this.currentStep - 1]
    this.selectedTemplateName = this.sessionPlanStepsList[this.currentStep - 1].templateName
  }

  // Compare two Object.
  compareFn(ob1: any, ob2: any): any {
    if (ob1 == null && ob2 == null) {
      return true
    }
    return ob1 && ob2 ? ob1.id == ob2.id : ob1 == ob2
  }

  // Go one step back.
  goBackOneStep(): void {
    this.currentStep = this.currentStep - 1
    this.showCurrentStepTemplate()
  }

  //********************************************* GET FUNCTIONS ************************************************************

  // Get all components.
  getAllComponents(): void {
    this.getBatchStatusList()
    this.getFacultyList()
  }

  // Get batch.
  getBatch(): void {
    this.spinnerService.loadingMessage = "Getting Batch Details..."
    this.batchService.getBatch(this.batchID).subscribe((response) => {
      this.batch = response.body
      this.getBatchSessionList()
      this.updateBatchForm()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get batch session list.
  getBatchSessionList(): void {
    this.spinnerService.loadingMessage = "Getting batch session plan..."
    let queryParams: any
    if (this.isFaculty) {
      queryParams = {
        facultyID: this.isFaculty ? this.loginID : null
      }
    }
    this.batchSessionService.getBatchSessionDates(this.batchID, queryParams).subscribe((response) => {
      this.batchSessionList = response.body

      // If batch session plan does not exist then create batch session plan.
      if (this.batchSessionList.length == 0) {
        this.showCurrentStepTemplate()
      }

      // If batch session plan exists then show the session plan.
      if (this.batchSessionList.length > 0) {
        this.currentStep = 7
        this.showCurrentStepTemplate()
      }
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get batch module list.
  getBatchModuleList(): void {
    // if (this.batchModuleList.length > 0) {
    //   return
    // }
    this.spinnerService.loadingMessage = "Getting Batch Modules..."
    // let queryParams: any = {
    //   limit: this.limitModules,
    //   offset: this.offsetModule
    // } 
    let queryParams: any = {
      limit: -1,
      offset: 0,
      facultyID: this.facultyID
    }
    console.log("queryParams batch topic", queryParams)
    this.batchService.getBatchModulesWithAllFields(this.batchID, queryParams).subscribe((response) => {
      this.batchModuleList = response.body
      console.log("batchModuleList", this.batchModuleList
      )
      this.totalModules = parseInt(response.headers.get('X-Total-Count'))
      this.batchModuleTabList = this.formatBatchModuleTabList(this.batchModuleTabList, this.batchModuleList, false)
      this.formatBatchModuleList()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get batch module.
  getBatchModule(): void {
    this.spinnerService.loadingMessage = "Getting Batch Modules..."
    let queryParams: any = {
      order: (this.batchModuleList[0].order - 1)
    }
    this.batchService.getBatchModules(this.batchID, queryParams).subscribe((response) => {
      this.previousBatchModule = response.body
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get batch status list.
  getBatchStatusList(): void {
    this.generalService.getGeneralTypeByType("batch_status").subscribe((respond: any) => {
      this.batchStatusList = respond
    }, (err) => {
      console.error(err.error.error)
    })
  }

  // Get faculty list.
  getFacultyList(): void {
    this.generalService.getFacultyList().subscribe((data: any) => {
      this.facultyList = this.facultyList.concat(data.body)
    }, (err) => {
      console.error(err)
    })
  }

  // Get day list.
  getDayList(): void {
    this.generalService.getDaysList().subscribe((response: any) => {
      this.dayList = response.body
      this.formatDayListFields()
    }, (err) => {
      console.error(err)
    })
  }

}
