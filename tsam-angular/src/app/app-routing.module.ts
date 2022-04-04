import { NgModule } from '@angular/core';
import { Routes, RouterModule, ExtraOptions } from '@angular/router';
import { LoginComponent } from './component/login/login.component';
import { DashboardComponent } from './component/dashboard/dashboard.component';
import { TestMasterComponent } from './component/test-master/test-master.component';
import { BatchMasterComponent } from './component/batch-master/batch-master.component';
import { TalentMasterComponent } from './component/talent-master/talent-master.component';
import { InterviewScheduleComponent } from './component/interview-schedule/interview-schedule.component';
import { CourseMasterComponent } from './component/course-master/course-master.component';
import { PageNotFoundComponent } from './component/page-not-found/page-not-found.component';
import { FacultyMasterComponent } from './component/faculty-master/faculty-master.component';
import { TalentEnquiryComponent } from './component/talent-enquiry/talent-enquiry.component';
import { CompanyMasterComponent } from './component/company-master/company-master.component';
import { CompanyEnquiryComponent } from './component/company-enquiry/company-enquiry.component';
import { AdminTechnologyComponent } from './component/admin-technology/admin-technology.component';
import { CompanyMasterRequirementComponent } from './component/company-master-requirement/company-master-requirement.component';
import { AdminRoleComponent } from './component/admin-role/admin-role.component';
import { CollegeCampusComponent } from './component/college-campus/college-campus.component';
import { SeminarComponent } from './component/seminar/seminar.component';
import { RouteGuard } from './providers/guard/route.guard';
import { CollegeMasterComponent } from './component/college-master/college-master.component';
import { CollegeBranchComponent } from './component/college-branch/college-branch.component';
import { AdminDegreeComponent } from './component/admin-degree/admin-degree.component';
import { AdminDesignationComponent } from './component/admin-designation/admin-designation.component';
import { CompanyBranchComponent } from './component/company-branch/company-branch.component';
import { CourseSessionComponent } from './component/course-session/course-session.component';
import { AdminUniversityComponent } from './component/admin-university/admin-university.component';
import { BatchSessionComponent } from './component/batch-session/batch-session.component';
import { CallingReportsComponent } from './component/calling-reports/calling-reports.component';
import { BatchSessionFeedbackComponent } from './component/batch-session-feedback/batch-session-feedback.component';
import { AdminFeedbackComponent } from './component/admin-feedback/admin-feedback.component';
import { BatchFeedbackComponent } from './component/batch-feedback/batch-feedback.component';
import { TargetCommunityComponent } from './component/target-community/target-community.component';
import { DepartmentComponent } from './component/department/department.component';
import { ProjectComponent } from './component/project/project.component';
import { AdminEmployeeComponent } from './component/admin-employee/admin-employee.component';
import { LifetimeValueReportComponent } from './component/lifetime-value-report/lifetime-value-report.component';
import { NextActionReportComponent } from './component/next-action-report/next-action-report.component';
import { AdminCareerObjectiveComponent } from './component/admin-career-objective/admin-career-objective.component';
import { AdminUserComponent } from './component/admin-user/admin-user.component';
import { AdminManagementComponent } from './component/admin-management/admin-management.component';
import { AdminGroupComponent } from './component/admin-group/admin-group.component';
import { AdminFeelingComponent } from './component/admin-feeling/admin-feeling.component';
import { TimesheetComponent } from './component/timesheet/timesheet.component';
import { AdminResourceComponent } from './component/admin-resource/admin-resource.component';
import { SalaryTrendComponent } from './component/salary-trend/salary-trend.component';
import { WaitingListReportComponent } from './component/waiting-list-report/waiting-list-report.component';
import { LoginReportComponent } from './component/login-report/login-report.component';
import { AdminSpeakerComponent } from './component/admin-speaker/admin-speaker.component';
import { ProSummaryReportComponent } from './component/pro-summary-report/pro-summary-report.component';
import { FresherSummaryComponent } from './component/fresher-summary/fresher-summary.component';
import { LogoutComponent } from './component/logout/logout.component';
import { AccountSettingsComponent } from './component/account-settings/account-settings.component';
import { ComingSoonComponent } from './component/coming-soon/coming-soon.component';
import { MyOpportunitiesComponent } from './component/my-opportunities/my-opportunities.component';
import { CompanyDetailsComponent } from './component/company-details/company-details.component';
import { MyCoursesComponent } from './component/my-courses/my-courses.component';
import { PackageSummaryComponent } from './component/package-summary/package-summary.component';
import { FacultyReportComponent } from './component/faculty-report/faculty-report.component';
import { MySessionsComponent } from './component/my-sessions/my-sessions.component';
import { ProblemOfTheDayComponent } from './component/problem-of-the-day/problem-of-the-day.component';
import { ProblemDetailsComponent } from './component/problem-details/problem-details.component';
import { PracticeComponent } from './component/practice/practice.component';
import { EventsComponent } from './component/events/events.component';
import { EventDetailsComponent } from './component/event-details/event-details.component';
import { TestComponent } from './component/test/test.component';
import { TestDetailsComponent } from './component/test-details/test-details.component';
import { TalentDashboardComponent } from './component/talent-dashboard/talent-dashboard.component';
import { ReferAndEarnComponent } from './component/refer-and-earn/refer-and-earn.component';
import { FacultyDashboardComponent } from './component/faculty-dashboard/faculty-dashboard.component';
import { CourseDetailsComponent } from './component/course-details/course-details.component';
import { ProgrammingQuestionComponent } from './component/programming-question/programming-question.component';
import { AdminProgrammingQuestionTypeComponent } from './component/admin-programming-question-type/admin-programming-question-type.component';
import { AdminEventComponent } from './component/admin-event/admin-event.component';
import { BatchProjectComponent } from './component/batch-project/batch-project.component';
import { ProgrammingQuestionTalentAnswerComponent } from './component/programming-question-talent-answer/programming-question-talent-answer.component';
import { ProgrammingQuestionTalentAnswerDetailsComponent } from './component/programming-question-talent-answer-details/programming-question-talent-answer-details.component';
import { ProgrammingConceptComponent } from './component/programming-concept/programming-concept.component';
import { PracticeProblemDetailsComponent } from './component/practice-problem-details/practice-problem-details.component';
import { BatchDetailsComponent } from './component/batch-details/batch-details.component';
import { CodeEditorComponent } from './component/code-editor/code-editor.component';
import { PracticeProblemListComponent } from './component/practice-problem-list/practice-problem-list.component';
import { ProgrammingLanguageComponent } from './component/programming-language/programming-language.component';
import { CkeEditorComponent } from './component/cke-editor/cke-editor.component';
import { QuestionComponent } from './component/question/question.component';
import { ProgrammingQuestionModalComponent } from './component/programming-question-modal/programming-question-modal.component';
import { BlogTopicComponent } from './component/blog-topic/blog-topic.component';
import { BlogComponent } from './component/blog/blog.component';
import { BlogVerificationComponent } from './component/blog-verification/blog-verification.component';
import { BlogDetailsComponent } from './component/blog-details/blog-details.component';
import { BrowseBlogsComponent } from './component/browse-blogs/browse-blogs.component';
import { CourseProgrammingAssignmentComponent } from './component/course-programming-assignment/course-programming-assignment.component';
import { ProgrammingAssignmentComponent } from './component/programming-assignment/programming-assignment.component';
import { CourseModuleComponent } from './component/course-module/course-module.component';
import { ModuleComponent } from './component/module/module.component';
import { BatchTopicAssignmentDetailsComponent } from './component/batch-topic-assignment-details/batch-topic-assignment-details.component';
import { CourseComponent } from './component/course/course.component';
import { BatchCompletionDetailsComponent } from './component/batch-completion-details/batch-completion-details.component';
import { TalentBatchDetailsComponent } from './component/talent-batch-details/talent-batch-details.component';
import { ProgrammingProjectComponent } from './component/programming-project/programming-project.component';
import { ModuleConceptComponent } from './component/module-concept/module-concept.component';
import { ModuleConceptScoreComponent } from './component/module-concept-score/module-concept-score.component';
import { BatchProjectDetailsComponent } from './component/batch-project-details/batch-project-details.component';
import { FacultyTalentFeedbackComponent } from './component/faculty-talent-feedback/faculty-talent-feedback.component';

