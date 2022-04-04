import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { UrlConstant } from 'src/app/service/constant';
import { CourseService, ICourseList } from 'src/app/service/course/course.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { CourseDetailsComponent } from '../course-details/course-details.component';
import { CourseModuleComponent } from '../course-module/course-module.component';
import { Location } from '@angular/common';

@Component({
  selector: 'app-course',
  templateUrl: './course.component.html',
  styleUrls: ['./course.component.css']
})
export class CourseComponent implements OnInit {


  //changed 
  @ViewChild(CourseModuleComponent) childC: CourseModuleComponent;
  @ViewChild(CourseDetailsComponent) courseDetailsComponent: CourseDetailsComponent;
  inputValue: string = "module-changed"
  // permission
  permission: IPermission

  //courseModule
  isCourseModule: boolean

  // totalCourseModules: number

  //course
  courseList: ICourseList[]
  isCourseLoaded: boolean

  // spinner

  isResourceLoading: boolean


  // modal
  modalRef: NgbModalRef

  //multi-step-selector
  progress1: any = 0;
  progress2: any = 0;

  //courseForm
  courseSelectionForm: FormGroup

  // course-details
  courseID: string
  courseName: string

  @ViewChild('deleteConfirmationModal') deleteConfirmationModal: any

  constructor(
    private _location: Location,
    private formBuilder: FormBuilder,
    private spinnerService: SpinnerService,
    private courseService: CourseService,
    private utilService: UtilityService,
    private urlConstant: UrlConstant,
    private router: Router,
    private route: ActivatedRoute,
  ) {
    this.initializeVariables()
    this.createForms()
    this.getAllCourseList()

  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {

  }

  initializeVariables(): void {
    this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_COURSE_MASTER)

    this.courseList = []
    this.extractCourseDetails()
    this.isCourseLoaded = false
    this.isCourseModule = false

  }

  changedCourse(): void {

    this.mergeSearchToQueryParams()
    if (this.courseDetailsComponent != null) {
      this.courseDetailsComponent.getQueryParm()
    }


  }

  createForms(): void {
    this.createCourseSelectionForm()

  }

  createCourseSelectionForm(): void {

    this.courseSelectionForm = this.formBuilder.group({
      courseID: new FormControl(null)
    })
  }

  changeProgress1() {
    console.log("in course change progress");

    if (this.progress1 == 0) {
      this.progress1 = 100;
      return;
    }

    if (this.progress1 == 100) {
      this.progress2 = 0;
      return;
    }
  }

  changeProgress2() {
    this.progress1 = 100;
    this.progress2 = 100;
  }

  changeProgress() {

    this.progress1 = 0;
    this.progress2 = 0;
  }


  mergeSearchToQueryParams(): void {
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: this.courseSelectionForm.value,
      queryParamsHandling: 'merge'
    });
  }

  toggleCourseModule(): void {
    this.isCourseModule = !this.isCourseModule
  }
  extractCourseDetails(): void {
    this.route.queryParamMap.subscribe(params => {
      this.courseID = params.get('courseID')
      this.courseName = params.get('courseName')
    })
  }

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
    this._location.back();
  }
  // =============================================================== CRUD ===============================================================   

  getAllCourseList(): void {
    this.spinnerService.loadingMessage = "Getting Course Modules"


    this.courseService.getCourseList().subscribe((response: any) => {

      this.courseList = response.body

      this.courseSelectionForm.get("courseID").setValue(this.courseID)
      this.mergeSearchToQueryParams()
      this.isCourseLoaded = true
    }, (error) => {
      console.error(error)
    }

    ).add(() => {


    })

  }
}

