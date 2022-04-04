import { NgModule } from '@angular/core';
import { DatePipe } from '@angular/common';
import { HashLocationStrategy, LocationStrategy } from '@angular/common'
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { BrowserModule } from '@angular/platform-browser';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { TouchedErrorStateMatcher } from './service/talent/touched-error-state.matcher';

// Directives
import { AllowNumbersWithDecimalDirective } from './directives/allow-numbers-with-decimal/allow-numbers-with-decimal.directive';
import { AllowNumbersOnlyDirective } from './directives/allow-numbers-only/allow-numbers-only.directive';
import { EmptyToNullDirective } from './directives/empty-to-null/empty-to-null.directive';


import { LoginComponent } from './component/login/login.component';
import { SubNavbarComponent } from './component/sub-navbar/sub-navbar.component';
import { DashboardComponent } from './component/dashboard/dashboard.component';
import { AdminRoleComponent } from './component/admin-role/admin-role.component';
import { TestMasterComponent } from './component/test-master/test-master.component';
import { BatchMasterComponent } from './component/batch-master/batch-master.component';
import { MasterNavbarComponent } from './component/master-navbar/master-navbar.component';
import { TalentMasterComponent } from './component/talent-master/talent-master.component';
import { MasterFooterComponent } from './component/master-footer/master-footer.component';
import { CourseMasterComponent } from './component/course-master/course-master.component';
import { PageNotFoundComponent } from './component/page-not-found/page-not-found.component';
import { AdminUniversityComponent } from './component/admin-university/admin-university.component';
import { InterviewScheduleComponent } from './component/interview-schedule/interview-schedule.component';
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
import { EnquiryEnrollmentComponent } from './component/enquiry-enrollment/enquiry-enrollment.component';
import { CampusRegistrationFormComponent } from './component/campus-registration-form/campus-registration-form.component';
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
import { TalentEnquiryComponent } from './component/talent-enquiry/talent-enquiry.component';
import { CompanyMasterComponent } from './component/company-master/company-master.component';
import { CollegeCampusComponent } from './component/college-campus/college-campus.component';
import { FacultyMasterComponent } from './component/faculty-master/faculty-master.component';
import { BatchTalentComponent } from './component/batch-talent/batch-talent.component';
import { CollegeMasterComponent } from './component/college-master/college-master.component';
import { CollegeBranchComponent } from './component/college-branch/college-branch.component';
import { AdminDegreeComponent } from './component/admin-degree/admin-degree.component';
import { AdminDesignationComponent } from './component/admin-designation/admin-designation.component';
import { CompanyBranchComponent } from './component/company-branch/company-branch.component';
import { EnquiryFormComponent } from './component/enquiry-form/enquiry-form.component';
import { CourseSessionComponent } from './component/course-session/course-session.component';
import { CompanyEnquiryComponent } from './component/company-enquiry/company-enquiry.component';
import { AdminTechnologyComponent } from './component/admin-technology/admin-technology.component';
import { CompanyMasterRequirementComponent } from './component/company-master-requirement/company-master-requirement.component';
import { FacultyBatchComponent } from './component/faculty-batch/faculty-batch.component';
import { SelectedTalentListComponent } from './component/selected-talent-list/selected-talent-list.component';
import { SeminarComponent } from './component/seminar/seminar.component';
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
import { QuestionComponent } from './component/question/question.component';
import { BatchSessionComponent } from './component/batch-session/batch-session.component'

// Angular material
import { MatStepperModule } from '@angular/material/stepper';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatListModule } from '@angular/material/list';
import { MatSelectModule } from '@angular/material/select';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { ErrorStateMatcher } from '@angular/material/core';
import { MatMenuModule } from '@angular/material/menu';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatTabsModule } from '@angular/material/tabs';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatProgressBarModule } from '@angular/material/progress-bar';

