import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProgrammingQuestionTalentAnswerComponent } from './programming-question-talent-answer.component';

describe('ProgrammingSolutionComponent', () => {
  let component: ProgrammingQuestionTalentAnswerComponent;
  let fixture: ComponentFixture<ProgrammingQuestionTalentAnswerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ProgrammingQuestionTalentAnswerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProgrammingQuestionTalentAnswerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
