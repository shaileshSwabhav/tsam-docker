import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProgrammingQuestionModalComponent } from './programming-question-modal.component';

describe('ProgrammingQuestionModalComponent', () => {
  let component: ProgrammingQuestionModalComponent;
  let fixture: ComponentFixture<ProgrammingQuestionModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ProgrammingQuestionModalComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProgrammingQuestionModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
