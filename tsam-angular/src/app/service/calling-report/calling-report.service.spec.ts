import { TestBed } from '@angular/core/testing';

import { CallingReportService } from './calling-report.service';

describe('CallingReportService', () => {
  let service: CallingReportService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(CallingReportService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
