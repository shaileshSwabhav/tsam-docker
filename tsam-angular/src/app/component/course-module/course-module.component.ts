import { Component, OnInit, Output, ViewChild, EventEmitter } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { CdkDragDrop, moveItemInArray } from '@angular/cdk/drag-drop';
import { MatTable } from '@angular/material/table';

import { UrlConstant } from 'src/app/service/constant';
import { CourseModuleService, ICourseModule } from 'src/app/service/course-module/course-module.service';
import { ModuleService } from 'src/app/service/module/module.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { StorageService } from 'src/app/service/storage/storage.service';

@Component({
  selector: 'app-course-module',
  templateUrl: './course-module.component.html',
  styleUrls: ['./course-module.component.css']
})
export class CourseModuleComponent implements OnInit {

  // form
  courseModuleForm: FormGroup
  moduleSearchForm: FormGroup

  // modules
  allModules: ICourseModule[]
  totalModules: number
  selectedModulesList: ICourseModule[]

  //course-modules
  courseModules: ICourseModule[]
  totalCourseModules: number

  //flags
  isCourseModule: boolean
  isallCourseModulesFetched: boolean

  // search
  searchFormValue: any
  isSearched: boolean

  // course-details
  courseID: string
  courseName: string

  // other


  paginationString: string

  // constant
  private readonly IGNORE_SEARCH_FIELD: string[] = ["limit", "offset", "courseID"]
  selectedModulesColumn: string[]

  // initialValue = {};

  @Output() changeToProgress1: EventEmitter<any>
  @Output() redirectToCourse: EventEmitter<any>

  @ViewChild('drawer') drawer: any
  @ViewChild('deleteModal') deleteModal: any
  @ViewChild('courseModule') courseModule: any
  @ViewChild('table') table: MatTable<any>;

  constructor(
    private formBuilder: FormBuilder,
    private moduleService: ModuleService,
    private courseModuleService: CourseModuleService,
    private urlConstant: UrlConstant,
    private spinnerService: SpinnerService,
    private router: Router,
    private route: ActivatedRoute,
    private utilService: UtilityService,
    private storageService: StorageService,
  ) {
    this.initializeVariables()
    this.createForms()
  }

  initializeVariables(): void {
    this.changeToProgress1 = new EventEmitter()
    this.redirectToCourse = new EventEmitter()

    this.allModules = []
    this.selectedModulesColumn = []
    this.selectedModulesList = []
    this.courseModules = []


    this.totalModules = 0

    this.searchFormValue = null

    this.isSearched = false

    this.totalCourseModules = 0
    this.isallCourseModulesFetched = false
    this.getAllModules()
  }

  createForms(): void {
    this.createSearchForm()
    this.createModuleForm()
  }

  getAllComponents(): void {
    this.getAllCourseModules()
  }

  extractCourseDetails(): void {
    this.route.queryParamMap.subscribe(params => {
      this.courseID = params.get('courseID')
      this.courseName = params.get('courseName')
    })
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
    this.extractCourseDetails()
    this.getAllComponents()
  }

  changeToCourseModule() {
    this.changeToProgress1.emit()
  }

  createSearchForm(): void {
    this.moduleSearchForm = this.formBuilder.group({
      moduleName: new FormControl(null),
      courseID: new FormControl(this.courseID),
      isActive: new FormControl(null),
      limit: new FormControl(5),
      offset: new FormControl(0)
    })
  }

  createModuleForm(): void {
    this.courseModuleForm = this.formBuilder.group({
      modules: new FormArray([]),
    })
  }

  get moduleControlArray(): FormArray {
    return this.courseModuleForm.get("modules") as FormArray
  }

  addModulesToForm(): void {
    this.moduleControlArray.push(new FormGroup({
      moduleID: new FormControl(null),
      order: new FormControl(null),
      isActive: new FormControl(true),
      courseID: new FormControl(this.courseID),
      module: new FormControl(null),
      isMarked: new FormControl(false)
    }))
  }

  patchValuesToForm(): void {
    for (let index = 0; index < this.allModules.length; index++) {
      this.addModulesToForm()
      this.moduleControlArray.at(index).get("isMarked").setValue(false)
      this.moduleControlArray.at(index).get("module").setValue(this.allModules[index])
      this.moduleControlArray.at(index).get("moduleID").setValue(this.allModules[index].id)
      this.moduleControlArray.at(index).get("order").setValue(0)

      let tempID = this.allModules[index].id
      if (this.isModulePresent(tempID)) {
        this.moduleControlArray.at(index).get("isMarked").setValue(true)
        this.moduleControlArray.at(index).get("order").setValue(this.selectedModulesList[this.getIndexOFModulePresentID(tempID)].order)
      }
    }
  }

  calculateTotalTime(moduleTopics: any): number {
    let totalTime: number = 0
    for (let index = 0; index < moduleTopics?.length; index++) {
      totalTime += moduleTopics[index].totalTime
    }
    return totalTime
  }

  setPaginationString() {
    this.paginationString = ''
    let limit = this.moduleSearchForm.get('limit').value
    let offset = this.moduleSearchForm.get('offset').value

    let start: number = limit * offset + 1
    let end: number = +limit + limit * offset
    if (this.totalModules < end) {
      end = this.totalModules
    }
    if (this.totalModules == 0) {
      this.paginationString = ''
      return
    }
    this.paginationString = `${start} - ${end}`
  }

