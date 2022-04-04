import { TestBed } from '@angular/core/testing';

import { BatchProjectService } from './batch-project.service';

describe('BatchProjectService', () => {
  let service: BatchProjectService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BatchProjectService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