// Other Dependencies
import { DragDropModule } from '@angular/cdk/drag-drop';
import { CKEditorModule } from 'ckeditor4-angular';
import { ChartsModule } from 'ng2-charts';
import { NgxPaginationModule } from 'ngx-pagination';
import { NgxSpinnerModule } from "ngx-bootstrap-spinner";
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { NgSelectModule } from '@ng-select/ng-select';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { CourseDetailsComponent } from './component/course-details/course-details.component';
import { ProgrammingQuestionComponent } from './component/programming-question/programming-question.component';
import { AdminProgrammingQuestionTypeComponent } from './component/admin-programming-question-type/admin-programming-question-type.component';
import { AdminEventComponent } from './component/admin-event/admin-event.component';
import { BatchProjectComponent } from './component/batch-project/batch-project.component';
import { ProgrammingQuestionTalentAnswerComponent } from './component/programming-question-talent-answer/programming-question-talent-answer.component';
import { ProgrammingQuestionTalentAnswerDetailsComponent } from './component/programming-question-talent-answer-details/programming-question-talent-answer-details.component';
import { ProgrammingConceptComponent } from './component/programming-concept/programming-concept.component';
import { PracticeProblemDetailsComponent } from './component/practice-problem-details/practice-problem-details.component'
import { BatchDetailsComponent } from './component/batch-details/batch-details.component';
import { CodeEditorComponent } from './component/code-editor/code-editor.component';
import { PracticeProblemListComponent } from './component/practice-problem-list/practice-problem-list.component';
import { ProgrammingLanguageComponent } from './component/programming-language/programming-language.component';
import { CkeEditorComponent } from './component/cke-editor/cke-editor.component';
import { BlogTopicComponent } from './component/blog-topic/blog-topic.component';
import { BlogComponent } from './component/blog/blog.component';
import { BlogVerificationComponent } from './component/blog-verification/blog-verification.component';
import { BrowseBlogsComponent } from './component/browse-blogs/browse-blogs.component';
import { BlogDetailsComponent } from './component/blog-details/blog-details.component';
import { ProgrammingQuestionModalComponent } from './component/programming-question-modal/programming-question-modal.component'
import { CourseProgrammingAssignmentComponent } from './component/course-programming-assignment/course-programming-assignment.component'
import { ProgrammingAssignmentComponent } from './component/programming-assignment/programming-assignment.component';
import { CourseModuleComponent } from './component/course-module/course-module.component';
import { ModuleComponent } from './component/module/module.component';
import { BatchTopicAssignmentDetailsComponent } from './component/batch-topic-assignment-details/batch-topic-assignment-details.component';
import { BatchTopicAssignmentScoreComponent } from './component/batch-topic-assignment-score/batch-topic-assignment-score.component';
import { BatchTopicComponent } from './component/batch-topic/batch-topic.component';
import { CourseComponent } from './component/course/course.component';
import { CourseModuleTopicComponent } from './component/course-module-topic/course-module-topic.component';
import { BatchCompletionDetailsComponent } from './component/batch-completion-details/batch-completion-details.component';
import { MinutePipe } from './pipe/minute/minute.pipe';
import { TalentBatchDetailsComponent } from './component/talent-batch-details/talent-batch-details.component';
import { TalentBatchDetailsFeedbackComponent } from './component/talent-batch-details-feedback/talent-batch-details-feedback.component';
import { TalentBatchSessionPlanComponent } from './component/talent-batch-session-plan/talent-batch-session-plan.component';
import { TalentBatchDetailsAssignmentComponent } from './component/talent-batch-details-assignment/talent-batch-details-assignment.component';
import { NgCircleProgressModule } from 'ng-circle-progress';
import { BatchSessionPlanComponent } from './component/batch-session-plan/batch-session-plan.component';
import { BatchSessionDetailsComponent } from './component/batch-session-details/batch-session-details.component';
import { ProgrammingProjectComponent } from './component/programming-project/programming-project.component';
import { TalentFeedbackModalComponent } from './component/talent-feedback-modal/talent-feedback-modal.component';
import { BatchTopicAssignmentComponent } from './component/batch-topic-assignment/batch-topic-assignment.component';
import { BatchSessionCompletionFeedbackComponent } from './component/batch-session-completion-feedback/batch-session-completion-feedback.component';
import { TalentBatchDetailsProjectComponent } from './component/talent-batch-details-project/talent-batch-details-project.component';
import { NgAudioRecorderModule } from 'ng-audio-recorder';
import { SpinnerComponent } from './component/spinner/spinner.component';
import { ModuleConceptNodeComponent } from './component/module-concept-node/module-concept-node.component';
import { ModuleConceptComponent } from './component/module-concept/module-concept.component';
import { BatchModulesComponent } from './component/batch-modules/batch-modules.component';
import { ModuleConceptScoreComponent } from './component/module-concept-score/module-concept-score.component';
import { ModuleConceptScoreNodeComponent } from './component/module-concept-score-node/module-concept-score-node.component';
import { ApiInterceptorService } from './service/api-interceptor/api-interceptor.service';
import { BatchProgressReportComponent } from './component/batch-progress-report/batch-progress-report.component';
import { BatchProjectDetailsComponent } from './component/batch-project-details/batch-project-details.component';
import { BatchProjectScoreComponent } from './component/batch-project-score/batch-project-score.component';
import { FocusFirstInvalidDirective } from './directives/focus-first-invalid/focus-first-invalid.directive';
import { FacultyTalentFeedbackComponent } from './component/faculty-talent-feedback/faculty-talent-feedback.component';
import { SessionPlanUpdateComponent } from './component/session-plan-update/session-plan-update.component';

