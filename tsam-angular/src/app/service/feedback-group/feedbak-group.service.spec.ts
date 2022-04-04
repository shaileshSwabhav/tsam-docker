import { TestBed } from '@angular/core/testing';

import { FeedbackGroupService } from './feedbak-group.service';

describe('FeedbackGroupService', () => {
  let service: FeedbackGroupService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(FeedbackGroupService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