  searchOrGetModule(): void {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.changePage(1)
      return
    }
    this.moduleSearchForm.patchValue(queryParams)
    this.searchModule()
  }

  changePage(pageNumber: number): void {
    this.moduleSearchForm.get("offset").setValue(pageNumber - 1)
    this.searchModule();
  }

  searchModule(): void {
    this.searchFormValue = { ...this.moduleSearchForm?.value }
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: this.searchFormValue,
    })
    let flag: boolean = true
    for (let field in this.searchFormValue) {
      if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
        delete this.searchFormValue[field];
      } else {
        if (!this.IGNORE_SEARCH_FIELD.includes(field)) {
          this.isSearched = true
        }
        flag = false
      }
    }
    // No API call on empty search.
    if (flag) {
      return
    }
    this.getAllModules()
  }

  getAllCourseModules(): void {
    this.courseModules = []
    this.totalCourseModules = 0
    this.spinnerService.loadingMessage = "Getting course module"

    let queryParams: any = {
      limit: -1,
      offset: 0
    }
    this.courseModuleService.getCourseModules(this.courseID, queryParams).subscribe((response: any) => {
      this.selectedModulesList = response.body
      this.courseModules.push(...response.body)
      this.totalCourseModules = response.headers.get("X-Total-Count")
      // console.log(response);
      this.getModuleFieldName()
    }, (err: any) => {
      this.totalModules = 0
      this.allModules = []
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    }).add(() => {
      this.searchOrGetModule()
      this.isallCourseModulesFetched = true

    })
  }

  getModuleFieldName(): void {
    this.selectedModulesColumn = ["id", "module", "logo", "moduleName", "moduleTopics", "order", "close"]
  }

  getAllModules(): void {
    this.allModules = []
    this.totalModules = 0
    this.isCourseModule = true
    this.spinnerService.loadingMessage = "Getting All modules"

    this.moduleService.getModule(this.searchFormValue).subscribe((response: any) => {
      this.allModules = response.body
      // console.log(this.allModules);

      this.totalModules = response.headers.get("X-Total-Count")
      this.createModuleForm()
      this.patchValuesToForm()
    }, (err: any) => {
      this.totalModules = 0
      this.allModules = []
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    }).add(() => {
      this.isCourseModule = true
      this.setPaginationString()

    })
  }

  processSelectedModuleList(): void {
    // Store the selected module list in session storage.
    this.storageService.setItem("selectedModuleList", JSON.stringify(this.selectedModulesList))
    this.storageService.setItem("courseModules", JSON.stringify(this.courseModules))
    this.changeToCourseModule()
  }

  resetSearchAndGetAll(): void {
    this.resetModuleSearchForm()
    this.isSearched = false
    this.router.navigate([this.urlConstant.COURSE_MODULE])
    this.changePage(1)
  }

  resetModuleSearchForm(): void {
    let limit = this.moduleSearchForm.get("limit").value
    let offset = this.moduleSearchForm.get("offset").value

    this.moduleSearchForm.reset()
    this.moduleSearchForm.get("limit").setValue(limit)
    this.moduleSearchForm.get("offset").setValue(offset)
    this.moduleSearchForm.get("courseID").setValue(this.courseID)
  }



  // Takes a list called selectedEnquiriesList and adds all the checked enquiries to list, also does not contain duplicate values.
  toggleModules(moduleControl: any): void {
    let id = moduleControl.get("module")?.value.id
    if (!this.isModulePresent(id)) {
      this.selectedModulesList.push(moduleControl.value)
      this.setOrderNoForPreview()
      return
    }
    let index = this.getIndexOFModulePresentID(id)
    this.selectedModulesList.splice(index, 1)
    this.setOrderNoForPreview()
  }

  isModulePresent(moduleID: string): boolean {
    for (let index = 0; index < this.selectedModulesList.length; index++) {
      let id = this.selectedModulesList[index].module.id
      if (id == moduleID) {
        return true
      }
    }
    return false
  }

  getIndexOFModulePresentID(moduleID: any): number {
    for (let index = 0; index < this.selectedModulesList.length; index++) {
      let id = this.selectedModulesList[index].module.id
      if (id == moduleID) {
        this.selectedModulesList[index].order = index + 1
        return index
      }
    }
    return -1
  }

  setOrderNoForPreview(): void {
    for (let index = 0; index < this.moduleControlArray.length; index++) {
      if (this.moduleControlArray.at(index).get("isMarked").value) {
        let ind = this.getIndexOFModulePresentID(this.moduleControlArray.at(index).get('module').value.id)
        this.moduleControlArray.at(index).get("order").setValue(ind + 1)
      }
    }
  }

  onDrop(event: CdkDragDrop<string[]>) {
    moveItemInArray(this.selectedModulesList, event.previousIndex, event.currentIndex);
    this.selectedModulesList.forEach((user, idx) => {
      user.order = idx + 1;
    });
  }

  onDropClose(moduleID: string, j: number): void {
    this.selectedModulesList.splice(j, 1)
    for (let index = 0; index < this.moduleControlArray.length; index++) {
      if (this.moduleControlArray.at(index).get("moduleID").value == moduleID) {
        this.moduleControlArray.at(index).get("isMarked").setValue(false)
        break
      }
    }
    this.setOrderNoForPreview()
  }

  redirectToCourseDetails(): void {
    this.redirectToCourse.emit()
  }
}