@NgModule({
      declarations: [
            AppComponent,
            LoginComponent,
            AdminRoleComponent,
            DashboardComponent,
            TestMasterComponent,
            BatchMasterComponent,
            PageNotFoundComponent,
            MasterNavbarComponent,
            TalentMasterComponent,
            MasterFooterComponent,
            CourseMasterComponent,
            TalentEnquiryComponent,
            CompanyMasterComponent,
            FacultyMasterComponent,
            CollegeCampusComponent,
            CompanyEnquiryComponent,
            AdminTechnologyComponent,
            CompanyMasterRequirementComponent,
            FacultyBatchComponent,
            SelectedTalentListComponent,
            SeminarComponent,
            EmptyToNullDirective,
            QuestionComponent,
            AllowNumbersOnlyDirective,
            BatchTalentComponent,
            CollegeMasterComponent,
            CollegeBranchComponent,
            AdminDegreeComponent,
            AdminDesignationComponent,
            CompanyBranchComponent,
            EnquiryFormComponent,
            SubNavbarComponent,
            CourseSessionComponent,
            AdminUniversityComponent,
            AllowNumbersWithDecimalDirective,
            BatchSessionComponent,
            InterviewScheduleComponent,
            CallingReportsComponent,
            BatchSessionFeedbackComponent,
            AdminFeedbackComponent,
            BatchFeedbackComponent,
            TargetCommunityComponent,
            DepartmentComponent,
            ProjectComponent,
            AdminEmployeeComponent,
            LifetimeValueReportComponent,
            NextActionReportComponent,
            AdminCareerObjectiveComponent,
            AdminUserComponent,
            AdminManagementComponent,
            AdminGroupComponent,
            AdminFeelingComponent,
            TimesheetComponent,
            AdminResourceComponent,
            EnquiryEnrollmentComponent,
            CampusRegistrationFormComponent,
            SalaryTrendComponent,
            WaitingListReportComponent,
            LoginReportComponent,
            AdminSpeakerComponent,
            ProSummaryReportComponent,
            FresherSummaryComponent,
            LogoutComponent,
            AccountSettingsComponent,
            ComingSoonComponent,
            MyOpportunitiesComponent,
            CompanyDetailsComponent,
            MyCoursesComponent,
            PackageSummaryComponent,
            FacultyReportComponent,
            MySessionsComponent,
            ProblemOfTheDayComponent,
            ProblemDetailsComponent,
            PracticeComponent,
            EventsComponent,
            EventDetailsComponent,
            TestComponent,
            TestDetailsComponent,
            TalentDashboardComponent,
            ReferAndEarnComponent,
            FacultyDashboardComponent,
            CourseDetailsComponent,
            ProgrammingQuestionComponent,
            AdminProgrammingQuestionTypeComponent,
            AdminEventComponent,
            BatchProjectComponent,
            ProgrammingQuestionTalentAnswerComponent,
            ProgrammingQuestionTalentAnswerDetailsComponent,
            ProgrammingConceptComponent,
            PracticeProblemDetailsComponent,
            BatchDetailsComponent,
            CodeEditorComponent,
            PracticeProblemListComponent,
            ProgrammingLanguageComponent,
            CkeEditorComponent,
            ProgrammingQuestionModalComponent,
            BlogTopicComponent,
            BlogComponent,
            BlogVerificationComponent,
            BrowseBlogsComponent,
            BlogDetailsComponent,
            ProgrammingAssignmentComponent,
            CourseProgrammingAssignmentComponent,
            CourseModuleComponent,
            ModuleComponent,
            BatchTopicAssignmentDetailsComponent,
            BatchTopicAssignmentScoreComponent,
            BatchTopicComponent,
            CourseComponent,
            CourseModuleTopicComponent,
            BatchCompletionDetailsComponent,
            MinutePipe,
            TalentBatchDetailsComponent,
            TalentBatchDetailsFeedbackComponent,
            TalentBatchSessionPlanComponent,
            TalentBatchDetailsAssignmentComponent,
            BatchSessionPlanComponent,
            BatchSessionDetailsComponent,
            ProgrammingProjectComponent,
            TalentFeedbackModalComponent,
            BatchTopicAssignmentComponent,
            BatchSessionCompletionFeedbackComponent,
            TalentBatchDetailsProjectComponent,
            SpinnerComponent,
            ModuleConceptNodeComponent,
            ModuleConceptComponent,
            BatchModulesComponent,
            ModuleConceptScoreComponent,
            ModuleConceptScoreNodeComponent,
            BatchProgressReportComponent,
            BatchProjectDetailsComponent,
            BatchProjectScoreComponent,
            FocusFirstInvalidDirective,
            FacultyTalentFeedbackComponent,
            SessionPlanUpdateComponent,
      ],
      // Deprecated(since 9.0) : Ivy will instantiate component 
      entryComponents: [ProgrammingQuestionModalComponent],
      imports: [
            FormsModule,
            BrowserModule,
            NgSelectModule,
            HttpClientModule,
            AppRoutingModule,
            ReactiveFormsModule,
            NgxSpinnerModule,
            BrowserAnimationsModule,
            NgxPaginationModule,
            NgbModule,
            MatStepperModule,
            MatInputModule,
            MatButtonModule,
            MatListModule,
            MatSelectModule,
            MatDatepickerModule,
            MatNativeDateModule,
            FontAwesomeModule,
            MatMenuModule,
            MatAutocompleteModule,
            MatSortModule,
            MatTableModule,
            MatSidenavModule,
            MatTabsModule,
            MatToolbarModule,
            ChartsModule,
            CKEditorModule,
            MatProgressBarModule,
            DragDropModule,
            NgAudioRecorderModule,
            NgCircleProgressModule.forRoot({})
      ],
      providers: [{ provide: LocationStrategy, useClass: HashLocationStrategy },
      { provide: ErrorStateMatcher, useClass: TouchedErrorStateMatcher },
      {
            provide: HTTP_INTERCEPTORS,
            useClass: ApiInterceptorService,
            multi: true,
      },
            DatePipe],
      bootstrap: [AppComponent]
})
export class AppModule { }
