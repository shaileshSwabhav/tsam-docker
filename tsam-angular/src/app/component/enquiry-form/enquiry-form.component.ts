import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { EnquiryFormService, IAcademic } from 'src/app/service/talent/enquiry-form.service';
import { ActivatedRoute } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { DegreeService } from 'src/app/service/degree/degree.service';

@Component({
  selector: 'app-enquiry-form',
  templateUrl: './enquiry-form.component.html',
  styleUrls: ['./enquiry-form.component.css']
})
export class EnquiryFormComponent implements OnInit {

  // Components.
  enquiryTypeList: any[]
  degreeList: any[]
  specializationList: any[]
  sourceList: any[]
  courseList: any[]

  // FormGroup type.
  enquiryForm: FormGroup

  // Enquiry reatled numbers.
  currentYear: number

  // Enquiry.
  enquiry: any

  // Source.
  sourceNameParam: string
  selectedSourceID: string

  // Course.
  courseCodeArray: string[]

  // Enquiry Type.
  enquiryTypeParam: string
  selectedEnquiryType: string

  // String.


  // Resume.
  docStatus: string
  displayedFileName: string
  isFileUploading: boolean

  // Flags.
  showForm: boolean
  showRedirectMessage: boolean

  constructor(
    private formBuilder: FormBuilder,
    private enquiryFormService: EnquiryFormService,
    private activatedRoute: ActivatedRoute,
    private spinnerService: SpinnerService,
    private fileOperationService: FileOperationService,
    private generalService: GeneralService,
    private degreeService: DegreeService,
  ) {
    this.initializeVariables()
    this.getAllComponent()
  }

  // ngAfterViewInit() {
  //   this.elementRef.nativeElement.ownerDocument.body.style.backgroundColor = 'rgb(29, 29, 77)'
  // }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize all global variables.
  initializeVariables(): void {

    // Components.
    this.enquiryTypeList = []
    this.degreeList = []
    this.specializationList = []
    this.sourceList = []
    this.courseList = []

    // Flags.
    this.showForm = true
    this.showRedirectMessage = false

    // Resume.
    this.isFileUploading = false
    this.docStatus = ""
    this.displayedFileName = "Select file"

    // Enquiry.
    this.enquiry = {}

    // Course.
    this.courseCodeArray = []

    // Enquiry Type.
    this.selectedEnquiryType = null

    // Date.
    this.currentYear = new Date().getFullYear()

    // Create forms.
    this.createEnquiryForm()

    // Get source, enquiry for and course from url params.
    let queryParams = this.activatedRoute.snapshot.queryParams

    // Source.
    this.sourceNameParam = queryParams['src']
    if (!this.sourceNameParam) {
      this.sourceNameParam = "web"
    }

    // Course.
    let courseCodeParam: string = queryParams['crs']
    if (courseCodeParam) {
      this.courseCodeArray = courseCodeParam.split(',')
    }

    // Enquiry type.
    this.enquiryTypeParam = queryParams['ef']

    // // Below will clear the routes query params
    // this.router.navigate(
    //   ['.'],
    //   { relativeTo: this.activatedRoute, queryParams: {} }
    // );

  }

  //*********************************************CREATE FORMS************************************************************
  // Create new enquiry form.
  createEnquiryForm(): void {
    this.enquiryForm = this.formBuilder.group({
      firstName: new FormControl(null, [Validators.pattern(/^[a-zA-Z]*$/), Validators.required]),
      lastName: new FormControl(null, [Validators.pattern(/^[a-zA-Z]*$/), Validators.required]),
      email: new FormControl(null, [Validators.pattern(/^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$/), Validators.required]),
      contact: new FormControl(null, [Validators.pattern(/^[6789]\d{9}$/), Validators.required]),
      degree: new FormControl(null, [Validators.required]),
      specialization: new FormControl({ value: null, disabled: true }),
      courses: new FormControl(null, [Validators.required]),
      passout: new FormControl(null, [Validators.required, Validators.min(1980), Validators.max(this.currentYear + 3)]),
      city: new FormControl(null, [Validators.pattern(/^[a-zA-Z ]*$/), Validators.required]),
      resume: new FormControl(null)
    })
  }

  //*********************************************FUNCTIONS FOR ALL FORMS************************************************************
  // Validate enquiry form.
  validateEnquiryForm(): void {
    if (this.enquiryForm.invalid) {
      // console.log(this.findInvalidControls())
      this.enquiryForm.markAllAsTouched()
      return
    }
    this.addEnquiry()
  }

  // Give specific specializations dropdwon by specific degree.
  showSpecificSpecializations(degree: any, specialization: any): boolean {
    if (!degree.value) {
      return false
    }
    if (degree.value.id === specialization.degreeID) {
      return true
    }
    return false
  }

  // On uplaoding resume.
  onResourceSelect(event: any): void {
    this.docStatus = ""
    let files = event.target.files
    if (files && files.length) {
      let file = files[0]

      // Upload resume if it is present.
      this.isFileUploading = true
      this.fileOperationService.uploadResume(file).subscribe((data: any) => {
        this.enquiryForm.markAsDirty()
        this.enquiryForm.patchValue({
          resume: data
        })
        this.displayedFileName = file.name
        this.isFileUploading = false
        this.docStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
      }, (error) => {
        this.isFileUploading = false
        this.docStatus = `<p><span>&#10060;</span> ${error}</p>`
      })
    }
  }

