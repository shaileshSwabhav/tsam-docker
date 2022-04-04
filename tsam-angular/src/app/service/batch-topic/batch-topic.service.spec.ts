import { TestBed } from '@angular/core/testing';

import { BatchTopicService } from './batch-topic.service';

describe('BatchTopicService', () => {
  let service: BatchTopicService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BatchTopicService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
