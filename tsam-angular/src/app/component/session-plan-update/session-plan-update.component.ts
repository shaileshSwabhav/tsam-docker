import { CdkDragDrop, moveItemInArray } from '@angular/cdk/drag-drop';
import { DatePipe } from '@angular/common';
import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { IBatchSessionModules } from 'src/app/models/batch/batch_session';
import { BatchSessionService } from 'src/app/service/batch-session/batch-session.service';
import { BatchService, IBatchModule } from 'src/app/service/batch/batch.service';
import { Role } from 'src/app/service/constant';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';
import { LocalService } from 'src/app/service/storage/local.service';

@Component({
  selector: 'app-session-plan-update',
  templateUrl: './session-plan-update.component.html',
  styleUrls: ['./session-plan-update.component.css']
})
export class SessionPlanUpdateComponent implements OnInit {
  // Course Module.
  courseModuleList: any[]
  selectedCourseModule: IBatchSessionModules[]
  allCourseModules: IBatchSessionModules[]
  selectedOrderedCourseModules: IBatchSessionModules[];
  selectedSubTopicsID: string[];
  courseModuleTabList: any[]
  batchModules: IBatchModule[]

  tabIndex: number
  totalBatchModules: number;

  batchID: string;
  loginID: string

  isReordering: boolean;
  isSessionPlanModificationStarted: boolean
  isFaculty: boolean;
  isAdmin: boolean;

  @Output() 
  changeTab:EventEmitter<number> = new EventEmitter();


  constructor(
    private spinnerService: SpinnerService,
    private batchSessionService: BatchSessionService,
    private route: ActivatedRoute,
    private localService: LocalService,
    private role: Role,
    private batchService: BatchService,
    private datePipe: DatePipe,

  ) { 
    this.extractID()
  }

  ngOnInit(): void {
    this.initializeVariables()
  }

  extractID(): void {
    this.route.queryParamMap.subscribe(params => {
      this.batchID = params.get("batchID")
      // this.courseID = params.get("courseID")
      // this.batchName = params.get("batchName")
    }, err => {
      console.error(err);
    })
  }

  initializeVariables(): void {

    this.loginID = this.localService.getJsonValue("loginID")

    this.isFaculty = (this.localService.getJsonValue("roleName") == this.role.FACULTY ? true : false)
    this.isAdmin = (this.localService.getJsonValue("roleName") == this.role.ADMIN ? true : false)

    this.courseModuleList = []
    this.selectedCourseModule = []
    this.allCourseModules = []
    this.selectedOrderedCourseModules = []
    this.selectedSubTopicsID = []
    this.courseModuleTabList = []
    this.batchModules = []

    this.tabIndex = 0
    this.totalBatchModules = 0

    this.isReordering = false
    this.isSessionPlanModificationStarted = false

    // this.changeTab= new EventEmitter
    this.changeTabs(1)
  }

