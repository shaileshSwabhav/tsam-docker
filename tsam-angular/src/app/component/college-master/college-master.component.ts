import { Component, OnInit, ViewChild } from '@angular/core'
import { FormGroup, FormBuilder, Validators, FormArray, FormControl } from '@angular/forms'
import { IState, ICountry, CollegeService, ICollege, IUniversity, ICollegeBranch } from 'src/app/service/college/college.service'
import { GeneralService } from 'src/app/service/general/general.service'
import { Constant, UrlConstant } from 'src/app/service/constant'
import { SpinnerService } from 'src/app/service/spinner/spinner.service';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap'
import { Router } from '@angular/router'
import { IPermission } from 'src/app/service/menu/menu.service'
import { UtilityService } from 'src/app/service/utility/utility.service'
import { LocalService } from 'src/app/service/storage/local.service'
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service'
import { UrlService } from 'src/app/service/url.service'

@Component({
      selector: 'app-college-master',
      templateUrl: './college-master.component.html',
      styleUrls: ['./college-master.component.css']
})
export class CollegeMasterComponent implements OnInit {

      // College.
      college: ICollege
      collegeList: any[]
      collegeForm: FormGroup
      showBranchForm: boolean
      showCollegeControls: boolean

      // Additional Components.
      salesPeople: any[]
      // stateList: IState[]
      stateList: IState[][]
      countryList: ICountry[]
      ratingList: number[]
      allUniversities: any[]
      // countryIDList conatins the list of countryIDs whose states data is present in stateList.
      countryIDSet: Set<string>

      // Search.
      searchCollegeForm: FormGroup
      searchFormValue: any
      searched: boolean

      // Pagination.
      limit: number
      offset: number
      currentPage: number
      totalColleges: number
      paginationString: string

      // Modal.
      selectedCollegeID: string
      modalHeader: string
      modalButton: string
      modalAction: () => void
      formHandler: () => void
      modalRef: any
      isOperationUpdate: boolean
      isUploadedSuccessfully: boolean

      //spinner


      // access
      permission: IPermission
      isViewMode: boolean

      // constants
      readonly NIL_UUID: string = this.constant.NIL_UUID
      readonly COLLEGE_EXCEL_DEMO_LINK = this.urlConstant.COLLEGE_DEMO

      @ViewChild("deleteConfirmationModal") deleteConfirmationModal: any

