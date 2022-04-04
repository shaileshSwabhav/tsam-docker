import { Component, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, FormControl, Validators, FormArray } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { faFacebookSquare, faInstagramSquare, faGithubSquare, faLinkedinIn } from '@fortawesome/free-brands-svg-icons';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { DegreeService } from 'src/app/service/degree/degree.service';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { EnquiryFormService, IEnquiry, IExamination, IScore, IMastersAbroad } from 'src/app/service/talent/enquiry-form.service';
import { TechnologyService } from 'src/app/service/technology/technology.service';

@Component({
  selector: 'app-enquiry-enrollment',
  templateUrl: './enquiry-enrollment.component.html',
  styleUrls: ['./enquiry-enrollment.component.css']
})
export class EnquiryEnrollmentComponent implements OnInit {

  // Components.
  technologyList: any[]
  enquiryTypeList: any[]
  countryList: any[]
  stateList: any[]
  degreeList: any[]
  specializationList: any[]
  academicYearList: any[]
  designationList: any[]
  examinationList: any[]
  universityList: any[]
  collegeBranchList: any[]
  sourceList: any[]
  yearOfMSList: number[]

  // tech component
  isTechLoading: boolean
  techLimit: number
  techOffset: number

  // FormGroup type.
  personalDetailsForm: FormGroup
  addressForm: FormGroup
  academicsForm: FormGroup
  experiencesForm: FormGroup
  socialMediaForm: FormGroup
  resumeForm: FormGroup
  mastersAbroadForm: FormGroup

  // Flags.
  showMastersAbroadForm: boolean
  showExperiencesForm: boolean
  isFileUploading: boolean
  showMBADegreeRequiredError: boolean

  // Enquiry reatled numbers.
  indexOfCurrentWorkingExp: number
  currentYear: number
  degreeIDSet: Set<string>

  // Social media icons.
  facebookIcon: any
  instagramIcon: any
  githubIcon: any
  linkedInIcon: any

  // Enquiry.
  enquiry: IEnquiry

  // Source.
  sourceName: any
  sourceID: string

  // Exam.
  greExam: IExamination
  gmatExam: IExamination
  toeflExam: IExamination
  ieltsExam: IExamination

  // String.

  docStatus: string
  displayedFileName: string

  // Map.
  universityMap: Map<string, any[]> = new Map()
  areUniversitiesLoading: boolean

