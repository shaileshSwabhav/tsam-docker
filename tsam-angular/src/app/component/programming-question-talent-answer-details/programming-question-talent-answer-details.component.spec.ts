import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProgrammingQuestionTalentAnswerDetailsComponent } from './programming-question-talent-answer-details.component';

describe('SolutionDetailsComponent', () => {
  let component: ProgrammingQuestionTalentAnswerDetailsComponent;
  let fixture: ComponentFixture<ProgrammingQuestionTalentAnswerDetailsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ProgrammingQuestionTalentAnswerDetailsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProgrammingQuestionTalentAnswerDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
