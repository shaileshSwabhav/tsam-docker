import { TestBed } from '@angular/core/testing';

import { BatchTalentService } from './batch-talent.service';

describe('BatchTalentService', () => {
  let service: BatchTalentService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BatchTalentService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
