import { DatePipe } from '@angular/common';
import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormArray, FormControl, Validators, AbstractControl, FormBuilder } from '@angular/forms';
import { IBatchModule, IBatch, BatchService } from 'src/app/service/batch/batch.service';
import { CourseModuleService, ICourseModule } from 'src/app/service/course-module/course-module.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { ActivatedRoute } from '@angular/router';
import { BatchSessionService } from 'src/app/service/batch-session/batch-session.service';
import { CdkDragDrop, moveItemInArray } from '@angular/cdk/drag-drop';
import { IBatchSessionModules } from 'src/app/models/batch/batch_session';
import { LocalService } from 'src/app/service/storage/local.service';
import { Role } from 'src/app/service/constant';


@Component({
  selector: 'app-batch-modules',
  templateUrl: './batch-modules.component.html',
  styleUrls: ['./batch-modules.component.css']
})
export class BatchModulesComponent implements OnInit {

  // course modules
  //courseModule
  isCourseModule: boolean
  courseModules: ICourseModule[]
  totalCourseModules: number
  moduleIDList: string[]
  daysList: any[]
  facultyList: any[];
  isCourseLoaded: boolean
  courseModuleTabList: any[]
  isReordering: boolean


  // batch modules
  batchModules: IBatchModule[]
  totalBatchModules: number
  batchModuleForm: FormGroup

  // module modal
  @Input() activeModuleTab: number

  //totalDaysSelected
  totalDaysSelected: number

  // batches
  totalBatch: number;
  batches: any[];
  // selectedBatch: any;
  batchID: string

  courseHeaderList: any[]

  moduleTimings: any;

  //multi-step-selector
  progress1: number = 0;
  progress2: number = 0;


  modalRef: NgbModalRef;
  entity: string
  waitingListMessage: string
  @Input() batch: any
  courseID: string
  batchName: string
  loginID: string


  // Course Module.
  courseModuleList: any[]
  selectedCourseModule: IBatchSessionModules[]
  allCourseModules: IBatchSessionModules[]
  tabIndex: number


  @ViewChild('deleteModal') deleteModal: any

  isSessionPlanModificationStarted: boolean
  selectedSubTopicsID: string[];
  selectedOrderedCourseModules: IBatchSessionModules[];
  isFaculty: boolean;
  isAdmin: boolean;

