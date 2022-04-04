import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { FacultyTalentFeedbackComponent } from './faculty-talent-feedback.component';

describe('FacultyTalentFeedbackComponent', () => {
  let component: FacultyTalentFeedbackComponent;
  let fixture: ComponentFixture<FacultyTalentFeedbackComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ FacultyTalentFeedbackComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(FacultyTalentFeedbackComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
