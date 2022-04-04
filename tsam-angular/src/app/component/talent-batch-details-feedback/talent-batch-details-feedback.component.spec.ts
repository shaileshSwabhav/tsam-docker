import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TalentBatchDetailsFeedbackComponent } from './talent-batch-details-feedback.component';

describe('TalentBatchDetailsFeedbackComponent', () => {
  let component: TalentBatchDetailsFeedbackComponent;
  let fixture: ComponentFixture<TalentBatchDetailsFeedbackComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TalentBatchDetailsFeedbackComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TalentBatchDetailsFeedbackComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
