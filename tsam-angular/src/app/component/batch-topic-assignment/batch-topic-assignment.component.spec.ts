import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BatchTopicAssignmentComponent } from './batch-topic-assignment.component';

describe('BatchTopicAssignmentComponent', () => {
  let component: BatchTopicAssignmentComponent;
  let fixture: ComponentFixture<BatchTopicAssignmentComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BatchTopicAssignmentComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BatchTopicAssignmentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
