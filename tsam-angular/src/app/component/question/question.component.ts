import { Component, OnInit } from '@angular/core';
import { FormArray, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { GeneralService } from 'src/app/service/general/general.service';
import { OptionService } from 'src/app/service/option/option.service';
import { QuestionService } from 'src/app/service/question/question.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-question',
  templateUrl: './question.component.html',
  styleUrls: ['./question.component.css']
})
export class QuestionComponent implements OnInit {

  questions: any[];
  options: any[];
  totalQues: number;
  limit: number;
  offset: number;
  currentpage: number;
  difficultyLevels: any[];
  technologies: any[];

  Id: string;

  updateHit: boolean = false;

  selectedId: string;
  selectedOptionId: string;
  selectedQuestionId: string;

  searched: boolean = false;

  modelButton: string;
  modelHeader: string;
  modelAction: () => void;
  formAction: (param?) => void;
  isQuesUpadteClick: boolean;


  questionForm: FormGroup;
  optionForm: FormGroup;
  searchForm: FormGroup;

  constructor(
    private quesService: QuestionService,
    private util: UtilityService,
    private optionService: OptionService,
    private generalService: GeneralService,
    private formBuilder: FormBuilder,
    private spinnerService: SpinnerService,
    private optService: OptionService
  ) {
    this.limit = 5;
    this.offset = 0;

  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {


    this.onAddQuesClick();
    this.createQuestionForm();
    this.createSearchForm();
    this.createOptionForm();
    this.formAction();
    this.loadData();

  }



  loadData() {
    this.getAllTechnologies();
    this.getDiffLevels();
    this.getQuestions();
  }


  // Pagination Controls

  changePage($event) {

    // $event will be the page number & offset will be 1 less than it.
    this.offset = $event - 1
    this.currentpage = $event;
    this.getQuestions();
  }

  // On Limit Change
  limitChange() {
    this.offset = 1;
    this.getQuestions();
  }

  // On Page Number Change.
  OffsetChange(param: number) {
    this.offset = param;
    this.getQuestions();
  }



  createQuestionForm() {
    this.questionForm = this.formBuilder.group(
      {
        id: [],
        question: [null, Validators.required],
        subject: [null],
        difficulty: [null],
        options: this.formBuilder.array([])
      }
    )
  }

  createSearchForm() {
    this.searchForm = this.formBuilder.group(
      {
        id: [],
        question: [null],
        subject: [null],
        difficulty: [null]
      }
    )
  }

  createOptionForm() {
    this.optionForm = this.formBuilder.group(
      {
        id: null,
        option: [null, Validators.required],
        status: [null, Validators.required]
      }
    )
  }

  //set new options
  get newOptions() {
    return this.questionForm.get('options') as FormArray;
  }


  //add option to new ques
  addOption() {
    this.newOptions.push(this.formBuilder.group(
      {
        id: [],
        option: [null],
        status: [null]
      }
    ));
  }



  //add a new question
  addQuestion() {


    let question = this.questionForm.value;
    console.log(question)
    this.quesService.addQuestion(question).subscribe
      (
        data => {
          alert("Question Added Successfully");
          this.reset();
          this.getQuestions();
        },
        err => {
          console.log(err);

        }
      )

  }




  ///reset Form
  reset() {
    this.questionForm.reset();
  }


  // Update Model Param.
  updateCommonParam(modelbutton, modelheader, modelAction, formAction, isQuesUpadteClick) {
    this.modelAction = modelAction;
    this.modelButton = modelbutton;
    this.modelHeader = modelheader;
    this.isQuesUpadteClick = isQuesUpadteClick;
    this.formAction = formAction;
  }




  //Add new question click
  onAddQuesClick() {
    this.createQuestionForm();
    this.updateCommonParam('Add Question', 'Add New Question',
      this.addQuestion, this.createQuestionForm, false);
  }








  // On update Button Click.
  onUpdateQuestionClick(updateID: string) {

    this.updateCommonParam('Update Question', 'Update Question',
      this.updateQues,
      this.updateQuestionForm, true);
    this.createQuestionForm();
    this.selectedId = updateID;
    let index = updateID;
    // this.createOptionForm();

    this.quesService.getQuestionByID(index).subscribe
      (
        data => {
          this.updateQuestionForm(data);
        }
      )
    // console.log(selectedData)
    // this.updateRoleForm(selectedData);
  }



  //update Form 
  updateQuestionForm(data) {

    data.options.forEach(t => {
      this.addOption();

    })
    console.log(this.questionForm)
    this.questionForm.patchValue(
      {
        question: data.question,
        difficulty: data.difficulty,
        subject: data.subject,
        options: data.options,
        id: this.selectedId

      }
    )



  }


  //update Question using Id
  updateQues() {


    let ques = this.questionForm.value;
    this.quesService.updateQuestionByID(this.selectedId, ques).
      subscribe(
        data => {
          alert("Updated Successfully");

          this.reset();
          this.getQuestions();
        }
      )


  }



  //delete

  assignID(id: string) {
    this.selectedId = id;
  }



  //delete
  deleteQuestion() {

    this.quesService.deleteQuestionByID(this.selectedId).subscribe
      (
        data => {
          alert("Question deleted successfully");
          this.getQuestions();
        }, err => {
          console.log(err.error.error);

        }
      )
  }



  //get all questions
  getQuestions() {

    this.quesService.getAllQuestions(this.limit, this.offset)
      .subscribe(
        data => {
          this.questions = [];
          this.questions = this.questions.concat(data.body)
          this.totalQues = parseInt(data.headers.get('X-Total-Count'));


        }, err => {
          console.log(this.util.getErrorString(err))

        })

  }


  //get options using question ID

  getOptionsforID(id: string) {
    this.options = [];
    this.selectedQuestionId = id;

    this.optionService.getOptionByQuestion(id).subscribe
      (data => {
        this.options = data;


      },
        error => {
          console.log(error);

        })
  }




  //view question modal

  onViewQuestion(id: string) {

    this.Id = id;
    this.createQuestionForm();
    this.quesService.getQuestionByID(id).subscribe
      (
        data => {
          this.updateQuestionForm(data);

        },
        err => {
          console.log(err.error.error);

        }
      )
  }


  //create option Form

  //updateOption form

  updateOptionForm(data) {
    this.optionForm.patchValue
      (
        {
          option: data.option,
          status: data.status,
          id: this.selectedOptionId
        }
      )
  }


  //on update click
  updateOptionClick(optionId: string) {
    this.selectedOptionId = optionId;
    let id = this.selectedOptionId;
    this.createOptionForm();
    this.optService.getOptionByID(id).subscribe
      (
        data => {
          console.log(data)
          this.updateOptionForm(data);
        }
      )
  }




  //update option
  updateOption() {

    let value = this.optionForm.value;
    let index = this.selectedOptionId;
    this.optService.updateOptionByID(index, value).subscribe
      (
        data => {
          alert("Option updated successfully");
          this.optionForm.reset();
          this.getOptionsforID(this.selectedQuestionId);
        }
      )
  }






  //assign option id

  assignOptionId(optionId: string) {
    this.selectedOptionId = optionId;
  }



  //delete Option

  deleteOption() {

    let id = this.selectedOptionId;
    this.optService.deleteOptionByID(id).subscribe
      (data => {
        console.log(data)
        alert("option deleted successfully");
        this.getOptionsforID(this.selectedId);
      }, err => {
        console.log(err.error.error)
      })
  }


  //search

  search() {

    this.limit = 5;
    this.offset = 0;
    let data = this.searchForm.value;
    this.quesService.searchQuestion(data, this.limit, this.offset).
      subscribe(
        res => {
          console.log(res)
          this.questions = [];
          this.questions = this.questions.concat(res.body);
          this.totalQues = parseInt(res.headers.get('X-Total-Count'));
          this.searchForm.reset();
          this.searched = true;

          console.log(this.questions)
        },
        error => console.log(error)
      )
  }


  resetSearch() {

    this.searchForm.reset();
    this.searched = false;
    this.getQuestions();

  }

  //all Data part

  getAllTechnologies() {
    this.generalService.getTechnologies().subscribe((respond: any[]) => {
      this.technologies = respond
    }
      , (err) => {
        console.log(err)
      })

  }



  //get data from general
  getDiffLevels() {
    this.generalService.getGeneralTypeByType("question_difficulty").subscribe((respond: any[]) => {
      this.difficultyLevels = respond
    }, (err) => {
      console.log(err)
    })
  }


}
