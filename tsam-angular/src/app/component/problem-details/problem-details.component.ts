import { Component, OnInit, QueryList, ViewChild, ViewChildren } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { IProgrammingQuestion, IProgrammingQuestionSolutionDTO, ProgrammingQuestionService } from 'src/app/service/programming-question/programming-question.service';
import { IProblemOfTheDayQuestion, ProblemOfTheDayService } from 'src/app/service/problem-of-the-day/problem-of-the-day.service';
import { IProgrammingQuestionTalentAnswer, IProgrammingQuestionSolutionIsViewed, ProgrammingQuestionTalentAnswerService } from 'src/app/service/programming-question-talent-answer/programming-question-talent-answer.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { CodeEditorComponent } from '../code-editor/code-editor.component';
import { NgbNav } from '@ng-bootstrap/ng-bootstrap';

// Default language.
const DEFAULT_LANGUAGE = "C"

@Component({
  selector: 'app-problem-details',
  templateUrl: './problem-details.component.html',
  styleUrls: ['./problem-details.component.css']
})
export class ProblemDetailsComponent implements OnInit {

  // Components.
  languageList: any[]
  solutionLanguageList: any[]

  // Problem.
  problemID: string
  problem: IProgrammingQuestion
  problemList: IProblemOfTheDayQuestion[]
  selectedProblem: IProgrammingQuestion

  // Spinner.



  // Sub test.
  subTestID: string
  showTimer: boolean

  // Option.
  selectedOption: any
  correctOption: any

  // Answer.
  answer: IProgrammingQuestionTalentAnswer
  talentAnswer: string

  // Flags.
  isSolution: boolean
  isSolutionViewed: boolean
  isOptionViewed: boolean
  isSelectedProblemChanged: boolean

  // Search.
  searchDateInString: string

  // Programming language.
  selecetedLanguage: any

  // Language and Code.
  languageAndCodeAnswer: any
  languageAndCodeSolution: any

  // Programming question solution.
  selectedQuestionSolution: IProgrammingQuestionSolutionDTO

  // Solution is viewed.
  solutionIsViewed: IProgrammingQuestionSolutionIsViewed

  // Nav bar.
  @ViewChild("nav") navBarElement: NgbNav

  // Child component.
  @ViewChildren(CodeEditorComponent) codeEditorList: QueryList<CodeEditorComponent>
  @ViewChild('codeEditorAnswer') codeEditorAnswer: CodeEditorComponent
  @ViewChild('codeEditorSolution') codeEditorSolution: CodeEditorComponent

