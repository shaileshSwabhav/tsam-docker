import { Component, OnInit, QueryList, ViewChild, ViewChildren } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbNav } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { Role, UrlConstant } from 'src/app/service/constant';
import { GeneralService } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { IProgrammingQuestionSolutionDTO } from 'src/app/service/programming-question/programming-question.service';
import { IProgrammingQuestionTalentAnswerWithFullQuestionDTO, ProgrammingQuestionTalentAnswerService } from 'src/app/service/programming-question-talent-answer/programming-question-talent-answer.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { CodeEditorComponent } from '../code-editor/code-editor.component';

@Component({
  selector: 'app-programming-question-talent-answer-details',
  templateUrl: './programming-question-talent-answer-details.component.html',
  styleUrls: ['./programming-question-talent-answer-details.component.css']
})
export class ProgrammingQuestionTalentAnswerDetailsComponent implements OnInit {

  // Components.
  languageList: any[]
  solutionLanguageList: any[]

  // Programming question talent answer.
  answerID: string
  answer: IProgrammingQuestionTalentAnswerWithFullQuestionDTO

  // Spinner.



  // Permission.
  permission: IPermission
  roleName: string

  // Programming language.
  selecetedLanguage: any

  // Language and Code.
  languageAndCodeAnswer: any
  languageAndCodeSolution: any

  // Programming question solution.
  selectedQuestionSolution: IProgrammingQuestionSolutionDTO

  // Nav bar.
  @ViewChild("nav") navBarElement: NgbNav

  // Score.
  scoreForm: FormGroup

  // Child component.
  @ViewChildren(CodeEditorComponent) codeEditorList: QueryList<CodeEditorComponent>
  @ViewChild('codeEditorAnswer') codeEditorAnswer: CodeEditorComponent
  @ViewChild('codeEditorSolution') codeEditorSolution: CodeEditorComponent

