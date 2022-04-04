import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BatchFeedbackComponent } from './batch-feedback.component';

describe('BatchFeedbackComponent', () => {
  let component: BatchFeedbackComponent;
  let fixture: ComponentFixture<BatchFeedbackComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BatchFeedbackComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BatchFeedbackComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