  constructor(
    private activatedRoute: ActivatedRoute,
    private spinnerService: SpinnerService,
    private questionService: ProgrammingQuestionService,
    private answerService: ProgrammingQuestionTalentAnswerService,
    private problemOfTheDayService: ProblemOfTheDayService,
    private router: Router,
    private localService: LocalService,
    private generalService: GeneralService,
  ) {
    this.initializeVariables()
    this.getProgrammingLanguageList()
    this.getProblemsDetails()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize global variables.
  initializeVariables(): void {

    // Components.
    this.languageList = []
    this.solutionLanguageList = []

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Problem Details"


    // Flags.
    this.isSolution = false
    this.isSolutionViewed = false
    this.isOptionViewed = false
    this.isSelectedProblemChanged = false

    // Sub test.
    this.showTimer = false
    this.subTestID = this.activatedRoute.snapshot.queryParamMap.get("subTestID")
    if (this.subTestID) {
      this.showTimer = true
    }

    // Problem.
    this.problemList = []
    this.problemID = this.activatedRoute.snapshot.queryParamMap.get("problemID")
    this.searchDateInString = this.activatedRoute.snapshot.queryParamMap.get("searchDate")
    if (this.searchDateInString) {
      this.getProblemsOfThePreviousDay()
      return
    }
    this.getProblemsOfTheDay()
  }

  // =============================================================FORMAT FUNCTIONS==========================================================================

  // Format fields of problem details.
  formatProblemFields(): void {

    // Set class for difficulty.
    if (this.problem.level == 1) {
      this.problem.levelName = "Easy"
      this.problem.levelClass = "easy"
    }
    if (this.problem.level == 2) {
      this.problem.levelName = "Medium"
      this.problem.levelClass = "medium"
    }
    if (this.problem.level == 3) {
      this.problem.levelName = "Hard"
      this.problem.levelClass = "hard"
    }

    // If problem has options then set the options' class.
    if (this.problem.hasOptions) {
      this.isSolution = false
      for (let i = 0; i < this.problem.options.length; i++) {
        if (this.problem.options[i].isCorrect) {
          this.correctOption = this.problem.options[i]
        }
        if (this.problem.isAnswered && this.problem.programmingQuestionOptionID == this.problem.options[i].id) {
          this.problem.options[i].optionClass = "card selected-option-style h-100"
          continue
        }
        this.problem.options[i].optionClass = "card option-style h-100"
      }
    }

    // If problem does not have options then make selectedOption as null.
    if (!this.problem.hasOptions) {
      this.selectedOption = null
      this.isSolution = true
    }

    // Set success ratio of problem.
    if (this.problem.attemptedByCount != 0) {
      this.problem.successRatio = ((this.problem.solvedByCount / this.problem.attemptedByCount) * 100).toFixed(2)
      if (this.problem.successRatio == "0.00") {
        this.problem.successRatio = "0"
      }
    }
    else {
      this.problem.successRatio = "0"
    }

    // If problem has options and solution is viewed then show the correct option.
    if (this.problem.hasOptions && this.problem.solutonIsViewed) {
      this.isOptionViewed = true
    }

    // If problem has options then dont create instance of code editor.
    if (this.problem.hasOptions) {
      return
    }

    // Set the solution language list to languages that the solutions of the problem have.
    this.solutionLanguageList = []
    if (this.problem.solutions?.length > 0) {
      for (let i = 0; i < this.problem.solutions.length; i++) {

        // Set the fisrt solution as the selected programming question solution.
        if (i == 0) {
          this.selectedQuestionSolution = this.problem.solutions[i]
        }
        this.solutionLanguageList.push(this.problem.solutions[i].programmingLanguage)
      }
    }

    // If code editor answer and code editor solution instance is not created then create it first.
    if (!this.codeEditorAnswer || !this.codeEditorSolution) {
      this.codeEditorList.changes.subscribe((tempCodeEditorList: QueryList<CodeEditorComponent>) => {
        // If only one code editor exists and search date string is null then it is code editor answer. 
        if (tempCodeEditorList.length == 1 && !this.searchDateInString) {
          this.codeEditorAnswer = tempCodeEditorList.first
        }

        // If only one code editor exists and search date string is not null and problem is answered then it is code editor solution. 
        if (tempCodeEditorList.length == 1 && this.searchDateInString && this.problem.isAnswered) {
          this.codeEditorAnswer = tempCodeEditorList.first
        }

        // If only one code editor exists and search date string is not null and problem is not answered then it is code editor solution. 
        if (tempCodeEditorList.length == 1 && this.searchDateInString && !this.problem.isAnswered) {
          this.codeEditorSolution = tempCodeEditorList.first
        }

        // If two code editors exist then it is code editor answer and code editor solution. 
        if (tempCodeEditorList.length == 2) {
          this.codeEditorAnswer = tempCodeEditorList.last
          this.codeEditorSolution = tempCodeEditorList.first
        }

        // If code editor answer exists then send values to code editor answer.
        if (this.codeEditorAnswer && !this.problem.hasOptions) {

          // // If selected problem is null then set the current problem as selected problem.
          // if (!this.selectedProblem){
          //   this.selectedProblem = this.problem
          //   this.isSelectedProblemChanged = false
          // }

          // // If previous problem is not same as current problem then set isSelectedProblemChanged as true, else false
          // if (this.selectedProblem && (this.selectedProblem.id != this.problem.id)){
          //   this.isSelectedProblemChanged = true
          //   this.selectedProblem = this.problem
          // }
          // else if (this.selectedProblem && this.selectedProblem.id == this.problem.id){
          //   this.isSelectedProblemChanged = false
          // }

          this.sendValuesToCodeEditorAnswer()
        }

        // If code editor solution instance exists and solution is viewed then call its function.
        if (this.codeEditorSolution && (this.problem.solutonIsViewed || this.solutionIsViewed) && !this.problem.hasOptions) {
          this.setProgrammingQuestionSolution()
        }
      })
    }

    // If selected problem is null then set the current problem as selected problem.
    if (!this.selectedProblem) {
      this.selectedProblem = this.problem
      this.isSelectedProblemChanged = false
    }

    // If previous problem is not same as current problem then set isSelectedProblemChanged as true, else false
    if (this.selectedProblem && (this.selectedProblem.id != this.problem.id)) {
      this.isSelectedProblemChanged = true
      this.selectedProblem = this.problem
    }
    else if (this.selectedProblem && this.selectedProblem.id == this.problem.id) {
      this.isSelectedProblemChanged = false
    }

    // If code editor answer instance exists and it is for todays date or if previous days' answer exists then call its function.
    if ((this.codeEditorAnswer && !this.searchDateInString) || (this.codeEditorAnswer && this.searchDateInString && this.problem.isAnswered)) {
      this.sendValuesToCodeEditorAnswer()
    }

    // If problem solutions is already viewed by talent then show the solutions.
    if (this.problem.solutonIsViewed) {
      this.isSolutionViewed = true

      // If nav bar exists then select the first tab as active tab.
      if (this.navBarElement) {
        this.navBarElement.select(this.solutionLanguageList[0].id)
      }
    }

    // If problem solutions is not viewed by talent then hide the solutions.
    if (!this.problem.solutonIsViewed) {
      this.isSolutionViewed = false
    }

    // // If code editor solution instance exists then call its function.
    // if (this.codeEditorSolution){
    //   this.setProgrammingSolutionQuestion()
    // }
  }

  // Format fields of problem list.
  formatProblemListFields(): void {

    for (let i = 0; i < this.problemList.length; i++) {

      // Set class for problem.
      if (this.problemList[i].id == this.problemID) {
        this.problemList[i].problemClass = "card problem-selected-card-style h-100"
        continue
      }
      this.problemList[i].problemClass = "card problem-disabled-card-style h-100"
    }
  }

  // =============================================================BUTTON CLICK FUNCTIONS==========================================================================

  // On clicking problem label.
  onProblemLabelClick(problemID: string): void {
    if (this.problemID == problemID) {
      return
    }
    this.selecetedLanguage = null
    this.router.navigate(['/problems-of-the-day/problem-details'], {
      queryParams: {
        "problemID": problemID,
      }
    }).catch(err => {
      console.error(err)
    })
    this.problemID = problemID
    this.formatProblemListFields()
    this.getProblemsDetails()
  }

  // On clicking option.
  onOptionCliick(option: any): void {
    if (this.problem.isAnswered) {
      return
    }
    if (this.selectedOption && this.selectedOption.id == option.id) {
      return
    }
    this.selectedOption = option
    for (let i = 0; i < this.problem.options.length; i++) {
      if (this.selectedOption.id == this.problem.options[i].id) {
        this.problem.options[i].optionClass = "card selected-option-style h-100"
        continue
      }
      this.problem.options[i].optionClass = "card option-style h-100"
    }
    return
  }

  // On clicking option submit button.
  onSubmitOptionButtonClick(): void {

    // Get programming question type id from programming question for practice type.
    let tempQuestionTypeID: string
    for (let i = 0; i < this.problem.programmingQuestionTypes.length; i++) {
      if (this.problem.programmingQuestionTypes[i].programmingType == "Practice") {
        tempQuestionTypeID = this.problem.programmingQuestionTypes[i].id
      }
    }

    this.answer = {
      answer: null,
      score: 0,
      programmingQuestionID: this.problem.id,
      talentID: this.localService.getJsonValue("loginID"),
      programmingQuestionOptionID: this.selectedOption.id,
      isCorrect: this.selectedOption.isCorrect,
      programmingLanguageID: null,
      programmingQuestionTypeID: tempQuestionTypeID
    }

    // If selected option is correcr then give full score to answer.
    if (this.selectedOption.isCorrect) {
      this.answer.score = this.problem.score
    }

    this.addAnswer()
  }

  // On clicking language tab.
  onLanguageTabClick(languageID: string): void {

    // If problem has solutions then get the solution for the language tab selected.
    if (this.problem.solutions?.length > 0) {
      for (let i = 0; i < this.problem.solutions.length; i++) {
        if (this.problem.solutions[i].programmingLanguage.id == languageID) {
          this.selectedQuestionSolution = this.problem.solutions[i]
          this.sendValuesToCodeEditorSolution()
          break
        }
      }
    }
  }

  // On clicking view solutions button.
  onViewSolutionsClick(): void {
    if (confirm("Once you view the solutions you will not get any score for your answer; Are you sure you want to view the solution(s) ?")) {

      // Add the solution is viewed and then show the solutions.
      this.solutionIsViewed = {
        programmingQuestionID: this.problem.id,
        talentID: this.localService.getJsonValue("loginID"),
      }
      this.addSolutionIsViewed()
    }
  }

  // On clicking view correct option button.
  onViewCorrectOptionClick(): void {
    if (confirm("Once you view the correct option you will not get any score for your answer; Are you sure you want to view the correct option ?")) {

      // Add the solution is viewed and then show the solutions.
      this.solutionIsViewed = {
        programmingQuestionID: this.problem.id,
        talentID: this.localService.getJsonValue("loginID"),
      }
      this.addSolutionIsViewed()
    }
  }

  // =============================================================REDIRECT FUNCTIONS==========================================================================

  // On clicking back button.
  redirectToProblemOfTheDay(): void {
    this.router.navigate(['/problems-of-the-day'], {
    }).catch(err => {
      console.error(err)
    })
  }

  // =============================================================ADD FUNCTIONS==========================================================================

  // Add new answer.
  addAnswer(): void {
    this.spinnerService.loadingMessage = "Submitting Your Answer"


    this.answerService.addAnswer(this.answer).subscribe((response: any) => {
      alert("Your answer has been submitted successfully")
      this.getProblemsDetails()
    }, (error) => {
      console.error(error)
      alert("There was some error while submitting your answer, try again later")
    })
  }

  // Add new solution is viewed.
  addSolutionIsViewed(): void {
    this.spinnerService.loadingMessage = "Showing Solution"


    this.answerService.addSolutionIsViewed(this.solutionIsViewed).subscribe((response: any) => {
      if (!this.problem.hasOptions) {
        this.isSolutionViewed = true
      }
      if (this.problem.hasOptions) {
        this.isOptionViewed = true
      }
    }, (error) => {
      console.error(error)
      alert("There was some error while showing solution, try again later")
    })
  }

  // =============================================================OTHER FUNCTIONS==========================================================================

  // Receive values from child component.
  receiveChildValues(tempLanguageAndCode: any) {

    // Give code to talent answer recieved from child component.
    this.talentAnswer = tempLanguageAndCode.code

    // Get language id from langauge name received from child component.
    for (let i = 0; i < this.languageList.length; i++) {
      if (tempLanguageAndCode.languageName == this.languageList[i].name) {
        this.selecetedLanguage = this.languageList[i]
        break
      }
    }

    // If talent answer is null, it means language has been changed in child component, get problem details for the language.
    if (this.talentAnswer == null) {
      this.getProblemsDetails()
      return
    }

    // Submit answer.
    this.onSubmittingAnswer()
  }

  // On submitting talent answer.
  onSubmittingAnswer(): void {

    // Get programming question type id from programming question for practice type.
    let tempQuestionTypeID: string
    for (let i = 0; i < this.problem.programmingQuestionTypes.length; i++) {
      if (this.problem.programmingQuestionTypes[i].programmingType == "Practice") {
        tempQuestionTypeID = this.problem.programmingQuestionTypes[i].id
      }
    }

    this.answer = {
      answer: this.talentAnswer,
      score: 0,
      programmingQuestionID: this.problem.id,
      talentID: this.localService.getJsonValue("loginID"),
      programmingQuestionOptionID: null,
      isCorrect: null,
      programmingLanguageID: this.selecetedLanguage.id,
      programmingQuestionTypeID: tempQuestionTypeID
    }
    this.addAnswer()
  }

  // Send values to code editor answer after getting problem details.
  sendValuesToCodeEditorAnswer(): void {

    let isPreviousDay: boolean = false
    if (this.searchDateInString) {
      isPreviousDay = true
    }

    // If problem is answered then send language and code to child component.
    if (this.problem.isAnswered) {
      this.languageAndCodeAnswer = {
        languageName: this.problem.programmingLanguage.name,
        code: this.problem.answer,
        isSolution: false,
        isReadOnly: isPreviousDay,
        isSelectedProblemChanged: this.isSelectedProblemChanged,
        testCases: this.problem.testCases
      }
      this.codeEditorAnswer.valuesFromParentChange(this.languageAndCodeAnswer)
    }

    // If problem is not answered then send language to child component.
    if (!this.problem.isAnswered) {
      let tempLanguageName: string
      if (this.selecetedLanguage) {
        tempLanguageName = this.selecetedLanguage.name
      }
      if (!this.selecetedLanguage) {
        tempLanguageName = DEFAULT_LANGUAGE
      }
      this.languageAndCodeAnswer = {
        languageName: tempLanguageName,
        code: null,
        isSolution: false,
        isPreviousDay: isPreviousDay,
        isSelectedProblemChanged: this.isSelectedProblemChanged,
        testCases: this.problem.testCases
      }
      this.codeEditorAnswer.valuesFromParentChange(this.languageAndCodeAnswer)
    }
  }

  // Send values to code editor solution after getting problem details and language tab change.
  sendValuesToCodeEditorSolution(): void {

    this.languageAndCodeSolution = {
      languageName: this.selectedQuestionSolution.programmingLanguage.name,
      code: this.selectedQuestionSolution.solution,
      isSolution: true,
      isReadOnly: true
    }
    this.codeEditorSolution.valuesFromParentChange(this.languageAndCodeSolution)
  }

  // Set the initial selected programming question solution.
  setProgrammingQuestionSolution(): void {

    // Set the soltion language list to languages that the solutions of the problem have.
    if (this.problem.solutions?.length > 0) {
      this.solutionLanguageList = []
      for (let i = 0; i < this.problem.solutions.length; i++) {

        // Set the fisrt solution as the selected programming question solution.
        if (i == 0) {
          this.selectedQuestionSolution = this.problem.solutions[i]
        }
        this.solutionLanguageList.push(this.problem.solutions[i].programmingLanguage)
      }

      // Send the values to code editor solution.
      if (this.selectedQuestionSolution) {
        this.sendValuesToCodeEditorSolution()
      }
    }
  }

  // =============================================================GET FUNCTIONS==========================================================================

  // Get problem details.
  getProblemsDetails(): void {
    this.spinnerService.loadingMessage = "Getting Problem Details"


    let queryParams: any = {
      "talentID": this.localService.getJsonValue("loginID")
    }
    if (this.selecetedLanguage) {
      queryParams.programmingLanguageID = this.selecetedLanguage.id
    }
    this.questionService.getProgrammingQuestion(this.problemID, queryParams).subscribe((response) => {
      this.problem = response
      this.formatProblemFields()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get problem of the day questions.
  getProblemsOfTheDay(): void {
    this.spinnerService.loadingMessage = "Getting Problem Details"


    this.problemOfTheDayService.getProblemsOfTheDay().subscribe((response) => {
      this.problemList = response
      this.formatProblemListFields()
    }, error => {
      console.error(error)
    })
  }

  // Get problem of the day questions.
  getProblemsOfThePreviousDay(): void {
    this.spinnerService.loadingMessage = "Getting Problems Of The Day"


    let queryParams: any = {
      "searchDate": this.searchDateInString
    }
    this.problemOfTheDayService.getProblemsOfThePreviousDay(queryParams).subscribe((response) => {
      this.problemList = response
      this.formatProblemListFields()
    }, error => {
      console.error(error)
    })
  }

  // Get programming language list.
  getProgrammingLanguageList(): void {
    this.generalService.getProgrammingLanguageList().subscribe((response: any) => {
      this.languageList = response
    }, (err: any) => {
      console.error(err)
    })
  }
}
