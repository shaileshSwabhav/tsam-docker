import { Component, OnInit } from '@angular/core';
import { Validators, FormBuilder, FormGroup } from '@angular/forms';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { DummyDataService } from 'src/app/service/dummydata/dummy-data.service';
declare var $: any;

@Component({
      selector: 'app-faculty-batch',
      templateUrl: './faculty-batch.component.html',
      styleUrls: ['./faculty-batch.component.css']
})
export class FacultyBatchComponent implements OnInit {

      student: any[];
      selectedStudent: any[];
      modelHeader: string;
      modelButton: string;
      courseform: FormGroup;
      isCourseUpdate: boolean;
      modelFunction: () => void;
      isCourseUpdateBody: boolean;
      isStudentBatchUpdate: boolean;
      technologylist: any[];
      selectedtechnology: any[];
      formHandler: (index: number) => void;
      resourseType: string;
      enquiry: any[];
      multipleSelect: boolean;
      multi: any[];
      view: any[] = [1030, 500];
      constructor(
            private formbuilder: FormBuilder,
            private util: UtilityService,
            private dummydata: DummyDataService,
      ) {
            this.selectedStudent = [];
            this.multipleSelect = false;
      }

      ngOnInit() {
            this.loadStudent();
            this.onCourseAddButtonClick()
            this.formHandler(0);
            this.isStudentBatchUpdate = false;
      }

      createFormForAddCourse(param: number) {
            this.courseform = this.formbuilder.group({
                  courseName: ['', Validators.required],
                  technology: ['', Validators.required],
                  courseType: ['', Validators.required],
                  eligibility: ['', Validators.required],
                  price: ['', Validators.required],
                  discount: ['', Validators.required],
                  courseDuration: ['', Validators.required],
                  session: ['', Validators.required],
                  salesPerson: ['', Validators.required]
            });
      }

      // Form Action.
      formAction() {
            this.modelFunction();
      }


      // Add Course.
      addCourse() {
            this.util.info("Add Function Call......")
      }


      // Update Course.
      updateCourse() {
            this.util.info("Update Function Call......")
      }

      deleteCourse() {
            this.util.info("Delete Function Call......")
      }

      // On Course Add Button Click
      onCourseAddButtonClick() {
            this.UpdateCommonVariable(this.createFormForAddCourse, this.addCourse, "Add New Batch", "Add Batch", false, true);
      }


      // On Update Course Button Click.
      onCourseUpdateButtonClick() {
            this.UpdateCommonVariable(this.updateCourseForm, this.updateCourse, "Update Batch", "Update Batch", true, false);
      }


      // On Serial Number Change On Model.
      onSerialNumberChange(serialnumber: number) {
            serialnumber = Number(serialnumber);
            this.formHandler(serialnumber);
      }

      onStudentBatchSelect(batchnumber: string) {
            this.isStudentBatchUpdate = true;
            this.util.info(batchnumber);
      }


      // Update Course Form.
      updateCourseForm(param: number) {
            this.courseform.patchValue({});
            this.isCourseUpdateBody = true;
      }


      // Change
      resourceTypeChange(value) {
            this.resourseType = value;
      }


      // 
      UpdateCommonVariable(formhandler, actionhandler, headername: any, buttonname: any, formupdate: any, formupdatebody: any) {
            this.formHandler = formhandler;
            this.modelFunction = actionhandler;
            this.modelHeader = headername;
            this.modelButton = buttonname;
            this.isCourseUpdate = formupdate;
            this.isCourseUpdateBody = formupdatebody
      }

      // Remove Student.
      removeStudent(index: number) {
            this.util.info("Called", typeof index, index);
            this.selectedStudent.splice(index, 1);
            this.util.info(this.selectedStudent);
      }

      modelToggle(closedmodelname: string, openmodelname) {
            this.modelHide(closedmodelname)
            $('#' + openmodelname).modal('toggle')
      }

      modelHide(name: string) {
            $('#' + name).modal('hide')
      }


