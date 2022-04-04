import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { CourseService } from 'src/app/service/course/course.service';
import { LocalService } from 'src/app/service/storage/local.service';

@Component({
  selector: 'app-my-courses',
  templateUrl: './my-courses.component.html',
  styleUrls: ['./my-courses.component.css']
})
export class MyCoursesComponent implements OnInit {

  // Course.
  courseList: any[]

  // Flags.
  hideCertificate: boolean

  // Pagination.
  limit: number
  currentPage: number
  offset: number
  paginationString: string
  totalCourses: number

  // Spinner.



  constructor(
    private spinnerService: SpinnerService,
    private localService: LocalService,
    private router: Router,
    private courseService: CourseService,
  ) {
    this.initializeVariables()
    this.getAllCourses()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize global variables.
  initializeVariables(): void {

    // Course.
    this.courseList = []

    // Flags.
    this.hideCertificate = false

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Courses"


    // Pagination.
    this.limit = 8
    this.offset = 0
    this.currentPage = 0
  }

  // Format fields of all courses.
  formatCourseFields(): void {
    for (let i = 0; i < this.courseList.length; i++) {
      this.courseList[i].completedPercentage = this.courseList[i].completedPercentage.toFixed(2);
    }
  }

  // Get all courses.
  getAllCourses(): void {
    this.spinnerService.loadingMessage = "Getting Courses"


    this.courseService.getMyCoursesByTalent(this.localService.getJsonValue("loginID")).subscribe((response) => {
      this.courseList = response
      this.totalCourses = this.courseList.length
      this.formatCourseFields()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.setPaginationString()
    })
  }

  // Set total course list on current page.
  setPaginationString() {
    this.paginationString = ''
    let start: number = this.limit * this.offset + 1
    let end: number = +this.limit + this.limit * this.offset
    if (this.totalCourses < end) {
      end = this.totalCourses
    }
    if (this.totalCourses == 0) {
      this.paginationString = ''
      return
    }
    this.paginationString = `${start} - ${end} of ${this.totalCourses}`
  }

  // On page change.
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
    this.getAllCourses()
  }

  // Redirect to my sessions page.
  redirectToMySessions(batchID: string, batchName: string): void {
    this.router.navigate(['my-batches/my-sessions'], {
      queryParams: {
        "batchID": batchID,
        "batchName": batchName
      }
    }).catch(err => {
      console.error(err)
    })
  }

  // Redirect to batch feedback page.
  redirectToBatchFeedback(batchID: string, batchName: string): void {
    this.router.navigate(['/batch/master/feedback'], {
      queryParams: {
        "batchID": batchID,
        "batchName": batchName,
        "talentID": this.localService.getJsonValue("loginID")
      }
    }).catch(err => {
      console.error(err)
    })
  }

}