  // On clicking module tab.
  async onModuleTabClick(event: any): Promise<void> {
    for (let i = 0; i < this.courseModuleList.length; i++) {
      console.log("index", event);

      if (i == event) {
        this.selectedCourseModule[0] = this.courseModuleList[i].module
        try {
          this.selectedCourseModule[0].moduleTopics = await this.getModuleTopics(this.courseModuleList[i].module?.id)

          if (this.isReordering) {
            this.selectedOrderedCourseModules = []
            this.processSubTopicsOrder()
            this.selectedOrderedCourseModules.push(this.courseModuleList[i]?.module)

          }

        } catch (error) {
          console.error(error)
        }
      }
    }
    this.checkAllPresentIDs()
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

  processSubTopicsOrder(): void {
    console.log("this.selectedSubTopicsID ", this.selectedSubTopicsID);
    console.log("this.courseModuleList ", this.courseModuleList);
    console.log("this.selectedOrderedCourseModules", this.selectedOrderedCourseModules);


    for (let index = 0; index < this.courseModuleList.length; index++) {
      let batchModule = this.courseModuleList[index].module;
      for (let j = 0; j < batchModule?.moduleTopics?.length; j++) {
        let moduleTopics = batchModule.moduleTopics[j];
        for (let k = 0; k < moduleTopics?.subTopics.length; k++) {
          let subTopic = moduleTopics?.subTopics[k];

          if (!this.isSubTopicMarked(subTopic?.id)) {
            moduleTopics?.subTopics?.splice(k, 1)
            k--
          }
          console.log(subTopic.topicName);

          if (moduleTopics.subTopics.length < 1) {
            batchModule.moduleTopics.splice(j, 1)
            j--
          }
        }
      }
    }
  }

  isSubTopicMarked(subTopicID: string): boolean {
    if (this.selectedSubTopicsID.includes(subTopicID)) {
      return true
    }
    return false
  }

  checkAllPresentIDs(): void {
    for (let index = 0; index < this.selectedCourseModule[0]?.moduleTopics.length; index++) {
      const topic = this.selectedCourseModule[0]?.moduleTopics[index];
      for (let index = 0; index < topic.subTopics.length; index++) {
        const subTopic = topic.subTopics[index];
        if (this.selectedSubTopicsID.includes(subTopic.id)) {
          subTopic.isChecked = true
        }
        this.isAnySubTopicMarked(topic)
      }
    }
  }

  isAnySubTopicMarked(topic: any): boolean {
    for (let index = 0; index < topic.subTopics.length; index++) {
      const element = topic.subTopics[index];
      if (this.selectedSubTopicsID.includes(element.id)) {
        topic.isChecked = true
        return true
      }
    }
    topic.isChecked = false
    return false
  }

  toggleAllSubTopics(topic: any) {
    console.log("this.isAllSubTopicMarked(topic)", this.isAnySubTopicMarked(topic));

    if (!topic.isChecked || !this.isAllSubTopicMarked(topic)) {
      this.markAllSubTopics(topic)
      return
    }
    this.unmarkAllSubTopics(topic)
  }

  isAllSubTopicMarked(topic: any): boolean {
    for (let index = 0; index < topic.subTopics.length; index++) {
      const element = topic.subTopics[index];
      if (!this.selectedSubTopicsID.includes(element.id)) {
        // topic.isChecked = false
        return false
      }
    }
    // topic.isChecked = false
    return true
  }

  unmarkAllSubTopics(topic: any) {
    topic.isChecked = false
    for (let index = 0; index < topic.subTopics.length; index++) {
      const element = topic.subTopics[index];
      if (this.selectedSubTopicsID.includes(element.id)) {
        element.isChecked = false
        this.selectedSubTopicsID.splice(this.selectedSubTopicsID.indexOf(element.id), 1)
      }
    }

  }

  markAllSubTopics(topic: any): void {
    topic.isChecked = true
    for (let index = 0; index < topic.subTopics.length; index++) {
      const element = topic.subTopics[index];
      if (this.selectedSubTopicsID.includes(element.id)) {
        element.isChecked = true
        continue
      }
      this.toggleSubTopics(element.id, element, topic)
    }

  }

  toggleSubTopics(subTopicID: string, subTopic: any, topic: any) {
    if (this.selectedSubTopicsID.includes(subTopicID)) {
      subTopic.isChecked = false
      this.selectedSubTopicsID.splice(this.selectedSubTopicsID.indexOf(subTopicID), 1)
      this.isAnySubTopicMarked(topic)
      return
    }
    subTopic.isChecked = true

    this.selectedSubTopicsID.push(subTopicID)
    this.isAnySubTopicMarked(topic)
  }

  async changeTabs(key: number): Promise<void> {
    switch (key) {
      case 1:
        console.log("in course change progress");
        try {
          this.courseModuleList = await this.getBatchModules(this.batchID)
          console.log("courseModuleList ", this.courseModuleList);
          await this.getAllBatchModules()
          this.formatCourseModuleList()
        } catch (error) {
          console.log(error);

        }

        break;
      case 2:
        this.isReordering = true
        this.tabIndex = 0
        this.onModuleTabClick(0)
        this.changeTab.emit(2)
        break;

      default:
        break;
    }
  }


  // Format the course modules list.
  async formatCourseModuleList(): Promise<void> {

    // Create course module tab list.
    this.courseModuleTabList = []
    for (let i = 0; i < this.courseModuleList.length; i++) {
      this.courseModuleTabList.push(
        {
          moduleName: this.courseModuleList[i].module.moduleName,
          module: this.courseModuleList[i].module,
        }
      )

      if (this.courseModuleList[i].module.logo) {
        this.courseModuleTabList[i].imageURL = this.courseModuleList[i].module.logo
      }
      else {
        this.courseModuleTabList[i].imageURL = "assets/icon/grey-icons/Score.png"
      }
    }

    if (this.courseModuleList.length > 0) {
      this.selectedCourseModule[0] = this.courseModuleList[0].module
      try {
        this.selectedCourseModule[0].moduleTopics = await this.getModuleTopics(this.courseModuleList[0].module?.id)
        console.log("this.selectedCourseModule", this.selectedCourseModule[0].moduleTopics);
      } catch (error) {
        console.error(error)
      }
    }
    this.setDefaultCheckedForSubTopics()
    // this.formatCourseModuleHeaderList()
  }

  setDefaultCheckedForSubTopics(): void {
    // this.selectedSubTopicsID = []
    console.log(this.courseModuleList);
    // alert("test")
    for (let index = 0; index < this.courseModuleList.length; index++) {
      let batchModule = this.courseModuleList[index].module;
      for (let j = 0; j < batchModule?.moduleTopics?.length; j++) {
        let moduleTopics = batchModule.moduleTopics[j];
        for (let k = 0; k < moduleTopics?.subTopics.length; k++) {
          let subTopic = moduleTopics?.subTopics[k];
          subTopic.isChecked = false

          if (subTopic?.batchSessionTopicIsCompleted != null) {
            this.toggleSubTopics(subTopic?.id, subTopic, moduleTopics)
            continue
          }
        }
        // this.isAnySubTopicMarked(moduleTopics)
      }
    }
    console.log("this.selectedSubtopicSId", this.selectedSubTopicsID);
    this.checkAllPresentIDs()

  }


  // On dragging and dropping topics table rows.
  onDropTopicsTableRow(event: CdkDragDrop<string[]>, tempTopics: any[]) {
    moveItemInArray(tempTopics, event.previousIndex, event.currentIndex)
    tempTopics.forEach((topic, index) => {
      topic.order = index + 1
    })
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

  async processBatchSessions(): Promise<void> {

    let errors: string[] = []
    let batchSessions: any[] = []
    for (let index = 0; index < this.courseModuleList.length; index++) {
      let batchModule = this.courseModuleList[index].module;
      let subTopicOrder: number = 0
      for (let j = 0; j < batchModule?.moduleTopics?.length; j++) {
        let moduleTopics = batchModule.moduleTopics[j];
        for (let k = 0; k < moduleTopics?.subTopics.length; k++) {
          let subTopic = moduleTopics?.subTopics[k];

          console.log("subTopic -> ", subTopic);

          let batchSession: any = {
            id: subTopic.batchSessionTopicID,
            batchID: this.batchID,
            moduleID: subTopic.moduleID,
            topicID: subTopic.topicID,
            subTopicID: subTopic.id,
            order: ++subTopicOrder,
            initialDate: this.datePipe.transform(subTopic.batchSessionInitialDate, "yyyy-MM-dd"),
            completedDate: subTopic.batchSessionCompletionDate,
            isCompleted: subTopic.batchSessionTopicIsCompleted,
            batchSessionID: subTopic.batchSessionID,
            totalTime: subTopic.totalTime,
          }
          // if (subTopic?.batchSessionTopicID != null) {
          //   batchSession.id = subTopic.batchSessionTopicID
          // }
          batchSessions.push(batchSession)

        }
      }
    }
    try {
      console.log(batchSessions);

      await this.updateBatchSessions(batchSessions)
    } catch (error) {
      errors.push(error)
    }
    let msgString = ""
    if (errors.length > 0) {
      for (let index = 0; index < errors.length; index++) {
        msgString += (index == 0 ? "" : "\n") + errors[index]
      }
      alert(msgString)
      return
    }
    msgString = 'Session Plan updated successfully'
    alert(msgString)
    this.changeTab.emit(3)

  }

  async getAllBatchModules(): Promise<void> {
    for (let i = 0; i < this.courseModuleList.length; i++) {
      this.allCourseModules[i] = this.courseModuleList[i].module

      try {
        this.allCourseModules[i].moduleTopics = await this.getModuleTopics(this.courseModuleList[i].module?.id)

        // this.processSubTopicsOrder()
      } catch (error) {
        console.error(error)
      }
    }
    console.log("this.allCourseModules", this.allCourseModules);
  }

  async updateBatchSessions(batchSessions: any): Promise<any> {
    try {
      return await new Promise<any>((resolve, reject) => {
        this.spinnerService.loadingMessage = "Updating sessions";
        let queryParams: any = {}
        if (this.isFaculty) {
          queryParams['facultyID'] = this.loginID
        }
        console.log(queryParams);

        this.batchSessionService.updateBatchSession(this.batchID, batchSessions, queryParams).subscribe((response) => {
          resolve(response)
        }, (err: any) => {
          console.error(err);
          if (err.statusText.includes('Unknown')) {
            reject("No connection to server. Check internet.");
            return;
          }
          reject(err.error.error);
        });
      })
    } finally {

    }

  }
}
