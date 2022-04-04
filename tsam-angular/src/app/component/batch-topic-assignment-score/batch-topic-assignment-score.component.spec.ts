import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BatchTopicAssignmentScoreComponent } from './batch-topic-assignment-score.component';

describe('BatchTopicAssignmentScoreComponent', () => {
  let component: BatchTopicAssignmentScoreComponent;
  let fixture: ComponentFixture<BatchTopicAssignmentScoreComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BatchTopicAssignmentScoreComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BatchTopicAssignmentScoreComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
