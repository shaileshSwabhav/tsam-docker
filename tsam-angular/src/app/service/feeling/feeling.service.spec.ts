import { TestBed } from '@angular/core/testing';

import { FeelingService } from './feeling.service';

describe('FeelingService', () => {
  let service: FeelingService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(FeelingService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
