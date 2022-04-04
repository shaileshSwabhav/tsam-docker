import { Injectable } from "@angular/core";
import { environment } from "src/environments/environment";

@Injectable({
      providedIn: 'root'
})
export class Constant {
      readonly NAVIGATION: string = "navigation"
      readonly ROLE: string = "Role"
      readonly READ: string = "Read"
      readonly WRITE: string = "Write"
      readonly EDIT: string = "Edit"
      readonly DELETE: string = "Delete"

      readonly BASE_URL: string = environment.BASE_URL
      readonly SERVER_URL: string = "https://swabhav-tsam.herokuapp.com/api/v1/tsam"
      readonly LOCAL_URL: string = "http://127.0.0.1:8080/api/v1/tsam"

      readonly TENANT_ID = "7ca2664b-f379-43db-bdf9-7fdd40707219"
      readonly NIL_UUID: string = "00000000-0000-0000-0000-000000000000"
}

@Injectable({
      providedIn: 'root'
})
export class Role {
      readonly SALES_PERSON: string = "SalesPerson"
      readonly FACULTY: string = "Faculty"
      readonly ADMIN: string = "Admin"
      readonly TALENT: string = "Talent"
      readonly COLLEGE: string = "College"
      readonly COMPANY: string = "Company"
      readonly DEVELOPER: string = "Developer"
      readonly MARKETING_EXECUTIVE: string = "Marketing Executive"
      readonly UI_UX_DESIGNER: string = "UI/UX Designer"
}

@Injectable({
      providedIn: 'root'
})
export class UrlConstant {

      // TALENT
      readonly TALENT_MASTER = "/talent/master"
      readonly MY_TALENT = "/my/talent"
      readonly TALENT_ENQUIRY = "/talent/enquiry"
      readonly INTERVIEW_SCHEDULE = "/talent/master/interview-schedule"
      readonly TALENT_TARGET_COMMUNITY = "/talent/target-community"
      readonly TALENT_CAREER_OBJECTVE = "/talent/career-objective"
      readonly TALENT_DEGREE = "/talent/degree"
      readonly TALENT_DESIGNATION = "/talent/designation"
      readonly TALENT_UNIVERSITY = "/talent/university"
      readonly TALENT_REPORT_CALLING = "/talent/report/calling"
      readonly TALENT_NEXT_ACTION = "/talent/report/next-action"
      readonly TALENT_REPORT_LIFETIME = "/talent/report/lifetime"
      readonly TALENT_REPORT_WAITING_LIST = "/talent/report/waiting-list"
      readonly TALENT_REPORT_PRO_SUMMARY = "/talent/report/professional-summary"
      readonly TALENT_REPORT_FRESHER_SUMMARY = "/talent/report/fresher-summary"
      readonly TALENT_REPORT_PACKAGE_SUMMARY = "/talent/report/package-summary"

      // COURSE
      readonly BANK_COURSE = "/bank/course"
      readonly BANK_COURSE_DETAILS = "/bank/course/details"
      readonly BANK_COURSE_PROGRAMMING_ASSIGNMENT = "/bank/course/programming-assignment"
      readonly TRAINING_COURSE_MASTER = "/training/course/master"
      readonly TRAINING_COURSE_MASTER_DETAILS = "/training/course/master/details"
      readonly TRAINING_COURSE_MASTER_PROGRAMMING_ASSIGNMENT = "/training/course/master/programming-assignment"

      // MODULE
      readonly TRAINING_MODULE_CONCEPT_TREE = "/training/module/module-concept"
      readonly TRAINING_MODULE = "/training/module"
      readonly BANK_MODULE_CONCEPT_TREE = "/bank/module/module-concept"
      readonly BANK_MODULE = "/bank/module"

