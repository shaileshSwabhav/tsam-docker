import { DatePipe } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { ILeaderBoard, IPerformer, IProblemOfTheDayQuestion, ProblemOfTheDayService } from 'src/app/service/problem-of-the-day/problem-of-the-day.service';
import { LocalService } from 'src/app/service/storage/local.service';

@Component({
  selector: 'app-problem-of-the-day',
  templateUrl: './problem-of-the-day.component.html',
  styleUrls: ['./problem-of-the-day.component.css']
})
export class ProblemOfTheDayComponent implements OnInit {

  // Questions.
  problemList: IProblemOfTheDayQuestion[]

  // Spinner.



  // Leader.
  leaderBoard: ILeaderBoard
  performerList: IPerformer[]

  // Date.
  problemOfTheDayDate: Date

  // Flags.
  isSearched: boolean

  // Search.
  searchDateInString: Date

  constructor(
    private spinnerService: SpinnerService,
    private router: Router,
    private problemOfTheDayService: ProblemOfTheDayService,
    private datePipe: DatePipe,
    private localService: LocalService,
  ) {
    this.initializeVariables()
    this.getProblemsOfTheDay()
    this.getLeaderBoard()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize global variables.
  initializeVariables(): void {

    // Date.
    this.problemOfTheDayDate = new Date()

    // Questions.
    this.problemList = []

    // Leader.
    this.performerList = []

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Problems"

  }

  // =============================================================SEARCH FUNCTIONS==========================================================================

  // On clicking search button.
  onSearchButtonClick(searchDateInString: any): void {

    // If there is no input value.
    if (searchDateInString == "" || searchDateInString == null) {
      return
    }

    this.searchDateInString = searchDateInString
    let searchDate: Date = new Date(searchDateInString)
    let todayDate: Date = new Date()
    let todayDateInString: string = this.datePipe.transform(todayDate, 'yyyy-MM-dd');
    let diffTime = Math.round((todayDate.getTime() - searchDate.getTime()) / (1000 * 60 * 60 * 24))

    // If search date is after today's date or is equal to todays's date then give alert.
    if (diffTime < 0) {
      alert("Please select a date before todays's date")
      return
    }

    // If search date is today's date then give today's problem of the day questions.
    if (searchDateInString == todayDateInString) {
      searchDateInString = null
      this.problemOfTheDayDate = new Date()
      this.getProblemsOfTheDay()
      this.getLeaderBoard()
      return
    }

    // If search date is before 31 days and more from today's date.
    if (diffTime > 30) {
      alert("Please select a date within 30 days from today's date")
      return
    }

    this.getProblemsOfThePreviousDay()
    this.getLeaderBoardPreviousDays()
    this.problemOfTheDayDate = searchDate
  }

  // =============================================================FORMAT FUNCTIONS==========================================================================

  // Format fields of problem of the day questions.
  formatProblemOfTheDayQuestionFields(): void {

    for (let i = 0; i < this.problemList.length; i++) {

      // Set class for difficulty.
      if (this.problemList[i].level == 1) {
        this.problemList[i].levelName = "Easy"
        this.problemList[i].levelClass = "card-sub-detail easy"
      }
      if (this.problemList[i].level == 2) {
        this.problemList[i].levelName = "Medium"
        this.problemList[i].levelClass = "card-sub-detail medium"
      }
      if (this.problemList[i].level == 3) {
        this.problemList[i].levelName = "Hard"
        this.problemList[i].levelClass = "card-sub-detail hard"
      }

      // Set success ratio.
      if (this.problemList[i].attemptedByCount != 0) {
        this.problemList[i].successRatio = ((this.problemList[i].solvedByCount / this.problemList[i].attemptedByCount) * 100).toFixed(2)
        if (this.problemList[i].successRatio == "0.00") {
          this.problemList[i].successRatio = "0"
        }
      }
      else {
        this.problemList[i].successRatio = "0"
      }
    }
  }

  // Format leader board.
  formatLeaderBoard(): void {

    // Put all performers in one list.
    this.performerList = this.leaderBoard.allPerformers

    // If self performer exists then push it in performer list.
    this.performerList.splice(0, 0, this.leaderBoard.selfPerformer)

    // Give default image to talents that dont have an image.
    for (let i = 0; i < this.performerList.length; i++) {
      if (!this.performerList[i].image) {
        this.performerList[i].image = "assets/images/default-profile-image.png"
      }
    }
  }

  // =============================================================REDIRECT FUNCTIONS==========================================================================

  // Redirect to problem details page.
  redirectToProblemDetails(problemID: string): void {
    this.router.navigate(['/problems-of-the-day/problem-details'], {
      queryParams: {
        "problemID": problemID,
        "searchDate": this.searchDateInString
      }
    }).catch(err => {
      console.error(err)
    })
  }

  // =============================================================GET FUNCTIONS==========================================================================

  // Get problem of the day questions.
  getProblemsOfTheDay(): void {
    this.spinnerService.loadingMessage = "Getting Problems Of The Day"


    this.problemOfTheDayService.getProblemsOfTheDay().subscribe((response) => {
      this.problemList = response
      this.formatProblemOfTheDayQuestionFields()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get problem of the day previous days questions.
  getProblemsOfThePreviousDay(): void {
    this.spinnerService.loadingMessage = "Getting Problems Of The Day"


    let queryParams: any = {
      "searchDate": this.searchDateInString
    }
    this.problemOfTheDayService.getProblemsOfThePreviousDay(queryParams).subscribe((response) => {
      this.problemList = response
      this.formatProblemOfTheDayQuestionFields()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get leader board.
  getLeaderBoard(): void {
    this.spinnerService.loadingMessage = "Getting Problems Of The Day"


    this.problemOfTheDayService.getLeaderBoardPotd(this.localService.getJsonValue("loginID")).subscribe((response) => {
      this.leaderBoard = response
      this.formatLeaderBoard()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get leader board of previous days.
  getLeaderBoardPreviousDays(): void {
    this.spinnerService.loadingMessage = "Getting Problems Of The Day"


    let queryParams: any = {
      "searchDate": this.searchDateInString
    }
    this.problemOfTheDayService.getLeaderBoardPotd(this.localService.getJsonValue("loginID"), queryParams).subscribe((response) => {
      this.leaderBoard = response
      this.formatLeaderBoard()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }
}
