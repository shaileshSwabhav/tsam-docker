import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProgrammingQuestionComponent } from './programming-question.component';

describe('ProgrammingQuestionComponent', () => {
  let component: ProgrammingQuestionComponent;
  let fixture: ComponentFixture<ProgrammingQuestionComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ProgrammingQuestionComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProgrammingQuestionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