  constructor(
    private formBuilder: FormBuilder,
    public utilService: UtilityService,
    private courseModuleService: CourseModuleService,
    private datePipe: DatePipe,
    private generalService: GeneralService,
    private batchService: BatchService,
    private modalService: NgbModal,
    private spinnerService: SpinnerService,
    private route: ActivatedRoute,
    private batchSessionService: BatchSessionService,
    private localService: LocalService,
    private role: Role,



  ) {
    this.extractID()
    this.initializeVariables()
    this.getAllComponents()

  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  getAllComponents(): void {
    this.getFacultyList()
    this.getDaysList()
  }
  initializeVariables(): void {

    if (this.localService.getJsonValue("roleName") == this.role.ADMIN) {
      // this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.TRAINING_RESOURCE)
    }
    if (this.localService.getJsonValue("roleName") == this.role.FACULTY) {
      // this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.BANK_RESOURCE)
    }
    this.loginID = this.localService.getJsonValue("loginID")

    this.isFaculty = (this.localService.getJsonValue("roleName") == this.role.FACULTY ? true : false)
    this.isAdmin = (this.localService.getJsonValue("roleName") == this.role.ADMIN ? true : false)

    this.isCourseModule = true
    this.daysList = []
    this.facultyList = []
    this.courseModuleTabList = []
    this.isReordering = false
    this.isSessionPlanModificationStarted = false
    this.selectedSubTopicsID = []
    this.selectedCourseModule = []
    this.selectedOrderedCourseModules = []
    this.allCourseModules = []
    this.tabIndex = 0
    this.showModules()

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
  async showModules(): Promise<void> {
    try {
      this.createBatchModuleForm()
      this.moduleIDList = []
      this.courseModules = await this.getCourseModules(this.courseID)
      console.log(this.courseModules);
      this.batchModules = await this.getBatchModules(this.batchID)
      console.log("batchModules", this.batchModules)
      this.handleBatchModuleForm(true)
      // this.openModal(this.moduleModal, "xl")
    } catch (error) {
      console.log(error);
    }
  }
  createBatchModuleForm(): void {
    this.batchModuleForm = this.formBuilder.group({
      batchModules: new FormArray([])
    })
  }

  get batchModuleControlArray(): FormArray {
    return this.batchModuleForm.get("batchModules") as FormArray
  }

  addBatchModulesToForm(): void {
    this.batchModuleControlArray.push(this.formBuilder.group({
      id: new FormControl(null),
      batchID: new FormControl(this.batchID),
      moduleID: new FormControl(null, [Validators.required]),
      facultyID: new FormControl(null, [Validators.required]),
      order: new FormControl(null, [Validators.required, Validators.min(1)]),
      startDate: new FormControl(null),
      isCompleted: new FormControl(false),
      isMarked: new FormControl(false),
      isModuleTimingMarked: new FormControl(false),
      moduleTimings: new FormArray([]),
      isApplyToAllSessions: new FormControl(false),
      estimatedEndDate: new FormControl(null),

    }))
  }

  getModuleTimingControlArray(index: number): FormArray {
    return this.batchModuleControlArray.at(index).get("moduleTimings") as FormArray
  }

  addModuleTimingToForm(index: number): void {
    this.getModuleTimingControlArray(index).push(this.formBuilder.group({
      id: new FormControl(null),
      dayID: new FormControl(null),//required
      day: new FormControl(null),
      moduleID: new FormControl(null),
      facultyID: new FormControl(null),
      isMarked: new FormControl(null),
      fromTime: new FormControl(null),//required
      toTime: new FormControl(null),//required
    }))
  }

  // On checking apply it to all sessions.
  onApplyToAllModulesClick(index: number): void {
    this.batchModuleControlArray.at(index).get('isApplyToAllSessions')
      .setValue(!this.batchModuleControlArray.at(index).get('isApplyToAllSessions').value)

    // If checkbox is unselected then enable all the batch timings.
    if (this.batchModuleControlArray.at(index).get('isApplyToAllSessions').value == false) {
      // 	for (let i = 1; i < this.getModuleTimingControlArray(index).length; i++) {
      // 		this.getModuleTimingControlArray(index).at(i).get("toTime").enable()
      // 		this.getModuleTimingControlArray(index).at(i).get("fromTime").enable()
      // 	}
      return
    }

    var element = -1
    for (let k = 0; k < this.getModuleTimingControlArray(index).length; k++) {
      if (this.getModuleTimingControlArray(index).at(k).get("isMarked").value && element == -1) {
        element = k;
      }
    }
    console.log(element);
    console.log(this.getModuleTimingControlArray(index).at(element).get("fromTime").value);
    console.log(this.getModuleTimingControlArray(index).at(element).get("toTime").value);

    // If first session day timings are empty then give alert.
    if (this.getModuleTimingControlArray(index).at(element).get("fromTime").value == null
      || this.getModuleTimingControlArray(index).at(element).get("toTime").value == null) {
      this.batchModuleControlArray.at(index).get('isApplyToAllSessions').setValue(false)
      alert("Please fill in the from time and to time of the first session day")
      return
    }


    // If start time and end time are same.
    if (this.getModuleTimingControlArray(index).at(element).get("fromTime").value == this.getModuleTimingControlArray(index).at(element).get("toTime").value) {
      this.batchModuleControlArray.at(index).get('isApplyToAllSessions').setValue(false)
      alert("Start time and end time cannot be same")
      return
    }

    // If there are more than one batch timings then set all the other batch timings as the value for the first 
    // batch timing.
    // let moduleTimingsControl = (this.batchModuleControlArray.at(index).get("moduleTimings") as FormArray).controls
    if (this.getModuleTimingControlArray(index).length > 1) {
      for (let i = 0; i < this.getModuleTimingControlArray(index).length; i++) {
        if (i == element) {
          console.log(i);
          continue
        }
        this.getModuleTimingControlArray(index).at(i).get('fromTime').setValue(this.getModuleTimingControlArray(index).at(element).get('fromTime').value)
        this.getModuleTimingControlArray(index).at(i).get('toTime').setValue(this.getModuleTimingControlArray(index).at(element).get('toTime').value)
        // this.getModuleTimingControlArray(index).at(i).get('fromTime').disable()
        // this.getModuleTimingControlArray(index).at(i).get('toTime').disable()

      }
    }

  }

  // async getAllBatchModules(): Promise<void> {
  //   for (let i = 0; i < this.courseModuleList.length; i++) {
  //     this.allCourseModules[i] = this.courseModuleList[i].module

  //     try {
  //       this.allCourseModules[i].moduleTopics = await this.getModuleTopics(this.courseModuleList[i].module?.id)

  //       // this.processSubTopicsOrder()
  //     } catch (error) {
  //       console.error(error)
  //     }
  //   }
  //   console.log("this.allCourseModules", this.allCourseModules);

  // }
  // async processBatchSessions(): Promise<void> {

  //   let errors: string[] = []
  //   let batchSessions: any[] = []
  //   for (let index = 0; index < this.courseModuleList.length; index++) {
  //     let batchModule = this.courseModuleList[index].module;
  //     let subTopicOrder: number = 0
  //     for (let j = 0; j < batchModule?.moduleTopics?.length; j++) {
  //       let moduleTopics = batchModule.moduleTopics[j];
  //       for (let k = 0; k < moduleTopics?.subTopics.length; k++) {
  //         let subTopic = moduleTopics?.subTopics[k];

  //         console.log("subTopic -> ", subTopic);

  //         let batchSession: any = {
  //           id: subTopic.batchSessionTopicID,
  //           batchID: this.batchID,
  //           moduleID: subTopic.moduleID,
  //           topicID: subTopic.topicID,
  //           subTopicID: subTopic.id,
  //           order: ++subTopicOrder,
  //           initialDate: this.datePipe.transform(subTopic.batchSessionInitialDate, "yyyy-MM-dd"),
  //           completedDate: subTopic.batchSessionCompletionDate,
  //           isCompleted: subTopic.batchSessionTopicIsCompleted,
  //           batchSessionID: subTopic.batchSessionID,
  //           totalTime: subTopic.totalTime,
  //         }
  //         // if (subTopic?.batchSessionTopicID != null) {
  //         //   batchSession.id = subTopic.batchSessionTopicID
  //         // }
  //         batchSessions.push(batchSession)

  //       }
  //     }
  //   }
  //   try {
  //     console.log(batchSessions);

  //     await this.updateBatchSessions(batchSessions)
  //   } catch (error) {
  //     errors.push(error)
  //   }
  //   let msgString = ""
  //   if (errors.length > 0) {
  //     for (let index = 0; index < errors.length; index++) {
  //       msgString += (index == 0 ? "" : "\n") + errors[index]
  //     }
  //     alert(msgString)
  //     return
  //   }
  //   msgString = 'Session Plan updated successfully'
  //   alert(msgString)
  //   this.activeModuleTab=2

  // }

  // // On dragging and dropping topics table rows.
  // onDropTopicsTableRow(event: CdkDragDrop<string[]>, tempTopics: any[]) {
  //   moveItemInArray(tempTopics, event.previousIndex, event.currentIndex)
  //   tempTopics.forEach((topic, index) => {
  //     topic.order = index + 1
  //   })
  // }

  toggleModuleTimingValidators(i: number): void {
    this.batchModuleControlArray.at(i).get("isModuleTimingMarked").
      setValue(!this.batchModuleControlArray.at(i).get("isModuleTimingMarked").value)

    if (this.batchModuleControlArray.at(i).get("isModuleTimingMarked").value) {
      this.addModuleTimingsValidators(this.batchModuleControlArray.at(i).get("moduleTimings"))
      return
    }

    this.removeModuleTimingsValidators(this.batchModuleControlArray.at(i).get("moduleTimings"))
  }

  addModuleTimingsValidators(formControl: AbstractControl): void {

    let moduleTimingsControl = (formControl as FormArray).controls
    for (let index = 0; index < moduleTimingsControl.length; index++) {
      const element = moduleTimingsControl[index];

      if (element.get("isMarked").value) {
        element.get("dayID").setValidators(Validators.required)
        element.get("fromTime").setValidators(Validators.required)
        element.get("toTime").setValidators(Validators.required)
        this.utilService.updateValueAndValiditors(element as FormGroup)
      }
    }

  }

  getTotalDays(i: number): number {
    this.totalDaysSelected = 0
    let moduleTimingsControl = (this.batchModuleControlArray.at(i).get("moduleTimings") as FormArray).controls
    for (let index = 0; index < moduleTimingsControl.length; index++) {
      const element = moduleTimingsControl[index];
      // console.log(element);
      if (element.get("isMarked").value) {
        // console.log(this.totalDaysSelected);
        this.totalDaysSelected++
        element.get("dayID").setValidators(Validators.required)
        element.get("fromTime").setValidators(Validators.required)
        element.get("toTime").setValidators(Validators.required)
        this.utilService.updateValueAndValiditors(element as FormGroup)
      }
    }
    return this.totalDaysSelected
  }

  removeModuleTimingsValidators(formControl: AbstractControl): void {
    let moduleTimingsControl = (formControl as FormArray).controls
    for (let index = 0; index < moduleTimingsControl.length; index++) {
      const element = moduleTimingsControl[index];

      element.get("isMarked").setValue(false)
      element.get("dayID").setValidators(null)
      element.get("fromTime").setValidators(null)
      element.get("toTime").setValidators(null)
      this.utilService.updateValueAndValiditors(element as FormGroup)
    }
    this.totalDaysSelected = 0
  }

  onModuleTimeChange(i: number, j: number, moduleTimingForm: FormGroup): void {
    moduleTimingForm.get('isMarked').setValue(!moduleTimingForm.get('isMarked').value)

    if (j > 0 && moduleTimingForm.get('isMarked').value) {
      let fromTime = this.getModuleTimingControlArray(i).at(0).get("fromTime").value
      let toTime = this.getModuleTimingControlArray(i).at(0).get("toTime").value

      console.log("fromtime -> ", fromTime, " toTime -> ", toTime);
      moduleTimingForm.get("fromTime").setValue(fromTime)
      moduleTimingForm.get("toTime").setValue(toTime)
    }
  }

  deleteModuleTiming(i: number, j: number): void {
    this.getModuleTimingControlArray(i).at(j).get("isMarked").setValue(false)
  }

  handleBatchModuleForm(isCourseAdded: boolean): void {
    if (isCourseAdded) {
      this.addCourseModuleToForm()
    }

    this.addBatchModuleIDToList()
    this.disableBatchModulesInForm()
  }

  toggleModule(batchModuleForm: FormGroup): void {
    // console.log(batchModuleForm);
    batchModuleForm.disable()
    this.removeFormValueAndValidators(batchModuleForm)
    batchModuleForm.get("isMarked").setValue(!batchModuleForm.get("isMarked").value)

    if (batchModuleForm.get("isMarked").value) {
      this.setModuleFormValidators(batchModuleForm)
      batchModuleForm.enable()
      return
    }
  }

  removeFormValueAndValidators(batchModuleForm: FormGroup): void {
    batchModuleForm.get("facultyID").clearValidators()
    batchModuleForm.get("moduleID").clearValidators()
    batchModuleForm.get("order").clearValidators()
    batchModuleForm.get("startDate").clearValidators()
    batchModuleForm.get("estimatedEndDate").clearValidators()

    batchModuleForm.get("facultyID").setValue(null)
    batchModuleForm.get("order").setValue(null)
    batchModuleForm.get("estimatedEndDate").setValue(null)
    batchModuleForm.get("startDate").setValue(null)

    this.utilService.updateValueAndValiditors(batchModuleForm)
  }

  setModuleFormValidators(batchModuleForm: FormGroup): void {
    batchModuleForm.get("facultyID").setValidators([Validators.required])
    batchModuleForm.get("moduleID").setValidators([Validators.required])
    batchModuleForm.get("order").setValidators([Validators.required, Validators.min(1)])
    // batchModuleForm.get("startDate").setValidators([Validators.required])

    this.utilService.updateValueAndValiditors(batchModuleForm)
  }

  async getModuleTopics(moduleID: any): Promise<any> {
    try {
      return await new Promise<any>((resolve, reject) => {
        this.spinnerService.loadingMessage = "Getting modules";

        this.batchSessionService.getModuleTopics(this.batchID, moduleID).subscribe((response) => {
          // console.log(response.body);
          resolve(response.body)
        },
          (err: any) => {
            console.error(err);
            if (err.statusText.includes('Unknown')) {
              reject("No connection to server. Check internet.");
              return;
            }
            reject(err.error.error);
          });

      })

    } finally {
      // this.spinnerService.stopSpinner()
    }
  }

  // async updateBatchSessions(batchSessions: any): Promise<any> {
  //   try {
  //     return await new Promise<any>((resolve, reject) => {
  //       this.spinnerService.loadingMessage = "Updating sessions";
  //       // this.spinnerService.startSpinner()
  //       let queryParams: any = {}
  //       if (this.isFaculty) {
  //         queryParams['facultyID'] = this.loginID
  //       }

  //       console.log(queryParams);

  //       this.batchSessionService.updateBatchSession(this.batchID, batchSessions, queryParams).subscribe((response) => {
  //         resolve(response)

  //       }, (err: any) => {
  //         console.error(err);
  //         if (err.statusText.includes('Unknown')) {
  //           reject("No connection to server. Check internet.");
  //           return;
  //         }
  //         reject(err.error.error);
  //       });

  //     })
  //   } finally {
  //     // this.spinnerService.stopSpinner()
  //   }
  // }
  async getCourseModules(courseID: string): Promise<ICourseModule[]> {
    try {
      return await new Promise<ICourseModule[]>((resolve, reject) => {
        this.spinnerService.loadingMessage = "Getting course modules";
        // this.spinnerService.startSpinner()

        this.courseModules = [];
        this.totalCourseModules = 0;
        let queryParams: any = {
          isActive: "1",
          limit: -1,
          offset: 0,
        };
        if (this.isFaculty) {
          queryParams['facultyID'] = this.loginID
        }
        this.courseModuleService.getCourseModules(courseID, queryParams).subscribe((response) => {
          // this.courseModules = response.body
          resolve(response.body)
          this.totalCourseModules = parseInt(response.headers.get("X-Total-Count"));
        }, (err: any) => {
          console.error(err);
          if (err.statusText.includes('Unknown')) {
            reject("No connection to server. Check internet.");
            return;
          }
          reject(err.error.error);
        });
      });
    } finally {
      // this.spinnerService.stopSpinner()

    }
  }

  disableBatchModulesInForm(): void {
    let batchModuleForm: IBatchModule[] = this.batchModuleForm.value.batchModules
    for (let index = 0; index < this.batchModules?.length; index++) {
      for (let j = 0; j < batchModuleForm?.length; j++) {
        if (this.batchModules[index].module.id == batchModuleForm[j].moduleID) {
          this.batchModuleControlArray.at(j).get("id").setValue(this.batchModules[index].id)

          this.batchModuleControlArray.at(j).get("facultyID").setValue(this.batchModules[index].faculty?.id)
          this.batchModuleControlArray.at(j).get("isCompleted").setValue(this.batchModules[index].isCompleted)
          this.batchModuleControlArray.at(j).get("startDate").
            setValue(this.datePipe.transform(this.batchModules[index].startDate, "yyyy-MM-dd"))
          // this.batchModuleControlArray.at(j).get("completedDate").
          //   setValue(this.datePipe.transform(this.batchModules[index].completedDate, "yyyy-MM-dd"))

          this.batchModuleControlArray.at(j).get("estimatedEndDate").
            setValue(this.datePipe.transform(this.batchModules[index].estimatedEndDate, "yyyy-MM-dd"))

          this.batchModuleControlArray.at(j).get("order").setValue(this.batchModules[index].order)
          this.batchModuleControlArray.at(j).get("isMarked").setValue(true)
          this.batchModuleControlArray.at(j).get("isModuleTimingMarked").setValue(false)

          if (this.batchModules[index]?.moduleTimings?.length > 0) {
            this.batchModuleControlArray.at(j).get("isModuleTimingMarked").setValue(true)
          }

          for (let k = 0; k < this.batchModules[index]?.moduleTimings?.length; k++) {
            const element = this.batchModules[index]?.moduleTimings[k];

            this.getModuleTimingControlArray(j).at(element.day.order - 1).get("id").setValue(element.id)
            this.getModuleTimingControlArray(j).at(element.day.order - 1).get("dayID").setValue(element.day.id)
            // this.getModuleTimingControlArray(j).at(element.day.order-1).get("batchID").setValue(element.batchID)
            this.getModuleTimingControlArray(j).at(element.day.order - 1).get("moduleID").setValue(element.moduleID)
            // this.getModuleTimingControlArray(j).at(element.day.order-1).get("facultyID").setValue(element.facultyID)
            this.getModuleTimingControlArray(j).at(element.day.order - 1).get("fromTime").setValue(element.fromTime)
            this.getModuleTimingControlArray(j).at(element.day.order - 1).get("toTime").setValue(element.toTime)
            this.getModuleTimingControlArray(j).at(element.day.order - 1).get("isMarked").setValue(true)
          }
          this.batchModuleControlArray.at(j).enable()
        }
      }
    }
  }

  addCourseModuleToForm(): void {
    this.createBatchModuleForm()
    for (let index = 0; index < this.courseModules.length; index++) {
      this.addBatchModulesToForm()
      let len = this.batchModuleControlArray.controls?.length
      this.batchModuleControlArray.at(len - 1).get("moduleID").setValue(this.courseModules[index].module?.id)
      this.addModuleTiming(index, this.courseModules[index].module?.id)
      this.batchModuleControlArray.at(len - 1).disable()
      // this.batchModuleControlArray.at(len - 1).get("order").setValue(this.courseModules[index].order)
    }
  }

  addModuleTiming(j: number, moduleID: string): void {
    for (let index = 0; index < this.daysList.length; index++) {
      this.addModuleTimingToForm(j)
      let len = this.getModuleTimingControlArray(j).controls?.length
      this.getModuleTimingControlArray(j).at(len - 1).get("dayID").setValue(this.daysList[index].id)
      this.getModuleTimingControlArray(j).at(len - 1).get("moduleID").setValue(moduleID)
      this.getModuleTimingControlArray(j).at(len - 1).get("day").setValue(this.daysList[index])
      this.getModuleTimingControlArray(j).at(len - 1).get("isMarked").setValue(false)
    }
  }

  async getBatchModules(batchID: string): Promise<IBatchModule[]> {
    try {
      return await new Promise<IBatchModule[]>((resolve, reject) => {
        this.spinnerService.loadingMessage = "Getting batch modules";

        this.batchModules = [];
        this.totalBatchModules = 0;
        let queryParams: any = {
          limit: -1,
          offset: 0
        }
        if (this.isFaculty) {
          queryParams['facultyID'] = this.loginID
        }
        this.batchService.getBatchModules(batchID, queryParams).subscribe((response) => {
          resolve(response.body);
          this.totalBatchModules = parseInt(response.headers.get("X-Total-Count"));
        }, (err: any) => {
          console.error(err);
          if (err.statusText.includes('Unknown')) {
            reject("No connection to server. Check internet.");
            return;
          }
          reject(err.error.error);
        });
      });
    } finally {

    }
  }

  addBatchModuleIDToList(): void {
    for (let index = 0; index < this.batchModules?.length; index++) {
      this.moduleIDList.push(this.batchModules[index].module.id)
    }
  }

  async addBatchModule(batchModule: IBatchModule, errors: string[] = []): Promise<void> {
    try {
      return await new Promise<any>((resolve, reject) => {

        this.spinnerService.loadingMessage = "Updating batch modules"

        this.batchService.addBatchModule(this.batchID, batchModule).subscribe((response) => {
          console.log(response);
          resolve(response)
        }, (err: any) => {
          console.error(err)
          reject(err)
          errors.push(err.error?.error)
        })
      })
    } finally {

    }
  }


  // // On clicking module tab.
  // async onModuleTabClick(event: any): Promise<void> {
  //   for (let i = 0; i < this.courseModuleList.length; i++) {
  //     console.log("index", event);

  //     if (i == event) {
  //       this.selectedCourseModule[0] = this.courseModuleList[i].module
  //       try {
  //         this.selectedCourseModule[0].moduleTopics = await this.getModuleTopics(this.courseModuleList[i].module?.id)

  //         if (this.progress2 == 100) {
  //           this.selectedOrderedCourseModules = []
  //           this.processSubTopicsOrder()
  //           this.selectedOrderedCourseModules.push(this.courseModuleList[i]?.module)

  //         }

  //       } catch (error) {
  //         console.error(error)
  //       }
  //     }
  //   }
  //   this.checkAllPresentIDs()
  // }

  // // Format the course modules list.
  // async formatCourseModuleList(): Promise<void> {

  //   // Create course module tab list.
  //   this.courseModuleTabList = []
  //   for (let i = 0; i < this.courseModuleList.length; i++) {
  //     this.courseModuleTabList.push(
  //       {
  //         moduleName: this.courseModuleList[i].module.moduleName,
  //         module: this.courseModuleList[i].module,
  //       }
  //     )

  //     if (this.courseModuleList[i].module.logo) {
  //       this.courseModuleTabList[i].imageURL = this.courseModuleList[i].module.logo
  //     }
  //     else {
  //       this.courseModuleTabList[i].imageURL = "assets/icon/grey-icons/Score.png"
  //     }
  //   }

  //   if (this.courseModuleList.length > 0) {
  //     this.selectedCourseModule[0] = this.courseModuleList[0].module
  //     try {
  //       this.selectedCourseModule[0].moduleTopics = await this.getModuleTopics(this.courseModuleList[0].module?.id)
  //       console.log("this.selectedCourseModule", this.selectedCourseModule[0].moduleTopics);
  //     } catch (error) {
  //       console.error(error)
  //     }
  //   }
  //   this.setDefaultCheckedForSubTopics()
  //   // this.formatCourseModuleHeaderList()
  // }


  // setDefaultCheckedForSubTopics(): void {
  //   // this.selectedSubTopicsID = []
  //   console.log(this.courseModuleList);
  //   // alert("test")
  //   for (let index = 0; index < this.courseModuleList.length; index++) {
  //     let batchModule = this.courseModuleList[index].module;
  //     for (let j = 0; j < batchModule?.moduleTopics?.length; j++) {
  //       let moduleTopics = batchModule.moduleTopics[j];
  //       for (let k = 0; k < moduleTopics?.subTopics.length; k++) {
  //         let subTopic = moduleTopics?.subTopics[k];
  //         subTopic.isChecked = false

  //         if (subTopic?.batchSessionTopicIsCompleted != null) {
  //           this.toggleSubTopics(subTopic?.id, subTopic, moduleTopics)
  //           continue
  //         }
  //       }
  //       // this.isAnySubTopicMarked(moduleTopics)
  //     }
  //   }
  //   console.log("this.selectedSubtopicSId",this.selectedSubTopicsID);
  //   this.checkAllPresentIDs()

  // }

  // checkAllPresentIDs(): void {
  //   for (let index = 0; index < this.selectedCourseModule[0]?.moduleTopics.length; index++) {
  //     const topic = this.selectedCourseModule[0]?.moduleTopics[index];
  //     for (let index = 0; index < topic.subTopics.length; index++) {
  //       const subTopic = topic.subTopics[index];
  //       if (this.selectedSubTopicsID.includes(subTopic.id)) {
  //         subTopic.isChecked = true
  //       }
  //       this.isAnySubTopicMarked(topic)
  //     }
  //   }
  // }
  // ============================================================

  async changeTabs(key: number): Promise<void> {
    switch (key) {
      case 0:
        this.progress1 = 0;
        this.progress2 = 0;
        break;
      case 1:
        console.log("in course change progress");
        if (this.progress1 == 0) {
          this.progress1 = 100;
        }
        if (this.progress1 == 100) {
          this.progress2 = 0;
        }
        // try {
        //   this.courseModuleList = await this.getBatchModules(this.batchID)
        //   console.log("courseModuleList ", this.courseModuleList);
        //   await this.getAllBatchModules()
        //   this.formatCourseModuleList()
        // } catch (error) {
        //   console.log(error);

        // }

        break;
      case 2:
        this.progress1 = 100;
        this.progress2 = 100;
        // this.isReordering = true
        // this.tabIndex = 0
        // this.onModuleTabClick(0)
        break;
      case 3:
       this.activeModuleTab=2
        break;

      default:
        this.resetToCourseModule()
        break;
    }
  }
  // processSubTopicsOrder(): void {
  //   console.log("this.selectedSubTopicsID ", this.selectedSubTopicsID);
  //   console.log("this.courseModuleList ", this.courseModuleList);
  //   console.log("this.selectedOrderedCourseModules", this.selectedOrderedCourseModules);


  //   for (let index = 0; index < this.courseModuleList.length; index++) {
  //     let batchModule = this.courseModuleList[index].module;
  //     for (let j = 0; j < batchModule?.moduleTopics?.length; j++) {
  //       let moduleTopics = batchModule.moduleTopics[j];
  //       for (let k = 0; k < moduleTopics?.subTopics.length; k++) {
  //         let subTopic = moduleTopics?.subTopics[k];

  //         if (!this.isSubTopicMarked(subTopic?.id)) {
  //           moduleTopics?.subTopics?.splice(k, 1)
  //           k--
  //         }
  //         console.log(subTopic.topicName);

  //         if (moduleTopics.subTopics.length < 1) {
  //           batchModule.moduleTopics.splice(j, 1)
  //           j--
  //         }
  //       }
  //     }
  //   }
  // }

  // toggleSubTopics(subTopicID: string, subTopic: any, topic: any) {
  //   if (this.selectedSubTopicsID.includes(subTopicID)) {
  //     subTopic.isChecked = false
  //     this.selectedSubTopicsID.splice(this.selectedSubTopicsID.indexOf(subTopicID), 1)
  //     this.isAnySubTopicMarked(topic)
  //     return
  //   }
  //   subTopic.isChecked = true

  //   this.selectedSubTopicsID.push(subTopicID)
  //   this.isAnySubTopicMarked(topic)
  // }

  // isAnySubTopicMarked(topic: any): boolean {
  //   for (let index = 0; index < topic.subTopics.length; index++) {
  //     const element = topic.subTopics[index];
  //     if (this.selectedSubTopicsID.includes(element.id)) {
  //       topic.isChecked = true
  //       return true
  //     }
  //   }
  //   topic.isChecked = false
  //   return false
  // }

  // isAllSubTopicMarked(topic: any): boolean {
  //   for (let index = 0; index < topic.subTopics.length; index++) {
  //     const element = topic.subTopics[index];
  //     if (!this.selectedSubTopicsID.includes(element.id)) {
  //       // topic.isChecked = false
  //       return false
  //     }
  //   }
  //   // topic.isChecked = false
  //   return true
  // }

  // isSubTopicMarked(subTopicID: string): boolean {
  //   if (this.selectedSubTopicsID.includes(subTopicID)) {
  //     return true
  //   }
  //   return false
  // }

  // toggleAllSubTopics(topic: any) {
  //   console.log("this.isAllSubTopicMarked(topic)", this.isAnySubTopicMarked(topic));

  //   if (!topic.isChecked || !this.isAllSubTopicMarked(topic)) {
  //     this.markAllSubTopics(topic)
  //     return
  //   }
  //   this.unmarkAllSubTopics(topic)
  // }

  // unmarkAllSubTopics(topic: any) {
  //   topic.isChecked = false
  //   for (let index = 0; index < topic.subTopics.length; index++) {
  //     const element = topic.subTopics[index];
  //     if (this.selectedSubTopicsID.includes(element.id)) {
  //       element.isChecked = false
  //       this.selectedSubTopicsID.splice(this.selectedSubTopicsID.indexOf(element.id), 1)
  //     }
  //   }

  // }

  // markAllSubTopics(topic: any): void {
  //   topic.isChecked = true
  //   for (let index = 0; index < topic.subTopics.length; index++) {
  //     const element = topic.subTopics[index];
  //     if (this.selectedSubTopicsID.includes(element.id)) {
  //       element.isChecked = true
  //       continue
  //     }
  //     this.toggleSubTopics(element.id, element, topic)
  //   }

  // }

  resetToCourseModule(): void {
    // this.isCourseModule=false
    this.progress1 = 0
    this.progress2 = 0
    this.toggleCourseModule()
  }
  redirectToCourseDetails(): void {
    this.progress1 = 0
    this.progress2 = 0
    this.toggleCourseModule()
    // this._location.back();
  }

  toggleCourseModule(): void {
    this.isCourseModule = !this.isCourseModule
  }
  async updateBatchModule(batchModule: IBatchModule, errors: string[] = []): Promise<void> {
    try {
      this.spinnerService.loadingMessage = "Updating batch modules"
      return await new Promise<any>((resolve, reject) => {

        this.batchService.updateBatchModule(this.batchID, batchModule).subscribe((response) => {
          console.log(response);
          resolve(response.body)
        }, (err: any) => {
          console.error(err)
          reject(err)
          errors.push(err.error?.error)
        })
      })

    } finally {
      // this.spinnerService.stopSpinner()
    }


  }

  validateBatchModule(): void {
    console.log(this.batchModuleForm.controls);

    if (!this.batchModuleForm.valid) {
      this.utilService.markGroupDirty(this.batchModuleForm)
      return
    }
    this.submitBatchModules()
  }


  // Get course module list.
  async getCourseModuleList(): Promise<any> {
    this.spinnerService.loadingMessage = "Getting Course Modules..."

    let queryParams: any = {
      limit: -1,
      offset: 0
    }

    if (this.isFaculty) {
      queryParams['facultyID'] = this.loginID
    }
    try {
      return await new Promise<any>((resolve, reject) => {

        this.courseModuleService.getCourseModules(this.courseID, queryParams).subscribe((response) => {
          // this.courseModuleList = response.body
          // this.formatCourseModuleList()
          resolve(response.body)
        }, error => {
          console.error(error)
          if (error.statusText.includes('Unknown')) {
            alert("No connection to server. Check internet.")
            reject(error)
          }
        })

      })
    } finally {

    }
  }

  async submitBatchModules(): Promise<void> {
    let batchModules: IBatchModule[] = this.batchModuleForm.value.batchModules
    let err: string[] = []
    console.log(batchModules);

    this.setModuleTimings(batchModules)

    for (let index = 0; index < batchModules?.length; index++) {
      if (batchModules[index]?.id) {
        try {
          await this.updateBatchModule(batchModules[index], err)
        } catch (error) {
          console.error(error);
        }
        continue
      }
      try {
        await this.addBatchModule(batchModules[index], err)
      } catch (error) {
        console.error(error);
      }
    }

    if (err.length == 0) {
      alert("Modules successufully updated")
      if (this.isFaculty) {
        this.changeTabs(1)
      }
      if(this.isAdmin){
        this.changeTabs(3)
      }
      // this.onTabChange(2)
    } else {
      let errorString = ""
      for (let index = 0; index < err.length; index++) {
        errorString += (index == 0 ? "" : "\n") + err[index]
      }
      alert(errorString)
      return
    }

    // this.searchBatch()
    // this.changeProgress1()
  }

  setModuleTimings(batchModules: IBatchModule[]): void {
    for (let index = 0; index < batchModules.length; index++) {
      for (let j = 0; j < batchModules[index].moduleTimings?.length; j++) {
        if (!batchModules[index].moduleTimings[j].isMarked) {
          batchModules[index].moduleTimings.splice(j, 1)
          j = j - 1
          continue
        }
        batchModules[index].moduleTimings[j].facultyID = batchModules[index].facultyID
      }
    }
  }

  onDeleteModuleClick(batchModuleID: string): void {
    this.waitingListMessage = null
    this.entity = null

    this.openModal(this.deleteModal, "md").result.then(() => {
      this.deleteBatchModule(batchModuleID)
    }, (err) => {
      console.error(err);
      return
    })
  }

  deleteBatchModule(batchModuleID: string): void {
    this.spinnerService.loadingMessage = "Deleting module"

    this.batchService.deleteBatchModule(this.batchID, batchModuleID).subscribe((response: any) => {
      console.log(response);
      alert("Batch module successfully deleted")
      this.showModules()
      // this.batchModules = await this.getBatchModules(this.selectedBatch.id)
      // this.changeTabs(1)
      // this.onTabChange(2)
    }, (err: any) => {
      console.error(err);
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.");
        return;
      }
      alert(err.error?.error)
    })
  }

  getDaysList(): void {
    this.generalService.getDaysList().subscribe((response: any) => {
      this.daysList = response.body
      // this.formatDayListFields()
    }, (err) => {
      console.error(err);
    })
  }

  // Format the fields of day list.
  // formatDayListFields(): void {

  // 	for (let i = 0; i < this.daysList.length; i++) {
  // 		this.daysList[i].isSelected = false
  // 		for (let j = 0; j < this.batchTiming.controls.length; j++) {
  // 			if (this.daysList[i].id == this.batchTiming.at(j).get('day').value.id) {
  // 				this.daysList[i].isSelected = true
  // 				continue
  // 			}
  // 		}
  // 	}
  // }

  openModal(modalContent: any, modalSize = "lg"): NgbModalRef {

    // this.resetBrochureUploadFields()
    // this.resetLogoUploadFields()

    let options: NgbModalOptions = {
      ariaLabelledBy: 'modal-basic-title', keyboard: false,
      backdrop: 'static', size: modalSize
    }
    this.modalRef = this.modalService.open(modalContent, options)
    return this.modalRef

    /*this.modalRef.result.subscribe((result) => {
    }, (reason) => {

    });*/
  }

  // Get All Faculty List
  getFacultyList(): void {
    this.generalService.getFacultyList().subscribe((data: any) => {
      this.facultyList = data.body
      console.log(this.facultyList);
      if (this.isFaculty) {
        this.processFaculties()
      }
    }, (err) => {
      console.error(err);
    })

  }
  processFaculties() {
    for (let index = 0; index < this.facultyList.length; index++) {
      const element = this.facultyList[index];
      if (element.id != this.loginID) {
        this.facultyList.splice(index, 1)
        index--
      }

    }
  }
}