      // BATCH
      readonly MY_BATCH = "/my/batch"
      readonly MY_BATCHES = "/my-batches"
      readonly MY_BATCH_SESSION = "/my/batch/session"
      readonly MY_BATCH_SESSION_DETAILS = "/my/batch/session/details"
      readonly MY_BATCH_ASSIGNMENT_DETAILS = "/my/batch/assignment-details"
      readonly MY_BATCH_PROJECT_DETAILS = "/my/batch/project-details"
      readonly MY_BATCH_SESSION_FEEDBACK = "/my/batch/session/feedback"
      readonly MY_BATCH_FEEDBACK = "/my/batch/feedback"
      readonly TRAINING_BATCH_MASTER = "/training/batch/master"
      readonly TRAINING_BATCH_MASTER_SESSION = "/training/batch/master/session"
      readonly TRAINING_BATCH_MASTER_SESSION_DETAILS = "/training/batch/master/session/details"
      readonly TRAINING_BATCH_MASTER_ASSIGNMENT_DETAILS = "/training/batch/master/assignment-details"
      readonly TRAINING_BATCH_MASTER_PROJECT_DETAILS = "/training/batch/master/project-details"
      readonly TRAINING_BATCH_MASTER_SESSION_FEEDBACK = "/training/batch/master/session/feedback"
      readonly TRAINING_BATCH_MASTER_FEEDBACK = "/training/batch/master/feedback"

      // COLLEGE
      readonly COLLEGE_MASTER = "/sourcing/college/master"
      readonly COLLEGE_BRANCH = "/sourcing/college/branch"
      readonly COLLEGE_CAMPUS_DRIVE = "/sourcing/college/campus-drive"
      readonly COLLEGE_SEMINAR = "/sourcing/college/seminar"
      readonly COLLEGE_WORKSHOP = "/sourcing/college/workshop"

      // COMPANY
      readonly COMPANY_MASTER = "/sourcing/company/master"
      readonly COMPANY_BRANCH = "/sourcing/company/branch"
      readonly COMPANY_ENQUIRY = "/sourcing/company/enquiry"
      readonly COMPANY_REQUIREMENT = "/sourcing/company/requirement"
      readonly COMPANY_SALARY_TREND = "/sourcing/company/salary-trend"

      // ADMIN
      readonly ADMIN_ALL_EMPLOYEE = "/admin/all-employee"
      readonly ADMIN_OTHER_EMPLOYEE = "/admin/other-employee"
      readonly ADMIN_FACULTY = "/admin/faculty"
      readonly ADMIN_DEPARTMENT = "/admin/department"
      readonly ADMIN_EMPLOYEE_PROJECT = "/admin/employee-project"
      readonly ADMIN_USER = "/admin/user"
      readonly ADMIN_LOGIN_REPORT = "/admin/report/login"

      // TRAINING
      readonly TRAINING_RESOURCE = "/training/resource"
      readonly TRAINING_PROJECT = "/training/project"
      readonly TRAINING_TECHNOLOGY = "/training/technology"
      readonly TRAINING_INPUT_LANGUAGE = "/training/input-language"
      readonly TRAINING_FACULTY = "/training/faculty"
      readonly TRAINING_PROGRAMMING_QUESTION_TYPE = "/training/programming-question-type"

      // BANK
      readonly BANK_RESOURCE = "/bank/resource"

      // MAIN MENU
      readonly RESOURCE = "/resource"
      readonly DASHBOARD = "/dashboard"
      readonly FACULTY_DASHBOARD = "/faculty-dashboard"
      readonly TALENT_DASHBOARD = "/talent-dashboard"

      // FEEDBACK
      readonly TRAINING_FEEDBACK = "/training/feedback-question"
      readonly TRAINING_GROUP = "/training/feedback-question-group"

      // PROGRAMMING QUESTION
      readonly TRAINING_PROGRAMMING_QUESTION = "/training/coding-question"
      readonly BANK_PROGRAMMING_QUESTION = "/bank/coding-question"

