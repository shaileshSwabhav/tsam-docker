import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalRef, NgbModalOptions } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BatchSessionService } from 'src/app/service/batch-session/batch-session.service';
import { BatchTopicAssignmentService } from 'src/app/service/batch-topic-assignment/batch-topic-assignment.service';
import { IBatchTopicAssignment } from 'src/app/service/batch-topic/batch-topic.service';
import { Role, UrlConstant } from 'src/app/service/constant';
import { CourseModuleService, ICourseModule } from 'src/app/service/course-module/course-module.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { IModule, ITopicProgrammingQuestion, ModuleService } from 'src/app/service/module/module.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { IModuleTopic } from 'src/app/models/course/module_topic';
import { LocalService } from 'src/app/service/storage/local.service';
import { BatchService, IBatchModule } from 'src/app/service/batch/batch.service';
import { HttpResponse } from '@angular/common/http';
import { MatTabChangeEvent } from '@angular/material/tabs';

@Component({
  selector: 'app-batch-topic-assignment',
  templateUrl: './batch-topic-assignment.component.html',
  styleUrls: ['./batch-topic-assignment.component.css']
})
export class BatchTopicAssignmentComponent implements OnInit {


  batchTopicAssignmentAddCount: number
  batchTopicAssignments: IBatchTopicAssignment[]
  selectedModuleID: string

  // batch-module
  newProgrammingQuestions: string[]
  selectedProgrammingQuestion: string[]
  totalModuleTopics: number

  batchID: string
  courseID: string

  // modal
  modalRef: any

  // inputs
  @Input() isView: boolean
  @Input() isAssign: boolean


  // modal
  @ViewChild("deleteModal") deleteModal: any

  // access
  permission: IPermission
  loginID: string
  isFaculty: boolean

  constructor(
    private moduleService: ModuleService,
    private batchTopicAssignmentService: BatchTopicAssignmentService,
    private batchService: BatchService,
    private utilService: UtilityService,
    private urlConstant: UrlConstant,
    private modalService: NgbModal,
    private route: ActivatedRoute,
    private spinnerService: SpinnerService,
    private localService: LocalService,
    private role: Role,
  ) {
    this.extractID()
    this.initializeVariables()
  }

  extractID(): void {
    let queryParams = this.route.snapshot.queryParamMap
    this.batchID = queryParams.get("batchID")
    this.courseID = queryParams.get("courseID")
  }

  initializeVariables(): void {
    if (this.localService.getJsonValue("roleName") == this.role.ADMIN ||
      this.localService.getJsonValue("roleName") == this.role.SALES_PERSON) {
      this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_BATCH_MASTER_SESSION_DETAILS)
    }

    this.isFaculty = this.localService.getJsonValue("roleName") == this.role.FACULTY

    if (this.isFaculty) {
      this.permission = this.utilService.getPermission(this.urlConstant.MY_BATCH_SESSION_DETAILS)
    }

    this.batchTopicAssignments = []
    this.batchModuleTopics = []
    this.selectedProgrammingQuestion = []
    this.newProgrammingQuestions = []

    this.batchTopicAssignmentAddCount = 0

