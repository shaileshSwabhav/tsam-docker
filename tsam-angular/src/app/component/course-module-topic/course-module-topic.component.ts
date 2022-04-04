import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { Role, UrlConstant } from 'src/app/service/constant';
import { CourseModuleService, ICourseModule } from 'src/app/service/course-module/course-module.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { StorageService } from 'src/app/service/storage/storage.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { IModuleTopic } from 'src/app/models/course/module_topic';
import { LocalService } from 'src/app/service/storage/local.service';

@Component({
  selector: 'app-course-module-topic',
  templateUrl: './course-module-topic.component.html',
  styleUrls: ['./course-module-topic.component.css']
})
export class CourseModuleTopicComponent implements OnInit {

  // Course module.
  selectedModuleList: ICourseModule[]
  courseModuleList: ICourseModule[]
  courseID: string

  // Module topics.
  selectedModuleTopicList: IModuleTopic[]

  // Access.
  permission: IPermission

  // extra

  isAssignmentVisible: boolean
	isFaculty: boolean

  // Input


  @Output() changeProgress2: EventEmitter<any>
  @Output() createCourseModuleClick: EventEmitter<any>
  @Output() changeToTopics: EventEmitter<any>
  @Output() changeToModules: EventEmitter<any>

  // Update course module.
  toBeDeletedModuleCount: number
  toBeUpdatedModuleCount: number
  toBeAddedModuleCount: number
  tobeDeletedModuleList: string[]
  tobeUpdatedModuleList: any[]
  toBeAddedModuleList: any[]
  isCourseModuleUpdateSuccessfulCount: number

  constructor(
    private utilService: UtilityService,
    private urlConstant: UrlConstant,
    private courseModuleService: CourseModuleService,
    private route: ActivatedRoute,
    private spinnerService: SpinnerService,
    private router: Router,
    private storageService: StorageService,
    private localService: LocalService,
		private role: Role,
  ) {
    // this.courseID = "ab62731f-f9b1-4bf8-9f0d-cbf691077f9a"
    this.initializeVariables()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
    this.extractQueryParams()
    if (this.storageService.getItem("selectedModuleList") && this.storageService.getItem("courseModules") && this.courseID) {
      this.selectedModuleList = JSON.parse(this.storageService.getItem("selectedModuleList"))
      this.courseModuleList = JSON.parse(this.storageService.getItem("courseModules"))
      this.formatCourseModuleList()
    }
  }

  // Initialize global variables.
  initializeVariables(): void {

    // Access.
    this.isFaculty = (this.localService.getJsonValue("roleName") == this.role.FACULTY ? true : false)
		if (this.isFaculty){
			this.permission = this.utilService.getPermission(this.urlConstant.BANK_COURSE)
		}
		if (!this.isFaculty){
			this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_COURSE_MASTER)
		}

    // Course module.
    this.selectedModuleList = []

    // Module topics.
    this.selectedModuleTopicList = []
    this.courseModuleList = []

    // extra

    this.isAssignmentVisible = false
    this.changeProgress2 = new EventEmitter()
    this.createCourseModuleClick = new EventEmitter()
    this.changeToModules = new EventEmitter()
    this.changeToTopics = new EventEmitter()

