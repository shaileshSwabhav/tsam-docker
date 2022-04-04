import { TestBed } from '@angular/core/testing';

import { SeminarService } from './seminar.service';

describe('SeminarService', () => {
  let service: SeminarService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(SeminarService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
