import { DatePipe } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormControl, Validators, FormBuilder, FormArray } from '@angular/forms';
import { NgbModalRef, NgbModalOptions, NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { AccountSettingsService, IPassswordChange } from 'src/app/service/account-settings/account-settings.service';
import { Role } from 'src/app/service/constant';
import { DegreeService } from 'src/app/service/degree/degree.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { TalentService } from 'src/app/service/talent/talent.service';
import { TechnologyService } from 'src/app/service/technology/technology.service';

@Component({
  selector: 'app-account-settings',
  templateUrl: './account-settings.component.html',
  styleUrls: ['./account-settings.component.css']
})
export class AccountSettingsComponent implements OnInit {

  // Components.
  countryList: any[]
  stateList: any[]
  degreeList: any[]
  specializationList: any[]
  collegeBranchList: any[]
  degreeIDSet: Set<string>
  academicYearList: any[]
  technologyList: any[]
  designationList: any[]
  examinationList: any[]
  universityList: any[]
  yearOfMSList: number[]

  // tech component
  isTechLoading: boolean
  techLimit: number
  techOffset: number

  // Password Change.
  passwordChange: IPassswordChange[]
  changePasswordForm: FormGroup

  // Talent.
  isTalent: boolean
  talent: any
  talentID: string
  basicDetailsForm: FormGroup
  socialMediaForm: FormGroup
  aboutMeForm: FormGroup
  mySkillsForm: FormGroup
  academicsForm: FormGroup
  experiencesForm: FormGroup
  mastersAbroadForm: FormGroup

  // Flags.
  noMatch: boolean
  isBasicDetailsFormUpdated: boolean
  isSocialMediaFormUpdated: boolean
  isAboutMeFormUpdated: boolean
  isMySkillsUpdated: boolean
  isAcademicsFormUpdated: boolean
  isExperiencesFormUpdated: boolean
  isMastersAbroadFormUpdated: boolean
  showMBADegreeRequiredError: boolean
  showMastersAbroad: boolean

  // Exam.
  greExam: IExamination
  gmatExam: IExamination
  toeflExam: IExamination
  ieltsExam: IExamination

  // Talent realted numbers.
  currentYear: number
  indexOfCurrentWorkingExp: number

  // Verify password.
  verifyPasswordForm: FormGroup

  // Modal.
  modalRef: any
  @ViewChild('changePasswordFormModal') changePasswordFormModal: any

  // Other variables.
  roleID: string
  email: string

  // Spinner.



  // Map.
  universityMap: Map<string, any[]> = new Map()
  areUniversitiesLoading: boolean

  constructor(
    private formBuilder: FormBuilder,
    private localService: LocalService,
    private spinnerService: SpinnerService,
    private accountSettingsService: AccountSettingsService,
    private modalService: NgbModal,
    private talentService: TalentService,
    private role: Role,
    private datePipe: DatePipe,
    private generalService: GeneralService,
    private techService: TechnologyService,
    private degreeService: DegreeService,
  ) {
    this.initializeVariables()
    if (this.isTalent) {
      this.getTalent()
      this.getAllComponent()
    }
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize all global variables.
  initializeVariables() {

    // Components.
    this.countryList = []
    this.stateList = []
    this.degreeList = []
    this.specializationList = []
    this.collegeBranchList = []
    this.degreeIDSet = new Set()
    this.academicYearList = []
    this.technologyList = []
    this.designationList = []
    this.examinationList = []
    this.universityList = []
    this.yearOfMSList = []

    // Other variables.
    this.roleID = this.localService.getJsonValue("roleID")
    this.email = this.localService.getJsonValue("email")

    // Talent.
    if (this.localService.getJsonValue("roleName") == this.role.TALENT) {
      this.isTalent = true
    }
    this.talent = {}
    this.talentID = this.localService.getJsonValue("loginID")

    // Flags.
    this.noMatch = false
    this.isBasicDetailsFormUpdated = true
    this.isSocialMediaFormUpdated = true
    this.isAboutMeFormUpdated = true
    this.isMySkillsUpdated = true
    this.isAcademicsFormUpdated = true
    this.isExperiencesFormUpdated = true
    this.isMastersAbroadFormUpdated = true
    this.showMBADegreeRequiredError = false
    this.areUniversitiesLoading = false
    this.showMastersAbroad = false
    this.isTechLoading = false

    this.techLimit = 10
    this.techOffset = 0

    // Talent realted numbers.
    this.currentYear = new Date().getFullYear()
    this.indexOfCurrentWorkingExp = -1

    // Spinner.
    this.spinnerService.loadingMessage = "Loading"


    // Create forms.
    this.createVerifyPasswordForm()
    this.createBasicDetailsForm()
    this.createSocialMediaForm()
    this.createAboutMeForm()
    this.createMySkillsForm()
    this.createAcademicsForm()
    this.createExperiencesForm()
    this.createMastersAbroadForm()
  }

  // ============================================================= CREATE FORMS ==========================================================================
  // Create verify password form.
  createVerifyPasswordForm(): void {
    this.verifyPasswordForm = this.formBuilder.group({
      password: new FormControl(null, [Validators.required, Validators.maxLength(100)]),
    })
  }

  // Create change password form.
  createChangePasswordForm(): void {
    this.changePasswordForm = this.formBuilder.group({
      password: new FormControl(null, [Validators.required,
      Validators.pattern(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,100}$/)]),
      confirmPassword: new FormControl(null, [Validators.required,
      Validators.pattern(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,100}$/)]),
    })
  }

  // Create basic details talent form.
  createBasicDetailsForm(): void {
    this.basicDetailsForm = this.formBuilder.group({
      firstName: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z]*$/), Validators.maxLength(50)]),
      lastName: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z]*$/), Validators.maxLength(50)]),
      contact: new FormControl(null, [Validators.required, Validators.pattern(/^[6789]\d{9}$/)]),
      address: new FormControl(null, [Validators.pattern(/^[.0-9a-zA-Z\s,-\/]+$/)]),
      city: new FormControl(null, [Validators.pattern(/^[a-zA-Z]+([a-zA-Z ]?)+$/), Validators.maxLength(50)]),
      state: new FormControl(null),
      country: new FormControl(null),
      pinCode: new FormControl(null, [Validators.pattern(/^[1-9][0-9]{5}$/)]),
      academicYear: new FormControl(null, [Validators.required]),
    })
  }

  // Create social media links talent form.
  createSocialMediaForm(): void {
    this.socialMediaForm = this.formBuilder.group({
      facebookUrl: new FormControl(null, [Validators.maxLength(200),
      Validators.pattern(/^(?:https?:\/\/)?(?:www\.)?facebook\.com\/.(?:(?:\w)*#!\/)?(?:pages\/)?(?:[\w\-]*\/)*([\w\-\.]*)$/)]),
      instagramUrl: new FormControl(null, [Validators.maxLength(200),
      Validators.pattern(/(?:(?:http|https):\/\/)?(?:www.)?(?:instagram.com|instagr.am)\/([A-Za-z0-9-_]+)/)]),
      githubUrl: new FormControl(null, [Validators.maxLength(200),
      Validators.pattern(/^(http(s?):\/\/)?(www\.)?github\.([a-z])+\/([A-Za-z0-9]{1,})+\/?$/)]),
      linkedInUrl: new FormControl(null, [Validators.maxLength(200),
      Validators.pattern(/(https?)?:?(\/\/)?(([w]{3}||\w\w)\.)?linkedin.com(\w+:{0,1}\w*@)?(\S+)(:([0-9])+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?/)]),
    })
  }

  // Create about me talent form.
  createAboutMeForm(): void {
    this.aboutMeForm = this.formBuilder.group({
      about: new FormControl(null, [Validators.pattern(/^[a-zA-Z ]*$/), Validators.maxLength(500)]),
    })
  }

  // Create my skilss talent form.
  createMySkillsForm(): void {
    this.mySkillsForm = this.formBuilder.group({
      mySkills: new FormControl(null, [Validators.pattern(/^[a-zA-Z ]*$/), Validators.maxLength(500)]),
      technologies: new FormControl(Array()),
    })
  }

  // Create new academics form.
  createAcademicsForm(): void {
    this.academicsForm = this.formBuilder.group({
      academics: this.formBuilder.array([])
    })
  }

  // Add new academic to academics form.
  addAcademic(): void {
    this.academicControlArray.push(this.formBuilder.group({
      degree: new FormControl(null, [Validators.required]),
      specialization: new FormControl(null, [Validators.required]),
      college: new FormControl(null, [Validators.required]),
      percentage: new FormControl(null, [Validators.min(1), Validators.pattern(/(^100([.]0{1,2})?)$|(^\d{1,2}([.]\d{1,2})?)$/i), Validators.required]),
      passout: new FormControl(null, [Validators.required, Validators.min(1980), Validators.max(this.currentYear + 3)])
    }))
  }

  // Create new experiences form.
  createExperiencesForm(): void {
    this.experiencesForm = this.formBuilder.group({
      experiences: this.formBuilder.array([])
    })
  }

  // Add new experience to experiences form.
  addExperience(): void {
    this.experienceControlArray.push(this.formBuilder.group({
      company: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z ]*$/)]),
      designation: new FormControl(null, [Validators.required]),
      technologies: new FormControl(null),
      isCurrentWorking: new FormControl(false),
      fromDate: new FormControl(null, [Validators.required]),
      toDate: new FormControl(null, [Validators.required]),
      package: new FormControl(null, [Validators.min(100000), Validators.max(100000000)]),
    }))
  }

  // Create new masters abroad form.
  createMastersAbroadForm(): void {
    this.mastersAbroadForm = this.formBuilder.group({
      id: new FormControl(null),
      degree: new FormControl(null, [Validators.required]),
      countries: new FormControl(null, [Validators.required]),
      universities: new FormControl(Array(), [Validators.required]),
      yearOfMS: new FormControl(null),
      greID: new FormControl(null),
      gre: new FormControl(null, [Validators.required, Validators.min(0)]),
      gmatID: new FormControl(null),
      gmat: new FormControl(null, [Validators.min(0)]),
      toeflID: new FormControl(null),
      toefl: new FormControl(null, [Validators.min(0)]),
      ieltsID: new FormControl(null),
      ielts: new FormControl(null, [Validators.min(0), Validators.pattern(/^(10|\d)(\.5)?$/)]),
      talentID: new FormControl(null),
      enquiryID: new FormControl(null)
    })
  }

  // ======================================= VERIFY PASSWORD FUNCTIONS FOR CHANGE PASSWORD TAB ==========================================

  // Reset verify password form.
  resetVerifyPasswordForm(): void {
    this.verifyPasswordForm.reset()
  }

  // On clicking the verify password button.
  onVerifyPasswordButtonClick(): void {
    this.spinnerService.loadingMessage = "Verifying Password"


    let passwordChange: IPassswordChange = {
      password: this.verifyPasswordForm.get('password').value,
      roleID: this.roleID,
      email: this.email
    }
    this.accountSettingsService.verifyPassword(passwordChange).subscribe((data) => {
      this.forOpeningChangePasswordFormModal()
      this.createVerifyPasswordForm()
    }, (error) => {
      console.error(error)
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

  // ======================================= CHANGE PASSWORD FUNCTIONS FOR CHANGE PASSWORD TAB ==========================================

  // For opening the change password form modal.
  forOpeningChangePasswordFormModal(): void {
    this.createChangePasswordForm()
    this.openModal(this.changePasswordFormModal, 'md')
  }

  // On clicking sumbit button in change password form.
  onChangePasswordFormSubmit(): void {
    if (this.changePasswordForm.invalid) {
      this.changePasswordForm.markAllAsTouched()
      return
    }
    if (this.changePasswordForm.get('password').value != this.changePasswordForm.get('confirmPassword').value) {
      alert("Passwords do not match")
      // this.noMatch = true
      return
    }
    this.changePassword()
  }

  // Change password.
  changePassword(): void {
    this.spinnerService.loadingMessage = "Changing Password"


    let passwordChange: IPassswordChange = {
      password: this.changePasswordForm.get('password').value,
      roleID: this.roleID,
      email: this.email
    }
    this.accountSettingsService.changePassword(passwordChange).subscribe((response: any) => {
      alert(response)
      this.modalRef.close()
    }, (error) => {
      console.error(error)
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

  // ======================================= FUNCTIONS FOR PROFILE TAB ==========================================

  //********************BASIC DETAILS RELATED FUNCTIONS********************* */

  // Get data from talent and prepopulate the basic details form.
  prepopulateBasicDetailsForm(): void {
    this.isBasicDetailsFormUpdated = false
    this.basicDetailsForm.markAsPristine()
    this.basicDetailsForm.patchValue(this.talent)
    this.basicDetailsForm.disable()
  }

  // On clicking basic details edit button.
  onBasicDetailsFormEditClick(): void {
    this.basicDetailsForm.enable()
    if (!this.talent.country) {
      this.basicDetailsForm.get('state').disable()
    }
  }

  // On clicking basic details save button.
  onBasicDetailsFormSaveClick(): void {

    this.isBasicDetailsFormUpdated = true

    // If form is invalid mark all fields as touched.
    if (this.basicDetailsForm.invalid) {
      this.basicDetailsForm.markAllAsTouched()
      return
    }

    // Fill the basic details of talent.
    this.talent.firstName = this.basicDetailsForm.get('firstName').value
    this.talent.lastName = this.basicDetailsForm.get('lastName').value
    this.talent.contact = this.basicDetailsForm.get('contact').value
    this.talent.country = this.basicDetailsForm.get('country').value
    this.talent.state = this.basicDetailsForm.get('state').value
    this.talent.city = this.basicDetailsForm.get('city').value
    this.talent.address = this.basicDetailsForm.get('address').value
    this.talent.pinCode = this.basicDetailsForm.get('pinCode').value
    this.talent.academicYear = this.basicDetailsForm.get('academicYear').value

    // Extract ID from objects in talent.
    this.patchIDFromObjectsForTalent()

    // Update talent.
    this.updateTalent()
  }

  //********************SOCIAL MEDIA RELATED FUNCTIONS********************* */

  // Get data from talent and prepopulate the social media form.
  prepopulateSocialMediaForm(): void {
    this.isSocialMediaFormUpdated = false
    this.socialMediaForm.markAsPristine()
    this.socialMediaForm.patchValue(this.talent)
    this.socialMediaForm.disable()
  }

  // On clicking social media edit button.
  onSocialMediaFormEditClick(): void {
    this.socialMediaForm.enable()
  }

  // On clicking social media save button.
  onSocialMediaFormSaveClick(): void {

    this.isSocialMediaFormUpdated = true

    // If form is invalid mark all fields as touched.
    if (this.socialMediaForm.invalid) {
      this.socialMediaForm.markAllAsTouched()
      return
    }

    // Fill the social media of talent.
    this.talent.facebookUrl = this.socialMediaForm.get('facebookUrl').value
    this.talent.instagramUrl = this.socialMediaForm.get('instagramUrl').value
    this.talent.githubUrl = this.socialMediaForm.get('githubUrl').value
    this.talent.linkedInUrl = this.socialMediaForm.get('linkedInUrl').value

    // Extract ID from objects in talent.
    this.patchIDFromObjectsForTalent()

    // Update talent.
    this.updateTalent()
  }

  //********************ABOUT ME RELATED FUNCTIONS********************* */

  // Get data from talent and prepopulate the about me form.
  prepopulateAboutMeForm(): void {
    this.isAboutMeFormUpdated = false
    this.aboutMeForm.markAsPristine()
    this.aboutMeForm.patchValue(this.talent)
    this.aboutMeForm.disable()
  }

  // On clicking about me edit button.
  onAboutMeFormEditClick(): void {
    this.aboutMeForm.enable()
  }

  // On clicking about me save button.
  onAboutMeFormSaveClick(): void {

    this.isAboutMeFormUpdated = true

    // If form is invalid mark all fields as touched.
    if (this.aboutMeForm.invalid) {
      this.aboutMeForm.markAllAsTouched()
      return
    }

    // Fill the about me of talent.
    this.talent.about = this.aboutMeForm.get('about').value

    // Extract ID from objects in talent.
    this.patchIDFromObjectsForTalent()

    // Update talent.
    this.updateTalent()
  }

  //********************MY SKILLS RELATED FUNCTIONS********************* */

  // Get data from talent and prepopulate the my skills form.
  prepopulateMySkillsForm(): void {
    this.isMySkillsUpdated = false
    this.mySkillsForm.markAsPristine()
    this.mySkillsForm.patchValue(this.talent)
    this.mySkillsForm.disable()
  }

  // On clicking my skills edit button.
  onMySkillsFormEditClick(): void {
    this.mySkillsForm.enable()
  }

  // On clicking my skills save button.
  onMySkillsFormSaveClick(): void {

    this.isMySkillsUpdated = true

    // If form is invalid mark all fields as touched.
    if (this.mySkillsForm.invalid) {
      this.mySkillsForm.markAllAsTouched()
      return
    }

    // Fill the my skills of talent.
    this.talent.about = this.mySkillsForm.get('mySkills').value
    this.talent.technologies = this.mySkillsForm.get('technologies').value

    // Extract ID from objects in talent.
    this.patchIDFromObjectsForTalent()

    // Update talent.
    this.updateTalent()
  }

  //********************ACADEMICS RELATED FUNCTIONS********************* */

  // Get data from talent and prepopulate the academics form.
  prepopulateAcademicsForm(): void {
    this.isAcademicsFormUpdated = false
    this.academicsForm.markAsPristine()
    this.academicsForm.setControl("academics", this.formBuilder.array([]))
    if (this.talent.academics && this.talent.academics.length > 0) {
      for (let i = 0; i < this.talent.academics.length; i++) {
        this.addAcademic()
        this.getSpecializationListByDegreeID(this.talent.academics[i].degree.id)
      }
    }
    this.academicsForm.get('academics').patchValue(this.talent.academics)
    this.academicsForm.disable()
  }

  // On clicking academics edit button.
  onAcademicsFormEditClick(): void {
    this.academicsForm.enable()
  }

  // On clicking academics save button.
  onAacdemicsFormSaveClick(): void {

    this.isAcademicsFormUpdated = true

    // If form is invalid mark all fields as touched.
    if (this.academicsForm.invalid) {
      this.academicsForm.markAllAsTouched()
      return
    }

    // Fill the academics of talent.
    this.talent.academics = this.academicsForm.get('academics').value

    // Extract ID from objects in talent.
    this.patchIDFromObjectsForTalent()

    // Update talent.
    this.updateTalent()
  }

  //********************EXPERIENCES RELATED FUNCTIONS********************* */

  // Get data from talent and prepopulate the experiences form.
  prepopulateExperiencesForm(): void {
    this.isExperiencesFormUpdated = false
    this.experiencesForm.markAsPristine()
    this.experiencesForm.setControl("experiences", this.formBuilder.array([]))
    if (this.talent.experiences && this.talent.experiences.length > 0) {
      for (let i = 0; i < this.talent.experiences.length; i++) {
        if (this.talent.experiences[i].toDate == null) {
          this.indexOfCurrentWorkingExp = i
        }
        this.addExperience()
      }
    }
    if (this.indexOfCurrentWorkingExp != -1) {
      this.experienceControlArray.at(this.indexOfCurrentWorkingExp).get('isCurrentWorking').setValue('true')
      let control = this.experienceControlArray.controls[this.indexOfCurrentWorkingExp]
      if (control instanceof FormGroup) {
        control.removeControl('toDate')
      }
    }

    this.experiencesForm.get('experiences').patchValue(this.talent.experiences)
    this.experiencesForm.disable()
  }

  // On clicking experiences edit button.
  onExperiencesFormEditClick(): void {
    this.experiencesForm.enable()
  }

  // On clicking experiences save button.
  onExperiencesFormSaveClick(): void {

    this.isExperiencesFormUpdated = true

    // If form is invalid mark all fields as touched.
    if (this.experiencesForm.invalid) {
      this.experiencesForm.markAllAsTouched()
      return
    }

    // Fill the experiences of talent.
    this.talent.experiences = this.experiencesForm.get('experiences').value

    // Extract ID from objects in talent.
    this.patchIDFromObjectsForTalent()

    // Update talent.
    this.updateTalent()
  }

  //********************MASTERS ABROAD RELATED FUNCTIONS********************* */

  // Get data from talent and prepopulate the masters abroad form.
  prepopulateMastersAbroadForm(): void {
    this.showMastersAbroad = false
    this.isMastersAbroadFormUpdated = false
    this.mastersAbroadForm.markAsPristine()

    let greScore: number = null
    let greID: string = null
    let gmatScore: any = null
    let gmatID: string = null
    let toeflScore: any = null
    let toeflID: string = null
    let ieltsScore: any = null
    let ieltsID: string = null

    if (this.talent.mastersAbroad) {
      for (let i = 0; i < this.talent.mastersAbroad.scores.length; i++) {
        if (this.talent.mastersAbroad.scores[i].examination.name == "GRE") {
          greScore = this.talent.mastersAbroad.scores[i].marksObtained
          greID = this.talent.mastersAbroad.scores[i].id
        }
        if (this.talent.mastersAbroad.scores[i].examination.name == "GMAT") {
          gmatScore = this.talent.mastersAbroad.scores[i].marksObtained
          gmatID = this.talent.mastersAbroad.scores[i].id
        }
        if (this.talent.mastersAbroad.scores[i].examination.name == "TOEFL") {
          toeflScore = this.talent.mastersAbroad.scores[i].marksObtained
          toeflID = this.talent.mastersAbroad.scores[i].id
        }
        if (this.talent.mastersAbroad.scores[i].examination.name == "IELTS") {
          ieltsScore = this.talent.mastersAbroad.scores[i].marksObtained
          ieltsID = this.talent.mastersAbroad.scores[i].id
        }
      }
      let tempMastersAbroad: any = {
        id: this.talent.mastersAbroad.id,
        degree: this.talent.mastersAbroad.degree,
        countries: this.talent.mastersAbroad.countries,
        universities: this.talent.mastersAbroad.universities,
        yearOfMS: this.talent.mastersAbroad.yearOfMS,
        gre: greScore,
        greID: greID,
        gmat: gmatScore,
        gmatID: gmatID,
        toefl: toeflScore,
        toeflID: toeflID,
        ielts: ieltsScore,
        ieltsID: ieltsID
      }
      this.mastersAbroadForm.patchValue(tempMastersAbroad)
      this.showMastersAbroad = true
    }
    this.mastersAbroadForm.disable()
  }

  // On clicking masters abroad edit button.
  onMastersAbroadFormEditClick(): void {
    this.showMastersAbroad = true
    this.mastersAbroadForm.enable()
  }

  // On clicking masters abroad save button.
  onMastersAbroadFormSaveClick(): void {

    this.isMastersAbroadFormUpdated = true

    // If form is invalid mark all fields as touched.
    if (this.mastersAbroadForm.invalid) {
      this.mastersAbroadForm.markAllAsTouched()
      return
    }

    // Fill the masters abroad of talent.
    if (!this.showMastersAbroad) {
      this.talent.mastersAbroad = null
      this.talent.isMastersAbroad = false
    }
    if (this.showMastersAbroad) {
      let scoresArray: any[] = []
      let mastersAbroad: any
      if (this.mastersAbroadForm.get('gre').value != null) {
        let score: any = {
          id: this.mastersAbroadForm.get('greID').value,
          marksObtained: this.mastersAbroadForm.get('gre').value,
          examinationID: this.greExam.id,
          mastersAbroadID: this.mastersAbroadForm.get('id').value
        }
        scoresArray.push(score)
      }
      if (this.mastersAbroadForm.get('gmat').value != null) {
        let score: any = {
          id: this.mastersAbroadForm.get('gmatID').value,
          marksObtained: this.mastersAbroadForm.get('gmat').value,
          examinationID: this.gmatExam.id,
          mastersAbroadID: this.mastersAbroadForm.get('id').value
        }
        scoresArray.push(score)
      }
      if (this.mastersAbroadForm.get('toefl').value != null) {
        let score: any = {
          id: this.mastersAbroadForm.get('toeflID').value,
          marksObtained: this.mastersAbroadForm.get('toefl').value,
          examinationID: this.toeflExam.id,
          mastersAbroadID: this.mastersAbroadForm.get('id').value
        }
        scoresArray.push(score)
      }
      if (this.mastersAbroadForm.get('ielts').value != null) {
        let score: any = {
          id: this.mastersAbroadForm.get('ieltsID').value,
          marksObtained: this.mastersAbroadForm.get('ielts').value,
          examinationID: this.ieltsExam.id,
          mastersAbroadID: this.mastersAbroadForm.get('id').value
        }
        scoresArray.push(score)
      }
      mastersAbroad = {
        scores: scoresArray,
        degreeID: this.mastersAbroadForm.get('degree').value.id,
        countries: this.mastersAbroadForm.get('countries').value,
        universities: this.mastersAbroadForm.get('universities').value,
        yearOfMS: this.mastersAbroadForm.get('yearOfMS').value,
        talentID: this.talent.id
      }
      if (this.talent.mastersAbroad) {
        mastersAbroad.id = this.talent.mastersAbroad.id
      }
      this.talent.mastersAbroad = mastersAbroad
    }

    // Extract ID from objects in talent.
    this.patchIDFromObjectsForTalent()

    // Update talent.
    this.updateTalent()
  }

  //*****************************OTHER FUNCTIONS******************************* */

  // Get data from talent and prepopulate all profile forms.
  prepolulteProfileForms(): void {
    if (this.isBasicDetailsFormUpdated) {
      this.prepopulateBasicDetailsForm()
    }
    if (this.isSocialMediaFormUpdated) {
      this.prepopulateSocialMediaForm()
    }
    if (this.isAboutMeFormUpdated) {
      this.prepopulateAboutMeForm()
    }
    if (this.isMySkillsUpdated) {
      this.prepopulateMySkillsForm()
    }
    if (this.isAcademicsFormUpdated) {
      this.prepopulateAcademicsForm()
    }
    if (this.isExperiencesFormUpdated) {
      this.prepopulateExperiencesForm()
    }
    if (this.isMastersAbroadFormUpdated) {
      this.prepopulateMastersAbroadForm()
    }
  }

  // Extract ID from objects in talent form.
  patchIDFromObjectsForTalent(): void {

    // Talent source.
    if (this.talent.talentSource) {
      this.talent.sourceID = this.talent.talentSource.id
      delete this.talent.talentSource
    }

    // Salesperson.
    if (this.talent.salesPerson) {
      this.talent.salesPersonID = this.talent.salesPerson.id
      delete this.talent.salesPerson
    }

    // Designation.
    for (let i = 0; i < this.talent.experiences.length; i++) {
      if (this.talent.experiences[i].designation) {
        this.talent.experiences[i].designationID = this.talent.experiences[i].designation.id
        delete this.talent.experiences[i].designation
      }
    }

    // Is experience.
    if (this.talent.experiences?.length == 0) {
      this.talent.isExperience = false
    }
    if (this.talent.experiences?.length > 0) {
      this.talent.isExperience = true
    }

    // Degree and specialization.
    for (let i = 0; i < this.talent.academics.length; i++) {
      if (this.talent.academics[i].degree) {
        this.talent.academics[i].degreeID = this.talent.academics[i].degree.id
        delete this.talent.academics[i].degree
      }
      if (this.talent.academics[i].specialization) {
        this.talent.academics[i].specializationID = this.talent.academics[i].specialization.id
        delete this.talent.academics[i].specialization
      }
    }

    // Experience in months.
    this.calculateTotalYearsOfExperinceOfTalent()
  }

  // Calculate total years of experience for talent in number before add or update of talent.
  calculateTotalYearsOfExperinceOfTalent(): string {
    let monthDiff: number = 0
    if ((this.talent.experiences) && (this.talent.experiences.length != 0)) {
      for (let i = 0; i < this.talent.experiences.length; i++) {
        monthDiff = this.calculateYearsOfExperience(this.talent.experiences[i]) + monthDiff
      }
      this.talent.experienceInMonths = monthDiff
      return
    }
    this.talent.experienceInMonths = monthDiff
    return
  }

  // Calculates the years of experience for each experience before add or update of talent.
  calculateYearsOfExperience(experience: any): number {
    let toDate
    let fromDate
    if (experience.fromDate != null && experience.fromDate != "") {
      if (experience.toDate == null || experience.toDate == "") {
        toDate = new Date()
      }
      else {
        toDate = new Date(experience.toDate)
      }
      fromDate = new Date(experience.fromDate)
      let monthDiff: number = toDate.getMonth() - fromDate.getMonth() + (12 * (toDate.getFullYear() - fromDate.getFullYear()))
      return monthDiff
    }
  }

  // Get academics array from academics form.
  get academicControlArray(): FormArray {
    return this.academicsForm.get("academics") as FormArray
  }

  // Get expereinces array from expereinces form.
  get experienceControlArray(): FormArray {
    return this.experiencesForm.get("experiences") as FormArray
  }

  // To aply validation to college name.
  addCollegeToList(option: any): Promise<any> {
    return new Promise((resolve) => {
      resolve(option)
    })
  }

  // Delete academic from academics form.
  deleteAcademic(index: number): void {
    if (confirm("Are you sure you want to delete this college?")) {
      this.academicsForm.enable()
      this.academicsForm.markAsDirty()
      this.academicControlArray.removeAt(index)
    }
  }

  // Delete experience from experiences form.
  deleteExperience(index: number): void {
    if (confirm("Are you sure you want to delete this experience ?")) {
      if (this.indexOfCurrentWorkingExp == index) {
        this.indexOfCurrentWorkingExp = -1
      }
      if (this.indexOfCurrentWorkingExp > index) {
        this.indexOfCurrentWorkingExp = this.indexOfCurrentWorkingExp - 1
      }
      this.experiencesForm.markAsDirty()
      this.experienceControlArray.removeAt(index)
    }
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

  // On degree change in masters abroad form.
  onDegreeChange(degree: any): void {
    if (degree != null && degree.name == "M.B.A.") {
      this.mastersAbroadForm.get('gmat').setValidators([Validators.required, Validators.min(0), Validators.max(this.gmatExam.totalMarks)])
      this.mastersAbroadForm.get('gmat').updateValueAndValidity()
      this.showMBADegreeRequiredError = true
      return
    }
    this.mastersAbroadForm.get('gmat').setValidators([Validators.min(0), Validators.max(this.gmatExam.totalMarks)])
    this.mastersAbroadForm.get('gmat').updateValueAndValidity()
    this.showMBADegreeRequiredError = false
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

  // Delete masters abroad of talent.
  deleteMastersAbroad(): void {
    if (confirm("Are you sure you want to delete this masters abroad details?")) {
      this.showMastersAbroad = false
      this.mastersAbroadForm.markAsDirty()
    }
  }

  // Add masters abroad of talent.
  addMastersAbroad(): void {
    this.createMastersAbroadForm()
    this.mastersAbroadForm.get('gre').setValidators([Validators.required, Validators.min(0), Validators.max(this.greExam?.totalMarks)])
    this.mastersAbroadForm.get('gre').updateValueAndValidity()

    this.mastersAbroadForm.get('gmat').setValidators([Validators.min(0), Validators.max(this.gmatExam.totalMarks)])
    this.mastersAbroadForm.get('gmat').updateValueAndValidity()

    this.mastersAbroadForm.get('toefl').setValidators([Validators.min(0), Validators.max(this.toeflExam.totalMarks)])
    this.mastersAbroadForm.get('toefl').updateValueAndValidity()

    this.mastersAbroadForm.get('ielts').setValidators([Validators.min(0), Validators.pattern(/^(10|\d)(\.\d{1,2})?$/),
    Validators.max(this.ieltsExam.totalMarks)])
    this.mastersAbroadForm.get('ielts').updateValueAndValidity()
    this.showMastersAbroad = true
  }

  // Update talent.
  updateTalent(): void {
    this.spinnerService.loadingMessage = "Updating Talent"


    this.talentService.updateTalent(this.talent).subscribe((response: any) => {
      this.getTalent()
      alert(response)
    }, (error) => {
      console.error(error)
      if (error.error?.error) {
        alert(error.error.error)
        return
      }
      alert("Check connection")
    })
  }

  // ======================================= OTHER FUNCTIONS FOR ACCOUNT SETTINGS ==========================================

  // On changing tab get the component lists.
  onTabChange(event: any) {
    if (event == 1) {
      this.isBasicDetailsFormUpdated = true
      this.isSocialMediaFormUpdated = true
      this.isAboutMeFormUpdated = true
      this.isMySkillsUpdated = true
      this.isAcademicsFormUpdated = true
      this.isExperiencesFormUpdated = true
      this.isMastersAbroadFormUpdated = true
      this.getTalent()
    }
    if (event == 2) {
      this.createVerifyPasswordForm()
    }
  }

  // Used to open modal.
  openModal(content: any, size?: string): NgbModalRef {
    if (!size) {
      size = 'lg'
    }
    let options: NgbModalOptions = {
      ariaLabelledBy: 'modal-basic-title', keyboard: false,
      backdrop: 'static', size: size
    }
    this.modalRef = this.modalService.open(content, options)
    return this.modalRef
  }

  // ======================================= GET FUNCTIONS ======================================================================

  // Populate all lists.
  getAllComponent(): void {
    this.getCountryList()
    this.getDegreeList()
    this.getCollegeBranchList()
    this.getAcademicYearList()
    this.getTechnologyList()
    this.getDesignationList()
    this.getExaminationList()
    this.getYearOfMSList()
  }

  // Get talent by id.
  getTalent(): void {
    this.spinnerService.loadingMessage = "Getting Profile"


    this.talentService.getTalent(this.talentID).subscribe(response => {
      this.talent = response
      if (this.talent.country != undefined) {
        this.getStateListByCountryID(this.talent.country)
      }
      if (this.talent.experiences && this.talent.experiences.length > 0) {
        for (let i = 0; i < this.talent.experiences.length; i++) {
          if (this.talent.experiences[i].toDate == null) {
            this.indexOfCurrentWorkingExp = i
          }
          let fromDate = this.talent.experiences[i]?.fromDate
          if (fromDate) {
            this.talent.experiences[i].fromDate = this.datePipe.transform(fromDate, 'yyyy-MM')
          }
          let toDate = this.talent.experiences[i]?.toDate
          if (toDate) {
            this.talent.experiences[i].toDate = this.datePipe.transform(toDate, 'yyyy-MM')
          }
          this.addExperience()
        }
      }
      this.prepolulteProfileForms()
    }, error => {
      console.error(error)
    })
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
  getSpecializationListByDegreeID(degreeID: string, index?: number, shouldClear?: boolean): void {
    if (index != undefined) {
      this.academicControlArray.at(index).get('specialization').setValue(null)
    }
    if (this.degreeIDSet.has(degreeID)) {
      return
    }
    this.degreeIDSet.add(degreeID)
    if (degreeID) {
      this.generalService.getSpecializationByDegreeID(degreeID).subscribe((response: any) => {
        this.specializationList = this.specializationList.concat(response.body)
      }, (err) => {
        console.error(err)
      })
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
      this.basicDetailsForm.get('state').setValue(null)
      this.basicDetailsForm.get('state').disable()
      return
    }
    if (country.name == "U.S.A." || country.name == "India") {
      this.basicDetailsForm.get('state').enable()
      this.generalService.getStatesByCountryID(country.id).subscribe((respond: any[]) => {
        this.stateList = respond;
      }, (err) => {
        console.error(err)
      });
    }
    else {
      this.basicDetailsForm.get('state').setValue(null)
      this.basicDetailsForm.get('state').disable()
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

  // Get examination list.
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

    this.mastersAbroadForm.get('gre').setValidators([Validators.required, Validators.min(0), Validators.max(this.greExam?.totalMarks)])
    this.mastersAbroadForm.get('gre').updateValueAndValidity()

    this.mastersAbroadForm.get('gmat').setValidators([Validators.min(0), Validators.max(this.gmatExam.totalMarks)])
    this.mastersAbroadForm.get('gmat').updateValueAndValidity()

    this.mastersAbroadForm.get('toefl').setValidators([Validators.min(0), Validators.max(this.toeflExam.totalMarks)])
    this.mastersAbroadForm.get('toefl').updateValueAndValidity()

    this.mastersAbroadForm.get('ielts').setValidators([Validators.min(0), Validators.pattern(/^(10|\d)(\.\d{1,2})?$/),
    Validators.max(this.ieltsExam.totalMarks)])
    this.mastersAbroadForm.get('ielts').updateValueAndValidity()
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

  // Create list of years.
  getYearOfMSList(): void {
    this.yearOfMSList.push(this.currentYear)
    this.yearOfMSList.push((this.currentYear + 1))
    this.yearOfMSList.push((this.currentYear + 2))
  }

}

// Interface for examination
export interface IExamination {
  id?: string
  name: string
  totalMarks: number
}
