import { TestBed } from '@angular/core/testing';

import { TargetCommunityService } from './target-community.service';

describe('TargetCommunityService', () => {
  let service: TargetCommunityService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TargetCommunityService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