  constructor(
    private formBuilder: FormBuilder,
    private enquiryFormService: EnquiryFormService,
    private activatedRoute: ActivatedRoute,
    private spinnerService: SpinnerService,
    private fileOperationService: FileOperationService,
    private generalService: GeneralService,
    private techService: TechnologyService,
    private degreeService: DegreeService,
  ) {
    this.initializeVariables()
    this.getAllComponent()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize all global variables.
  initializeVariables(): void {

    // Components.
    this.technologyList = []
    this.enquiryTypeList = []
    this.countryList = []
    this.stateList = []
    this.degreeList = []
    this.specializationList = []
    this.academicYearList = []
    this.designationList = []
    this.examinationList = []
    this.universityList = []
    this.collegeBranchList = []
    this.sourceList = []
    this.yearOfMSList = []
    this.degreeIDSet = new Set()

    // Experience.
    this.indexOfCurrentWorkingExp = -1

    // Resume.
    this.isFileUploading = false
    this.docStatus = ""
    this.displayedFileName = "Select file"

    // Date.
    this.currentYear = new Date().getFullYear()

    // Social media.
    this.facebookIcon = faFacebookSquare
    this.instagramIcon = faInstagramSquare
    this.githubIcon = faGithubSquare
    this.linkedInIcon = faLinkedinIn

    // Forms' visibility.
    this.showMastersAbroadForm = false
    this.showExperiencesForm = false
    this.showMBADegreeRequiredError = false

    // Masters abroad.
    this.areUniversitiesLoading = false
    this.universityList = []

    // tech component
    this.isTechLoading = false
    this.techLimit = 10
    this.techOffset = 0

    // Get source from url params.
    this.activatedRoute.queryParams.subscribe(params => {
      this.sourceName = params['source']
      if (this.sourceName == undefined) {
        this.sourceName = "web"
      }
      this.getSourceList()
    })

    // Create forms.
    this.createAllForms()
  }

  //*********************************************Create Forms************************************************************
  // Create new personal details form.
  createPersonalDetailsForm(): void {
    this.personalDetailsForm = this.formBuilder.group({
      firstName: new FormControl(null, [Validators.pattern(/^[a-zA-Z]*$/), Validators.required]),
      lastName: new FormControl(null, [Validators.pattern(/^[a-zA-Z]*$/), Validators.required]),
      email: new FormControl(null, [Validators.pattern(/^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$/), Validators.required]),
      contact: new FormControl(null, [Validators.pattern(/^[6789]\d{9}$/), Validators.required]),
      technologies: new FormControl(Array()),
      enquiryType: new FormControl(null, [Validators.required]),
      academicYear: new FormControl(null, [Validators.required]),
    })
  }

  // Create new address form.
  createAddressForm(): void {
    this.addressForm = this.formBuilder.group({
      address: new FormControl(null, [Validators.pattern(/^[.0-9a-zA-Z\s,-\/]+$/)]),
      city: new FormControl(null, [Validators.pattern(/^[a-zA-Z]+([a-zA-Z ]?)+$/)]),
      state: new FormControl({ value: null, disabled: true }),
      country: new FormControl(null),
      pinCode: new FormControl(null, [Validators.pattern(/^[1-9][0-9]{5}$/)]),
    })
  }

  // Create new academics form.
  createAcademicsForm(): void {
    this.academicsForm = this.formBuilder.group({
      academics: this.formBuilder.array([])
    })
    this.addAcademic()
  }

  // Create new experiences form.
  createExperiencesForm(): void {
    this.experiencesForm = this.formBuilder.group({
      numberOfExperiences: new FormControl(1, [Validators.required, Validators.min(1)]),
      resume: new FormControl(null),
      experiences: this.formBuilder.array([])
    })
    this.addExperience()
  }

  // Create new social media form.
  createSocialMediaForm(): void {
    this.socialMediaForm = this.formBuilder.group({
      facebookUrl: new FormControl(null, [Validators.pattern(/^(?:https?:\/\/)?(?:www\.)?facebook\.com\/.(?:(?:\w)*#!\/)?(?:pages\/)?(?:[\w\-]*\/)*([\w\-\.]*)$/)]),
      instagramUrl: new FormControl(null, [Validators.pattern(/(?:(?:http|https):\/\/)?(?:www.)?(?:instagram.com|instagr.am)\/([A-Za-z0-9-_]+)/)]),
      githubUrl: new FormControl(null, [Validators.pattern(/^(http(s?):\/\/)?(www\.)?github\.([a-z])+\/([A-Za-z0-9]{1,})+\/?$/)]),
      linkedInUrl: new FormControl(null, [Validators.pattern(/(https?)?:?(\/\/)?(([w]{3}||\w\w)\.)?linkedin.com(\w+:{0,1}\w*@)?(\S+)(:([0-9])+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?/)]),
    })
  }

  // Create new masters abroad form.
  createMastersAbroadForm(): void {
    this.mastersAbroadForm = this.formBuilder.group({
      id: new FormControl(null),
      degree: new FormControl(null),
      countries: new FormControl(null),
      universities: new FormControl({ value: Array(), disabled: true }),
      gre: new FormControl(null, [Validators.required, Validators.min(0), Validators.max(this.greExam.totalMarks)]),
      gmat: new FormControl(null, [Validators.min(0), Validators.max(this.gmatExam.totalMarks)]),
      toefl: new FormControl(null, [Validators.min(0), Validators.max(this.toeflExam.totalMarks)]),
      ielts: new FormControl(null, [Validators.min(0), Validators.pattern(/^(10|\d)(\.\d{1,2})?$/),
      Validators.max(this.ieltsExam.totalMarks)]),
      yearOfMS: new FormControl(null),
    })
  }

  // Add new experience to experiences form.
  addExperience(): void {
    this.experienceControlArray.push(this.formBuilder.group({
      company: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z ]*$/)]),
      designation: new FormControl(null, [Validators.required]),
      technologies: new FormControl(null, [Validators.required]),
      isCurrentWorking: new FormControl(false),
      fromDate: new FormControl(null, [Validators.required]),
      toDate: new FormControl(null, [Validators.required]),
    }))
  }

  // Add new academic to academics form.
  addAcademic(): void {
    this.academicControlArray.push(this.formBuilder.group({
      degree: new FormControl(null, [Validators.required]),
      specialization: new FormControl(null),
      college: new FormControl(null, [Validators.required]),
      percentage: new FormControl(null, [Validators.min(1), Validators.pattern(/(^100([.]0{1,2})?)$|(^\d{1,2}([.]\d{1,2})?)$/i), Validators.required]),
      passout: new FormControl(null, [Validators.required, Validators.min(1980), Validators.max(this.currentYear + 3)])
    }))
  }

  //*********************************************Functions for all forms************************************************************
  // Get expereinces array from expereinces form.
  get experienceControlArray(): FormArray {
    return this.experiencesForm.get("experiences") as FormArray
  }

  // Get academics array from academics form.
  get academicControlArray(): FormArray {
    return this.academicsForm.get("academics") as FormArray
  }

  // Validate personal details form.
  validatePersonalDetailsForm(): void {
    if (this.personalDetailsForm.invalid) {
      this.personalDetailsForm.markAllAsTouched()
      return
    }
  }

  // Validate academics form.
  validateAacdemicsForm(): void {
    if (this.academicsForm.invalid) {
      this.academicsForm.markAllAsTouched()
      return
    }
  }

  // Validate experiences form.
  validateExperiencesForm(): void {
    if (this.experiencesForm.invalid) {
      this.experiencesForm.markAllAsTouched()
      return
    }
  }

  // Validate masters abroad form.
  validateMastersAboradForm(): void {
    if (this.mastersAbroadForm.invalid) {
      this.mastersAbroadForm.markAllAsTouched()
      return
    }
  }

  // Delete academic from academics form.
  deleteAcademic(index: number): void {
    if (confirm("Are you sure you want to delete all academic detials?")) {
      this.academicsForm.markAsDirty()
      this.academicControlArray.removeAt(index)
    }
  }

  // Delete experience from experiences form.
  deleteExperience(index: number): void {
    if (confirm("Are you sure you want to delete all experience detials?")) {
      if (this.indexOfCurrentWorkingExp == index) {
        this.indexOfCurrentWorkingExp = -1
      }
      if (this.indexOfCurrentWorkingExp > index) {
        this.indexOfCurrentWorkingExp = this.indexOfCurrentWorkingExp - 1
      }
      this.experiencesForm.markAsDirty()
      this.experienceControlArray.removeAt(index)
      this.experiencesForm.get('numberOfExperiences').setValue(this.experienceControlArray.length)
    }
  }

  // Give specific specializations dropdwon by specific degree.
  showSpecificSpecializations(academic: any, specialization: any): boolean {
    if (!academic.get('degree').value) {
      return false
    }
    if (academic.get('degree').value.id === specialization.degreeID) {
      return true
    }
    return false
  }

  // On clicking currently working or not.
  isWorkingClicked(isWorking: string, index: number, expereince: any): void {
    if (isWorking == "false") {
      expereince.addControl('toDate', this.formBuilder.control(null, [Validators.required]))
      this.indexOfCurrentWorkingExp = -1
      return
    }
    expereince.removeControl('toDate')
    this.indexOfCurrentWorkingExp = index
  }

  // Reset all forms on submit.
  createAllForms(): void {
    this.createPersonalDetailsForm()
    this.createAddressForm()
    this.createAcademicsForm()
    this.createSocialMediaForm()
  }

  // Show forms based on enquiry type.
  assignSelectedEnquiryType(enquiryType: string): void {
    if (enquiryType == "Placement") {//MS in Us
      this.showMastersAbroadForm = true
      this.showExperiencesForm = false
      this.createMastersAbroadForm()
      return
    }
    if (enquiryType == "Training And Placement") {//Job
      this.showMastersAbroadForm = false
      this.showExperiencesForm = true
      this.createExperiencesForm()
      return
    }
    this.showMastersAbroadForm = false
    this.showExperiencesForm = false
  }

  // On uploading resume.
  onResourceSelect(event: any): void {
    this.docStatus = ""
    let files = event.target.files
    if (files && files.length) {
      let file = files[0]
      // Upload resume if it is present.
      this.isFileUploading = true
      this.fileOperationService.uploadResume(file).subscribe((data: any) => {
        if (this.showMastersAbroadForm) {
          this.mastersAbroadForm.patchValue({
            resume: data
          })
        }
        if (this.showExperiencesForm) {
          this.experiencesForm.patchValue({
            resume: data
          })
        }
        this.displayedFileName = file.name
        this.isFileUploading = false
        //this.isResumeUploadedToServer = true
        this.docStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
      }, (error) => {
        this.isFileUploading = false
        this.docStatus = `<p><span>&#10060;</span> ${error}</p>`
      })
    }
  }

  // On degree change in masters abroad form.
  onDegreeChange(degree: any): void {
    if (degree != null && degree.name == "MBA") {
      this.mastersAbroadForm.get('gmat').setValidators([Validators.required, Validators.min(0), Validators.max(this.gmatExam.totalMarks)])
      this.mastersAbroadForm.get('gmat').updateValueAndValidity()
      this.showMBADegreeRequiredError = true
      return
    }
    this.mastersAbroadForm.get('gmat').setValidators([Validators.min(0), Validators.max(this.gmatExam.totalMarks)])
    this.mastersAbroadForm.get('gmat').updateValueAndValidity()
    this.showMBADegreeRequiredError = false
  }

  // Validate the masters abroad form or experiences form.
  validateExperiencesOrMastersAbroad(): void {
    if (this.showExperiencesForm) {
      this.validateExperiencesForm()
      return
    }
    if (this.showMastersAbroadForm) {
      this.validateMastersAboradForm()
    }
  }

  // On numbers of companies change.
  numberOfCompaniesChange(numberOfExperiences: any): void {
    if (numberOfExperiences.value) {
      if (numberOfExperiences.value == this.experienceControlArray.length) {
        return
      }
      if (numberOfExperiences.value > this.experienceControlArray.length) {
        let diff: number = numberOfExperiences.value - this.experienceControlArray.length
        for (let i = 0; i < diff; i++) {
          this.addExperience()
        }
        return
      }
      if (numberOfExperiences.value < 0) {
        return
      }
      if (numberOfExperiences.value < this.experienceControlArray.length) {
        this.experiencesForm.get('numberOfExperiences').setValue(this.experienceControlArray.length)
        alert("Please delete one or more experiences")
        return
      }
      this.experiencesForm.get('numberOfExperiences').setValue(this.experienceControlArray.length)
    }
  }

  // On clearing all country selections in masters abroad form.
  onClearCountries(): void {
    this.universityList = []
    this.mastersAbroadForm.get('universities').reset()
    this.mastersAbroadForm.get('universities').disable()
    for (let universities of Array.from(this.universityMap.values())) {
      for (let university of universities) {
        university.isVisible = false
      }
    }
  }

  // On removing one country from masters abroad form.
  onRemoveCountry(id: any): void {
    this.universityList = []
    let selectedUniversities: any[] = this.mastersAbroadForm.get('universities').value
    this.mastersAbroadForm.get('universities').reset()
    for (let [countryID, universities] of this.universityMap) {
      if (countryID == id) {
        for (let university of universities) {
          university.isVisible = false
        }
        continue
      }
      for (let university of universities) {
        if (!university.isVisible) {
          break
        }
        this.universityList = this.universityList.concat(university)
      }
    }
  }

  // On adding one country in masters abroad form.
  onAddCountry(countryID: string): void {
    this.mastersAbroadForm.get('universities').enable()
    if (this.universityMap.has(countryID)) {
      let allUniversities = this.universityMap.get(countryID)
      for (let university of allUniversities) {
        university.isVisible = true
        this.universityList = this.universityList.concat(university)
      }
      return
    }
    this.areUniversitiesLoading = true
    this.generalService.getUniversityByCountryID(countryID).subscribe((response: any[]) => {
      for (let i = 0; i < response.length; i++) {
        response[i].isVisible = true
        response[i].countryName = response[i].country?.name
      }
      this.areUniversitiesLoading = false
      this.universityMap.set(countryID, response)
      this.universityList = this.universityList.concat(response)
    }, err => {
      this.areUniversitiesLoading = false
      console.error(err)
    })
  }

  // To aply validation to college name.
  addCollegeToList(option: any): Promise<any> {
    return new Promise((resolve) => {
      resolve(option)
    }
    )
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
      if (this.sourceList[i].name == this.sourceName) {
        this.sourceID = this.sourceList[i].id
      }
    }
  }

  // Extract ID from objects in enquiry form.
  patchIDFromObjectsForEnquiry(): void {
    // Source.
    this.enquiry.sourceID = this.sourceID

    // Designation.
    if (this.experiencesForm) {
      for (let i = 0; i < this.experienceControlArray.length; i++) {
        if (this.experienceControlArray.at(i).get('designation').value) {
          this.enquiry.experiences[i].designationID = this.experienceControlArray.at(i).get('designation').value.id
          delete this.enquiry.experiences['designation']
        }
      }
    }

    // Degree and specialization.
    for (let i = 0; i < this.academicControlArray.length; i++) {
      if (this.academicControlArray.at(i).get('degree').value) {
        this.enquiry.academics[i].degreeID = this.academicControlArray.at(i).get('degree').value.id
        delete this.enquiry.academics['degree']
      }
      if (this.academicControlArray.at(i).get('specialization').value) {
        this.enquiry.academics[i].specializationID = this.academicControlArray.at(i).get('specialization').value.id
        delete this.enquiry.academics['specialization']
      }
    }
  }

  //*********************************************Add functions for enquiry************************************************************
  // Add New Enquiry.
  addEnquiry(): void {

    this.spinnerService.loadingMessage = "Sending enquiry...Please wait"
    this.assignFormDetailsToEnquiry()
    let enquiryJson: string = JSON.stringify(this.enquiry)
    this.enquiryFormService.addEnquiry(this.enquiry).subscribe((data) => {
      this.createAllForms()

      alert("Enquiry sent successfully")
    }, (error) => {
      console.error(error)

      alert("Enquiry could not be sent...Please try again")
    })
  }

  // Assign all forms details to enquiry.
  assignFormDetailsToEnquiry(): void {
    // If MS in US then give MS in US form fields.
    if (this.showMastersAbroadForm) {
      this.enquiry = {
        ...this.personalDetailsForm.value,
        ...this.addressForm.value,
        ...this.academicsForm.value,
        ...this.socialMediaForm.value,
      }
      // To create masters abroad and assign it to enquiry.
      let scoresArray: IScore[] = []
      let mastersAbroad: IMastersAbroad
      if (this.mastersAbroadForm.get('gre').value != null) {
        let score: IScore = {
          marksObtained: this.mastersAbroadForm.get('gre').value,
          examinationID: this.greExam.id,
        }
        scoresArray.push(score)
      }
      if (this.mastersAbroadForm.get('gmat').value != null) {
        let score: IScore = {
          marksObtained: this.mastersAbroadForm.get('gmat').value,
          examinationID: this.gmatExam.id,
        }
        scoresArray.push(score)
      }
      if (this.mastersAbroadForm.get('toefl').value != null) {
        let score: IScore = {
          marksObtained: this.mastersAbroadForm.get('toefl').value,
          examinationID: this.toeflExam.id,
        }
        scoresArray.push(score)
      }
      if (this.mastersAbroadForm.get('ielts').value != null) {
        let score: IScore = {
          marksObtained: this.mastersAbroadForm.get('ielts').value,
          examinationID: this.ieltsExam.id,
        }
        scoresArray.push(score)
      }
      mastersAbroad = {
        scores: scoresArray,
        degreeID: this.mastersAbroadForm.get('degree').value.id,
        countries: this.mastersAbroadForm.get('countries').value,
        universities: this.mastersAbroadForm.get('universities').value,
        yearOfMS: this.mastersAbroadForm.get('yearOfMS').value,
      }
      this.enquiry.mastersAbroad = mastersAbroad
      this.enquiry.isMastersAbroad = true
    }

    // If experienced then give experiences form fields.
    if (this.showExperiencesForm) {
      this.enquiry = {
        ...this.personalDetailsForm.value,
        ...this.addressForm.value,
        ...this.academicsForm.value,
        ...this.experiencesForm.value,
        ...this.socialMediaForm.value,
      }
      this.enquiry.isExperience = true
    }
    this.patchIDFromObjectsForEnquiry()
  }

  //*********************************************Get functions************************************************************
  // Populate all lists.
  getAllComponent(): void {
    this.getTechnologyList()
    this.getEnquiryTypeList()
    this.getCountryList()
    this.getDegreeList()
    this.getAcademicYearList()
    this.getDesignationList()
    this.getExaminationList()
    this.getCollegeBranchList()
    this.getYearOfMSList()
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

  // Get technology list.
  getTechnologyList(event?: any): void {
    // this.generalService.getTechnologies().subscribe(
    //   data => {
    //     this.technologyList = data
    //   }
    // ), err => {
    //   console.error(err)
    // }
    let queryParams: any = {}
    if (event && event?.term != "") {
      queryParams.language = event?.term
    }
    this.isTechLoading = true
    this.techService.getAllTechnologies(this.techLimit, this.techOffset, queryParams).subscribe((response) => {
      // console.log("getTechnology -> ", response);
      this.technologyList = []
      this.technologyList = this.technologyList.concat(response.body)
    }, (err) => {
      console.error(err)
    }).add(() => {
      this.isTechLoading = false
    })
  }

  // Get enquiry type list.
  getEnquiryTypeList(): void {
    this.generalService.getGeneralTypeByType("talent_enquiry_type").subscribe(
      data => {
        this.enquiryTypeList = data
      }
    ), err => {
      console.error(err)
    }
  }

  // Get country list.
  getCountryList(): void {
    this.generalService.getCountries().subscribe(
      data => {
        this.countryList = data
      }
    ), err => {
      console.error(err)
    }
  }

  // Get state list by country ID.
  getStateListByCountryID(country: any) {
    if (country == null) {
      this.stateList = [];
      this.addressForm.get('state').setValue(null)
      this.addressForm.get('state').disable()
      return
    }
    if (country.name == "U.S.A." || country.name == "India") {
      this.addressForm.get('state').enable()
      this.generalService.getStatesByCountryID(country.id).subscribe((respond: any[]) => {
        this.stateList = respond;
      }, (err) => {
        console.error(err)
      });
    }
    else {
      this.addressForm.get('state').setValue(null)
      this.addressForm.get('state').disable()
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

  // Get specialization list by degree id.
  getSpecializationListByDegreeID(degreeID: string, shouldClear?: boolean): void {
    if (this.degreeIDSet.has(degreeID)) {
      return
    }
    this.degreeIDSet.add(degreeID)
    if (degreeID) {
      this.generalService.getSpecializationByDegreeID(degreeID).subscribe((response: any) => {
        this.specializationList = this.specializationList.concat(response.body);
      }, (err) => {
        console.error(err)
      })
    }
  }

  // Get academic year list.
  getAcademicYearList(): void {
    this.generalService.getGeneralTypeByType("academic_year").subscribe(
      data => {
        this.academicYearList = data
      }
    ), err => {
      console.error(err)
    }
  }

  // Get designation list.
  getDesignationList(): void {
    this.generalService.getDesignations().subscribe(
      data => {
        this.designationList = data
      }
    ), err => {
      console.error(err)
    }
  }

  //Get examination list.
  getExaminationList(): void {
    this.generalService.getExaminationList().subscribe(
      data => {
        this.examinationList = data
        this.setExaminationsTotalMarks()
      }
    ), err => {
      console.error(err)
    }
  }

  // Getsuniversity list by country ID.
  getUniversityListByCountryID(country: any): void {
    this.mastersAbroadForm.get('universities').setValue(Array())
    if (country == null) {
      this.universityList = []
      return
    }
    this.generalService.getUniversityByCountryID(country.id).subscribe(response => {
      this.universityList = response
    }, err => {
      console.error(err)
    })
  }

  // Set the examinations name for all scores in masters abroad form.
  setExaminationsTotalMarks(): void {
    for (let i = 0; i < this.examinationList.length; i++) {
      if (this.examinationList[i].name == "GRE") {
        this.greExam = this.examinationList[i]
      }
      if (this.examinationList[i].name == "GMAT") {
        this.gmatExam = this.examinationList[i]
      }
      if (this.examinationList[i].name == "TOEFL") {
        this.toeflExam = this.examinationList[i]
      }
      if (this.examinationList[i].name == "IELTS") {
        this.ieltsExam = this.examinationList[i]
      }
    }
  }

  // Get college branch list.
  getCollegeBranchList(): void {
    this.generalService.getCollegeBranchList().subscribe(response => {
      this.collegeBranchList = response
    }, err => {
      console.error(err)
    })
  }

  // Create list of years.
  getYearOfMSList(): void {
    this.yearOfMSList.push(this.currentYear)
    this.yearOfMSList.push((this.currentYear + 1))
    this.yearOfMSList.push((this.currentYear + 2))
  }

}