      loadStudent() {
            this.student = [
                  {
                        name: 'Akash Rai',
                        rollNo: 12,
                  },
                  {
                        name: 'Abhishek Singh',
                        rollNo: 121,
                  },
                  {
                        name: 'Sonam Singh',
                        rollNo: 130,
                  },
                  {
                        name: 'Akash Rai',
                        rollNo: 12,
                  },
                  {
                        name: 'Abhishek Singh',
                        rollNo: 121,
                  },
                  {
                        name: 'Sonam Singh',
                        rollNo: 130,
                  },
                  {
                        name: 'Akash Rai',
                        rollNo: 12,
                  },
                  {
                        name: 'Abhishek Singh',
                        rollNo: 121,
                  },
                  {
                        name: 'Sonam Singh',
                        rollNo: 130,
                  }
            ];

            this.selectedStudent = [
                  {
                        name: 'Akash Rai',
                        rollNo: 12,
                  },
                  {
                        name: 'Abhishek Singh',
                        rollNo: 121,
                  },
                  {
                        name: 'Sonam Singh',
                        rollNo: 130,
                  }
            ];

            this.multi = [
                  {
                        "name": "Student 1",
                        "series": [
                              {
                                    "name": "Introduction",
                                    "value": 11
                              },
                              {
                                    "name": "Variables, Types and Constants",
                                    "value": 12
                              },
                              {
                                    "name": "Functions and Packages",
                                    "value": 13
                              },
                              {
                                    "name": "Conditional Statements and Loops",
                                    "value": 14
                              },
                              {
                                    "name": "Arrays, Slices and Variadic Functions",
                                    "value": 15
                              },
                              {
                                    "name": "More types",
                                    "value": 16
                              },
                              {
                                    "name": "Pointers, Structs and Methods",
                                    "value": 17
                              },
                              {
                                    "name": "Interfaces",
                                    "value": 18
                              },
                              {
                                    "name": "Concurrency",
                                    "value": 19
                              }
                        ]
                  },

                  {
                        "name": "Student 2",
                        "series": [
                              {
                                    "name": "Introduction",
                                    "value": 8
                              },
                              {
                                    "name": "Variables, Types and Constants",
                                    "value": 9
                              },
                              {
                                    "name": "Functions and Packages",
                                    "value": 10
                              },
                              {
                                    "name": "Conditional Statements and Loops",
                                    "value": 11
                              },
                              {
                                    "name": "Arrays, Slices and Variadic Functions",
                                    "value": 12
                              },
                              {
                                    "name": "More types",
                                    "value": 13
                              },
                              {
                                    "name": "Pointers, Structs and Methods",
                                    "value": 14
                              },
                              {
                                    "name": "Interfaces",
                                    "value": 15
                              },
                              {
                                    "name": "Concurrency",
                                    "value": 16
                              }
                        ]
                  },

                  {
                        "name": "student 3",
                        "series": [
                              {
                                    "name": "Introduction",
                                    "value": 20
                              },
                              {
                                    "name": "Variables, Types and Constants",
                                    "value": 15
                              },
                              {
                                    "name": "Functions and Packages",
                                    "value": 17
                              },
                              {
                                    "name": "Conditional Statements and Loops",
                                    "value": 17
                              },
                              {
                                    "name": "Arrays, Slices and Variadic Functions",
                                    "value": 15
                              },
                              {
                                    "name": "More types",
                                    "value": 17
                              },
                              {
                                    "name": "Pointers, Structs and Methods",
                                    "value": 17
                              },
                              {
                                    "name": "Interfaces",
                                    "value": 17
                              },
                              {
                                    "name": "Concurrency",
                                    "value": 17
                              }
                        ]
                  },
                  {
                        "name": "Student 4",
                        "series": [
                              {
                                    "name": "Introduction",
                                    "value": 14
                              },
                              {
                                    "name": "Variables, Types and Constants",
                                    "value": 20
                              },
                              {
                                    "name": "Functions and Packages",
                                    "value": 3
                              },
                              {
                                    "name": "Conditional Statements and Loops",
                                    "value": 1
                              },
                              {
                                    "name": "Arrays, Slices and Variadic Functions",
                                    "value": 5
                              },
                              {
                                    "name": "More types",
                                    "value": 7
                              },
                              {
                                    "name": "Pointers, Structs and Methods",
                                    "value": 7
                              },
                              {
                                    "name": "Interfaces",
                                    "value": 17
                              },
                              {
                                    "name": "Concurrency",
                                    "value": 17
                              }
                        ]
                  }
            ];

            // this.getTechnology();
            // this.getEnquiryasStudentList();
      }

      // getEnquiryasStudentList() {
      //       this.dummydata.getEnquiries().then((respond: any[]) => {
      //             this.enquiry = respond;
      //          
      //       })
      // }

      printSelected() {
            this.util.info(this.selectedStudent)
      }

      // getTechnology() {
      //       this.dummydata.getTechnology().then((respond: any[]) => {
      //             this.technologylist = respond;
      //          
      //       });
      // }


      // options
      legend: boolean = true;
      showLabels: boolean = true;
      xAxis: boolean = true;
      yAxis: boolean = true;
      showYAxisLabel: boolean = true;
      showXAxisLabel: boolean = true;
      xAxisLabel: string = 'Session';
      yAxisLabel: string = 'Marks';
      timeline: boolean = true;

      colorScheme = {
            domain: ['#5AA454', '#E44D25', '#CFC0BB', '#7aa3e5', '#a8385d', '#aae3f5']
      };

}