    // Update course module.
    this.toBeAddedModuleCount = 0
    this.toBeUpdatedModuleCount = 0
    this.toBeDeletedModuleCount = 0
    this.tobeDeletedModuleList = []
    this.tobeUpdatedModuleList = []
    this.toBeAddedModuleList = []
    this.isCourseModuleUpdateSuccessfulCount = 0
  }

  //********************************************* COURSE MODULE FUNCTIONS ************************************************************

  // On clicking module on side tab.
  onModuleClick(courseModule: ICourseModule): void {
    this.selectedModuleList.filter((value) => value.isSelected = false)
    courseModule.isSelected = true
    this.selectedModuleTopicList = courseModule.module.moduleTopics
    this.setIsClickedModuleTopic()
  }

  onNextClick(): void {
    this.isAssignmentVisible = true
    this.changeProgress2.emit()
  }

  onCreateClick(): void {
    // this.createCourseModuleClick.emit()
    this.processSelectedModuleList()
  }

  // Format course module list fields.
  formatCourseModuleList(): void {

    // Give default logo to course modules.
    for (let i = 0; i < this.selectedModuleList.length; i++) {
      if (this.selectedModuleList[i].module?.logo == null) {
        this.selectedModuleList[i].module.logo = "assets/icon/grey-icons/Score.png"
      }
    }

    // Get module topics of the first course module.
    if (this.selectedModuleList?.length > 0) {
      this.onModuleClick(this.selectedModuleList[0])
    }
  }

  // Redirect to course details page.
  redirectToCourseDetails(): void {
    this.router.navigate(['/course/master/details'], {
      queryParams: {
        "courseID": this.courseID,
      }
    }).catch(err => {
      console.error(err)
    })
  }

  resetToModules(): void {

    this.changeToModules.emit()
  }
  resetToTopics(): void {
    this.isAssignmentVisible = false;
    this.changeToTopics.emit()
  }
  //********************************************* MODULE TOPIC FUNCTIONS ************************************************************

  // Set isTopicClicked and isSubtopicClicked as false initially.
  setIsClickedModuleTopic(): void {
    for (let i = 0; i < this.selectedModuleTopicList.length; i++) {
      this.selectedModuleTopicList[i]['isTopicClicked'] = false
      if (this.selectedModuleTopicList[i].subTopics?.length > 0) {
        for (let j = 0; j < this.selectedModuleTopicList[i].subTopics.length; j++) {
          this.selectedModuleTopicList[i].subTopics['isTopicClicked'] = false
        }
      }
    }
  }

  //********************************************* COURSE MODULE UPDATE FUNCTIONS ************************************************************

  // Add, update or delete selected modules.
  processSelectedModuleList(): void {
    for (let i = 0; i < this.courseModuleList.length; i++) {
      let tempCourseModule: any = this.isModulePresentInSelectedModulesList((this.courseModuleList[i].module.id))

      // For counting to be updated modules count.  
      if (tempCourseModule != null) {
        this.tobeUpdatedModuleList.push(tempCourseModule)
        this.toBeUpdatedModuleCount = this.toBeUpdatedModuleCount + 1
      }

      // For counting to be deleted modules count.  
      if (tempCourseModule == null) {
        this.tobeDeletedModuleList.push(this.courseModuleList[i].id)
        this.toBeDeletedModuleCount = this.toBeDeletedModuleCount + 1
      }
    }

    // For counting to be added modules count.    
    for (let i = 0; i < this.selectedModuleList.length; i++) {
      if (!this.isModulePresentInCourseModules(this.selectedModuleList[i].module.id)) {
        this.toBeAddedModuleList.push(this.selectedModuleList[i])
        this.toBeAddedModuleCount = this.toBeAddedModuleCount + 1
      }
    }

    // Set the count for successful operations to sum of all operations.
    this.isCourseModuleUpdateSuccessfulCount = this.toBeDeletedModuleCount + this.toBeUpdatedModuleCount + this.toBeAddedModuleCount

    // console.log("tobe deleted: ", this.toBeDeletedModuleCount)
    // console.log("tobe updates: ", this.toBeUpdatedModuleCount)
    // console.log("tobe added: ", this.toBeAddedModuleCount)

    // console.log("tobe deleted: ", this.tobeDeletedModuleList)
    // console.log("tobe updates: ", this.tobeUpdatedModuleList)
    // console.log("tobe added: ", this.toBeAddedModuleList)

    this.deleteAllCourseModules()
    this.updateAllCourseModules()
    this.addAllCourseModules()
  }

  // Delete all course modules.
  deleteAllCourseModules(): void {

    // If to be deleted modules count is greater than one then call all delete course modules API calls.
    if (this.toBeDeletedModuleCount > 0) {
      for (let index = 0; index < this.tobeDeletedModuleList.length; index++) {
        this.deleteCourseModule(this.tobeDeletedModuleList[index])
      }
    }
  }

  // Update all course modules.
  updateAllCourseModules(): void {

    // If to be deleted modules count is 0 and to be updated modules count is greater than one then call all
    // update course modules API calls.
    if (this.toBeDeletedModuleCount == 0 && this.toBeUpdatedModuleCount > 0) {
      for (let index = 0; index < this.tobeUpdatedModuleList.length; index++) {
        this.updateCourseModule(this.tobeUpdatedModuleList[index])
      }
    }
  }

  // Add all course modules.
  addAllCourseModules(): void {

    // If to be deleted modules count is 0, to be updated count is 0 and to be added modules count is greater than one then call all
    // add course modules API calls.
    if (this.toBeDeletedModuleCount == 0 && this.toBeUpdatedModuleCount == 0 && this.toBeAddedModuleCount > 0) {
      for (let index = 0; index < this.toBeAddedModuleList.length; index++) {
        this.addCourseModule(this.toBeAddedModuleList[index])
      }
    }
  }

  // Check if module is present in selected module list.
  isModulePresentInSelectedModulesList(moduleID: string): any {
    for (let index = 0; index < this.selectedModuleList.length; index++) {
      if (this.selectedModuleList[index].module.id == moduleID) {
        return this.selectedModuleList[index]
      }
    }
    return null
  }

  // Check if module is present in course module list.
  isModulePresentInCourseModules(moduleID: string): boolean {
    for (let index = 0; index < this.courseModuleList.length; index++) {
      if (this.courseModuleList[index].module.id == moduleID) {
        return true
      }
    }
    return false
  }

  // Add course module.
  addCourseModule(courseModule: any): void {
    let errors: string[] = []
    this.spinnerService.loadingMessage = "Updating Course Modules"

    this.courseModuleService.addCourseModule(this.courseID, courseModule).subscribe((response: any) => {
      this.isCourseModuleUpdateSuccessfulCount = this.isCourseModuleUpdateSuccessfulCount - 1
      this.alertOnSuccessForCourseModuleUpdate()
    }, (err: any) => {
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    }).add(() => {
      this.toBeAddedModuleCount = this.toBeAddedModuleCount - 1
      this.stopSpinnerForCourseModuleUpdate()
    })
  }

  // Update course module.
  updateCourseModule(courseModule: ICourseModule): void {
    this.spinnerService.loadingMessage = "Updating Course Modules"

    courseModule.moduleID = courseModule.module.id
    this.courseModuleService.updateCourseModule(this.courseID, courseModule).subscribe((response: any) => {
      this.isCourseModuleUpdateSuccessfulCount = this.isCourseModuleUpdateSuccessfulCount - 1
      this.alertOnSuccessForCourseModuleUpdate()
    }, (err: any) => {
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    }).add(() => {
      this.toBeUpdatedModuleCount = this.toBeUpdatedModuleCount - 1
      if (this.toBeUpdatedModuleCount == 0 && this.toBeAddedModuleCount > 0) {
        this.addAllCourseModules()
        return
      }
      this.stopSpinnerForCourseModuleUpdate()
    })
  }

  // Delete course module.
  deleteCourseModule(courseModuleID: string): void {
    this.spinnerService.loadingMessage = "Updating Course Modules"

    this.courseModuleService.deleteCourseModule(this.courseID, courseModuleID).subscribe((response: any) => {
      this.isCourseModuleUpdateSuccessfulCount = this.isCourseModuleUpdateSuccessfulCount - 1
      this.alertOnSuccessForCourseModuleUpdate()
    }, (err: any) => {
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    }).add(() => {
      this.toBeDeletedModuleCount = this.toBeDeletedModuleCount - 1
      if (this.toBeDeletedModuleCount == 0 && this.toBeUpdatedModuleCount > 0) {
        this.updateAllCourseModules()
        return
      }
      if (this.toBeDeletedModuleCount == 0 && this.toBeUpdatedModuleCount == 0) {
        this.addAllCourseModules()
        return
      }
      this.stopSpinnerForCourseModuleUpdate()
    })
  }

  // Stop spinner for course module update.
  stopSpinnerForCourseModuleUpdate(): void {

    if (this.toBeDeletedModuleCount == 0 && this.toBeUpdatedModuleCount == 0
      && this.toBeAddedModuleCount == 0) {

    }
  }

  // Alert messgae on successful course module update.
  alertOnSuccessForCourseModuleUpdate(): void {

    if (this.isCourseModuleUpdateSuccessfulCount == 0) {
      alert("Course modules updated successfully")
      this.createCourseModuleClick.emit()

    }
  }

  //********************************************* OTHER FUNCTIONS ************************************************************

  // Extract query params from url.
  extractQueryParams(): void {
    this.courseID = this.route.snapshot.queryParamMap.get("courseID")
    // this.courseID = "ab62731f-f9b1-4bf8-9f0d-cbf691077f9a"
  }

  //********************************************* GET FUNCTIONS ************************************************************

  // // Get course module list.
  // getCourseModuleList(): void {
  //   this.spinnerService.loadingMessage = "Getting course modules"
  //   
  //   let queryParams: any = {
  //     limit: -1,
  //     offset: 0
  //   }
  //   this.courseModuleService.getCourseModules(this.courseID, queryParams).subscribe((response: any) => {
  //     this.courseModuleList = response.body
  //     this.formatCourseModuleList()
  //   },(err : any) => {
  //     console.error(err);
  //  
  //     if (err.statusText.includes('Unknown')) {
  //       alert("No connection to server. Check internet.")
  //       return
  //     }
  //     alert(err.error.error)
  //   }).add(() => {
  //     
  //   })
  // }

}
