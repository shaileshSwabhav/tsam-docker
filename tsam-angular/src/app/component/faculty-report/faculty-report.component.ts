import { DatePipe } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { IBatch, IFacultyReport, IWorkingHours, ReportService } from 'src/app/service/report/report.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-faculty-report',
  templateUrl: './faculty-report.component.html',
  styleUrls: ['./faculty-report.component.css']
})
export class FacultyReportComponent implements OnInit {

  // faculty-report
  facultyReport: IFacultyReport[]

  // flags
  isReportLoaded: boolean

  // spinner


  weekIndex: string

  // date
  date: string
  userDate: string

  readonly MONDAY: string = "Monday"
  readonly TUESDAY: string = "Tueday"
  readonly WEDNESDAY: string = "Wednesday"
  readonly THURSDAY: string = "Thursday"
  readonly FRIDAY: string = "Friday"
  readonly SATURDAY: string = "Saturday"
  readonly SUNDAY: string = "Sunday"

  constructor(
    private spinnerService: SpinnerService,
    private reportService: ReportService,
    public utilService: UtilityService,
    private router: Router,
    private route: ActivatedRoute,
    private datePipe: DatePipe,
  ) {
    this.initializeVariables()
    this.addWeekIndex()
  }

  initializeVariables(): void {
    this.facultyReport = []
    this.isReportLoaded = true
    this.spinnerService.loadingMessage = "Getting faculty report"
    this.weekIndex = "0"
    this.date = this.getMonday().toISOString()
    this.userDate = this.datePipe.transform(this.getMonday(), 'yyyy-MM-dd')
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  doesWorkingHoursExist(workingHours: Map<string, IWorkingHours>): boolean {
    return Object.keys(workingHours).length != 0
  }

  onWeekSelectClick(): void {
    this.userDate = null
    this.userDate = this.datePipe.transform(this.getMonday(), 'yyyy-MM-dd')
    this.date = new Date(this.userDate).toISOString()
    this.addQueryParams()
  }

  onDateSelectClick(): void {
    this.weekIndex = null
    let currentDate = this.getMonday()
    currentDate = new Date(this.userDate)
    // console.log(currentDate);
    this.date = currentDate.toISOString()

    this.addQueryParams()
  }

  addQueryParams(): void {
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: {
        date: this.datePipe.transform(new Date(this.date), "yyyy-MM-dd")
      },
    })
    this.getAllFacultyReport()
  }

  addWeekIndex(): void {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.date = this.getMonday().toISOString()
      this.userDate = this.date
      this.onWeekSelectClick()
      return
    }
    this.date = new Date(queryParams.date).toISOString()
    this.userDate = queryParams.date
    this.getAllFacultyReport()
  }

  /**
  * getMonday
  * @returns date depending on the weekIndex value, 0 is for this week,
  *  -7 is for previous week and 7 for next week.
  */
  getMonday(): Date {
    let date = new Date();
    let day = date.getDay();
    let finalDate = new Date();

    let offset = 1 - day;
    if (offset > 0) {
      offset = -6;
    }

    offset += (+this.weekIndex)
    finalDate.setDate(new Date().getDate() + offset);
    return finalDate
  }

  getAllFacultyReport(): void {
    this.spinnerService.loadingMessage = "Getting faculty report"

    this.facultyReport = []
    this.isReportLoaded = true
    // console.log(this.date);

    let queryParams: any = {
      date: this.date
    }
    // console.log(queryParams);

    this.reportService.getAllFacultyReport(queryParams).subscribe((response: any) => {
      this.facultyReport = response.body
      this.isReportLoaded = true
      if (this.facultyReport.length == 0) {
        this.isReportLoaded = false
      }
      // console.log(response.body);
    }, (err: any) => {
      console.error(err);
      this.isReportLoaded = false
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }


  // ============================================= CALCULATE TIME DIFFERENCE =============================================

  assignTotalBatchHours(): number {
    let hours: number = 0
    for (let k = 0; k < this.facultyReport.length; k++) {
      Object.keys(this.facultyReport[k].workingHours).map(key => {
        hours += this.facultyReport[k].workingHours[key].totalHours
      });
    }
    return hours
  }

  assignTotalHours(): number {
    let hours: number = 0
    for (let k = 0; k < this.facultyReport.length; k++) {
      hours += this.facultyReport[k].totalTrainingHours
    }
    return hours
  }

  assignTotalWorkingHours(day: string): string {
    let hours: number = 0

    for (let k = 0; k < this.facultyReport.length; k++) {

      if (day == this.MONDAY) {
        hours += this.calculateTimeDifference(this.facultyReport[k]?.monday)
      }

      if (day == this.TUESDAY) {
        hours += this.calculateTimeDifference(this.facultyReport[k]?.tuesday)
      }

      if (day == this.WEDNESDAY) {
        hours += this.calculateTimeDifference(this.facultyReport[k]?.wednesday)
      }

      if (day == this.THURSDAY) {
        hours += this.calculateTimeDifference(this.facultyReport[k]?.thursday)
      }

      if (day == this.FRIDAY) {
        hours += this.calculateTimeDifference(this.facultyReport[k]?.friday)
      }

      if (day == this.SATURDAY) {
        hours += this.calculateTimeDifference(this.facultyReport[k]?.saturday)
      }

      if (day == this.SUNDAY) {
        hours += this.calculateTimeDifference(this.facultyReport[k]?.sunday)
      }
    }
    return hours.toString()
  }

  calculateTimeDifference(batches: IBatch[]): number {
    let hourDiff: number = 0
    if (batches) {
      for (let k = 0; k < batches?.length; k++) {
        hourDiff += batches[k].totalDailyHours
      }
    }
    return hourDiff
  }

}