    this.loginID = this.localService.getJsonValue("loginID")
  }

  getAllComponents(): void {
    // this.getCourseModule()
    this.getBatchModules()
  }

  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
    this.getAllComponents()
  }

  onModuleTabClick(matTab: MatTabChangeEvent): void {
    // console.log(matTab);
    if (matTab.index >= 0) {
      this.selectedModuleID = this.batchModules[matTab.index].module.id
      this.batchModuleTopics = this.batchModules[matTab.index].module.moduleTopics
      this.getModuleTopics(this.batchModules[matTab.index].module.id)
    }
  }

  onDeleteAssignmentClick(subTopicID: string, topicAssignmentID: string): void {
    // console.log(subTopicID, topicAssignmentID);
    this.openModal(this.deleteModal, 'md').result.then(() => {
      this.deleteBatchTopicAssignment(subTopicID, topicAssignmentID)
    }, (err) => {
      console.error(err);
      return
    })
  }

  deleteBatchTopicAssignment(subTopicID: string, topicAssignmentID: string): void {
    this.spinnerService.loadingMessage = "Deleting assignment"

    this.batchTopicAssignmentService.deleteBatchTopicAssignment(this.batchID, subTopicID, topicAssignmentID).subscribe((response: any) => {
      // console.log(response);
      this.getBatchModules()
    }, (err) => {
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(this.utilService.getErrorString(err))
      console.error(err)
    })
  }

  getSelectedProgrammingQuestion(): void {

    for (let index = 0; index < this.batchModules.length; index++) {
      for (let j = 0; j < this.batchModules[index].module.moduleTopics.length; j++) {
        let moduleTopics: IModuleTopic[] = this.batchModules[index].module.moduleTopics
        for (let k = 0; k < moduleTopics[j].topicProgrammingQuestions.length; k++) {
          if (!this.selectedProgrammingQuestion.includes(moduleTopics[j].topicProgrammingQuestions[k].programmingQuestion.id)) {
            this.selectedProgrammingQuestion.push(moduleTopics[j].topicProgrammingQuestions[k].programmingQuestion.id)
          }
        }
      }
    }

    // console.log("selectedProgrammingQuestion -> ", this.selectedProgrammingQuestion);
    this.markSelectedProgrammingQuestion()
  }

  markSelectedProgrammingQuestion(): void {
    // console.log("newProgrammingQuestions -> ", this.newProgrammingQuestions);

    for (let i = 0; i < this.moduleTopics.length; i++) {
      for (let j = 0; j < this.moduleTopics[i].topicProgrammingQuestions.length; j++) {
        this.moduleTopics[i].topicProgrammingQuestions[j].isMarked = false

        if (this.selectedProgrammingQuestion.includes(this.moduleTopics[i].topicProgrammingQuestions[j].programmingQuestion.id)) {
          this.moduleTopics[i].topicProgrammingQuestions[j].isMarked = true
          continue
        }
        if (this.newProgrammingQuestions.includes(this.moduleTopics[i].topicProgrammingQuestions[j].programmingQuestion.id)) {
          this.moduleTopics[i].topicProgrammingQuestions[j].isMarked = true
        }
      }
    }
  }


  toggleProgrammingQuestion(topicID: string, topicProgrammingQuestion: ITopicProgrammingQuestion): void {
    let isProgrammingQuestionFound: boolean = false

    // console.log(topicProgrammingQuestion);

    for (let index = 0; index < this.batchTopicAssignments.length; index++) {
      if (this.batchTopicAssignments[index].programmingQuestionID == topicProgrammingQuestion.programmingQuestion.id) {
        if (!isProgrammingQuestionFound) {
          isProgrammingQuestionFound = true
        }

        if (topicID == this.batchTopicAssignments[index].topicID) {
          this.batchTopicAssignments.splice(index, 1)
          topicProgrammingQuestion.isMarked = false

          this.newProgrammingQuestions.splice(this.newProgrammingQuestions.indexOf(topicProgrammingQuestion.programmingQuestion.id), 1)
          break
        }
      }
    }

    if (!isProgrammingQuestionFound) {
      let batchTopicAssignment: IBatchTopicAssignment = {
        batchID: this.batchID,
        programmingQuestionID: topicProgrammingQuestion.programmingQuestion.id,
        topicID: topicID,
        moduleID: this.selectedModuleID
      }

      this.batchTopicAssignments.push(batchTopicAssignment)
      this.newProgrammingQuestions.push(topicProgrammingQuestion.programmingQuestion.id)
      topicProgrammingQuestion.isMarked = true
      // console.log("newProgrammingQuestions -> ", this.newProgrammingQuestions);
    }
    // console.log(this.batchTopicAssignments);
  }

  moduleTopics: IModuleTopic[] = []
  batchModules: IBatchModule[] = []
  batchModuleTopics: IModuleTopic[] = []
  totalBatchModules: number

  getBatchModules(): void {
    this.spinnerService.loadingMessage = "Getting modules"

    let queryParams: any = {
      limit: -1,
      offset: 0,
      field: ["ModuleTopics", "TopicProgrammingQuestions"]
    }

    if (this.isFaculty) {
      queryParams.facultyID = this.loginID
    }

    this.batchModules = []
    this.batchModuleTopics = []
    this.totalBatchModules = 0

    console.log(queryParams);

    this.batchService.getBatchModules(this.batchID, queryParams).subscribe((response: HttpResponse<IBatchModule[]>) => {
      // this.batchService.getBatchModulesWithAllFields(this.batchID, queryParams).subscribe((response: HttpResponse<IBatchModule[]>) => {
      this.batchModules = response.body
      // console.log(this.batchModules)
      this.totalBatchModules = parseInt(response.headers.get("X-Total-Count"))
      this.batchModuleTopics = this.batchModules[0].module.moduleTopics
      // this.getModuleTopics(this.batchModules[0].module.id)
    }, (err: any) => {
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error?.error)
    })
  }

  getModuleTopics(moduleID: string): void {
    this.spinnerService.loadingMessage = "Getting topics"

    this.moduleTopics = []

    let queryParams: any = {
      limit: -1,
      offset: 0,
      batchID: this.batchID
    }

    this.moduleService.getModuleTopic(moduleID, queryParams).subscribe((response: HttpResponse<IModuleTopic[]>) => {
      this.moduleTopics = response.body
      // console.log(this.moduleTopics);
      this.getSelectedProgrammingQuestion()
    }, (err: any) => {
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error?.error)
    })
  }

  // Create session plan.
  async onAssignAssignmentClick(): Promise<void> {
    let errors: string[] = []
    for (let index = 0; index < this.batchTopicAssignments.length; index++) {
      // console.log(this.batchTopicAssignments[index]);
      try {
        await this.addBatchTopicAssignment(this.batchTopicAssignments[index])
      } catch (err) {
        errors.push(err)
      }
    }

    this.batchTopicAssignments = []

    if (errors.length > 0) {
      alert(this.createErrorString(errors))
      return
    }
    alert("Assignments successfully assigned")
    this.getBatchModules()
  }

  // Add batch topic assignment.
  addBatchTopicAssignment(batchTopicAssignment: IBatchTopicAssignment): Promise<any> {
    try {
      return new Promise<any>((response, reject) => {
        this.spinnerService.loadingMessage = "Adding assignments"

        this.batchTopicAssignmentService.addBatchTopicAssignment(this.batchID, batchTopicAssignment.topicID, batchTopicAssignment).subscribe((res: any) => {
          // console.log(res);
          response(res)
        }, (err) => {
          reject(err?.error?.error)
        })
      })
    } finally {
    }
  }


  openModal(modalContent: any, modalSize = "lg"): NgbModalRef {
    let options: NgbModalOptions = {
      ariaLabelledBy: 'modal-basic-title', backdrop: 'static',
      size: modalSize, keyboard: false
    }

    this.modalRef = this.modalService.open(modalContent, options)
    return this.modalRef
  }

  createErrorString(errors: string[]): string {
    let errorString = ""
    for (let index = 0; index < errors.length; index++) {
      errorString += (index == 0 ? "" : "\n") + errors[index]
    }
    return errorString
  }
}
