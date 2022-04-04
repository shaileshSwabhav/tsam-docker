import { TestBed } from '@angular/core/testing';

import { TargetCommunityFunctionService } from './target-community-function.service';

describe('TargetCommunityFunctionService', () => {
  let service: TargetCommunityFunctionService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TargetCommunityFunctionService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