      // PROGRAMMING SOLUTION
      readonly ADMIN_PROGRAMMING_QUESTION_TALENT_ANSWER = "/admin/coding-problems/answers"
      readonly ADMIN_PROGRAMMING_QUESTION_TALENT_ANSWER_DETAILS = "/admin/coding-problems/answers/solution-details"
      readonly PROGRAMMING_QUESTION_TALENT_ANSWER = "/coding-problems/answers"
      readonly PROGRAMMING_QUESTION_TALENT_ANSWER_DETAILS = "/coding-problems/answers/solution-details"

      // PROGRAMMING ASSIGNMENT
      readonly ADMIN_PROGRAMMING_ASSIGNMENT = "/admin/coding-problems/assignment"
      readonly FACULTY_PROGRAMMING_ASSIGNMENT = "/coding-problems/assignment"

      // PROFILE
      readonly TIMESHEET = "/timesheet"
      readonly MY_TIMESHEET = "/my/timesheet"
      readonly ACCOUNT_SETTINGS = "/profile/account-settings"

      // CONCEPT
      readonly TRAINING_PROGRAMMING_CONCEPT = "/training/concept"
      readonly BANK_PROGRAMMING_CONCEPT = "/bank/concept"

      // FILE UPLOADS
      readonly UNIVERSITY_DEMO = `${environment.FILE_UPLOAD_LOACTION}/demos/universityDemo.xlsx`
      readonly COLLEGE_DEMO = `${environment.FILE_UPLOAD_LOACTION}/demos/collegeDemo.xlsx`
      readonly TALENT_BASIC_DEMO = `${environment.FILE_UPLOAD_LOACTION}/demos/talentDemo.xlsx`
      readonly ENQUIRY_BASIC_DEMO = `${environment.FILE_UPLOAD_LOACTION}/demos/enquiryDemo.xlsx`

      // NOT IN USE!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
      readonly COURSE_SESSION = "/course/master/session"
      readonly COURSE_MODULE = "/course/master/module"
      readonly TEST_MASTER = "/test-master"
}

@Injectable({
      providedIn: 'root'
})
export class UploadConstant {
      readonly ALLOWED_DOCUMENT_EXTENSION = ".pdf, .doc, .docx, .ppt, .pptx, .txt"
      readonly ALLOWED_AUDIO_EXTENSION = ".mp3, .m4a, .wav"
      readonly ALLOWED_IMAGE_EXTENSION = ".png, .jpeg, .gif"
}

@Injectable({
      providedIn: 'root'
})
export class FeedbackType {
      readonly TALENT_BATCH_FEEDBACK = "Talent_Batch_Feedback"
      readonly TALENT_SESSION_FEEDBACK = "Talent_Session_Feedback"
      readonly FACULTY_SESSION_FEEDBACK = "Faculty_Session_Feedback"
      readonly FACULTY_BATCH_FEEDBACK = "Faculty_Batch_Feedback"
      readonly AHA_MOMENT_FEEDBACK = "Aha_Moment_Feedback"
      readonly FACULTY_ASSESSMENT = "Faculty_Assessment"
}

@Injectable({
      providedIn: 'root'
})
export class UploadFolderPath {

}

@Injectable({
      providedIn: 'root'
})
export class BackNavigationUrl {
      // Please maintain this in alphabetic order.
      readonly BTA_DETAILS_TO_BTA_SCORES = "bta-details-bta-scores"
      readonly BP_DETAILS_TO_BP_SCORES = "bp-details-bp-scores"
}

@Injectable({
      providedIn: 'root'
})
export class AccessLevel {
      readonly ONLY_TALENT = {
            isTalent: true,
            isFaculty: false,
            isAdmin: false,
            isSalesPerson: false,
      }
      readonly ONLY_FACULTY = {
            isTalent: false,
            isFaculty: true,
            isAdmin: false,
            isSalesPerson: false,
      }
      readonly ONLY_ADMIN = {
            isTalent: false,
            isFaculty: false,
            isAdmin: true,
            isSalesPerson: false,
      }
      readonly ADMIN_AND_SALESPERSON = {
            isTalent: false,
            isFaculty: false,
            isAdmin: true,
            isSalesPerson: true,
      }
}   