import { TestBed } from '@angular/core/testing';

import { BatchTopicAssignmentService } from './batch-topic-assignment.service';

describe('BatchTopicAssignmentService', () => {
  let service: BatchTopicAssignmentService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BatchTopicAssignmentService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