  // Compare for select option field.
  compareFn(optionOne: any, optionTwo: any): boolean {
    if (optionOne == null && optionTwo == null) {
      return true
    }
    if (optionTwo != undefined && optionOne != undefined) {
      return optionOne.id === optionTwo.id
    }
    return false
  }

  // Get the source id by the source name.
  getSourceIDBySourceName(): void {
    for (let i = 0; i < this.sourceList.length; i++) {
      if (this.sourceList[i].name == this.sourceNameParam) {
        this.selectedSourceID = this.sourceList[i].id
        return
      }
    }
  }

  // Get the course by course code.
  getCourseByCourseCode(): void {
    let selectedCourseList: any[] = []
    for (let i = 0; i < this.courseCodeArray.length; i++) {
      for (let j = 0; j < this.courseList.length; j++) {
        if (this.courseCodeArray[i] == this.courseList[j].code) {
          selectedCourseList.push(this.courseList[j])
        }
      }
    }
    if (selectedCourseList.length > 0) {
      this.enquiryForm.get('courses').setValue(selectedCourseList)
    }
  }

  redirectToHomePage(): void {
    // this.elementRef.nativeElement.ownerDocument.body.style.backgroundColor = 'white'
    window.location.href = 'https://swabhavtechlabs.com/'
  }

  // Get all invalid controls in talent form.
  public findInvalidControls(): any[] {
    const invalid = []
    const controls = this.enquiryForm.controls
    for (const name in controls) {
      if (controls[name].invalid) {
        invalid.push(name)
      }
    }
    return invalid
  }

  // Get the enquiry type by enquiry type param.
  getEnquiryTypeByEnquiryTypeName(): void {
    if (this.enquiryTypeParam) {
      for (let i = 0; i < this.enquiryTypeList.length; i++) {
        if (this.enquiryTypeList[i]?.value?.toLowerCase() === this.enquiryTypeParam?.toLowerCase()) {
          this.selectedEnquiryType = this.enquiryTypeList[i].value
          return
        }
      }
    }
  }

  //*********************************************ADD FUNCTIONS FOR ENQUIRY************************************************************
  // Add New Enquiry.
  addEnquiry(): void {

    this.spinnerService.loadingMessage = "Sending enquiry...Please wait"
    this.assignEnquiry()
    this.enquiryFormService.addEnquiry(this.enquiry).subscribe((data) => {

      this.showForm = false
      this.showRedirectMessage = true
    }, (error) => {
      console.error(error)

      alert("Enquiry could not be sent. Please try again!")
    })
  }

  // Assign all fields to enquiry from enquiry form and params.
  assignEnquiry(): void {
    this.enquiry.firstName = this.enquiryForm.get('firstName').value
    this.enquiry.lastName = this.enquiryForm.get('lastName').value
    this.enquiry.email = this.enquiryForm.get('email').value
    this.enquiry.contact = this.enquiryForm.get('contact').value
    if (this.enquiryForm.get('resume').value) {
      this.enquiry.resume = this.enquiryForm.get('resume').value
    }
    this.enquiry.courses = this.enquiryForm.get('courses').value
    this.enquiry.sourceID = this.selectedSourceID
    if (this.enquiryTypeParam) {
      this.enquiry.enquiryType = this.selectedEnquiryType
    }
    let specializationID: string = null
    if (this.enquiryForm.get('specialization').value) {
      specializationID = this.enquiryForm.get('specialization').value.id
    }
    let academic: IAcademic = {
      college: null,
      collegeID: null,
      enquiryID: null,
      percentage: null,
      degreeID: this.enquiryForm.get('degree').value.id,
      passout: this.enquiryForm.get('passout').value,
      specializationID: specializationID,
    }
    this.enquiry.academics = []
    this.enquiry.academics.push(academic)
  }

  //*********************************************GET FUNCTIONS************************************************************
  // Populate all lists.
  getAllComponent(): void {
    this.getDegreeList()
    this.getSourceList()
    this.getCourseList()
    this.getEnquiryTypeList()
  }

  // Get technology list.
  getSourceList(): void {
    this.generalService.getSources().subscribe(
      data => {
        this.sourceList = data
        this.getSourceIDBySourceName()
      }
    ), err => {
      console.error(err)
    }
  }

  // Get enquiry type list.
  getEnquiryTypeList(): void {
    this.generalService.getGeneralTypeByType("talent_enquiry_type").subscribe(
      response => {
        this.enquiryTypeList = response
        this.getEnquiryTypeByEnquiryTypeName()
      }
    ), err => {
      console.error(err)
    }
  }

  // Get degree list.
  getDegreeList(): void {
    let queryParams: any = {
      limit: -1,
      offset: 0,
    }
    this.degreeService.getAllDegrees(queryParams).subscribe(
      data => {
        this.degreeList = data.body
      }
    ), err => {
      console.error(err)
    }
  }

  // Get course list.
  getCourseList(): void {
    this.generalService.getCourseList().subscribe(
      response => {
        this.courseList = response.body
        this.getCourseByCourseCode()
      }
    ), err => {
      console.error(err)
    }
  }

  // Get specialization list by degree id.
  getSpecializationListByDegreeID(degree: any): void {
    this.enquiryForm.get('specialization').setValue(null)
    this.specializationList = []
    if (!degree) {
      this.enquiryForm.get('specialization').setValue(null)
      this.specializationList = []
      this.enquiryForm.get('specialization').disable()
      return
    }
    this.generalService.getSpecializationByDegreeID(degree.id).subscribe((response: any) => {
      this.specializationList = response.body
      this.enquiryForm.get('specialization').enable()
    }, (err) => {
      console.error(err)
    })
  }
}


