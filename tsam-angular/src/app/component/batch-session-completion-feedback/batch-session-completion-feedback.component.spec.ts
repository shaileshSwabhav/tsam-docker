import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BatchSessionCompletionFeedbackComponent } from './batch-session-completion-feedback.component';

describe('BatchSessionCompletionFeedbackComponent', () => {
  let component: BatchSessionCompletionFeedbackComponent;
  let fixture: ComponentFixture<BatchSessionCompletionFeedbackComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BatchSessionCompletionFeedbackComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BatchSessionCompletionFeedbackComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
