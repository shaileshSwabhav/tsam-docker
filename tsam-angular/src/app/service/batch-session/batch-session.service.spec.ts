import { TestBed } from '@angular/core/testing';

import { BatchSessionService } from './batch-session.service';

describe('BatchSessionService', () => {
  let service: BatchSessionService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(BatchSessionService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
