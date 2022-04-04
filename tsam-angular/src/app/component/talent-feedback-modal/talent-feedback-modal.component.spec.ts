import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TalentFeedbackModalComponent } from './talent-feedback-modal.component';

describe('TalentFeedbackModalComponent', () => {
  let component: TalentFeedbackModalComponent;
  let fixture: ComponentFixture<TalentFeedbackModalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TalentFeedbackModalComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TalentFeedbackModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
