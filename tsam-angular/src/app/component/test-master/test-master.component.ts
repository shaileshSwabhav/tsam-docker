import { Component, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, Validators, FormArray, FormControl } from '@angular/forms';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { TalentService } from 'src/app/service/talent/talent.service';

@Component({
      selector: 'app-test-master',
      templateUrl: './test-master.component.html',
      styleUrls: ['./test-master.component.css']
})
export class TestMasterComponent implements OnInit {

      fieldlabel: string;
      flags: any[2];
      subjects: any[];
      questionlist: any[];
      modelHeader: string;
      questionTypes: any[];
      modelAction: () => void;
      questionform: FormGroup
      modelActionType: string;
      questiondifficulty: any[];
      formaction: (questionid?) => void;
      constructor(
            private formbuilder: FormBuilder,
            private util: UtilityService,
      ) {
            this.onQuestionAddButtonClick();
            this.fieldlabel = 'Option';
      }

      ngOnInit() {
            this.loadData();
            this.formaction();
            // this.createQuestionForm();
      }

      // Create Question Form.
      createQuestionForm() {
            this.questionform = this.formbuilder.group({
                  questionID: ['', Validators.required],
                  questionType: ['', Validators.required],
                  difficulty: ['', Validators.required],
                  subject: ['', Validators.required],
                  question: ['', Validators.required],
                  options: this.formbuilder.array([])
            });
            this.addNewOption();
      }

      // Update Question Form with Value.
      updateQuestionForm(index: number) {
            this.questionform.patchValue(this.questionlist[index]);
            this.updateOption(index)
      }

      updateOption(index) {
            let length = this.questionlist[index].options.length;
            for (let i = 0; i < length; i++) {
                  this.optionControl.at(i).patchValue(this.questionlist[index].options[i]);
                  if (length - 1 > i && this.optionControl.length < length) {
                        this.addNewOption();
                  }
            }
            let diff = this.optionControl.length;
            for (let i = length; i < diff; i++) {
                  this.optionControl.removeAt(i);
            }
      }

      // updateTestCase(index) {
      //       let length = this.questionlist[index].options.length;
      //       for (let i = 0; i < length; i++) {
      //             this.optionControl.at(i).patchValue(this.questionlist[index].options[i]);
      //             if (length - 1 > i && this.optionControl.length < length) {
      //                   this.addNewOption();
      //             }
      //       }
      //       let diff = this.optionControl.length;
      //       for (let i = length; i < diff; i++) {
      //             this.optionControl.removeAt(i);
      //       }
      // }

      isQuestionTypeProgram(): boolean {
            if (this.formControl.questionType.value.type == "Program") {
                  this.fieldlabel = "Test Case"
                  return true;
            }
            this.fieldlabel = "Option"
            return false;
      }

      // Add New Option To Options Field in Form.
      addNewOption() {
            this.optionControl.push(this.formbuilder.group({
                  option: ['', Validators.required],
                  status: ['', Validators.required]
            }))
      }

      //// Add New Testcase To Options Field in Form when Question type is Program..
      addNewTestCase() {
            this.optionControl.push(this.formbuilder.group({
                  testcase: ['', Validators.required]
            }))
      }

      get optionControl() {
            return this.questionform.get('options') as FormArray
      }

      get testCaseControl() {
            return this.questionform.get('testcase') as FormArray
      }


      get formControl() {
            return this.questionform.controls;
      }

      // Question Add To Database.
      addQuestion() {
            this.util.info(this.questionform.value);
      }

      formAction() {
            this.modelAction();
      }

      // Question Update
      updateQuestion() {
            this.util.info(this.questionform.value);
      }

      // Get Question index from question ID.
      getQuestionByQuestionID(questionid: string) {
            for (let index = 0; index < this.questionlist.length; index++) {
                  if (this.questionlist[index].questionID == questionid) {
                        return index;
                  }
            }
      }

      onQuestionAddButtonClick() {
            this.updateFlags(this.createQuestionForm, 'Add New Question', 'Add Question', this.addQuestion);
            this.formaction();
      }

      onQuestionEditButtonClick(questionid: string) {
            let index = this.getQuestionByQuestionID(questionid)
            this.updateFlags(this.updateQuestionForm, 'Update Question', 'Update Question', this.updateQuestion);
            this.formaction(index);
      }

      // Compare c1 and c2
      compareFn(c1, c2): boolean {
            return c1 && c2 ? c1.id == c2.id : c1 == c2;
      }


      // Load All Data.
      loadData() {
            this.questionlist = [
                  {
                        questionID: 'QN001',
                        question: 'What is Serialization ?',
                        options: [
                              {
                                    option: 'Option1',
                                    status: {
                                          id: '2',
                                          status: false
                                    }
                              },
                              {
                                    option: 'Option2',
                                    status: {
                                          id: '2',
                                          status: false
                                    }
                              },
                              {
                                    option: 'Option3',
                                    status: {
                                          id: '2',
                                          status: false
                                    }
                              },
                              {
                                    option: 'Option5',
                                    status: {
                                          id: '1',
                                          status: true
                                    }
                              }
                        ],
                        questionType: {
                              id: '1',
                              type: 'MCQ'
                        },
                        difficulty: {
                              id: '1',
                              level: 'medium'
                        },
                        subject: {
                              id: '2',
                              subject: "Java"
                        }
                  },
                  {
                        questionID: 'QN002',
                        question: 'What is Golang Struct ?',
                        options: [
                              {
                                    option: 'Option',
                                    status: {
                                          id: '2',
                                          status: false
                                    }
                              },
                              {
                                    option: 'Option',
                                    status: {
                                          id: '2',
                                          status: false
                                    }
                              },
                              {
                                    option: 'Option',
                                    status: {
                                          id: '2',
                                          status: false
                                    }
                              },
                              {
                                    option: 'Option',
                                    status: {
                                          id: '1',
                                          status: true
                                    }
                              }
                        ],
                        questionType: {
                              id: '1',
                              type: 'MCQ'
                        },
                        difficulty: {
                              id: '1',
                              level: 'medium'
                        },
                        subject: {
                              id: '1',
                              subject: "Golang"
                        }
                  }
            ];

            this.questionTypes = [
                  {
                        id: '1',
                        type: 'MCQ'
                  },
                  {
                        id: '2',
                        type: 'Program'
                  }
            ];

            this.questiondifficulty = [
                  {
                        id: '1',
                        level: 'Easy'
                  },
                  {
                        id: '2',
                        level: 'Medium'
                  },
                  {
                        id: '3',
                        level: 'Difficult'
                  }
            ];

            this.subjects = [
                  {
                        id: '1',
                        subject: 'Golang'
                  },
                  {
                        id: '2',
                        subject: 'Java'
                  },
                  {
                        id: '3',
                        subject: 'Javascript'
                  },
                  {
                        id: '4',
                        subject: 'Julia'
                  }
            ];

            this.flags = [
                  {
                        id: '1',
                        status: true
                  },
                  {
                        id: '2',
                        status: false
                  }
            ];
      }

      // Update Flags.
      updateFlags(formaction, modelheader, modelactiontype, modelaction) {
            this.formaction = formaction;
            this.modelAction = modelaction;
            this.modelActionType = modelactiontype;
            this.modelHeader = modelheader;
      }

}