      constructor(
            private formBuilder: FormBuilder,
            private collegeService: CollegeService,
            private generalService: GeneralService,
            private constant: Constant,
            private urlConstant: UrlConstant,
            private spinnerService: SpinnerService,
            private modalService: NgbModal,
            private router: Router,
            private utilService: UtilityService,
            private localService: LocalService,
            private fileOps: FileOperationService,
            private urlService: UrlService,
      ) {
            this.limit = 5
            this.offset = 0
            this.spinnerService.loadingMessage = "Getting colleges"
            this.isOperationUpdate = false
            this.searched = false
            this.isViewMode = false
            this.countryIDSet = new Set()
            // this.stateList = [] as IState[]
            // this.NIL_UUID = this.constant.NIL_UUID
            this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.COLLEGE_MASTER)
      }


      get ongoingOperations() {
            return this.spinnerService.ongoingOperations
      }

      ngOnInit() {

            this.getAllComponents()
            this.createAddCollegeForm()
            this.createSearchCollegeForm()
            this.getAllColleges()
      }

      // Handles pagination
      changePage($event: any): void {
            // $event will be the page number & offset will be 1 less than it.

            this.offset = $event - 1
            this.currentPage = $event
            this.getAllColleges()
      }

      // gets all required components like all states,countries etc
      getAllComponents(): void {
            this.getAllCountries()
            this.getAllRatings()
            this.getAllSalesPeople()
            this.getAllUniversities()
      }

      // getStateList(countryID: string, shouldClear?: boolean): void {
      //       if (this.countryIDSet.has(countryID)) {
      //             return
      //       }
      //       this.countryIDSet.add(countryID)

      //       if (countryID) {
      //             this.generalService.getStatesByCountryID(countryID).subscribe((data: any) => {
      //                   this.stateList = this.stateList.concat(data);
      //                   console.log("State List :", this.stateList);
      //             }, (err) => {
      //                   console.log(err)
      //             })
      //       }
      // }

      getStateList(countryID: string, index: number): void {
            // if (this.countryIDSet.has(countryID)) {
            //       return
            // }
            // this.countryIDSet.add(countryID)
            this.stateList[index] = []
            this.collegeBranch.at(index).get('address').get('state')?.reset()
            this.collegeBranch.at(index).get('address').get('state')?.disable()

            if (countryID) {
                  this.generalService.getStatesByCountryID(countryID).subscribe((data: any) => {
                        // this.stateList = this.stateList.concat(data);
                        this.stateList[index] = data
                        console.log("State List :", this.stateList)
                        if (this.stateList[index].length > 0 && !this.isViewMode) {
                              this.collegeBranch.at(index).get('address').get('state')?.enable()
                        }
                  }, (err) => {
                        console.log(err)
                  })
            }
      }

      getAllRatings(): void {
            this.ratingList = []
            this.generalService.getGeneralTypeByType("college_rating").subscribe((data: any[]) => {
                  for (let rating of data) {
                        this.ratingList.push(+rating.value)
                  }
                  this.ratingList = this.ratingList.sort()
            }, (err) => {
                  console.log(err)
            })
      }

      getAllCountries(): void {
            this.generalService.getCountries().subscribe((data: ICountry[]) => {
                  this.countryList = data.sort()
            }, (err) => {
                  console.log(err)
            })
      }

      getAllUniversities(): void {
            this.generalService.getAllUniversities().subscribe((data: any) => {
                  this.allUniversities = data
            }, (err) => {
                  console.log(err)
            })
      }

      getAllSalesPeople(): void {
            this.generalService.getSalesPersonList().subscribe((data: any) => {
                  this.salesPeople = data.body
            }, (err) => {
                  console.log(err)
            })
      }

      showSpecificStates(branch: any, state: IState): boolean {
            console.log(branch.get('address').get('country').value);
            console.log(state);

            if (!branch.get('address').get('country').value) {
                  return false
            }
            if (branch.get('address').get('country').value.id === state.countryID) {
                  return true
            }
            return false
      }

      compareFn(c1: any, c2: any): boolean {
            return c1 && c2 ? c1.id === c2.id : c1 === c2
      }

      // Create search form.
      createSearchCollegeForm(): void {
            this.searchCollegeForm = this.formBuilder.group({
                  collegeName: new FormControl(null, [Validators.maxLength(100), Validators.pattern("^[a-zA-Z]+([a-zA-Z. ]?)+")]),
            })
      }

      // Create college form.
      createAddCollegeForm(): void {
            this.createUpdateCollegeForm()
            this.stateList = []
            this.addNewBranchInForm()
      }

      createUpdateCollegeForm(): void {
            this.collegeForm = this.formBuilder.group({
                  id: new FormControl(null),
                  code: new FormControl(null),
                  collegeName: new FormControl(null, [Validators.required, Validators.maxLength(100), Validators.pattern("^[a-zA-Z]+([a-zA-Z. ]?)+")]),
                  chairmanName: new FormControl(null, [Validators.maxLength(70), Validators.pattern("^[a-zA-Z ]*$")]),
                  chairmanContact: new FormControl(null, [Validators.pattern(/^[6789]\d{9}$/)]),
                  collegeBranches: this.formBuilder.array([])
            })
      }

      // Gets collegeBranches formBuilder array.
      get collegeBranch() {
            return this.collegeForm.get('collegeBranches') as FormArray
      }

      // Add New Branch.
      addNewBranchInForm(): void {
            this.collegeBranch.push(this.formBuilder.group({
                  id: new FormControl(null),
                  collegeID: new FormControl(null),
                  university: new FormControl(null, [Validators.required]),
                  branchName: new FormControl(null, [Validators.required, Validators.maxLength(150)]),
                  code: new FormControl(null),
                  salesPerson: new FormControl(null),
                  tpoName: new FormControl(null, [Validators.maxLength(70), Validators.pattern("^[a-zA-Z ]*$")]),
                  tpoContact: new FormControl(null, [Validators.pattern(/^[6789]\d{9}$/)]),
                  tpoAlternateContact: new FormControl(null, [Validators.pattern(/^[6789]\d{9}$/)]),
                  tpoEmail: new FormControl(null, [Validators.pattern(/^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$/)]),
                  collegeRating: new FormControl(null),
                  allIndiaRanking: new FormControl(null, [Validators.max(100000), Validators.pattern("^[1-9][0-9]*$")]),
                  email: new FormControl(null, [Validators.pattern(/^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$/)]),
                  address: this.formBuilder.group({
                        state: new FormControl(null, [Validators.required]),
                        country: new FormControl(null, [Validators.required]),
                        address: new FormControl(null, [Validators.required, Validators.maxLength(100)]),
                        city: new FormControl(null, [Validators.required, Validators.maxLength(50), Validators.pattern(/^[a-zA-Z]+([a-zA-Z ]?)+$/)]),
                        pinCode: new FormControl(null, [Validators.required, Validators.pattern(/^[1-9][0-9]{5}$/)])
                  }),
            }))
            // console.log(this.stateListTest)
            this.stateList.push([])
            this.collegeBranch.at(this.collegeBranch.controls.length - 1).get('address')?.get('state')?.disable()
      }

      // Delete College Branch from college Form.
      deleteBranchInForm(index: number): void {
            if (!confirm("This action will delete the branch.\nAre you sure?")) {
                  return
            }
            this.collegeForm.markAsDirty()
            this.collegeBranch.removeAt(index)
            this.stateList.splice(index, 1)
      }

      onAddNewCollegeButtonClick(collegeDetailModal: any): void {
            this.isOperationUpdate = false
            this.isViewMode = false

            this.showBranchForm = true
            this.showCollegeControls = true
            this.opencollegeDetailModal(collegeDetailModal)
            this.setVariableForAddCollege()
            this.formHandler()
      }

      onAddNewCollegeBranchButtonClick(collegeID: string, collegeDetailModal: any): void {
            this.showBranchForm = true
            this.isViewMode = false
            this.showCollegeControls = false
            this.opencollegeDetailModal(collegeDetailModal)
            this.setVariableForAddCollegeBranch()
            this.formHandler()
            this.collegeForm.removeControl('collegeName')
            this.selectedCollegeID = collegeID
      }

      onCollegeUpdateClick() {
            this.isOperationUpdate = true
            this.isViewMode = false;
            this.collegeForm.enable()
            this.setVariableForUpdateCollege()
      }

      onCollegeDeleteClick(collegeID: string): void {
            this.selectedCollegeID = collegeID
            this.opencollegeDetailModal(this.deleteConfirmationModal, "md")
      }

      deleteCollege(): void {
            this.modalRef.close()
            this.spinnerService.loadingMessage = "Deleting college"

            this.collegeService.deleteCollege(this.selectedCollegeID).subscribe((data: any) => {
                  this.getAllColleges()
                  alert(data)
            }, (error) => {
                  console.log(error)

                  alert(error.error)
            })
      }

      // Fields may be added in future.
      searchColleges(): void {
            console.log("In search");

            this.searchFormValue = { ...this.searchCollegeForm.value }
            let flag: boolean = true
            for (let field in this.searchFormValue) {
                  if (!this.searchFormValue[field]) {
                        delete this.searchFormValue[field]
                  } else {
                        this.searched = true
                        flag = false
                  }
            }
            // No API call on empty search.
            if (flag) {
                  return
            }
            this.spinnerService.loadingMessage = "Searching colleges"
            this.changePage(1)
      }


      setPaginationString() {
            this.paginationString = ''
            let start: number = this.limit * this.offset + 1
            let end: number = +this.limit + this.limit * this.offset
            if (this.totalColleges < end) {
                  end = this.totalColleges
            }
            if (this.totalColleges == 0) {
                  this.paginationString = ''
                  return
            }
            this.paginationString = `${start} - ${end}`
      }


      getAllColleges(): void {
            this.spinnerService.loadingMessage = "Getting Colleges"
            this.collegeService.getAllColleges(this.limit, this.offset, this.searchFormValue).subscribe(data => {

                  this.collegeList = data.body
                  this.totalColleges = parseInt(data.headers.get('X-Total-Count'))
                  this.setPaginationString()
            }, (error) => {
                  this.totalColleges = 0
                  console.log(error)

                  if (error.statusText.includes('Unknown')) {
                        alert("No connection to server. Check internet.")
                  }
            })
      }

      getBranchesByCollege(collegeID: string): void {
            console.log("college id" + collegeID)
            this.urlService.setUrlFrom(this.constructor.name, this.router.url)
            this.router.navigate([this.urlConstant.COLLEGE_BRANCH], {
                  queryParams: {
                        "collegeID": collegeID
                  }
            }).catch(err => {
                  console.log(err)

            })
      }

      resetSearchAndGetAll(): void {
            this.spinnerService.loadingMessage = "Getting colleges"

            this.searchCollegeForm.reset()
            this.searchFormValue = null
            this.changePage(1)
            this.getAllColleges()
            this.searched = false
      }

      // Add College.
      addCollege(): void {
            this.spinnerService.loadingMessage = "Adding college"

            this.collegeService.addCollege(this.collegeForm.value).subscribe((data: string) => {
                  this.modalRef.close()
                  this.getAllColleges()
                  this.collegeForm.reset()
                  alert("college added with id:" + data)
            }, (error) => {
                  console.log(error)

                  if (error.error) {
                        if (error.error.error) {
                              alert(error.error.error)
                              return
                        }
                        alert(error.error)
                        return
                  }
                  alert(error.statusText)
            })
      }

      // Update College.
      updateCollege(): void {
            this.spinnerService.loadingMessage = "Updating college"

            this.collegeService.updateCollege(this.collegeForm.value).subscribe((data: string) => {
                  this.modalRef.close()
                  this.getAllColleges()
                  this.collegeForm.reset()
                  alert(data)
            }, (error) => {
                  console.log(error)

                  if (error.error) {
                        alert(error.error)
                        return
                  }
                  alert(error.statusText)
            })
      }

      onViewCollegeClick(college: ICollege): void {
            this.modalHeader = "View College"
            this.isViewMode = true
            this.showBranchForm = false
            this.showCollegeControls = true
            this.college = college
            this.updateForm()
            this.collegeForm.disable()
      }

      // ===========================================START TEST=====================================================
      excelUploadedColleges: ICollegeExcel[]
      reformedUploadedColleges: ICollege[]

      // Use service.
      onFileChange(event: any, excelFile: any): void {
            const files = event.target.files;
            if (files.length === 0) {
                  return
            }
            if (files.length !== 1) {
                  alert('only 1 file should be uploaded');
                  return
            }
            const file = files[0]
            this.spinnerService.loadingMessage = "Uploading excel"

            this.fileOps.uploadExcel(file).subscribe((uploadedColleges: ICollegeExcel[]) => {
                  if (this.validateColleges(uploadedColleges)) {
                        this.excelUploadedColleges = uploadedColleges
                        this.isUploadedSuccessfully = true
                        this.reformCollegesFromExcel()
                  } else {
                        excelFile.value = "";
                  }

            }, (err) => {

                  alert(err)
                  excelFile.value = "";
            })
      }

      addMultipleColleges() {
            console.log("in multiple add");
            console.log("colleges:", this.reformedUploadedColleges);
            this.spinnerService.loadingMessage = "Adding colleges"

            this.collegeService.addMultipleColleges(this.reformedUploadedColleges).subscribe((data: any) => {
                  this.modalRef.close()
                  this.getAllColleges()
                  this.collegeForm.reset()
                  alert("colleges added")
                  console.log(data);

            }, (error) => {
                  console.log(error)

                  if (error.error) {
                        if (error.error.error) {
                              alert(error.error.error)
                              return
                        }
                        alert(error.error)
                        return
                  }
                  alert(error.statusText)
            })

      }

      validateColleges(colleges: ICollegeExcel[]): boolean {
            this.isUploadedSuccessfully = true
            console.log("colleges:", colleges);
            console.log("in validate", colleges);

            if (!colleges || colleges.length == 0) {
                  alert("No colleges")
                  return false
            }
            for (let i = 0; i < colleges.length; i++) {
                  if (!colleges[i].collegeName || colleges[i].collegeName === "") {
                        alert(`college name on row ${i + 2} not specified`)
                        return false
                  }
                  if (!colleges[i].branchName || colleges[i].branchName === "") {
                        alert(`branch name of college: ${colleges[i].collegeName} is not specified`)
                        return false
                  }
                  if (!colleges[i].universityName || colleges[i].universityName === "") {
                        alert(`university name of college: ${colleges[i].collegeName} is not specified`)
                        return false
                  }
                  if (!colleges[i].countryName || colleges[i].countryName === "") {
                        alert(`country name of college: ${colleges[i].collegeName} is not specified`)
                        return false
                  }
                  if (!colleges[i].stateName || colleges[i].stateName === "") {
                        alert(`state name of college: ${colleges[i].collegeName} is not specified`)
                        return false
                  }
                  if (!colleges[i].city || colleges[i].city === "") {
                        alert(`city of college: ${colleges[i].collegeName} is not specified`)
                        return false
                  }
                  if (!colleges[i].pinCode || colleges[i].pinCode === "") {
                        alert(`PIN code of college: ${colleges[i].collegeName} is not specified`)
                        return false
                  }
                  if (!colleges[i].address || colleges[i].address === "") {
                        alert(`address of college: ${colleges[i].collegeName} is not specified`)
                        return false
                  }
            }
            return true
      }

      reformCollegesFromExcel() {
            // this.excelUploadedColleges.sort((a, b) => {
            //       let num = a.collegeName.localeCompare(b.collegeName)
            //       console.log(num);
            //       return num;
            // })
            // console.log(this.excelUploadedColleges);


            // this.excelUploadedColleges.forEach(excelCollege => {
            //       let college : ICollege
            //       college.collegeName = excelCollege.collegeName
            //       while ()
            // });
            // let allRoutes = this.excelUploadedColleges.reduce(function (routes, e) {
            //       if (!routes[e.collegeName]) routes[e.collegeName] = [];
            //       routes[e.collegeName].push({
            //             branchName: e.branchName,
            //             city: e.city,
            //             address: e.address,
            //             countryName: e.countryName,
            //             stateName: e.stateName,
            //             pinCode: e.pinCode
            //       });
            //       return routes;
            // }, {})
            let collegeMap: Map<string, ICollegeBranch[]> = new Map()
            this.excelUploadedColleges.forEach(college => {
                  let branches = collegeMap.get(college.collegeName)
                  if (branches) {
                        branches.push({
                              branchName: college.branchName,
                              university: {
                                    id: null,
                                    universityName: college.universityName
                              },
                              address: {
                                    address: college.address,
                                    city: college.city,
                                    pinCode: college.pinCode,
                                    country: {
                                          id: null,
                                          name: college.countryName
                                    },
                                    state: {
                                          id: null,
                                          name: college.stateName
                                    }
                              }
                        })
                  } else {
                        collegeMap.set(college.collegeName, [{
                              branchName: college.branchName,
                              university: {
                                    id: null,
                                    universityName: college.universityName
                              },
                              address: {
                                    address: college.address,
                                    city: college.city,
                                    pinCode: college.pinCode,
                                    country: {
                                          id: null,
                                          name: college.countryName
                                    },
                                    state: {
                                          id: null,
                                          name: college.stateName
                                    }
                              }
                        }])
                  }

            })

            // console.log("test value:", allRoutes);
            console.log("final value:", collegeMap);

            this.reformedUploadedColleges = Array.from(collegeMap, ([collegeName, collegeBranches]) =>
                  ({ collegeName, collegeBranches }));

            console.log("God:", this.reformedUploadedColleges);

            // let plotRoute = function (r: any) { console.log('plotting route: ', r) }
            // Object.values(allRoutes).map(plotRoute)
      }


      // ==========================================END TEST======================================================

      // On Update Button Click in Table.
      updateCollegeForm(college: ICollege, collegeDetailModal: any): void {
            this.showBranchForm = false
            this.showCollegeControls = true
            this.college = college
            this.opencollegeDetailModal(collegeDetailModal)
            this.setVariableForUpdateCollege()
            this.updateForm()
      }

      // Add College.
      addCollegeBranch(): void {
            let collegeBranch: ICollegeBranch = this.collegeForm.value.collegeBranches[0]
            collegeBranch.collegeID = this.selectedCollegeID
            console.log(collegeBranch)
            this.spinnerService.loadingMessage = "Adding college branch"

            this.collegeService.addCollegeBranch(collegeBranch).subscribe((data: string) => {
                  this.modalRef.close()
                  this.getAllColleges()
                  this.collegeForm.reset()
                  alert("college branch added with id:" + data)
            }, (error) => {
                  console.log(error)

                  if (error.error) {
                        if (error.error.error) {
                              alert(error.error.error)
                              return
                        }
                        alert(error.error)
                        return
                  }
                  alert(error.statusText)
            })
      }

      updateForm(): void {
            this.createUpdateCollegeForm()
            this.collegeForm.patchValue(this.college)
      }

      setVariableForAddCollege(): void {
            this.updateVariable(this.createAddCollegeForm, "New College", "Add college", this.addCollege)
      }

      setVariableForAddCollegeBranch(): void {
            this.updateVariable(this.createAddCollegeForm, "New College branch", "Add branch", this.addCollegeBranch)
      }

      setVariableForUpdateCollege(): void {
            this.updateVariable(this.updateForm, "Update College", "Update College", this.updateCollege)
      }


      updateVariable(formaction: () => void, modalheader: string, modalbutton: string, modalaction: () => void): void {
            this.formHandler = formaction
            this.modalHeader = modalheader
            this.modalButton = modalbutton
            this.modalAction = modalaction
      }

      // ***************** Use one open modal function and pass properties as args***********
      // modal for add/update form
      opencollegeDetailModal(collegeDetailModal: any, modalSize?: string): void {
            this.isUploadedSuccessfully = false

            if (modalSize == undefined) {
                  modalSize = "xl"
            }

            this.modalRef = this.modalService.open(collegeDetailModal, {
                  ariaLabelledBy: 'modal-basic-title', keyboard: false,
                  backdrop: 'static', size: modalSize
            })
            /*this.modalRef.result.then((result) => {
            }, (reason) => {
            })*/
      }

      //check form validation
      validate(): void {
            console.log(this.collegeForm.controls)

            if (this.collegeForm.invalid) {
                  console.log(this.findInvalidControls())
                  this.collegeForm.markAllAsTouched()
                  return
            }
            this.modalAction()
      }

      public findInvalidControls() {
            const invalid = []
            const controls = this.collegeForm.controls
            for (const name in controls) {
                  if (controls[name].invalid) {
                        invalid.push(name)
                  }
            }
            return invalid
      }


      // onFileChange(event: any) {
      //       let workBook = null;
      //       let jsonData = null;
      //       const reader = new FileReader();
      //       const file = event.target.files[0];
      //       reader.onload = () => {
      //             const data = reader.result;
      //             // will try type string
      //             workBook = XLSX.read(data, { type: 'binary' })
      //             console.log(workBook);

      //             jsonData = workBook.SheetNames.reduce((initial, name) => {
      //                   const sheet = workBook.Sheets[name];
      //                   initial[name] = XLSX.utils.sheet_to_json(sheet, { blankrows: true });
      //                   //  /** Output format */
      //                   //     header?: "A"|number|string[];

      //                   //     /** Override worksheet range */
      //                   //     range?: any;

      //                   //     /** Include or omit blank lines in the output */
      //                   //     blankrows?: boolean;

      //                   //     /** Default value for null/undefined values */
      //                   //     defval?: any;

      //                   //     /** if true, return raw data; if false, return formatted text */
      //                   //     raw?: boolean;

      //                   //     /** if true, return raw numbers; if false, return formatted numbers */
      //                   //     rawNumbers?: boolean;
      //                   return initial;
      //             }, {});
      //             const dataString = JSON.stringify(jsonData);
      //             console.log(dataString.slice(0, 300).concat("..."))
      //       }
      //       reader.readAsBinaryString(file)
      // }
}

interface ICollegeExcel {
      collegeName: string
      branchName: string
      countryName: string
      stateName: string
      universityName: string
      city: string
      pinCode: string
      address: string

}