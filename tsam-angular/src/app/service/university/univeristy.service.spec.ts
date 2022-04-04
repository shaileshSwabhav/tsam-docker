import { TestBed } from '@angular/core/testing';

import { UniveristyService } from './univeristy.service';

describe('UniveristyService', () => {
  let service: UniveristyService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(UniveristyService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