  constructor(
    private activatedRoute: ActivatedRoute,
    private spinnerService: SpinnerService,
    private answerService: ProgrammingQuestionTalentAnswerService,
    private router: Router,
    private localService: LocalService,
    private role: Role,
    public utilService: UtilityService,
    private urlConstant: UrlConstant,
    private generalService: GeneralService,
    private formBuilder: FormBuilder,
  ) {
    this.initializeVariables()
    this.getAnswerDetails()
    this.getProgrammingLanguageList()
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
    this.spinnerService.loadingMessage = "Getting Talent Answer Details"


    // Problem.
    this.answerID = this.activatedRoute.snapshot.queryParamMap.get("answerID")

    // Score.
    this.createScoreForm()

    // Permision.
    // Get permissions from menus using utilityService function.
    if (this.localService.getJsonValue("roleName") == this.role.ADMIN) {
      this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.ADMIN_PROGRAMMING_QUESTION_TALENT_ANSWER_DETAILS)
    }
    else {
      this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.PROGRAMMING_QUESTION_TALENT_ANSWER_DETAILS)
    }
    this.roleName = this.localService.getJsonValue("roleName")
  }

  // =========================================CREATE FORMS=================================================

  // Create score form.
  createScoreForm(): void {
    this.scoreForm = this.formBuilder.group({
      score: new FormControl(null, [Validators.required, Validators.min(0)]),
    })
  }

  // =============================================================FORMAT FUNCTIONS==========================================================================

  // Format fields of problem details.
  formatProblemFields(): void {

    // Set class for difficulty.
    if (this.answer.programmingQuestion.level == 1) {
      this.answer.programmingQuestion.levelName = "Easy"
      this.answer.programmingQuestion.levelClass = "easy"
    }
    if (this.answer.programmingQuestion.level == 2) {
      this.answer.programmingQuestion.levelName = "Medium"
      this.answer.programmingQuestion.levelClass = "medium"
    }
    if (this.answer.programmingQuestion.level == 3) {
      this.answer.programmingQuestion.levelName = "Hard"
      this.answer.programmingQuestion.levelClass = "hard"
    }

    // Give max limit to score in score form.
    this.scoreForm.get("score").setValue(this.answer.score)
    this.scoreForm.get("score").setValidators([Validators.required, Validators.min(0), Validators.max(this.answer.programmingQuestion.score)])
    this.utilService.updateValueAndValiditors(this.scoreForm)

    // If solution is viewed by talent for this programming question then disable score form.
    if (this.answer.programmingQuestion.solutonIsViewed) {
      this.scoreForm.disable()
    }

    // Set the solution language list to languages that the solutions of the problem have.
    this.solutionLanguageList = []
    if (this.answer.programmingQuestion.solutions?.length > 0) {
      for (let i = 0; i < this.answer.programmingQuestion.solutions.length; i++) {

        // Set the first solution as the selected programming question solution.
        if (i == 0) {
          this.selectedQuestionSolution = this.answer.programmingQuestion.solutions[i]
        }
        this.solutionLanguageList.push(this.answer.programmingQuestion.solutions[i].programmingLanguage)
      }
    }

    // If code editor answer and code editor solution instance is not created then create it first.
    if (!this.codeEditorAnswer || !this.codeEditorSolution) {
      this.codeEditorList.changes.subscribe((tempCodeEditorList: QueryList<CodeEditorComponent>) => {
        // If only one code editor exists then it is code editor answer. 
        if (tempCodeEditorList.length == 1) {
          this.codeEditorAnswer = tempCodeEditorList.first
        }

        // If two code editors exist then it is code editor answer and code editor solution. 
        if (tempCodeEditorList.length == 2) {
          this.codeEditorAnswer = tempCodeEditorList.last
          this.codeEditorSolution = tempCodeEditorList.first
        }

        // If code editor answer exists then send values to code editor answer.
        if (this.codeEditorAnswer) {
          this.sendValuesToCodeEditorAnswer()
        }

        // If code editor solution instance exists and solution is viewed then call its function.
        if (this.codeEditorSolution) {
          this.setProgrammingSolutionQuestion()
        }
      })
    }

    // If code editor answer instance exists or if answer exists then call its function.
    if (this.codeEditorAnswer) {
      this.sendValuesToCodeEditorAnswer()
    }

    // If nav bar exists then select the first tab as active tab.
    if (this.navBarElement) {
      this.navBarElement.select(this.solutionLanguageList[0].id)
    }
  }

  // =============================================================BUTTON CLICK FUNCTIONS==========================================================================

  // On clicking language tab.
  onLanguageTabClick(languageID: string): void {

    // If problem has solutions then get the solution for the language tab selected.
    if (this.answer.programmingQuestion.solutions?.length > 0) {
      for (let i = 0; i < this.answer.programmingQuestion.solutions.length; i++) {
        if (this.answer.programmingQuestion.solutions[i].programmingLanguage.id == languageID) {
          this.selectedQuestionSolution = this.answer.programmingQuestion.solutions[i]
          this.sendValuesToCodeEditorSolution()
          break
        }
      }
    }
  }

  // On clicking the submit button.
  validateScoreForm(): void {
    if (this.scoreForm.invalid) {
      this.scoreForm.markAllAsTouched()
      return
    }
    if (confirm("Are sure you want to submit the score ?")) {
      this.updateAnswerScoreAndIsCorrect()
    }
  }

  // =============================================================REDIRECT FUNCTIONS==========================================================================

  // On clicking back button.
  redirectToProblemOfTheDay(): void {
    if (this.roleName == this.role.ADMIN) {
      this.router.navigate(['/admin/coding-problems/answers'], {
      }).catch(err => {
        console.error(err)
      })
    } else {
      this.router.navigate(['/coding-problems/answers'], {
      }).catch(err => {
        console.error(err)
      })
    }
  }

  // =============================================================CRUD FUNCTIONS==========================================================================

  // Update score and is correct of programming question talent answer.
  updateAnswerScoreAndIsCorrect(): void {
    this.spinnerService.loadingMessage = "Submitting the score"


    this.answer.score = this.scoreForm.get('score').value
    let isCorrect: boolean = false
    if (this.answer.score == this.answer.programmingQuestion.score) {
      isCorrect = true
    }
    let answerScore: any = {
      id: this.answer.id,
      score: this.answer.score,
      isCorrect: isCorrect
    }
    this.answerService.updateAnswerScore(answerScore).subscribe((response: any) => {
      alert(response)
      this.getAnswerDetails()
    }, (error) => {
      console.error(error)
      if (error.error?.error) {
        alert(error.error?.error)
        return
      }
      alert(error.statusText)
    })
  }

  // =============================================================OTHER FUNCTIONS==========================================================================

  // // Receive values from child component.
  // receiveChildValues(tempLanguageAndCode: any) {

  //   // Give code to talent answer recieved from child component.
  //   this.talentAnswer = tempLanguageAndCode.code 

  //   // Get language id from langauge name received from child component.
  //   for (let i = 0; i < this.languageList.length; i++){
  //     if (tempLanguageAndCode.languageName == this.languageList[i].name){
  //       this.selecetedLanguage = this.languageList[i]
  //       break
  //     }
  //   }

  //   // Submit answer.
  //   this.getSolutionsDetails()
  // }

  // Send values to code editor answer after getting problem details.
  sendValuesToCodeEditorAnswer(): void {

    // If problem is answered then send language and code to child component.
    this.languageAndCodeAnswer = {
      languageName: this.answer.programmingLanguage.name,
      code: this.answer.answer,
      isSolution: true,
      isReadOnly: true
    }
    this.codeEditorAnswer.valuesFromParentChange(this.languageAndCodeAnswer)
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
  setProgrammingSolutionQuestion(): void {

    // Set the solution language list to languages that the solutions of the problem have.
    if (this.answer.programmingQuestion.solutions?.length > 0) {
      this.solutionLanguageList = []
      for (let i = 0; i < this.answer.programmingQuestion.solutions.length; i++) {

        // Set the fisrt solution as the selected programming question solution.
        if (i == 0) {
          this.selectedQuestionSolution = this.answer.programmingQuestion.solutions[i]
        }
        this.solutionLanguageList.push(this.answer.programmingQuestion.solutions[i].programmingLanguage)
      }

      // Send the values to code editor solution.
      if (this.selectedQuestionSolution) {
        this.sendValuesToCodeEditorSolution()
      }
    }
  }

  // =============================================================GET FUNCTIONS==========================================================================

  // Get programming question talent answer details.
  getAnswerDetails(): void {
    this.spinnerService.loadingMessage = "Getting Talent Answer Details"


    this.answerService.getAnswer(this.answerID).subscribe((response) => {
      this.answer = response
      this.formatProblemFields()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
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
