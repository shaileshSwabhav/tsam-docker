import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BatchSessionFeedbackComponent } from './batch-session-feedback.component';

describe('BatchSessionFeedbackComponent', () => {
  let component: BatchSessionFeedbackComponent;
  let fixture: ComponentFixture<BatchSessionFeedbackComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BatchSessionFeedbackComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BatchSessionFeedbackComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
