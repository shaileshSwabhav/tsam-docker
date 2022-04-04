import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AdminProgrammingQuestionTypeComponent } from './admin-programming-question-type.component';

describe('AdminProgrammingQuestionTypeComponent', () => {
  let component: AdminProgrammingQuestionTypeComponent;
  let fixture: ComponentFixture<AdminProgrammingQuestionTypeComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AdminProgrammingQuestionTypeComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AdminProgrammingQuestionTypeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