const routes: Routes = [

      // *************************** NEW ***************************************************************
      {
            path: "login", component: LoginComponent
      },

      {
            path: 'profile',
            children: [
                  { path: '', redirectTo: 'account-settings', pathMatch: 'full' },
                  { path: 'account-settings', component: AccountSettingsComponent, canActivate: [RouteGuard] },
                  { path: 'logout', component: LogoutComponent },
            ]
      },

      // // ********************************************* FACULTY **************************************************

      {
            path: 'faculty-dashboard', component: FacultyDashboardComponent, canActivate: [RouteGuard],
      },

      {
            path: 'my',
            children: [
                  { path: '', redirectTo: 'timesheet', pathMatch: 'full' },
                  { path: 'timesheet', component: TimesheetComponent, canActivate: [RouteGuard] },
                  { path: 'batch', component: BatchMasterComponent, canActivate: [RouteGuard] },
                  { path: 'talent', component: TalentMasterComponent, canActivate: [RouteGuard] },
            ]
      },

      {
            path: 'bank',
            children: [
                  { path: '', redirectTo: 'coding-question', pathMatch: 'full' },
                  { path: 'coding-question', component: ProgrammingQuestionComponent, canActivate: [RouteGuard] },
                  { path: 'concept', component: ProgrammingConceptComponent, canActivate: [RouteGuard] },
                  { path: 'resource', component: AdminResourceComponent, canActivate: [RouteGuard] },
                  { path: 'module', component: ModuleComponent, canActivate: [RouteGuard] },
                  { path: 'project', component: ProgrammingProjectComponent, canActivate: [RouteGuard] },
                  { path: 'course', component: CourseMasterComponent, canActivate: [RouteGuard] },

            ]
      },
      //feedback component
      {
            path: 'faculty-talent-feedback', component: FacultyTalentFeedbackComponent,
      },
      {
            path: 'my/batch/session/details', component: BatchDetailsComponent, canActivate: [RouteGuard],
      },
      {
            path: 'my/batch/session/details/concept-tree', component: ModuleConceptScoreComponent, canActivate: [RouteGuard],
      },
      {
            path: 'my/batch/session', component: BatchSessionComponent, canActivate: [RouteGuard],
      },
      {
            path: 'my/batch/session/completion', component: BatchCompletionDetailsComponent, canActivate: [RouteGuard],
      },
      {
            path: 'my/batch/session/feedback', component: BatchSessionFeedbackComponent, canActivate: [RouteGuard],
      },
      {
            path: 'my/batch/feedback', component: BatchFeedbackComponent, canActivate: [RouteGuard],
      },
      {
            path: 'my/batch/assignment-details', component: BatchTopicAssignmentDetailsComponent, canActivate: [RouteGuard],
      },
      {
            path: 'my/batch/project-details', component: BatchTopicAssignmentDetailsComponent, canActivate: [RouteGuard],
      },
      {
            path: "bank/course/details", component: CourseComponent, canActivate: [RouteGuard],
      },
      {
            path: 'bank/course/programming-assignment', component: CourseProgrammingAssignmentComponent, canActivate: [RouteGuard]
      },
      {
            path: 'bank/module/module-concept', component: ModuleConceptComponent, canActivate: [RouteGuard]
      },

      // //***********************************************ADMIN AND SALESPERSON***************************************************
      
      {
            path: "dashboard", component: DashboardComponent, canActivate: [RouteGuard]
      },
      
      {
            path: 'talent',
            children: [
                  { path: 'enquiry', component: TalentEnquiryComponent, canActivate: [RouteGuard] },
                  { path: 'target-community', component: TargetCommunityComponent, canActivate: [RouteGuard] },
                  { path: 'master', component: TalentMasterComponent, canActivate: [RouteGuard] },
                  { path: 'career-objective', component: AdminCareerObjectiveComponent, canActivate: [RouteGuard] },
                  { path: 'degree', component: AdminDegreeComponent, canActivate: [RouteGuard] },
                  { path: 'designation', component: AdminDesignationComponent, canActivate: [RouteGuard] },
                  { path: 'university', component: AdminUniversityComponent, canActivate: [RouteGuard] },
                  { path: 'report/calling', component: CallingReportsComponent, canActivate: [RouteGuard] },
                  { path: 'report/next-action', component: NextActionReportComponent, canActivate: [RouteGuard] },
                  { path: 'report/lifetime', component: LifetimeValueReportComponent, canActivate: [RouteGuard] },
                  { path: 'report/waiting-list', component: WaitingListReportComponent, canActivate: [RouteGuard] },
                  { path: 'report/professional-summary', component: ProSummaryReportComponent, canActivate: [RouteGuard] },
                  { path: 'report/fresher-summary', component: FresherSummaryComponent, canActivate: [RouteGuard] },
                  { path: 'report/package-summary', component: PackageSummaryComponent, canActivate: [RouteGuard] },

            ]
      },

      {
            path: 'training',
            children: [
                  { path: 'batch/master', component: BatchMasterComponent, canActivate: [RouteGuard] },
                  { path: 'course/master', component: CourseMasterComponent, canActivate: [RouteGuard] },
                  { path: 'module', component: ModuleComponent, canActivate: [RouteGuard] },
                  { path: 'resource', component: AdminResourceComponent, canActivate: [RouteGuard] },
                  { path: 'project', component: ProgrammingProjectComponent, canActivate: [RouteGuard] },
                  { path: 'coding-question', component: ProgrammingQuestionComponent, canActivate: [RouteGuard] },
                  { path: 'concept', component: ProgrammingConceptComponent, canActivate: [RouteGuard] },
                  { path: 'technology', component: AdminTechnologyComponent, canActivate: [RouteGuard] },
                  { path: 'feedback-question', component: AdminFeedbackComponent, canActivate: [RouteGuard] },
                  { path: 'feedback-question-group', component: AdminGroupComponent, canActivate: [RouteGuard] },
                  { path: 'input-language', component: ProgrammingLanguageComponent, canActivate: [RouteGuard] },
                  { path: 'faculty', component: FacultyReportComponent, canActivate: [RouteGuard] },
            ]
      },

      {
            path: 'sourcing',
            children: [
                  {
                        path: 'company',
                        children: [
                              { path: '', redirectTo: 'enquiry', pathMatch: 'full' },
                              { path: 'enquiry', component: CompanyEnquiryComponent, canActivate: [RouteGuard] },
                              { path: 'master', component: CompanyMasterComponent, canActivate: [RouteGuard] },
                              { path: 'requirement', component: CompanyMasterRequirementComponent, canActivate: [RouteGuard] },
                              { path: 'branch', component: CompanyBranchComponent, canActivate: [RouteGuard] },
                              { path: 'salary-trend', component: SalaryTrendComponent, canActivate: [RouteGuard] },
                        ]
                  },
                  {
                        path :'college',
                        children: [
                              { path: '', redirectTo: 'database', pathMatch: 'full' },
                              { path: 'master', component: CollegeMasterComponent, canActivate: [RouteGuard] },
                              { path: 'campus-drive', component: CollegeCampusComponent, canActivate: [RouteGuard] },
                              { path: 'branch', component: CollegeBranchComponent, canActivate: [RouteGuard] },
                              { path: 'seminar', component: SeminarComponent, canActivate: [RouteGuard] },
                              { path: 'workshop', component: AdminEventComponent, canActivate: [RouteGuard] },

                        ]
                  },
            ]
      },

      {
            path: 'admin',
            children: [
                  { path: 'all-employee', component: AdminManagementComponent, canActivate: [RouteGuard] },
                  { path: 'other-employee', component: AdminEmployeeComponent, canActivate: [RouteGuard] },
                  { path: 'faculty', component: FacultyMasterComponent, canActivate: [RouteGuard] },
                  { path: 'department', component: DepartmentComponent, canActivate: [RouteGuard] },
                  { path: 'employee-project', component: ProjectComponent, canActivate: [RouteGuard] },
                  { path: 'user', component: AdminUserComponent, canActivate: [RouteGuard] },
                  { path: 'report/login', component: LoginReportComponent, canActivate: [RouteGuard] },
            ]
      },
      {
            path: 'training/batch/master/session', component: BatchSessionComponent, canActivate: [RouteGuard],
      },
      {
            path: 'training/batch/master/session/details/concept-tree', component: ModuleConceptScoreComponent, canActivate: [RouteGuard],
      },
      {
            path: 'training/batch/master/session/details', component: BatchDetailsComponent, canActivate: [RouteGuard],
      },
      {
            path: 'training/batch/master/session/completion', component: BatchCompletionDetailsComponent, canActivate: [RouteGuard],
      },
      {
            path: 'training/batch/master/session/feedback', component: BatchSessionFeedbackComponent, canActivate: [RouteGuard],
      },
      {
            path: 'training/batch/master/feedback', component: BatchFeedbackComponent, canActivate: [RouteGuard],
      },
      {
            path: 'training/batch/master/assignment-details', component: BatchTopicAssignmentDetailsComponent, canActivate: [RouteGuard],
      },
      {
            path: 'training/batch/master/project-details', component: BatchTopicAssignmentDetailsComponent, canActivate: [RouteGuard],
      },
      {
            path: "training/course/master/details", component: CourseComponent, canActivate: [RouteGuard],
      },
      {
            path: 'training/course/master/programming-assignment', component: CourseProgrammingAssignmentComponent, canActivate: [RouteGuard]
      },
      {
            path: "training/module/module-concept", component: ModuleConceptComponent, canActivate: [RouteGuard],
      },
      {
            path: 'talent/master/interview-schedule', component: InterviewScheduleComponent, canActivate: [RouteGuard],
      },   
      {
            path: 'training/batch/master/project-details', component: BatchProjectDetailsComponent, canActivate: [RouteGuard],
      },

      //***********************************************TALENT***************************************************
      
      {
            path: 'talent-dashboard', component: TalentDashboardComponent, canActivate: [RouteGuard],
      },
      {
            path: 'my-batches/concept-tree', component: ModuleConceptScoreComponent, canActivate: [RouteGuard],
      },
      {
            path: "my-batches", component: TalentBatchDetailsComponent, canActivate: [RouteGuard]
      },

      //***********************************************DEVELOPER***************************************************
      {
            path: 'resource', component: AdminResourceComponent, canActivate: [RouteGuard],
      },
      {
            path: 'timesheet', component: TimesheetComponent, canActivate: [RouteGuard],
      },

      //************************************TEST***************************************** */
      // {
      //       path: "code", component: CodeEditorComponent
      // },
      // {
      //       path: "cke-editor", component: CkeEditorComponent
      // },
      // {
      //       path: "question-modal", component: ProgrammingQuestionModalComponent
      // },

      // *******************************************COMING SOON**************************************************
      // {
      //       path: "coming-soon", component: ComingSoonComponent, canActivate: [RouteGuard]
      // },

      //*************************DECIDE LATER************************************** */
      // {
      //       path: "talent/master/:param", component: TalentMasterComponent, canActivate: [RouteGuard]
      // },
      // {
      //       path: "batch-talent", component: BatchTalentComponent, canActivate: [RouteGuard]
      // },
      // {
      //       path: "faculty-batch", component: FacultyBatchComponent, canActivate: [RouteGuard]
      // },
      // {
      //       path: "campus-registration-form", component: CampusRegistrationFormComponent
      // },
      // {
      //       path: "enquiry-form", component: EnquiryFormComponent
      // },
      // {
      //       path: "enquiry-enrollment", component: EnquiryEnrollmentComponent
      // },
      //***************************************************************************** */

      // *************************** OLD ***************************************************************

      // // *******************************************LOGIN**************************************************
      // {
      //       path: "login", component: LoginComponent
      // },
      // {
      //       path: "code", component: CodeEditorComponent
      // },
      // {
      //       path: "cke-editor", component: CkeEditorComponent
      // },
      // // {
      // //       path: "question-modal", component: ProgrammingQuestionModalComponent
      // // },

      // // *******************************************DASHBOARD**************************************************
      // {
      //       path: "dashboard", component: DashboardComponent, canActivate: [RouteGuard]
      // },

      // // *******************************************TALENT**************************************************
      // {
      //       path: 'talent',
      //       children: [
      //             { path: '', redirectTo: 'master', pathMatch: 'full' },
      //             { path: 'master', component: TalentMasterComponent, canActivate: [RouteGuard] },
      //             // children:[
      //             //       { path: 'interview-schedule', component: InterviewScheduleComponent, canActivate: [RouteGuard]},
      //             // ]},
      //             { path: 'enquiry', component: TalentEnquiryComponent, canActivate: [RouteGuard] },
      //       ]
      // },
      // {
      //       path: 'talent/master/interview-schedule', component: InterviewScheduleComponent, canActivate: [RouteGuard],
      // },

      // // *******************************************COURSE**************************************************
      // {
      //       path: 'course',
      //       children: [
      //             { path: '', redirectTo: 'master', pathMatch: 'full' },
      //             { path: 'master', component: CourseMasterComponent, canActivate: [RouteGuard] },
      //       ]
      // },
      // {
      //       path: 'course/master/session', component: CourseSessionComponent, canActivate: [RouteGuard]
      // },
      // // {
      // //       path: "course/master/details", component: CourseDetailsComponent, canActivate: [RouteGuard]
      // // },
      // {
      //       path: 'course/master/programming-assignment', component: CourseProgrammingAssignmentComponent, canActivate: [RouteGuard]
      // },
      // {
      //       path: "module", component: ModuleComponent, canActivate: [RouteGuard],
      // },
      // {
      //       path: "module/module-concept", component: ModuleConceptComponent, canActivate: [RouteGuard],
      // },
      // {
      //       path: 'my-batches/concept-tree', component: ModuleConceptScoreComponent, canActivate: [RouteGuard],
      // },
      // {
      //       path: 'batch/master/session/details/concept-tree', component: ModuleConceptScoreComponent, canActivate: [RouteGuard],
      // },
      // {
      //       path: "course/master/details", component: CourseComponent, canActivate: [RouteGuard],
      // },
      // // *******************************************BATCH**************************************************
      // {
      //       path: 'batch',
      //       children: [
      //             { path: '', redirectTo: 'master', pathMatch: 'full' },
      //             { path: 'master', component: BatchMasterComponent, canActivate: [RouteGuard] },
      //             { path: 'project', component: ProgrammingProjectComponent, canActivate: [RouteGuard] },
      //       ]
      // },
      // // {
      // //       path: 'batch/master/details', component: BatchDetailsComponent, canActivate: [RouteGuard],
      // // },

      // {
      //       path: 'batch/master/session/details', component: BatchDetailsComponent, canActivate: [RouteGuard],
      // },
      // {
      //       path: 'batch/master/session/completion', component: BatchCompletionDetailsComponent, canActivate: [RouteGuard],
      // },

      // {
      //       path: 'batch/master/session', component: BatchSessionComponent, canActivate: [RouteGuard],
      // },
      // {
      //       path: 'batch/master/session/feedback', component: BatchSessionFeedbackComponent, canActivate: [RouteGuard],
      // },
      // {
      //       path: 'batch/master/feedback', component: BatchFeedbackComponent, canActivate: [RouteGuard],
      // },
      // {
      //       path: 'batch/master/assignment-details', component: BatchTopicAssignmentDetailsComponent, canActivate: [RouteGuard],
      // },
      // {
      //       path: 'batch/master/project-details', component: BatchProjectDetailsComponent, canActivate: [RouteGuard],
      // },
      // // *******************************************COMPANY**************************************************
      // {
      //       path: 'company',
      //       children: [
      //             { path: '', redirectTo: 'master', pathMatch: 'full' },
      //             { path: 'master', component: CompanyMasterComponent, canActivate: [RouteGuard] },
      //             { path: 'branch', component: CompanyBranchComponent, canActivate: [RouteGuard] },
      //             { path: 'enquiry', component: CompanyEnquiryComponent, canActivate: [RouteGuard] },
      //             { path: 'requirement', component: CompanyMasterRequirementComponent, canActivate: [RouteGuard] },
      //       ]
      // },

      // // *******************************************COLLEGE**************************************************
      // {
      //       path: 'college',
      //       children: [
      //             { path: '', redirectTo: 'master', pathMatch: 'full' },
      //             { path: 'master', component: CollegeMasterComponent, canActivate: [RouteGuard] },
      //             { path: 'branch', component: CollegeBranchComponent, canActivate: [RouteGuard] },
      //             { path: 'campus', component: CollegeCampusComponent, canActivate: [RouteGuard] },
      //             { path: 'seminar', component: SeminarComponent, canActivate: [RouteGuard] },
      //       ]
      // },

      // // *******************************************ADMIN**************************************************
      // {
      //       path: 'admin',
      //       children: [
      //             { path: '', redirectTo: 'management', pathMatch: 'full' },
      //             { path: 'management', component: AdminManagementComponent, canActivate: [RouteGuard] },
      //             { path: 'salary-trend', component: SalaryTrendComponent, canActivate: [RouteGuard] },
      //             { path: 'career-objective', component: AdminCareerObjectiveComponent, canActivate: [RouteGuard] },
      //             { path: 'target-community', component: TargetCommunityComponent, canActivate: [RouteGuard] },
      //             {
      //                   path: 'models',
      //                   children: [
      //                         { path: '', redirectTo: 'degree', pathMatch: 'full' },
      //                         { path: 'degree', component: AdminDegreeComponent, canActivate: [RouteGuard] },
      //                         { path: 'designation', component: AdminDesignationComponent, canActivate: [RouteGuard] },
      //                         { path: 'department', component: DepartmentComponent, canActivate: [RouteGuard] },
      //                         { path: 'technology', component: AdminTechnologyComponent, canActivate: [RouteGuard] },
      //                         { path: 'role', component: AdminRoleComponent, canActivate: [RouteGuard] },
      //                         { path: 'project', component: ProjectComponent, canActivate: [RouteGuard] },
      //                         { path: 'university', component: AdminUniversityComponent, canActivate: [RouteGuard] },
      //                         { path: 'feeling', component: AdminFeelingComponent, canActivate: [RouteGuard] },
      //                         { path: 'resource', component: AdminResourceComponent, canActivate: [RouteGuard] },
      //                         { path: 'speaker', component: AdminSpeakerComponent, canActivate: [RouteGuard] },
      //                         { path: 'programming-question-type', component: AdminProgrammingQuestionTypeComponent, canActivate: [RouteGuard] },
      //                         { path: 'admin-event', component: AdminEventComponent, canActivate: [RouteGuard] },
      //                         { path: 'programming-concept', component: ProgrammingConceptComponent, canActivate: [RouteGuard] },
      //                         { path: 'programming-language', component: ProgrammingLanguageComponent, canActivate: [RouteGuard] },
      //                         { path: 'blog-topic', component: BlogTopicComponent, canActivate: [RouteGuard] },
      //                   ]
      //             },
      //             {
      //                   path: 'employee',
      //                   children: [
      //                         { path: '', redirectTo: 'user', pathMatch: 'full' },
      //                         { path: 'user', component: AdminUserComponent, canActivate: [RouteGuard] },
      //                         { path: 'faculty', component: FacultyMasterComponent, canActivate: [RouteGuard] },
      //                         { path: 'other', component: AdminEmployeeComponent, canActivate: [RouteGuard] },
      //                   ]
      //             },
      //             {
      //                   path: 'feedback',
      //                   children: [
      //                         { path: '', redirectTo: 'question', pathMatch: 'full' },
      //                         { path: 'question', component: AdminFeedbackComponent, canActivate: [RouteGuard] },
      //                         { path: 'group', component: AdminGroupComponent, canActivate: [RouteGuard] },
      //                   ]
      //             },
      //             {
      //                   path: 'coding-problems',
      //                   children: [
      //                         { path: '', redirectTo: 'questions', pathMatch: 'full' },
      //                         { path: 'questions', component: ProgrammingQuestionComponent, canActivate: [RouteGuard] },
      //                         { path: 'answers', component: ProgrammingQuestionTalentAnswerComponent, canActivate: [RouteGuard] },
      //                         { path: 'assignment', component: ProgrammingAssignmentComponent, canActivate: [RouteGuard] },
      //                   ]
      //             },
      //             {
      //                   path: 'community',
      //                   children: [
      //                         { path: '', redirectTo: 'blog', pathMatch: 'full' },
      //                         { path: 'blog', component: BlogVerificationComponent, canActivate: [RouteGuard] },
      //                   ]
      //             },
      //       ]
      // },

      // // ******************************************* REPORT **************************************************
      // {
      //       path: 'report',
      //       children: [
      //             {
      //                   path: 'talent',
      //                   children: [
      //                         { path: '', redirectTo: 'calling', pathMatch: 'full' },
      //                         { path: 'lifetime-value', component: LifetimeValueReportComponent, canActivate: [RouteGuard] },
      //                         { path: 'waiting-list', component: WaitingListReportComponent, canActivate: [RouteGuard] },
      //                         { path: 'pro-summary', component: ProSummaryReportComponent, canActivate: [RouteGuard] },
      //                         { path: 'calling', component: CallingReportsComponent, canActivate: [RouteGuard] },
      //                         { path: 'next-action', component: NextActionReportComponent, canActivate: [RouteGuard] },
      //                         { path: 'fresher-summary', component: FresherSummaryComponent, canActivate: [RouteGuard] },
      //                         { path: 'package-summary', component: PackageSummaryComponent, canActivate: [RouteGuard] },
      //                   ]
      //             },
      //             {
      //                   path: 'faculty',
      //                   children: [
      //                         { path: '', redirectTo: 'faculty-report', pathMatch: 'full' },
      //                         { path: 'faculty-report', component: FacultyReportComponent, canActivate: [RouteGuard] },
      //                   ]
      //             },
      //             {
      //                   path: 'admin',
      //                   children: [
      //                         { path: '', redirectTo: 'login-report', pathMatch: 'full' },
      //                         { path: 'login-report', component: LoginReportComponent, canActivate: [RouteGuard] },
      //                   ]
      //             },
      //       ]
      // },

      // // *******************************************TEST**************************************************
      // {
      //       path: "test-master", component: TestMasterComponent, canActivate: [RouteGuard]
      // },

      // // *******************************************TARGET COMMUNITY**************************************************
      // {
      //       path: "target-community", component: TargetCommunityComponent, canActivate: [RouteGuard]
      // },

      // // *******************************************PROFILE**************************************************
      // {
      //       path: 'profile',
      //       children: [
      //             { path: '', redirectTo: 'timesheet', pathMatch: 'full' },
      //             { path: 'timesheet', component: TimesheetComponent, canActivate: [RouteGuard] },
      //             { path: 'account-settings', component: AccountSettingsComponent, canActivate: [RouteGuard] },
      //             { path: 'logout', component: LogoutComponent },
      //             { path: 'refer-and-earn', component: ReferAndEarnComponent, canActivate: [RouteGuard] },
      //       ]
      // },

      // // *******************************************RESOURCE**************************************************
      // {
      //       path: "resource", component: AdminResourceComponent, canActivate: [RouteGuard]
      // },

      // // *******************************************COMING SOON**************************************************
      // {
      //       path: "coming-soon", component: ComingSoonComponent, canActivate: [RouteGuard]
      // },

      // // *******************************************MY-OPPURTUNITIES**************************************************
      // {
      //       path: 'my-opportunities', component: MyOpportunitiesComponent, canActivate: [RouteGuard],
      //       // children: [
      //       //       { path: 'company-details', component: CompanyDetailsComponent, canActivate: [RouteGuard]},
      //       // ]
      // },

      // {
      //       path: "my-opportunities/company-details", component: CompanyDetailsComponent, canActivate: [RouteGuard]
      // },

      // // *******************************************MY-COURSES**************************************************
      // // {
      // //       path: 'my-courses', component: MyCoursesComponent, canActivate: [RouteGuard],
      // //       // children: [
      // //       //       { path: 'company-details', component: CompanyDetailsComponent, canActivate: [RouteGuard]},
      // //       // ]
      // // },

      // // {
      // //       path: "my-courses/my-sessions", component: MySessionsComponent, canActivate: [RouteGuard]
      // // },

      // {
      //       path: "my-batches", component: TalentBatchDetailsComponent, canActivate: [RouteGuard]
      // },

      // // *******************************************PROBLEMS-OF-THE-DAY**************************************************
      // {
      //       path: 'problems-of-the-day', component: ProblemOfTheDayComponent, canActivate: [RouteGuard],
      // },

      // {
      //       path: "problems-of-the-day/problem-details", component: ProblemDetailsComponent, canActivate: [RouteGuard]
      // },

      // // *******************************************PRACTISE**************************************************
      // {
      //       path: 'practice', component: PracticeComponent, canActivate: [RouteGuard],
      // },

      // {
      //       path: "practice/problem-details", component: PracticeProblemDetailsComponent, canActivate: [RouteGuard]
      // },

      // // *******************************************EVENTS**************************************************
      // {
      //       path: 'events', component: EventsComponent, canActivate: [RouteGuard],
      // },

      // {
      //       path: "events/event-details", component: EventDetailsComponent, canActivate: [RouteGuard]
      // },
      // // {
      // //       path: "events/batch-details", component: CourseDetailsComponent, canActivate: [RouteGuard]
      // // },
      // // {
      // //       path: "course-details", component: CourseDetailsComponent, canActivate: [RouteGuard]
      // // },

      // // *******************************************TEST**************************************************
      // {
      //       path: 'test', component: TestComponent, canActivate: [RouteGuard],
      // },

      // {
      //       path: "test/test-details", component: TestDetailsComponent, canActivate: [RouteGuard]
      // },

      // {
      //       path: "test/test-details/problem-details", component: ProblemDetailsComponent, canActivate: [RouteGuard]
      // },

      // // *******************************************TALENT DASHBOARD**************************************************
      // {
      //       path: 'talent-dashboard', component: TalentDashboardComponent, canActivate: [RouteGuard],
      // },

      // // *******************************************TALENT BATCH DETAILS**************************************************

      // { path: 'talent/batch/details', component: TalentBatchDetailsComponent, canActivate: [RouteGuard] },


      // // ******************************************* FACULTY DASHBOARD **************************************************
      // {
      //       path: 'faculty-dashboard', component: FacultyDashboardComponent, canActivate: [RouteGuard],
      // },

      // // *******************************************CODING PROBLEMS**************************************************
      // {
      //       path: 'coding-problems',
      //       children: [
      //             { path: '', redirectTo: 'questions', pathMatch: 'full' },
      //             { path: 'questions', component: ProgrammingQuestionComponent, canActivate: [RouteGuard] },
      //             { path: 'answers', component: ProgrammingQuestionTalentAnswerComponent, canActivate: [RouteGuard] },
      //             { path: 'assignment', component: ProgrammingAssignmentComponent, canActivate: [RouteGuard] },
      //       ]
      // },

      // // *******************************************ANSWER DETAILS**************************************************

      // {
      //       path: 'admin/coding-problems/answers/answer-details', component: ProgrammingQuestionTalentAnswerDetailsComponent, canActivate: [RouteGuard],
      // },

      // {
      //       path: 'coding-problems/answers/answer-details', component: ProgrammingQuestionTalentAnswerDetailsComponent, canActivate: [RouteGuard],
      // },

      // // *******************************************PROGRAMMING CONCEPT**************************************************

      // // {
      // //       path: 'entity/programming-concept', component: ProgrammingConceptComponent, canActivate: [RouteGuard],
      // // },

      // {
      //       path: 'entity',
      //       children: [
      //             { path: '', redirectTo: 'programming-concept', pathMatch: 'full' },
      //             { path: 'programming-concept', component: ProgrammingConceptComponent, canActivate: [RouteGuard] },
      //             { path: 'blog-topic', component: BlogTopicComponent, canActivate: [RouteGuard] },
      //       ]
      // },

      // // *******************************************PRACTICE PROBLEM LIST**************************************************
      // {
      //       path: 'practice/practice-problem-list', component: PracticeProblemListComponent, canActivate: [RouteGuard],
      // },

      // // *******************************************COMMUNITY**************************************************

      {
            path: 'community',
            children: [
                  {
                        path: 'blog',
                        children: [
                              { path: '', redirectTo: 'my-blogs', pathMatch: 'full' },
                              { path: 'my-blogs', component: BlogComponent },
                              { path: 'browse-topics', component: BrowseBlogsComponent, canActivate: [RouteGuard] },
                        ]
                  },
            ]
      },

      // // *******************************************BLOG DETAILS**************************************************
      {
            path: 'community/blog/browse-topics/blog-details', component: BlogDetailsComponent
      },

      // //*************************DECIDE LATER************************************** */
      // // {
      // //       path: "talent/master/:param", component: TalentMasterComponent, canActivate: [RouteGuard]
      // // },
      // // {
      // //       path: "batch-talent", component: BatchTalentComponent, canActivate: [RouteGuard]
      // // },
      // // {
      // //       path: "faculty-batch", component: FacultyBatchComponent, canActivate: [RouteGuard]
      // // },
      // // {
      // //       path: "campus-registration-form", component: CampusRegistrationFormComponent
      // // },
      // // {
      // //       path: "enquiry-form", component: EnquiryFormComponent
      // // },
      // // {
      // //       path: "enquiry-enrollment", component: EnquiryEnrollmentComponent
      // // },
      // //***************************************************************************** */

      {
            path: "",
            pathMatch: "full",
            redirectTo: "/login"
      },
      {
            path: "**", component: PageNotFoundComponent
      },
];

const routerOptions: ExtraOptions = {
      // scrollPositionRestoration: 'enabled',
      anchorScrolling: 'enabled',
      useHash: true,
};

@NgModule({
      imports: [RouterModule.forRoot(routes, routerOptions)],
      exports: [RouterModule]
})
export class AppRoutingModule { }
